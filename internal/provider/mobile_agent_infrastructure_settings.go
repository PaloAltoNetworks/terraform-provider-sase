package provider

import (
	"context"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	sefSZSA "github.com/paloaltonetworks/sase-go/netsec/service/v1/mobileagent/infrastructuresettings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Data source (listing).
var (
	_ datasource.DataSource              = &mobileAgentInfrastructureSettingsListDataSource{}
	_ datasource.DataSourceWithConfigure = &mobileAgentInfrastructureSettingsListDataSource{}
)

func NewMobileAgentInfrastructureSettingsListDataSource() datasource.DataSource {
	return &mobileAgentInfrastructureSettingsListDataSource{}
}

type mobileAgentInfrastructureSettingsListDataSource struct {
	client *sase.Client
}

type mobileAgentInfrastructureSettingsListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data   []mobileAgentInfrastructureSettingsListDsModelConfig `tfsdk:"data"`
	Limit  types.Int64                                          `tfsdk:"limit"`
	Offset types.Int64                                          `tfsdk:"offset"`
	Total  types.Int64                                          `tfsdk:"total"`
}

type mobileAgentInfrastructureSettingsListDsModelConfig struct {
	DnsServers     []mobileAgentInfrastructureSettingsListDsModelDnsServersObject   `tfsdk:"dns_servers"`
	EnableWins     *mobileAgentInfrastructureSettingsListDsModelEnableWinsObject    `tfsdk:"enable_wins"`
	IpPools        []mobileAgentInfrastructureSettingsListDsModelIpPoolsObject      `tfsdk:"ip_pools"`
	Ipv6           types.Bool                                                       `tfsdk:"ipv6"`
	Name           types.String                                                     `tfsdk:"name"`
	PortalHostname mobileAgentInfrastructureSettingsListDsModelPortalHostnameObject `tfsdk:"portal_hostname"`
	RegionIpv6     mobileAgentInfrastructureSettingsListDsModelRegionIpv6Object     `tfsdk:"region_ipv6"`
	UdpQueries     *mobileAgentInfrastructureSettingsListDsModelUdpQueriesObject    `tfsdk:"udp_queries"`
}

type mobileAgentInfrastructureSettingsListDsModelDnsServersObject struct {
	DnsSuffix          []types.String                                                        `tfsdk:"dns_suffix"`
	InternalDnsMatch   []mobileAgentInfrastructureSettingsListDsModelInternalDnsMatchObject  `tfsdk:"internal_dns_match"`
	Name               types.String                                                          `tfsdk:"name"`
	PrimaryPublicDns   *mobileAgentInfrastructureSettingsListDsModelPrimaryPublicDnsObject   `tfsdk:"primary_public_dns"`
	SecondaryPublicDns *mobileAgentInfrastructureSettingsListDsModelSecondaryPublicDnsObject `tfsdk:"secondary_public_dns"`
}

type mobileAgentInfrastructureSettingsListDsModelInternalDnsMatchObject struct {
	DomainList []types.String                                               `tfsdk:"domain_list"`
	Name       types.String                                                 `tfsdk:"name"`
	Primary    *mobileAgentInfrastructureSettingsListDsModelPrimaryObject   `tfsdk:"primary"`
	Secondary  *mobileAgentInfrastructureSettingsListDsModelSecondaryObject `tfsdk:"secondary"`
}

type mobileAgentInfrastructureSettingsListDsModelPrimaryObject struct {
	DnsServer       types.Bool `tfsdk:"dns_server"`
	UseCloudDefault types.Bool `tfsdk:"use_cloud_default"`
}

type mobileAgentInfrastructureSettingsListDsModelSecondaryObject struct {
	DnsServer       types.Bool `tfsdk:"dns_server"`
	UseCloudDefault types.Bool `tfsdk:"use_cloud_default"`
}

type mobileAgentInfrastructureSettingsListDsModelPrimaryPublicDnsObject struct {
	DnsServer types.String `tfsdk:"dns_server"`
}

