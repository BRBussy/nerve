package login

import (
	zx303ServerMessage "gitlab.com/iotTracker/nerve/server/zx303/message"
	zx303ServerMessageException "gitlab.com/iotTracker/nerve/server/zx303/message/exception"
	zx303ServerMessageHandler "gitlab.com/iotTracker/nerve/server/zx303/message/handler"
)

type handler struct {
}

func New() zx303ServerMessageHandler.Handler {
	return &handler{}
}

func (h *handler) ValidateMessage(message *zx303ServerMessage.Message) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return zx303ServerMessageException.Invalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (h *handler) Handle(message *zx303ServerMessage.Message) (*zx303ServerMessage.Message, error) {
	if err := h.ValidateMessage(message); err != nil {
		return nil, err
	}

	return &zx303ServerMessage.Message{}, nil
}
