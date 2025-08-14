package asyncapiv3

// OAuthFlow is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#oauthFlowObject
type OAuthFlow struct {
	// --- AsyncAPI fields -----------------------------------------------------

	AuthorizationURL string            `json:"authorizationUrl"`
	TokenURL         string            `json:"tokenUrl"`
	RefreshURL       string            `json:"refreshUrl,omitempty"`
	AvailableScopes  map[string]string `json:"availableScopes"`

	// --- Non AsyncAPI fields -------------------------------------------------
}

// OAuthFlows is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#oauthFlowsObject
type OAuthFlows struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Implicit          OAuthFlow `json:"implicit,omitzero"`
	Password          OAuthFlow `json:"password,omitzero"`
	ClientCredentials OAuthFlow `json:"clientCredential,omitzero"`
	AuthorizationCode OAuthFlow `json:"authorizationCode,omitzero"`

	// --- Non AsyncAPI fields -------------------------------------------------
}

// SecurityScheme is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#securitySchemeObject
type SecurityScheme struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Type             string     `json:"type"`
	Description      string     `json:"description,omitempty"`
	Name             string     `json:"name,omitempty"`
	In               string     `json:"in,omitempty"`
	Scheme           string     `json:"scheme,omitempty"`
	BearerFormat     string     `json:"bearerFormat,omitempty"`
	Flows            OAuthFlows `json:"flows,omitzero"`
	OpenIDConnectURL string     `json:"openIdConnectUrl,omitempty"`
	Scopes           []string   `json:"scopes,omitempty"`
	Reference        string     `json:"$ref,omitempty"`

	// --- Non AsyncAPI fields -------------------------------------------------

	ReferenceTo *SecurityScheme `json:"-"`
}

// SecuritySchemeType is type of a security scheme.
type SecuritySchemeType string

// Security scheme types supported by the AsyncAPI specification.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#securitySchemeObject
const (
	SecuritySchemeUserPassword         SecuritySchemeType = "userPassword"
	SecuritySchemeAPIKey               SecuritySchemeType = "apiKey"
	SecuritySchemeX509                 SecuritySchemeType = "X509"
	SecuritySchemeSymmetricEncryption  SecuritySchemeType = "symmetricEncryption"
	SecuritySchemeASymmetricEncryption SecuritySchemeType = "asymmetricEncryption"
	SecuritySchemeHTTPAPIKey           SecuritySchemeType = "httpApiKey"
	SecuritySchemeHTTP                 SecuritySchemeType = "http"
	SecuritySchemeOAuth2               SecuritySchemeType = "oauth2"
	SecuritySchemeOpenIDConnect        SecuritySchemeType = "openIdConnect"
	SecuritySchemePlain                SecuritySchemeType = "plain"
	SecuritySchemeScramSha256          SecuritySchemeType = "scramSha256"
	SecuritySchemeScramSha512          SecuritySchemeType = "scramSha512"
	SecuritySchemeGSSAPI               SecuritySchemeType = "gssapi"
)

// generateMetadata generates metadata for the SecurityScheme.
func (s *SecurityScheme) generateMetadata(parentName, name string, number *int) {
	// Prevent modification if nil
	if s == nil {
		return
	}

	// Set name
	s.Name = generateFullName(parentName, name, "Security_Scheme", number)
}

func (s *SecurityScheme) setDependencies(spec Specification) error {
	// Prevent modification if nil
	if s == nil {
		return nil
	}

	// Add pointer to reference if there is one
	if s.Reference != "" {
		refTo, err := spec.ReferenceSecurity(s.Reference)
		if err != nil {
			return err
		}
		s.ReferenceTo = refTo
	}

	return nil
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
