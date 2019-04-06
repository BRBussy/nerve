package exception

import (
	"strings"
)

type Creation struct {
	Reasons []string
}

func (e Creation) Error() string {
	return "message creation error: " + strings.Join(e.Reasons, "; ")
}
