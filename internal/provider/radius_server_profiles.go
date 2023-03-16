package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	ieMayhq "github.com/paloaltonetworks/sase-go/netsec/schema/radius/server/profiles"
	bVmbuOb "github.com/paloaltonetworks/sase-go/netsec/service/v1/radiusserverprofiles"

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
	_ datasource.DataSource              = &radiusServerProfilesListDataSource{}
	_ datasource.DataSourceWithConfigure = &radiusServerProfilesListDataSource{}
)

func NewRadiusServerProfilesListDataSource() datasource.DataSource {
	return &radiusServerProfilesListDataSource{}
}

type radiusServerProfilesListDataSource struct {
	client *sase.Client
}

type radiusServerProfilesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Folder types.String `tfsdk:"folder"`
	Name   types.String `tfsdk:"name"`

	// Output.
	Data []radiusServerProfilesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type radiusServerProfilesListDsModelConfig struct {
	ObjectId types.String                                   `tfsdk:"object_id"`
	Protocol *radiusServerProfilesListDsModelProtocolObject `tfsdk:"protocol"`
	Retries  types.Int64                                    `tfsdk:"retries"`
	Server   []radiusServerProfilesListDsModelServerObject  `tfsdk:"server"`
	Timeout  types.Int64                                    `tfsdk:"timeout"`
}

type radiusServerProfilesListDsModelProtocolObject struct {
	CHAP           types.Bool                                           `tfsdk:"c_h_a_p"`
	EAPTTLSWithPAP *radiusServerProfilesListDsModelEAPTTLSWithPAPObject `tfsdk:"e_a_p_t_t_l_s_with_p_a_p"`
	PAP            types.Bool                                           `tfsdk:"p_a_p"`
	PEAPMSCHAPv2   *radiusServerProfilesListDsModelPEAPMSCHAPv2Object   `tfsdk:"p_e_a_p_m_s_c_h_a_pv2"`
	PEAPWithGTC    *radiusServerProfilesListDsModelPEAPWithGTCObject    `tfsdk:"p_e_a_p_with_g_t_c"`
}

type radiusServerProfilesListDsModelEAPTTLSWithPAPObject struct {
	AnonOuterId       types.Bool   `tfsdk:"anon_outer_id"`
	RadiusCertProfile types.String `tfsdk:"radius_cert_profile"`
}

type radiusServerProfilesListDsModelPEAPMSCHAPv2Object struct {
	AllowPwdChange    types.Bool   `tfsdk:"allow_pwd_change"`
	AnonOuterId       types.Bool   `tfsdk:"anon_outer_id"`
	RadiusCertProfile types.String `tfsdk:"radius_cert_profile"`
}

type radiusServerProfilesListDsModelPEAPWithGTCObject struct {
	AnonOuterId       types.Bool   `tfsdk:"anon_outer_id"`
	RadiusCertProfile types.String `tfsdk:"radius_cert_profile"`
}

type radiusServerProfilesListDsModelServerObject struct {
	IpAddress types.String `tfsdk:"ip_address"`
	Name      types.String `tfsdk:"name"`
	Port      types.Int64  `tfsdk:"port"`
	Secret    types.String `tfsdk:"secret"`
}

// Metadata returns the data source type name.
func (d *radiusServerProfilesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_radius_server_profiles_list"
}

