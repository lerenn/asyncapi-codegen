package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/lerenn/asyncapi-codegen/pkg/codegen"
)

var (
	ErrInvalidGenerate = errors.New("invalid generate argument")
)

type Flags struct {
	InputPath   string
	OutputPath  string
	PackageName string
	Generate    string
	Broker      string
}

func ProcessFlags() Flags {
	var f Flags

	flag.StringVar(&f.InputPath, "i", "asyncapi.yaml", "AsyncAPI specification file to use")
	flag.StringVar(&f.OutputPath, "o", "asyncapi.gen.go", "Destination file")
	flag.StringVar(&f.PackageName, "p", "asyncapi", "Golang package name")
	flag.StringVar(&f.Generate, "g", "client,application,broker,types", "Generation options")

	flag.Parse()

	return f
}

func (f Flags) ToCodegenOptions() (codegen.Options, error) {
	opt := codegen.Options{
		OutputPath:  f.OutputPath,
		PackageName: f.PackageName,
	}

	if f.Generate != "" {
		gens := strings.Split(f.Generate, ",")
		for _, v := range gens {
			switch v {
			// Generic code
			case "application":
				opt.Generate.Application = true
			case "broker":
				opt.Generate.Broker = true
			case "client":
				opt.Generate.Client = true
			case "types":
				opt.Generate.Types = true
			// Broker implementations
			case "nats":
				opt.Generate.NATS = true
			default:
				return opt, fmt.Errorf("%w: %q", ErrInvalidGenerate, v)
			}
		}
	}

	return opt, nil
}
