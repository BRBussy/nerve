package handler

import (
	ServerMessage "gitlab.com/iotTracker/nerve/server/message"
)

type Handler interface {
	Handle(request *HandleRequest) (*HandleResponse, error)
}

type HandleRequest struct {
	Message  ServerMessage.Message
	LoggedIn bool
}

type HandleResponse struct {
	Message ServerMessage.Message
}
