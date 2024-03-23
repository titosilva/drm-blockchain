package udp

import (
	packet "drm-blockchain/src/networking/transport"
	errorutils "drm-blockchain/src/utils/error"
	"errors"
	"net"
)

type Server struct {
	conn    *net.UDPConn
	closed  bool
	Addr    *net.UDPAddr
	Packets chan packet.Packet
}

func Open(addr string) (*Server, error) {
	resolved, err := net.ResolveUDPAddr("udp", addr)

	if err != nil {
		return nil, errorutils.NewfWithInner(err, "failed UDP address resolution for \"%s\"", addr)
	}

	conn, err := net.ListenUDP("udp", resolved)

	if err != nil {
		return nil, errorutils.NewfWithInner(err, "failed to listen on UDP for \"%s\"", addr)
	}

	server := new(Server)
	server.conn = conn
	server.conn.SetReadBuffer(packet.PacketMaxSize)
	server.Addr = resolved
	server.Packets = make(chan packet.Packet)

	go server.listen()

	return server, nil
}

func (server *Server) listen() error {
	if server.closed {
		panic("Server closed!")
	}

	for {
		var data [packet.PacketMaxSize]byte
		sz, addr, err := server.conn.ReadFromUDP(data[:])

		if err != nil {
			return err
		}

		pkt, err := packet.NewPacket(addr, data[:sz])

		if err != nil {
			return err
		}

		server.Packets <- pkt
	}
}

func (server *Server) SendPkt(pkt packet.Packet) error {
	return server.Send(pkt.Data[:], pkt.Addr)
}

func (server *Server) Send(data []byte, addr net.Addr) error {
	if addr.Network() != "udp" {
		return errors.New("expected UDP address")
	}

	_, err := server.conn.WriteToUDP(data[:], addr.(*net.UDPAddr))
	return err
}

func (server *Server) Close() {
	close(server.Packets)
	server.conn.Close()
	server.closed = true
}
