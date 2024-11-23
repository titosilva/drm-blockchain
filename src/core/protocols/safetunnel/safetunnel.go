package safetunnel

import (
	"drm-blockchain/src/networking/tcp"
	"drm-blockchain/src/networking/tunnel"
)

// SafeTunnel implements Tunnel
type SafeTunnel struct {
	unsafe_tunnel tcp.TcpTunnel
}

// Close implements tunnel.Tunnel.
func (s *SafeTunnel) Close() error {
	return s.unsafe_tunnel.Close()
}

// ReceivePkt implements tunnel.Tunnel.
func (s *SafeTunnel) ReceivePkt() <-chan tunnel.Packet {
	return s.unsafe_tunnel.ReceivePkt()
}

// SendPkt implements tunnel.Tunnel.
func (s *SafeTunnel) SendPkt(pkt tunnel.Packet) error {
	return s.unsafe_tunnel.SendPkt(pkt)
}

// WaitClose implements tunnel.Tunnel.
func (s *SafeTunnel) WaitClose() <-chan any {
	return s.unsafe_tunnel.WaitClose()
}

// Static implementation check
var _ tunnel.Tunnel = (*SafeTunnel)(nil)
