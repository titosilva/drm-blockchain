package tcp

import (
	"drm-blockchain/src/networking/tunnel"
	"net"
)

type TcpTunnel struct {
	conn        net.Conn
	recv        chan tunnel.Packet
	closed      bool
	closed_chan chan any
}

func NewTunnel(conn net.Conn) *TcpTunnel {
	r := new(TcpTunnel)
	r.conn = conn
	r.recv = make(chan tunnel.Packet)

	return r
}

func (conn *TcpTunnel) SendPkt(pkt tunnel.Packet) error {
	return conn.Send(pkt.Data[:])
}

func (conn *TcpTunnel) Send(data []byte) error {
	_, err := conn.conn.Write(data[:])
	return err
}

func (conn *TcpTunnel) ReceivePkt() <-chan tunnel.Packet {
	return conn.recv
}

func (conn *TcpTunnel) Close() error {
	if conn.closed {
		return nil
	}

	conn.closed = true
	conn.closed_chan <- true
	return conn.conn.Close()
}

func (conn *TcpTunnel) WaitClose() <-chan any {
	return conn.closed_chan
}

// Static implementation check
var _ tunnel.Tunnel = (*TcpTunnel)(nil)
