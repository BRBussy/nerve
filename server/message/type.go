package message

type Type string

func (t Type) String() string {
	return string(t)
}

const Login Type = "01"
const Heartbeat Type = "08"
const GPSPosition Type = "10"
const GPSPosition2 Type = "11"
const Status Type = "13"
const Hibernation Type = "14"
const FactorySettingsRestored Type = "15"
const WhiteListTotal Type = "16"
const OfflineWIFIData Type = "17"
const TimeSynchronisation Type = "30"
const SetRemoteListeningCellNumber Type = "40"
const SetSOSNumber Type = "41"
const SetDadNumber Type = "42"
const SetMomNumber Type = "43"
const StopDataUpload Type = "44"
const SetGPSOffTime Type = "46"
const SetDoNotDisturb Type = "47"
const RestartDevice Type = "48"
const FindDevice Type = "49"
const SetAlarmClock Type = "50"
const TurnOffAlarmClock Type = "56"
const DeviceSetup Type = "57"
const WhiteListSynchronisation Type = "58"
const TurnOnLightSensorSwitch Type = "61"
const SetServerIPAndPort Type = "66"
const SetRecoveryPassword Type = "67"
const WIFIPosition Type = "69"
