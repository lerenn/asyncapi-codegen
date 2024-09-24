package options

// GeneratorOptions are the options to activate some parts of code generation.
type GeneratorOptions struct {
	// Application should be true for application code generation to be generated
	Application bool
	// User should be true for user code generation to be generated
	User bool
	// Types should be true for type code (or common code) generation to be generated
	Types bool
}

// Options is the struct that gather configuration of codegen.
type Options struct {
	// OutputPath is the path to the generated code file
	OutputPath string

	// PackageName is the package name of the generated code
	PackageName string

	// Generate contains options regarding which golang code should be generated
	Generate GeneratorOptions

	// DisableFormatting states if the formatting should be disabled when
	// writing the generated code
	DisableFormatting bool

	// ConvertKeys defines a schema property keys conversion strategy.
	// Supported values: snake, camel, kebab, none
	ConvertKeys string

	// NamingScheme defines the naming scheme for generated golang structs
	// Supported values: camel, none
	NamingScheme string

	// IgnoreStringFormat states whether the properties' format (date, date-time) should impact the type in types
	IgnoreStringFormat bool

	// ForcePointers can be used to force all struct fields to be generated as pointers
	ForcePointers bool
}
