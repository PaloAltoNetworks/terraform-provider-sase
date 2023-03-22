package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	fcnKgqA "github.com/paloaltonetworks/sase-go/netsec/schema/dns/security/profiles"
	uSsfsLd "github.com/paloaltonetworks/sase-go/netsec/service/v1/dnssecurityprofiles"

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
	_ datasource.DataSource              = &dnsSecurityProfilesListDataSource{}
	_ datasource.DataSourceWithConfigure = &dnsSecurityProfilesListDataSource{}
)

func NewDnsSecurityProfilesListDataSource() datasource.DataSource {
	return &dnsSecurityProfilesListDataSource{}
}

type dnsSecurityProfilesListDataSource struct {
	client *sase.Client
}

type dnsSecurityProfilesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []dnsSecurityProfilesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type dnsSecurityProfilesListDsModelConfig struct {
	BotnetDomains *dnsSecurityProfilesListDsModelBotnetDomainsObject `tfsdk:"botnet_domains"`
	Description   types.String                                       `tfsdk:"description"`
	ObjectId      types.String                                       `tfsdk:"object_id"`
	Name          types.String                                       `tfsdk:"name"`
}

type dnsSecurityProfilesListDsModelBotnetDomainsObject struct {
	DnsSecurityCategories []dnsSecurityProfilesListDsModelDnsSecurityCategoriesObject `tfsdk:"dns_security_categories"`
	Lists                 []dnsSecurityProfilesListDsModelListsObject                 `tfsdk:"lists"`
	Sinkhole              *dnsSecurityProfilesListDsModelSinkholeObject               `tfsdk:"sinkhole"`
	Whitelist             []dnsSecurityProfilesListDsModelWhitelistObject             `tfsdk:"whitelist"`
}

type dnsSecurityProfilesListDsModelDnsSecurityCategoriesObject struct {
	Action        types.String `tfsdk:"action"`
	LogLevel      types.String `tfsdk:"log_level"`
	Name          types.String `tfsdk:"name"`
	PacketCapture types.String `tfsdk:"packet_capture"`
}

type dnsSecurityProfilesListDsModelListsObject struct {
	Action        *dnsSecurityProfilesListDsModelActionObject `tfsdk:"action"`
	Name          types.String                                `tfsdk:"name"`
	PacketCapture types.String                                `tfsdk:"packet_capture"`
}

type dnsSecurityProfilesListDsModelActionObject struct {
	Alert    types.Bool `tfsdk:"alert"`
	Allow    types.Bool `tfsdk:"allow"`
	Block    types.Bool `tfsdk:"block"`
	Sinkhole types.Bool `tfsdk:"sinkhole"`
}

type dnsSecurityProfilesListDsModelSinkholeObject struct {
	Ipv4Address types.String `tfsdk:"ipv4_address"`
	Ipv6Address types.String `tfsdk:"ipv6_address"`
}

type dnsSecurityProfilesListDsModelWhitelistObject struct {
	Description types.String `tfsdk:"description"`
	Name        types.String `tfsdk:"name"`
}

// Metadata returns the data source type name.
func (d *dnsSecurityProfilesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_security_profiles_list"
}

