package codegen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/ghodss/yaml"
	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
	asyncapiv2 "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v2"
	asyncapiv3 "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v3"
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
	// Check that the version is correct
	version, err := versionFromJSON(data)
	if err != nil {
		return CodeGen{}, err
	}

	// Use a different specification based on the AsyncAPI version
	// NOTE: version should already be correct at this moment
	var spec asyncapi.Specification
	switch version[:1] {
	case "2":
		spec = &asyncapiv2.Specification{}
	case "3":
		spec = &asyncapiv3.Specification{}
	default:
		return CodeGen{}, fmt.Errorf("unknown version (%q): this should not have happened", version)
	}

	// Parse JSON
	if err := json.Unmarshal(data, &spec); err != nil {
		return CodeGen{}, err
	}

	// Process specification
	spec.Process()

	// Return a new codegen with this spec
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
