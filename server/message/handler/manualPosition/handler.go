package manualPosition

import (
	"gitlab.com/iotTracker/nerve/log"
	clientSession "gitlab.com/iotTracker/nerve/server/client/session"
	serverMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler"
	serverMessageHandlerException "gitlab.com/iotTracker/nerve/server/message/handler/exception"
)

type handler struct {
}

func New() serverMessageHandler.Handler {
	return &handler{}
}

func (h *handler) ValidateHandleRequest(request *serverMessageHandler.HandleRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return serverMessageHandlerException.MessageInvalid{Reasons: reasonsInvalid, Message: request.Message}
	}
	return nil
}

func (h *handler) Handle(clientSession *clientSession.Session, request *serverMessageHandler.HandleRequest) (*serverMessageHandler.HandleResponse, error) {
	if err := h.ValidateHandleRequest(request); err != nil {
		return nil, err
	}

	log.Info("Manual Position")

	return &serverMessageHandler.HandleResponse{}, nil
}
