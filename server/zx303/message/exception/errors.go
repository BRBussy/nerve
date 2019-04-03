package exception

import "strings"

type Creation struct {
	Reasons []string
}

func (e Creation) Error() string {
	return "message creation error: " + strings.Join(e.Reasons, "; ")
}

type Invalid struct {
	Reasons []string
}

func (e Invalid) Error() string {
	return "message invalid: " + strings.Join(e.Reasons, "; ")
}
