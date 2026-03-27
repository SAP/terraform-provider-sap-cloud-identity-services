package provider

import (
	"context"
	"fmt"
	"reflect"

	corporateidps "github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/corporateIdps"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/generic"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type sloEndpointData struct {
	BindingName      types.String `tfsdk:"binding_name"`
	Location         types.String `tfsdk:"location"`
	ResponseLocation types.String `tfsdk:"response_location"`
	Default          types.Bool   `tfsdk:"default"`
}

type signingCertificateData struct {
	Base64Certificate types.String `tfsdk:"base64_certificate"`
	Dn                types.String `tfsdk:"dn"`
	Default           types.Bool   `tfsdk:"default"`
	ValidFrom         types.String `tfsdk:"valid_from"`
	ValidTo           types.String `tfsdk:"valid_to"`
}

type saml2ConfigData struct {
	SamlMetadataUrl     types.String `tfsdk:"saml_metadata_url"    json:"samlMetadataUrl"`
	DigestAlgorithm     types.String `tfsdk:"digest_algorithm"     json:"digestAlgorithm"`
	IncludeScoping      types.Bool   `tfsdk:"include_scoping"      json:"includeScoping"`
	NameIdFormat        types.String `tfsdk:"name_id_format"       json:"defaultNameIdFormat"`
	AllowCreate         types.String `tfsdk:"allow_create"         json:"allowCreate"`
	AssertionAttributes types.List   `tfsdk:"assertion_attributes" json:"assertionAttributes"`
	SigningCertificates types.List   `tfsdk:"signing_certificates" json:"certificatesForSigning"`
	SsoEndpoints        types.List   `tfsdk:"sso_endpoints"        json:"ssoEndpoints"`
	SloEndpoints        types.List   `tfsdk:"slo_endpoints"        json:"sloEndpoints"`
}

type oidcAdditionalConfigData struct {
	EnforceNonce             types.Bool `tfsdk:"enforce_nonce"               json:"enforceNonce"`
	EnforceIssuerCheck       types.Bool `tfsdk:"enforce_issuer_check"        json:"enforceIssuerCheck"`
	DisableLogoutIdTokenHint types.Bool `tfsdk:"disable_logout_id_token_hint" json:"omitIDTokenHintForLogout"`
}

type oidcConfigData struct {
	DiscoveryUrl             types.String `tfsdk:"discovery_url"               json:"discoveryUrl"`
	ClientId                 types.String `tfsdk:"client_id"                   json:"clientId"`
	ClientSecret             types.String `tfsdk:"client_secret"               json:"clientSecret"`
	SubjectNameIdentifier    types.String `tfsdk:"subject_name_identifier"     json:"subjectNameIdentifier"`
	TokenEndpointAuthMethod  types.String `tfsdk:"token_endpoint_auth_method"  json:"tokenEndpointAuthMethod"`
	Scopes                   types.Set    `tfsdk:"scopes"                      json:"scopes"`
	PkceEnabled              types.Bool   `tfsdk:"enable_pkce"                 json:"pkceEnabled"`
	AdditionalConfig         types.Object `tfsdk:"additional_config"           json:"additionalConfig"`
	Issuer                   types.String `tfsdk:"issuer"                      json:"issuer"`
	JwksUri                  types.String `tfsdk:"jwks_uri"                    json:"jwksUri"`
	Jwks                     types.String `tfsdk:"jwks"                        json:"jwkSetPlain"`
	TokenEndpoint            types.String `tfsdk:"token_endpoint"              json:"tokenEndpoint"`
	AuthorizationEndpoint    types.String `tfsdk:"authorization_endpoint"      json:"authorizationEndpoint"`
	LogoutEndpoint           types.String `tfsdk:"logout_endpoint"             json:"endSessionEndpoint"`
	UserInfoEndpoint         types.String `tfsdk:"user_info_endpoint"          json:"userInfoEndpoint"`
	IsClientSecretConfigured types.Bool   `tfsdk:"is_client_secret_configured" json:"isClientSecretConfigured"`
}

type loginHintConfigData struct {
	LoginHintType types.String `tfsdk:"login_hint_type" json:"loginHintType"`
	SendMethod    types.String `tfsdk:"send_method"     json:"sendMethod"`
}

