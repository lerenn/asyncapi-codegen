package ci

const (
	// RedPandaImage is the image used for redpanda.
	RedPandaImage = "docker.redpanda.com/redpandadata/redpanda:v23.3.11"
	// KafkaImage is the image used for kafka.
	KafkaImage = "bitnami/kafka:3.5.1"
	// GolangImage is the image used for golang execution.
	GolangImage = "golang:1.21.4"
	// LinterImage is the image used for linter.
	LinterImage = "golangci/golangci-lint:v1.55"
	// NATSImage is the image used for NATS.
	NATSImage = "nats:2.10"
)
