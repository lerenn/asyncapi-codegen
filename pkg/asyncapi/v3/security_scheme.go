package asyncapiv3

import "github.com/lerenn/asyncapi-codegen/pkg/utils"

// OAuthFlow is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#oauthFlowObject
type OAuthFlow struct {
	// --- AsyncAPI fields -----------------------------------------------------

	AuthorizationURL string            `json:"authorizationUrl"`
	TokenURL         string            `json:"tokenUrl"`
	RefreshURL       string            `json:"refreshUrl"`
	AvailableScopes  map[string]string `json:"availableScopes"`

	// --- Non AsyncAPI fields -------------------------------------------------
}

// OAuthFlows is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#oauthFlowsObject
type OAuthFlows struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Implicit          OAuthFlow `json:"implicit"`
	Password          OAuthFlow `json:"password"`
	ClientCredentials OAuthFlow `json:"clientCredential"`
	AuthorizationCode OAuthFlow `json:"authorizationCode"`

	// --- Non AsyncAPI fields -------------------------------------------------
}

// SecurityScheme is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#securitySchemeObject
type SecurityScheme struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Type             string     `json:"type"`
	Description      string     `json:"description"`
	Name             string     `json:"name"`
	In               string     `json:"in"`
	Scheme           string     `json:"scheme"`
	BearerFormat     string     `json:"bearerFormat"`
	Flows            OAuthFlows `json:"flows"`
	OpenIDConnectURL string     `json:"openIdConnectUrl"`
	Scopes           []string   `json:"scopes"`
	Reference        string     `json:"$ref"`

	// --- Non AsyncAPI fields -------------------------------------------------

	ReferenceTo *SecurityScheme `json:"-"`
}

// Process processes the SecurityScheme to make it ready for code generation.
func (s *SecurityScheme) Process(name string, spec Specification) {
	// Prevent modification if nil
	if s == nil {
		return
	}

	// Set name
	s.Name = utils.UpperFirstLetter(name)

	// Add pointer to reference if there is one
	if s.Reference != "" {
		s.ReferenceTo = spec.ReferenceSecurity(s.Reference)
	}
}

// RemoveDuplicateSecuritySchemes removes the security schemes that have the same
// name, keeping the first occurrence.
func RemoveDuplicateSecuritySchemes(securities []*SecurityScheme) []*SecurityScheme {
	newList := make([]*SecurityScheme, 0)
	for _, s := range securities {
		present := false
		for _, ps := range newList {
			if ps.Name == s.Name {
				present = true
				break
			}
		}

		if !present {
			newList = append(newList, s)
		}
	}
	return newList
}