type identityFederationData struct {
	UseLocalUserStore        types.Bool `tfsdk:"use_local_user_store"        json:"useLocalUserStore"`
	AllowLocalUsersOnly      types.Bool `tfsdk:"allow_local_users_only"      json:"allowLocalUsersOnly"`
	ApplyLocalIdPAuthnChecks types.Bool `tfsdk:"apply_local_idp_auth_and_checks" json:"applyLocalIdPAuthnChecks"`
	RequiredGroups           types.Set  `tfsdk:"required_groups"             json:"requiredGroups"`
}

type corporateIdPData struct {
	Id                    types.String `tfsdk:"id"                       json:"id"`
	Name                  types.String `tfsdk:"name"                     json:"name"`
	DisplayName           types.String `tfsdk:"display_name"             json:"displayName"`
	Type                  types.String `tfsdk:"type"                     json:"type"`
	LogoutUrl             types.String `tfsdk:"logout_url"               json:"logoutUrl"`
	ForwardAllSsoRequests types.Bool   `tfsdk:"forward_all_sso_requests" json:"forwardAllSsoRequests"`
	IdentityFederation    types.Object `tfsdk:"identity_federation"      json:"identityFederation"`
	LoginHintConfig       types.Object `tfsdk:"login_hint_config"        json:"loginHintConfiguration"`
	Saml2Config           types.Object `tfsdk:"saml2_config"             json:"saml2Configuration"`
	OidcConfig            types.Object `tfsdk:"oidc_config"              json:"oidcConfiguration"`
}

