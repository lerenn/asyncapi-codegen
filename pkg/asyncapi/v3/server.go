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

// generateMetadata generates metadata for the Server.
func (srv *Server) generateMetadata(name string) {
	// Prevent modification if nil
	if srv == nil {
		return
	}

	// Set name
	srv.Name = template.Namify(name)

	// Generate variables metadata
	for n, s := range srv.Variables {
		s.generateMetadata(n + "Variable")
	}

	// Generate security metadata
	srv.Security.generateMetadata(srv.Name + "Security")

	// Generate tags metadata
	for i, t := range srv.Tags {
		t.generateMetadata(fmt.Sprintf("%sTag%d", srv.Name, i))
	}

	// Generate external documentation metadata
	srv.ExternalDocs.generateMetadata(srv.Name + ExternalDocsNameSuffix)

	// Generate Bindings metadata
	srv.Bindings.generateMetadata(srv.Name + BindingsSuffix)
}

// setDependencies sets dependencies between the different elements of the Server.
func (srv *Server) setDependencies(spec Specification) error {
	// Prevent modification if nil
	if srv == nil {
		return nil
	}

	// Set references
	if err := srv.setReference(spec); err != nil {
		return err
	}

	// Set variables dependencies
	for _, s := range srv.Variables {
		if err := s.setDependencies(spec); err != nil {
			return err
		}
	}

	// Set security dependencies
	if err := srv.Security.setDependencies(spec); err != nil {
		return err
	}

	// Set tags dependencies
	for _, t := range srv.Tags {
		if err := t.setDependencies(spec); err != nil {
			return err
		}
	}

	// Set external documentation dependencies
	if err := srv.ExternalDocs.setDependencies(spec); err != nil {
		return err
	}

	// Set Bindings dependencies
	if err := srv.Bindings.setDependencies(spec); err != nil {
		return err
	}

	return nil
}

func (srv *Server) setReference(spec Specification) error {
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
