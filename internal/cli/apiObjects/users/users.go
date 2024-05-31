package users

type Name struct{
	FamilyName 		string 		`json:"familyName"`
	GivenName 		string 		`json:"givenName"`
	Formatted 		string 		`json:"formatted,omitempty"`
	MiddleName		string 		`json:"middleName,omitempty"`
	HonoricPrefix 	string 		`json:"honoricPrefix,omitempty"`
	HonoricSuffix 	string 		`json:"honoricSuffix,omitempty"`
}

type Email struct{
	Value 		string 		`json:"value"`
	Type 		string 		`json:"type"`
	Display 	string 		`json:"display,omitempty"`
	Primary 	bool 		`json:"primary,omitempty"`
}

type PhoneNumber struct{
	Value 		string 		`json:"value"`
	Type 		string 		`json:"type"`
	Display 	string 		`json:"display,omitempty"`
	Primary 	bool 		`json:"primary,omitempty"`
}

type Photo struct{
	Value 		string 		`json:"value"`
	Type 		string 		`json:"type"`
	Display 	string 		`json:"display,omitempty"`
	Primary 	bool 		`json:"primary,omitempty"`
}

type Enititlement struct{
	Value 		string 		`json:"value"`
	Type 		string 		`json:"type"`
	Primary 	bool 		`json:"primary,omitempty"`
}

type Role struct{
	Value 		string 		`json:"value"`
	Type 		string 		`json:"type"`
	Primary 	bool 		`json:"primary,omitempty"`
}

type Meta struct{
	Description		string		`json:"description,omitempty"`
	Created			string		`json:"created,omitempty"`
	LastModified	string		`json:"lastModified,omitempty"`
	Location		string		`json:"location,omitempty"`
	ResourceType	string		`json:"resourceType,omitempty"`
	Version 		string		`json:"version,omitempty"`
}

type Address struct{
	Formatted 		string 		`json:"formatted,omitempty"`
	Primary 		bool 		`json:"primary,omitempty"`
	Country 		string 		`json:"country,omitempty"`
	Locality 		string 		`json:"locality,omitempty"`
	PostalCode 		string 		`json:"postalCode,omitempty"`
	Region 			string 		`json:"region,omitempty"`
	StreetAddress 	string 		`json:"streetAddress,omitempty"`
	Type 			string 		`json:"type"`
}

type Resource struct {
	Id 					string 			`json:"id"`
	ExternalId			string			`json:"externalId,omitempty"`
	Name 				Name 			`json:"name,omitempty"`
	Emails 				[]Email			`json:"emails"`
	Meta				Meta			`json:"meta,omitempty"`
	Schemas 			[]string 		`json:"schemas"`
	UserName 			string 			`json:"userName"`
	Password 			string 			`json:"password,omitempty"`
	DisplayName 		string 			`json:"displayName,omitempty"`
	NickName 			string 			`json:"nickName,omitempty"`
	ProfileUrl 			string 			`json:"profileUrl,omitempty"`
	Title 				string 			`json:"title,omitempty"`
	UserType 			string 			`json:"userType,omitempty"`
	PreferredLanguage	string			`json:"preferredLanguage,omitempty"`
	Locale 				string 			`json:"locale,omitempty"`
	TimeZone 			string 			`json:"timeZone,omitempty"`
	Active 				bool 			`json:"active,omitempty"`
	PhoneNumbers 		[]PhoneNumber	`json:"phoneNumbers,omitempty"`
	Photo 				[]Photo 		`json:"photos,omitempty"`
	Addresses 			[]Address 		`json:"addresses,omitempty"`
	Entitlements 		[]Enititlement 	`json:"entitlements,omitempty"`
	Roles 				[]Role 			`json:"roles,omitempty"`
	
}

type UserGet struct {
	Schemas 		[]string 	`json:"schemas,omitempty"`
	Resources 		[]Resource 	`json:"resources,omitempty"`
	TotalResult		int			`json:"totalResults,omitempty"`
	ItemsPerPage	int			`json:"itemsPerPage,omitempty"`
	StartIndex		int			`json:"startIndex,omitempty"`
	StartId 		string		`json:"startId,omitempty"`
	NextId 			string		`json:"nextId,omitempty"`
}

type UserReq struct {
	Schemas []string `json:"schemas"`
	UserName string `json:"userName"`
	Name Name `json:"name"`
	Emails []Email `json:"emails"`
}

type UserPost struct {
	Id string `json:"id"`
	Schemas []string `json:"schemas"`
	UserName string `json:"userName"`
	Name Name `json:"name"`
	Emails []Email `json:"emails"`
}