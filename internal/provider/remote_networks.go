package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	jhtSIUK "github.com/paloaltonetworks/sase-go/netsec/schema/remote/networks"
	xsuBWMo "github.com/paloaltonetworks/sase-go/netsec/service/v1/remotenetworks"

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
	_ datasource.DataSource              = &remoteNetworksListDataSource{}
	_ datasource.DataSourceWithConfigure = &remoteNetworksListDataSource{}
)

func NewRemoteNetworksListDataSource() datasource.DataSource {
	return &remoteNetworksListDataSource{}
}

type remoteNetworksListDataSource struct {
	client *sase.Client
}

type remoteNetworksListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []remoteNetworksListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type remoteNetworksListDsModelConfig struct {
	EcmpLoadBalancing    types.String                                 `tfsdk:"ecmp_load_balancing"`
	EcmpTunnels          []remoteNetworksListDsModelEcmpTunnelsObject `tfsdk:"ecmp_tunnels"`
	ObjectId             types.String                                 `tfsdk:"object_id"`
	IpsecTunnel          types.String                                 `tfsdk:"ipsec_tunnel"`
	LicenseType          types.String                                 `tfsdk:"license_type"`
	Name                 types.String                                 `tfsdk:"name"`
	Protocol             *remoteNetworksListDsModelProtocolObject     `tfsdk:"protocol"`
	Region               types.String                                 `tfsdk:"region"`
	SecondaryIpsecTunnel types.String                                 `tfsdk:"secondary_ipsec_tunnel"`
	SpnName              types.String                                 `tfsdk:"spn_name"`
	Subnets              []types.String                               `tfsdk:"subnets"`
}

type remoteNetworksListDsModelEcmpTunnelsObject struct {
	DoNotExportRoutes         types.Bool   `tfsdk:"do_not_export_routes"`
	IpsecTunnel               types.String `tfsdk:"ipsec_tunnel"`
	LocalIpAddress            types.String `tfsdk:"local_ip_address"`
	Name                      types.String `tfsdk:"name"`
	OriginateDefaultRoute     types.Bool   `tfsdk:"originate_default_route"`
	PeerAs                    types.String `tfsdk:"peer_as"`
	PeerIpAddress             types.String `tfsdk:"peer_ip_address"`
	PeeringType               types.String `tfsdk:"peering_type"`
	Secret                    types.String `tfsdk:"secret"`
	SummarizeMobileUserRoutes types.Bool   `tfsdk:"summarize_mobile_user_routes"`
}

type remoteNetworksListDsModelProtocolObject struct {
	Bgp     *remoteNetworksListDsModelBgpObject     `tfsdk:"bgp"`
	BgpPeer *remoteNetworksListDsModelBgpPeerObject `tfsdk:"bgp_peer"`
}

type remoteNetworksListDsModelBgpObject struct {
	DoNotExportRoutes         types.Bool   `tfsdk:"do_not_export_routes"`
	Enable                    types.Bool   `tfsdk:"enable"`
	LocalIpAddress            types.String `tfsdk:"local_ip_address"`
	OriginateDefaultRoute     types.Bool   `tfsdk:"originate_default_route"`
	PeerAs                    types.String `tfsdk:"peer_as"`
	PeerIpAddress             types.String `tfsdk:"peer_ip_address"`
	PeeringType               types.String `tfsdk:"peering_type"`
	Secret                    types.String `tfsdk:"secret"`
	SummarizeMobileUserRoutes types.Bool   `tfsdk:"summarize_mobile_user_routes"`
}

type remoteNetworksListDsModelBgpPeerObject struct {
	LocalIpAddress types.String `tfsdk:"local_ip_address"`
	PeerIpAddress  types.String `tfsdk:"peer_ip_address"`
	Secret         types.String `tfsdk:"secret"`
}

// Metadata returns the data source type name.
func (d *remoteNetworksListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote_networks_list"
}

