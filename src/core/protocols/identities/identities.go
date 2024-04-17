package identities

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/rand"
	"drm-blockchain/src/core/protocols/identities/address"
	"drm-blockchain/src/core/protocols/identities/identitykeys"
	"errors"
)

type Identity struct {
	publicKey *ecdsa.PublicKey

	// Will be filled only when referring to their own address
	// and will be nil when referring to other addresses
	privateKey *ecdsa.PrivateKey
}

func Generate() *Identity {
	privKey := identitykeys.GeneratePrivateKey()

	id := new(Identity)
	id.loadPrivateKey(privKey)

	return id
}

func FromPrivateKey(privKeyBs []byte) (*Identity, error) {
	privKey, err := identitykeys.DecodeIdentityPrivateKey(privKeyBs)

	if err != nil {
		return nil, err
	}

	id := new(Identity)
	id.loadPrivateKey(privKey)

	return id, nil
}

func FromAddress(addr string) (*Identity, error) {
	pubKey, err := address.ComputePublicKeyFromAddress(addr)

	if err != nil {
		return nil, err
	}

	id := new(Identity)
	id.loadPublicKey(pubKey)

	return id, nil
}

func (id *Identity) loadPrivateKey(privKey *ecdsa.PrivateKey) {
	id.privateKey = privKey
	id.publicKey = &privKey.PublicKey
}

func (id *Identity) loadPublicKey(pubKey *ecdsa.PublicKey) {
	id.privateKey = nil
	id.publicKey = pubKey
}

func (id *Identity) GetAddress() string {
	return address.ComputeAddressFromPublicKey(id.publicKey)
}

func (id *Identity) ExportPrivateKey() ([]byte, error) {
	if id.privateKey == nil {
		return nil, errors.New("this identity does not have a private key loaded")
	}

	return identitykeys.EncodeIdentityPrivateKey(id.privateKey)
}

func (id *Identity) Sign(data []byte) ([]byte, error) {
	if id.privateKey == nil {
		return nil, errors.New("this identity does not have a private key loaded")
	}

	return ecdsa.SignASN1(rand.Reader, id.privateKey, data)
}

func (id *Identity) DeriveSecret(key *ecdh.PrivateKey) ([]byte, error) {
	pub, err := id.publicKey.ECDH()

	if err != nil {
		return nil, err
	}

	return key.ECDH(pub)
}
