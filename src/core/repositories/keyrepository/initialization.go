package keyrepository

import (
	"drm-blockchain/src/core/protocols/identities"
)

func (kr *KeyRepository) loadOrCreateSelfIdentity() error {
	exists, err := kr.bs.Exists(selfPrivKeyPath)

	if err != nil {
		return err
	}

	var id *identities.Identity

	if exists {
		id, err = kr.createIdentity()
	} else {
		id, err = kr.loadIdentity()
	}

	if err != nil {
		return err
	}

	kr.id = id
	return nil
}

func (kr KeyRepository) createIdentity() (*identities.Identity, error) {
	id := identities.Generate()
	keyBytes, err := id.ExportPrivateKey()
	if err != nil {
		return nil, err
	}

	err = kr.bs.Write(selfPrivKeyPath, keyBytes)
	if err != nil {
		return nil, err
	}

	return id, err
}

func (kr KeyRepository) loadIdentity() (*identities.Identity, error) {
	keyBytes, err := kr.bs.Get(selfPrivKeyPath)
	if err != nil {
		return nil, err
	}

	id, err := identities.FromPrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}

	return id, nil
}
