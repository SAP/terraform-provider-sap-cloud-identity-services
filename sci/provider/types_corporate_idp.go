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

// appendPatch resolves the JSON path for fieldName, builds a replace patch, and appends it to reqs.
// On error it records diagnostics and returns reqs unchanged.
func appendPatch(reqs []generic.PatchRequest, diagnostics *diag.Diagnostics, fieldName, path string, value any, t reflect.Type) []generic.PatchRequest {
	patchReq, diags := utils.GetPatchRequest(fieldName, path, value, t)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return reqs
	}
	return append(reqs, patchReq)
}

func diffIdentityFederation(ctx context.Context, plan, state corporateIdPData, basePath string) ([]generic.PatchRequest, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	reqs := []generic.PatchRequest{}

	var planIdF, stateIdF identityFederationData
	diagnostics.Append(plan.IdentityFederation.As(ctx, &planIdF, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})...)
	diagnostics.Append(state.IdentityFederation.As(ctx, &stateIdF, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})...)
	if diagnostics.HasError() {
		return reqs, diagnostics
	}

	idfType := reflect.TypeFor[identityFederationData]()

	if !planIdF.UseLocalUserStore.Equal(stateIdF.UseLocalUserStore) {
		reqs = appendPatch(reqs, &diagnostics, "UseLocalUserStore", basePath, planIdF.UseLocalUserStore.ValueBool(), idfType)
	}
	if !planIdF.AllowLocalUsersOnly.Equal(stateIdF.AllowLocalUsersOnly) {
		reqs = appendPatch(reqs, &diagnostics, "AllowLocalUsersOnly", basePath, planIdF.AllowLocalUsersOnly.ValueBool(), idfType)
	}
	if !planIdF.ApplyLocalIdPAuthnChecks.Equal(stateIdF.ApplyLocalIdPAuthnChecks) {
		reqs = appendPatch(reqs, &diagnostics, "ApplyLocalIdPAuthnChecks", basePath, planIdF.ApplyLocalIdPAuthnChecks.ValueBool(), idfType)
	}
	if !planIdF.RequiredGroups.Equal(stateIdF.RequiredGroups) {
		val := []string{}
		if !planIdF.RequiredGroups.IsNull() {
			diagnostics.Append(planIdF.RequiredGroups.ElementsAs(ctx, &val, true)...)
		}
		reqs = appendPatch(reqs, &diagnostics, "RequiredGroups", basePath, val, idfType)
	}

	return reqs, diagnostics
}

func diffLoginHintConfig(ctx context.Context, plan, state corporateIdPData, basePath string) ([]generic.PatchRequest, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	reqs := []generic.PatchRequest{}

	var planLhc, stateLhc loginHintConfigData
	diagnostics.Append(plan.LoginHintConfig.As(ctx, &planLhc, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})...)
	diagnostics.Append(state.LoginHintConfig.As(ctx, &stateLhc, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})...)
	if diagnostics.HasError() {
		return reqs, diagnostics
	}

	lhcType := reflect.TypeFor[loginHintConfigData]()

	if !planLhc.LoginHintType.Equal(stateLhc.LoginHintType) {
		reqs = appendPatch(reqs, &diagnostics, "LoginHintType", basePath, planLhc.LoginHintType.ValueString(), lhcType)
	}
	if !planLhc.SendMethod.Equal(stateLhc.SendMethod) {
		reqs = appendPatch(reqs, &diagnostics, "SendMethod", basePath, planLhc.SendMethod.ValueString(), lhcType)
	}

	return reqs, diagnostics
}

