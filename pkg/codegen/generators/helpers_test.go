package generators

import (
	"testing"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
	"github.com/stretchr/testify/assert"
)

// TestGenerateValidateTagsEnum ensures enum validation directives are generated
// for non-string enums (integers, booleans) and not only for strings (issue
// #137). JSON numbers are unmarshaled as float64, which is what we test here.
func TestGenerateValidateTagsEnum(t *testing.T) {
	cases := map[string]struct {
		enum     []any
		expected string
	}{
		"string enum": {
			enum:     []any{"foo", "bar baz"},
			expected: "oneof='foo' 'bar baz'",
		},
		"integer enum": {
			enum:     []any{float64(1), float64(2), float64(3)},
			expected: "oneof=1 2 3",
		},
		"boolean enum": {
			enum:     []any{true, false},
			expected: "oneof=true false",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			tag := GenerateValidateTags[any](asyncapi.Validations[any]{
				Enum: tc.enum,
			}, false, "string")
			assert.Contains(t, tag, tc.expected)
		})
	}
}

// TestGenerateValidateTagsConst ensures const validation directives are
// generated for non-string consts too (issue #137).
func TestGenerateValidateTagsConst(t *testing.T) {
	tag := GenerateValidateTags[any](asyncapi.Validations[any]{
		Const: float64(42),
	}, false, "integer")
	assert.Contains(t, tag, "eq=42")
}
