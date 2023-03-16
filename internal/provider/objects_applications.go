package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	hIzciTY "github.com/paloaltonetworks/sase-go/netsec/schema/objects/applications"
	rrePbcM "github.com/paloaltonetworks/sase-go/netsec/service/v1/applications"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	_ datasource.DataSource              = &objectsApplicationsListDataSource{}
	_ datasource.DataSourceWithConfigure = &objectsApplicationsListDataSource{}
)

func NewObjectsApplicationsListDataSource() datasource.DataSource {
	return &objectsApplicationsListDataSource{}
}

type objectsApplicationsListDataSource struct {
	client *sase.Client
}

type objectsApplicationsListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []objectsApplicationsListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type objectsApplicationsListDsModelConfig struct {
	AbleToTransferFile     types.Bool                                      `tfsdk:"able_to_transfer_file"`
	AlgDisableCapability   types.String                                    `tfsdk:"alg_disable_capability"`
	Category               types.String                                    `tfsdk:"category"`
	ConsumeBigBandwidth    types.Bool                                      `tfsdk:"consume_big_bandwidth"`
	DataIdent              types.Bool                                      `tfsdk:"data_ident"`
	Default                *objectsApplicationsListDsModelDefaultObject    `tfsdk:"default"`
	Description            types.String                                    `tfsdk:"description"`
	EvasiveBehavior        types.Bool                                      `tfsdk:"evasive_behavior"`
	FileTypeIdent          types.Bool                                      `tfsdk:"file_type_ident"`
	HasKnownVulnerability  types.Bool                                      `tfsdk:"has_known_vulnerability"`
	ObjectId               types.String                                    `tfsdk:"object_id"`
	Name                   types.String                                    `tfsdk:"name"`
	NoAppidCaching         types.Bool                                      `tfsdk:"no_appid_caching"`
	ParentApp              types.String                                    `tfsdk:"parent_app"`
	PervasiveUse           types.Bool                                      `tfsdk:"pervasive_use"`
	ProneToMisuse          types.Bool                                      `tfsdk:"prone_to_misuse"`
	Risk                   types.Int64                                     `tfsdk:"risk"`
	Signature              []objectsApplicationsListDsModelSignatureObject `tfsdk:"signature"`
	Subcategory            types.String                                    `tfsdk:"subcategory"`
	TcpHalfClosedTimeout   types.Int64                                     `tfsdk:"tcp_half_closed_timeout"`
	TcpTimeWaitTimeout     types.Int64                                     `tfsdk:"tcp_time_wait_timeout"`
	TcpTimeout             types.Int64                                     `tfsdk:"tcp_timeout"`
	Technology             types.String                                    `tfsdk:"technology"`
	Timeout                types.Int64                                     `tfsdk:"timeout"`
	TunnelApplications     types.Bool                                      `tfsdk:"tunnel_applications"`
	TunnelOtherApplication types.Bool                                      `tfsdk:"tunnel_other_application"`
	UdpTimeout             types.Int64                                     `tfsdk:"udp_timeout"`
	UsedByMalware          types.Bool                                      `tfsdk:"used_by_malware"`
	VirusIdent             types.Bool                                      `tfsdk:"virus_ident"`
}

type objectsApplicationsListDsModelDefaultObject struct {
	IdentByIcmp6Type  *objectsApplicationsListDsModelIdentByIcmp6TypeObject `tfsdk:"ident_by_icmp6_type"`
	IdentByIcmpType   *objectsApplicationsListDsModelIdentByIcmpTypeObject  `tfsdk:"ident_by_icmp_type"`
	IdentByIpProtocol types.String                                          `tfsdk:"ident_by_ip_protocol"`
	Port              []types.String                                        `tfsdk:"port"`
}

type objectsApplicationsListDsModelIdentByIcmp6TypeObject struct {
	Code types.String `tfsdk:"code"`
	Type types.String `tfsdk:"type"`
}

type objectsApplicationsListDsModelIdentByIcmpTypeObject struct {
	Code types.String `tfsdk:"code"`
	Type types.String `tfsdk:"type"`
}

type objectsApplicationsListDsModelSignatureObject struct {
	AndCondition []objectsApplicationsListDsModelAndConditionObject `tfsdk:"and_condition"`
	Comment      types.String                                       `tfsdk:"comment"`
	Name         types.String                                       `tfsdk:"name"`
	OrderFree    types.Bool                                         `tfsdk:"order_free"`
	Scope        types.String                                       `tfsdk:"scope"`
}

type objectsApplicationsListDsModelAndConditionObject struct {
	Name        types.String                                      `tfsdk:"name"`
	OrCondition []objectsApplicationsListDsModelOrConditionObject `tfsdk:"or_condition"`
}

type objectsApplicationsListDsModelOrConditionObject struct {
	Name     types.String                                 `tfsdk:"name"`
	Operator objectsApplicationsListDsModelOperatorObject `tfsdk:"operator"`
}

type objectsApplicationsListDsModelOperatorObject struct {
	EqualTo      *objectsApplicationsListDsModelEqualToObject      `tfsdk:"equal_to"`
	GreaterThan  *objectsApplicationsListDsModelGreaterThanObject  `tfsdk:"greater_than"`
	LessThan     *objectsApplicationsListDsModelLessThanObject     `tfsdk:"less_than"`
	PatternMatch *objectsApplicationsListDsModelPatternMatchObject `tfsdk:"pattern_match"`
}

type objectsApplicationsListDsModelEqualToObject struct {
	Context  types.String `tfsdk:"context"`
	Mask     types.String `tfsdk:"mask"`
	Position types.String `tfsdk:"position"`
	Value    types.String `tfsdk:"value"`
}

type objectsApplicationsListDsModelGreaterThanObject struct {
	Context   types.String                                    `tfsdk:"context"`
	Qualifier []objectsApplicationsListDsModelQualifierObject `tfsdk:"qualifier"`
	Value     types.Int64                                     `tfsdk:"value"`
}

type objectsApplicationsListDsModelQualifierObject struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type objectsApplicationsListDsModelLessThanObject struct {
	Context   types.String                                    `tfsdk:"context"`
	Qualifier []objectsApplicationsListDsModelQualifierObject `tfsdk:"qualifier"`
	Value     types.Int64                                     `tfsdk:"value"`
}

type objectsApplicationsListDsModelPatternMatchObject struct {
	Context   types.String                                    `tfsdk:"context"`
	Pattern   types.String                                    `tfsdk:"pattern"`
	Qualifier []objectsApplicationsListDsModelQualifierObject `tfsdk:"qualifier"`
}

// Metadata returns the data source type name.
func (d *objectsApplicationsListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_applications_list"
}