func diffSaml2Config(ctx context.Context, plan, state corporateIdPData, basePath string) ([]generic.PatchRequest, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	reqs := []generic.PatchRequest{}

	var planSaml, stateSaml saml2ConfigData
	diagnostics.Append(plan.Saml2Config.As(ctx, &planSaml, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})...)
	diagnostics.Append(state.Saml2Config.As(ctx, &stateSaml, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})...)
	if diagnostics.HasError() {
		return reqs, diagnostics
	}

	samlType := reflect.TypeFor[saml2ConfigData]()

	if !planSaml.SamlMetadataUrl.Equal(stateSaml.SamlMetadataUrl) {
		reqs = appendPatch(reqs, &diagnostics, "SamlMetadataUrl", basePath, planSaml.SamlMetadataUrl.ValueString(), samlType)
	}
	if !planSaml.DigestAlgorithm.Equal(stateSaml.DigestAlgorithm) {
		reqs = appendPatch(reqs, &diagnostics, "DigestAlgorithm", basePath, planSaml.DigestAlgorithm.ValueString(), samlType)
	}
	if !planSaml.IncludeScoping.Equal(stateSaml.IncludeScoping) {
		reqs = appendPatch(reqs, &diagnostics, "IncludeScoping", basePath, planSaml.IncludeScoping.ValueBool(), samlType)
	}
	if !planSaml.NameIdFormat.Equal(stateSaml.NameIdFormat) {
		reqs = appendPatch(reqs, &diagnostics, "NameIdFormat", basePath, planSaml.NameIdFormat.ValueString(), samlType)
	}
	if !planSaml.AllowCreate.Equal(stateSaml.AllowCreate) {
		reqs = appendPatch(reqs, &diagnostics, "AllowCreate", basePath, planSaml.AllowCreate.ValueString(), samlType)
	}

	if !planSaml.AssertionAttributes.Equal(stateSaml.AssertionAttributes) {
		val := []corporateidps.AssertionAttribute{}
		if !planSaml.AssertionAttributes.IsNull() {
			diagnostics.Append(planSaml.AssertionAttributes.ElementsAs(ctx, &val, true)...)
		}
		reqs = appendPatch(reqs, &diagnostics, "AssertionAttributes", basePath, val, samlType)
	}
	if !planSaml.SigningCertificates.Equal(stateSaml.SigningCertificates) {
		val := []corporateidps.SigningCertificateData{}
		if !planSaml.SigningCertificates.IsNull() {
			diagnostics.Append(planSaml.SigningCertificates.ElementsAs(ctx, &val, true)...)
		}
		reqs = appendPatch(reqs, &diagnostics, "SigningCertificates", basePath, val, samlType)
	}
	if !planSaml.SsoEndpoints.Equal(stateSaml.SsoEndpoints) {
		val := []corporateidps.SAML2SSOEndpoint{}
		if !planSaml.SsoEndpoints.IsNull() {
			diagnostics.Append(planSaml.SsoEndpoints.ElementsAs(ctx, &val, true)...)
		}
		reqs = appendPatch(reqs, &diagnostics, "SsoEndpoints", basePath, val, samlType)
	}
	if !planSaml.SloEndpoints.Equal(stateSaml.SloEndpoints) {
		val := []corporateidps.SAML2SLOEndpoint{}
		if !planSaml.SloEndpoints.IsNull() {
			diagnostics.Append(planSaml.SloEndpoints.ElementsAs(ctx, &val, true)...)
		}
		reqs = appendPatch(reqs, &diagnostics, "SloEndpoints", basePath, val, samlType)
	}

	return reqs, diagnostics
}

func diffOidcConfig(ctx context.Context, plan, state corporateIdPData, basePath string) ([]generic.PatchRequest, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	reqs := []generic.PatchRequest{}

	var planOidc, stateOidc oidcConfigData
	diagnostics.Append(plan.OidcConfig.As(ctx, &planOidc, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})...)
	diagnostics.Append(state.OidcConfig.As(ctx, &stateOidc, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})...)
	if diagnostics.HasError() {
		return reqs, diagnostics
	}

	oidcType := reflect.TypeFor[oidcConfigData]()

	if !planOidc.DiscoveryUrl.Equal(stateOidc.DiscoveryUrl) {
		reqs = appendPatch(reqs, &diagnostics, "DiscoveryUrl", basePath, planOidc.DiscoveryUrl.ValueString(), oidcType)
	}
	if !planOidc.ClientId.Equal(stateOidc.ClientId) {
		reqs = appendPatch(reqs, &diagnostics, "ClientId", basePath, planOidc.ClientId.ValueString(), oidcType)
	}
	if !planOidc.ClientSecret.Equal(stateOidc.ClientSecret) {
		reqs = appendPatch(reqs, &diagnostics, "ClientSecret", basePath, planOidc.ClientSecret.ValueString(), oidcType)
	}
	if !planOidc.SubjectNameIdentifier.Equal(stateOidc.SubjectNameIdentifier) {
		reqs = appendPatch(reqs, &diagnostics, "SubjectNameIdentifier", basePath, planOidc.SubjectNameIdentifier.ValueString(), oidcType)
	}
	if !planOidc.TokenEndpointAuthMethod.Equal(stateOidc.TokenEndpointAuthMethod) {
		reqs = appendPatch(reqs, &diagnostics, "TokenEndpointAuthMethod", basePath, planOidc.TokenEndpointAuthMethod.ValueString(), oidcType)
	}
	if !planOidc.PkceEnabled.Equal(stateOidc.PkceEnabled) {
		reqs = appendPatch(reqs, &diagnostics, "PkceEnabled", basePath, planOidc.PkceEnabled.ValueBool(), oidcType)
	}

	if !planOidc.Scopes.Equal(stateOidc.Scopes) {
		val := []string{}
		if !planOidc.Scopes.IsNull() {
			diagnostics.Append(planOidc.Scopes.ElementsAs(ctx, &val, true)...)
		}
		reqs = appendPatch(reqs, &diagnostics, "Scopes", basePath, val, oidcType)
	}

	if !planOidc.AdditionalConfig.Equal(stateOidc.AdditionalConfig) {
		additionalConfigTag, diags := utils.GetAttributeTag("AdditionalConfig", oidcType)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return reqs, diagnostics
		}
		additionalConfigPath := fmt.Sprintf("%s/%s", basePath, additionalConfigTag)

		var planAdditional, stateAdditional oidcAdditionalConfigData
		diagnostics.Append(planOidc.AdditionalConfig.As(ctx, &planAdditional, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})...)
		diagnostics.Append(stateOidc.AdditionalConfig.As(ctx, &stateAdditional, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})...)
		if diagnostics.HasError() {
			return reqs, diagnostics
		}

		additionalConfigType := reflect.TypeFor[oidcAdditionalConfigData]()

		if !planAdditional.EnforceNonce.Equal(stateAdditional.EnforceNonce) {
			reqs = appendPatch(reqs, &diagnostics, "EnforceNonce", additionalConfigPath, planAdditional.EnforceNonce.ValueBool(), additionalConfigType)
		}
		if !planAdditional.EnforceIssuerCheck.Equal(stateAdditional.EnforceIssuerCheck) {
			reqs = appendPatch(reqs, &diagnostics, "EnforceIssuerCheck", additionalConfigPath, planAdditional.EnforceIssuerCheck.ValueBool(), additionalConfigType)
		}
		if !planAdditional.DisableLogoutIdTokenHint.Equal(stateAdditional.DisableLogoutIdTokenHint) {
			reqs = appendPatch(reqs, &diagnostics, "DisableLogoutIdTokenHint", additionalConfigPath, planAdditional.DisableLogoutIdTokenHint.ValueBool(), additionalConfigType)
		}
	}

	return reqs, diagnostics
}

