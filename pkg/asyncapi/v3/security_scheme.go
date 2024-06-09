package asyncapiv3

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
