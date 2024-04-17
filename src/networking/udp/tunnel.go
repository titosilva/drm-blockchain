package udp

import (
	"drm-blockchain/src/networking/tunnel"
	"drm-blockchain/src/utils/multichannel"
	"net"
)

type UdpTunnel struct {
	server     *Server
	addr       *net.UDPAddr
	Recv       *multichannel.MultiChannel[tunnel.Packet]
	closedChan *multichannel.MultiChannel[any]
	closed     bool
}

func NewTunnel(server *Server, addr string) *UdpTunnel {
	r := new(UdpTunnel)
	r.server = server
	var err error
	r.addr, err = net.ResolveUDPAddr("udp", addr)

	if err != nil {
		panic("cannot resolve address")
	}

	r.Recv = multichannel.New[tunnel.Packet]()
	r.closedChan = multichannel.New[any]()

	return r
}

func (conn *UdpTunnel) SendPkt(pkt tunnel.Packet) error {
	return conn.Send(pkt.Data[:])
}

func (conn *UdpTunnel) Send(data []byte) error {
	if conn.closed {
		panic("tunnel closed!")
	}

	err := conn.server.Send(data[:], conn.addr)
	return err
}

func (conn *UdpTunnel) ReceivePkt() <-chan tunnel.Packet {
	if conn.closed {
		panic("tunnel closed!")
	}

	return conn.Recv.Subscribe()
}

// WaitClose implements tunnel.Tunnel.
func (tunnel *UdpTunnel) WaitClose() <-chan any {
	if tunnel.closed {
		panic("tunnel closed!")
	}

	return tunnel.closedChan.Subscribe()
}

func (conn *UdpTunnel) Close() error {
	if conn.closed {
		return nil
	}

	conn.closed = true
	conn.Recv.Close()
	conn.closedChan.Notify(0)
	conn.closedChan.Close()
	return nil
}

var _ tunnel.Tunnel = new(UdpTunnel)
