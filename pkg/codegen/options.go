package codegen

import (
	"errors"

	"github.com/lerenn/asyncapi-codegen/pkg/codegen/generators"
)

var (
	ErrInvalidBroker = errors.New("invalid broker")
)

type Options struct {
	OutputPath  string
	PackageName string
	Generate    generators.Options
}
