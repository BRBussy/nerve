package login

import (
	"gitlab.com/iotTracker/nerve/log"
	serverMessage "gitlab.com/iotTracker/nerve/server/message"
	serverMessageException "gitlab.com/iotTracker/nerve/server/message/exception"
	serverMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler"
)

const LoginSuccessData = "01"
const LoginFailureData = "44"

type handler struct {
}

func New() serverMessageHandler.Handler {
	return &handler{}
}

func (h *handler) ValidateMessage(message *serverMessage.Message) error {
	reasonsInvalid := make([]string, 0)

	if len(message.Data) < 16 {
		reasonsInvalid = append(reasonsInvalid, "login message data not long enough")
	}

	if len(reasonsInvalid) > 0 {
		return serverMessageException.Invalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (h *handler) Handle(message *serverMessage.Message) (*serverMessage.Message, error) {
	if err := h.ValidateMessage(message); err != nil {
		return nil, err
	}

	log.Info("Log in Device with IMEI: ", message.Data[:16])

	outMessage := serverMessage.Message{
		Type:       serverMessage.Login,
		Data:       LoginFailureData,
		DataLength: "01",
	}
	// determine if device is allowed to log in
	if true {
		outMessage.Data = LoginSuccessData
	}

	return &outMessage, nil
}
