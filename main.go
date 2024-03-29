package main

import (
	"context"
	services "drm-blockchain/src/core"
	"drm-blockchain/src/core/protocols/handshake"
	"fmt"
)

func main() {
	fmt.Println("DRM Blockchain")

	diCtx := services.SetupServices()
	services.InitializeServices(diCtx)

	cancellation := context.Background()
	handshake.Open("127.0.0.1:8080", cancellation, diCtx)
}
