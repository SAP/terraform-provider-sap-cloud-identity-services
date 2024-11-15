package provider

import (
	"context"
	"fmt"
	"terraform-provider-ias/internal/cli"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"terraform-provider-ias/internal/cli/apiObjects/users"
)

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
				Computed: true,
				Validators: []validator.String{
					ValidUUID(),
				},
			},
			"schemas": schema.SetAttribute{
				ElementType: types.StringType,
				Required: true,
			},
			"user_name": schema.StringAttribute{
				Required: true,
			},
			"name": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"family_name": schema.StringAttribute{
						Required: true,
					},
					"given_name": schema.StringAttribute{
						Required: true,
					},
					"formatted": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"middle_name": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"honoric_prefix": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"honoric_suffix": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
				},
			},
			"emails": schema.SetNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"value": schema.StringAttribute{
							Required: true,
						},
						"type": schema.StringAttribute{
							Required: true,
						},
						"display": schema.StringAttribute{
							Computed: true,
							Optional: true,
						},
						"primary": schema.BoolAttribute{
							Computed: true,
							Optional: true,
						},
					},
				},
			},
			"password": schema.StringAttribute{
				Optional: true,
				Sensitive: true,
				//regex to check validity
			},
			"display_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"title": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"user_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"active": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"send_mail": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"mail_verified": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"status": schema.StringAttribute{
				Optional: true,
				Computed: true,
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

	res, err := r.cli.User.Create(ctx, args)

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

func (rs *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getUserRequest(ctx context.Context, plan userData) (*users.User, diag.Diagnostics){

	var diagnostics  diag.Diagnostics

	var schemas []string
	diags := plan.Schemas.ElementsAs(ctx, &schemas, true)
	diagnostics.Append(diags...)

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
		SAPExtension: users.SAPExtension{
			SendMail: plan.SendMail.ValueBool(),
			MailVerified: plan.MailVerified.ValueBool(),
			Status: plan.Status.ValueString(),
		},
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