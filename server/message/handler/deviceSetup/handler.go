package deviceSetup

import (
	"fmt"
	hexPadding "gitlab.com/iotTracker/nerve/hex/padding"
	"gitlab.com/iotTracker/nerve/log"
	clientSession "gitlab.com/iotTracker/nerve/server/client/session"
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

	if len(reasonsInvalid) > 0 {
		return serverMessageHandlerException.MessageInvalid{Reasons: reasonsInvalid, Message: request.Message}
	}
	return nil
}

func (h *handler) Handle(clientSession *clientSession.Session, request *serverMessageHandler.HandleRequest) (*serverMessageHandler.HandleResponse, error) {
	if err := h.ValidateHandleRequest(request); err != nil {
		return nil, err
	}

	log.Info("Device Setup")

	uploadInterval := "0060"
	deviceSwitchBitString := fmt.Sprintf("%d%d%d%d%d%d%d%d",
		0, // n/a
		0, // n/a
		0, // sensor switch
		0, // light sense
		0, // bluetooth
		0, // vibration alarm
		0, // step
		1, // gps
	)
	deviceSwitchInt, err := strconv.ParseInt(deviceSwitchBitString, 2, 9)
	if err != nil {
		return nil, serverMessageHandlerException.Handling{Reasons: []string{"device bit string to int parse", err.Error()}}
	}
	deviceSwitchHexString := hexPadding.Pad(fmt.Sprintf("%x", deviceSwitchInt), 2)

	return &serverMessageHandler.HandleResponse{Messages: []serverMessage.Message{{
		Type:       request.Message.Type,
		DataLength: 31,
		Data: fmt.Sprintf("%s%s0000000000000000000000000000000000000000000000003B3B3B",
			uploadInterval,
			deviceSwitchHexString,
		),
	}},
	}, nil

	// 78 78 03 69 16 04 13 03 18 49 14 75 90 5B D3 0E 25 00 1E 10 BB F7 63 5D 14 75 90 06 E6 26 56 05 01 CC 00 28 66 0F 21 32 28 66 0F 1F 28 28 66 0E A8 1E 28 66 10 73 14 28 66 0F 20 14 0D 0A
	// 78 78 05 57 78 78 5c a4 0d0a
}
