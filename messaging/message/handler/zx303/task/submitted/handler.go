package submitted

import (
	"gitlab.com/iotTracker/brain/search/identifier/id"
	zx303TaskAdministrator "gitlab.com/iotTracker/brain/tracker/zx303/task/administrator"
	zx303TaskStep "gitlab.com/iotTracker/brain/tracker/zx303/task/step"
	messagingClient "gitlab.com/iotTracker/messaging/client"
	messagingException "gitlab.com/iotTracker/messaging/exception"
	messagingHub "gitlab.com/iotTracker/messaging/hub"
	messagingMessage "gitlab.com/iotTracker/messaging/message"
	messagingMessageHandler "gitlab.com/iotTracker/messaging/message/handler"
	messagingMessageHandlerException "gitlab.com/iotTracker/messaging/message/handler/exception"
	zx303TaskSubmittedMessage "gitlab.com/iotTracker/messaging/message/zx303/task/submitted"
	nerveException "gitlab.com/iotTracker/nerve/exception"
	zx303Client "gitlab.com/iotTracker/nerve/server/client"
)

type handler struct {
	MessagingHub      messagingHub.Hub
	taskAdministrator zx303TaskAdministrator.Administrator
}

func New(
	MessagingHub messagingHub.Hub,
	taskAdministrator zx303TaskAdministrator.Administrator,
) messagingMessageHandler.Handler {
	return &handler{
		MessagingHub:      MessagingHub,
		taskAdministrator: taskAdministrator,
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
		// TODO: this should not fail. the client is not in hub, may be registered to another nerve instance
		// TODO: remove this once there is something retrying tasks
		// fail the task at an indeterminate step
		if _, err := h.taskAdministrator.FailTask(&zx303TaskAdministrator.FailTaskRequest{
			ZX303TaskIdentifier: id.Identifier{
				Id: taskSubmittedMessage.Task.Id,
			},
			FailedStepIdx: -1,
		}); err != nil {
			return nerveException.Unexpected{Reasons: []string{
				"could not fail task",
				err.Error(),
				"could not get client",
			}}
		}
		return nerveException.Unexpected{Reasons: []string{
			"could not get client",
			err.Error(),
		}}
	}

	// cast to xz303 server client
	zx303ServerClient, ok := client.(*zx303Client.Client)
	if !ok {
		// fail the task at an indeterminate step
		if _, err := h.taskAdministrator.FailTask(&zx303TaskAdministrator.FailTaskRequest{
			ZX303TaskIdentifier: id.Identifier{
				Id: taskSubmittedMessage.Task.Id,
			},
			FailedStepIdx: -1,
		}); err != nil {
			return nerveException.Unexpected{Reasons: []string{"could not cast client to zx303Client.Client", "could not fail task"}}
		}
		return nerveException.Unexpected{Reasons: []string{"could not cast client to zx303Client.Client"}}
	}

	// get the pending step that should be handled by the client
	pendingStep, pendingStepIdx, err := taskSubmittedMessage.Task.PendingStep()
	if err != nil {
		// fail the task at an indeterminate step
		if _, err := h.taskAdministrator.FailTask(&zx303TaskAdministrator.FailTaskRequest{
			ZX303TaskIdentifier: id.Identifier{
				Id: taskSubmittedMessage.Task.Id,
			},
			FailedStepIdx: pendingStepIdx,
		}); err != nil {
			return nerveException.Unexpected{Reasons: []string{
				"could not fail task",
				err.Error(),
				"could not get tasks pending step",
			}}
		}
		return nerveException.Unexpected{Reasons: []string{"could not get tasks pending step", err.Error()}}
	}

	// give the pending step to the client to be handled
	if err := zx303ServerClient.HandleTaskStep(*pendingStep); err != nil {
		// fail task at this step
		if _, err := h.taskAdministrator.FailTask(&zx303TaskAdministrator.FailTaskRequest{
			ZX303TaskIdentifier: id.Identifier{
				Id: taskSubmittedMessage.Task.Id,
			},
			FailedStepIdx: pendingStepIdx,
		}); err != nil {
			return nerveException.Unexpected{Reasons: []string{
				"could not fail task",
				err.Error(),
				"handling step failed",
			}}
		}
		return messagingMessageHandlerException.Handling{Reasons: []string{"handling step", err.Error()}}
	}

	// transition step straight to finished
	if _, err := h.taskAdministrator.TransitionTask(&zx303TaskAdministrator.TransitionTaskRequest{
		ZX303TaskIdentifier: id.Identifier{
			Id: taskSubmittedMessage.Task.Id,
		},
		StepIdx:       pendingStepIdx,
		NewStepStatus: zx303TaskStep.Finished,
	}); err != nil {
		return nerveException.Unexpected{Reasons: []string{"transitioning step to finished", err.Error()}}
	}

	return nil
}
