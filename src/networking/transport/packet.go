package packet

import (
	"errors"
	"net"
)

const (
	PacketMaxSize = 8192
)

type Packet struct {
	Addr net.Addr
	Data []byte
}

func NewPacket(addr net.Addr, data []byte) (Packet, error) {
	if len(data) > PacketMaxSize {
		return Packet{}, errors.New("max packet size exceeded")
	}

	return Packet{
		Addr: addr,
		Data: data,
	}, nil
}
