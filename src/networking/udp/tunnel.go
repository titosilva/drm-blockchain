package udp

import (
	"drm-blockchain/src/networking/tunnel"
	"drm-blockchain/src/utils/multichannel"
	"net"
)

type UdpTunnel struct {
	server *Server
	addr   *net.UDPAddr
	Recv   *multichannel.MultiChannel[tunnel.Packet]
	closed bool
}

func NewTunnel(server *Server, addr *net.UDPAddr) *UdpTunnel {
	r := new(UdpTunnel)
	r.server = server
	r.addr = addr

	r.Recv = multichannel.New[tunnel.Packet]()

	return r
}

func (conn *UdpTunnel) SendPkt(pkt tunnel.Packet) error {
	return conn.Send(pkt.Data[:])
}

func (conn *UdpTunnel) Send(data []byte) error {
	err := conn.server.Send(data[:], conn.addr)
	return err
}

func (conn *UdpTunnel) ReceivePkt() <-chan tunnel.Packet {
	return conn.Recv.Subscribe()
}

func (conn *UdpTunnel) Close() error {
	if conn.closed {
		return nil
	}

	conn.closed = true
	conn.Recv.Close()
	return nil
}
