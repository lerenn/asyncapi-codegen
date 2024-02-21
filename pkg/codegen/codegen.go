package codegen

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
	asyncapiv2 "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v2"
	generatorsv2 "github.com/lerenn/asyncapi-codegen/pkg/codegen/generators/v2"
	"github.com/lerenn/asyncapi-codegen/pkg/codegen/options"
	"golang.org/x/tools/imports"
)

// CodeGen is the main structure for the code generation.
type CodeGen struct {
	specification asyncapi.Specification
	modulePath    string
	moduleVersion string
}

// New creates a new code generation structure that can be used to generate code.
func New(spec asyncapi.Specification) (CodeGen, error) {
	modulePath, moduleVersion := modulePathVersion()

	return CodeGen{
		specification: spec,
		modulePath:    modulePath,
		moduleVersion: moduleVersion,
	}, nil
}

func modulePathVersion() (path, version string) {
	path = "unknown module path"
	version = "unknown version"
	if bi, ok := debug.ReadBuildInfo(); ok {
		if bi.Main.Path != "" {
			path = bi.Main.Path
		}
		if bi.Main.Version != "" {
			version = bi.Main.Version
		}
	}

	return path, version
}

// Generate generates code from the code generation structure, that have already
// processed the AsyncAPI file when creating it.
func (cg CodeGen) Generate(opt options.Options) error {
	content, err := cg.generateContent(opt)
	if err != nil {
		return err
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

func (cg CodeGen) generateContent(opt options.Options) (string, error) {
	version := cg.specification.AsyncAPIVersion()

	// NOTE: version should already be correct at this moment
	switch version[:2] {
	case "v2":
		spec, ok := cg.specification.(asyncapiv2.Specification)
		if !ok {
			return "", fmt.Errorf("unknown spec format: this should not have happened")
		}

		return generatorsv2.Generator{
			Specification: spec,
			Options:       opt,
			ModulePath:    cg.modulePath,
			ModuleVersion: cg.moduleVersion,
		}.Generate()
	default:
		return "", fmt.Errorf("unkown version (%q): this should not have happened", version)
	}
}
