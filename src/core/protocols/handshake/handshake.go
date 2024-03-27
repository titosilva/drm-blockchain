package handshake

import (
	"context"
	"drm-blockchain/src/core/protocols/handshake/executor"
	"drm-blockchain/src/core/protocols/handshake/messages"
	"drm-blockchain/src/networking/tunnel"
	"drm-blockchain/src/networking/udp"
	"drm-blockchain/src/utils"
)

type HandshakeHost struct {
	udpServer    *udp.Server
	closed       bool
	cancellation context.Context
	executors    map[string]*executor.HandshakeExecutor
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

func (host *HandshakeHost) processPacket(pkt tunnel.Packet) {
	capsule, err := messages.Decode(pkt.Data)

	if err != nil {
		return
	}

	_, typeName, err := messages.Disassemble(capsule)
	if err != nil {
		return
	}

	if typeName == utils.TypeToString[messages.Hello]() {
		tunnel = host.udpServer.Tunnel()
	}
}

func (host *HandshakeHost) Close() {
	host.closed = true
	host.udpServer.Close()
}
