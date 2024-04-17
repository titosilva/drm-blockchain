package handshake

import (
	"context"
	"drm-blockchain/src/core/protocols/handshake/messages"
	"drm-blockchain/src/core/repositories/keyrepository"
	"drm-blockchain/src/di"
	"drm-blockchain/src/networking/tunnel"
	"drm-blockchain/src/networking/udp"
	"drm-blockchain/src/utils"
	"net"
)

type HandshakeHost struct {
	udpServer    *udp.Server
	closed       bool
	di           *di.DIContext
	cancellation context.Context
	keyRepo      keyrepository.IKeyRepository
}

func Open(addr string, cancellation context.Context, diCtx *di.DIContext) (*HandshakeHost, error) {
	udpServer, err := udp.Open(addr)
	if err != nil {
		return nil, err
	}

	host := new(HandshakeHost)
	host.udpServer = udpServer
	host.cancellation = cancellation
	host.di = diCtx
	host.keyRepo = di.GetInterfaceService[keyrepository.IKeyRepository](diCtx)
	go host.listen()

	return host, nil
}

func (host *HandshakeHost) GetNodeAddress() string {
	return host.keyRepo.GetSelfIdentity().GetAddress()
}

func (host *HandshakeHost) Greet(nodeAddr string, addr string) {
	assembly, _ := messages.Assemble(messages.Hello{
		DestinationAddress: nodeAddr,
		SourceAddress:      host.GetNodeAddress(),
	})

	udpAddr, _ := net.ResolveUDPAddr("udp", addr)
	data, _ := messages.Encode(assembly)
	host.udpServer.Send(data, udpAddr)
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
