package hostdirectory

import (
	"drm-blockchain/src/di"
	"drm-blockchain/src/networking/tcp"
)

type HostDirectory struct {
	di    *di.DIContext
	known map[string]tcp.TcpTunnel
}

func DIFactory(di *di.DIContext) *HostDirectory {
	hd := NewHostDirectory(di)
	return hd
}

func NewHostDirectory(di *di.DIContext) *HostDirectory {
	hd := new(HostDirectory)
	hd.di = di
	hd.known = make(map[string]tcp.TcpTunnel)
	return hd
}
