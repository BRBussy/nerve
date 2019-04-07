package client

import (
	"bufio"
	"fmt"
	"gitlab.com/iotTracker/nerve/log"
	clientException "gitlab.com/iotTracker/nerve/server/client/exception"
	serverMessage "gitlab.com/iotTracker/nerve/server/message"
	serverMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler"
	"net"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	WriteWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	HeartbeatWait = 30 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	HeartbeatPeriod = (HeartbeatWait * 9) / 10
	// Maximum message size allowed from peer.
	MaxMessageSize = 1024
)

type client struct {
	socket           net.Conn
	outgoingMessages chan serverMessage.Message
	messageHandlers  map[serverMessage.Type]serverMessageHandler.Handler
}

func New(
	socket net.Conn,
	messageHandlers map[serverMessage.Type]serverMessageHandler.Handler,
) *client {
	return &client{
		socket:           socket,
		outgoingMessages: make(chan serverMessage.Message),
		messageHandlers:  messageHandlers,
	}
}

func (c *client) Send(message serverMessage.Message) error {
	messageBytes, err := message.Bytes()
	if err != nil {
		return clientException.MessageConversion{Reasons: []string{"message to bytes", err.Error()}}
	}
	if _, err = c.socket.Write(messageBytes); err != nil {
		return clientException.SendingMessage{Message: message, Reasons: []string{err.Error()}}
	}
	return nil
}

func (c *client) HandleTX() {
	//heartbeatTicker := time.NewTicker(HeartbeatPeriod)

	defer func() {
		c.socket.Close()
		// heartbeatTicker.Stop()
	}()

	for {
		select {
		case outMessage, ok := <-c.outgoingMessages:
			if !ok {
				log.Info("the outgoing messages channel has been closed")
				return
			}
			outMessageBytes, err := outMessage.Bytes()
			if err != nil {
				log.Warn(clientException.MessageConversion{Reasons: []string{"message to bytes", err.Error()}}.Error())
			}
			c.socket.SetWriteDeadline(time.Now().Add(WriteWait))
			if _, err = c.socket.Write(outMessageBytes); err != nil {
				log.Warn(clientException.SendingMessage{Message: outMessage, Reasons: []string{err.Error()}}.Error())
			}
			log.Info("OUT: ", outMessage.String())

			//case <-heartbeatTicker.C:
			//	heartbeatMessage := serverMessage.Message{
			//		Type:       serverMessage.Heartbeat,
			//		DataLength: 1,
			//	}
			//	heartbeatMessageBytes, err := heartbeatMessage.Bytes()
			//	if err != nil {
			//		log.Warn(clientException.MessageConversion{Reasons: []string{"heartbeat message to bytes", err.Error()}}.Error())
			//	}
			//	c.socket.SetWriteDeadline(time.Now().Add(WriteWait))
			//	if _, err = c.socket.Write(heartbeatMessageBytes); err != nil {
			//		log.Warn(clientException.SendingMessage{Message: heartbeatMessage, Reasons: []string{err.Error()}}.Error())
			//	}
			//	log.Info("OUT: ", heartbeatMessage.String())

		}
	}
}

func (c *client) HandleRX() {
	defer func() {
		c.socket.Close()
	}()

	log.Info(fmt.Sprintf("serving %s", c.socket.RemoteAddr().String()))
	reader := bufio.NewReaderSize(c.socket, MaxMessageSize)
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
			log.Info("IN: ", inMessage.String())
			// handle the message
			if c.messageHandlers[inMessage.Type] == nil {
				log.Warn(clientException.NoHandler{Message: *inMessage}.Error())
			}
			response, err := c.messageHandlers[inMessage.Type].Handle(&serverMessageHandler.HandleRequest{
				Message: *inMessage,
			})
			if err != nil {
				log.Warn(err.Error())
			}
			// send back any messages if required
			for msgIdx := range response.Messages {
				c.outgoingMessages <- response.Messages[msgIdx]
			}
		}
		// check to see if scanner stopped with an error
		if scr.Err() != nil {
			log.Warn("scanning stopped with an error:", scr.Err().Error())
			break CommLoop
		}
	}

}
