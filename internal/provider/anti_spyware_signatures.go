package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	iblCTtp "github.com/paloaltonetworks/sase-go/netsec/service/v1/antispywaresignatures"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Data source (listing).
var (
	_ datasource.DataSource              = &antiSpywareSignaturesListDataSource{}
	_ datasource.DataSourceWithConfigure = &antiSpywareSignaturesListDataSource{}
)

func NewAntiSpywareSignaturesListDataSource() datasource.DataSource {
	return &antiSpywareSignaturesListDataSource{}
}

type antiSpywareSignaturesListDataSource struct {
	client *sase.Client
}

type antiSpywareSignaturesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []antiSpywareSignaturesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type antiSpywareSignaturesListDsModelConfig struct {
	Bugtraq       []types.String                                       `tfsdk:"bugtraq"`
	Comment       types.String                                         `tfsdk:"comment"`
	Cve           []types.String                                       `tfsdk:"cve"`
	DefaultAction *antiSpywareSignaturesListDsModelDefaultActionObject `tfsdk:"default_action"`
	Direction     types.String                                         `tfsdk:"direction"`
	ObjectId      types.String                                         `tfsdk:"object_id"`
	Reference     []types.String                                       `tfsdk:"reference"`
	Severity      types.String                                         `tfsdk:"severity"`
	Signature     *antiSpywareSignaturesListDsModelSignatureObject     `tfsdk:"signature"`
	ThreatId      types.Int64                                          `tfsdk:"threat_id"`
	Threatname    types.String                                         `tfsdk:"threatname"`
	Vendor        []types.String                                       `tfsdk:"vendor"`
}

type antiSpywareSignaturesListDsModelDefaultActionObject struct {
	Alert       types.Bool                                     `tfsdk:"alert"`
	Allow       types.Bool                                     `tfsdk:"allow"`
	BlockIp     *antiSpywareSignaturesListDsModelBlockIpObject `tfsdk:"block_ip"`
	Drop        types.Bool                                     `tfsdk:"drop"`
	ResetBoth   types.Bool                                     `tfsdk:"reset_both"`
	ResetClient types.Bool                                     `tfsdk:"reset_client"`
	ResetServer types.Bool                                     `tfsdk:"reset_server"`
}

type antiSpywareSignaturesListDsModelBlockIpObject struct {
	Duration types.Int64  `tfsdk:"duration"`
	TrackBy  types.String `tfsdk:"track_by"`
}

type antiSpywareSignaturesListDsModelSignatureObject struct {
	Combination *antiSpywareSignaturesListDsModelCombinationObject `tfsdk:"combination"`
	Standard    []antiSpywareSignaturesListDsModelStandardObject   `tfsdk:"standard"`
}

type antiSpywareSignaturesListDsModelCombinationObject struct {
	AndCondition  []antiSpywareSignaturesListDsModelAndConditionObject `tfsdk:"and_condition"`
	OrderFree     types.Bool                                           `tfsdk:"order_free"`
	TimeAttribute *antiSpywareSignaturesListDsModelTimeAttributeObject `tfsdk:"time_attribute"`
}

type antiSpywareSignaturesListDsModelAndConditionObject struct {
	Name        types.String                                        `tfsdk:"name"`
	OrCondition []antiSpywareSignaturesListDsModelOrConditionObject `tfsdk:"or_condition"`
}

type antiSpywareSignaturesListDsModelOrConditionObject struct {
	Name     types.String `tfsdk:"name"`
	ThreatId types.String `tfsdk:"threat_id"`
}

type antiSpywareSignaturesListDsModelTimeAttributeObject struct {
	Interval  types.Int64  `tfsdk:"interval"`
	Threshold types.Int64  `tfsdk:"threshold"`
	TrackBy   types.String `tfsdk:"track_by"`
}

type antiSpywareSignaturesListDsModelStandardObject struct {
	AndCondition []antiSpywareSignaturesListDsModelAndConditionObject `tfsdk:"and_condition"`
	Comment      types.String                                         `tfsdk:"comment"`
	Name         types.String                                         `tfsdk:"name"`
	OrderFree    types.Bool                                           `tfsdk:"order_free"`
	Scope        types.String                                         `tfsdk:"scope"`
}

