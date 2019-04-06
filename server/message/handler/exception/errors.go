package exception

import (
	"fmt"
	serverMessage "gitlab.com/iotTracker/nerve/server/message"
	"strings"
)

type UnsupportedMessage struct {
	Reasons []string
	Message serverMessage.Message
}

func (e UnsupportedMessage) Error() string {
	return fmt.Sprintf("unsupported message:\n%s\n%s", e.Message.String(), strings.Join(e.Reasons, "; "))
}

type Handling struct {
	Reasons []string
}

func (e Handling) Error() string {
	return "handling error: " + strings.Join(e.Reasons, "; ")
}

type MessageInvalid struct {
	Message serverMessage.Message
	Reasons []string
}

func (e MessageInvalid) Error() string {
	return fmt.Sprintf("message invalid: %s, %s", e.Message, strings.Join(e.Reasons, "; "))
}
