package notifications

import "errors"

var (
	// ErrInvalidMailDriver is an error given when an incorrect mail driver is selected
	ErrInvalidMailDriver = errors.New("Invalid mail driver")

	// ErrInvalidSMSDriver is an error given when an incorrect sms driver is selected
	ErrInvalidSMSDriver = errors.New("Invalid sms driver")

	// ErrInvalidPayloadType is an error given when an invalid payload is used
	ErrInvalidPayloadType = errors.New("Invalid notification payload type")
)
