// tfsdk tags have been added to fields of certain structs to help
// with the conversion of the terraform config to the API request [refer function getApplicationRequest]

package applications

import corporateidps "github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/corporateIdps"

type Meta struct {
	Created          string `json:"created,omitempty"`
	CreatedBy        string `json:"createdBy,omitempty"`
	LastModified     string `json:"lastModified,omitempty"`
	ModifiedBy       string `json:"modifiedBy,omitempty"`
	ModifiedOnBehalf string `json:"modifiedOnBehalf,omitempty"`
	Location         string `json:"location,omitempty"`
	Type             string `json:"type,omitempty"`
	Version          string `json:"version,omitempty"`
}

type Theme struct {
	Type string `json:"type,omitempty"`
	// Basic
	Advanced string `json:"advanced,omitempty"`
}

type Logo struct {
	ResourceId  string  `json:"resourceId"`
	Version     int     `json:"version,omitempty"`
	AspectRatio float32 `json:"aspectRatio,omitempty"`
}

type Branding struct {
	DisplayName                  string `json:"displayName,omitempty"`
	ShowDisplayNameOnLogonScreen bool   `json:"showDisplayNameOnLogonScreen,omitempty"`
	RememberMeVisible            bool   `json:"rememberMeVisible,omitempty"`
	RememberMeChecked            bool   `json:"rememberMeChecked,omitempty"`
	RefreshParent                bool   `json:"refreshParent,omitempty"`
	TokenUrlEmbedCharacter       string `json:"tokenUrlEmbedCharacter,omitempty"`
	EmailTemplateSet             string `json:"emailTemplateSet,omitempty"`
	Theme                        Theme  `json:"theme,omitempty"`
	Logo                         Logo   `json:"logo,omitempty"`
}

type UserAttribute struct {
	UserAttributeName string `json:"userAttributeName"`
	IsRequired        bool   `json:"isRequired"`
}

type UserAccess struct {
	Type                    string        `json:"type,omitempty"`
	UserAttributesForAccess UserAttribute `json:"userAttributesForAccess"`
}

type AuthorizationScope string

type ApiCertificateData struct {
	Id                  string               `json:"id"`
	Dn                  string               `json:"dn"`
	Description         string               `json:"description"`
	ApiNames            []string             `json:"apiNames"`
	AllApisAccess       bool                 `json:"allApisAccess"`
	AuthorizationScopes []AuthorizationScope `json:"authorizationScopes"`
	Base64Certificate   string               `json:"base64Certificate"`
}

type JwtClientAuthCredential struct {
	Id                  string               `json:"id"`
	Subject             string               `json:"subject"`
	Description         string               `json:"description"`
	IdentityProviderId  string               `json:"identityProviderId"`
	AuthorizationScopes []AuthorizationScope `json:"authorizationScopes"`
	ApiNames            []string             `json:"apiNames"`
	AllApisAccess       bool                 `json:"allApisAccess"`
}

type AssertionAttribute struct {
	AssertionAttributeName string `json:"assertionAttributeName,omitempty" tfsdk:"attribute_name"`
	UserAttributeName      string `json:"userAttributeName,omitempty" tfsdk:"attribute_value"`
	Inherited              bool   `json:"inherited" tfsdk:"inherited"`
}

type AdvancedAssertionAttribute struct {
	AttributeName  string `json:"attributeName,omitempty" tfsdk:"attribute_name"`
	AttributeValue string `json:"attributeValue,omitempty" tfsdk:"attribute_value"`
	Inherited      bool   `json:"inherited" tfsdk:"inherited"`
}

type DisabledInheritedProperties struct {
	AssertionAttributes         []AssertionAttribute         `json:"assertionAttributes,omitempty"`
	AdvancedAssertionAttributes []AdvancedAssertionAttribute `json:"advancedAssertionAttributes,omitempty"`
}

type AuthenicationRule struct {
	UserType           string `json:"userType,omitempty" tfsdk:"user_type"`
	UserGroup          string `json:"userGroup,omitempty" tfsdk:"user_group"`
	UserEmailDomain    string `json:"userEmailDomain,omitempty" tfsdk:"user_email_domain"`
	IdentityProviderId string `json:"identityProviderId,omitempty" tfsdk:"identity_provider_id"`
	IpNetworkRange     string `json:"ipNetworkRange,omitempty" tfsdk:"ip_network_range"`
}

