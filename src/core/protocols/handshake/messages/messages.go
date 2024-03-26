package messages

import (
	"drm-blockchain/src/utils"
	errorutils "drm-blockchain/src/utils/error"
	"encoding/asn1"
)

type HandshakeMessageCapsule struct {
	TypeName string
	Content  []byte
}

func Assemble(msg any) (HandshakeMessageCapsule, error) {
	bs, err := asn1.Marshal(msg)

	if err != nil {
		return HandshakeMessageCapsule{}, err
	}

	return HandshakeMessageCapsule{
		TypeName: utils.TypeOf(msg),
		Content:  bs,
	}, nil
}

func Disassemble(capsule HandshakeMessageCapsule) (content any, typeName string, err error) {
	switch capsule.TypeName {
	case utils.TypeToString[Hello]():
		content = new(Hello)
	default:
		err = errorutils.Newf("unkwnown prelude '%s' on handshake", capsule.TypeName)
	}

	if err != nil {
		return nil, "", err
	}

	_, err = asn1.Unmarshal(capsule.Content, content)
	return content, capsule.TypeName, err
}

func Encode(hmsg HandshakeMessageCapsule) ([]byte, error) {
	return asn1.Marshal(hmsg)
}

func Decode(bs []byte) (HandshakeMessageCapsule, error) {
	hmsg := new(HandshakeMessageCapsule)
	_, err := asn1.Unmarshal(bs, hmsg)

	if err != nil {
		return HandshakeMessageCapsule{}, err
	}

	return *hmsg, err
}

type Hello struct {
	DestinationAddress string
	SourceAddress      string
}
