package udp

import (
	"drm-blockchain/src/collections/structures/concurrent/safemap"
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
	tunnels *safemap.SafeMap[string, *UdpTunnel]
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
	server.tunnels = safemap.New[string, *UdpTunnel]()
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

		pkt, err := tunnel.NewPacket(addr.String(), data[:sz])

		if err != nil {
			return err
		}

		server.Packets <- pkt
		if tunnel, found := server.tunnels.Get(addr.String()); found {
			tunnel.Recv.Notify(pkt)
		}
	}
}

func (server *Server) Send(data []byte, addr net.Addr) error {
	if server.closed {
		panic("Server closed!")
	}

	if addr.Network() != "udp" {
		return errors.New("expected UDP address")
	}

	_, err := server.conn.WriteToUDP(data[:], addr.(*net.UDPAddr))
	return err
}

func (server *Server) Tunnel(addr string) *UdpTunnel {
	if server.closed {
		panic("Server closed!")
	}

	tunnel := NewTunnel(server, addr)
	server.registerTunnel(addr, tunnel)
	return tunnel
}

func (server *Server) registerTunnel(addr string, tunnel *UdpTunnel) {
	if server.closed {
		panic("Server closed!")
	}

	server.tunnels.Set(addr, tunnel)
	go func(addr0 string, tunnel0 *UdpTunnel, server0 *Server) {
		for range tunnel0.WaitClose() {
			server0.unregisterTunnel(addr0)
			break
		}
	}(addr, tunnel, server)
}

func (server *Server) unregisterTunnel(addr string) {
	if server.closed {
		panic("Server closed!")
	}

	server.tunnels.Delete(addr)
}

func (server *Server) Close() {
	close(server.Packets)
	server.conn.Close()
	server.closed = true
}
