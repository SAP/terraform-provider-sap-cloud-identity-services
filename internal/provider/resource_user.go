package provider

import (
	"context"
	"fmt"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/users"
)

var (
	defaultUserSchemas = []attr.Value{
		types.StringValue("urn:ietf:params:scim:schemas:core:2.0:User"),
		types.StringValue("urn:ietf:params:scim:schemas:extension:sap:2.0:User"),
	}

	emailTypeValues = []string{"work", "home", "other"}
	userTypeValues  = []string{"public", "partner", "customer", "external", "onboardee", "employee"}
	activeValues    = []string{"active", "inactive", "new"}
)

func newUserResource() resource.Resource {
	return &userResource{}
}

type userResource struct {
	cli *cli.SciClient
}

func (d *userResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.cli = req.ProviderData.(*cli.SciClient)
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Creates a user in the SAP Cloud Identity Services.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the user.",
				Computed:            true,
				Validators: []validator.String{
					utils.ValidUUID(),
				},
			},
			"schemas": schema.SetAttribute{
				MarkdownDescription: "List of SCIM schemas to configure users. The attribute is configured with default values :\n" +
					utils.PrintDefaultSchemas(defaultUserSchemas) +
					fmt.Sprintln("\n \t If the attribute must be overridden with custom values, the default schemas must be provided in addition to the custom schemas."),
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.StringType,
						defaultUserSchemas,
					),
				),
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					utils.DefaultValuesChecker(defaultUserSchemas),
				},
			},
			"user_name": schema.StringAttribute{
				MarkdownDescription: "Unique user name of the user.",
				Required:            true,
			},
			"emails": schema.SetNestedAttribute{
				MarkdownDescription: "Emails of the user.",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"value": schema.StringAttribute{
							MarkdownDescription: "Value of the user's email.",
							Required:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the user's email. " + utils.ValidValuesString(emailTypeValues),
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf(emailTypeValues...),
							},
						},
						"primary": schema.BoolAttribute{
							MarkdownDescription: "Set the email to be primary or not.",
							Computed:            true,
							Optional:            true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
			"name": schema.SingleNestedAttribute{
				MarkdownDescription: "Name of the user",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"family_name": schema.StringAttribute{
						MarkdownDescription: "Last name of the user. The following characters: <, >, : are not allowed.",
						Optional:            true,
					},
					"given_name": schema.StringAttribute{
						MarkdownDescription: "First name of the user. The following characters: <, >, : are not allowed.",
						Optional:            true,
					},
					"honorific_prefix": schema.StringAttribute{
						MarkdownDescription: "HonorificPrefix is part of the Master Data attributes and have canonical values. The specific values for this attribute can be found on `<tenantUrl>/service/md/salutations`",
						Optional:            true,
					},
				},
			},
			"initial_password": schema.StringAttribute{
				MarkdownDescription: "The initial password to be configured for the user. If this attribute is configured, the password will have to be changed by the user at the first login.",
				Optional:            true,
				Sensitive:           true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The name to be displayed for the user.",
				Optional:            true,
			},
			"user_type": schema.StringAttribute{
				MarkdownDescription: "Specifies the type of the user. The default type is \"public\". " + utils.ValidValuesString(userTypeValues),
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(userTypeValues...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Determines whether the user is active or not. The default value for the attribute is false.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"sap_extension_user": schema.SingleNestedAttribute{
				MarkdownDescription: "Configure attributes particular to the schema `" + defaultUserSchemas[1].String() + "`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"send_mail": schema.BoolAttribute{
						MarkdownDescription: "Specifies if an activation mail should be sent. The value of the attribute only matters when creating the user.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"mail_verified": schema.BoolAttribute{
						MarkdownDescription: "The attribute specifies if the e-mail of the newly created user is verified or not. So if the values of the \"mail_verified\" and \"send_mail\" attributes are true, the user will receive an e-mail and they will be able to log on. On the other hand, if the \"send_mail\" is true, but the \"mail_verified\" is false, the user will receive e-mail and they have to click the verification link in the e-mail. If the attribute \"mail_verified\" is not configured, the default value is false.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"status": schema.StringAttribute{
						MarkdownDescription: "Specifies if the user is created as active, inactive or new. If the attribute \"active\" is not configured, the default value is inactive. " +
							utils.ValidValuesString(activeValues),
						Optional: true,
						Computed: true,
						Validators: []validator.String{
							stringvalidator.OneOf(activeValues...),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"custom_schemas": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Furthur enhance your user with custom schemas. The attribute must configured as a valid JSON string.",
				Validators: []validator.String{
					utils.ValidJSON(),
				},
			},
		},
	}
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan userData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	args, customSchemas, diags := getUserRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, _, err := r.cli.User.Create(ctx, customSchemas, args)
	if err != nil {
		resp.Diagnostics.AddError("Error creating user", fmt.Sprintf("%s", err))
		return
	}

	state, diags := userValueFrom(ctx, res, customSchemas)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// the initial password is not returned in the response, hence it must be read from the plan
	state.InitialPassword = plan.InitialPassword

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var config userData
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, customSchemasRes, err := r.cli.User.GetByUserId(ctx, config.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving user", fmt.Sprintf("%s", err))
		return
	}

	state, diags := userValueFrom(ctx, res, customSchemasRes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// the initial password is not returned in the response, hence it must be read from the state
	state.InitialPassword = config.InitialPassword

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan userData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state userData
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Id.IsNull() {
		resp.Diagnostics.AddError("User ID is missing", "Please provide a valid ID")
		return
	}

	args, customSchemas, diags := getUserRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	args.Id = state.Id.ValueString()

	res, _, err := r.cli.User.Update(ctx, customSchemas, args)
	if err != nil {
		resp.Diagnostics.AddError("Error updating application", fmt.Sprintf("%s", err))
		return
	}

	updatedState, diags := userValueFrom(ctx, res, customSchemas)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// the initial password is not returned in the response, hence it must be read from the plan
	updatedState.InitialPassword = plan.InitialPassword

	diags = resp.State.Set(ctx, &updatedState)
	resp.Diagnostics.Append(diags...)
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var config userData
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Id.IsNull() {
		resp.Diagnostics.AddError("User ID is missing", "Please provide a valid ID")
		return
	}

	err := r.cli.User.Delete(ctx, config.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Error deleting user", fmt.Sprintf("%s", err))
		return
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
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

	if plan.Name != nil {
		name := plan.Name

		if !name.FamilyName.IsNull() {
			args.Name.FamilyName = name.FamilyName.ValueString()
		}
		if !name.GivenName.IsNull() {
			args.Name.GivenName = name.GivenName.ValueString()
		}
		if !name.HonorificPrefix.IsNull() {
			args.Name.HonorificPrefix = name.HonorificPrefix.ValueString()
		}
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

		var sapExtensionUser sapExtensionUserData
		diags = plan.SapExtensionUser.As(ctx, &sapExtensionUser, basetypes.ObjectAsOptions{})
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return nil, "", diagnostics
		}

		if !sapExtensionUser.SendMail.IsNull() && !sapExtensionUser.SendMail.IsUnknown() {
			args.SAPExtension.SendMail = sapExtensionUser.SendMail.ValueBool()
		}

		if !sapExtensionUser.MailVerified.IsNull() && !sapExtensionUser.MailVerified.IsUnknown() {
			args.SAPExtension.MailVerified = sapExtensionUser.MailVerified.ValueBool()
		}

		if !sapExtensionUser.Status.IsNull() && !sapExtensionUser.Status.IsUnknown() {
			args.SAPExtension.Status = sapExtensionUser.Status.ValueString()
		}

	}

	var customSchemas string
	if !plan.CustomSchemas.IsNull() {
		customSchemas = plan.CustomSchemas.ValueString()
	}

	return args, customSchemas, diagnostics
}
