package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	fZwFwyb "github.com/paloaltonetworks/sase-go/netsec/schema/anti/spyware/profiles"
	iGpoRYz "github.com/paloaltonetworks/sase-go/netsec/service/v1/antispywareprofiles"

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
	_ datasource.DataSource              = &antiSpywareProfilesListDataSource{}
	_ datasource.DataSourceWithConfigure = &antiSpywareProfilesListDataSource{}
)

func NewAntiSpywareProfilesListDataSource() datasource.DataSource {
	return &antiSpywareProfilesListDataSource{}
}

type antiSpywareProfilesListDataSource struct {
	client *sase.Client
}

type antiSpywareProfilesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []antiSpywareProfilesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type antiSpywareProfilesListDsModelConfig struct {
	Description     types.String                                          `tfsdk:"description"`
	ObjectId        types.String                                          `tfsdk:"object_id"`
	Name            types.String                                          `tfsdk:"name"`
	Rules           []antiSpywareProfilesListDsModelRulesObject           `tfsdk:"rules"`
	ThreatException []antiSpywareProfilesListDsModelThreatExceptionObject `tfsdk:"threat_exception"`
}

type antiSpywareProfilesListDsModelRulesObject struct {
	Action        *antiSpywareProfilesListDsModelActionObject `tfsdk:"action"`
	Category      types.String                                `tfsdk:"category"`
	Name          types.String                                `tfsdk:"name"`
	PacketCapture types.String                                `tfsdk:"packet_capture"`
	Severity      []types.String                              `tfsdk:"severity"`
	ThreatName    types.String                                `tfsdk:"threat_name"`
}

type antiSpywareProfilesListDsModelActionObject struct {
	Alert       types.Bool                                   `tfsdk:"alert"`
	Allow       types.Bool                                   `tfsdk:"allow"`
	BlockIp     *antiSpywareProfilesListDsModelBlockIpObject `tfsdk:"block_ip"`
	Drop        types.Bool                                   `tfsdk:"drop"`
	ResetBoth   types.Bool                                   `tfsdk:"reset_both"`
	ResetClient types.Bool                                   `tfsdk:"reset_client"`
	ResetServer types.Bool                                   `tfsdk:"reset_server"`
}

type antiSpywareProfilesListDsModelBlockIpObject struct {
	Duration types.Int64  `tfsdk:"duration"`
	TrackBy  types.String `tfsdk:"track_by"`
}

type antiSpywareProfilesListDsModelThreatExceptionObject struct {
	Action        *antiSpywareProfilesListDsModelActionObject1   `tfsdk:"action"`
	ExemptIp      []antiSpywareProfilesListDsModelExemptIpObject `tfsdk:"exempt_ip"`
	Name          types.String                                   `tfsdk:"name"`
	Notes         types.String                                   `tfsdk:"notes"`
	PacketCapture types.String                                   `tfsdk:"packet_capture"`
}

type antiSpywareProfilesListDsModelActionObject1 struct {
	Alert       types.Bool                                   `tfsdk:"alert"`
	Allow       types.Bool                                   `tfsdk:"allow"`
	BlockIp     *antiSpywareProfilesListDsModelBlockIpObject `tfsdk:"block_ip"`
	Default     types.Bool                                   `tfsdk:"default"`
	Drop        types.Bool                                   `tfsdk:"drop"`
	ResetBoth   types.Bool                                   `tfsdk:"reset_both"`
	ResetClient types.Bool                                   `tfsdk:"reset_client"`
	ResetServer types.Bool                                   `tfsdk:"reset_server"`
}

type antiSpywareProfilesListDsModelExemptIpObject struct {
	Name types.String `tfsdk:"name"`
}

// Metadata returns the data source type name.
func (d *antiSpywareProfilesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_anti_spyware_profiles_list"
}

