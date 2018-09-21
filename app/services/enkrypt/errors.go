package enkrypt

import "errors"

var (
	// ErrInvalidKeySize is an error given when the key size incorrect
	ErrInvalidKeySize = errors.New("key must be 16, 24, or 32 bytes")

	// ErrInvalidAlgorithm is an error given when the selected algorithm is not available for use
	ErrInvalidAlgorithm = errors.New("encryption algorithm is not supported")
)
