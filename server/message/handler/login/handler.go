package login

import (
	"gitlab.com/iotTracker/nerve/log"
	ServerMessage "gitlab.com/iotTracker/nerve/server/message"
	ServerMessageHandler "gitlab.com/iotTracker/nerve/server/messageer"
	ServerMessageException "gitlab.com/iotTracker/nerve/server/messagetion"
)

type handler struct {
}

func New() ServerMessageHandler.Handler {
	return &handler{}
}

func (h *handler) ValidateMessage(message *ServerMessage.Message) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return ServerMessageException.Invalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (h *handler) Handle(message *ServerMessage.Message) (*ServerMessage.Message, error) {
	if err := h.ValidateMessage(message); err != nil {
		return nil, err
	}

	log.Info("Handling Login")

	return &ServerMessage.Message{}, nil
}
