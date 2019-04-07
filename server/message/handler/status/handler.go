package status

import (
	"fmt"
	hexPadding "gitlab.com/iotTracker/nerve/hex/padding"
	"gitlab.com/iotTracker/nerve/log"
	serverMessage "gitlab.com/iotTracker/nerve/server/message"
	serverMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler"
	serverMessageHandlerException "gitlab.com/iotTracker/nerve/server/message/handler/exception"
	"strconv"
)

type handler struct {
}

func New() serverMessageHandler.Handler {
	return &handler{}
}

func (h *handler) ValidateHandleRequest(request *serverMessageHandler.HandleRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(request.Message.Data) < 9 {
		reasonsInvalid = append(reasonsInvalid, "data not long enough")
	}

	if len(reasonsInvalid) > 0 {
		return serverMessageHandlerException.MessageInvalid{Reasons: reasonsInvalid, Message: request.Message}
	}
	return nil
}

func (h *handler) Handle(request *serverMessageHandler.HandleRequest) (*serverMessageHandler.HandleResponse, error) {
	if err := h.ValidateHandleRequest(request); err != nil {
		return nil, err
	}

	batteryPercentage, err := strconv.ParseInt(request.Message.Data[0:2], 16, 0)
	if err != nil {
		return nil, serverMessageHandlerException.Handling{Reasons: []string{"parsing battery percentage", err.Error()}}
	}
	softwareVersion, err := strconv.ParseInt(request.Message.Data[2:4], 16, 0)
	if err != nil {
		return nil, serverMessageHandlerException.Handling{Reasons: []string{"parsing software version", err.Error()}}
	}
	timezone, err := strconv.ParseInt(request.Message.Data[4:6], 16, 0)
	if err != nil {
		return nil, serverMessageHandlerException.Handling{Reasons: []string{"parsing time zone", err.Error()}}
	}
	uploadInterval, err := strconv.ParseInt(request.Message.Data[6:8], 16, 0)
	if err != nil {
		return nil, serverMessageHandlerException.Handling{Reasons: []string{"parsing upload interval", err.Error()}}
	}
	otherThing, err := strconv.ParseInt(request.Message.Data[8:], 16, 0)
	if err != nil {
		return nil, serverMessageHandlerException.Handling{Reasons: []string{"parsing otherThing", err.Error()}}
	}

	log.Info(fmt.Sprintf("Status: Battery Percentage: %d%%, Software V%d, Timezone: %d, Upload Interval: %d, Other: %d",
		batteryPercentage,
		softwareVersion,
		timezone,
		uploadInterval,
		otherThing,
	))

	return &serverMessageHandler.HandleResponse{Messages: []serverMessage.Message{{
		Type:       request.Message.Type,
		DataLength: 2,
		Data:       hexPadding.Pad(fmt.Sprintf("%x", 60), 2),
	}}}, nil
}
