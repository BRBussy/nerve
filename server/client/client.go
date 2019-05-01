package client

import (
	"bufio"
	"fmt"
	nerveException "gitlab.com/iotTracker/nerve/exception"
	"gitlab.com/iotTracker/nerve/log"
	clientException "gitlab.com/iotTracker/nerve/server/client/exception"
	serverMessage "gitlab.com/iotTracker/nerve/server/message"
	serverMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler"
	serverSession "gitlab.com/iotTracker/nerve/server/session"
	"net"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	WriteWait = 10 * time.Second
	// Time allowed between heartbeats
	// if no heartbeat received in after this time the connection
	// is terminated
	HeartbeatWait = 180 * time.Second
	// Maximum message size allowed from peer.
	MaxMessageSize = 1024
)

type Client struct {
	socket           net.Conn
	outgoingMessages chan serverMessage.Message
	messageHandlers  map[serverMessage.Type]serverMessageHandler.Handler
	serverSession    serverSession.Session
	heartbeat        chan bool
	stop             chan bool
	stopTX           chan bool
	stopRX           bool
}

func New(
	socket net.Conn,
	messageHandlers map[serverMessage.Type]serverMessageHandler.Handler,
) *Client {
	return &Client{
		socket:           socket,
		outgoingMessages: make(chan serverMessage.Message),
		messageHandlers:  messageHandlers,
		heartbeat:        make(chan bool),
		stopTX:           make(chan bool),
		stopRX:           false,
		stop:             make(chan bool),
	}
}

func (c *Client) Send(message serverMessage.Message) error {
	messageBytes, err := message.Bytes()
	if err != nil {
		return clientException.MessageConversion{Reasons: []string{"message to bytes", err.Error()}}
	}
	if _, err = c.socket.Write(messageBytes); err != nil {
		return clientException.SendingMessage{Message: message, Reasons: []string{err.Error()}}
	}
	return nil
}

func (c *Client) HandleLifeCycle() {
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

		case <-c.heartbeat:
			heartbeatCountdownTimer.Reset(HeartbeatWait)

		case <-c.stop:
			c.stopTX <- true
			c.stopRX = true
			c.socket.Close()
			break LifeCycle
		}
	}
	log.Info(fmt.Sprintf("%s lifecycle ended", c.socket.RemoteAddr().String()))
}

func (c *Client) HandleTX() {
	for {
		select {
		case outMessage, ok := <-c.outgoingMessages:
			if !ok {
				log.Warn("the outgoing messages channel has been closed")
				c.stop <- true
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

func (c *Client) HandleRX() {

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

			// if the client is not logged in and this message is not of type Login terminate the connection
			if !(c.serverSession.LoggedIn || inMessage.Type == serverMessage.Login) {
				log.Warn(clientException.UnauthenticatedCommunication{Reasons: []string{"device not logged in"}}.Error())
				c.stop <- true
				continue
			}

			var response *serverMessageHandler.HandleResponse

			// if this is a log in message then we do not need to check that the client
			// is logged in before handling the message
			switch inMessage.Type {
			case serverMessage.Login:
				// handle the login message
				response, err = c.messageHandlers[inMessage.Type].Handle(
					&c.serverSession,
					&serverMessageHandler.HandleRequest{
						Message: *inMessage,
					})
				if err != nil {
					log.Warn(err.Error())
					c.stop <- true
					continue
				}
				// if the client has not been set to logged in stop the client connection
				if !c.serverSession.LoggedIn {
					log.Warn(nerveException.Unexpected{Reasons: []string{
						"client still not logged in",
					}}.Error())
					c.stop <- true
					continue
				}

			case serverMessage.Heartbeat:
				// notify the lifecycle monitor of the heartbeat
				c.heartbeat <- true
				// handle the heartbeat message
				response, err = c.messageHandlers[inMessage.Type].Handle(
					&c.serverSession,
					&serverMessageHandler.HandleRequest{
						Message: *inMessage,
					})
				if err != nil {
					log.Warn(err.Error())
					c.stop <- true
					continue
				}

			default:
				// it is not a Login Message and the device is logged in
				if c.messageHandlers[inMessage.Type] == nil {
					// if there is no handler log a warning and carry on
					log.Warn(clientException.NoHandler{Message: *inMessage}.Error())
					continue
				}

				// otherwise handle the message
				response, err = c.messageHandlers[inMessage.Type].Handle(
					&c.serverSession,
					&serverMessageHandler.HandleRequest{
						Message: *inMessage,
					})
				if err != nil {
					log.Warn(err.Error())
					continue
				}
			}

			// if the response is nil do not attempt to send back response messages
			if response == nil {
				log.Warn(nerveException.Unexpected{Reasons: []string{"handler response is nil"}}.Error())
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
