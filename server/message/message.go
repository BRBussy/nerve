package message

import (
	"encoding/hex"
	"fmt"
	messageException "gitlab.com/iotTracker/nerve/server/message/exception"
)

const StartMarker = "7878"
const EndMarker = "0d0a"

type Message struct {
	Type       Type
	DataLength string
	Data       string
}

func New(rawMessage string) (*Message, error) {
	var newMessage Message
	if len(rawMessage) < 4 {
		return nil, messageException.Creation{Reasons: []string{"raw message string not long enough", rawMessage}}
	}

	newMessage.DataLength = rawMessage[:2]

	if len(rawMessage) == 4 {
		newMessage.Type = Type(rawMessage[2:])
		newMessage.Data = ""
	} else {
		newMessage.Type = Type(rawMessage[2:4])
		newMessage.Data = rawMessage[4:]
	}

	return &newMessage, nil
}

func (m Message) Bytes() ([]byte, error) {
	return hex.DecodeString(fmt.Sprintf(
		"%s%s%s%s%s",
		StartMarker,
		m.DataLength,
		m.Type,
		m.Data,
		EndMarker,
	))
}

func (m Message) String() string {
	switch m.Type {
	case Login:
		return fmt.Sprintf("[type: Login, Data: %s]", m.Data)
	case Heartbeat:
		return fmt.Sprintf("[type: Login, Data: %s]", m.Data)
	case GPSPosition:
		return fmt.Sprintf("[type: GPSPositioning, Data: %s]", m.Data)
	default:

	}
	return ""
}
