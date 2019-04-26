package handler

import (
	serverClient "gitlab.com/iotTracker/nerve/server/client"
	serverMessage "gitlab.com/iotTracker/nerve/server/message"
)

type Handler interface {
	Handle(request *HandleRequest) (*HandleResponse, error)
}

type HandleRequest struct {
	Client  *serverClient.Client
	Message serverMessage.Message
}

type HandleResponse struct {
	Messages []serverMessage.Message
}