func corporateIdPValueFrom(ctx context.Context, c corporateidps.IdentityProvider) (corporateIdPData, diag.Diagnostics) {

	var diags, diagnostics diag.Diagnostics

	corporateIdP := corporateIdPData{
		Id:                    types.StringValue(c.Id),
		DisplayName:           types.StringValue(c.DisplayName),
		Type:                  types.StringValue(c.Type),
		ForwardAllSsoRequests: types.BoolValue(c.ForwardAllSsoRequests),
	}

	// Identity Federation
	var idf identityFederationData
	idf.UseLocalUserStore = types.BoolValue(c.IdentityFederation.UseLocalUserStore)
	idf.AllowLocalUsersOnly = types.BoolValue(c.IdentityFederation.AllowLocalUsersOnly)
	idf.ApplyLocalIdPAuthnChecks = types.BoolValue(c.IdentityFederation.ApplyLocalIdPAuthnChecks)

	if len(c.IdentityFederation.RequiredGroups) > 0 {
		idf.RequiredGroups, diags = types.SetValueFrom(ctx, types.StringType, c.IdentityFederation.RequiredGroups)
		diagnostics.Append(diags...)

		if diagnostics.HasError() {
			return corporateIdP, diagnostics
		}

	} else {
		idf.RequiredGroups = types.SetNull(types.StringType)
	}

	corporateIdP.IdentityFederation, diags = types.ObjectValueFrom(ctx, identityFederationObjType.AttrTypes, idf)
	diagnostics.Append(diags...)

	if diagnostics.HasError() {
		return corporateIdP, diagnostics
	}

	// Login Hint Configuration
	var loginHintConfig loginHintConfigData
	loginHintConfig.LoginHintType = types.StringValue(c.LoginHintConfiguration.LoginHintType)
	loginHintConfig.SendMethod = types.StringValue(c.LoginHintConfiguration.SendMethod)

	corporateIdP.LoginHintConfig, diags = types.ObjectValueFrom(ctx, loginHintConfigObjType.AttrTypes, loginHintConfig)
	diagnostics.Append(diags...)

	if diagnostics.HasError() {
		return corporateIdP, diagnostics
	}

	if len(c.LogoutUrl) > 0 {
		corporateIdP.LogoutUrl = types.StringValue(c.LogoutUrl)
	}

	if len(c.Name) > 0 {
		corporateIdP.Name = types.StringValue(c.Name)
	}

	//SAML2 Configuration
	// check that type is not OIDC
	if c.Type != idpTypeValues[3] {
		saml2Config := &saml2ConfigData{
			DigestAlgorithm: types.StringValue(c.Saml2Configuration.DigestAlgorithm),
			IncludeScoping:  types.BoolValue(c.Saml2Configuration.IncludeScoping),
			NameIdFormat:    types.StringValue(c.Saml2Configuration.DefaultNameIdFormat),
			AllowCreate:     types.StringValue(c.Saml2Configuration.AllowCreate),
		}

		if len(c.Saml2Configuration.SamlMetadataUrl) > 0 {
			saml2Config.SamlMetadataUrl = types.StringValue(c.Saml2Configuration.SamlMetadataUrl)
		}

		// SAML2 Assertion Attributes
		if len(c.Saml2Configuration.AssertionAttributes) > 0 {
			attributes, diags := types.ListValueFrom(ctx, saml2AssertionAttributeObjType, c.Saml2Configuration.AssertionAttributes)
			diagnostics.Append(diags...)

			if diags.HasError() {
				return corporateIdP, diagnostics
			}

			saml2Config.AssertionAttributes = attributes
		} else {
			saml2Config.AssertionAttributes = types.ListNull(saml2AssertionAttributeObjType)
		}

		// SAML2 Certificates
		// the mapping is done manually in order to handle empty strings
		if len(c.Saml2Configuration.CertificatesForSigning) > 0 {

			var certificatesData []signingCertificateData

			for _, certificate := range c.Saml2Configuration.CertificatesForSigning {

				certificateData := signingCertificateData{
					Base64Certificate: types.StringValue(certificate.Base64Certificate),
					Default:           types.BoolValue(certificate.IsDefault),
				}

				if len(certificate.Dn) > 0 {
					certificateData.Dn = types.StringValue(certificate.Dn)
				}

				if len(certificate.ValidFrom) > 0 {
					certificateData.ValidFrom = types.StringValue(certificate.ValidFrom)
				}

				if len(certificate.ValidTo) > 0 {
					certificateData.ValidTo = types.StringValue(certificate.ValidTo)
				}

				certificatesData = append(certificatesData, certificateData)
			}

			certificates, diags := types.ListValueFrom(ctx, saml2SigningCertificateObjType, certificatesData)
			diagnostics.Append(diags...)

			if diagnostics.HasError() {
				return corporateIdP, diagnostics
			}
			saml2Config.SigningCertificates = certificates
		} else {
			saml2Config.SigningCertificates = types.ListNull(saml2SigningCertificateObjType)
		}

		// SAML2 SSO Endpoints
		if len(c.Saml2Configuration.SsoEndpoints) > 0 {

			endpoints, diags := types.ListValueFrom(ctx, saml2SsoEndpointObjType, c.Saml2Configuration.SsoEndpoints)
			diagnostics.Append(diags...)

			if diagnostics.HasError() {
				return corporateIdP, diagnostics
			}

			saml2Config.SsoEndpoints = endpoints
		} else {
			saml2Config.SsoEndpoints = types.ListNull(saml2SsoEndpointObjType)
		}

		// SAML2 SLO Endpoints
		// the mapping is done manually in order to handle empty strings
		if len(c.Saml2Configuration.SloEndpoints) > 0 {

			var endpointsData []sloEndpointData
			for _, endpoint := range c.Saml2Configuration.SloEndpoints {

				endpointData := sloEndpointData{
					BindingName: types.StringValue(endpoint.BindingName),
					Location:    types.StringValue(endpoint.Location),
					Default:     types.BoolValue(endpoint.IsDefault),
				}

				if len(endpoint.ResponseLocation) > 0 {
					endpointData.ResponseLocation = types.StringValue(endpoint.ResponseLocation)
				} else {
					endpointData.ResponseLocation = types.StringNull()
				}

				endpointsData = append(endpointsData, endpointData)
			}

			endpoints, diags := types.ListValueFrom(ctx, saml2SloEndpointObjType, endpointsData)
			diagnostics.Append(diags...)

			if diagnostics.HasError() {
				return corporateIdP, diagnostics
			}

			saml2Config.SloEndpoints = endpoints
		} else {
			saml2Config.SloEndpoints = types.ListNull(saml2SloEndpointObjType)
		}

		corporateIdP.Saml2Config, diags = types.ObjectValueFrom(ctx, IdPSaml2ConfigObjType.AttrTypes, saml2Config)
		diagnostics.Append(diags...)
	} else {
		corporateIdP.Saml2Config = types.ObjectNull(IdPSaml2ConfigObjType.AttrTypes)
	}

	// OIDC Configuration
	if c.Type == idpTypeValues[3] {

		oidcConfig := &oidcConfigData{
			DiscoveryUrl:             types.StringValue(c.OidcConfiguration.DiscoveryUrl),
			ClientId:                 types.StringValue(c.OidcConfiguration.ClientId),
			SubjectNameIdentifier:    types.StringValue(c.OidcConfiguration.SubjectNameIdentifier),
			TokenEndpointAuthMethod:  types.StringValue(c.OidcConfiguration.TokenEndpointAuthMethod),
			PkceEnabled:              types.BoolValue(c.OidcConfiguration.PkceEnabled),
			Issuer:                   types.StringValue(c.OidcConfiguration.Issuer),
			JwksUri:                  types.StringValue(c.OidcConfiguration.JwksUri),
			Jwks:                     types.StringValue(c.OidcConfiguration.JwkSetPlain),
			IsClientSecretConfigured: types.BoolValue(c.OidcConfiguration.IsClientSecretConfigured),
		}

		// OIDC Scopes
		scopes, diags := types.SetValueFrom(ctx, types.StringType, c.OidcConfiguration.Scopes)
		diagnostics.Append(diags...)

		if diagnostics.HasError() {
			return corporateIdP, diagnostics
		}

		oidcConfig.Scopes = scopes

		//OIDC Additional Config
		if c.OidcConfiguration.AdditionalConfig != nil {

			additionalConfig, diags := types.ObjectValueFrom(ctx, OidcCAdditionalConfigObjType.AttrTypes, c.OidcConfiguration.AdditionalConfig)
			diagnostics.Append(diags...)

			if diagnostics.HasError() {
				return corporateIdP, diagnostics
			}

			oidcConfig.AdditionalConfig = additionalConfig
		} else {
			oidcConfig.AdditionalConfig = types.ObjectNull(OidcCAdditionalConfigObjType.AttrTypes)
		}

		//OIDC Endpoints
		if len(c.OidcConfiguration.TokenEndpoint) > 0 {
			oidcConfig.TokenEndpoint = types.StringValue(c.OidcConfiguration.TokenEndpoint)
		}

		if len(c.OidcConfiguration.AuthorizationEndpoint) > 0 {
			oidcConfig.AuthorizationEndpoint = types.StringValue(c.OidcConfiguration.AuthorizationEndpoint)
		}

		if len(c.OidcConfiguration.EndSessionEndpoint) > 0 {
			oidcConfig.LogoutEndpoint = types.StringValue(c.OidcConfiguration.EndSessionEndpoint)
		}

		if len(c.OidcConfiguration.UserInfoEndpoint) > 0 {
			oidcConfig.UserInfoEndpoint = types.StringValue(c.OidcConfiguration.UserInfoEndpoint)
		}

		corporateIdP.OidcConfig, diags = types.ObjectValueFrom(ctx, oidcConfigObjType.AttrTypes, oidcConfig)
		diagnostics.Append(diags...)
	} else {
		corporateIdP.OidcConfig = types.ObjectNull(oidcConfigObjType.AttrTypes)
	}

	return corporateIdP, diagnostics
}

