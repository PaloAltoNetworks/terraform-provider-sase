package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	ngmdzgb "github.com/paloaltonetworks/sase-go/netsec/schema/qos/profiles"
	qCqdYhf "github.com/paloaltonetworks/sase-go/netsec/service/v1/qosprofiles"

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
	_ datasource.DataSource              = &qosProfilesListDataSource{}
	_ datasource.DataSourceWithConfigure = &qosProfilesListDataSource{}
)

func NewQosProfilesListDataSource() datasource.DataSource {
	return &qosProfilesListDataSource{}
}

type qosProfilesListDataSource struct {
	client *sase.Client
}

type qosProfilesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []qosProfilesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type qosProfilesListDsModelConfig struct {
	AggregateBandwidth *qosProfilesListDsModelAggregateBandwidthObject `tfsdk:"aggregate_bandwidth"`
	ClassBandwidthType *qosProfilesListDsModelClassBandwidthTypeObject `tfsdk:"class_bandwidth_type"`
	ObjectId           types.String                                    `tfsdk:"object_id"`
	Name               types.String                                    `tfsdk:"name"`
}

type qosProfilesListDsModelAggregateBandwidthObject struct {
	EgressGuaranteed types.Int64 `tfsdk:"egress_guaranteed"`
	EgressMax        types.Int64 `tfsdk:"egress_max"`
}

type qosProfilesListDsModelClassBandwidthTypeObject struct {
	Mbps       *qosProfilesListDsModelMbpsObject       `tfsdk:"mbps"`
	Percentage *qosProfilesListDsModelPercentageObject `tfsdk:"percentage"`
}

type qosProfilesListDsModelMbpsObject struct {
	Class []qosProfilesListDsModelClassObject `tfsdk:"class"`
}

type qosProfilesListDsModelClassObject struct {
	ClassBandwidth *qosProfilesListDsModelClassBandwidthObject `tfsdk:"class_bandwidth"`
	Name           types.String                                `tfsdk:"name"`
	Priority       types.String                                `tfsdk:"priority"`
}

type qosProfilesListDsModelClassBandwidthObject struct {
	EgressGuaranteed types.Int64 `tfsdk:"egress_guaranteed"`
	EgressMax        types.Int64 `tfsdk:"egress_max"`
}

type qosProfilesListDsModelPercentageObject struct {
	Class []qosProfilesListDsModelClassObject `tfsdk:"class"`
}

// Metadata returns the data source type name.
func (d *qosProfilesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_qos_profiles_list"
}

