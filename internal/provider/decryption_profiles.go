package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	vMYBRZK "github.com/paloaltonetworks/sase-go/netsec/schema/decryption/profiles"
	bpgvUeD "github.com/paloaltonetworks/sase-go/netsec/service/v1/decryptionprofiles"

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
	_ datasource.DataSource              = &decryptionProfilesListDataSource{}
	_ datasource.DataSourceWithConfigure = &decryptionProfilesListDataSource{}
)

func NewDecryptionProfilesListDataSource() datasource.DataSource {
	return &decryptionProfilesListDataSource{}
}

type decryptionProfilesListDataSource struct {
	client *sase.Client
}

type decryptionProfilesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []decryptionProfilesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type decryptionProfilesListDsModelConfig struct {
	ObjectId            types.String                                            `tfsdk:"object_id"`
	Name                types.String                                            `tfsdk:"name"`
	SslForwardProxy     *decryptionProfilesListDsModelSslForwardProxyObject     `tfsdk:"ssl_forward_proxy"`
	SslInboundProxy     *decryptionProfilesListDsModelSslInboundProxyObject     `tfsdk:"ssl_inbound_proxy"`
	SslNoProxy          *decryptionProfilesListDsModelSslNoProxyObject          `tfsdk:"ssl_no_proxy"`
	SslProtocolSettings *decryptionProfilesListDsModelSslProtocolSettingsObject `tfsdk:"ssl_protocol_settings"`
}

type decryptionProfilesListDsModelSslForwardProxyObject struct {
	AutoIncludeAltname            types.Bool `tfsdk:"auto_include_altname"`
	BlockClientCert               types.Bool `tfsdk:"block_client_cert"`
	BlockExpiredCertificate       types.Bool `tfsdk:"block_expired_certificate"`
	BlockTimeoutCert              types.Bool `tfsdk:"block_timeout_cert"`
	BlockTls13DowngradeNoResource types.Bool `tfsdk:"block_tls13_downgrade_no_resource"`
	BlockUnknownCert              types.Bool `tfsdk:"block_unknown_cert"`
	BlockUnsupportedCipher        types.Bool `tfsdk:"block_unsupported_cipher"`
	BlockUnsupportedVersion       types.Bool `tfsdk:"block_unsupported_version"`
	BlockUntrustedIssuer          types.Bool `tfsdk:"block_untrusted_issuer"`
	RestrictCertExts              types.Bool `tfsdk:"restrict_cert_exts"`
	StripAlpn                     types.Bool `tfsdk:"strip_alpn"`
}

type decryptionProfilesListDsModelSslInboundProxyObject struct {
	BlockIfHsmUnavailable   types.Bool `tfsdk:"block_if_hsm_unavailable"`
	BlockIfNoResource       types.Bool `tfsdk:"block_if_no_resource"`
	BlockUnsupportedCipher  types.Bool `tfsdk:"block_unsupported_cipher"`
	BlockUnsupportedVersion types.Bool `tfsdk:"block_unsupported_version"`
}

type decryptionProfilesListDsModelSslNoProxyObject struct {
	BlockExpiredCertificate types.Bool `tfsdk:"block_expired_certificate"`
	BlockUntrustedIssuer    types.Bool `tfsdk:"block_untrusted_issuer"`
}

type decryptionProfilesListDsModelSslProtocolSettingsObject struct {
	AuthAlgoMd5             types.Bool   `tfsdk:"auth_algo_md5"`
	AuthAlgoSha1            types.Bool   `tfsdk:"auth_algo_sha1"`
	AuthAlgoSha256          types.Bool   `tfsdk:"auth_algo_sha256"`
	AuthAlgoSha384          types.Bool   `tfsdk:"auth_algo_sha384"`
	EncAlgo3des             types.Bool   `tfsdk:"enc_algo3des"`
	EncAlgoAes128Cbc        types.Bool   `tfsdk:"enc_algo_aes128_cbc"`
	EncAlgoAes128Gcm        types.Bool   `tfsdk:"enc_algo_aes128_gcm"`
	EncAlgoAes256Cbc        types.Bool   `tfsdk:"enc_algo_aes256_cbc"`
	EncAlgoAes256Gcm        types.Bool   `tfsdk:"enc_algo_aes256_gcm"`
	EncAlgoChacha20Poly1305 types.Bool   `tfsdk:"enc_algo_chacha20_poly1305"`
	EncAlgoRc4              types.Bool   `tfsdk:"enc_algo_rc4"`
	KeyxchgAlgoDhe          types.Bool   `tfsdk:"keyxchg_algo_dhe"`
	KeyxchgAlgoEcdhe        types.Bool   `tfsdk:"keyxchg_algo_ecdhe"`
	KeyxchgAlgoRsa          types.Bool   `tfsdk:"keyxchg_algo_rsa"`
	MaxVersion              types.String `tfsdk:"max_version"`
	MinVersion              types.String `tfsdk:"min_version"`
}

// Metadata returns the data source type name.
func (d *decryptionProfilesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_decryption_profiles_list"
}

