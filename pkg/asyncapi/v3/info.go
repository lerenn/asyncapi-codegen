package asyncapiv3

import "fmt"

// Info is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#infoObject
type Info struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Title          string                 `json:"title"`
	Version        string                 `json:"version"`
	Description    string                 `json:"description"`
	TermsOfService string                 `json:"termsOfService"`
	Contact        *Contact               `json:"contact"`
	License        *License               `json:"license"`
	Tags           []*Tag                 `json:"tags"`
	ExternalDocs   *ExternalDocumentation `json:"externalDocs"`

	// --- Non AsyncAPI fields -------------------------------------------------
}

// Process processes the Info to make it ready for code generation.
func (info *Info) Process(spec Specification) error { // Prevent modification if nil
	if info == nil {
		return nil
	}

	// Process tags
	for i, t := range info.Tags {
		if err := t.Process(fmt.Sprintf("InfoTag%d", i), spec); err != nil {
			return err
		}
	}

	// Process external documentation
	if err := info.ExternalDocs.Process("InfoExternalDocs", spec); err != nil {
		return err
	}

	return nil
}
