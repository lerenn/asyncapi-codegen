package codegen

import (
	"go/format"
	"os"
	"runtime/debug"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
	"github.com/lerenn/asyncapi-codegen/pkg/codegen/generators"
)

type CodeGen struct {
	Specification asyncapi.Specification
	ModulePath    string
	ModuleVersion string
}

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

	buf, err := format.Source([]byte(content))
	if err != nil {
		return err
	}

	return os.WriteFile(opt.OutputPath, buf, 0755)
}

func (cg CodeGen) generateImports(opts Options) (string, error) {
	return generators.ImportsGenerator{
		Options:       opts.Generate,
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

	// Generate application input
	subscriber, err := generators.NewAppSubscriberGenerator(cg.Specification).Generate()
	if err != nil {
		return "", err
	}
	content += subscriber

	// Generate application output
	controller, err := generators.NewAppControllerGenerator(cg.Specification).Generate()
	if err != nil {
		return "", err
	}
	content += controller

	return content, nil
}

func (cg CodeGen) generateClient(opts Options) (string, error) {
	var content string

	// Generate client input
	subscriber, err := generators.NewClientSubscriberGenerator(cg.Specification).Generate()
	if err != nil {
		return "", err
	}
	content += subscriber

	// Generate client output
	controller, err := generators.NewClientControllerGenerator(cg.Specification).Generate()
	if err != nil {
		return "", err
	}
	content += controller

	return content, nil
}
