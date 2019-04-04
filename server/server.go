package server

import (
	"bufio"
	"fmt"
	"gitlab.com/iotTracker/nerve/log"
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
		go s.handleConnection(c)
	}
}

func (s *server) handleConnection(c net.Conn) {
	// TODO: use heart beat packets to determine when to drop the connection

	log.Info(fmt.Sprintf("serving %s", c.RemoteAddr().String()))
	reader := bufio.NewReaderSize(c, 1024)
	scr := bufio.NewScanner(reader)
	scr.Split(splitFunc)
CommLoop:
	for {
		// scan advances the scanner to the next token
		// which in this case is a complete message from the device
		// it returns false when the scan stops by reaching the end
		// of the input or an error
		for scr.Scan() {
			// create message from data token
			inMessage, err := serverMessage.New(string(scr.Bytes()))
			if err != nil {
				log.Warn(err.Error())
				continue
			}
			inMessageBytes, _ := inMessage.Bytes()
			log.Info("IN -->", inMessageBytes)

			// handle the message
			outMessage, err := s.handleMessage(inMessage)
			if err != nil {
				log.Warn(err.Error())
				continue
			}
			// if a message needs to be returned, return it
			if outMessage != nil {
				// send the message back
				outMessageBytes, err := outMessage.Bytes()
				if err != nil {
					log.Warn("error converting message to bytes:" + err.Error())
					break CommLoop
				}
				log.Info("OUT -->", outMessageBytes)
				if _, err = c.Write(outMessageBytes); err != nil {
					log.Warn("error sending message to device:" + err.Error())
					break CommLoop
				}
			}
		}
		// check to see if scanner stopped with an error
		if scr.Err() != nil {
			log.Warn("scanning stopped with an error:", scr.Err().Error())
			break CommLoop
		}
	}
	log.Info(fmt.Sprintf("%s disconnected", c.RemoteAddr().String()))
}

func (s *server) handleMessage(message *serverMessage.Message) (*serverMessage.Message, error) {
	if s.MessageHandlers[message.Type] == nil {
		return nil, serverException.NoHandler{MessageType: message.Type}
	}

	return s.MessageHandlers[message.Type].Handle(message)
}
