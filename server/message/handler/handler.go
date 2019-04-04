package handler

import (
	ServerMessage "gitlab.com/iotTracker/nerve/server/message"
)

type Handler interface {
	Handle(message *ServerMessage.Message) (*ServerMessage.Message, error)
}
