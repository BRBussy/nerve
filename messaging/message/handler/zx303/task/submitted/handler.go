package submitted

import (
	messagingClient "gitlab.com/iotTracker/messaging/client"
	messagingException "gitlab.com/iotTracker/messaging/exception"
	messagingHub "gitlab.com/iotTracker/messaging/hub"
	messagingMessage "gitlab.com/iotTracker/messaging/message"
	messagingMessageHandler "gitlab.com/iotTracker/messaging/message/handler"
	messageHandlerException "gitlab.com/iotTracker/messaging/message/handler/exception"
	zx303TaskSubmittedMessage "gitlab.com/iotTracker/messaging/message/zx303/task/submitted"
	nerveException "gitlab.com/iotTracker/nerve/exception"
	zx303Client "gitlab.com/iotTracker/nerve/server/client"
)

type handler struct {
	MessagingHub messagingHub.Hub
}

func New(
	MessagingHub messagingHub.Hub,
) messagingMessageHandler.Handler {
	return &handler{
		MessagingHub: MessagingHub,
	}
}

func (h *handler) WantsMessage(message messagingMessage.Message) bool {
	return message.Type() == messagingMessage.ZX303TaskSubmitted
}

func (*handler) ValidateMessage(message messagingMessage.Message) error {
	reasonsInvalid := make([]string, 0)

	if message == nil {
		reasonsInvalid = append(reasonsInvalid, "message is nil")
	} else {
		if _, ok := message.(zx303TaskSubmittedMessage.Message); !ok {
			reasonsInvalid = append(reasonsInvalid, "cannot cast message to zx303GPSReadingMessage.Message")
		}
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

	// get client from messaging hub
	client, err := h.MessagingHub.GetClient(messagingClient.Identifier{
		Type: messagingClient.ZX303,
		Id:   taskSubmittedMessage.Task.DeviceId.Id,
	})
	if err != nil {
		return messageHandlerException.Handling{Reasons: []string{"getting client", err.Error()}}
	}

	// cast to xz303 server client
	zx303ServerClient, ok := client.(*zx303Client.Client)
	if !ok {
		return nerveException.Unexpected{Reasons: []string{"could not cast client to zx303Client.Client"}}
	}

	pendingStep, err := taskSubmittedMessage.Task.PendingStep()
	if err != nil {
		return nerveException.Unexpected{Reasons: []string{"could not get tasks pending step", err.Error()}}
	}

	if err := zx303ServerClient.HandleTaskStep(*pendingStep); err != nil {
		return messageHandlerException.Handling{Reasons: []string{"handling task step", err.Error()}}
	}

	return nil
}
