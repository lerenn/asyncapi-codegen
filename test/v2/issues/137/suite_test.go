//go:generate go run ../../../../cmd/asyncapi-codegen -g types -p issue137 -i asyncapi.yaml -o ./asyncapi.gen.go
package issue137

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

func TestValidateOneOfExist(t *testing.T) {
	var schema AuditSchema
	field, ok := reflect.TypeOf(schema).FieldByName("Channel")
	assert.True(t, ok)

	validateTag, ok := field.Tag.Lookup("validate")
	assert.True(t, ok)

	var oneOfTag string
	for _, tag := range strings.Split(validateTag, ",") {
		if strings.HasPrefix(tag, "oneof=") {
			oneOfTag = tag
		}
	}

	assert.Equal(t, "oneof=API0 API1 API2 API3 API4", oneOfTag)
}
