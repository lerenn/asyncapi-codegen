package asyncapiv3

import (
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/utils"
)

// Server is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#serverObject
type Server struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Host            string                     `json:"host"`
	Protocol        string                     `json:"protocol"`
	ProtocolVersion string                     `json:"protocolVersion"`
	PathName        string                     `json:"pathname"`
	Description     string                     `json:"description"`
	Title           string                     `json:"title"`
	Summary         string                     `json:"summary"`
	Variables       map[string]*ServerVariable `json:"variables"`
	Security        *SecurityScheme            `json:"security"`
	Tags            []*Tag                     `json:"tags"`
	ExternalDocs    *ExternalDocumentation     `json:"externalDocs"`
	Bindings        *ServerBindings            `json:"bindings"`
	Reference       string                     `json:"$ref"`

	// --- Non AsyncAPI fields -------------------------------------------------

	Name        string  `json:"-"`
	ReferenceTo *Server `json:"-"`
}

// Process processes the Server to make it ready for code generation.
func (srv *Server) Process(name string, spec Specification) {
	// Prevent modification if nil
	if srv == nil {
		return
	}

	// Set name
	srv.Name = utils.UpperFirstLetter(name)

	// Add pointer to reference if there is one
	if srv.Reference != "" {
		srv.ReferenceTo = spec.ReferenceServer(srv.Reference)
	}

	// Process variables
	for n, s := range srv.Variables {
		s.Process(n+"Variable", spec)
	}

	// Process security
	srv.Security.Process(srv.Name+"Security", spec)

	// Process tags
	for i, t := range srv.Tags {
		t.Process(fmt.Sprintf("%sTag%d", srv.Name, i), spec)
	}

	// Process external documentation
	srv.ExternalDocs.Process(srv.Name+ExternalDocsNameSuffix, spec)

	// Process Bindings
	srv.Bindings.Process(srv.Name+BindingsSuffix, spec)
}
