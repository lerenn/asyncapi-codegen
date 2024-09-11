//go:generate go run ../../../../cmd/asyncapi-codegen --ignore-string-format -p issue255ignoredates -i ./asyncapi.yaml -o ./ignoredates/asyncapi.gen.go
//go:generate go run ../../../../cmd/asyncapi-codegen -p issue255default -i ./asyncapi.yaml -o ./default/asyncapi.gen.go

package issue255

import (
	"testing"
	"time"

	"cloud.google.com/go/civil"
	issue255default "github.com/lerenn/asyncapi-codegen/test/v2/issues/255/default"
	ignoredates "github.com/lerenn/asyncapi-codegen/test/v2/issues/255/ignoredates"
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

func (suite *Suite) TestDefault() {
	t := time.Now()
	_ = issue255default.TestMessagePayload{
		DateProp: &civil.Date{
			Year:  2024,
			Month: 12,
			Day:   12,
		},
		DateTimeProp: &t,
	}

	s := "toto"
	_ = ignoredates.TestMessagePayload{
		DateProp:     &s,
		DateTimeProp: &s,
	}
}
