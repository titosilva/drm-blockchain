package udp_test

import (
	"drm-blockchain/src/networking/transport/udp"
	"net"
	"testing"
)

const (
	testMsg   = "any sufficiently advanced technology is indistinguishable from magic"
	testAddr1 = "127.0.0.1:45000"
	testAddr2 = "127.0.0.1:45001"
)

func sendUdpMsgTo(msg string, addr string) {
	// Resolve the string address to a UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", addr)

	if err != nil {
		panic("failed to resolve address on udp send")
	}

	// Dial to the address with UDP
	conn, err := net.DialUDP("udp", nil, udpAddr)

	if err != nil {
		panic("failed to open connection on udp send")
	}

	// Send a message to the server
	_, err = conn.Write([]byte(msg))
	if err != nil {
		panic("failed to send message")
	}
}

func Test__UDPServer__ShouldReceiveData__OnLoopbackAddr(t *testing.T) {
	server, err := udp.Open(testAddr1)
	if err != nil {
		t.Error(err)
	}

	defer server.Close()

	sendUdpMsgTo(testMsg, testAddr1)
	m := <-server.Channel
	if string(m.Data) != testMsg {
		t.Error("Did not receive expected UDP message")
	}
}

func Test__UDPServer__ShouldWriteData__OnLoopbackAddr(t *testing.T) {
	server1, err := udp.Open(testAddr1)
	if err != nil {
		t.Error(err)
	}
	defer server1.Close()

	server2, err := udp.Open(testAddr2)
	if err != nil {
		t.Error(err)
	}
	defer server2.Close()

	server1.Send([]byte(testMsg), server2.Addr)
	m := <-server2.Channel

	if string(m.Data) != testMsg {
		t.Error("Did not receive expected UDP message")
	}
}