// Schema defines the schema for this listing data source.
func (d *decryptionProfilesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"object_id": dsschema.StringAttribute{
							Description:         "The `object_id` parameter.",
							MarkdownDescription: "The `object_id` parameter.",
							Computed:            true,
						},
						"name": dsschema.StringAttribute{
							Description:         "The `name` parameter.",
							MarkdownDescription: "The `name` parameter.",
							Computed:            true,
						},
						"ssl_forward_proxy": dsschema.SingleNestedAttribute{
							Description:         "The `ssl_forward_proxy` parameter.",
							MarkdownDescription: "The `ssl_forward_proxy` parameter.",
							Computed:            true,
							Attributes: map[string]dsschema.Attribute{
								"auto_include_altname": dsschema.BoolAttribute{
									Description:         "The `auto_include_altname` parameter.",
									MarkdownDescription: "The `auto_include_altname` parameter.",
									Computed:            true,
								},
								"block_client_cert": dsschema.BoolAttribute{
									Description:         "The `block_client_cert` parameter.",
									MarkdownDescription: "The `block_client_cert` parameter.",
									Computed:            true,
								},
								"block_expired_certificate": dsschema.BoolAttribute{
									Description:         "The `block_expired_certificate` parameter.",
									MarkdownDescription: "The `block_expired_certificate` parameter.",
									Computed:            true,
								},
								"block_timeout_cert": dsschema.BoolAttribute{
									Description:         "The `block_timeout_cert` parameter.",
									MarkdownDescription: "The `block_timeout_cert` parameter.",
									Computed:            true,
								},
								"block_tls13_downgrade_no_resource": dsschema.BoolAttribute{
									Description:         "The `block_tls13_downgrade_no_resource` parameter.",
									MarkdownDescription: "The `block_tls13_downgrade_no_resource` parameter.",
									Computed:            true,
								},
								"block_unknown_cert": dsschema.BoolAttribute{
									Description:         "The `block_unknown_cert` parameter.",
									MarkdownDescription: "The `block_unknown_cert` parameter.",
									Computed:            true,
								},
								"block_unsupported_cipher": dsschema.BoolAttribute{
									Description:         "The `block_unsupported_cipher` parameter.",
									MarkdownDescription: "The `block_unsupported_cipher` parameter.",
									Computed:            true,
								},
								"block_unsupported_version": dsschema.BoolAttribute{
									Description:         "The `block_unsupported_version` parameter.",
									MarkdownDescription: "The `block_unsupported_version` parameter.",
									Computed:            true,
								},
								"block_untrusted_issuer": dsschema.BoolAttribute{
									Description:         "The `block_untrusted_issuer` parameter.",
									MarkdownDescription: "The `block_untrusted_issuer` parameter.",
									Computed:            true,
								},
								"restrict_cert_exts": dsschema.BoolAttribute{
									Description:         "The `restrict_cert_exts` parameter.",
									MarkdownDescription: "The `restrict_cert_exts` parameter.",
									Computed:            true,
								},
								"strip_alpn": dsschema.BoolAttribute{
									Description:         "The `strip_alpn` parameter.",
									MarkdownDescription: "The `strip_alpn` parameter.",
									Computed:            true,
								},
							},
						},
						"ssl_inbound_proxy": dsschema.SingleNestedAttribute{
							Description:         "The `ssl_inbound_proxy` parameter.",
							MarkdownDescription: "The `ssl_inbound_proxy` parameter.",
							Computed:            true,
							Attributes: map[string]dsschema.Attribute{
								"block_if_hsm_unavailable": dsschema.BoolAttribute{
									Description:         "The `block_if_hsm_unavailable` parameter.",
									MarkdownDescription: "The `block_if_hsm_unavailable` parameter.",
									Computed:            true,
								},
								"block_if_no_resource": dsschema.BoolAttribute{
									Description:         "The `block_if_no_resource` parameter.",
									MarkdownDescription: "The `block_if_no_resource` parameter.",
									Computed:            true,
								},
								"block_unsupported_cipher": dsschema.BoolAttribute{
									Description:         "The `block_unsupported_cipher` parameter.",
									MarkdownDescription: "The `block_unsupported_cipher` parameter.",
									Computed:            true,
								},
								"block_unsupported_version": dsschema.BoolAttribute{
									Description:         "The `block_unsupported_version` parameter.",
									MarkdownDescription: "The `block_unsupported_version` parameter.",
									Computed:            true,
								},
							},
						},
						"ssl_no_proxy": dsschema.SingleNestedAttribute{
							Description:         "The `ssl_no_proxy` parameter.",
							MarkdownDescription: "The `ssl_no_proxy` parameter.",
							Computed:            true,
							Attributes: map[string]dsschema.Attribute{
								"block_expired_certificate": dsschema.BoolAttribute{
									Description:         "The `block_expired_certificate` parameter.",
									MarkdownDescription: "The `block_expired_certificate` parameter.",
									Computed:            true,
								},
								"block_untrusted_issuer": dsschema.BoolAttribute{
									Description:         "The `block_untrusted_issuer` parameter.",
									MarkdownDescription: "The `block_untrusted_issuer` parameter.",
									Computed:            true,
								},
							},
						},
						"ssl_protocol_settings": dsschema.SingleNestedAttribute{
							Description:         "The `ssl_protocol_settings` parameter.",
							MarkdownDescription: "The `ssl_protocol_settings` parameter.",
							Computed:            true,
							Attributes: map[string]dsschema.Attribute{
								"auth_algo_md5": dsschema.BoolAttribute{
									Description:         "The `auth_algo_md5` parameter.",
									MarkdownDescription: "The `auth_algo_md5` parameter.",
									Computed:            true,
								},
								"auth_algo_sha1": dsschema.BoolAttribute{
									Description:         "The `auth_algo_sha1` parameter.",
									MarkdownDescription: "The `auth_algo_sha1` parameter.",
									Computed:            true,
								},
								"auth_algo_sha256": dsschema.BoolAttribute{
									Description:         "The `auth_algo_sha256` parameter.",
									MarkdownDescription: "The `auth_algo_sha256` parameter.",
									Computed:            true,
								},
								"auth_algo_sha384": dsschema.BoolAttribute{
									Description:         "The `auth_algo_sha384` parameter.",
									MarkdownDescription: "The `auth_algo_sha384` parameter.",
									Computed:            true,
								},
								"enc_algo3des": dsschema.BoolAttribute{
									Description:         "The `enc_algo3des` parameter.",
									MarkdownDescription: "The `enc_algo3des` parameter.",
									Computed:            true,
								},
								"enc_algo_aes128_cbc": dsschema.BoolAttribute{
									Description:         "The `enc_algo_aes128_cbc` parameter.",
									MarkdownDescription: "The `enc_algo_aes128_cbc` parameter.",
									Computed:            true,
								},
								"enc_algo_aes128_gcm": dsschema.BoolAttribute{
									Description:         "The `enc_algo_aes128_gcm` parameter.",
									MarkdownDescription: "The `enc_algo_aes128_gcm` parameter.",
									Computed:            true,
								},
								"enc_algo_aes256_cbc": dsschema.BoolAttribute{
									Description:         "The `enc_algo_aes256_cbc` parameter.",
									MarkdownDescription: "The `enc_algo_aes256_cbc` parameter.",
									Computed:            true,
								},
								"enc_algo_aes256_gcm": dsschema.BoolAttribute{
									Description:         "The `enc_algo_aes256_gcm` parameter.",
									MarkdownDescription: "The `enc_algo_aes256_gcm` parameter.",
									Computed:            true,
								},
								"enc_algo_chacha20_poly1305": dsschema.BoolAttribute{
									Description:         "The `enc_algo_chacha20_poly1305` parameter.",
									MarkdownDescription: "The `enc_algo_chacha20_poly1305` parameter.",
									Computed:            true,
								},
								"enc_algo_rc4": dsschema.BoolAttribute{
									Description:         "The `enc_algo_rc4` parameter.",
									MarkdownDescription: "The `enc_algo_rc4` parameter.",
									Computed:            true,
								},
								"keyxchg_algo_dhe": dsschema.BoolAttribute{
									Description:         "The `keyxchg_algo_dhe` parameter.",
									MarkdownDescription: "The `keyxchg_algo_dhe` parameter.",
									Computed:            true,
								},
								"keyxchg_algo_ecdhe": dsschema.BoolAttribute{
									Description:         "The `keyxchg_algo_ecdhe` parameter.",
									MarkdownDescription: "The `keyxchg_algo_ecdhe` parameter.",
									Computed:            true,
								},
								"keyxchg_algo_rsa": dsschema.BoolAttribute{
									Description:         "The `keyxchg_algo_rsa` parameter.",
									MarkdownDescription: "The `keyxchg_algo_rsa` parameter.",
									Computed:            true,
								},
								"max_version": dsschema.StringAttribute{
									Description:         "The `max_version` parameter.",
									MarkdownDescription: "The `max_version` parameter.",
									Computed:            true,
								},
								"min_version": dsschema.StringAttribute{
									Description:         "The `min_version` parameter.",
									MarkdownDescription: "The `min_version` parameter.",
									Computed:            true,
								},
							},
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
func (d *decryptionProfilesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *decryptionProfilesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state decryptionProfilesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_decryption_profiles_list",
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
	svc := bpgvUeD.NewClient(d.client)
	input := bpgvUeD.ListInput{
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
	var var0 []decryptionProfilesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]decryptionProfilesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 decryptionProfilesListDsModelConfig
			var var3 *decryptionProfilesListDsModelSslForwardProxyObject
			if var1.SslForwardProxy != nil {
				var3 = &decryptionProfilesListDsModelSslForwardProxyObject{}
				var3.AutoIncludeAltname = types.BoolValue(var1.SslForwardProxy.AutoIncludeAltname)
				var3.BlockClientCert = types.BoolValue(var1.SslForwardProxy.BlockClientCert)
				var3.BlockExpiredCertificate = types.BoolValue(var1.SslForwardProxy.BlockExpiredCertificate)
				var3.BlockTimeoutCert = types.BoolValue(var1.SslForwardProxy.BlockTimeoutCert)
				var3.BlockTls13DowngradeNoResource = types.BoolValue(var1.SslForwardProxy.BlockTls13DowngradeNoResource)
				var3.BlockUnknownCert = types.BoolValue(var1.SslForwardProxy.BlockUnknownCert)
				var3.BlockUnsupportedCipher = types.BoolValue(var1.SslForwardProxy.BlockUnsupportedCipher)
				var3.BlockUnsupportedVersion = types.BoolValue(var1.SslForwardProxy.BlockUnsupportedVersion)
				var3.BlockUntrustedIssuer = types.BoolValue(var1.SslForwardProxy.BlockUntrustedIssuer)
				var3.RestrictCertExts = types.BoolValue(var1.SslForwardProxy.RestrictCertExts)
				var3.StripAlpn = types.BoolValue(var1.SslForwardProxy.StripAlpn)
			}
			var var4 *decryptionProfilesListDsModelSslInboundProxyObject
			if var1.SslInboundProxy != nil {
				var4 = &decryptionProfilesListDsModelSslInboundProxyObject{}
				var4.BlockIfHsmUnavailable = types.BoolValue(var1.SslInboundProxy.BlockIfHsmUnavailable)
				var4.BlockIfNoResource = types.BoolValue(var1.SslInboundProxy.BlockIfNoResource)
				var4.BlockUnsupportedCipher = types.BoolValue(var1.SslInboundProxy.BlockUnsupportedCipher)
				var4.BlockUnsupportedVersion = types.BoolValue(var1.SslInboundProxy.BlockUnsupportedVersion)
			}
			var var5 *decryptionProfilesListDsModelSslNoProxyObject
			if var1.SslNoProxy != nil {
				var5 = &decryptionProfilesListDsModelSslNoProxyObject{}
				var5.BlockExpiredCertificate = types.BoolValue(var1.SslNoProxy.BlockExpiredCertificate)
				var5.BlockUntrustedIssuer = types.BoolValue(var1.SslNoProxy.BlockUntrustedIssuer)
			}
			var var6 *decryptionProfilesListDsModelSslProtocolSettingsObject
			if var1.SslProtocolSettings != nil {
				var6 = &decryptionProfilesListDsModelSslProtocolSettingsObject{}
				var6.AuthAlgoMd5 = types.BoolValue(var1.SslProtocolSettings.AuthAlgoMd5)
				var6.AuthAlgoSha1 = types.BoolValue(var1.SslProtocolSettings.AuthAlgoSha1)
				var6.AuthAlgoSha256 = types.BoolValue(var1.SslProtocolSettings.AuthAlgoSha256)
				var6.AuthAlgoSha384 = types.BoolValue(var1.SslProtocolSettings.AuthAlgoSha384)
				var6.EncAlgo3des = types.BoolValue(var1.SslProtocolSettings.EncAlgo3des)
				var6.EncAlgoAes128Cbc = types.BoolValue(var1.SslProtocolSettings.EncAlgoAes128Cbc)
				var6.EncAlgoAes128Gcm = types.BoolValue(var1.SslProtocolSettings.EncAlgoAes128Gcm)
				var6.EncAlgoAes256Cbc = types.BoolValue(var1.SslProtocolSettings.EncAlgoAes256Cbc)
				var6.EncAlgoAes256Gcm = types.BoolValue(var1.SslProtocolSettings.EncAlgoAes256Gcm)
				var6.EncAlgoChacha20Poly1305 = types.BoolValue(var1.SslProtocolSettings.EncAlgoChacha20Poly1305)
				var6.EncAlgoRc4 = types.BoolValue(var1.SslProtocolSettings.EncAlgoRc4)
				var6.KeyxchgAlgoDhe = types.BoolValue(var1.SslProtocolSettings.KeyxchgAlgoDhe)
				var6.KeyxchgAlgoEcdhe = types.BoolValue(var1.SslProtocolSettings.KeyxchgAlgoEcdhe)
				var6.KeyxchgAlgoRsa = types.BoolValue(var1.SslProtocolSettings.KeyxchgAlgoRsa)
				var6.MaxVersion = types.StringValue(var1.SslProtocolSettings.MaxVersion)
				var6.MinVersion = types.StringValue(var1.SslProtocolSettings.MinVersion)
			}
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Name = types.StringValue(var1.Name)
			var2.SslForwardProxy = var3
			var2.SslInboundProxy = var4
			var2.SslNoProxy = var5
			var2.SslProtocolSettings = var6
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
	_ datasource.DataSource              = &decryptionProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &decryptionProfilesDataSource{}
)

func NewDecryptionProfilesDataSource() datasource.DataSource {
	return &decryptionProfilesDataSource{}
}

type decryptionProfilesDataSource struct {
	client *sase.Client
}

type decryptionProfilesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`
	Folder   types.String `tfsdk:"folder"`

	// Output.
	// Ref: #/components/schemas/decryption-profiles
	// input omit: ObjectId
	Name                types.String                                        `tfsdk:"name"`
	SslForwardProxy     *decryptionProfilesDsModelSslForwardProxyObject     `tfsdk:"ssl_forward_proxy"`
	SslInboundProxy     *decryptionProfilesDsModelSslInboundProxyObject     `tfsdk:"ssl_inbound_proxy"`
	SslNoProxy          *decryptionProfilesDsModelSslNoProxyObject          `tfsdk:"ssl_no_proxy"`
	SslProtocolSettings *decryptionProfilesDsModelSslProtocolSettingsObject `tfsdk:"ssl_protocol_settings"`
}

type decryptionProfilesDsModelSslForwardProxyObject struct {
	AutoIncludeAltname            types.Bool `tfsdk:"auto_include_altname"`
	BlockClientCert               types.Bool `tfsdk:"block_client_cert"`
	BlockExpiredCertificate       types.Bool `tfsdk:"block_expired_certificate"`
	BlockTimeoutCert              types.Bool `tfsdk:"block_timeout_cert"`
	BlockTls13DowngradeNoResource types.Bool `tfsdk:"block_tls13_downgrade_no_resource"`
	BlockUnknownCert              types.Bool `tfsdk:"block_unknown_cert"`
	BlockUnsupportedCipher        types.Bool `tfsdk:"block_unsupported_cipher"`
	BlockUnsupportedVersion       types.Bool `tfsdk:"block_unsupported_version"`
	BlockUntrustedIssuer          types.Bool `tfsdk:"block_untrusted_issuer"`
	RestrictCertExts              types.Bool `tfsdk:"restrict_cert_exts"`
	StripAlpn                     types.Bool `tfsdk:"strip_alpn"`
}

type decryptionProfilesDsModelSslInboundProxyObject struct {
	BlockIfHsmUnavailable   types.Bool `tfsdk:"block_if_hsm_unavailable"`
	BlockIfNoResource       types.Bool `tfsdk:"block_if_no_resource"`
	BlockUnsupportedCipher  types.Bool `tfsdk:"block_unsupported_cipher"`
	BlockUnsupportedVersion types.Bool `tfsdk:"block_unsupported_version"`
}

type decryptionProfilesDsModelSslNoProxyObject struct {
	BlockExpiredCertificate types.Bool `tfsdk:"block_expired_certificate"`
	BlockUntrustedIssuer    types.Bool `tfsdk:"block_untrusted_issuer"`
}

type decryptionProfilesDsModelSslProtocolSettingsObject struct {
	AuthAlgoMd5             types.Bool   `tfsdk:"auth_algo_md5"`
	AuthAlgoSha1            types.Bool   `tfsdk:"auth_algo_sha1"`
	AuthAlgoSha256          types.Bool   `tfsdk:"auth_algo_sha256"`
	AuthAlgoSha384          types.Bool   `tfsdk:"auth_algo_sha384"`
	EncAlgo3des             types.Bool   `tfsdk:"enc_algo3des"`
	EncAlgoAes128Cbc        types.Bool   `tfsdk:"enc_algo_aes128_cbc"`
	EncAlgoAes128Gcm        types.Bool   `tfsdk:"enc_algo_aes128_gcm"`
	EncAlgoAes256Cbc        types.Bool   `tfsdk:"enc_algo_aes256_cbc"`
	EncAlgoAes256Gcm        types.Bool   `tfsdk:"enc_algo_aes256_gcm"`
	EncAlgoChacha20Poly1305 types.Bool   `tfsdk:"enc_algo_chacha20_poly1305"`
	EncAlgoRc4              types.Bool   `tfsdk:"enc_algo_rc4"`
	KeyxchgAlgoDhe          types.Bool   `tfsdk:"keyxchg_algo_dhe"`
	KeyxchgAlgoEcdhe        types.Bool   `tfsdk:"keyxchg_algo_ecdhe"`
	KeyxchgAlgoRsa          types.Bool   `tfsdk:"keyxchg_algo_rsa"`
	MaxVersion              types.String `tfsdk:"max_version"`
	MinVersion              types.String `tfsdk:"min_version"`
}

// Metadata returns the data source type name.
func (d *decryptionProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_decryption_profiles"
}

// Schema defines the schema for this listing data source.
func (d *decryptionProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"name": dsschema.StringAttribute{
				Description:         "The `name` parameter.",
				MarkdownDescription: "The `name` parameter.",
				Computed:            true,
			},
			"ssl_forward_proxy": dsschema.SingleNestedAttribute{
				Description:         "The `ssl_forward_proxy` parameter.",
				MarkdownDescription: "The `ssl_forward_proxy` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"auto_include_altname": dsschema.BoolAttribute{
						Description:         "The `auto_include_altname` parameter.",
						MarkdownDescription: "The `auto_include_altname` parameter.",
						Computed:            true,
					},
					"block_client_cert": dsschema.BoolAttribute{
						Description:         "The `block_client_cert` parameter.",
						MarkdownDescription: "The `block_client_cert` parameter.",
						Computed:            true,
					},
					"block_expired_certificate": dsschema.BoolAttribute{
						Description:         "The `block_expired_certificate` parameter.",
						MarkdownDescription: "The `block_expired_certificate` parameter.",
						Computed:            true,
					},
					"block_timeout_cert": dsschema.BoolAttribute{
						Description:         "The `block_timeout_cert` parameter.",
						MarkdownDescription: "The `block_timeout_cert` parameter.",
						Computed:            true,
					},
					"block_tls13_downgrade_no_resource": dsschema.BoolAttribute{
						Description:         "The `block_tls13_downgrade_no_resource` parameter.",
						MarkdownDescription: "The `block_tls13_downgrade_no_resource` parameter.",
						Computed:            true,
					},
					"block_unknown_cert": dsschema.BoolAttribute{
						Description:         "The `block_unknown_cert` parameter.",
						MarkdownDescription: "The `block_unknown_cert` parameter.",
						Computed:            true,
					},
					"block_unsupported_cipher": dsschema.BoolAttribute{
						Description:         "The `block_unsupported_cipher` parameter.",
						MarkdownDescription: "The `block_unsupported_cipher` parameter.",
						Computed:            true,
					},
					"block_unsupported_version": dsschema.BoolAttribute{
						Description:         "The `block_unsupported_version` parameter.",
						MarkdownDescription: "The `block_unsupported_version` parameter.",
						Computed:            true,
					},
					"block_untrusted_issuer": dsschema.BoolAttribute{
						Description:         "The `block_untrusted_issuer` parameter.",
						MarkdownDescription: "The `block_untrusted_issuer` parameter.",
						Computed:            true,
					},
					"restrict_cert_exts": dsschema.BoolAttribute{
						Description:         "The `restrict_cert_exts` parameter.",
						MarkdownDescription: "The `restrict_cert_exts` parameter.",
						Computed:            true,
					},
					"strip_alpn": dsschema.BoolAttribute{
						Description:         "The `strip_alpn` parameter.",
						MarkdownDescription: "The `strip_alpn` parameter.",
						Computed:            true,
					},
				},
			},
			"ssl_inbound_proxy": dsschema.SingleNestedAttribute{
				Description:         "The `ssl_inbound_proxy` parameter.",
				MarkdownDescription: "The `ssl_inbound_proxy` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"block_if_hsm_unavailable": dsschema.BoolAttribute{
						Description:         "The `block_if_hsm_unavailable` parameter.",
						MarkdownDescription: "The `block_if_hsm_unavailable` parameter.",
						Computed:            true,
					},
					"block_if_no_resource": dsschema.BoolAttribute{
						Description:         "The `block_if_no_resource` parameter.",
						MarkdownDescription: "The `block_if_no_resource` parameter.",
						Computed:            true,
					},
					"block_unsupported_cipher": dsschema.BoolAttribute{
						Description:         "The `block_unsupported_cipher` parameter.",
						MarkdownDescription: "The `block_unsupported_cipher` parameter.",
						Computed:            true,
					},
					"block_unsupported_version": dsschema.BoolAttribute{
						Description:         "The `block_unsupported_version` parameter.",
						MarkdownDescription: "The `block_unsupported_version` parameter.",
						Computed:            true,
					},
				},
			},
			"ssl_no_proxy": dsschema.SingleNestedAttribute{
				Description:         "The `ssl_no_proxy` parameter.",
				MarkdownDescription: "The `ssl_no_proxy` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"block_expired_certificate": dsschema.BoolAttribute{
						Description:         "The `block_expired_certificate` parameter.",
						MarkdownDescription: "The `block_expired_certificate` parameter.",
						Computed:            true,
					},
					"block_untrusted_issuer": dsschema.BoolAttribute{
						Description:         "The `block_untrusted_issuer` parameter.",
						MarkdownDescription: "The `block_untrusted_issuer` parameter.",
						Computed:            true,
					},
				},
			},
			"ssl_protocol_settings": dsschema.SingleNestedAttribute{
				Description:         "The `ssl_protocol_settings` parameter.",
				MarkdownDescription: "The `ssl_protocol_settings` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"auth_algo_md5": dsschema.BoolAttribute{
						Description:         "The `auth_algo_md5` parameter.",
						MarkdownDescription: "The `auth_algo_md5` parameter.",
						Computed:            true,
					},
					"auth_algo_sha1": dsschema.BoolAttribute{
						Description:         "The `auth_algo_sha1` parameter.",
						MarkdownDescription: "The `auth_algo_sha1` parameter.",
						Computed:            true,
					},
					"auth_algo_sha256": dsschema.BoolAttribute{
						Description:         "The `auth_algo_sha256` parameter.",
						MarkdownDescription: "The `auth_algo_sha256` parameter.",
						Computed:            true,
					},
					"auth_algo_sha384": dsschema.BoolAttribute{
						Description:         "The `auth_algo_sha384` parameter.",
						MarkdownDescription: "The `auth_algo_sha384` parameter.",
						Computed:            true,
					},
					"enc_algo3des": dsschema.BoolAttribute{
						Description:         "The `enc_algo3des` parameter.",
						MarkdownDescription: "The `enc_algo3des` parameter.",
						Computed:            true,
					},
					"enc_algo_aes128_cbc": dsschema.BoolAttribute{
						Description:         "The `enc_algo_aes128_cbc` parameter.",
						MarkdownDescription: "The `enc_algo_aes128_cbc` parameter.",
						Computed:            true,
					},
					"enc_algo_aes128_gcm": dsschema.BoolAttribute{
						Description:         "The `enc_algo_aes128_gcm` parameter.",
						MarkdownDescription: "The `enc_algo_aes128_gcm` parameter.",
						Computed:            true,
					},
					"enc_algo_aes256_cbc": dsschema.BoolAttribute{
						Description:         "The `enc_algo_aes256_cbc` parameter.",
						MarkdownDescription: "The `enc_algo_aes256_cbc` parameter.",
						Computed:            true,
					},
					"enc_algo_aes256_gcm": dsschema.BoolAttribute{
						Description:         "The `enc_algo_aes256_gcm` parameter.",
						MarkdownDescription: "The `enc_algo_aes256_gcm` parameter.",
						Computed:            true,
					},
					"enc_algo_chacha20_poly1305": dsschema.BoolAttribute{
						Description:         "The `enc_algo_chacha20_poly1305` parameter.",
						MarkdownDescription: "The `enc_algo_chacha20_poly1305` parameter.",
						Computed:            true,
					},
					"enc_algo_rc4": dsschema.BoolAttribute{
						Description:         "The `enc_algo_rc4` parameter.",
						MarkdownDescription: "The `enc_algo_rc4` parameter.",
						Computed:            true,
					},
					"keyxchg_algo_dhe": dsschema.BoolAttribute{
						Description:         "The `keyxchg_algo_dhe` parameter.",
						MarkdownDescription: "The `keyxchg_algo_dhe` parameter.",
						Computed:            true,
					},
					"keyxchg_algo_ecdhe": dsschema.BoolAttribute{
						Description:         "The `keyxchg_algo_ecdhe` parameter.",
						MarkdownDescription: "The `keyxchg_algo_ecdhe` parameter.",
						Computed:            true,
					},
					"keyxchg_algo_rsa": dsschema.BoolAttribute{
						Description:         "The `keyxchg_algo_rsa` parameter.",
						MarkdownDescription: "The `keyxchg_algo_rsa` parameter.",
						Computed:            true,
					},
					"max_version": dsschema.StringAttribute{
						Description:         "The `max_version` parameter.",
						MarkdownDescription: "The `max_version` parameter.",
						Computed:            true,
					},
					"min_version": dsschema.StringAttribute{
						Description:         "The `min_version` parameter.",
						MarkdownDescription: "The `min_version` parameter.",
						Computed:            true,
					},
				},
			},
		},
	}
}

// Configure prepares the struct.
func (d *decryptionProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *decryptionProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state decryptionProfilesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_decryption_profiles",
		"object_id":                   state.ObjectId.ValueString(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := bpgvUeD.NewClient(d.client)
	input := bpgvUeD.ReadInput{
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
	var var0 *decryptionProfilesDsModelSslForwardProxyObject
	if ans.SslForwardProxy != nil {
		var0 = &decryptionProfilesDsModelSslForwardProxyObject{}
		var0.AutoIncludeAltname = types.BoolValue(ans.SslForwardProxy.AutoIncludeAltname)
		var0.BlockClientCert = types.BoolValue(ans.SslForwardProxy.BlockClientCert)
		var0.BlockExpiredCertificate = types.BoolValue(ans.SslForwardProxy.BlockExpiredCertificate)
		var0.BlockTimeoutCert = types.BoolValue(ans.SslForwardProxy.BlockTimeoutCert)
		var0.BlockTls13DowngradeNoResource = types.BoolValue(ans.SslForwardProxy.BlockTls13DowngradeNoResource)
		var0.BlockUnknownCert = types.BoolValue(ans.SslForwardProxy.BlockUnknownCert)
		var0.BlockUnsupportedCipher = types.BoolValue(ans.SslForwardProxy.BlockUnsupportedCipher)
		var0.BlockUnsupportedVersion = types.BoolValue(ans.SslForwardProxy.BlockUnsupportedVersion)
		var0.BlockUntrustedIssuer = types.BoolValue(ans.SslForwardProxy.BlockUntrustedIssuer)
		var0.RestrictCertExts = types.BoolValue(ans.SslForwardProxy.RestrictCertExts)
		var0.StripAlpn = types.BoolValue(ans.SslForwardProxy.StripAlpn)
	}
	var var1 *decryptionProfilesDsModelSslInboundProxyObject
	if ans.SslInboundProxy != nil {
		var1 = &decryptionProfilesDsModelSslInboundProxyObject{}
		var1.BlockIfHsmUnavailable = types.BoolValue(ans.SslInboundProxy.BlockIfHsmUnavailable)
		var1.BlockIfNoResource = types.BoolValue(ans.SslInboundProxy.BlockIfNoResource)
		var1.BlockUnsupportedCipher = types.BoolValue(ans.SslInboundProxy.BlockUnsupportedCipher)
		var1.BlockUnsupportedVersion = types.BoolValue(ans.SslInboundProxy.BlockUnsupportedVersion)
	}
	var var2 *decryptionProfilesDsModelSslNoProxyObject
	if ans.SslNoProxy != nil {
		var2 = &decryptionProfilesDsModelSslNoProxyObject{}
		var2.BlockExpiredCertificate = types.BoolValue(ans.SslNoProxy.BlockExpiredCertificate)
		var2.BlockUntrustedIssuer = types.BoolValue(ans.SslNoProxy.BlockUntrustedIssuer)
	}
	var var3 *decryptionProfilesDsModelSslProtocolSettingsObject
	if ans.SslProtocolSettings != nil {
		var3 = &decryptionProfilesDsModelSslProtocolSettingsObject{}
		var3.AuthAlgoMd5 = types.BoolValue(ans.SslProtocolSettings.AuthAlgoMd5)
		var3.AuthAlgoSha1 = types.BoolValue(ans.SslProtocolSettings.AuthAlgoSha1)
		var3.AuthAlgoSha256 = types.BoolValue(ans.SslProtocolSettings.AuthAlgoSha256)
		var3.AuthAlgoSha384 = types.BoolValue(ans.SslProtocolSettings.AuthAlgoSha384)
		var3.EncAlgo3des = types.BoolValue(ans.SslProtocolSettings.EncAlgo3des)
		var3.EncAlgoAes128Cbc = types.BoolValue(ans.SslProtocolSettings.EncAlgoAes128Cbc)
		var3.EncAlgoAes128Gcm = types.BoolValue(ans.SslProtocolSettings.EncAlgoAes128Gcm)
		var3.EncAlgoAes256Cbc = types.BoolValue(ans.SslProtocolSettings.EncAlgoAes256Cbc)
		var3.EncAlgoAes256Gcm = types.BoolValue(ans.SslProtocolSettings.EncAlgoAes256Gcm)
		var3.EncAlgoChacha20Poly1305 = types.BoolValue(ans.SslProtocolSettings.EncAlgoChacha20Poly1305)
		var3.EncAlgoRc4 = types.BoolValue(ans.SslProtocolSettings.EncAlgoRc4)
		var3.KeyxchgAlgoDhe = types.BoolValue(ans.SslProtocolSettings.KeyxchgAlgoDhe)
		var3.KeyxchgAlgoEcdhe = types.BoolValue(ans.SslProtocolSettings.KeyxchgAlgoEcdhe)
		var3.KeyxchgAlgoRsa = types.BoolValue(ans.SslProtocolSettings.KeyxchgAlgoRsa)
		var3.MaxVersion = types.StringValue(ans.SslProtocolSettings.MaxVersion)
		var3.MinVersion = types.StringValue(ans.SslProtocolSettings.MinVersion)
	}
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.SslForwardProxy = var0
	state.SslInboundProxy = var1
	state.SslNoProxy = var2
	state.SslProtocolSettings = var3

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &decryptionProfilesResource{}
	_ resource.ResourceWithConfigure   = &decryptionProfilesResource{}
	_ resource.ResourceWithImportState = &decryptionProfilesResource{}
)

func NewDecryptionProfilesResource() resource.Resource {
	return &decryptionProfilesResource{}
}

type decryptionProfilesResource struct {
	client *sase.Client
}

type decryptionProfilesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/decryption-profiles
	ObjectId            types.String                                        `tfsdk:"object_id"`
	Name                types.String                                        `tfsdk:"name"`
	SslForwardProxy     *decryptionProfilesRsModelSslForwardProxyObject     `tfsdk:"ssl_forward_proxy"`
	SslInboundProxy     *decryptionProfilesRsModelSslInboundProxyObject     `tfsdk:"ssl_inbound_proxy"`
	SslNoProxy          *decryptionProfilesRsModelSslNoProxyObject          `tfsdk:"ssl_no_proxy"`
	SslProtocolSettings *decryptionProfilesRsModelSslProtocolSettingsObject `tfsdk:"ssl_protocol_settings"`
}

type decryptionProfilesRsModelSslForwardProxyObject struct {
	AutoIncludeAltname            types.Bool `tfsdk:"auto_include_altname"`
	BlockClientCert               types.Bool `tfsdk:"block_client_cert"`
	BlockExpiredCertificate       types.Bool `tfsdk:"block_expired_certificate"`
	BlockTimeoutCert              types.Bool `tfsdk:"block_timeout_cert"`
	BlockTls13DowngradeNoResource types.Bool `tfsdk:"block_tls13_downgrade_no_resource"`
	BlockUnknownCert              types.Bool `tfsdk:"block_unknown_cert"`
	BlockUnsupportedCipher        types.Bool `tfsdk:"block_unsupported_cipher"`
	BlockUnsupportedVersion       types.Bool `tfsdk:"block_unsupported_version"`
	BlockUntrustedIssuer          types.Bool `tfsdk:"block_untrusted_issuer"`
	RestrictCertExts              types.Bool `tfsdk:"restrict_cert_exts"`
	StripAlpn                     types.Bool `tfsdk:"strip_alpn"`
}

type decryptionProfilesRsModelSslInboundProxyObject struct {
	BlockIfHsmUnavailable   types.Bool `tfsdk:"block_if_hsm_unavailable"`
	BlockIfNoResource       types.Bool `tfsdk:"block_if_no_resource"`
	BlockUnsupportedCipher  types.Bool `tfsdk:"block_unsupported_cipher"`
	BlockUnsupportedVersion types.Bool `tfsdk:"block_unsupported_version"`
}

type decryptionProfilesRsModelSslNoProxyObject struct {
	BlockExpiredCertificate types.Bool `tfsdk:"block_expired_certificate"`
	BlockUntrustedIssuer    types.Bool `tfsdk:"block_untrusted_issuer"`
}

type decryptionProfilesRsModelSslProtocolSettingsObject struct {
	AuthAlgoMd5             types.Bool   `tfsdk:"auth_algo_md5"`
	AuthAlgoSha1            types.Bool   `tfsdk:"auth_algo_sha1"`
	AuthAlgoSha256          types.Bool   `tfsdk:"auth_algo_sha256"`
	AuthAlgoSha384          types.Bool   `tfsdk:"auth_algo_sha384"`
	EncAlgo3des             types.Bool   `tfsdk:"enc_algo3des"`
	EncAlgoAes128Cbc        types.Bool   `tfsdk:"enc_algo_aes128_cbc"`
	EncAlgoAes128Gcm        types.Bool   `tfsdk:"enc_algo_aes128_gcm"`
	EncAlgoAes256Cbc        types.Bool   `tfsdk:"enc_algo_aes256_cbc"`
	EncAlgoAes256Gcm        types.Bool   `tfsdk:"enc_algo_aes256_gcm"`
	EncAlgoChacha20Poly1305 types.Bool   `tfsdk:"enc_algo_chacha20_poly1305"`
	EncAlgoRc4              types.Bool   `tfsdk:"enc_algo_rc4"`
	KeyxchgAlgoDhe          types.Bool   `tfsdk:"keyxchg_algo_dhe"`
	KeyxchgAlgoEcdhe        types.Bool   `tfsdk:"keyxchg_algo_ecdhe"`
	KeyxchgAlgoRsa          types.Bool   `tfsdk:"keyxchg_algo_rsa"`
	MaxVersion              types.String `tfsdk:"max_version"`
	MinVersion              types.String `tfsdk:"min_version"`
}

// Metadata returns the data source type name.
func (r *decryptionProfilesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_decryption_profiles"
}

// Schema defines the schema for this listing data source.
func (r *decryptionProfilesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

			"object_id": rsschema.StringAttribute{
				Description:         "The `object_id` parameter.",
				MarkdownDescription: "The `object_id` parameter.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": rsschema.StringAttribute{
				Description:         "The `name` parameter.",
				MarkdownDescription: "The `name` parameter.",
				Required:            true,
			},
			"ssl_forward_proxy": rsschema.SingleNestedAttribute{
				Description:         "The `ssl_forward_proxy` parameter.",
				MarkdownDescription: "The `ssl_forward_proxy` parameter.",
				Optional:            true,
				Attributes: map[string]rsschema.Attribute{
					"auto_include_altname": rsschema.BoolAttribute{
						Description:         "The `auto_include_altname` parameter. Default: `false`.",
						MarkdownDescription: "The `auto_include_altname` parameter. Default: `false`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"block_client_cert": rsschema.BoolAttribute{
						Description:         "The `block_client_cert` parameter. Default: `false`.",
						MarkdownDescription: "The `block_client_cert` parameter. Default: `false`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"block_expired_certificate": rsschema.BoolAttribute{
						Description:         "The `block_expired_certificate` parameter. Default: `false`.",
						MarkdownDescription: "The `block_expired_certificate` parameter. Default: `false`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"block_timeout_cert": rsschema.BoolAttribute{
						Description:         "The `block_timeout_cert` parameter. Default: `false`.",
						MarkdownDescription: "The `block_timeout_cert` parameter. Default: `false`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"block_tls13_downgrade_no_resource": rsschema.BoolAttribute{
						Description:         "The `block_tls13_downgrade_no_resource` parameter. Default: `false`.",
						MarkdownDescription: "The `block_tls13_downgrade_no_resource` parameter. Default: `false`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"block_unknown_cert": rsschema.BoolAttribute{
						Description:         "The `block_unknown_cert` parameter. Default: `false`.",
						MarkdownDescription: "The `block_unknown_cert` parameter. Default: `false`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"block_unsupported_cipher": rsschema.BoolAttribute{
						Description:         "The `block_unsupported_cipher` parameter. Default: `false`.",
						MarkdownDescription: "The `block_unsupported_cipher` parameter. Default: `false`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"block_unsupported_version": rsschema.BoolAttribute{
						Description:         "The `block_unsupported_version` parameter. Default: `false`.",
						MarkdownDescription: "The `block_unsupported_version` parameter. Default: `false`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"block_untrusted_issuer": rsschema.BoolAttribute{
						Description:         "The `block_untrusted_issuer` parameter. Default: `false`.",
						MarkdownDescription: "The `block_untrusted_issuer` parameter. Default: `false`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"restrict_cert_exts": rsschema.BoolAttribute{
						Description:         "The `restrict_cert_exts` parameter. Default: `false`.",
						MarkdownDescription: "The `restrict_cert_exts` parameter. Default: `false`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"strip_alpn": rsschema.BoolAttribute{
						Description:         "The `strip_alpn` parameter. Default: `false`.",
						MarkdownDescription: "The `strip_alpn` parameter. Default: `false`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
				},
			},
			"ssl_inbound_proxy": rsschema.SingleNestedAttribute{
				Description:         "The `ssl_inbound_proxy` parameter.",
				MarkdownDescription: "The `ssl_inbound_proxy` parameter.",
				Optional:            true,
				Attributes: map[string]rsschema.Attribute{
					"block_if_hsm_unavailable": rsschema.BoolAttribute{
						Description:         "The `block_if_hsm_unavailable` parameter. Default: `false`.",
						MarkdownDescription: "The `block_if_hsm_unavailable` parameter. Default: `false`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"block_if_no_resource": rsschema.BoolAttribute{
						Description:         "The `block_if_no_resource` parameter. Default: `false`.",
						MarkdownDescription: "The `block_if_no_resource` parameter. Default: `false`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"block_unsupported_cipher": rsschema.BoolAttribute{
						Description:         "The `block_unsupported_cipher` parameter. Default: `false`.",
						MarkdownDescription: "The `block_unsupported_cipher` parameter. Default: `false`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"block_unsupported_version": rsschema.BoolAttribute{
						Description:         "The `block_unsupported_version` parameter. Default: `false`.",
						MarkdownDescription: "The `block_unsupported_version` parameter. Default: `false`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
				},
			},
			"ssl_no_proxy": rsschema.SingleNestedAttribute{
				Description:         "The `ssl_no_proxy` parameter.",
				MarkdownDescription: "The `ssl_no_proxy` parameter.",
				Optional:            true,
				Attributes: map[string]rsschema.Attribute{
					"block_expired_certificate": rsschema.BoolAttribute{
						Description:         "The `block_expired_certificate` parameter. Default: `false`.",
						MarkdownDescription: "The `block_expired_certificate` parameter. Default: `false`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"block_untrusted_issuer": rsschema.BoolAttribute{
						Description:         "The `block_untrusted_issuer` parameter. Default: `false`.",
						MarkdownDescription: "The `block_untrusted_issuer` parameter. Default: `false`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
				},
			},
			"ssl_protocol_settings": rsschema.SingleNestedAttribute{
				Description:         "The `ssl_protocol_settings` parameter.",
				MarkdownDescription: "The `ssl_protocol_settings` parameter.",
				Optional:            true,
				Attributes: map[string]rsschema.Attribute{
					"auth_algo_md5": rsschema.BoolAttribute{
						Description:         "The `auth_algo_md5` parameter. Default: `true`.",
						MarkdownDescription: "The `auth_algo_md5` parameter. Default: `true`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(true),
						},
					},
					"auth_algo_sha1": rsschema.BoolAttribute{
						Description:         "The `auth_algo_sha1` parameter. Default: `true`.",
						MarkdownDescription: "The `auth_algo_sha1` parameter. Default: `true`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(true),
						},
					},
					"auth_algo_sha256": rsschema.BoolAttribute{
						Description:         "The `auth_algo_sha256` parameter. Default: `true`.",
						MarkdownDescription: "The `auth_algo_sha256` parameter. Default: `true`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(true),
						},
					},
					"auth_algo_sha384": rsschema.BoolAttribute{
						Description:         "The `auth_algo_sha384` parameter. Default: `true`.",
						MarkdownDescription: "The `auth_algo_sha384` parameter. Default: `true`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(true),
						},
					},
					"enc_algo3des": rsschema.BoolAttribute{
						Description:         "The `enc_algo3des` parameter. Default: `true`.",
						MarkdownDescription: "The `enc_algo3des` parameter. Default: `true`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(true),
						},
					},
					"enc_algo_aes128_cbc": rsschema.BoolAttribute{
						Description:         "The `enc_algo_aes128_cbc` parameter. Default: `true`.",
						MarkdownDescription: "The `enc_algo_aes128_cbc` parameter. Default: `true`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(true),
						},
					},
					"enc_algo_aes128_gcm": rsschema.BoolAttribute{
						Description:         "The `enc_algo_aes128_gcm` parameter. Default: `true`.",
						MarkdownDescription: "The `enc_algo_aes128_gcm` parameter. Default: `true`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(true),
						},
					},
					"enc_algo_aes256_cbc": rsschema.BoolAttribute{
						Description:         "The `enc_algo_aes256_cbc` parameter. Default: `true`.",
						MarkdownDescription: "The `enc_algo_aes256_cbc` parameter. Default: `true`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(true),
						},
					},
					"enc_algo_aes256_gcm": rsschema.BoolAttribute{
						Description:         "The `enc_algo_aes256_gcm` parameter. Default: `true`.",
						MarkdownDescription: "The `enc_algo_aes256_gcm` parameter. Default: `true`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(true),
						},
					},
					"enc_algo_chacha20_poly1305": rsschema.BoolAttribute{
						Description:         "The `enc_algo_chacha20_poly1305` parameter. Default: `true`.",
						MarkdownDescription: "The `enc_algo_chacha20_poly1305` parameter. Default: `true`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(true),
						},
					},
					"enc_algo_rc4": rsschema.BoolAttribute{
						Description:         "The `enc_algo_rc4` parameter. Default: `true`.",
						MarkdownDescription: "The `enc_algo_rc4` parameter. Default: `true`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(true),
						},
					},
					"keyxchg_algo_dhe": rsschema.BoolAttribute{
						Description:         "The `keyxchg_algo_dhe` parameter. Default: `true`.",
						MarkdownDescription: "The `keyxchg_algo_dhe` parameter. Default: `true`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(true),
						},
					},
					"keyxchg_algo_ecdhe": rsschema.BoolAttribute{
						Description:         "The `keyxchg_algo_ecdhe` parameter. Default: `true`.",
						MarkdownDescription: "The `keyxchg_algo_ecdhe` parameter. Default: `true`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(true),
						},
					},
					"keyxchg_algo_rsa": rsschema.BoolAttribute{
						Description:         "The `keyxchg_algo_rsa` parameter. Default: `true`.",
						MarkdownDescription: "The `keyxchg_algo_rsa` parameter. Default: `true`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(true),
						},
					},
					"max_version": rsschema.StringAttribute{
						Description:         "The `max_version` parameter. Default: `\"tls1-2\"`. Value must be one of: `\"sslv3\"`, `\"tls1-0\"`, `\"tls1-1\"`, `\"tls1-2\"`, `\"tls1-3\"`, `\"max\"`.",
						MarkdownDescription: "The `max_version` parameter. Default: `\"tls1-2\"`. Value must be one of: `\"sslv3\"`, `\"tls1-0\"`, `\"tls1-1\"`, `\"tls1-2\"`, `\"tls1-3\"`, `\"max\"`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							DefaultString("tls1-2"),
						},
						Validators: []validator.String{
							stringvalidator.OneOf("sslv3", "tls1-0", "tls1-1", "tls1-2", "tls1-3", "max"),
						},
					},
					"min_version": rsschema.StringAttribute{
						Description:         "The `min_version` parameter. Default: `\"tls1-0\"`. Value must be one of: `\"sslv3\"`, `\"tls1-0\"`, `\"tls1-1\"`, `\"tls1-2\"`, `\"tls1-3\"`.",
						MarkdownDescription: "The `min_version` parameter. Default: `\"tls1-0\"`. Value must be one of: `\"sslv3\"`, `\"tls1-0\"`, `\"tls1-1\"`, `\"tls1-2\"`, `\"tls1-3\"`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							DefaultString("tls1-0"),
						},
						Validators: []validator.String{
							stringvalidator.OneOf("sslv3", "tls1-0", "tls1-1", "tls1-2", "tls1-3"),
						},
					},
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *decryptionProfilesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *decryptionProfilesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state decryptionProfilesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_decryption_profiles",
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := bpgvUeD.NewClient(r.client)
	input := bpgvUeD.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 vMYBRZK.Config
	var0.Name = state.Name.ValueString()
	var var1 *vMYBRZK.SslForwardProxyObject
	if state.SslForwardProxy != nil {
		var1 = &vMYBRZK.SslForwardProxyObject{}
		var1.AutoIncludeAltname = state.SslForwardProxy.AutoIncludeAltname.ValueBool()
		var1.BlockClientCert = state.SslForwardProxy.BlockClientCert.ValueBool()
		var1.BlockExpiredCertificate = state.SslForwardProxy.BlockExpiredCertificate.ValueBool()
		var1.BlockTimeoutCert = state.SslForwardProxy.BlockTimeoutCert.ValueBool()
		var1.BlockTls13DowngradeNoResource = state.SslForwardProxy.BlockTls13DowngradeNoResource.ValueBool()
		var1.BlockUnknownCert = state.SslForwardProxy.BlockUnknownCert.ValueBool()
		var1.BlockUnsupportedCipher = state.SslForwardProxy.BlockUnsupportedCipher.ValueBool()
		var1.BlockUnsupportedVersion = state.SslForwardProxy.BlockUnsupportedVersion.ValueBool()
		var1.BlockUntrustedIssuer = state.SslForwardProxy.BlockUntrustedIssuer.ValueBool()
		var1.RestrictCertExts = state.SslForwardProxy.RestrictCertExts.ValueBool()
		var1.StripAlpn = state.SslForwardProxy.StripAlpn.ValueBool()
	}
	var0.SslForwardProxy = var1
	var var2 *vMYBRZK.SslInboundProxyObject
	if state.SslInboundProxy != nil {
		var2 = &vMYBRZK.SslInboundProxyObject{}
		var2.BlockIfHsmUnavailable = state.SslInboundProxy.BlockIfHsmUnavailable.ValueBool()
		var2.BlockIfNoResource = state.SslInboundProxy.BlockIfNoResource.ValueBool()
		var2.BlockUnsupportedCipher = state.SslInboundProxy.BlockUnsupportedCipher.ValueBool()
		var2.BlockUnsupportedVersion = state.SslInboundProxy.BlockUnsupportedVersion.ValueBool()
	}
	var0.SslInboundProxy = var2
	var var3 *vMYBRZK.SslNoProxyObject
	if state.SslNoProxy != nil {
		var3 = &vMYBRZK.SslNoProxyObject{}
		var3.BlockExpiredCertificate = state.SslNoProxy.BlockExpiredCertificate.ValueBool()
		var3.BlockUntrustedIssuer = state.SslNoProxy.BlockUntrustedIssuer.ValueBool()
	}
	var0.SslNoProxy = var3
	var var4 *vMYBRZK.SslProtocolSettingsObject
	if state.SslProtocolSettings != nil {
		var4 = &vMYBRZK.SslProtocolSettingsObject{}
		var4.AuthAlgoMd5 = state.SslProtocolSettings.AuthAlgoMd5.ValueBool()
		var4.AuthAlgoSha1 = state.SslProtocolSettings.AuthAlgoSha1.ValueBool()
		var4.AuthAlgoSha256 = state.SslProtocolSettings.AuthAlgoSha256.ValueBool()
		var4.AuthAlgoSha384 = state.SslProtocolSettings.AuthAlgoSha384.ValueBool()
		var4.EncAlgo3des = state.SslProtocolSettings.EncAlgo3des.ValueBool()
		var4.EncAlgoAes128Cbc = state.SslProtocolSettings.EncAlgoAes128Cbc.ValueBool()
		var4.EncAlgoAes128Gcm = state.SslProtocolSettings.EncAlgoAes128Gcm.ValueBool()
		var4.EncAlgoAes256Cbc = state.SslProtocolSettings.EncAlgoAes256Cbc.ValueBool()
		var4.EncAlgoAes256Gcm = state.SslProtocolSettings.EncAlgoAes256Gcm.ValueBool()
		var4.EncAlgoChacha20Poly1305 = state.SslProtocolSettings.EncAlgoChacha20Poly1305.ValueBool()
		var4.EncAlgoRc4 = state.SslProtocolSettings.EncAlgoRc4.ValueBool()
		var4.KeyxchgAlgoDhe = state.SslProtocolSettings.KeyxchgAlgoDhe.ValueBool()
		var4.KeyxchgAlgoEcdhe = state.SslProtocolSettings.KeyxchgAlgoEcdhe.ValueBool()
		var4.KeyxchgAlgoRsa = state.SslProtocolSettings.KeyxchgAlgoRsa.ValueBool()
		var4.MaxVersion = state.SslProtocolSettings.MaxVersion.ValueString()
		var4.MinVersion = state.SslProtocolSettings.MinVersion.ValueString()
	}
	var0.SslProtocolSettings = var4
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
	var var5 *decryptionProfilesRsModelSslForwardProxyObject
	if ans.SslForwardProxy != nil {
		var5 = &decryptionProfilesRsModelSslForwardProxyObject{}
		var5.AutoIncludeAltname = types.BoolValue(ans.SslForwardProxy.AutoIncludeAltname)
		var5.BlockClientCert = types.BoolValue(ans.SslForwardProxy.BlockClientCert)
		var5.BlockExpiredCertificate = types.BoolValue(ans.SslForwardProxy.BlockExpiredCertificate)
		var5.BlockTimeoutCert = types.BoolValue(ans.SslForwardProxy.BlockTimeoutCert)
		var5.BlockTls13DowngradeNoResource = types.BoolValue(ans.SslForwardProxy.BlockTls13DowngradeNoResource)
		var5.BlockUnknownCert = types.BoolValue(ans.SslForwardProxy.BlockUnknownCert)
		var5.BlockUnsupportedCipher = types.BoolValue(ans.SslForwardProxy.BlockUnsupportedCipher)
		var5.BlockUnsupportedVersion = types.BoolValue(ans.SslForwardProxy.BlockUnsupportedVersion)
		var5.BlockUntrustedIssuer = types.BoolValue(ans.SslForwardProxy.BlockUntrustedIssuer)
		var5.RestrictCertExts = types.BoolValue(ans.SslForwardProxy.RestrictCertExts)
		var5.StripAlpn = types.BoolValue(ans.SslForwardProxy.StripAlpn)
	}
	var var6 *decryptionProfilesRsModelSslInboundProxyObject
	if ans.SslInboundProxy != nil {
		var6 = &decryptionProfilesRsModelSslInboundProxyObject{}
		var6.BlockIfHsmUnavailable = types.BoolValue(ans.SslInboundProxy.BlockIfHsmUnavailable)
		var6.BlockIfNoResource = types.BoolValue(ans.SslInboundProxy.BlockIfNoResource)
		var6.BlockUnsupportedCipher = types.BoolValue(ans.SslInboundProxy.BlockUnsupportedCipher)
		var6.BlockUnsupportedVersion = types.BoolValue(ans.SslInboundProxy.BlockUnsupportedVersion)
	}
	var var7 *decryptionProfilesRsModelSslNoProxyObject
	if ans.SslNoProxy != nil {
		var7 = &decryptionProfilesRsModelSslNoProxyObject{}
		var7.BlockExpiredCertificate = types.BoolValue(ans.SslNoProxy.BlockExpiredCertificate)
		var7.BlockUntrustedIssuer = types.BoolValue(ans.SslNoProxy.BlockUntrustedIssuer)
	}
	var var8 *decryptionProfilesRsModelSslProtocolSettingsObject
	if ans.SslProtocolSettings != nil {
		var8 = &decryptionProfilesRsModelSslProtocolSettingsObject{}
		var8.AuthAlgoMd5 = types.BoolValue(ans.SslProtocolSettings.AuthAlgoMd5)
		var8.AuthAlgoSha1 = types.BoolValue(ans.SslProtocolSettings.AuthAlgoSha1)
		var8.AuthAlgoSha256 = types.BoolValue(ans.SslProtocolSettings.AuthAlgoSha256)
		var8.AuthAlgoSha384 = types.BoolValue(ans.SslProtocolSettings.AuthAlgoSha384)
		var8.EncAlgo3des = types.BoolValue(ans.SslProtocolSettings.EncAlgo3des)
		var8.EncAlgoAes128Cbc = types.BoolValue(ans.SslProtocolSettings.EncAlgoAes128Cbc)
		var8.EncAlgoAes128Gcm = types.BoolValue(ans.SslProtocolSettings.EncAlgoAes128Gcm)
		var8.EncAlgoAes256Cbc = types.BoolValue(ans.SslProtocolSettings.EncAlgoAes256Cbc)
		var8.EncAlgoAes256Gcm = types.BoolValue(ans.SslProtocolSettings.EncAlgoAes256Gcm)
		var8.EncAlgoChacha20Poly1305 = types.BoolValue(ans.SslProtocolSettings.EncAlgoChacha20Poly1305)
		var8.EncAlgoRc4 = types.BoolValue(ans.SslProtocolSettings.EncAlgoRc4)
		var8.KeyxchgAlgoDhe = types.BoolValue(ans.SslProtocolSettings.KeyxchgAlgoDhe)
		var8.KeyxchgAlgoEcdhe = types.BoolValue(ans.SslProtocolSettings.KeyxchgAlgoEcdhe)
		var8.KeyxchgAlgoRsa = types.BoolValue(ans.SslProtocolSettings.KeyxchgAlgoRsa)
		var8.MaxVersion = types.StringValue(ans.SslProtocolSettings.MaxVersion)
		var8.MinVersion = types.StringValue(ans.SslProtocolSettings.MinVersion)
	}
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.SslForwardProxy = var5
	state.SslInboundProxy = var6
	state.SslNoProxy = var7
	state.SslProtocolSettings = var8

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *decryptionProfilesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state decryptionProfilesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_decryption_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := bpgvUeD.NewClient(r.client)
	input := bpgvUeD.ReadInput{
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
	var var0 *decryptionProfilesRsModelSslForwardProxyObject
	if ans.SslForwardProxy != nil {
		var0 = &decryptionProfilesRsModelSslForwardProxyObject{}
		var0.AutoIncludeAltname = types.BoolValue(ans.SslForwardProxy.AutoIncludeAltname)
		var0.BlockClientCert = types.BoolValue(ans.SslForwardProxy.BlockClientCert)
		var0.BlockExpiredCertificate = types.BoolValue(ans.SslForwardProxy.BlockExpiredCertificate)
		var0.BlockTimeoutCert = types.BoolValue(ans.SslForwardProxy.BlockTimeoutCert)
		var0.BlockTls13DowngradeNoResource = types.BoolValue(ans.SslForwardProxy.BlockTls13DowngradeNoResource)
		var0.BlockUnknownCert = types.BoolValue(ans.SslForwardProxy.BlockUnknownCert)
		var0.BlockUnsupportedCipher = types.BoolValue(ans.SslForwardProxy.BlockUnsupportedCipher)
		var0.BlockUnsupportedVersion = types.BoolValue(ans.SslForwardProxy.BlockUnsupportedVersion)
		var0.BlockUntrustedIssuer = types.BoolValue(ans.SslForwardProxy.BlockUntrustedIssuer)
		var0.RestrictCertExts = types.BoolValue(ans.SslForwardProxy.RestrictCertExts)
		var0.StripAlpn = types.BoolValue(ans.SslForwardProxy.StripAlpn)
	}
	var var1 *decryptionProfilesRsModelSslInboundProxyObject
	if ans.SslInboundProxy != nil {
		var1 = &decryptionProfilesRsModelSslInboundProxyObject{}
		var1.BlockIfHsmUnavailable = types.BoolValue(ans.SslInboundProxy.BlockIfHsmUnavailable)
		var1.BlockIfNoResource = types.BoolValue(ans.SslInboundProxy.BlockIfNoResource)
		var1.BlockUnsupportedCipher = types.BoolValue(ans.SslInboundProxy.BlockUnsupportedCipher)
		var1.BlockUnsupportedVersion = types.BoolValue(ans.SslInboundProxy.BlockUnsupportedVersion)
	}
	var var2 *decryptionProfilesRsModelSslNoProxyObject
	if ans.SslNoProxy != nil {
		var2 = &decryptionProfilesRsModelSslNoProxyObject{}
		var2.BlockExpiredCertificate = types.BoolValue(ans.SslNoProxy.BlockExpiredCertificate)
		var2.BlockUntrustedIssuer = types.BoolValue(ans.SslNoProxy.BlockUntrustedIssuer)
	}
	var var3 *decryptionProfilesRsModelSslProtocolSettingsObject
	if ans.SslProtocolSettings != nil {
		var3 = &decryptionProfilesRsModelSslProtocolSettingsObject{}
		var3.AuthAlgoMd5 = types.BoolValue(ans.SslProtocolSettings.AuthAlgoMd5)
		var3.AuthAlgoSha1 = types.BoolValue(ans.SslProtocolSettings.AuthAlgoSha1)
		var3.AuthAlgoSha256 = types.BoolValue(ans.SslProtocolSettings.AuthAlgoSha256)
		var3.AuthAlgoSha384 = types.BoolValue(ans.SslProtocolSettings.AuthAlgoSha384)
		var3.EncAlgo3des = types.BoolValue(ans.SslProtocolSettings.EncAlgo3des)
		var3.EncAlgoAes128Cbc = types.BoolValue(ans.SslProtocolSettings.EncAlgoAes128Cbc)
		var3.EncAlgoAes128Gcm = types.BoolValue(ans.SslProtocolSettings.EncAlgoAes128Gcm)
		var3.EncAlgoAes256Cbc = types.BoolValue(ans.SslProtocolSettings.EncAlgoAes256Cbc)
		var3.EncAlgoAes256Gcm = types.BoolValue(ans.SslProtocolSettings.EncAlgoAes256Gcm)
		var3.EncAlgoChacha20Poly1305 = types.BoolValue(ans.SslProtocolSettings.EncAlgoChacha20Poly1305)
		var3.EncAlgoRc4 = types.BoolValue(ans.SslProtocolSettings.EncAlgoRc4)
		var3.KeyxchgAlgoDhe = types.BoolValue(ans.SslProtocolSettings.KeyxchgAlgoDhe)
		var3.KeyxchgAlgoEcdhe = types.BoolValue(ans.SslProtocolSettings.KeyxchgAlgoEcdhe)
		var3.KeyxchgAlgoRsa = types.BoolValue(ans.SslProtocolSettings.KeyxchgAlgoRsa)
		var3.MaxVersion = types.StringValue(ans.SslProtocolSettings.MaxVersion)
		var3.MinVersion = types.StringValue(ans.SslProtocolSettings.MinVersion)
	}
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.SslForwardProxy = var0
	state.SslInboundProxy = var1
	state.SslNoProxy = var2
	state.SslProtocolSettings = var3

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *decryptionProfilesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state decryptionProfilesRsModel
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
		"resource_name":               "sase_decryption_profiles",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := bpgvUeD.NewClient(r.client)
	input := bpgvUeD.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 vMYBRZK.Config
	var0.Name = plan.Name.ValueString()
	var var1 *vMYBRZK.SslForwardProxyObject
	if plan.SslForwardProxy != nil {
		var1 = &vMYBRZK.SslForwardProxyObject{}
		var1.AutoIncludeAltname = plan.SslForwardProxy.AutoIncludeAltname.ValueBool()
		var1.BlockClientCert = plan.SslForwardProxy.BlockClientCert.ValueBool()
		var1.BlockExpiredCertificate = plan.SslForwardProxy.BlockExpiredCertificate.ValueBool()
		var1.BlockTimeoutCert = plan.SslForwardProxy.BlockTimeoutCert.ValueBool()
		var1.BlockTls13DowngradeNoResource = plan.SslForwardProxy.BlockTls13DowngradeNoResource.ValueBool()
		var1.BlockUnknownCert = plan.SslForwardProxy.BlockUnknownCert.ValueBool()
		var1.BlockUnsupportedCipher = plan.SslForwardProxy.BlockUnsupportedCipher.ValueBool()
		var1.BlockUnsupportedVersion = plan.SslForwardProxy.BlockUnsupportedVersion.ValueBool()
		var1.BlockUntrustedIssuer = plan.SslForwardProxy.BlockUntrustedIssuer.ValueBool()
		var1.RestrictCertExts = plan.SslForwardProxy.RestrictCertExts.ValueBool()
		var1.StripAlpn = plan.SslForwardProxy.StripAlpn.ValueBool()
	}
	var0.SslForwardProxy = var1
	var var2 *vMYBRZK.SslInboundProxyObject
	if plan.SslInboundProxy != nil {
		var2 = &vMYBRZK.SslInboundProxyObject{}
		var2.BlockIfHsmUnavailable = plan.SslInboundProxy.BlockIfHsmUnavailable.ValueBool()
		var2.BlockIfNoResource = plan.SslInboundProxy.BlockIfNoResource.ValueBool()
		var2.BlockUnsupportedCipher = plan.SslInboundProxy.BlockUnsupportedCipher.ValueBool()
		var2.BlockUnsupportedVersion = plan.SslInboundProxy.BlockUnsupportedVersion.ValueBool()
	}
	var0.SslInboundProxy = var2
	var var3 *vMYBRZK.SslNoProxyObject
	if plan.SslNoProxy != nil {
		var3 = &vMYBRZK.SslNoProxyObject{}
		var3.BlockExpiredCertificate = plan.SslNoProxy.BlockExpiredCertificate.ValueBool()
		var3.BlockUntrustedIssuer = plan.SslNoProxy.BlockUntrustedIssuer.ValueBool()
	}
	var0.SslNoProxy = var3
	var var4 *vMYBRZK.SslProtocolSettingsObject
	if plan.SslProtocolSettings != nil {
		var4 = &vMYBRZK.SslProtocolSettingsObject{}
		var4.AuthAlgoMd5 = plan.SslProtocolSettings.AuthAlgoMd5.ValueBool()
		var4.AuthAlgoSha1 = plan.SslProtocolSettings.AuthAlgoSha1.ValueBool()
		var4.AuthAlgoSha256 = plan.SslProtocolSettings.AuthAlgoSha256.ValueBool()
		var4.AuthAlgoSha384 = plan.SslProtocolSettings.AuthAlgoSha384.ValueBool()
		var4.EncAlgo3des = plan.SslProtocolSettings.EncAlgo3des.ValueBool()
		var4.EncAlgoAes128Cbc = plan.SslProtocolSettings.EncAlgoAes128Cbc.ValueBool()
		var4.EncAlgoAes128Gcm = plan.SslProtocolSettings.EncAlgoAes128Gcm.ValueBool()
		var4.EncAlgoAes256Cbc = plan.SslProtocolSettings.EncAlgoAes256Cbc.ValueBool()
		var4.EncAlgoAes256Gcm = plan.SslProtocolSettings.EncAlgoAes256Gcm.ValueBool()
		var4.EncAlgoChacha20Poly1305 = plan.SslProtocolSettings.EncAlgoChacha20Poly1305.ValueBool()
		var4.EncAlgoRc4 = plan.SslProtocolSettings.EncAlgoRc4.ValueBool()
		var4.KeyxchgAlgoDhe = plan.SslProtocolSettings.KeyxchgAlgoDhe.ValueBool()
		var4.KeyxchgAlgoEcdhe = plan.SslProtocolSettings.KeyxchgAlgoEcdhe.ValueBool()
		var4.KeyxchgAlgoRsa = plan.SslProtocolSettings.KeyxchgAlgoRsa.ValueBool()
		var4.MaxVersion = plan.SslProtocolSettings.MaxVersion.ValueString()
		var4.MinVersion = plan.SslProtocolSettings.MinVersion.ValueString()
	}
	var0.SslProtocolSettings = var4
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var5 *decryptionProfilesRsModelSslForwardProxyObject
	if ans.SslForwardProxy != nil {
		var5 = &decryptionProfilesRsModelSslForwardProxyObject{}
		var5.AutoIncludeAltname = types.BoolValue(ans.SslForwardProxy.AutoIncludeAltname)
		var5.BlockClientCert = types.BoolValue(ans.SslForwardProxy.BlockClientCert)
		var5.BlockExpiredCertificate = types.BoolValue(ans.SslForwardProxy.BlockExpiredCertificate)
		var5.BlockTimeoutCert = types.BoolValue(ans.SslForwardProxy.BlockTimeoutCert)
		var5.BlockTls13DowngradeNoResource = types.BoolValue(ans.SslForwardProxy.BlockTls13DowngradeNoResource)
		var5.BlockUnknownCert = types.BoolValue(ans.SslForwardProxy.BlockUnknownCert)
		var5.BlockUnsupportedCipher = types.BoolValue(ans.SslForwardProxy.BlockUnsupportedCipher)
		var5.BlockUnsupportedVersion = types.BoolValue(ans.SslForwardProxy.BlockUnsupportedVersion)
		var5.BlockUntrustedIssuer = types.BoolValue(ans.SslForwardProxy.BlockUntrustedIssuer)
		var5.RestrictCertExts = types.BoolValue(ans.SslForwardProxy.RestrictCertExts)
		var5.StripAlpn = types.BoolValue(ans.SslForwardProxy.StripAlpn)
	}
	var var6 *decryptionProfilesRsModelSslInboundProxyObject
	if ans.SslInboundProxy != nil {
		var6 = &decryptionProfilesRsModelSslInboundProxyObject{}
		var6.BlockIfHsmUnavailable = types.BoolValue(ans.SslInboundProxy.BlockIfHsmUnavailable)
		var6.BlockIfNoResource = types.BoolValue(ans.SslInboundProxy.BlockIfNoResource)
		var6.BlockUnsupportedCipher = types.BoolValue(ans.SslInboundProxy.BlockUnsupportedCipher)
		var6.BlockUnsupportedVersion = types.BoolValue(ans.SslInboundProxy.BlockUnsupportedVersion)
	}
	var var7 *decryptionProfilesRsModelSslNoProxyObject
	if ans.SslNoProxy != nil {
		var7 = &decryptionProfilesRsModelSslNoProxyObject{}
		var7.BlockExpiredCertificate = types.BoolValue(ans.SslNoProxy.BlockExpiredCertificate)
		var7.BlockUntrustedIssuer = types.BoolValue(ans.SslNoProxy.BlockUntrustedIssuer)
	}
	var var8 *decryptionProfilesRsModelSslProtocolSettingsObject
	if ans.SslProtocolSettings != nil {
		var8 = &decryptionProfilesRsModelSslProtocolSettingsObject{}
		var8.AuthAlgoMd5 = types.BoolValue(ans.SslProtocolSettings.AuthAlgoMd5)
		var8.AuthAlgoSha1 = types.BoolValue(ans.SslProtocolSettings.AuthAlgoSha1)
		var8.AuthAlgoSha256 = types.BoolValue(ans.SslProtocolSettings.AuthAlgoSha256)
		var8.AuthAlgoSha384 = types.BoolValue(ans.SslProtocolSettings.AuthAlgoSha384)
		var8.EncAlgo3des = types.BoolValue(ans.SslProtocolSettings.EncAlgo3des)
		var8.EncAlgoAes128Cbc = types.BoolValue(ans.SslProtocolSettings.EncAlgoAes128Cbc)
		var8.EncAlgoAes128Gcm = types.BoolValue(ans.SslProtocolSettings.EncAlgoAes128Gcm)
		var8.EncAlgoAes256Cbc = types.BoolValue(ans.SslProtocolSettings.EncAlgoAes256Cbc)
		var8.EncAlgoAes256Gcm = types.BoolValue(ans.SslProtocolSettings.EncAlgoAes256Gcm)
		var8.EncAlgoChacha20Poly1305 = types.BoolValue(ans.SslProtocolSettings.EncAlgoChacha20Poly1305)
		var8.EncAlgoRc4 = types.BoolValue(ans.SslProtocolSettings.EncAlgoRc4)
		var8.KeyxchgAlgoDhe = types.BoolValue(ans.SslProtocolSettings.KeyxchgAlgoDhe)
		var8.KeyxchgAlgoEcdhe = types.BoolValue(ans.SslProtocolSettings.KeyxchgAlgoEcdhe)
		var8.KeyxchgAlgoRsa = types.BoolValue(ans.SslProtocolSettings.KeyxchgAlgoRsa)
		var8.MaxVersion = types.StringValue(ans.SslProtocolSettings.MaxVersion)
		var8.MinVersion = types.StringValue(ans.SslProtocolSettings.MinVersion)
	}
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.SslForwardProxy = var5
	state.SslInboundProxy = var6
	state.SslNoProxy = var7
	state.SslProtocolSettings = var8

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *decryptionProfilesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name":               "sase_decryption_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	svc := bpgvUeD.NewClient(r.client)
	input := bpgvUeD.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *decryptionProfilesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