type Saml2AcsEndpoint struct {
	BindingName string `json:"bindingName" tfsdk:"binding_name"`
	Location    string `json:"location" tfsdk:"location"`
	Index       int32  `json:"index" tfsdk:"index"`
	IsDefault   bool   `json:"isDefault,omitempty" tfsdk:"default"`
}

type Saml2SLOEndpoint struct {
	BindingName      string `json:"bindingName" tfsdk:"binding_name"`
	Location         string `json:"location" tfsdk:"location"`
	ResponseLocation string `json:"responseLocation,omitempty" tfsdk:"response_location"`
}

type EncryptionCertificateData struct {
	Dn                string `json:"dn,omitempty" tfsdk:"dn"`
	Base64Certificate string `json:"base64Certificate" tfsdk:"base64_certificate"`
	ValidFrom         string `json:"validFrom,omitempty" tfsdk:"valid_from"`
	ValidTo           string `json:"validTo,omitempty" tfsdk:"valid_to"`
}

type ProxyAuthnRequest struct {
	AuthenticationContext string `json:"authenticationContext,omitempty" tfsdk:"auth_context"`
	IssuerNameSuffix      string `json:"issuerNameSuffix,omitempty" tfsdk:"issuer_name_suffix"`
}

type SamlConfiguration struct {
	SamlMetadataUrl           string                                 `json:"samlMetadataUrl,omitempty" tfsdk:"saml_metadata_url"`
	DefaultNameIdFormat       string                                 `json:"defaultNameIdFormat,omitempty" tfsdk:"default_name_id_format"`
	AcsEndpoints              []Saml2AcsEndpoint                     `json:"acsEndpoints,omitempty" tfsdk:"acs_endpoints"`
	SloEndpoints              []Saml2SLOEndpoint                     `json:"sloEndpoints,omitempty" tfsdk:"slo_endpoints"`
	SignSLOMessages           bool                                   `json:"signSLOMessages" tfsdk:"sign_slo_messages"`
	RequireSignedSLOMessages  bool                                   `json:"requireSignedSLOMessages" tfsdk:"require_signed_slo_messages"`
	RequireSignedAuthnRequest bool                                   `json:"requireSignedAuthnRequest" tfsdk:"require_signed_auth_requests"`
	SignAssertions            bool                                   `json:"signAssertions" tfsdk:"sign_assertions"`
	SignAuthnResponses        bool                                   `json:"signAuthnResponses" tfsdk:"sign_auth_responses"`
	ResponseElementsToEncrypt string                                 `json:"responseElementsToEncrypt,omitempty" tfsdk:"response_elements_to_encrypt"`
	CertificatesForSigning    []corporateidps.SigningCertificateData `json:"certificatesForSigning,omitempty" tfsdk:"signing_certificates"`
	CertificateForEncryption  *EncryptionCertificateData             `json:"certificateForEncryption,omitempty" tfsdk:"encryption_certificate"`
	DigestAlgorithm           string                                 `json:"digestAlgorithm,omitempty" tfsdk:"digest_algorithm"`
	// ProxyAuthnRequest *ProxyAuthnRequest `json:"proxyAuthnRequest,omitempty" tfsdk:"proxy_auth_request"`
}

type ConsumedService struct {
	ServiceInstanceId string `json:"serviceInstanceId,omitempty"`
	AppId             string `json:"appId"`
	ClientId          string `json:"clientId"`
	PlanName          string `json:"planName,omitempty"`
	Inherit           string `json:"inherit"`
}

type ProvidedApi struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type ConsumedApi struct {
	AppId    string `json:"appId"`
	ClientId string `json:"clientId,omitempty"`
	ApiName  string `json:"apiName"`
	Name     string `json:"name"`
}

