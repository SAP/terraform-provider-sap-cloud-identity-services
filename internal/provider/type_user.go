package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-sci/internal/cli/apiObjects/users"
)

type sapExtensionUserData struct {
	SendMail     types.Bool   `tfsdk:"send_mail"`
	MailVerified types.Bool   `tfsdk:"mail_verified"`
	Status       types.String `tfsdk:"status"`
}

type emailData struct {
	Value   types.String `tfsdk:"value"`
	Type    types.String `tfsdk:"type"`
	Primary types.Bool   `tfsdk:"primary"`
}

type nameData struct {
	FamilyName      types.String `tfsdk:"family_name"`
	GivenName       types.String `tfsdk:"given_name"`
	HonorificPrefix types.String `tfsdk:"honorific_prefix"`
}

type userData struct {
	Id               types.String `tfsdk:"id"`
	Schemas          types.Set    `tfsdk:"schemas"`
	UserName         types.String `tfsdk:"user_name"`
	Name             *nameData    `tfsdk:"name"`
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

	if len(u.DisplayName) > 0 {
		user.DisplayName = types.StringValue(u.DisplayName)
	}

	user.Schemas, diags = types.SetValueFrom(ctx, types.StringType, u.Schemas)
	diagnostics.Append(diags...)

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
		user.Name = nil
	} else {
		user.Name = &name
	}

	userEmails := []emailData{}
	for _, emailRes := range u.Emails {
		userEmail := emailData{
			Value:   types.StringValue(emailRes.Value),
			Type:    types.StringValue(emailRes.Type),
			Primary: types.BoolValue(emailRes.Primary),
		}
		userEmails = append(userEmails, userEmail)
	}

	user.Emails, diags = types.SetValueFrom(ctx, emailObjType, userEmails)
	diagnostics.Append(diags...)

	sapExtensionUser := sapExtensionUserData{
		SendMail:     types.BoolValue(u.SAPExtension.SendMail),
		MailVerified: types.BoolValue(u.SAPExtension.MailVerified),
		Status:       types.StringValue(u.SAPExtension.Status),
	}

	user.SapExtensionUser, diags = types.ObjectValueFrom(ctx, sapExtensionUserObjType, sapExtensionUser)
	diagnostics.Append(diags...)

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
