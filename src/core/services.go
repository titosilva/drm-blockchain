package services

import (
	"drm-blockchain/src/core/blobstore"
	"drm-blockchain/src/core/blobstore/localblobstore"
	"drm-blockchain/src/core/repositories/keyrepository"
	"drm-blockchain/src/di"
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
