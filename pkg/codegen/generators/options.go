package generators

// Options are the options to activate some parts of code generation
type Options struct {
	// Application should be true for application code generation to be generated
	Application bool
	// User should be true for user code generation to be generated
	User bool
	// Types should be true for type code (or common code) generation to be generated
	Types bool
}
