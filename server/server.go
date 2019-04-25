package server

import (
	"fmt"
	zx303DeviceAdministrator "gitlab.com/iotTracker/brain/tracker/device/zx303/administrator"
	"gitlab.com/iotTracker/nerve/log"
	"gitlab.com/iotTracker/nerve/server/client"
	serverException "gitlab.com/iotTracker/nerve/server/exception"
	serverMessage "gitlab.com/iotTracker/nerve/server/message"
	serverMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler"
	"net"
)

type server struct {
	Port                     string
	IPAddress                string
	done                     chan bool
	listener                 net.Listener
	MessageHandlers          map[serverMessage.Type]serverMessageHandler.Handler
	zx303DeviceAdministrator zx303DeviceAdministrator.Administrator
}

func New(
	Port string,
	IPAddress string,
	zx303DeviceAdministrator zx303DeviceAdministrator.Administrator,
) *server {

	return &server{
		Port:                     Port,
		IPAddress:                IPAddress,
		MessageHandlers:          make(map[serverMessage.Type]serverMessageHandler.Handler),
		done:                     make(chan bool),
		zx303DeviceAdministrator: zx303DeviceAdministrator,
	}
}

func (s *server) RegisterMessageHandler(messageType serverMessage.Type, handler serverMessageHandler.Handler) {
	s.MessageHandlers[messageType] = handler
}

func (s *server) Start() error {
	log.Info(fmt.Sprintf("Starting  Server listening at %s:%s", s.IPAddress, s.Port))
	var err error
	s.listener, err = net.Listen("tcp4", fmt.Sprintf("%s:%s", s.IPAddress, s.Port))
	if err != nil {
		return serverException.Listen{Reasons: []string{"", err.Error()}}
	}
	defer func() {
		log.Info("closing listener")
		s.listener.Close()
	}()

	for {
		c, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.done:
				return nil
			default:
				log.Error(serverException.AcceptConnection{Reasons: []string{"", err.Error()}}.Error())
			}
		}

		newClient := client.New(
			s.zx303DeviceAdministrator,
			c,
			s.MessageHandlers,
		)
		go newClient.HandleRX()
		go newClient.HandleTX()
		go newClient.HandleLifeCycle()
	}
	log.Info("Stopping Socket Server")
	return nil
}

func (s *server) Stop() {
	s.done <- true
	s.listener.Close()
}
