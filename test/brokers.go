package asyncapi_test

import (
	"testing"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers"
	"github.com/nats-io/nats.go"

	"github.com/stretchr/testify/assert"
)

// BrokerControllers returns a list of BrokerController to test based on the
// docker-compose file of the project.
func BrokerControllers(t *testing.T) []extensions.BrokerController {
	// NATS
	nc, err := nats.Connect("nats://localhost:4222")
	assert.NoError(t, err)
	natsController := brokers.NewNATSController(nc)

	// Kafka
	kafkaController := brokers.NewKafkaController([]string{"127.0.0.1:9092"})

	return []extensions.BrokerController{
		natsController,
		kafkaController,
	}
}
