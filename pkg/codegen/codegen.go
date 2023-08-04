package codegen

import (
	"os"
	"runtime/debug"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
	"github.com/lerenn/asyncapi-codegen/pkg/codegen/generators"
	"golang.org/x/tools/imports"
)

// CodeGen is the main structure for the code generation
type CodeGen struct {
	Specification asyncapi.Specification
	ModulePath    string
	ModuleVersion string
}

// New creates a new code generation structure that can be used to generate code
func New(spec asyncapi.Specification) CodeGen {
	modulePath := "unknown module path"
	moduleVersion := "unknown version"
	if bi, ok := debug.ReadBuildInfo(); ok {
		if bi.Main.Path != "" {
			modulePath = bi.Main.Path
		}
		if bi.Main.Version != "" {
			moduleVersion = bi.Main.Version
		}
	}

	return CodeGen{
		Specification: spec,
		ModulePath:    modulePath,
		ModuleVersion: moduleVersion,
	}
}

// Generate generates code from the code generation structure, that have already
// processed the AsyncAPI file when creating it
func (cg CodeGen) Generate(opt Options) error {
	content, err := cg.generateImports(opt)
	if err != nil {
		return err
	}

	for remainingParts, part := true, ""; remainingParts; part = "" {
		switch {
		case opt.Generate.Application:
			part, err = cg.generateApp(opt)
			opt.Generate.Application = false
		case opt.Generate.Client:
			part, err = cg.generateClient(opt)
			opt.Generate.Client = false
		case opt.Generate.Broker:
			part, err = generators.BrokerControllerGenerator{}.Generate()
			opt.Generate.Broker = false
		case opt.Generate.Types:
			part, err = cg.generateTypes()
			opt.Generate.Types = false
		case opt.Generate.NATS:
			part, err = generators.BrokerNATSGenerator{}.Generate()
			opt.Generate.NATS = false
		default:
			remainingParts = false
		}

		if err != nil {
			return err
		}

		content += part
	}

	var fileContent []byte
	if !opt.DisableFormatting {
		fileContent, err = imports.Process("", []byte(content), &imports.Options{
			TabWidth:  8,
			TabIndent: true,
			Comments:  true,
			Fragment:  true,
		})
		if err != nil {
			return err
		}
	} else {
		fileContent = []byte(content)
	}

	return os.WriteFile(opt.OutputPath, fileContent, 0644)
}

func (cg CodeGen) generateImports(opts Options) (string, error) {
	return generators.ImportsGenerator{
		PackageName:   opts.PackageName,
		ModuleVersion: cg.ModuleVersion,
		ModuleName:    cg.ModulePath,
	}.Generate()
}

func (cg CodeGen) generateTypes() (string, error) {
	return generators.TypesGenerator{Specification: cg.Specification}.Generate()
}

func (cg CodeGen) generateApp(opts Options) (string, error) {
	var content string

	// Generate application subscriber
	subscriber, err := generators.NewSubscriberGenerator(
		generators.SideIsApplication,
		cg.Specification,
	).Generate()
	if err != nil {
		return "", err
	}
	content += subscriber

	// Generate application controller
	controller, err := generators.NewControllerGenerator(
		generators.SideIsApplication,
		cg.Specification,
	).Generate()
	if err != nil {
		return "", err
	}
	content += controller

	return content, nil
}

func (cg CodeGen) generateClient(opts Options) (string, error) {
	var content string

	// Generate client subscriber
	subscriber, err := generators.NewSubscriberGenerator(
		generators.SideIsClient,
		cg.Specification,
	).Generate()
	if err != nil {
		return "", err
	}
	content += subscriber
	// Generate client controller
	controller, err := generators.NewControllerGenerator(
		generators.SideIsClient,
		cg.Specification,
	).Generate()
	if err != nil {
		return "", err
	}
	content += controller

	return content, nil
}
