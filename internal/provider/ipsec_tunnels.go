package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	mvZFtQR "github.com/paloaltonetworks/sase-go/netsec/schema/ipsec/tunnels"
	fVAkWHS "github.com/paloaltonetworks/sase-go/netsec/service/v1/ipsectunnels"

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
	_ datasource.DataSource              = &ipsecTunnelsListDataSource{}
	_ datasource.DataSourceWithConfigure = &ipsecTunnelsListDataSource{}
)

func NewIpsecTunnelsListDataSource() datasource.DataSource {
	return &ipsecTunnelsListDataSource{}
}

type ipsecTunnelsListDataSource struct {
	client *sase.Client
}

type ipsecTunnelsListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []ipsecTunnelsListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type ipsecTunnelsListDsModelConfig struct {
	AntiReplay             types.Bool                                  `tfsdk:"anti_replay"`
	AutoKey                ipsecTunnelsListDsModelAutoKeyObject        `tfsdk:"auto_key"`
	CopyTos                types.Bool                                  `tfsdk:"copy_tos"`
	EnableGreEncapsulation types.Bool                                  `tfsdk:"enable_gre_encapsulation"`
	ObjectId               types.String                                `tfsdk:"object_id"`
	Name                   types.String                                `tfsdk:"name"`
	TunnelMonitor          *ipsecTunnelsListDsModelTunnelMonitorObject `tfsdk:"tunnel_monitor"`
}

type ipsecTunnelsListDsModelAutoKeyObject struct {
	IkeGateway         []ipsecTunnelsListDsModelIkeGatewayObject `tfsdk:"ike_gateway"`
	IpsecCryptoProfile types.String                              `tfsdk:"ipsec_crypto_profile"`
	ProxyId            []ipsecTunnelsListDsModelProxyIdObject    `tfsdk:"proxy_id"`
}

type ipsecTunnelsListDsModelIkeGatewayObject struct {
	Name types.String `tfsdk:"name"`
}

type ipsecTunnelsListDsModelProxyIdObject struct {
	Local    types.String                           `tfsdk:"local"`
	Name     types.String                           `tfsdk:"name"`
	Protocol *ipsecTunnelsListDsModelProtocolObject `tfsdk:"protocol"`
	Remote   types.String                           `tfsdk:"remote"`
}

type ipsecTunnelsListDsModelProtocolObject struct {
	Number types.Int64                       `tfsdk:"number"`
	Tcp    *ipsecTunnelsListDsModelTcpObject `tfsdk:"tcp"`
	Udp    *ipsecTunnelsListDsModelUdpObject `tfsdk:"udp"`
}

type ipsecTunnelsListDsModelTcpObject struct {
	LocalPort  types.Int64 `tfsdk:"local_port"`
	RemotePort types.Int64 `tfsdk:"remote_port"`
}

type ipsecTunnelsListDsModelUdpObject struct {
	LocalPort  types.Int64 `tfsdk:"local_port"`
	RemotePort types.Int64 `tfsdk:"remote_port"`
}

type ipsecTunnelsListDsModelTunnelMonitorObject struct {
	DestinationIp types.String `tfsdk:"destination_ip"`
	Enable        types.Bool   `tfsdk:"enable"`
	ProxyId       types.String `tfsdk:"proxy_id"`
}

// Metadata returns the data source type name.
func (d *ipsecTunnelsListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipsec_tunnels_list"
}

