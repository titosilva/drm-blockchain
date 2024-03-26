package handshake

import (
	"context"
	"drm-blockchain/src/core/protocols/handshake/messages"
	packet "drm-blockchain/src/networking/transport"
	"drm-blockchain/src/networking/transport/udp"
	"drm-blockchain/src/utils"
)

type HandshakeHost struct {
	udpServer    *udp.Server
	closed       bool
	cancellation context.Context
}

func Open(addr string, cancellation context.Context) (*HandshakeHost, error) {
	udpServer, err := udp.Open(addr)
	if err != nil {
		return nil, err
	}

	host := new(HandshakeHost)
	host.udpServer = udpServer
	host.cancellation = cancellation
	go host.listen()

	return host, nil
}

func (host *HandshakeHost) listen() {
	if host.closed {
		panic("Handshake host closed!")
	}

	for {
		select {
		case pkt := <-host.udpServer.Packets:
			go host.processPacket(pkt)
		case <-host.cancellation.Done():
			return
		}
	}
}

func (host *HandshakeHost) processPacket(pkt packet.Packet) {
	capsule, err := messages.Decode(pkt.Data)

	if err != nil {
		return
	}

	_, typeName, err := messages.Disassemble(capsule)
	switch typeName {
	case utils.TypeToString[messages.Hello]():
	}
}

func (host *HandshakeHost) Close() {
	host.closed = true
	host.udpServer.Close()
}
