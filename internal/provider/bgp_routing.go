package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	fhcUKOQ "github.com/paloaltonetworks/sase-go/netsec/service/v1/bgprouting"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Data source (listing).
var (
	_ datasource.DataSource              = &bgpRoutingListDataSource{}
	_ datasource.DataSourceWithConfigure = &bgpRoutingListDataSource{}
)

func NewBgpRoutingListDataSource() datasource.DataSource {
	return &bgpRoutingListDataSource{}
}

type bgpRoutingListDataSource struct {
	client *sase.Client
}

type bgpRoutingListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []bgpRoutingListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type bgpRoutingListDsModelConfig struct {
	AcceptRouteOverSC         types.Bool                                    `tfsdk:"accept_route_over_s_c"`
	AddHostRouteToIkePeer     types.Bool                                    `tfsdk:"add_host_route_to_ike_peer"`
	BackboneRouting           types.String                                  `tfsdk:"backbone_routing"`
	OutboundRoutesForServices []types.String                                `tfsdk:"outbound_routes_for_services"`
	RoutingPreference         *bgpRoutingListDsModelRoutingPreferenceObject `tfsdk:"routing_preference"`
	WithdrawStaticRoute       types.Bool                                    `tfsdk:"withdraw_static_route"`
}

type bgpRoutingListDsModelRoutingPreferenceObject struct {
	Default          types.Bool `tfsdk:"default"`
	HotPotatoRouting types.Bool `tfsdk:"hot_potato_routing"`
}

// Metadata returns the data source type name.
func (d *bgpRoutingListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bgp_routing_list"
}

// Schema defines the schema for this listing data source.
func (d *bgpRoutingListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"accept_route_over_s_c": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"add_host_route_to_ike_peer": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"backbone_routing": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"outbound_routes_for_services": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"routing_preference": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"default": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"hot_potato_routing": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
							},
						},
						"withdraw_static_route": dsschema.BoolAttribute{
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
func (d *bgpRoutingListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *bgpRoutingListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state bgpRoutingListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_bgp_routing_list",
		"limit":            state.Limit.ValueInt64(),
		"has_limit":        !state.Limit.IsNull(),
		"offset":           state.Offset.ValueInt64(),
		"has_offset":       !state.Offset.IsNull(),
		"folder":           state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := fhcUKOQ.NewClient(d.client)
	input := fhcUKOQ.ListInput{
		Folder: state.Folder.ValueString(),
	}
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
	idBuilder.WriteString(IdSeparator)
	idBuilder.WriteString(input.Folder)
	state.Id = types.StringValue(idBuilder.String())
	var var0 []bgpRoutingListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]bgpRoutingListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 bgpRoutingListDsModelConfig
			var var3 *bgpRoutingListDsModelRoutingPreferenceObject
			if var1.RoutingPreference != nil {
				var3 = &bgpRoutingListDsModelRoutingPreferenceObject{}
				if var1.RoutingPreference.Default != nil {
					var3.Default = types.BoolValue(true)
				}
				if var1.RoutingPreference.HotPotatoRouting != nil {
					var3.HotPotatoRouting = types.BoolValue(true)
				}
			}
			var2.AcceptRouteOverSC = types.BoolValue(var1.AcceptRouteOverSC)
			var2.AddHostRouteToIkePeer = types.BoolValue(var1.AddHostRouteToIkePeer)
			var2.BackboneRouting = types.StringValue(var1.BackboneRouting)
			var2.OutboundRoutesForServices = EncodeStringSlice(var1.OutboundRoutesForServices)
			var2.RoutingPreference = var3
			var2.WithdrawStaticRoute = types.BoolValue(var1.WithdrawStaticRoute)
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
