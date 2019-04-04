package exception

import (
	"fmt"
	serverMessage "gitlab.com/iotTracker/nerve/server/message"
	"strings"
)

type Listen struct {
	Reasons []string
}

func (e Listen) Error() string {
	return "listening error: " + strings.Join(e.Reasons, "; ")
}

type AcceptConnection struct {
	Reasons []string
}

func (e AcceptConnection) Error() string {
	return "accept connection error: " + strings.Join(e.Reasons, "; ")
}

type StartOfMessageNotFound struct {
	Data    string
	Reasons []string
}

func (e StartOfMessageNotFound) Error() string {
	return fmt.Sprintf("start of message not found in given data:\n'%s'\n%s", e.Data, strings.Join(e.Reasons, "; "))
}

type DecodingError struct {
	Reasons []string
}

func (e DecodingError) Error() string {
	return "decoding error: " + strings.Join(e.Reasons, "; ")
}

type NoHandler struct {
	MessageType serverMessage.Type
}

func (e NoHandler) Error() string {
	return "no handler for given message type: " + e.MessageType.String()
}
