package services

import (
	"drm-blockchain/src/di"
	"drm-blockchain/src/services/blobstore"
	"drm-blockchain/src/services/blobstore/localblobstore"
	"drm-blockchain/src/services/repositories/keyrepository"
)

func SetupServices() *di.DIContext {
	ctx := di.NewContext()

	// Register singletons
	di.AddSingleton[keyrepository.KeyRepository](ctx, keyrepository.DIFactory)

	// Register factories
	di.AddInterfaceFactory[blobstore.BlobStore](ctx, localblobstore.DIFactory)

	return ctx
}

func InitializeServices(ctx *di.DIContext) {
	err := di.GetService[keyrepository.KeyRepository](ctx).Initialize()

	if err != nil {
		panic(err)
	}
}
