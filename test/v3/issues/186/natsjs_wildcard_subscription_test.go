//go:generate go run ../../../../cmd/asyncapi-codegen -p issue186 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue186

import (
	"context"
	"sync"
	"testing"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/natsjetstream"
	testutil "github.com/lerenn/asyncapi-codegen/pkg/utils/test"
	natsio "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestWildcardSubscription(t *testing.T) {
	name := "issue186"

	brokerAddrParams := testutil.BrokerAddressParams{
		Schema:         "nats",
		DockerizedAddr: "nats-jetstream",
		DockerizedPort: "4222",
		LocalPort:      "4225",
	}

	nats, err := natsio.Connect(testutil.BrokerAddress(brokerAddrParams))
	require.NoError(t, err)

	natsJSBroker, err := natsjetstream.NewController(
		testutil.BrokerAddress(brokerAddrParams),
		natsjetstream.WithStreamConfig(jetstream.StreamConfig{
			Name:     name,
			Subjects: ChannelsPaths,
		}),
		natsjetstream.WithConsumerConfig(jetstream.ConsumerConfig{
			Name:    name,
			Durable: name, // make it durable so msgs aren't reprocessed on rerun
		}),
	)

	assert.NoError(t, err, "new controller should not return error")
	defer natsJSBroker.Close()
	suite.Run(t, newSuite(natsJSBroker, nats))
}

type Suite struct {
	broker extensions.BrokerController
	nats   *natsio.Conn
	app    *AppController
	suite.Suite
}

func newSuite(broker extensions.BrokerController, nats *natsio.Conn) *Suite {
	return &Suite{
		broker: broker,
		nats:   nats,
	}
}

func (suite *Suite) SetupTest() {
	app, err := NewAppController(suite.broker)
	suite.Require().NoError(err)
	suite.app = app
}

func (suite *Suite) TearDownTest() {
	suite.app.Close(context.Background())
}

func (suite *Suite) TestIssue186Star() {
	var wg sync.WaitGroup

	const payload = "some test msg"

	// Register consumer
	err := suite.app.SubscribeToStarRequestOperation(context.Background(),
		func(ctx context.Context, msg StarMessage) error {
			suite.Require().Equal(payload, msg.Payload)
			wg.Done()
			return nil
		})
	suite.Require().NoError(err)
	defer suite.app.UnsubscribeFromStarRequestOperation(context.Background())

	// Publish a message
	wg.Add(1)
	err = suite.nats.Publish("v2.issue186.star.dynamic.msg", []byte(payload))
	suite.Require().NoError(err)

	wg.Wait()
}

func (suite *Suite) TestIssue186Angle() {
	var wg sync.WaitGroup

	const payload = "some test msg"

	// Register consumer
	err := suite.app.SubscribeToAngleRequestOperation(context.Background(),
		func(ctx context.Context, msg AngleMessage) error {
			suite.Require().Equal(payload, msg.Payload)
			wg.Done()
			return nil
		})
	suite.Require().NoError(err)
	defer suite.app.UnsubscribeFromAngleRequestOperation(context.Background())

	// Publish a message
	wg.Add(1)
	err = suite.nats.Publish("v2.issue186.angle.dynamic.msg", []byte(payload))
	suite.Require().NoError(err)

	wg.Wait()
}
