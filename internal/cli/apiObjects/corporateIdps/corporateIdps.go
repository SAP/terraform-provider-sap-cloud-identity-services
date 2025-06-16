package corporateidps

type SAML2SSOEndpoint struct {
	BindingName string `json:"bindingName"`
	IsDefault   bool   `json:"isDefault,omitempty"`
	Location    string `json:"location"`
}

type SAML2SLOEndpoint struct {
	BindingName      string `json:"bindingName"`
	IsDefault        bool   `json:"isDefault,omitempty"`
	Location         string `json:"location"`
	ResponseLocation string `json:"responseLocation,omitempty"`
}

type SigningCertificateData struct {
	Base64Certificate string `json:"base64Certificate"`
	Dn                string `json:"dn,omitempty"`
	IsDefault         bool   `json:"isDefault"`
	ValidFrom         string `json:"validFrom,omitempty"`
	ValidTo           string `json:"validTo,omitempty"`
}

type AssertionAttribute struct {
	Name  string `json:"name" tfsdk:"name"`
	Value string `json:"value" tfsdk:"value"`
}

type SAML2Configuration struct {
	Alias                  string                   `json:"alias,omitempty"`
	AllowCreate            string                   `json:"allowCreate,omitempty"`
	AssertionAttributes    []AssertionAttribute     `json:"assertionAttributes,omitempty"`
	CertificatesForSigning []SigningCertificateData `json:"certificatesForSigning,omitempty"`
	DefaultNameIdFormat    string                   `json:"defaultNameIdFormat,omitempty"`
	DigestAlgorithm        string                   `json:"digestAlgorithm,omitempty"`
	IncludeScoping         bool                     `json:"includeScoping,omitempty"`
	SamlMetadataUrl        string                   `json:"samlMetadataUrl,omitempty"`
	SloEndpoints           []SAML2SLOEndpoint       `json:"sloEndpoints,omitempty"`
	SsoEndpoints           []SAML2SSOEndpoint       `json:"ssoEndpoints,omitempty"`
}

type OIDCAdditionalConfig struct {
	EnforceIssuerCheck       bool `json:"enforceIssuerCheck,omitempty"`
	EnforceNonce             bool `json:"enforceNonce,omitempty"`
	OmitIDTokenHintForLogout bool `json:"omitIDTokenHintForLogout,omitempty"`
}

type OIDCConfiguration struct {
	Acrs                     []string             `json:"acrs,omitempty"`
	AdditionalConfig         OIDCAdditionalConfig `json:"additionalConfig,omitempty"`
	AuthorizationEndpoint    string               `json:"authorizationEndpoint,omitempty"`
	ClientId                 string               `json:"clientId,omitempty"`
	ClientSecret             string               `json:"clientSecret,omitempty"`
	DiscoveryUrl             string               `json:"discoveryUrl,omitempty"`
	EndSessionEndpoint       string               `json:"endSessionEndpoint,omitempty"`
	IsClientSecretConfigured bool                 `json:"isClientSecretConfigured,omitempty"`
	Issuer                   string               `json:"issuer,omitempty"`
	JwkSetPlain              string               `json:"jwkSetPlain,omitempty"`
	JwksUri                  string               `json:"jwksUri,omitempty"`
	PkceEnabled              bool                 `json:"pkceEnabled,omitempty"`
	RefreshDelay             int                  `json:"refreshDelay,omitempty"`
	Scopes                   []string             `json:"scopes,omitempty"`
	SubjectNameIdentifier    string               `json:"subjectNameIdentifier,omitempty"`
	TokenEndpoint            string               `json:"tokenEndpoint,omitempty"`
	TokenEndpointAuthMethod  string               `json:"tokenEndpointAuthMethod,omitempty"`
	UserInfoEndpoint         string               `json:"userInfoEndpoint,omitempty"`
}

type Meta struct {
	Created          string `json:"created,omitempty"`
	CreatedBy        string `json:"createdBy,omitempty"`
	LastModified     string `json:"lastModified,omitempty"`
	Location         string `json:"location,omitempty"`
	ModifiedBy       string `json:"modifiedBy,omitempty"`
	ModifiedOnBehalf string `json:"modifiedOnBehalf,omitempty"`
}

type LoginHintConfiguration struct {
	LoginHintType string `json:"loginHintType"`
	SendMethod    string `json:"sendMethod,omitempty"`
}

type IdentityFederation struct {
	AllowLocalUsersOnly      bool     `json:"allowLocalUsersOnly,omitempty"`
	ApplyLocalIdPAuthnChecks bool     `json:"applyLocalIdPAuthnChecks,omitempty"`
	RequiredGroups           []string `json:"requiredGroups,omitempty"`
	UseLocalUserStore        bool     `json:"useLocalUserStore,omitempty"`
}

type IdentityProvider struct {
	AutomaticRedirect      bool                   `json:"automaticRedirect,omitempty"`
	CompanyId              string                 `json:"companyId,omitempty"`
	DisplayName            string                 `json:"displayName"`
	ForwardAllSsoRequests  bool                   `json:"forwardAllSsoRequests,omitempty"`
	Id                     string                 `json:"id,omitempty"`
	IdentityFederation     IdentityFederation     `json:"identityFederation,omitempty"`
	LoginHintConfiguration LoginHintConfiguration `json:"loginHintConfiguration,omitempty"`
	LogoutUrl              string                 `json:"logoutUrl,omitempty"`
	// Meta                   Meta                   `json:"meta,omitempty"`
	Name                   string                 `json:"name,omitempty"`
	OidcConfiguration      *OIDCConfiguration      `json:"oidcConfiguration"`
	Saml2Configuration     SAML2Configuration     `json:"saml2Configuration,omitempty"`
	Type                   string                 `json:"type,omitempty"`
}

type IdentityProvidersResponse struct {
	IdentityProviders []IdentityProvider `json:"identityProviders,omitempty"`
	ItemsPerPage      int32              `json:"itemsPerPage,omitempty"`
	NextCursor        string             `json:"nextCursor,omitempty"`
	TotalResults      int32              `json:"totalResults,omitempty"`
}
