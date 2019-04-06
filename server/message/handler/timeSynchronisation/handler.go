package timeSynchronisation

import (
	"fmt"
	hexPadding "gitlab.com/iotTracker/nerve/hex/padding"
	"gitlab.com/iotTracker/nerve/log"
	serverMessage "gitlab.com/iotTracker/nerve/server/message"
	serverMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler"
	serverMessageHandlerException "gitlab.com/iotTracker/nerve/server/message/handler/exception"
	"time"
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

func (h *handler) Handle(request *serverMessageHandler.HandleRequest) (*serverMessageHandler.HandleResponse, error) {
	if err := h.ValidateHandleRequest(request); err != nil {
		return nil, err
	}

	log.Info("Time Synchronisation")

	timeNow := time.Now().UTC()
	return &serverMessageHandler.HandleResponse{Messages: []serverMessage.Message{{
		Type:       request.Message.Type,
		DataLength: 7,
		Data: fmt.Sprintf("%s%s%s%s%s%s",
			hexPadding.Pad(fmt.Sprintf("%x", int(timeNow.Year())), 4),
			hexPadding.Pad(fmt.Sprintf("%x", int(timeNow.Month())), 2),
			hexPadding.Pad(fmt.Sprintf("%x", int(timeNow.Day())), 2),
			hexPadding.Pad(fmt.Sprintf("%x", int(timeNow.Hour())), 2),
			hexPadding.Pad(fmt.Sprintf("%x", int(timeNow.Minute())), 2),
			hexPadding.Pad(fmt.Sprintf("%x", int(timeNow.Second())), 2),
		),
	}},
	}, nil
}
