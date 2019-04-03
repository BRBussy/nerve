package message

type Type string

func (t Type) String() string {
	return string(t)
}

const Login Type = "01"
const HeartBeat Type = "08"
const GPSPositioning Type = "10"
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

var ValidType = map[Type]bool{
	Login:                  true,
	HeartBeat:              true,
	GPSPositioning:         true,
	Status:                 true,
	DeviceHibernation:      true,
	RestoreFactorySettings: true,
	WhiteListTotal:         true,
	OfflineWIFIData:        true,
	TimeSynchronisation:    true,
	RemoteListeningNumber:  true,
	SOSNumber:              true,
	DadNumber:              true,
	MomNumber:              true,
	StopDataUpload:         true,
}
