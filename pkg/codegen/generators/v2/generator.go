package v2

import (
	"fmt"

	asyncapiv2 "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v2"
	"github.com/lerenn/asyncapi-codegen/pkg/codegen/options"
)

type Generator struct {
	Options       options.Options
	Specification asyncapiv2.Specification
	ModulePath    string
	ModuleVersion string
}

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

	return importsGenerator{
		PackageName:   opts.PackageName,
		ModuleVersion: g.ModuleVersion,
		ModuleName:    g.ModulePath,
		CustomImports: imps,
	}.Generate()
}

func (g Generator) generateTypes() (string, error) {
	return typesGenerator{Specification: g.Specification}.Generate()
}

func (g Generator) generateApp() (string, error) {
	var content string

	// Generate application subscriber
	subscriber, err := newSubscriberGenerator(
		SideIsApplication,
		g.Specification,
	).Generate()
	if err != nil {
		return "", err
	}
	content += subscriber

	// Generate application controller
	controller, err := newControllerGenerator(
		SideIsApplication,
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

	// Generate user subscriber
	subscriber, err := newSubscriberGenerator(
		SideIsUser,
		g.Specification,
	).Generate()
	if err != nil {
		return "", err
	}
	content += subscriber
	// Generate user controller
	controller, err := newControllerGenerator(
		SideIsUser,
		g.Specification,
	).Generate()
	if err != nil {
		return "", err
	}
	content += controller

	return content, nil
}
