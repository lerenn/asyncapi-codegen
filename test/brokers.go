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
func BrokerControllers(t *testing.T) map[string]extensions.BrokerController {
	t.Helper() // Set this function as a helper

	// Set a specific queueGroupeID to avoid collision between tests
	queueGroupID := fmt.Sprintf("test-%d", time.Now().UnixNano())

	return map[string]extensions.BrokerController{
		"NATS":  nats.NewController("nats://localhost:4222", nats.WithQueueGroup(queueGroupID)),
		"Kafka": kafka.NewController([]string{"localhost:9094"}, kafka.WithGroupID(queueGroupID)),
	}
}
