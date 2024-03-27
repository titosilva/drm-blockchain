package tunnel

type Tunnel interface {
	SendPkt(Packet) error
	ReceivePkt() <-chan Packet
	WaitClose() <-chan any
	Close() error
}
