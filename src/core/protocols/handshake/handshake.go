package handshake

import "drm-blockchain/src/networking/transport/udp"

type HandshakeHost struct {
	udp *udp.Server
}