func corporateIdPsValueFrom(ctx context.Context, c corporateidps.IdentityProvidersResponse) []corporateIdPData {

	idps := []corporateIdPData{}

	for _, res := range c.IdentityProviders {
		idp, _ := corporateIdPValueFrom(ctx, res)
		idps = append(idps, idp)
	}

	return idps
}

func (r *corporateIdPResource) getCorporateIdPRequest(ctx context.Context, plan corporateIdPData) (*corporateidps.IdentityProvider, diag.Diagnostics) {
	var diags, diagnostics diag.Diagnostics

	corporateIdP := &corporateidps.IdentityProvider{
		DisplayName: plan.DisplayName.ValueString(),
	}

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		corporateIdP.Name = plan.Name.ValueString()
	}

	if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
		corporateIdP.Type = plan.Type.ValueString()
	}

	if !plan.LogoutUrl.IsNull() && !plan.LogoutUrl.IsUnknown() {
		corporateIdP.LogoutUrl = plan.LogoutUrl.ValueString()
	}

	if !plan.ForwardAllSsoRequests.IsNull() && !plan.ForwardAllSsoRequests.IsUnknown() {
		corporateIdP.ForwardAllSsoRequests = plan.ForwardAllSsoRequests.ValueBool()
	}

	if !plan.IdentityFederation.IsNull() && !plan.IdentityFederation.IsUnknown() {

		var idF corporateidps.IdentityFederation
		diags = plan.IdentityFederation.As(ctx, &idF, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return nil, diagnostics
		}

		corporateIdP.IdentityFederation = &idF
	}

	if !plan.LoginHintConfig.IsNull() && !plan.LoginHintConfig.IsUnknown() {

		var loginHintConfig corporateidps.LoginHintConfiguration
		diags = plan.LoginHintConfig.As(ctx, &loginHintConfig, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return &corporateidps.IdentityProvider{}, diagnostics
		}

		corporateIdP.LoginHintConfiguration = &loginHintConfig

	}

	if !plan.Saml2Config.IsNull() && !plan.Saml2Config.IsUnknown() {

		var saml2Config corporateidps.SAML2Configuration
		diags = plan.Saml2Config.As(ctx, &saml2Config, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})
		diagnostics.Append(diags...)

		if diagnostics.HasError() {
			return nil, diagnostics
		}

		corporateIdP.Saml2Configuration = &saml2Config
	}

	if !plan.OidcConfig.IsNull() && !plan.OidcConfig.IsUnknown() {

		var oidcConfig corporateidps.OIDCConfiguration
		diags = plan.OidcConfig.As(ctx, &oidcConfig, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})
		diagnostics.Append(diags...)

		corporateIdP.OidcConfiguration = &oidcConfig

	}

	return corporateIdP, diagnostics
}

