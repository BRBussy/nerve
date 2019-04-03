package handler

import (
	zx303ServerMessage "gitlab.com/iotTracker/nerve/server/zx303/message"
)

type Handler interface {
	Handle(message *zx303ServerMessage.Message) (*zx303ServerMessage.Message, error)
}
