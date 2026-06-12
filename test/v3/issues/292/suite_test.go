//go:generate go run ../../../../cmd/asyncapi-codegen -p issue292 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue292

import (
	"context"
	"testing"
	"time"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/utils"
	"github.com/stretchr/testify/suite"
)

// silentBroker is a broker that never delivers any message. It lets us exercise
// the request-reply wait loop without a real broker: a reply never arrives, so
// the only way out is the request context being cancelled.
type silentBroker struct{}

func (silentBroker) Publish(_ context.Context, _ string, _ extensions.BrokerMessage) error {
	return nil
}

func (silentBroker) Subscribe(_ context.Context, _ string) (extensions.BrokerChannelSubscription, error) {
	sub := extensions.NewBrokerChannelSubscription(
		make(chan extensions.AcknowledgeableBrokerMessage, 1),
		make(chan any, 1),
	)
	sub.WaitForCancellationAsync(func() {})
	return sub, nil
}

func TestSuite(t *testing.T) {
	suite.Run(t, NewSuite())
}

type Suite struct {
	suite.Suite
}

func NewSuite() *Suite {
	return &Suite{}
}

// TestRequestReturnsOnContextTimeout ensures a request-reply call returns when
// its context times out instead of looping forever (issue #292). Before the
// fix, a context error from the wait loop was logged and ignored, so the loop
// spun endlessly and this test would hang.
func (suite *Suite) TestRequestReturnsOnContextTimeout() {
	user, err := NewUserController(silentBroker{})
	suite.Require().NoError(err)
	defer user.Close(context.Background())

	var msg PingMessage
	msg.Payload.Event = utils.ToPointer("ping")

	// The request must return shortly after the context deadline, not hang.
	done := make(chan error, 1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		_, reqErr := user.RequestToPingOperation(ctx, msg)
		done <- reqErr
	}()

	select {
	case reqErr := <-done:
		suite.Require().Error(reqErr)
		suite.Require().ErrorIs(reqErr, extensions.ErrContextCanceled)
	case <-time.After(5 * time.Second):
		suite.Fail("RequestToPingOperation did not return after context timeout (infinite loop)")
	}
}
