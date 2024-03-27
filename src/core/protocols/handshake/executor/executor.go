package executor

import (
	"drm-blockchain/src/core/protocols/handshake/messages"
	"drm-blockchain/src/core/protocols/identities"
	"drm-blockchain/src/core/repositories/keyrepository"
	"drm-blockchain/src/networking/tunnel"
	"errors"
)

type HandshakeExecutor struct {
	keyRepo *keyrepository.KeyRepository
}

func (ex *HandshakeExecutor) Execute(helloMsg messages.Hello, tunnel *tunnel.Tunnel) {
	_, err := ex.verifyMsg(helloMsg)
	if err != nil {
		return
	}

}

func (ex *HandshakeExecutor) verifyMsg(helloMsg messages.Hello) (*identities.Identity, error) {
	id := ex.keyRepo.GetSelfIdentity()
	if id.GetAddress() != helloMsg.DestinationAddress {
		return nil, errors.New("destination address does not match self address")
	}

	srcId, err := identities.FromAddress(helloMsg.SourceAddress)
	if err != nil {
		return nil, errors.New("provided source address cannot be parsed to identity")
	}

	return srcId, nil
}
