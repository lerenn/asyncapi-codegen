package generators

import (
	"strings"

	"github.com/lerenn/asyncapi-codegen/pkg/utils"
	"github.com/lerenn/asyncapi-codegen/pkg/utils/template"
)

// These mirror the asyncapi v2/v3 SchemaType values. They are duplicated here as
// plain strings so that this package stays version-agnostic (it must not import
// the concrete v2/v3 packages, which would create an import cycle with the
// template packages that depend on it).
const (
	schemaTypeObject = "object"
	schemaTypeArray  = "array"
)

// messageFieldHeader mirrors the v2/v3 MessageFieldIsHeader value. A reference
// segment equal to it is rendered as the generated "headers" struct field.
const messageFieldHeader = "header"

// TemplateSchema is the minimal view of an asyncapi Schema (v2 or v3) needed by
// the template helpers shared between both versions. The type parameter T is the
// concrete Schema type, so child schemas keep their concrete type. Both
// asyncapi/v2.Schema and asyncapi/v3.Schema implement TemplateSchema of themselves.
type TemplateSchema[T any] interface {
	GetType() string
	GetIsRequired() bool
	IsFieldRequired(field string) bool
	GetProperties() map[string]*T
	GetItems() *T
	GetAdditionalProperties() *T
}

// GetChildrenObjectSchemas returns all the children object schemas of a schema,
// only from the first level and without AnyOf, AllOf and OneOf.
func GetChildrenObjectSchemas[T TemplateSchema[T]](s T) []*T {
	allSchemas := utils.MapToList(s.GetProperties())

	if items := s.GetItems(); items != nil {
		allSchemas = append(allSchemas, items)
	}

	if additional := s.GetAdditionalProperties(); additional != nil {
		allSchemas = append(allSchemas, additional)
	}

	// Only keep object schemas (and arrays of objects).
	filteredSchemas := make([]*T, 0, len(allSchemas))
	for _, child := range allSchemas {
		schema := *child // concrete T exposes the TemplateSchema methods
		switch {
		case schema.GetType() == schemaTypeObject:
			filteredSchemas = append(filteredSchemas, child)
		case schema.GetType() == schemaTypeArray:
			if items := schema.GetItems(); items != nil && (*items).GetType() == schemaTypeObject {
				filteredSchemas = append(filteredSchemas, items)
			}
		}
	}

	return filteredSchemas
}

// referenceToSlicePath converts a reference to a slice where each element is a
// step of the path.
func referenceToSlicePath(ref string) []string {
	ref = strings.ReplaceAll(ref, ".", "/")
	ref = strings.ReplaceAll(ref, "#", "")
	return strings.Split(ref, "/")[1:]
}

// ReferenceToStructAttributePath converts a reference to a struct attribute path
// in the form of "a.b.c" where a, b and c are struct attributes in the form of
// golang conventional type names.
func ReferenceToStructAttributePath(ref string) string {
	path := referenceToSlicePath(ref)

	for k, v := range path {
		// If this is concerning the header, then it will be named "headers"
		if v == messageFieldHeader {
			v = "headers"
		}

		path[k] = template.Namify(v)
	}

	return strings.Join(path, ".")
}

// IsRequired checks if a field is required in an asyncapi schema.
func IsRequired[T TemplateSchema[T]](schema T, field string) bool {
	return schema.IsFieldRequired(field)
}

// forcePointers, when true, makes every non-array struct field be generated as a
// pointer. It is set explicitly on every generation run (see SetForcePointers),
// so the setting does not leak between successive generations sharing the same
// process.
var forcePointers bool

// SetForcePointers configures whether all struct fields (except arrays) should be
// generated as pointers.
func SetForcePointers(force bool) {
	forcePointers = force
}

// IsFieldPointer reports whether a struct field should be generated as a pointer.
func IsFieldPointer[T TemplateSchema[T]](parent T, field string, schema T) bool {
	if forcePointers {
		return schema.GetType() != schemaTypeArray
	}

	return !(IsRequired(parent, field) || schema.GetIsRequired()) && schema.GetType() != schemaTypeArray
}
