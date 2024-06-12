package applications

type Meta struct{
	Created					string		`json:"created,omitempty"`
	CreatedBy				string		`json:"createdBy,omitempty"`
	LastModified			string		`json:"lastModified,omitempty"`
	ModifiedBy				string		`json:"modifiedBy,omitempty"`
	ModifiedOnBehalf		string		`json:"modifiedOnBehalf,omitempty"`
	Location				string		`json:"location,omitempty"`
	Type					string		`json:"type,omitempty"`
	Version 				string		`json:"version,omitempty"`
}

type Theme struct{
	Type 		string 		`json:"type,omitempty"`
	// Basic
	Advanced	string 		`json:"advanced,omitempty"`
}

type Logo struct{	
	ResourceId		string 		`json:"resourceId"`
	Version 		int 		`json:"version,omitempty"`
	AspectRatio  	float32		`json:"aspectRatio,omitempty"`
}


type Branding struct{
	DisplayName 					string 		`json:"displayName,omitempty"`
	ShowDisplayNameOnLogonScreen	bool 		`json:"showDisplayNameOnLogonScreen,omitempty"`
	RememberMeVisible 				bool 		`json:"rememberMeVisible,omitempty"`
	RememberMeChecked				bool 		`json:"rememberMeChecked,omitempty"`
	RefreshParent					bool 		`json:"refreshParent,omitempty"`
	TokenUrlEmbedCharacter			string 		`json:"tokenUrlEmbedCharacter,omitempty"`
	EmailTemplateSet				string 		`json:"emailTemplateSet,omitempty"`
	Theme 							Theme		`json:"theme,omitempty"`
	Logo 							Logo 		`json:"logo,omitempty"`
}

type UserAttribute struct{
	UserAttributeName		string 		`json:"userAttributeName"`
	IsRequired				bool 		`json:"isRequired"`
}

type UserAccess struct{
	Type						string 			`json:"type"`
	UserAttributesForAccess		UserAttribute	`json:"userAttributesForAccess"`
}

type AuthorizationScope string

type ApiCertificateData struct{
	Id 						string 					`json:"id"`
	Dn						string 					`json:"dn"`
	Description				string 					`json:"description"`
	ApiNames				[]string 				`json:"apiNames"`
	AllApisAccess			bool 					`json:"allApisAccess"`
	AuthorizationScopes		[]AuthorizationScope	`json:"authorizationScopes"`
	Base64Certificate		string 					`json:"base64Certificate"`
}

type JwtClientAuthCredential struct{
	Id 						string 					`json:"id"`
	Subject					string 					`json:"subject"`
	Description				string 					`json:"description"`
	IdentityProviderId		string 					`json:"identityProviderId"`
	AuthorizationScopes		[]AuthorizationScope	`json:"authorizationScopes"`
	ApiNames				[]string 				`json:"apiNames"`
	AllApisAccess			bool 					`json:"allApisAccess"`
}	

type AssertionAttribute struct{
	AssertionAttributeName		string 		`json:"assertionAttributeName"`
	UserAttributeName			string 		`json:"userAttributeName"`
	Inherited					bool 		`json:"inherited"`
}

type AdvancedAssertionAttribute struct{
	AttributeName		string 		`json:"attributeName,omitempty"`
	AttributeValue		string 		`json:"attributeValue,omitempty"`
	Inherited			bool 		`json:"inherited"`
}

type DisabledInheritedProperties struct {
	AssertionAttributes				[]AssertionAttribute			`json:"assertionAttributes,omitempty"`
	AdvancedAssertionAttributes		[]AdvancedAssertionAttribute	`json:"advancedAssertionAttributes,omitempty"`
}

type AuthenticationSchema struct{
	SsoType 							string 							`json:"ssoType"`
	HomeUrl								string 							`json:"homeUrl"`
	SubjectNameIdentifier				string 							`json:"subjectNameIdentifier"`
	FallbackSubjectNameIdentifier		string 							`json:"fallbackSubjectNameIdentifier"`
	SubjectNameIdentifierFunction		string 							`json:"subjectNameIdentifierFunction"`
	RememberMeExpirationTimeInMonths	string 							`json:"rememberMeExpirationTimeInMonths"`
	PasswordPolicy						string 							`json:"passwordPolicy"`
	UserAccess							UserAccess 						`json:"userAccess"`
	CompanyId							string 							`json:"companyId"`
	ClientId							string 							`json:"clientId"`
	ApiCertificates						[]ApiCertificateData			`json:"apiCertificates"`
	JwtClientAuthCredentials			[]JwtClientAuthCredential		`json:"jwtClientAuthCredentials"`
	AssertionAttributes					[]AssertionAttribute			`json:"assertionAttributes"`	
	AdvancedAssertionAttributes			[]AdvancedAssertionAttribute 	`json:"advancedAssertionAttributes"`
	DisabledInheritedProperties			DisabledInheritedProperties		`json:"disabledInheritedProperties,omitempty"`
	SocialSignOn						bool 							`json:"socialSignOn,omitempty"`
	
}

type ApplicationResponse struct{
	Id						string 		`json:"id"`
	Meta 					Meta 		`json:"meta,omitempty"`
	Name 					string 		`json:"name"`
	Description				string 		`json:"description,omitempty"`
	ParentApplicationId 	string 		`json:"parentApplicationId,omitempty"`
	MultiTenantApp 			bool 		`json:"multiTenantApp,omitempty"`	//only for SAP internal use
	PrivacyPolicy 			string 		`json:"privacyPolicy,omitempty"`
	TermsOfUse 				string 		`json:"termsOfUse,omitempty"`
	GlobalAccount 			string 		`json:"globalAccount,omitempty"`
	Schemas 				[]string 	`json:"schemas,omitempty"`
	Branding 				Branding 	`json:"branding,omitempty"`
}

type ApplicationsResponse struct{
	TotalResults	int 					`json:"totalResults,omitempty"`
	ItemsPerPage	int 					`json:"itemsPerPage,omitempty"`
	NextCursor		string 					`json:"nextCursor,omitempty"`
	Applications	[]ApplicationResponse	`json:"applications,omitempty"`
}