// Schema defines the schema for this listing data source.
func (d *objectsApplicationsListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Retrieves a listing of config items.",

		Attributes: map[string]dsschema.Attribute{
			"id": dsschema.StringAttribute{
				Description: "The object ID.",
				Computed:    true,
			},

			// Input.
			"limit": dsschema.Int64Attribute{
				Description: "The max count in result entry (count per page)",
				Optional:    true,
				Computed:    true,
			},
			"offset": dsschema.Int64Attribute{
				Description: "The offset of the result entry",
				Optional:    true,
				Computed:    true,
			},
			"name": dsschema.StringAttribute{
				Description: "The name of the entry",
				Optional:    true,
			},
			"folder": dsschema.StringAttribute{
				Description: "The folder of the entry",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Service Connections", "Mobile Users Container", "Mobile Users Explicit Proxy"),
				},
			},

			// Output.
			"data": dsschema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"able_to_transfer_file": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"alg_disable_capability": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"category": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"consume_big_bandwidth": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"data_ident": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"default": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"ident_by_icmp6_type": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"code": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"type": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
								"ident_by_icmp_type": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"code": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"type": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
								"ident_by_ip_protocol": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
								"port": dsschema.ListAttribute{
									Description: "",
									Computed:    true,
									ElementType: types.StringType,
								},
							},
						},
						"description": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"evasive_behavior": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"file_type_ident": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"has_known_vulnerability": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"object_id": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"no_appid_caching": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"parent_app": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"pervasive_use": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"prone_to_misuse": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"risk": dsschema.Int64Attribute{
							Description: "",
							Computed:    true,
						},
						"signature": dsschema.ListNestedAttribute{
							Description: "",
							Computed:    true,
							NestedObject: dsschema.NestedAttributeObject{
								Attributes: map[string]dsschema.Attribute{
									"and_condition": dsschema.ListNestedAttribute{
										Description: "",
										Computed:    true,
										NestedObject: dsschema.NestedAttributeObject{
											Attributes: map[string]dsschema.Attribute{
												"name": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"or_condition": dsschema.ListNestedAttribute{
													Description: "",
													Computed:    true,
													NestedObject: dsschema.NestedAttributeObject{
														Attributes: map[string]dsschema.Attribute{
															"name": dsschema.StringAttribute{
																Description: "",
																Computed:    true,
															},
															"operator": dsschema.SingleNestedAttribute{
																Description: "",
																Computed:    true,
																Attributes: map[string]dsschema.Attribute{
																	"equal_to": dsschema.SingleNestedAttribute{
																		Description: "",
																		Computed:    true,
																		Attributes: map[string]dsschema.Attribute{
																			"context": dsschema.StringAttribute{
																				Description: "",
																				Computed:    true,
																			},
																			"mask": dsschema.StringAttribute{
																				Description: "",
																				Computed:    true,
																			},
																			"position": dsschema.StringAttribute{
																				Description: "",
																				Computed:    true,
																			},
																			"value": dsschema.StringAttribute{
																				Description: "",
																				Computed:    true,
																			},
																		},
																	},
																	"greater_than": dsschema.SingleNestedAttribute{
																		Description: "",
																		Computed:    true,
																		Attributes: map[string]dsschema.Attribute{
																			"context": dsschema.StringAttribute{
																				Description: "",
																				Computed:    true,
																			},
																			"qualifier": dsschema.ListNestedAttribute{
																				Description: "",
																				Computed:    true,
																				NestedObject: dsschema.NestedAttributeObject{
																					Attributes: map[string]dsschema.Attribute{
																						"name": dsschema.StringAttribute{
																							Description: "",
																							Computed:    true,
																						},
																						"value": dsschema.StringAttribute{
																							Description: "",
																							Computed:    true,
																						},
																					},
																				},
																			},
																			"value": dsschema.Int64Attribute{
																				Description: "",
																				Computed:    true,
																			},
																		},
																	},
																	"less_than": dsschema.SingleNestedAttribute{
																		Description: "",
																		Computed:    true,
																		Attributes: map[string]dsschema.Attribute{
																			"context": dsschema.StringAttribute{
																				Description: "",
																				Computed:    true,
																			},
																			"qualifier": dsschema.ListNestedAttribute{
																				Description: "",
																				Computed:    true,
																				NestedObject: dsschema.NestedAttributeObject{
																					Attributes: map[string]dsschema.Attribute{
																						"name": dsschema.StringAttribute{
																							Description: "",
																							Computed:    true,
																						},
																						"value": dsschema.StringAttribute{
																							Description: "",
																							Computed:    true,
																						},
																					},
																				},
																			},
																			"value": dsschema.Int64Attribute{
																				Description: "",
																				Computed:    true,
																			},
																		},
																	},
																	"pattern_match": dsschema.SingleNestedAttribute{
																		Description: "",
																		Computed:    true,
																		Attributes: map[string]dsschema.Attribute{
																			"context": dsschema.StringAttribute{
																				Description: "",
																				Computed:    true,
																			},
																			"pattern": dsschema.StringAttribute{
																				Description: "",
																				Computed:    true,
																			},
																			"qualifier": dsschema.ListNestedAttribute{
																				Description: "",
																				Computed:    true,
																				NestedObject: dsschema.NestedAttributeObject{
																					Attributes: map[string]dsschema.Attribute{
																						"name": dsschema.StringAttribute{
																							Description: "",
																							Computed:    true,
																						},
																						"value": dsschema.StringAttribute{
																							Description: "",
																							Computed:    true,
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
											},
										},
									},
									"comment": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"name": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"order_free": dsschema.BoolAttribute{
										Description: "",
										Computed:    true,
									},
									"scope": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
						},
						"subcategory": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"tcp_half_closed_timeout": dsschema.Int64Attribute{
							Description: "",
							Computed:    true,
						},
						"tcp_time_wait_timeout": dsschema.Int64Attribute{
							Description: "",
							Computed:    true,
						},
						"tcp_timeout": dsschema.Int64Attribute{
							Description: "",
							Computed:    true,
						},
						"technology": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"timeout": dsschema.Int64Attribute{
							Description: "",
							Computed:    true,
						},
						"tunnel_applications": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"tunnel_other_application": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"udp_timeout": dsschema.Int64Attribute{
							Description: "",
							Computed:    true,
						},
						"used_by_malware": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"virus_ident": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
					},
				},
			},
			"total": dsschema.Int64Attribute{
				Description: "",
				Computed:    true,
			},
		},
	}
}

// Configure prepares the struct.
func (d *objectsApplicationsListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *objectsApplicationsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state objectsApplicationsListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_objects_applications_list",
		"limit":            state.Limit.ValueInt64(),
		"has_limit":        !state.Limit.IsNull(),
		"offset":           state.Offset.ValueInt64(),
		"has_offset":       !state.Offset.IsNull(),
		"name":             state.Name.ValueString(),
		"has_name":         !state.Name.IsNull(),
		"folder":           state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := rrePbcM.NewClient(d.client)
	input := rrePbcM.ListInput{
		Folder: state.Folder.ValueString(),
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
	state.Id = types.StringValue(strings.Join([]string{strconv.FormatInt(*input.Limit, 10), strconv.FormatInt(*input.Offset, 10), *input.Name, input.Folder}, IdSeparator))
	var var0 []objectsApplicationsListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]objectsApplicationsListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 objectsApplicationsListDsModelConfig
			var var3 *objectsApplicationsListDsModelDefaultObject
			if var1.Default != nil {
				var3 = &objectsApplicationsListDsModelDefaultObject{}
				var var4 *objectsApplicationsListDsModelIdentByIcmp6TypeObject
				if var1.Default.IdentByIcmp6Type != nil {
					var4 = &objectsApplicationsListDsModelIdentByIcmp6TypeObject{}
					var4.Code = types.StringValue(var1.Default.IdentByIcmp6Type.Code)
					var4.Type = types.StringValue(var1.Default.IdentByIcmp6Type.Type)
				}
				var var5 *objectsApplicationsListDsModelIdentByIcmpTypeObject
				if var1.Default.IdentByIcmpType != nil {
					var5 = &objectsApplicationsListDsModelIdentByIcmpTypeObject{}
					var5.Code = types.StringValue(var1.Default.IdentByIcmpType.Code)
					var5.Type = types.StringValue(var1.Default.IdentByIcmpType.Type)
				}
				var3.IdentByIcmp6Type = var4
				var3.IdentByIcmpType = var5
				var3.IdentByIpProtocol = types.StringValue(var1.Default.IdentByIpProtocol)
				var3.Port = EncodeStringSlice(var1.Default.Port)
			}
			var var6 []objectsApplicationsListDsModelSignatureObject
			if len(var1.Signature) != 0 {
				var6 = make([]objectsApplicationsListDsModelSignatureObject, 0, len(var1.Signature))
				for var7Index := range var1.Signature {
					var7 := var1.Signature[var7Index]
					var var8 objectsApplicationsListDsModelSignatureObject
					var var9 []objectsApplicationsListDsModelAndConditionObject
					if len(var7.AndCondition) != 0 {
						var9 = make([]objectsApplicationsListDsModelAndConditionObject, 0, len(var7.AndCondition))
						for var10Index := range var7.AndCondition {
							var10 := var7.AndCondition[var10Index]
							var var11 objectsApplicationsListDsModelAndConditionObject
							var var12 []objectsApplicationsListDsModelOrConditionObject
							if len(var10.OrCondition) != 0 {
								var12 = make([]objectsApplicationsListDsModelOrConditionObject, 0, len(var10.OrCondition))
								for var13Index := range var10.OrCondition {
									var13 := var10.OrCondition[var13Index]
									var var14 objectsApplicationsListDsModelOrConditionObject
									var var15 objectsApplicationsListDsModelOperatorObject
									var var16 *objectsApplicationsListDsModelEqualToObject
									if var13.Operator.EqualTo != nil {
										var16 = &objectsApplicationsListDsModelEqualToObject{}
										var16.Context = types.StringValue(var13.Operator.EqualTo.Context)
										var16.Mask = types.StringValue(var13.Operator.EqualTo.Mask)
										var16.Position = types.StringValue(var13.Operator.EqualTo.Position)
										var16.Value = types.StringValue(var13.Operator.EqualTo.Value)
									}
									var var17 *objectsApplicationsListDsModelGreaterThanObject
									if var13.Operator.GreaterThan != nil {
										var17 = &objectsApplicationsListDsModelGreaterThanObject{}
										var var18 []objectsApplicationsListDsModelQualifierObject
										if len(var13.Operator.GreaterThan.Qualifier) != 0 {
											var18 = make([]objectsApplicationsListDsModelQualifierObject, 0, len(var13.Operator.GreaterThan.Qualifier))
											for var19Index := range var13.Operator.GreaterThan.Qualifier {
												var19 := var13.Operator.GreaterThan.Qualifier[var19Index]
												var var20 objectsApplicationsListDsModelQualifierObject
												var20.Name = types.StringValue(var19.Name)
												var20.Value = types.StringValue(var19.Value)
												var18 = append(var18, var20)
											}
										}
										var17.Context = types.StringValue(var13.Operator.GreaterThan.Context)
										var17.Qualifier = var18
										var17.Value = types.Int64Value(var13.Operator.GreaterThan.Value)
									}
									var var21 *objectsApplicationsListDsModelLessThanObject
									if var13.Operator.LessThan != nil {
										var21 = &objectsApplicationsListDsModelLessThanObject{}
										var var22 []objectsApplicationsListDsModelQualifierObject
										if len(var13.Operator.LessThan.Qualifier) != 0 {
											var22 = make([]objectsApplicationsListDsModelQualifierObject, 0, len(var13.Operator.LessThan.Qualifier))
											for var23Index := range var13.Operator.LessThan.Qualifier {
												var23 := var13.Operator.LessThan.Qualifier[var23Index]
												var var24 objectsApplicationsListDsModelQualifierObject
												var24.Name = types.StringValue(var23.Name)
												var24.Value = types.StringValue(var23.Value)
												var22 = append(var22, var24)
											}
										}
										var21.Context = types.StringValue(var13.Operator.LessThan.Context)
										var21.Qualifier = var22
										var21.Value = types.Int64Value(var13.Operator.LessThan.Value)
									}
									var var25 *objectsApplicationsListDsModelPatternMatchObject
									if var13.Operator.PatternMatch != nil {
										var25 = &objectsApplicationsListDsModelPatternMatchObject{}
										var var26 []objectsApplicationsListDsModelQualifierObject
										if len(var13.Operator.PatternMatch.Qualifier) != 0 {
											var26 = make([]objectsApplicationsListDsModelQualifierObject, 0, len(var13.Operator.PatternMatch.Qualifier))
											for var27Index := range var13.Operator.PatternMatch.Qualifier {
												var27 := var13.Operator.PatternMatch.Qualifier[var27Index]
												var var28 objectsApplicationsListDsModelQualifierObject
												var28.Name = types.StringValue(var27.Name)
												var28.Value = types.StringValue(var27.Value)
												var26 = append(var26, var28)
											}
										}
										var25.Context = types.StringValue(var13.Operator.PatternMatch.Context)
										var25.Pattern = types.StringValue(var13.Operator.PatternMatch.Pattern)
										var25.Qualifier = var26
									}
									var15.EqualTo = var16
									var15.GreaterThan = var17
									var15.LessThan = var21
									var15.PatternMatch = var25
									var14.Name = types.StringValue(var13.Name)
									var14.Operator = var15
									var12 = append(var12, var14)
								}
							}
							var11.Name = types.StringValue(var10.Name)
							var11.OrCondition = var12
							var9 = append(var9, var11)
						}
					}
					var8.AndCondition = var9
					var8.Comment = types.StringValue(var7.Comment)
					var8.Name = types.StringValue(var7.Name)
					var8.OrderFree = types.BoolValue(var7.OrderFree)
					var8.Scope = types.StringValue(var7.Scope)
					var6 = append(var6, var8)
				}
			}
			var2.AbleToTransferFile = types.BoolValue(var1.AbleToTransferFile)
			var2.AlgDisableCapability = types.StringValue(var1.AlgDisableCapability)
			var2.Category = types.StringValue(var1.Category)
			var2.ConsumeBigBandwidth = types.BoolValue(var1.ConsumeBigBandwidth)
			var2.DataIdent = types.BoolValue(var1.DataIdent)
			var2.Default = var3
			var2.Description = types.StringValue(var1.Description)
			var2.EvasiveBehavior = types.BoolValue(var1.EvasiveBehavior)
			var2.FileTypeIdent = types.BoolValue(var1.FileTypeIdent)
			var2.HasKnownVulnerability = types.BoolValue(var1.HasKnownVulnerability)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Name = types.StringValue(var1.Name)
			var2.NoAppidCaching = types.BoolValue(var1.NoAppidCaching)
			var2.ParentApp = types.StringValue(var1.ParentApp)
			var2.PervasiveUse = types.BoolValue(var1.PervasiveUse)
			var2.ProneToMisuse = types.BoolValue(var1.ProneToMisuse)
			var2.Risk = types.Int64Value(var1.Risk)
			var2.Signature = var6
			var2.Subcategory = types.StringValue(var1.Subcategory)
			var2.TcpHalfClosedTimeout = types.Int64Value(var1.TcpHalfClosedTimeout)
			var2.TcpTimeWaitTimeout = types.Int64Value(var1.TcpTimeWaitTimeout)
			var2.TcpTimeout = types.Int64Value(var1.TcpTimeout)
			var2.Technology = types.StringValue(var1.Technology)
			var2.Timeout = types.Int64Value(var1.Timeout)
			var2.TunnelApplications = types.BoolValue(var1.TunnelApplications)
			var2.TunnelOtherApplication = types.BoolValue(var1.TunnelOtherApplication)
			var2.UdpTimeout = types.Int64Value(var1.UdpTimeout)
			var2.UsedByMalware = types.BoolValue(var1.UsedByMalware)
			var2.VirusIdent = types.BoolValue(var1.VirusIdent)
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
	_ datasource.DataSource              = &objectsApplicationsDataSource{}
	_ datasource.DataSourceWithConfigure = &objectsApplicationsDataSource{}
)

func NewObjectsApplicationsDataSource() datasource.DataSource {
	return &objectsApplicationsDataSource{}
}

type objectsApplicationsDataSource struct {
	client *sase.Client
}

type objectsApplicationsDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`

	// Output.
	// Ref: #/components/schemas/objects-applications
	AbleToTransferFile    types.Bool                               `tfsdk:"able_to_transfer_file"`
	AlgDisableCapability  types.String                             `tfsdk:"alg_disable_capability"`
	Category              types.String                             `tfsdk:"category"`
	ConsumeBigBandwidth   types.Bool                               `tfsdk:"consume_big_bandwidth"`
	DataIdent             types.Bool                               `tfsdk:"data_ident"`
	Default               *objectsApplicationsDsModelDefaultObject `tfsdk:"default"`
	Description           types.String                             `tfsdk:"description"`
	EvasiveBehavior       types.Bool                               `tfsdk:"evasive_behavior"`
	FileTypeIdent         types.Bool                               `tfsdk:"file_type_ident"`
	HasKnownVulnerability types.Bool                               `tfsdk:"has_known_vulnerability"`
	// input omit: ObjectId
	Name                   types.String                                `tfsdk:"name"`
	NoAppidCaching         types.Bool                                  `tfsdk:"no_appid_caching"`
	ParentApp              types.String                                `tfsdk:"parent_app"`
	PervasiveUse           types.Bool                                  `tfsdk:"pervasive_use"`
	ProneToMisuse          types.Bool                                  `tfsdk:"prone_to_misuse"`
	Risk                   types.Int64                                 `tfsdk:"risk"`
	Signature              []objectsApplicationsDsModelSignatureObject `tfsdk:"signature"`
	Subcategory            types.String                                `tfsdk:"subcategory"`
	TcpHalfClosedTimeout   types.Int64                                 `tfsdk:"tcp_half_closed_timeout"`
	TcpTimeWaitTimeout     types.Int64                                 `tfsdk:"tcp_time_wait_timeout"`
	TcpTimeout             types.Int64                                 `tfsdk:"tcp_timeout"`
	Technology             types.String                                `tfsdk:"technology"`
	Timeout                types.Int64                                 `tfsdk:"timeout"`
	TunnelApplications     types.Bool                                  `tfsdk:"tunnel_applications"`
	TunnelOtherApplication types.Bool                                  `tfsdk:"tunnel_other_application"`
	UdpTimeout             types.Int64                                 `tfsdk:"udp_timeout"`
	UsedByMalware          types.Bool                                  `tfsdk:"used_by_malware"`
	VirusIdent             types.Bool                                  `tfsdk:"virus_ident"`
}

type objectsApplicationsDsModelDefaultObject struct {
	IdentByIcmp6Type  *objectsApplicationsDsModelIdentByIcmp6TypeObject `tfsdk:"ident_by_icmp6_type"`
	IdentByIcmpType   *objectsApplicationsDsModelIdentByIcmpTypeObject  `tfsdk:"ident_by_icmp_type"`
	IdentByIpProtocol types.String                                      `tfsdk:"ident_by_ip_protocol"`
	Port              []types.String                                    `tfsdk:"port"`
}

type objectsApplicationsDsModelIdentByIcmp6TypeObject struct {
	Code types.String `tfsdk:"code"`
	Type types.String `tfsdk:"type"`
}

type objectsApplicationsDsModelIdentByIcmpTypeObject struct {
	Code types.String `tfsdk:"code"`
	Type types.String `tfsdk:"type"`
}

type objectsApplicationsDsModelSignatureObject struct {
	AndCondition []objectsApplicationsDsModelAndConditionObject `tfsdk:"and_condition"`
	Comment      types.String                                   `tfsdk:"comment"`
	Name         types.String                                   `tfsdk:"name"`
	OrderFree    types.Bool                                     `tfsdk:"order_free"`
	Scope        types.String                                   `tfsdk:"scope"`
}

type objectsApplicationsDsModelAndConditionObject struct {
	Name        types.String                                  `tfsdk:"name"`
	OrCondition []objectsApplicationsDsModelOrConditionObject `tfsdk:"or_condition"`
}

type objectsApplicationsDsModelOrConditionObject struct {
	Name     types.String                             `tfsdk:"name"`
	Operator objectsApplicationsDsModelOperatorObject `tfsdk:"operator"`
}

type objectsApplicationsDsModelOperatorObject struct {
	EqualTo      *objectsApplicationsDsModelEqualToObject      `tfsdk:"equal_to"`
	GreaterThan  *objectsApplicationsDsModelGreaterThanObject  `tfsdk:"greater_than"`
	LessThan     *objectsApplicationsDsModelLessThanObject     `tfsdk:"less_than"`
	PatternMatch *objectsApplicationsDsModelPatternMatchObject `tfsdk:"pattern_match"`
}

type objectsApplicationsDsModelEqualToObject struct {
	Context  types.String `tfsdk:"context"`
	Mask     types.String `tfsdk:"mask"`
	Position types.String `tfsdk:"position"`
	Value    types.String `tfsdk:"value"`
}

type objectsApplicationsDsModelGreaterThanObject struct {
	Context   types.String                                `tfsdk:"context"`
	Qualifier []objectsApplicationsDsModelQualifierObject `tfsdk:"qualifier"`
	Value     types.Int64                                 `tfsdk:"value"`
}

type objectsApplicationsDsModelQualifierObject struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type objectsApplicationsDsModelLessThanObject struct {
	Context   types.String                                `tfsdk:"context"`
	Qualifier []objectsApplicationsDsModelQualifierObject `tfsdk:"qualifier"`
	Value     types.Int64                                 `tfsdk:"value"`
}

type objectsApplicationsDsModelPatternMatchObject struct {
	Context   types.String                                `tfsdk:"context"`
	Pattern   types.String                                `tfsdk:"pattern"`
	Qualifier []objectsApplicationsDsModelQualifierObject `tfsdk:"qualifier"`
}

// Metadata returns the data source type name.
func (d *objectsApplicationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_applications"
}

// Schema defines the schema for this listing data source.
func (d *objectsApplicationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Retrieves config for a specific item.",

		Attributes: map[string]dsschema.Attribute{
			"id": dsschema.StringAttribute{
				Description: "The object ID.",
				Computed:    true,
			},

			// Input.
			"object_id": dsschema.StringAttribute{
				Description: "The uuid of the resource",
				Required:    true,
			},

			// Output.
			"able_to_transfer_file": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"alg_disable_capability": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"category": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"consume_big_bandwidth": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"data_ident": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"default": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"ident_by_icmp6_type": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"code": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"type": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
						},
					},
					"ident_by_icmp_type": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"code": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"type": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
						},
					},
					"ident_by_ip_protocol": dsschema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"port": dsschema.ListAttribute{
						Description: "",
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
			"description": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"evasive_behavior": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"file_type_ident": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"has_known_vulnerability": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"no_appid_caching": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"parent_app": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"pervasive_use": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"prone_to_misuse": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"risk": dsschema.Int64Attribute{
				Description: "",
				Computed:    true,
			},
			"signature": dsschema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"and_condition": dsschema.ListNestedAttribute{
							Description: "",
							Computed:    true,
							NestedObject: dsschema.NestedAttributeObject{
								Attributes: map[string]dsschema.Attribute{
									"name": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"or_condition": dsschema.ListNestedAttribute{
										Description: "",
										Computed:    true,
										NestedObject: dsschema.NestedAttributeObject{
											Attributes: map[string]dsschema.Attribute{
												"name": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"operator": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"equal_to": dsschema.SingleNestedAttribute{
															Description: "",
															Computed:    true,
															Attributes: map[string]dsschema.Attribute{
																"context": dsschema.StringAttribute{
																	Description: "",
																	Computed:    true,
																},
																"mask": dsschema.StringAttribute{
																	Description: "",
																	Computed:    true,
																},
																"position": dsschema.StringAttribute{
																	Description: "",
																	Computed:    true,
																},
																"value": dsschema.StringAttribute{
																	Description: "",
																	Computed:    true,
																},
															},
														},
														"greater_than": dsschema.SingleNestedAttribute{
															Description: "",
															Computed:    true,
															Attributes: map[string]dsschema.Attribute{
																"context": dsschema.StringAttribute{
																	Description: "",
																	Computed:    true,
																},
																"qualifier": dsschema.ListNestedAttribute{
																	Description: "",
																	Computed:    true,
																	NestedObject: dsschema.NestedAttributeObject{
																		Attributes: map[string]dsschema.Attribute{
																			"name": dsschema.StringAttribute{
																				Description: "",
																				Computed:    true,
																			},
																			"value": dsschema.StringAttribute{
																				Description: "",
																				Computed:    true,
																			},
																		},
																	},
																},
																"value": dsschema.Int64Attribute{
																	Description: "",
																	Computed:    true,
																},
															},
														},
														"less_than": dsschema.SingleNestedAttribute{
															Description: "",
															Computed:    true,
															Attributes: map[string]dsschema.Attribute{
																"context": dsschema.StringAttribute{
																	Description: "",
																	Computed:    true,
																},
																"qualifier": dsschema.ListNestedAttribute{
																	Description: "",
																	Computed:    true,
																	NestedObject: dsschema.NestedAttributeObject{
																		Attributes: map[string]dsschema.Attribute{
																			"name": dsschema.StringAttribute{
																				Description: "",
																				Computed:    true,
																			},
																			"value": dsschema.StringAttribute{
																				Description: "",
																				Computed:    true,
																			},
																		},
																	},
																},
																"value": dsschema.Int64Attribute{
																	Description: "",
																	Computed:    true,
																},
															},
														},
														"pattern_match": dsschema.SingleNestedAttribute{
															Description: "",
															Computed:    true,
															Attributes: map[string]dsschema.Attribute{
																"context": dsschema.StringAttribute{
																	Description: "",
																	Computed:    true,
																},
																"pattern": dsschema.StringAttribute{
																	Description: "",
																	Computed:    true,
																},
																"qualifier": dsschema.ListNestedAttribute{
																	Description: "",
																	Computed:    true,
																	NestedObject: dsschema.NestedAttributeObject{
																		Attributes: map[string]dsschema.Attribute{
																			"name": dsschema.StringAttribute{
																				Description: "",
																				Computed:    true,
																			},
																			"value": dsschema.StringAttribute{
																				Description: "",
																				Computed:    true,
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
								},
							},
						},
						"comment": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"order_free": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"scope": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
					},
				},
			},
			"subcategory": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"tcp_half_closed_timeout": dsschema.Int64Attribute{
				Description: "",
				Computed:    true,
			},
			"tcp_time_wait_timeout": dsschema.Int64Attribute{
				Description: "",
				Computed:    true,
			},
			"tcp_timeout": dsschema.Int64Attribute{
				Description: "",
				Computed:    true,
			},
			"technology": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"timeout": dsschema.Int64Attribute{
				Description: "",
				Computed:    true,
			},
			"tunnel_applications": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"tunnel_other_application": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"udp_timeout": dsschema.Int64Attribute{
				Description: "",
				Computed:    true,
			},
			"used_by_malware": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"virus_ident": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
		},
	}
}

// Configure prepares the struct.
func (d *objectsApplicationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *objectsApplicationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state objectsApplicationsDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_objects_applications",
		"object_id":        state.ObjectId.ValueString(),
	})

	// Prepare to run the command.
	svc := rrePbcM.NewClient(d.client)
	input := rrePbcM.ReadInput{
		ObjectId: state.ObjectId.ValueString(),
	}

	// Perform the operation.
	ans, err := svc.Read(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error getting singleton", err.Error())
		return
	}

	// Store the answer to state.
	state.Id = types.StringValue(strings.Join([]string{input.ObjectId}, IdSeparator))
	var var0 *objectsApplicationsDsModelDefaultObject
	if ans.Default != nil {
		var0 = &objectsApplicationsDsModelDefaultObject{}
		var var1 *objectsApplicationsDsModelIdentByIcmp6TypeObject
		if ans.Default.IdentByIcmp6Type != nil {
			var1 = &objectsApplicationsDsModelIdentByIcmp6TypeObject{}
			var1.Code = types.StringValue(ans.Default.IdentByIcmp6Type.Code)
			var1.Type = types.StringValue(ans.Default.IdentByIcmp6Type.Type)
		}
		var var2 *objectsApplicationsDsModelIdentByIcmpTypeObject
		if ans.Default.IdentByIcmpType != nil {
			var2 = &objectsApplicationsDsModelIdentByIcmpTypeObject{}
			var2.Code = types.StringValue(ans.Default.IdentByIcmpType.Code)
			var2.Type = types.StringValue(ans.Default.IdentByIcmpType.Type)
		}
		var0.IdentByIcmp6Type = var1
		var0.IdentByIcmpType = var2
		var0.IdentByIpProtocol = types.StringValue(ans.Default.IdentByIpProtocol)
		var0.Port = EncodeStringSlice(ans.Default.Port)
	}
	var var3 []objectsApplicationsDsModelSignatureObject
	if len(ans.Signature) != 0 {
		var3 = make([]objectsApplicationsDsModelSignatureObject, 0, len(ans.Signature))
		for var4Index := range ans.Signature {
			var4 := ans.Signature[var4Index]
			var var5 objectsApplicationsDsModelSignatureObject
			var var6 []objectsApplicationsDsModelAndConditionObject
			if len(var4.AndCondition) != 0 {
				var6 = make([]objectsApplicationsDsModelAndConditionObject, 0, len(var4.AndCondition))
				for var7Index := range var4.AndCondition {
					var7 := var4.AndCondition[var7Index]
					var var8 objectsApplicationsDsModelAndConditionObject
					var var9 []objectsApplicationsDsModelOrConditionObject
					if len(var7.OrCondition) != 0 {
						var9 = make([]objectsApplicationsDsModelOrConditionObject, 0, len(var7.OrCondition))
						for var10Index := range var7.OrCondition {
							var10 := var7.OrCondition[var10Index]
							var var11 objectsApplicationsDsModelOrConditionObject
							var var12 objectsApplicationsDsModelOperatorObject
							var var13 *objectsApplicationsDsModelEqualToObject
							if var10.Operator.EqualTo != nil {
								var13 = &objectsApplicationsDsModelEqualToObject{}
								var13.Context = types.StringValue(var10.Operator.EqualTo.Context)
								var13.Mask = types.StringValue(var10.Operator.EqualTo.Mask)
								var13.Position = types.StringValue(var10.Operator.EqualTo.Position)
								var13.Value = types.StringValue(var10.Operator.EqualTo.Value)
							}
							var var14 *objectsApplicationsDsModelGreaterThanObject
							if var10.Operator.GreaterThan != nil {
								var14 = &objectsApplicationsDsModelGreaterThanObject{}
								var var15 []objectsApplicationsDsModelQualifierObject
								if len(var10.Operator.GreaterThan.Qualifier) != 0 {
									var15 = make([]objectsApplicationsDsModelQualifierObject, 0, len(var10.Operator.GreaterThan.Qualifier))
									for var16Index := range var10.Operator.GreaterThan.Qualifier {
										var16 := var10.Operator.GreaterThan.Qualifier[var16Index]
										var var17 objectsApplicationsDsModelQualifierObject
										var17.Name = types.StringValue(var16.Name)
										var17.Value = types.StringValue(var16.Value)
										var15 = append(var15, var17)
									}
								}
								var14.Context = types.StringValue(var10.Operator.GreaterThan.Context)
								var14.Qualifier = var15
								var14.Value = types.Int64Value(var10.Operator.GreaterThan.Value)
							}
							var var18 *objectsApplicationsDsModelLessThanObject
							if var10.Operator.LessThan != nil {
								var18 = &objectsApplicationsDsModelLessThanObject{}
								var var19 []objectsApplicationsDsModelQualifierObject
								if len(var10.Operator.LessThan.Qualifier) != 0 {
									var19 = make([]objectsApplicationsDsModelQualifierObject, 0, len(var10.Operator.LessThan.Qualifier))
									for var20Index := range var10.Operator.LessThan.Qualifier {
										var20 := var10.Operator.LessThan.Qualifier[var20Index]
										var var21 objectsApplicationsDsModelQualifierObject
										var21.Name = types.StringValue(var20.Name)
										var21.Value = types.StringValue(var20.Value)
										var19 = append(var19, var21)
									}
								}
								var18.Context = types.StringValue(var10.Operator.LessThan.Context)
								var18.Qualifier = var19
								var18.Value = types.Int64Value(var10.Operator.LessThan.Value)
							}
							var var22 *objectsApplicationsDsModelPatternMatchObject
							if var10.Operator.PatternMatch != nil {
								var22 = &objectsApplicationsDsModelPatternMatchObject{}
								var var23 []objectsApplicationsDsModelQualifierObject
								if len(var10.Operator.PatternMatch.Qualifier) != 0 {
									var23 = make([]objectsApplicationsDsModelQualifierObject, 0, len(var10.Operator.PatternMatch.Qualifier))
									for var24Index := range var10.Operator.PatternMatch.Qualifier {
										var24 := var10.Operator.PatternMatch.Qualifier[var24Index]
										var var25 objectsApplicationsDsModelQualifierObject
										var25.Name = types.StringValue(var24.Name)
										var25.Value = types.StringValue(var24.Value)
										var23 = append(var23, var25)
									}
								}
								var22.Context = types.StringValue(var10.Operator.PatternMatch.Context)
								var22.Pattern = types.StringValue(var10.Operator.PatternMatch.Pattern)
								var22.Qualifier = var23
							}
							var12.EqualTo = var13
							var12.GreaterThan = var14
							var12.LessThan = var18
							var12.PatternMatch = var22
							var11.Name = types.StringValue(var10.Name)
							var11.Operator = var12
							var9 = append(var9, var11)
						}
					}
					var8.Name = types.StringValue(var7.Name)
					var8.OrCondition = var9
					var6 = append(var6, var8)
				}
			}
			var5.AndCondition = var6
			var5.Comment = types.StringValue(var4.Comment)
			var5.Name = types.StringValue(var4.Name)
			var5.OrderFree = types.BoolValue(var4.OrderFree)
			var5.Scope = types.StringValue(var4.Scope)
			var3 = append(var3, var5)
		}
	}
	state.AbleToTransferFile = types.BoolValue(ans.AbleToTransferFile)
	state.AlgDisableCapability = types.StringValue(ans.AlgDisableCapability)
	state.Category = types.StringValue(ans.Category)
	state.ConsumeBigBandwidth = types.BoolValue(ans.ConsumeBigBandwidth)
	state.DataIdent = types.BoolValue(ans.DataIdent)
	state.Default = var0
	state.Description = types.StringValue(ans.Description)
	state.EvasiveBehavior = types.BoolValue(ans.EvasiveBehavior)
	state.FileTypeIdent = types.BoolValue(ans.FileTypeIdent)
	state.HasKnownVulnerability = types.BoolValue(ans.HasKnownVulnerability)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.NoAppidCaching = types.BoolValue(ans.NoAppidCaching)
	state.ParentApp = types.StringValue(ans.ParentApp)
	state.PervasiveUse = types.BoolValue(ans.PervasiveUse)
	state.ProneToMisuse = types.BoolValue(ans.ProneToMisuse)
	state.Risk = types.Int64Value(ans.Risk)
	state.Signature = var3
	state.Subcategory = types.StringValue(ans.Subcategory)
	state.TcpHalfClosedTimeout = types.Int64Value(ans.TcpHalfClosedTimeout)
	state.TcpTimeWaitTimeout = types.Int64Value(ans.TcpTimeWaitTimeout)
	state.TcpTimeout = types.Int64Value(ans.TcpTimeout)
	state.Technology = types.StringValue(ans.Technology)
	state.Timeout = types.Int64Value(ans.Timeout)
	state.TunnelApplications = types.BoolValue(ans.TunnelApplications)
	state.TunnelOtherApplication = types.BoolValue(ans.TunnelOtherApplication)
	state.UdpTimeout = types.Int64Value(ans.UdpTimeout)
	state.UsedByMalware = types.BoolValue(ans.UsedByMalware)
	state.VirusIdent = types.BoolValue(ans.VirusIdent)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &objectsApplicationsResource{}
	_ resource.ResourceWithConfigure   = &objectsApplicationsResource{}
	_ resource.ResourceWithImportState = &objectsApplicationsResource{}
)

func NewObjectsApplicationsResource() resource.Resource {
	return &objectsApplicationsResource{}
}

type objectsApplicationsResource struct {
	client *sase.Client
}

type objectsApplicationsRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/objects-applications
	AbleToTransferFile     types.Bool                                  `tfsdk:"able_to_transfer_file"`
	AlgDisableCapability   types.String                                `tfsdk:"alg_disable_capability"`
	Category               types.String                                `tfsdk:"category"`
	ConsumeBigBandwidth    types.Bool                                  `tfsdk:"consume_big_bandwidth"`
	DataIdent              types.Bool                                  `tfsdk:"data_ident"`
	Default                *objectsApplicationsRsModelDefaultObject    `tfsdk:"default"`
	Description            types.String                                `tfsdk:"description"`
	EvasiveBehavior        types.Bool                                  `tfsdk:"evasive_behavior"`
	FileTypeIdent          types.Bool                                  `tfsdk:"file_type_ident"`
	HasKnownVulnerability  types.Bool                                  `tfsdk:"has_known_vulnerability"`
	ObjectId               types.String                                `tfsdk:"object_id"`
	Name                   types.String                                `tfsdk:"name"`
	NoAppidCaching         types.Bool                                  `tfsdk:"no_appid_caching"`
	ParentApp              types.String                                `tfsdk:"parent_app"`
	PervasiveUse           types.Bool                                  `tfsdk:"pervasive_use"`
	ProneToMisuse          types.Bool                                  `tfsdk:"prone_to_misuse"`
	Risk                   types.Int64                                 `tfsdk:"risk"`
	Signature              []objectsApplicationsRsModelSignatureObject `tfsdk:"signature"`
	Subcategory            types.String                                `tfsdk:"subcategory"`
	TcpHalfClosedTimeout   types.Int64                                 `tfsdk:"tcp_half_closed_timeout"`
	TcpTimeWaitTimeout     types.Int64                                 `tfsdk:"tcp_time_wait_timeout"`
	TcpTimeout             types.Int64                                 `tfsdk:"tcp_timeout"`
	Technology             types.String                                `tfsdk:"technology"`
	Timeout                types.Int64                                 `tfsdk:"timeout"`
	TunnelApplications     types.Bool                                  `tfsdk:"tunnel_applications"`
	TunnelOtherApplication types.Bool                                  `tfsdk:"tunnel_other_application"`
	UdpTimeout             types.Int64                                 `tfsdk:"udp_timeout"`
	UsedByMalware          types.Bool                                  `tfsdk:"used_by_malware"`
	VirusIdent             types.Bool                                  `tfsdk:"virus_ident"`
}

type objectsApplicationsRsModelDefaultObject struct {
	IdentByIcmp6Type  *objectsApplicationsRsModelIdentByIcmp6TypeObject `tfsdk:"ident_by_icmp6_type"`
	IdentByIcmpType   *objectsApplicationsRsModelIdentByIcmpTypeObject  `tfsdk:"ident_by_icmp_type"`
	IdentByIpProtocol types.String                                      `tfsdk:"ident_by_ip_protocol"`
	Port              []types.String                                    `tfsdk:"port"`
}

type objectsApplicationsRsModelIdentByIcmp6TypeObject struct {
	Code types.String `tfsdk:"code"`
	Type types.String `tfsdk:"type"`
}

type objectsApplicationsRsModelIdentByIcmpTypeObject struct {
	Code types.String `tfsdk:"code"`
	Type types.String `tfsdk:"type"`
}

type objectsApplicationsRsModelSignatureObject struct {
	AndCondition []objectsApplicationsRsModelAndConditionObject `tfsdk:"and_condition"`
	Comment      types.String                                   `tfsdk:"comment"`
	Name         types.String                                   `tfsdk:"name"`
	OrderFree    types.Bool                                     `tfsdk:"order_free"`
	Scope        types.String                                   `tfsdk:"scope"`
}

type objectsApplicationsRsModelAndConditionObject struct {
	Name        types.String                                  `tfsdk:"name"`
	OrCondition []objectsApplicationsRsModelOrConditionObject `tfsdk:"or_condition"`
}

type objectsApplicationsRsModelOrConditionObject struct {
	Name     types.String                             `tfsdk:"name"`
	Operator objectsApplicationsRsModelOperatorObject `tfsdk:"operator"`
}

type objectsApplicationsRsModelOperatorObject struct {
	EqualTo      *objectsApplicationsRsModelEqualToObject      `tfsdk:"equal_to"`
	GreaterThan  *objectsApplicationsRsModelGreaterThanObject  `tfsdk:"greater_than"`
	LessThan     *objectsApplicationsRsModelLessThanObject     `tfsdk:"less_than"`
	PatternMatch *objectsApplicationsRsModelPatternMatchObject `tfsdk:"pattern_match"`
}

type objectsApplicationsRsModelEqualToObject struct {
	Context  types.String `tfsdk:"context"`
	Mask     types.String `tfsdk:"mask"`
	Position types.String `tfsdk:"position"`
	Value    types.String `tfsdk:"value"`
}

type objectsApplicationsRsModelGreaterThanObject struct {
	Context   types.String                                `tfsdk:"context"`
	Qualifier []objectsApplicationsRsModelQualifierObject `tfsdk:"qualifier"`
	Value     types.Int64                                 `tfsdk:"value"`
}

type objectsApplicationsRsModelQualifierObject struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type objectsApplicationsRsModelLessThanObject struct {
	Context   types.String                                `tfsdk:"context"`
	Qualifier []objectsApplicationsRsModelQualifierObject `tfsdk:"qualifier"`
	Value     types.Int64                                 `tfsdk:"value"`
}

type objectsApplicationsRsModelPatternMatchObject struct {
	Context   types.String                                `tfsdk:"context"`
	Pattern   types.String                                `tfsdk:"pattern"`
	Qualifier []objectsApplicationsRsModelQualifierObject `tfsdk:"qualifier"`
}

// Metadata returns the data source type name.
func (r *objectsApplicationsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_applications"
}

// Schema defines the schema for this listing data source.
func (r *objectsApplicationsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = rsschema.Schema{
		Description: "Retrieves config for a specific item.",

		Attributes: map[string]rsschema.Attribute{
			"id": rsschema.StringAttribute{
				Description: "The object ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			// Input.
			"folder": rsschema.StringAttribute{
				Description: "The folder of the entry",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Service Connections", "Mobile Users Container", "Mobile Users Explicit Proxy"),
				},
			},

			"able_to_transfer_file": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"alg_disable_capability": rsschema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(127),
				},
			},
			"category": rsschema.StringAttribute{
				Description: "",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
			},
			"consume_big_bandwidth": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"data_ident": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"default": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"ident_by_icmp6_type": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"code": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"type": rsschema.StringAttribute{
								Description: "",
								Required:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
						},
					},
					"ident_by_icmp_type": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"code": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"type": rsschema.StringAttribute{
								Description: "",
								Required:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
						},
					},
					"ident_by_ip_protocol": rsschema.StringAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							DefaultString(""),
						},
					},
					"port": rsschema.ListAttribute{
						Description: "",
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
			"description": rsschema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 1023),
				},
			},
			"evasive_behavior": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"file_type_ident": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"has_known_vulnerability": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"object_id": rsschema.StringAttribute{
				Description: "",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": rsschema.StringAttribute{
				Description: "",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(31),
				},
			},
			"no_appid_caching": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"parent_app": rsschema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(127),
				},
			},
			"pervasive_use": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"prone_to_misuse": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"risk": rsschema.Int64Attribute{
				Description: "",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					DefaultInt64(0),
				},
				Validators: []validator.Int64{
					int64validator.Between(1, 5),
				},
			},
			"signature": rsschema.ListNestedAttribute{
				Description: "",
				Optional:    true,
				NestedObject: rsschema.NestedAttributeObject{
					Attributes: map[string]rsschema.Attribute{
						"and_condition": rsschema.ListNestedAttribute{
							Description: "",
							Optional:    true,
							NestedObject: rsschema.NestedAttributeObject{
								Attributes: map[string]rsschema.Attribute{
									"name": rsschema.StringAttribute{
										Description: "",
										Required:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthAtMost(31),
										},
									},
									"or_condition": rsschema.ListNestedAttribute{
										Description: "",
										Optional:    true,
										NestedObject: rsschema.NestedAttributeObject{
											Attributes: map[string]rsschema.Attribute{
												"name": rsschema.StringAttribute{
													Description: "",
													Required:    true,
													PlanModifiers: []planmodifier.String{
														DefaultString(""),
													},
													Validators: []validator.String{
														stringvalidator.LengthAtMost(31),
													},
												},
												"operator": rsschema.SingleNestedAttribute{
													Description: "",
													Required:    true,
													Attributes: map[string]rsschema.Attribute{
														"equal_to": rsschema.SingleNestedAttribute{
															Description: "",
															Optional:    true,
															Attributes: map[string]rsschema.Attribute{
																"context": rsschema.StringAttribute{
																	Description: "",
																	Required:    true,
																	PlanModifiers: []planmodifier.String{
																		DefaultString(""),
																	},
																},
																"mask": rsschema.StringAttribute{
																	Description: "",
																	Optional:    true,
																	Computed:    true,
																	PlanModifiers: []planmodifier.String{
																		DefaultString(""),
																	},
																	Validators: []validator.String{
																		stringvalidator.LengthAtMost(10),
																	},
																},
																"position": rsschema.StringAttribute{
																	Description: "",
																	Optional:    true,
																	Computed:    true,
																	PlanModifiers: []planmodifier.String{
																		DefaultString(""),
																	},
																	Validators: []validator.String{
																		stringvalidator.LengthAtMost(127),
																	},
																},
																"value": rsschema.StringAttribute{
																	Description: "",
																	Required:    true,
																	PlanModifiers: []planmodifier.String{
																		DefaultString(""),
																	},
																	Validators: []validator.String{
																		stringvalidator.LengthAtMost(10),
																	},
																},
															},
														},
														"greater_than": rsschema.SingleNestedAttribute{
															Description: "",
															Optional:    true,
															Attributes: map[string]rsschema.Attribute{
																"context": rsschema.StringAttribute{
																	Description: "",
																	Required:    true,
																	PlanModifiers: []planmodifier.String{
																		DefaultString(""),
																	},
																	Validators: []validator.String{
																		stringvalidator.LengthAtMost(127),
																	},
																},
																"qualifier": rsschema.ListNestedAttribute{
																	Description: "",
																	Optional:    true,
																	NestedObject: rsschema.NestedAttributeObject{
																		Attributes: map[string]rsschema.Attribute{
																			"name": rsschema.StringAttribute{
																				Description: "",
																				Required:    true,
																				PlanModifiers: []planmodifier.String{
																					DefaultString(""),
																				},
																				Validators: []validator.String{
																					stringvalidator.LengthAtMost(31),
																				},
																			},
																			"value": rsschema.StringAttribute{
																				Description: "",
																				Required:    true,
																				PlanModifiers: []planmodifier.String{
																					DefaultString(""),
																				},
																			},
																		},
																	},
																},
																"value": rsschema.Int64Attribute{
																	Description: "",
																	Required:    true,
																	PlanModifiers: []planmodifier.Int64{
																		DefaultInt64(0),
																	},
																	Validators: []validator.Int64{
																		int64validator.Between(0, 4294967295),
																	},
																},
															},
														},
														"less_than": rsschema.SingleNestedAttribute{
															Description: "",
															Optional:    true,
															Attributes: map[string]rsschema.Attribute{
																"context": rsschema.StringAttribute{
																	Description: "",
																	Required:    true,
																	PlanModifiers: []planmodifier.String{
																		DefaultString(""),
																	},
																	Validators: []validator.String{
																		stringvalidator.LengthAtMost(127),
																	},
																},
																"qualifier": rsschema.ListNestedAttribute{
																	Description: "",
																	Optional:    true,
																	NestedObject: rsschema.NestedAttributeObject{
																		Attributes: map[string]rsschema.Attribute{
																			"name": rsschema.StringAttribute{
																				Description: "",
																				Required:    true,
																				PlanModifiers: []planmodifier.String{
																					DefaultString(""),
																				},
																				Validators: []validator.String{
																					stringvalidator.LengthAtMost(31),
																				},
																			},
																			"value": rsschema.StringAttribute{
																				Description: "",
																				Required:    true,
																				PlanModifiers: []planmodifier.String{
																					DefaultString(""),
																				},
																			},
																		},
																	},
																},
																"value": rsschema.Int64Attribute{
																	Description: "",
																	Required:    true,
																	PlanModifiers: []planmodifier.Int64{
																		DefaultInt64(0),
																	},
																	Validators: []validator.Int64{
																		int64validator.Between(0, 4294967295),
																	},
																},
															},
														},
														"pattern_match": rsschema.SingleNestedAttribute{
															Description: "",
															Optional:    true,
															Attributes: map[string]rsschema.Attribute{
																"context": rsschema.StringAttribute{
																	Description: "",
																	Required:    true,
																	PlanModifiers: []planmodifier.String{
																		DefaultString(""),
																	},
																	Validators: []validator.String{
																		stringvalidator.LengthAtMost(127),
																	},
																},
																"pattern": rsschema.StringAttribute{
																	Description: "",
																	Required:    true,
																	PlanModifiers: []planmodifier.String{
																		DefaultString(""),
																	},
																	Validators: []validator.String{
																		stringvalidator.LengthAtMost(127),
																	},
																},
																"qualifier": rsschema.ListNestedAttribute{
																	Description: "",
																	Optional:    true,
																	NestedObject: rsschema.NestedAttributeObject{
																		Attributes: map[string]rsschema.Attribute{
																			"name": rsschema.StringAttribute{
																				Description: "",
																				Required:    true,
																				PlanModifiers: []planmodifier.String{
																					DefaultString(""),
																				},
																				Validators: []validator.String{
																					stringvalidator.LengthAtMost(31),
																				},
																			},
																			"value": rsschema.StringAttribute{
																				Description: "",
																				Required:    true,
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
									},
								},
							},
						},
						"comment": rsschema.StringAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
							Validators: []validator.String{
								stringvalidator.LengthBetween(0, 256),
							},
						},
						"name": rsschema.StringAttribute{
							Description: "",
							Required:    true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
							Validators: []validator.String{
								stringvalidator.LengthAtMost(31),
							},
						},
						"order_free": rsschema.BoolAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.Bool{
								DefaultBool(false),
							},
						},
						"scope": rsschema.StringAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								DefaultString("protocol-data-unit"),
							},
							Validators: []validator.String{
								stringvalidator.OneOf("protocol-data-unit", "session"),
							},
						},
					},
				},
			},
			"subcategory": rsschema.StringAttribute{
				Description: "",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(63),
				},
			},
			"tcp_half_closed_timeout": rsschema.Int64Attribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					DefaultInt64(0),
				},
				Validators: []validator.Int64{
					int64validator.Between(1, 604800),
				},
			},
			"tcp_time_wait_timeout": rsschema.Int64Attribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					DefaultInt64(0),
				},
				Validators: []validator.Int64{
					int64validator.Between(1, 600),
				},
			},
			"tcp_timeout": rsschema.Int64Attribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					DefaultInt64(0),
				},
				Validators: []validator.Int64{
					int64validator.Between(0, 604800),
				},
			},
			"technology": rsschema.StringAttribute{
				Description: "",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(63),
				},
			},
			"timeout": rsschema.Int64Attribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					DefaultInt64(0),
				},
				Validators: []validator.Int64{
					int64validator.Between(0, 604800),
				},
			},
			"tunnel_applications": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"tunnel_other_application": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"udp_timeout": rsschema.Int64Attribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					DefaultInt64(0),
				},
				Validators: []validator.Int64{
					int64validator.Between(0, 604800),
				},
			},
			"used_by_malware": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"virus_ident": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *objectsApplicationsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *objectsApplicationsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state objectsApplicationsRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "sase_objects_applications",
		"folder":        state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := rrePbcM.NewClient(r.client)
	input := rrePbcM.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 hIzciTY.Config
	var0.AbleToTransferFile = state.AbleToTransferFile.ValueBool()
	var0.AlgDisableCapability = state.AlgDisableCapability.ValueString()
	var0.Category = state.Category.ValueString()
	var0.ConsumeBigBandwidth = state.ConsumeBigBandwidth.ValueBool()
	var0.DataIdent = state.DataIdent.ValueBool()
	var var1 *hIzciTY.DefaultObject
	if state.Default != nil {
		var1 = &hIzciTY.DefaultObject{}
		var var2 *hIzciTY.IdentByIcmp6TypeObject
		if state.Default.IdentByIcmp6Type != nil {
			var2 = &hIzciTY.IdentByIcmp6TypeObject{}
			var2.Code = state.Default.IdentByIcmp6Type.Code.ValueString()
			var2.Type = state.Default.IdentByIcmp6Type.Type.ValueString()
		}
		var1.IdentByIcmp6Type = var2
		var var3 *hIzciTY.IdentByIcmpTypeObject
		if state.Default.IdentByIcmpType != nil {
			var3 = &hIzciTY.IdentByIcmpTypeObject{}
			var3.Code = state.Default.IdentByIcmpType.Code.ValueString()
			var3.Type = state.Default.IdentByIcmpType.Type.ValueString()
		}
		var1.IdentByIcmpType = var3
		var1.IdentByIpProtocol = state.Default.IdentByIpProtocol.ValueString()
		var1.Port = DecodeStringSlice(state.Default.Port)
	}
	var0.Default = var1
	var0.Description = state.Description.ValueString()
	var0.EvasiveBehavior = state.EvasiveBehavior.ValueBool()
	var0.FileTypeIdent = state.FileTypeIdent.ValueBool()
	var0.HasKnownVulnerability = state.HasKnownVulnerability.ValueBool()
	var0.Name = state.Name.ValueString()
	var0.NoAppidCaching = state.NoAppidCaching.ValueBool()
	var0.ParentApp = state.ParentApp.ValueString()
	var0.PervasiveUse = state.PervasiveUse.ValueBool()
	var0.ProneToMisuse = state.ProneToMisuse.ValueBool()
	var0.Risk = state.Risk.ValueInt64()
	var var4 []hIzciTY.SignatureObject
	if len(state.Signature) != 0 {
		var4 = make([]hIzciTY.SignatureObject, 0, len(state.Signature))
		for var5Index := range state.Signature {
			var5 := state.Signature[var5Index]
			var var6 hIzciTY.SignatureObject
			var var7 []hIzciTY.AndConditionObject
			if len(var5.AndCondition) != 0 {
				var7 = make([]hIzciTY.AndConditionObject, 0, len(var5.AndCondition))
				for var8Index := range var5.AndCondition {
					var8 := var5.AndCondition[var8Index]
					var var9 hIzciTY.AndConditionObject
					var9.Name = var8.Name.ValueString()
					var var10 []hIzciTY.OrConditionObject
					if len(var8.OrCondition) != 0 {
						var10 = make([]hIzciTY.OrConditionObject, 0, len(var8.OrCondition))
						for var11Index := range var8.OrCondition {
							var11 := var8.OrCondition[var11Index]
							var var12 hIzciTY.OrConditionObject
							var12.Name = var11.Name.ValueString()
							var var13 hIzciTY.OperatorObject
							var var14 *hIzciTY.EqualToObject
							if var11.Operator.EqualTo != nil {
								var14 = &hIzciTY.EqualToObject{}
								var14.Context = var11.Operator.EqualTo.Context.ValueString()
								var14.Mask = var11.Operator.EqualTo.Mask.ValueString()
								var14.Position = var11.Operator.EqualTo.Position.ValueString()
								var14.Value = var11.Operator.EqualTo.Value.ValueString()
							}
							var13.EqualTo = var14
							var var15 *hIzciTY.GreaterThanObject
							if var11.Operator.GreaterThan != nil {
								var15 = &hIzciTY.GreaterThanObject{}
								var15.Context = var11.Operator.GreaterThan.Context.ValueString()
								var var16 []hIzciTY.QualifierObject
								if len(var11.Operator.GreaterThan.Qualifier) != 0 {
									var16 = make([]hIzciTY.QualifierObject, 0, len(var11.Operator.GreaterThan.Qualifier))
									for var17Index := range var11.Operator.GreaterThan.Qualifier {
										var17 := var11.Operator.GreaterThan.Qualifier[var17Index]
										var var18 hIzciTY.QualifierObject
										var18.Name = var17.Name.ValueString()
										var18.Value = var17.Value.ValueString()
										var16 = append(var16, var18)
									}
								}
								var15.Qualifier = var16
								var15.Value = var11.Operator.GreaterThan.Value.ValueInt64()
							}
							var13.GreaterThan = var15
							var var19 *hIzciTY.LessThanObject
							if var11.Operator.LessThan != nil {
								var19 = &hIzciTY.LessThanObject{}
								var19.Context = var11.Operator.LessThan.Context.ValueString()
								var var20 []hIzciTY.QualifierObject
								if len(var11.Operator.LessThan.Qualifier) != 0 {
									var20 = make([]hIzciTY.QualifierObject, 0, len(var11.Operator.LessThan.Qualifier))
									for var21Index := range var11.Operator.LessThan.Qualifier {
										var21 := var11.Operator.LessThan.Qualifier[var21Index]
										var var22 hIzciTY.QualifierObject
										var22.Name = var21.Name.ValueString()
										var22.Value = var21.Value.ValueString()
										var20 = append(var20, var22)
									}
								}
								var19.Qualifier = var20
								var19.Value = var11.Operator.LessThan.Value.ValueInt64()
							}
							var13.LessThan = var19
							var var23 *hIzciTY.PatternMatchObject
							if var11.Operator.PatternMatch != nil {
								var23 = &hIzciTY.PatternMatchObject{}
								var23.Context = var11.Operator.PatternMatch.Context.ValueString()
								var23.Pattern = var11.Operator.PatternMatch.Pattern.ValueString()
								var var24 []hIzciTY.QualifierObject
								if len(var11.Operator.PatternMatch.Qualifier) != 0 {
									var24 = make([]hIzciTY.QualifierObject, 0, len(var11.Operator.PatternMatch.Qualifier))
									for var25Index := range var11.Operator.PatternMatch.Qualifier {
										var25 := var11.Operator.PatternMatch.Qualifier[var25Index]
										var var26 hIzciTY.QualifierObject
										var26.Name = var25.Name.ValueString()
										var26.Value = var25.Value.ValueString()
										var24 = append(var24, var26)
									}
								}
								var23.Qualifier = var24
							}
							var13.PatternMatch = var23
							var12.Operator = var13
							var10 = append(var10, var12)
						}
					}
					var9.OrCondition = var10
					var7 = append(var7, var9)
				}
			}
			var6.AndCondition = var7
			var6.Comment = var5.Comment.ValueString()
			var6.Name = var5.Name.ValueString()
			var6.OrderFree = var5.OrderFree.ValueBool()
			var6.Scope = var5.Scope.ValueString()
			var4 = append(var4, var6)
		}
	}
	var0.Signature = var4
	var0.Subcategory = state.Subcategory.ValueString()
	var0.TcpHalfClosedTimeout = state.TcpHalfClosedTimeout.ValueInt64()
	var0.TcpTimeWaitTimeout = state.TcpTimeWaitTimeout.ValueInt64()
	var0.TcpTimeout = state.TcpTimeout.ValueInt64()
	var0.Technology = state.Technology.ValueString()
	var0.Timeout = state.Timeout.ValueInt64()
	var0.TunnelApplications = state.TunnelApplications.ValueBool()
	var0.TunnelOtherApplication = state.TunnelOtherApplication.ValueBool()
	var0.UdpTimeout = state.UdpTimeout.ValueInt64()
	var0.UsedByMalware = state.UsedByMalware.ValueBool()
	var0.VirusIdent = state.VirusIdent.ValueBool()
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Create(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in create", err.Error())
		return
	}

	// Store the answer to state.
	state.Id = types.StringValue(strings.Join([]string{input.Folder, ans.ObjectId}, IdSeparator))
	var var27 *objectsApplicationsRsModelDefaultObject
	if ans.Default != nil {
		var27 = &objectsApplicationsRsModelDefaultObject{}
		var var28 *objectsApplicationsRsModelIdentByIcmp6TypeObject
		if ans.Default.IdentByIcmp6Type != nil {
			var28 = &objectsApplicationsRsModelIdentByIcmp6TypeObject{}
			var28.Code = types.StringValue(ans.Default.IdentByIcmp6Type.Code)
			var28.Type = types.StringValue(ans.Default.IdentByIcmp6Type.Type)
		}
		var var29 *objectsApplicationsRsModelIdentByIcmpTypeObject
		if ans.Default.IdentByIcmpType != nil {
			var29 = &objectsApplicationsRsModelIdentByIcmpTypeObject{}
			var29.Code = types.StringValue(ans.Default.IdentByIcmpType.Code)
			var29.Type = types.StringValue(ans.Default.IdentByIcmpType.Type)
		}
		var27.IdentByIcmp6Type = var28
		var27.IdentByIcmpType = var29
		var27.IdentByIpProtocol = types.StringValue(ans.Default.IdentByIpProtocol)
		var27.Port = EncodeStringSlice(ans.Default.Port)
	}
	var var30 []objectsApplicationsRsModelSignatureObject
	if len(ans.Signature) != 0 {
		var30 = make([]objectsApplicationsRsModelSignatureObject, 0, len(ans.Signature))
		for var31Index := range ans.Signature {
			var31 := ans.Signature[var31Index]
			var var32 objectsApplicationsRsModelSignatureObject
			var var33 []objectsApplicationsRsModelAndConditionObject
			if len(var31.AndCondition) != 0 {
				var33 = make([]objectsApplicationsRsModelAndConditionObject, 0, len(var31.AndCondition))
				for var34Index := range var31.AndCondition {
					var34 := var31.AndCondition[var34Index]
					var var35 objectsApplicationsRsModelAndConditionObject
					var var36 []objectsApplicationsRsModelOrConditionObject
					if len(var34.OrCondition) != 0 {
						var36 = make([]objectsApplicationsRsModelOrConditionObject, 0, len(var34.OrCondition))
						for var37Index := range var34.OrCondition {
							var37 := var34.OrCondition[var37Index]
							var var38 objectsApplicationsRsModelOrConditionObject
							var var39 objectsApplicationsRsModelOperatorObject
							var var40 *objectsApplicationsRsModelEqualToObject
							if var37.Operator.EqualTo != nil {
								var40 = &objectsApplicationsRsModelEqualToObject{}
								var40.Context = types.StringValue(var37.Operator.EqualTo.Context)
								var40.Mask = types.StringValue(var37.Operator.EqualTo.Mask)
								var40.Position = types.StringValue(var37.Operator.EqualTo.Position)
								var40.Value = types.StringValue(var37.Operator.EqualTo.Value)
							}
							var var41 *objectsApplicationsRsModelGreaterThanObject
							if var37.Operator.GreaterThan != nil {
								var41 = &objectsApplicationsRsModelGreaterThanObject{}
								var var42 []objectsApplicationsRsModelQualifierObject
								if len(var37.Operator.GreaterThan.Qualifier) != 0 {
									var42 = make([]objectsApplicationsRsModelQualifierObject, 0, len(var37.Operator.GreaterThan.Qualifier))
									for var43Index := range var37.Operator.GreaterThan.Qualifier {
										var43 := var37.Operator.GreaterThan.Qualifier[var43Index]
										var var44 objectsApplicationsRsModelQualifierObject
										var44.Name = types.StringValue(var43.Name)
										var44.Value = types.StringValue(var43.Value)
										var42 = append(var42, var44)
									}
								}
								var41.Context = types.StringValue(var37.Operator.GreaterThan.Context)
								var41.Qualifier = var42
								var41.Value = types.Int64Value(var37.Operator.GreaterThan.Value)
							}
							var var45 *objectsApplicationsRsModelLessThanObject
							if var37.Operator.LessThan != nil {
								var45 = &objectsApplicationsRsModelLessThanObject{}
								var var46 []objectsApplicationsRsModelQualifierObject
								if len(var37.Operator.LessThan.Qualifier) != 0 {
									var46 = make([]objectsApplicationsRsModelQualifierObject, 0, len(var37.Operator.LessThan.Qualifier))
									for var47Index := range var37.Operator.LessThan.Qualifier {
										var47 := var37.Operator.LessThan.Qualifier[var47Index]
										var var48 objectsApplicationsRsModelQualifierObject
										var48.Name = types.StringValue(var47.Name)
										var48.Value = types.StringValue(var47.Value)
										var46 = append(var46, var48)
									}
								}
								var45.Context = types.StringValue(var37.Operator.LessThan.Context)
								var45.Qualifier = var46
								var45.Value = types.Int64Value(var37.Operator.LessThan.Value)
							}
							var var49 *objectsApplicationsRsModelPatternMatchObject
							if var37.Operator.PatternMatch != nil {
								var49 = &objectsApplicationsRsModelPatternMatchObject{}
								var var50 []objectsApplicationsRsModelQualifierObject
								if len(var37.Operator.PatternMatch.Qualifier) != 0 {
									var50 = make([]objectsApplicationsRsModelQualifierObject, 0, len(var37.Operator.PatternMatch.Qualifier))
									for var51Index := range var37.Operator.PatternMatch.Qualifier {
										var51 := var37.Operator.PatternMatch.Qualifier[var51Index]
										var var52 objectsApplicationsRsModelQualifierObject
										var52.Name = types.StringValue(var51.Name)
										var52.Value = types.StringValue(var51.Value)
										var50 = append(var50, var52)
									}
								}
								var49.Context = types.StringValue(var37.Operator.PatternMatch.Context)
								var49.Pattern = types.StringValue(var37.Operator.PatternMatch.Pattern)
								var49.Qualifier = var50
							}
							var39.EqualTo = var40
							var39.GreaterThan = var41
							var39.LessThan = var45
							var39.PatternMatch = var49
							var38.Name = types.StringValue(var37.Name)
							var38.Operator = var39
							var36 = append(var36, var38)
						}
					}
					var35.Name = types.StringValue(var34.Name)
					var35.OrCondition = var36
					var33 = append(var33, var35)
				}
			}
			var32.AndCondition = var33
			var32.Comment = types.StringValue(var31.Comment)
			var32.Name = types.StringValue(var31.Name)
			var32.OrderFree = types.BoolValue(var31.OrderFree)
			var32.Scope = types.StringValue(var31.Scope)
			var30 = append(var30, var32)
		}
	}
	state.AbleToTransferFile = types.BoolValue(ans.AbleToTransferFile)
	state.AlgDisableCapability = types.StringValue(ans.AlgDisableCapability)
	state.Category = types.StringValue(ans.Category)
	state.ConsumeBigBandwidth = types.BoolValue(ans.ConsumeBigBandwidth)
	state.DataIdent = types.BoolValue(ans.DataIdent)
	state.Default = var27
	state.Description = types.StringValue(ans.Description)
	state.EvasiveBehavior = types.BoolValue(ans.EvasiveBehavior)
	state.FileTypeIdent = types.BoolValue(ans.FileTypeIdent)
	state.HasKnownVulnerability = types.BoolValue(ans.HasKnownVulnerability)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.NoAppidCaching = types.BoolValue(ans.NoAppidCaching)
	state.ParentApp = types.StringValue(ans.ParentApp)
	state.PervasiveUse = types.BoolValue(ans.PervasiveUse)
	state.ProneToMisuse = types.BoolValue(ans.ProneToMisuse)
	state.Risk = types.Int64Value(ans.Risk)
	state.Signature = var30
	state.Subcategory = types.StringValue(ans.Subcategory)
	state.TcpHalfClosedTimeout = types.Int64Value(ans.TcpHalfClosedTimeout)
	state.TcpTimeWaitTimeout = types.Int64Value(ans.TcpTimeWaitTimeout)
	state.TcpTimeout = types.Int64Value(ans.TcpTimeout)
	state.Technology = types.StringValue(ans.Technology)
	state.Timeout = types.Int64Value(ans.Timeout)
	state.TunnelApplications = types.BoolValue(ans.TunnelApplications)
	state.TunnelOtherApplication = types.BoolValue(ans.TunnelOtherApplication)
	state.UdpTimeout = types.Int64Value(ans.UdpTimeout)
	state.UsedByMalware = types.BoolValue(ans.UsedByMalware)
	state.VirusIdent = types.BoolValue(ans.VirusIdent)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *objectsApplicationsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var idType types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &idType)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := idType.ValueString()
	tokens := strings.Split(id, IdSeparator)
	if len(tokens) != 2 {
		resp.Diagnostics.AddError("Error in resource ID format", "Expected 2 tokens")
		return
	}

	var state objectsApplicationsRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "sase_objects_applications",
		"locMap":        map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":        tokens,
	})

	// Prepare to read the config.
	svc := rrePbcM.NewClient(r.client)
	input := rrePbcM.ReadInput{
		ObjectId: tokens[1],
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
	state.Id = idType
	var var0 *objectsApplicationsRsModelDefaultObject
	if ans.Default != nil {
		var0 = &objectsApplicationsRsModelDefaultObject{}
		var var1 *objectsApplicationsRsModelIdentByIcmp6TypeObject
		if ans.Default.IdentByIcmp6Type != nil {
			var1 = &objectsApplicationsRsModelIdentByIcmp6TypeObject{}
			var1.Code = types.StringValue(ans.Default.IdentByIcmp6Type.Code)
			var1.Type = types.StringValue(ans.Default.IdentByIcmp6Type.Type)
		}
		var var2 *objectsApplicationsRsModelIdentByIcmpTypeObject
		if ans.Default.IdentByIcmpType != nil {
			var2 = &objectsApplicationsRsModelIdentByIcmpTypeObject{}
			var2.Code = types.StringValue(ans.Default.IdentByIcmpType.Code)
			var2.Type = types.StringValue(ans.Default.IdentByIcmpType.Type)
		}
		var0.IdentByIcmp6Type = var1
		var0.IdentByIcmpType = var2
		var0.IdentByIpProtocol = types.StringValue(ans.Default.IdentByIpProtocol)
		var0.Port = EncodeStringSlice(ans.Default.Port)
	}
	var var3 []objectsApplicationsRsModelSignatureObject
	if len(ans.Signature) != 0 {
		var3 = make([]objectsApplicationsRsModelSignatureObject, 0, len(ans.Signature))
		for var4Index := range ans.Signature {
			var4 := ans.Signature[var4Index]
			var var5 objectsApplicationsRsModelSignatureObject
			var var6 []objectsApplicationsRsModelAndConditionObject
			if len(var4.AndCondition) != 0 {
				var6 = make([]objectsApplicationsRsModelAndConditionObject, 0, len(var4.AndCondition))
				for var7Index := range var4.AndCondition {
					var7 := var4.AndCondition[var7Index]
					var var8 objectsApplicationsRsModelAndConditionObject
					var var9 []objectsApplicationsRsModelOrConditionObject
					if len(var7.OrCondition) != 0 {
						var9 = make([]objectsApplicationsRsModelOrConditionObject, 0, len(var7.OrCondition))
						for var10Index := range var7.OrCondition {
							var10 := var7.OrCondition[var10Index]
							var var11 objectsApplicationsRsModelOrConditionObject
							var var12 objectsApplicationsRsModelOperatorObject
							var var13 *objectsApplicationsRsModelEqualToObject
							if var10.Operator.EqualTo != nil {
								var13 = &objectsApplicationsRsModelEqualToObject{}
								var13.Context = types.StringValue(var10.Operator.EqualTo.Context)
								var13.Mask = types.StringValue(var10.Operator.EqualTo.Mask)
								var13.Position = types.StringValue(var10.Operator.EqualTo.Position)
								var13.Value = types.StringValue(var10.Operator.EqualTo.Value)
							}
							var var14 *objectsApplicationsRsModelGreaterThanObject
							if var10.Operator.GreaterThan != nil {
								var14 = &objectsApplicationsRsModelGreaterThanObject{}
								var var15 []objectsApplicationsRsModelQualifierObject
								if len(var10.Operator.GreaterThan.Qualifier) != 0 {
									var15 = make([]objectsApplicationsRsModelQualifierObject, 0, len(var10.Operator.GreaterThan.Qualifier))
									for var16Index := range var10.Operator.GreaterThan.Qualifier {
										var16 := var10.Operator.GreaterThan.Qualifier[var16Index]
										var var17 objectsApplicationsRsModelQualifierObject
										var17.Name = types.StringValue(var16.Name)
										var17.Value = types.StringValue(var16.Value)
										var15 = append(var15, var17)
									}
								}
								var14.Context = types.StringValue(var10.Operator.GreaterThan.Context)
								var14.Qualifier = var15
								var14.Value = types.Int64Value(var10.Operator.GreaterThan.Value)
							}
							var var18 *objectsApplicationsRsModelLessThanObject
							if var10.Operator.LessThan != nil {
								var18 = &objectsApplicationsRsModelLessThanObject{}
								var var19 []objectsApplicationsRsModelQualifierObject
								if len(var10.Operator.LessThan.Qualifier) != 0 {
									var19 = make([]objectsApplicationsRsModelQualifierObject, 0, len(var10.Operator.LessThan.Qualifier))
									for var20Index := range var10.Operator.LessThan.Qualifier {
										var20 := var10.Operator.LessThan.Qualifier[var20Index]
										var var21 objectsApplicationsRsModelQualifierObject
										var21.Name = types.StringValue(var20.Name)
										var21.Value = types.StringValue(var20.Value)
										var19 = append(var19, var21)
									}
								}
								var18.Context = types.StringValue(var10.Operator.LessThan.Context)
								var18.Qualifier = var19
								var18.Value = types.Int64Value(var10.Operator.LessThan.Value)
							}
							var var22 *objectsApplicationsRsModelPatternMatchObject
							if var10.Operator.PatternMatch != nil {
								var22 = &objectsApplicationsRsModelPatternMatchObject{}
								var var23 []objectsApplicationsRsModelQualifierObject
								if len(var10.Operator.PatternMatch.Qualifier) != 0 {
									var23 = make([]objectsApplicationsRsModelQualifierObject, 0, len(var10.Operator.PatternMatch.Qualifier))
									for var24Index := range var10.Operator.PatternMatch.Qualifier {
										var24 := var10.Operator.PatternMatch.Qualifier[var24Index]
										var var25 objectsApplicationsRsModelQualifierObject
										var25.Name = types.StringValue(var24.Name)
										var25.Value = types.StringValue(var24.Value)
										var23 = append(var23, var25)
									}
								}
								var22.Context = types.StringValue(var10.Operator.PatternMatch.Context)
								var22.Pattern = types.StringValue(var10.Operator.PatternMatch.Pattern)
								var22.Qualifier = var23
							}
							var12.EqualTo = var13
							var12.GreaterThan = var14
							var12.LessThan = var18
							var12.PatternMatch = var22
							var11.Name = types.StringValue(var10.Name)
							var11.Operator = var12
							var9 = append(var9, var11)
						}
					}
					var8.Name = types.StringValue(var7.Name)
					var8.OrCondition = var9
					var6 = append(var6, var8)
				}
			}
			var5.AndCondition = var6
			var5.Comment = types.StringValue(var4.Comment)
			var5.Name = types.StringValue(var4.Name)
			var5.OrderFree = types.BoolValue(var4.OrderFree)
			var5.Scope = types.StringValue(var4.Scope)
			var3 = append(var3, var5)
		}
	}
	state.AbleToTransferFile = types.BoolValue(ans.AbleToTransferFile)
	state.AlgDisableCapability = types.StringValue(ans.AlgDisableCapability)
	state.Category = types.StringValue(ans.Category)
	state.ConsumeBigBandwidth = types.BoolValue(ans.ConsumeBigBandwidth)
	state.DataIdent = types.BoolValue(ans.DataIdent)
	state.Default = var0
	state.Description = types.StringValue(ans.Description)
	state.EvasiveBehavior = types.BoolValue(ans.EvasiveBehavior)
	state.FileTypeIdent = types.BoolValue(ans.FileTypeIdent)
	state.HasKnownVulnerability = types.BoolValue(ans.HasKnownVulnerability)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.NoAppidCaching = types.BoolValue(ans.NoAppidCaching)
	state.ParentApp = types.StringValue(ans.ParentApp)
	state.PervasiveUse = types.BoolValue(ans.PervasiveUse)
	state.ProneToMisuse = types.BoolValue(ans.ProneToMisuse)
	state.Risk = types.Int64Value(ans.Risk)
	state.Signature = var3
	state.Subcategory = types.StringValue(ans.Subcategory)
	state.TcpHalfClosedTimeout = types.Int64Value(ans.TcpHalfClosedTimeout)
	state.TcpTimeWaitTimeout = types.Int64Value(ans.TcpTimeWaitTimeout)
	state.TcpTimeout = types.Int64Value(ans.TcpTimeout)
	state.Technology = types.StringValue(ans.Technology)
	state.Timeout = types.Int64Value(ans.Timeout)
	state.TunnelApplications = types.BoolValue(ans.TunnelApplications)
	state.TunnelOtherApplication = types.BoolValue(ans.TunnelOtherApplication)
	state.UdpTimeout = types.Int64Value(ans.UdpTimeout)
	state.UsedByMalware = types.BoolValue(ans.UsedByMalware)
	state.VirusIdent = types.BoolValue(ans.VirusIdent)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *objectsApplicationsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state objectsApplicationsRsModel
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
		"resource_name": "sase_objects_applications",
		"object_id":     state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := rrePbcM.NewClient(r.client)
	input := rrePbcM.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 hIzciTY.Config
	var0.AbleToTransferFile = plan.AbleToTransferFile.ValueBool()
	var0.AlgDisableCapability = plan.AlgDisableCapability.ValueString()
	var0.Category = plan.Category.ValueString()
	var0.ConsumeBigBandwidth = plan.ConsumeBigBandwidth.ValueBool()
	var0.DataIdent = plan.DataIdent.ValueBool()
	var var1 *hIzciTY.DefaultObject
	if plan.Default != nil {
		var1 = &hIzciTY.DefaultObject{}
		var var2 *hIzciTY.IdentByIcmp6TypeObject
		if plan.Default.IdentByIcmp6Type != nil {
			var2 = &hIzciTY.IdentByIcmp6TypeObject{}
			var2.Code = plan.Default.IdentByIcmp6Type.Code.ValueString()
			var2.Type = plan.Default.IdentByIcmp6Type.Type.ValueString()
		}
		var1.IdentByIcmp6Type = var2
		var var3 *hIzciTY.IdentByIcmpTypeObject
		if plan.Default.IdentByIcmpType != nil {
			var3 = &hIzciTY.IdentByIcmpTypeObject{}
			var3.Code = plan.Default.IdentByIcmpType.Code.ValueString()
			var3.Type = plan.Default.IdentByIcmpType.Type.ValueString()
		}
		var1.IdentByIcmpType = var3
		var1.IdentByIpProtocol = plan.Default.IdentByIpProtocol.ValueString()
		var1.Port = DecodeStringSlice(plan.Default.Port)
	}
	var0.Default = var1
	var0.Description = plan.Description.ValueString()
	var0.EvasiveBehavior = plan.EvasiveBehavior.ValueBool()
	var0.FileTypeIdent = plan.FileTypeIdent.ValueBool()
	var0.HasKnownVulnerability = plan.HasKnownVulnerability.ValueBool()
	var0.Name = plan.Name.ValueString()
	var0.NoAppidCaching = plan.NoAppidCaching.ValueBool()
	var0.ParentApp = plan.ParentApp.ValueString()
	var0.PervasiveUse = plan.PervasiveUse.ValueBool()
	var0.ProneToMisuse = plan.ProneToMisuse.ValueBool()
	var0.Risk = plan.Risk.ValueInt64()
	var var4 []hIzciTY.SignatureObject
	if len(plan.Signature) != 0 {
		var4 = make([]hIzciTY.SignatureObject, 0, len(plan.Signature))
		for var5Index := range plan.Signature {
			var5 := plan.Signature[var5Index]
			var var6 hIzciTY.SignatureObject
			var var7 []hIzciTY.AndConditionObject
			if len(var5.AndCondition) != 0 {
				var7 = make([]hIzciTY.AndConditionObject, 0, len(var5.AndCondition))
				for var8Index := range var5.AndCondition {
					var8 := var5.AndCondition[var8Index]
					var var9 hIzciTY.AndConditionObject
					var9.Name = var8.Name.ValueString()
					var var10 []hIzciTY.OrConditionObject
					if len(var8.OrCondition) != 0 {
						var10 = make([]hIzciTY.OrConditionObject, 0, len(var8.OrCondition))
						for var11Index := range var8.OrCondition {
							var11 := var8.OrCondition[var11Index]
							var var12 hIzciTY.OrConditionObject
							var12.Name = var11.Name.ValueString()
							var var13 hIzciTY.OperatorObject
							var var14 *hIzciTY.EqualToObject
							if var11.Operator.EqualTo != nil {
								var14 = &hIzciTY.EqualToObject{}
								var14.Context = var11.Operator.EqualTo.Context.ValueString()
								var14.Mask = var11.Operator.EqualTo.Mask.ValueString()
								var14.Position = var11.Operator.EqualTo.Position.ValueString()
								var14.Value = var11.Operator.EqualTo.Value.ValueString()
							}
							var13.EqualTo = var14
							var var15 *hIzciTY.GreaterThanObject
							if var11.Operator.GreaterThan != nil {
								var15 = &hIzciTY.GreaterThanObject{}
								var15.Context = var11.Operator.GreaterThan.Context.ValueString()
								var var16 []hIzciTY.QualifierObject
								if len(var11.Operator.GreaterThan.Qualifier) != 0 {
									var16 = make([]hIzciTY.QualifierObject, 0, len(var11.Operator.GreaterThan.Qualifier))
									for var17Index := range var11.Operator.GreaterThan.Qualifier {
										var17 := var11.Operator.GreaterThan.Qualifier[var17Index]
										var var18 hIzciTY.QualifierObject
										var18.Name = var17.Name.ValueString()
										var18.Value = var17.Value.ValueString()
										var16 = append(var16, var18)
									}
								}
								var15.Qualifier = var16
								var15.Value = var11.Operator.GreaterThan.Value.ValueInt64()
							}
							var13.GreaterThan = var15
							var var19 *hIzciTY.LessThanObject
							if var11.Operator.LessThan != nil {
								var19 = &hIzciTY.LessThanObject{}
								var19.Context = var11.Operator.LessThan.Context.ValueString()
								var var20 []hIzciTY.QualifierObject
								if len(var11.Operator.LessThan.Qualifier) != 0 {
									var20 = make([]hIzciTY.QualifierObject, 0, len(var11.Operator.LessThan.Qualifier))
									for var21Index := range var11.Operator.LessThan.Qualifier {
										var21 := var11.Operator.LessThan.Qualifier[var21Index]
										var var22 hIzciTY.QualifierObject
										var22.Name = var21.Name.ValueString()
										var22.Value = var21.Value.ValueString()
										var20 = append(var20, var22)
									}
								}
								var19.Qualifier = var20
								var19.Value = var11.Operator.LessThan.Value.ValueInt64()
							}
							var13.LessThan = var19
							var var23 *hIzciTY.PatternMatchObject
							if var11.Operator.PatternMatch != nil {
								var23 = &hIzciTY.PatternMatchObject{}
								var23.Context = var11.Operator.PatternMatch.Context.ValueString()
								var23.Pattern = var11.Operator.PatternMatch.Pattern.ValueString()
								var var24 []hIzciTY.QualifierObject
								if len(var11.Operator.PatternMatch.Qualifier) != 0 {
									var24 = make([]hIzciTY.QualifierObject, 0, len(var11.Operator.PatternMatch.Qualifier))
									for var25Index := range var11.Operator.PatternMatch.Qualifier {
										var25 := var11.Operator.PatternMatch.Qualifier[var25Index]
										var var26 hIzciTY.QualifierObject
										var26.Name = var25.Name.ValueString()
										var26.Value = var25.Value.ValueString()
										var24 = append(var24, var26)
									}
								}
								var23.Qualifier = var24
							}
							var13.PatternMatch = var23
							var12.Operator = var13
							var10 = append(var10, var12)
						}
					}
					var9.OrCondition = var10
					var7 = append(var7, var9)
				}
			}
			var6.AndCondition = var7
			var6.Comment = var5.Comment.ValueString()
			var6.Name = var5.Name.ValueString()
			var6.OrderFree = var5.OrderFree.ValueBool()
			var6.Scope = var5.Scope.ValueString()
			var4 = append(var4, var6)
		}
	}
	var0.Signature = var4
	var0.Subcategory = plan.Subcategory.ValueString()
	var0.TcpHalfClosedTimeout = plan.TcpHalfClosedTimeout.ValueInt64()
	var0.TcpTimeWaitTimeout = plan.TcpTimeWaitTimeout.ValueInt64()
	var0.TcpTimeout = plan.TcpTimeout.ValueInt64()
	var0.Technology = plan.Technology.ValueString()
	var0.Timeout = plan.Timeout.ValueInt64()
	var0.TunnelApplications = plan.TunnelApplications.ValueBool()
	var0.TunnelOtherApplication = plan.TunnelOtherApplication.ValueBool()
	var0.UdpTimeout = plan.UdpTimeout.ValueInt64()
	var0.UsedByMalware = plan.UsedByMalware.ValueBool()
	var0.VirusIdent = plan.VirusIdent.ValueBool()
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var27 *objectsApplicationsRsModelDefaultObject
	if ans.Default != nil {
		var27 = &objectsApplicationsRsModelDefaultObject{}
		var var28 *objectsApplicationsRsModelIdentByIcmp6TypeObject
		if ans.Default.IdentByIcmp6Type != nil {
			var28 = &objectsApplicationsRsModelIdentByIcmp6TypeObject{}
			var28.Code = types.StringValue(ans.Default.IdentByIcmp6Type.Code)
			var28.Type = types.StringValue(ans.Default.IdentByIcmp6Type.Type)
		}
		var var29 *objectsApplicationsRsModelIdentByIcmpTypeObject
		if ans.Default.IdentByIcmpType != nil {
			var29 = &objectsApplicationsRsModelIdentByIcmpTypeObject{}
			var29.Code = types.StringValue(ans.Default.IdentByIcmpType.Code)
			var29.Type = types.StringValue(ans.Default.IdentByIcmpType.Type)
		}
		var27.IdentByIcmp6Type = var28
		var27.IdentByIcmpType = var29
		var27.IdentByIpProtocol = types.StringValue(ans.Default.IdentByIpProtocol)
		var27.Port = EncodeStringSlice(ans.Default.Port)
	}
	var var30 []objectsApplicationsRsModelSignatureObject
	if len(ans.Signature) != 0 {
		var30 = make([]objectsApplicationsRsModelSignatureObject, 0, len(ans.Signature))
		for var31Index := range ans.Signature {
			var31 := ans.Signature[var31Index]
			var var32 objectsApplicationsRsModelSignatureObject
			var var33 []objectsApplicationsRsModelAndConditionObject
			if len(var31.AndCondition) != 0 {
				var33 = make([]objectsApplicationsRsModelAndConditionObject, 0, len(var31.AndCondition))
				for var34Index := range var31.AndCondition {
					var34 := var31.AndCondition[var34Index]
					var var35 objectsApplicationsRsModelAndConditionObject
					var var36 []objectsApplicationsRsModelOrConditionObject
					if len(var34.OrCondition) != 0 {
						var36 = make([]objectsApplicationsRsModelOrConditionObject, 0, len(var34.OrCondition))
						for var37Index := range var34.OrCondition {
							var37 := var34.OrCondition[var37Index]
							var var38 objectsApplicationsRsModelOrConditionObject
							var var39 objectsApplicationsRsModelOperatorObject
							var var40 *objectsApplicationsRsModelEqualToObject
							if var37.Operator.EqualTo != nil {
								var40 = &objectsApplicationsRsModelEqualToObject{}
								var40.Context = types.StringValue(var37.Operator.EqualTo.Context)
								var40.Mask = types.StringValue(var37.Operator.EqualTo.Mask)
								var40.Position = types.StringValue(var37.Operator.EqualTo.Position)
								var40.Value = types.StringValue(var37.Operator.EqualTo.Value)
							}
							var var41 *objectsApplicationsRsModelGreaterThanObject
							if var37.Operator.GreaterThan != nil {
								var41 = &objectsApplicationsRsModelGreaterThanObject{}
								var var42 []objectsApplicationsRsModelQualifierObject
								if len(var37.Operator.GreaterThan.Qualifier) != 0 {
									var42 = make([]objectsApplicationsRsModelQualifierObject, 0, len(var37.Operator.GreaterThan.Qualifier))
									for var43Index := range var37.Operator.GreaterThan.Qualifier {
										var43 := var37.Operator.GreaterThan.Qualifier[var43Index]
										var var44 objectsApplicationsRsModelQualifierObject
										var44.Name = types.StringValue(var43.Name)
										var44.Value = types.StringValue(var43.Value)
										var42 = append(var42, var44)
									}
								}
								var41.Context = types.StringValue(var37.Operator.GreaterThan.Context)
								var41.Qualifier = var42
								var41.Value = types.Int64Value(var37.Operator.GreaterThan.Value)
							}
							var var45 *objectsApplicationsRsModelLessThanObject
							if var37.Operator.LessThan != nil {
								var45 = &objectsApplicationsRsModelLessThanObject{}
								var var46 []objectsApplicationsRsModelQualifierObject
								if len(var37.Operator.LessThan.Qualifier) != 0 {
									var46 = make([]objectsApplicationsRsModelQualifierObject, 0, len(var37.Operator.LessThan.Qualifier))
									for var47Index := range var37.Operator.LessThan.Qualifier {
										var47 := var37.Operator.LessThan.Qualifier[var47Index]
										var var48 objectsApplicationsRsModelQualifierObject
										var48.Name = types.StringValue(var47.Name)
										var48.Value = types.StringValue(var47.Value)
										var46 = append(var46, var48)
									}
								}
								var45.Context = types.StringValue(var37.Operator.LessThan.Context)
								var45.Qualifier = var46
								var45.Value = types.Int64Value(var37.Operator.LessThan.Value)
							}
							var var49 *objectsApplicationsRsModelPatternMatchObject
							if var37.Operator.PatternMatch != nil {
								var49 = &objectsApplicationsRsModelPatternMatchObject{}
								var var50 []objectsApplicationsRsModelQualifierObject
								if len(var37.Operator.PatternMatch.Qualifier) != 0 {
									var50 = make([]objectsApplicationsRsModelQualifierObject, 0, len(var37.Operator.PatternMatch.Qualifier))
									for var51Index := range var37.Operator.PatternMatch.Qualifier {
										var51 := var37.Operator.PatternMatch.Qualifier[var51Index]
										var var52 objectsApplicationsRsModelQualifierObject
										var52.Name = types.StringValue(var51.Name)
										var52.Value = types.StringValue(var51.Value)
										var50 = append(var50, var52)
									}
								}
								var49.Context = types.StringValue(var37.Operator.PatternMatch.Context)
								var49.Pattern = types.StringValue(var37.Operator.PatternMatch.Pattern)
								var49.Qualifier = var50
							}
							var39.EqualTo = var40
							var39.GreaterThan = var41
							var39.LessThan = var45
							var39.PatternMatch = var49
							var38.Name = types.StringValue(var37.Name)
							var38.Operator = var39
							var36 = append(var36, var38)
						}
					}
					var35.Name = types.StringValue(var34.Name)
					var35.OrCondition = var36
					var33 = append(var33, var35)
				}
			}
			var32.AndCondition = var33
			var32.Comment = types.StringValue(var31.Comment)
			var32.Name = types.StringValue(var31.Name)
			var32.OrderFree = types.BoolValue(var31.OrderFree)
			var32.Scope = types.StringValue(var31.Scope)
			var30 = append(var30, var32)
		}
	}
	state.AbleToTransferFile = types.BoolValue(ans.AbleToTransferFile)
	state.AlgDisableCapability = types.StringValue(ans.AlgDisableCapability)
	state.Category = types.StringValue(ans.Category)
	state.ConsumeBigBandwidth = types.BoolValue(ans.ConsumeBigBandwidth)
	state.DataIdent = types.BoolValue(ans.DataIdent)
	state.Default = var27
	state.Description = types.StringValue(ans.Description)
	state.EvasiveBehavior = types.BoolValue(ans.EvasiveBehavior)
	state.FileTypeIdent = types.BoolValue(ans.FileTypeIdent)
	state.HasKnownVulnerability = types.BoolValue(ans.HasKnownVulnerability)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.NoAppidCaching = types.BoolValue(ans.NoAppidCaching)
	state.ParentApp = types.StringValue(ans.ParentApp)
	state.PervasiveUse = types.BoolValue(ans.PervasiveUse)
	state.ProneToMisuse = types.BoolValue(ans.ProneToMisuse)
	state.Risk = types.Int64Value(ans.Risk)
	state.Signature = var30
	state.Subcategory = types.StringValue(ans.Subcategory)
	state.TcpHalfClosedTimeout = types.Int64Value(ans.TcpHalfClosedTimeout)
	state.TcpTimeWaitTimeout = types.Int64Value(ans.TcpTimeWaitTimeout)
	state.TcpTimeout = types.Int64Value(ans.TcpTimeout)
	state.Technology = types.StringValue(ans.Technology)
	state.Timeout = types.Int64Value(ans.Timeout)
	state.TunnelApplications = types.BoolValue(ans.TunnelApplications)
	state.TunnelOtherApplication = types.BoolValue(ans.TunnelOtherApplication)
	state.UdpTimeout = types.Int64Value(ans.UdpTimeout)
	state.UsedByMalware = types.BoolValue(ans.UsedByMalware)
	state.VirusIdent = types.BoolValue(ans.VirusIdent)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *objectsApplicationsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var idType types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &idType)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := idType.ValueString()
	tokens := strings.Split(id, IdSeparator)
	if len(tokens) != 2 {
		resp.Diagnostics.AddError("Error in resource ID format", "Expected 2 tokens")
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource delete", map[string]any{
		"resource_name": "sase_objects_applications",
		"locMap":        map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":        tokens,
	})

	svc := rrePbcM.NewClient(r.client)
	input := rrePbcM.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *objectsApplicationsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