// Schema defines the schema for this listing data source.
func (d *antiSpywareProfilesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"description": dsschema.StringAttribute{
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
						"rules": dsschema.ListNestedAttribute{
							Description: "",
							Computed:    true,
							NestedObject: dsschema.NestedAttributeObject{
								Attributes: map[string]dsschema.Attribute{
									"action": dsschema.SingleNestedAttribute{
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
									"category": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"name": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"packet_capture": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"severity": dsschema.ListAttribute{
										Description: "",
										Computed:    true,
										ElementType: types.StringType,
									},
									"threat_name": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
						},
						"threat_exception": dsschema.ListNestedAttribute{
							Description: "",
							Computed:    true,
							NestedObject: dsschema.NestedAttributeObject{
								Attributes: map[string]dsschema.Attribute{
									"action": dsschema.SingleNestedAttribute{
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
											"default": dsschema.BoolAttribute{
												Description: "",
												Computed:    true,
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
									"exempt_ip": dsschema.ListNestedAttribute{
										Description: "",
										Computed:    true,
										NestedObject: dsschema.NestedAttributeObject{
											Attributes: map[string]dsschema.Attribute{
												"name": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
											},
										},
									},
									"name": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"notes": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"packet_capture": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
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
func (d *antiSpywareProfilesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *antiSpywareProfilesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state antiSpywareProfilesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_anti_spyware_profiles_list",
		"limit":            state.Limit.ValueInt64(),
		"has_limit":        !state.Limit.IsNull(),
		"offset":           state.Offset.ValueInt64(),
		"has_offset":       !state.Offset.IsNull(),
		"name":             state.Name.ValueString(),
		"has_name":         !state.Name.IsNull(),
		"folder":           state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := iGpoRYz.NewClient(d.client)
	input := iGpoRYz.ListInput{
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
	var var0 []antiSpywareProfilesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]antiSpywareProfilesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 antiSpywareProfilesListDsModelConfig
			var var3 []antiSpywareProfilesListDsModelRulesObject
			if len(var1.Rules) != 0 {
				var3 = make([]antiSpywareProfilesListDsModelRulesObject, 0, len(var1.Rules))
				for var4Index := range var1.Rules {
					var4 := var1.Rules[var4Index]
					var var5 antiSpywareProfilesListDsModelRulesObject
					var var6 *antiSpywareProfilesListDsModelActionObject
					if var4.Action != nil {
						var6 = &antiSpywareProfilesListDsModelActionObject{}
						var var7 *antiSpywareProfilesListDsModelBlockIpObject
						if var4.Action.BlockIp != nil {
							var7 = &antiSpywareProfilesListDsModelBlockIpObject{}
							var7.Duration = types.Int64Value(var4.Action.BlockIp.Duration)
							var7.TrackBy = types.StringValue(var4.Action.BlockIp.TrackBy)
						}
						if var4.Action.Alert != nil {
							var6.Alert = types.BoolValue(true)
						}
						if var4.Action.Allow != nil {
							var6.Allow = types.BoolValue(true)
						}
						var6.BlockIp = var7
						if var4.Action.Drop != nil {
							var6.Drop = types.BoolValue(true)
						}
						if var4.Action.ResetBoth != nil {
							var6.ResetBoth = types.BoolValue(true)
						}
						if var4.Action.ResetClient != nil {
							var6.ResetClient = types.BoolValue(true)
						}
						if var4.Action.ResetServer != nil {
							var6.ResetServer = types.BoolValue(true)
						}
					}
					var5.Action = var6
					var5.Category = types.StringValue(var4.Category)
					var5.Name = types.StringValue(var4.Name)
					var5.PacketCapture = types.StringValue(var4.PacketCapture)
					var5.Severity = EncodeStringSlice(var4.Severity)
					var5.ThreatName = types.StringValue(var4.ThreatName)
					var3 = append(var3, var5)
				}
			}
			var var8 []antiSpywareProfilesListDsModelThreatExceptionObject
			if len(var1.ThreatException) != 0 {
				var8 = make([]antiSpywareProfilesListDsModelThreatExceptionObject, 0, len(var1.ThreatException))
				for var9Index := range var1.ThreatException {
					var9 := var1.ThreatException[var9Index]
					var var10 antiSpywareProfilesListDsModelThreatExceptionObject
					var var11 *antiSpywareProfilesListDsModelActionObject1
					if var9.Action != nil {
						var11 = &antiSpywareProfilesListDsModelActionObject1{}
						var var12 *antiSpywareProfilesListDsModelBlockIpObject
						if var9.Action.BlockIp != nil {
							var12 = &antiSpywareProfilesListDsModelBlockIpObject{}
							var12.Duration = types.Int64Value(var9.Action.BlockIp.Duration)
							var12.TrackBy = types.StringValue(var9.Action.BlockIp.TrackBy)
						}
						if var9.Action.Alert != nil {
							var11.Alert = types.BoolValue(true)
						}
						if var9.Action.Allow != nil {
							var11.Allow = types.BoolValue(true)
						}
						var11.BlockIp = var12
						if var9.Action.Default != nil {
							var11.Default = types.BoolValue(true)
						}
						if var9.Action.Drop != nil {
							var11.Drop = types.BoolValue(true)
						}
						if var9.Action.ResetBoth != nil {
							var11.ResetBoth = types.BoolValue(true)
						}
						if var9.Action.ResetClient != nil {
							var11.ResetClient = types.BoolValue(true)
						}
						if var9.Action.ResetServer != nil {
							var11.ResetServer = types.BoolValue(true)
						}
					}
					var var13 []antiSpywareProfilesListDsModelExemptIpObject
					if len(var9.ExemptIp) != 0 {
						var13 = make([]antiSpywareProfilesListDsModelExemptIpObject, 0, len(var9.ExemptIp))
						for var14Index := range var9.ExemptIp {
							var14 := var9.ExemptIp[var14Index]
							var var15 antiSpywareProfilesListDsModelExemptIpObject
							var15.Name = types.StringValue(var14.Name)
							var13 = append(var13, var15)
						}
					}
					var10.Action = var11
					var10.ExemptIp = var13
					var10.Name = types.StringValue(var9.Name)
					var10.Notes = types.StringValue(var9.Notes)
					var10.PacketCapture = types.StringValue(var9.PacketCapture)
					var8 = append(var8, var10)
				}
			}
			var2.Description = types.StringValue(var1.Description)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Name = types.StringValue(var1.Name)
			var2.Rules = var3
			var2.ThreatException = var8
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
	_ datasource.DataSource              = &antiSpywareProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &antiSpywareProfilesDataSource{}
)

func NewAntiSpywareProfilesDataSource() datasource.DataSource {
	return &antiSpywareProfilesDataSource{}
}

type antiSpywareProfilesDataSource struct {
	client *sase.Client
}

type antiSpywareProfilesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`

	// Output.
	// Ref: #/components/schemas/anti-spyware-profiles
	Description types.String `tfsdk:"description"`
	// input omit: ObjectId
	Name            types.String                                      `tfsdk:"name"`
	Rules           []antiSpywareProfilesDsModelRulesObject           `tfsdk:"rules"`
	ThreatException []antiSpywareProfilesDsModelThreatExceptionObject `tfsdk:"threat_exception"`
}

type antiSpywareProfilesDsModelRulesObject struct {
	Action        *antiSpywareProfilesDsModelActionObject `tfsdk:"action"`
	Category      types.String                            `tfsdk:"category"`
	Name          types.String                            `tfsdk:"name"`
	PacketCapture types.String                            `tfsdk:"packet_capture"`
	Severity      []types.String                          `tfsdk:"severity"`
	ThreatName    types.String                            `tfsdk:"threat_name"`
}

type antiSpywareProfilesDsModelActionObject struct {
	Alert       types.Bool                               `tfsdk:"alert"`
	Allow       types.Bool                               `tfsdk:"allow"`
	BlockIp     *antiSpywareProfilesDsModelBlockIpObject `tfsdk:"block_ip"`
	Drop        types.Bool                               `tfsdk:"drop"`
	ResetBoth   types.Bool                               `tfsdk:"reset_both"`
	ResetClient types.Bool                               `tfsdk:"reset_client"`
	ResetServer types.Bool                               `tfsdk:"reset_server"`
}

type antiSpywareProfilesDsModelBlockIpObject struct {
	Duration types.Int64  `tfsdk:"duration"`
	TrackBy  types.String `tfsdk:"track_by"`
}

type antiSpywareProfilesDsModelThreatExceptionObject struct {
	Action        *antiSpywareProfilesDsModelActionObject1   `tfsdk:"action"`
	ExemptIp      []antiSpywareProfilesDsModelExemptIpObject `tfsdk:"exempt_ip"`
	Name          types.String                               `tfsdk:"name"`
	Notes         types.String                               `tfsdk:"notes"`
	PacketCapture types.String                               `tfsdk:"packet_capture"`
}

type antiSpywareProfilesDsModelActionObject1 struct {
	Alert       types.Bool                               `tfsdk:"alert"`
	Allow       types.Bool                               `tfsdk:"allow"`
	BlockIp     *antiSpywareProfilesDsModelBlockIpObject `tfsdk:"block_ip"`
	Default     types.Bool                               `tfsdk:"default"`
	Drop        types.Bool                               `tfsdk:"drop"`
	ResetBoth   types.Bool                               `tfsdk:"reset_both"`
	ResetClient types.Bool                               `tfsdk:"reset_client"`
	ResetServer types.Bool                               `tfsdk:"reset_server"`
}

type antiSpywareProfilesDsModelExemptIpObject struct {
	Name types.String `tfsdk:"name"`
}

// Metadata returns the data source type name.
func (d *antiSpywareProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_anti_spyware_profiles"
}

// Schema defines the schema for this listing data source.
func (d *antiSpywareProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"description": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"rules": dsschema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"action": dsschema.SingleNestedAttribute{
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
						"category": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"packet_capture": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"severity": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"threat_name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
					},
				},
			},
			"threat_exception": dsschema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"action": dsschema.SingleNestedAttribute{
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
								"default": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
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
						"exempt_ip": dsschema.ListNestedAttribute{
							Description: "",
							Computed:    true,
							NestedObject: dsschema.NestedAttributeObject{
								Attributes: map[string]dsschema.Attribute{
									"name": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
						},
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"notes": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"packet_capture": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure prepares the struct.
func (d *antiSpywareProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *antiSpywareProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state antiSpywareProfilesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_anti_spyware_profiles",
		"object_id":        state.ObjectId.ValueString(),
	})

	// Prepare to run the command.
	svc := iGpoRYz.NewClient(d.client)
	input := iGpoRYz.ReadInput{
		ObjectId: state.ObjectId.ValueString(),
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
	state.Id = types.StringValue(idBuilder.String())
	var var0 []antiSpywareProfilesDsModelRulesObject
	if len(ans.Rules) != 0 {
		var0 = make([]antiSpywareProfilesDsModelRulesObject, 0, len(ans.Rules))
		for var1Index := range ans.Rules {
			var1 := ans.Rules[var1Index]
			var var2 antiSpywareProfilesDsModelRulesObject
			var var3 *antiSpywareProfilesDsModelActionObject
			if var1.Action != nil {
				var3 = &antiSpywareProfilesDsModelActionObject{}
				var var4 *antiSpywareProfilesDsModelBlockIpObject
				if var1.Action.BlockIp != nil {
					var4 = &antiSpywareProfilesDsModelBlockIpObject{}
					var4.Duration = types.Int64Value(var1.Action.BlockIp.Duration)
					var4.TrackBy = types.StringValue(var1.Action.BlockIp.TrackBy)
				}
				if var1.Action.Alert != nil {
					var3.Alert = types.BoolValue(true)
				}
				if var1.Action.Allow != nil {
					var3.Allow = types.BoolValue(true)
				}
				var3.BlockIp = var4
				if var1.Action.Drop != nil {
					var3.Drop = types.BoolValue(true)
				}
				if var1.Action.ResetBoth != nil {
					var3.ResetBoth = types.BoolValue(true)
				}
				if var1.Action.ResetClient != nil {
					var3.ResetClient = types.BoolValue(true)
				}
				if var1.Action.ResetServer != nil {
					var3.ResetServer = types.BoolValue(true)
				}
			}
			var2.Action = var3
			var2.Category = types.StringValue(var1.Category)
			var2.Name = types.StringValue(var1.Name)
			var2.PacketCapture = types.StringValue(var1.PacketCapture)
			var2.Severity = EncodeStringSlice(var1.Severity)
			var2.ThreatName = types.StringValue(var1.ThreatName)
			var0 = append(var0, var2)
		}
	}
	var var5 []antiSpywareProfilesDsModelThreatExceptionObject
	if len(ans.ThreatException) != 0 {
		var5 = make([]antiSpywareProfilesDsModelThreatExceptionObject, 0, len(ans.ThreatException))
		for var6Index := range ans.ThreatException {
			var6 := ans.ThreatException[var6Index]
			var var7 antiSpywareProfilesDsModelThreatExceptionObject
			var var8 *antiSpywareProfilesDsModelActionObject1
			if var6.Action != nil {
				var8 = &antiSpywareProfilesDsModelActionObject1{}
				var var9 *antiSpywareProfilesDsModelBlockIpObject
				if var6.Action.BlockIp != nil {
					var9 = &antiSpywareProfilesDsModelBlockIpObject{}
					var9.Duration = types.Int64Value(var6.Action.BlockIp.Duration)
					var9.TrackBy = types.StringValue(var6.Action.BlockIp.TrackBy)
				}
				if var6.Action.Alert != nil {
					var8.Alert = types.BoolValue(true)
				}
				if var6.Action.Allow != nil {
					var8.Allow = types.BoolValue(true)
				}
				var8.BlockIp = var9
				if var6.Action.Default != nil {
					var8.Default = types.BoolValue(true)
				}
				if var6.Action.Drop != nil {
					var8.Drop = types.BoolValue(true)
				}
				if var6.Action.ResetBoth != nil {
					var8.ResetBoth = types.BoolValue(true)
				}
				if var6.Action.ResetClient != nil {
					var8.ResetClient = types.BoolValue(true)
				}
				if var6.Action.ResetServer != nil {
					var8.ResetServer = types.BoolValue(true)
				}
			}
			var var10 []antiSpywareProfilesDsModelExemptIpObject
			if len(var6.ExemptIp) != 0 {
				var10 = make([]antiSpywareProfilesDsModelExemptIpObject, 0, len(var6.ExemptIp))
				for var11Index := range var6.ExemptIp {
					var11 := var6.ExemptIp[var11Index]
					var var12 antiSpywareProfilesDsModelExemptIpObject
					var12.Name = types.StringValue(var11.Name)
					var10 = append(var10, var12)
				}
			}
			var7.Action = var8
			var7.ExemptIp = var10
			var7.Name = types.StringValue(var6.Name)
			var7.Notes = types.StringValue(var6.Notes)
			var7.PacketCapture = types.StringValue(var6.PacketCapture)
			var5 = append(var5, var7)
		}
	}
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.Rules = var0
	state.ThreatException = var5

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &antiSpywareProfilesResource{}
	_ resource.ResourceWithConfigure   = &antiSpywareProfilesResource{}
	_ resource.ResourceWithImportState = &antiSpywareProfilesResource{}
)

func NewAntiSpywareProfilesResource() resource.Resource {
	return &antiSpywareProfilesResource{}
}

type antiSpywareProfilesResource struct {
	client *sase.Client
}

type antiSpywareProfilesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/anti-spyware-profiles
	Description     types.String                                      `tfsdk:"description"`
	ObjectId        types.String                                      `tfsdk:"object_id"`
	Name            types.String                                      `tfsdk:"name"`
	Rules           []antiSpywareProfilesRsModelRulesObject           `tfsdk:"rules"`
	ThreatException []antiSpywareProfilesRsModelThreatExceptionObject `tfsdk:"threat_exception"`
}

type antiSpywareProfilesRsModelRulesObject struct {
	Action        *antiSpywareProfilesRsModelActionObject `tfsdk:"action"`
	Category      types.String                            `tfsdk:"category"`
	Name          types.String                            `tfsdk:"name"`
	PacketCapture types.String                            `tfsdk:"packet_capture"`
	Severity      []types.String                          `tfsdk:"severity"`
	ThreatName    types.String                            `tfsdk:"threat_name"`
}

type antiSpywareProfilesRsModelActionObject struct {
	Alert       types.Bool                               `tfsdk:"alert"`
	Allow       types.Bool                               `tfsdk:"allow"`
	BlockIp     *antiSpywareProfilesRsModelBlockIpObject `tfsdk:"block_ip"`
	Drop        types.Bool                               `tfsdk:"drop"`
	ResetBoth   types.Bool                               `tfsdk:"reset_both"`
	ResetClient types.Bool                               `tfsdk:"reset_client"`
	ResetServer types.Bool                               `tfsdk:"reset_server"`
}

type antiSpywareProfilesRsModelBlockIpObject struct {
	Duration types.Int64  `tfsdk:"duration"`
	TrackBy  types.String `tfsdk:"track_by"`
}

type antiSpywareProfilesRsModelThreatExceptionObject struct {
	Action        *antiSpywareProfilesRsModelActionObject1   `tfsdk:"action"`
	ExemptIp      []antiSpywareProfilesRsModelExemptIpObject `tfsdk:"exempt_ip"`
	Name          types.String                               `tfsdk:"name"`
	Notes         types.String                               `tfsdk:"notes"`
	PacketCapture types.String                               `tfsdk:"packet_capture"`
}

type antiSpywareProfilesRsModelActionObject1 struct {
	Alert       types.Bool                               `tfsdk:"alert"`
	Allow       types.Bool                               `tfsdk:"allow"`
	BlockIp     *antiSpywareProfilesRsModelBlockIpObject `tfsdk:"block_ip"`
	Default     types.Bool                               `tfsdk:"default"`
	Drop        types.Bool                               `tfsdk:"drop"`
	ResetBoth   types.Bool                               `tfsdk:"reset_both"`
	ResetClient types.Bool                               `tfsdk:"reset_client"`
	ResetServer types.Bool                               `tfsdk:"reset_server"`
}

type antiSpywareProfilesRsModelExemptIpObject struct {
	Name types.String `tfsdk:"name"`
}

// Metadata returns the data source type name.
func (r *antiSpywareProfilesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_anti_spyware_profiles"
}

// Schema defines the schema for this listing data source.
func (r *antiSpywareProfilesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

			"description": rsschema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
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
			},
			"rules": rsschema.ListNestedAttribute{
				Description: "",
				Optional:    true,
				NestedObject: rsschema.NestedAttributeObject{
					Attributes: map[string]rsschema.Attribute{
						"action": rsschema.SingleNestedAttribute{
							Description: "",
							Optional:    true,
							Attributes: map[string]rsschema.Attribute{
								"alert": rsschema.BoolAttribute{
									Description: "",
									Optional:    true,
								},
								"allow": rsschema.BoolAttribute{
									Description: "",
									Optional:    true,
								},
								"block_ip": rsschema.SingleNestedAttribute{
									Description: "",
									Optional:    true,
									Attributes: map[string]rsschema.Attribute{
										"duration": rsschema.Int64Attribute{
											Description: "",
											Optional:    true,
											Computed:    true,
											PlanModifiers: []planmodifier.Int64{
												DefaultInt64(0),
											},
											Validators: []validator.Int64{
												int64validator.Between(1, 3600),
											},
										},
										"track_by": rsschema.StringAttribute{
											Description: "",
											Optional:    true,
											Computed:    true,
											PlanModifiers: []planmodifier.String{
												DefaultString(""),
											},
											Validators: []validator.String{
												stringvalidator.OneOf("source-and-destination", "source"),
											},
										},
									},
								},
								"drop": rsschema.BoolAttribute{
									Description: "",
									Optional:    true,
								},
								"reset_both": rsschema.BoolAttribute{
									Description: "",
									Optional:    true,
								},
								"reset_client": rsschema.BoolAttribute{
									Description: "",
									Optional:    true,
								},
								"reset_server": rsschema.BoolAttribute{
									Description: "",
									Optional:    true,
								},
							},
						},
						"category": rsschema.StringAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
							Validators: []validator.String{
								stringvalidator.OneOf("dns-proxy", "backdoor", "data-theft", "autogen", "spyware", "dns-security", "downloader", "dns-phishing", "phishing-kit", "cryptominer", "hacktool", "dns-benign", "dns-wildfire", "botnet", "dns-grayware", "inline-cloud-c2", "keylogger", "p2p-communication", "domain-edl", "webshell", "command-and-control", "dns-ddns", "net-worm", "any", "tls-fingerprint", "dns-new-domain", "dns", "fraud", "dns-c2", "adware", "post-exploitation", "dns-malware", "browser-hijack", "dns-parked"),
							},
						},
						"name": rsschema.StringAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
						},
						"packet_capture": rsschema.StringAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
							Validators: []validator.String{
								stringvalidator.OneOf("disable", "single-packet", "extended-capture"),
							},
						},
						"severity": rsschema.ListAttribute{
							Description: "",
							Optional:    true,
							ElementType: types.StringType,
						},
						"threat_name": rsschema.StringAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(4),
							},
						},
					},
				},
			},
			"threat_exception": rsschema.ListNestedAttribute{
				Description: "",
				Optional:    true,
				NestedObject: rsschema.NestedAttributeObject{
					Attributes: map[string]rsschema.Attribute{
						"action": rsschema.SingleNestedAttribute{
							Description: "",
							Optional:    true,
							Attributes: map[string]rsschema.Attribute{
								"alert": rsschema.BoolAttribute{
									Description: "",
									Optional:    true,
								},
								"allow": rsschema.BoolAttribute{
									Description: "",
									Optional:    true,
								},
								"block_ip": rsschema.SingleNestedAttribute{
									Description: "",
									Optional:    true,
									Attributes: map[string]rsschema.Attribute{
										"duration": rsschema.Int64Attribute{
											Description: "",
											Optional:    true,
											Computed:    true,
											PlanModifiers: []planmodifier.Int64{
												DefaultInt64(0),
											},
											Validators: []validator.Int64{
												int64validator.Between(1, 3600),
											},
										},
										"track_by": rsschema.StringAttribute{
											Description: "",
											Optional:    true,
											Computed:    true,
											PlanModifiers: []planmodifier.String{
												DefaultString(""),
											},
											Validators: []validator.String{
												stringvalidator.OneOf("source-and-destination", "source"),
											},
										},
									},
								},
								"default": rsschema.BoolAttribute{
									Description: "",
									Optional:    true,
								},
								"drop": rsschema.BoolAttribute{
									Description: "",
									Optional:    true,
								},
								"reset_both": rsschema.BoolAttribute{
									Description: "",
									Optional:    true,
								},
								"reset_client": rsschema.BoolAttribute{
									Description: "",
									Optional:    true,
								},
								"reset_server": rsschema.BoolAttribute{
									Description: "",
									Optional:    true,
								},
							},
						},
						"exempt_ip": rsschema.ListNestedAttribute{
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
									},
								},
							},
						},
						"name": rsschema.StringAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
						},
						"notes": rsschema.StringAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
						},
						"packet_capture": rsschema.StringAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
							Validators: []validator.String{
								stringvalidator.OneOf("disable", "single-packet", "extended-capture"),
							},
						},
					},
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *antiSpywareProfilesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *antiSpywareProfilesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state antiSpywareProfilesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "sase_anti_spyware_profiles",
		"folder":        state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := iGpoRYz.NewClient(r.client)
	input := iGpoRYz.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 fZwFwyb.Config
	var0.Description = state.Description.ValueString()
	var0.Name = state.Name.ValueString()
	var var1 []fZwFwyb.RulesObject
	if len(state.Rules) != 0 {
		var1 = make([]fZwFwyb.RulesObject, 0, len(state.Rules))
		for var2Index := range state.Rules {
			var2 := state.Rules[var2Index]
			var var3 fZwFwyb.RulesObject
			var var4 *fZwFwyb.ActionObject
			if var2.Action != nil {
				var4 = &fZwFwyb.ActionObject{}
				if var2.Action.Alert.ValueBool() {
					var4.Alert = struct{}{}
				}
				if var2.Action.Allow.ValueBool() {
					var4.Allow = struct{}{}
				}
				var var5 *fZwFwyb.BlockIpObject
				if var2.Action.BlockIp != nil {
					var5 = &fZwFwyb.BlockIpObject{}
					var5.Duration = var2.Action.BlockIp.Duration.ValueInt64()
					var5.TrackBy = var2.Action.BlockIp.TrackBy.ValueString()
				}
				var4.BlockIp = var5
				if var2.Action.Drop.ValueBool() {
					var4.Drop = struct{}{}
				}
				if var2.Action.ResetBoth.ValueBool() {
					var4.ResetBoth = struct{}{}
				}
				if var2.Action.ResetClient.ValueBool() {
					var4.ResetClient = struct{}{}
				}
				if var2.Action.ResetServer.ValueBool() {
					var4.ResetServer = struct{}{}
				}
			}
			var3.Action = var4
			var3.Category = var2.Category.ValueString()
			var3.Name = var2.Name.ValueString()
			var3.PacketCapture = var2.PacketCapture.ValueString()
			var3.Severity = DecodeStringSlice(var2.Severity)
			var3.ThreatName = var2.ThreatName.ValueString()
			var1 = append(var1, var3)
		}
	}
	var0.Rules = var1
	var var6 []fZwFwyb.ThreatExceptionObject
	if len(state.ThreatException) != 0 {
		var6 = make([]fZwFwyb.ThreatExceptionObject, 0, len(state.ThreatException))
		for var7Index := range state.ThreatException {
			var7 := state.ThreatException[var7Index]
			var var8 fZwFwyb.ThreatExceptionObject
			var var9 *fZwFwyb.ActionObject1
			if var7.Action != nil {
				var9 = &fZwFwyb.ActionObject1{}
				if var7.Action.Alert.ValueBool() {
					var9.Alert = struct{}{}
				}
				if var7.Action.Allow.ValueBool() {
					var9.Allow = struct{}{}
				}
				var var10 *fZwFwyb.BlockIpObject
				if var7.Action.BlockIp != nil {
					var10 = &fZwFwyb.BlockIpObject{}
					var10.Duration = var7.Action.BlockIp.Duration.ValueInt64()
					var10.TrackBy = var7.Action.BlockIp.TrackBy.ValueString()
				}
				var9.BlockIp = var10
				if var7.Action.Default.ValueBool() {
					var9.Default = struct{}{}
				}
				if var7.Action.Drop.ValueBool() {
					var9.Drop = struct{}{}
				}
				if var7.Action.ResetBoth.ValueBool() {
					var9.ResetBoth = struct{}{}
				}
				if var7.Action.ResetClient.ValueBool() {
					var9.ResetClient = struct{}{}
				}
				if var7.Action.ResetServer.ValueBool() {
					var9.ResetServer = struct{}{}
				}
			}
			var8.Action = var9
			var var11 []fZwFwyb.ExemptIpObject
			if len(var7.ExemptIp) != 0 {
				var11 = make([]fZwFwyb.ExemptIpObject, 0, len(var7.ExemptIp))
				for var12Index := range var7.ExemptIp {
					var12 := var7.ExemptIp[var12Index]
					var var13 fZwFwyb.ExemptIpObject
					var13.Name = var12.Name.ValueString()
					var11 = append(var11, var13)
				}
			}
			var8.ExemptIp = var11
			var8.Name = var7.Name.ValueString()
			var8.Notes = var7.Notes.ValueString()
			var8.PacketCapture = var7.PacketCapture.ValueString()
			var6 = append(var6, var8)
		}
	}
	var0.ThreatException = var6
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
	idBuilder.WriteString(ans.ObjectId)
	state.Id = types.StringValue(idBuilder.String())
	var var14 []antiSpywareProfilesRsModelRulesObject
	if len(ans.Rules) != 0 {
		var14 = make([]antiSpywareProfilesRsModelRulesObject, 0, len(ans.Rules))
		for var15Index := range ans.Rules {
			var15 := ans.Rules[var15Index]
			var var16 antiSpywareProfilesRsModelRulesObject
			var var17 *antiSpywareProfilesRsModelActionObject
			if var15.Action != nil {
				var17 = &antiSpywareProfilesRsModelActionObject{}
				var var18 *antiSpywareProfilesRsModelBlockIpObject
				if var15.Action.BlockIp != nil {
					var18 = &antiSpywareProfilesRsModelBlockIpObject{}
					var18.Duration = types.Int64Value(var15.Action.BlockIp.Duration)
					var18.TrackBy = types.StringValue(var15.Action.BlockIp.TrackBy)
				}
				if var15.Action.Alert != nil {
					var17.Alert = types.BoolValue(true)
				}
				if var15.Action.Allow != nil {
					var17.Allow = types.BoolValue(true)
				}
				var17.BlockIp = var18
				if var15.Action.Drop != nil {
					var17.Drop = types.BoolValue(true)
				}
				if var15.Action.ResetBoth != nil {
					var17.ResetBoth = types.BoolValue(true)
				}
				if var15.Action.ResetClient != nil {
					var17.ResetClient = types.BoolValue(true)
				}
				if var15.Action.ResetServer != nil {
					var17.ResetServer = types.BoolValue(true)
				}
			}
			var16.Action = var17
			var16.Category = types.StringValue(var15.Category)
			var16.Name = types.StringValue(var15.Name)
			var16.PacketCapture = types.StringValue(var15.PacketCapture)
			var16.Severity = EncodeStringSlice(var15.Severity)
			var16.ThreatName = types.StringValue(var15.ThreatName)
			var14 = append(var14, var16)
		}
	}
	var var19 []antiSpywareProfilesRsModelThreatExceptionObject
	if len(ans.ThreatException) != 0 {
		var19 = make([]antiSpywareProfilesRsModelThreatExceptionObject, 0, len(ans.ThreatException))
		for var20Index := range ans.ThreatException {
			var20 := ans.ThreatException[var20Index]
			var var21 antiSpywareProfilesRsModelThreatExceptionObject
			var var22 *antiSpywareProfilesRsModelActionObject1
			if var20.Action != nil {
				var22 = &antiSpywareProfilesRsModelActionObject1{}
				var var23 *antiSpywareProfilesRsModelBlockIpObject
				if var20.Action.BlockIp != nil {
					var23 = &antiSpywareProfilesRsModelBlockIpObject{}
					var23.Duration = types.Int64Value(var20.Action.BlockIp.Duration)
					var23.TrackBy = types.StringValue(var20.Action.BlockIp.TrackBy)
				}
				if var20.Action.Alert != nil {
					var22.Alert = types.BoolValue(true)
				}
				if var20.Action.Allow != nil {
					var22.Allow = types.BoolValue(true)
				}
				var22.BlockIp = var23
				if var20.Action.Default != nil {
					var22.Default = types.BoolValue(true)
				}
				if var20.Action.Drop != nil {
					var22.Drop = types.BoolValue(true)
				}
				if var20.Action.ResetBoth != nil {
					var22.ResetBoth = types.BoolValue(true)
				}
				if var20.Action.ResetClient != nil {
					var22.ResetClient = types.BoolValue(true)
				}
				if var20.Action.ResetServer != nil {
					var22.ResetServer = types.BoolValue(true)
				}
			}
			var var24 []antiSpywareProfilesRsModelExemptIpObject
			if len(var20.ExemptIp) != 0 {
				var24 = make([]antiSpywareProfilesRsModelExemptIpObject, 0, len(var20.ExemptIp))
				for var25Index := range var20.ExemptIp {
					var25 := var20.ExemptIp[var25Index]
					var var26 antiSpywareProfilesRsModelExemptIpObject
					var26.Name = types.StringValue(var25.Name)
					var24 = append(var24, var26)
				}
			}
			var21.Action = var22
			var21.ExemptIp = var24
			var21.Name = types.StringValue(var20.Name)
			var21.Notes = types.StringValue(var20.Notes)
			var21.PacketCapture = types.StringValue(var20.PacketCapture)
			var19 = append(var19, var21)
		}
	}
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.Rules = var14
	state.ThreatException = var19

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *antiSpywareProfilesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state antiSpywareProfilesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "sase_anti_spyware_profiles",
		"locMap":        map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":        tokens,
	})

	// Prepare to read the config.
	svc := iGpoRYz.NewClient(r.client)
	input := iGpoRYz.ReadInput{
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
	var var0 []antiSpywareProfilesRsModelRulesObject
	if len(ans.Rules) != 0 {
		var0 = make([]antiSpywareProfilesRsModelRulesObject, 0, len(ans.Rules))
		for var1Index := range ans.Rules {
			var1 := ans.Rules[var1Index]
			var var2 antiSpywareProfilesRsModelRulesObject
			var var3 *antiSpywareProfilesRsModelActionObject
			if var1.Action != nil {
				var3 = &antiSpywareProfilesRsModelActionObject{}
				var var4 *antiSpywareProfilesRsModelBlockIpObject
				if var1.Action.BlockIp != nil {
					var4 = &antiSpywareProfilesRsModelBlockIpObject{}
					var4.Duration = types.Int64Value(var1.Action.BlockIp.Duration)
					var4.TrackBy = types.StringValue(var1.Action.BlockIp.TrackBy)
				}
				if var1.Action.Alert != nil {
					var3.Alert = types.BoolValue(true)
				}
				if var1.Action.Allow != nil {
					var3.Allow = types.BoolValue(true)
				}
				var3.BlockIp = var4
				if var1.Action.Drop != nil {
					var3.Drop = types.BoolValue(true)
				}
				if var1.Action.ResetBoth != nil {
					var3.ResetBoth = types.BoolValue(true)
				}
				if var1.Action.ResetClient != nil {
					var3.ResetClient = types.BoolValue(true)
				}
				if var1.Action.ResetServer != nil {
					var3.ResetServer = types.BoolValue(true)
				}
			}
			var2.Action = var3
			var2.Category = types.StringValue(var1.Category)
			var2.Name = types.StringValue(var1.Name)
			var2.PacketCapture = types.StringValue(var1.PacketCapture)
			var2.Severity = EncodeStringSlice(var1.Severity)
			var2.ThreatName = types.StringValue(var1.ThreatName)
			var0 = append(var0, var2)
		}
	}
	var var5 []antiSpywareProfilesRsModelThreatExceptionObject
	if len(ans.ThreatException) != 0 {
		var5 = make([]antiSpywareProfilesRsModelThreatExceptionObject, 0, len(ans.ThreatException))
		for var6Index := range ans.ThreatException {
			var6 := ans.ThreatException[var6Index]
			var var7 antiSpywareProfilesRsModelThreatExceptionObject
			var var8 *antiSpywareProfilesRsModelActionObject1
			if var6.Action != nil {
				var8 = &antiSpywareProfilesRsModelActionObject1{}
				var var9 *antiSpywareProfilesRsModelBlockIpObject
				if var6.Action.BlockIp != nil {
					var9 = &antiSpywareProfilesRsModelBlockIpObject{}
					var9.Duration = types.Int64Value(var6.Action.BlockIp.Duration)
					var9.TrackBy = types.StringValue(var6.Action.BlockIp.TrackBy)
				}
				if var6.Action.Alert != nil {
					var8.Alert = types.BoolValue(true)
				}
				if var6.Action.Allow != nil {
					var8.Allow = types.BoolValue(true)
				}
				var8.BlockIp = var9
				if var6.Action.Default != nil {
					var8.Default = types.BoolValue(true)
				}
				if var6.Action.Drop != nil {
					var8.Drop = types.BoolValue(true)
				}
				if var6.Action.ResetBoth != nil {
					var8.ResetBoth = types.BoolValue(true)
				}
				if var6.Action.ResetClient != nil {
					var8.ResetClient = types.BoolValue(true)
				}
				if var6.Action.ResetServer != nil {
					var8.ResetServer = types.BoolValue(true)
				}
			}
			var var10 []antiSpywareProfilesRsModelExemptIpObject
			if len(var6.ExemptIp) != 0 {
				var10 = make([]antiSpywareProfilesRsModelExemptIpObject, 0, len(var6.ExemptIp))
				for var11Index := range var6.ExemptIp {
					var11 := var6.ExemptIp[var11Index]
					var var12 antiSpywareProfilesRsModelExemptIpObject
					var12.Name = types.StringValue(var11.Name)
					var10 = append(var10, var12)
				}
			}
			var7.Action = var8
			var7.ExemptIp = var10
			var7.Name = types.StringValue(var6.Name)
			var7.Notes = types.StringValue(var6.Notes)
			var7.PacketCapture = types.StringValue(var6.PacketCapture)
			var5 = append(var5, var7)
		}
	}
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.Rules = var0
	state.ThreatException = var5

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *antiSpywareProfilesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state antiSpywareProfilesRsModel
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
		"resource_name": "sase_anti_spyware_profiles",
		"object_id":     state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := iGpoRYz.NewClient(r.client)
	input := iGpoRYz.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 fZwFwyb.Config
	var0.Description = plan.Description.ValueString()
	var0.Name = plan.Name.ValueString()
	var var1 []fZwFwyb.RulesObject
	if len(plan.Rules) != 0 {
		var1 = make([]fZwFwyb.RulesObject, 0, len(plan.Rules))
		for var2Index := range plan.Rules {
			var2 := plan.Rules[var2Index]
			var var3 fZwFwyb.RulesObject
			var var4 *fZwFwyb.ActionObject
			if var2.Action != nil {
				var4 = &fZwFwyb.ActionObject{}
				if var2.Action.Alert.ValueBool() {
					var4.Alert = struct{}{}
				}
				if var2.Action.Allow.ValueBool() {
					var4.Allow = struct{}{}
				}
				var var5 *fZwFwyb.BlockIpObject
				if var2.Action.BlockIp != nil {
					var5 = &fZwFwyb.BlockIpObject{}
					var5.Duration = var2.Action.BlockIp.Duration.ValueInt64()
					var5.TrackBy = var2.Action.BlockIp.TrackBy.ValueString()
				}
				var4.BlockIp = var5
				if var2.Action.Drop.ValueBool() {
					var4.Drop = struct{}{}
				}
				if var2.Action.ResetBoth.ValueBool() {
					var4.ResetBoth = struct{}{}
				}
				if var2.Action.ResetClient.ValueBool() {
					var4.ResetClient = struct{}{}
				}
				if var2.Action.ResetServer.ValueBool() {
					var4.ResetServer = struct{}{}
				}
			}
			var3.Action = var4
			var3.Category = var2.Category.ValueString()
			var3.Name = var2.Name.ValueString()
			var3.PacketCapture = var2.PacketCapture.ValueString()
			var3.Severity = DecodeStringSlice(var2.Severity)
			var3.ThreatName = var2.ThreatName.ValueString()
			var1 = append(var1, var3)
		}
	}
	var0.Rules = var1
	var var6 []fZwFwyb.ThreatExceptionObject
	if len(plan.ThreatException) != 0 {
		var6 = make([]fZwFwyb.ThreatExceptionObject, 0, len(plan.ThreatException))
		for var7Index := range plan.ThreatException {
			var7 := plan.ThreatException[var7Index]
			var var8 fZwFwyb.ThreatExceptionObject
			var var9 *fZwFwyb.ActionObject1
			if var7.Action != nil {
				var9 = &fZwFwyb.ActionObject1{}
				if var7.Action.Alert.ValueBool() {
					var9.Alert = struct{}{}
				}
				if var7.Action.Allow.ValueBool() {
					var9.Allow = struct{}{}
				}
				var var10 *fZwFwyb.BlockIpObject
				if var7.Action.BlockIp != nil {
					var10 = &fZwFwyb.BlockIpObject{}
					var10.Duration = var7.Action.BlockIp.Duration.ValueInt64()
					var10.TrackBy = var7.Action.BlockIp.TrackBy.ValueString()
				}
				var9.BlockIp = var10
				if var7.Action.Default.ValueBool() {
					var9.Default = struct{}{}
				}
				if var7.Action.Drop.ValueBool() {
					var9.Drop = struct{}{}
				}
				if var7.Action.ResetBoth.ValueBool() {
					var9.ResetBoth = struct{}{}
				}
				if var7.Action.ResetClient.ValueBool() {
					var9.ResetClient = struct{}{}
				}
				if var7.Action.ResetServer.ValueBool() {
					var9.ResetServer = struct{}{}
				}
			}
			var8.Action = var9
			var var11 []fZwFwyb.ExemptIpObject
			if len(var7.ExemptIp) != 0 {
				var11 = make([]fZwFwyb.ExemptIpObject, 0, len(var7.ExemptIp))
				for var12Index := range var7.ExemptIp {
					var12 := var7.ExemptIp[var12Index]
					var var13 fZwFwyb.ExemptIpObject
					var13.Name = var12.Name.ValueString()
					var11 = append(var11, var13)
				}
			}
			var8.ExemptIp = var11
			var8.Name = var7.Name.ValueString()
			var8.Notes = var7.Notes.ValueString()
			var8.PacketCapture = var7.PacketCapture.ValueString()
			var6 = append(var6, var8)
		}
	}
	var0.ThreatException = var6
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var14 []antiSpywareProfilesRsModelRulesObject
	if len(ans.Rules) != 0 {
		var14 = make([]antiSpywareProfilesRsModelRulesObject, 0, len(ans.Rules))
		for var15Index := range ans.Rules {
			var15 := ans.Rules[var15Index]
			var var16 antiSpywareProfilesRsModelRulesObject
			var var17 *antiSpywareProfilesRsModelActionObject
			if var15.Action != nil {
				var17 = &antiSpywareProfilesRsModelActionObject{}
				var var18 *antiSpywareProfilesRsModelBlockIpObject
				if var15.Action.BlockIp != nil {
					var18 = &antiSpywareProfilesRsModelBlockIpObject{}
					var18.Duration = types.Int64Value(var15.Action.BlockIp.Duration)
					var18.TrackBy = types.StringValue(var15.Action.BlockIp.TrackBy)
				}
				if var15.Action.Alert != nil {
					var17.Alert = types.BoolValue(true)
				}
				if var15.Action.Allow != nil {
					var17.Allow = types.BoolValue(true)
				}
				var17.BlockIp = var18
				if var15.Action.Drop != nil {
					var17.Drop = types.BoolValue(true)
				}
				if var15.Action.ResetBoth != nil {
					var17.ResetBoth = types.BoolValue(true)
				}
				if var15.Action.ResetClient != nil {
					var17.ResetClient = types.BoolValue(true)
				}
				if var15.Action.ResetServer != nil {
					var17.ResetServer = types.BoolValue(true)
				}
			}
			var16.Action = var17
			var16.Category = types.StringValue(var15.Category)
			var16.Name = types.StringValue(var15.Name)
			var16.PacketCapture = types.StringValue(var15.PacketCapture)
			var16.Severity = EncodeStringSlice(var15.Severity)
			var16.ThreatName = types.StringValue(var15.ThreatName)
			var14 = append(var14, var16)
		}
	}
	var var19 []antiSpywareProfilesRsModelThreatExceptionObject
	if len(ans.ThreatException) != 0 {
		var19 = make([]antiSpywareProfilesRsModelThreatExceptionObject, 0, len(ans.ThreatException))
		for var20Index := range ans.ThreatException {
			var20 := ans.ThreatException[var20Index]
			var var21 antiSpywareProfilesRsModelThreatExceptionObject
			var var22 *antiSpywareProfilesRsModelActionObject1
			if var20.Action != nil {
				var22 = &antiSpywareProfilesRsModelActionObject1{}
				var var23 *antiSpywareProfilesRsModelBlockIpObject
				if var20.Action.BlockIp != nil {
					var23 = &antiSpywareProfilesRsModelBlockIpObject{}
					var23.Duration = types.Int64Value(var20.Action.BlockIp.Duration)
					var23.TrackBy = types.StringValue(var20.Action.BlockIp.TrackBy)
				}
				if var20.Action.Alert != nil {
					var22.Alert = types.BoolValue(true)
				}
				if var20.Action.Allow != nil {
					var22.Allow = types.BoolValue(true)
				}
				var22.BlockIp = var23
				if var20.Action.Default != nil {
					var22.Default = types.BoolValue(true)
				}
				if var20.Action.Drop != nil {
					var22.Drop = types.BoolValue(true)
				}
				if var20.Action.ResetBoth != nil {
					var22.ResetBoth = types.BoolValue(true)
				}
				if var20.Action.ResetClient != nil {
					var22.ResetClient = types.BoolValue(true)
				}
				if var20.Action.ResetServer != nil {
					var22.ResetServer = types.BoolValue(true)
				}
			}
			var var24 []antiSpywareProfilesRsModelExemptIpObject
			if len(var20.ExemptIp) != 0 {
				var24 = make([]antiSpywareProfilesRsModelExemptIpObject, 0, len(var20.ExemptIp))
				for var25Index := range var20.ExemptIp {
					var25 := var20.ExemptIp[var25Index]
					var var26 antiSpywareProfilesRsModelExemptIpObject
					var26.Name = types.StringValue(var25.Name)
					var24 = append(var24, var26)
				}
			}
			var21.Action = var22
			var21.ExemptIp = var24
			var21.Name = types.StringValue(var20.Name)
			var21.Notes = types.StringValue(var20.Notes)
			var21.PacketCapture = types.StringValue(var20.PacketCapture)
			var19 = append(var19, var21)
		}
	}
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.Rules = var14
	state.ThreatException = var19

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *antiSpywareProfilesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name": "sase_anti_spyware_profiles",
		"locMap":        map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":        tokens,
	})

	svc := iGpoRYz.NewClient(r.client)
	input := iGpoRYz.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *antiSpywareProfilesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
