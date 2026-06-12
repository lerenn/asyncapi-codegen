//go:generate go run ../../../../cmd/asyncapi-codegen -p issue162 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue162

import (
	"context"
	"testing"

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

// pingSubscriber implements the generated AppSubscriber interface using the
// x-go-received-func overridden method name (issue #162). It would not compile
// if the callback were still named PingReceived.
type pingSubscriber struct{}

func (pingSubscriber) OnPing(_ context.Context, _ PingMessage) error { return nil }

var _ AppSubscriber = pingSubscriber{}

// TestOverriddenFunctionNamesExist verifies, at compile time, that the
// functions generated for the operation use the names provided through the
// x-go-*-func extensions instead of the default derived names (issue #162).
func (suite *Suite) TestOverriddenFunctionNamesExist() {
	// App side: subscribe / reply use the overridden names.
	var _ = (*AppController).ListenForPing
	var _ = (*AppController).AnswerPing

	// User side: send / request use the overridden names.
	var _ = (*UserController).PublishPing
	var _ = (*UserController).AskPing

	suite.Require().True(true)
}
