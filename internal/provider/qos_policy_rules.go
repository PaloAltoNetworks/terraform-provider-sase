package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	lNsAvVs "github.com/paloaltonetworks/sase-go/netsec/schema/qos/policy/rules"
	tzldypq "github.com/paloaltonetworks/sase-go/netsec/service/v1/qospolicyrules"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Data source (listing).
var (
	_ datasource.DataSource              = &qosPolicyRulesListDataSource{}
	_ datasource.DataSourceWithConfigure = &qosPolicyRulesListDataSource{}
)

func NewQosPolicyRulesListDataSource() datasource.DataSource {
	return &qosPolicyRulesListDataSource{}
}

type qosPolicyRulesListDataSource struct {
	client *sase.Client
}

type qosPolicyRulesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit    types.Int64  `tfsdk:"limit"`
	Offset   types.Int64  `tfsdk:"offset"`
	Name     types.String `tfsdk:"name"`
	Folder   types.String `tfsdk:"folder"`
	Position types.String `tfsdk:"position"`

	// Output.
	Data []qosPolicyRulesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type qosPolicyRulesListDsModelConfig struct {
	Action      qosPolicyRulesListDsModelActionObject   `tfsdk:"action"`
	Description types.String                            `tfsdk:"description"`
	DscpTos     *qosPolicyRulesListDsModelDscpTosObject `tfsdk:"dscp_tos"`
	ObjectId    types.String                            `tfsdk:"object_id"`
	Name        types.String                            `tfsdk:"name"`
	Schedule    types.String                            `tfsdk:"schedule"`
}

type qosPolicyRulesListDsModelActionObject struct {
	Class types.String `tfsdk:"class"`
}

type qosPolicyRulesListDsModelDscpTosObject struct {
	Codepoints []qosPolicyRulesListDsModelCodepointsObject `tfsdk:"codepoints"`
}

type qosPolicyRulesListDsModelCodepointsObject struct {
	Name types.String                         `tfsdk:"name"`
	Type *qosPolicyRulesListDsModelTypeObject `tfsdk:"type"`
}

type qosPolicyRulesListDsModelTypeObject struct {
	Af     *qosPolicyRulesListDsModelAfObject     `tfsdk:"af"`
	Cs     *qosPolicyRulesListDsModelCsObject     `tfsdk:"cs"`
	Custom *qosPolicyRulesListDsModelCustomObject `tfsdk:"custom"`
	Ef     types.Bool                             `tfsdk:"ef"`
	Tos    *qosPolicyRulesListDsModelTosObject    `tfsdk:"tos"`
}

type qosPolicyRulesListDsModelAfObject struct {
	Codepoint types.String `tfsdk:"codepoint"`
}

type qosPolicyRulesListDsModelCsObject struct {
	Codepoint types.String `tfsdk:"codepoint"`
}

type qosPolicyRulesListDsModelCustomObject struct {
	Codepoint *qosPolicyRulesListDsModelCodepointObject `tfsdk:"codepoint"`
}

type qosPolicyRulesListDsModelCodepointObject struct {
	BinaryValue   types.String `tfsdk:"binary_value"`
	CodepointName types.String `tfsdk:"codepoint_name"`
}

type qosPolicyRulesListDsModelTosObject struct {
	Codepoint types.String `tfsdk:"codepoint"`
}

// Metadata returns the data source type name.
func (d *qosPolicyRulesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_qos_policy_rules_list"
}

