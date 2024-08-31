package csrf

import (
	"golang.org/x/net/xsrftoken"
)

const (
	CsrfTokenHeader = "X-CSRF-TOKEN"
)

func Generate(actionID string) string {
	return xsrftoken.Generate(Key(), "none", actionID)
}

func Valid(token, actionID string) bool {
	return xsrftoken.Valid(token, Key(), "none", actionID)
}
