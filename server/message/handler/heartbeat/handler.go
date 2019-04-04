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

func (h *handler) ValidateMessage(message *serverMessage.Message) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return serverMessageException.Invalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (h *handler) Handle(message *serverMessage.Message) (*serverMessage.Message, error) {
	if err := h.ValidateMessage(message); err != nil {
		return nil, err
	}

	log.Info("Handling Heartbeat")

	return &serverMessage.Message{}, nil
}