type mobileAgentInfrastructureSettingsListDsModelSecondaryPublicDnsObject struct {
	DnsServer types.String `tfsdk:"dns_server"`
}

type mobileAgentInfrastructureSettingsListDsModelEnableWinsObject struct {
	No  types.Bool                                             `tfsdk:"no"`
	Yes *mobileAgentInfrastructureSettingsListDsModelYesObject `tfsdk:"yes"`
}

type mobileAgentInfrastructureSettingsListDsModelYesObject struct {
	WinsServers []mobileAgentInfrastructureSettingsListDsModelWinsServersObject `tfsdk:"wins_servers"`
}

type mobileAgentInfrastructureSettingsListDsModelWinsServersObject struct {
	Name      types.String `tfsdk:"name"`
	Primary   types.String `tfsdk:"primary"`
	Secondary types.String `tfsdk:"secondary"`
}

type mobileAgentInfrastructureSettingsListDsModelIpPoolsObject struct {
	IpPool []types.String `tfsdk:"ip_pool"`
	Name   types.String   `tfsdk:"name"`
}

type mobileAgentInfrastructureSettingsListDsModelPortalHostnameObject struct {
	CustomDomain  *mobileAgentInfrastructureSettingsListDsModelCustomDomainObject  `tfsdk:"custom_domain"`
	DefaultDomain *mobileAgentInfrastructureSettingsListDsModelDefaultDomainObject `tfsdk:"default_domain"`
}

type mobileAgentInfrastructureSettingsListDsModelCustomDomainObject struct {
	Cname                types.String `tfsdk:"cname"`
	Hostname             types.String `tfsdk:"hostname"`
	SslTlsServiceProfile types.String `tfsdk:"ssl_tls_service_profile"`
}

type mobileAgentInfrastructureSettingsListDsModelDefaultDomainObject struct {
	Hostname types.String `tfsdk:"hostname"`
}

type mobileAgentInfrastructureSettingsListDsModelRegionIpv6Object struct {
	Region []mobileAgentInfrastructureSettingsListDsModelRegionObject `tfsdk:"region"`
}

type mobileAgentInfrastructureSettingsListDsModelRegionObject struct {
	Locations []types.String `tfsdk:"locations"`
	Name      types.String   `tfsdk:"name"`
}

type mobileAgentInfrastructureSettingsListDsModelUdpQueriesObject struct {
	Retries *mobileAgentInfrastructureSettingsListDsModelRetriesObject `tfsdk:"retries"`
}

type mobileAgentInfrastructureSettingsListDsModelRetriesObject struct {
	Attempts types.Int64 `tfsdk:"attempts"`
	Interval types.Int64 `tfsdk:"interval"`
}

// Metadata returns the data source type name.
func (d *mobileAgentInfrastructureSettingsListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mobile_agent_infrastructure_settings_list"
}

