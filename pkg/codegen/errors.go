package codegen

import "errors"

var (
	// ErrInvalidBroker is an error raised when using an unknown broker.
	ErrInvalidBroker = errors.New("invalid broker")

	// ErrInvalidFileFormat is returned when using an invalid format for AsyncAPI specification.
	ErrInvalidFileFormat = errors.New("invalid file format")
)
