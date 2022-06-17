package platform

import (
	"errors"
)

var (
	ErrContextValue = errors.New("context value is missing or has wrong type")
	ErrResponseType = errors.New("response value has wrong type")
)
