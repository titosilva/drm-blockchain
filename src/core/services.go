package services

import (
	"drm-blockchain/src/core/blobstore"
	"drm-blockchain/src/core/blobstore/localblobstore"
	"drm-blockchain/src/core/protocols/handshake"
	"drm-blockchain/src/core/repositories/keyrepository"
	"drm-blockchain/src/di"
)

func SetupServices() *di.DIContext {
	ctx := di.NewContext()

	// Register singletons
	di.AddInterfaceSingleton[keyrepository.IKeyRepository](ctx, keyrepository.DIFactory)

	// Register factories
	di.AddFactory[handshake.Executor](ctx, handshake.ExecutorDIFactory)
	di.AddInterfaceFactory[blobstore.BlobStore](ctx, localblobstore.DIFactory)

	return ctx
}

func InitializeServices(ctx *di.DIContext) {
	err := di.GetInterfaceService[keyrepository.IKeyRepository](ctx).Initialize()

	if err != nil {
		panic(err)
	}
}
