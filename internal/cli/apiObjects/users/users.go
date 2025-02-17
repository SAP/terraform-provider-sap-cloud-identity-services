package users

type CorporateGroup struct {
	Value string `json:"value"`
}

type SAPExtension struct {
	// LoginTime				string 		`json:"loginTime"`
	SourceSystem       int    `json:"sourceSystem"`
	SourceSystemId     string `json:"sourceSystemId"`
	ApplicationId      string `json:"applicationId"`
	EmailTemplateSetId string `json:"emailTemplateSetId"`
	SendMail           bool   `json:"sendMail"`
	TargetUrl          string `json:"targetUrl"`
	MailVerified       bool   `json:"mailVerified"`
	// UserUuid				string 		`json:"userUuid"`
	// UserUuidHistory			string 		`json:"userUuidHistory"` 		//read only
	UserId string `json:"userId"`
	// SapUserName				string 		`json:"sapUserName"`
	Status         string `json:"status,omitempty"`
	TotpEnabled    bool   `json:"totpEnabled"`
	WebAuthEnabled bool   `json:"webAuthEnabled"`
	// Industry				string 		`json:"industry"`
	// CompanyRelationship		string 		`json:"companyRelationship"`
	MfaEnabled bool `json:"mfaEnabled"`
	// contactPreferences
	// socialIdentities
	// passwordDetails
	// emails
	// addresses
	CorporateGroups []CorporateGroup `json:"corporateGroups"`
	// ValidFrom				string 		`json:"validFrom"`
	// ValidTo					string 		`json:"validTo"`
}

type Manager struct {
	DisplayName string `json:"displayName"`
	Value       string `json:"value"`
}

type EnterpriseUser struct {
	Division       string  `json:"division"`
	CostCenter     string  `json:"costCenter"`
	Organization   string  `json:"organization"`
	Department     string  `json:"department"`
	EmployeeNumber string  `json:"employeeNumber"`
	Manager        Manager `json:"manager"`
}

type Name struct {
	FamilyName    string `json:"familyName"`
	GivenName     string `json:"givenName"`
	Formatted     string `json:"formatted,omitempty"`
	MiddleName    string `json:"middleName,omitempty"`
	HonoricPrefix string `json:"honoricPrefix,omitempty"`
	HonoricSuffix string `json:"honoricSuffix,omitempty"`
}

type Email struct {
	Value   string `json:"value"`
	Type    string `json:"type"`
	Display string `json:"display,omitempty"`
	Primary bool   `json:"primary,omitempty"`
}

type PhoneNumber struct {
	Value   string `json:"value"`
	Type    string `json:"type"`
	Display string `json:"display,omitempty"`
	Primary bool   `json:"primary,omitempty"`
}

type Photo struct {
	Value   string `json:"value"`
	Type    string `json:"type"`
	Display string `json:"display,omitempty"`
	Primary bool   `json:"primary,omitempty"`
}

type Enititlement struct {
	Value   string `json:"value"`
	Type    string `json:"type"`
	Primary bool   `json:"primary,omitempty"`
}

type Role struct {
	Value   string `json:"value"`
	Type    string `json:"type"`
	Primary bool   `json:"primary,omitempty"`
}

type Meta struct {
	Description  string `json:"description,omitempty"`
	Created      string `json:"created,omitempty"`
	LastModified string `json:"lastModified,omitempty"`
	Location     string `json:"location,omitempty"`
	ResourceType string `json:"resourceType,omitempty"`
	Version      string `json:"version,omitempty"`
}

type Address struct {
	Formatted     string `json:"formatted,omitempty"`
	Primary       bool   `json:"primary,omitempty"`
	Country       string `json:"country,omitempty"`
	Locality      string `json:"locality,omitempty"`
	PostalCode    string `json:"postalCode,omitempty"`
	Region        string `json:"region,omitempty"`
	StreetAddress string `json:"streetAddress,omitempty"`
	Type          string `json:"type"`
}

type User struct {
	Id					string			`json:"id,omitempty"`
	ExternalId			string         	`json:"externalId,omitempty"`
	Meta             	Meta           	`json:"meta,omitempty"`
	Schemas           	[]string       	`json:"schemas"`
	UserName          	string         	`json:"userName"`
	Password         	string         	`json:"password,omitempty"`
	Name              	Name           	`json:"name,omitempty"`
	DisplayName       	string         	`json:"displayName,omitempty"`
	NickName          	string         	`json:"nickName,omitempty"`
	ProfileUrl        	string         	`json:"profileUrl,omitempty"`
	Title             	string         	`json:"title,omitempty"`
	UserType          	string         	`json:"userType,omitempty"`
	PreferredLanguage 	string        	`json:"preferredLanguage,omitempty"`
	Locale            	string         	`json:"locale,omitempty"`
	TimeZone          	string         	`json:"timeZone,omitempty"`
	Active            	bool           	`json:"active,omitempty"`
	Emails            	[]Email        	`json:"emails"`
	PhoneNumbers      	[]PhoneNumber  	`json:"phoneNumbers,omitempty"`
	Photo             	[]Photo        	`json:"photos,omitempty"`
	Addresses         	[]Address      	`json:"addresses,omitempty"`
	Entitlements      	[]Enititlement 	`json:"entitlements,omitempty"`
	Roles             	[]Role         	`json:"roles,omitempty"`
	// EnterpriseUser    	EnterpriseUser 	`json:"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User,omitempty"`
	SAPExtension 	  	SAPExtension 	`json:"urn:ietf:params:scim:schemas:extension:sap:2.0:User"`
}

type UsersResponse struct {
	Schemas      []string `json:"schemas,omitempty"`
	Resources    []User   `json:"Resources,omitempty"`
	TotalResult  int      `json:"totalResults,omitempty"`
	ItemsPerPage int      `json:"itemsPerPage,omitempty"`
	StartIndex   int      `json:"startIndex,omitempty"`
	StartId      string   `json:"startId,omitempty"`
	NextId       string   `json:"nextId,omitempty"`
}
