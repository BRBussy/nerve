package submitted

import (
	messagingException "gitlab.com/iotTracker/messaging/exception"
	messagingMessage "gitlab.com/iotTracker/messaging/message"
	messagingMessageHandler "gitlab.com/iotTracker/messaging/message/handler"
	zx303TaskSubmittedMessage "gitlab.com/iotTracker/messaging/message/zx303/task/submitted"
	nerveException "gitlab.com/iotTracker/nerve/exception"
	"gitlab.com/iotTracker/nerve/log"
)

type handler struct {
}

func New() messagingMessageHandler.Handler {
	return &handler{}
}

func (h *handler) WantsMessage(message messagingMessage.Message) bool {
	return message.Type() == messagingMessage.ZX303TaskSubmitted
}

func (*handler) ValidateMessage(message messagingMessage.Message) error {
	reasonsInvalid := make([]string, 0)

	if _, ok := message.(zx303TaskSubmittedMessage.Message); !ok {
		reasonsInvalid = append(reasonsInvalid, "cannot cast message to zx303GPSReadingMessage.Message")
	}

	if len(reasonsInvalid) > 0 {
		return messagingException.InvalidMessage{Reasons: reasonsInvalid}
	}

	return nil
}

func (h *handler) HandleMessage(message messagingMessage.Message) error {
	if err := h.ValidateMessage(message); err != nil {
		return err
	}
	taskSubmittedMessage, ok := message.(zx303TaskSubmittedMessage.Message)
	if !ok {
		return nerveException.Unexpected{Reasons: []string{"cannot cast message to zx303TaskSubmittedMessage.Message"}}
	}

	log.Info("handle task submitted message!", taskSubmittedMessage)

	return nil
}
