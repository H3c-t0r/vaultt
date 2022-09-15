package cert

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestCRLFetch(t *testing.T) {
	storage := &logical.InmemStorage{}

	lb, err := Factory(context.Background(), &logical.BackendConfig{
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: 300 * time.Second,
			MaxLeaseTTLVal:     1800 * time.Second,
		},
		StorageView: storage,
	})

	require.NoError(t, err)
	b := lb.(*backend)
	closeChan := make(chan bool)
	go func() {
		t := time.NewTicker(50 * time.Millisecond)
		for {
			select {
			case <-t.C:
				b.PeriodicFunc(context.Background(), &logical.Request{Storage: storage})
			case <-closeChan:
				break
			}
		}
	}()
	defer close(closeChan)

	if err != nil {
		t.Fatalf("error: %s", err)
	}
	connState, err := testConnState("test-fixtures/keys/cert.pem",
		"test-fixtures/keys/key.pem", "test-fixtures/root/rootcacert.pem")
	require.NoError(t, err)
	caPEM, err := ioutil.ReadFile("test-fixtures/root/rootcacert.pem")
	require.NoError(t, err)
	caKeyPEM, err := ioutil.ReadFile("test-fixtures/keys/key.pem")
	require.NoError(t, err)
	certPEM, err := ioutil.ReadFile("test-fixtures/keys/cert.pem")

	caBundle, err := certutil.ParsePEMBundle(string(caPEM))
	require.NoError(t, err)
	bundle, err := certutil.ParsePEMBundle(string(certPEM) + "\n" + string(caKeyPEM))
	require.NoError(t, err)
	//  Entry with one cert first

	revocationListTemplate := &x509.RevocationList{
		RevokedCertificates: []pkix.RevokedCertificate{
			{
				SerialNumber:   big.NewInt(1),
				RevocationTime: time.Now(),
			},
		},
		Number:             big.NewInt(1),
		ThisUpdate:         time.Now(),
		NextUpdate:         time.Now().Add(50 * time.Millisecond),
		SignatureAlgorithm: x509.SHA1WithRSA,
	}

	crlBytes, err := x509.CreateRevocationList(rand.Reader, revocationListTemplate, caBundle.Certificate, bundle.PrivateKey)
	require.NoError(t, err)

	var serverURL *url.URL
	crlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host == serverURL.Host {
			w.Write(crlBytes)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	serverURL, _ = url.Parse(crlServer.URL)

	req := &logical.Request{
		Connection: &logical.Connection{
			ConnState: &connState,
		},
		Storage: storage,
		Auth:    &logical.Auth{},
	}

	fd := &framework.FieldData{
		Raw: map[string]interface{}{
			"name":        "test",
			"certificate": string(caPEM),
			"policies":    "foo,bar",
		},
		Schema: pathCerts(b).Fields,
	}

	resp, err := b.pathCertWrite(context.Background(), req, fd)
	if err != nil {
		t.Fatal(err)
	}

	empty_login_fd := &framework.FieldData{
		Raw:    map[string]interface{}{},
		Schema: pathLogin(b).Fields,
	}
	resp, err = b.pathLogin(context.Background(), req, empty_login_fd)
	if err != nil {
		t.Fatal(err)
	}
	if resp.IsError() {
		t.Fatalf("got error: %#v", *resp)
	}

	// Set a bad CRL
	fd = &framework.FieldData{
		Raw: map[string]interface{}{
			"name": "testcrl",
			"url":  "http://wrongserver.com",
		},
		Schema: pathCRLs(b).Fields,
	}
	resp, err = b.pathCRLWrite(context.Background(), req, fd)
	if err == nil {
		t.Fatal(err)
	}
	if resp.IsError() {
		t.Fatalf("got error: %#v", *resp)
	}

	// Set good CRL
	fd = &framework.FieldData{
		Raw: map[string]interface{}{
			"name": "testcrl",
			"url":  crlServer.URL,
		},
		Schema: pathCRLs(b).Fields,
	}
	resp, err = b.pathCRLWrite(context.Background(), req, fd)
	if err != nil {
		t.Fatal(err)
	}
	if resp.IsError() {
		t.Fatalf("got error: %#v", *resp)
	}

	if len(b.crls["testcrl"].Serials) != 1 {
		t.Fatalf("wrong number of certs in CRL")
	}

	// Add a cert to the CRL, then wait to see if it gets automatically picked up
	revocationListTemplate.RevokedCertificates = []pkix.RevokedCertificate{
		{
			SerialNumber:   big.NewInt(1),
			RevocationTime: revocationListTemplate.RevokedCertificates[0].RevocationTime,
		},
		{
			SerialNumber:   big.NewInt(2),
			RevocationTime: time.Now(),
		},
	}
	revocationListTemplate.ThisUpdate = time.Now()
	revocationListTemplate.NextUpdate = time.Now().Add(1 * time.Minute)
	revocationListTemplate.Number = big.NewInt(2)

	crlBytes, err = x509.CreateRevocationList(rand.Reader, revocationListTemplate, caBundle.Certificate, bundle.PrivateKey)
	require.NoError(t, err)
	time.Sleep(60 * time.Millisecond)
	if len(b.crls["testcrl"].Serials) != 2 {
		t.Fatalf("wrong number of certs in CRL")
	}
}
