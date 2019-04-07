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

func (m Message) HexString() string {
	dataLengthHex := fmt.Sprintf("%x", m.DataLength)
	if len(dataLengthHex) == 1 {
		dataLengthHex = "0" + dataLengthHex
	}
	return fmt.Sprintf(
		"%s %s %s %s %s",
		StartMarker,
		dataLengthHex,
		m.Type,
		m.Data,
		EndMarker,
	)
}

// 14:54:04
// 14:56:04

func (m Message) String() string {
	switch m.Type {
	case Login:
		return fmt.Sprintf("Type: Login - %s", m.HexString())
	case Heartbeat:
		return fmt.Sprintf("Type: Heartbeat - %s", m.HexString())
	case GPSPosition:
		return fmt.Sprintf("Type: GPS Position - %s", m.HexString())
	case GPSPosition2:
		return fmt.Sprintf("Type: GPS Position2 - %s", m.HexString())
	case Status:
		return fmt.Sprintf("Type: Status - %s", m.HexString())
	case Hibernation:
		return fmt.Sprintf("Type: Device Hibernation - %s", m.HexString())
	case FactorySettingsRestored:
		return fmt.Sprintf("Type: : Factory Settings - %s", m.HexString())
	case WhiteListTotal:
		return fmt.Sprintf("Type: : White List Total - %s", m.HexString())
	case OfflineWIFIData:
		return fmt.Sprintf("Type: Offline WIFI Data - %s", m.HexString())
	case TimeSynchronisation:
		return fmt.Sprintf("Type: Time Syncronisation - %s", m.HexString())
	case SetRemoteListeningCellNumber:
		return fmt.Sprintf("Type: Set Remote Listening Number - %s", m.HexString())
	case SetSOSNumber:
		return fmt.Sprintf("Type: Set SOS Number - %s", m.HexString())
	case SetDadNumber:
		return fmt.Sprintf("Type: Set Dad Number - %s", m.HexString())
	case SetMomNumber:
		return fmt.Sprintf("Type: Set Mom Number - %s", m.HexString())
	case StopDataUpload:
		return fmt.Sprintf("Type: Stop Data Upload - %s", m.HexString())
	case SetGPSOffTime:
		return fmt.Sprintf("Type: Set GPS Off Time - %s", m.HexString())
	case SetDoNotDisturb:
		return fmt.Sprintf("Type: Set Do Not Disturb - %s", m.HexString())
	case RestartDevice:
		return fmt.Sprintf("Type: Restart Device - %s", m.HexString())
	case FindDevice:
		return fmt.Sprintf("Type: Find Device - %s", m.HexString())
	case SetAlarmClock:
		return fmt.Sprintf("Type: Set Alarm Clock - %s", m.HexString())
	case TurnOffAlarmClock:
		return fmt.Sprintf("Type: Turn Off Alarm Clock - %s", m.HexString())
	case DeviceSetup:
		return fmt.Sprintf("Type: Device Setup - %s", m.HexString())
	case WhiteListSynchronisation:
		return fmt.Sprintf("Type: White List Synchronisation - %s", m.HexString())
	case TurnOnLightSensorSwitch:
		return fmt.Sprintf("Type: Turn On Light Sensor Switch - %s", m.HexString())
	case SetServerIPAndPort:
		return fmt.Sprintf("Type: Set Server IP and Port - %s", m.HexString())
	case SetRecoveryPassword:
		return fmt.Sprintf("Type: Set Recovery Password - %s", m.HexString())
	case WIFIPosition:
		return fmt.Sprintf("Type: WIFI Position - %s", m.HexString())
	case ManualPosition:
		return fmt.Sprintf("Type: Manual Position - %s", m.HexString())
	case ChargeComplete:
		return fmt.Sprintf("Type: Charge Complete - %s", m.HexString())
	case ChargerConnected:
		return fmt.Sprintf("Type: Charger Connected - %s", m.HexString())
	case ChargerDisconnected:
		return fmt.Sprintf("Type: Charger Disconnected - %s", m.HexString())
	case SetUploadInterval:
		return fmt.Sprintf("Type: Set Upload Interval - %s", m.HexString())

	default:
		return fmt.Sprintf("Type: unknown '%s' - %s", m.Type, m.HexString())
	}
}
