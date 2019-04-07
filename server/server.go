package server

import (
	"fmt"
	"gitlab.com/iotTracker/nerve/log"
	"gitlab.com/iotTracker/nerve/server/client"
	serverException "gitlab.com/iotTracker/nerve/server/exception"
	serverMessage "gitlab.com/iotTracker/nerve/server/message"
	serverMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler"
	"net"
)

type server struct {
	Port            string
	IPAddress       string
	MessageHandlers map[serverMessage.Type]serverMessageHandler.Handler
}

func New(
	Port string,
	IPAddress string,
) *server {

	return &server{
		Port:            Port,
		IPAddress:       IPAddress,
		MessageHandlers: make(map[serverMessage.Type]serverMessageHandler.Handler),
	}
}

func (s *server) RegisterMessageHandler(messageType serverMessage.Type, handler serverMessageHandler.Handler) {
	s.MessageHandlers[messageType] = handler
}

func (s *server) Start() error {
	log.Info(fmt.Sprintf("Starting  Server listening at %s:%s", s.IPAddress, s.Port))
	listener, err := net.Listen("tcp4", fmt.Sprintf("%s:%s", s.IPAddress, s.Port))
	if err != nil {
		return serverException.Listen{Reasons: []string{"", err.Error()}}
	}
	defer listener.Close()

	for {
		c, err := listener.Accept()
		if err != nil {
			return serverException.AcceptConnection{Reasons: []string{"", err.Error()}}
		}

		newClient := client.New(c, s.MessageHandlers)
		go newClient.HandleRX()
		go newClient.HandleTX()
	}
}
