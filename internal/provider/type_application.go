package provider

import (
	"context"

	// "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-ias/internal/cli/apiObjects/applications"
)

type applicationData struct {
	//INPUT
	Id 						types.String `tfsdk:"id"`
	//OUTPUT
	Name        			types.String `tfsdk:"name"`
	Description 			types.String `tfsdk:"description"`
	ParentApplicationId		types.String `tfsdk:"parent_application_id"`
	MultiTenantApp			types.Bool	 `tfsdk:"multi_tenant_app"`
	GlobalAccount 			types.String `tfsdk:"global_account"`
}

func applicationValueFrom(_ context.Context, a applications.ApplicationResponse) applicationData {
	return applicationData{
		Id:          types.StringValue(a.Id),
		Name:        types.StringValue(a.Name),
		Description: types.StringValue(a.Description),
		ParentApplicationId: types.StringValue(a.ParentApplicationId),
		MultiTenantApp: types.BoolValue(a.MultiTenantApp),
		GlobalAccount: types.StringValue(a.GlobalAccount),
	}
}

func applicationsValueFrom(ctx context.Context, a applications.ApplicationsResponse) []applicationData {
	apps := []applicationData{}

	for _, appRes := range a.Applications {

		app := applicationValueFrom(ctx, appRes)
		apps = append(apps, app)

	}

	return apps
}
