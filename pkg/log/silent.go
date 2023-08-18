package log

// Silent is a logger that does not log anything
type Silent struct {
}

// Info logs information based on a message and key-value elements
func (s Silent) Info(_ Context, _ string, _ ...AdditionalInfo) {}

// Error logs error based on a message and key-value elements
func (s Silent) Error(_ Context, _ string, _ ...AdditionalInfo) {}
