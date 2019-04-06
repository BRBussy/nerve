package heartbeat

import (
	"gitlab.com/iotTracker/nerve/log"
	serverMessage "gitlab.com/iotTracker/nerve/server/message"
	serverMessageException "gitlab.com/iotTracker/nerve/server/message/exception"
	serverMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler"
)

type handler struct {
}

func New() serverMessageHandler.Handler {
	return &handler{}
}

func (h *handler) ValidateHandleRequest(request *handler.HandleRequest *serverMessage.Message) error {
	reasonsInvalid := make([]string, 0)

	if err := message.IsValid(); err != nil {
		reasonsInvalid = append(reasonsInvalid, "invalid message: "+err.Error())
	}

	if len(reasonsInvalid) > 0 {
		return serverMessageException.Invalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (h *handler) Handle(message *serverMessage.Message) (*serverMessage.Message, error) {
	if err := h.ValidateHandleRequest(request *handler.HandleRequest); err != nil {
		return nil, err
	}

	log.Info("Heartbeat")

	return nil, nil
}
