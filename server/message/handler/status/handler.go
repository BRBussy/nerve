package status

import (
	"fmt"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	zx303StatusReading "gitlab.com/iotTracker/brain/tracker/zx303/reading/status"
	zx303StatusReadingMessage "gitlab.com/iotTracker/messaging/message/zx303/reading/status"
	messagingProducer "gitlab.com/iotTracker/messaging/producer"
	hexPadding "gitlab.com/iotTracker/nerve/hex/padding"
	serverSession "gitlab.com/iotTracker/nerve/server/client/session"
	serverMessage "gitlab.com/iotTracker/nerve/server/message"
	serverMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler"
	serverMessageHandlerException "gitlab.com/iotTracker/nerve/server/message/handler/exception"
	"strconv"
	"time"
)

type handler struct {
	brainQueueProducer messagingProducer.Producer
}

func New(
	brainQueueProducer messagingProducer.Producer,
) serverMessageHandler.Handler {
	return &handler{
		brainQueueProducer: brainQueueProducer,
	}
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

func (h *handler) Handle(serverSession *serverSession.Session, request *serverMessageHandler.HandleRequest) (*serverMessageHandler.HandleResponse, error) {
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
	//otherThing, err := strconv.ParseInt(request.Message.Data[8:], 16, 0)
	//if err != nil {
	//	return nil, serverMessageHandlerException.Handling{Reasons: []string{"parsing otherThing", err.Error()}}
	//}

	if err := h.brainQueueProducer.Produce(zx303StatusReadingMessage.Message{
		Reading: zx303StatusReading.Reading{
			DeviceId: id.Identifier{
				Id: serverSession.ZX303Device.Id,
			},
			OwnerPartyType:    serverSession.ZX303Device.OwnerPartyType,
			OwnerId:           serverSession.ZX303Device.OwnerId,
			AssignedPartyType: serverSession.ZX303Device.AssignedPartyType,
			AssignedId:        serverSession.ZX303Device.AssignedId,
			Timestamp:         time.Now().UTC().Unix(),
			BatteryPercentage: batteryPercentage,
			UploadInterval:    uploadInterval,
			SoftwareVersion:   softwareVersion,
			Timezone:          timezone,
		},
	}); err != nil {
		return nil, serverMessageHandlerException.MessageProduction{Reasons: []string{err.Error()}}
	}

	return &serverMessageHandler.HandleResponse{Messages: []serverMessage.Message{{
		Type:       request.Message.Type,
		DataLength: 2,
		Data:       hexPadding.Pad(fmt.Sprintf("%x", 60), 2),
	}}}, nil
}
