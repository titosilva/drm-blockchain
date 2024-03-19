package udpserver

import (
	udppacket "drm-blockchain/src/networking/transport"
	errorutils "drm-blockchain/src/utils/error"
	"fmt"
	"net"
)

type UDPServer struct {
	conn    *net.UDPConn
	closed  bool
	Addr    *net.UDPAddr
	Channel chan udppacket.UDPPacket
}

const (
	UDPServerPacketBufferSize = 256
)

func Open(addr string) (*UDPServer, error) {
	resolved, err := net.ResolveUDPAddr("udp", addr)

	if err != nil {
		return nil, errorutils.NewfWithInner(err, "Failed UDP address resolution for \"%s\"", addr)
	}

	conn, err := net.ListenUDP("udp", resolved)

	if err != nil {
		return nil, errorutils.NewfWithInner(err, "Failed to listen on UDP for \"%s\"", addr)
	}

	server := new(UDPServer)
	server.conn = conn
	server.conn.SetReadBuffer(udppacket.UDPPacketSize)
	server.Addr = resolved
	server.Channel = make(chan udppacket.UDPPacket, UDPServerPacketBufferSize)

	go server.listen()

	return server, nil
}

func (server *UDPServer) listen() error {
	if server.closed {
		fmt.Println("Server closed!")
		panic("Server closed!")
	}

	for {
		var data [udppacket.UDPPacketSize]byte
		sz, addr, err := server.conn.ReadFromUDP(data[:])

		if err != nil {
			return err
		}

		server.Channel <- udppacket.New(addr, data[:sz])
	}
}

func (server *UDPServer) SendPkt(pkt udppacket.UDPPacket) {
	server.Send(pkt.Data[:], pkt.Addr)
}

func (server *UDPServer) Send(data []byte, addr *net.UDPAddr) {
	server.conn.WriteToUDP(data[:], addr)
}

func (server *UDPServer) Close() {
	close(server.Channel)
	server.conn.Close()
	server.closed = true
}
