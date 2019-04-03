package exception

import (
	"fmt"
	"strings"
)

type StartOfMessageNotFound struct {
	Data    string
	Reasons []string
}

func (e StartOfMessageNotFound) Error() string {
	return fmt.Sprintf("start of message not found in given data:\n'%s'\n%s", e.Data, strings.Join(e.Reasons, "; "))
}

type DecodingError struct {
	Reasons []string
}

func (e DecodingError) Error() string {
	return "decoding error: " + strings.Join(e.Reasons, "; ")
}