// Schema defines the schema for this listing data source.
func (d *qosPolicyRulesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Retrieves a listing of config items.",

		Attributes: map[string]dsschema.Attribute{
			"id": dsschema.StringAttribute{
				Description:         "The object ID.",
				MarkdownDescription: "The object ID.",
				Computed:            true,
			},

			// Input.
			"limit": dsschema.Int64Attribute{
				Description:         "The max count in result entry (count per page).",
				MarkdownDescription: "The max count in result entry (count per page).",
				Optional:            true,
				Computed:            true,
			},
			"offset": dsschema.Int64Attribute{
				Description:         "The offset of the result entry.",
				MarkdownDescription: "The offset of the result entry.",
				Optional:            true,
				Computed:            true,
			},
			"name": dsschema.StringAttribute{
				Description:         "The name of the entry.",
				MarkdownDescription: "The name of the entry.",
				Optional:            true,
			},
			"folder": dsschema.StringAttribute{
				Description:         "The folder of the entry. Value must be one of: `\"Shared\"`, `\"Mobile Users\"`, `\"Remote Networks\"`, `\"Service Connections\"`, `\"Mobile Users Container\"`, `\"Mobile Users Explicit Proxy\"`.",
				MarkdownDescription: "The folder of the entry. Value must be one of: `\"Shared\"`, `\"Mobile Users\"`, `\"Remote Networks\"`, `\"Service Connections\"`, `\"Mobile Users Container\"`, `\"Mobile Users Explicit Proxy\"`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Service Connections", "Mobile Users Container", "Mobile Users Explicit Proxy"),
				},
			},
			"position": dsschema.StringAttribute{
				Description:         "The position of a security rule. Value must be one of: `\"pre\"`, `\"post\"`.",
				MarkdownDescription: "The position of a security rule. Value must be one of: `\"pre\"`, `\"post\"`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("pre", "post"),
				},
			},

			// Output.
			"data": dsschema.ListNestedAttribute{
				Description:         "The `data` parameter.",
				MarkdownDescription: "The `data` parameter.",
				Computed:            true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"action": dsschema.SingleNestedAttribute{
							Description:         "The `action` parameter.",
							MarkdownDescription: "The `action` parameter.",
							Computed:            true,
							Attributes: map[string]dsschema.Attribute{
								"class": dsschema.StringAttribute{
									Description:         "The `class` parameter.",
									MarkdownDescription: "The `class` parameter.",
									Computed:            true,
								},
							},
						},
						"description": dsschema.StringAttribute{
							Description:         "The `description` parameter.",
							MarkdownDescription: "The `description` parameter.",
							Computed:            true,
						},
						"dscp_tos": dsschema.SingleNestedAttribute{
							Description:         "The `dscp_tos` parameter.",
							MarkdownDescription: "The `dscp_tos` parameter.",
							Computed:            true,
							Attributes: map[string]dsschema.Attribute{
								"codepoints": dsschema.ListNestedAttribute{
									Description:         "The `codepoints` parameter.",
									MarkdownDescription: "The `codepoints` parameter.",
									Computed:            true,
									NestedObject: dsschema.NestedAttributeObject{
										Attributes: map[string]dsschema.Attribute{
											"name": dsschema.StringAttribute{
												Description:         "The `name` parameter.",
												MarkdownDescription: "The `name` parameter.",
												Computed:            true,
											},
											"type": dsschema.SingleNestedAttribute{
												Description:         "The `type` parameter.",
												MarkdownDescription: "The `type` parameter.",
												Computed:            true,
												Attributes: map[string]dsschema.Attribute{
													"af": dsschema.SingleNestedAttribute{
														Description:         "The `af` parameter.",
														MarkdownDescription: "The `af` parameter.",
														Computed:            true,
														Attributes: map[string]dsschema.Attribute{
															"codepoint": dsschema.StringAttribute{
																Description:         "The `codepoint` parameter.",
																MarkdownDescription: "The `codepoint` parameter.",
																Computed:            true,
															},
														},
													},
													"cs": dsschema.SingleNestedAttribute{
														Description:         "The `cs` parameter.",
														MarkdownDescription: "The `cs` parameter.",
														Computed:            true,
														Attributes: map[string]dsschema.Attribute{
															"codepoint": dsschema.StringAttribute{
																Description:         "The `codepoint` parameter.",
																MarkdownDescription: "The `codepoint` parameter.",
																Computed:            true,
															},
														},
													},
													"custom": dsschema.SingleNestedAttribute{
														Description:         "The `custom` parameter.",
														MarkdownDescription: "The `custom` parameter.",
														Computed:            true,
														Attributes: map[string]dsschema.Attribute{
															"codepoint": dsschema.SingleNestedAttribute{
																Description:         "The `codepoint` parameter.",
																MarkdownDescription: "The `codepoint` parameter.",
																Computed:            true,
																Attributes: map[string]dsschema.Attribute{
																	"binary_value": dsschema.StringAttribute{
																		Description:         "The `binary_value` parameter.",
																		MarkdownDescription: "The `binary_value` parameter.",
																		Computed:            true,
																	},
																	"codepoint_name": dsschema.StringAttribute{
																		Description:         "The `codepoint_name` parameter.",
																		MarkdownDescription: "The `codepoint_name` parameter.",
																		Computed:            true,
																	},
																},
															},
														},
													},
													"ef": dsschema.BoolAttribute{
														Description:         "The `ef` parameter.",
														MarkdownDescription: "The `ef` parameter.",
														Computed:            true,
													},
													"tos": dsschema.SingleNestedAttribute{
														Description:         "The `tos` parameter.",
														MarkdownDescription: "The `tos` parameter.",
														Computed:            true,
														Attributes: map[string]dsschema.Attribute{
															"codepoint": dsschema.StringAttribute{
																Description:         "The `codepoint` parameter.",
																MarkdownDescription: "The `codepoint` parameter.",
																Computed:            true,
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
						"object_id": dsschema.StringAttribute{
							Description:         "The `object_id` parameter.",
							MarkdownDescription: "The `object_id` parameter.",
							Computed:            true,
						},
						"name": dsschema.StringAttribute{
							Description:         "The `name` parameter.",
							MarkdownDescription: "The `name` parameter.",
							Computed:            true,
						},
						"schedule": dsschema.StringAttribute{
							Description:         "The `schedule` parameter.",
							MarkdownDescription: "The `schedule` parameter.",
							Computed:            true,
						},
					},
				},
			},
			"total": dsschema.Int64Attribute{
				Description:         "The `total` parameter.",
				MarkdownDescription: "The `total` parameter.",
				Computed:            true,
			},
		},
	}
}

// Configure prepares the struct.
func (d *qosPolicyRulesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *qosPolicyRulesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state qosPolicyRulesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_qos_policy_rules_list",
		"terraform_provider_function": "Read",
		"limit":                       state.Limit.ValueInt64(),
		"has_limit":                   !state.Limit.IsNull(),
		"offset":                      state.Offset.ValueInt64(),
		"has_offset":                  !state.Offset.IsNull(),
		"name":                        state.Name.ValueString(),
		"has_name":                    !state.Name.IsNull(),
		"folder":                      state.Folder.ValueString(),
		"position":                    state.Position.ValueString(),
	})

	// Prepare to run the command.
	svc := tzldypq.NewClient(d.client)
	input := tzldypq.ListInput{
		Folder:   state.Folder.ValueString(),
		Position: state.Position.ValueString(),
	}
	if !state.Limit.IsNull() {
		input.Limit = api.Int(state.Limit.ValueInt64())
	}
	if !state.Offset.IsNull() {
		input.Offset = api.Int(state.Offset.ValueInt64())
	}
	if !state.Name.IsNull() {
		input.Name = api.String(state.Name.ValueString())
	}

	// Perform the operation.
	ans, err := svc.List(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error getting listing", err.Error())
		return
	}

	// Store the answer to state.
	var idBuilder strings.Builder
	if input.Limit != nil {
		idBuilder.WriteString(strconv.FormatInt(*input.Limit, 10))
	} else {
		idBuilder.WriteString("0")
	}
	idBuilder.WriteString(IdSeparator)
	if input.Offset != nil {
		idBuilder.WriteString(strconv.FormatInt(*input.Offset, 10))
	} else {
		idBuilder.WriteString("0")
	}
	idBuilder.WriteString(IdSeparator)
	if input.Name != nil {
		idBuilder.WriteString(*input.Name)
	}
	idBuilder.WriteString(IdSeparator)
	idBuilder.WriteString(input.Folder)
	idBuilder.WriteString(IdSeparator)
	idBuilder.WriteString(input.Position)
	state.Id = types.StringValue(idBuilder.String())
	var var0 []qosPolicyRulesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]qosPolicyRulesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 qosPolicyRulesListDsModelConfig
			var var3 qosPolicyRulesListDsModelActionObject
			var3.Class = types.StringValue(var1.Action.Class)
			var var4 *qosPolicyRulesListDsModelDscpTosObject
			if var1.DscpTos != nil {
				var4 = &qosPolicyRulesListDsModelDscpTosObject{}
				var var5 []qosPolicyRulesListDsModelCodepointsObject
				if len(var1.DscpTos.Codepoints) != 0 {
					var5 = make([]qosPolicyRulesListDsModelCodepointsObject, 0, len(var1.DscpTos.Codepoints))
					for var6Index := range var1.DscpTos.Codepoints {
						var6 := var1.DscpTos.Codepoints[var6Index]
						var var7 qosPolicyRulesListDsModelCodepointsObject
						var var8 *qosPolicyRulesListDsModelTypeObject
						if var6.Type != nil {
							var8 = &qosPolicyRulesListDsModelTypeObject{}
							var var9 *qosPolicyRulesListDsModelAfObject
							if var6.Type.Af != nil {
								var9 = &qosPolicyRulesListDsModelAfObject{}
								var9.Codepoint = types.StringValue(var6.Type.Af.Codepoint)
							}
							var var10 *qosPolicyRulesListDsModelCsObject
							if var6.Type.Cs != nil {
								var10 = &qosPolicyRulesListDsModelCsObject{}
								var10.Codepoint = types.StringValue(var6.Type.Cs.Codepoint)
							}
							var var11 *qosPolicyRulesListDsModelCustomObject
							if var6.Type.Custom != nil {
								var11 = &qosPolicyRulesListDsModelCustomObject{}
								var var12 *qosPolicyRulesListDsModelCodepointObject
								if var6.Type.Custom.Codepoint != nil {
									var12 = &qosPolicyRulesListDsModelCodepointObject{}
									var12.BinaryValue = types.StringValue(var6.Type.Custom.Codepoint.BinaryValue)
									var12.CodepointName = types.StringValue(var6.Type.Custom.Codepoint.CodepointName)
								}
								var11.Codepoint = var12
							}
							var var13 *qosPolicyRulesListDsModelTosObject
							if var6.Type.Tos != nil {
								var13 = &qosPolicyRulesListDsModelTosObject{}
								var13.Codepoint = types.StringValue(var6.Type.Tos.Codepoint)
							}
							var8.Af = var9
							var8.Cs = var10
							var8.Custom = var11
							if var6.Type.Ef != nil {
								var8.Ef = types.BoolValue(true)
							}
							var8.Tos = var13
						}
						var7.Name = types.StringValue(var6.Name)
						var7.Type = var8
						var5 = append(var5, var7)
					}
				}
				var4.Codepoints = var5
			}
			var2.Action = var3
			var2.Description = types.StringValue(var1.Description)
			var2.DscpTos = var4
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Name = types.StringValue(var1.Name)
			var2.Schedule = types.StringValue(var1.Schedule)
			var0 = append(var0, var2)
		}
	}
	state.Data = var0
	if !state.Limit.IsNull() {
		state.Limit = types.Int64Value(ans.Limit)
	}
	if !state.Offset.IsNull() {
		state.Offset = types.Int64Value(ans.Offset)
	}
	state.Total = types.Int64Value(ans.Total)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Data source.
var (
	_ datasource.DataSource              = &qosPolicyRulesDataSource{}
	_ datasource.DataSourceWithConfigure = &qosPolicyRulesDataSource{}
)

func NewQosPolicyRulesDataSource() datasource.DataSource {
	return &qosPolicyRulesDataSource{}
}

type qosPolicyRulesDataSource struct {
	client *sase.Client
}

type qosPolicyRulesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`
	Folder   types.String `tfsdk:"folder"`

	// Output.
	// Ref: #/components/schemas/qos-policy-rules
	Action      qosPolicyRulesDsModelActionObject   `tfsdk:"action"`
	Description types.String                        `tfsdk:"description"`
	DscpTos     *qosPolicyRulesDsModelDscpTosObject `tfsdk:"dscp_tos"`
	// input omit: ObjectId
	Name     types.String `tfsdk:"name"`
	Schedule types.String `tfsdk:"schedule"`
}

type qosPolicyRulesDsModelActionObject struct {
	Class types.String `tfsdk:"class"`
}

type qosPolicyRulesDsModelDscpTosObject struct {
	Codepoints []qosPolicyRulesDsModelCodepointsObject `tfsdk:"codepoints"`
}

type qosPolicyRulesDsModelCodepointsObject struct {
	Name types.String                     `tfsdk:"name"`
	Type *qosPolicyRulesDsModelTypeObject `tfsdk:"type"`
}

type qosPolicyRulesDsModelTypeObject struct {
	Af     *qosPolicyRulesDsModelAfObject     `tfsdk:"af"`
	Cs     *qosPolicyRulesDsModelCsObject     `tfsdk:"cs"`
	Custom *qosPolicyRulesDsModelCustomObject `tfsdk:"custom"`
	Ef     types.Bool                         `tfsdk:"ef"`
	Tos    *qosPolicyRulesDsModelTosObject    `tfsdk:"tos"`
}

type qosPolicyRulesDsModelAfObject struct {
	Codepoint types.String `tfsdk:"codepoint"`
}

type qosPolicyRulesDsModelCsObject struct {
	Codepoint types.String `tfsdk:"codepoint"`
}

type qosPolicyRulesDsModelCustomObject struct {
	Codepoint *qosPolicyRulesDsModelCodepointObject `tfsdk:"codepoint"`
}

type qosPolicyRulesDsModelCodepointObject struct {
	BinaryValue   types.String `tfsdk:"binary_value"`
	CodepointName types.String `tfsdk:"codepoint_name"`
}

type qosPolicyRulesDsModelTosObject struct {
	Codepoint types.String `tfsdk:"codepoint"`
}

// Metadata returns the data source type name.
func (d *qosPolicyRulesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_qos_policy_rules"
}

// Schema defines the schema for this listing data source.
func (d *qosPolicyRulesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Retrieves config for a specific item.",

		Attributes: map[string]dsschema.Attribute{
			"id": dsschema.StringAttribute{
				Description:         "The object ID.",
				MarkdownDescription: "The object ID.",
				Computed:            true,
			},

			// Input.
			"object_id": dsschema.StringAttribute{
				Description:         "The uuid of the resource.",
				MarkdownDescription: "The uuid of the resource.",
				Required:            true,
			},
			"folder": dsschema.StringAttribute{
				Description:         "The folder of the entry. Value must be one of: `\"Shared\"`, `\"Mobile Users\"`, `\"Remote Networks\"`, `\"Service Connections\"`, `\"Mobile Users Container\"`, `\"Mobile Users Explicit Proxy\"`.",
				MarkdownDescription: "The folder of the entry. Value must be one of: `\"Shared\"`, `\"Mobile Users\"`, `\"Remote Networks\"`, `\"Service Connections\"`, `\"Mobile Users Container\"`, `\"Mobile Users Explicit Proxy\"`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Service Connections", "Mobile Users Container", "Mobile Users Explicit Proxy"),
				},
			},

			// Output.
			"action": dsschema.SingleNestedAttribute{
				Description:         "The `action` parameter.",
				MarkdownDescription: "The `action` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"class": dsschema.StringAttribute{
						Description:         "The `class` parameter.",
						MarkdownDescription: "The `class` parameter.",
						Computed:            true,
					},
				},
			},
			"description": dsschema.StringAttribute{
				Description:         "The `description` parameter.",
				MarkdownDescription: "The `description` parameter.",
				Computed:            true,
			},
			"dscp_tos": dsschema.SingleNestedAttribute{
				Description:         "The `dscp_tos` parameter.",
				MarkdownDescription: "The `dscp_tos` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"codepoints": dsschema.ListNestedAttribute{
						Description:         "The `codepoints` parameter.",
						MarkdownDescription: "The `codepoints` parameter.",
						Computed:            true,
						NestedObject: dsschema.NestedAttributeObject{
							Attributes: map[string]dsschema.Attribute{
								"name": dsschema.StringAttribute{
									Description:         "The `name` parameter.",
									MarkdownDescription: "The `name` parameter.",
									Computed:            true,
								},
								"type": dsschema.SingleNestedAttribute{
									Description:         "The `type` parameter.",
									MarkdownDescription: "The `type` parameter.",
									Computed:            true,
									Attributes: map[string]dsschema.Attribute{
										"af": dsschema.SingleNestedAttribute{
											Description:         "The `af` parameter.",
											MarkdownDescription: "The `af` parameter.",
											Computed:            true,
											Attributes: map[string]dsschema.Attribute{
												"codepoint": dsschema.StringAttribute{
													Description:         "The `codepoint` parameter.",
													MarkdownDescription: "The `codepoint` parameter.",
													Computed:            true,
												},
											},
										},
										"cs": dsschema.SingleNestedAttribute{
											Description:         "The `cs` parameter.",
											MarkdownDescription: "The `cs` parameter.",
											Computed:            true,
											Attributes: map[string]dsschema.Attribute{
												"codepoint": dsschema.StringAttribute{
													Description:         "The `codepoint` parameter.",
													MarkdownDescription: "The `codepoint` parameter.",
													Computed:            true,
												},
											},
										},
										"custom": dsschema.SingleNestedAttribute{
											Description:         "The `custom` parameter.",
											MarkdownDescription: "The `custom` parameter.",
											Computed:            true,
											Attributes: map[string]dsschema.Attribute{
												"codepoint": dsschema.SingleNestedAttribute{
													Description:         "The `codepoint` parameter.",
													MarkdownDescription: "The `codepoint` parameter.",
													Computed:            true,
													Attributes: map[string]dsschema.Attribute{
														"binary_value": dsschema.StringAttribute{
															Description:         "The `binary_value` parameter.",
															MarkdownDescription: "The `binary_value` parameter.",
															Computed:            true,
														},
														"codepoint_name": dsschema.StringAttribute{
															Description:         "The `codepoint_name` parameter.",
															MarkdownDescription: "The `codepoint_name` parameter.",
															Computed:            true,
														},
													},
												},
											},
										},
										"ef": dsschema.BoolAttribute{
											Description:         "The `ef` parameter.",
											MarkdownDescription: "The `ef` parameter.",
											Computed:            true,
										},
										"tos": dsschema.SingleNestedAttribute{
											Description:         "The `tos` parameter.",
											MarkdownDescription: "The `tos` parameter.",
											Computed:            true,
											Attributes: map[string]dsschema.Attribute{
												"codepoint": dsschema.StringAttribute{
													Description:         "The `codepoint` parameter.",
													MarkdownDescription: "The `codepoint` parameter.",
													Computed:            true,
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
			"name": dsschema.StringAttribute{
				Description:         "The `name` parameter.",
				MarkdownDescription: "The `name` parameter.",
				Computed:            true,
			},
			"schedule": dsschema.StringAttribute{
				Description:         "The `schedule` parameter.",
				MarkdownDescription: "The `schedule` parameter.",
				Computed:            true,
			},
		},
	}
}

// Configure prepares the struct.
func (d *qosPolicyRulesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *qosPolicyRulesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state qosPolicyRulesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_qos_policy_rules",
		"object_id":                   state.ObjectId.ValueString(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := tzldypq.NewClient(d.client)
	input := tzldypq.ReadInput{
		ObjectId: state.ObjectId.ValueString(),
		Folder:   state.Folder.ValueString(),
	}

	// Perform the operation.
	ans, err := svc.Read(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error getting singleton", err.Error())
		return
	}

	// Store the answer to state.
	var idBuilder strings.Builder
	idBuilder.WriteString(input.ObjectId)
	idBuilder.WriteString(IdSeparator)
	idBuilder.WriteString(input.Folder)
	state.Id = types.StringValue(idBuilder.String())
	var var0 qosPolicyRulesDsModelActionObject
	var0.Class = types.StringValue(ans.Action.Class)
	var var1 *qosPolicyRulesDsModelDscpTosObject
	if ans.DscpTos != nil {
		var1 = &qosPolicyRulesDsModelDscpTosObject{}
		var var2 []qosPolicyRulesDsModelCodepointsObject
		if len(ans.DscpTos.Codepoints) != 0 {
			var2 = make([]qosPolicyRulesDsModelCodepointsObject, 0, len(ans.DscpTos.Codepoints))
			for var3Index := range ans.DscpTos.Codepoints {
				var3 := ans.DscpTos.Codepoints[var3Index]
				var var4 qosPolicyRulesDsModelCodepointsObject
				var var5 *qosPolicyRulesDsModelTypeObject
				if var3.Type != nil {
					var5 = &qosPolicyRulesDsModelTypeObject{}
					var var6 *qosPolicyRulesDsModelAfObject
					if var3.Type.Af != nil {
						var6 = &qosPolicyRulesDsModelAfObject{}
						var6.Codepoint = types.StringValue(var3.Type.Af.Codepoint)
					}
					var var7 *qosPolicyRulesDsModelCsObject
					if var3.Type.Cs != nil {
						var7 = &qosPolicyRulesDsModelCsObject{}
						var7.Codepoint = types.StringValue(var3.Type.Cs.Codepoint)
					}
					var var8 *qosPolicyRulesDsModelCustomObject
					if var3.Type.Custom != nil {
						var8 = &qosPolicyRulesDsModelCustomObject{}
						var var9 *qosPolicyRulesDsModelCodepointObject
						if var3.Type.Custom.Codepoint != nil {
							var9 = &qosPolicyRulesDsModelCodepointObject{}
							var9.BinaryValue = types.StringValue(var3.Type.Custom.Codepoint.BinaryValue)
							var9.CodepointName = types.StringValue(var3.Type.Custom.Codepoint.CodepointName)
						}
						var8.Codepoint = var9
					}
					var var10 *qosPolicyRulesDsModelTosObject
					if var3.Type.Tos != nil {
						var10 = &qosPolicyRulesDsModelTosObject{}
						var10.Codepoint = types.StringValue(var3.Type.Tos.Codepoint)
					}
					var5.Af = var6
					var5.Cs = var7
					var5.Custom = var8
					if var3.Type.Ef != nil {
						var5.Ef = types.BoolValue(true)
					}
					var5.Tos = var10
				}
				var4.Name = types.StringValue(var3.Name)
				var4.Type = var5
				var2 = append(var2, var4)
			}
		}
		var1.Codepoints = var2
	}
	state.Action = var0
	state.Description = types.StringValue(ans.Description)
	state.DscpTos = var1
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.Schedule = types.StringValue(ans.Schedule)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &qosPolicyRulesResource{}
	_ resource.ResourceWithConfigure   = &qosPolicyRulesResource{}
	_ resource.ResourceWithImportState = &qosPolicyRulesResource{}
)

func NewQosPolicyRulesResource() resource.Resource {
	return &qosPolicyRulesResource{}
}

type qosPolicyRulesResource struct {
	client *sase.Client
}

type qosPolicyRulesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder   types.String `tfsdk:"folder"`
	Position types.String `tfsdk:"position"`

	// Request body input.
	// Ref: #/components/schemas/qos-policy-rules
	Action      qosPolicyRulesRsModelActionObject   `tfsdk:"action"`
	Description types.String                        `tfsdk:"description"`
	DscpTos     *qosPolicyRulesRsModelDscpTosObject `tfsdk:"dscp_tos"`
	ObjectId    types.String                        `tfsdk:"object_id"`
	Name        types.String                        `tfsdk:"name"`
	Schedule    types.String                        `tfsdk:"schedule"`
}

type qosPolicyRulesRsModelActionObject struct {
	Class types.String `tfsdk:"class"`
}

type qosPolicyRulesRsModelDscpTosObject struct {
	Codepoints []qosPolicyRulesRsModelCodepointsObject `tfsdk:"codepoints"`
}

type qosPolicyRulesRsModelCodepointsObject struct {
	Name types.String                     `tfsdk:"name"`
	Type *qosPolicyRulesRsModelTypeObject `tfsdk:"type"`
}

type qosPolicyRulesRsModelTypeObject struct {
	Af     *qosPolicyRulesRsModelAfObject     `tfsdk:"af"`
	Cs     *qosPolicyRulesRsModelCsObject     `tfsdk:"cs"`
	Custom *qosPolicyRulesRsModelCustomObject `tfsdk:"custom"`
	Ef     types.Bool                         `tfsdk:"ef"`
	Tos    *qosPolicyRulesRsModelTosObject    `tfsdk:"tos"`
}

type qosPolicyRulesRsModelAfObject struct {
	Codepoint types.String `tfsdk:"codepoint"`
}

type qosPolicyRulesRsModelCsObject struct {
	Codepoint types.String `tfsdk:"codepoint"`
}

type qosPolicyRulesRsModelCustomObject struct {
	Codepoint *qosPolicyRulesRsModelCodepointObject `tfsdk:"codepoint"`
}

type qosPolicyRulesRsModelCodepointObject struct {
	BinaryValue   types.String `tfsdk:"binary_value"`
	CodepointName types.String `tfsdk:"codepoint_name"`
}

type qosPolicyRulesRsModelTosObject struct {
	Codepoint types.String `tfsdk:"codepoint"`
}

// Metadata returns the data source type name.
func (r *qosPolicyRulesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_qos_policy_rules"
}

// Schema defines the schema for this listing data source.
func (r *qosPolicyRulesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = rsschema.Schema{
		Description: "Retrieves config for a specific item.",

		Attributes: map[string]rsschema.Attribute{
			"id": rsschema.StringAttribute{
				Description:         "The object ID.",
				MarkdownDescription: "The object ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			// Input.
			"folder": rsschema.StringAttribute{
				Description:         "The folder of the entry. Value must be one of: `\"Shared\"`, `\"Mobile Users\"`, `\"Remote Networks\"`, `\"Service Connections\"`, `\"Mobile Users Container\"`, `\"Mobile Users Explicit Proxy\"`.",
				MarkdownDescription: "The folder of the entry. Value must be one of: `\"Shared\"`, `\"Mobile Users\"`, `\"Remote Networks\"`, `\"Service Connections\"`, `\"Mobile Users Container\"`, `\"Mobile Users Explicit Proxy\"`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Service Connections", "Mobile Users Container", "Mobile Users Explicit Proxy"),
				},
			},
			"position": rsschema.StringAttribute{
				Description:         "The position of a security rule. Value must be one of: `\"pre\"`, `\"post\"`.",
				MarkdownDescription: "The position of a security rule. Value must be one of: `\"pre\"`, `\"post\"`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("pre", "post"),
				},
			},

			"action": rsschema.SingleNestedAttribute{
				Description:         "The `action` parameter.",
				MarkdownDescription: "The `action` parameter.",
				Required:            true,
				Attributes: map[string]rsschema.Attribute{
					"class": rsschema.StringAttribute{
						Description:         "The `class` parameter.",
						MarkdownDescription: "The `class` parameter.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							DefaultString(""),
						},
					},
				},
			},
			"description": rsschema.StringAttribute{
				Description:         "The `description` parameter.",
				MarkdownDescription: "The `description` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
			},
			"dscp_tos": rsschema.SingleNestedAttribute{
				Description:         "The `dscp_tos` parameter.",
				MarkdownDescription: "The `dscp_tos` parameter.",
				Optional:            true,
				Attributes: map[string]rsschema.Attribute{
					"codepoints": rsschema.ListNestedAttribute{
						Description:         "The `codepoints` parameter.",
						MarkdownDescription: "The `codepoints` parameter.",
						Optional:            true,
						NestedObject: rsschema.NestedAttributeObject{
							Attributes: map[string]rsschema.Attribute{
								"name": rsschema.StringAttribute{
									Description:         "The `name` parameter.",
									MarkdownDescription: "The `name` parameter.",
									Optional:            true,
									Computed:            true,
									PlanModifiers: []planmodifier.String{
										DefaultString(""),
									},
								},
								"type": rsschema.SingleNestedAttribute{
									Description:         "The `type` parameter.",
									MarkdownDescription: "The `type` parameter.",
									Optional:            true,
									Attributes: map[string]rsschema.Attribute{
										"af": rsschema.SingleNestedAttribute{
											Description:         "The `af` parameter.",
											MarkdownDescription: "The `af` parameter.",
											Optional:            true,
											Attributes: map[string]rsschema.Attribute{
												"codepoint": rsschema.StringAttribute{
													Description:         "The `codepoint` parameter.",
													MarkdownDescription: "The `codepoint` parameter.",
													Optional:            true,
													Computed:            true,
													PlanModifiers: []planmodifier.String{
														DefaultString(""),
													},
												},
											},
										},
										"cs": rsschema.SingleNestedAttribute{
											Description:         "The `cs` parameter.",
											MarkdownDescription: "The `cs` parameter.",
											Optional:            true,
											Attributes: map[string]rsschema.Attribute{
												"codepoint": rsschema.StringAttribute{
													Description:         "The `codepoint` parameter.",
													MarkdownDescription: "The `codepoint` parameter.",
													Optional:            true,
													Computed:            true,
													PlanModifiers: []planmodifier.String{
														DefaultString(""),
													},
												},
											},
										},
										"custom": rsschema.SingleNestedAttribute{
											Description:         "The `custom` parameter.",
											MarkdownDescription: "The `custom` parameter.",
											Optional:            true,
											Attributes: map[string]rsschema.Attribute{
												"codepoint": rsschema.SingleNestedAttribute{
													Description:         "The `codepoint` parameter.",
													MarkdownDescription: "The `codepoint` parameter.",
													Optional:            true,
													Attributes: map[string]rsschema.Attribute{
														"binary_value": rsschema.StringAttribute{
															Description:         "The `binary_value` parameter.",
															MarkdownDescription: "The `binary_value` parameter.",
															Optional:            true,
															Computed:            true,
															PlanModifiers: []planmodifier.String{
																DefaultString(""),
															},
														},
														"codepoint_name": rsschema.StringAttribute{
															Description:         "The `codepoint_name` parameter.",
															MarkdownDescription: "The `codepoint_name` parameter.",
															Optional:            true,
															Computed:            true,
															PlanModifiers: []planmodifier.String{
																DefaultString(""),
															},
														},
													},
												},
											},
										},
										"ef": rsschema.BoolAttribute{
											Description:         "The `ef` parameter.",
											MarkdownDescription: "The `ef` parameter.",
											Optional:            true,
										},
										"tos": rsschema.SingleNestedAttribute{
											Description:         "The `tos` parameter.",
											MarkdownDescription: "The `tos` parameter.",
											Optional:            true,
											Attributes: map[string]rsschema.Attribute{
												"codepoint": rsschema.StringAttribute{
													Description:         "The `codepoint` parameter.",
													MarkdownDescription: "The `codepoint` parameter.",
													Optional:            true,
													Computed:            true,
													PlanModifiers: []planmodifier.String{
														DefaultString(""),
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
			"object_id": rsschema.StringAttribute{
				Description:         "The `object_id` parameter.",
				MarkdownDescription: "The `object_id` parameter.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": rsschema.StringAttribute{
				Description:         "The `name` parameter.",
				MarkdownDescription: "The `name` parameter.",
				Required:            true,
			},
			"schedule": rsschema.StringAttribute{
				Description:         "The `schedule` parameter.",
				MarkdownDescription: "The `schedule` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *qosPolicyRulesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *qosPolicyRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state qosPolicyRulesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_qos_policy_rules",
		"folder":                      state.Folder.ValueString(),
		"position":                    state.Position.ValueString(),
	})

	// Prepare to create the config.
	svc := tzldypq.NewClient(r.client)
	input := tzldypq.CreateInput{
		Folder:   state.Folder.ValueString(),
		Position: state.Position.ValueString(),
	}
	var var0 lNsAvVs.Config
	var var1 lNsAvVs.ActionObject
	var1.Class = state.Action.Class.ValueString()
	var0.Action = var1
	var0.Description = state.Description.ValueString()
	var var2 *lNsAvVs.DscpTosObject
	if state.DscpTos != nil {
		var2 = &lNsAvVs.DscpTosObject{}
		var var3 []lNsAvVs.CodepointsObject
		if len(state.DscpTos.Codepoints) != 0 {
			var3 = make([]lNsAvVs.CodepointsObject, 0, len(state.DscpTos.Codepoints))
			for var4Index := range state.DscpTos.Codepoints {
				var4 := state.DscpTos.Codepoints[var4Index]
				var var5 lNsAvVs.CodepointsObject
				var5.Name = var4.Name.ValueString()
				var var6 *lNsAvVs.TypeObject
				if var4.Type != nil {
					var6 = &lNsAvVs.TypeObject{}
					var var7 *lNsAvVs.AfObject
					if var4.Type.Af != nil {
						var7 = &lNsAvVs.AfObject{}
						var7.Codepoint = var4.Type.Af.Codepoint.ValueString()
					}
					var6.Af = var7
					var var8 *lNsAvVs.CsObject
					if var4.Type.Cs != nil {
						var8 = &lNsAvVs.CsObject{}
						var8.Codepoint = var4.Type.Cs.Codepoint.ValueString()
					}
					var6.Cs = var8
					var var9 *lNsAvVs.CustomObject
					if var4.Type.Custom != nil {
						var9 = &lNsAvVs.CustomObject{}
						var var10 *lNsAvVs.CodepointObject
						if var4.Type.Custom.Codepoint != nil {
							var10 = &lNsAvVs.CodepointObject{}
							var10.BinaryValue = var4.Type.Custom.Codepoint.BinaryValue.ValueString()
							var10.CodepointName = var4.Type.Custom.Codepoint.CodepointName.ValueString()
						}
						var9.Codepoint = var10
					}
					var6.Custom = var9
					if var4.Type.Ef.ValueBool() {
						var6.Ef = struct{}{}
					}
					var var11 *lNsAvVs.TosObject
					if var4.Type.Tos != nil {
						var11 = &lNsAvVs.TosObject{}
						var11.Codepoint = var4.Type.Tos.Codepoint.ValueString()
					}
					var6.Tos = var11
				}
				var5.Type = var6
				var3 = append(var3, var5)
			}
		}
		var2.Codepoints = var3
	}
	var0.DscpTos = var2
	var0.Name = state.Name.ValueString()
	var0.Schedule = state.Schedule.ValueString()
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Create(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in create", err.Error())
		return
	}

	// Store the answer to state.
	var idBuilder strings.Builder
	idBuilder.WriteString(input.Folder)
	idBuilder.WriteString(IdSeparator)
	idBuilder.WriteString(input.Position)
	idBuilder.WriteString(IdSeparator)
	idBuilder.WriteString(ans.ObjectId)
	state.Id = types.StringValue(idBuilder.String())
	var var12 qosPolicyRulesRsModelActionObject
	var12.Class = types.StringValue(ans.Action.Class)
	var var13 *qosPolicyRulesRsModelDscpTosObject
	if ans.DscpTos != nil {
		var13 = &qosPolicyRulesRsModelDscpTosObject{}
		var var14 []qosPolicyRulesRsModelCodepointsObject
		if len(ans.DscpTos.Codepoints) != 0 {
			var14 = make([]qosPolicyRulesRsModelCodepointsObject, 0, len(ans.DscpTos.Codepoints))
			for var15Index := range ans.DscpTos.Codepoints {
				var15 := ans.DscpTos.Codepoints[var15Index]
				var var16 qosPolicyRulesRsModelCodepointsObject
				var var17 *qosPolicyRulesRsModelTypeObject
				if var15.Type != nil {
					var17 = &qosPolicyRulesRsModelTypeObject{}
					var var18 *qosPolicyRulesRsModelAfObject
					if var15.Type.Af != nil {
						var18 = &qosPolicyRulesRsModelAfObject{}
						var18.Codepoint = types.StringValue(var15.Type.Af.Codepoint)
					}
					var var19 *qosPolicyRulesRsModelCsObject
					if var15.Type.Cs != nil {
						var19 = &qosPolicyRulesRsModelCsObject{}
						var19.Codepoint = types.StringValue(var15.Type.Cs.Codepoint)
					}
					var var20 *qosPolicyRulesRsModelCustomObject
					if var15.Type.Custom != nil {
						var20 = &qosPolicyRulesRsModelCustomObject{}
						var var21 *qosPolicyRulesRsModelCodepointObject
						if var15.Type.Custom.Codepoint != nil {
							var21 = &qosPolicyRulesRsModelCodepointObject{}
							var21.BinaryValue = types.StringValue(var15.Type.Custom.Codepoint.BinaryValue)
							var21.CodepointName = types.StringValue(var15.Type.Custom.Codepoint.CodepointName)
						}
						var20.Codepoint = var21
					}
					var var22 *qosPolicyRulesRsModelTosObject
					if var15.Type.Tos != nil {
						var22 = &qosPolicyRulesRsModelTosObject{}
						var22.Codepoint = types.StringValue(var15.Type.Tos.Codepoint)
					}
					var17.Af = var18
					var17.Cs = var19
					var17.Custom = var20
					if var15.Type.Ef != nil {
						var17.Ef = types.BoolValue(true)
					}
					var17.Tos = var22
				}
				var16.Name = types.StringValue(var15.Name)
				var16.Type = var17
				var14 = append(var14, var16)
			}
		}
		var13.Codepoints = var14
	}
	state.Action = var12
	state.Description = types.StringValue(ans.Description)
	state.DscpTos = var13
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.Schedule = types.StringValue(ans.Schedule)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *qosPolicyRulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var idType types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &idType)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := idType.ValueString()
	tokens := strings.Split(id, IdSeparator)
	if len(tokens) != 3 {
		resp.Diagnostics.AddError("Error in resource ID format", "Expected 3 tokens")
		return
	}

	var state qosPolicyRulesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_qos_policy_rules",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 2, "Position": 1},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := tzldypq.NewClient(r.client)
	input := tzldypq.ReadInput{
		ObjectId: tokens[2],
		Folder:   tokens[0],
	}

	// Perform the operation.
	ans, err := svc.Read(ctx, input)
	if err != nil {
		if IsObjectNotFound(err) {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError("Error reading config", err.Error())
		}
		return
	}

	// Store the answer to state.
	state.Folder = types.StringValue(tokens[0])
	state.Position = types.StringValue(tokens[1])
	state.Id = idType
	var var0 qosPolicyRulesRsModelActionObject
	var0.Class = types.StringValue(ans.Action.Class)
	var var1 *qosPolicyRulesRsModelDscpTosObject
	if ans.DscpTos != nil {
		var1 = &qosPolicyRulesRsModelDscpTosObject{}
		var var2 []qosPolicyRulesRsModelCodepointsObject
		if len(ans.DscpTos.Codepoints) != 0 {
			var2 = make([]qosPolicyRulesRsModelCodepointsObject, 0, len(ans.DscpTos.Codepoints))
			for var3Index := range ans.DscpTos.Codepoints {
				var3 := ans.DscpTos.Codepoints[var3Index]
				var var4 qosPolicyRulesRsModelCodepointsObject
				var var5 *qosPolicyRulesRsModelTypeObject
				if var3.Type != nil {
					var5 = &qosPolicyRulesRsModelTypeObject{}
					var var6 *qosPolicyRulesRsModelAfObject
					if var3.Type.Af != nil {
						var6 = &qosPolicyRulesRsModelAfObject{}
						var6.Codepoint = types.StringValue(var3.Type.Af.Codepoint)
					}
					var var7 *qosPolicyRulesRsModelCsObject
					if var3.Type.Cs != nil {
						var7 = &qosPolicyRulesRsModelCsObject{}
						var7.Codepoint = types.StringValue(var3.Type.Cs.Codepoint)
					}
					var var8 *qosPolicyRulesRsModelCustomObject
					if var3.Type.Custom != nil {
						var8 = &qosPolicyRulesRsModelCustomObject{}
						var var9 *qosPolicyRulesRsModelCodepointObject
						if var3.Type.Custom.Codepoint != nil {
							var9 = &qosPolicyRulesRsModelCodepointObject{}
							var9.BinaryValue = types.StringValue(var3.Type.Custom.Codepoint.BinaryValue)
							var9.CodepointName = types.StringValue(var3.Type.Custom.Codepoint.CodepointName)
						}
						var8.Codepoint = var9
					}
					var var10 *qosPolicyRulesRsModelTosObject
					if var3.Type.Tos != nil {
						var10 = &qosPolicyRulesRsModelTosObject{}
						var10.Codepoint = types.StringValue(var3.Type.Tos.Codepoint)
					}
					var5.Af = var6
					var5.Cs = var7
					var5.Custom = var8
					if var3.Type.Ef != nil {
						var5.Ef = types.BoolValue(true)
					}
					var5.Tos = var10
				}
				var4.Name = types.StringValue(var3.Name)
				var4.Type = var5
				var2 = append(var2, var4)
			}
		}
		var1.Codepoints = var2
	}
	state.Action = var0
	state.Description = types.StringValue(ans.Description)
	state.DscpTos = var1
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.Schedule = types.StringValue(ans.Schedule)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *qosPolicyRulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state qosPolicyRulesRsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource update", map[string]any{
		"terraform_provider_function": "Update",
		"resource_name":               "sase_qos_policy_rules",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := tzldypq.NewClient(r.client)
	input := tzldypq.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 lNsAvVs.Config
	var var1 lNsAvVs.ActionObject
	var1.Class = plan.Action.Class.ValueString()
	var0.Action = var1
	var0.Description = plan.Description.ValueString()
	var var2 *lNsAvVs.DscpTosObject
	if plan.DscpTos != nil {
		var2 = &lNsAvVs.DscpTosObject{}
		var var3 []lNsAvVs.CodepointsObject
		if len(plan.DscpTos.Codepoints) != 0 {
			var3 = make([]lNsAvVs.CodepointsObject, 0, len(plan.DscpTos.Codepoints))
			for var4Index := range plan.DscpTos.Codepoints {
				var4 := plan.DscpTos.Codepoints[var4Index]
				var var5 lNsAvVs.CodepointsObject
				var5.Name = var4.Name.ValueString()
				var var6 *lNsAvVs.TypeObject
				if var4.Type != nil {
					var6 = &lNsAvVs.TypeObject{}
					var var7 *lNsAvVs.AfObject
					if var4.Type.Af != nil {
						var7 = &lNsAvVs.AfObject{}
						var7.Codepoint = var4.Type.Af.Codepoint.ValueString()
					}
					var6.Af = var7
					var var8 *lNsAvVs.CsObject
					if var4.Type.Cs != nil {
						var8 = &lNsAvVs.CsObject{}
						var8.Codepoint = var4.Type.Cs.Codepoint.ValueString()
					}
					var6.Cs = var8
					var var9 *lNsAvVs.CustomObject
					if var4.Type.Custom != nil {
						var9 = &lNsAvVs.CustomObject{}
						var var10 *lNsAvVs.CodepointObject
						if var4.Type.Custom.Codepoint != nil {
							var10 = &lNsAvVs.CodepointObject{}
							var10.BinaryValue = var4.Type.Custom.Codepoint.BinaryValue.ValueString()
							var10.CodepointName = var4.Type.Custom.Codepoint.CodepointName.ValueString()
						}
						var9.Codepoint = var10
					}
					var6.Custom = var9
					if var4.Type.Ef.ValueBool() {
						var6.Ef = struct{}{}
					}
					var var11 *lNsAvVs.TosObject
					if var4.Type.Tos != nil {
						var11 = &lNsAvVs.TosObject{}
						var11.Codepoint = var4.Type.Tos.Codepoint.ValueString()
					}
					var6.Tos = var11
				}
				var5.Type = var6
				var3 = append(var3, var5)
			}
		}
		var2.Codepoints = var3
	}
	var0.DscpTos = var2
	var0.Name = plan.Name.ValueString()
	var0.Schedule = plan.Schedule.ValueString()
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var12 qosPolicyRulesRsModelActionObject
	var12.Class = types.StringValue(ans.Action.Class)
	var var13 *qosPolicyRulesRsModelDscpTosObject
	if ans.DscpTos != nil {
		var13 = &qosPolicyRulesRsModelDscpTosObject{}
		var var14 []qosPolicyRulesRsModelCodepointsObject
		if len(ans.DscpTos.Codepoints) != 0 {
			var14 = make([]qosPolicyRulesRsModelCodepointsObject, 0, len(ans.DscpTos.Codepoints))
			for var15Index := range ans.DscpTos.Codepoints {
				var15 := ans.DscpTos.Codepoints[var15Index]
				var var16 qosPolicyRulesRsModelCodepointsObject
				var var17 *qosPolicyRulesRsModelTypeObject
				if var15.Type != nil {
					var17 = &qosPolicyRulesRsModelTypeObject{}
					var var18 *qosPolicyRulesRsModelAfObject
					if var15.Type.Af != nil {
						var18 = &qosPolicyRulesRsModelAfObject{}
						var18.Codepoint = types.StringValue(var15.Type.Af.Codepoint)
					}
					var var19 *qosPolicyRulesRsModelCsObject
					if var15.Type.Cs != nil {
						var19 = &qosPolicyRulesRsModelCsObject{}
						var19.Codepoint = types.StringValue(var15.Type.Cs.Codepoint)
					}
					var var20 *qosPolicyRulesRsModelCustomObject
					if var15.Type.Custom != nil {
						var20 = &qosPolicyRulesRsModelCustomObject{}
						var var21 *qosPolicyRulesRsModelCodepointObject
						if var15.Type.Custom.Codepoint != nil {
							var21 = &qosPolicyRulesRsModelCodepointObject{}
							var21.BinaryValue = types.StringValue(var15.Type.Custom.Codepoint.BinaryValue)
							var21.CodepointName = types.StringValue(var15.Type.Custom.Codepoint.CodepointName)
						}
						var20.Codepoint = var21
					}
					var var22 *qosPolicyRulesRsModelTosObject
					if var15.Type.Tos != nil {
						var22 = &qosPolicyRulesRsModelTosObject{}
						var22.Codepoint = types.StringValue(var15.Type.Tos.Codepoint)
					}
					var17.Af = var18
					var17.Cs = var19
					var17.Custom = var20
					if var15.Type.Ef != nil {
						var17.Ef = types.BoolValue(true)
					}
					var17.Tos = var22
				}
				var16.Name = types.StringValue(var15.Name)
				var16.Type = var17
				var14 = append(var14, var16)
			}
		}
		var13.Codepoints = var14
	}
	state.Action = var12
	state.Description = types.StringValue(ans.Description)
	state.DscpTos = var13
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.Schedule = types.StringValue(ans.Schedule)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *qosPolicyRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var idType types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &idType)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := idType.ValueString()
	tokens := strings.Split(id, IdSeparator)
	if len(tokens) != 3 {
		resp.Diagnostics.AddError("Error in resource ID format", "Expected 3 tokens")
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource delete", map[string]any{
		"terraform_provider_function": "Delete",
		"resource_name":               "sase_qos_policy_rules",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 2, "Position": 1},
		"tokens":                      tokens,
	})

	svc := tzldypq.NewClient(r.client)
	input := tzldypq.DeleteInput{
		ObjectId: tokens[2],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *qosPolicyRulesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
