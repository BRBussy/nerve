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
const DeviceHibernation Type = "14"
const RestoreFactorySettings Type = "15"
const WhiteListTotal Type = "16"
const OfflineWIFIData Type = "17"
const TimeSynchronisation Type = "30"
const RemoteListeningNumber Type = "40"
const SOSNumber Type = "41"
const DadNumber Type = "42"
const MomNumber Type = "43"
const StopDataUpload Type = "44"
