package tunnel

import (
	"errors"
)

const (
	PacketMaxSize = 8192
)

type Packet struct {
	Address string
	Data    []byte
}

func NewPacket(addr string, data []byte) (Packet, error) {
	if len(data) > PacketMaxSize {
		return Packet{}, errors.New("max packet size exceeded")
	}

	return Packet{
		Data: data,
	}, nil
}
