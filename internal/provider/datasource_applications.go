package provider

import (
	"context"
	"fmt"
	"terraform-provider-ias/internal/cli"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func newApplicationsDataSource() datasource.DataSource {
	return &applicationsDataSource{}
}

type applicationsDataSource struct {
	cli *cli.IasClient
}

type applicationsData struct{
	// Id 		types.String	`tfsdk:"id"`
	Values 	types.List 		`tfsdk:"values"`
}

var subjectNameIdentitfierObjType = map[string]attr.Type{
	"source" : types.StringType,
	"value" : types.StringType,
}

var advancedAssertionAttributesObjType = types.ObjectType{
	AttrTypes : map[string]attr.Type {
		"source" : types.StringType,
		"attribute_name" : types.StringType,
		"attribute_value" : types.StringType,
		"inherited" : types.BoolType,
	},
}

var assertionAttributesObjType = types.ObjectType {
	AttrTypes: map[string]attr.Type{
		"attribute_name" : types.StringType,
		"attribute_value" : types.StringType,
		"inherited" : types.BoolType,
	},
}

var authenticationRulesObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"user_type" : types.StringType,
		"user_group" : types.StringType,
		"user_email_domain" : types.StringType,
		"identity_provider_id" : types.StringType,
		"ip_network_range" : types.StringType,
	},
}

var appObjType = types.ObjectType {
	AttrTypes: map[string]attr.Type{
		"id": types.StringType,
		"name": types.StringType,
		"description": types.StringType,
		"parent_application_id": types.StringType,
		"multi_tenant_app": types.BoolType,
		"global_account": types.StringType,
		"sso_type": types.StringType,
		"subject_name_identifier": types.ObjectType{
			AttrTypes: subjectNameIdentitfierObjType,
		},
		"assertion_attributes": types.ListType{
			ElemType:assertionAttributesObjType,
		},
		"advanced_assertion_attributes" : types.ListType{
			ElemType: advancedAssertionAttributesObjType,
		},
		"default_authenticating_idp" : types.StringType,
		"authentication_rules" : types.ListType{
			ElemType: authenticationRulesObjType,
		},
	},
}


func (d *applicationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.cli = req.ProviderData.(*cli.IasClient)
}

func (d *applicationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_applications"
}

func (d *applicationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{

		Attributes: map[string]schema.Attribute{
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Id of the application",
							Computed: true,
							Validators: []validator.String{
								ValidUUID(),
							},
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Free text description of the Application",
							Computed: true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description for the application",
							Computed: true,
						},
						"parent_application_id": schema.StringAttribute{
							MarkdownDescription: "ID of the parent, from which the application will inherit its configurations",
							Computed: true,
							Validators: []validator.String{
								ValidUUID(),
							},
						},
						"multi_tenant_app": schema.BoolAttribute{
							MarkdownDescription: "Only for Internal Use",
							Computed: true,
						},
						"global_account": schema.StringAttribute{
							// MarkdownDescription: "",
							Computed: true,
						},
						"sso_type": schema.StringAttribute{
							MarkdownDescription: "The preferred protocol for the application",
							Computed: true,
						},
						"subject_name_identifier" : schema.SingleNestedAttribute{
							MarkdownDescription: "The attribute by which the application uses to identify the users. Identity Authentication sends the attribute to the application as subject in OpenID Connect tokens.",
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"source": schema.StringAttribute{
									MarkdownDescription: "Acceptable values: \"Identity Directory\", \"Corporate Idenity Provider\", \"Expression\"",
									Computed: true,
								},
								"value": schema.StringAttribute{
									MarkdownDescription: "If the source is Identity Directory, the only acceptable values are \" none, uid, mail, loginName, displayName, personnelNumber, userUuid\"",
									Computed: true,
								},
							},
						},
						"assertion_attributes": schema.ListNestedAttribute{
							MarkdownDescription: "User attributes to be sent to the application. The Source of these attributes is always the Identity Directory, thus only valid attribute values will be accepted.",
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"attribute_name": schema.StringAttribute{
										MarkdownDescription: "Name of the attribute",
										Computed: true,
									},
									"attribute_value": schema.StringAttribute{
										MarkdownDescription: "Value of the attribute.",
										Computed: true,
									},
									"inherited": schema.BoolAttribute{
										MarkdownDescription: "Indicates whether the attribute has been inherited from a parent application.",
										Computed: true,
									},
								},
							},
						},
						"advanced_assertion_attributes" : schema.ListNestedAttribute{
							MarkdownDescription: "Identical to the assertion attributes, except that the assertion attributes can come from other Sources.",
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"source": schema.StringAttribute{
										MarkdownDescription: "Acceptable values: \"Corporate Idenity Provider\", \"Expression\"",
										Computed: true,
									},
									"attribute_name": schema.StringAttribute{
										MarkdownDescription: "Name of the attribute",
										Computed: true,
									},
									"attribute_value": schema.StringAttribute{
										MarkdownDescription: "Value of the attribute",
										Computed: true,
									},
									"inherited": schema.BoolAttribute{
										MarkdownDescription: "Indicates whether the attribute has been inherited from a parent application.",
										Computed: true,
									},
								},
							},
						},
						"default_authenticating_idp" : schema.StringAttribute{
							MarkdownDescription: "A default identity provider can be used for users with any user domain, group and type. This identity provider is used when none of the defined authentication rules meets the criteria.",
							Computed: true,
						},
						"authentication_rules": schema.ListNestedAttribute{
							MarkdownDescription: "Rules to manage authentication. Each rule is evaluated by priority until the criteria of a rule are fulfilled.",
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"identity_provider_id": schema.StringAttribute{
										MarkdownDescription: "The identity provider to delegate authentication to when all the defined conditions are met.",
										Computed: true,
									},
									"user_type": schema.StringAttribute{
										MarkdownDescription: "The type of user to be authenticated.",
										Computed: true,
									},
									"user_group": schema.StringAttribute{
										MarkdownDescription: "The user group to be authenticated.",
										Computed: true,
									},
									"user_email_domain": schema.StringAttribute{
										MarkdownDescription: "Valid email domain to be authenticated.",
										Computed: true,
									},
									"ip_network_range": schema.StringAttribute{
										MarkdownDescription: "Valid IP range to be authenticated.",
										Computed: true,
									},
								},
							},
						},
					},
				},
				Computed: true,
			},

		},
	}
}

func (d *applicationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var config applicationsData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	res, err := d.cli.ApplicationConfiguration.Application.Get(ctx)

	if err != nil {
		resp.Diagnostics.AddError("Error retrieving application", fmt.Sprintf("%s", err))
		return
	}

	resApps := applicationsValueFrom(ctx, res)

	config.Values, diags = types.ListValueFrom(ctx, appObjType, resApps)
	resp.Diagnostics.Append(diags...)
	
	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
