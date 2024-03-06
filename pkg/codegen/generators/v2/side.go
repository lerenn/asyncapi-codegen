package generatorv2

// Side represents the side of the code generation based on asyncapi documentation,
// i.e in front (user) or behind (application) asyncapi specification.
type Side string

const (
	// SideIsApplication is the application side based on asyncapi documentation,
	// i.e. the side that stand behind of the asyncapi specification.
	SideIsApplication Side = "app"
	// SideIsUser is the user side based on asyncapi documentation,
	// i.e. the side that use the asyncapi specification.
	SideIsUser Side = "user"
)
