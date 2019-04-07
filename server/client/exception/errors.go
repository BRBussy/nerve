package exception

import (
	"fmt"
	serverMessage "gitlab.com/iotTracker/nerve/server/message"
	"strings"
)

type MessageConversion struct {
	Reasons []string
}

func (e MessageConversion) Error() string {
	return "error converting message: " + strings.Join(e.Reasons, "; ")
}

type SendingMessage struct {
	Reasons []string
	Message serverMessage.Message
}

func (e SendingMessage) Error() string {
	return fmt.Sprintf("error sending message: %s : %s", e.Message, strings.Join(e.Reasons, "; "))
}

type NoHandler struct {
	Message serverMessage.Message
}

func (e NoHandler) Error() string {
	return fmt.Sprintf("no handler for given message %s", e.Message)
}
