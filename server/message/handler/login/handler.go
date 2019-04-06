package login

import (
	"gitlab.com/iotTracker/nerve/log"
	serverMessage "gitlab.com/iotTracker/nerve/server/message"
	serverMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler"
	serverMessageHandlerException "gitlab.com/iotTracker/nerve/server/message/handler/exception"
)

const SuccessData = ""
const FailureData = "44"

type handler struct {
}

func New() serverMessageHandler.Handler {
	return &handler{}
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

func (h *handler) Handle(request *serverMessageHandler.HandleRequest) (*serverMessageHandler.HandleResponse, error) {
	if err := h.ValidateHandleRequest(request); err != nil {
		return nil, err
	}

	log.Info("Log in Device with IMEI: ", request.Message.Data[:16])

	outMessage := serverMessage.Message{
		Type:       serverMessage.Login,
		Data:       FailureData,
		DataLength: 1,
	}
	// determine if device is allowed to log in
	if true {
		outMessage.Data = SuccessData
	}

	return &serverMessageHandler.HandleResponse{
		Messages: []serverMessage.Message{outMessage},
	}, nil
}
