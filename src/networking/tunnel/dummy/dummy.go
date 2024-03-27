package dummy

import (
	"drm-blockchain/src/networking/tunnel"
	"drm-blockchain/src/utils/multichannel"
)

type DummyTunnel struct {
	Send  chan tunnel.Packet
	Recv  *multichannel.MultiChannel[tunnel.Packet]
	close *multichannel.MultiChannel[any]
}

func New() *DummyTunnel {
	d := new(DummyTunnel)

	d.Recv = multichannel.New[tunnel.Packet]()
	d.Send = make(chan tunnel.Packet)
	d.close = multichannel.New[any]()

	return d
}

// ReceivePkt implements tunnel.Tunnel.
func (d DummyTunnel) ReceivePkt() <-chan tunnel.Packet {
	return d.Recv.Subscribe()
}

// SendPkt implements tunnel.Tunnel.
func (d DummyTunnel) SendPkt(pkt tunnel.Packet) error {
	d.Send <- pkt
	return nil
}

func (d DummyTunnel) WaitClose() <-chan any {
	return d.close.Subscribe()
}

// Close implements tunnel.Tunnel.
func (d DummyTunnel) Close() error {
	d.Recv.Close()
	close(d.Send)
	d.close.Notify(0)
	return nil
}

var _ tunnel.Tunnel = DummyTunnel{}