// Schema defines the schema for this listing data source.
func (d *mobileAgentInfrastructureSettingsListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Retrieves a listing of config items.",

		Attributes: map[string]dsschema.Attribute{
			"id": dsschema.StringAttribute{
				Description: "The object ID.",
				Computed:    true,
			},

			// Input.
			"folder": dsschema.StringAttribute{
				Description: "",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Mobile Users"),
				},
			},

			// Output.
			"data": dsschema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"dns_servers": dsschema.ListNestedAttribute{
							Description: "",
							Computed:    true,
							NestedObject: dsschema.NestedAttributeObject{
								Attributes: map[string]dsschema.Attribute{
									"dns_suffix": dsschema.ListAttribute{
										Description: "",
										Computed:    true,
										ElementType: types.StringType,
									},
									"internal_dns_match": dsschema.ListNestedAttribute{
										Description: "",
										Computed:    true,
										NestedObject: dsschema.NestedAttributeObject{
											Attributes: map[string]dsschema.Attribute{
												"domain_list": dsschema.ListAttribute{
													Description: "",
													Computed:    true,
													ElementType: types.StringType,
												},
												"name": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"primary": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"dns_server": dsschema.BoolAttribute{
															Description: "",
															Computed:    true,
														},
														"use_cloud_default": dsschema.BoolAttribute{
															Description: "",
															Computed:    true,
														},
													},
												},
												"secondary": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"dns_server": dsschema.BoolAttribute{
															Description: "",
															Computed:    true,
														},
														"use_cloud_default": dsschema.BoolAttribute{
															Description: "",
															Computed:    true,
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
									"primary_public_dns": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"dns_server": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
										},
									},
									"secondary_public_dns": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"dns_server": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
										},
									},
								},
							},
						},
						"enable_wins": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"no": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"yes": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"wins_servers": dsschema.ListNestedAttribute{
											Description: "",
											Computed:    true,
											NestedObject: dsschema.NestedAttributeObject{
												Attributes: map[string]dsschema.Attribute{
													"name": dsschema.StringAttribute{
														Description: "",
														Computed:    true,
													},
													"primary": dsschema.StringAttribute{
														Description: "",
														Computed:    true,
													},
													"secondary": dsschema.StringAttribute{
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
						"ip_pools": dsschema.ListNestedAttribute{
							Description: "",
							Computed:    true,
							NestedObject: dsschema.NestedAttributeObject{
								Attributes: map[string]dsschema.Attribute{
									"ip_pool": dsschema.ListAttribute{
										Description: "",
										Computed:    true,
										ElementType: types.StringType,
									},
									"name": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
						},
						"ipv6": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"portal_hostname": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"custom_domain": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"cname": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"hostname": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"ssl_tls_service_profile": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
								"default_domain": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"hostname": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
							},
						},
						"region_ipv6": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"region": dsschema.ListNestedAttribute{
									Description: "",
									Computed:    true,
									NestedObject: dsschema.NestedAttributeObject{
										Attributes: map[string]dsschema.Attribute{
											"locations": dsschema.ListAttribute{
												Description: "",
												Computed:    true,
												ElementType: types.StringType,
											},
											"name": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
										},
									},
								},
							},
						},
						"udp_queries": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"retries": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"attempts": dsschema.Int64Attribute{
											Description: "",
											Computed:    true,
										},
										"interval": dsschema.Int64Attribute{
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
			"limit": dsschema.Int64Attribute{
				Description: "",
				Computed:    true,
			},
			"offset": dsschema.Int64Attribute{
				Description: "",
				Computed:    true,
			},
			"total": dsschema.Int64Attribute{
				Description: "",
				Computed:    true,
			},
		},
	}
}

