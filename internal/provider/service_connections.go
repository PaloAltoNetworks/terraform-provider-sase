package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	yaiLoaU "github.com/paloaltonetworks/sase-go/netsec/service/v1/serviceconnections"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Data source (listing).
var (
	_ datasource.DataSource              = &serviceConnectionsListDataSource{}
	_ datasource.DataSourceWithConfigure = &serviceConnectionsListDataSource{}
)

func NewServiceConnectionsListDataSource() datasource.DataSource {
	return &serviceConnectionsListDataSource{}
}

type serviceConnectionsListDataSource struct {
	client *sase.Client
}

type serviceConnectionsListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []serviceConnectionsListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type serviceConnectionsListDsModelConfig struct {
	BackupSC             types.String                                 `tfsdk:"backup_s_c"`
	BgpPeer              *serviceConnectionsListDsModelBgpPeerObject  `tfsdk:"bgp_peer"`
	IpsecTunnel          types.String                                 `tfsdk:"ipsec_tunnel"`
	Name                 types.String                                 `tfsdk:"name"`
	NatPool              types.String                                 `tfsdk:"nat_pool"`
	NoExportCommunity    types.String                                 `tfsdk:"no_export_community"`
	OnboardingType       types.String                                 `tfsdk:"onboarding_type"`
	Protocol             *serviceConnectionsListDsModelProtocolObject `tfsdk:"protocol"`
	Qos                  *serviceConnectionsListDsModelQosObject      `tfsdk:"qos"`
	Region               types.String                                 `tfsdk:"region"`
	SecondaryIpsecTunnel types.String                                 `tfsdk:"secondary_ipsec_tunnel"`
	SourceNat            types.Bool                                   `tfsdk:"source_nat"`
	Subnets              []types.String                               `tfsdk:"subnets"`
}

type serviceConnectionsListDsModelBgpPeerObject struct {
	LocalIpAddress   types.String `tfsdk:"local_ip_address"`
	LocalIpv6Address types.String `tfsdk:"local_ipv6_address"`
	PeerIpAddress    types.String `tfsdk:"peer_ip_address"`
	PeerIpv6Address  types.String `tfsdk:"peer_ipv6_address"`
	SameAsPrimary    types.Bool   `tfsdk:"same_as_primary"`
	Secret           types.String `tfsdk:"secret"`
}

type serviceConnectionsListDsModelProtocolObject struct {
	Bgp *serviceConnectionsListDsModelBgpObject `tfsdk:"bgp"`
}

type serviceConnectionsListDsModelBgpObject struct {
	DoNotExportRoutes         types.Bool   `tfsdk:"do_not_export_routes"`
	Enable                    types.Bool   `tfsdk:"enable"`
	FastFailover              types.Bool   `tfsdk:"fast_failover"`
	LocalIpAddress            types.String `tfsdk:"local_ip_address"`
	OriginateDefaultRoute     types.Bool   `tfsdk:"originate_default_route"`
	PeerAs                    types.String `tfsdk:"peer_as"`
	PeerIpAddress             types.String `tfsdk:"peer_ip_address"`
	Secret                    types.String `tfsdk:"secret"`
	SummarizeMobileUserRoutes types.Bool   `tfsdk:"summarize_mobile_user_routes"`
}

type serviceConnectionsListDsModelQosObject struct {
	Enable     types.Bool   `tfsdk:"enable"`
	QosProfile types.String `tfsdk:"qos_profile"`
}

// Metadata returns the data source type name.
func (d *serviceConnectionsListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_connections_list"
}

