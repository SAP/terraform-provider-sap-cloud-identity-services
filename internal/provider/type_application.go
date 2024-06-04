package provider

import (
	"context"

	// "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-ias/internal/cli/apiObjects/applications"
)

type applicationData struct {
	//INPUT
	Id types.String `tfsdk:"id"`
	//OUTPUT
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func applicationValueFrom(ctx context.Context, a applications.ApplicationResponse) applicationData {
	return applicationData{
		Id:          types.StringValue(a.Id),
		Name:        types.StringValue(a.Name),
		Description: types.StringValue(a.Description),
	}
}
