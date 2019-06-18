package raft

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	fmt "fmt"
	"math/big"
	mathrand "math/rand"
	"net"
	"sync"
	"time"

	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/vault/cluster"
)

// RaftTLSKey is a single TLS keypair in the Keyring
type RaftTLSKey struct {
	// ID is a unique identifier for this Key
	ID string `json:"id"`

	// KeyType defines the algorighm used to generate the private keys
	KeyType string `json:"key_type"`

	// AppliedIndex is the earliest known raft index that safely contains this
	// key.
	AppliedIndex uint64 `json:"applied_index"`

	// CertBytes is the marshaled certificate.
	CertBytes []byte `json:"cluster_cert"`

	// KeyParams is the marshaled private key.
	KeyParams *certutil.ClusterKeyParams `json:"cluster_key_params"`

	// CreatedTime is the time this key was generated. This value is useful in
	// determining when the next rotation should be.
	CreatedTime time.Time `json:"created_time"`

	parsedCert *x509.Certificate
	parsedKey  *ecdsa.PrivateKey
}

// RaftTLSKeyring is the set of keys that raft uses for network communication.
// Only one key is used to dial at a time but both keys will be used to accept
// connections.
type RaftTLSKeyring struct {
	// Keys is the set of available key pairs
	Keys []*RaftTLSKey `json:"keys"`

	// AppliedIndex is the earliest known raft index that safely contains the
	// latest key in the keyring.
	AppliedIndex uint64 `json:"applied_index"`

	// Term is an incrementing identifier value used to quickly determine if two
	// states of the keyring are different.
	Term uint64 `json:"term"`

	// ActiveKey is the key that is active in the keyring. Only
	// the active key is used for dialing.
	ActiveKey *RaftTLSKey `json:"active_key"`
}

func GenerateTLSKey() (*RaftTLSKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return nil, err
	}

	host, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	host = fmt.Sprintf("raft-%s", host)
	template := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: host,
		},
		DNSNames: []string{host},
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement | x509.KeyUsageCertSign,
		SerialNumber: big.NewInt(mathrand.Int63()),
		NotBefore:    time.Now().Add(-30 * time.Second),
		// 30 years of single-active uptime ought to be enough for anybody
		NotAfter:              time.Now().Add(262980 * time.Hour),
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, key.Public(), key)
	if err != nil {
		return nil, errwrap.Wrapf("unable to generate local cluster certificate: {{err}}", err)
	}

	return &RaftTLSKey{
		ID:        host,
		KeyType:   certutil.PrivateKeyTypeP521,
		CertBytes: certBytes,
		KeyParams: &certutil.ClusterKeyParams{
			Type: certutil.PrivateKeyTypeP521,
			X:    key.PublicKey.X,
			Y:    key.PublicKey.Y,
			D:    key.D,
		},
		CreatedTime: time.Now(),
	}, nil
}

// Make sure raftLayer satisfies the raft.StreamLayer interface
var _ raft.StreamLayer = (*raftLayer)(nil)

// Make sure raftLayer satisfies the cluster.Handler and cluster.Client
// interfaces
var _ cluster.Handler = (*raftLayer)(nil)
var _ cluster.Client = (*raftLayer)(nil)

// RaftLayer implements the raft.StreamLayer interface,
// so that we can use a single RPC layer for Raft and Vault
type raftLayer struct {
	// Addr is the listener address to return
	addr net.Addr

	// connCh is used to accept connections
	connCh chan net.Conn

	// Tracks if we are closed
	closed    bool
	closeCh   chan struct{}
	closeLock sync.Mutex

	logger log.Logger

	dialerFunc func(string, time.Duration) (net.Conn, error)

	// TLS config
	keyring       *RaftTLSKeyring
	baseTLSConfig *tls.Config
}

// NewRaftLayer creates a new raftLayer object. It parses the TLS information
// from the network config.
func NewRaftLayer(logger log.Logger, raftTLSKeyring *RaftTLSKeyring, clusterAddr net.Addr, baseTLSConfig *tls.Config) (*raftLayer, error) {
	switch {
	case clusterAddr == nil:
		// Clustering disabled on the server, don't try to look for params
		return nil, errors.New("no raft addr found")
	}

	layer := &raftLayer{
		addr:          clusterAddr,
		connCh:        make(chan net.Conn),
		closeCh:       make(chan struct{}),
		logger:        logger,
		baseTLSConfig: baseTLSConfig,
	}

	if err := layer.setTLSKeyring(raftTLSKeyring); err != nil {
		return nil, err
	}

	return layer, nil
}

