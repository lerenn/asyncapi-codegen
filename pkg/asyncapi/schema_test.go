package asyncapi

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalType(t *testing.T) {
	cases := []struct {
		in       string
		typ      string
		nullable bool
		err      bool
	}{
		{in: ``},
		{in: `"string"`, typ: "string"},
		{in: `["string","null"]`, typ: "string", nullable: true},
		{in: `["null","integer"]`, typ: "integer", nullable: true},
		{in: `["string"]`, typ: "string"},
		{in: `["string","integer"]`, err: true},
		{in: `42`, err: true},
	}

	for _, c := range cases {
		typ, nullable, err := UnmarshalType(json.RawMessage(c.in))
		if c.err {
			assert.Error(t, err, c.in)
			continue
		}
		assert.NoError(t, err, c.in)
		assert.Equal(t, c.typ, typ, c.in)
		assert.Equal(t, c.nullable, nullable, c.in)
	}
}
