package keyrepository

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"drm-blockchain/src/core/blobstore"
	"drm-blockchain/src/di"
)

type KeyRepository struct {
	// Dependencies
	bs blobstore.BlobStore

	ecdsaStaticPrivateKey *ecdsa.PrivateKey
	ecdsaStaticPublicKey  *ecdsa.PublicKey
}

const (
	staticSignaturePrivKeyPath string = "keys/static-priv"
)

func DIFactory(ctx *di.DIContext) *KeyRepository {
	kr := new(KeyRepository)
	kr.bs = di.GetInterfaceService[blobstore.BlobStore](ctx)
	return kr
}

func (kr *KeyRepository) Initialize() error {
	if err := kr.loadOrCreateStaticKeys(); err != nil {
		return err
	}

	return nil
}

func (kr *KeyRepository) loadOrCreateStaticKeys() error {
	exists, err := kr.bs.Exists(staticSignaturePrivKeyPath)

	if err != nil {
		return err
	}

	var priv *ecdsa.PrivateKey

	if exists {
		priv, err = kr.loadStaticKeys()
	} else {
		priv, err = kr.createStaticSignatureKeys()
	}

	if err != nil {
		return err
	}

	kr.ecdsaStaticPrivateKey = priv
	kr.ecdsaStaticPublicKey = &priv.PublicKey
	return nil
}

func (kr KeyRepository) createStaticSignatureKeys() (*ecdsa.PrivateKey, error) {
	priv, err := ecdsa.GenerateKey(kr.staticKeysCurve(), rand.Reader)
	if err != nil {
		return nil, err
	}

	keyBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, err
	}

	err = kr.bs.Write(staticSignaturePrivKeyPath, keyBytes)
	if err != nil {
		return nil, err
	}

	return priv, err
}

func (kr KeyRepository) loadStaticKeys() (*ecdsa.PrivateKey, error) {
	keyBytes, err := kr.bs.Get(staticSignaturePrivKeyPath)
	if err != nil {
		return nil, err
	}

	key, err := x509.ParsePKCS8PrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}

	return key.(*ecdsa.PrivateKey), nil
}

func (kr KeyRepository) staticKeysCurve() elliptic.Curve {
	return elliptic.P256()
}
