package keyrepository

import (
	"drm-blockchain/src/core/blobstore"
	"drm-blockchain/src/core/protocols/identities"
	"drm-blockchain/src/di"
)

type KeyRepository struct {
	// Dependencies
	bs blobstore.BlobStore
	id *identities.Identity
}

const (
	selfPrivKeyPath string = "keys/static-priv"
)

func DIFactory(ctx *di.DIContext) *KeyRepository {
	kr := new(KeyRepository)
	kr.bs = di.GetInterfaceService[blobstore.BlobStore](ctx)
	return kr
}

func (kr *KeyRepository) Initialize() error {
	if err := kr.loadOrCreateSelfIdentity(); err != nil {
		return err
	}

	return nil
}

func (kr *KeyRepository) GetSelfIdentity() *identities.Identity {
	return kr.id
}
