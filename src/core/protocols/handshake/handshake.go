package handshake

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"drm-blockchain/src/core/protocols/handshake/messages"
	"drm-blockchain/src/core/protocols/identities"
	"drm-blockchain/src/core/protocols/identities/identitykeys"
	"drm-blockchain/src/core/repositories/keyrepository"
	"drm-blockchain/src/di"
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
	tun := host.udpServer.Tunnel(udpAddr.String())
	tun.Send(data)

	challengePkt := <-tun.ReceivePkt()
	challengeMsg, _ := messages.Decode(challengePkt.Data)
	challenge, _, _ := messages.Disassemble(challengeMsg)

	nonce := challenge.(*messages.Challenge).Nonce

	ephKey, _ := identitykeys.GetECDHCurve().GenerateKey(rand.Reader)
	nodeId, _ := identities.FromAddress(nodeAddr)

	challengeData := append(nonce, ephKey.PublicKey().Bytes()...)
	challengeData = append(challengeData, []byte(nodeAddr)...)

	digest := sha256.New()
	digest.Write(challengeData)
	challengeSum := digest.Sum(nil)

	signature, _ := host.keyRepo.Sign(challengeSum)

	challengeResp := messages.ChallengeResponse{
		EphemeralPubKey: ephKey.PublicKey().Bytes(),
		Signature:       signature,
	}

	challengeRespMsg, _ := messages.Assemble(challengeResp)
	challengeRespData, _ := messages.Encode(challengeRespMsg)
	tun.Send(challengeRespData)

	secret, _ := nodeId.DeriveSecret(ephKey)
	print(secret)
}

func (host *HandshakeHost) listen() {
	if host.closed {
		panic("Handshake host closed!")
	}

	for {
		select {
		case dg := <-host.udpServer.Datagrams:
			go host.processPacket(dg)
		case <-host.cancellation.Done():
			return
		}
	}
}

func (host *HandshakeHost) processPacket(dg udp.Datagram) {
	capsule, err := messages.Decode(dg.Data)

	if err != nil {
		return
	}

	content, typeName, err := messages.Disassemble(capsule)
	if err != nil {
		return
	}

	if typeName == utils.TypeToString[messages.Hello]() {
		tunnel := host.udpServer.Tunnel(dg.Addr.String())
		executor := di.GetService[Executor](host.di)

		executor.Execute(content.(*messages.Hello), tunnel)
	}
}

func (host *HandshakeHost) Close() {
	host.closed = true
	host.udpServer.Close()
}
