package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	zAHtTyI "github.com/paloaltonetworks/sase-go/netsec/schema/ike/gateways"
	fGoRZph "github.com/paloaltonetworks/sase-go/netsec/service/v1/ikegateways"

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
	_ datasource.DataSource              = &ikeGatewaysListDataSource{}
	_ datasource.DataSourceWithConfigure = &ikeGatewaysListDataSource{}
)

func NewIkeGatewaysListDataSource() datasource.DataSource {
	return &ikeGatewaysListDataSource{}
}

type ikeGatewaysListDataSource struct {
	client *sase.Client
}

type ikeGatewaysListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []ikeGatewaysListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type ikeGatewaysListDsModelConfig struct {
	Authentication ikeGatewaysListDsModelAuthenticationObject  `tfsdk:"authentication"`
	ObjectId       types.String                                `tfsdk:"object_id"`
	LocalId        *ikeGatewaysListDsModelLocalIdObject        `tfsdk:"local_id"`
	Name           types.String                                `tfsdk:"name"`
	PeerAddress    ikeGatewaysListDsModelPeerAddressObject     `tfsdk:"peer_address"`
	PeerId         *ikeGatewaysListDsModelPeerIdObject         `tfsdk:"peer_id"`
	Protocol       ikeGatewaysListDsModelProtocolObject        `tfsdk:"protocol"`
	ProtocolCommon *ikeGatewaysListDsModelProtocolCommonObject `tfsdk:"protocol_common"`
}

type ikeGatewaysListDsModelAuthenticationObject struct {
	AllowIdPayloadMismatch     types.Bool                                    `tfsdk:"allow_id_payload_mismatch"`
	CertificateProfile         types.String                                  `tfsdk:"certificate_profile"`
	LocalCertificate           *ikeGatewaysListDsModelLocalCertificateObject `tfsdk:"local_certificate"`
	PreSharedKey               *ikeGatewaysListDsModelPreSharedKeyObject     `tfsdk:"pre_shared_key"`
	StrictValidationRevocation types.Bool                                    `tfsdk:"strict_validation_revocation"`
	UseManagementAsSource      types.Bool                                    `tfsdk:"use_management_as_source"`
}

type ikeGatewaysListDsModelLocalCertificateObject struct {
	LocalCertificateName types.String `tfsdk:"local_certificate_name"`
}

type ikeGatewaysListDsModelPreSharedKeyObject struct {
	Key types.String `tfsdk:"key"`
}

type ikeGatewaysListDsModelLocalIdObject struct {
	ObjectId types.String `tfsdk:"object_id"`
	Type     types.String `tfsdk:"type"`
}

type ikeGatewaysListDsModelPeerAddressObject struct {
	DynamicValue types.Bool   `tfsdk:"dynamic_value"`
	Fqdn         types.String `tfsdk:"fqdn"`
	Ip           types.String `tfsdk:"ip"`
}

type ikeGatewaysListDsModelPeerIdObject struct {
	ObjectId types.String `tfsdk:"object_id"`
	Type     types.String `tfsdk:"type"`
}

type ikeGatewaysListDsModelProtocolObject struct {
	Ikev1   *ikeGatewaysListDsModelIkev1Object `tfsdk:"ikev1"`
	Ikev2   *ikeGatewaysListDsModelIkev2Object `tfsdk:"ikev2"`
	Version types.String                       `tfsdk:"version"`
}

type ikeGatewaysListDsModelIkev1Object struct {
	Dpd              *ikeGatewaysListDsModelDpdObject `tfsdk:"dpd"`
	IkeCryptoProfile types.String                     `tfsdk:"ike_crypto_profile"`
}

type ikeGatewaysListDsModelDpdObject struct {
	Enable types.Bool `tfsdk:"enable"`
}

type ikeGatewaysListDsModelIkev2Object struct {
	Dpd              *ikeGatewaysListDsModelDpdObject `tfsdk:"dpd"`
	IkeCryptoProfile types.String                     `tfsdk:"ike_crypto_profile"`
}

type ikeGatewaysListDsModelProtocolCommonObject struct {
	Fragmentation *ikeGatewaysListDsModelFragmentationObject `tfsdk:"fragmentation"`
	NatTraversal  *ikeGatewaysListDsModelNatTraversalObject  `tfsdk:"nat_traversal"`
	PassiveMode   types.Bool                                 `tfsdk:"passive_mode"`
}

type ikeGatewaysListDsModelFragmentationObject struct {
	Enable types.Bool `tfsdk:"enable"`
}

type ikeGatewaysListDsModelNatTraversalObject struct {
	Enable types.Bool `tfsdk:"enable"`
}

// Metadata returns the data source type name.
func (d *ikeGatewaysListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ike_gateways_list"
}

