package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	gADaUcy "github.com/paloaltonetworks/sase-go/netsec/schema/tls/service/profiles"
	qUVHRkq "github.com/paloaltonetworks/sase-go/netsec/service/v1/tlsserviceprofiles"

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
	_ datasource.DataSource              = &tlsServiceProfilesListDataSource{}
	_ datasource.DataSourceWithConfigure = &tlsServiceProfilesListDataSource{}
)

func NewTlsServiceProfilesListDataSource() datasource.DataSource {
	return &tlsServiceProfilesListDataSource{}
}

type tlsServiceProfilesListDataSource struct {
	client *sase.Client
}

type tlsServiceProfilesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []tlsServiceProfilesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type tlsServiceProfilesListDsModelConfig struct {
	Certificate      types.String                                        `tfsdk:"certificate"`
	ObjectId         types.String                                        `tfsdk:"object_id"`
	Name             types.String                                        `tfsdk:"name"`
	ProtocolSettings tlsServiceProfilesListDsModelProtocolSettingsObject `tfsdk:"protocol_settings"`
}

type tlsServiceProfilesListDsModelProtocolSettingsObject struct {
	AuthAlgoSha1     types.Bool   `tfsdk:"auth_algo_sha1"`
	AuthAlgoSha256   types.Bool   `tfsdk:"auth_algo_sha256"`
	AuthAlgoSha384   types.Bool   `tfsdk:"auth_algo_sha384"`
	EncAlgo3des      types.Bool   `tfsdk:"enc_algo3des"`
	EncAlgoAes128Cbc types.Bool   `tfsdk:"enc_algo_aes128_cbc"`
	EncAlgoAes128Gcm types.Bool   `tfsdk:"enc_algo_aes128_gcm"`
	EncAlgoAes256Cbc types.Bool   `tfsdk:"enc_algo_aes256_cbc"`
	EncAlgoAes256Gcm types.Bool   `tfsdk:"enc_algo_aes256_gcm"`
	EncAlgoRc4       types.Bool   `tfsdk:"enc_algo_rc4"`
	KeyxchgAlgoDhe   types.Bool   `tfsdk:"keyxchg_algo_dhe"`
	KeyxchgAlgoEcdhe types.Bool   `tfsdk:"keyxchg_algo_ecdhe"`
	KeyxchgAlgoRsa   types.Bool   `tfsdk:"keyxchg_algo_rsa"`
	MaxVersion       types.String `tfsdk:"max_version"`
	MinVersion       types.String `tfsdk:"min_version"`
}

// Metadata returns the data source type name.
func (d *tlsServiceProfilesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tls_service_profiles_list"
}

