package executor

import (
	"drm-blockchain/src/core/protocols/handshake/messages"
	"drm-blockchain/src/core/repositories/keyrepository"
	"net"
)

type HandshakeExecutor struct {
	keyRepo *keyrepository.KeyRepository
}

func (ex *HandshakeExecutor) Start(remote *net.Addr, helloMsg messages.Hello) {

}

func (ex *HandshakeExecutor) verifyMsg(helloMsg messages.Hello) {

}
