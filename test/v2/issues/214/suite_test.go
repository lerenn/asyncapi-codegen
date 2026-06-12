//go:generate go run ../../../../cmd/asyncapi-codegen -p issue214 -g types -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue214

import (
	"encoding/json"
	"testing"

	"github.com/lerenn/asyncapi-codegen/pkg/utils"
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

// TestGoNameOverride ensures the x-go-name extension overrides the generated
// type name (the schema is named `Event`, not `EventSchema`) and the struct
// field name (`EventID`, not `EventId`), while the JSON tag keeps the original
// spec key (issue #214).
func (suite *Suite) TestGoNameOverride() {
	// Compile-time proof of the overridden type and field names.
	value := Event{
		EventID: utils.ToPointer(int64(42)),
	}

	// The JSON wire format must still use the original property key.
	data, err := json.Marshal(value)
	suite.Require().NoError(err)
	suite.Require().JSONEq(`{"event_id":42}`, string(data))
}
