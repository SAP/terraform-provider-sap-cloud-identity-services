package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/generic"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/users"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/utils"
)

type nameData struct {
	FamilyName      types.String `tfsdk:"family_name"`
	GivenName       types.String `tfsdk:"given_name"`
	HonorificPrefix types.String `tfsdk:"honorific_prefix"`
}

type userData struct {
	Id               types.String `tfsdk:"id"`
	Schemas          types.Set    `tfsdk:"schemas" json:"schemas"`
	UserName         types.String `tfsdk:"user_name" json:"userName"`
	Name             types.Object `tfsdk:"name" json:"name"`
	DisplayName      types.String `tfsdk:"display_name" json:"displayName"`
	Emails           types.Set    `tfsdk:"emails" json:"emails"`
	InitialPassword  types.String `tfsdk:"initial_password" json:"password"`
	UserType         types.String `tfsdk:"user_type" json:"userType"`
	Active           types.Bool   `tfsdk:"active" json:"active"`
	SapExtensionUser types.Object `tfsdk:"sap_extension_user" json:"urn:ietf:params:scim:schemas:extension:sap:2.0:User"`
	CustomSchemas    types.String `tfsdk:"custom_schemas"`
	Groups           types.List   `tfsdk:"groups" json:"groups"`
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
	} else {
		user.CustomSchemas = types.StringNull()
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

func getUserUpdateRequest(ctx context.Context, plan userData, state userData) ([]generic.PatchRequest, diag.Diagnostics) {

	var diags diag.Diagnostics
	reqs := []generic.PatchRequest{}

	argsType := reflect.TypeFor[userData]()

	if !plan.Schemas.Equal(state.Schemas) {

		schemas := []string{}

		if !plan.Schemas.IsNull() {
			diags = plan.Schemas.ElementsAs(ctx, &schemas, true)
			if diags.HasError() {
				return reqs, diags
			}
		}

		patchReq, diags := utils.GetScimPatchRequest("Schemas", "", schemas, argsType)
		if diags.HasError() {
			return reqs, diags
		}
		reqs = append(reqs, patchReq)
	}

	if !plan.UserName.Equal(state.UserName) {
		patchReq, diags := utils.GetScimPatchRequest("UserName", "", plan.UserName.ValueString(), argsType)
		if diags.HasError() {
			return reqs, diags
		}
		reqs = append(reqs, patchReq)
	}

	if !plan.DisplayName.Equal(state.DisplayName) {
		var displayName string
		if !plan.DisplayName.IsNull() {
			displayName = plan.DisplayName.ValueString()
		}
		patchReq, diags := utils.GetScimPatchRequest("DisplayName", "", displayName, argsType)
		if diags.HasError() {
			return reqs, diags
		}
		reqs = append(reqs, patchReq)
	}

	if !plan.Emails.Equal(state.Emails) {
		emails := []users.Email{}

		if !plan.Emails.IsNull() {
			diags = plan.Emails.ElementsAs(ctx, &emails, true)
			if diags.HasError() {
				return reqs, diags
			}
		}

		patchReq, diags := utils.GetScimPatchRequest("Emails", "", emails, argsType)
		if diags.HasError() {
			return reqs, diags
		}
		reqs = append(reqs, patchReq)
	}

	if !plan.InitialPassword.Equal(state.InitialPassword) {
		var password string
		if !plan.InitialPassword.IsNull() {
			password = plan.InitialPassword.ValueString()
		}
		patchReq, diags := utils.GetScimPatchRequest("InitialPassword", "", password, argsType)
		if diags.HasError() {
			return reqs, diags
		}
		reqs = append(reqs, patchReq)
	}

	if !plan.UserType.Equal(state.UserType) {
		var userType string
		if !plan.UserType.IsNull() && !plan.UserType.IsUnknown() {
			userType = plan.UserType.ValueString()
		}
		patchReq, diags := utils.GetScimPatchRequest("UserType", "", userType, argsType)
		if diags.HasError() {
			return reqs, diags
		}
		reqs = append(reqs, patchReq)
	}

	if !plan.Active.Equal(state.Active) {
		var active bool
		if !plan.Active.IsNull() && !plan.Active.IsUnknown() {
			active = plan.Active.ValueBool()
		}
		patchReq, diags := utils.GetScimPatchRequest("Active", "", active, argsType)
		if diags.HasError() {
			return reqs, diags
		}
		reqs = append(reqs, patchReq)
	}

	if !plan.Name.Equal(state.Name) {
		name := users.Name{}

		if !plan.Name.IsNull() {
			diags = plan.Name.As(ctx, &name, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    true,
				UnhandledUnknownAsEmpty: true,
			})
			if diags.HasError() {
				return reqs, diags
			}
		}

		patchRequest, diags := utils.GetScimPatchRequest("Name", "", name, argsType)

		if diags.HasError() {
			return reqs, diags
		}

		reqs = append(reqs, patchRequest)
	}

	if !plan.Groups.Equal(state.Groups) {
		groups := groupData{}

		if !plan.Groups.IsNull() {
			diags = plan.Groups.ElementsAs(ctx, &groups, true)
			if diags.HasError() {
				return reqs, diags
			}
		}

		patchReq, diags := utils.GetScimPatchRequest("Groups", "", groups, argsType)
		if diags.HasError() {
			return reqs, diags
		}

		reqs = append(reqs, patchReq)
	}

	if !plan.SapExtensionUser.Equal(state.SapExtensionUser) {
		sapExtensionPath, diags := utils.GetAttributeTag("SapExtensionUser", argsType)
		if diags.HasError() {
			return reqs, diags
		}

		sapExtensionArgsType := reflect.TypeFor[users.SAPExtension]()

		var planSapExtension, stateSapExtension users.SAPExtension

		diags = plan.SapExtensionUser.As(ctx, &planSapExtension, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})
		if diags.HasError() {
			return reqs, diags
		}

		diags = state.SapExtensionUser.As(ctx, &stateSapExtension, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})
		if diags.HasError() {
			return reqs, diags
		}

		if planSapExtension.SendMail != stateSapExtension.SendMail {
			patchReq, diags := utils.GetScimPatchRequest("SendMail", sapExtensionPath, planSapExtension.SendMail, sapExtensionArgsType)
			if diags.HasError() {
				return reqs, diags
			}
			reqs = append(reqs, patchReq)
		}

		if planSapExtension.MailVerified != stateSapExtension.MailVerified {
			patchReq, diags := utils.GetScimPatchRequest("MailVerified", sapExtensionPath, planSapExtension.MailVerified, sapExtensionArgsType)
			if diags.HasError() {
				return reqs, diags
			}
			reqs = append(reqs, patchReq)
		}

		if planSapExtension.Status != stateSapExtension.Status {
			patchReq, diags := utils.GetScimPatchRequest("Status", sapExtensionPath, planSapExtension.Status, sapExtensionArgsType)
			if diags.HasError() {
				return reqs, diags
			}
			reqs = append(reqs, patchReq)
		}
	}

	if !plan.CustomSchemas.Equal(state.CustomSchemas) {

		planCustomSchemas := ""

		if !plan.CustomSchemas.IsNull() {

			planCustomSchemas = plan.CustomSchemas.ValueString()

			var planCustomSchemasMap map[string]map[string]any
			if err := json.Unmarshal([]byte(planCustomSchemas), &planCustomSchemasMap); err != nil {
				diags.AddError("Failed to unmarshal custom schemas", err.Error())
				return reqs, diags
			}

			for schema, attributesMap := range planCustomSchemasMap {

				for k, v := range attributesMap {

					path := fmt.Sprintf("%s:%s", schema, k)
					patchReq := utils.GeneratePatchRequest(path, v)
					reqs = append(reqs, patchReq)

				}

			}
		}

	}

	return reqs, diags
}
