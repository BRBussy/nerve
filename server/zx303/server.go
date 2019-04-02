package zx303

import (
	"fmt"
	nerveServer "gitlab.com/iotTracker/nerve/server"
	serverException "gitlab.com/iotTracker/nerve/server/exception"
	"net"
)

type server struct {
	Port      string
	IPAddress string
}

func (s *server) Start(request *nerveServer.StartRequest) error {
	s.Port = request.Port
	s.IPAddress = request.IPAddress

	listener, err := net.Listen("tcp4", fmt.Sprintf("%s:%s", s.IPAddress, s.Port))
	if err != nil {
		return serverException.Listen{Reasons: []string{"zx303", err.Error()}}
	}
	defer listener.Close()

	for {
		c, err := listener.Accept()
		if err != nil {
			return serverException.AcceptConnection{Reasons: []string{"zx303", err.Error()}}
		}
		go s.handleConnection(c)
	}
}

func (s *server) handleConnection(c net.Conn) {

}
