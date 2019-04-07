package exception

import (
	"strings"
)

type Listen struct {
	Reasons []string
}

func (e Listen) Error() string {
	return "listening error: " + strings.Join(e.Reasons, "; ")
}

type AcceptConnection struct {
	Reasons []string
}

func (e AcceptConnection) Error() string {
	return "accept connection error: " + strings.Join(e.Reasons, "; ")
}
