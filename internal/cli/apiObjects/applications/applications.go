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

type ApplicationResponse struct{
	Id						string 		`json:"id"`
	Meta 					Meta 		`json:"meta,omitempty"`
	Name 					string 		`json:"name"`
	Description				string 		`json:"string,omitempty"`
	ParentApplicationId 	string 		`json:"parentApplicationId,omitempty"`
	MultiTenantApp 			bool 		`json:"multiTenantApp,omitempty"`	//only for SAP internal use
	PrivacyPolicy 			string 		`json:"privacyPolicy,omitempty"`
	TermsOfUse 				string 		`json:"termsOfUse,omitempty"`
	GlobalAccount 			string 		`json:"globalAccount,omitempty"`
	Schemas 				[]string 	`json:"schemas,omitempty"`
	Branding 				Branding 	`json:"branding,omitempty"`
}

type ApplicationsResponse struct{
	TotalResults	int 				`json:"totalResults,omitempty"`
	ItemsPerPage	int 				`json:"itemsPerPage,omitempty"`
	NextCursor		string 				`json:"nextCursor,omitempty"`
	Applications	ApplicationResponse	`json:"applications,omitempty"`
}

