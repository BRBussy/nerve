package message

import (
	"fmt"
	messageException "gitlab.com/iotTracker/nerve/server/message/exception"
	"strconv"
)

type Message struct {
	Type       Type
	DataLength uint64
	Data       string
}

func New(rawMessage string) (*Message, error) {
	var newMessage Message
	var err error
	if len(rawMessage) < 4 {
		return nil, messageException.Creation{Reasons: []string{"raw message string not long enough", rawMessage}}
	}

	newMessage.DataLength, err = strconv.ParseUint(rawMessage[:2], 16, 32)
	if err != nil {
		return nil, messageException.Creation{Reasons: []string{"data length parsing", err.Error()}}
	}

	if len(rawMessage) == 4 {
		newMessage.Type = Type(rawMessage[2:])
		newMessage.Data = ""
	} else {
		newMessage.Type = Type(rawMessage[2:4])
		newMessage.Data = rawMessage[4:]
	}
	if !ValidType[newMessage.Type] {
		return nil, messageException.Creation{Reasons: []string{"invalid type", newMessage.Type.String()}}
	}

	return &newMessage, nil
}

func (m Message) String() string {
	switch m.Type {
	case Login:
		return fmt.Sprintf("Login: IMEI: %s, SoftV: %s", m.Data[:16], m.Data[16:])
	case Heartbeat:
		return fmt.Sprintf("Heartbeat - ")
	case GPSPositioning:
		return fmt.Sprintf("GPSPositioning - ")
	default:

	}
	return ""
}
