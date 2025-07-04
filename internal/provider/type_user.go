package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/users"
)

type nameData struct {
	FamilyName      types.String `tfsdk:"family_name"`
	GivenName       types.String `tfsdk:"given_name"`
	HonorificPrefix types.String `tfsdk:"honorific_prefix"`
}

type userData struct {
	Id               types.String `tfsdk:"id"`
	Schemas          types.Set    `tfsdk:"schemas"`
	UserName         types.String `tfsdk:"user_name"`
	Name             types.Object `tfsdk:"name"`
	DisplayName      types.String `tfsdk:"display_name"`
	Emails           types.Set    `tfsdk:"emails"`
	InitialPassword  types.String `tfsdk:"initial_password"`
	UserType         types.String `tfsdk:"user_type"`
	Active           types.Bool   `tfsdk:"active"`
	SapExtensionUser types.Object `tfsdk:"sap_extension_user"`
	CustomSchemas    types.String `tfsdk:"custom_schemas"`
}

func userValueFrom(ctx context.Context, u users.User, cS string) (userData, diag.Diagnostics) {
	var diagnostics, diags diag.Diagnostics

	user := userData{
		Id:       types.StringValue(u.Id),
		UserName: types.StringValue(u.UserName),
		UserType: types.StringValue(u.UserType),
		Active:   types.BoolValue(u.Active),
	}

	// Display Name
	if len(u.DisplayName) > 0 {
		user.DisplayName = types.StringValue(u.DisplayName)
	}

	// Schemas
	user.Schemas, diags = types.SetValueFrom(ctx, types.StringType, u.Schemas)
	diagnostics.Append(diags...)

	// Emails
	userEmails, diags := types.SetValueFrom(ctx, emailObjType, u.Emails)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return user, diagnostics
	}

	user.Emails = userEmails

	// Name
	// mapping is done manually to handle null values
	if u.Name != nil {
		name := nameData{}

		if len(u.Name.FamilyName) > 0 {
			name.FamilyName = types.StringValue(u.Name.FamilyName)
		}
		if len(u.Name.GivenName) > 0 {
			name.GivenName = types.StringValue(u.Name.GivenName)
		}
		if len(u.Name.HonorificPrefix) > 0 {
			name.HonorificPrefix = types.StringValue(u.Name.HonorificPrefix)
		}

		if name == (nameData{}) {
			user.Name = types.ObjectNull(nameObjType)
		} else {
			userData, diags := types.ObjectValueFrom(ctx, nameObjType, name)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return user, diagnostics
			}

			user.Name = userData
		}
	} else {
		user.Name = types.ObjectNull(nameObjType)
	}

	// SAP Extension User
	sapExtensionUser, diags := types.ObjectValueFrom(ctx, sapExtensionUserObjType, u.SAPExtension)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return user, diagnostics
	}

	user.SapExtensionUser = sapExtensionUser

	if len(cS) > 0 {
		user.CustomSchemas = types.StringValue(cS)
	}

	return user, diagnostics
}

func usersValueFrom(ctx context.Context, u users.UsersResponse, customSchemas map[int]string) []userData {
	users := []userData{}

	for i, userRes := range u.Resources {

		user, _ := userValueFrom(ctx, userRes, customSchemas[i])
		users = append(users, user)

	}

	return users
}
