package main

import (
	services "drm-blockchain/src/core"
	"fmt"
)

func main() {
	fmt.Println("DRM Blockchain")

	ctx := services.SetupServices()
	services.InitializeServices(ctx)
}
