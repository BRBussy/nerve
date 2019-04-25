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

type UnauthenticatedCommunication struct {
	Reasons []string
}

func (e UnauthenticatedCommunication) Error() string {
	return "unauthenticated communication: " + strings.Join(e.Reasons, "; ")
}

type AuthenticationError struct {
	Reasons []string
}

func (e AuthenticationError) Error() string {
	return "authentication error: " + strings.Join(e.Reasons, "; ")
}