// Configure prepares the struct.
func (d *mobileAgentInfrastructureSettingsListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *mobileAgentInfrastructureSettingsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state mobileAgentInfrastructureSettingsListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_mobile_agent_infrastructure_settings_list",
		"folder":           state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := sefSZSA.NewClient(d.client)
	input := sefSZSA.ListInput{
		Folder: state.Folder.ValueString(),
	}

	// Perform the operation.
	ans, err := svc.List(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error getting listing", err.Error())
		return
	}

	// Store the answer to state.
	var idBuilder strings.Builder
	idBuilder.WriteString(input.Folder)
	state.Id = types.StringValue(idBuilder.String())
	var var0 []mobileAgentInfrastructureSettingsListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]mobileAgentInfrastructureSettingsListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 mobileAgentInfrastructureSettingsListDsModelConfig
			var var3 []mobileAgentInfrastructureSettingsListDsModelDnsServersObject
			if len(var1.DnsServers) != 0 {
				var3 = make([]mobileAgentInfrastructureSettingsListDsModelDnsServersObject, 0, len(var1.DnsServers))
				for var4Index := range var1.DnsServers {
					var4 := var1.DnsServers[var4Index]
					var var5 mobileAgentInfrastructureSettingsListDsModelDnsServersObject
					var var6 []mobileAgentInfrastructureSettingsListDsModelInternalDnsMatchObject
					if len(var4.InternalDnsMatch) != 0 {
						var6 = make([]mobileAgentInfrastructureSettingsListDsModelInternalDnsMatchObject, 0, len(var4.InternalDnsMatch))
						for var7Index := range var4.InternalDnsMatch {
							var7 := var4.InternalDnsMatch[var7Index]
							var var8 mobileAgentInfrastructureSettingsListDsModelInternalDnsMatchObject
							var var9 *mobileAgentInfrastructureSettingsListDsModelPrimaryObject
							if var7.Primary != nil {
								var9 = &mobileAgentInfrastructureSettingsListDsModelPrimaryObject{}
								if var7.Primary.DnsServer != nil {
									var9.DnsServer = types.BoolValue(true)
								}
								if var7.Primary.UseCloudDefault != nil {
									var9.UseCloudDefault = types.BoolValue(true)
								}
							}
							var var10 *mobileAgentInfrastructureSettingsListDsModelSecondaryObject
							if var7.Secondary != nil {
								var10 = &mobileAgentInfrastructureSettingsListDsModelSecondaryObject{}
								if var7.Secondary.DnsServer != nil {
									var10.DnsServer = types.BoolValue(true)
								}
								if var7.Secondary.UseCloudDefault != nil {
									var10.UseCloudDefault = types.BoolValue(true)
								}
							}
							var8.DomainList = EncodeStringSlice(var7.DomainList)
							var8.Name = types.StringValue(var7.Name)
							var8.Primary = var9
							var8.Secondary = var10
							var6 = append(var6, var8)
						}
					}
					var var11 *mobileAgentInfrastructureSettingsListDsModelPrimaryPublicDnsObject
					if var4.PrimaryPublicDns != nil {
						var11 = &mobileAgentInfrastructureSettingsListDsModelPrimaryPublicDnsObject{}
						var11.DnsServer = types.StringValue(var4.PrimaryPublicDns.DnsServer)
					}
					var var12 *mobileAgentInfrastructureSettingsListDsModelSecondaryPublicDnsObject
					if var4.SecondaryPublicDns != nil {
						var12 = &mobileAgentInfrastructureSettingsListDsModelSecondaryPublicDnsObject{}
						var12.DnsServer = types.StringValue(var4.SecondaryPublicDns.DnsServer)
					}
					var5.DnsSuffix = EncodeStringSlice(var4.DnsSuffix)
					var5.InternalDnsMatch = var6
					var5.Name = types.StringValue(var4.Name)
					var5.PrimaryPublicDns = var11
					var5.SecondaryPublicDns = var12
					var3 = append(var3, var5)
				}
			}
			var var13 *mobileAgentInfrastructureSettingsListDsModelEnableWinsObject
			if var1.EnableWins != nil {
				var13 = &mobileAgentInfrastructureSettingsListDsModelEnableWinsObject{}
				var var14 *mobileAgentInfrastructureSettingsListDsModelYesObject
				if var1.EnableWins.Yes != nil {
					var14 = &mobileAgentInfrastructureSettingsListDsModelYesObject{}
					var var15 []mobileAgentInfrastructureSettingsListDsModelWinsServersObject
					if len(var1.EnableWins.Yes.WinsServers) != 0 {
						var15 = make([]mobileAgentInfrastructureSettingsListDsModelWinsServersObject, 0, len(var1.EnableWins.Yes.WinsServers))
						for var16Index := range var1.EnableWins.Yes.WinsServers {
							var16 := var1.EnableWins.Yes.WinsServers[var16Index]
							var var17 mobileAgentInfrastructureSettingsListDsModelWinsServersObject
							var17.Name = types.StringValue(var16.Name)
							var17.Primary = types.StringValue(var16.Primary)
							var17.Secondary = types.StringValue(var16.Secondary)
							var15 = append(var15, var17)
						}
					}
					var14.WinsServers = var15
				}
				if var1.EnableWins.No != nil {
					var13.No = types.BoolValue(true)
				}
				var13.Yes = var14
			}
			var var18 []mobileAgentInfrastructureSettingsListDsModelIpPoolsObject
			if len(var1.IpPools) != 0 {
				var18 = make([]mobileAgentInfrastructureSettingsListDsModelIpPoolsObject, 0, len(var1.IpPools))
				for var19Index := range var1.IpPools {
					var19 := var1.IpPools[var19Index]
					var var20 mobileAgentInfrastructureSettingsListDsModelIpPoolsObject
					var20.IpPool = EncodeStringSlice(var19.IpPool)
					var20.Name = types.StringValue(var19.Name)
					var18 = append(var18, var20)
				}
			}
			var var21 mobileAgentInfrastructureSettingsListDsModelPortalHostnameObject
			var var22 *mobileAgentInfrastructureSettingsListDsModelCustomDomainObject
			if var1.PortalHostname.CustomDomain != nil {
				var22 = &mobileAgentInfrastructureSettingsListDsModelCustomDomainObject{}
				var22.Cname = types.StringValue(var1.PortalHostname.CustomDomain.Cname)
				var22.Hostname = types.StringValue(var1.PortalHostname.CustomDomain.Hostname)
				var22.SslTlsServiceProfile = types.StringValue(var1.PortalHostname.CustomDomain.SslTlsServiceProfile)
			}
			var var23 *mobileAgentInfrastructureSettingsListDsModelDefaultDomainObject
			if var1.PortalHostname.DefaultDomain != nil {
				var23 = &mobileAgentInfrastructureSettingsListDsModelDefaultDomainObject{}
				var23.Hostname = types.StringValue(var1.PortalHostname.DefaultDomain.Hostname)
			}
			var21.CustomDomain = var22
			var21.DefaultDomain = var23
			var var24 mobileAgentInfrastructureSettingsListDsModelRegionIpv6Object
			var var25 []mobileAgentInfrastructureSettingsListDsModelRegionObject
			if len(var1.RegionIpv6.Region) != 0 {
				var25 = make([]mobileAgentInfrastructureSettingsListDsModelRegionObject, 0, len(var1.RegionIpv6.Region))
				for var26Index := range var1.RegionIpv6.Region {
					var26 := var1.RegionIpv6.Region[var26Index]
					var var27 mobileAgentInfrastructureSettingsListDsModelRegionObject
					var27.Locations = EncodeStringSlice(var26.Locations)
					var27.Name = types.StringValue(var26.Name)
					var25 = append(var25, var27)
				}
			}
			var24.Region = var25
			var var28 *mobileAgentInfrastructureSettingsListDsModelUdpQueriesObject
			if var1.UdpQueries != nil {
				var28 = &mobileAgentInfrastructureSettingsListDsModelUdpQueriesObject{}
				var var29 *mobileAgentInfrastructureSettingsListDsModelRetriesObject
				if var1.UdpQueries.Retries != nil {
					var29 = &mobileAgentInfrastructureSettingsListDsModelRetriesObject{}
					var29.Attempts = types.Int64Value(var1.UdpQueries.Retries.Attempts)
					var29.Interval = types.Int64Value(var1.UdpQueries.Retries.Interval)
				}
				var28.Retries = var29
			}
			var2.DnsServers = var3
			var2.EnableWins = var13
			var2.IpPools = var18
			var2.Ipv6 = types.BoolValue(var1.Ipv6)
			var2.Name = types.StringValue(var1.Name)
			var2.PortalHostname = var21
			var2.RegionIpv6 = var24
			var2.UdpQueries = var28
			var0 = append(var0, var2)
		}
	}
	state.Data = var0
	state.Limit = types.Int64Value(ans.Limit)
	state.Offset = types.Int64Value(ans.Offset)
	state.Total = types.Int64Value(ans.Total)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