// Schema defines the schema for this listing data source.
func (d *qosProfilesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"aggregate_bandwidth": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"egress_guaranteed": dsschema.Int64Attribute{
									Description: "",
									Computed:    true,
								},
								"egress_max": dsschema.Int64Attribute{
									Description: "",
									Computed:    true,
								},
							},
						},
						"class_bandwidth_type": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"mbps": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"class": dsschema.ListNestedAttribute{
											Description: "",
											Computed:    true,
											NestedObject: dsschema.NestedAttributeObject{
												Attributes: map[string]dsschema.Attribute{
													"class_bandwidth": dsschema.SingleNestedAttribute{
														Description: "",
														Computed:    true,
														Attributes: map[string]dsschema.Attribute{
															"egress_guaranteed": dsschema.Int64Attribute{
																Description: "",
																Computed:    true,
															},
															"egress_max": dsschema.Int64Attribute{
																Description: "",
																Computed:    true,
															},
														},
													},
													"name": dsschema.StringAttribute{
														Description: "",
														Computed:    true,
													},
													"priority": dsschema.StringAttribute{
														Description: "",
														Computed:    true,
													},
												},
											},
										},
									},
								},
								"percentage": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"class": dsschema.ListNestedAttribute{
											Description: "",
											Computed:    true,
											NestedObject: dsschema.NestedAttributeObject{
												Attributes: map[string]dsschema.Attribute{
													"class_bandwidth": dsschema.SingleNestedAttribute{
														Description: "",
														Computed:    true,
														Attributes: map[string]dsschema.Attribute{
															"egress_guaranteed": dsschema.Int64Attribute{
																Description: "",
																Computed:    true,
															},
															"egress_max": dsschema.Int64Attribute{
																Description: "",
																Computed:    true,
															},
														},
													},
													"name": dsschema.StringAttribute{
														Description: "",
														Computed:    true,
													},
													"priority": dsschema.StringAttribute{
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
						"object_id": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"name": dsschema.StringAttribute{
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
func (d *qosProfilesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *qosProfilesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state qosProfilesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_qos_profiles_list",
		"limit":            state.Limit.ValueInt64(),
		"has_limit":        !state.Limit.IsNull(),
		"offset":           state.Offset.ValueInt64(),
		"has_offset":       !state.Offset.IsNull(),
		"name":             state.Name.ValueString(),
		"has_name":         !state.Name.IsNull(),
		"folder":           state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := qCqdYhf.NewClient(d.client)
	input := qCqdYhf.ListInput{
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
	var var0 []qosProfilesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]qosProfilesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 qosProfilesListDsModelConfig
			var var3 *qosProfilesListDsModelAggregateBandwidthObject
			if var1.AggregateBandwidth != nil {
				var3 = &qosProfilesListDsModelAggregateBandwidthObject{}
				var3.EgressGuaranteed = types.Int64Value(var1.AggregateBandwidth.EgressGuaranteed)
				var3.EgressMax = types.Int64Value(var1.AggregateBandwidth.EgressMax)
			}
			var var4 *qosProfilesListDsModelClassBandwidthTypeObject
			if var1.ClassBandwidthType != nil {
				var4 = &qosProfilesListDsModelClassBandwidthTypeObject{}
				var var5 *qosProfilesListDsModelMbpsObject
				if var1.ClassBandwidthType.Mbps != nil {
					var5 = &qosProfilesListDsModelMbpsObject{}
					var var6 []qosProfilesListDsModelClassObject
					if len(var1.ClassBandwidthType.Mbps.Class) != 0 {
						var6 = make([]qosProfilesListDsModelClassObject, 0, len(var1.ClassBandwidthType.Mbps.Class))
						for var7Index := range var1.ClassBandwidthType.Mbps.Class {
							var7 := var1.ClassBandwidthType.Mbps.Class[var7Index]
							var var8 qosProfilesListDsModelClassObject
							var var9 *qosProfilesListDsModelClassBandwidthObject
							if var7.ClassBandwidth != nil {
								var9 = &qosProfilesListDsModelClassBandwidthObject{}
								var9.EgressGuaranteed = types.Int64Value(var7.ClassBandwidth.EgressGuaranteed)
								var9.EgressMax = types.Int64Value(var7.ClassBandwidth.EgressMax)
							}
							var8.ClassBandwidth = var9
							var8.Name = types.StringValue(var7.Name)
							var8.Priority = types.StringValue(var7.Priority)
							var6 = append(var6, var8)
						}
					}
					var5.Class = var6
				}
				var var10 *qosProfilesListDsModelPercentageObject
				if var1.ClassBandwidthType.Percentage != nil {
					var10 = &qosProfilesListDsModelPercentageObject{}
					var var11 []qosProfilesListDsModelClassObject
					if len(var1.ClassBandwidthType.Percentage.Class) != 0 {
						var11 = make([]qosProfilesListDsModelClassObject, 0, len(var1.ClassBandwidthType.Percentage.Class))
						for var12Index := range var1.ClassBandwidthType.Percentage.Class {
							var12 := var1.ClassBandwidthType.Percentage.Class[var12Index]
							var var13 qosProfilesListDsModelClassObject
							var var14 *qosProfilesListDsModelClassBandwidthObject
							if var12.ClassBandwidth != nil {
								var14 = &qosProfilesListDsModelClassBandwidthObject{}
								var14.EgressGuaranteed = types.Int64Value(var12.ClassBandwidth.EgressGuaranteed)
								var14.EgressMax = types.Int64Value(var12.ClassBandwidth.EgressMax)
							}
							var13.ClassBandwidth = var14
							var13.Name = types.StringValue(var12.Name)
							var13.Priority = types.StringValue(var12.Priority)
							var11 = append(var11, var13)
						}
					}
					var10.Class = var11
				}
				var4.Mbps = var5
				var4.Percentage = var10
			}
			var2.AggregateBandwidth = var3
			var2.ClassBandwidthType = var4
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Name = types.StringValue(var1.Name)
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
	_ datasource.DataSource              = &qosProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &qosProfilesDataSource{}
)

func NewQosProfilesDataSource() datasource.DataSource {
	return &qosProfilesDataSource{}
}

type qosProfilesDataSource struct {
	client *sase.Client
}

type qosProfilesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`

	// Output.
	// Ref: #/components/schemas/qos-profiles
	AggregateBandwidth *qosProfilesDsModelAggregateBandwidthObject `tfsdk:"aggregate_bandwidth"`
	ClassBandwidthType *qosProfilesDsModelClassBandwidthTypeObject `tfsdk:"class_bandwidth_type"`
	// input omit: ObjectId
	Name types.String `tfsdk:"name"`
}

type qosProfilesDsModelAggregateBandwidthObject struct {
	EgressGuaranteed types.Int64 `tfsdk:"egress_guaranteed"`
	EgressMax        types.Int64 `tfsdk:"egress_max"`
}

type qosProfilesDsModelClassBandwidthTypeObject struct {
	Mbps       *qosProfilesDsModelMbpsObject       `tfsdk:"mbps"`
	Percentage *qosProfilesDsModelPercentageObject `tfsdk:"percentage"`
}

type qosProfilesDsModelMbpsObject struct {
	Class []qosProfilesDsModelClassObject `tfsdk:"class"`
}

type qosProfilesDsModelClassObject struct {
	ClassBandwidth *qosProfilesDsModelClassBandwidthObject `tfsdk:"class_bandwidth"`
	Name           types.String                            `tfsdk:"name"`
	Priority       types.String                            `tfsdk:"priority"`
}

type qosProfilesDsModelClassBandwidthObject struct {
	EgressGuaranteed types.Int64 `tfsdk:"egress_guaranteed"`
	EgressMax        types.Int64 `tfsdk:"egress_max"`
}

type qosProfilesDsModelPercentageObject struct {
	Class []qosProfilesDsModelClassObject `tfsdk:"class"`
}

// Metadata returns the data source type name.
func (d *qosProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_qos_profiles"
}

// Schema defines the schema for this listing data source.
func (d *qosProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"aggregate_bandwidth": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"egress_guaranteed": dsschema.Int64Attribute{
						Description: "",
						Computed:    true,
					},
					"egress_max": dsschema.Int64Attribute{
						Description: "",
						Computed:    true,
					},
				},
			},
			"class_bandwidth_type": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"mbps": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"class": dsschema.ListNestedAttribute{
								Description: "",
								Computed:    true,
								NestedObject: dsschema.NestedAttributeObject{
									Attributes: map[string]dsschema.Attribute{
										"class_bandwidth": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"egress_guaranteed": dsschema.Int64Attribute{
													Description: "",
													Computed:    true,
												},
												"egress_max": dsschema.Int64Attribute{
													Description: "",
													Computed:    true,
												},
											},
										},
										"name": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"priority": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
							},
						},
					},
					"percentage": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"class": dsschema.ListNestedAttribute{
								Description: "",
								Computed:    true,
								NestedObject: dsschema.NestedAttributeObject{
									Attributes: map[string]dsschema.Attribute{
										"class_bandwidth": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"egress_guaranteed": dsschema.Int64Attribute{
													Description: "",
													Computed:    true,
												},
												"egress_max": dsschema.Int64Attribute{
													Description: "",
													Computed:    true,
												},
											},
										},
										"name": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"priority": dsschema.StringAttribute{
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
			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
		},
	}
}

// Configure prepares the struct.
func (d *qosProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *qosProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state qosProfilesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_qos_profiles",
		"object_id":        state.ObjectId.ValueString(),
	})

	// Prepare to run the command.
	svc := qCqdYhf.NewClient(d.client)
	input := qCqdYhf.ReadInput{
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
	var var0 *qosProfilesDsModelAggregateBandwidthObject
	if ans.AggregateBandwidth != nil {
		var0 = &qosProfilesDsModelAggregateBandwidthObject{}
		var0.EgressGuaranteed = types.Int64Value(ans.AggregateBandwidth.EgressGuaranteed)
		var0.EgressMax = types.Int64Value(ans.AggregateBandwidth.EgressMax)
	}
	var var1 *qosProfilesDsModelClassBandwidthTypeObject
	if ans.ClassBandwidthType != nil {
		var1 = &qosProfilesDsModelClassBandwidthTypeObject{}
		var var2 *qosProfilesDsModelMbpsObject
		if ans.ClassBandwidthType.Mbps != nil {
			var2 = &qosProfilesDsModelMbpsObject{}
			var var3 []qosProfilesDsModelClassObject
			if len(ans.ClassBandwidthType.Mbps.Class) != 0 {
				var3 = make([]qosProfilesDsModelClassObject, 0, len(ans.ClassBandwidthType.Mbps.Class))
				for var4Index := range ans.ClassBandwidthType.Mbps.Class {
					var4 := ans.ClassBandwidthType.Mbps.Class[var4Index]
					var var5 qosProfilesDsModelClassObject
					var var6 *qosProfilesDsModelClassBandwidthObject
					if var4.ClassBandwidth != nil {
						var6 = &qosProfilesDsModelClassBandwidthObject{}
						var6.EgressGuaranteed = types.Int64Value(var4.ClassBandwidth.EgressGuaranteed)
						var6.EgressMax = types.Int64Value(var4.ClassBandwidth.EgressMax)
					}
					var5.ClassBandwidth = var6
					var5.Name = types.StringValue(var4.Name)
					var5.Priority = types.StringValue(var4.Priority)
					var3 = append(var3, var5)
				}
			}
			var2.Class = var3
		}
		var var7 *qosProfilesDsModelPercentageObject
		if ans.ClassBandwidthType.Percentage != nil {
			var7 = &qosProfilesDsModelPercentageObject{}
			var var8 []qosProfilesDsModelClassObject
			if len(ans.ClassBandwidthType.Percentage.Class) != 0 {
				var8 = make([]qosProfilesDsModelClassObject, 0, len(ans.ClassBandwidthType.Percentage.Class))
				for var9Index := range ans.ClassBandwidthType.Percentage.Class {
					var9 := ans.ClassBandwidthType.Percentage.Class[var9Index]
					var var10 qosProfilesDsModelClassObject
					var var11 *qosProfilesDsModelClassBandwidthObject
					if var9.ClassBandwidth != nil {
						var11 = &qosProfilesDsModelClassBandwidthObject{}
						var11.EgressGuaranteed = types.Int64Value(var9.ClassBandwidth.EgressGuaranteed)
						var11.EgressMax = types.Int64Value(var9.ClassBandwidth.EgressMax)
					}
					var10.ClassBandwidth = var11
					var10.Name = types.StringValue(var9.Name)
					var10.Priority = types.StringValue(var9.Priority)
					var8 = append(var8, var10)
				}
			}
			var7.Class = var8
		}
		var1.Mbps = var2
		var1.Percentage = var7
	}
	state.AggregateBandwidth = var0
	state.ClassBandwidthType = var1
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &qosProfilesResource{}
	_ resource.ResourceWithConfigure   = &qosProfilesResource{}
	_ resource.ResourceWithImportState = &qosProfilesResource{}
)

func NewQosProfilesResource() resource.Resource {
	return &qosProfilesResource{}
}

type qosProfilesResource struct {
	client *sase.Client
}

type qosProfilesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/qos-profiles
	AggregateBandwidth *qosProfilesRsModelAggregateBandwidthObject `tfsdk:"aggregate_bandwidth"`
	ClassBandwidthType *qosProfilesRsModelClassBandwidthTypeObject `tfsdk:"class_bandwidth_type"`
	ObjectId           types.String                                `tfsdk:"object_id"`
	Name               types.String                                `tfsdk:"name"`
}

type qosProfilesRsModelAggregateBandwidthObject struct {
	EgressGuaranteed types.Int64 `tfsdk:"egress_guaranteed"`
	EgressMax        types.Int64 `tfsdk:"egress_max"`
}

type qosProfilesRsModelClassBandwidthTypeObject struct {
	Mbps       *qosProfilesRsModelMbpsObject       `tfsdk:"mbps"`
	Percentage *qosProfilesRsModelPercentageObject `tfsdk:"percentage"`
}

type qosProfilesRsModelMbpsObject struct {
	Class []qosProfilesRsModelClassObject `tfsdk:"class"`
}

type qosProfilesRsModelClassObject struct {
	ClassBandwidth *qosProfilesRsModelClassBandwidthObject `tfsdk:"class_bandwidth"`
	Name           types.String                            `tfsdk:"name"`
	Priority       types.String                            `tfsdk:"priority"`
}

type qosProfilesRsModelClassBandwidthObject struct {
	EgressGuaranteed types.Int64 `tfsdk:"egress_guaranteed"`
	EgressMax        types.Int64 `tfsdk:"egress_max"`
}

type qosProfilesRsModelPercentageObject struct {
	Class []qosProfilesRsModelClassObject `tfsdk:"class"`
}

// Metadata returns the data source type name.
func (r *qosProfilesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_qos_profiles"
}

// Schema defines the schema for this listing data source.
func (r *qosProfilesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

			"aggregate_bandwidth": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"egress_guaranteed": rsschema.Int64Attribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							DefaultInt64(0),
						},
						Validators: []validator.Int64{
							int64validator.Between(0, 16000),
						},
					},
					"egress_max": rsschema.Int64Attribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							DefaultInt64(0),
						},
						Validators: []validator.Int64{
							int64validator.Between(0, 60000),
						},
					},
				},
			},
			"class_bandwidth_type": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"mbps": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"class": rsschema.ListNestedAttribute{
								Description: "",
								Optional:    true,
								NestedObject: rsschema.NestedAttributeObject{
									Attributes: map[string]rsschema.Attribute{
										"class_bandwidth": rsschema.SingleNestedAttribute{
											Description: "",
											Optional:    true,
											Attributes: map[string]rsschema.Attribute{
												"egress_guaranteed": rsschema.Int64Attribute{
													Description: "",
													Optional:    true,
													Computed:    true,
													PlanModifiers: []planmodifier.Int64{
														DefaultInt64(0),
													},
													Validators: []validator.Int64{
														int64validator.Between(0, 60000),
													},
												},
												"egress_max": rsschema.Int64Attribute{
													Description: "",
													Optional:    true,
													Computed:    true,
													PlanModifiers: []planmodifier.Int64{
														DefaultInt64(0),
													},
													Validators: []validator.Int64{
														int64validator.Between(0, 60000),
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
											Validators: []validator.String{
												stringvalidator.LengthAtMost(31),
											},
										},
										"priority": rsschema.StringAttribute{
											Description: "",
											Optional:    true,
											Computed:    true,
											PlanModifiers: []planmodifier.String{
												DefaultString("medium"),
											},
											Validators: []validator.String{
												stringvalidator.OneOf("real-time", "high", "medium", "low"),
											},
										},
									},
								},
							},
						},
					},
					"percentage": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"class": rsschema.ListNestedAttribute{
								Description: "",
								Optional:    true,
								NestedObject: rsschema.NestedAttributeObject{
									Attributes: map[string]rsschema.Attribute{
										"class_bandwidth": rsschema.SingleNestedAttribute{
											Description: "",
											Optional:    true,
											Attributes: map[string]rsschema.Attribute{
												"egress_guaranteed": rsschema.Int64Attribute{
													Description: "",
													Optional:    true,
													Computed:    true,
													PlanModifiers: []planmodifier.Int64{
														DefaultInt64(0),
													},
													Validators: []validator.Int64{
														int64validator.Between(0, 60000),
													},
												},
												"egress_max": rsschema.Int64Attribute{
													Description: "",
													Optional:    true,
													Computed:    true,
													PlanModifiers: []planmodifier.Int64{
														DefaultInt64(0),
													},
													Validators: []validator.Int64{
														int64validator.Between(0, 60000),
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
											Validators: []validator.String{
												stringvalidator.LengthAtMost(31),
											},
										},
										"priority": rsschema.StringAttribute{
											Description: "",
											Optional:    true,
											Computed:    true,
											PlanModifiers: []planmodifier.String{
												DefaultString("medium"),
											},
											Validators: []validator.String{
												stringvalidator.OneOf("real-time", "high", "medium", "low"),
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
		},
	}
}

// Configure prepares the struct.
func (r *qosProfilesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *qosProfilesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state qosProfilesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "sase_qos_profiles",
		"folder":        state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := qCqdYhf.NewClient(r.client)
	input := qCqdYhf.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 ngmdzgb.Config
	var var1 *ngmdzgb.AggregateBandwidthObject
	if state.AggregateBandwidth != nil {
		var1 = &ngmdzgb.AggregateBandwidthObject{}
		var1.EgressGuaranteed = state.AggregateBandwidth.EgressGuaranteed.ValueInt64()
		var1.EgressMax = state.AggregateBandwidth.EgressMax.ValueInt64()
	}
	var0.AggregateBandwidth = var1
	var var2 *ngmdzgb.ClassBandwidthTypeObject
	if state.ClassBandwidthType != nil {
		var2 = &ngmdzgb.ClassBandwidthTypeObject{}
		var var3 *ngmdzgb.MbpsObject
		if state.ClassBandwidthType.Mbps != nil {
			var3 = &ngmdzgb.MbpsObject{}
			var var4 []ngmdzgb.ClassObject
			if len(state.ClassBandwidthType.Mbps.Class) != 0 {
				var4 = make([]ngmdzgb.ClassObject, 0, len(state.ClassBandwidthType.Mbps.Class))
				for var5Index := range state.ClassBandwidthType.Mbps.Class {
					var5 := state.ClassBandwidthType.Mbps.Class[var5Index]
					var var6 ngmdzgb.ClassObject
					var var7 *ngmdzgb.ClassBandwidthObject
					if var5.ClassBandwidth != nil {
						var7 = &ngmdzgb.ClassBandwidthObject{}
						var7.EgressGuaranteed = var5.ClassBandwidth.EgressGuaranteed.ValueInt64()
						var7.EgressMax = var5.ClassBandwidth.EgressMax.ValueInt64()
					}
					var6.ClassBandwidth = var7
					var6.Name = var5.Name.ValueString()
					var6.Priority = var5.Priority.ValueString()
					var4 = append(var4, var6)
				}
			}
			var3.Class = var4
		}
		var2.Mbps = var3
		var var8 *ngmdzgb.PercentageObject
		if state.ClassBandwidthType.Percentage != nil {
			var8 = &ngmdzgb.PercentageObject{}
			var var9 []ngmdzgb.ClassObject
			if len(state.ClassBandwidthType.Percentage.Class) != 0 {
				var9 = make([]ngmdzgb.ClassObject, 0, len(state.ClassBandwidthType.Percentage.Class))
				for var10Index := range state.ClassBandwidthType.Percentage.Class {
					var10 := state.ClassBandwidthType.Percentage.Class[var10Index]
					var var11 ngmdzgb.ClassObject
					var var12 *ngmdzgb.ClassBandwidthObject
					if var10.ClassBandwidth != nil {
						var12 = &ngmdzgb.ClassBandwidthObject{}
						var12.EgressGuaranteed = var10.ClassBandwidth.EgressGuaranteed.ValueInt64()
						var12.EgressMax = var10.ClassBandwidth.EgressMax.ValueInt64()
					}
					var11.ClassBandwidth = var12
					var11.Name = var10.Name.ValueString()
					var11.Priority = var10.Priority.ValueString()
					var9 = append(var9, var11)
				}
			}
			var8.Class = var9
		}
		var2.Percentage = var8
	}
	var0.ClassBandwidthType = var2
	var0.Name = state.Name.ValueString()
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Create(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in create", err.Error())
		return
	}

	// Store the answer to state.
	state.Id = types.StringValue(strings.Join([]string{input.Folder, ans.ObjectId}, IdSeparator))
	var var13 *qosProfilesRsModelAggregateBandwidthObject
	if ans.AggregateBandwidth != nil {
		var13 = &qosProfilesRsModelAggregateBandwidthObject{}
		var13.EgressGuaranteed = types.Int64Value(ans.AggregateBandwidth.EgressGuaranteed)
		var13.EgressMax = types.Int64Value(ans.AggregateBandwidth.EgressMax)
	}
	var var14 *qosProfilesRsModelClassBandwidthTypeObject
	if ans.ClassBandwidthType != nil {
		var14 = &qosProfilesRsModelClassBandwidthTypeObject{}
		var var15 *qosProfilesRsModelMbpsObject
		if ans.ClassBandwidthType.Mbps != nil {
			var15 = &qosProfilesRsModelMbpsObject{}
			var var16 []qosProfilesRsModelClassObject
			if len(ans.ClassBandwidthType.Mbps.Class) != 0 {
				var16 = make([]qosProfilesRsModelClassObject, 0, len(ans.ClassBandwidthType.Mbps.Class))
				for var17Index := range ans.ClassBandwidthType.Mbps.Class {
					var17 := ans.ClassBandwidthType.Mbps.Class[var17Index]
					var var18 qosProfilesRsModelClassObject
					var var19 *qosProfilesRsModelClassBandwidthObject
					if var17.ClassBandwidth != nil {
						var19 = &qosProfilesRsModelClassBandwidthObject{}
						var19.EgressGuaranteed = types.Int64Value(var17.ClassBandwidth.EgressGuaranteed)
						var19.EgressMax = types.Int64Value(var17.ClassBandwidth.EgressMax)
					}
					var18.ClassBandwidth = var19
					var18.Name = types.StringValue(var17.Name)
					var18.Priority = types.StringValue(var17.Priority)
					var16 = append(var16, var18)
				}
			}
			var15.Class = var16
		}
		var var20 *qosProfilesRsModelPercentageObject
		if ans.ClassBandwidthType.Percentage != nil {
			var20 = &qosProfilesRsModelPercentageObject{}
			var var21 []qosProfilesRsModelClassObject
			if len(ans.ClassBandwidthType.Percentage.Class) != 0 {
				var21 = make([]qosProfilesRsModelClassObject, 0, len(ans.ClassBandwidthType.Percentage.Class))
				for var22Index := range ans.ClassBandwidthType.Percentage.Class {
					var22 := ans.ClassBandwidthType.Percentage.Class[var22Index]
					var var23 qosProfilesRsModelClassObject
					var var24 *qosProfilesRsModelClassBandwidthObject
					if var22.ClassBandwidth != nil {
						var24 = &qosProfilesRsModelClassBandwidthObject{}
						var24.EgressGuaranteed = types.Int64Value(var22.ClassBandwidth.EgressGuaranteed)
						var24.EgressMax = types.Int64Value(var22.ClassBandwidth.EgressMax)
					}
					var23.ClassBandwidth = var24
					var23.Name = types.StringValue(var22.Name)
					var23.Priority = types.StringValue(var22.Priority)
					var21 = append(var21, var23)
				}
			}
			var20.Class = var21
		}
		var14.Mbps = var15
		var14.Percentage = var20
	}
	state.AggregateBandwidth = var13
	state.ClassBandwidthType = var14
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *qosProfilesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state qosProfilesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "sase_qos_profiles",
		"locMap":        map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":        tokens,
	})

	// Prepare to read the config.
	svc := qCqdYhf.NewClient(r.client)
	input := qCqdYhf.ReadInput{
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
	var var0 *qosProfilesRsModelAggregateBandwidthObject
	if ans.AggregateBandwidth != nil {
		var0 = &qosProfilesRsModelAggregateBandwidthObject{}
		var0.EgressGuaranteed = types.Int64Value(ans.AggregateBandwidth.EgressGuaranteed)
		var0.EgressMax = types.Int64Value(ans.AggregateBandwidth.EgressMax)
	}
	var var1 *qosProfilesRsModelClassBandwidthTypeObject
	if ans.ClassBandwidthType != nil {
		var1 = &qosProfilesRsModelClassBandwidthTypeObject{}
		var var2 *qosProfilesRsModelMbpsObject
		if ans.ClassBandwidthType.Mbps != nil {
			var2 = &qosProfilesRsModelMbpsObject{}
			var var3 []qosProfilesRsModelClassObject
			if len(ans.ClassBandwidthType.Mbps.Class) != 0 {
				var3 = make([]qosProfilesRsModelClassObject, 0, len(ans.ClassBandwidthType.Mbps.Class))
				for var4Index := range ans.ClassBandwidthType.Mbps.Class {
					var4 := ans.ClassBandwidthType.Mbps.Class[var4Index]
					var var5 qosProfilesRsModelClassObject
					var var6 *qosProfilesRsModelClassBandwidthObject
					if var4.ClassBandwidth != nil {
						var6 = &qosProfilesRsModelClassBandwidthObject{}
						var6.EgressGuaranteed = types.Int64Value(var4.ClassBandwidth.EgressGuaranteed)
						var6.EgressMax = types.Int64Value(var4.ClassBandwidth.EgressMax)
					}
					var5.ClassBandwidth = var6
					var5.Name = types.StringValue(var4.Name)
					var5.Priority = types.StringValue(var4.Priority)
					var3 = append(var3, var5)
				}
			}
			var2.Class = var3
		}
		var var7 *qosProfilesRsModelPercentageObject
		if ans.ClassBandwidthType.Percentage != nil {
			var7 = &qosProfilesRsModelPercentageObject{}
			var var8 []qosProfilesRsModelClassObject
			if len(ans.ClassBandwidthType.Percentage.Class) != 0 {
				var8 = make([]qosProfilesRsModelClassObject, 0, len(ans.ClassBandwidthType.Percentage.Class))
				for var9Index := range ans.ClassBandwidthType.Percentage.Class {
					var9 := ans.ClassBandwidthType.Percentage.Class[var9Index]
					var var10 qosProfilesRsModelClassObject
					var var11 *qosProfilesRsModelClassBandwidthObject
					if var9.ClassBandwidth != nil {
						var11 = &qosProfilesRsModelClassBandwidthObject{}
						var11.EgressGuaranteed = types.Int64Value(var9.ClassBandwidth.EgressGuaranteed)
						var11.EgressMax = types.Int64Value(var9.ClassBandwidth.EgressMax)
					}
					var10.ClassBandwidth = var11
					var10.Name = types.StringValue(var9.Name)
					var10.Priority = types.StringValue(var9.Priority)
					var8 = append(var8, var10)
				}
			}
			var7.Class = var8
		}
		var1.Mbps = var2
		var1.Percentage = var7
	}
	state.AggregateBandwidth = var0
	state.ClassBandwidthType = var1
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *qosProfilesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state qosProfilesRsModel
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
		"resource_name": "sase_qos_profiles",
		"object_id":     state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := qCqdYhf.NewClient(r.client)
	input := qCqdYhf.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 ngmdzgb.Config
	var var1 *ngmdzgb.AggregateBandwidthObject
	if plan.AggregateBandwidth != nil {
		var1 = &ngmdzgb.AggregateBandwidthObject{}
		var1.EgressGuaranteed = plan.AggregateBandwidth.EgressGuaranteed.ValueInt64()
		var1.EgressMax = plan.AggregateBandwidth.EgressMax.ValueInt64()
	}
	var0.AggregateBandwidth = var1
	var var2 *ngmdzgb.ClassBandwidthTypeObject
	if plan.ClassBandwidthType != nil {
		var2 = &ngmdzgb.ClassBandwidthTypeObject{}
		var var3 *ngmdzgb.MbpsObject
		if plan.ClassBandwidthType.Mbps != nil {
			var3 = &ngmdzgb.MbpsObject{}
			var var4 []ngmdzgb.ClassObject
			if len(plan.ClassBandwidthType.Mbps.Class) != 0 {
				var4 = make([]ngmdzgb.ClassObject, 0, len(plan.ClassBandwidthType.Mbps.Class))
				for var5Index := range plan.ClassBandwidthType.Mbps.Class {
					var5 := plan.ClassBandwidthType.Mbps.Class[var5Index]
					var var6 ngmdzgb.ClassObject
					var var7 *ngmdzgb.ClassBandwidthObject
					if var5.ClassBandwidth != nil {
						var7 = &ngmdzgb.ClassBandwidthObject{}
						var7.EgressGuaranteed = var5.ClassBandwidth.EgressGuaranteed.ValueInt64()
						var7.EgressMax = var5.ClassBandwidth.EgressMax.ValueInt64()
					}
					var6.ClassBandwidth = var7
					var6.Name = var5.Name.ValueString()
					var6.Priority = var5.Priority.ValueString()
					var4 = append(var4, var6)
				}
			}
			var3.Class = var4
		}
		var2.Mbps = var3
		var var8 *ngmdzgb.PercentageObject
		if plan.ClassBandwidthType.Percentage != nil {
			var8 = &ngmdzgb.PercentageObject{}
			var var9 []ngmdzgb.ClassObject
			if len(plan.ClassBandwidthType.Percentage.Class) != 0 {
				var9 = make([]ngmdzgb.ClassObject, 0, len(plan.ClassBandwidthType.Percentage.Class))
				for var10Index := range plan.ClassBandwidthType.Percentage.Class {
					var10 := plan.ClassBandwidthType.Percentage.Class[var10Index]
					var var11 ngmdzgb.ClassObject
					var var12 *ngmdzgb.ClassBandwidthObject
					if var10.ClassBandwidth != nil {
						var12 = &ngmdzgb.ClassBandwidthObject{}
						var12.EgressGuaranteed = var10.ClassBandwidth.EgressGuaranteed.ValueInt64()
						var12.EgressMax = var10.ClassBandwidth.EgressMax.ValueInt64()
					}
					var11.ClassBandwidth = var12
					var11.Name = var10.Name.ValueString()
					var11.Priority = var10.Priority.ValueString()
					var9 = append(var9, var11)
				}
			}
			var8.Class = var9
		}
		var2.Percentage = var8
	}
	var0.ClassBandwidthType = var2
	var0.Name = plan.Name.ValueString()
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var13 *qosProfilesRsModelAggregateBandwidthObject
	if ans.AggregateBandwidth != nil {
		var13 = &qosProfilesRsModelAggregateBandwidthObject{}
		var13.EgressGuaranteed = types.Int64Value(ans.AggregateBandwidth.EgressGuaranteed)
		var13.EgressMax = types.Int64Value(ans.AggregateBandwidth.EgressMax)
	}
	var var14 *qosProfilesRsModelClassBandwidthTypeObject
	if ans.ClassBandwidthType != nil {
		var14 = &qosProfilesRsModelClassBandwidthTypeObject{}
		var var15 *qosProfilesRsModelMbpsObject
		if ans.ClassBandwidthType.Mbps != nil {
			var15 = &qosProfilesRsModelMbpsObject{}
			var var16 []qosProfilesRsModelClassObject
			if len(ans.ClassBandwidthType.Mbps.Class) != 0 {
				var16 = make([]qosProfilesRsModelClassObject, 0, len(ans.ClassBandwidthType.Mbps.Class))
				for var17Index := range ans.ClassBandwidthType.Mbps.Class {
					var17 := ans.ClassBandwidthType.Mbps.Class[var17Index]
					var var18 qosProfilesRsModelClassObject
					var var19 *qosProfilesRsModelClassBandwidthObject
					if var17.ClassBandwidth != nil {
						var19 = &qosProfilesRsModelClassBandwidthObject{}
						var19.EgressGuaranteed = types.Int64Value(var17.ClassBandwidth.EgressGuaranteed)
						var19.EgressMax = types.Int64Value(var17.ClassBandwidth.EgressMax)
					}
					var18.ClassBandwidth = var19
					var18.Name = types.StringValue(var17.Name)
					var18.Priority = types.StringValue(var17.Priority)
					var16 = append(var16, var18)
				}
			}
			var15.Class = var16
		}
		var var20 *qosProfilesRsModelPercentageObject
		if ans.ClassBandwidthType.Percentage != nil {
			var20 = &qosProfilesRsModelPercentageObject{}
			var var21 []qosProfilesRsModelClassObject
			if len(ans.ClassBandwidthType.Percentage.Class) != 0 {
				var21 = make([]qosProfilesRsModelClassObject, 0, len(ans.ClassBandwidthType.Percentage.Class))
				for var22Index := range ans.ClassBandwidthType.Percentage.Class {
					var22 := ans.ClassBandwidthType.Percentage.Class[var22Index]
					var var23 qosProfilesRsModelClassObject
					var var24 *qosProfilesRsModelClassBandwidthObject
					if var22.ClassBandwidth != nil {
						var24 = &qosProfilesRsModelClassBandwidthObject{}
						var24.EgressGuaranteed = types.Int64Value(var22.ClassBandwidth.EgressGuaranteed)
						var24.EgressMax = types.Int64Value(var22.ClassBandwidth.EgressMax)
					}
					var23.ClassBandwidth = var24
					var23.Name = types.StringValue(var22.Name)
					var23.Priority = types.StringValue(var22.Priority)
					var21 = append(var21, var23)
				}
			}
			var20.Class = var21
		}
		var14.Mbps = var15
		var14.Percentage = var20
	}
	state.AggregateBandwidth = var13
	state.ClassBandwidthType = var14
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *qosProfilesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name": "sase_qos_profiles",
		"locMap":        map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":        tokens,
	})

	svc := qCqdYhf.NewClient(r.client)
	input := qCqdYhf.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *qosProfilesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