// Metadata returns the data source type name.
func (d *antiSpywareSignaturesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_anti_spyware_signatures_list"
}

// Schema defines the schema for this listing data source.
func (d *antiSpywareSignaturesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"bugtraq": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"comment": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"cve": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"default_action": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"alert": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"allow": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"block_ip": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"duration": dsschema.Int64Attribute{
											Description: "",
											Computed:    true,
										},
										"track_by": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
								"drop": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"reset_both": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"reset_client": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"reset_server": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
							},
						},
						"direction": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"object_id": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"reference": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"severity": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"signature": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"combination": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
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
																"threat_id": dsschema.StringAttribute{
																	Description: "",
																	Computed:    true,
																},
															},
														},
													},
												},
											},
										},
										"order_free": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
										"time_attribute": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"interval": dsschema.Int64Attribute{
													Description: "",
													Computed:    true,
												},
												"threshold": dsschema.Int64Attribute{
													Description: "",
													Computed:    true,
												},
												"track_by": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
											},
										},
									},
								},
								"standard": dsschema.ListNestedAttribute{
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
																	"threat_id": dsschema.StringAttribute{
																		Description: "",
																		Computed:    true,
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
							},
						},
						"threat_id": dsschema.Int64Attribute{
							Description: "",
							Computed:    true,
						},
						"threatname": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"vendor": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
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
func (d *antiSpywareSignaturesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *antiSpywareSignaturesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state antiSpywareSignaturesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_anti_spyware_signatures_list",
		"terraform_provider_function": "Read",
		"limit":                       state.Limit.ValueInt64(),
		"has_limit":                   !state.Limit.IsNull(),
		"offset":                      state.Offset.ValueInt64(),
		"has_offset":                  !state.Offset.IsNull(),
		"name":                        state.Name.ValueString(),
		"has_name":                    !state.Name.IsNull(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := iblCTtp.NewClient(d.client)
	input := iblCTtp.ListInput{
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
	state.Id = types.StringValue(idBuilder.String())
	var var0 []antiSpywareSignaturesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]antiSpywareSignaturesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 antiSpywareSignaturesListDsModelConfig
			var var3 *antiSpywareSignaturesListDsModelDefaultActionObject
			if var1.DefaultAction != nil {
				var3 = &antiSpywareSignaturesListDsModelDefaultActionObject{}
				var var4 *antiSpywareSignaturesListDsModelBlockIpObject
				if var1.DefaultAction.BlockIp != nil {
					var4 = &antiSpywareSignaturesListDsModelBlockIpObject{}
					var4.Duration = types.Int64Value(var1.DefaultAction.BlockIp.Duration)
					var4.TrackBy = types.StringValue(var1.DefaultAction.BlockIp.TrackBy)
				}
				if var1.DefaultAction.Alert != nil {
					var3.Alert = types.BoolValue(true)
				}
				if var1.DefaultAction.Allow != nil {
					var3.Allow = types.BoolValue(true)
				}
				var3.BlockIp = var4
				if var1.DefaultAction.Drop != nil {
					var3.Drop = types.BoolValue(true)
				}
				if var1.DefaultAction.ResetBoth != nil {
					var3.ResetBoth = types.BoolValue(true)
				}
				if var1.DefaultAction.ResetClient != nil {
					var3.ResetClient = types.BoolValue(true)
				}
				if var1.DefaultAction.ResetServer != nil {
					var3.ResetServer = types.BoolValue(true)
				}
			}
			var var5 *antiSpywareSignaturesListDsModelSignatureObject
			if var1.Signature != nil {
				var5 = &antiSpywareSignaturesListDsModelSignatureObject{}
				var var6 *antiSpywareSignaturesListDsModelCombinationObject
				if var1.Signature.Combination != nil {
					var6 = &antiSpywareSignaturesListDsModelCombinationObject{}
					var var7 []antiSpywareSignaturesListDsModelAndConditionObject
					if len(var1.Signature.Combination.AndCondition) != 0 {
						var7 = make([]antiSpywareSignaturesListDsModelAndConditionObject, 0, len(var1.Signature.Combination.AndCondition))
						for var8Index := range var1.Signature.Combination.AndCondition {
							var8 := var1.Signature.Combination.AndCondition[var8Index]
							var var9 antiSpywareSignaturesListDsModelAndConditionObject
							var var10 []antiSpywareSignaturesListDsModelOrConditionObject
							if len(var8.OrCondition) != 0 {
								var10 = make([]antiSpywareSignaturesListDsModelOrConditionObject, 0, len(var8.OrCondition))
								for var11Index := range var8.OrCondition {
									var11 := var8.OrCondition[var11Index]
									var var12 antiSpywareSignaturesListDsModelOrConditionObject
									var12.Name = types.StringValue(var11.Name)
									var12.ThreatId = types.StringValue(var11.ThreatId)
									var10 = append(var10, var12)
								}
							}
							var9.Name = types.StringValue(var8.Name)
							var9.OrCondition = var10
							var7 = append(var7, var9)
						}
					}
					var var13 *antiSpywareSignaturesListDsModelTimeAttributeObject
					if var1.Signature.Combination.TimeAttribute != nil {
						var13 = &antiSpywareSignaturesListDsModelTimeAttributeObject{}
						var13.Interval = types.Int64Value(var1.Signature.Combination.TimeAttribute.Interval)
						var13.Threshold = types.Int64Value(var1.Signature.Combination.TimeAttribute.Threshold)
						var13.TrackBy = types.StringValue(var1.Signature.Combination.TimeAttribute.TrackBy)
					}
					var6.AndCondition = var7
					var6.OrderFree = types.BoolValue(var1.Signature.Combination.OrderFree)
					var6.TimeAttribute = var13
				}
				var var14 []antiSpywareSignaturesListDsModelStandardObject
				if len(var1.Signature.Standard) != 0 {
					var14 = make([]antiSpywareSignaturesListDsModelStandardObject, 0, len(var1.Signature.Standard))
					for var15Index := range var1.Signature.Standard {
						var15 := var1.Signature.Standard[var15Index]
						var var16 antiSpywareSignaturesListDsModelStandardObject
						var var17 []antiSpywareSignaturesListDsModelAndConditionObject
						if len(var15.AndCondition) != 0 {
							var17 = make([]antiSpywareSignaturesListDsModelAndConditionObject, 0, len(var15.AndCondition))
							for var18Index := range var15.AndCondition {
								var18 := var15.AndCondition[var18Index]
								var var19 antiSpywareSignaturesListDsModelAndConditionObject
								var var20 []antiSpywareSignaturesListDsModelOrConditionObject
								if len(var18.OrCondition) != 0 {
									var20 = make([]antiSpywareSignaturesListDsModelOrConditionObject, 0, len(var18.OrCondition))
									for var21Index := range var18.OrCondition {
										var21 := var18.OrCondition[var21Index]
										var var22 antiSpywareSignaturesListDsModelOrConditionObject
										var22.Name = types.StringValue(var21.Name)
										var22.ThreatId = types.StringValue(var21.ThreatId)
										var20 = append(var20, var22)
									}
								}
								var19.Name = types.StringValue(var18.Name)
								var19.OrCondition = var20
								var17 = append(var17, var19)
							}
						}
						var16.AndCondition = var17
						var16.Comment = types.StringValue(var15.Comment)
						var16.Name = types.StringValue(var15.Name)
						var16.OrderFree = types.BoolValue(var15.OrderFree)
						var16.Scope = types.StringValue(var15.Scope)
						var14 = append(var14, var16)
					}
				}
				var5.Combination = var6
				var5.Standard = var14
			}
			var2.Bugtraq = EncodeStringSlice(var1.Bugtraq)
			var2.Comment = types.StringValue(var1.Comment)
			var2.Cve = EncodeStringSlice(var1.Cve)
			var2.DefaultAction = var3
			var2.Direction = types.StringValue(var1.Direction)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Reference = EncodeStringSlice(var1.Reference)
			var2.Severity = types.StringValue(var1.Severity)
			var2.Signature = var5
			var2.ThreatId = types.Int64Value(var1.ThreatId)
			var2.Threatname = types.StringValue(var1.Threatname)
			var2.Vendor = EncodeStringSlice(var1.Vendor)
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
