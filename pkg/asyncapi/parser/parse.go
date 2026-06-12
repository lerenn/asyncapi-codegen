package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/ghodss/yaml"
	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
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
		return nil, fmt.Errorf("%w: unsupported file extension %q", ErrInvalidFileFormat, filepath.Ext(params.Path))
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
	// An OpenAPI document can be used as a schema-only dependency: only its
	// "components" are kept, so that its schemas can be referenced from an
	// AsyncAPI spec even when the document contains fields AsyncAPI does not
	// allow (e.g. an array-typed "servers"). It is parsed with the same major
	// version as the spec that depends on it.
	if isOpenAPISpecification(params.Data) {
		return fromOpenAPISpecification(params.Data, params.MajorVersion)
	}

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
		return nil, fmt.Errorf("unknown version (%d): this should not have happened", majorVersion)
	}

	// Parse JSON
	if err := json.Unmarshal(params.Data, &spec); err != nil {
		return nil, err
	}

	return spec, nil
}

// isOpenAPISpecification reports whether the given JSON is an OpenAPI document
// (it has an "openapi" field) rather than an AsyncAPI one.
func isOpenAPISpecification(data []byte) bool {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return false
	}

	_, hasOpenAPI := m["openapi"]
	_, hasAsyncAPI := m["asyncapi"]

	return hasOpenAPI && !hasAsyncAPI
}

// fromOpenAPISpecification parses only the "components" of an OpenAPI document.
// Everything else (servers, paths, ...) is ignored, so the document can be used
// purely as a source of referenceable schemas. The components are parsed as an
// AsyncAPI specification of the given major version (defaulting to 3), so that
// the dependency can be added to a spec of the same version.
//
//nolint:ireturn,nolintlint
func fromOpenAPISpecification(data []byte, majorVersion int) (asyncapi.Specification, error) {
	var doc struct {
		Components json.RawMessage `json:"components"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, err
	}

	if majorVersion == 0 {
		majorVersion = 3
	}

	var spec asyncapi.Specification
	switch majorVersion {
	case 2:
		spec = asyncapiv2.NewSpecification()
	case 3:
		spec = asyncapiv3.NewSpecification()
	default:
		return nil, fmt.Errorf("%w: unsupported major version (%d) for OpenAPI dependency", ErrInvalidVersion, majorVersion)
	}

	if len(doc.Components) == 0 {
		return spec, nil
	}

	// Reuse the regular components unmarshaling by wrapping them in a minimal
	// specification.
	wrapper := []byte(`{"components":`)
	wrapper = append(wrapper, doc.Components...)
	wrapper = append(wrapper, '}')

	if err := json.Unmarshal(wrapper, &spec); err != nil {
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
