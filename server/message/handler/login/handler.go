package login

import (
	zx303DeviceAuthenticator "gitlab.com/iotTracker/brain/tracker/device/zx303/authenticator"
	"gitlab.com/iotTracker/nerve/log"
	serverMessage "gitlab.com/iotTracker/nerve/server/message"
	serverMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler"
	serverMessageHandlerException "gitlab.com/iotTracker/nerve/server/message/handler/exception"
	serverSession "gitlab.com/iotTracker/nerve/server/session"
)

const SuccessData = "01"
const FailureData = "44"

type handler struct {
	zx303DeviceAuthenticator zx303DeviceAuthenticator.Authenticator
}

func New(
	zx303DeviceAuthenticator zx303DeviceAuthenticator.Authenticator,
) serverMessageHandler.Handler {
	return &handler{
		zx303DeviceAuthenticator: zx303DeviceAuthenticator,
	}
}

func (h *handler) ValidateHandleRequest(request *serverMessageHandler.HandleRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(request.Message.Data) < 16 {
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

	log.Info("Log in Device with IMEI: ", request.Message.Data[:16])

	responseMessages := make([]serverMessage.Message, 0)

	if true {
		responseMessages = append(responseMessages, serverMessage.Message{
			Type:       serverMessage.Login,
			Data:       SuccessData,
			DataLength: 1,
		})
		responseMessages = append(responseMessages, serverMessage.Message{
			Type:       serverMessage.ManualPosition,
			Data:       "",
			DataLength: 1,
		})
	} else {
		responseMessages = append(responseMessages, serverMessage.Message{
			Type:       serverMessage.Login,
			Data:       FailureData,
			DataLength: 1,
		})
	}

	// 7878 05 57 78 78 5c a4 0d0a
	// 7878 1F 57 00 60 01 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 3B 3B 3B 0D 0A
	return &serverMessageHandler.HandleResponse{
		Messages: responseMessages,
	}, nil
}
