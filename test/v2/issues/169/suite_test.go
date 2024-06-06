//go:generate go run ../../../../cmd/asyncapi-codegen -p issue169 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue169

import (
	"context"
	"crypto/tls"
	"sync"
	"testing"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/kafka"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/natsjetstream"
	testutil "github.com/lerenn/asyncapi-codegen/pkg/utils/test"
	natsio "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	name := "issue169"

	// NATS Core with TLS and basic auth
	natsBrokerTLSBasicAuth, err := nats.NewController(
		testutil.BrokerAddress(testutil.BrokerAddressParams{
			Schema:         "nats",
			DockerizedAddr: "nats-tls-basic-auth",
			DockerizedPort: "4222",
			LocalPort:      "4224",
		}),
		nats.WithQueueGroup(name),
		nats.WithConnectionOpts(natsio.Secure(&tls.Config{InsecureSkipVerify: true}),
			natsio.UserInfo("user", "password"),
		),
	)
	assert.NoError(t, err, "new controller should not return error")
	defer natsBrokerTLSBasicAuth.Close()
	suite.Run(t, NewSuite(natsBrokerTLSBasicAuth))

	// NATS jetstream with TLS and basic auth
	natsJSBrokerTLSBasicAuth, err := natsjetstream.NewController(
		testutil.BrokerAddress(testutil.BrokerAddressParams{
			Schema:         "nats",
			DockerizedAddr: "nats-jetstream-tls-basic-auth",
			DockerizedPort: "4222",
			LocalPort:      "4227",
		}),
		natsjetstream.WithStreamConfig(jetstream.StreamConfig{
			Name:     name,
			Subjects: ChannelsPaths,
		}),
		natsjetstream.WithConsumerConfig(jetstream.ConsumerConfig{Name: name}),
		natsjetstream.WithConnectionOpts(natsio.Secure(&tls.Config{InsecureSkipVerify: true}),
			natsio.UserInfo("user", "password"),
		),
	)
	assert.NoError(t, err, "new controller should not return error")
	defer natsJSBrokerTLSBasicAuth.Close()
	suite.Run(t, NewSuite(natsJSBrokerTLSBasicAuth))

	// Kafka with TLS and basic auth
	sha512Mechanism, err := scram.Mechanism(scram.SHA512, "user", "password")
	assert.NoError(t, err, "new scram.SHA512 should not return a error")
	kafkaBrokerTLSBasicAuth, err := kafka.NewController(
		[]string{
			testutil.BrokerAddress(testutil.BrokerAddressParams{
				DockerizedAddr: "kafka-tls-basic-auth",
				DockerizedPort: "9092",
				LocalPort:      "9096",
			}),
		},
		kafka.WithGroupID(name),
		kafka.WithTLS(&tls.Config{InsecureSkipVerify: true}),
		kafka.WithSasl(sha512Mechanism),
	)
	assert.NoError(t, err, "new controller should not return error")
	suite.Run(t, NewSuite(kafkaBrokerTLSBasicAuth))
}

type Suite struct {
	broker extensions.BrokerController
	app    *AppController
	user   *UserController
	suite.Suite
}

func NewSuite(broker extensions.BrokerController) *Suite {
	return &Suite{
		broker: broker,
	}
}

func (suite *Suite) SetupTest() {
	// Create app
	app, err := NewAppController(suite.broker)
	suite.Require().NoError(err)
	suite.app = app

	// Create user
	user, err := NewUserController(suite.broker)
	suite.Require().NoError(err)
	suite.user = user
}

func (suite *Suite) TearDownTest() {
	suite.app.Close(context.Background())
	suite.user.Close(context.Background())
}

func (suite *Suite) TestIssue169() {
	var wg sync.WaitGroup

	// Test message
	sent := V2Issue169MsgMessage{
		Payload: "some test msg",
	}

	// Validate msg
	err := suite.app.SubscribeV2Issue169Msg(context.Background(),
		func(ctx context.Context, msg V2Issue169MsgMessage) error {
			suite.Require().Equal(sent.Payload, msg.Payload)
			wg.Done()
			return nil
		})
	suite.Require().NoError(err)
	defer suite.app.UnsubscribeV2Issue169Msg(context.Background())

	// Publish the message
	wg.Add(1)
	err = suite.user.PublishV2Issue169Msg(context.Background(), sent)
	suite.Require().NoError(err)

	wg.Wait()
}
