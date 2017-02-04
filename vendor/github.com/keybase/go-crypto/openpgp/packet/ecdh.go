package packet

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
)

type ecdhPrivateKey struct {
	ecdsa.PublicKey
	x *big.Int
}

func (e *ecdhPrivateKey) Decrypt(b []byte) ([]byte, error) {
	// TODO(maxtaco): compute the shared secret, run the KDF and
	// recover the decrypted shard key.
	return nil, errors.New("ECDH decrypt unimplemented")
}
