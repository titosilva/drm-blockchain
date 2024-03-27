package udp

import (
	"drm-blockchain/src/networking/tunnel"
	errorutils "drm-blockchain/src/utils/error"
	"errors"
	"net"
)

type Server struct {
	conn    *net.UDPConn
	closed  bool
	Addr    *net.UDPAddr
	Packets chan tunnel.Packet
	tunnels map[string]*UdpTunnel
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
	server.conn.SetReadBuffer(tunnel.PacketMaxSize)
	server.Addr = resolved
	server.Packets = make(chan tunnel.Packet, 256)
	server.closed = false

	go server.listen()

	return server, nil
}

func (server *Server) listen() error {
	if server.closed {
		panic("Server closed!")
	}

	for {
		var data [tunnel.PacketMaxSize]byte
		sz, addr, err := server.conn.ReadFromUDP(data[:])

		if err != nil {
			return err
		}

		pkt, err := tunnel.NewPacket(data[:sz])

		if err != nil {
			return err
		}

		server.Packets <- pkt
		if tunnel, found := server.tunnels[addr.String()]; found {
			tunnel.Recv.Notify(pkt)
		}
	}
}

func (server *Server) Send(data []byte, addr net.Addr) error {
	if addr.Network() != "udp" {
		return errors.New("expected UDP address")
	}

	_, err := server.conn.WriteToUDP(data[:], addr.(*net.UDPAddr))
	return err
}

func (server *Server) Tunnel(addr *net.UDPAddr) *UdpTunnel {
	return NewTunnel(server, addr)
}

func (server *Server) Close() {
	close(server.Packets)
	server.conn.Close()
	server.closed = true
}
