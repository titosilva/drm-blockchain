package tcp

import (
	packet "drm-blockchain/src/networking/transport"
	"errors"
	"net"
)

type Connection struct {
	conn    net.Conn
	closed  bool
	Packets chan packet.Packet
}

func NewConnection(conn net.Conn) *Connection {
	r := new(Connection)
	r.conn = conn
	r.Packets = make(chan packet.Packet)

	return r
}

func (conn *Connection) SendPkt(pkt packet.Packet) error {
	addr := pkt.Addr

	if addr.Network() != "tcp" {
		return errors.New("expected UDP address")
	}

	if conn.conn.RemoteAddr() != addr {
		return errors.New("this connection cannot send packets to the requested address")
	}

	return conn.Send(pkt.Data[:])
}

func (conn *Connection) Send(data []byte) error {
	_, err := conn.conn.Write(data[:])
	return err
}

func (conn *Connection) Close() error {
	conn.closed = true
	return conn.conn.Close()
}
