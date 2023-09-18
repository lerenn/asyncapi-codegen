package codegen

import (
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
)

var (
	// ErrInvalidBroker is an error raised when using an unknown broker.
	ErrInvalidBroker = fmt.Errorf("%w: invalid broker", extensions.ErrAsyncAPI)

	// ErrInvalidFileFormat is returned when using an invalid format for AsyncAPI specification.
	ErrInvalidFileFormat = fmt.Errorf("%w: invalid file format", extensions.ErrAsyncAPI)
)
