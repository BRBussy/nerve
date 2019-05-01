package gpsPosition

import (
	"fmt"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	zx303GPSReading "gitlab.com/iotTracker/brain/tracker/zx303/reading/gps"
	zx303GPSReadingMessage "gitlab.com/iotTracker/messaging/message/zx303/reading/gps"
	messagingProducer "gitlab.com/iotTracker/messaging/producer"
	"gitlab.com/iotTracker/nerve/log"
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

	if len(request.Message.Data) < 36 {
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
	// parse timestamp
	year, err := strconv.ParseInt(request.Message.Data[:2], 16, 0)
	if err != nil {
		return nil, serverMessageHandlerException.Handling{Reasons: []string{"parsing year", err.Error()}}
	}
	month, err := strconv.ParseInt(request.Message.Data[2:4], 16, 0)
	if err != nil {
		return nil, serverMessageHandlerException.Handling{Reasons: []string{"parsing month", err.Error()}}
	}
	day, err := strconv.ParseInt(request.Message.Data[4:6], 16, 0)
	if err != nil {
		return nil, serverMessageHandlerException.Handling{Reasons: []string{"parsing day", err.Error()}}
	}
	hour, err := strconv.ParseInt(request.Message.Data[6:8], 16, 0)
	if err != nil {
		return nil, serverMessageHandlerException.Handling{Reasons: []string{"parsing hour", err.Error()}}
	}
	minute, err := strconv.ParseInt(request.Message.Data[8:10], 16, 0)
	if err != nil {
		return nil, serverMessageHandlerException.Handling{Reasons: []string{"parsing minute", err.Error()}}
	}
	second, err := strconv.ParseInt(request.Message.Data[10:12], 16, 0)
	if err != nil {
		return nil, serverMessageHandlerException.Handling{Reasons: []string{"parsing second", err.Error()}}
	}
	//gpsInformationLength := request.Message.Data[12]
	noSatellites, err := strconv.ParseInt(string(request.Message.Data[13]), 16, 0)
	if err != nil {
		return nil, serverMessageHandlerException.Handling{Reasons: []string{"parsing no satellites", err.Error()}}
	}
	gpsLatInt, err := strconv.ParseInt(request.Message.Data[14:22], 16, 0)
	if err != nil {
		return nil, serverMessageHandlerException.Handling{Reasons: []string{"parsing gps lat int", err.Error()}}
	}
	gpsLatitude := float32(gpsLatInt) / (30000 * 60)

	gpsLongInt, err := strconv.ParseInt(request.Message.Data[22:30], 16, 0)
	if err != nil {
		return nil, serverMessageHandlerException.Handling{Reasons: []string{"parsing gps lon int", err.Error()}}
	}
	gpsLongitude := float32(gpsLongInt) / (30000 * 60)

	speed, err := strconv.ParseInt(request.Message.Data[30:32], 16, 0)
	if err != nil {
		return nil, serverMessageHandlerException.Handling{Reasons: []string{"parsing speed", err.Error()}}
	}

	gpsFlagsInt, err := strconv.ParseInt(request.Message.Data[32:36], 16, 0)
	if err != nil {
		return nil, serverMessageHandlerException.Handling{Reasons: []string{"heading info", err.Error()}}
	}
	gpsFlagsBitString := fmt.Sprintf("%b", gpsFlagsInt)
	paddingBits := 16 - len(gpsFlagsBitString)
	for i := 0; i < paddingBits; i++ {
		gpsFlagsBitString = "0" + gpsFlagsBitString
	}
	//positioning, err := strconv.ParseBool(string(gpsFlagsBitString[3]))
	//if err != nil {
	//	return nil, serverMessageHandlerException.Handling{Reasons: []string{"positioning bool", err.Error()}}
	//}
	west, err := strconv.ParseBool(string(gpsFlagsBitString[4]))
	if err != nil {
		return nil, serverMessageHandlerException.Handling{Reasons: []string{"west bool", err.Error()}}
	}
	north, err := strconv.ParseBool(string(gpsFlagsBitString[5]))
	if err != nil {
		return nil, serverMessageHandlerException.Handling{Reasons: []string{"north bool", err.Error()}}
	}
	if !north {
		gpsLatitude = gpsLatitude * -1
	}
	if west {
		gpsLongitude = gpsLongitude * -1
	}

	heading, err := strconv.ParseInt(gpsFlagsBitString[6:], 2, 0)
	if err != nil {
		return nil, serverMessageHandlerException.Handling{Reasons: []string{"heading", err.Error()}}
	}

	log.Info(fmt.Sprintf("GPS Position: Timestamp: 20%d %d/%d %d:%d:%d,  No. Satellites %d, Coordinates %f, %f, Speed %d km/h, Heading %d° ",
		year,
		month,
		day,
		hour,
		minute,
		second,
		noSatellites,
		gpsLatitude,
		gpsLongitude,
		speed,
		heading,
	))

	if err := h.brainQueueProducer.Produce(zx303GPSReadingMessage.Message{
		Reading: zx303GPSReading.Reading{
			DeviceId: id.Identifier{
				Id: serverSession.ZX303Device.Id,
			},
			OwnerPartyType:    serverSession.ZX303Device.OwnerPartyType,
			OwnerId:           serverSession.ZX303Device.OwnerId,
			AssignedPartyType: serverSession.ZX303Device.AssignedPartyType,
			AssignedId:        serverSession.ZX303Device.AssignedId,
			NoSatellites:      noSatellites,
			TimeStamp:         time.Now().UTC().Unix(),
			Latitude:          gpsLatitude,
			Longitude:         gpsLongitude,
			Speed:             speed,
			Heading:           heading,
		},
	}); err != nil {
		return nil, serverMessageHandlerException.MessageProduction{Reasons: []string{err.Error()}}
	}

	return &serverMessageHandler.HandleResponse{
		Messages: []serverMessage.Message{{
			Type:       request.Message.Type,
			DataLength: 0,
			Data:       request.Message.Data[:12],
		}},
	}, nil
}
