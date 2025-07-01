package provider

import (
	"context"

	corporateidps "github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/corporateIdps"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	SamlMetadataUrl     types.String `tfsdk:"saml_metadata_url"`
	DigestAlgorithm     types.String `tfsdk:"digest_algorithm"`
	IncludeScoping      types.Bool   `tfsdk:"include_scoping"`
	NameIdFormat        types.String `tfsdk:"name_id_format"`
	AllowCreate         types.String `tfsdk:"allow_create"`
	AssertionAttributes types.List   `tfsdk:"assertion_attributes"`
	SigningCertificates types.List   `tfsdk:"signing_certificates"`
	SsoEndpoints        types.List   `tfsdk:"sso_endpoints"`
	SloEndpoints        types.List   `tfsdk:"slo_endpoints"`
}

type loginHintConfigData struct {
	LoginHintType types.String `tfsdk:"login_hint_type"`
	SendMethod    types.String `tfsdk:"send_method"`
}

type identityFederationData struct {
	UseLocalUserStore        types.Bool `tfsdk:"use_local_user_store"`
	AllowLocalUsersOnly      types.Bool `tfsdk:"allow_local_users_only"`
	ApplyLocalIdPAuthnChecks types.Bool `tfsdk:"apply_local_idp_auth_and_checks"`
	RequiredGroups           types.Set  `tfsdk:"required_groups"`
}

type corporateIdPData struct {
	Id                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	DisplayName           types.String `tfsdk:"display_name"`
	Type                  types.String `tfsdk:"type"`
	LogoutUrl             types.String `tfsdk:"logout_url"`
	ForwardAllSsoRequests types.Bool   `tfsdk:"forward_all_sso_requests"`
	IdentityFederation    types.Object `tfsdk:"identity_federation"`
	LoginHintConfig       types.Object `tfsdk:"login_hint_config"`
	Saml2Config           types.Object `tfsdk:"saml2_config"`
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

		corporateIdP.Saml2Config, diags = types.ObjectValueFrom(ctx, saml2ConfigObjType.AttrTypes, saml2Config)
		diagnostics.Append(diags...)
	} else {
		corporateIdP.Saml2Config = types.ObjectNull(saml2ConfigObjType.AttrTypes)
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
