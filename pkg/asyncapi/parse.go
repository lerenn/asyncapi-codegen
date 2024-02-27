package asyncapi

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/ghodss/yaml"
	asyncapiv2 "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v2"
	asyncapiv3 "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v3"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
)

var (
	// ErrInvalidFileFormat is returned when using an invalid format for AsyncAPI specification.
	ErrInvalidFileFormat = fmt.Errorf("%w: invalid file format", extensions.ErrAsyncAPI)

	// ErrInvalidVersion is returned when the version is either unsupported or invalid.
	ErrInvalidVersion = fmt.Errorf("%w: unsupported/invalid version", extensions.ErrAsyncAPI)
)

// FromFile parses the AsyncAPI specification either from a YAML file or a JSON file.
//
// NOTE: It returns the Specification with filled fields, but this doesn't link
// references, apply traits, etc. You have to call method `Process` for this.
//
//nolint:ireturn,nolintlint
func FromFile(path string) (Specification, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	switch filepath.Ext(path) {
	case ".yaml", ".yml":
		return FromYAML(data)
	case ".json":
		return FromJSON(data)
	default:
		return nil, fmt.Errorf("%w: %q", ErrInvalidFileFormat, path)
	}
}

// FromYAML parses the AsyncAPI specification from a YAML file.
//
// NOTE: It returns the Specification with filled fields, but this doesn't link
// references, apply traits, etc. You have to call method `Process` for this.
//
//nolint:ireturn,nolintlint
func FromYAML(data []byte) (Specification, error) {
	// Change YAML to JSON
	data, err := yaml.YAMLToJSON(data)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	return FromJSON(data)
}

// FromJSON parses the AsyncAPI specification from a JSON file.
//
// NOTE: It returns the Specification with filled fields, but this doesn't link
// references, apply traits, etc. You have to call method `Process` for this.
//
//nolint:ireturn,nolintlint
func FromJSON(data []byte) (Specification, error) {
	// Check that the version is correct
	version, err := versionFromJSON(data)
	if err != nil {
		return nil, err
	}

	// Use a different specification based on the AsyncAPI version
	// NOTE: version should already be correct at this moment
	var spec Specification
	switch version[:1] {
	case "2":
		spec = &asyncapiv2.Specification{}
	case "3":
		spec = &asyncapiv3.Specification{}
	default:
		return nil, fmt.Errorf("unknown version (%q): this should not have happened", version)
	}

	// Parse JSON
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, err
	}

	return spec, nil
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
