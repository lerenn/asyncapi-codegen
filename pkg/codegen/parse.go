package codegen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/asyncapi/parser-go/pkg/parser"
	"github.com/ghodss/yaml"
	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
)

// FromFile parses the AsyncAPI specification either from a YAML file or a JSON file.
func FromFile(path string) (CodeGen, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return CodeGen{}, err
	}

	switch filepath.Ext(path) {
	case ".yaml", ".yml":
		return FromYAML(data)
	case ".json":
		return FromJSON(data)
	default:
		return CodeGen{}, fmt.Errorf("%w: %q", ErrInvalidFileFormat, path)
	}
}

// FromYAML parses the AsyncAPI specification from a YAML file.
func FromYAML(data []byte) (CodeGen, error) {
	// Verify specification
	if err := verifySpecificationData(data); err != nil {
		return CodeGen{}, err
	}

	// Change YAML to JSON
	data, err := yaml.YAMLToJSON(data)
	if err != nil {
		return CodeGen{}, err
	}

	// Parse JSON
	return FromJSON(data)
}

// FromJSON parses the AsyncAPI specification from a JSON file.
func FromJSON(data []byte) (CodeGen, error) {
	var spec asyncapi.Specification

	// Verify specification
	if err := verifySpecificationData(data); err != nil {
		return CodeGen{}, err
	}

	// Parse JSON
	if err := json.Unmarshal(data, &spec); err != nil {
		return CodeGen{}, err
	}

	// Process specification
	spec.Process()

	return New(spec), nil
}

func verifySpecificationData(data []byte) error {
	// Create a new parser
	p, err := parser.New()
	if err != nil {
		return err
	}

	// Return the result of the parsing
	return p(bytes.NewReader(data), os.Stdout)
}