// Schema defines the schema for this listing data source.
func (d *remoteNetworksListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

			// Output.
			"data": dsschema.ListNestedAttribute{
				Description:         "The `data` parameter.",
				MarkdownDescription: "The `data` parameter.",
				Computed:            true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"ecmp_load_balancing": dsschema.StringAttribute{
							Description:         "The `ecmp_load_balancing` parameter.",
							MarkdownDescription: "The `ecmp_load_balancing` parameter.",
							Computed:            true,
						},
						"ecmp_tunnels": dsschema.ListNestedAttribute{
							Description:         "The `ecmp_tunnels` parameter.",
							MarkdownDescription: "The `ecmp_tunnels` parameter.",
							Computed:            true,
							NestedObject: dsschema.NestedAttributeObject{
								Attributes: map[string]dsschema.Attribute{
									"do_not_export_routes": dsschema.BoolAttribute{
										Description:         "The `do_not_export_routes` parameter.",
										MarkdownDescription: "The `do_not_export_routes` parameter.",
										Computed:            true,
									},
									"ipsec_tunnel": dsschema.StringAttribute{
										Description:         "The `ipsec_tunnel` parameter.",
										MarkdownDescription: "The `ipsec_tunnel` parameter.",
										Computed:            true,
									},
									"local_ip_address": dsschema.StringAttribute{
										Description:         "The `local_ip_address` parameter.",
										MarkdownDescription: "The `local_ip_address` parameter.",
										Computed:            true,
									},
									"name": dsschema.StringAttribute{
										Description:         "The `name` parameter.",
										MarkdownDescription: "The `name` parameter.",
										Computed:            true,
									},
									"originate_default_route": dsschema.BoolAttribute{
										Description:         "The `originate_default_route` parameter.",
										MarkdownDescription: "The `originate_default_route` parameter.",
										Computed:            true,
									},
									"peer_as": dsschema.StringAttribute{
										Description:         "The `peer_as` parameter.",
										MarkdownDescription: "The `peer_as` parameter.",
										Computed:            true,
									},
									"peer_ip_address": dsschema.StringAttribute{
										Description:         "The `peer_ip_address` parameter.",
										MarkdownDescription: "The `peer_ip_address` parameter.",
										Computed:            true,
									},
									"peering_type": dsschema.StringAttribute{
										Description:         "The `peering_type` parameter.",
										MarkdownDescription: "The `peering_type` parameter.",
										Computed:            true,
									},
									"secret": dsschema.StringAttribute{
										Description:         "The `secret` parameter.",
										MarkdownDescription: "The `secret` parameter.",
										Computed:            true,
									},
									"summarize_mobile_user_routes": dsschema.BoolAttribute{
										Description:         "The `summarize_mobile_user_routes` parameter.",
										MarkdownDescription: "The `summarize_mobile_user_routes` parameter.",
										Computed:            true,
									},
								},
							},
						},
						"object_id": dsschema.StringAttribute{
							Description:         "The `object_id` parameter.",
							MarkdownDescription: "The `object_id` parameter.",
							Computed:            true,
						},
						"ipsec_tunnel": dsschema.StringAttribute{
							Description:         "The `ipsec_tunnel` parameter.",
							MarkdownDescription: "The `ipsec_tunnel` parameter.",
							Computed:            true,
						},
						"license_type": dsschema.StringAttribute{
							Description:         "The `license_type` parameter.",
							MarkdownDescription: "The `license_type` parameter.",
							Computed:            true,
						},
						"name": dsschema.StringAttribute{
							Description:         "The `name` parameter.",
							MarkdownDescription: "The `name` parameter.",
							Computed:            true,
						},
						"protocol": dsschema.SingleNestedAttribute{
							Description:         "The `protocol` parameter.",
							MarkdownDescription: "The `protocol` parameter.",
							Computed:            true,
							Attributes: map[string]dsschema.Attribute{
								"bgp": dsschema.SingleNestedAttribute{
									Description:         "The `bgp` parameter.",
									MarkdownDescription: "The `bgp` parameter.",
									Computed:            true,
									Attributes: map[string]dsschema.Attribute{
										"do_not_export_routes": dsschema.BoolAttribute{
											Description:         "The `do_not_export_routes` parameter.",
											MarkdownDescription: "The `do_not_export_routes` parameter.",
											Computed:            true,
										},
										"enable": dsschema.BoolAttribute{
											Description:         "The `enable` parameter.",
											MarkdownDescription: "The `enable` parameter.",
											Computed:            true,
										},
										"local_ip_address": dsschema.StringAttribute{
											Description:         "The `local_ip_address` parameter.",
											MarkdownDescription: "The `local_ip_address` parameter.",
											Computed:            true,
										},
										"originate_default_route": dsschema.BoolAttribute{
											Description:         "The `originate_default_route` parameter.",
											MarkdownDescription: "The `originate_default_route` parameter.",
											Computed:            true,
										},
										"peer_as": dsschema.StringAttribute{
											Description:         "The `peer_as` parameter.",
											MarkdownDescription: "The `peer_as` parameter.",
											Computed:            true,
										},
										"peer_ip_address": dsschema.StringAttribute{
											Description:         "The `peer_ip_address` parameter.",
											MarkdownDescription: "The `peer_ip_address` parameter.",
											Computed:            true,
										},
										"peering_type": dsschema.StringAttribute{
											Description:         "The `peering_type` parameter.",
											MarkdownDescription: "The `peering_type` parameter.",
											Computed:            true,
										},
										"secret": dsschema.StringAttribute{
											Description:         "The `secret` parameter.",
											MarkdownDescription: "The `secret` parameter.",
											Computed:            true,
										},
										"summarize_mobile_user_routes": dsschema.BoolAttribute{
											Description:         "The `summarize_mobile_user_routes` parameter.",
											MarkdownDescription: "The `summarize_mobile_user_routes` parameter.",
											Computed:            true,
										},
									},
								},
								"bgp_peer": dsschema.SingleNestedAttribute{
									Description:         "The `bgp_peer` parameter.",
									MarkdownDescription: "The `bgp_peer` parameter.",
									Computed:            true,
									Attributes: map[string]dsschema.Attribute{
										"local_ip_address": dsschema.StringAttribute{
											Description:         "The `local_ip_address` parameter.",
											MarkdownDescription: "The `local_ip_address` parameter.",
											Computed:            true,
										},
										"peer_ip_address": dsschema.StringAttribute{
											Description:         "The `peer_ip_address` parameter.",
											MarkdownDescription: "The `peer_ip_address` parameter.",
											Computed:            true,
										},
										"secret": dsschema.StringAttribute{
											Description:         "The `secret` parameter.",
											MarkdownDescription: "The `secret` parameter.",
											Computed:            true,
										},
									},
								},
							},
						},
						"region": dsschema.StringAttribute{
							Description:         "The `region` parameter.",
							MarkdownDescription: "The `region` parameter.",
							Computed:            true,
						},
						"secondary_ipsec_tunnel": dsschema.StringAttribute{
							Description:         "The `secondary_ipsec_tunnel` parameter.",
							MarkdownDescription: "The `secondary_ipsec_tunnel` parameter.",
							Computed:            true,
						},
						"spn_name": dsschema.StringAttribute{
							Description:         "The `spn_name` parameter.",
							MarkdownDescription: "The `spn_name` parameter.",
							Computed:            true,
						},
						"subnets": dsschema.ListAttribute{
							Description:         "The `subnets` parameter.",
							MarkdownDescription: "The `subnets` parameter.",
							Computed:            true,
							ElementType:         types.StringType,
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
func (d *remoteNetworksListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *remoteNetworksListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state remoteNetworksListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_remote_networks_list",
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
	svc := xsuBWMo.NewClient(d.client)
	input := xsuBWMo.ListInput{
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
	var var0 []remoteNetworksListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]remoteNetworksListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 remoteNetworksListDsModelConfig
			var var3 []remoteNetworksListDsModelEcmpTunnelsObject
			if len(var1.EcmpTunnels) != 0 {
				var3 = make([]remoteNetworksListDsModelEcmpTunnelsObject, 0, len(var1.EcmpTunnels))
				for var4Index := range var1.EcmpTunnels {
					var4 := var1.EcmpTunnels[var4Index]
					var var5 remoteNetworksListDsModelEcmpTunnelsObject
					var5.DoNotExportRoutes = types.BoolValue(var4.DoNotExportRoutes)
					var5.IpsecTunnel = types.StringValue(var4.IpsecTunnel)
					var5.LocalIpAddress = types.StringValue(var4.LocalIpAddress)
					var5.Name = types.StringValue(var4.Name)
					var5.OriginateDefaultRoute = types.BoolValue(var4.OriginateDefaultRoute)
					var5.PeerAs = types.StringValue(var4.PeerAs)
					var5.PeerIpAddress = types.StringValue(var4.PeerIpAddress)
					var5.PeeringType = types.StringValue(var4.PeeringType)
					var5.Secret = types.StringValue(var4.Secret)
					var5.SummarizeMobileUserRoutes = types.BoolValue(var4.SummarizeMobileUserRoutes)
					var3 = append(var3, var5)
				}
			}
			var var6 *remoteNetworksListDsModelProtocolObject
			if var1.Protocol != nil {
				var6 = &remoteNetworksListDsModelProtocolObject{}
				var var7 *remoteNetworksListDsModelBgpObject
				if var1.Protocol.Bgp != nil {
					var7 = &remoteNetworksListDsModelBgpObject{}
					var7.DoNotExportRoutes = types.BoolValue(var1.Protocol.Bgp.DoNotExportRoutes)
					var7.Enable = types.BoolValue(var1.Protocol.Bgp.Enable)
					var7.LocalIpAddress = types.StringValue(var1.Protocol.Bgp.LocalIpAddress)
					var7.OriginateDefaultRoute = types.BoolValue(var1.Protocol.Bgp.OriginateDefaultRoute)
					var7.PeerAs = types.StringValue(var1.Protocol.Bgp.PeerAs)
					var7.PeerIpAddress = types.StringValue(var1.Protocol.Bgp.PeerIpAddress)
					var7.PeeringType = types.StringValue(var1.Protocol.Bgp.PeeringType)
					var7.Secret = types.StringValue(var1.Protocol.Bgp.Secret)
					var7.SummarizeMobileUserRoutes = types.BoolValue(var1.Protocol.Bgp.SummarizeMobileUserRoutes)
				}
				var var8 *remoteNetworksListDsModelBgpPeerObject
				if var1.Protocol.BgpPeer != nil {
					var8 = &remoteNetworksListDsModelBgpPeerObject{}
					var8.LocalIpAddress = types.StringValue(var1.Protocol.BgpPeer.LocalIpAddress)
					var8.PeerIpAddress = types.StringValue(var1.Protocol.BgpPeer.PeerIpAddress)
					var8.Secret = types.StringValue(var1.Protocol.BgpPeer.Secret)
				}
				var6.Bgp = var7
				var6.BgpPeer = var8
			}
			var2.EcmpLoadBalancing = types.StringValue(var1.EcmpLoadBalancing)
			var2.EcmpTunnels = var3
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.IpsecTunnel = types.StringValue(var1.IpsecTunnel)
			var2.LicenseType = types.StringValue(var1.LicenseType)
			var2.Name = types.StringValue(var1.Name)
			var2.Protocol = var6
			var2.Region = types.StringValue(var1.Region)
			var2.SecondaryIpsecTunnel = types.StringValue(var1.SecondaryIpsecTunnel)
			var2.SpnName = types.StringValue(var1.SpnName)
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

// Data source.
var (
	_ datasource.DataSource              = &remoteNetworksDataSource{}
	_ datasource.DataSourceWithConfigure = &remoteNetworksDataSource{}
)

func NewRemoteNetworksDataSource() datasource.DataSource {
	return &remoteNetworksDataSource{}
}

type remoteNetworksDataSource struct {
	client *sase.Client
}

type remoteNetworksDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`
	Folder   types.String `tfsdk:"folder"`

	// Output.
	// Ref: #/components/schemas/remote-networks
	EcmpLoadBalancing types.String                             `tfsdk:"ecmp_load_balancing"`
	EcmpTunnels       []remoteNetworksDsModelEcmpTunnelsObject `tfsdk:"ecmp_tunnels"`
	// input omit: ObjectId
	IpsecTunnel          types.String                         `tfsdk:"ipsec_tunnel"`
	LicenseType          types.String                         `tfsdk:"license_type"`
	Name                 types.String                         `tfsdk:"name"`
	Protocol             *remoteNetworksDsModelProtocolObject `tfsdk:"protocol"`
	Region               types.String                         `tfsdk:"region"`
	SecondaryIpsecTunnel types.String                         `tfsdk:"secondary_ipsec_tunnel"`
	SpnName              types.String                         `tfsdk:"spn_name"`
	Subnets              []types.String                       `tfsdk:"subnets"`
}

type remoteNetworksDsModelEcmpTunnelsObject struct {
	DoNotExportRoutes         types.Bool   `tfsdk:"do_not_export_routes"`
	IpsecTunnel               types.String `tfsdk:"ipsec_tunnel"`
	LocalIpAddress            types.String `tfsdk:"local_ip_address"`
	Name                      types.String `tfsdk:"name"`
	OriginateDefaultRoute     types.Bool   `tfsdk:"originate_default_route"`
	PeerAs                    types.String `tfsdk:"peer_as"`
	PeerIpAddress             types.String `tfsdk:"peer_ip_address"`
	PeeringType               types.String `tfsdk:"peering_type"`
	Secret                    types.String `tfsdk:"secret"`
	SummarizeMobileUserRoutes types.Bool   `tfsdk:"summarize_mobile_user_routes"`
}

type remoteNetworksDsModelProtocolObject struct {
	Bgp     *remoteNetworksDsModelBgpObject     `tfsdk:"bgp"`
	BgpPeer *remoteNetworksDsModelBgpPeerObject `tfsdk:"bgp_peer"`
}

type remoteNetworksDsModelBgpObject struct {
	DoNotExportRoutes         types.Bool   `tfsdk:"do_not_export_routes"`
	Enable                    types.Bool   `tfsdk:"enable"`
	LocalIpAddress            types.String `tfsdk:"local_ip_address"`
	OriginateDefaultRoute     types.Bool   `tfsdk:"originate_default_route"`
	PeerAs                    types.String `tfsdk:"peer_as"`
	PeerIpAddress             types.String `tfsdk:"peer_ip_address"`
	PeeringType               types.String `tfsdk:"peering_type"`
	Secret                    types.String `tfsdk:"secret"`
	SummarizeMobileUserRoutes types.Bool   `tfsdk:"summarize_mobile_user_routes"`
}

type remoteNetworksDsModelBgpPeerObject struct {
	LocalIpAddress types.String `tfsdk:"local_ip_address"`
	PeerIpAddress  types.String `tfsdk:"peer_ip_address"`
	Secret         types.String `tfsdk:"secret"`
}

// Metadata returns the data source type name.
func (d *remoteNetworksDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote_networks"
}

// Schema defines the schema for this listing data source.
func (d *remoteNetworksDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"ecmp_load_balancing": dsschema.StringAttribute{
				Description:         "The `ecmp_load_balancing` parameter.",
				MarkdownDescription: "The `ecmp_load_balancing` parameter.",
				Computed:            true,
			},
			"ecmp_tunnels": dsschema.ListNestedAttribute{
				Description:         "The `ecmp_tunnels` parameter.",
				MarkdownDescription: "The `ecmp_tunnels` parameter.",
				Computed:            true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"do_not_export_routes": dsschema.BoolAttribute{
							Description:         "The `do_not_export_routes` parameter.",
							MarkdownDescription: "The `do_not_export_routes` parameter.",
							Computed:            true,
						},
						"ipsec_tunnel": dsschema.StringAttribute{
							Description:         "The `ipsec_tunnel` parameter.",
							MarkdownDescription: "The `ipsec_tunnel` parameter.",
							Computed:            true,
						},
						"local_ip_address": dsschema.StringAttribute{
							Description:         "The `local_ip_address` parameter.",
							MarkdownDescription: "The `local_ip_address` parameter.",
							Computed:            true,
						},
						"name": dsschema.StringAttribute{
							Description:         "The `name` parameter.",
							MarkdownDescription: "The `name` parameter.",
							Computed:            true,
						},
						"originate_default_route": dsschema.BoolAttribute{
							Description:         "The `originate_default_route` parameter.",
							MarkdownDescription: "The `originate_default_route` parameter.",
							Computed:            true,
						},
						"peer_as": dsschema.StringAttribute{
							Description:         "The `peer_as` parameter.",
							MarkdownDescription: "The `peer_as` parameter.",
							Computed:            true,
						},
						"peer_ip_address": dsschema.StringAttribute{
							Description:         "The `peer_ip_address` parameter.",
							MarkdownDescription: "The `peer_ip_address` parameter.",
							Computed:            true,
						},
						"peering_type": dsschema.StringAttribute{
							Description:         "The `peering_type` parameter.",
							MarkdownDescription: "The `peering_type` parameter.",
							Computed:            true,
						},
						"secret": dsschema.StringAttribute{
							Description:         "The `secret` parameter.",
							MarkdownDescription: "The `secret` parameter.",
							Computed:            true,
						},
						"summarize_mobile_user_routes": dsschema.BoolAttribute{
							Description:         "The `summarize_mobile_user_routes` parameter.",
							MarkdownDescription: "The `summarize_mobile_user_routes` parameter.",
							Computed:            true,
						},
					},
				},
			},
			"ipsec_tunnel": dsschema.StringAttribute{
				Description:         "The `ipsec_tunnel` parameter.",
				MarkdownDescription: "The `ipsec_tunnel` parameter.",
				Computed:            true,
			},
			"license_type": dsschema.StringAttribute{
				Description:         "The `license_type` parameter.",
				MarkdownDescription: "The `license_type` parameter.",
				Computed:            true,
			},
			"name": dsschema.StringAttribute{
				Description:         "The `name` parameter.",
				MarkdownDescription: "The `name` parameter.",
				Computed:            true,
			},
			"protocol": dsschema.SingleNestedAttribute{
				Description:         "The `protocol` parameter.",
				MarkdownDescription: "The `protocol` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"bgp": dsschema.SingleNestedAttribute{
						Description:         "The `bgp` parameter.",
						MarkdownDescription: "The `bgp` parameter.",
						Computed:            true,
						Attributes: map[string]dsschema.Attribute{
							"do_not_export_routes": dsschema.BoolAttribute{
								Description:         "The `do_not_export_routes` parameter.",
								MarkdownDescription: "The `do_not_export_routes` parameter.",
								Computed:            true,
							},
							"enable": dsschema.BoolAttribute{
								Description:         "The `enable` parameter.",
								MarkdownDescription: "The `enable` parameter.",
								Computed:            true,
							},
							"local_ip_address": dsschema.StringAttribute{
								Description:         "The `local_ip_address` parameter.",
								MarkdownDescription: "The `local_ip_address` parameter.",
								Computed:            true,
							},
							"originate_default_route": dsschema.BoolAttribute{
								Description:         "The `originate_default_route` parameter.",
								MarkdownDescription: "The `originate_default_route` parameter.",
								Computed:            true,
							},
							"peer_as": dsschema.StringAttribute{
								Description:         "The `peer_as` parameter.",
								MarkdownDescription: "The `peer_as` parameter.",
								Computed:            true,
							},
							"peer_ip_address": dsschema.StringAttribute{
								Description:         "The `peer_ip_address` parameter.",
								MarkdownDescription: "The `peer_ip_address` parameter.",
								Computed:            true,
							},
							"peering_type": dsschema.StringAttribute{
								Description:         "The `peering_type` parameter.",
								MarkdownDescription: "The `peering_type` parameter.",
								Computed:            true,
							},
							"secret": dsschema.StringAttribute{
								Description:         "The `secret` parameter.",
								MarkdownDescription: "The `secret` parameter.",
								Computed:            true,
							},
							"summarize_mobile_user_routes": dsschema.BoolAttribute{
								Description:         "The `summarize_mobile_user_routes` parameter.",
								MarkdownDescription: "The `summarize_mobile_user_routes` parameter.",
								Computed:            true,
							},
						},
					},
					"bgp_peer": dsschema.SingleNestedAttribute{
						Description:         "The `bgp_peer` parameter.",
						MarkdownDescription: "The `bgp_peer` parameter.",
						Computed:            true,
						Attributes: map[string]dsschema.Attribute{
							"local_ip_address": dsschema.StringAttribute{
								Description:         "The `local_ip_address` parameter.",
								MarkdownDescription: "The `local_ip_address` parameter.",
								Computed:            true,
							},
							"peer_ip_address": dsschema.StringAttribute{
								Description:         "The `peer_ip_address` parameter.",
								MarkdownDescription: "The `peer_ip_address` parameter.",
								Computed:            true,
							},
							"secret": dsschema.StringAttribute{
								Description:         "The `secret` parameter.",
								MarkdownDescription: "The `secret` parameter.",
								Computed:            true,
							},
						},
					},
				},
			},
			"region": dsschema.StringAttribute{
				Description:         "The `region` parameter.",
				MarkdownDescription: "The `region` parameter.",
				Computed:            true,
			},
			"secondary_ipsec_tunnel": dsschema.StringAttribute{
				Description:         "The `secondary_ipsec_tunnel` parameter.",
				MarkdownDescription: "The `secondary_ipsec_tunnel` parameter.",
				Computed:            true,
			},
			"spn_name": dsschema.StringAttribute{
				Description:         "The `spn_name` parameter.",
				MarkdownDescription: "The `spn_name` parameter.",
				Computed:            true,
			},
			"subnets": dsschema.ListAttribute{
				Description:         "The `subnets` parameter.",
				MarkdownDescription: "The `subnets` parameter.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

// Configure prepares the struct.
func (d *remoteNetworksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *remoteNetworksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state remoteNetworksDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_remote_networks",
		"object_id":                   state.ObjectId.ValueString(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := xsuBWMo.NewClient(d.client)
	input := xsuBWMo.ReadInput{
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
	var var0 []remoteNetworksDsModelEcmpTunnelsObject
	if len(ans.EcmpTunnels) != 0 {
		var0 = make([]remoteNetworksDsModelEcmpTunnelsObject, 0, len(ans.EcmpTunnels))
		for var1Index := range ans.EcmpTunnels {
			var1 := ans.EcmpTunnels[var1Index]
			var var2 remoteNetworksDsModelEcmpTunnelsObject
			var2.DoNotExportRoutes = types.BoolValue(var1.DoNotExportRoutes)
			var2.IpsecTunnel = types.StringValue(var1.IpsecTunnel)
			var2.LocalIpAddress = types.StringValue(var1.LocalIpAddress)
			var2.Name = types.StringValue(var1.Name)
			var2.OriginateDefaultRoute = types.BoolValue(var1.OriginateDefaultRoute)
			var2.PeerAs = types.StringValue(var1.PeerAs)
			var2.PeerIpAddress = types.StringValue(var1.PeerIpAddress)
			var2.PeeringType = types.StringValue(var1.PeeringType)
			var2.Secret = types.StringValue(var1.Secret)
			var2.SummarizeMobileUserRoutes = types.BoolValue(var1.SummarizeMobileUserRoutes)
			var0 = append(var0, var2)
		}
	}
	var var3 *remoteNetworksDsModelProtocolObject
	if ans.Protocol != nil {
		var3 = &remoteNetworksDsModelProtocolObject{}
		var var4 *remoteNetworksDsModelBgpObject
		if ans.Protocol.Bgp != nil {
			var4 = &remoteNetworksDsModelBgpObject{}
			var4.DoNotExportRoutes = types.BoolValue(ans.Protocol.Bgp.DoNotExportRoutes)
			var4.Enable = types.BoolValue(ans.Protocol.Bgp.Enable)
			var4.LocalIpAddress = types.StringValue(ans.Protocol.Bgp.LocalIpAddress)
			var4.OriginateDefaultRoute = types.BoolValue(ans.Protocol.Bgp.OriginateDefaultRoute)
			var4.PeerAs = types.StringValue(ans.Protocol.Bgp.PeerAs)
			var4.PeerIpAddress = types.StringValue(ans.Protocol.Bgp.PeerIpAddress)
			var4.PeeringType = types.StringValue(ans.Protocol.Bgp.PeeringType)
			var4.Secret = types.StringValue(ans.Protocol.Bgp.Secret)
			var4.SummarizeMobileUserRoutes = types.BoolValue(ans.Protocol.Bgp.SummarizeMobileUserRoutes)
		}
		var var5 *remoteNetworksDsModelBgpPeerObject
		if ans.Protocol.BgpPeer != nil {
			var5 = &remoteNetworksDsModelBgpPeerObject{}
			var5.LocalIpAddress = types.StringValue(ans.Protocol.BgpPeer.LocalIpAddress)
			var5.PeerIpAddress = types.StringValue(ans.Protocol.BgpPeer.PeerIpAddress)
			var5.Secret = types.StringValue(ans.Protocol.BgpPeer.Secret)
		}
		var3.Bgp = var4
		var3.BgpPeer = var5
	}
	state.EcmpLoadBalancing = types.StringValue(ans.EcmpLoadBalancing)
	state.EcmpTunnels = var0
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.IpsecTunnel = types.StringValue(ans.IpsecTunnel)
	state.LicenseType = types.StringValue(ans.LicenseType)
	state.Name = types.StringValue(ans.Name)
	state.Protocol = var3
	state.Region = types.StringValue(ans.Region)
	state.SecondaryIpsecTunnel = types.StringValue(ans.SecondaryIpsecTunnel)
	state.SpnName = types.StringValue(ans.SpnName)
	state.Subnets = EncodeStringSlice(ans.Subnets)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &remoteNetworksResource{}
	_ resource.ResourceWithConfigure   = &remoteNetworksResource{}
	_ resource.ResourceWithImportState = &remoteNetworksResource{}
)

func NewRemoteNetworksResource() resource.Resource {
	return &remoteNetworksResource{}
}

type remoteNetworksResource struct {
	client *sase.Client
}

type remoteNetworksRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/remote-networks
	EcmpLoadBalancing    types.String                             `tfsdk:"ecmp_load_balancing"`
	EcmpTunnels          []remoteNetworksRsModelEcmpTunnelsObject `tfsdk:"ecmp_tunnels"`
	ObjectId             types.String                             `tfsdk:"object_id"`
	IpsecTunnel          types.String                             `tfsdk:"ipsec_tunnel"`
	LicenseType          types.String                             `tfsdk:"license_type"`
	Name                 types.String                             `tfsdk:"name"`
	Protocol             *remoteNetworksRsModelProtocolObject     `tfsdk:"protocol"`
	Region               types.String                             `tfsdk:"region"`
	SecondaryIpsecTunnel types.String                             `tfsdk:"secondary_ipsec_tunnel"`
	SpnName              types.String                             `tfsdk:"spn_name"`
	Subnets              []types.String                           `tfsdk:"subnets"`
}

type remoteNetworksRsModelEcmpTunnelsObject struct {
	DoNotExportRoutes         types.Bool   `tfsdk:"do_not_export_routes"`
	IpsecTunnel               types.String `tfsdk:"ipsec_tunnel"`
	LocalIpAddress            types.String `tfsdk:"local_ip_address"`
	Name                      types.String `tfsdk:"name"`
	OriginateDefaultRoute     types.Bool   `tfsdk:"originate_default_route"`
	PeerAs                    types.String `tfsdk:"peer_as"`
	PeerIpAddress             types.String `tfsdk:"peer_ip_address"`
	PeeringType               types.String `tfsdk:"peering_type"`
	Secret                    types.String `tfsdk:"secret"`
	SummarizeMobileUserRoutes types.Bool   `tfsdk:"summarize_mobile_user_routes"`
}

type remoteNetworksRsModelProtocolObject struct {
	Bgp     *remoteNetworksRsModelBgpObject     `tfsdk:"bgp"`
	BgpPeer *remoteNetworksRsModelBgpPeerObject `tfsdk:"bgp_peer"`
}

type remoteNetworksRsModelBgpObject struct {
	DoNotExportRoutes         types.Bool   `tfsdk:"do_not_export_routes"`
	Enable                    types.Bool   `tfsdk:"enable"`
	LocalIpAddress            types.String `tfsdk:"local_ip_address"`
	OriginateDefaultRoute     types.Bool   `tfsdk:"originate_default_route"`
	PeerAs                    types.String `tfsdk:"peer_as"`
	PeerIpAddress             types.String `tfsdk:"peer_ip_address"`
	PeeringType               types.String `tfsdk:"peering_type"`
	Secret                    types.String `tfsdk:"secret"`
	SummarizeMobileUserRoutes types.Bool   `tfsdk:"summarize_mobile_user_routes"`
}

type remoteNetworksRsModelBgpPeerObject struct {
	LocalIpAddress types.String `tfsdk:"local_ip_address"`
	PeerIpAddress  types.String `tfsdk:"peer_ip_address"`
	Secret         types.String `tfsdk:"secret"`
}

// Metadata returns the data source type name.
func (r *remoteNetworksResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote_networks"
}

// Schema defines the schema for this listing data source.
func (r *remoteNetworksResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

			"ecmp_load_balancing": rsschema.StringAttribute{
				Description:         "The `ecmp_load_balancing` parameter. Default: `%!q(*string=0xc0006100d0)`. Value must be one of: `\"enable\"`, `\"disable\"`.",
				MarkdownDescription: "The `ecmp_load_balancing` parameter. Default: `%!q(*string=0xc0006100d0)`. Value must be one of: `\"enable\"`, `\"disable\"`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					DefaultString("disable"),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("enable", "disable"),
				},
			},
			"ecmp_tunnels": rsschema.ListNestedAttribute{
				Description:         "The `ecmp_tunnels` parameter.",
				MarkdownDescription: "The `ecmp_tunnels` parameter.",
				Optional:            true,
				NestedObject: rsschema.NestedAttributeObject{
					Attributes: map[string]rsschema.Attribute{
						"do_not_export_routes": rsschema.BoolAttribute{
							Description:         "The `do_not_export_routes` parameter.",
							MarkdownDescription: "The `do_not_export_routes` parameter.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Bool{
								DefaultBool(false),
							},
						},
						"ipsec_tunnel": rsschema.StringAttribute{
							Description:         "The `ipsec_tunnel` parameter.",
							MarkdownDescription: "The `ipsec_tunnel` parameter.",
							Required:            true,
						},
						"local_ip_address": rsschema.StringAttribute{
							Description:         "The `local_ip_address` parameter.",
							MarkdownDescription: "The `local_ip_address` parameter.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
						},
						"name": rsschema.StringAttribute{
							Description:         "The `name` parameter.",
							MarkdownDescription: "The `name` parameter.",
							Required:            true,
						},
						"originate_default_route": rsschema.BoolAttribute{
							Description:         "The `originate_default_route` parameter.",
							MarkdownDescription: "The `originate_default_route` parameter.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Bool{
								DefaultBool(false),
							},
						},
						"peer_as": rsschema.StringAttribute{
							Description:         "The `peer_as` parameter.",
							MarkdownDescription: "The `peer_as` parameter.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
						},
						"peer_ip_address": rsschema.StringAttribute{
							Description:         "The `peer_ip_address` parameter.",
							MarkdownDescription: "The `peer_ip_address` parameter.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
						},
						"peering_type": rsschema.StringAttribute{
							Description:         "The `peering_type` parameter. Value must be one of: `\"exchange-v4-over-v4\"`, `\"exchange-v4-v6-over-v4\"`, `\"exchange-v4-over-v4-v6-over-v6\"`, `\"exchange-v6-over-v6\"`.",
							MarkdownDescription: "The `peering_type` parameter. Value must be one of: `\"exchange-v4-over-v4\"`, `\"exchange-v4-v6-over-v4\"`, `\"exchange-v4-over-v4-v6-over-v6\"`, `\"exchange-v6-over-v6\"`.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
							Validators: []validator.String{
								stringvalidator.OneOf("exchange-v4-over-v4", "exchange-v4-v6-over-v4", "exchange-v4-over-v4-v6-over-v6", "exchange-v6-over-v6"),
							},
						},
						"secret": rsschema.StringAttribute{
							Description:         "The `secret` parameter.",
							MarkdownDescription: "The `secret` parameter.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
						},
						"summarize_mobile_user_routes": rsschema.BoolAttribute{
							Description:         "The `summarize_mobile_user_routes` parameter.",
							MarkdownDescription: "The `summarize_mobile_user_routes` parameter.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Bool{
								DefaultBool(false),
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
			"ipsec_tunnel": rsschema.StringAttribute{
				Description:         "The `ipsec_tunnel` parameter.",
				MarkdownDescription: "The `ipsec_tunnel` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
			},
			"license_type": rsschema.StringAttribute{
				Description:         "The `license_type` parameter. String length must be at least 1.",
				MarkdownDescription: "The `license_type` parameter. String length must be at least 1.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": rsschema.StringAttribute{
				Description:         "The `name` parameter. String length must be at most 63.",
				MarkdownDescription: "The `name` parameter. String length must be at most 63.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(63),
				},
			},
			"protocol": rsschema.SingleNestedAttribute{
				Description:         "The `protocol` parameter.",
				MarkdownDescription: "The `protocol` parameter.",
				Optional:            true,
				Attributes: map[string]rsschema.Attribute{
					"bgp": rsschema.SingleNestedAttribute{
						Description:         "The `bgp` parameter.",
						MarkdownDescription: "The `bgp` parameter.",
						Optional:            true,
						Attributes: map[string]rsschema.Attribute{
							"do_not_export_routes": rsschema.BoolAttribute{
								Description:         "The `do_not_export_routes` parameter.",
								MarkdownDescription: "The `do_not_export_routes` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.Bool{
									DefaultBool(false),
								},
							},
							"enable": rsschema.BoolAttribute{
								Description:         "The `enable` parameter.",
								MarkdownDescription: "The `enable` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.Bool{
									DefaultBool(false),
								},
							},
							"local_ip_address": rsschema.StringAttribute{
								Description:         "The `local_ip_address` parameter.",
								MarkdownDescription: "The `local_ip_address` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"originate_default_route": rsschema.BoolAttribute{
								Description:         "The `originate_default_route` parameter.",
								MarkdownDescription: "The `originate_default_route` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.Bool{
									DefaultBool(false),
								},
							},
							"peer_as": rsschema.StringAttribute{
								Description:         "The `peer_as` parameter.",
								MarkdownDescription: "The `peer_as` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"peer_ip_address": rsschema.StringAttribute{
								Description:         "The `peer_ip_address` parameter.",
								MarkdownDescription: "The `peer_ip_address` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"peering_type": rsschema.StringAttribute{
								Description:         "The `peering_type` parameter. Value must be one of: `\"exchange-v4-over-v4\"`, `\"exchange-v4-v6-over-v4\"`, `\"exchange-v4-over-v4-v6-over-v6\"`, `\"exchange-v6-over-v6\"`.",
								MarkdownDescription: "The `peering_type` parameter. Value must be one of: `\"exchange-v4-over-v4\"`, `\"exchange-v4-v6-over-v4\"`, `\"exchange-v4-over-v4-v6-over-v6\"`, `\"exchange-v6-over-v6\"`.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
								Validators: []validator.String{
									stringvalidator.OneOf("exchange-v4-over-v4", "exchange-v4-v6-over-v4", "exchange-v4-over-v4-v6-over-v6", "exchange-v6-over-v6"),
								},
							},
							"secret": rsschema.StringAttribute{
								Description:         "The `secret` parameter.",
								MarkdownDescription: "The `secret` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"summarize_mobile_user_routes": rsschema.BoolAttribute{
								Description:         "The `summarize_mobile_user_routes` parameter.",
								MarkdownDescription: "The `summarize_mobile_user_routes` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.Bool{
									DefaultBool(false),
								},
							},
						},
					},
					"bgp_peer": rsschema.SingleNestedAttribute{
						Description:         "The `bgp_peer` parameter.",
						MarkdownDescription: "The `bgp_peer` parameter.",
						Optional:            true,
						Attributes: map[string]rsschema.Attribute{
							"local_ip_address": rsschema.StringAttribute{
								Description:         "The `local_ip_address` parameter.",
								MarkdownDescription: "The `local_ip_address` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"peer_ip_address": rsschema.StringAttribute{
								Description:         "The `peer_ip_address` parameter.",
								MarkdownDescription: "The `peer_ip_address` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"secret": rsschema.StringAttribute{
								Description:         "The `secret` parameter.",
								MarkdownDescription: "The `secret` parameter.",
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
			"region": rsschema.StringAttribute{
				Description:         "The `region` parameter. String length must be at least 1.",
				MarkdownDescription: "The `region` parameter. String length must be at least 1.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"secondary_ipsec_tunnel": rsschema.StringAttribute{
				Description:         "The `secondary_ipsec_tunnel` parameter.",
				MarkdownDescription: "The `secondary_ipsec_tunnel` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
			},
			"spn_name": rsschema.StringAttribute{
				Description:         "The `spn_name` parameter.",
				MarkdownDescription: "The `spn_name` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
			},
			"subnets": rsschema.ListAttribute{
				Description:         "The `subnets` parameter.",
				MarkdownDescription: "The `subnets` parameter.",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

// Configure prepares the struct.
func (r *remoteNetworksResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *remoteNetworksResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state remoteNetworksRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_remote_networks",
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := xsuBWMo.NewClient(r.client)
	input := xsuBWMo.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 jhtSIUK.Config
	var0.EcmpLoadBalancing = state.EcmpLoadBalancing.ValueString()
	var var1 []jhtSIUK.EcmpTunnelsObject
	if len(state.EcmpTunnels) != 0 {
		var1 = make([]jhtSIUK.EcmpTunnelsObject, 0, len(state.EcmpTunnels))
		for var2Index := range state.EcmpTunnels {
			var2 := state.EcmpTunnels[var2Index]
			var var3 jhtSIUK.EcmpTunnelsObject
			var3.DoNotExportRoutes = var2.DoNotExportRoutes.ValueBool()
			var3.IpsecTunnel = var2.IpsecTunnel.ValueString()
			var3.LocalIpAddress = var2.LocalIpAddress.ValueString()
			var3.Name = var2.Name.ValueString()
			var3.OriginateDefaultRoute = var2.OriginateDefaultRoute.ValueBool()
			var3.PeerAs = var2.PeerAs.ValueString()
			var3.PeerIpAddress = var2.PeerIpAddress.ValueString()
			var3.PeeringType = var2.PeeringType.ValueString()
			var3.Secret = var2.Secret.ValueString()
			var3.SummarizeMobileUserRoutes = var2.SummarizeMobileUserRoutes.ValueBool()
			var1 = append(var1, var3)
		}
	}
	var0.EcmpTunnels = var1
	var0.IpsecTunnel = state.IpsecTunnel.ValueString()
	var0.LicenseType = state.LicenseType.ValueString()
	var0.Name = state.Name.ValueString()
	var var4 *jhtSIUK.ProtocolObject
	if state.Protocol != nil {
		var4 = &jhtSIUK.ProtocolObject{}
		var var5 *jhtSIUK.BgpObject
		if state.Protocol.Bgp != nil {
			var5 = &jhtSIUK.BgpObject{}
			var5.DoNotExportRoutes = state.Protocol.Bgp.DoNotExportRoutes.ValueBool()
			var5.Enable = state.Protocol.Bgp.Enable.ValueBool()
			var5.LocalIpAddress = state.Protocol.Bgp.LocalIpAddress.ValueString()
			var5.OriginateDefaultRoute = state.Protocol.Bgp.OriginateDefaultRoute.ValueBool()
			var5.PeerAs = state.Protocol.Bgp.PeerAs.ValueString()
			var5.PeerIpAddress = state.Protocol.Bgp.PeerIpAddress.ValueString()
			var5.PeeringType = state.Protocol.Bgp.PeeringType.ValueString()
			var5.Secret = state.Protocol.Bgp.Secret.ValueString()
			var5.SummarizeMobileUserRoutes = state.Protocol.Bgp.SummarizeMobileUserRoutes.ValueBool()
		}
		var4.Bgp = var5
		var var6 *jhtSIUK.BgpPeerObject
		if state.Protocol.BgpPeer != nil {
			var6 = &jhtSIUK.BgpPeerObject{}
			var6.LocalIpAddress = state.Protocol.BgpPeer.LocalIpAddress.ValueString()
			var6.PeerIpAddress = state.Protocol.BgpPeer.PeerIpAddress.ValueString()
			var6.Secret = state.Protocol.BgpPeer.Secret.ValueString()
		}
		var4.BgpPeer = var6
	}
	var0.Protocol = var4
	var0.Region = state.Region.ValueString()
	var0.SecondaryIpsecTunnel = state.SecondaryIpsecTunnel.ValueString()
	var0.SpnName = state.SpnName.ValueString()
	var0.Subnets = DecodeStringSlice(state.Subnets)
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
	var var7 []remoteNetworksRsModelEcmpTunnelsObject
	if len(ans.EcmpTunnels) != 0 {
		var7 = make([]remoteNetworksRsModelEcmpTunnelsObject, 0, len(ans.EcmpTunnels))
		for var8Index := range ans.EcmpTunnels {
			var8 := ans.EcmpTunnels[var8Index]
			var var9 remoteNetworksRsModelEcmpTunnelsObject
			var9.DoNotExportRoutes = types.BoolValue(var8.DoNotExportRoutes)
			var9.IpsecTunnel = types.StringValue(var8.IpsecTunnel)
			var9.LocalIpAddress = types.StringValue(var8.LocalIpAddress)
			var9.Name = types.StringValue(var8.Name)
			var9.OriginateDefaultRoute = types.BoolValue(var8.OriginateDefaultRoute)
			var9.PeerAs = types.StringValue(var8.PeerAs)
			var9.PeerIpAddress = types.StringValue(var8.PeerIpAddress)
			var9.PeeringType = types.StringValue(var8.PeeringType)
			var9.Secret = types.StringValue(var8.Secret)
			var9.SummarizeMobileUserRoutes = types.BoolValue(var8.SummarizeMobileUserRoutes)
			var7 = append(var7, var9)
		}
	}
	var var10 *remoteNetworksRsModelProtocolObject
	if ans.Protocol != nil {
		var10 = &remoteNetworksRsModelProtocolObject{}
		var var11 *remoteNetworksRsModelBgpObject
		if ans.Protocol.Bgp != nil {
			var11 = &remoteNetworksRsModelBgpObject{}
			var11.DoNotExportRoutes = types.BoolValue(ans.Protocol.Bgp.DoNotExportRoutes)
			var11.Enable = types.BoolValue(ans.Protocol.Bgp.Enable)
			var11.LocalIpAddress = types.StringValue(ans.Protocol.Bgp.LocalIpAddress)
			var11.OriginateDefaultRoute = types.BoolValue(ans.Protocol.Bgp.OriginateDefaultRoute)
			var11.PeerAs = types.StringValue(ans.Protocol.Bgp.PeerAs)
			var11.PeerIpAddress = types.StringValue(ans.Protocol.Bgp.PeerIpAddress)
			var11.PeeringType = types.StringValue(ans.Protocol.Bgp.PeeringType)
			var11.Secret = types.StringValue(ans.Protocol.Bgp.Secret)
			var11.SummarizeMobileUserRoutes = types.BoolValue(ans.Protocol.Bgp.SummarizeMobileUserRoutes)
		}
		var var12 *remoteNetworksRsModelBgpPeerObject
		if ans.Protocol.BgpPeer != nil {
			var12 = &remoteNetworksRsModelBgpPeerObject{}
			var12.LocalIpAddress = types.StringValue(ans.Protocol.BgpPeer.LocalIpAddress)
			var12.PeerIpAddress = types.StringValue(ans.Protocol.BgpPeer.PeerIpAddress)
			var12.Secret = types.StringValue(ans.Protocol.BgpPeer.Secret)
		}
		var10.Bgp = var11
		var10.BgpPeer = var12
	}
	state.EcmpLoadBalancing = types.StringValue(ans.EcmpLoadBalancing)
	state.EcmpTunnels = var7
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.IpsecTunnel = types.StringValue(ans.IpsecTunnel)
	state.LicenseType = types.StringValue(ans.LicenseType)
	state.Name = types.StringValue(ans.Name)
	state.Protocol = var10
	state.Region = types.StringValue(ans.Region)
	state.SecondaryIpsecTunnel = types.StringValue(ans.SecondaryIpsecTunnel)
	state.SpnName = types.StringValue(ans.SpnName)
	state.Subnets = EncodeStringSlice(ans.Subnets)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *remoteNetworksResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state remoteNetworksRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_remote_networks",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := xsuBWMo.NewClient(r.client)
	input := xsuBWMo.ReadInput{
		ObjectId: tokens[1],
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
	state.Id = idType
	var var0 []remoteNetworksRsModelEcmpTunnelsObject
	if len(ans.EcmpTunnels) != 0 {
		var0 = make([]remoteNetworksRsModelEcmpTunnelsObject, 0, len(ans.EcmpTunnels))
		for var1Index := range ans.EcmpTunnels {
			var1 := ans.EcmpTunnels[var1Index]
			var var2 remoteNetworksRsModelEcmpTunnelsObject
			var2.DoNotExportRoutes = types.BoolValue(var1.DoNotExportRoutes)
			var2.IpsecTunnel = types.StringValue(var1.IpsecTunnel)
			var2.LocalIpAddress = types.StringValue(var1.LocalIpAddress)
			var2.Name = types.StringValue(var1.Name)
			var2.OriginateDefaultRoute = types.BoolValue(var1.OriginateDefaultRoute)
			var2.PeerAs = types.StringValue(var1.PeerAs)
			var2.PeerIpAddress = types.StringValue(var1.PeerIpAddress)
			var2.PeeringType = types.StringValue(var1.PeeringType)
			var2.Secret = types.StringValue(var1.Secret)
			var2.SummarizeMobileUserRoutes = types.BoolValue(var1.SummarizeMobileUserRoutes)
			var0 = append(var0, var2)
		}
	}
	var var3 *remoteNetworksRsModelProtocolObject
	if ans.Protocol != nil {
		var3 = &remoteNetworksRsModelProtocolObject{}
		var var4 *remoteNetworksRsModelBgpObject
		if ans.Protocol.Bgp != nil {
			var4 = &remoteNetworksRsModelBgpObject{}
			var4.DoNotExportRoutes = types.BoolValue(ans.Protocol.Bgp.DoNotExportRoutes)
			var4.Enable = types.BoolValue(ans.Protocol.Bgp.Enable)
			var4.LocalIpAddress = types.StringValue(ans.Protocol.Bgp.LocalIpAddress)
			var4.OriginateDefaultRoute = types.BoolValue(ans.Protocol.Bgp.OriginateDefaultRoute)
			var4.PeerAs = types.StringValue(ans.Protocol.Bgp.PeerAs)
			var4.PeerIpAddress = types.StringValue(ans.Protocol.Bgp.PeerIpAddress)
			var4.PeeringType = types.StringValue(ans.Protocol.Bgp.PeeringType)
			var4.Secret = types.StringValue(ans.Protocol.Bgp.Secret)
			var4.SummarizeMobileUserRoutes = types.BoolValue(ans.Protocol.Bgp.SummarizeMobileUserRoutes)
		}
		var var5 *remoteNetworksRsModelBgpPeerObject
		if ans.Protocol.BgpPeer != nil {
			var5 = &remoteNetworksRsModelBgpPeerObject{}
			var5.LocalIpAddress = types.StringValue(ans.Protocol.BgpPeer.LocalIpAddress)
			var5.PeerIpAddress = types.StringValue(ans.Protocol.BgpPeer.PeerIpAddress)
			var5.Secret = types.StringValue(ans.Protocol.BgpPeer.Secret)
		}
		var3.Bgp = var4
		var3.BgpPeer = var5
	}
	state.EcmpLoadBalancing = types.StringValue(ans.EcmpLoadBalancing)
	state.EcmpTunnels = var0
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.IpsecTunnel = types.StringValue(ans.IpsecTunnel)
	state.LicenseType = types.StringValue(ans.LicenseType)
	state.Name = types.StringValue(ans.Name)
	state.Protocol = var3
	state.Region = types.StringValue(ans.Region)
	state.SecondaryIpsecTunnel = types.StringValue(ans.SecondaryIpsecTunnel)
	state.SpnName = types.StringValue(ans.SpnName)
	state.Subnets = EncodeStringSlice(ans.Subnets)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *remoteNetworksResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state remoteNetworksRsModel
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
		"resource_name":               "sase_remote_networks",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := xsuBWMo.NewClient(r.client)
	input := xsuBWMo.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 jhtSIUK.Config
	var0.EcmpLoadBalancing = plan.EcmpLoadBalancing.ValueString()
	var var1 []jhtSIUK.EcmpTunnelsObject
	if len(plan.EcmpTunnels) != 0 {
		var1 = make([]jhtSIUK.EcmpTunnelsObject, 0, len(plan.EcmpTunnels))
		for var2Index := range plan.EcmpTunnels {
			var2 := plan.EcmpTunnels[var2Index]
			var var3 jhtSIUK.EcmpTunnelsObject
			var3.DoNotExportRoutes = var2.DoNotExportRoutes.ValueBool()
			var3.IpsecTunnel = var2.IpsecTunnel.ValueString()
			var3.LocalIpAddress = var2.LocalIpAddress.ValueString()
			var3.Name = var2.Name.ValueString()
			var3.OriginateDefaultRoute = var2.OriginateDefaultRoute.ValueBool()
			var3.PeerAs = var2.PeerAs.ValueString()
			var3.PeerIpAddress = var2.PeerIpAddress.ValueString()
			var3.PeeringType = var2.PeeringType.ValueString()
			var3.Secret = var2.Secret.ValueString()
			var3.SummarizeMobileUserRoutes = var2.SummarizeMobileUserRoutes.ValueBool()
			var1 = append(var1, var3)
		}
	}
	var0.EcmpTunnels = var1
	var0.IpsecTunnel = plan.IpsecTunnel.ValueString()
	var0.LicenseType = plan.LicenseType.ValueString()
	var0.Name = plan.Name.ValueString()
	var var4 *jhtSIUK.ProtocolObject
	if plan.Protocol != nil {
		var4 = &jhtSIUK.ProtocolObject{}
		var var5 *jhtSIUK.BgpObject
		if plan.Protocol.Bgp != nil {
			var5 = &jhtSIUK.BgpObject{}
			var5.DoNotExportRoutes = plan.Protocol.Bgp.DoNotExportRoutes.ValueBool()
			var5.Enable = plan.Protocol.Bgp.Enable.ValueBool()
			var5.LocalIpAddress = plan.Protocol.Bgp.LocalIpAddress.ValueString()
			var5.OriginateDefaultRoute = plan.Protocol.Bgp.OriginateDefaultRoute.ValueBool()
			var5.PeerAs = plan.Protocol.Bgp.PeerAs.ValueString()
			var5.PeerIpAddress = plan.Protocol.Bgp.PeerIpAddress.ValueString()
			var5.PeeringType = plan.Protocol.Bgp.PeeringType.ValueString()
			var5.Secret = plan.Protocol.Bgp.Secret.ValueString()
			var5.SummarizeMobileUserRoutes = plan.Protocol.Bgp.SummarizeMobileUserRoutes.ValueBool()
		}
		var4.Bgp = var5
		var var6 *jhtSIUK.BgpPeerObject
		if plan.Protocol.BgpPeer != nil {
			var6 = &jhtSIUK.BgpPeerObject{}
			var6.LocalIpAddress = plan.Protocol.BgpPeer.LocalIpAddress.ValueString()
			var6.PeerIpAddress = plan.Protocol.BgpPeer.PeerIpAddress.ValueString()
			var6.Secret = plan.Protocol.BgpPeer.Secret.ValueString()
		}
		var4.BgpPeer = var6
	}
	var0.Protocol = var4
	var0.Region = plan.Region.ValueString()
	var0.SecondaryIpsecTunnel = plan.SecondaryIpsecTunnel.ValueString()
	var0.SpnName = plan.SpnName.ValueString()
	var0.Subnets = DecodeStringSlice(plan.Subnets)
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var7 []remoteNetworksRsModelEcmpTunnelsObject
	if len(ans.EcmpTunnels) != 0 {
		var7 = make([]remoteNetworksRsModelEcmpTunnelsObject, 0, len(ans.EcmpTunnels))
		for var8Index := range ans.EcmpTunnels {
			var8 := ans.EcmpTunnels[var8Index]
			var var9 remoteNetworksRsModelEcmpTunnelsObject
			var9.DoNotExportRoutes = types.BoolValue(var8.DoNotExportRoutes)
			var9.IpsecTunnel = types.StringValue(var8.IpsecTunnel)
			var9.LocalIpAddress = types.StringValue(var8.LocalIpAddress)
			var9.Name = types.StringValue(var8.Name)
			var9.OriginateDefaultRoute = types.BoolValue(var8.OriginateDefaultRoute)
			var9.PeerAs = types.StringValue(var8.PeerAs)
			var9.PeerIpAddress = types.StringValue(var8.PeerIpAddress)
			var9.PeeringType = types.StringValue(var8.PeeringType)
			var9.Secret = types.StringValue(var8.Secret)
			var9.SummarizeMobileUserRoutes = types.BoolValue(var8.SummarizeMobileUserRoutes)
			var7 = append(var7, var9)
		}
	}
	var var10 *remoteNetworksRsModelProtocolObject
	if ans.Protocol != nil {
		var10 = &remoteNetworksRsModelProtocolObject{}
		var var11 *remoteNetworksRsModelBgpObject
		if ans.Protocol.Bgp != nil {
			var11 = &remoteNetworksRsModelBgpObject{}
			var11.DoNotExportRoutes = types.BoolValue(ans.Protocol.Bgp.DoNotExportRoutes)
			var11.Enable = types.BoolValue(ans.Protocol.Bgp.Enable)
			var11.LocalIpAddress = types.StringValue(ans.Protocol.Bgp.LocalIpAddress)
			var11.OriginateDefaultRoute = types.BoolValue(ans.Protocol.Bgp.OriginateDefaultRoute)
			var11.PeerAs = types.StringValue(ans.Protocol.Bgp.PeerAs)
			var11.PeerIpAddress = types.StringValue(ans.Protocol.Bgp.PeerIpAddress)
			var11.PeeringType = types.StringValue(ans.Protocol.Bgp.PeeringType)
			var11.Secret = types.StringValue(ans.Protocol.Bgp.Secret)
			var11.SummarizeMobileUserRoutes = types.BoolValue(ans.Protocol.Bgp.SummarizeMobileUserRoutes)
		}
		var var12 *remoteNetworksRsModelBgpPeerObject
		if ans.Protocol.BgpPeer != nil {
			var12 = &remoteNetworksRsModelBgpPeerObject{}
			var12.LocalIpAddress = types.StringValue(ans.Protocol.BgpPeer.LocalIpAddress)
			var12.PeerIpAddress = types.StringValue(ans.Protocol.BgpPeer.PeerIpAddress)
			var12.Secret = types.StringValue(ans.Protocol.BgpPeer.Secret)
		}
		var10.Bgp = var11
		var10.BgpPeer = var12
	}
	state.EcmpLoadBalancing = types.StringValue(ans.EcmpLoadBalancing)
	state.EcmpTunnels = var7
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.IpsecTunnel = types.StringValue(ans.IpsecTunnel)
	state.LicenseType = types.StringValue(ans.LicenseType)
	state.Name = types.StringValue(ans.Name)
	state.Protocol = var10
	state.Region = types.StringValue(ans.Region)
	state.SecondaryIpsecTunnel = types.StringValue(ans.SecondaryIpsecTunnel)
	state.SpnName = types.StringValue(ans.SpnName)
	state.Subnets = EncodeStringSlice(ans.Subnets)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *remoteNetworksResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"terraform_provider_function": "Delete",
		"resource_name":               "sase_remote_networks",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	svc := xsuBWMo.NewClient(r.client)
	input := xsuBWMo.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *remoteNetworksResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
