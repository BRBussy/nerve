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
	case FactorySettingsRestored:
		return fmt.Sprintf("[type : Factory Settings, Data: %s]", m.Data)
	case WhiteListTotal:
		return fmt.Sprintf("[type : White List Total, Data: %s]", m.Data)
	case OfflineWIFIData:
		return fmt.Sprintf("[type: Offline WIFI Data, Data: %s]", m.Data)
	case TimeSynchronisation:
		return fmt.Sprintf("[type: Time Syncronisation, Data: %s]", m.Data)
	case SetRemoteListeningCellNumber:
		return fmt.Sprintf("[type: Set Remote Listening Number, Data: %s]", m.Data)
	case SetSOSNumber:
		return fmt.Sprintf("[type: Set SOS Number, Data: %s]", m.Data)
	case SetDadNumber:
		return fmt.Sprintf("[type: Set Dad Number, Data: %s]", m.Data)
	case SetMomNumber:
		return fmt.Sprintf("[type: Set Mom Number, Data: %s]", m.Data)
	case StopDataUpload:
		return fmt.Sprintf("[type: Stop Data Upload, Data: %s]", m.Data)
	case SetGPSOffTime:
		return fmt.Sprintf("[type: Set GPS Off Time, Data: %s]", m.Data)
	case SetDoNotDisturb:
		return fmt.Sprintf("[type: Set Do Not Disturb, Data: %s]", m.Data)
	case RestartDevice:
		return fmt.Sprintf("[type: Restart Device, Data: %s]", m.Data)
	case FindDevice:
		return fmt.Sprintf("[type: Find Device, Data: %s]", m.Data)
	case SetAlarmClock:
		return fmt.Sprintf("[type: Set Alarm Clock, Data: %s]", m.Data)
	case TurnOffAlarmClock:
		return fmt.Sprintf("[type: Turn Off Alarm Clock, Data: %s]", m.Data)
	case DeviceSetup:
		return fmt.Sprintf("[type: Device Setup, Data: %s]", m.Data)
	case WhiteListSynchronisation:
		return fmt.Sprintf("[type: White List Synchronisation, Data: %s]", m.Data)
	case TurnOnLightSensorSwitch:
		return fmt.Sprintf("[type: Turn On Light Sensor Switch, Data: %s]", m.Data)
	case SetServerIPAndPort:
		return fmt.Sprintf("[type: Set Server IP and Port, Data: %s]", m.Data)
	case SetRecoveryPassword:
		return fmt.Sprintf("[type: Set Recovery Password, Data: %s]", m.Data)
	case WIFIPosition:
		return fmt.Sprintf("[type: WIFI Position, Data: %s]", m.Data)

	default:
		return fmt.Sprintf("[type: unknown, Data: %s]", m.Data)
	}
}
