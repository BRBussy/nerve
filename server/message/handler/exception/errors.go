package exception

import (
	"fmt"
	ServerMessage "gitlab.com/iotTracker/nerve/server/message"
	"strings"
)

type UnsupportedMessage struct {
	Reasons []string
	Message ServerMessage.Message
}

func (e UnsupportedMessage) Error() string {
	return fmt.Sprintf("unsupported message:\n%s\n%s", e.Message.String(), strings.Join(e.Reasons, "; "))
}