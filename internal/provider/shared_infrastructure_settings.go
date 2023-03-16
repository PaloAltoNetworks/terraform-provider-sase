package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	wHGMtfb "github.com/paloaltonetworks/sase-go/netsec/service/v1/sharedinfrastructuresettings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Data source (listing).
var (
	_ datasource.DataSource              = &sharedInfrastructureSettingsListDataSource{}
	_ datasource.DataSourceWithConfigure = &sharedInfrastructureSettingsListDataSource{}
)

func NewSharedInfrastructureSettingsListDataSource() datasource.DataSource {
	return &sharedInfrastructureSettingsListDataSource{}
}

type sharedInfrastructureSettingsListDataSource struct {
	client *sase.Client
}

type sharedInfrastructureSettingsListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64 `tfsdk:"limit"`
	Offset types.Int64 `tfsdk:"offset"`

	// Output.
	Data []sharedInfrastructureSettingsListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type sharedInfrastructureSettingsListDsModelConfig struct {
	ApiKey                         types.String   `tfsdk:"api_key"`
	CaptivePortalRedirectIpAddress types.String   `tfsdk:"captive_portal_redirect_ip_address"`
	EgressIpNotificationUrl        types.String   `tfsdk:"egress_ip_notification_url"`
	InfraBgpAs                     types.String   `tfsdk:"infra_bgp_as"`
	InfrastructureSubnet           types.String   `tfsdk:"infrastructure_subnet"`
	InfrastructureSubnetIpv6       types.String   `tfsdk:"infrastructure_subnet_ipv6"`
	Ipv6                           types.Bool     `tfsdk:"ipv6"`
	LoopbackIps                    []types.String `tfsdk:"loopback_ips"`
	TunnelMonitorIpAddress         types.String   `tfsdk:"tunnel_monitor_ip_address"`
}

// Metadata returns the data source type name.
func (d *sharedInfrastructureSettingsListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_shared_infrastructure_settings_list"
}

// Schema defines the schema for this listing data source.
func (d *sharedInfrastructureSettingsListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"api_key": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"captive_portal_redirect_ip_address": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"egress_ip_notification_url": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"infra_bgp_as": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"infrastructure_subnet": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"infrastructure_subnet_ipv6": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"ipv6": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"loopback_ips": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"tunnel_monitor_ip_address": dsschema.StringAttribute{
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
func (d *sharedInfrastructureSettingsListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *sharedInfrastructureSettingsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sharedInfrastructureSettingsListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_shared_infrastructure_settings_list",
		"limit":            state.Limit.ValueInt64(),
		"has_limit":        !state.Limit.IsNull(),
		"offset":           state.Offset.ValueInt64(),
		"has_offset":       !state.Offset.IsNull(),
	})

	// Prepare to run the command.
	svc := wHGMtfb.NewClient(d.client)
	input := wHGMtfb.ListInput{}
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
	state.Id = types.StringValue(strings.Join([]string{strconv.FormatInt(*input.Limit, 10), strconv.FormatInt(*input.Offset, 10)}, IdSeparator))
	var var0 []sharedInfrastructureSettingsListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]sharedInfrastructureSettingsListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 sharedInfrastructureSettingsListDsModelConfig
			var2.ApiKey = types.StringValue(var1.ApiKey)
			var2.CaptivePortalRedirectIpAddress = types.StringValue(var1.CaptivePortalRedirectIpAddress)
			var2.EgressIpNotificationUrl = types.StringValue(var1.EgressIpNotificationUrl)
			var2.InfraBgpAs = types.StringValue(var1.InfraBgpAs)
			var2.InfrastructureSubnet = types.StringValue(var1.InfrastructureSubnet)
			var2.InfrastructureSubnetIpv6 = types.StringValue(var1.InfrastructureSubnetIpv6)
			var2.Ipv6 = types.BoolValue(var1.Ipv6)
			var2.LoopbackIps = EncodeStringSlice(var1.LoopbackIps)
			var2.TunnelMonitorIpAddress = types.StringValue(var1.TunnelMonitorIpAddress)
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
