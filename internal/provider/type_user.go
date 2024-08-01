package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-ias/internal/cli/apiObjects/users"
)

type email struct {
	Value 		types.String 	`tfsdk:"value"`
	Type 		types.String 	`tfsdk:"type"`
	Display 	types.String 	`tfsdk:"display"`
	Primary 	types.Bool 		`tfsdk:"primary"`
}

type name struct {
	FamilyName 		types.String 	`tfsdk:"family_name"`
	GivenName 		types.String 	`tfsdk:"given_name"`
	Formatted 		types.String 	`tfsdk:"formatted"`
	MiddleName 		types.String 	`tfsdk:"middle_name"`
	HonoricPrefix 	types.String 	`tfsdk:"honoric_prefix"`
	HonoricSuffix   types.String 	`tfsdk:"honoric_suffix"`
}

type userData struct {
	//INPUT
	Id 				types.String `tfsdk:"id"`
	//OUTPUT
	Schemas 		types.List 	 `tfsdk:"schemas"`
	UserName 		types.String `tfsdk:"user_name"`
	Name        	types.Object `tfsdk:"name"`
	Emails			types.List 	 `tfsdk:"emails"` 				
}	

func userValueFrom(ctx context.Context, u users.User) (userData, diag.Diagnostics) {
	var diagnostics, diags diag.Diagnostics

	user := userData{
		Id:          types.StringValue(u.Id),
		UserName:    types.StringValue(u.UserName),
	}

	user.Schemas, diags = types.ListValueFrom(ctx, types.StringType, u.Schemas)
	diagnostics.Append(diags...)

	userName := name{
		FamilyName: types.StringValue(u.Name.FamilyName),
		GivenName: types.StringValue(u.Name.GivenName),
		Formatted: types.StringValue(u.Name.Formatted),
		MiddleName: types.StringValue(u.Name.MiddleName),
		HonoricPrefix: types.StringValue(u.Name.HonoricPrefix),
		HonoricSuffix: types.StringValue(u.Name.HonoricSuffix),
	}

	user.Name, diags = types.ObjectValueFrom(ctx, nameObjType, userName)
	diagnostics.Append(diags...)

	userEmails := []email{}
	for _, emailRes := range u.Emails {
		userEmail := email{
			Value: types.StringValue(emailRes.Value),
			Type: types.StringValue(emailRes.Type),
			Display: types.StringValue(emailRes.Display),
			Primary: types.BoolValue(emailRes.Primary),
		}
		userEmails = append(userEmails, userEmail)
	}

	user.Emails, diags = types.ListValueFrom(ctx, emailObjType, userEmails)
	diagnostics.Append(diags...)

	return user, diagnostics
}

func usersValueFrom(ctx context.Context, u users.UsersResponse) []userData {
	users := []userData{}

	for _, userRes := range u.Resources {

		user, _ := userValueFrom(ctx, userRes)
		users = append(users, user)

	}

	return users
}
