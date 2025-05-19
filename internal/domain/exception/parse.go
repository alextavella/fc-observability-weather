package exception

import "errors"

var (
	ErrParseBody = errors.New("invalid request body")
)
