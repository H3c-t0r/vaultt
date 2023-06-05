// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package transit

import (
	"context"
	"reflect"

	"github.com/hashicorp/vault/sdk/logical"

	"testing"
)

func TestTransit_SignCSR(t *testing.T) {
	// NOTE: Use an existing CSR or generate one here?
	templateCsr := `
-----BEGIN CERTIFICATE REQUEST-----
MIICRTCCAS0CAQAwADCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAM49
McW7u3ILuAJfSFLUtGOMGBytHmMFcjTiX+5JcajFj0Uszb+HQ7eIsJJNXhVc/7fg
Z01DZvcCqb9ChEWE3xi4GEkPMXay7p7G1ooSLnQp6Z0lL5CuIFfMVOTvjfhTwRaJ
l9v2mMlm80BeiAUBqeoyGVrIh5fKASxaE0jrhjAxhGzqrXdDnL8A4na6ArprV4iS
aEAziODd2WmplSKgUwEaFdeG1t1bJf3o5ZQRCnKNtQcAk8UmgtvFEO8ohGMln/Fj
O7u7s6iRhOGf1g1NCAP5pGqxNx3bjz5f/CUcTSIGAReEomg41QTIhD9muCTL8qnm
6lS87wkGTv7qbeIGB7sCAwEAAaAAMA0GCSqGSIb3DQEBCwUAA4IBAQAfjE+jNqIk
4V1tL3g5XPjxr2+QcwddPf8opmbAzgt0+TiIHcDGBAxsXyi7sC9E5AFfFp7W07Zv
r5+v4i529K9q0BgGtHFswoEnhd4dC8Ye53HtSoEtXkBpZMDrtbS7eZa9WccT6zNx
4taTkpptZVrmvPj+jLLFkpKJJ3d+Gbrp6hiORPadT+igLKkqvTeocnhOdAtt427M
RXTVgN14pV3tqO+5MXzNw5tGNPcwWARWwPH9eCRxLwLUuxE4Qu73pUeEFjDEfGkN
iBnlTsTXBOMqSGryEkmRaZslWDvblvYeObYw+uc3kCbJ7jRy9soVwkbb5FueF/yC
O1aQIm23HrrG
-----END CERTIFICATE REQUEST-----
`

	testTransit_SignCSR(t, "rsa-2048", templateCsr)
	testTransit_SignCSR(t, "rsa-3072", templateCsr)
	testTransit_SignCSR(t, "rsa-4096", templateCsr)
	testTransit_SignCSR(t, "ecdsa-p256", templateCsr)
	testTransit_SignCSR(t, "ecdsa-p384", templateCsr)
	testTransit_SignCSR(t, "ecdsa-p521", templateCsr)
	testTransit_SignCSR(t, "ed25519", templateCsr)
	testTransit_SignCSR(t, "aes256-gcm96", templateCsr)
}

func testTransit_SignCSR(t *testing.T, keyType, pemTemplateCsr string) {
	var resp *logical.Response
	var err error
	b, s := createBackendWithStorage(t)

	// Create the policy
	policyReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "keys/test-key",
		Storage:   s,
		Data: map[string]interface{}{
			"type": keyType,
		},
	}
	resp, err = b.HandleRequest(context.Background(), policyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}

	csrSignReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "keys/test-key/csr",
		Storage:   s,
		Data: map[string]interface{}{
			"csr": pemTemplateCsr,
		},
	}

	resp, err = b.HandleRequest(context.Background(), csrSignReq)

	switch keyType {
	case "rsa-2048", "rsa-3072", "rsa-4096", "ecdsa-p256", "ecdsa-p384", "ecdsa-p521", "ed25519":
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("failed to sign CSR, err:%v resp:%#v", err, resp)
		}

		signedCsrBytes, ok := resp.Data["csr"]
		if !ok {
			t.Fatal("expected response data to hold a 'csr' field")
		}

		signedCsr, err := parseCsr(string(signedCsrBytes.([]byte)))
		if err != nil {
			t.Errorf("failed to parse returned csr, err:%v", err)
		}

		templateCsr, err := parseCsr(pemTemplateCsr)
		if err != nil {
			t.Errorf("failed to parse returned template csr, err:%v", err)
		}

		// NOTE: Check other fields?
		if !reflect.DeepEqual(signedCsr.Subject, templateCsr.Subject) {
			t.Errorf("subjects should have matched, err:%v", err)
		}

	default:
		if err == nil || (resp != nil && !resp.IsError()) {
			t.Fatalf("should have failed to sign CSR, provided key type does not support signing")
		}
	}
}

func testSetCertificateChain(t *testing.T) {
}
