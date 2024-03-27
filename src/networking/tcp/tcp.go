package tcp

import (
	errorutils "drm-blockchain/src/utils/error"
	"net"
)

type Server struct {
	listener    *net.TCPListener
	Connections chan net.Conn
	closed      bool
}

func Open(addr string) (*Server, error) {
	resolved, err := net.ResolveTCPAddr("tcp", addr)

	if err != nil {
		return nil, errorutils.NewfWithInner(err, "Failed TCP address resolution for \"%s\"", addr)
	}

	listener, err := net.ListenTCP("tcp", resolved)

	if err != nil {
		return nil, errorutils.NewfWithInner(err, "Failed to listen on TCP for \"%s\"", addr)
	}

	server := new(Server)
	server.listener = listener
	server.Connections = make(chan net.Conn)
	server.closed = false

	go server.listen()

	return server, nil
}

func (server *Server) listen() error {
	if server.closed {
		panic("Server closed!")
	}

	for {
		conn, err := server.listener.Accept()

		if err != nil {
			return err
		}

		server.Connections <- conn
	}
}

func (server *Server) Close() {
	close(server.Connections)
	server.listener.Close()
	server.closed = true
}
