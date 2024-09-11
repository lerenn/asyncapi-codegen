package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/TheSadlig/asyncapi-codegen/pkg/codegen/options"
	"github.com/spf13/cobra"
)

var (
	// ErrInvalidGenerate happens when using an invalid generation argument.
	ErrInvalidGenerate = errors.New("invalid generate argument")
)

// Flags contains all command line flags.
type Flags struct {
	// InputPaths are the path of the AsyncAPI specification file and its dependencies
	InputPaths []string

	// OutputPath is the path of the generated code file
	OutputPath string

	// PackageName is the package name of the generated code
	PackageName string

	// Generate contains options, separated by commas, regarding which
	// golang code should be generated
	Generate string

	// Broker contains the broker name whose code should be generated
	Broker string

	// DisableFormatting states if the formatting should be disabled when
	// writing the generated code
	DisableFormatting bool

	// ConvertKeys defines a schema property keys conversion strategy.
	// Supported values: snake, camel, kebab, none
	ConvertKeys string

	// NamingScheme defines the naming case for generated golang structs
	// Supported values: camel, none
	NamingScheme string

	// IgnoreStringFormat states whether the properties' format (date, date-time) should impact the type in types
	IgnoreStringFormat bool
}

// SetToCommand adds the flags to a cobra command.
func (f *Flags) SetToCommand(cmd *cobra.Command) {
	cmd.Flags().StringSliceVarP(
		&f.InputPaths, "input", "i", []string{"asyncapi.yaml"},
		"AsyncAPI specification file to use, and its dependencies")
	cmd.Flags().StringVarP(&f.OutputPath, "output", "o", "asyncapi.gen.go", "Destination file")
	cmd.Flags().StringVarP(&f.PackageName, "package", "p", "asyncapi", "Golang package name")
	cmd.Flags().StringVarP(&f.Generate, "generate", "g", "user,application,types", "Generation options")
	cmd.Flags().BoolVarP(&f.DisableFormatting, "disable-formatting", "f", false, "Disables the code generation formatting")
	cmd.Flags().StringVarP(&f.ConvertKeys, "convert-keys", "c", "none",
		"Schema property key names conversion strategy.\nSupported values: snake, camel, kebab, none.")
	cmd.Flags().StringVarP(&f.NamingScheme, "naming-scheme", "n", "none",
		"Naming scheme for generated golang elements.\nSupported values: camel, none.")
	cmd.Flags().BoolVar(&f.IgnoreStringFormat, "ignore-string-format", false,
		"Ignores the format (date, date-time) on string properties, generating golang string, instead of dates")
}

// ToCodegenOptions processes command line flags structure to code generation tool options.
func (f Flags) ToCodegenOptions() (options.Options, error) {
	opt := options.Options{
		OutputPath:         f.OutputPath,
		PackageName:        f.PackageName,
		DisableFormatting:  f.DisableFormatting,
		ConvertKeys:        f.ConvertKeys,
		NamingScheme:       f.NamingScheme,
		IgnoreStringFormat: f.IgnoreStringFormat,
	}

	if f.Generate != "" {
		gens := strings.Split(f.Generate, ",")
		for _, v := range gens {
			switch v {
			case "application":
				opt.Generate.Application = true
			case "user":
				opt.Generate.User = true
			case "types":
				opt.Generate.Types = true
			default:
				return opt, fmt.Errorf("%w: %q", ErrInvalidGenerate, v)
			}
		}
	}

	return opt, nil
}