// Schema defines the schema for this listing data source.
func (d *tlsServiceProfilesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"certificate": dsschema.StringAttribute{
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
						"protocol_settings": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"auth_algo_sha1": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"auth_algo_sha256": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"auth_algo_sha384": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"enc_algo3des": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"enc_algo_aes128_cbc": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"enc_algo_aes128_gcm": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"enc_algo_aes256_cbc": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"enc_algo_aes256_gcm": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"enc_algo_rc4": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"keyxchg_algo_dhe": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"keyxchg_algo_ecdhe": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"keyxchg_algo_rsa": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"max_version": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
								"min_version": dsschema.StringAttribute{
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
func (d *tlsServiceProfilesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *tlsServiceProfilesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state tlsServiceProfilesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_tls_service_profiles_list",
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
	svc := qUVHRkq.NewClient(d.client)
	input := qUVHRkq.ListInput{
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
	var var0 []tlsServiceProfilesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]tlsServiceProfilesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 tlsServiceProfilesListDsModelConfig
			var var3 tlsServiceProfilesListDsModelProtocolSettingsObject
			var3.AuthAlgoSha1 = types.BoolValue(var1.ProtocolSettings.AuthAlgoSha1)
			var3.AuthAlgoSha256 = types.BoolValue(var1.ProtocolSettings.AuthAlgoSha256)
			var3.AuthAlgoSha384 = types.BoolValue(var1.ProtocolSettings.AuthAlgoSha384)
			var3.EncAlgo3des = types.BoolValue(var1.ProtocolSettings.EncAlgo3des)
			var3.EncAlgoAes128Cbc = types.BoolValue(var1.ProtocolSettings.EncAlgoAes128Cbc)
			var3.EncAlgoAes128Gcm = types.BoolValue(var1.ProtocolSettings.EncAlgoAes128Gcm)
			var3.EncAlgoAes256Cbc = types.BoolValue(var1.ProtocolSettings.EncAlgoAes256Cbc)
			var3.EncAlgoAes256Gcm = types.BoolValue(var1.ProtocolSettings.EncAlgoAes256Gcm)
			var3.EncAlgoRc4 = types.BoolValue(var1.ProtocolSettings.EncAlgoRc4)
			var3.KeyxchgAlgoDhe = types.BoolValue(var1.ProtocolSettings.KeyxchgAlgoDhe)
			var3.KeyxchgAlgoEcdhe = types.BoolValue(var1.ProtocolSettings.KeyxchgAlgoEcdhe)
			var3.KeyxchgAlgoRsa = types.BoolValue(var1.ProtocolSettings.KeyxchgAlgoRsa)
			var3.MaxVersion = types.StringValue(var1.ProtocolSettings.MaxVersion)
			var3.MinVersion = types.StringValue(var1.ProtocolSettings.MinVersion)
			var2.Certificate = types.StringValue(var1.Certificate)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Name = types.StringValue(var1.Name)
			var2.ProtocolSettings = var3
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
	_ datasource.DataSource              = &tlsServiceProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &tlsServiceProfilesDataSource{}
)

func NewTlsServiceProfilesDataSource() datasource.DataSource {
	return &tlsServiceProfilesDataSource{}
}

type tlsServiceProfilesDataSource struct {
	client *sase.Client
}

type tlsServiceProfilesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`

	// Output.
	// Ref: #/components/schemas/tls-service-profiles
	Certificate types.String `tfsdk:"certificate"`
	// input omit: ObjectId
	Name             types.String                                    `tfsdk:"name"`
	ProtocolSettings tlsServiceProfilesDsModelProtocolSettingsObject `tfsdk:"protocol_settings"`
}

type tlsServiceProfilesDsModelProtocolSettingsObject struct {
	AuthAlgoSha1     types.Bool   `tfsdk:"auth_algo_sha1"`
	AuthAlgoSha256   types.Bool   `tfsdk:"auth_algo_sha256"`
	AuthAlgoSha384   types.Bool   `tfsdk:"auth_algo_sha384"`
	EncAlgo3des      types.Bool   `tfsdk:"enc_algo3des"`
	EncAlgoAes128Cbc types.Bool   `tfsdk:"enc_algo_aes128_cbc"`
	EncAlgoAes128Gcm types.Bool   `tfsdk:"enc_algo_aes128_gcm"`
	EncAlgoAes256Cbc types.Bool   `tfsdk:"enc_algo_aes256_cbc"`
	EncAlgoAes256Gcm types.Bool   `tfsdk:"enc_algo_aes256_gcm"`
	EncAlgoRc4       types.Bool   `tfsdk:"enc_algo_rc4"`
	KeyxchgAlgoDhe   types.Bool   `tfsdk:"keyxchg_algo_dhe"`
	KeyxchgAlgoEcdhe types.Bool   `tfsdk:"keyxchg_algo_ecdhe"`
	KeyxchgAlgoRsa   types.Bool   `tfsdk:"keyxchg_algo_rsa"`
	MaxVersion       types.String `tfsdk:"max_version"`
	MinVersion       types.String `tfsdk:"min_version"`
}

// Metadata returns the data source type name.
func (d *tlsServiceProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tls_service_profiles"
}

// Schema defines the schema for this listing data source.
func (d *tlsServiceProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"certificate": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"protocol_settings": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"auth_algo_sha1": dsschema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"auth_algo_sha256": dsschema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"auth_algo_sha384": dsschema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"enc_algo3des": dsschema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"enc_algo_aes128_cbc": dsschema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"enc_algo_aes128_gcm": dsschema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"enc_algo_aes256_cbc": dsschema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"enc_algo_aes256_gcm": dsschema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"enc_algo_rc4": dsschema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"keyxchg_algo_dhe": dsschema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"keyxchg_algo_ecdhe": dsschema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"keyxchg_algo_rsa": dsschema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"max_version": dsschema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"min_version": dsschema.StringAttribute{
						Description: "",
						Computed:    true,
					},
				},
			},
		},
	}
}

// Configure prepares the struct.
func (d *tlsServiceProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *tlsServiceProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state tlsServiceProfilesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_tls_service_profiles",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to run the command.
	svc := qUVHRkq.NewClient(d.client)
	input := qUVHRkq.ReadInput{
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
	var var0 tlsServiceProfilesDsModelProtocolSettingsObject
	var0.AuthAlgoSha1 = types.BoolValue(ans.ProtocolSettings.AuthAlgoSha1)
	var0.AuthAlgoSha256 = types.BoolValue(ans.ProtocolSettings.AuthAlgoSha256)
	var0.AuthAlgoSha384 = types.BoolValue(ans.ProtocolSettings.AuthAlgoSha384)
	var0.EncAlgo3des = types.BoolValue(ans.ProtocolSettings.EncAlgo3des)
	var0.EncAlgoAes128Cbc = types.BoolValue(ans.ProtocolSettings.EncAlgoAes128Cbc)
	var0.EncAlgoAes128Gcm = types.BoolValue(ans.ProtocolSettings.EncAlgoAes128Gcm)
	var0.EncAlgoAes256Cbc = types.BoolValue(ans.ProtocolSettings.EncAlgoAes256Cbc)
	var0.EncAlgoAes256Gcm = types.BoolValue(ans.ProtocolSettings.EncAlgoAes256Gcm)
	var0.EncAlgoRc4 = types.BoolValue(ans.ProtocolSettings.EncAlgoRc4)
	var0.KeyxchgAlgoDhe = types.BoolValue(ans.ProtocolSettings.KeyxchgAlgoDhe)
	var0.KeyxchgAlgoEcdhe = types.BoolValue(ans.ProtocolSettings.KeyxchgAlgoEcdhe)
	var0.KeyxchgAlgoRsa = types.BoolValue(ans.ProtocolSettings.KeyxchgAlgoRsa)
	var0.MaxVersion = types.StringValue(ans.ProtocolSettings.MaxVersion)
	var0.MinVersion = types.StringValue(ans.ProtocolSettings.MinVersion)
	state.Certificate = types.StringValue(ans.Certificate)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.ProtocolSettings = var0

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &tlsServiceProfilesResource{}
	_ resource.ResourceWithConfigure   = &tlsServiceProfilesResource{}
	_ resource.ResourceWithImportState = &tlsServiceProfilesResource{}
)

func NewTlsServiceProfilesResource() resource.Resource {
	return &tlsServiceProfilesResource{}
}

type tlsServiceProfilesResource struct {
	client *sase.Client
}

type tlsServiceProfilesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/tls-service-profiles
	Certificate      types.String                                    `tfsdk:"certificate"`
	ObjectId         types.String                                    `tfsdk:"object_id"`
	Name             types.String                                    `tfsdk:"name"`
	ProtocolSettings tlsServiceProfilesRsModelProtocolSettingsObject `tfsdk:"protocol_settings"`
}

type tlsServiceProfilesRsModelProtocolSettingsObject struct {
	AuthAlgoSha1     types.Bool   `tfsdk:"auth_algo_sha1"`
	AuthAlgoSha256   types.Bool   `tfsdk:"auth_algo_sha256"`
	AuthAlgoSha384   types.Bool   `tfsdk:"auth_algo_sha384"`
	EncAlgo3des      types.Bool   `tfsdk:"enc_algo3des"`
	EncAlgoAes128Cbc types.Bool   `tfsdk:"enc_algo_aes128_cbc"`
	EncAlgoAes128Gcm types.Bool   `tfsdk:"enc_algo_aes128_gcm"`
	EncAlgoAes256Cbc types.Bool   `tfsdk:"enc_algo_aes256_cbc"`
	EncAlgoAes256Gcm types.Bool   `tfsdk:"enc_algo_aes256_gcm"`
	EncAlgoRc4       types.Bool   `tfsdk:"enc_algo_rc4"`
	KeyxchgAlgoDhe   types.Bool   `tfsdk:"keyxchg_algo_dhe"`
	KeyxchgAlgoEcdhe types.Bool   `tfsdk:"keyxchg_algo_ecdhe"`
	KeyxchgAlgoRsa   types.Bool   `tfsdk:"keyxchg_algo_rsa"`
	MaxVersion       types.String `tfsdk:"max_version"`
	MinVersion       types.String `tfsdk:"min_version"`
}

// Metadata returns the data source type name.
func (r *tlsServiceProfilesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tls_service_profiles"
}

// Schema defines the schema for this listing data source.
func (r *tlsServiceProfilesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

			"certificate": rsschema.StringAttribute{
				Description: "",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
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
					stringvalidator.LengthAtMost(127),
				},
			},
			"protocol_settings": rsschema.SingleNestedAttribute{
				Description: "",
				Required:    true,
				Attributes: map[string]rsschema.Attribute{
					"auth_algo_sha1": rsschema.BoolAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"auth_algo_sha256": rsschema.BoolAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"auth_algo_sha384": rsschema.BoolAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"enc_algo3des": rsschema.BoolAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"enc_algo_aes128_cbc": rsschema.BoolAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"enc_algo_aes128_gcm": rsschema.BoolAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"enc_algo_aes256_cbc": rsschema.BoolAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"enc_algo_aes256_gcm": rsschema.BoolAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"enc_algo_rc4": rsschema.BoolAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"keyxchg_algo_dhe": rsschema.BoolAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"keyxchg_algo_ecdhe": rsschema.BoolAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"keyxchg_algo_rsa": rsschema.BoolAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"max_version": rsschema.StringAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							DefaultString("max"),
						},
						Validators: []validator.String{
							stringvalidator.OneOf("tls1-0", "tls1-1", "tls1-2", "tls1-3", "max"),
						},
					},
					"min_version": rsschema.StringAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							DefaultString("tls1-0"),
						},
						Validators: []validator.String{
							stringvalidator.OneOf("tls1-0", "tls1-1", "tls1-2"),
						},
					},
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *tlsServiceProfilesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *tlsServiceProfilesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state tlsServiceProfilesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_tls_service_profiles",
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := qUVHRkq.NewClient(r.client)
	input := qUVHRkq.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 gADaUcy.Config
	var0.Certificate = state.Certificate.ValueString()
	var0.Name = state.Name.ValueString()
	var var1 gADaUcy.ProtocolSettingsObject
	var1.AuthAlgoSha1 = state.ProtocolSettings.AuthAlgoSha1.ValueBool()
	var1.AuthAlgoSha256 = state.ProtocolSettings.AuthAlgoSha256.ValueBool()
	var1.AuthAlgoSha384 = state.ProtocolSettings.AuthAlgoSha384.ValueBool()
	var1.EncAlgo3des = state.ProtocolSettings.EncAlgo3des.ValueBool()
	var1.EncAlgoAes128Cbc = state.ProtocolSettings.EncAlgoAes128Cbc.ValueBool()
	var1.EncAlgoAes128Gcm = state.ProtocolSettings.EncAlgoAes128Gcm.ValueBool()
	var1.EncAlgoAes256Cbc = state.ProtocolSettings.EncAlgoAes256Cbc.ValueBool()
	var1.EncAlgoAes256Gcm = state.ProtocolSettings.EncAlgoAes256Gcm.ValueBool()
	var1.EncAlgoRc4 = state.ProtocolSettings.EncAlgoRc4.ValueBool()
	var1.KeyxchgAlgoDhe = state.ProtocolSettings.KeyxchgAlgoDhe.ValueBool()
	var1.KeyxchgAlgoEcdhe = state.ProtocolSettings.KeyxchgAlgoEcdhe.ValueBool()
	var1.KeyxchgAlgoRsa = state.ProtocolSettings.KeyxchgAlgoRsa.ValueBool()
	var1.MaxVersion = state.ProtocolSettings.MaxVersion.ValueString()
	var1.MinVersion = state.ProtocolSettings.MinVersion.ValueString()
	var0.ProtocolSettings = var1
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
	var var2 tlsServiceProfilesRsModelProtocolSettingsObject
	var2.AuthAlgoSha1 = types.BoolValue(ans.ProtocolSettings.AuthAlgoSha1)
	var2.AuthAlgoSha256 = types.BoolValue(ans.ProtocolSettings.AuthAlgoSha256)
	var2.AuthAlgoSha384 = types.BoolValue(ans.ProtocolSettings.AuthAlgoSha384)
	var2.EncAlgo3des = types.BoolValue(ans.ProtocolSettings.EncAlgo3des)
	var2.EncAlgoAes128Cbc = types.BoolValue(ans.ProtocolSettings.EncAlgoAes128Cbc)
	var2.EncAlgoAes128Gcm = types.BoolValue(ans.ProtocolSettings.EncAlgoAes128Gcm)
	var2.EncAlgoAes256Cbc = types.BoolValue(ans.ProtocolSettings.EncAlgoAes256Cbc)
	var2.EncAlgoAes256Gcm = types.BoolValue(ans.ProtocolSettings.EncAlgoAes256Gcm)
	var2.EncAlgoRc4 = types.BoolValue(ans.ProtocolSettings.EncAlgoRc4)
	var2.KeyxchgAlgoDhe = types.BoolValue(ans.ProtocolSettings.KeyxchgAlgoDhe)
	var2.KeyxchgAlgoEcdhe = types.BoolValue(ans.ProtocolSettings.KeyxchgAlgoEcdhe)
	var2.KeyxchgAlgoRsa = types.BoolValue(ans.ProtocolSettings.KeyxchgAlgoRsa)
	var2.MaxVersion = types.StringValue(ans.ProtocolSettings.MaxVersion)
	var2.MinVersion = types.StringValue(ans.ProtocolSettings.MinVersion)
	state.Certificate = types.StringValue(ans.Certificate)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.ProtocolSettings = var2

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *tlsServiceProfilesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state tlsServiceProfilesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_tls_service_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := qUVHRkq.NewClient(r.client)
	input := qUVHRkq.ReadInput{
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
	var var0 tlsServiceProfilesRsModelProtocolSettingsObject
	var0.AuthAlgoSha1 = types.BoolValue(ans.ProtocolSettings.AuthAlgoSha1)
	var0.AuthAlgoSha256 = types.BoolValue(ans.ProtocolSettings.AuthAlgoSha256)
	var0.AuthAlgoSha384 = types.BoolValue(ans.ProtocolSettings.AuthAlgoSha384)
	var0.EncAlgo3des = types.BoolValue(ans.ProtocolSettings.EncAlgo3des)
	var0.EncAlgoAes128Cbc = types.BoolValue(ans.ProtocolSettings.EncAlgoAes128Cbc)
	var0.EncAlgoAes128Gcm = types.BoolValue(ans.ProtocolSettings.EncAlgoAes128Gcm)
	var0.EncAlgoAes256Cbc = types.BoolValue(ans.ProtocolSettings.EncAlgoAes256Cbc)
	var0.EncAlgoAes256Gcm = types.BoolValue(ans.ProtocolSettings.EncAlgoAes256Gcm)
	var0.EncAlgoRc4 = types.BoolValue(ans.ProtocolSettings.EncAlgoRc4)
	var0.KeyxchgAlgoDhe = types.BoolValue(ans.ProtocolSettings.KeyxchgAlgoDhe)
	var0.KeyxchgAlgoEcdhe = types.BoolValue(ans.ProtocolSettings.KeyxchgAlgoEcdhe)
	var0.KeyxchgAlgoRsa = types.BoolValue(ans.ProtocolSettings.KeyxchgAlgoRsa)
	var0.MaxVersion = types.StringValue(ans.ProtocolSettings.MaxVersion)
	var0.MinVersion = types.StringValue(ans.ProtocolSettings.MinVersion)
	state.Certificate = types.StringValue(ans.Certificate)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.ProtocolSettings = var0

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *tlsServiceProfilesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state tlsServiceProfilesRsModel
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
		"resource_name":               "sase_tls_service_profiles",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := qUVHRkq.NewClient(r.client)
	input := qUVHRkq.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 gADaUcy.Config
	var0.Certificate = plan.Certificate.ValueString()
	var0.Name = plan.Name.ValueString()
	var var1 gADaUcy.ProtocolSettingsObject
	var1.AuthAlgoSha1 = plan.ProtocolSettings.AuthAlgoSha1.ValueBool()
	var1.AuthAlgoSha256 = plan.ProtocolSettings.AuthAlgoSha256.ValueBool()
	var1.AuthAlgoSha384 = plan.ProtocolSettings.AuthAlgoSha384.ValueBool()
	var1.EncAlgo3des = plan.ProtocolSettings.EncAlgo3des.ValueBool()
	var1.EncAlgoAes128Cbc = plan.ProtocolSettings.EncAlgoAes128Cbc.ValueBool()
	var1.EncAlgoAes128Gcm = plan.ProtocolSettings.EncAlgoAes128Gcm.ValueBool()
	var1.EncAlgoAes256Cbc = plan.ProtocolSettings.EncAlgoAes256Cbc.ValueBool()
	var1.EncAlgoAes256Gcm = plan.ProtocolSettings.EncAlgoAes256Gcm.ValueBool()
	var1.EncAlgoRc4 = plan.ProtocolSettings.EncAlgoRc4.ValueBool()
	var1.KeyxchgAlgoDhe = plan.ProtocolSettings.KeyxchgAlgoDhe.ValueBool()
	var1.KeyxchgAlgoEcdhe = plan.ProtocolSettings.KeyxchgAlgoEcdhe.ValueBool()
	var1.KeyxchgAlgoRsa = plan.ProtocolSettings.KeyxchgAlgoRsa.ValueBool()
	var1.MaxVersion = plan.ProtocolSettings.MaxVersion.ValueString()
	var1.MinVersion = plan.ProtocolSettings.MinVersion.ValueString()
	var0.ProtocolSettings = var1
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var2 tlsServiceProfilesRsModelProtocolSettingsObject
	var2.AuthAlgoSha1 = types.BoolValue(ans.ProtocolSettings.AuthAlgoSha1)
	var2.AuthAlgoSha256 = types.BoolValue(ans.ProtocolSettings.AuthAlgoSha256)
	var2.AuthAlgoSha384 = types.BoolValue(ans.ProtocolSettings.AuthAlgoSha384)
	var2.EncAlgo3des = types.BoolValue(ans.ProtocolSettings.EncAlgo3des)
	var2.EncAlgoAes128Cbc = types.BoolValue(ans.ProtocolSettings.EncAlgoAes128Cbc)
	var2.EncAlgoAes128Gcm = types.BoolValue(ans.ProtocolSettings.EncAlgoAes128Gcm)
	var2.EncAlgoAes256Cbc = types.BoolValue(ans.ProtocolSettings.EncAlgoAes256Cbc)
	var2.EncAlgoAes256Gcm = types.BoolValue(ans.ProtocolSettings.EncAlgoAes256Gcm)
	var2.EncAlgoRc4 = types.BoolValue(ans.ProtocolSettings.EncAlgoRc4)
	var2.KeyxchgAlgoDhe = types.BoolValue(ans.ProtocolSettings.KeyxchgAlgoDhe)
	var2.KeyxchgAlgoEcdhe = types.BoolValue(ans.ProtocolSettings.KeyxchgAlgoEcdhe)
	var2.KeyxchgAlgoRsa = types.BoolValue(ans.ProtocolSettings.KeyxchgAlgoRsa)
	var2.MaxVersion = types.StringValue(ans.ProtocolSettings.MaxVersion)
	var2.MinVersion = types.StringValue(ans.ProtocolSettings.MinVersion)
	state.Certificate = types.StringValue(ans.Certificate)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.ProtocolSettings = var2

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *tlsServiceProfilesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name":               "sase_tls_service_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	svc := qUVHRkq.NewClient(r.client)
	input := qUVHRkq.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *tlsServiceProfilesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