// Schema defines the schema for this listing data source.
func (d *radiusServerProfilesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"name": dsschema.StringAttribute{
				Description: "The name of the entry",
				Optional:    true,
			},

			// Output.
			"data": dsschema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"object_id": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"protocol": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"c_h_a_p": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"e_a_p_t_t_l_s_with_p_a_p": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"anon_outer_id": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
										"radius_cert_profile": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
								"p_a_p": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"p_e_a_p_m_s_c_h_a_pv2": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"allow_pwd_change": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
										"anon_outer_id": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
										"radius_cert_profile": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
								"p_e_a_p_with_g_t_c": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"anon_outer_id": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
										"radius_cert_profile": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
							},
						},
						"retries": dsschema.Int64Attribute{
							Description: "",
							Computed:    true,
						},
						"server": dsschema.ListNestedAttribute{
							Description: "",
							Computed:    true,
							NestedObject: dsschema.NestedAttributeObject{
								Attributes: map[string]dsschema.Attribute{
									"ip_address": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"name": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"port": dsschema.Int64Attribute{
										Description: "",
										Computed:    true,
									},
									"secret": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
						},
						"timeout": dsschema.Int64Attribute{
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
func (d *radiusServerProfilesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *radiusServerProfilesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state radiusServerProfilesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_radius_server_profiles_list",
		"limit":            state.Limit.ValueInt64(),
		"has_limit":        !state.Limit.IsNull(),
		"offset":           state.Offset.ValueInt64(),
		"has_offset":       !state.Offset.IsNull(),
		"folder":           state.Folder.ValueString(),
		"name":             state.Name.ValueString(),
		"has_name":         !state.Name.IsNull(),
	})

	// Prepare to run the command.
	svc := bVmbuOb.NewClient(d.client)
	input := bVmbuOb.ListInput{
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
	idBuilder.WriteString(input.Folder)
	idBuilder.WriteString(IdSeparator)
	if input.Name != nil {
		idBuilder.WriteString(*input.Name)
	}
	state.Id = types.StringValue(idBuilder.String())
	var var0 []radiusServerProfilesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]radiusServerProfilesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 radiusServerProfilesListDsModelConfig
			var var3 *radiusServerProfilesListDsModelProtocolObject
			if var1.Protocol != nil {
				var3 = &radiusServerProfilesListDsModelProtocolObject{}
				var var4 *radiusServerProfilesListDsModelEAPTTLSWithPAPObject
				if var1.Protocol.EAPTTLSWithPAP != nil {
					var4 = &radiusServerProfilesListDsModelEAPTTLSWithPAPObject{}
					var4.AnonOuterId = types.BoolValue(var1.Protocol.EAPTTLSWithPAP.AnonOuterId)
					var4.RadiusCertProfile = types.StringValue(var1.Protocol.EAPTTLSWithPAP.RadiusCertProfile)
				}
				var var5 *radiusServerProfilesListDsModelPEAPMSCHAPv2Object
				if var1.Protocol.PEAPMSCHAPv2 != nil {
					var5 = &radiusServerProfilesListDsModelPEAPMSCHAPv2Object{}
					var5.AllowPwdChange = types.BoolValue(var1.Protocol.PEAPMSCHAPv2.AllowPwdChange)
					var5.AnonOuterId = types.BoolValue(var1.Protocol.PEAPMSCHAPv2.AnonOuterId)
					var5.RadiusCertProfile = types.StringValue(var1.Protocol.PEAPMSCHAPv2.RadiusCertProfile)
				}
				var var6 *radiusServerProfilesListDsModelPEAPWithGTCObject
				if var1.Protocol.PEAPWithGTC != nil {
					var6 = &radiusServerProfilesListDsModelPEAPWithGTCObject{}
					var6.AnonOuterId = types.BoolValue(var1.Protocol.PEAPWithGTC.AnonOuterId)
					var6.RadiusCertProfile = types.StringValue(var1.Protocol.PEAPWithGTC.RadiusCertProfile)
				}
				if var1.Protocol.CHAP != nil {
					var3.CHAP = types.BoolValue(true)
				}
				var3.EAPTTLSWithPAP = var4
				if var1.Protocol.PAP != nil {
					var3.PAP = types.BoolValue(true)
				}
				var3.PEAPMSCHAPv2 = var5
				var3.PEAPWithGTC = var6
			}
			var var7 []radiusServerProfilesListDsModelServerObject
			if len(var1.Server) != 0 {
				var7 = make([]radiusServerProfilesListDsModelServerObject, 0, len(var1.Server))
				for var8Index := range var1.Server {
					var8 := var1.Server[var8Index]
					var var9 radiusServerProfilesListDsModelServerObject
					var9.IpAddress = types.StringValue(var8.IpAddress)
					var9.Name = types.StringValue(var8.Name)
					var9.Port = types.Int64Value(var8.Port)
					var9.Secret = types.StringValue(var8.Secret)
					var7 = append(var7, var9)
				}
			}
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Protocol = var3
			var2.Retries = types.Int64Value(var1.Retries)
			var2.Server = var7
			var2.Timeout = types.Int64Value(var1.Timeout)
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
	_ datasource.DataSource              = &radiusServerProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &radiusServerProfilesDataSource{}
)

func NewRadiusServerProfilesDataSource() datasource.DataSource {
	return &radiusServerProfilesDataSource{}
}

type radiusServerProfilesDataSource struct {
	client *sase.Client
}

type radiusServerProfilesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`

	// Output.
	// Ref: #/components/schemas/radius-server-profiles
	// input omit: ObjectId
	Protocol *radiusServerProfilesDsModelProtocolObject `tfsdk:"protocol"`
	Retries  types.Int64                                `tfsdk:"retries"`
	Server   []radiusServerProfilesDsModelServerObject  `tfsdk:"server"`
	Timeout  types.Int64                                `tfsdk:"timeout"`
}

type radiusServerProfilesDsModelProtocolObject struct {
	CHAP           types.Bool                                       `tfsdk:"c_h_a_p"`
	EAPTTLSWithPAP *radiusServerProfilesDsModelEAPTTLSWithPAPObject `tfsdk:"e_a_p_t_t_l_s_with_p_a_p"`
	PAP            types.Bool                                       `tfsdk:"p_a_p"`
	PEAPMSCHAPv2   *radiusServerProfilesDsModelPEAPMSCHAPv2Object   `tfsdk:"p_e_a_p_m_s_c_h_a_pv2"`
	PEAPWithGTC    *radiusServerProfilesDsModelPEAPWithGTCObject    `tfsdk:"p_e_a_p_with_g_t_c"`
}

type radiusServerProfilesDsModelEAPTTLSWithPAPObject struct {
	AnonOuterId       types.Bool   `tfsdk:"anon_outer_id"`
	RadiusCertProfile types.String `tfsdk:"radius_cert_profile"`
}

type radiusServerProfilesDsModelPEAPMSCHAPv2Object struct {
	AllowPwdChange    types.Bool   `tfsdk:"allow_pwd_change"`
	AnonOuterId       types.Bool   `tfsdk:"anon_outer_id"`
	RadiusCertProfile types.String `tfsdk:"radius_cert_profile"`
}

type radiusServerProfilesDsModelPEAPWithGTCObject struct {
	AnonOuterId       types.Bool   `tfsdk:"anon_outer_id"`
	RadiusCertProfile types.String `tfsdk:"radius_cert_profile"`
}

type radiusServerProfilesDsModelServerObject struct {
	IpAddress types.String `tfsdk:"ip_address"`
	Name      types.String `tfsdk:"name"`
	Port      types.Int64  `tfsdk:"port"`
	Secret    types.String `tfsdk:"secret"`
}

// Metadata returns the data source type name.
func (d *radiusServerProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_radius_server_profiles"
}

// Schema defines the schema for this listing data source.
func (d *radiusServerProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"protocol": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"c_h_a_p": dsschema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"e_a_p_t_t_l_s_with_p_a_p": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"anon_outer_id": dsschema.BoolAttribute{
								Description: "",
								Computed:    true,
							},
							"radius_cert_profile": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
						},
					},
					"p_a_p": dsschema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"p_e_a_p_m_s_c_h_a_pv2": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"allow_pwd_change": dsschema.BoolAttribute{
								Description: "",
								Computed:    true,
							},
							"anon_outer_id": dsschema.BoolAttribute{
								Description: "",
								Computed:    true,
							},
							"radius_cert_profile": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
						},
					},
					"p_e_a_p_with_g_t_c": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"anon_outer_id": dsschema.BoolAttribute{
								Description: "",
								Computed:    true,
							},
							"radius_cert_profile": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
						},
					},
				},
			},
			"retries": dsschema.Int64Attribute{
				Description: "",
				Computed:    true,
			},
			"server": dsschema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"ip_address": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"port": dsschema.Int64Attribute{
							Description: "",
							Computed:    true,
						},
						"secret": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
					},
				},
			},
			"timeout": dsschema.Int64Attribute{
				Description: "",
				Computed:    true,
			},
		},
	}
}

// Configure prepares the struct.
func (d *radiusServerProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *radiusServerProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state radiusServerProfilesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_radius_server_profiles",
		"object_id":        state.ObjectId.ValueString(),
	})

	// Prepare to run the command.
	svc := bVmbuOb.NewClient(d.client)
	input := bVmbuOb.ReadInput{
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
	var var0 *radiusServerProfilesDsModelProtocolObject
	if ans.Protocol != nil {
		var0 = &radiusServerProfilesDsModelProtocolObject{}
		var var1 *radiusServerProfilesDsModelEAPTTLSWithPAPObject
		if ans.Protocol.EAPTTLSWithPAP != nil {
			var1 = &radiusServerProfilesDsModelEAPTTLSWithPAPObject{}
			var1.AnonOuterId = types.BoolValue(ans.Protocol.EAPTTLSWithPAP.AnonOuterId)
			var1.RadiusCertProfile = types.StringValue(ans.Protocol.EAPTTLSWithPAP.RadiusCertProfile)
		}
		var var2 *radiusServerProfilesDsModelPEAPMSCHAPv2Object
		if ans.Protocol.PEAPMSCHAPv2 != nil {
			var2 = &radiusServerProfilesDsModelPEAPMSCHAPv2Object{}
			var2.AllowPwdChange = types.BoolValue(ans.Protocol.PEAPMSCHAPv2.AllowPwdChange)
			var2.AnonOuterId = types.BoolValue(ans.Protocol.PEAPMSCHAPv2.AnonOuterId)
			var2.RadiusCertProfile = types.StringValue(ans.Protocol.PEAPMSCHAPv2.RadiusCertProfile)
		}
		var var3 *radiusServerProfilesDsModelPEAPWithGTCObject
		if ans.Protocol.PEAPWithGTC != nil {
			var3 = &radiusServerProfilesDsModelPEAPWithGTCObject{}
			var3.AnonOuterId = types.BoolValue(ans.Protocol.PEAPWithGTC.AnonOuterId)
			var3.RadiusCertProfile = types.StringValue(ans.Protocol.PEAPWithGTC.RadiusCertProfile)
		}
		if ans.Protocol.CHAP != nil {
			var0.CHAP = types.BoolValue(true)
		}
		var0.EAPTTLSWithPAP = var1
		if ans.Protocol.PAP != nil {
			var0.PAP = types.BoolValue(true)
		}
		var0.PEAPMSCHAPv2 = var2
		var0.PEAPWithGTC = var3
	}
	var var4 []radiusServerProfilesDsModelServerObject
	if len(ans.Server) != 0 {
		var4 = make([]radiusServerProfilesDsModelServerObject, 0, len(ans.Server))
		for var5Index := range ans.Server {
			var5 := ans.Server[var5Index]
			var var6 radiusServerProfilesDsModelServerObject
			var6.IpAddress = types.StringValue(var5.IpAddress)
			var6.Name = types.StringValue(var5.Name)
			var6.Port = types.Int64Value(var5.Port)
			var6.Secret = types.StringValue(var5.Secret)
			var4 = append(var4, var6)
		}
	}
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Protocol = var0
	state.Retries = types.Int64Value(ans.Retries)
	state.Server = var4
	state.Timeout = types.Int64Value(ans.Timeout)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &radiusServerProfilesResource{}
	_ resource.ResourceWithConfigure   = &radiusServerProfilesResource{}
	_ resource.ResourceWithImportState = &radiusServerProfilesResource{}
)

func NewRadiusServerProfilesResource() resource.Resource {
	return &radiusServerProfilesResource{}
}

type radiusServerProfilesResource struct {
	client *sase.Client
}

type radiusServerProfilesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/radius-server-profiles
	ObjectId types.String                               `tfsdk:"object_id"`
	Protocol *radiusServerProfilesRsModelProtocolObject `tfsdk:"protocol"`
	Retries  types.Int64                                `tfsdk:"retries"`
	Server   []radiusServerProfilesRsModelServerObject  `tfsdk:"server"`
	Timeout  types.Int64                                `tfsdk:"timeout"`
}

type radiusServerProfilesRsModelProtocolObject struct {
	CHAP           types.Bool                                       `tfsdk:"c_h_a_p"`
	EAPTTLSWithPAP *radiusServerProfilesRsModelEAPTTLSWithPAPObject `tfsdk:"e_a_p_t_t_l_s_with_p_a_p"`
	PAP            types.Bool                                       `tfsdk:"p_a_p"`
	PEAPMSCHAPv2   *radiusServerProfilesRsModelPEAPMSCHAPv2Object   `tfsdk:"p_e_a_p_m_s_c_h_a_pv2"`
	PEAPWithGTC    *radiusServerProfilesRsModelPEAPWithGTCObject    `tfsdk:"p_e_a_p_with_g_t_c"`
}

type radiusServerProfilesRsModelEAPTTLSWithPAPObject struct {
	AnonOuterId       types.Bool   `tfsdk:"anon_outer_id"`
	RadiusCertProfile types.String `tfsdk:"radius_cert_profile"`
}

type radiusServerProfilesRsModelPEAPMSCHAPv2Object struct {
	AllowPwdChange    types.Bool   `tfsdk:"allow_pwd_change"`
	AnonOuterId       types.Bool   `tfsdk:"anon_outer_id"`
	RadiusCertProfile types.String `tfsdk:"radius_cert_profile"`
}

type radiusServerProfilesRsModelPEAPWithGTCObject struct {
	AnonOuterId       types.Bool   `tfsdk:"anon_outer_id"`
	RadiusCertProfile types.String `tfsdk:"radius_cert_profile"`
}

type radiusServerProfilesRsModelServerObject struct {
	IpAddress types.String `tfsdk:"ip_address"`
	Name      types.String `tfsdk:"name"`
	Port      types.Int64  `tfsdk:"port"`
	Secret    types.String `tfsdk:"secret"`
}

// Metadata returns the data source type name.
func (r *radiusServerProfilesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_radius_server_profiles"
}

// Schema defines the schema for this listing data source.
func (r *radiusServerProfilesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

			"object_id": rsschema.StringAttribute{
				Description: "",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"protocol": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"c_h_a_p": rsschema.BoolAttribute{
						Description: "",
						Optional:    true,
					},
					"e_a_p_t_t_l_s_with_p_a_p": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"anon_outer_id": rsschema.BoolAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Bool{
									DefaultBool(false),
								},
							},
							"radius_cert_profile": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
						},
					},
					"p_a_p": rsschema.BoolAttribute{
						Description: "",
						Optional:    true,
					},
					"p_e_a_p_m_s_c_h_a_pv2": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"allow_pwd_change": rsschema.BoolAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Bool{
									DefaultBool(false),
								},
							},
							"anon_outer_id": rsschema.BoolAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Bool{
									DefaultBool(false),
								},
							},
							"radius_cert_profile": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
						},
					},
					"p_e_a_p_with_g_t_c": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"anon_outer_id": rsschema.BoolAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Bool{
									DefaultBool(false),
								},
							},
							"radius_cert_profile": rsschema.StringAttribute{
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
			"retries": rsschema.Int64Attribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					DefaultInt64(0),
				},
				Validators: []validator.Int64{
					int64validator.Between(1, 5),
				},
			},
			"server": rsschema.ListNestedAttribute{
				Description: "",
				Required:    true,
				NestedObject: rsschema.NestedAttributeObject{
					Attributes: map[string]rsschema.Attribute{
						"ip_address": rsschema.StringAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
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
						"port": rsschema.Int64Attribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.Int64{
								DefaultInt64(0),
							},
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
						"secret": rsschema.StringAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
							Validators: []validator.String{
								stringvalidator.LengthAtMost(64),
							},
						},
					},
				},
			},
			"timeout": rsschema.Int64Attribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					DefaultInt64(0),
				},
				Validators: []validator.Int64{
					int64validator.Between(1, 120),
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *radiusServerProfilesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *radiusServerProfilesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state radiusServerProfilesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "sase_radius_server_profiles",
		"folder":        state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := bVmbuOb.NewClient(r.client)
	input := bVmbuOb.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 ieMayhq.Config
	var var1 *ieMayhq.ProtocolObject
	if state.Protocol != nil {
		var1 = &ieMayhq.ProtocolObject{}
		if state.Protocol.CHAP.ValueBool() {
			var1.CHAP = struct{}{}
		}
		var var2 *ieMayhq.EAPTTLSWithPAPObject
		if state.Protocol.EAPTTLSWithPAP != nil {
			var2 = &ieMayhq.EAPTTLSWithPAPObject{}
			var2.AnonOuterId = state.Protocol.EAPTTLSWithPAP.AnonOuterId.ValueBool()
			var2.RadiusCertProfile = state.Protocol.EAPTTLSWithPAP.RadiusCertProfile.ValueString()
		}
		var1.EAPTTLSWithPAP = var2
		if state.Protocol.PAP.ValueBool() {
			var1.PAP = struct{}{}
		}
		var var3 *ieMayhq.PEAPMSCHAPv2Object
		if state.Protocol.PEAPMSCHAPv2 != nil {
			var3 = &ieMayhq.PEAPMSCHAPv2Object{}
			var3.AllowPwdChange = state.Protocol.PEAPMSCHAPv2.AllowPwdChange.ValueBool()
			var3.AnonOuterId = state.Protocol.PEAPMSCHAPv2.AnonOuterId.ValueBool()
			var3.RadiusCertProfile = state.Protocol.PEAPMSCHAPv2.RadiusCertProfile.ValueString()
		}
		var1.PEAPMSCHAPv2 = var3
		var var4 *ieMayhq.PEAPWithGTCObject
		if state.Protocol.PEAPWithGTC != nil {
			var4 = &ieMayhq.PEAPWithGTCObject{}
			var4.AnonOuterId = state.Protocol.PEAPWithGTC.AnonOuterId.ValueBool()
			var4.RadiusCertProfile = state.Protocol.PEAPWithGTC.RadiusCertProfile.ValueString()
		}
		var1.PEAPWithGTC = var4
	}
	var0.Protocol = var1
	var0.Retries = state.Retries.ValueInt64()
	var var5 []ieMayhq.ServerObject
	if len(state.Server) != 0 {
		var5 = make([]ieMayhq.ServerObject, 0, len(state.Server))
		for var6Index := range state.Server {
			var6 := state.Server[var6Index]
			var var7 ieMayhq.ServerObject
			var7.IpAddress = var6.IpAddress.ValueString()
			var7.Name = var6.Name.ValueString()
			var7.Port = var6.Port.ValueInt64()
			var7.Secret = var6.Secret.ValueString()
			var5 = append(var5, var7)
		}
	}
	var0.Server = var5
	var0.Timeout = state.Timeout.ValueInt64()
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
	var var8 *radiusServerProfilesRsModelProtocolObject
	if ans.Protocol != nil {
		var8 = &radiusServerProfilesRsModelProtocolObject{}
		var var9 *radiusServerProfilesRsModelEAPTTLSWithPAPObject
		if ans.Protocol.EAPTTLSWithPAP != nil {
			var9 = &radiusServerProfilesRsModelEAPTTLSWithPAPObject{}
			var9.AnonOuterId = types.BoolValue(ans.Protocol.EAPTTLSWithPAP.AnonOuterId)
			var9.RadiusCertProfile = types.StringValue(ans.Protocol.EAPTTLSWithPAP.RadiusCertProfile)
		}
		var var10 *radiusServerProfilesRsModelPEAPMSCHAPv2Object
		if ans.Protocol.PEAPMSCHAPv2 != nil {
			var10 = &radiusServerProfilesRsModelPEAPMSCHAPv2Object{}
			var10.AllowPwdChange = types.BoolValue(ans.Protocol.PEAPMSCHAPv2.AllowPwdChange)
			var10.AnonOuterId = types.BoolValue(ans.Protocol.PEAPMSCHAPv2.AnonOuterId)
			var10.RadiusCertProfile = types.StringValue(ans.Protocol.PEAPMSCHAPv2.RadiusCertProfile)
		}
		var var11 *radiusServerProfilesRsModelPEAPWithGTCObject
		if ans.Protocol.PEAPWithGTC != nil {
			var11 = &radiusServerProfilesRsModelPEAPWithGTCObject{}
			var11.AnonOuterId = types.BoolValue(ans.Protocol.PEAPWithGTC.AnonOuterId)
			var11.RadiusCertProfile = types.StringValue(ans.Protocol.PEAPWithGTC.RadiusCertProfile)
		}
		if ans.Protocol.CHAP != nil {
			var8.CHAP = types.BoolValue(true)
		}
		var8.EAPTTLSWithPAP = var9
		if ans.Protocol.PAP != nil {
			var8.PAP = types.BoolValue(true)
		}
		var8.PEAPMSCHAPv2 = var10
		var8.PEAPWithGTC = var11
	}
	var var12 []radiusServerProfilesRsModelServerObject
	if len(ans.Server) != 0 {
		var12 = make([]radiusServerProfilesRsModelServerObject, 0, len(ans.Server))
		for var13Index := range ans.Server {
			var13 := ans.Server[var13Index]
			var var14 radiusServerProfilesRsModelServerObject
			var14.IpAddress = types.StringValue(var13.IpAddress)
			var14.Name = types.StringValue(var13.Name)
			var14.Port = types.Int64Value(var13.Port)
			var14.Secret = types.StringValue(var13.Secret)
			var12 = append(var12, var14)
		}
	}
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Protocol = var8
	state.Retries = types.Int64Value(ans.Retries)
	state.Server = var12
	state.Timeout = types.Int64Value(ans.Timeout)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *radiusServerProfilesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state radiusServerProfilesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "sase_radius_server_profiles",
		"locMap":        map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":        tokens,
	})

	// Prepare to read the config.
	svc := bVmbuOb.NewClient(r.client)
	input := bVmbuOb.ReadInput{
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
	var var0 *radiusServerProfilesRsModelProtocolObject
	if ans.Protocol != nil {
		var0 = &radiusServerProfilesRsModelProtocolObject{}
		var var1 *radiusServerProfilesRsModelEAPTTLSWithPAPObject
		if ans.Protocol.EAPTTLSWithPAP != nil {
			var1 = &radiusServerProfilesRsModelEAPTTLSWithPAPObject{}
			var1.AnonOuterId = types.BoolValue(ans.Protocol.EAPTTLSWithPAP.AnonOuterId)
			var1.RadiusCertProfile = types.StringValue(ans.Protocol.EAPTTLSWithPAP.RadiusCertProfile)
		}
		var var2 *radiusServerProfilesRsModelPEAPMSCHAPv2Object
		if ans.Protocol.PEAPMSCHAPv2 != nil {
			var2 = &radiusServerProfilesRsModelPEAPMSCHAPv2Object{}
			var2.AllowPwdChange = types.BoolValue(ans.Protocol.PEAPMSCHAPv2.AllowPwdChange)
			var2.AnonOuterId = types.BoolValue(ans.Protocol.PEAPMSCHAPv2.AnonOuterId)
			var2.RadiusCertProfile = types.StringValue(ans.Protocol.PEAPMSCHAPv2.RadiusCertProfile)
		}
		var var3 *radiusServerProfilesRsModelPEAPWithGTCObject
		if ans.Protocol.PEAPWithGTC != nil {
			var3 = &radiusServerProfilesRsModelPEAPWithGTCObject{}
			var3.AnonOuterId = types.BoolValue(ans.Protocol.PEAPWithGTC.AnonOuterId)
			var3.RadiusCertProfile = types.StringValue(ans.Protocol.PEAPWithGTC.RadiusCertProfile)
		}
		if ans.Protocol.CHAP != nil {
			var0.CHAP = types.BoolValue(true)
		}
		var0.EAPTTLSWithPAP = var1
		if ans.Protocol.PAP != nil {
			var0.PAP = types.BoolValue(true)
		}
		var0.PEAPMSCHAPv2 = var2
		var0.PEAPWithGTC = var3
	}
	var var4 []radiusServerProfilesRsModelServerObject
	if len(ans.Server) != 0 {
		var4 = make([]radiusServerProfilesRsModelServerObject, 0, len(ans.Server))
		for var5Index := range ans.Server {
			var5 := ans.Server[var5Index]
			var var6 radiusServerProfilesRsModelServerObject
			var6.IpAddress = types.StringValue(var5.IpAddress)
			var6.Name = types.StringValue(var5.Name)
			var6.Port = types.Int64Value(var5.Port)
			var6.Secret = types.StringValue(var5.Secret)
			var4 = append(var4, var6)
		}
	}
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Protocol = var0
	state.Retries = types.Int64Value(ans.Retries)
	state.Server = var4
	state.Timeout = types.Int64Value(ans.Timeout)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *radiusServerProfilesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state radiusServerProfilesRsModel
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
		"resource_name": "sase_radius_server_profiles",
		"object_id":     state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := bVmbuOb.NewClient(r.client)
	input := bVmbuOb.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 ieMayhq.Config
	var var1 *ieMayhq.ProtocolObject
	if plan.Protocol != nil {
		var1 = &ieMayhq.ProtocolObject{}
		if plan.Protocol.CHAP.ValueBool() {
			var1.CHAP = struct{}{}
		}
		var var2 *ieMayhq.EAPTTLSWithPAPObject
		if plan.Protocol.EAPTTLSWithPAP != nil {
			var2 = &ieMayhq.EAPTTLSWithPAPObject{}
			var2.AnonOuterId = plan.Protocol.EAPTTLSWithPAP.AnonOuterId.ValueBool()
			var2.RadiusCertProfile = plan.Protocol.EAPTTLSWithPAP.RadiusCertProfile.ValueString()
		}
		var1.EAPTTLSWithPAP = var2
		if plan.Protocol.PAP.ValueBool() {
			var1.PAP = struct{}{}
		}
		var var3 *ieMayhq.PEAPMSCHAPv2Object
		if plan.Protocol.PEAPMSCHAPv2 != nil {
			var3 = &ieMayhq.PEAPMSCHAPv2Object{}
			var3.AllowPwdChange = plan.Protocol.PEAPMSCHAPv2.AllowPwdChange.ValueBool()
			var3.AnonOuterId = plan.Protocol.PEAPMSCHAPv2.AnonOuterId.ValueBool()
			var3.RadiusCertProfile = plan.Protocol.PEAPMSCHAPv2.RadiusCertProfile.ValueString()
		}
		var1.PEAPMSCHAPv2 = var3
		var var4 *ieMayhq.PEAPWithGTCObject
		if plan.Protocol.PEAPWithGTC != nil {
			var4 = &ieMayhq.PEAPWithGTCObject{}
			var4.AnonOuterId = plan.Protocol.PEAPWithGTC.AnonOuterId.ValueBool()
			var4.RadiusCertProfile = plan.Protocol.PEAPWithGTC.RadiusCertProfile.ValueString()
		}
		var1.PEAPWithGTC = var4
	}
	var0.Protocol = var1
	var0.Retries = plan.Retries.ValueInt64()
	var var5 []ieMayhq.ServerObject
	if len(plan.Server) != 0 {
		var5 = make([]ieMayhq.ServerObject, 0, len(plan.Server))
		for var6Index := range plan.Server {
			var6 := plan.Server[var6Index]
			var var7 ieMayhq.ServerObject
			var7.IpAddress = var6.IpAddress.ValueString()
			var7.Name = var6.Name.ValueString()
			var7.Port = var6.Port.ValueInt64()
			var7.Secret = var6.Secret.ValueString()
			var5 = append(var5, var7)
		}
	}
	var0.Server = var5
	var0.Timeout = plan.Timeout.ValueInt64()
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var8 *radiusServerProfilesRsModelProtocolObject
	if ans.Protocol != nil {
		var8 = &radiusServerProfilesRsModelProtocolObject{}
		var var9 *radiusServerProfilesRsModelEAPTTLSWithPAPObject
		if ans.Protocol.EAPTTLSWithPAP != nil {
			var9 = &radiusServerProfilesRsModelEAPTTLSWithPAPObject{}
			var9.AnonOuterId = types.BoolValue(ans.Protocol.EAPTTLSWithPAP.AnonOuterId)
			var9.RadiusCertProfile = types.StringValue(ans.Protocol.EAPTTLSWithPAP.RadiusCertProfile)
		}
		var var10 *radiusServerProfilesRsModelPEAPMSCHAPv2Object
		if ans.Protocol.PEAPMSCHAPv2 != nil {
			var10 = &radiusServerProfilesRsModelPEAPMSCHAPv2Object{}
			var10.AllowPwdChange = types.BoolValue(ans.Protocol.PEAPMSCHAPv2.AllowPwdChange)
			var10.AnonOuterId = types.BoolValue(ans.Protocol.PEAPMSCHAPv2.AnonOuterId)
			var10.RadiusCertProfile = types.StringValue(ans.Protocol.PEAPMSCHAPv2.RadiusCertProfile)
		}
		var var11 *radiusServerProfilesRsModelPEAPWithGTCObject
		if ans.Protocol.PEAPWithGTC != nil {
			var11 = &radiusServerProfilesRsModelPEAPWithGTCObject{}
			var11.AnonOuterId = types.BoolValue(ans.Protocol.PEAPWithGTC.AnonOuterId)
			var11.RadiusCertProfile = types.StringValue(ans.Protocol.PEAPWithGTC.RadiusCertProfile)
		}
		if ans.Protocol.CHAP != nil {
			var8.CHAP = types.BoolValue(true)
		}
		var8.EAPTTLSWithPAP = var9
		if ans.Protocol.PAP != nil {
			var8.PAP = types.BoolValue(true)
		}
		var8.PEAPMSCHAPv2 = var10
		var8.PEAPWithGTC = var11
	}
	var var12 []radiusServerProfilesRsModelServerObject
	if len(ans.Server) != 0 {
		var12 = make([]radiusServerProfilesRsModelServerObject, 0, len(ans.Server))
		for var13Index := range ans.Server {
			var13 := ans.Server[var13Index]
			var var14 radiusServerProfilesRsModelServerObject
			var14.IpAddress = types.StringValue(var13.IpAddress)
			var14.Name = types.StringValue(var13.Name)
			var14.Port = types.Int64Value(var13.Port)
			var14.Secret = types.StringValue(var13.Secret)
			var12 = append(var12, var14)
		}
	}
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Protocol = var8
	state.Retries = types.Int64Value(ans.Retries)
	state.Server = var12
	state.Timeout = types.Int64Value(ans.Timeout)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *radiusServerProfilesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name": "sase_radius_server_profiles",
		"locMap":        map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":        tokens,
	})

	svc := bVmbuOb.NewClient(r.client)
	input := bVmbuOb.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *radiusServerProfilesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
