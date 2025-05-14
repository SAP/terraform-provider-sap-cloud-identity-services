package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/users"
)

type sapExtensionUserData struct {
	SendMail     types.Bool   `tfsdk:"send_mail"`
	MailVerified types.Bool   `tfsdk:"mail_verified"`
	Status       types.String `tfsdk:"status"`
}

type emailData struct {
	Value   types.String `tfsdk:"value"`
	Type    types.String `tfsdk:"type"`
	Display types.String `tfsdk:"display"`
	Primary types.Bool   `tfsdk:"primary"`
}

type nameData struct {
	FamilyName    types.String `tfsdk:"family_name"`
	GivenName     types.String `tfsdk:"given_name"`
	Formatted     types.String `tfsdk:"formatted"`
	MiddleName    types.String `tfsdk:"middle_name"`
	HonoricPrefix types.String `tfsdk:"honoric_prefix"`
	HonoricSuffix types.String `tfsdk:"honoric_suffix"`
}

type userData struct {
	Id               types.String `tfsdk:"id"`
	Schemas          types.Set    `tfsdk:"schemas"`
	UserName         types.String `tfsdk:"user_name"`
	Name             types.Object `tfsdk:"name"`
	DisplayName      types.String `tfsdk:"display_name"`
	Emails           types.Set    `tfsdk:"emails"`
	Password         types.String `tfsdk:"password"`
	Title            types.String `tfsdk:"title"`
	UserType         types.String `tfsdk:"user_type"`
	Active           types.Bool   `tfsdk:"active"`
	SapExtensionUser types.Object `tfsdk:"sap_extension_user"`
	CustomSchemas    types.String `tfsdk:"custom_schemas"`
}

func userValueFrom(ctx context.Context, u users.User, cS string) (userData, diag.Diagnostics) {
	var diagnostics, diags diag.Diagnostics

	user := userData{
		Id:          types.StringValue(u.Id),
		UserName:    types.StringValue(u.UserName),
		DisplayName: types.StringValue(u.DisplayName),
		Title:       types.StringValue(u.Title),
		UserType:    types.StringValue(u.UserType),
		Active:      types.BoolValue(u.Active),
	}

	user.Schemas, diags = types.SetValueFrom(ctx, types.StringType, u.Schemas)
	diagnostics.Append(diags...)

	userName := nameData{
		FamilyName:    types.StringValue(u.Name.FamilyName),
		GivenName:     types.StringValue(u.Name.GivenName),
		Formatted:     types.StringValue(u.Name.Formatted),
		MiddleName:    types.StringValue(u.Name.MiddleName),
		HonoricPrefix: types.StringValue(u.Name.HonoricPrefix),
		HonoricSuffix: types.StringValue(u.Name.HonoricSuffix),
	}

	user.Name, diags = types.ObjectValueFrom(ctx, nameObjType, userName)
	diagnostics.Append(diags...)

	userEmails := []emailData{}
	for _, emailRes := range u.Emails {
		userEmail := emailData{
			Value:   types.StringValue(emailRes.Value),
			Type:    types.StringValue(emailRes.Type),
			Display: types.StringValue(emailRes.Display),
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
