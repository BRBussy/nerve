package handler

import (
	serverMessage "gitlab.com/iotTracker/nerve/server/message"
	serverSession "gitlab.com/iotTracker/nerve/server/session"
)

type Handler interface {
	Handle(serverSession *serverSession.Session, request *HandleRequest) (*HandleResponse, error)
}

type HandleRequest struct {
	Message serverMessage.Message
}

type HandleResponse struct {
	Messages []serverMessage.Message
}
