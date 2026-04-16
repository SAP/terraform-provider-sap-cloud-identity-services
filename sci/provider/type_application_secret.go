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
	ApiNames            types.Set    `tfsdk:"api_names"`
}

func applicationSecretValueFrom(ctx context.Context, s applications.ApplicationSecret) (applicationSecretData, diag.Diagnostics) {
	var diags diag.Diagnostics

	scopes, d := types.SetValueFrom(ctx, types.StringType, s.AuthorizationScopes)
	diags.Append(d...)

	apiNames, d := types.SetValueFrom(ctx, types.StringType, s.ApiNames)
	diags.Append(d...)

	allApisAccess := types.BoolNull()
	if s.AllApisAccess != nil {
		allApisAccess = types.BoolValue(*s.AllApisAccess)
	}

	return applicationSecretData{
		Id:                  types.StringValue(s.Id),
		ClientId:            types.StringValue(s.ClientId),
		Secret:              types.StringValue(s.Secret),
		Hint:                types.StringValue(s.Hint),
		Description:         types.StringValue(s.Description),
		ValidTo:             types.StringValue(s.ValidTo),
		AuthorizationScopes: scopes,
		AllApisAccess:       allApisAccess,
		ApiNames:            apiNames,
	}, diags
}

func getApplicationSecretRequest(ctx context.Context, plan applicationSecretData) (applications.ApplicationSecretRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	var scopes []string
	if !plan.AuthorizationScopes.IsUnknown() {
		diags.Append(plan.AuthorizationScopes.ElementsAs(ctx, &scopes, false)...)
	}

	var apiNames []string
	if !plan.ApiNames.IsUnknown() {
		diags.Append(plan.ApiNames.ElementsAs(ctx, &apiNames, false)...)
	}

	var allApisAccess *bool
	if !plan.AllApisAccess.IsNull() && !plan.AllApisAccess.IsUnknown() {
		v := plan.AllApisAccess.ValueBool()
		allApisAccess = &v
	}

	return applications.ApplicationSecretRequest{
		Description:         plan.Description.ValueString(),
		ValidTo:             plan.ValidTo.ValueString(),
		AuthorizationScopes: scopes,
		AllApisAccess:       allApisAccess,
		ApiNames:            apiNames,
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

	if !plan.AuthorizationScopes.Equal(state.AuthorizationScopes) && !plan.AuthorizationScopes.IsUnknown() {
		var scopes []string
		diags.Append(plan.AuthorizationScopes.ElementsAs(ctx, &scopes, false)...)
		ops = append(ops, utils.GenerateReplacePatchRequest("/authorizationScopes", scopes))
	}

	if !plan.ApiNames.Equal(state.ApiNames) && !plan.ApiNames.IsUnknown() {
		var apiNames []string
		diags.Append(plan.ApiNames.ElementsAs(ctx, &apiNames, false)...)
		ops = append(ops, utils.GenerateReplacePatchRequest("/apiNames", apiNames))
	}

	return ops, diags
}
