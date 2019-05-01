package handler

import (
	clientSession "gitlab.com/iotTracker/nerve/server/client/session"
	serverMessage "gitlab.com/iotTracker/nerve/server/message"
)

type Handler interface {
	Handle(clientSession *clientSession.Session, request *HandleRequest) (*HandleResponse, error)
}

type HandleRequest struct {
	Message serverMessage.Message
}

type HandleResponse struct {
	Messages []serverMessage.Message
}
