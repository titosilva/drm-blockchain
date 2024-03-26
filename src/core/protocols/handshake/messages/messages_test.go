package messages_test

import (
	"drm-blockchain/src/core/protocols/handshake/messages"
	"drm-blockchain/src/utils"
	ez "drm-blockchain/src/utils/test"
	"testing"
)

func Test__AssembleThenDisassemble__Hello(t *testing.T) {
	ez := ez.New(t)

	hello := messages.Hello{
		DestinationAddress: "123",
		SourceAddress:      "456",
	}

	msg, err := messages.Assemble(hello)
	ez.AssertNoError(err)

	decoded, typeName, err := messages.Disassemble(msg)
	ez.AssertNoError(err)
	ez.AssertAreEqual(typeName, utils.TypeToString[messages.Hello]())

	decHello := decoded.(*messages.Hello)
	ez.AssertAreEqual(hello, *decHello)
}

func Test__EncodeThenDecode__Hello(t *testing.T) {
	ez := ez.New(t)

	hello := messages.Hello{
		DestinationAddress: "123",
		SourceAddress:      "456",
	}

	msg, err := messages.Assemble(hello)
	ez.AssertNoError(err)
	bs, err := messages.Encode(msg)
	ez.AssertNoError(err)

	decoded, err := messages.Decode(bs)
	ez.AssertNoError(err)
	recovered, typeName, err := messages.Disassemble(decoded)
	ez.AssertNoError(err)
	ez.AssertAreEqual(typeName, utils.TypeToString[messages.Hello]())

	recHello := recovered.(*messages.Hello)
	ez.AssertAreEqual(hello, *recHello)
}
