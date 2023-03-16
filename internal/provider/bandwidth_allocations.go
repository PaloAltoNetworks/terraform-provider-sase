package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	snSEbPJ "github.com/paloaltonetworks/sase-go/netsec/service/v1/bandwidthallocations"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Data source (listing).
var (
	_ datasource.DataSource              = &bandwidthAllocationsListDataSource{}
	_ datasource.DataSourceWithConfigure = &bandwidthAllocationsListDataSource{}
)

func NewBandwidthAllocationsListDataSource() datasource.DataSource {
	return &bandwidthAllocationsListDataSource{}
}

type bandwidthAllocationsListDataSource struct {
	client *sase.Client
}

type bandwidthAllocationsListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64 `tfsdk:"limit"`
	Offset types.Int64 `tfsdk:"offset"`

	// Output.
	Data []bandwidthAllocationsListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type bandwidthAllocationsListDsModelConfig struct {
	AllocatedBandwidth types.Int64                               `tfsdk:"allocated_bandwidth"`
	ObjectId           types.String                              `tfsdk:"object_id"`
	Name               types.String                              `tfsdk:"name"`
	Qos                *bandwidthAllocationsListDsModelQosObject `tfsdk:"qos"`
	SpnNameList        []types.String                            `tfsdk:"spn_name_list"`
}

type bandwidthAllocationsListDsModelQosObject struct {
	Customized      types.Bool   `tfsdk:"customized"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	GuaranteedRatio types.Int64  `tfsdk:"guaranteed_ratio"`
	Profile         types.String `tfsdk:"profile"`
}

// Metadata returns the data source type name.
func (d *bandwidthAllocationsListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bandwidth_allocations_list"
}

// Schema defines the schema for this listing data source.
func (d *bandwidthAllocationsListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

			// Output.
			"data": dsschema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"allocated_bandwidth": dsschema.Int64Attribute{
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
						"qos": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"customized": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"enabled": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"guaranteed_ratio": dsschema.Int64Attribute{
									Description: "",
									Computed:    true,
								},
								"profile": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
							},
						},
						"spn_name_list": dsschema.ListAttribute{
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
func (d *bandwidthAllocationsListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *bandwidthAllocationsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state bandwidthAllocationsListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_bandwidth_allocations_list",
		"limit":            state.Limit.ValueInt64(),
		"has_limit":        !state.Limit.IsNull(),
		"offset":           state.Offset.ValueInt64(),
		"has_offset":       !state.Offset.IsNull(),
	})

	// Prepare to run the command.
	svc := snSEbPJ.NewClient(d.client)
	input := snSEbPJ.ListInput{}
	if !state.Limit.IsNull() {
		input.Limit = api.Int(state.Limit.ValueInt64())
	}
	if !state.Offset.IsNull() {
		input.Offset = api.Int(state.Offset.ValueInt64())
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
	state.Id = types.StringValue(idBuilder.String())
	var var0 []bandwidthAllocationsListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]bandwidthAllocationsListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 bandwidthAllocationsListDsModelConfig
			var var3 *bandwidthAllocationsListDsModelQosObject
			if var1.Qos != nil {
				var3 = &bandwidthAllocationsListDsModelQosObject{}
				var3.Customized = types.BoolValue(var1.Qos.Customized)
				var3.Enabled = types.BoolValue(var1.Qos.Enabled)
				var3.GuaranteedRatio = types.Int64Value(var1.Qos.GuaranteedRatio)
				var3.Profile = types.StringValue(var1.Qos.Profile)
			}
			var2.AllocatedBandwidth = types.Int64Value(var1.AllocatedBandwidth)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Name = types.StringValue(var1.Name)
			var2.Qos = var3
			var2.SpnNameList = EncodeStringSlice(var1.SpnNameList)
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