func getCorporateIdPUpdateRequest(ctx context.Context, plan corporateIdPData, state corporateIdPData) ([]generic.PatchRequest, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	reqs := []generic.PatchRequest{}

	idpType := reflect.TypeFor[corporateIdPData]()

	if !plan.DisplayName.Equal(state.DisplayName) {
		reqs = appendPatch(reqs, &diagnostics, "DisplayName", "", plan.DisplayName.ValueString(), idpType)
	}

	if !plan.Name.Equal(state.Name) {
		if plan.Name.IsNull() || plan.Name.IsUnknown() {
			reqs = append(reqs, utils.GenerateDeletePatchRequest("/name"))
		} else {
			reqs = appendPatch(reqs, &diagnostics, "Name", "", plan.Name.ValueString(), idpType)
		}
	}

	if !plan.Type.Equal(state.Type) {
		reqs = appendPatch(reqs, &diagnostics, "Type", "", plan.Type.ValueString(), idpType)
	}

	if !plan.LogoutUrl.Equal(state.LogoutUrl) {
		if plan.LogoutUrl.IsNull() || plan.LogoutUrl.IsUnknown() {
			reqs = append(reqs, utils.GenerateDeletePatchRequest("/logoutUrl"))
		} else {
			reqs = appendPatch(reqs, &diagnostics, "LogoutUrl", "", plan.LogoutUrl.ValueString(), idpType)
		}
	}

	if !plan.ForwardAllSsoRequests.Equal(state.ForwardAllSsoRequests) {
		reqs = appendPatch(reqs, &diagnostics, "ForwardAllSsoRequests", "", plan.ForwardAllSsoRequests.ValueBool(), idpType)
	}

	if diagnostics.HasError() {
		return reqs, diagnostics
	}

	type subObjDiff struct {
		changed bool
		field   string
		diff    func(path string) ([]generic.PatchRequest, diag.Diagnostics)
	}

	subObjs := []subObjDiff{
		{
			changed: !plan.IdentityFederation.Equal(state.IdentityFederation),
			field:   "IdentityFederation",
			diff: func(path string) ([]generic.PatchRequest, diag.Diagnostics) {
				return diffIdentityFederation(ctx, plan, state, path)
			},
		},
		{
			changed: !plan.LoginHintConfig.Equal(state.LoginHintConfig),
			field:   "LoginHintConfig",
			diff: func(path string) ([]generic.PatchRequest, diag.Diagnostics) {
				return diffLoginHintConfig(ctx, plan, state, path)
			},
		},
		{
			changed: !plan.Saml2Config.Equal(state.Saml2Config),
			field:   "Saml2Config",
			diff: func(path string) ([]generic.PatchRequest, diag.Diagnostics) {
				return diffSaml2Config(ctx, plan, state, path)
			},
		},
		{
			changed: !plan.OidcConfig.Equal(state.OidcConfig),
			field:   "OidcConfig",
			diff: func(path string) ([]generic.PatchRequest, diag.Diagnostics) {
				return diffOidcConfig(ctx, plan, state, path)
			},
		},
	}

	for _, s := range subObjs {
		if !s.changed {
			continue
		}
		path, diags := utils.GetAttributeTag(s.field, idpType)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return reqs, diagnostics
		}
		subReqs, diags := s.diff(path)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return reqs, diagnostics
		}
		reqs = append(reqs, subReqs...)
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
