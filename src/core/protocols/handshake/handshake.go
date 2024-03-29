package handshake

import (
	"context"
	"drm-blockchain/src/core/protocols/handshake/messages"
	"drm-blockchain/src/di"
	"drm-blockchain/src/networking/tunnel"
	"drm-blockchain/src/networking/udp"
	"drm-blockchain/src/utils"
)

type HandshakeHost struct {
	udpServer    *udp.Server
	closed       bool
	di           *di.DIContext
	cancellation context.Context
}

func Open(addr string, cancellation context.Context, di *di.DIContext) (*HandshakeHost, error) {
	udpServer, err := udp.Open(addr)
	if err != nil {
		return nil, err
	}

	host := new(HandshakeHost)
	host.udpServer = udpServer
	host.cancellation = cancellation
	host.di = di
	go host.listen()

	return host, nil
}

func Greet(addr string) {

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

	content, typeName, err := messages.Disassemble(capsule)
	if err != nil {
		return
	}

	if typeName == utils.TypeToString[messages.Hello]() {
		tunnel := host.udpServer.Tunnel(pkt.Address)
		executor := di.GetService[Executor](host.di)

		executor.Execute(content.(*messages.Hello), tunnel)
	}
}

func (host *HandshakeHost) Close() {
	host.closed = true
	host.udpServer.Close()
}
