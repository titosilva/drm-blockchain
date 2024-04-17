package tunnel_test

import (
	"drm-blockchain/src/networking/tunnel"
	"drm-blockchain/src/networking/tunnel/dummy"
	ez "drm-blockchain/src/utils/test"
	"testing"
)

func Test__Receiver(t *testing.T) {
	ez := ez.New(t)
	d := dummy.New()
	pkt, _ := tunnel.NewPacket([]byte{1, 2, 3, 4, 5})

	c := d.ReceivePkt()
	go func() {
		d.Recv.Notify(pkt)
	}()

	p := <-c
	ez.AssertAreEqual(p, pkt)
}