// Schema defines the schema for this listing data source.
func (d *ikeGatewaysListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"authentication": dsschema.SingleNestedAttribute{
							Description:         "The `authentication` parameter.",
							MarkdownDescription: "The `authentication` parameter.",
							Computed:            true,
							Attributes: map[string]dsschema.Attribute{
								"allow_id_payload_mismatch": dsschema.BoolAttribute{
									Description:         "The `allow_id_payload_mismatch` parameter.",
									MarkdownDescription: "The `allow_id_payload_mismatch` parameter.",
									Computed:            true,
								},
								"certificate_profile": dsschema.StringAttribute{
									Description:         "The `certificate_profile` parameter.",
									MarkdownDescription: "The `certificate_profile` parameter.",
									Computed:            true,
								},
								"local_certificate": dsschema.SingleNestedAttribute{
									Description:         "The `local_certificate` parameter.",
									MarkdownDescription: "The `local_certificate` parameter.",
									Computed:            true,
									Attributes: map[string]dsschema.Attribute{
										"local_certificate_name": dsschema.StringAttribute{
											Description:         "The `local_certificate_name` parameter.",
											MarkdownDescription: "The `local_certificate_name` parameter.",
											Computed:            true,
										},
									},
								},
								"pre_shared_key": dsschema.SingleNestedAttribute{
									Description:         "The `pre_shared_key` parameter.",
									MarkdownDescription: "The `pre_shared_key` parameter.",
									Computed:            true,
									Attributes: map[string]dsschema.Attribute{
										"key": dsschema.StringAttribute{
											Description:         "The `key` parameter.",
											MarkdownDescription: "The `key` parameter.",
											Computed:            true,
										},
									},
								},
								"strict_validation_revocation": dsschema.BoolAttribute{
									Description:         "The `strict_validation_revocation` parameter.",
									MarkdownDescription: "The `strict_validation_revocation` parameter.",
									Computed:            true,
								},
								"use_management_as_source": dsschema.BoolAttribute{
									Description:         "The `use_management_as_source` parameter.",
									MarkdownDescription: "The `use_management_as_source` parameter.",
									Computed:            true,
								},
							},
						},
						"object_id": dsschema.StringAttribute{
							Description:         "The `object_id` parameter.",
							MarkdownDescription: "The `object_id` parameter.",
							Computed:            true,
						},
						"local_id": dsschema.SingleNestedAttribute{
							Description:         "The `local_id` parameter.",
							MarkdownDescription: "The `local_id` parameter.",
							Computed:            true,
							Attributes: map[string]dsschema.Attribute{
								"object_id": dsschema.StringAttribute{
									Description:         "The `object_id` parameter.",
									MarkdownDescription: "The `object_id` parameter.",
									Computed:            true,
								},
								"type": dsschema.StringAttribute{
									Description:         "The `type` parameter.",
									MarkdownDescription: "The `type` parameter.",
									Computed:            true,
								},
							},
						},
						"name": dsschema.StringAttribute{
							Description:         "The `name` parameter.",
							MarkdownDescription: "The `name` parameter.",
							Computed:            true,
						},
						"peer_address": dsschema.SingleNestedAttribute{
							Description:         "The `peer_address` parameter.",
							MarkdownDescription: "The `peer_address` parameter.",
							Computed:            true,
							Attributes: map[string]dsschema.Attribute{
								"dynamic_value": dsschema.BoolAttribute{
									Description:         "The `dynamic_value` parameter.",
									MarkdownDescription: "The `dynamic_value` parameter.",
									Computed:            true,
								},
								"fqdn": dsschema.StringAttribute{
									Description:         "The `fqdn` parameter.",
									MarkdownDescription: "The `fqdn` parameter.",
									Computed:            true,
								},
								"ip": dsschema.StringAttribute{
									Description:         "The `ip` parameter.",
									MarkdownDescription: "The `ip` parameter.",
									Computed:            true,
								},
							},
						},
						"peer_id": dsschema.SingleNestedAttribute{
							Description:         "The `peer_id` parameter.",
							MarkdownDescription: "The `peer_id` parameter.",
							Computed:            true,
							Attributes: map[string]dsschema.Attribute{
								"object_id": dsschema.StringAttribute{
									Description:         "The `object_id` parameter.",
									MarkdownDescription: "The `object_id` parameter.",
									Computed:            true,
								},
								"type": dsschema.StringAttribute{
									Description:         "The `type` parameter.",
									MarkdownDescription: "The `type` parameter.",
									Computed:            true,
								},
							},
						},
						"protocol": dsschema.SingleNestedAttribute{
							Description:         "The `protocol` parameter.",
							MarkdownDescription: "The `protocol` parameter.",
							Computed:            true,
							Attributes: map[string]dsschema.Attribute{
								"ikev1": dsschema.SingleNestedAttribute{
									Description:         "The `ikev1` parameter.",
									MarkdownDescription: "The `ikev1` parameter.",
									Computed:            true,
									Attributes: map[string]dsschema.Attribute{
										"dpd": dsschema.SingleNestedAttribute{
											Description:         "The `dpd` parameter.",
											MarkdownDescription: "The `dpd` parameter.",
											Computed:            true,
											Attributes: map[string]dsschema.Attribute{
												"enable": dsschema.BoolAttribute{
													Description:         "The `enable` parameter.",
													MarkdownDescription: "The `enable` parameter.",
													Computed:            true,
												},
											},
										},
										"ike_crypto_profile": dsschema.StringAttribute{
											Description:         "The `ike_crypto_profile` parameter.",
											MarkdownDescription: "The `ike_crypto_profile` parameter.",
											Computed:            true,
										},
									},
								},
								"ikev2": dsschema.SingleNestedAttribute{
									Description:         "The `ikev2` parameter.",
									MarkdownDescription: "The `ikev2` parameter.",
									Computed:            true,
									Attributes: map[string]dsschema.Attribute{
										"dpd": dsschema.SingleNestedAttribute{
											Description:         "The `dpd` parameter.",
											MarkdownDescription: "The `dpd` parameter.",
											Computed:            true,
											Attributes: map[string]dsschema.Attribute{
												"enable": dsschema.BoolAttribute{
													Description:         "The `enable` parameter.",
													MarkdownDescription: "The `enable` parameter.",
													Computed:            true,
												},
											},
										},
										"ike_crypto_profile": dsschema.StringAttribute{
											Description:         "The `ike_crypto_profile` parameter.",
											MarkdownDescription: "The `ike_crypto_profile` parameter.",
											Computed:            true,
										},
									},
								},
								"version": dsschema.StringAttribute{
									Description:         "The `version` parameter.",
									MarkdownDescription: "The `version` parameter.",
									Computed:            true,
								},
							},
						},
						"protocol_common": dsschema.SingleNestedAttribute{
							Description:         "The `protocol_common` parameter.",
							MarkdownDescription: "The `protocol_common` parameter.",
							Computed:            true,
							Attributes: map[string]dsschema.Attribute{
								"fragmentation": dsschema.SingleNestedAttribute{
									Description:         "The `fragmentation` parameter.",
									MarkdownDescription: "The `fragmentation` parameter.",
									Computed:            true,
									Attributes: map[string]dsschema.Attribute{
										"enable": dsschema.BoolAttribute{
											Description:         "The `enable` parameter.",
											MarkdownDescription: "The `enable` parameter.",
											Computed:            true,
										},
									},
								},
								"nat_traversal": dsschema.SingleNestedAttribute{
									Description:         "The `nat_traversal` parameter.",
									MarkdownDescription: "The `nat_traversal` parameter.",
									Computed:            true,
									Attributes: map[string]dsschema.Attribute{
										"enable": dsschema.BoolAttribute{
											Description:         "The `enable` parameter.",
											MarkdownDescription: "The `enable` parameter.",
											Computed:            true,
										},
									},
								},
								"passive_mode": dsschema.BoolAttribute{
									Description:         "The `passive_mode` parameter.",
									MarkdownDescription: "The `passive_mode` parameter.",
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
func (d *ikeGatewaysListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *ikeGatewaysListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ikeGatewaysListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_ike_gateways_list",
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
	svc := fGoRZph.NewClient(d.client)
	input := fGoRZph.ListInput{
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
	var var0 []ikeGatewaysListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]ikeGatewaysListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 ikeGatewaysListDsModelConfig
			var var3 ikeGatewaysListDsModelAuthenticationObject
			var var4 *ikeGatewaysListDsModelLocalCertificateObject
			if var1.Authentication.LocalCertificate != nil {
				var4 = &ikeGatewaysListDsModelLocalCertificateObject{}
				var4.LocalCertificateName = types.StringValue(var1.Authentication.LocalCertificate.LocalCertificateName)
			}
			var var5 *ikeGatewaysListDsModelPreSharedKeyObject
			if var1.Authentication.PreSharedKey != nil {
				var5 = &ikeGatewaysListDsModelPreSharedKeyObject{}
				var5.Key = types.StringValue(var1.Authentication.PreSharedKey.Key)
			}
			var3.AllowIdPayloadMismatch = types.BoolValue(var1.Authentication.AllowIdPayloadMismatch)
			var3.CertificateProfile = types.StringValue(var1.Authentication.CertificateProfile)
			var3.LocalCertificate = var4
			var3.PreSharedKey = var5
			var3.StrictValidationRevocation = types.BoolValue(var1.Authentication.StrictValidationRevocation)
			var3.UseManagementAsSource = types.BoolValue(var1.Authentication.UseManagementAsSource)
			var var6 *ikeGatewaysListDsModelLocalIdObject
			if var1.LocalId != nil {
				var6 = &ikeGatewaysListDsModelLocalIdObject{}
				var6.ObjectId = types.StringValue(var1.LocalId.ObjectId)
				var6.Type = types.StringValue(var1.LocalId.Type)
			}
			var var7 ikeGatewaysListDsModelPeerAddressObject
			if var1.PeerAddress.DynamicValue != nil {
				var7.DynamicValue = types.BoolValue(true)
			}
			var7.Fqdn = types.StringValue(var1.PeerAddress.Fqdn)
			var7.Ip = types.StringValue(var1.PeerAddress.Ip)
			var var8 *ikeGatewaysListDsModelPeerIdObject
			if var1.PeerId != nil {
				var8 = &ikeGatewaysListDsModelPeerIdObject{}
				var8.ObjectId = types.StringValue(var1.PeerId.ObjectId)
				var8.Type = types.StringValue(var1.PeerId.Type)
			}
			var var9 ikeGatewaysListDsModelProtocolObject
			var var10 *ikeGatewaysListDsModelIkev1Object
			if var1.Protocol.Ikev1 != nil {
				var10 = &ikeGatewaysListDsModelIkev1Object{}
				var var11 *ikeGatewaysListDsModelDpdObject
				if var1.Protocol.Ikev1.Dpd != nil {
					var11 = &ikeGatewaysListDsModelDpdObject{}
					var11.Enable = types.BoolValue(var1.Protocol.Ikev1.Dpd.Enable)
				}
				var10.Dpd = var11
				var10.IkeCryptoProfile = types.StringValue(var1.Protocol.Ikev1.IkeCryptoProfile)
			}
			var var12 *ikeGatewaysListDsModelIkev2Object
			if var1.Protocol.Ikev2 != nil {
				var12 = &ikeGatewaysListDsModelIkev2Object{}
				var var13 *ikeGatewaysListDsModelDpdObject
				if var1.Protocol.Ikev2.Dpd != nil {
					var13 = &ikeGatewaysListDsModelDpdObject{}
					var13.Enable = types.BoolValue(var1.Protocol.Ikev2.Dpd.Enable)
				}
				var12.Dpd = var13
				var12.IkeCryptoProfile = types.StringValue(var1.Protocol.Ikev2.IkeCryptoProfile)
			}
			var9.Ikev1 = var10
			var9.Ikev2 = var12
			var9.Version = types.StringValue(var1.Protocol.Version)
			var var14 *ikeGatewaysListDsModelProtocolCommonObject
			if var1.ProtocolCommon != nil {
				var14 = &ikeGatewaysListDsModelProtocolCommonObject{}
				var var15 *ikeGatewaysListDsModelFragmentationObject
				if var1.ProtocolCommon.Fragmentation != nil {
					var15 = &ikeGatewaysListDsModelFragmentationObject{}
					var15.Enable = types.BoolValue(var1.ProtocolCommon.Fragmentation.Enable)
				}
				var var16 *ikeGatewaysListDsModelNatTraversalObject
				if var1.ProtocolCommon.NatTraversal != nil {
					var16 = &ikeGatewaysListDsModelNatTraversalObject{}
					var16.Enable = types.BoolValue(var1.ProtocolCommon.NatTraversal.Enable)
				}
				var14.Fragmentation = var15
				var14.NatTraversal = var16
				var14.PassiveMode = types.BoolValue(var1.ProtocolCommon.PassiveMode)
			}
			var2.Authentication = var3
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.LocalId = var6
			var2.Name = types.StringValue(var1.Name)
			var2.PeerAddress = var7
			var2.PeerId = var8
			var2.Protocol = var9
			var2.ProtocolCommon = var14
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
	_ datasource.DataSource              = &ikeGatewaysDataSource{}
	_ datasource.DataSourceWithConfigure = &ikeGatewaysDataSource{}
)

func NewIkeGatewaysDataSource() datasource.DataSource {
	return &ikeGatewaysDataSource{}
}

type ikeGatewaysDataSource struct {
	client *sase.Client
}

type ikeGatewaysDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`
	Folder   types.String `tfsdk:"folder"`

	// Output.
	// Ref: #/components/schemas/ike-gateways
	Authentication ikeGatewaysDsModelAuthenticationObject `tfsdk:"authentication"`
	// input omit: ObjectId
	LocalId        *ikeGatewaysDsModelLocalIdObject        `tfsdk:"local_id"`
	Name           types.String                            `tfsdk:"name"`
	PeerAddress    ikeGatewaysDsModelPeerAddressObject     `tfsdk:"peer_address"`
	PeerId         *ikeGatewaysDsModelPeerIdObject         `tfsdk:"peer_id"`
	Protocol       ikeGatewaysDsModelProtocolObject        `tfsdk:"protocol"`
	ProtocolCommon *ikeGatewaysDsModelProtocolCommonObject `tfsdk:"protocol_common"`
}

type ikeGatewaysDsModelAuthenticationObject struct {
	AllowIdPayloadMismatch     types.Bool                                `tfsdk:"allow_id_payload_mismatch"`
	CertificateProfile         types.String                              `tfsdk:"certificate_profile"`
	LocalCertificate           *ikeGatewaysDsModelLocalCertificateObject `tfsdk:"local_certificate"`
	PreSharedKey               *ikeGatewaysDsModelPreSharedKeyObject     `tfsdk:"pre_shared_key"`
	StrictValidationRevocation types.Bool                                `tfsdk:"strict_validation_revocation"`
	UseManagementAsSource      types.Bool                                `tfsdk:"use_management_as_source"`
}

type ikeGatewaysDsModelLocalCertificateObject struct {
	LocalCertificateName types.String `tfsdk:"local_certificate_name"`
}

type ikeGatewaysDsModelPreSharedKeyObject struct {
	Key types.String `tfsdk:"key"`
}

type ikeGatewaysDsModelLocalIdObject struct {
	ObjectId types.String `tfsdk:"object_id"`
	Type     types.String `tfsdk:"type"`
}

type ikeGatewaysDsModelPeerAddressObject struct {
	DynamicValue types.Bool   `tfsdk:"dynamic_value"`
	Fqdn         types.String `tfsdk:"fqdn"`
	Ip           types.String `tfsdk:"ip"`
}

type ikeGatewaysDsModelPeerIdObject struct {
	ObjectId types.String `tfsdk:"object_id"`
	Type     types.String `tfsdk:"type"`
}

type ikeGatewaysDsModelProtocolObject struct {
	Ikev1   *ikeGatewaysDsModelIkev1Object `tfsdk:"ikev1"`
	Ikev2   *ikeGatewaysDsModelIkev2Object `tfsdk:"ikev2"`
	Version types.String                   `tfsdk:"version"`
}

type ikeGatewaysDsModelIkev1Object struct {
	Dpd              *ikeGatewaysDsModelDpdObject `tfsdk:"dpd"`
	IkeCryptoProfile types.String                 `tfsdk:"ike_crypto_profile"`
}

type ikeGatewaysDsModelDpdObject struct {
	Enable types.Bool `tfsdk:"enable"`
}

type ikeGatewaysDsModelIkev2Object struct {
	Dpd              *ikeGatewaysDsModelDpdObject `tfsdk:"dpd"`
	IkeCryptoProfile types.String                 `tfsdk:"ike_crypto_profile"`
}

type ikeGatewaysDsModelProtocolCommonObject struct {
	Fragmentation *ikeGatewaysDsModelFragmentationObject `tfsdk:"fragmentation"`
	NatTraversal  *ikeGatewaysDsModelNatTraversalObject  `tfsdk:"nat_traversal"`
	PassiveMode   types.Bool                             `tfsdk:"passive_mode"`
}

type ikeGatewaysDsModelFragmentationObject struct {
	Enable types.Bool `tfsdk:"enable"`
}

type ikeGatewaysDsModelNatTraversalObject struct {
	Enable types.Bool `tfsdk:"enable"`
}

// Metadata returns the data source type name.
func (d *ikeGatewaysDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ike_gateways"
}

// Schema defines the schema for this listing data source.
func (d *ikeGatewaysDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"authentication": dsschema.SingleNestedAttribute{
				Description:         "The `authentication` parameter.",
				MarkdownDescription: "The `authentication` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"allow_id_payload_mismatch": dsschema.BoolAttribute{
						Description:         "The `allow_id_payload_mismatch` parameter.",
						MarkdownDescription: "The `allow_id_payload_mismatch` parameter.",
						Computed:            true,
					},
					"certificate_profile": dsschema.StringAttribute{
						Description:         "The `certificate_profile` parameter.",
						MarkdownDescription: "The `certificate_profile` parameter.",
						Computed:            true,
					},
					"local_certificate": dsschema.SingleNestedAttribute{
						Description:         "The `local_certificate` parameter.",
						MarkdownDescription: "The `local_certificate` parameter.",
						Computed:            true,
						Attributes: map[string]dsschema.Attribute{
							"local_certificate_name": dsschema.StringAttribute{
								Description:         "The `local_certificate_name` parameter.",
								MarkdownDescription: "The `local_certificate_name` parameter.",
								Computed:            true,
							},
						},
					},
					"pre_shared_key": dsschema.SingleNestedAttribute{
						Description:         "The `pre_shared_key` parameter.",
						MarkdownDescription: "The `pre_shared_key` parameter.",
						Computed:            true,
						Attributes: map[string]dsschema.Attribute{
							"key": dsschema.StringAttribute{
								Description:         "The `key` parameter.",
								MarkdownDescription: "The `key` parameter.",
								Computed:            true,
							},
						},
					},
					"strict_validation_revocation": dsschema.BoolAttribute{
						Description:         "The `strict_validation_revocation` parameter.",
						MarkdownDescription: "The `strict_validation_revocation` parameter.",
						Computed:            true,
					},
					"use_management_as_source": dsschema.BoolAttribute{
						Description:         "The `use_management_as_source` parameter.",
						MarkdownDescription: "The `use_management_as_source` parameter.",
						Computed:            true,
					},
				},
			},
			"local_id": dsschema.SingleNestedAttribute{
				Description:         "The `local_id` parameter.",
				MarkdownDescription: "The `local_id` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"object_id": dsschema.StringAttribute{
						Description:         "The `object_id` parameter.",
						MarkdownDescription: "The `object_id` parameter.",
						Computed:            true,
					},
					"type": dsschema.StringAttribute{
						Description:         "The `type` parameter.",
						MarkdownDescription: "The `type` parameter.",
						Computed:            true,
					},
				},
			},
			"name": dsschema.StringAttribute{
				Description:         "The `name` parameter.",
				MarkdownDescription: "The `name` parameter.",
				Computed:            true,
			},
			"peer_address": dsschema.SingleNestedAttribute{
				Description:         "The `peer_address` parameter.",
				MarkdownDescription: "The `peer_address` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"dynamic_value": dsschema.BoolAttribute{
						Description:         "The `dynamic_value` parameter.",
						MarkdownDescription: "The `dynamic_value` parameter.",
						Computed:            true,
					},
					"fqdn": dsschema.StringAttribute{
						Description:         "The `fqdn` parameter.",
						MarkdownDescription: "The `fqdn` parameter.",
						Computed:            true,
					},
					"ip": dsschema.StringAttribute{
						Description:         "The `ip` parameter.",
						MarkdownDescription: "The `ip` parameter.",
						Computed:            true,
					},
				},
			},
			"peer_id": dsschema.SingleNestedAttribute{
				Description:         "The `peer_id` parameter.",
				MarkdownDescription: "The `peer_id` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"object_id": dsschema.StringAttribute{
						Description:         "The `object_id` parameter.",
						MarkdownDescription: "The `object_id` parameter.",
						Computed:            true,
					},
					"type": dsschema.StringAttribute{
						Description:         "The `type` parameter.",
						MarkdownDescription: "The `type` parameter.",
						Computed:            true,
					},
				},
			},
			"protocol": dsschema.SingleNestedAttribute{
				Description:         "The `protocol` parameter.",
				MarkdownDescription: "The `protocol` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"ikev1": dsschema.SingleNestedAttribute{
						Description:         "The `ikev1` parameter.",
						MarkdownDescription: "The `ikev1` parameter.",
						Computed:            true,
						Attributes: map[string]dsschema.Attribute{
							"dpd": dsschema.SingleNestedAttribute{
								Description:         "The `dpd` parameter.",
								MarkdownDescription: "The `dpd` parameter.",
								Computed:            true,
								Attributes: map[string]dsschema.Attribute{
									"enable": dsschema.BoolAttribute{
										Description:         "The `enable` parameter.",
										MarkdownDescription: "The `enable` parameter.",
										Computed:            true,
									},
								},
							},
							"ike_crypto_profile": dsschema.StringAttribute{
								Description:         "The `ike_crypto_profile` parameter.",
								MarkdownDescription: "The `ike_crypto_profile` parameter.",
								Computed:            true,
							},
						},
					},
					"ikev2": dsschema.SingleNestedAttribute{
						Description:         "The `ikev2` parameter.",
						MarkdownDescription: "The `ikev2` parameter.",
						Computed:            true,
						Attributes: map[string]dsschema.Attribute{
							"dpd": dsschema.SingleNestedAttribute{
								Description:         "The `dpd` parameter.",
								MarkdownDescription: "The `dpd` parameter.",
								Computed:            true,
								Attributes: map[string]dsschema.Attribute{
									"enable": dsschema.BoolAttribute{
										Description:         "The `enable` parameter.",
										MarkdownDescription: "The `enable` parameter.",
										Computed:            true,
									},
								},
							},
							"ike_crypto_profile": dsschema.StringAttribute{
								Description:         "The `ike_crypto_profile` parameter.",
								MarkdownDescription: "The `ike_crypto_profile` parameter.",
								Computed:            true,
							},
						},
					},
					"version": dsschema.StringAttribute{
						Description:         "The `version` parameter.",
						MarkdownDescription: "The `version` parameter.",
						Computed:            true,
					},
				},
			},
			"protocol_common": dsschema.SingleNestedAttribute{
				Description:         "The `protocol_common` parameter.",
				MarkdownDescription: "The `protocol_common` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"fragmentation": dsschema.SingleNestedAttribute{
						Description:         "The `fragmentation` parameter.",
						MarkdownDescription: "The `fragmentation` parameter.",
						Computed:            true,
						Attributes: map[string]dsschema.Attribute{
							"enable": dsschema.BoolAttribute{
								Description:         "The `enable` parameter.",
								MarkdownDescription: "The `enable` parameter.",
								Computed:            true,
							},
						},
					},
					"nat_traversal": dsschema.SingleNestedAttribute{
						Description:         "The `nat_traversal` parameter.",
						MarkdownDescription: "The `nat_traversal` parameter.",
						Computed:            true,
						Attributes: map[string]dsschema.Attribute{
							"enable": dsschema.BoolAttribute{
								Description:         "The `enable` parameter.",
								MarkdownDescription: "The `enable` parameter.",
								Computed:            true,
							},
						},
					},
					"passive_mode": dsschema.BoolAttribute{
						Description:         "The `passive_mode` parameter.",
						MarkdownDescription: "The `passive_mode` parameter.",
						Computed:            true,
					},
				},
			},
		},
	}
}

// Configure prepares the struct.
func (d *ikeGatewaysDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *ikeGatewaysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ikeGatewaysDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_ike_gateways",
		"object_id":                   state.ObjectId.ValueString(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := fGoRZph.NewClient(d.client)
	input := fGoRZph.ReadInput{
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
	var var0 ikeGatewaysDsModelAuthenticationObject
	var var1 *ikeGatewaysDsModelLocalCertificateObject
	if ans.Authentication.LocalCertificate != nil {
		var1 = &ikeGatewaysDsModelLocalCertificateObject{}
		var1.LocalCertificateName = types.StringValue(ans.Authentication.LocalCertificate.LocalCertificateName)
	}
	var var2 *ikeGatewaysDsModelPreSharedKeyObject
	if ans.Authentication.PreSharedKey != nil {
		var2 = &ikeGatewaysDsModelPreSharedKeyObject{}
		var2.Key = types.StringValue(ans.Authentication.PreSharedKey.Key)
	}
	var0.AllowIdPayloadMismatch = types.BoolValue(ans.Authentication.AllowIdPayloadMismatch)
	var0.CertificateProfile = types.StringValue(ans.Authentication.CertificateProfile)
	var0.LocalCertificate = var1
	var0.PreSharedKey = var2
	var0.StrictValidationRevocation = types.BoolValue(ans.Authentication.StrictValidationRevocation)
	var0.UseManagementAsSource = types.BoolValue(ans.Authentication.UseManagementAsSource)
	var var3 *ikeGatewaysDsModelLocalIdObject
	if ans.LocalId != nil {
		var3 = &ikeGatewaysDsModelLocalIdObject{}
		var3.ObjectId = types.StringValue(ans.LocalId.ObjectId)
		var3.Type = types.StringValue(ans.LocalId.Type)
	}
	var var4 ikeGatewaysDsModelPeerAddressObject
	if ans.PeerAddress.DynamicValue != nil {
		var4.DynamicValue = types.BoolValue(true)
	}
	var4.Fqdn = types.StringValue(ans.PeerAddress.Fqdn)
	var4.Ip = types.StringValue(ans.PeerAddress.Ip)
	var var5 *ikeGatewaysDsModelPeerIdObject
	if ans.PeerId != nil {
		var5 = &ikeGatewaysDsModelPeerIdObject{}
		var5.ObjectId = types.StringValue(ans.PeerId.ObjectId)
		var5.Type = types.StringValue(ans.PeerId.Type)
	}
	var var6 ikeGatewaysDsModelProtocolObject
	var var7 *ikeGatewaysDsModelIkev1Object
	if ans.Protocol.Ikev1 != nil {
		var7 = &ikeGatewaysDsModelIkev1Object{}
		var var8 *ikeGatewaysDsModelDpdObject
		if ans.Protocol.Ikev1.Dpd != nil {
			var8 = &ikeGatewaysDsModelDpdObject{}
			var8.Enable = types.BoolValue(ans.Protocol.Ikev1.Dpd.Enable)
		}
		var7.Dpd = var8
		var7.IkeCryptoProfile = types.StringValue(ans.Protocol.Ikev1.IkeCryptoProfile)
	}
	var var9 *ikeGatewaysDsModelIkev2Object
	if ans.Protocol.Ikev2 != nil {
		var9 = &ikeGatewaysDsModelIkev2Object{}
		var var10 *ikeGatewaysDsModelDpdObject
		if ans.Protocol.Ikev2.Dpd != nil {
			var10 = &ikeGatewaysDsModelDpdObject{}
			var10.Enable = types.BoolValue(ans.Protocol.Ikev2.Dpd.Enable)
		}
		var9.Dpd = var10
		var9.IkeCryptoProfile = types.StringValue(ans.Protocol.Ikev2.IkeCryptoProfile)
	}
	var6.Ikev1 = var7
	var6.Ikev2 = var9
	var6.Version = types.StringValue(ans.Protocol.Version)
	var var11 *ikeGatewaysDsModelProtocolCommonObject
	if ans.ProtocolCommon != nil {
		var11 = &ikeGatewaysDsModelProtocolCommonObject{}
		var var12 *ikeGatewaysDsModelFragmentationObject
		if ans.ProtocolCommon.Fragmentation != nil {
			var12 = &ikeGatewaysDsModelFragmentationObject{}
			var12.Enable = types.BoolValue(ans.ProtocolCommon.Fragmentation.Enable)
		}
		var var13 *ikeGatewaysDsModelNatTraversalObject
		if ans.ProtocolCommon.NatTraversal != nil {
			var13 = &ikeGatewaysDsModelNatTraversalObject{}
			var13.Enable = types.BoolValue(ans.ProtocolCommon.NatTraversal.Enable)
		}
		var11.Fragmentation = var12
		var11.NatTraversal = var13
		var11.PassiveMode = types.BoolValue(ans.ProtocolCommon.PassiveMode)
	}
	state.Authentication = var0
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.LocalId = var3
	state.Name = types.StringValue(ans.Name)
	state.PeerAddress = var4
	state.PeerId = var5
	state.Protocol = var6
	state.ProtocolCommon = var11

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &ikeGatewaysResource{}
	_ resource.ResourceWithConfigure   = &ikeGatewaysResource{}
	_ resource.ResourceWithImportState = &ikeGatewaysResource{}
)

func NewIkeGatewaysResource() resource.Resource {
	return &ikeGatewaysResource{}
}

type ikeGatewaysResource struct {
	client *sase.Client
}

type ikeGatewaysRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/ike-gateways
	Authentication ikeGatewaysRsModelAuthenticationObject  `tfsdk:"authentication"`
	ObjectId       types.String                            `tfsdk:"object_id"`
	LocalId        *ikeGatewaysRsModelLocalIdObject        `tfsdk:"local_id"`
	Name           types.String                            `tfsdk:"name"`
	PeerAddress    ikeGatewaysRsModelPeerAddressObject     `tfsdk:"peer_address"`
	PeerId         *ikeGatewaysRsModelPeerIdObject         `tfsdk:"peer_id"`
	Protocol       ikeGatewaysRsModelProtocolObject        `tfsdk:"protocol"`
	ProtocolCommon *ikeGatewaysRsModelProtocolCommonObject `tfsdk:"protocol_common"`
}

type ikeGatewaysRsModelAuthenticationObject struct {
	AllowIdPayloadMismatch     types.Bool                                `tfsdk:"allow_id_payload_mismatch"`
	CertificateProfile         types.String                              `tfsdk:"certificate_profile"`
	LocalCertificate           *ikeGatewaysRsModelLocalCertificateObject `tfsdk:"local_certificate"`
	PreSharedKey               *ikeGatewaysRsModelPreSharedKeyObject     `tfsdk:"pre_shared_key"`
	StrictValidationRevocation types.Bool                                `tfsdk:"strict_validation_revocation"`
	UseManagementAsSource      types.Bool                                `tfsdk:"use_management_as_source"`
}

type ikeGatewaysRsModelLocalCertificateObject struct {
	LocalCertificateName types.String `tfsdk:"local_certificate_name"`
}

type ikeGatewaysRsModelPreSharedKeyObject struct {
	Key types.String `tfsdk:"key"`
}

type ikeGatewaysRsModelLocalIdObject struct {
	ObjectId types.String `tfsdk:"object_id"`
	Type     types.String `tfsdk:"type"`
}

type ikeGatewaysRsModelPeerAddressObject struct {
	DynamicValue types.Bool   `tfsdk:"dynamic_value"`
	Fqdn         types.String `tfsdk:"fqdn"`
	Ip           types.String `tfsdk:"ip"`
}

type ikeGatewaysRsModelPeerIdObject struct {
	ObjectId types.String `tfsdk:"object_id"`
	Type     types.String `tfsdk:"type"`
}

type ikeGatewaysRsModelProtocolObject struct {
	Ikev1   *ikeGatewaysRsModelIkev1Object `tfsdk:"ikev1"`
	Ikev2   *ikeGatewaysRsModelIkev2Object `tfsdk:"ikev2"`
	Version types.String                   `tfsdk:"version"`
}

type ikeGatewaysRsModelIkev1Object struct {
	Dpd              *ikeGatewaysRsModelDpdObject `tfsdk:"dpd"`
	IkeCryptoProfile types.String                 `tfsdk:"ike_crypto_profile"`
}

type ikeGatewaysRsModelDpdObject struct {
	Enable types.Bool `tfsdk:"enable"`
}

type ikeGatewaysRsModelIkev2Object struct {
	Dpd              *ikeGatewaysRsModelDpdObject `tfsdk:"dpd"`
	IkeCryptoProfile types.String                 `tfsdk:"ike_crypto_profile"`
}

type ikeGatewaysRsModelProtocolCommonObject struct {
	Fragmentation *ikeGatewaysRsModelFragmentationObject `tfsdk:"fragmentation"`
	NatTraversal  *ikeGatewaysRsModelNatTraversalObject  `tfsdk:"nat_traversal"`
	PassiveMode   types.Bool                             `tfsdk:"passive_mode"`
}

type ikeGatewaysRsModelFragmentationObject struct {
	Enable types.Bool `tfsdk:"enable"`
}

type ikeGatewaysRsModelNatTraversalObject struct {
	Enable types.Bool `tfsdk:"enable"`
}

// Metadata returns the data source type name.
func (r *ikeGatewaysResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ike_gateways"
}

// Schema defines the schema for this listing data source.
func (r *ikeGatewaysResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

			"authentication": rsschema.SingleNestedAttribute{
				Description:         "The `authentication` parameter.",
				MarkdownDescription: "The `authentication` parameter.",
				Required:            true,
				Attributes: map[string]rsschema.Attribute{
					"allow_id_payload_mismatch": rsschema.BoolAttribute{
						Description:         "The `allow_id_payload_mismatch` parameter.",
						MarkdownDescription: "The `allow_id_payload_mismatch` parameter.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"certificate_profile": rsschema.StringAttribute{
						Description:         "The `certificate_profile` parameter.",
						MarkdownDescription: "The `certificate_profile` parameter.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							DefaultString(""),
						},
					},
					"local_certificate": rsschema.SingleNestedAttribute{
						Description:         "The `local_certificate` parameter.",
						MarkdownDescription: "The `local_certificate` parameter.",
						Optional:            true,
						Attributes: map[string]rsschema.Attribute{
							"local_certificate_name": rsschema.StringAttribute{
								Description:         "The `local_certificate_name` parameter.",
								MarkdownDescription: "The `local_certificate_name` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
						},
					},
					"pre_shared_key": rsschema.SingleNestedAttribute{
						Description:         "The `pre_shared_key` parameter.",
						MarkdownDescription: "The `pre_shared_key` parameter.",
						Optional:            true,
						Attributes: map[string]rsschema.Attribute{
							"key": rsschema.StringAttribute{
								Description:         "The `key` parameter.",
								MarkdownDescription: "The `key` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
						},
					},
					"strict_validation_revocation": rsschema.BoolAttribute{
						Description:         "The `strict_validation_revocation` parameter.",
						MarkdownDescription: "The `strict_validation_revocation` parameter.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"use_management_as_source": rsschema.BoolAttribute{
						Description:         "The `use_management_as_source` parameter.",
						MarkdownDescription: "The `use_management_as_source` parameter.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
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
			"local_id": rsschema.SingleNestedAttribute{
				Description:         "The `local_id` parameter.",
				MarkdownDescription: "The `local_id` parameter.",
				Optional:            true,
				Attributes: map[string]rsschema.Attribute{
					"object_id": rsschema.StringAttribute{
						Description:         "The `object_id` parameter. String length must be between 1 and 1024.",
						MarkdownDescription: "The `object_id` parameter. String length must be between 1 and 1024.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							DefaultString(""),
						},
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 1024),
						},
					},
					"type": rsschema.StringAttribute{
						Description:         "The `type` parameter.",
						MarkdownDescription: "The `type` parameter.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							DefaultString(""),
						},
					},
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
			"peer_address": rsschema.SingleNestedAttribute{
				Description:         "The `peer_address` parameter.",
				MarkdownDescription: "The `peer_address` parameter.",
				Required:            true,
				Attributes: map[string]rsschema.Attribute{
					"dynamic_value": rsschema.BoolAttribute{
						Description:         "The `dynamic_value` parameter.",
						MarkdownDescription: "The `dynamic_value` parameter.",
						Optional:            true,
					},
					"fqdn": rsschema.StringAttribute{
						Description:         "The `fqdn` parameter. String length must be at most 255.",
						MarkdownDescription: "The `fqdn` parameter. String length must be at most 255.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							DefaultString(""),
						},
						Validators: []validator.String{
							stringvalidator.LengthAtMost(255),
						},
					},
					"ip": rsschema.StringAttribute{
						Description:         "The `ip` parameter.",
						MarkdownDescription: "The `ip` parameter.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							DefaultString(""),
						},
					},
				},
			},
			"peer_id": rsschema.SingleNestedAttribute{
				Description:         "The `peer_id` parameter.",
				MarkdownDescription: "The `peer_id` parameter.",
				Optional:            true,
				Attributes: map[string]rsschema.Attribute{
					"object_id": rsschema.StringAttribute{
						Description:         "The `object_id` parameter. String length must be between 1 and 1024.",
						MarkdownDescription: "The `object_id` parameter. String length must be between 1 and 1024.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							DefaultString(""),
						},
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 1024),
						},
					},
					"type": rsschema.StringAttribute{
						Description:         "The `type` parameter. Value must be one of: `\"ipaddr\"`, `\"keyid\"`, `\"fqdn\"`, `\"ufqdn\"`.",
						MarkdownDescription: "The `type` parameter. Value must be one of: `\"ipaddr\"`, `\"keyid\"`, `\"fqdn\"`, `\"ufqdn\"`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							DefaultString(""),
						},
						Validators: []validator.String{
							stringvalidator.OneOf("ipaddr", "keyid", "fqdn", "ufqdn"),
						},
					},
				},
			},
			"protocol": rsschema.SingleNestedAttribute{
				Description:         "The `protocol` parameter.",
				MarkdownDescription: "The `protocol` parameter.",
				Required:            true,
				Attributes: map[string]rsschema.Attribute{
					"ikev1": rsschema.SingleNestedAttribute{
						Description:         "The `ikev1` parameter.",
						MarkdownDescription: "The `ikev1` parameter.",
						Optional:            true,
						Attributes: map[string]rsschema.Attribute{
							"dpd": rsschema.SingleNestedAttribute{
								Description:         "The `dpd` parameter.",
								MarkdownDescription: "The `dpd` parameter.",
								Optional:            true,
								Attributes: map[string]rsschema.Attribute{
									"enable": rsschema.BoolAttribute{
										Description:         "The `enable` parameter.",
										MarkdownDescription: "The `enable` parameter.",
										Optional:            true,
										Computed:            true,
										PlanModifiers: []planmodifier.Bool{
											DefaultBool(false),
										},
									},
								},
							},
							"ike_crypto_profile": rsschema.StringAttribute{
								Description:         "The `ike_crypto_profile` parameter.",
								MarkdownDescription: "The `ike_crypto_profile` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
						},
					},
					"ikev2": rsschema.SingleNestedAttribute{
						Description:         "The `ikev2` parameter.",
						MarkdownDescription: "The `ikev2` parameter.",
						Optional:            true,
						Attributes: map[string]rsschema.Attribute{
							"dpd": rsschema.SingleNestedAttribute{
								Description:         "The `dpd` parameter.",
								MarkdownDescription: "The `dpd` parameter.",
								Optional:            true,
								Attributes: map[string]rsschema.Attribute{
									"enable": rsschema.BoolAttribute{
										Description:         "The `enable` parameter.",
										MarkdownDescription: "The `enable` parameter.",
										Optional:            true,
										Computed:            true,
										PlanModifiers: []planmodifier.Bool{
											DefaultBool(false),
										},
									},
								},
							},
							"ike_crypto_profile": rsschema.StringAttribute{
								Description:         "The `ike_crypto_profile` parameter.",
								MarkdownDescription: "The `ike_crypto_profile` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
						},
					},
					"version": rsschema.StringAttribute{
						Description:         "The `version` parameter. Default: `\"ikev2-preferred\"`. Value must be one of: `\"ikev2-preferred\"`, `\"ikev1\"`, `\"ikev2\"`.",
						MarkdownDescription: "The `version` parameter. Default: `\"ikev2-preferred\"`. Value must be one of: `\"ikev2-preferred\"`, `\"ikev1\"`, `\"ikev2\"`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							DefaultString("ikev2-preferred"),
						},
						Validators: []validator.String{
							stringvalidator.OneOf("ikev2-preferred", "ikev1", "ikev2"),
						},
					},
				},
			},
			"protocol_common": rsschema.SingleNestedAttribute{
				Description:         "The `protocol_common` parameter.",
				MarkdownDescription: "The `protocol_common` parameter.",
				Optional:            true,
				Attributes: map[string]rsschema.Attribute{
					"fragmentation": rsschema.SingleNestedAttribute{
						Description:         "The `fragmentation` parameter.",
						MarkdownDescription: "The `fragmentation` parameter.",
						Optional:            true,
						Attributes: map[string]rsschema.Attribute{
							"enable": rsschema.BoolAttribute{
								Description:         "The `enable` parameter. Default: `false`.",
								MarkdownDescription: "The `enable` parameter. Default: `false`.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.Bool{
									DefaultBool(false),
								},
							},
						},
					},
					"nat_traversal": rsschema.SingleNestedAttribute{
						Description:         "The `nat_traversal` parameter.",
						MarkdownDescription: "The `nat_traversal` parameter.",
						Optional:            true,
						Attributes: map[string]rsschema.Attribute{
							"enable": rsschema.BoolAttribute{
								Description:         "The `enable` parameter.",
								MarkdownDescription: "The `enable` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.Bool{
									DefaultBool(false),
								},
							},
						},
					},
					"passive_mode": rsschema.BoolAttribute{
						Description:         "The `passive_mode` parameter.",
						MarkdownDescription: "The `passive_mode` parameter.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *ikeGatewaysResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *ikeGatewaysResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state ikeGatewaysRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_ike_gateways",
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := fGoRZph.NewClient(r.client)
	input := fGoRZph.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 zAHtTyI.Config
	var var1 zAHtTyI.AuthenticationObject
	var1.AllowIdPayloadMismatch = state.Authentication.AllowIdPayloadMismatch.ValueBool()
	var1.CertificateProfile = state.Authentication.CertificateProfile.ValueString()
	var var2 *zAHtTyI.LocalCertificateObject
	if state.Authentication.LocalCertificate != nil {
		var2 = &zAHtTyI.LocalCertificateObject{}
		var2.LocalCertificateName = state.Authentication.LocalCertificate.LocalCertificateName.ValueString()
	}
	var1.LocalCertificate = var2
	var var3 *zAHtTyI.PreSharedKeyObject
	if state.Authentication.PreSharedKey != nil {
		var3 = &zAHtTyI.PreSharedKeyObject{}
		var3.Key = state.Authentication.PreSharedKey.Key.ValueString()
	}
	var1.PreSharedKey = var3
	var1.StrictValidationRevocation = state.Authentication.StrictValidationRevocation.ValueBool()
	var1.UseManagementAsSource = state.Authentication.UseManagementAsSource.ValueBool()
	var0.Authentication = var1
	var var4 *zAHtTyI.LocalIdObject
	if state.LocalId != nil {
		var4 = &zAHtTyI.LocalIdObject{}
		var4.ObjectId = state.LocalId.ObjectId.ValueString()
		var4.Type = state.LocalId.Type.ValueString()
	}
	var0.LocalId = var4
	var0.Name = state.Name.ValueString()
	var var5 zAHtTyI.PeerAddressObject
	if state.PeerAddress.DynamicValue.ValueBool() {
		var5.DynamicValue = struct{}{}
	}
	var5.Fqdn = state.PeerAddress.Fqdn.ValueString()
	var5.Ip = state.PeerAddress.Ip.ValueString()
	var0.PeerAddress = var5
	var var6 *zAHtTyI.PeerIdObject
	if state.PeerId != nil {
		var6 = &zAHtTyI.PeerIdObject{}
		var6.ObjectId = state.PeerId.ObjectId.ValueString()
		var6.Type = state.PeerId.Type.ValueString()
	}
	var0.PeerId = var6
	var var7 zAHtTyI.ProtocolObject
	var var8 *zAHtTyI.Ikev1Object
	if state.Protocol.Ikev1 != nil {
		var8 = &zAHtTyI.Ikev1Object{}
		var var9 *zAHtTyI.DpdObject
		if state.Protocol.Ikev1.Dpd != nil {
			var9 = &zAHtTyI.DpdObject{}
			var9.Enable = state.Protocol.Ikev1.Dpd.Enable.ValueBool()
		}
		var8.Dpd = var9
		var8.IkeCryptoProfile = state.Protocol.Ikev1.IkeCryptoProfile.ValueString()
	}
	var7.Ikev1 = var8
	var var10 *zAHtTyI.Ikev2Object
	if state.Protocol.Ikev2 != nil {
		var10 = &zAHtTyI.Ikev2Object{}
		var var11 *zAHtTyI.DpdObject
		if state.Protocol.Ikev2.Dpd != nil {
			var11 = &zAHtTyI.DpdObject{}
			var11.Enable = state.Protocol.Ikev2.Dpd.Enable.ValueBool()
		}
		var10.Dpd = var11
		var10.IkeCryptoProfile = state.Protocol.Ikev2.IkeCryptoProfile.ValueString()
	}
	var7.Ikev2 = var10
	var7.Version = state.Protocol.Version.ValueString()
	var0.Protocol = var7
	var var12 *zAHtTyI.ProtocolCommonObject
	if state.ProtocolCommon != nil {
		var12 = &zAHtTyI.ProtocolCommonObject{}
		var var13 *zAHtTyI.FragmentationObject
		if state.ProtocolCommon.Fragmentation != nil {
			var13 = &zAHtTyI.FragmentationObject{}
			var13.Enable = state.ProtocolCommon.Fragmentation.Enable.ValueBool()
		}
		var12.Fragmentation = var13
		var var14 *zAHtTyI.NatTraversalObject
		if state.ProtocolCommon.NatTraversal != nil {
			var14 = &zAHtTyI.NatTraversalObject{}
			var14.Enable = state.ProtocolCommon.NatTraversal.Enable.ValueBool()
		}
		var12.NatTraversal = var14
		var12.PassiveMode = state.ProtocolCommon.PassiveMode.ValueBool()
	}
	var0.ProtocolCommon = var12
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
	var var15 ikeGatewaysRsModelAuthenticationObject
	var var16 *ikeGatewaysRsModelLocalCertificateObject
	if ans.Authentication.LocalCertificate != nil {
		var16 = &ikeGatewaysRsModelLocalCertificateObject{}
		var16.LocalCertificateName = types.StringValue(ans.Authentication.LocalCertificate.LocalCertificateName)
	}
	var var17 *ikeGatewaysRsModelPreSharedKeyObject
	if ans.Authentication.PreSharedKey != nil {
		var17 = &ikeGatewaysRsModelPreSharedKeyObject{}
		var17.Key = types.StringValue(ans.Authentication.PreSharedKey.Key)
	}
	var15.AllowIdPayloadMismatch = types.BoolValue(ans.Authentication.AllowIdPayloadMismatch)
	var15.CertificateProfile = types.StringValue(ans.Authentication.CertificateProfile)
	var15.LocalCertificate = var16
	var15.PreSharedKey = var17
	var15.StrictValidationRevocation = types.BoolValue(ans.Authentication.StrictValidationRevocation)
	var15.UseManagementAsSource = types.BoolValue(ans.Authentication.UseManagementAsSource)
	var var18 *ikeGatewaysRsModelLocalIdObject
	if ans.LocalId != nil {
		var18 = &ikeGatewaysRsModelLocalIdObject{}
		var18.ObjectId = types.StringValue(ans.LocalId.ObjectId)
		var18.Type = types.StringValue(ans.LocalId.Type)
	}
	var var19 ikeGatewaysRsModelPeerAddressObject
	if ans.PeerAddress.DynamicValue != nil {
		var19.DynamicValue = types.BoolValue(true)
	}
	var19.Fqdn = types.StringValue(ans.PeerAddress.Fqdn)
	var19.Ip = types.StringValue(ans.PeerAddress.Ip)
	var var20 *ikeGatewaysRsModelPeerIdObject
	if ans.PeerId != nil {
		var20 = &ikeGatewaysRsModelPeerIdObject{}
		var20.ObjectId = types.StringValue(ans.PeerId.ObjectId)
		var20.Type = types.StringValue(ans.PeerId.Type)
	}
	var var21 ikeGatewaysRsModelProtocolObject
	var var22 *ikeGatewaysRsModelIkev1Object
	if ans.Protocol.Ikev1 != nil {
		var22 = &ikeGatewaysRsModelIkev1Object{}
		var var23 *ikeGatewaysRsModelDpdObject
		if ans.Protocol.Ikev1.Dpd != nil {
			var23 = &ikeGatewaysRsModelDpdObject{}
			var23.Enable = types.BoolValue(ans.Protocol.Ikev1.Dpd.Enable)
		}
		var22.Dpd = var23
		var22.IkeCryptoProfile = types.StringValue(ans.Protocol.Ikev1.IkeCryptoProfile)
	}
	var var24 *ikeGatewaysRsModelIkev2Object
	if ans.Protocol.Ikev2 != nil {
		var24 = &ikeGatewaysRsModelIkev2Object{}
		var var25 *ikeGatewaysRsModelDpdObject
		if ans.Protocol.Ikev2.Dpd != nil {
			var25 = &ikeGatewaysRsModelDpdObject{}
			var25.Enable = types.BoolValue(ans.Protocol.Ikev2.Dpd.Enable)
		}
		var24.Dpd = var25
		var24.IkeCryptoProfile = types.StringValue(ans.Protocol.Ikev2.IkeCryptoProfile)
	}
	var21.Ikev1 = var22
	var21.Ikev2 = var24
	var21.Version = types.StringValue(ans.Protocol.Version)
	var var26 *ikeGatewaysRsModelProtocolCommonObject
	if ans.ProtocolCommon != nil {
		var26 = &ikeGatewaysRsModelProtocolCommonObject{}
		var var27 *ikeGatewaysRsModelFragmentationObject
		if ans.ProtocolCommon.Fragmentation != nil {
			var27 = &ikeGatewaysRsModelFragmentationObject{}
			var27.Enable = types.BoolValue(ans.ProtocolCommon.Fragmentation.Enable)
		}
		var var28 *ikeGatewaysRsModelNatTraversalObject
		if ans.ProtocolCommon.NatTraversal != nil {
			var28 = &ikeGatewaysRsModelNatTraversalObject{}
			var28.Enable = types.BoolValue(ans.ProtocolCommon.NatTraversal.Enable)
		}
		var26.Fragmentation = var27
		var26.NatTraversal = var28
		var26.PassiveMode = types.BoolValue(ans.ProtocolCommon.PassiveMode)
	}
	state.Authentication = var15
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.LocalId = var18
	state.Name = types.StringValue(ans.Name)
	state.PeerAddress = var19
	state.PeerId = var20
	state.Protocol = var21
	state.ProtocolCommon = var26

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *ikeGatewaysResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state ikeGatewaysRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_ike_gateways",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := fGoRZph.NewClient(r.client)
	input := fGoRZph.ReadInput{
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
	var var0 ikeGatewaysRsModelAuthenticationObject
	var var1 *ikeGatewaysRsModelLocalCertificateObject
	if ans.Authentication.LocalCertificate != nil {
		var1 = &ikeGatewaysRsModelLocalCertificateObject{}
		var1.LocalCertificateName = types.StringValue(ans.Authentication.LocalCertificate.LocalCertificateName)
	}
	var var2 *ikeGatewaysRsModelPreSharedKeyObject
	if ans.Authentication.PreSharedKey != nil {
		var2 = &ikeGatewaysRsModelPreSharedKeyObject{}
		var2.Key = types.StringValue(ans.Authentication.PreSharedKey.Key)
	}
	var0.AllowIdPayloadMismatch = types.BoolValue(ans.Authentication.AllowIdPayloadMismatch)
	var0.CertificateProfile = types.StringValue(ans.Authentication.CertificateProfile)
	var0.LocalCertificate = var1
	var0.PreSharedKey = var2
	var0.StrictValidationRevocation = types.BoolValue(ans.Authentication.StrictValidationRevocation)
	var0.UseManagementAsSource = types.BoolValue(ans.Authentication.UseManagementAsSource)
	var var3 *ikeGatewaysRsModelLocalIdObject
	if ans.LocalId != nil {
		var3 = &ikeGatewaysRsModelLocalIdObject{}
		var3.ObjectId = types.StringValue(ans.LocalId.ObjectId)
		var3.Type = types.StringValue(ans.LocalId.Type)
	}
	var var4 ikeGatewaysRsModelPeerAddressObject
	if ans.PeerAddress.DynamicValue != nil {
		var4.DynamicValue = types.BoolValue(true)
	}
	var4.Fqdn = types.StringValue(ans.PeerAddress.Fqdn)
	var4.Ip = types.StringValue(ans.PeerAddress.Ip)
	var var5 *ikeGatewaysRsModelPeerIdObject
	if ans.PeerId != nil {
		var5 = &ikeGatewaysRsModelPeerIdObject{}
		var5.ObjectId = types.StringValue(ans.PeerId.ObjectId)
		var5.Type = types.StringValue(ans.PeerId.Type)
	}
	var var6 ikeGatewaysRsModelProtocolObject
	var var7 *ikeGatewaysRsModelIkev1Object
	if ans.Protocol.Ikev1 != nil {
		var7 = &ikeGatewaysRsModelIkev1Object{}
		var var8 *ikeGatewaysRsModelDpdObject
		if ans.Protocol.Ikev1.Dpd != nil {
			var8 = &ikeGatewaysRsModelDpdObject{}
			var8.Enable = types.BoolValue(ans.Protocol.Ikev1.Dpd.Enable)
		}
		var7.Dpd = var8
		var7.IkeCryptoProfile = types.StringValue(ans.Protocol.Ikev1.IkeCryptoProfile)
	}
	var var9 *ikeGatewaysRsModelIkev2Object
	if ans.Protocol.Ikev2 != nil {
		var9 = &ikeGatewaysRsModelIkev2Object{}
		var var10 *ikeGatewaysRsModelDpdObject
		if ans.Protocol.Ikev2.Dpd != nil {
			var10 = &ikeGatewaysRsModelDpdObject{}
			var10.Enable = types.BoolValue(ans.Protocol.Ikev2.Dpd.Enable)
		}
		var9.Dpd = var10
		var9.IkeCryptoProfile = types.StringValue(ans.Protocol.Ikev2.IkeCryptoProfile)
	}
	var6.Ikev1 = var7
	var6.Ikev2 = var9
	var6.Version = types.StringValue(ans.Protocol.Version)
	var var11 *ikeGatewaysRsModelProtocolCommonObject
	if ans.ProtocolCommon != nil {
		var11 = &ikeGatewaysRsModelProtocolCommonObject{}
		var var12 *ikeGatewaysRsModelFragmentationObject
		if ans.ProtocolCommon.Fragmentation != nil {
			var12 = &ikeGatewaysRsModelFragmentationObject{}
			var12.Enable = types.BoolValue(ans.ProtocolCommon.Fragmentation.Enable)
		}
		var var13 *ikeGatewaysRsModelNatTraversalObject
		if ans.ProtocolCommon.NatTraversal != nil {
			var13 = &ikeGatewaysRsModelNatTraversalObject{}
			var13.Enable = types.BoolValue(ans.ProtocolCommon.NatTraversal.Enable)
		}
		var11.Fragmentation = var12
		var11.NatTraversal = var13
		var11.PassiveMode = types.BoolValue(ans.ProtocolCommon.PassiveMode)
	}
	state.Authentication = var0
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.LocalId = var3
	state.Name = types.StringValue(ans.Name)
	state.PeerAddress = var4
	state.PeerId = var5
	state.Protocol = var6
	state.ProtocolCommon = var11

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *ikeGatewaysResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ikeGatewaysRsModel
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
		"resource_name":               "sase_ike_gateways",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := fGoRZph.NewClient(r.client)
	input := fGoRZph.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 zAHtTyI.Config
	var var1 zAHtTyI.AuthenticationObject
	var1.AllowIdPayloadMismatch = plan.Authentication.AllowIdPayloadMismatch.ValueBool()
	var1.CertificateProfile = plan.Authentication.CertificateProfile.ValueString()
	var var2 *zAHtTyI.LocalCertificateObject
	if plan.Authentication.LocalCertificate != nil {
		var2 = &zAHtTyI.LocalCertificateObject{}
		var2.LocalCertificateName = plan.Authentication.LocalCertificate.LocalCertificateName.ValueString()
	}
	var1.LocalCertificate = var2
	var var3 *zAHtTyI.PreSharedKeyObject
	if plan.Authentication.PreSharedKey != nil {
		var3 = &zAHtTyI.PreSharedKeyObject{}
		var3.Key = plan.Authentication.PreSharedKey.Key.ValueString()
	}
	var1.PreSharedKey = var3
	var1.StrictValidationRevocation = plan.Authentication.StrictValidationRevocation.ValueBool()
	var1.UseManagementAsSource = plan.Authentication.UseManagementAsSource.ValueBool()
	var0.Authentication = var1
	var var4 *zAHtTyI.LocalIdObject
	if plan.LocalId != nil {
		var4 = &zAHtTyI.LocalIdObject{}
		var4.ObjectId = plan.LocalId.ObjectId.ValueString()
		var4.Type = plan.LocalId.Type.ValueString()
	}
	var0.LocalId = var4
	var0.Name = plan.Name.ValueString()
	var var5 zAHtTyI.PeerAddressObject
	if plan.PeerAddress.DynamicValue.ValueBool() {
		var5.DynamicValue = struct{}{}
	}
	var5.Fqdn = plan.PeerAddress.Fqdn.ValueString()
	var5.Ip = plan.PeerAddress.Ip.ValueString()
	var0.PeerAddress = var5
	var var6 *zAHtTyI.PeerIdObject
	if plan.PeerId != nil {
		var6 = &zAHtTyI.PeerIdObject{}
		var6.ObjectId = plan.PeerId.ObjectId.ValueString()
		var6.Type = plan.PeerId.Type.ValueString()
	}
	var0.PeerId = var6
	var var7 zAHtTyI.ProtocolObject
	var var8 *zAHtTyI.Ikev1Object
	if plan.Protocol.Ikev1 != nil {
		var8 = &zAHtTyI.Ikev1Object{}
		var var9 *zAHtTyI.DpdObject
		if plan.Protocol.Ikev1.Dpd != nil {
			var9 = &zAHtTyI.DpdObject{}
			var9.Enable = plan.Protocol.Ikev1.Dpd.Enable.ValueBool()
		}
		var8.Dpd = var9
		var8.IkeCryptoProfile = plan.Protocol.Ikev1.IkeCryptoProfile.ValueString()
	}
	var7.Ikev1 = var8
	var var10 *zAHtTyI.Ikev2Object
	if plan.Protocol.Ikev2 != nil {
		var10 = &zAHtTyI.Ikev2Object{}
		var var11 *zAHtTyI.DpdObject
		if plan.Protocol.Ikev2.Dpd != nil {
			var11 = &zAHtTyI.DpdObject{}
			var11.Enable = plan.Protocol.Ikev2.Dpd.Enable.ValueBool()
		}
		var10.Dpd = var11
		var10.IkeCryptoProfile = plan.Protocol.Ikev2.IkeCryptoProfile.ValueString()
	}
	var7.Ikev2 = var10
	var7.Version = plan.Protocol.Version.ValueString()
	var0.Protocol = var7
	var var12 *zAHtTyI.ProtocolCommonObject
	if plan.ProtocolCommon != nil {
		var12 = &zAHtTyI.ProtocolCommonObject{}
		var var13 *zAHtTyI.FragmentationObject
		if plan.ProtocolCommon.Fragmentation != nil {
			var13 = &zAHtTyI.FragmentationObject{}
			var13.Enable = plan.ProtocolCommon.Fragmentation.Enable.ValueBool()
		}
		var12.Fragmentation = var13
		var var14 *zAHtTyI.NatTraversalObject
		if plan.ProtocolCommon.NatTraversal != nil {
			var14 = &zAHtTyI.NatTraversalObject{}
			var14.Enable = plan.ProtocolCommon.NatTraversal.Enable.ValueBool()
		}
		var12.NatTraversal = var14
		var12.PassiveMode = plan.ProtocolCommon.PassiveMode.ValueBool()
	}
	var0.ProtocolCommon = var12
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var15 ikeGatewaysRsModelAuthenticationObject
	var var16 *ikeGatewaysRsModelLocalCertificateObject
	if ans.Authentication.LocalCertificate != nil {
		var16 = &ikeGatewaysRsModelLocalCertificateObject{}
		var16.LocalCertificateName = types.StringValue(ans.Authentication.LocalCertificate.LocalCertificateName)
	}
	var var17 *ikeGatewaysRsModelPreSharedKeyObject
	if ans.Authentication.PreSharedKey != nil {
		var17 = &ikeGatewaysRsModelPreSharedKeyObject{}
		var17.Key = types.StringValue(ans.Authentication.PreSharedKey.Key)
	}
	var15.AllowIdPayloadMismatch = types.BoolValue(ans.Authentication.AllowIdPayloadMismatch)
	var15.CertificateProfile = types.StringValue(ans.Authentication.CertificateProfile)
	var15.LocalCertificate = var16
	var15.PreSharedKey = var17
	var15.StrictValidationRevocation = types.BoolValue(ans.Authentication.StrictValidationRevocation)
	var15.UseManagementAsSource = types.BoolValue(ans.Authentication.UseManagementAsSource)
	var var18 *ikeGatewaysRsModelLocalIdObject
	if ans.LocalId != nil {
		var18 = &ikeGatewaysRsModelLocalIdObject{}
		var18.ObjectId = types.StringValue(ans.LocalId.ObjectId)
		var18.Type = types.StringValue(ans.LocalId.Type)
	}
	var var19 ikeGatewaysRsModelPeerAddressObject
	if ans.PeerAddress.DynamicValue != nil {
		var19.DynamicValue = types.BoolValue(true)
	}
	var19.Fqdn = types.StringValue(ans.PeerAddress.Fqdn)
	var19.Ip = types.StringValue(ans.PeerAddress.Ip)
	var var20 *ikeGatewaysRsModelPeerIdObject
	if ans.PeerId != nil {
		var20 = &ikeGatewaysRsModelPeerIdObject{}
		var20.ObjectId = types.StringValue(ans.PeerId.ObjectId)
		var20.Type = types.StringValue(ans.PeerId.Type)
	}
	var var21 ikeGatewaysRsModelProtocolObject
	var var22 *ikeGatewaysRsModelIkev1Object
	if ans.Protocol.Ikev1 != nil {
		var22 = &ikeGatewaysRsModelIkev1Object{}
		var var23 *ikeGatewaysRsModelDpdObject
		if ans.Protocol.Ikev1.Dpd != nil {
			var23 = &ikeGatewaysRsModelDpdObject{}
			var23.Enable = types.BoolValue(ans.Protocol.Ikev1.Dpd.Enable)
		}
		var22.Dpd = var23
		var22.IkeCryptoProfile = types.StringValue(ans.Protocol.Ikev1.IkeCryptoProfile)
	}
	var var24 *ikeGatewaysRsModelIkev2Object
	if ans.Protocol.Ikev2 != nil {
		var24 = &ikeGatewaysRsModelIkev2Object{}
		var var25 *ikeGatewaysRsModelDpdObject
		if ans.Protocol.Ikev2.Dpd != nil {
			var25 = &ikeGatewaysRsModelDpdObject{}
			var25.Enable = types.BoolValue(ans.Protocol.Ikev2.Dpd.Enable)
		}
		var24.Dpd = var25
		var24.IkeCryptoProfile = types.StringValue(ans.Protocol.Ikev2.IkeCryptoProfile)
	}
	var21.Ikev1 = var22
	var21.Ikev2 = var24
	var21.Version = types.StringValue(ans.Protocol.Version)
	var var26 *ikeGatewaysRsModelProtocolCommonObject
	if ans.ProtocolCommon != nil {
		var26 = &ikeGatewaysRsModelProtocolCommonObject{}
		var var27 *ikeGatewaysRsModelFragmentationObject
		if ans.ProtocolCommon.Fragmentation != nil {
			var27 = &ikeGatewaysRsModelFragmentationObject{}
			var27.Enable = types.BoolValue(ans.ProtocolCommon.Fragmentation.Enable)
		}
		var var28 *ikeGatewaysRsModelNatTraversalObject
		if ans.ProtocolCommon.NatTraversal != nil {
			var28 = &ikeGatewaysRsModelNatTraversalObject{}
			var28.Enable = types.BoolValue(ans.ProtocolCommon.NatTraversal.Enable)
		}
		var26.Fragmentation = var27
		var26.NatTraversal = var28
		var26.PassiveMode = types.BoolValue(ans.ProtocolCommon.PassiveMode)
	}
	state.Authentication = var15
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.LocalId = var18
	state.Name = types.StringValue(ans.Name)
	state.PeerAddress = var19
	state.PeerId = var20
	state.Protocol = var21
	state.ProtocolCommon = var26

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *ikeGatewaysResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name":               "sase_ike_gateways",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	svc := fGoRZph.NewClient(r.client)
	input := fGoRZph.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *ikeGatewaysResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
