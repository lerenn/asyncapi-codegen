package loggers

// Level is the minimum severity a logger will output. Messages with a severity
// below the configured level are discarded.
type Level int

const (
	// LevelInfo outputs info, warning and error logs. This is the default.
	LevelInfo Level = iota
	// LevelWarning outputs warning and error logs only.
	LevelWarning
	// LevelError outputs error logs only.
	LevelError
)

// levelSetter is implemented by the loggers that support a configurable level.
type levelSetter interface {
	setLevel(level Level)
}

// Option is a functional option shared by the loggers of this package.
type Option func(levelSetter)

// WithLevel sets the minimum severity level that the logger will output.
// Messages below this level are silently discarded.
func WithLevel(level Level) Option {
	return func(ls levelSetter) {
		ls.setLevel(level)
	}
}
