package exception

import (
	"fmt"
	zx303ServerMessage "gitlab.com/iotTracker/nerve/server/zx303/message"
	"strings"
)

type UnsupportedMessage struct {
	Reasons []string
	Message zx303ServerMessage.Message
}

func (e UnsupportedMessage) Error() string {
	return fmt.Sprintf("unsupported message:\n%s\n%s", e.Message.String(), strings.Join(e.Reasons, "; "))
}
