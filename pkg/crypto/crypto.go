package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

var (
	// ErrInvalidPEM error
	ErrInvalidPEM = errors.New("invalid PEM encoded key file")
	// ErrInvalidBlockT error
	ErrInvalidBlockT = errors.New("invalid PEM key block type")
)

// PrivateKey struct
type PrivateKey struct {
	key *rsa.PrivateKey
}

// Import to load key from filesystem
func Import(path string) (*PrivateKey, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, ErrInvalidPEM
	}
	if got, exp := block.Type, "RSA PRIVATE KEY"; got != exp {
		return nil, ErrInvalidBlockT
	}
	k, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return &PrivateKey{k}, nil
}

// Generate to generate rsa key
func Generate(opts ...int) (*PrivateKey, error) {
	var bits = 2048
	for _, opt := range opts {
		bits = opt
	}
	private, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	return &PrivateKey{private}, nil
}

// Signer returns private key signer
func (v *PrivateKey) Signer() (ssh.Signer, error) {
	signer, err := ssh.NewSignerFromSigner(v.key)
	if err != nil {
		return nil, err
	}
	return signer, nil
}
