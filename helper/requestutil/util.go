package requestutil

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/url"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/vault/helper/compressutil"
	"github.com/hashicorp/vault/helper/jsonutil"
)

type bufCloser struct {
	*bytes.Buffer
}

func (b bufCloser) Close() error {
	b.Reset()
	return nil
}

// GenerateForwardedRequest generates a new http.Request that contains the
// original requests's information in the new request's body.
func GenerateForwardedRequest(req *http.Request, addr string) (*http.Request, error) {
	fq := ForwardedRequest{
		Method:        req.Method,
		HeaderEntries: make(map[string]*HeaderEntry, len(req.Header)),
		Host:          req.Host,
		RemoteAddr:    req.RemoteAddr,
	}

	reqURL := req.URL
	fq.Url = &URL{
		Scheme:   reqURL.Scheme,
		Opaque:   reqURL.Opaque,
		Host:     reqURL.Host,
		Path:     reqURL.Path,
		RawPath:  reqURL.RawPath,
		RawQuery: reqURL.RawQuery,
		Fragment: reqURL.Fragment,
	}

	for k, v := range req.Header {
		fq.HeaderEntries[k] = &HeaderEntry{
			Values: v,
		}
	}

	if req.TLS != nil && req.TLS.PeerCertificates != nil && len(req.TLS.PeerCertificates) > 0 {
		fq.PeerCertificates = make([][]byte, len(req.TLS.PeerCertificates))
		for i, cert := range req.TLS.PeerCertificates {
			fq.PeerCertificates[i] = cert.Raw
		}
	}

	buf := bytes.NewBuffer(nil)
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		return nil, err
	}
	fq.Body = buf.Bytes()

	var newBody []byte
	switch os.Getenv("VAULT_MESSAGE_TYPE") {
	case "json":
		newBody, err = jsonutil.EncodeJSON(&fq)
	case "json_compress":
		newBody, err = jsonutil.EncodeJSONAndCompress(&fq, &compressutil.CompressionConfig{
			Type: compressutil.CompressionTypeLzw,
		})
	case "proto3":
		fallthrough
	default:
		newBody, err = proto.Marshal(&fq)
	}
	if err != nil {
		return nil, err
	}

	ret, err := http.NewRequest("POST", addr, bytes.NewBuffer(newBody))
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// ParseForwardedRequest generates a new http.Request that is comprised of the
// values in the given request's body, assuming it correctly parses into a
// ForwardedRequest.
func ParseForwardedRequest(req *http.Request) (*http.Request, error) {
	buf := bufCloser{
		Buffer: bytes.NewBuffer(nil),
	}
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		return nil, err
	}

	fq := new(ForwardedRequest)
	switch os.Getenv("VAULT_MESSAGE_TYPE") {
	case "json", "json_compress":
		err = jsonutil.DecodeJSON(buf.Bytes(), fq)
	default:
		err = proto.Unmarshal(buf.Bytes(), fq)
	}
	if err != nil {
		return nil, err
	}

	buf.Reset()
	_, err = buf.Write(fq.Body)
	if err != nil {
		return nil, err
	}

	ret := &http.Request{
		Method:     fq.Method,
		Header:     make(map[string][]string, len(fq.HeaderEntries)),
		Body:       buf,
		Host:       fq.Host,
		RemoteAddr: fq.RemoteAddr,
	}

	ret.URL = &url.URL{
		Scheme:   fq.Url.Scheme,
		Opaque:   fq.Url.Opaque,
		Host:     fq.Url.Host,
		Path:     fq.Url.Path,
		RawPath:  fq.Url.RawPath,
		RawQuery: fq.Url.RawQuery,
		Fragment: fq.Url.Fragment,
	}

	for k, v := range fq.HeaderEntries {
		ret.Header[k] = v.Values
	}

	if fq.PeerCertificates != nil && len(fq.PeerCertificates) > 0 {
		ret.TLS = &tls.ConnectionState{
			PeerCertificates: make([]*x509.Certificate, len(fq.PeerCertificates)),
		}
		for i, certBytes := range fq.PeerCertificates {
			cert, err := x509.ParseCertificate(certBytes)
			if err != nil {
				return nil, err
			}
			ret.TLS.PeerCertificates[i] = cert
		}
	}

	return ret, nil
}
