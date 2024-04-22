package test

import (
	"fmt"
	"testing"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/kafka"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	testutil "github.com/lerenn/asyncapi-codegen/pkg/utils/test"
)

// BrokerAddressParams is the parameters for the BrokerAddress function.
type BrokerAddressParams struct {
	Schema string
	Port   string

	DockerizedAddr string
	DockerizedPort string

	LocalAddr string
	LocalPort string
}

// BrokerControllers returns a list of BrokerController to test based on the
// docker-compose file of the project.
func BrokerControllers(t *testing.T) ([]extensions.BrokerController, func()) {
	t.Helper() // Set this function as a helper

	// Set a specific queueGroupID to avoid collision between tests
	queueGroupID := fmt.Sprintf("test-%s", t.Name())
	fmt.Println(queueGroupID)

	// Add NATS broker
	natsController, err := nats.NewController(
		testutil.BrokerAddress(testutil.BrokerAddressParams{
			Schema:         "nats",
			DockerizedAddr: "nats",
			Port:           "4222",
		}),
		nats.WithQueueGroup(queueGroupID))
	if err != nil {
		panic(err)
	}

	// Add Kafka broker
	kafkaController, err := kafka.NewController(
		[]string{
			testutil.BrokerAddress(testutil.BrokerAddressParams{
				DockerizedAddr: "kafka",
				Port:           "9092",
			}),
		},
		kafka.WithGroupID(queueGroupID))
	if err != nil {
		panic(err)
	}

	// Return brokers with their cleanup functions
	return []extensions.BrokerController{
			natsController,
			kafkaController,
		}, func() {
			natsController.Close()
		}
}
