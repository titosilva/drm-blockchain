package main

import (
	"drm-blockchain/src/services"
	"fmt"
)

func main() {
	fmt.Println("DRM Blockchain")

	ctx := services.SetupServices()
	services.InitializeServices(ctx)
}
