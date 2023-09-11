package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/lerenn/asyncapi-codegen/pkg/codegen"
)

var (
	// ErrInvalidGenerate happens when using an invalid generation argument
	ErrInvalidGenerate = errors.New("invalid generate argument")
)

// Flags contains all command line flags
type Flags struct {
	// InputPath is the path of the AsyncAPI specification file
	InputPath string

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
}

// ProcessFlags processes command line flags and fill the Flags structure with them
func ProcessFlags() Flags {
	var f Flags

	flag.StringVar(&f.InputPath, "i", "asyncapi.yaml", "AsyncAPI specification file to use")
	flag.StringVar(&f.OutputPath, "o", "asyncapi.gen.go", "Destination file")
	flag.StringVar(&f.PackageName, "p", "asyncapi", "Golang package name")
	flag.StringVar(&f.Generate, "g", "user,application,types", "Generation options")
	flag.BoolVar(&f.DisableFormatting, "disable-formatting", false, "Disables the code generation formatting")

	flag.Parse()

	return f
}

// ToCodegenOptions processes command line flags structure to code generation tool options
func (f Flags) ToCodegenOptions() (codegen.Options, error) {
	opt := codegen.Options{
		OutputPath:        f.OutputPath,
		PackageName:       f.PackageName,
		DisableFormatting: f.DisableFormatting,
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
