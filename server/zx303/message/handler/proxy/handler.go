package proxy

import (
	zx303ServerMessage "gitlab.com/iotTracker/nerve/server/zx303/message"
	zx303ServerMessageException "gitlab.com/iotTracker/nerve/server/zx303/message/exception"
	zx303ServerMessageHandler "gitlab.com/iotTracker/nerve/server/zx303/message/handler"
	zx303ServerMessageHandlerException "gitlab.com/iotTracker/nerve/server/zx303/message/handler/exception"
)

type handler struct {
	handlers map[zx303ServerMessage.Type]zx303ServerMessageHandler.Handler
}

func New(
	handlers map[zx303ServerMessage.Type]zx303ServerMessageHandler.Handler,
) zx303ServerMessageHandler.Handler {
	return &handler{
		handlers: handlers,
	}
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

	if h.handlers[message.Type] == nil {
		return nil, zx303ServerMessageHandlerException.UnsupportedMessage{
			Reasons: []string{"no handler for message"},
			Message: *message,
		}
	}

	return h.handlers[message.Type].Handle(message)
}
