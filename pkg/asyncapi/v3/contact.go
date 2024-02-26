package asyncapiv3

// Contact is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#contactObject
type Contact struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Name  string `json:"name"`
	URL   string `json:"url"`
	Email string `json:"email"`

	// --- Non AsyncAPI fields -------------------------------------------------
}
