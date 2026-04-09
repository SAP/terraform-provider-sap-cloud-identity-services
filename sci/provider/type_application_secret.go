package provider

import (
	"context"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/applications"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/generic"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type applicationSecretData struct {
	Id                  types.String `tfsdk:"id"`
	ApplicationId       types.String `tfsdk:"application_id"`
	ClientId            types.String `tfsdk:"client_id"`
	Secret              types.String `tfsdk:"secret"`
	Hint                types.String `tfsdk:"hint"`
	Description         types.String `tfsdk:"description"`
	ValidTo             types.String `tfsdk:"valid_to"`
	AuthorizationScopes types.Set    `tfsdk:"authorization_scopes"`
	AllApisAccess       types.Bool   `tfsdk:"all_apis_access"`
}

func applicationSecretValueFrom(ctx context.Context, s applications.ApplicationSecret) (applicationSecretData, diag.Diagnostics) {
	var diags diag.Diagnostics

	scopes, d := types.SetValueFrom(ctx, types.StringType, s.AuthorizationScopes)
	diags.Append(d...)

	return applicationSecretData{
		Id:                  types.StringValue(s.Id),
		ClientId:            types.StringValue(s.ClientId),
		Secret:              types.StringValue(s.Secret),
		Hint:                types.StringValue(s.Hint),
		Description:         types.StringValue(s.Description),
		ValidTo:             types.StringValue(s.ValidTo),
		AuthorizationScopes: scopes,
		AllApisAccess:       types.BoolValue(s.AllApisAccess),
	}, diags
}

func getApplicationSecretRequest(ctx context.Context, plan applicationSecretData) (applications.ApplicationSecretRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	var scopes []string
	diags.Append(plan.AuthorizationScopes.ElementsAs(ctx, &scopes, false)...)

	return applications.ApplicationSecretRequest{
		Description:         plan.Description.ValueString(),
		ValidTo:             plan.ValidTo.ValueString(),
		AuthorizationScopes: scopes,
		AllApisAccess:       plan.AllApisAccess.ValueBool(),
	}, diags
}

func getApplicationSecretUpdateRequest(ctx context.Context, plan, state applicationSecretData) ([]generic.PatchRequest, diag.Diagnostics) {
	var diags diag.Diagnostics
	var ops []generic.PatchRequest

	if !plan.Description.Equal(state.Description) {
		ops = append(ops, utils.GenerateReplacePatchRequest("/description", plan.Description.ValueString()))
	}

	if !plan.ValidTo.Equal(state.ValidTo) {
		ops = append(ops, utils.GenerateReplacePatchRequest("/validTo", plan.ValidTo.ValueString()))
	}

	if !plan.AllApisAccess.Equal(state.AllApisAccess) {
		ops = append(ops, utils.GenerateReplacePatchRequest("/allApisAccess", plan.AllApisAccess.ValueBool()))
	}

	if !plan.AuthorizationScopes.Equal(state.AuthorizationScopes) {
		var scopes []string
		diags.Append(plan.AuthorizationScopes.ElementsAs(ctx, &scopes, false)...)
		ops = append(ops, utils.GenerateReplacePatchRequest("/authorizationScopes", scopes))
	}

	return ops, diags
}
