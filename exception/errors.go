package exception

import "strings"

type Unexpected struct {
	Reasons []string
}

func (e Unexpected) Error() string {
	return "unexpected error: " + strings.Join(e.Reasons, "; ")
}
