package handshake_test

import (
	"context"
	"drm-blockchain/src/core/blobstore/localblobstore"
	"drm-blockchain/src/core/protocols/handshake"
	"drm-blockchain/src/core/repositories/keyrepository"
	"drm-blockchain/src/di"
	"testing"
)

func Test__Hello(t *testing.T) {
	cancellation := context.TODO()

	diCtx := di.NewContext()
	di.AddInterfaceSingleton(diCtx, keyrepository.DIFactory)
	di.AddInterfaceSingleton(diCtx, localblobstore.DIFactory)
	di.AddSingleton(diCtx, handshake.ExecutorDIFactory)

	host1, err := handshake.NewHost("localhost:8080", cancellation, diCtx)

	if err != nil {
		t.Error(err)
	}

	host2, err := handshake.NewHost("localhost:8081", cancellation, diCtx)
	host2_addr := host2.GetNodeAddress()

	if err != nil {
		t.Error(err)
	}

	secret := host1.Greet(host2_addr, "localhost:8081")
	if secret == "" {
		t.Error("Handshake failed")
	}
}
