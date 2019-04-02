package server

type Server interface {
	Start(request *StartRequest) error
}

type StartRequest struct {
	Port      string
	IPAddress string
}
