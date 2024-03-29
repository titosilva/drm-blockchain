package handshake

import (
	"drm-blockchain/src/core/protocols/handshake/messages"
	"drm-blockchain/src/core/protocols/identities"
	"drm-blockchain/src/core/repositories/keyrepository"
	"drm-blockchain/src/di"
	"drm-blockchain/src/networking/tunnel"
	"errors"
)

type Executor struct {
	keyRepo keyrepository.IKeyRepository
}

func ExecutorDIFactory(diCtx *di.DIContext) *Executor {
	return New(di.GetInterfaceService[keyrepository.IKeyRepository](diCtx))
}

func New(keyRepo keyrepository.IKeyRepository) *Executor {
	ex := new(Executor)

	ex.keyRepo = keyRepo

	return ex
}

func (ex *Executor) Execute(helloMsg *messages.Hello, tunnel tunnel.Tunnel) {
	_, err := ex.verifyMsg(helloMsg)
	if err != nil {
		return
	}

}

func (ex *Executor) verifyMsg(helloMsg *messages.Hello) (*identities.Identity, error) {
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
