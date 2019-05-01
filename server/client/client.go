package client

import (
	"bufio"
	"fmt"
	zx303TaskStep "gitlab.com/iotTracker/brain/tracker/zx303/task/step"
	messagingClient "gitlab.com/iotTracker/messaging/client"
	messagingHub "gitlab.com/iotTracker/messaging/hub"
	messagingMessage "gitlab.com/iotTracker/messaging/message"
	zx303TransmitMessage "gitlab.com/iotTracker/messaging/message/zx303/transmit"
	nerveException "gitlab.com/iotTracker/nerve/exception"
	"gitlab.com/iotTracker/nerve/log"
	clientException "gitlab.com/iotTracker/nerve/server/client/exception"
	clientSession "gitlab.com/iotTracker/nerve/server/client/session"
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
	HeartbeatWait = 180 * time.Second
	// Maximum message size allowed from peer.
	MaxMessageSize = 1024
)

type Client struct {
	messagingHub     messagingHub.Hub
	socket           net.Conn
	outgoingMessages chan serverMessage.Message
	messageHandlers  map[serverMessage.Type]serverMessageHandler.Handler
	clientSession    clientSession.Session
	heartbeat        chan bool
	stop             chan bool
	stopTX           chan bool
	stopRX           bool
}

func New(
	socket net.Conn,
	messageHandlers map[serverMessage.Type]serverMessageHandler.Handler,
	messagingHub messagingHub.Hub,
) *Client {
	return &Client{
		socket:           socket,
		outgoingMessages: make(chan serverMessage.Message),
		messageHandlers:  messageHandlers,
		heartbeat:        make(chan bool),
		stopTX:           make(chan bool),
		stopRX:           false,
		stop:             make(chan bool),
		messagingHub:     messagingHub,
	}
}

func (c *Client) Send(message messagingMessage.Message) error {
	if message.Type() != messagingMessage.ZX303Transmit {
		return nerveException.Unexpected{Reasons: []string{"invalid message type provided to zx303 client send"}}
	}

	nerveServerMessage, ok := message.(zx303TransmitMessage.Message)
	if !ok {
		return nerveException.Unexpected{Reasons: []string{"could not cast messagingMessage to zx303TransmitMessage.Message"}}
	}

	c.outgoingMessages <- nerveServerMessage.Message

	return nil
}

func (c *Client) Stop() error {
	c.stop <- true
	return nil
}

func (c *Client) IdentifiedBy(identifier messagingClient.Identifier) bool {
	return messagingClient.Identifier{
		Type: messagingClient.ZX303,
		Id:   c.clientSession.ZX303Device.Id,
	} == identifier
}

func (c *Client) Identifier() messagingClient.Identifier {
	return messagingClient.Identifier{
		Type: messagingClient.ZX303,
		Id:   c.clientSession.ZX303Device.Id,
	}
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
			c.messagingHub.DeRegisterClient(c)
			break LifeCycle

		case <-c.heartbeat:
			heartbeatCountdownTimer.Reset(HeartbeatWait)

		case <-c.stop:
			c.stopTX <- true
			c.stopRX = true
			c.socket.Close()
			c.messagingHub.DeRegisterClient(c)
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
			if !(c.clientSession.LoggedIn || inMessage.Type == serverMessage.Login) {
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
					&c.clientSession,
					&serverMessageHandler.HandleRequest{
						Message: *inMessage,
					})
				if err != nil {
					log.Warn(err.Error())
					c.stop <- true
					continue
				}
				// if the client has not been set to logged in stop the client connection
				if !c.clientSession.LoggedIn {
					log.Warn(nerveException.Unexpected{Reasons: []string{
						"client still not logged in",
					}}.Error())
					c.stop <- true
					continue
				}

				// register client with messaging hub
				if err := c.messagingHub.RegisterClient(c); err != nil {
					log.Warn(nerveException.Unexpected{Reasons: []string{
						"registering client with hub",
						err.Error(),
					}})
					c.stop <- true
					continue
				}

			case serverMessage.Heartbeat:
				// notify the lifecycle monitor of the heartbeat
				c.heartbeat <- true
				// handle the heartbeat message
				response, err = c.messageHandlers[inMessage.Type].Handle(
					&c.clientSession,
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
					&c.clientSession,
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

func (c *Client) HandleTaskStep(step zx303TaskStep.Step) error {
	switch step.Type {
	case zx303TaskStep.SendResetCommand:
		log.Info("send reset command!!!")
	default:
		return clientException.HandlingTaskStep{Reasons: []string{"invalid step type", string(step.Type)}}
	}
	return nil
}
