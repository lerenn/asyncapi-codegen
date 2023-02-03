package asyncapi

type Any struct {
	Type        string         `json:"type"`
	Description string         `json:"description"`
	Format      string         `json:"format"`
	Properties  map[string]Any `json:"properties"`
	Items       Items          `json:"items"`
	Reference   string         `json:"$ref"`

	// Non AsyncAPI fields
	Name string `json:"-"`
}
