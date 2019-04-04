package proxy

import (
	ServerMessage "gitlab.com/iotTracker/nerve/server/message"
	ServerMessageHandler "gitlab.com/iotTracker/nerve/server/messageer"
	ServerMessageHandlerException "gitlab.com/iotTracker/nerve/server/messageer/exception"
	ServerMessageException "gitlab.com/iotTracker/nerve/server/messagetion"
)

type handler struct {
	handlers map[ServerMessage.Type]ServerMessageHandler.Handler
}

func New(
	handlers map[ServerMessage.Type]ServerMessageHandler.Handler,
) ServerMessageHandler.Handler {
	return &handler{
		handlers: handlers,
	}
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

	if h.handlers[message.Type] == nil {
		return nil, ServerMessageHandlerException.UnsupportedMessage{
			Reasons: []string{"no handler for message"},
			Message: *message,
		}
	}

	return h.handlers[message.Type].Handle(message)
}
