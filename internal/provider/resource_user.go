package provider
// Attributes: map[string]schema.Attribute{
// 	"id": schema.StringAttribute{
// 		Required: true,
// 		Validators: []validator.String{
// 			ValidUUID(),
// 		},
// 	},
// 	"schemas": schema.ListAttribute{
// 		ElementType: types.StringType,
// 		Required: true,
// 	},
// 	"user_name": schema.StringAttribute{
// 		Required: true,
// 	},
// 	"name": schema.SingleNestedAttribute{
// 		Attributes: map[string]schema.Attribute{
// 			"family_name": schema.StringAttribute{
// 				Required: true,
// 			},
// 			"given_name": schema.StringAttribute{
// 				Required: true,
// 			},
// 			"formatted": schema.StringAttribute{
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"middle_name": schema.StringAttribute{
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"honoric_prefix": schema.StringAttribute{
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"honoric_suffix": schema.StringAttribute{
// 				Optional: true,
// 				Computed: true,
// 			},
// 		},
// 		Required: true,
// 	},
// 	"emails": schema.ListNestedAttribute{
// 		NestedObject: schema.NestedAttributeObject{
// 			Attributes: map[string]schema.Attribute{
// 				"value": schema.StringAttribute{
// 					Required: true,
// 				},
// 				"type": schema.StringAttribute{
// 					Required: true,
// 				},
// 				"display": schema.StringAttribute{
// 					Computed: true,
// 				},
// 				"primary": schema.BoolAttribute{
// 					Computed: true,
// 				},
// 			},
// 		},
// 		Required: true,
// 	},
// },	