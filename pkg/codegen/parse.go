package codegen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
)

// FromFile parses the AsyncAPI specification either from a YAML file or a JSON file
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

// FromYAML parses the AsyncAPI specification from a YAML file
func FromYAML(data []byte) (CodeGen, error) {
	data, err := yaml.YAMLToJSON(data)
	if err != nil {
		return CodeGen{}, err
	}

	return FromJSON(data)
}

// FromJSON parses the AsyncAPI specification from a JSON file
func FromJSON(data []byte) (CodeGen, error) {
	var spec asyncapi.Specification

	if err := json.Unmarshal(data, &spec); err != nil {
		return CodeGen{}, err
	}

	spec.Process()

	return New(spec), nil
}
