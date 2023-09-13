package codegen

import (
	"github.com/lerenn/asyncapi-codegen/pkg/codegen/generators"
)

// Options is the struct that gather configuration of codegen.
type Options struct {
	// OutputPath is the path to the generated code file
	OutputPath string

	// PackageName is the package name of the generated code
	PackageName string

	// Generate contains options regarding which golang code should be generated
	Generate generators.Options

	// DisableFormatting states if the formatting should be disabled when
	// writing the generated code
	DisableFormatting bool
}
