package provider

import (
	"context"
	"fmt"
	"terraform-provider-ias/internal/cli"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"terraform-provider-ias/internal/cli/apiObjects/users"
)

var emailTypeValues = []string{"work", "home", "other"}
var userTypeValues  = []string{"public", "partner", "customer", "external", "onboardee", "employee"}
var activeValues    = []string{"active", "inactive", "new"}

func newUserResource() resource.Resource {
	return &userResource{}
}

type userResource struct {
	cli *cli.IasClient
}

func (d *userResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.cli = req.ProviderData.(*cli.IasClient)
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique ID of the resource.",
				Computed: true,
				Validators: []validator.String{
					ValidUUID(),
				},
			},
			"schemas": schema.SetAttribute{
				// MarkdownDescription: "",
				ElementType: types.StringType,
				Required: true,
			},
			"user_name": schema.StringAttribute{
				MarkdownDescription: "Unique user name of the user.",
				Required: true,
			},
			"name": schema.SingleNestedAttribute{
				MarkdownDescription: "Name of the user",
				Required: true,
				Attributes: map[string]schema.Attribute{
					"family_name": schema.StringAttribute{
						MarkdownDescription: "The following characters: <, >, : are not allowed.",
						Required: true,
					},
					"given_name": schema.StringAttribute{
						MarkdownDescription: "The following characters: <, >, : are not allowed.",
						Required: true,
					},
					"formatted": schema.StringAttribute{
						// MarkdownDescription: ,
						Optional: true,
						Computed: true,
					},
					"middle_name": schema.StringAttribute{
						// MarkdownDescription: ,
						Optional: true,
						Computed: true,
					},
					"honoric_prefix": schema.StringAttribute{
						// MarkdownDescription: ,
						Optional: true,
						Computed: true,
					},
					"honoric_suffix": schema.StringAttribute{
						// MarkdownDescription: ,
						Optional: true,
						Computed: true,
					},
				},
			},
			"emails": schema.SetNestedAttribute{
				MarkdownDescription: "Email of the user.",
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"value": schema.StringAttribute{
							MarkdownDescription: "Value of the email of the user.",
							Required: true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the email of the user.",
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf(emailTypeValues...),
							},
						},
						"display": schema.StringAttribute{
							// MarkdownDescription: "",
							Computed: true,
							Optional: true,
						},
						"primary": schema.BoolAttribute{
							MarkdownDescription: "Set the email to be primary",
							Computed: true,
							Optional: true,
						},
					},
				},
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password to be set for the user.",
				Optional: true,
				Sensitive: true,
				//regex to check validity, if check is added, add a test
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The name to be displayed for the user.",
				Optional: true,
				Computed: true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "The title to be given for the user.",
				Optional: true,
				Computed: true,
			},
			"user_type": schema.StringAttribute{
				MarkdownDescription: "Specifies the type of the user.The default type is \"public\".",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf(userTypeValues...),
				},
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Determines whether the user is active or not.The default value for the attribute is false.",
				Optional: true,
				Computed: true,
			},
			"sap_extension_user": schema.SingleNestedAttribute{
				// MarkdownDescription:
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"send_mail": schema.BoolAttribute{
						MarkdownDescription: "Specifies if an activation mail should be sent. The value of the attribute only matters when creating the user.",
						Optional: true,
						Computed: true,
					},
					"mail_verified": schema.BoolAttribute{
						MarkdownDescription: "The attribute specifies if the e-mail of the newly created user is verified or not. So if the values of the \"mail_verified\" and \"send_mail\" attributes are true, the user will receive e-mail and they will be able to log on. On the other hand, if the \"send_mail\" is true, but the \"mail_verified\" is false, the user will receive e-mail and they have to click the verification link in the e-mail. If the attribute \"verified\" is not passed in the request body, the default value of \"mail_erified\" is false.",
						Optional: true,
						Computed: true,
					},
					"status": schema.StringAttribute{
						MarkdownDescription: "Specifies if the user is created as active, inactive or new. If the attribute \"active\" is not passed in the request body, the default value of the attribute \"status\" is inactive.",
						Optional: true,
						Computed: true,
						Validators: []validator.String{
							stringvalidator.OneOf(activeValues...),
						},
					},
				},
			},
		},	
	}
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan userData

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	args, diags := getUserRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.cli.User.Create(ctx, args)
	if resp.Diagnostics.HasError() {
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Error creating user", fmt.Sprintf("%s", err))
		return
	}

	state, diags := userValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)

	state.Password = plan.Password
	state.Schemas = plan.Schemas

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var config userData
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	
	res, err := r.cli.User.GetByUserId(ctx, config.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Error retrieving user", fmt.Sprintf("%s", err))
		return
	}

	state, diags := userValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)

	state.Password = config.Password
	state.Schemas = config.Schemas

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	var state userData
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if state.Id.IsNull() {
		resp.Diagnostics.AddError("User ID is missing", "Please provide a valid ID")
		return
	}

	args, diags := getUserRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)

	args.Id = state.Id.ValueString()

	res, err := r.cli.User.Update(ctx, args)
	if err != nil {
		resp.Diagnostics.AddError("Error updating application", fmt.Sprintf("%s", err))
		return
	}

	updatedState, diags := userValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)

	updatedState.Password = plan.Password
	updatedState.Schemas = plan.Schemas

	diags = resp.State.Set(ctx, &updatedState)
	resp.Diagnostics.Append(diags...)
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var config userData
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if config.Id.IsNull() {
		resp.Diagnostics.AddError("User ID is missing", "Please provide a valid ID")
		return
	}

	err := r.cli.User.Delete(ctx, config.Id.ValueString())

	if err!=nil{
		resp.Diagnostics.AddError("Error deleting user", fmt.Sprintf("%s", err))
		return
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getUserRequest(ctx context.Context, plan userData) (*users.User, diag.Diagnostics){

	var diagnostics  diag.Diagnostics

	var schemas []string
	diags := plan.Schemas.ElementsAs(ctx, &schemas, true)
	diagnostics.Append(diags...)

	if len(schemas) == 0{
		diagnostics.AddError("The Schemas attribute cannot be Null or Empty","Provide a valid value for \"schemas\"")
		return nil, diagnostics
	}

	var name nameData
	diags = plan.Name.As(ctx, &name, basetypes.ObjectAsOptions{})
	diagnostics.Append(diags...)
	
	var emails []emailData
	diags = plan.Emails.ElementsAs(ctx, &emails, true)
	diagnostics.Append(diags...)

	args := &users.User{
		Schemas: schemas,
		UserName: plan.UserName.ValueString(),
		Name: users.Name{
			FamilyName: name.FamilyName.ValueString(),
			GivenName: name.GivenName.ValueString(),
			Formatted: name.Formatted.ValueString(),
			MiddleName: name.MiddleName.ValueString(),
			HonoricPrefix: name.HonoricPrefix.ValueString(),
			HonoricSuffix: name.HonoricSuffix.ValueString(),
		},
		DisplayName: plan.DisplayName.ValueString(),
		Password: plan.Password.ValueString(),
		Title: plan.Title.ValueString(),
		UserType: plan.UserType.ValueString(),
		Active: plan.Active.ValueBool(),
	}

	if !plan.SapExtensionUser.IsNull() && !plan.SapExtensionUser.IsUnknown() {

		var sapExtensionUser sapExtensionUserData
		diags = plan.SapExtensionUser.As(ctx, &sapExtensionUser, basetypes.ObjectAsOptions{})
		diagnostics.Append(diags...)

		args.SAPExtension.SendMail = sapExtensionUser.SendMail.ValueBool()
		args.SAPExtension.MailVerified = sapExtensionUser.MailVerified.ValueBool()
		args.SAPExtension.Status = sapExtensionUser.Status.ValueString()
	} 

	for _, email := range emails{
		userEmail := users.Email{
			Value: email.Value.ValueString(),
			Type: email.Type.ValueString(),
			Display: email.Display.ValueString(),
			Primary: email.Primary.ValueBool(),
		}
		args.Emails = append([]users.Email{userEmail},args.Emails...)
	}

	return args, diagnostics
}