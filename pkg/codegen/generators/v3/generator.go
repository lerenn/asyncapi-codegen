package generatorv3

import (
	"fmt"

	asyncapi "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v3"
	"github.com/lerenn/asyncapi-codegen/pkg/codegen/generators"
	"github.com/lerenn/asyncapi-codegen/pkg/codegen/options"
)

// Generator is the structure that contains information to generate the code from
// the specification.
type Generator struct {
	Options       options.Options
	Specification asyncapi.Specification
	ModulePath    string
	ModuleVersion string
}

// Generate generates the source code from the specification.
func (g Generator) Generate() (string, error) {
	content, err := g.generateImports(g.Options)
	if err != nil {
		return "", err
	}

	for remainingParts, part := true, ""; remainingParts; part = "" {
		switch {
		case g.Options.Generate.Application:
			part, err = g.generateApp()
			g.Options.Generate.Application = false
		case g.Options.Generate.User:
			part, err = g.generateUser()
			g.Options.Generate.User = false
		case g.Options.Generate.Types:
			part, err = g.generateTypes()
			g.Options.Generate.Types = false
		default:
			remainingParts = false
		}

		if err != nil {
			return "", err
		}

		content += part
	}

	return content, nil
}

func (g Generator) generateImports(opts options.Options) (string, error) {
	imps, err := g.Specification.CustomImports()
	if err != nil {
		return "", fmt.Errorf("failed to generate custom imports: %w", err)
	}

	var requiredImports []string

	if generators.DateTimeFormatInSpec(g.Specification) {
		requiredImports = append(requiredImports, "\"cloud.google.com/go/civil\"")
	}

	return ImportsGenerator{
		PackageName:     opts.PackageName,
		ModuleVersion:   g.ModuleVersion,
		ModuleName:      g.ModulePath,
		RequiredImports: requiredImports,
		CustomImports:   imps,
	}.Generate()
}

func (g Generator) generateTypes() (string, error) {
	return TypesGenerator{Specification: g.Specification}.Generate()
}

func (g Generator) generateApp() (string, error) {
	var content string

	// Generate application listener
	listener, err := NewSubscriberGenerator(
		generators.SideIsApplication,
		g.Specification,
	).Generate()
	if err != nil {
		return "", err
	}
	content += listener

	// Generate application controller
	controller, err := NewControllerGenerator(
		generators.SideIsApplication,
		g.Specification,
	).Generate()
	if err != nil {
		return "", err
	}
	content += controller

	return content, nil
}

func (g Generator) generateUser() (string, error) {
	var content string

	// Generate user listener
	listener, err := NewSubscriberGenerator(
		generators.SideIsUser,
		g.Specification,
	).Generate()
	if err != nil {
		return "", err
	}
	content += listener
	// Generate user controller
	controller, err := NewControllerGenerator(
		generators.SideIsUser,
		g.Specification,
	).Generate()
	if err != nil {
		return "", err
	}
	content += controller

	return content, nil
}
