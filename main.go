package main

import (
	"context"
	services "drm-blockchain/src/core"
	"drm-blockchain/src/core/protocols/handshake"
	"fmt"
	"time"
)

func main() {
	fmt.Println("DRM Blockchain")

	diCtx := services.SetupServices()
	services.InitializeServices(diCtx)

	cancellation := context.Background()

	h1Addr := "127.0.0.1:8080"
	h2Addr := "127.0.0.1:8081"
	h1, _ := handshake.NewHost(h1Addr, cancellation, diCtx)
	defer h1.Close()

	h2, _ := handshake.NewHost(h2Addr, cancellation, diCtx)
	defer h2.Close()

	h1.Greet(h2.GetNodeAddress(), h2Addr)

	for {
		time.Sleep(time.Second)
	}
}
