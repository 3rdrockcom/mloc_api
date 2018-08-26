package payments

import "errors"

var (
	// ErrInvalidPayloadType is an error given when an invalid payload is used
	ErrInvalidPayloadType = errors.New("Invalid payment payload type")
)
