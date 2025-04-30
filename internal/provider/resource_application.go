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

	// "github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	// "github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	// "github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

var (
	sourceValues                        = []string{"Identity Directory", "Corporate Identity Provider", "Expression"}
	ssoValues                           = []string{"openIdConnect", "saml2"}
	usersTypeValues                     = []string{"public", "employee", "customer", "partner", "external", "onboardee"}
	subjectNameIdentifierFunctionValues = []string{"none", "upperCase", "lowerCase"}
	actionValues                        = []string{"allow", "tfa", "deny", "rsaTfa", "smsTfa", "webTfa", "emailTfa"}
	groupTypeValues                     = []string{"cloud", "onPremise"}
	authMethodValues                    = []string{"cert", "spnego", "uidPw", "token", "socialIdentity", "trustedIdpSamlAssertion"}
)

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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"multi_tenant_app": schema.BoolAttribute{
				MarkdownDescription: "Only for Internal Use",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"global_account": schema.StringAttribute{
				// MarkdownDescription: "",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
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
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"source": schema.StringAttribute{
								MarkdownDescription: "Acceptable values: \"Identity Directory\", \"Corporate Idenity Provider\", \"Expression\"",
								Optional:            true,
								// Computed:            true,
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
							// listvalidator.SizeAtLeast(1),
						},
						PlanModifiers: []planmodifier.List{
							utils.UpdateUnknown(),
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
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								"attribute_value": schema.StringAttribute{
									MarkdownDescription: "Value of the attribute.",
									Optional:            true,
									Computed:            true,
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 255),
									},
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								"inherited": schema.BoolAttribute{
									MarkdownDescription: "Indicates whether the attribute has been inherited from a parent application.",
									Computed:            true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
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
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"conditional_authentication": schema.ListNestedAttribute{
						MarkdownDescription: "Define rules for authenticating identity provider according to email domain, user type, user group, and IP range. Each rule is evaluated by priority until the criteria of a rule are fulfilled.",
						Optional:            true,
						Validators: []validator.List{
							listvalidator.AlsoRequires(
								path.MatchRoot("authentication_schema").AtName("conditional_authentication").AtAnyListIndex().AtName("identity_provider_id"),
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
											path.MatchRoot("authentication_schema").AtName("conditional_authentication").AtAnyListIndex().AtName("user_group"),
											path.MatchRoot("authentication_schema").AtName("conditional_authentication").AtAnyListIndex().AtName("user_email_domain"),
											path.MatchRoot("authentication_schema").AtName("conditional_authentication").AtAnyListIndex().AtName("ip_network_range"),
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
					"risk_based_authentication": schema.SingleNestedAttribute{
						MarkdownDescription: "Define rules for authentication according to IP range, group membership, authentication method, and type of the authenticating user.",
						// the sub-attribute default_action can be Computed, hence the root-attribute is configured as both Optional and Computed
						// the sub-attribute rules is only configurable by the user, thus it is configured as Optional
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.Object{
							// objectplanmodifier.UseStateForUnknown(),
							utils.UpdateAction(),
						},
						Attributes: map[string]schema.Attribute{
							"default_action": schema.ListAttribute{
								MarkdownDescription: `
								Set a default action for any IP range, group and authentication method. This rule is used when none of the defined authentication rules meets the criteria. If there are no rules configured, the risk-based authentication configuration on tenant level will be applied.
								Valid values that can be configured: allow , tfa , deny, rsaTfa, smsTfa, webTfa, emailTfa
								`,
								// this field has a default value if not explicitly set by the user, hence it is both Optional and Computed
								Optional:    true,
								Computed:    true,
								ElementType: types.StringType,
								Validators: []validator.List{
									listvalidator.ValueStringsAre(
										stringvalidator.OneOf(actionValues...),
									),
									listvalidator.SizeAtLeast(1),
								},
								// PlanModifiers: []planmodifier.List{
									// listplanmodifier.UseStateForUnknown(),
									
								// },
								// Default: listdefault.StaticValue(
								// 	types.ListValueMust(types.StringType, []attr.Value{
								// 		types.StringValue("allow"),
								// 	}),
								// ),
							},
							"rules": schema.ListNestedAttribute{
								MarkdownDescription: "Configure rules to manage authentication. Each rule is evaluated by priority until the criteria of a rule are fulfilled.",
								Optional:            true,
								// Computed: 		  	 true,
								Validators: []validator.List{
									listvalidator.AlsoRequires(
										path.MatchRoot("authentication_schema").AtName("risk_based_authentication").AtName("rules").AtAnyListIndex().AtName("actions"),
									),
									listvalidator.SizeAtLeast(1),
								},
								// PlanModifiers: []planmodifier.List{
								// 	listplanmodifier.UseStateForUnknown(),
								// },
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"actions": schema.ListAttribute{
											MarkdownDescription: "Action for the authentication of the user when all conditions specified are met",
											Optional:            true,
											ElementType:         types.StringType,
											Validators: []validator.List{
												listvalidator.ValueStringsAre(
													stringvalidator.OneOf(actionValues...),
												),
											},
										},
										"ip_network_range": schema.StringAttribute{
											MarkdownDescription: "Valid IP range to be authenticated",
											Optional:            true,
											Validators: []validator.String{
												stringvalidator.AtLeastOneOf(
													path.MatchRoot("authentication_schema").AtName("risk_based_authentication").AtName("rules").AtAnyListIndex().AtName("ip_forward_range"),
													path.MatchRoot("authentication_schema").AtName("risk_based_authentication").AtName("rules").AtAnyListIndex().AtName("group"),
													path.MatchRoot("authentication_schema").AtName("risk_based_authentication").AtName("rules").AtAnyListIndex().AtName("group_type"),
													path.MatchRoot("authentication_schema").AtName("risk_based_authentication").AtName("rules").AtAnyListIndex().AtName("auth_method"),
													path.MatchRoot("authentication_schema").AtName("risk_based_authentication").AtName("rules").AtAnyListIndex().AtName("user_type"),
													path.MatchRoot("authentication_schema").AtName("risk_based_authentication").AtName("rules").AtAnyListIndex().AtName("corporate_idp_attribute"),
												),
												utils.ValidIPAddress(),
											},
										},
										"ip_forward_range": schema.StringAttribute{
											MarkdownDescription: "Valid IP range to be authenticated",
											Optional:            true,
											Validators: []validator.String{
												utils.ValidIPAddress(),
											},
										},
										"group": schema.StringAttribute{
											MarkdownDescription: "User group to be authenticated",
											Optional:            true,
											Validators: []validator.String{
												utils.ValidUUID(),
											},
										},
										"group_type": schema.StringAttribute{
											MarkdownDescription: "Type of the group to be authenticated",
											Optional:            true,
											Validators: []validator.String{
												stringvalidator.OneOf(groupTypeValues...),
											},
										},
										"auth_method": schema.StringAttribute{
											MarkdownDescription: "Authentication method to be authenticated",
											Optional:            true,
											Validators: []validator.String{
												stringvalidator.OneOf(authMethodValues...),
											},
										},
										"user_type": schema.StringAttribute{
											MarkdownDescription: "Type of the user to be authenticated",
											Optional:            true,
											Validators: []validator.String{
												stringvalidator.OneOf(usersTypeValues...),
											},
										},
										"corporate_idp_attribute": schema.SingleNestedAttribute{
											// MarkdownDescription: ,
											Optional: true,
											Validators: []validator.Object{
												objectvalidator.AlsoRequires(
													path.MatchRoot("authentication_schema").AtName("risk_based_authentication").AtName("rules").AtAnyListIndex().AtName("corporate_idp_attribute").AtName("name"),
													path.MatchRoot("authentication_schema").AtName("risk_based_authentication").AtName("rules").AtAnyListIndex().AtName("corporate_idp_attribute").AtName("value"),
												),
											},
											Attributes: map[string]schema.Attribute{
												"name": schema.StringAttribute{
													MarkdownDescription: "Name of the attribute",
													Optional:            true,
												},
												"value": schema.StringAttribute{
													MarkdownDescription: "Value of the attribute",
													Optional:            true,
												},
											},
										},
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

	if config.AuthenticationSchema != nil && config.AuthenticationSchema.RBAConfiguration != nil {
		diags = compareRbaRules(ctx, config.AuthenticationSchema.RBAConfiguration.Rules, &state.AuthenticationSchema.RBAConfiguration.Rules)
		resp.Diagnostics.Append(diags...)
	}

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

	if config.AuthenticationSchema != nil && config.AuthenticationSchema.RBAConfiguration != nil {
		diags = compareRbaRules(ctx, config.AuthenticationSchema.RBAConfiguration.Rules, &state.AuthenticationSchema.RBAConfiguration.Rules)
		resp.Diagnostics.Append(diags...)

		// state.AuthenticationSchema.RBAConfiguration.DefaultAction = config.AuthenticationSchema.RBAConfiguration.DefaultAction
	}

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

	if config.AuthenticationSchema != nil && config.AuthenticationSchema.RBAConfiguration != nil {
		diags = compareRbaRules(ctx, config.AuthenticationSchema.RBAConfiguration.Rules, &updatedState.AuthenticationSchema.RBAConfiguration.Rules)
		resp.Diagnostics.Append(diags...)

		// updatedState.AuthenticationSchema.RBAConfiguration.DefaultAction = state.AuthenticationSchema.RBAConfiguration.DefaultAction
	}

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

	if plan.AuthenticationSchema != nil {

		authenticationSchema := plan.AuthenticationSchema

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

			var attributes []applications.AssertionAttribute
			diags := authenticationSchema.AssertionAttributes.ElementsAs(ctx, &attributes, true)
			diagnostics.Append(diags...)

			args.AuthenticationSchema.AssertionAttributes = &attributes
		}

		if !authenticationSchema.AdvancedAssertionAttributes.IsNull() {

			var advancedAssertionAttributes []advancedAssertionAttributesData
			diags := authenticationSchema.AdvancedAssertionAttributes.ElementsAs(ctx, &advancedAssertionAttributes, true)
			diagnostics.Append(diags...)

			for _, attribute := range advancedAssertionAttributes {

				assertionAttribute := applications.AdvancedAssertionAttribute{
					AttributeName: attribute.AttributeName.ValueString(),
				}

				if attribute.Source == types.StringValue("Corporate Identity Provider") {
					assertionAttribute.AttributeValue = "${corporateIdP." + attribute.AttributeValue.ValueString() + "}"
				} else {
					assertionAttribute.AttributeValue = attribute.AttributeValue.ValueString()
				}

				args.AuthenticationSchema.AdvancedAssertionAttributes = append(args.AuthenticationSchema.AdvancedAssertionAttributes, assertionAttribute)
			}
		}

		if !authenticationSchema.AuthenticationRules.IsNull() {

			var authrules []applications.AuthenicationRule
			diags = authenticationSchema.AuthenticationRules.ElementsAs(ctx, &authrules, true)
			diagnostics.Append(diags...)

			args.AuthenticationSchema.ConditionalAuthentication = authrules
		}

		if authenticationSchema.RBAConfiguration != nil {

			var defaultActions []string
			diags = authenticationSchema.RBAConfiguration.DefaultAction.ElementsAs(ctx, &defaultActions, true)
			diagnostics.Append(diags...)

			var rbaRules []applications.RBARule
			diags = authenticationSchema.RBAConfiguration.Rules.ElementsAs(ctx, &rbaRules, true)
			diagnostics.Append(diags...)

			args.AuthenticationSchema.RiskBasedAuthentication = &applications.RBAConfiguration{
				DefaultAction: defaultActions,
				Rules:         rbaRules,
			}
		}
	}

	return args, diagnostics
}

func compareRbaRules(ctx context.Context, configRules types.List, stateRules *types.List) diag.Diagnostics {

	var diagnostics diag.Diagnostics

	var configRulesData []rbaRulesData
	diags := configRules.ElementsAs(ctx, &configRulesData, true)
	diagnostics.Append(diags...)

	var stateRulesData []rbaRulesData
	diags = stateRules.ElementsAs(ctx, &stateRulesData, true)
	diagnostics.Append(diags...)

	for i, rule := range configRulesData {

		if rule.IpNetworkRange.IsNull() && stateRulesData[i].IpNetworkRange.Equal(types.StringValue("")) {
			stateRulesData[i].IpNetworkRange = types.StringNull()
		}

		if rule.IpForwardRange.IsNull() && stateRulesData[i].IpForwardRange.Equal(types.StringValue("")) {
			stateRulesData[i].IpForwardRange = types.StringNull()
		}

		if rule.CorporateIdpAttribute == nil && stateRulesData[i].CorporateIdpAttribute != nil {
			stateRulesData[i].CorporateIdpAttribute = nil
		}
	}

	*stateRules, diags = types.ListValueFrom(ctx, rbaRuleObjType, stateRulesData)
	diagnostics.Append(diags...)

	return diagnostics

}
