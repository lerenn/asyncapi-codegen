package asyncapi_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/kafka"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
)

// BrokerControllers returns a list of BrokerController to test based on the
// docker-compose file of the project.
func BrokerControllers(t *testing.T) (brokers []extensions.BrokerController, cleanup func()) {
	t.Helper() // Set this function as a helper

	// Initialize returned values
	brokers = make([]extensions.BrokerController, 0)

	// Set a specific queueGroupeID to avoid collision between tests
	queueGroupID := fmt.Sprintf("test-%d", time.Now().UnixNano())

	// Add NATS broker
	nb := nats.NewController("nats://localhost:4222", nats.WithQueueGroup(queueGroupID))
	brokers = append(brokers, nb)

	// Add kafka broker
	kb := kafka.NewController([]string{"localhost:9094"}, kafka.WithGroupID(queueGroupID))
	brokers = append(brokers, kb)

	// Return brokers with their cleanup functions
	return brokers, func() {
		// Clean up NATS
		nb.Close()
	}
}
