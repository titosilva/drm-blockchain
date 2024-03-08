package udppacket

import "net"

const (
	UDPPacketSize int = 8192
)

type UDPPacket struct {
	Addr *net.UDPAddr
	Data []byte
}

func New(addr *net.UDPAddr, data []byte) UDPPacket {
	return UDPPacket{
		Addr: addr,
		Data: data,
	}
}
