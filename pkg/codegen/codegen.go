package codegen

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/TheSadlig/asyncapi-codegen/pkg/asyncapi"
	"github.com/TheSadlig/asyncapi-codegen/pkg/asyncapi/parser"
	asyncapiv2 "github.com/TheSadlig/asyncapi-codegen/pkg/asyncapi/v2"
	asyncapiv3 "github.com/TheSadlig/asyncapi-codegen/pkg/asyncapi/v3"
	generatorv2 "github.com/TheSadlig/asyncapi-codegen/pkg/codegen/generators/v2"
	generatorv3 "github.com/TheSadlig/asyncapi-codegen/pkg/codegen/generators/v3"
	"github.com/TheSadlig/asyncapi-codegen/pkg/codegen/options"
	"github.com/TheSadlig/asyncapi-codegen/pkg/utils/template"
	"golang.org/x/tools/imports"
)

// CodeGen is the main structure for the code generation.
type CodeGen struct {
	specification asyncapi.Specification
	modulePath    string
	moduleVersion string
}

// FromFile returns a code generator from a specification file path.
func FromFile(path string, dependencies ...string) (CodeGen, error) {
	// Get specification from file
	spec, err := parser.FromFile(parser.FromFileParams{
		Path: path,
	})
	if err != nil {
		return CodeGen{}, err
	}

	// Get dependencies
	for _, path := range dependencies {
		dep, err := parser.FromFile(parser.FromFileParams{
			Path:         path,
			MajorVersion: spec.MajorVersion(),
		})
		if err != nil {
			return CodeGen{}, err
		}

		if err := spec.AddDependency(path, dep); err != nil {
			return CodeGen{}, err
		}
	}

	return New(spec)
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
	if err := template.SetConvertKeyFn(opt.ConvertKeys); err != nil {
		return err
	}

	if err := template.SetNamifyFn(opt.NamingScheme); err != nil {
		return err
	}

	// Process specification
	if err := cg.specification.Process(); err != nil {
		return err
	}

	// Generate content
	content, err := cg.generateContent(opt)
	if err != nil {
		return err
	}

	// Format content if not disabled
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

	// Write to file
	return os.WriteFile(opt.OutputPath, fileContent, 0644)
}

func (cg CodeGen) generateContent(opt options.Options) (string, error) {
	version := cg.specification.MajorVersion()
	switch version {
	case 2:
		spec, err := asyncapiv2.FromUnknownVersion(cg.specification)
		if err != nil {
			return "", err
		}

		return generatorv2.Generator{
			Specification: *spec,
			Options:       opt,
			ModulePath:    cg.modulePath,
			ModuleVersion: cg.moduleVersion,
		}.Generate()
	case 3:
		spec, err := asyncapiv3.FromUnknownVersion(cg.specification)
		if err != nil {
			return "", err
		}

		return generatorv3.Generator{
			Specification: *spec,
			Options:       opt,
			ModulePath:    cg.modulePath,
			ModuleVersion: cg.moduleVersion,
		}.Generate()
	default:
		return "", fmt.Errorf("unsupported major version (%q)", version)
	}
}