// Schema defines the schema for this listing data source.
func (d *serviceConnectionsListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"backup_s_c": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"bgp_peer": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"local_ip_address": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
								"local_ipv6_address": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
								"peer_ip_address": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
								"peer_ipv6_address": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
								"same_as_primary": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"secret": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
							},
						},
						"ipsec_tunnel": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"nat_pool": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"no_export_community": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"onboarding_type": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"protocol": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"bgp": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"do_not_export_routes": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
										"enable": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
										"fast_failover": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
										"local_ip_address": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"originate_default_route": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
										"peer_as": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"peer_ip_address": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"secret": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"summarize_mobile_user_routes": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
							},
						},
						"qos": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"enable": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"qos_profile": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
							},
						},
						"region": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"secondary_ipsec_tunnel": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"source_nat": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"subnets": dsschema.ListAttribute{
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
func (d *serviceConnectionsListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *serviceConnectionsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serviceConnectionsListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_service_connections_list",
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
	svc := yaiLoaU.NewClient(d.client)
	input := yaiLoaU.ListInput{
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
	var var0 []serviceConnectionsListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]serviceConnectionsListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 serviceConnectionsListDsModelConfig
			var var3 *serviceConnectionsListDsModelBgpPeerObject
			if var1.BgpPeer != nil {
				var3 = &serviceConnectionsListDsModelBgpPeerObject{}
				var3.LocalIpAddress = types.StringValue(var1.BgpPeer.LocalIpAddress)
				var3.LocalIpv6Address = types.StringValue(var1.BgpPeer.LocalIpv6Address)
				var3.PeerIpAddress = types.StringValue(var1.BgpPeer.PeerIpAddress)
				var3.PeerIpv6Address = types.StringValue(var1.BgpPeer.PeerIpv6Address)
				var3.SameAsPrimary = types.BoolValue(var1.BgpPeer.SameAsPrimary)
				var3.Secret = types.StringValue(var1.BgpPeer.Secret)
			}
			var var4 *serviceConnectionsListDsModelProtocolObject
			if var1.Protocol != nil {
				var4 = &serviceConnectionsListDsModelProtocolObject{}
				var var5 *serviceConnectionsListDsModelBgpObject
				if var1.Protocol.Bgp != nil {
					var5 = &serviceConnectionsListDsModelBgpObject{}
					var5.DoNotExportRoutes = types.BoolValue(var1.Protocol.Bgp.DoNotExportRoutes)
					var5.Enable = types.BoolValue(var1.Protocol.Bgp.Enable)
					var5.FastFailover = types.BoolValue(var1.Protocol.Bgp.FastFailover)
					var5.LocalIpAddress = types.StringValue(var1.Protocol.Bgp.LocalIpAddress)
					var5.OriginateDefaultRoute = types.BoolValue(var1.Protocol.Bgp.OriginateDefaultRoute)
					var5.PeerAs = types.StringValue(var1.Protocol.Bgp.PeerAs)
					var5.PeerIpAddress = types.StringValue(var1.Protocol.Bgp.PeerIpAddress)
					var5.Secret = types.StringValue(var1.Protocol.Bgp.Secret)
					var5.SummarizeMobileUserRoutes = types.BoolValue(var1.Protocol.Bgp.SummarizeMobileUserRoutes)
				}
				var4.Bgp = var5
			}
			var var6 *serviceConnectionsListDsModelQosObject
			if var1.Qos != nil {
				var6 = &serviceConnectionsListDsModelQosObject{}
				var6.Enable = types.BoolValue(var1.Qos.Enable)
				var6.QosProfile = types.StringValue(var1.Qos.QosProfile)
			}
			var2.BackupSC = types.StringValue(var1.BackupSC)
			var2.BgpPeer = var3
			var2.IpsecTunnel = types.StringValue(var1.IpsecTunnel)
			var2.Name = types.StringValue(var1.Name)
			var2.NatPool = types.StringValue(var1.NatPool)
			var2.NoExportCommunity = types.StringValue(var1.NoExportCommunity)
			var2.OnboardingType = types.StringValue(var1.OnboardingType)
			var2.Protocol = var4
			var2.Qos = var6
			var2.Region = types.StringValue(var1.Region)
			var2.SecondaryIpsecTunnel = types.StringValue(var1.SecondaryIpsecTunnel)
			var2.SourceNat = types.BoolValue(var1.SourceNat)
			var2.Subnets = EncodeStringSlice(var1.Subnets)
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