// Schema defines the schema for this listing data source.
func (d *ipsecTunnelsListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"anti_replay": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"auto_key": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"ike_gateway": dsschema.ListNestedAttribute{
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
								"ipsec_crypto_profile": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
								"proxy_id": dsschema.ListNestedAttribute{
									Description: "",
									Computed:    true,
									NestedObject: dsschema.NestedAttributeObject{
										Attributes: map[string]dsschema.Attribute{
											"local": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
											"name": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
											"protocol": dsschema.SingleNestedAttribute{
												Description: "",
												Computed:    true,
												Attributes: map[string]dsschema.Attribute{
													"number": dsschema.Int64Attribute{
														Description: "",
														Computed:    true,
													},
													"tcp": dsschema.SingleNestedAttribute{
														Description: "",
														Computed:    true,
														Attributes: map[string]dsschema.Attribute{
															"local_port": dsschema.Int64Attribute{
																Description: "",
																Computed:    true,
															},
															"remote_port": dsschema.Int64Attribute{
																Description: "",
																Computed:    true,
															},
														},
													},
													"udp": dsschema.SingleNestedAttribute{
														Description: "",
														Computed:    true,
														Attributes: map[string]dsschema.Attribute{
															"local_port": dsschema.Int64Attribute{
																Description: "",
																Computed:    true,
															},
															"remote_port": dsschema.Int64Attribute{
																Description: "",
																Computed:    true,
															},
														},
													},
												},
											},
											"remote": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
										},
									},
								},
							},
						},
						"copy_tos": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"enable_gre_encapsulation": dsschema.BoolAttribute{
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
						"tunnel_monitor": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"destination_ip": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
								"enable": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"proxy_id": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
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
func (d *ipsecTunnelsListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *ipsecTunnelsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ipsecTunnelsListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_ipsec_tunnels_list",
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
	svc := fVAkWHS.NewClient(d.client)
	input := fVAkWHS.ListInput{
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
	var var0 []ipsecTunnelsListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]ipsecTunnelsListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 ipsecTunnelsListDsModelConfig
			var var3 ipsecTunnelsListDsModelAutoKeyObject
			var var4 []ipsecTunnelsListDsModelIkeGatewayObject
			if len(var1.AutoKey.IkeGateway) != 0 {
				var4 = make([]ipsecTunnelsListDsModelIkeGatewayObject, 0, len(var1.AutoKey.IkeGateway))
				for var5Index := range var1.AutoKey.IkeGateway {
					var5 := var1.AutoKey.IkeGateway[var5Index]
					var var6 ipsecTunnelsListDsModelIkeGatewayObject
					var6.Name = types.StringValue(var5.Name)
					var4 = append(var4, var6)
				}
			}
			var var7 []ipsecTunnelsListDsModelProxyIdObject
			if len(var1.AutoKey.ProxyId) != 0 {
				var7 = make([]ipsecTunnelsListDsModelProxyIdObject, 0, len(var1.AutoKey.ProxyId))
				for var8Index := range var1.AutoKey.ProxyId {
					var8 := var1.AutoKey.ProxyId[var8Index]
					var var9 ipsecTunnelsListDsModelProxyIdObject
					var var10 *ipsecTunnelsListDsModelProtocolObject
					if var8.Protocol != nil {
						var10 = &ipsecTunnelsListDsModelProtocolObject{}
						var var11 *ipsecTunnelsListDsModelTcpObject
						if var8.Protocol.Tcp != nil {
							var11 = &ipsecTunnelsListDsModelTcpObject{}
							var11.LocalPort = types.Int64Value(var8.Protocol.Tcp.LocalPort)
							var11.RemotePort = types.Int64Value(var8.Protocol.Tcp.RemotePort)
						}
						var var12 *ipsecTunnelsListDsModelUdpObject
						if var8.Protocol.Udp != nil {
							var12 = &ipsecTunnelsListDsModelUdpObject{}
							var12.LocalPort = types.Int64Value(var8.Protocol.Udp.LocalPort)
							var12.RemotePort = types.Int64Value(var8.Protocol.Udp.RemotePort)
						}
						var10.Number = types.Int64Value(var8.Protocol.Number)
						var10.Tcp = var11
						var10.Udp = var12
					}
					var9.Local = types.StringValue(var8.Local)
					var9.Name = types.StringValue(var8.Name)
					var9.Protocol = var10
					var9.Remote = types.StringValue(var8.Remote)
					var7 = append(var7, var9)
				}
			}
			var3.IkeGateway = var4
			var3.IpsecCryptoProfile = types.StringValue(var1.AutoKey.IpsecCryptoProfile)
			var3.ProxyId = var7
			var var13 *ipsecTunnelsListDsModelTunnelMonitorObject
			if var1.TunnelMonitor != nil {
				var13 = &ipsecTunnelsListDsModelTunnelMonitorObject{}
				var13.DestinationIp = types.StringValue(var1.TunnelMonitor.DestinationIp)
				var13.Enable = types.BoolValue(var1.TunnelMonitor.Enable)
				var13.ProxyId = types.StringValue(var1.TunnelMonitor.ProxyId)
			}
			var2.AntiReplay = types.BoolValue(var1.AntiReplay)
			var2.AutoKey = var3
			var2.CopyTos = types.BoolValue(var1.CopyTos)
			var2.EnableGreEncapsulation = types.BoolValue(var1.EnableGreEncapsulation)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Name = types.StringValue(var1.Name)
			var2.TunnelMonitor = var13
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
	_ datasource.DataSource              = &ipsecTunnelsDataSource{}
	_ datasource.DataSourceWithConfigure = &ipsecTunnelsDataSource{}
)

func NewIpsecTunnelsDataSource() datasource.DataSource {
	return &ipsecTunnelsDataSource{}
}

type ipsecTunnelsDataSource struct {
	client *sase.Client
}

type ipsecTunnelsDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`
	Folder   types.String `tfsdk:"folder"`

	// Output.
	// Ref: #/components/schemas/ipsec-tunnels
	AntiReplay             types.Bool                       `tfsdk:"anti_replay"`
	AutoKey                ipsecTunnelsDsModelAutoKeyObject `tfsdk:"auto_key"`
	CopyTos                types.Bool                       `tfsdk:"copy_tos"`
	EnableGreEncapsulation types.Bool                       `tfsdk:"enable_gre_encapsulation"`
	// input omit: ObjectId
	Name          types.String                            `tfsdk:"name"`
	TunnelMonitor *ipsecTunnelsDsModelTunnelMonitorObject `tfsdk:"tunnel_monitor"`
}

type ipsecTunnelsDsModelAutoKeyObject struct {
	IkeGateway         []ipsecTunnelsDsModelIkeGatewayObject `tfsdk:"ike_gateway"`
	IpsecCryptoProfile types.String                          `tfsdk:"ipsec_crypto_profile"`
	ProxyId            []ipsecTunnelsDsModelProxyIdObject    `tfsdk:"proxy_id"`
}

type ipsecTunnelsDsModelIkeGatewayObject struct {
	Name types.String `tfsdk:"name"`
}

type ipsecTunnelsDsModelProxyIdObject struct {
	Local    types.String                       `tfsdk:"local"`
	Name     types.String                       `tfsdk:"name"`
	Protocol *ipsecTunnelsDsModelProtocolObject `tfsdk:"protocol"`
	Remote   types.String                       `tfsdk:"remote"`
}

type ipsecTunnelsDsModelProtocolObject struct {
	Number types.Int64                   `tfsdk:"number"`
	Tcp    *ipsecTunnelsDsModelTcpObject `tfsdk:"tcp"`
	Udp    *ipsecTunnelsDsModelUdpObject `tfsdk:"udp"`
}

type ipsecTunnelsDsModelTcpObject struct {
	LocalPort  types.Int64 `tfsdk:"local_port"`
	RemotePort types.Int64 `tfsdk:"remote_port"`
}

type ipsecTunnelsDsModelUdpObject struct {
	LocalPort  types.Int64 `tfsdk:"local_port"`
	RemotePort types.Int64 `tfsdk:"remote_port"`
}

type ipsecTunnelsDsModelTunnelMonitorObject struct {
	DestinationIp types.String `tfsdk:"destination_ip"`
	Enable        types.Bool   `tfsdk:"enable"`
	ProxyId       types.String `tfsdk:"proxy_id"`
}

// Metadata returns the data source type name.
func (d *ipsecTunnelsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipsec_tunnels"
}

// Schema defines the schema for this listing data source.
func (d *ipsecTunnelsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"folder": dsschema.StringAttribute{
				Description: "The folder of the entry",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Service Connections", "Mobile Users Container", "Mobile Users Explicit Proxy"),
				},
			},

			// Output.
			"anti_replay": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"auto_key": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"ike_gateway": dsschema.ListNestedAttribute{
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
					"ipsec_crypto_profile": dsschema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"proxy_id": dsschema.ListNestedAttribute{
						Description: "",
						Computed:    true,
						NestedObject: dsschema.NestedAttributeObject{
							Attributes: map[string]dsschema.Attribute{
								"local": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
								"name": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
								"protocol": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"number": dsschema.Int64Attribute{
											Description: "",
											Computed:    true,
										},
										"tcp": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"local_port": dsschema.Int64Attribute{
													Description: "",
													Computed:    true,
												},
												"remote_port": dsschema.Int64Attribute{
													Description: "",
													Computed:    true,
												},
											},
										},
										"udp": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"local_port": dsschema.Int64Attribute{
													Description: "",
													Computed:    true,
												},
												"remote_port": dsschema.Int64Attribute{
													Description: "",
													Computed:    true,
												},
											},
										},
									},
								},
								"remote": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
							},
						},
					},
				},
			},
			"copy_tos": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"enable_gre_encapsulation": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"tunnel_monitor": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"destination_ip": dsschema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"enable": dsschema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"proxy_id": dsschema.StringAttribute{
						Description: "",
						Computed:    true,
					},
				},
			},
		},
	}
}

// Configure prepares the struct.
func (d *ipsecTunnelsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *ipsecTunnelsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ipsecTunnelsDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_ipsec_tunnels",
		"object_id":                   state.ObjectId.ValueString(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := fVAkWHS.NewClient(d.client)
	input := fVAkWHS.ReadInput{
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
	var var0 ipsecTunnelsDsModelAutoKeyObject
	var var1 []ipsecTunnelsDsModelIkeGatewayObject
	if len(ans.AutoKey.IkeGateway) != 0 {
		var1 = make([]ipsecTunnelsDsModelIkeGatewayObject, 0, len(ans.AutoKey.IkeGateway))
		for var2Index := range ans.AutoKey.IkeGateway {
			var2 := ans.AutoKey.IkeGateway[var2Index]
			var var3 ipsecTunnelsDsModelIkeGatewayObject
			var3.Name = types.StringValue(var2.Name)
			var1 = append(var1, var3)
		}
	}
	var var4 []ipsecTunnelsDsModelProxyIdObject
	if len(ans.AutoKey.ProxyId) != 0 {
		var4 = make([]ipsecTunnelsDsModelProxyIdObject, 0, len(ans.AutoKey.ProxyId))
		for var5Index := range ans.AutoKey.ProxyId {
			var5 := ans.AutoKey.ProxyId[var5Index]
			var var6 ipsecTunnelsDsModelProxyIdObject
			var var7 *ipsecTunnelsDsModelProtocolObject
			if var5.Protocol != nil {
				var7 = &ipsecTunnelsDsModelProtocolObject{}
				var var8 *ipsecTunnelsDsModelTcpObject
				if var5.Protocol.Tcp != nil {
					var8 = &ipsecTunnelsDsModelTcpObject{}
					var8.LocalPort = types.Int64Value(var5.Protocol.Tcp.LocalPort)
					var8.RemotePort = types.Int64Value(var5.Protocol.Tcp.RemotePort)
				}
				var var9 *ipsecTunnelsDsModelUdpObject
				if var5.Protocol.Udp != nil {
					var9 = &ipsecTunnelsDsModelUdpObject{}
					var9.LocalPort = types.Int64Value(var5.Protocol.Udp.LocalPort)
					var9.RemotePort = types.Int64Value(var5.Protocol.Udp.RemotePort)
				}
				var7.Number = types.Int64Value(var5.Protocol.Number)
				var7.Tcp = var8
				var7.Udp = var9
			}
			var6.Local = types.StringValue(var5.Local)
			var6.Name = types.StringValue(var5.Name)
			var6.Protocol = var7
			var6.Remote = types.StringValue(var5.Remote)
			var4 = append(var4, var6)
		}
	}
	var0.IkeGateway = var1
	var0.IpsecCryptoProfile = types.StringValue(ans.AutoKey.IpsecCryptoProfile)
	var0.ProxyId = var4
	var var10 *ipsecTunnelsDsModelTunnelMonitorObject
	if ans.TunnelMonitor != nil {
		var10 = &ipsecTunnelsDsModelTunnelMonitorObject{}
		var10.DestinationIp = types.StringValue(ans.TunnelMonitor.DestinationIp)
		var10.Enable = types.BoolValue(ans.TunnelMonitor.Enable)
		var10.ProxyId = types.StringValue(ans.TunnelMonitor.ProxyId)
	}
	state.AntiReplay = types.BoolValue(ans.AntiReplay)
	state.AutoKey = var0
	state.CopyTos = types.BoolValue(ans.CopyTos)
	state.EnableGreEncapsulation = types.BoolValue(ans.EnableGreEncapsulation)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.TunnelMonitor = var10

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &ipsecTunnelsResource{}
	_ resource.ResourceWithConfigure   = &ipsecTunnelsResource{}
	_ resource.ResourceWithImportState = &ipsecTunnelsResource{}
)

func NewIpsecTunnelsResource() resource.Resource {
	return &ipsecTunnelsResource{}
}

type ipsecTunnelsResource struct {
	client *sase.Client
}

type ipsecTunnelsRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/ipsec-tunnels
	AntiReplay             types.Bool                              `tfsdk:"anti_replay"`
	AutoKey                ipsecTunnelsRsModelAutoKeyObject        `tfsdk:"auto_key"`
	CopyTos                types.Bool                              `tfsdk:"copy_tos"`
	EnableGreEncapsulation types.Bool                              `tfsdk:"enable_gre_encapsulation"`
	ObjectId               types.String                            `tfsdk:"object_id"`
	Name                   types.String                            `tfsdk:"name"`
	TunnelMonitor          *ipsecTunnelsRsModelTunnelMonitorObject `tfsdk:"tunnel_monitor"`
}

type ipsecTunnelsRsModelAutoKeyObject struct {
	IkeGateway         []ipsecTunnelsRsModelIkeGatewayObject `tfsdk:"ike_gateway"`
	IpsecCryptoProfile types.String                          `tfsdk:"ipsec_crypto_profile"`
	ProxyId            []ipsecTunnelsRsModelProxyIdObject    `tfsdk:"proxy_id"`
}

type ipsecTunnelsRsModelIkeGatewayObject struct {
	Name types.String `tfsdk:"name"`
}

type ipsecTunnelsRsModelProxyIdObject struct {
	Local    types.String                       `tfsdk:"local"`
	Name     types.String                       `tfsdk:"name"`
	Protocol *ipsecTunnelsRsModelProtocolObject `tfsdk:"protocol"`
	Remote   types.String                       `tfsdk:"remote"`
}

type ipsecTunnelsRsModelProtocolObject struct {
	Number types.Int64                   `tfsdk:"number"`
	Tcp    *ipsecTunnelsRsModelTcpObject `tfsdk:"tcp"`
	Udp    *ipsecTunnelsRsModelUdpObject `tfsdk:"udp"`
}

type ipsecTunnelsRsModelTcpObject struct {
	LocalPort  types.Int64 `tfsdk:"local_port"`
	RemotePort types.Int64 `tfsdk:"remote_port"`
}

type ipsecTunnelsRsModelUdpObject struct {
	LocalPort  types.Int64 `tfsdk:"local_port"`
	RemotePort types.Int64 `tfsdk:"remote_port"`
}

type ipsecTunnelsRsModelTunnelMonitorObject struct {
	DestinationIp types.String `tfsdk:"destination_ip"`
	Enable        types.Bool   `tfsdk:"enable"`
	ProxyId       types.String `tfsdk:"proxy_id"`
}

// Metadata returns the data source type name.
func (r *ipsecTunnelsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipsec_tunnels"
}

// Schema defines the schema for this listing data source.
func (r *ipsecTunnelsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

			"anti_replay": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"auto_key": rsschema.SingleNestedAttribute{
				Description: "",
				Required:    true,
				Attributes: map[string]rsschema.Attribute{
					"ike_gateway": rsschema.ListNestedAttribute{
						Description: "",
						Required:    true,
						NestedObject: rsschema.NestedAttributeObject{
							Attributes: map[string]rsschema.Attribute{
								"name": rsschema.StringAttribute{
									Description: "",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.String{
										DefaultString(""),
									},
								},
							},
						},
					},
					"ipsec_crypto_profile": rsschema.StringAttribute{
						Description: "",
						Required:    true,
					},
					"proxy_id": rsschema.ListNestedAttribute{
						Description: "",
						Optional:    true,
						NestedObject: rsschema.NestedAttributeObject{
							Attributes: map[string]rsschema.Attribute{
								"local": rsschema.StringAttribute{
									Description: "",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.String{
										DefaultString(""),
									},
								},
								"name": rsschema.StringAttribute{
									Description: "",
									Required:    true,
								},
								"protocol": rsschema.SingleNestedAttribute{
									Description: "",
									Optional:    true,
									Attributes: map[string]rsschema.Attribute{
										"number": rsschema.Int64Attribute{
											Description: "",
											Optional:    true,
											Computed:    true,
											PlanModifiers: []planmodifier.Int64{
												DefaultInt64(0),
											},
											Validators: []validator.Int64{
												int64validator.Between(1, 254),
											},
										},
										"tcp": rsschema.SingleNestedAttribute{
											Description: "",
											Optional:    true,
											Attributes: map[string]rsschema.Attribute{
												"local_port": rsschema.Int64Attribute{
													Description: "",
													Optional:    true,
													Computed:    true,
													PlanModifiers: []planmodifier.Int64{
														DefaultInt64(0),
													},
													Validators: []validator.Int64{
														int64validator.Between(0, 65535),
													},
												},
												"remote_port": rsschema.Int64Attribute{
													Description: "",
													Optional:    true,
													Computed:    true,
													PlanModifiers: []planmodifier.Int64{
														DefaultInt64(0),
													},
													Validators: []validator.Int64{
														int64validator.Between(0, 65535),
													},
												},
											},
										},
										"udp": rsschema.SingleNestedAttribute{
											Description: "",
											Optional:    true,
											Attributes: map[string]rsschema.Attribute{
												"local_port": rsschema.Int64Attribute{
													Description: "",
													Optional:    true,
													Computed:    true,
													PlanModifiers: []planmodifier.Int64{
														DefaultInt64(0),
													},
													Validators: []validator.Int64{
														int64validator.Between(0, 65535),
													},
												},
												"remote_port": rsschema.Int64Attribute{
													Description: "",
													Optional:    true,
													Computed:    true,
													PlanModifiers: []planmodifier.Int64{
														DefaultInt64(0),
													},
													Validators: []validator.Int64{
														int64validator.Between(0, 65535),
													},
												},
											},
										},
									},
								},
								"remote": rsschema.StringAttribute{
									Description: "",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.String{
										DefaultString(""),
									},
								},
							},
						},
					},
				},
			},
			"copy_tos": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"enable_gre_encapsulation": rsschema.BoolAttribute{
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
				Validators: []validator.String{
					stringvalidator.LengthAtMost(63),
				},
			},
			"tunnel_monitor": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"destination_ip": rsschema.StringAttribute{
						Description: "",
						Required:    true,
					},
					"enable": rsschema.BoolAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(true),
						},
					},
					"proxy_id": rsschema.StringAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							DefaultString(""),
						},
					},
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *ipsecTunnelsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *ipsecTunnelsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state ipsecTunnelsRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_ipsec_tunnels",
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := fVAkWHS.NewClient(r.client)
	input := fVAkWHS.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 mvZFtQR.Config
	var0.AntiReplay = state.AntiReplay.ValueBool()
	var var1 mvZFtQR.AutoKeyObject
	var var2 []mvZFtQR.IkeGatewayObject
	if len(state.AutoKey.IkeGateway) != 0 {
		var2 = make([]mvZFtQR.IkeGatewayObject, 0, len(state.AutoKey.IkeGateway))
		for var3Index := range state.AutoKey.IkeGateway {
			var3 := state.AutoKey.IkeGateway[var3Index]
			var var4 mvZFtQR.IkeGatewayObject
			var4.Name = var3.Name.ValueString()
			var2 = append(var2, var4)
		}
	}
	var1.IkeGateway = var2
	var1.IpsecCryptoProfile = state.AutoKey.IpsecCryptoProfile.ValueString()
	var var5 []mvZFtQR.ProxyIdObject
	if len(state.AutoKey.ProxyId) != 0 {
		var5 = make([]mvZFtQR.ProxyIdObject, 0, len(state.AutoKey.ProxyId))
		for var6Index := range state.AutoKey.ProxyId {
			var6 := state.AutoKey.ProxyId[var6Index]
			var var7 mvZFtQR.ProxyIdObject
			var7.Local = var6.Local.ValueString()
			var7.Name = var6.Name.ValueString()
			var var8 *mvZFtQR.ProtocolObject
			if var6.Protocol != nil {
				var8 = &mvZFtQR.ProtocolObject{}
				var8.Number = var6.Protocol.Number.ValueInt64()
				var var9 *mvZFtQR.TcpObject
				if var6.Protocol.Tcp != nil {
					var9 = &mvZFtQR.TcpObject{}
					var9.LocalPort = var6.Protocol.Tcp.LocalPort.ValueInt64()
					var9.RemotePort = var6.Protocol.Tcp.RemotePort.ValueInt64()
				}
				var8.Tcp = var9
				var var10 *mvZFtQR.UdpObject
				if var6.Protocol.Udp != nil {
					var10 = &mvZFtQR.UdpObject{}
					var10.LocalPort = var6.Protocol.Udp.LocalPort.ValueInt64()
					var10.RemotePort = var6.Protocol.Udp.RemotePort.ValueInt64()
				}
				var8.Udp = var10
			}
			var7.Protocol = var8
			var7.Remote = var6.Remote.ValueString()
			var5 = append(var5, var7)
		}
	}
	var1.ProxyId = var5
	var0.AutoKey = var1
	var0.CopyTos = state.CopyTos.ValueBool()
	var0.EnableGreEncapsulation = state.EnableGreEncapsulation.ValueBool()
	var0.Name = state.Name.ValueString()
	var var11 *mvZFtQR.TunnelMonitorObject
	if state.TunnelMonitor != nil {
		var11 = &mvZFtQR.TunnelMonitorObject{}
		var11.DestinationIp = state.TunnelMonitor.DestinationIp.ValueString()
		var11.Enable = state.TunnelMonitor.Enable.ValueBool()
		var11.ProxyId = state.TunnelMonitor.ProxyId.ValueString()
	}
	var0.TunnelMonitor = var11
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
	var var12 ipsecTunnelsRsModelAutoKeyObject
	var var13 []ipsecTunnelsRsModelIkeGatewayObject
	if len(ans.AutoKey.IkeGateway) != 0 {
		var13 = make([]ipsecTunnelsRsModelIkeGatewayObject, 0, len(ans.AutoKey.IkeGateway))
		for var14Index := range ans.AutoKey.IkeGateway {
			var14 := ans.AutoKey.IkeGateway[var14Index]
			var var15 ipsecTunnelsRsModelIkeGatewayObject
			var15.Name = types.StringValue(var14.Name)
			var13 = append(var13, var15)
		}
	}
	var var16 []ipsecTunnelsRsModelProxyIdObject
	if len(ans.AutoKey.ProxyId) != 0 {
		var16 = make([]ipsecTunnelsRsModelProxyIdObject, 0, len(ans.AutoKey.ProxyId))
		for var17Index := range ans.AutoKey.ProxyId {
			var17 := ans.AutoKey.ProxyId[var17Index]
			var var18 ipsecTunnelsRsModelProxyIdObject
			var var19 *ipsecTunnelsRsModelProtocolObject
			if var17.Protocol != nil {
				var19 = &ipsecTunnelsRsModelProtocolObject{}
				var var20 *ipsecTunnelsRsModelTcpObject
				if var17.Protocol.Tcp != nil {
					var20 = &ipsecTunnelsRsModelTcpObject{}
					var20.LocalPort = types.Int64Value(var17.Protocol.Tcp.LocalPort)
					var20.RemotePort = types.Int64Value(var17.Protocol.Tcp.RemotePort)
				}
				var var21 *ipsecTunnelsRsModelUdpObject
				if var17.Protocol.Udp != nil {
					var21 = &ipsecTunnelsRsModelUdpObject{}
					var21.LocalPort = types.Int64Value(var17.Protocol.Udp.LocalPort)
					var21.RemotePort = types.Int64Value(var17.Protocol.Udp.RemotePort)
				}
				var19.Number = types.Int64Value(var17.Protocol.Number)
				var19.Tcp = var20
				var19.Udp = var21
			}
			var18.Local = types.StringValue(var17.Local)
			var18.Name = types.StringValue(var17.Name)
			var18.Protocol = var19
			var18.Remote = types.StringValue(var17.Remote)
			var16 = append(var16, var18)
		}
	}
	var12.IkeGateway = var13
	var12.IpsecCryptoProfile = types.StringValue(ans.AutoKey.IpsecCryptoProfile)
	var12.ProxyId = var16
	var var22 *ipsecTunnelsRsModelTunnelMonitorObject
	if ans.TunnelMonitor != nil {
		var22 = &ipsecTunnelsRsModelTunnelMonitorObject{}
		var22.DestinationIp = types.StringValue(ans.TunnelMonitor.DestinationIp)
		var22.Enable = types.BoolValue(ans.TunnelMonitor.Enable)
		var22.ProxyId = types.StringValue(ans.TunnelMonitor.ProxyId)
	}
	state.AntiReplay = types.BoolValue(ans.AntiReplay)
	state.AutoKey = var12
	state.CopyTos = types.BoolValue(ans.CopyTos)
	state.EnableGreEncapsulation = types.BoolValue(ans.EnableGreEncapsulation)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.TunnelMonitor = var22

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *ipsecTunnelsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state ipsecTunnelsRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_ipsec_tunnels",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := fVAkWHS.NewClient(r.client)
	input := fVAkWHS.ReadInput{
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
	var var0 ipsecTunnelsRsModelAutoKeyObject
	var var1 []ipsecTunnelsRsModelIkeGatewayObject
	if len(ans.AutoKey.IkeGateway) != 0 {
		var1 = make([]ipsecTunnelsRsModelIkeGatewayObject, 0, len(ans.AutoKey.IkeGateway))
		for var2Index := range ans.AutoKey.IkeGateway {
			var2 := ans.AutoKey.IkeGateway[var2Index]
			var var3 ipsecTunnelsRsModelIkeGatewayObject
			var3.Name = types.StringValue(var2.Name)
			var1 = append(var1, var3)
		}
	}
	var var4 []ipsecTunnelsRsModelProxyIdObject
	if len(ans.AutoKey.ProxyId) != 0 {
		var4 = make([]ipsecTunnelsRsModelProxyIdObject, 0, len(ans.AutoKey.ProxyId))
		for var5Index := range ans.AutoKey.ProxyId {
			var5 := ans.AutoKey.ProxyId[var5Index]
			var var6 ipsecTunnelsRsModelProxyIdObject
			var var7 *ipsecTunnelsRsModelProtocolObject
			if var5.Protocol != nil {
				var7 = &ipsecTunnelsRsModelProtocolObject{}
				var var8 *ipsecTunnelsRsModelTcpObject
				if var5.Protocol.Tcp != nil {
					var8 = &ipsecTunnelsRsModelTcpObject{}
					var8.LocalPort = types.Int64Value(var5.Protocol.Tcp.LocalPort)
					var8.RemotePort = types.Int64Value(var5.Protocol.Tcp.RemotePort)
				}
				var var9 *ipsecTunnelsRsModelUdpObject
				if var5.Protocol.Udp != nil {
					var9 = &ipsecTunnelsRsModelUdpObject{}
					var9.LocalPort = types.Int64Value(var5.Protocol.Udp.LocalPort)
					var9.RemotePort = types.Int64Value(var5.Protocol.Udp.RemotePort)
				}
				var7.Number = types.Int64Value(var5.Protocol.Number)
				var7.Tcp = var8
				var7.Udp = var9
			}
			var6.Local = types.StringValue(var5.Local)
			var6.Name = types.StringValue(var5.Name)
			var6.Protocol = var7
			var6.Remote = types.StringValue(var5.Remote)
			var4 = append(var4, var6)
		}
	}
	var0.IkeGateway = var1
	var0.IpsecCryptoProfile = types.StringValue(ans.AutoKey.IpsecCryptoProfile)
	var0.ProxyId = var4
	var var10 *ipsecTunnelsRsModelTunnelMonitorObject
	if ans.TunnelMonitor != nil {
		var10 = &ipsecTunnelsRsModelTunnelMonitorObject{}
		var10.DestinationIp = types.StringValue(ans.TunnelMonitor.DestinationIp)
		var10.Enable = types.BoolValue(ans.TunnelMonitor.Enable)
		var10.ProxyId = types.StringValue(ans.TunnelMonitor.ProxyId)
	}
	state.AntiReplay = types.BoolValue(ans.AntiReplay)
	state.AutoKey = var0
	state.CopyTos = types.BoolValue(ans.CopyTos)
	state.EnableGreEncapsulation = types.BoolValue(ans.EnableGreEncapsulation)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.TunnelMonitor = var10

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *ipsecTunnelsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ipsecTunnelsRsModel
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
		"resource_name":               "sase_ipsec_tunnels",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := fVAkWHS.NewClient(r.client)
	input := fVAkWHS.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 mvZFtQR.Config
	var0.AntiReplay = plan.AntiReplay.ValueBool()
	var var1 mvZFtQR.AutoKeyObject
	var var2 []mvZFtQR.IkeGatewayObject
	if len(plan.AutoKey.IkeGateway) != 0 {
		var2 = make([]mvZFtQR.IkeGatewayObject, 0, len(plan.AutoKey.IkeGateway))
		for var3Index := range plan.AutoKey.IkeGateway {
			var3 := plan.AutoKey.IkeGateway[var3Index]
			var var4 mvZFtQR.IkeGatewayObject
			var4.Name = var3.Name.ValueString()
			var2 = append(var2, var4)
		}
	}
	var1.IkeGateway = var2
	var1.IpsecCryptoProfile = plan.AutoKey.IpsecCryptoProfile.ValueString()
	var var5 []mvZFtQR.ProxyIdObject
	if len(plan.AutoKey.ProxyId) != 0 {
		var5 = make([]mvZFtQR.ProxyIdObject, 0, len(plan.AutoKey.ProxyId))
		for var6Index := range plan.AutoKey.ProxyId {
			var6 := plan.AutoKey.ProxyId[var6Index]
			var var7 mvZFtQR.ProxyIdObject
			var7.Local = var6.Local.ValueString()
			var7.Name = var6.Name.ValueString()
			var var8 *mvZFtQR.ProtocolObject
			if var6.Protocol != nil {
				var8 = &mvZFtQR.ProtocolObject{}
				var8.Number = var6.Protocol.Number.ValueInt64()
				var var9 *mvZFtQR.TcpObject
				if var6.Protocol.Tcp != nil {
					var9 = &mvZFtQR.TcpObject{}
					var9.LocalPort = var6.Protocol.Tcp.LocalPort.ValueInt64()
					var9.RemotePort = var6.Protocol.Tcp.RemotePort.ValueInt64()
				}
				var8.Tcp = var9
				var var10 *mvZFtQR.UdpObject
				if var6.Protocol.Udp != nil {
					var10 = &mvZFtQR.UdpObject{}
					var10.LocalPort = var6.Protocol.Udp.LocalPort.ValueInt64()
					var10.RemotePort = var6.Protocol.Udp.RemotePort.ValueInt64()
				}
				var8.Udp = var10
			}
			var7.Protocol = var8
			var7.Remote = var6.Remote.ValueString()
			var5 = append(var5, var7)
		}
	}
	var1.ProxyId = var5
	var0.AutoKey = var1
	var0.CopyTos = plan.CopyTos.ValueBool()
	var0.EnableGreEncapsulation = plan.EnableGreEncapsulation.ValueBool()
	var0.Name = plan.Name.ValueString()
	var var11 *mvZFtQR.TunnelMonitorObject
	if plan.TunnelMonitor != nil {
		var11 = &mvZFtQR.TunnelMonitorObject{}
		var11.DestinationIp = plan.TunnelMonitor.DestinationIp.ValueString()
		var11.Enable = plan.TunnelMonitor.Enable.ValueBool()
		var11.ProxyId = plan.TunnelMonitor.ProxyId.ValueString()
	}
	var0.TunnelMonitor = var11
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var12 ipsecTunnelsRsModelAutoKeyObject
	var var13 []ipsecTunnelsRsModelIkeGatewayObject
	if len(ans.AutoKey.IkeGateway) != 0 {
		var13 = make([]ipsecTunnelsRsModelIkeGatewayObject, 0, len(ans.AutoKey.IkeGateway))
		for var14Index := range ans.AutoKey.IkeGateway {
			var14 := ans.AutoKey.IkeGateway[var14Index]
			var var15 ipsecTunnelsRsModelIkeGatewayObject
			var15.Name = types.StringValue(var14.Name)
			var13 = append(var13, var15)
		}
	}
	var var16 []ipsecTunnelsRsModelProxyIdObject
	if len(ans.AutoKey.ProxyId) != 0 {
		var16 = make([]ipsecTunnelsRsModelProxyIdObject, 0, len(ans.AutoKey.ProxyId))
		for var17Index := range ans.AutoKey.ProxyId {
			var17 := ans.AutoKey.ProxyId[var17Index]
			var var18 ipsecTunnelsRsModelProxyIdObject
			var var19 *ipsecTunnelsRsModelProtocolObject
			if var17.Protocol != nil {
				var19 = &ipsecTunnelsRsModelProtocolObject{}
				var var20 *ipsecTunnelsRsModelTcpObject
				if var17.Protocol.Tcp != nil {
					var20 = &ipsecTunnelsRsModelTcpObject{}
					var20.LocalPort = types.Int64Value(var17.Protocol.Tcp.LocalPort)
					var20.RemotePort = types.Int64Value(var17.Protocol.Tcp.RemotePort)
				}
				var var21 *ipsecTunnelsRsModelUdpObject
				if var17.Protocol.Udp != nil {
					var21 = &ipsecTunnelsRsModelUdpObject{}
					var21.LocalPort = types.Int64Value(var17.Protocol.Udp.LocalPort)
					var21.RemotePort = types.Int64Value(var17.Protocol.Udp.RemotePort)
				}
				var19.Number = types.Int64Value(var17.Protocol.Number)
				var19.Tcp = var20
				var19.Udp = var21
			}
			var18.Local = types.StringValue(var17.Local)
			var18.Name = types.StringValue(var17.Name)
			var18.Protocol = var19
			var18.Remote = types.StringValue(var17.Remote)
			var16 = append(var16, var18)
		}
	}
	var12.IkeGateway = var13
	var12.IpsecCryptoProfile = types.StringValue(ans.AutoKey.IpsecCryptoProfile)
	var12.ProxyId = var16
	var var22 *ipsecTunnelsRsModelTunnelMonitorObject
	if ans.TunnelMonitor != nil {
		var22 = &ipsecTunnelsRsModelTunnelMonitorObject{}
		var22.DestinationIp = types.StringValue(ans.TunnelMonitor.DestinationIp)
		var22.Enable = types.BoolValue(ans.TunnelMonitor.Enable)
		var22.ProxyId = types.StringValue(ans.TunnelMonitor.ProxyId)
	}
	state.AntiReplay = types.BoolValue(ans.AntiReplay)
	state.AutoKey = var12
	state.CopyTos = types.BoolValue(ans.CopyTos)
	state.EnableGreEncapsulation = types.BoolValue(ans.EnableGreEncapsulation)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.TunnelMonitor = var22

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *ipsecTunnelsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name":               "sase_ipsec_tunnels",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	svc := fVAkWHS.NewClient(r.client)
	input := fVAkWHS.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *ipsecTunnelsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
