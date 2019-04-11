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
	// Time allowed between heartbeats
	// if no heartbeat received in after this time the connection
	// is terminated
	HeartbeatWait = 30 * time.Second
	// Maximum message size allowed from peer.
	MaxMessageSize = 1024
)

type client struct {
	socket           net.Conn
	outgoingMessages chan serverMessage.Message
	messageHandlers  map[serverMessage.Type]serverMessageHandler.Handler
	loggedIn         bool
	stop             chan bool
	stopTX           chan bool
	stopRX           bool
}

func New(
	socket net.Conn,
	messageHandlers map[serverMessage.Type]serverMessageHandler.Handler,
) *client {
	return &client{
		socket:           socket,
		outgoingMessages: make(chan serverMessage.Message),
		messageHandlers:  messageHandlers,
		stopTX:           make(chan bool),
		stopRX:           false,
		stop:             make(chan bool),
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

func (c *client) HandleLifeCycle() {
	heartbeatCountdownTimer := time.NewTimer(HeartbeatWait)
LifeCycle:
	for {
		select {
		case <-heartbeatCountdownTimer.C:
			log.Info(fmt.Sprintf("timeout waiting for heartbeat from %s", c.socket.RemoteAddr().String()))
			c.stopTX <- true
			c.stopRX = true
			c.socket.Close()
			break LifeCycle

		case <-c.stop:
			c.stopTX <- true
			c.stopRX = true
			c.socket.Close()
			break LifeCycle
		}
	}
	log.Info(fmt.Sprintf("%s lifecycle ended", c.socket.RemoteAddr().String()))
}

func (c *client) HandleTX() {
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
				continue
			}
			c.socket.SetWriteDeadline(time.Now().Add(WriteWait))
			if _, err = c.socket.Write(outMessageBytes); err != nil {
				log.Warn(clientException.SendingMessage{Message: outMessage, Reasons: []string{err.Error()}}.Error())
				c.stop <- true
				continue
			}
			log.Info("OUT: ", outMessage.String())

		case <-c.stopTX:
			log.Info(fmt.Sprintf("stopping TX with %s", c.socket.RemoteAddr().String()))
			return
		}
	}
}

func (c *client) HandleRX() {

	log.Info(fmt.Sprintf("serving %s", c.socket.RemoteAddr().String()))
	reader := bufio.NewReaderSize(c.socket, MaxMessageSize)
	scr := bufio.NewScanner(reader)
	scr.Split(splitFunc)

Comms:
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
				continue
			}
			// send back any messages if required
			for msgIdx := range response.Messages {
				c.outgoingMessages <- response.Messages[msgIdx]
			}
		}
		// check to see if scanner stopped with an error
		if scr.Err() != nil {
			if !c.stopRX {
				log.Warn("scanning stopped with error:", scr.Err().Error())
				c.stop <- true
			}
			break Comms
		}
	}
	log.Info(fmt.Sprintf("connection with %s terminated", c.socket.RemoteAddr().String()))
}