func (l *raftLayer) setTLSKeyring(keyring *RaftTLSKeyring) error {
	// Fast path a noop update
	if l.keyring != nil && l.keyring.Term == keyring.Term {
		return nil
	}

	for _, key := range keyring.Keys {
		switch {
		case key.KeyParams == nil:
			return errors.New("no raft cluster key params found")

		case key.KeyParams.X == nil, key.KeyParams.Y == nil, key.KeyParams.D == nil:
			return errors.New("failed to parse raft cluster key")

		case key.KeyParams.Type != certutil.PrivateKeyTypeP521:
			return errors.New("failed to find valid raft cluster key type")

		case len(key.CertBytes) == 0:
			return errors.New("no cluster cert found")
		}

		parsedCert, err := x509.ParseCertificate(key.CertBytes)
		if err != nil {
			return errwrap.Wrapf("error parsing raft cluster certificate: {{err}}", err)
		}

		key.parsedCert = parsedCert
		key.parsedKey = &ecdsa.PrivateKey{
			PublicKey: ecdsa.PublicKey{
				Curve: elliptic.P521(),
				X:     key.KeyParams.X,
				Y:     key.KeyParams.Y,
			},
			D: key.KeyParams.D,
		}
	}

	if keyring.ActiveKey == nil {
		return errors.New("expected one active key to be present in the keyring")
	}

	l.keyring = keyring

	return nil
}

func (l *raftLayer) ClientLookup(ctx context.Context, requestInfo *tls.CertificateRequestInfo) (*tls.Certificate, error) {
	for _, subj := range requestInfo.AcceptableCAs {
		for _, key := range l.keyring.Keys {
			if bytes.Equal(subj, key.parsedCert.RawIssuer) {
				localCert := make([]byte, len(key.CertBytes))
				copy(localCert, key.CertBytes)

				return &tls.Certificate{
					Certificate: [][]byte{localCert},
					PrivateKey:  key.parsedKey,
					Leaf:        key.parsedCert,
				}, nil
			}
		}
	}

	return nil, nil
}

func (l *raftLayer) ServerLookup(ctx context.Context, clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	if l.keyring == nil {
		return nil, errors.New("got raft connection but no local cert")
	}

	for _, key := range l.keyring.Keys {
		if clientHello.ServerName == key.ID {
			localCert := make([]byte, len(key.CertBytes))
			copy(localCert, key.CertBytes)

			return &tls.Certificate{
				Certificate: [][]byte{localCert},
				PrivateKey:  key.parsedKey,
				Leaf:        key.parsedCert,
			}, nil
		}
	}

	return nil, nil
}

// CALookup returns the CA to use when validating this connection.
func (l *raftLayer) CALookup(context.Context) ([]*x509.Certificate, error) {
	ret := make([]*x509.Certificate, len(l.keyring.Keys))
	for i, key := range l.keyring.Keys {
		ret[i] = key.parsedCert
	}
	return ret, nil
}

// Stop shutsdown the raft layer.
func (l *raftLayer) Stop() error {
	l.Close()
	return nil
}

// Handoff is used to hand off a connection to the
// RaftLayer. This allows it to be Accept()'ed
func (l *raftLayer) Handoff(ctx context.Context, wg *sync.WaitGroup, quit chan struct{}, conn *tls.Conn) error {
	if l.closed {
		return errors.New("raft is shutdown")
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case l.connCh <- conn:
		case <-l.closeCh:
		case <-ctx.Done():
		case <-quit:
		}
	}()

	return nil
}

// Accept is used to return connection which are
// dialed to be used with the Raft layer
func (l *raftLayer) Accept() (net.Conn, error) {
	select {
	case conn := <-l.connCh:
		return conn, nil
	case <-l.closeCh:
		return nil, fmt.Errorf("Raft RPC layer closed")
	}
}

// Close is used to stop listening for Raft connections
func (l *raftLayer) Close() error {
	l.closeLock.Lock()
	defer l.closeLock.Unlock()

	if !l.closed {
		l.closed = true
		close(l.closeCh)
	}
	return nil
}

// Addr is used to return the address of the listener
func (l *raftLayer) Addr() net.Addr {
	return l.addr
}

// Dial is used to create a new outgoing connection
func (l *raftLayer) Dial(address raft.ServerAddress, timeout time.Duration) (net.Conn, error) {

	tlsConfig := l.baseTLSConfig.Clone()

	key := l.keyring.ActiveKey
	if key == nil {
		return nil, errors.New("no active key")
	}

	tlsConfig.NextProtos = []string{consts.RaftStorageALPN}
	tlsConfig.ServerName = key.parsedCert.Subject.CommonName

	l.logger.Debug("creating rpc dialer", "host", tlsConfig.ServerName)

	pool := x509.NewCertPool()
	pool.AddCert(key.parsedCert)
	tlsConfig.RootCAs = pool
	tlsConfig.ClientCAs = pool

	dialer := &net.Dialer{
		Timeout: timeout,
	}
	return tls.DialWithDialer(dialer, "tcp", string(address), tlsConfig)
}
