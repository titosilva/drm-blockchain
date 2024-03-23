package messages

import (
	"drm-blockchain/src/utils"
	errorutils "drm-blockchain/src/utils/error"
	"encoding/asn1"
)

type HandshakeMessage struct {
	Prelude string
	Content []byte
}

const PreludeLength = 10

func Assemble(msg any) (HandshakeMessage, error) {
	bs, err := asn1.Marshal(msg)

	if err != nil {
		return HandshakeMessage{}, err
	}

	return HandshakeMessage{
		Prelude: utils.TypeOf(msg),
		Content: bs,
	}, nil
}

func Disassemble(msg HandshakeMessage) (content any, typeName string, err error) {
	switch msg.Prelude {
	case utils.TypeToString[Hello]():
		content = new(Hello)
	default:
		err = errorutils.Newf("unkwnown prelude '%s' on handshake", msg.Prelude)
	}

	if err != nil {
		return nil, "", err
	}

	_, err = asn1.Unmarshal(msg.Content, content)
	return content, msg.Prelude, err
}

func Encode(hmsg HandshakeMessage) ([]byte, error) {
	return asn1.Marshal(hmsg)
}

func Decode(bs []byte) (HandshakeMessage, error) {
	hmsg := new(HandshakeMessage)
	_, err := asn1.Unmarshal(bs, hmsg)

	if err != nil {
		return HandshakeMessage{}, err
	}

	return *hmsg, err
}

type Hello struct {
	DestinationAddress string
	SourceAddress      string
	SourcePublicKey    []byte
}