func getCorporateIdPUpdateRequest(ctx context.Context, plan corporateIdPData, state corporateIdPData) ([]generic.PatchRequest, diag.Diagnostics) {

	var diagnostics diag.Diagnostics
	reqs := []generic.PatchRequest{}

	idpType := reflect.TypeFor[corporateIdPData]()

	if !plan.DisplayName.Equal(state.DisplayName) {
		patchReq, diags := utils.GetPatchRequest("DisplayName", "", plan.DisplayName.ValueString(), idpType)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return reqs, diagnostics
		}
		reqs = append(reqs, patchReq)
	}

	if !plan.Name.Equal(state.Name) {
		if plan.Name.IsNull() || plan.Name.IsUnknown() {
			reqs = append(reqs, utils.GenerateDeletePatchRequest("/name"))
		} else {
			patchReq, diags := utils.GetPatchRequest("Name", "", plan.Name.ValueString(), idpType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			reqs = append(reqs, patchReq)
		}
	}

	if !plan.Type.Equal(state.Type) {
		patchReq, diags := utils.GetPatchRequest("Type", "", plan.Type.ValueString(), idpType)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return reqs, diagnostics
		}
		reqs = append(reqs, patchReq)
	}

	if !plan.LogoutUrl.Equal(state.LogoutUrl) {
		if plan.LogoutUrl.IsNull() || plan.LogoutUrl.IsUnknown() {
			reqs = append(reqs, utils.GenerateDeletePatchRequest("/logoutUrl"))
		} else {
			patchReq, diags := utils.GetPatchRequest("LogoutUrl", "", plan.LogoutUrl.ValueString(), idpType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			reqs = append(reqs, patchReq)
		}
	}

	if !plan.ForwardAllSsoRequests.Equal(state.ForwardAllSsoRequests) {
		patchReq, diags := utils.GetPatchRequest("ForwardAllSsoRequests", "", plan.ForwardAllSsoRequests.ValueBool(), idpType)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return reqs, diagnostics
		}
		reqs = append(reqs, patchReq)
	}

	// IdentityFederation — drill into leaf fields
	if !plan.IdentityFederation.Equal(state.IdentityFederation) {
		idfPath, diags := utils.GetAttributeTag("IdentityFederation", idpType)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return reqs, diagnostics
		}

		var planIdF, stateIdF identityFederationData
		diags = plan.IdentityFederation.As(ctx, &planIdF, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
		diagnostics.Append(diags...)
		diags = state.IdentityFederation.As(ctx, &stateIdF, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return reqs, diagnostics
		}

		idfType := reflect.TypeFor[identityFederationData]()

		if !planIdF.UseLocalUserStore.Equal(stateIdF.UseLocalUserStore) {
			patchReq, diags := utils.GetPatchRequest("UseLocalUserStore", idfPath, planIdF.UseLocalUserStore.ValueBool(), idfType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			reqs = append(reqs, patchReq)
		}

		if !planIdF.AllowLocalUsersOnly.Equal(stateIdF.AllowLocalUsersOnly) {
			patchReq, diags := utils.GetPatchRequest("AllowLocalUsersOnly", idfPath, planIdF.AllowLocalUsersOnly.ValueBool(), idfType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			reqs = append(reqs, patchReq)
		}

		if !planIdF.ApplyLocalIdPAuthnChecks.Equal(stateIdF.ApplyLocalIdPAuthnChecks) {
			patchReq, diags := utils.GetPatchRequest("ApplyLocalIdPAuthnChecks", idfPath, planIdF.ApplyLocalIdPAuthnChecks.ValueBool(), idfType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			reqs = append(reqs, patchReq)
		}

		if !planIdF.RequiredGroups.Equal(stateIdF.RequiredGroups) {
			val := []string{}
			if !planIdF.RequiredGroups.IsNull() {
				diags = planIdF.RequiredGroups.ElementsAs(ctx, &val, true)
				diagnostics.Append(diags...)
				if diagnostics.HasError() {
					return reqs, diagnostics
				}
			}
			patchReq, diags := utils.GetPatchRequest("RequiredGroups", idfPath, val, idfType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			reqs = append(reqs, patchReq)
		}
	}

	// LoginHintConfig — drill into leaf fields
	if !plan.LoginHintConfig.Equal(state.LoginHintConfig) {
		lhcPath, diags := utils.GetAttributeTag("LoginHintConfig", idpType)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return reqs, diagnostics
		}

		var planLhc, stateLhc loginHintConfigData
		diags = plan.LoginHintConfig.As(ctx, &planLhc, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
		diagnostics.Append(diags...)
		diags = state.LoginHintConfig.As(ctx, &stateLhc, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return reqs, diagnostics
		}

		lhcType := reflect.TypeFor[loginHintConfigData]()

		if !planLhc.LoginHintType.Equal(stateLhc.LoginHintType) {
			patchReq, diags := utils.GetPatchRequest("LoginHintType", lhcPath, planLhc.LoginHintType.ValueString(), lhcType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			reqs = append(reqs, patchReq)
		}

		if !planLhc.SendMethod.Equal(stateLhc.SendMethod) {
			patchReq, diags := utils.GetPatchRequest("SendMethod", lhcPath, planLhc.SendMethod.ValueString(), lhcType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			reqs = append(reqs, patchReq)
		}
	}

	// Saml2Config — drill into leaf fields
	if !plan.Saml2Config.Equal(state.Saml2Config) {
		samlPath, diags := utils.GetAttributeTag("Saml2Config", idpType)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return reqs, diagnostics
		}

		var planSaml, stateSaml saml2ConfigData
		diags = plan.Saml2Config.As(ctx, &planSaml, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
		diagnostics.Append(diags...)
		diags = state.Saml2Config.As(ctx, &stateSaml, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return reqs, diagnostics
		}

		samlType := reflect.TypeFor[saml2ConfigData]()

		if !planSaml.SamlMetadataUrl.Equal(stateSaml.SamlMetadataUrl) {
			patchReq, diags := utils.GetPatchRequest("SamlMetadataUrl", samlPath, planSaml.SamlMetadataUrl.ValueString(), samlType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			reqs = append(reqs, patchReq)
		}

		if !planSaml.DigestAlgorithm.Equal(stateSaml.DigestAlgorithm) {
			patchReq, diags := utils.GetPatchRequest("DigestAlgorithm", samlPath, planSaml.DigestAlgorithm.ValueString(), samlType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			reqs = append(reqs, patchReq)
		}

		if !planSaml.IncludeScoping.Equal(stateSaml.IncludeScoping) {
			patchReq, diags := utils.GetPatchRequest("IncludeScoping", samlPath, planSaml.IncludeScoping.ValueBool(), samlType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			reqs = append(reqs, patchReq)
		}

		if !planSaml.NameIdFormat.Equal(stateSaml.NameIdFormat) {
			patchReq, diags := utils.GetPatchRequest("NameIdFormat", samlPath, planSaml.NameIdFormat.ValueString(), samlType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			reqs = append(reqs, patchReq)
		}

		if !planSaml.AllowCreate.Equal(stateSaml.AllowCreate) {
			patchReq, diags := utils.GetPatchRequest("AllowCreate", samlPath, planSaml.AllowCreate.ValueString(), samlType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			reqs = append(reqs, patchReq)
		}

		if !planSaml.AssertionAttributes.Equal(stateSaml.AssertionAttributes) {
			val := []corporateidps.AssertionAttribute{}
			if !planSaml.AssertionAttributes.IsNull() {
				diags = planSaml.AssertionAttributes.ElementsAs(ctx, &val, true)
				diagnostics.Append(diags...)
				if diagnostics.HasError() {
					return reqs, diagnostics
				}
			}
			patchReq, diags := utils.GetPatchRequest("AssertionAttributes", samlPath, val, samlType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			reqs = append(reqs, patchReq)
		}

		if !planSaml.SigningCertificates.Equal(stateSaml.SigningCertificates) {
			val := []corporateidps.SigningCertificateData{}
			if !planSaml.SigningCertificates.IsNull() {
				diags = planSaml.SigningCertificates.ElementsAs(ctx, &val, true)
				diagnostics.Append(diags...)
				if diagnostics.HasError() {
					return reqs, diagnostics
				}
			}
			patchReq, diags := utils.GetPatchRequest("SigningCertificates", samlPath, val, samlType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			reqs = append(reqs, patchReq)
		}

		if !planSaml.SsoEndpoints.Equal(stateSaml.SsoEndpoints) {
			val := []corporateidps.SAML2SSOEndpoint{}
			if !planSaml.SsoEndpoints.IsNull() {
				diags = planSaml.SsoEndpoints.ElementsAs(ctx, &val, true)
				diagnostics.Append(diags...)
				if diagnostics.HasError() {
					return reqs, diagnostics
				}
			}
			patchReq, diags := utils.GetPatchRequest("SsoEndpoints", samlPath, val, samlType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			reqs = append(reqs, patchReq)
		}

		if !planSaml.SloEndpoints.Equal(stateSaml.SloEndpoints) {
			val := []corporateidps.SAML2SLOEndpoint{}
			if !planSaml.SloEndpoints.IsNull() {
				diags = planSaml.SloEndpoints.ElementsAs(ctx, &val, true)
				diagnostics.Append(diags...)
				if diagnostics.HasError() {
					return reqs, diagnostics
				}
			}
			patchReq, diags := utils.GetPatchRequest("SloEndpoints", samlPath, val, samlType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			reqs = append(reqs, patchReq)
		}
	}

	// OidcConfig — drill into leaf fields
	if !plan.OidcConfig.Equal(state.OidcConfig) {
		oidcPath, diags := utils.GetAttributeTag("OidcConfig", idpType)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return reqs, diagnostics
		}

		var planOidc, stateOidc oidcConfigData
		diags = plan.OidcConfig.As(ctx, &planOidc, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
		diagnostics.Append(diags...)
		diags = state.OidcConfig.As(ctx, &stateOidc, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return reqs, diagnostics
		}

		oidcType := reflect.TypeFor[oidcConfigData]()

		if !planOidc.DiscoveryUrl.Equal(stateOidc.DiscoveryUrl) {
			patchReq, diags := utils.GetPatchRequest("DiscoveryUrl", oidcPath, planOidc.DiscoveryUrl.ValueString(), oidcType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			reqs = append(reqs, patchReq)
		}

		if !planOidc.ClientId.Equal(stateOidc.ClientId) {
			patchReq, diags := utils.GetPatchRequest("ClientId", oidcPath, planOidc.ClientId.ValueString(), oidcType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			reqs = append(reqs, patchReq)
		}

		if !planOidc.ClientSecret.Equal(stateOidc.ClientSecret) {
			patchReq, diags := utils.GetPatchRequest("ClientSecret", oidcPath, planOidc.ClientSecret.ValueString(), oidcType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			reqs = append(reqs, patchReq)
		}

		if !planOidc.SubjectNameIdentifier.Equal(stateOidc.SubjectNameIdentifier) {
			patchReq, diags := utils.GetPatchRequest("SubjectNameIdentifier", oidcPath, planOidc.SubjectNameIdentifier.ValueString(), oidcType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			reqs = append(reqs, patchReq)
		}

		if !planOidc.TokenEndpointAuthMethod.Equal(stateOidc.TokenEndpointAuthMethod) {
			patchReq, diags := utils.GetPatchRequest("TokenEndpointAuthMethod", oidcPath, planOidc.TokenEndpointAuthMethod.ValueString(), oidcType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			reqs = append(reqs, patchReq)
		}

		if !planOidc.Scopes.Equal(stateOidc.Scopes) {
			val := []string{}
			if !planOidc.Scopes.IsNull() {
				diags = planOidc.Scopes.ElementsAs(ctx, &val, true)
				diagnostics.Append(diags...)
				if diagnostics.HasError() {
					return reqs, diagnostics
				}
			}
			patchReq, diags := utils.GetPatchRequest("Scopes", oidcPath, val, oidcType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			reqs = append(reqs, patchReq)
		}

		if !planOidc.PkceEnabled.Equal(stateOidc.PkceEnabled) {
			patchReq, diags := utils.GetPatchRequest("PkceEnabled", oidcPath, planOidc.PkceEnabled.ValueBool(), oidcType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			reqs = append(reqs, patchReq)
		}

		// AdditionalConfig — drill one more level
		if !planOidc.AdditionalConfig.Equal(stateOidc.AdditionalConfig) {
			additionalConfigTag, diags := utils.GetAttributeTag("AdditionalConfig", oidcType)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}
			additionalConfigPath := fmt.Sprintf("%s/%s", oidcPath, additionalConfigTag)

			var planAdditional, stateAdditional oidcAdditionalConfigData
			diags = planOidc.AdditionalConfig.As(ctx, &planAdditional, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
			diagnostics.Append(diags...)
			diags = stateOidc.AdditionalConfig.As(ctx, &stateAdditional, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return reqs, diagnostics
			}

			additionalConfigType := reflect.TypeFor[oidcAdditionalConfigData]()

			if !planAdditional.EnforceNonce.Equal(stateAdditional.EnforceNonce) {
				patchReq, diags := utils.GetPatchRequest("EnforceNonce", additionalConfigPath, planAdditional.EnforceNonce.ValueBool(), additionalConfigType)
				diagnostics.Append(diags...)
				if diagnostics.HasError() {
					return reqs, diagnostics
				}
				reqs = append(reqs, patchReq)
			}

			if !planAdditional.EnforceIssuerCheck.Equal(stateAdditional.EnforceIssuerCheck) {
				patchReq, diags := utils.GetPatchRequest("EnforceIssuerCheck", additionalConfigPath, planAdditional.EnforceIssuerCheck.ValueBool(), additionalConfigType)
				diagnostics.Append(diags...)
				if diagnostics.HasError() {
					return reqs, diagnostics
				}
				reqs = append(reqs, patchReq)
			}

			if !planAdditional.DisableLogoutIdTokenHint.Equal(stateAdditional.DisableLogoutIdTokenHint) {
				patchReq, diags := utils.GetPatchRequest("DisableLogoutIdTokenHint", additionalConfigPath, planAdditional.DisableLogoutIdTokenHint.ValueBool(), additionalConfigType)
				diagnostics.Append(diags...)
				if diagnostics.HasError() {
					return reqs, diagnostics
				}
				reqs = append(reqs, patchReq)
			}
		}
	}

	return reqs, diagnostics
}

func mapOidcClientSecret(ctx context.Context, plan corporateIdPData, state *corporateIdPData) diag.Diagnostics {

	var oidcPlan oidcConfigData
	diags := plan.OidcConfig.As(ctx, &oidcPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})
	if diags.HasError() {
		return diags
	}

	if !oidcPlan.ClientSecret.IsNull() && !oidcPlan.ClientSecret.IsUnknown() {

		var oidcState oidcConfigData
		diags = state.OidcConfig.As(ctx, &oidcState, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})
		if diags.HasError() {
			return diags
		}

		oidcState.ClientSecret = oidcPlan.ClientSecret

		state.OidcConfig, diags = types.ObjectValueFrom(ctx, oidcConfigObjType.AttrTypes, oidcState)
		if diags.HasError() {
			return diags
		}
	}

	return nil
}
