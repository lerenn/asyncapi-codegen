package asyncapiv3

import (
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/utils/template"
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
func (srv *Server) Process(name string, spec Specification) error {
	// Prevent modification if nil
	if srv == nil {
		return nil
	}

	// Set name
	srv.Name = template.Namify(name)

	// Process references
	if err := srv.processReference(spec); err != nil {
		return err
	}

	// Process variables
	for n, s := range srv.Variables {
		if err := s.Process(n+"Variable", spec); err != nil {
			return err
		}
	}

	// Process security
	if err := srv.Security.Process(srv.Name+"Security", spec); err != nil {
		return err
	}

	// Process tags
	for i, t := range srv.Tags {
		if err := t.Process(fmt.Sprintf("%sTag%d", srv.Name, i), spec); err != nil {
			return err
		}
	}

	// Process external documentation
	if err := srv.ExternalDocs.Process(srv.Name+ExternalDocsNameSuffix, spec); err != nil {
		return err
	}

	// Process Bindings
	if err := srv.Bindings.Process(srv.Name+BindingsSuffix, spec); err != nil {
		return err
	}

	return nil
}

func (srv *Server) processReference(spec Specification) error {
	// check reference exists
	if srv.Reference == "" {
		return nil
	}

	// Get reference
	refTo, err := spec.ReferenceServer(srv.Reference)
	if err != nil {
		return err
	}

	// Set reference
	srv.ReferenceTo = refTo

	return nil
}
