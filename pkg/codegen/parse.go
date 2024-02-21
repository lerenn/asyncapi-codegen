package codegen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/ghodss/yaml"
	asyncapi "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v2"
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

	// Check that the version is correct
	_, err := versionFromJSON(data)
	if err != nil {
		return CodeGen{}, err
	}

	// Parse JSON
	if err := json.Unmarshal(data, &spec); err != nil {
		return CodeGen{}, err
	}

	// Process specification
	spec.Process()

	return New(spec)
}

func versionFromJSON(data []byte) (string, error) {
	var m map[string]any

	// Parse JSON
	if err := json.Unmarshal(data, &m); err != nil {
		return "", err
	}

	// Get version
	version, exists := m["asyncapi"]
	if !exists {
		return "", fmt.Errorf("%w: there is no 'asyncapi' field in specification", ErrInvalidVersion)
	}

	// Get stringed version
	versionStr, ok := version.(string)
	if !ok {
		return "", fmt.Errorf("%w: 'asyncapi' field is no string (%q)", ErrInvalidVersion, reflect.TypeOf(version))
	}

	// Check versions
	if !IsVersionSupported(versionStr) {
		return "", fmt.Errorf("%w: %q", ErrInvalidVersion, versionStr)
	}

	return versionStr, nil
}
