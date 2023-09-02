package broker

type Field string

func (f Field) String() string {
	return string(f)
}

const (
	// CorrelationIDField is the name of the field that will contain the correlation ID
	CorrelationIDField Field = "correlation_id"
)
