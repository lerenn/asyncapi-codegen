package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/TheSadlig/asyncapi-codegen/pkg/asyncapi"
	asyncapiv2 "github.com/TheSadlig/asyncapi-codegen/pkg/asyncapi/v2"
	asyncapiv3 "github.com/TheSadlig/asyncapi-codegen/pkg/asyncapi/v3"
	"github.com/TheSadlig/asyncapi-codegen/pkg/extensions"
	"github.com/ghodss/yaml"
)

var (
	// ErrInvalidFileFormat is returned when using an invalid format for AsyncAPI specification.
	ErrInvalidFileFormat = fmt.Errorf("%w: invalid file format", extensions.ErrAsyncAPI)

	// ErrInvalidVersion is returned when the version is either unsupported or invalid.
	ErrInvalidVersion = fmt.Errorf("%w: unsupported/invalid version", extensions.ErrAsyncAPI)
)

// FromFileParams are the parameters to parse an AsyncAPI specification from a file.
type FromFileParams struct {
	// Path to the file that contains the AsyncAPI specification.
	Path string
	// MajorVersion is the major version of the AsyncAPI specification.
	// If it is 0, it will try to get it from the specification.
	MajorVersion int
}

// FromFile parses the AsyncAPI specification either from a YAML file or a JSON file.
// If there is no version provided, it will try to get it from the specification.
//
// NOTE: It returns the Specification with filled fields, but this doesn't
// generate metadata, link references, apply traits, etc. You have to call
// method `Process` for this.
//
//nolint:ireturn,nolintlint
func FromFile(params FromFileParams) (asyncapi.Specification, error) {
	data, err := os.ReadFile(params.Path)
	if err != nil {
		return nil, err
	}

	switch filepath.Ext(params.Path) {
	case ".yaml", ".yml":
		return FromYAML(FromYAMLParams{
			Data:         data,
			MajorVersion: params.MajorVersion,
		})
	case ".json":
		return FromJSON(FromJSONParams{
			Data:         data,
			MajorVersion: params.MajorVersion,
		})
	default:
		return nil, fmt.Errorf("%w: %q", ErrInvalidFileFormat, params.MajorVersion)
	}
}

// FromYAMLParams are the parameters to parse an AsyncAPI specification from a YAML file.
type FromYAMLParams struct {
	// Cata is the content of the YAML that contains the AsyncAPI specification.
	Data []byte
	// MajorVersion is the major version of the AsyncAPI specification.
	// If it is 0, it will try to get it from the specification.
	MajorVersion int
}

// FromYAML parses the AsyncAPI specification from a YAML file.
// If there is no version provided, it will try to get it from the specification.
//
// NOTE: It returns the Specification with filled fields, but this doesn't link
// references, apply traits, etc. You have to call method `Process` for this.
//
//nolint:ireturn,nolintlint
func FromYAML(params FromYAMLParams) (asyncapi.Specification, error) {
	// Change YAML to JSON
	data, err := yaml.YAMLToJSON(params.Data)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	return FromJSON(FromJSONParams{
		Data:         data,
		MajorVersion: params.MajorVersion,
	})
}

// FromJSONParams are the parameters to parse an AsyncAPI specification from a JSON file.
type FromJSONParams struct {
	// Cata is the content of the JSON that contains the AsyncAPI specification.
	Data []byte
	// MajorVersion is the major version of the AsyncAPI specification.
	// If it is 0, it will try to get it from the specification.
	MajorVersion int
}

// FromJSON parses the AsyncAPI specification from a JSON file.
// If there is no version provided, it will try to get it from the specification.
//
// NOTE: It returns the Specification with filled fields, but this doesn't link
// references, apply traits, etc. You have to call method `Process` for this.
//
//nolint:ireturn,nolintlint
func FromJSON(params FromJSONParams) (asyncapi.Specification, error) {
	// Check that the version is correct
	majorVersion := params.MajorVersion
	if majorVersion == 0 {
		v, err := majorVersionFromJSON(params.Data)
		if err != nil {
			return nil, err
		}
		majorVersion = v
	}

	// Use a different specification based on the AsyncAPI version
	// NOTE: version should already be correct at this moment
	var spec asyncapi.Specification
	switch majorVersion {
	case 2:
		spec = asyncapiv2.NewSpecification()
	case 3:
		spec = asyncapiv3.NewSpecification()
	default:
		return nil, fmt.Errorf("unknown version (%q): this should not have happened", majorVersion)
	}

	// Parse JSON
	if err := json.Unmarshal(params.Data, &spec); err != nil {
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
	if !asyncapi.IsVersionSupported(versionStr) {
		return "", fmt.Errorf("%w: %q", ErrInvalidVersion, versionStr)
	}

	return versionStr, nil
}

func majorVersionFromJSON(data []byte) (int, error) {
	versionStr, err := versionFromJSON(data)
	if err != nil {
		return 0, err
	}

	v, err := strconv.Atoi(versionStr[:1])
	if err != nil {
		return 0, err
	}

	return v, nil
}