type AuthenticationSchema struct {
	SsoType                       string                       `json:"ssoType,omitempty" tfsdk:"sso_type"`
	SubjectNameIdentifier         string                       `json:"subjectNameIdentifier,omitempty" tfsdk:"subject_name_identifier"`
	SubjectNameIdentifierFunction string                       `json:"subjectNameIdentifierFunction,omitempty" tfsdk:"subject_name_identifier_function"`
	AssertionAttributes           []AssertionAttribute         `json:"assertionAttributes" tfsdk:"assertion_attributes"`
	AdvancedAssertionAttributes   []AdvancedAssertionAttribute `json:"advancedAssertionAttributes,omitempty" tfsdk:"advanced_assertion_attributes"`
	DefaultAuthenticatingIdpId    string                       `json:"defaultAuthenticatingIdpId,omitempty" tfsdk:"default_authenticating_idp"`
	ConditionalAuthentication     []AuthenicationRule          `json:"conditionalAuthentication,omitempty" tfsdk:"conditional_authentication"`
	Saml2Configuration            *SamlConfiguration           `json:"saml2Configuration,omitempty" tfsdk:""`
	// RiskBasedAuthentication       RBAConfiguration            `json:"riskBasedAuthentication"`
	// HomeUrl								string 							`json:"homeUrl"`
	// FallbackSubjectNameIdentifier		string 							`json:"fallbackSubjectNameIdentifier,omitempty"`
	// RememberMeExpirationTimeInMonths	string 							`json:"rememberMeExpirationTimeInMonths,omitempty"`
	// PasswordPolicy						string 							`json:"passwordPolicy"`
	// UserAccess							UserAccess 						`json:"userAccess,omitempty"`
	// CompanyId							string 							`json:"companyId"`
	// ClientId							string 							`json:"clientId"`
	// ApiCertificates						[]ApiCertificateData			`json:"apiCertificates"`
	// JwtClientAuthCredentials			[]JwtClientAuthCredential		`json:"jwtClientAuthCredentials"`
	// DisabledInheritedProperties			DisabledInheritedProperties		`json:"disabledInheritedProperties,omitempty"`
	// SocialSignOn						bool 							`json:"socialSignOn,omitempty"`
	// SpnegoEnabled						bool 							`json:"spnegoEnabled,omitempty"`
	// BiometricAuthenticationEnabled		bool 							`json:"biometricAuthenticationEnabled,omitempty"`
	// ForceAuthentication					bool 							`json:"forceAuthentication,omitempty"`
	// ConcurrentAccess					[]string						`json:"concurrentAccess,omitempty"`
	// TrustAllCorporateIdentityProviders	bool 							`json:"trustAllCorporateIdentityProviders,omitempty"`
	// AllowSciUsers						bool 							`json:"allowIaUsers,omitempty"`
	// ConsumedServices					[]ConsumedService				`json:"consumedServices,omitempty"`
	// ProvidedApis						[]ProvidedApi 					`json:"providedApis,omitempty"`
	// ConsumedApis						[]ConsumedApi 					`json:"consumedApis,omitempty"`
	// smsVerificationConfig
	// captchaConfig
	// openIdConnectConfiguration
	// sapManagedAttributes
	// idpCertificateSerialNumber
	// restApiAuthentication
}

type Application struct {
	Id                   string                `json:"id"`
	Name                 string                `json:"name"`
	Description          string                `json:"description,omitempty"`
	ParentApplicationId  string                `json:"parentApplicationId,omitempty"`
	MultiTenantApp       bool                  `json:"multiTenantApp,omitempty"` //only for SAP internal use
	Schemas              []string              `json:"schemas,omitempty"`
	AuthenticationSchema *AuthenticationSchema `json:"urn:sap:identity:application:schemas:extension:sci:1.0:Authentication"`
	// GlobalAccount        string               `json:"globalAccount,omitempty"`
	// Meta 					Meta 					`json:"meta,omitempty"`
	// PrivacyPolicy 			string 					`json:"privacyPolicy,omitempty"`
	// TermsOfUse 				string 					`json:"termsOfUse,omitempty"`
	// Branding 				Branding 				`json:"branding,omitempty"`
}

type ApplicationsResponse struct {
	TotalResults int           `json:"totalResults,omitempty"`
	ItemsPerPage int           `json:"itemsPerPage,omitempty"`
	NextCursor   string        `json:"nextCursor,omitempty"`
	Applications []Application `json:"applications,omitempty"`
}
