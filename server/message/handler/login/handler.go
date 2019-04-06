package login

import (
	"gitlab.com/iotTracker/nerve/log"
	serverMessage "gitlab.com/iotTracker/nerve/server/message"
	serverMessageException "gitlab.com/iotTracker/nerve/server/message/exception"
	serverMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler"
)

const SuccessData = "01"
const FailureData = "44"

type handler struct {
}

func New() serverMessageHandler.Handler {
	return &handler{}
}

func (h *handler) ValidateHandleRequest(request *serverMessageHandler.HandleRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return serverMessageException.Invalid{Reasons: reasonsInvalid}
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
		DataLength: "01",
	}
	// determine if device is allowed to log in
	if true {
		outMessage.Data = SuccessData
	}

	return &serverMessageHandler.HandleResponse{
		Messages: []serverMessage.Message{outMessage},
	}, nil
}
