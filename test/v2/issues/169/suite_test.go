//go:generate go run ../../../../cmd/asyncapi-codegen -p issue169 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue169

import (
	"context"
	"crypto/tls"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/kafka"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/natsjetstream"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"

	natsio "github.com/nats-io/nats.go"
)

func TestSuite(t *testing.T) {
	name := "issue169"

	// core nats with TLS and basic auth
	natsBrokerTLSBasicAuth, err := nats.NewController("nats://nats-tls-basic-auth:4222", nats.WithQueueGroup(name),
		nats.WithConnectionOpts(natsio.Secure(&tls.Config{InsecureSkipVerify: true}),
			natsio.UserInfo("user", "password"),
		),
	)
	assert.NoError(t, err, "new controller should not return error")
	defer natsBrokerTLSBasicAuth.Close()
	suite.Run(t, NewSuite(natsBrokerTLSBasicAuth))

	// nats jetstream with TLS and basic auth
	natsJSBrokerTLSBasicAuth, err := natsjetstream.NewController(
		"nats://nats-jetstream-tls-basic-auth:4222",
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

	// kafka with TLS and basic auth
	sha512Mechanism, err := scram.Mechanism(scram.SHA512, "user", "password")
	assert.NoError(t, err, "new scram.SHA512 should not return a error")
	kafkaBrokerTLSBasicAuth, err := kafka.NewController([]string{"kafka-tls-basic-auth:9092"},
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

	wg sync.WaitGroup
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

func (suite *Suite) TestIssue169App() {
	// Test message
	sent := Issue169MsgMessage{
		Payload: "some test msg",
	}

	// validate msg
	//nolint:contextcheck
	err := suite.app.SubscribeIssue169Msg(context.Background(), func(_ context.Context, msg Issue169MsgMessage) error {
		//nolint:contextcheck
		suite.app.UnsubscribeIssue169Msg(context.Background())
		suite.Require().Equal(sent, msg)
		suite.wg.Done()
		return nil
	})
	suite.Require().NoError(err)

	suite.wg.Add(1)

	// Publish the message
	err = suite.app.PublishIssue169Msg(context.Background(), sent)
	suite.Require().NoError(err)

	// Wait for errorhandler is called
	suite.wg.Wait()
}

func (suite *Suite) TestIssue169User() {
	// Test message
	sent := Issue169MsgMessage{
		Payload: "some test msg",
	}

	// validate message
	err := suite.user.SubscribeIssue169Msg(context.Background(), func(_ context.Context, msg Issue169MsgMessage) error {
		suite.user.UnsubscribeIssue169Msg(context.Background())
		suite.Require().Equal(sent, msg)
		suite.wg.Done()
		return nil
	})
	suite.Require().NoError(err)

	suite.wg.Add(1)

	// Publish the message
	err = suite.user.PublishIssue169Msg(context.Background(), sent)
	suite.Require().NoError(err)

	// Wait for errorhandler is called
	suite.wg.Wait()
}
