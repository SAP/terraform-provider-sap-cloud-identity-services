package provider

import (
	"context"
	"fmt"
	"terraform-provider-ias/internal/cli"
	"terraform-provider-ias/internal/cli/apiObjects/applications"
	"terraform-provider-ias/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

var sourceValues = []string{"Identity Directory", "Corporate Identity Provider", "Expression"}
var ssoValues = []string{"openIdConnect", "saml2"}
var usersTypeValues = []string{"public", "employee", "customer", "partner", "external", "onboardee"}
var subjectNameIdentifierFunctionValues = []string{"none", "upperCase", "lowerCase"}

func newApplicationResource() resource.Resource {
	return &applicationResource{}
}

type applicationResource struct {
	cli *cli.IasClient
}

func (d *applicationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.cli = req.ProviderData.(*cli.IasClient)
}

func (r *applicationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (r *applicationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Creates an application in the SAP Cloud Identity Services.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the application",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					utils.ValidUUID(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the application",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Free text description of the Application",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"parent_application_id": schema.StringAttribute{
				MarkdownDescription: "ID of the parent, from which the application will inherit its configurations",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					utils.ValidUUID(),
				},
			},
			"multi_tenant_app": schema.BoolAttribute{
				MarkdownDescription: "Only for Internal Use",
				Optional:            true,
				Computed:            true,
			},
			"global_account": schema.StringAttribute{
				// MarkdownDescription: "",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"authentication_schema": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"sso_type": schema.StringAttribute{
						MarkdownDescription: "The preferred protocol for the application. Acceptable values: \"openIdConnect\", \"saml2\"",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(ssoValues...),
						},
					},
					"subject_name_identifier": schema.SingleNestedAttribute{
						MarkdownDescription: "The attribute by which the application uses to identify the users. Identity Authentication sends the attribute to the application as subject in OpenID Connect tokens.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.Object{
							objectvalidator.AlsoRequires(
								path.MatchRoot("authentication_schema").AtName("subject_name_identifier").AtName("source"),
								path.MatchRoot("authentication_schema").AtName("subject_name_identifier").AtName("value"),
							),
						},
						Attributes: map[string]schema.Attribute{
							"source": schema.StringAttribute{
								MarkdownDescription: "Acceptable values: \"Identity Directory\", \"Corporate Idenity Provider\", \"Expression\"",
								Optional:            true,
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.OneOf(sourceValues...),
								},
							},
							"value": schema.StringAttribute{
								MarkdownDescription: "If the source is Identity Directory, the only acceptable values are \" none\", \"uid\", \"mail\", \"loginName\", \"displayName\", \"personnelNumber\", \"userUuid\"",
								Optional:            true,
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.LengthBetween(1, 255),
								},
							},
						},
					},
					"subject_name_identifier_function": schema.StringAttribute{
						MarkdownDescription: "Convert the subject name identifier to uppercase or lowercase. The only acceptable values are \"none\", \"upperCase\", \"lowerCase\"",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(subjectNameIdentifierFunctionValues...),
						},
					},
					"assertion_attributes": schema.ListNestedAttribute{
						MarkdownDescription: "User attributes to be sent to the application. The Source of these attributes is always the Identity Directory, thus only valid attribute values will be accepted.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.List{
							listvalidator.AlsoRequires(
								path.MatchRoot("authentication_schema").AtName("assertion_attributes").AtAnyListIndex().AtName("attribute_name"),
								path.MatchRoot("authentication_schema").AtName("assertion_attributes").AtAnyListIndex().AtName("attribute_value"),
							),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"attribute_name": schema.StringAttribute{
									MarkdownDescription: "Name of the attribute",
									Optional:            true,
									Computed:            true,
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 255),
									},
								},
								"attribute_value": schema.StringAttribute{
									MarkdownDescription: "Value of the attribute.",
									Optional:            true,
									Computed:            true,
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 255),
									},
								},
								"inherited": schema.BoolAttribute{
									MarkdownDescription: "Indicates whether the attribute has been inherited from a parent application.",
									Computed:            true,
								},
							},
						},
					},
					"advanced_assertion_attributes": schema.ListNestedAttribute{
						MarkdownDescription: "Identical to the assertion attributes, except that the assertion attributes can come from other Sources.",
						Optional:            true,
						Validators: []validator.List{
							listvalidator.AlsoRequires(
								path.MatchRoot("authentication_schema").AtName("advanced_assertion_attributes").AtAnyListIndex().AtName("source"),
								path.MatchRoot("authentication_schema").AtName("advanced_assertion_attributes").AtAnyListIndex().AtName("attribute_name"),
								path.MatchRoot("authentication_schema").AtName("advanced_assertion_attributes").AtAnyListIndex().AtName("attribute_value"),
							),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"source": schema.StringAttribute{
									MarkdownDescription: "Acceptable values: \"Corporate Idenity Provider\", \"Expression\"",
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.OneOf(sourceValues[1:]...),
									},
								},
								"attribute_name": schema.StringAttribute{
									MarkdownDescription: "Name of the attribute",
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 512),
									},
								},
								"attribute_value": schema.StringAttribute{
									MarkdownDescription: "Value of the attribute",
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 512),
									},
								},
								"inherited": schema.BoolAttribute{
									MarkdownDescription: "Indicates whether the attribute has been inherited from a parent application.",
									Computed:            true,
								},
							},
						},
					},
					"default_authenticating_idp": schema.StringAttribute{
						MarkdownDescription: "A default identity provider can be used for users with any user domain, group and type. This identity provider is used when none of the defined authentication rules meets the criteria.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 128),
						},
					},
					"authentication_rules": schema.ListNestedAttribute{
						MarkdownDescription: "Rules to manage authentication. Each rule is evaluated by priority until the criteria of a rule are fulfilled.",
						Optional:            true,
						Validators: []validator.List{
							listvalidator.AlsoRequires(
								path.MatchRoot("authentication_schema").AtName("authentication_rules").AtAnyListIndex().AtName("identity_provider_id"),
							),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"identity_provider_id": schema.StringAttribute{
									MarkdownDescription: "The identity provider to delegate authentication to when all the defined conditions are met.",
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 255),
									},
								},
								"user_type": schema.StringAttribute{
									MarkdownDescription: "The type of user to be authenticated.",
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.OneOf(usersTypeValues...),
										stringvalidator.AtLeastOneOf(
											path.MatchRoot("authentication_schema").AtName("authentication_rules").AtAnyListIndex().AtName("user_type"),
											path.MatchRoot("authentication_schema").AtName("authentication_rules").AtAnyListIndex().AtName("user_group"),
											path.MatchRoot("authentication_schema").AtName("authentication_rules").AtAnyListIndex().AtName("user_email_domain"),
											path.MatchRoot("authentication_schema").AtName("authentication_rules").AtAnyListIndex().AtName("ip_network_range"),
										),
									},
								},
								"user_group": schema.StringAttribute{
									MarkdownDescription: "The user group to be authenticated.",
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 255),
									},
								},
								"user_email_domain": schema.StringAttribute{
									MarkdownDescription: "Valid email domain to be authenticated.",
									Optional:            true,
									Validators: []validator.String{
										utils.ValidEmailDomain(),
									},
								},
								"ip_network_range": schema.StringAttribute{
									MarkdownDescription: "Valid IP range to be authenticated.",
									Optional:            true,
									Validators: []validator.String{
										utils.ValidIPAddress(),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *applicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var config applicationData

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	args, diags := getApplicationRequest(ctx, config)
	resp.Diagnostics.Append(diags...)

	res, _, err := r.cli.Application.Create(ctx, args)

	if err != nil {
		resp.Diagnostics.AddError("Error creating application", fmt.Sprintf("%s", err))
		return
	}

	state, diags := applicationValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *applicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var config applicationData
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	res, _, err := r.cli.Application.GetByAppId(ctx, config.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Error retrieving application", fmt.Sprintf("%s", err))
		return
	}

	state, diags := applicationValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *applicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var config applicationData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	// Retrieve the current state to get the existing application ID
	var state applicationData
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if state.Id.IsNull() {
		resp.Diagnostics.AddError("Application ID is missing", "Please provide a valid ID")
		return
	}

	// Update the application details
	args, diags := getApplicationRequest(ctx, config)
	resp.Diagnostics.Append(diags...)

	args.Id = state.Id.ValueString()

	res, _, err := r.cli.Application.Update(ctx, args)
	if err != nil {
		resp.Diagnostics.AddError("Error updating application", fmt.Sprintf("%s", err))
		return
	}
	updatedState, diags := applicationValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &updatedState)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *applicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var config applicationData
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if config.Id.IsNull() {
		resp.Diagnostics.AddError("Application ID is missing", "Please provide a valid ID")
		return
	}

	err := r.cli.Application.Delete(ctx, config.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Error deleting application", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *applicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getApplicationRequest(ctx context.Context, plan applicationData) (*applications.Application, diag.Diagnostics) {

	var diagnostics, diags diag.Diagnostics

	args := &applications.Application{
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		MultiTenantApp: plan.MultiTenantApp.ValueBool(),
		GlobalAccount:  plan.GlobalAccount.ValueString(),
	}

	diagnostics.Append(diags...)
	if !plan.ParentApplicationId.IsNull() {
		args.ParentApplicationId = plan.ParentApplicationId.ValueString()
	}

	if !plan.AuthenticationSchema.IsNull() && !plan.AuthenticationSchema.IsUnknown() {

		var authenticationSchema authenticationSchemaData
		diags = plan.AuthenticationSchema.As(ctx, &authenticationSchema, basetypes.ObjectAsOptions{})
		diagnostics.Append(diags...)

		if !authenticationSchema.SsoType.IsNull() {
			args.AuthenticationSchema.SsoType = authenticationSchema.SsoType.ValueString()
		}

		if authenticationSchema.SubjectNameIdentifier != nil && !authenticationSchema.SubjectNameIdentifier.Source.IsNull() {

			if authenticationSchema.SubjectNameIdentifier.Source.ValueString() == "Identity Directory" || authenticationSchema.SubjectNameIdentifier.Source.ValueString() == "Expression" {
				args.AuthenticationSchema.SubjectNameIdentifier = authenticationSchema.SubjectNameIdentifier.Value.ValueString()
			} else {
				args.AuthenticationSchema.SubjectNameIdentifier = "${corporateIdP." + authenticationSchema.SubjectNameIdentifier.Value.ValueString() + "}"
			}
		}

		if !authenticationSchema.SubjectNameIdentifierFunction.IsNull() {
			args.AuthenticationSchema.SubjectNameIdentifierFunction = authenticationSchema.SubjectNameIdentifierFunction.ValueString()
		}

		if !authenticationSchema.AssertionAttributes.IsNull() {

			var attributes []assertionAttributesData
			diags := authenticationSchema.AssertionAttributes.ElementsAs(ctx, &attributes, true)
			diagnostics.Append(diags...)

			for _, attribute := range attributes {

				assertionAttribute := applications.AssertionAttribute{
					AssertionAttributeName: attribute.AttributeName.ValueString(),
					UserAttributeName:      attribute.AttributeValue.ValueString(),
				}
				args.AuthenticationSchema.AssertionAttributes = append(args.AuthenticationSchema.AssertionAttributes, assertionAttribute)

			}
		}

		if !authenticationSchema.AdvancedAssertionAttributes.IsNull() {

			var advancedAssertionAttributes []advancedAssertionAttributesData
			diags := authenticationSchema.AdvancedAssertionAttributes.ElementsAs(ctx, &advancedAssertionAttributes, true)
			diagnostics.Append(diags...)

			for _, attribute := range advancedAssertionAttributes {

				if attribute.Source == types.StringValue("Corporate Identity Provider") {

					assertionAttribute := applications.AdvancedAssertionAttribute{
						AttributeName:  attribute.AttributeName.ValueString(),
						AttributeValue: "${corporateIdP." + attribute.AttributeValue.ValueString() + "}",
					}
					args.AuthenticationSchema.AdvancedAssertionAttributes = append(args.AuthenticationSchema.AdvancedAssertionAttributes, assertionAttribute)

				} else {

					assertionAttribute := applications.AdvancedAssertionAttribute{
						AttributeName:  attribute.AttributeName.ValueString(),
						AttributeValue: attribute.AttributeValue.ValueString(),
					}
					args.AuthenticationSchema.AdvancedAssertionAttributes = append(args.AuthenticationSchema.AdvancedAssertionAttributes, assertionAttribute)

				}
			}
		}

		if !authenticationSchema.AuthenticationRules.IsNull() {

			var authrules []authenticationRulesData
			diags = authenticationSchema.AuthenticationRules.ElementsAs(ctx, &authrules, true)
			diagnostics.Append(diags...)

			for _, rule := range authrules {

				authrule := applications.AuthenicationRule{
					UserType:           rule.UserType.ValueString(),
					UserGroup:          rule.UserGroup.ValueString(),
					UserEmailDomain:    rule.UserEmailDomain.ValueString(),
					IdentityProviderId: rule.IdentityProviderId.ValueString(),
					IpNetworkRange:     rule.IpNetworkRange.ValueString(),
				}

				args.AuthenticationSchema.ConditionalAuthentication = append(args.AuthenticationSchema.ConditionalAuthentication, authrule)
			}
		}
	}

	return args, diagnostics
}
