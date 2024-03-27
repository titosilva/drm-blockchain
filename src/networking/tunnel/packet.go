package tunnel

import (
	"errors"
)

const (
	PacketMaxSize = 8192
)

type Packet struct {
	Data []byte
}

func NewPacket(data []byte) (Packet, error) {
	if len(data) > PacketMaxSize {
		return Packet{}, errors.New("max packet size exceeded")
	}

	return Packet{
		Data: data,
	}, nil
}
