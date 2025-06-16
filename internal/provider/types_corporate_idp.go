package provider

import (
	"context"

	corporateidps "github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/corporateIdps"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type saml2ConfigData struct {
	SamlMetadataUrl     types.String `tfsdk:"saml_metadata_url"`
	DigestAlgorithm     types.String `tfsdk:"digest_algorithm"`
	IncludeScoping      types.Bool   `tfsdk:"include_scoping"`
	NameIdFormat        types.String `tfsdk:"name_id_format"`
	AllowCreate         types.String `tfsdk:"allow_create"`
	AssertionAttributes types.List   `tfsdk:"assertion_attributes"`
}

type loginHintConfigData struct {
	LoginHintType types.String `tfsdk:"login_hint_type"`
	SendMethod    types.String `tfsdk:"send_method"`
}

type identityFederationData struct {
	UseLocalUserStore        types.Bool `tfsdk:"use_local_user_store"`
	AllowLocalUsersOnly      types.Bool `tfsdk:"allow_local_users_only"`
	ApplyLocalIdPAuthnChecks types.Bool `tfsdk:"apply_local_idp_auth_and_checks"`
	RequiredGroups           types.List `tfsdk:"required_groups"`
}

type corporateIdPData struct {
	Id                    types.String            `tfsdk:"id"`
	Name                  types.String            `tfsdk:"name"`
	DisplayName           types.String            `tfsdk:"display_name"`
	Type                  types.String            `tfsdk:"type"`
	LogoutUrl             types.String            `tfsdk:"logout_url"`
	ForwardAllSsoRequests types.Bool              `tfsdk:"forward_all_sso_requests"`
	IdentityFederation    *identityFederationData `tfsdk:"identity_federation"`
	LoginHintConfig       *loginHintConfigData    `tfsdk:"login_hint_config"`
	Saml2Config           *saml2ConfigData        `tfsdk:"saml2_config"`
}

func corporateIdPValueFrom(ctx context.Context, c corporateidps.IdentityProvider) (corporateIdPData, diag.Diagnostics) {

	var diags, diagnostics diag.Diagnostics

	corporateIdP := corporateIdPData{
		Id:                    types.StringValue(c.Id),
		DisplayName:           types.StringValue(c.DisplayName),
		Type:                  types.StringValue(c.Type),
		ForwardAllSsoRequests: types.BoolValue(c.ForwardAllSsoRequests),
		IdentityFederation: &identityFederationData{
			UseLocalUserStore:        types.BoolValue(c.IdentityFederation.UseLocalUserStore),
			AllowLocalUsersOnly:      types.BoolValue(c.IdentityFederation.AllowLocalUsersOnly),
			ApplyLocalIdPAuthnChecks: types.BoolValue(c.IdentityFederation.ApplyLocalIdPAuthnChecks),
		},
		LoginHintConfig: &loginHintConfigData{
			LoginHintType: types.StringValue(c.LoginHintConfiguration.LoginHintType),
			SendMethod:    types.StringValue(c.LoginHintConfiguration.SendMethod),
		},
	}

	if len(c.IdentityFederation.RequiredGroups) > 0 {
		corporateIdP.IdentityFederation.RequiredGroups, diags = types.ListValueFrom(ctx, types.StringType, c.IdentityFederation.RequiredGroups)
		diagnostics.Append(diags...)
	} else {
		corporateIdP.IdentityFederation.RequiredGroups = types.ListNull(types.StringType)
	}

	if len(c.LogoutUrl) > 0 {
		corporateIdP.LogoutUrl = types.StringValue(c.LogoutUrl)
	}

	if len(c.Name) > 0 {
		corporateIdP.Name = types.StringValue(c.Name)
	}

	//SAML2 Configuration
	saml2Config := &saml2ConfigData{
		DigestAlgorithm: types.StringValue(c.Saml2Configuration.DigestAlgorithm),
		IncludeScoping:  types.BoolValue(c.Saml2Configuration.IncludeScoping),
		NameIdFormat:    types.StringValue(c.Saml2Configuration.DefaultNameIdFormat),
		AllowCreate:     types.StringValue(c.Saml2Configuration.AllowCreate),
	}

	if len(c.Saml2Configuration.SamlMetadataUrl) > 0 {
		saml2Config.SamlMetadataUrl = types.StringValue(c.Saml2Configuration.SamlMetadataUrl)
	}

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

	corporateIdP.Saml2Config = saml2Config

	return corporateIdP, diagnostics
}
