package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

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
	Groups           types.List   `tfsdk:"groups"`
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

	// Groups
	if len(u.Groups) > 0 {
		groups, diags := types.ListValueFrom(ctx, groupListObjType, u.Groups)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return user, diagnostics
		}

		user.Groups = groups
	} else {
		user.Groups = types.ListNull(groupListObjType)
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

func getUserRequest(ctx context.Context, plan userData) (*users.User, string, diag.Diagnostics) {

	var diagnostics diag.Diagnostics

	var emails []users.Email
	diags := plan.Emails.ElementsAs(ctx, &emails, true)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return nil, "", diagnostics
	}

	var schemas []string
	diags = plan.Schemas.ElementsAs(ctx, &schemas, true)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return nil, "", diagnostics
	}

	args := &users.User{
		UserName: plan.UserName.ValueString(),
		Emails:   emails,
		Schemas:  schemas,
	}

	if !plan.DisplayName.IsNull() {
		args.DisplayName = plan.DisplayName.ValueString()
	}

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		var name users.Name
		diags = plan.Name.As(ctx, &name, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})

		diagnostics.Append(diags...)

		if diagnostics.HasError() {
			return nil, "", diagnostics
		}

		args.Name = &name
	}

	if !plan.InitialPassword.IsNull() {
		args.Password = plan.InitialPassword.ValueString()
	}

	if !plan.UserType.IsNull() && !plan.UserType.IsUnknown() {
		args.UserType = plan.UserType.ValueString()
	}

	if !plan.Active.IsNull() && !plan.Active.IsUnknown() {
		args.Active = plan.Active.ValueBool()
	}

	if !plan.SapExtensionUser.IsNull() && !plan.SapExtensionUser.IsUnknown() {

		var sapExtensionUser users.SAPExtension
		diags = plan.SapExtensionUser.As(ctx, &sapExtensionUser, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return nil, "", diagnostics
		}

		args.SAPExtension = &sapExtensionUser
	}

	var customSchemas string
	if !plan.CustomSchemas.IsNull() {
		customSchemas = plan.CustomSchemas.ValueString()
	}

	return args, customSchemas, diagnostics
}
