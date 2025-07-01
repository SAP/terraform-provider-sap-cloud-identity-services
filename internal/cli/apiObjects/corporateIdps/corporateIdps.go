package corporateidps

type SAML2SSOEndpoint struct {
	BindingName string `json:"bindingName" tfsdk:"binding_name"`
	IsDefault   bool   `json:"isDefault,omitempty" tfsdk:"default"`
	Location    string `json:"location" tfsdk:"location"`
}

type SAML2SLOEndpoint struct {
	BindingName      string `json:"bindingName" tfsdk:"binding_name"`
	IsDefault        bool   `json:"isDefault,omitempty" tfsdk:"default"`
	Location         string `json:"location" tfsdk:"location"`
	ResponseLocation string `json:"responseLocation,omitempty" tfsdk:"response_location"`
}

type SigningCertificateData struct {
	Base64Certificate string `json:"base64Certificate" tfsdk:"base64_certificate"`
	Dn                string `json:"dn,omitempty" tfsdk:"dn"`
	IsDefault         bool   `json:"isDefault" tfsdk:"default"`
	ValidFrom         string `json:"validFrom,omitempty" tfsdk:"valid_from"`
	ValidTo           string `json:"validTo,omitempty" tfsdk:"valid_to"`
}

type AssertionAttribute struct {
	Name  string `json:"name" tfsdk:"name"`
	Value string `json:"value" tfsdk:"value"`
}

type SAML2Configuration struct {
	// Alias                  string                   `json:"alias,omitempty" tfsdk:"alias"`
	AllowCreate            string                   `json:"allowCreate,omitempty" tfsdk:"allow_create"`
	AssertionAttributes    []AssertionAttribute     `json:"assertionAttributes,omitempty" tfsdk:"assertion_attributes"`
	CertificatesForSigning []SigningCertificateData `json:"certificatesForSigning,omitempty" tfsdk:"signing_certificates"`
	DefaultNameIdFormat    string                   `json:"defaultNameIdFormat,omitempty" tfsdk:"name_id_format"`
	DigestAlgorithm        string                   `json:"digestAlgorithm,omitempty" tfsdk:"digest_algorithm"`
	IncludeScoping         bool                     `json:"includeScoping,omitempty" tfsdk:"include_scoping"`
	SamlMetadataUrl        string                   `json:"samlMetadataUrl,omitempty" tfsdk:"saml_metadata_url"`
	SloEndpoints           []SAML2SLOEndpoint       `json:"sloEndpoints,omitempty" tfsdk:"slo_endpoints"`
	SsoEndpoints           []SAML2SSOEndpoint       `json:"ssoEndpoints,omitempty" tfsdk:"sso_endpoints"`
}

type OIDCAdditionalConfig struct {
	EnforceIssuerCheck       bool `json:"enforceIssuerCheck,omitempty" tfsdk:"enforce_issuer_check"`
	EnforceNonce             bool `json:"enforceNonce,omitempty" tfsdk:"enforce_nonce"`
	OmitIDTokenHintForLogout bool `json:"omitIDTokenHintForLogout,omitempty" tfsdk:"disable_logout_id_token_hint"`
}

type OIDCConfiguration struct {
	// Acrs                     []string              `json:"acrs,omitempty"`
	AdditionalConfig         *OIDCAdditionalConfig `json:"additionalConfig,omitempty" tfsdk:"additional_config"`
	AuthorizationEndpoint    string                `json:"authorizationEndpoint,omitempty" tfsdk:"authorization_endpoint"`
	ClientId                 string                `json:"clientId,omitempty" tfsdk:"client_id"`
	ClientSecret             string                `json:"clientSecret,omitempty" tfsdk:"client_secret"`
	DiscoveryUrl             string                `json:"discoveryUrl,omitempty" tfsdk:"discovery_url"`
	EndSessionEndpoint       string                `json:"endSessionEndpoint,omitempty" tfsdk:"logout_endpoint"`
	IsClientSecretConfigured bool                  `json:"isClientSecretConfigured,omitempty" tfsdk:"is_client_secret_configured"`
	Issuer                   string                `json:"issuer,omitempty" tfsdk:"issuer"`
	JwkSetPlain              string                `json:"jwkSetPlain,omitempty" tfsdk:"jwks"`
	JwksUri                  string                `json:"jwksUri,omitempty" tfsdk:"jwks_uri"`
	PkceEnabled              bool                  `json:"pkceEnabled,omitempty" tfsdk:"enable_pkce"`
	// RefreshDelay             int                   `json:"refreshDelay,omitempty"`
	Scopes                  []string `json:"scopes,omitempty" tfsdk:"scopes"`
	SubjectNameIdentifier   string   `json:"subjectNameIdentifier,omitempty" tfsdk:"subject_name_identifier"`
	TokenEndpoint           string   `json:"tokenEndpoint,omitempty" tfsdk:"token_endpoint"`
	TokenEndpointAuthMethod string   `json:"tokenEndpointAuthMethod,omitempty" tfsdk:"token_endpoint_auth_method"`
	UserInfoEndpoint        string   `json:"userInfoEndpoint,omitempty" tfsdk:"user_info_endpoint"`
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
	LoginHintType string `json:"loginHintType" tfsdk:"login_hint_type"`
	SendMethod    string `json:"sendMethod,omitempty" tfsdk:"send_method"`
}

type IdentityFederation struct {
	AllowLocalUsersOnly      bool     `json:"allowLocalUsersOnly,omitempty" tfsdk:"allow_local_users_only"`
	ApplyLocalIdPAuthnChecks bool     `json:"applyLocalIdPAuthnChecks,omitempty" tfsdk:"apply_local_idp_auth_and_checks"`
	RequiredGroups           []string `json:"requiredGroups,omitempty" tfsdk:"required_groups"`
	UseLocalUserStore        bool     `json:"useLocalUserStore,omitempty" tfsdk:"use_local_user_store"`
}

type IdentityProvider struct {
	AutomaticRedirect      bool                    `json:"automaticRedirect,omitempty"`
	CompanyId              string                  `json:"companyId,omitempty"`
	DisplayName            string                  `json:"displayName"`
	ForwardAllSsoRequests  bool                    `json:"forwardAllSsoRequests,omitempty"`
	Id                     string                  `json:"id,omitempty"`
	IdentityFederation     *IdentityFederation     `json:"identityFederation,omitempty"`
	LoginHintConfiguration *LoginHintConfiguration `json:"loginHintConfiguration,omitempty"`
	LogoutUrl              string                  `json:"logoutUrl,omitempty"`
	// Meta                   Meta                   `json:"meta,omitempty"`
	Name               string              `json:"name,omitempty"`
	OidcConfiguration  *OIDCConfiguration  `json:"oidcConfiguration"`
	Saml2Configuration *SAML2Configuration `json:"saml2Configuration,omitempty"`
	Type               string              `json:"type,omitempty"`
}

type IdentityProvidersResponse struct {
	IdentityProviders []IdentityProvider `json:"identityProviders,omitempty"`
	ItemsPerPage      int32              `json:"itemsPerPage,omitempty"`
	NextCursor        string             `json:"nextCursor,omitempty"`
	TotalResults      int32              `json:"totalResults,omitempty"`
}
