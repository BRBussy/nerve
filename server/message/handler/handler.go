package handler

import (
	serverSession "gitlab.com/iotTracker/nerve/server/client/session"
	serverMessage "gitlab.com/iotTracker/nerve/server/message"
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
