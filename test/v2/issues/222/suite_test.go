//go:generate go run ../../../../cmd/asyncapi-codegen -p issue222 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue222

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, NewSuite())
}

type Suite struct {
	suite.Suite
}

func NewSuite() *Suite {
	return &Suite{}
}

const val = `
{
	"DateProp": "2024-02-02",
	"DateTimeProp": "2024-06-05T13:45:30.0000Z"
}
`

func (suite *Suite) TestMarshalling() {
	var res TestSchema
	err := json.Unmarshal([]byte(val), &res)
	assert.NoError(suite.T(), err)
}
