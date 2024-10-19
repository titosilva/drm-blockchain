package handshake_test

import (
	"context"
	"drm-blockchain/src/core/protocols/handshake"
	"drm-blockchain/src/di"
	"testing"
)

func Test__Hello(t *testing.T) {
	ctx := context.TODO()
	di := di.NewContext()
	host, err := handshake.Open("localhost:8080", ctx, di)

	if err != nil {
		t.Error(err)
	}

	host.Greet("localhost:8081")
}
