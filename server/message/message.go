package message

import (
	"encoding/hex"
	"fmt"
	messageException "gitlab.com/iotTracker/nerve/server/message/exception"
	"strconv"
)

const StartMarker = "7878"
const EndMarker = "0d0a"

type Message struct {
	Type       Type
	DataLength int64
	Data       string
}

func New(rawMessage string) (*Message, error) {
	var newMessage Message
	var err error
	if len(rawMessage) < 4 {
		return nil, messageException.Creation{Reasons: []string{"raw message string not long enough", rawMessage}}
	}

	newMessage.DataLength, err = strconv.ParseInt(rawMessage[:2], 16, 0)
	if err != nil {
		return nil, messageException.Creation{Reasons: []string{"length hex decoding", err.Error()}}
	}

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
	dataLengthHex := fmt.Sprintf("%x", m.DataLength)
	if len(dataLengthHex) == 1 {
		dataLengthHex = "0" + dataLengthHex
	}
	return hex.DecodeString(fmt.Sprintf(
		"%s%s%s%s%s",
		StartMarker,
		dataLengthHex,
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
		return fmt.Sprintf("[type: Heartbeat, Data: %s]", m.Data)
	case GPSPosition:
		return fmt.Sprintf("[type: GPS Position, Data: %s]", m.Data)
	case Status:
		return fmt.Sprintf("[type: Status, Data: %s]", m.Data)
	case Hibernation:
		return fmt.Sprintf("[type: Device Hibernation, Data: %s]", m.Data)
	default:

	}
	return ""
}
