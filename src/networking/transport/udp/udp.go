package udp

import (
	errorutils "drm-blockchain/src/utils/error"
	"fmt"
	"net"
)

type Packet struct {
	Addr *net.UDPAddr
	Data []byte
}

func NewPacket(addr *net.UDPAddr, data []byte) Packet {
	return Packet{
		Addr: addr,
		Data: data,
	}
}

type Server struct {
	conn    *net.UDPConn
	closed  bool
	Addr    *net.UDPAddr
	Channel chan Packet
}

const (
	ServerPacketBufferSize = 256
	PacketSize             = 8192
)

func Open(addr string) (*Server, error) {
	resolved, err := net.ResolveUDPAddr("udp", addr)

	if err != nil {
		return nil, errorutils.NewfWithInner(err, "Failed UDP address resolution for \"%s\"", addr)
	}

	conn, err := net.ListenUDP("udp", resolved)

	if err != nil {
		return nil, errorutils.NewfWithInner(err, "Failed to listen on UDP for \"%s\"", addr)
	}

	server := new(Server)
	server.conn = conn
	server.conn.SetReadBuffer(PacketSize)
	server.Addr = resolved
	server.Channel = make(chan Packet, ServerPacketBufferSize)

	go server.listen()

	return server, nil
}

func (server *Server) listen() error {
	if server.closed {
		fmt.Println("Server closed!")
		panic("Server closed!")
	}

	for {
		var data [PacketSize]byte
		sz, addr, err := server.conn.ReadFromUDP(data[:])

		if err != nil {
			return err
		}

		server.Channel <- NewPacket(addr, data[:sz])
	}
}

func (server *Server) SendPkt(pkt Packet) {
	server.Send(pkt.Data[:], pkt.Addr)
}

func (server *Server) Send(data []byte, addr *net.UDPAddr) {
	server.conn.WriteToUDP(data[:], addr)
}

func (server *Server) Close() {
	close(server.Channel)
	server.conn.Close()
	server.closed = true
}