// Schema defines the schema for this listing data source.
func (d *dnsSecurityProfilesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"botnet_domains": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"dns_security_categories": dsschema.ListNestedAttribute{
									Description: "",
									Computed:    true,
									NestedObject: dsschema.NestedAttributeObject{
										Attributes: map[string]dsschema.Attribute{
											"action": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
											"log_level": dsschema.StringAttribute{
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
										},
									},
								},
								"lists": dsschema.ListNestedAttribute{
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
													"block": dsschema.BoolAttribute{
														Description: "",
														Computed:    true,
													},
													"sinkhole": dsschema.BoolAttribute{
														Description: "",
														Computed:    true,
													},
												},
											},
											"name": dsschema.StringAttribute{
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
								"sinkhole": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"ipv4_address": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"ipv6_address": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
								"whitelist": dsschema.ListNestedAttribute{
									Description: "",
									Computed:    true,
									NestedObject: dsschema.NestedAttributeObject{
										Attributes: map[string]dsschema.Attribute{
											"description": dsschema.StringAttribute{
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
							},
						},
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
func (d *dnsSecurityProfilesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *dnsSecurityProfilesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state dnsSecurityProfilesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_dns_security_profiles_list",
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
	svc := uSsfsLd.NewClient(d.client)
	input := uSsfsLd.ListInput{
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
	var var0 []dnsSecurityProfilesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]dnsSecurityProfilesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 dnsSecurityProfilesListDsModelConfig
			var var3 *dnsSecurityProfilesListDsModelBotnetDomainsObject
			if var1.BotnetDomains != nil {
				var3 = &dnsSecurityProfilesListDsModelBotnetDomainsObject{}
				var var4 []dnsSecurityProfilesListDsModelDnsSecurityCategoriesObject
				if len(var1.BotnetDomains.DnsSecurityCategories) != 0 {
					var4 = make([]dnsSecurityProfilesListDsModelDnsSecurityCategoriesObject, 0, len(var1.BotnetDomains.DnsSecurityCategories))
					for var5Index := range var1.BotnetDomains.DnsSecurityCategories {
						var5 := var1.BotnetDomains.DnsSecurityCategories[var5Index]
						var var6 dnsSecurityProfilesListDsModelDnsSecurityCategoriesObject
						var6.Action = types.StringValue(var5.Action)
						var6.LogLevel = types.StringValue(var5.LogLevel)
						var6.Name = types.StringValue(var5.Name)
						var6.PacketCapture = types.StringValue(var5.PacketCapture)
						var4 = append(var4, var6)
					}
				}
				var var7 []dnsSecurityProfilesListDsModelListsObject
				if len(var1.BotnetDomains.Lists) != 0 {
					var7 = make([]dnsSecurityProfilesListDsModelListsObject, 0, len(var1.BotnetDomains.Lists))
					for var8Index := range var1.BotnetDomains.Lists {
						var8 := var1.BotnetDomains.Lists[var8Index]
						var var9 dnsSecurityProfilesListDsModelListsObject
						var var10 *dnsSecurityProfilesListDsModelActionObject
						if var8.Action != nil {
							var10 = &dnsSecurityProfilesListDsModelActionObject{}
							if var8.Action.Alert != nil {
								var10.Alert = types.BoolValue(true)
							}
							if var8.Action.Allow != nil {
								var10.Allow = types.BoolValue(true)
							}
							if var8.Action.Block != nil {
								var10.Block = types.BoolValue(true)
							}
							if var8.Action.Sinkhole != nil {
								var10.Sinkhole = types.BoolValue(true)
							}
						}
						var9.Action = var10
						var9.Name = types.StringValue(var8.Name)
						var9.PacketCapture = types.StringValue(var8.PacketCapture)
						var7 = append(var7, var9)
					}
				}
				var var11 *dnsSecurityProfilesListDsModelSinkholeObject
				if var1.BotnetDomains.Sinkhole != nil {
					var11 = &dnsSecurityProfilesListDsModelSinkholeObject{}
					var11.Ipv4Address = types.StringValue(var1.BotnetDomains.Sinkhole.Ipv4Address)
					var11.Ipv6Address = types.StringValue(var1.BotnetDomains.Sinkhole.Ipv6Address)
				}
				var var12 []dnsSecurityProfilesListDsModelWhitelistObject
				if len(var1.BotnetDomains.Whitelist) != 0 {
					var12 = make([]dnsSecurityProfilesListDsModelWhitelistObject, 0, len(var1.BotnetDomains.Whitelist))
					for var13Index := range var1.BotnetDomains.Whitelist {
						var13 := var1.BotnetDomains.Whitelist[var13Index]
						var var14 dnsSecurityProfilesListDsModelWhitelistObject
						var14.Description = types.StringValue(var13.Description)
						var14.Name = types.StringValue(var13.Name)
						var12 = append(var12, var14)
					}
				}
				var3.DnsSecurityCategories = var4
				var3.Lists = var7
				var3.Sinkhole = var11
				var3.Whitelist = var12
			}
			var2.BotnetDomains = var3
			var2.Description = types.StringValue(var1.Description)
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
	_ datasource.DataSource              = &dnsSecurityProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &dnsSecurityProfilesDataSource{}
)

func NewDnsSecurityProfilesDataSource() datasource.DataSource {
	return &dnsSecurityProfilesDataSource{}
}

type dnsSecurityProfilesDataSource struct {
	client *sase.Client
}

type dnsSecurityProfilesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`
	Folder   types.String `tfsdk:"folder"`

	// Output.
	// Ref: #/components/schemas/dns-security-profiles
	BotnetDomains *dnsSecurityProfilesDsModelBotnetDomainsObject `tfsdk:"botnet_domains"`
	Description   types.String                                   `tfsdk:"description"`
	// input omit: ObjectId
	Name types.String `tfsdk:"name"`
}

type dnsSecurityProfilesDsModelBotnetDomainsObject struct {
	DnsSecurityCategories []dnsSecurityProfilesDsModelDnsSecurityCategoriesObject `tfsdk:"dns_security_categories"`
	Lists                 []dnsSecurityProfilesDsModelListsObject                 `tfsdk:"lists"`
	Sinkhole              *dnsSecurityProfilesDsModelSinkholeObject               `tfsdk:"sinkhole"`
	Whitelist             []dnsSecurityProfilesDsModelWhitelistObject             `tfsdk:"whitelist"`
}

type dnsSecurityProfilesDsModelDnsSecurityCategoriesObject struct {
	Action        types.String `tfsdk:"action"`
	LogLevel      types.String `tfsdk:"log_level"`
	Name          types.String `tfsdk:"name"`
	PacketCapture types.String `tfsdk:"packet_capture"`
}

type dnsSecurityProfilesDsModelListsObject struct {
	Action        *dnsSecurityProfilesDsModelActionObject `tfsdk:"action"`
	Name          types.String                            `tfsdk:"name"`
	PacketCapture types.String                            `tfsdk:"packet_capture"`
}

type dnsSecurityProfilesDsModelActionObject struct {
	Alert    types.Bool `tfsdk:"alert"`
	Allow    types.Bool `tfsdk:"allow"`
	Block    types.Bool `tfsdk:"block"`
	Sinkhole types.Bool `tfsdk:"sinkhole"`
}

type dnsSecurityProfilesDsModelSinkholeObject struct {
	Ipv4Address types.String `tfsdk:"ipv4_address"`
	Ipv6Address types.String `tfsdk:"ipv6_address"`
}

type dnsSecurityProfilesDsModelWhitelistObject struct {
	Description types.String `tfsdk:"description"`
	Name        types.String `tfsdk:"name"`
}

// Metadata returns the data source type name.
func (d *dnsSecurityProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_security_profiles"
}

// Schema defines the schema for this listing data source.
func (d *dnsSecurityProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"botnet_domains": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"dns_security_categories": dsschema.ListNestedAttribute{
						Description: "",
						Computed:    true,
						NestedObject: dsschema.NestedAttributeObject{
							Attributes: map[string]dsschema.Attribute{
								"action": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
								"log_level": dsschema.StringAttribute{
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
							},
						},
					},
					"lists": dsschema.ListNestedAttribute{
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
										"block": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
										"sinkhole": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
								"name": dsschema.StringAttribute{
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
					"sinkhole": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"ipv4_address": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"ipv6_address": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
						},
					},
					"whitelist": dsschema.ListNestedAttribute{
						Description: "",
						Computed:    true,
						NestedObject: dsschema.NestedAttributeObject{
							Attributes: map[string]dsschema.Attribute{
								"description": dsschema.StringAttribute{
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
				},
			},
			"description": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
		},
	}
}

// Configure prepares the struct.
func (d *dnsSecurityProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *dnsSecurityProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state dnsSecurityProfilesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_dns_security_profiles",
		"object_id":                   state.ObjectId.ValueString(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := uSsfsLd.NewClient(d.client)
	input := uSsfsLd.ReadInput{
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
	var var0 *dnsSecurityProfilesDsModelBotnetDomainsObject
	if ans.BotnetDomains != nil {
		var0 = &dnsSecurityProfilesDsModelBotnetDomainsObject{}
		var var1 []dnsSecurityProfilesDsModelDnsSecurityCategoriesObject
		if len(ans.BotnetDomains.DnsSecurityCategories) != 0 {
			var1 = make([]dnsSecurityProfilesDsModelDnsSecurityCategoriesObject, 0, len(ans.BotnetDomains.DnsSecurityCategories))
			for var2Index := range ans.BotnetDomains.DnsSecurityCategories {
				var2 := ans.BotnetDomains.DnsSecurityCategories[var2Index]
				var var3 dnsSecurityProfilesDsModelDnsSecurityCategoriesObject
				var3.Action = types.StringValue(var2.Action)
				var3.LogLevel = types.StringValue(var2.LogLevel)
				var3.Name = types.StringValue(var2.Name)
				var3.PacketCapture = types.StringValue(var2.PacketCapture)
				var1 = append(var1, var3)
			}
		}
		var var4 []dnsSecurityProfilesDsModelListsObject
		if len(ans.BotnetDomains.Lists) != 0 {
			var4 = make([]dnsSecurityProfilesDsModelListsObject, 0, len(ans.BotnetDomains.Lists))
			for var5Index := range ans.BotnetDomains.Lists {
				var5 := ans.BotnetDomains.Lists[var5Index]
				var var6 dnsSecurityProfilesDsModelListsObject
				var var7 *dnsSecurityProfilesDsModelActionObject
				if var5.Action != nil {
					var7 = &dnsSecurityProfilesDsModelActionObject{}
					if var5.Action.Alert != nil {
						var7.Alert = types.BoolValue(true)
					}
					if var5.Action.Allow != nil {
						var7.Allow = types.BoolValue(true)
					}
					if var5.Action.Block != nil {
						var7.Block = types.BoolValue(true)
					}
					if var5.Action.Sinkhole != nil {
						var7.Sinkhole = types.BoolValue(true)
					}
				}
				var6.Action = var7
				var6.Name = types.StringValue(var5.Name)
				var6.PacketCapture = types.StringValue(var5.PacketCapture)
				var4 = append(var4, var6)
			}
		}
		var var8 *dnsSecurityProfilesDsModelSinkholeObject
		if ans.BotnetDomains.Sinkhole != nil {
			var8 = &dnsSecurityProfilesDsModelSinkholeObject{}
			var8.Ipv4Address = types.StringValue(ans.BotnetDomains.Sinkhole.Ipv4Address)
			var8.Ipv6Address = types.StringValue(ans.BotnetDomains.Sinkhole.Ipv6Address)
		}
		var var9 []dnsSecurityProfilesDsModelWhitelistObject
		if len(ans.BotnetDomains.Whitelist) != 0 {
			var9 = make([]dnsSecurityProfilesDsModelWhitelistObject, 0, len(ans.BotnetDomains.Whitelist))
			for var10Index := range ans.BotnetDomains.Whitelist {
				var10 := ans.BotnetDomains.Whitelist[var10Index]
				var var11 dnsSecurityProfilesDsModelWhitelistObject
				var11.Description = types.StringValue(var10.Description)
				var11.Name = types.StringValue(var10.Name)
				var9 = append(var9, var11)
			}
		}
		var0.DnsSecurityCategories = var1
		var0.Lists = var4
		var0.Sinkhole = var8
		var0.Whitelist = var9
	}
	state.BotnetDomains = var0
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &dnsSecurityProfilesResource{}
	_ resource.ResourceWithConfigure   = &dnsSecurityProfilesResource{}
	_ resource.ResourceWithImportState = &dnsSecurityProfilesResource{}
)

func NewDnsSecurityProfilesResource() resource.Resource {
	return &dnsSecurityProfilesResource{}
}

type dnsSecurityProfilesResource struct {
	client *sase.Client
}

type dnsSecurityProfilesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/dns-security-profiles
	BotnetDomains *dnsSecurityProfilesRsModelBotnetDomainsObject `tfsdk:"botnet_domains"`
	Description   types.String                                   `tfsdk:"description"`
	ObjectId      types.String                                   `tfsdk:"object_id"`
	Name          types.String                                   `tfsdk:"name"`
}

type dnsSecurityProfilesRsModelBotnetDomainsObject struct {
	DnsSecurityCategories []dnsSecurityProfilesRsModelDnsSecurityCategoriesObject `tfsdk:"dns_security_categories"`
	Lists                 []dnsSecurityProfilesRsModelListsObject                 `tfsdk:"lists"`
	Sinkhole              *dnsSecurityProfilesRsModelSinkholeObject               `tfsdk:"sinkhole"`
	Whitelist             []dnsSecurityProfilesRsModelWhitelistObject             `tfsdk:"whitelist"`
}

type dnsSecurityProfilesRsModelDnsSecurityCategoriesObject struct {
	Action        types.String `tfsdk:"action"`
	LogLevel      types.String `tfsdk:"log_level"`
	Name          types.String `tfsdk:"name"`
	PacketCapture types.String `tfsdk:"packet_capture"`
}

type dnsSecurityProfilesRsModelListsObject struct {
	Action        *dnsSecurityProfilesRsModelActionObject `tfsdk:"action"`
	Name          types.String                            `tfsdk:"name"`
	PacketCapture types.String                            `tfsdk:"packet_capture"`
}

type dnsSecurityProfilesRsModelActionObject struct {
	Alert    types.Bool `tfsdk:"alert"`
	Allow    types.Bool `tfsdk:"allow"`
	Block    types.Bool `tfsdk:"block"`
	Sinkhole types.Bool `tfsdk:"sinkhole"`
}

type dnsSecurityProfilesRsModelSinkholeObject struct {
	Ipv4Address types.String `tfsdk:"ipv4_address"`
	Ipv6Address types.String `tfsdk:"ipv6_address"`
}

type dnsSecurityProfilesRsModelWhitelistObject struct {
	Description types.String `tfsdk:"description"`
	Name        types.String `tfsdk:"name"`
}

// Metadata returns the data source type name.
func (r *dnsSecurityProfilesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_security_profiles"
}

// Schema defines the schema for this listing data source.
func (r *dnsSecurityProfilesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

			"botnet_domains": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"dns_security_categories": rsschema.ListNestedAttribute{
						Description: "",
						Optional:    true,
						NestedObject: rsschema.NestedAttributeObject{
							Attributes: map[string]rsschema.Attribute{
								"action": rsschema.StringAttribute{
									Description: "",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.String{
										DefaultString("default"),
									},
									Validators: []validator.String{
										stringvalidator.OneOf("default", "allow", "block", "sinkhole"),
									},
								},
								"log_level": rsschema.StringAttribute{
									Description: "",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.String{
										DefaultString("default"),
									},
									Validators: []validator.String{
										stringvalidator.OneOf("default", "none", "low", "informational", "medium", "high", "critical"),
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
							},
						},
					},
					"lists": rsschema.ListNestedAttribute{
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
										"block": rsschema.BoolAttribute{
											Description: "",
											Optional:    true,
										},
										"sinkhole": rsschema.BoolAttribute{
											Description: "",
											Optional:    true,
										},
									},
								},
								"name": rsschema.StringAttribute{
									Description: "",
									Required:    true,
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
					"sinkhole": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"ipv4_address": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
								Validators: []validator.String{
									stringvalidator.OneOf("127.0.0.1", "pan-sinkhole-default-ip"),
								},
							},
							"ipv6_address": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
								Validators: []validator.String{
									stringvalidator.OneOf("::1"),
								},
							},
						},
					},
					"whitelist": rsschema.ListNestedAttribute{
						Description: "",
						Optional:    true,
						NestedObject: rsschema.NestedAttributeObject{
							Attributes: map[string]rsschema.Attribute{
								"description": rsschema.StringAttribute{
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
									PlanModifiers: []planmodifier.String{
										DefaultString(""),
									},
								},
							},
						},
					},
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
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *dnsSecurityProfilesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *dnsSecurityProfilesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state dnsSecurityProfilesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_dns_security_profiles",
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := uSsfsLd.NewClient(r.client)
	input := uSsfsLd.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 fcnKgqA.Config
	var var1 *fcnKgqA.BotnetDomainsObject
	if state.BotnetDomains != nil {
		var1 = &fcnKgqA.BotnetDomainsObject{}
		var var2 []fcnKgqA.DnsSecurityCategoriesObject
		if len(state.BotnetDomains.DnsSecurityCategories) != 0 {
			var2 = make([]fcnKgqA.DnsSecurityCategoriesObject, 0, len(state.BotnetDomains.DnsSecurityCategories))
			for var3Index := range state.BotnetDomains.DnsSecurityCategories {
				var3 := state.BotnetDomains.DnsSecurityCategories[var3Index]
				var var4 fcnKgqA.DnsSecurityCategoriesObject
				var4.Action = var3.Action.ValueString()
				var4.LogLevel = var3.LogLevel.ValueString()
				var4.Name = var3.Name.ValueString()
				var4.PacketCapture = var3.PacketCapture.ValueString()
				var2 = append(var2, var4)
			}
		}
		var1.DnsSecurityCategories = var2
		var var5 []fcnKgqA.ListsObject
		if len(state.BotnetDomains.Lists) != 0 {
			var5 = make([]fcnKgqA.ListsObject, 0, len(state.BotnetDomains.Lists))
			for var6Index := range state.BotnetDomains.Lists {
				var6 := state.BotnetDomains.Lists[var6Index]
				var var7 fcnKgqA.ListsObject
				var var8 *fcnKgqA.ActionObject
				if var6.Action != nil {
					var8 = &fcnKgqA.ActionObject{}
					if var6.Action.Alert.ValueBool() {
						var8.Alert = struct{}{}
					}
					if var6.Action.Allow.ValueBool() {
						var8.Allow = struct{}{}
					}
					if var6.Action.Block.ValueBool() {
						var8.Block = struct{}{}
					}
					if var6.Action.Sinkhole.ValueBool() {
						var8.Sinkhole = struct{}{}
					}
				}
				var7.Action = var8
				var7.Name = var6.Name.ValueString()
				var7.PacketCapture = var6.PacketCapture.ValueString()
				var5 = append(var5, var7)
			}
		}
		var1.Lists = var5
		var var9 *fcnKgqA.SinkholeObject
		if state.BotnetDomains.Sinkhole != nil {
			var9 = &fcnKgqA.SinkholeObject{}
			var9.Ipv4Address = state.BotnetDomains.Sinkhole.Ipv4Address.ValueString()
			var9.Ipv6Address = state.BotnetDomains.Sinkhole.Ipv6Address.ValueString()
		}
		var1.Sinkhole = var9
		var var10 []fcnKgqA.WhitelistObject
		if len(state.BotnetDomains.Whitelist) != 0 {
			var10 = make([]fcnKgqA.WhitelistObject, 0, len(state.BotnetDomains.Whitelist))
			for var11Index := range state.BotnetDomains.Whitelist {
				var11 := state.BotnetDomains.Whitelist[var11Index]
				var var12 fcnKgqA.WhitelistObject
				var12.Description = var11.Description.ValueString()
				var12.Name = var11.Name.ValueString()
				var10 = append(var10, var12)
			}
		}
		var1.Whitelist = var10
	}
	var0.BotnetDomains = var1
	var0.Description = state.Description.ValueString()
	var0.Name = state.Name.ValueString()
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
	var var13 *dnsSecurityProfilesRsModelBotnetDomainsObject
	if ans.BotnetDomains != nil {
		var13 = &dnsSecurityProfilesRsModelBotnetDomainsObject{}
		var var14 []dnsSecurityProfilesRsModelDnsSecurityCategoriesObject
		if len(ans.BotnetDomains.DnsSecurityCategories) != 0 {
			var14 = make([]dnsSecurityProfilesRsModelDnsSecurityCategoriesObject, 0, len(ans.BotnetDomains.DnsSecurityCategories))
			for var15Index := range ans.BotnetDomains.DnsSecurityCategories {
				var15 := ans.BotnetDomains.DnsSecurityCategories[var15Index]
				var var16 dnsSecurityProfilesRsModelDnsSecurityCategoriesObject
				var16.Action = types.StringValue(var15.Action)
				var16.LogLevel = types.StringValue(var15.LogLevel)
				var16.Name = types.StringValue(var15.Name)
				var16.PacketCapture = types.StringValue(var15.PacketCapture)
				var14 = append(var14, var16)
			}
		}
		var var17 []dnsSecurityProfilesRsModelListsObject
		if len(ans.BotnetDomains.Lists) != 0 {
			var17 = make([]dnsSecurityProfilesRsModelListsObject, 0, len(ans.BotnetDomains.Lists))
			for var18Index := range ans.BotnetDomains.Lists {
				var18 := ans.BotnetDomains.Lists[var18Index]
				var var19 dnsSecurityProfilesRsModelListsObject
				var var20 *dnsSecurityProfilesRsModelActionObject
				if var18.Action != nil {
					var20 = &dnsSecurityProfilesRsModelActionObject{}
					if var18.Action.Alert != nil {
						var20.Alert = types.BoolValue(true)
					}
					if var18.Action.Allow != nil {
						var20.Allow = types.BoolValue(true)
					}
					if var18.Action.Block != nil {
						var20.Block = types.BoolValue(true)
					}
					if var18.Action.Sinkhole != nil {
						var20.Sinkhole = types.BoolValue(true)
					}
				}
				var19.Action = var20
				var19.Name = types.StringValue(var18.Name)
				var19.PacketCapture = types.StringValue(var18.PacketCapture)
				var17 = append(var17, var19)
			}
		}
		var var21 *dnsSecurityProfilesRsModelSinkholeObject
		if ans.BotnetDomains.Sinkhole != nil {
			var21 = &dnsSecurityProfilesRsModelSinkholeObject{}
			var21.Ipv4Address = types.StringValue(ans.BotnetDomains.Sinkhole.Ipv4Address)
			var21.Ipv6Address = types.StringValue(ans.BotnetDomains.Sinkhole.Ipv6Address)
		}
		var var22 []dnsSecurityProfilesRsModelWhitelistObject
		if len(ans.BotnetDomains.Whitelist) != 0 {
			var22 = make([]dnsSecurityProfilesRsModelWhitelistObject, 0, len(ans.BotnetDomains.Whitelist))
			for var23Index := range ans.BotnetDomains.Whitelist {
				var23 := ans.BotnetDomains.Whitelist[var23Index]
				var var24 dnsSecurityProfilesRsModelWhitelistObject
				var24.Description = types.StringValue(var23.Description)
				var24.Name = types.StringValue(var23.Name)
				var22 = append(var22, var24)
			}
		}
		var13.DnsSecurityCategories = var14
		var13.Lists = var17
		var13.Sinkhole = var21
		var13.Whitelist = var22
	}
	state.BotnetDomains = var13
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *dnsSecurityProfilesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state dnsSecurityProfilesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_dns_security_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := uSsfsLd.NewClient(r.client)
	input := uSsfsLd.ReadInput{
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
	var var0 *dnsSecurityProfilesRsModelBotnetDomainsObject
	if ans.BotnetDomains != nil {
		var0 = &dnsSecurityProfilesRsModelBotnetDomainsObject{}
		var var1 []dnsSecurityProfilesRsModelDnsSecurityCategoriesObject
		if len(ans.BotnetDomains.DnsSecurityCategories) != 0 {
			var1 = make([]dnsSecurityProfilesRsModelDnsSecurityCategoriesObject, 0, len(ans.BotnetDomains.DnsSecurityCategories))
			for var2Index := range ans.BotnetDomains.DnsSecurityCategories {
				var2 := ans.BotnetDomains.DnsSecurityCategories[var2Index]
				var var3 dnsSecurityProfilesRsModelDnsSecurityCategoriesObject
				var3.Action = types.StringValue(var2.Action)
				var3.LogLevel = types.StringValue(var2.LogLevel)
				var3.Name = types.StringValue(var2.Name)
				var3.PacketCapture = types.StringValue(var2.PacketCapture)
				var1 = append(var1, var3)
			}
		}
		var var4 []dnsSecurityProfilesRsModelListsObject
		if len(ans.BotnetDomains.Lists) != 0 {
			var4 = make([]dnsSecurityProfilesRsModelListsObject, 0, len(ans.BotnetDomains.Lists))
			for var5Index := range ans.BotnetDomains.Lists {
				var5 := ans.BotnetDomains.Lists[var5Index]
				var var6 dnsSecurityProfilesRsModelListsObject
				var var7 *dnsSecurityProfilesRsModelActionObject
				if var5.Action != nil {
					var7 = &dnsSecurityProfilesRsModelActionObject{}
					if var5.Action.Alert != nil {
						var7.Alert = types.BoolValue(true)
					}
					if var5.Action.Allow != nil {
						var7.Allow = types.BoolValue(true)
					}
					if var5.Action.Block != nil {
						var7.Block = types.BoolValue(true)
					}
					if var5.Action.Sinkhole != nil {
						var7.Sinkhole = types.BoolValue(true)
					}
				}
				var6.Action = var7
				var6.Name = types.StringValue(var5.Name)
				var6.PacketCapture = types.StringValue(var5.PacketCapture)
				var4 = append(var4, var6)
			}
		}
		var var8 *dnsSecurityProfilesRsModelSinkholeObject
		if ans.BotnetDomains.Sinkhole != nil {
			var8 = &dnsSecurityProfilesRsModelSinkholeObject{}
			var8.Ipv4Address = types.StringValue(ans.BotnetDomains.Sinkhole.Ipv4Address)
			var8.Ipv6Address = types.StringValue(ans.BotnetDomains.Sinkhole.Ipv6Address)
		}
		var var9 []dnsSecurityProfilesRsModelWhitelistObject
		if len(ans.BotnetDomains.Whitelist) != 0 {
			var9 = make([]dnsSecurityProfilesRsModelWhitelistObject, 0, len(ans.BotnetDomains.Whitelist))
			for var10Index := range ans.BotnetDomains.Whitelist {
				var10 := ans.BotnetDomains.Whitelist[var10Index]
				var var11 dnsSecurityProfilesRsModelWhitelistObject
				var11.Description = types.StringValue(var10.Description)
				var11.Name = types.StringValue(var10.Name)
				var9 = append(var9, var11)
			}
		}
		var0.DnsSecurityCategories = var1
		var0.Lists = var4
		var0.Sinkhole = var8
		var0.Whitelist = var9
	}
	state.BotnetDomains = var0
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *dnsSecurityProfilesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state dnsSecurityProfilesRsModel
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
		"resource_name":               "sase_dns_security_profiles",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := uSsfsLd.NewClient(r.client)
	input := uSsfsLd.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 fcnKgqA.Config
	var var1 *fcnKgqA.BotnetDomainsObject
	if plan.BotnetDomains != nil {
		var1 = &fcnKgqA.BotnetDomainsObject{}
		var var2 []fcnKgqA.DnsSecurityCategoriesObject
		if len(plan.BotnetDomains.DnsSecurityCategories) != 0 {
			var2 = make([]fcnKgqA.DnsSecurityCategoriesObject, 0, len(plan.BotnetDomains.DnsSecurityCategories))
			for var3Index := range plan.BotnetDomains.DnsSecurityCategories {
				var3 := plan.BotnetDomains.DnsSecurityCategories[var3Index]
				var var4 fcnKgqA.DnsSecurityCategoriesObject
				var4.Action = var3.Action.ValueString()
				var4.LogLevel = var3.LogLevel.ValueString()
				var4.Name = var3.Name.ValueString()
				var4.PacketCapture = var3.PacketCapture.ValueString()
				var2 = append(var2, var4)
			}
		}
		var1.DnsSecurityCategories = var2
		var var5 []fcnKgqA.ListsObject
		if len(plan.BotnetDomains.Lists) != 0 {
			var5 = make([]fcnKgqA.ListsObject, 0, len(plan.BotnetDomains.Lists))
			for var6Index := range plan.BotnetDomains.Lists {
				var6 := plan.BotnetDomains.Lists[var6Index]
				var var7 fcnKgqA.ListsObject
				var var8 *fcnKgqA.ActionObject
				if var6.Action != nil {
					var8 = &fcnKgqA.ActionObject{}
					if var6.Action.Alert.ValueBool() {
						var8.Alert = struct{}{}
					}
					if var6.Action.Allow.ValueBool() {
						var8.Allow = struct{}{}
					}
					if var6.Action.Block.ValueBool() {
						var8.Block = struct{}{}
					}
					if var6.Action.Sinkhole.ValueBool() {
						var8.Sinkhole = struct{}{}
					}
				}
				var7.Action = var8
				var7.Name = var6.Name.ValueString()
				var7.PacketCapture = var6.PacketCapture.ValueString()
				var5 = append(var5, var7)
			}
		}
		var1.Lists = var5
		var var9 *fcnKgqA.SinkholeObject
		if plan.BotnetDomains.Sinkhole != nil {
			var9 = &fcnKgqA.SinkholeObject{}
			var9.Ipv4Address = plan.BotnetDomains.Sinkhole.Ipv4Address.ValueString()
			var9.Ipv6Address = plan.BotnetDomains.Sinkhole.Ipv6Address.ValueString()
		}
		var1.Sinkhole = var9
		var var10 []fcnKgqA.WhitelistObject
		if len(plan.BotnetDomains.Whitelist) != 0 {
			var10 = make([]fcnKgqA.WhitelistObject, 0, len(plan.BotnetDomains.Whitelist))
			for var11Index := range plan.BotnetDomains.Whitelist {
				var11 := plan.BotnetDomains.Whitelist[var11Index]
				var var12 fcnKgqA.WhitelistObject
				var12.Description = var11.Description.ValueString()
				var12.Name = var11.Name.ValueString()
				var10 = append(var10, var12)
			}
		}
		var1.Whitelist = var10
	}
	var0.BotnetDomains = var1
	var0.Description = plan.Description.ValueString()
	var0.Name = plan.Name.ValueString()
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var13 *dnsSecurityProfilesRsModelBotnetDomainsObject
	if ans.BotnetDomains != nil {
		var13 = &dnsSecurityProfilesRsModelBotnetDomainsObject{}
		var var14 []dnsSecurityProfilesRsModelDnsSecurityCategoriesObject
		if len(ans.BotnetDomains.DnsSecurityCategories) != 0 {
			var14 = make([]dnsSecurityProfilesRsModelDnsSecurityCategoriesObject, 0, len(ans.BotnetDomains.DnsSecurityCategories))
			for var15Index := range ans.BotnetDomains.DnsSecurityCategories {
				var15 := ans.BotnetDomains.DnsSecurityCategories[var15Index]
				var var16 dnsSecurityProfilesRsModelDnsSecurityCategoriesObject
				var16.Action = types.StringValue(var15.Action)
				var16.LogLevel = types.StringValue(var15.LogLevel)
				var16.Name = types.StringValue(var15.Name)
				var16.PacketCapture = types.StringValue(var15.PacketCapture)
				var14 = append(var14, var16)
			}
		}
		var var17 []dnsSecurityProfilesRsModelListsObject
		if len(ans.BotnetDomains.Lists) != 0 {
			var17 = make([]dnsSecurityProfilesRsModelListsObject, 0, len(ans.BotnetDomains.Lists))
			for var18Index := range ans.BotnetDomains.Lists {
				var18 := ans.BotnetDomains.Lists[var18Index]
				var var19 dnsSecurityProfilesRsModelListsObject
				var var20 *dnsSecurityProfilesRsModelActionObject
				if var18.Action != nil {
					var20 = &dnsSecurityProfilesRsModelActionObject{}
					if var18.Action.Alert != nil {
						var20.Alert = types.BoolValue(true)
					}
					if var18.Action.Allow != nil {
						var20.Allow = types.BoolValue(true)
					}
					if var18.Action.Block != nil {
						var20.Block = types.BoolValue(true)
					}
					if var18.Action.Sinkhole != nil {
						var20.Sinkhole = types.BoolValue(true)
					}
				}
				var19.Action = var20
				var19.Name = types.StringValue(var18.Name)
				var19.PacketCapture = types.StringValue(var18.PacketCapture)
				var17 = append(var17, var19)
			}
		}
		var var21 *dnsSecurityProfilesRsModelSinkholeObject
		if ans.BotnetDomains.Sinkhole != nil {
			var21 = &dnsSecurityProfilesRsModelSinkholeObject{}
			var21.Ipv4Address = types.StringValue(ans.BotnetDomains.Sinkhole.Ipv4Address)
			var21.Ipv6Address = types.StringValue(ans.BotnetDomains.Sinkhole.Ipv6Address)
		}
		var var22 []dnsSecurityProfilesRsModelWhitelistObject
		if len(ans.BotnetDomains.Whitelist) != 0 {
			var22 = make([]dnsSecurityProfilesRsModelWhitelistObject, 0, len(ans.BotnetDomains.Whitelist))
			for var23Index := range ans.BotnetDomains.Whitelist {
				var23 := ans.BotnetDomains.Whitelist[var23Index]
				var var24 dnsSecurityProfilesRsModelWhitelistObject
				var24.Description = types.StringValue(var23.Description)
				var24.Name = types.StringValue(var23.Name)
				var22 = append(var22, var24)
			}
		}
		var13.DnsSecurityCategories = var14
		var13.Lists = var17
		var13.Sinkhole = var21
		var13.Whitelist = var22
	}
	state.BotnetDomains = var13
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *dnsSecurityProfilesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name":               "sase_dns_security_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	svc := uSsfsLd.NewClient(r.client)
	input := uSsfsLd.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *dnsSecurityProfilesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
