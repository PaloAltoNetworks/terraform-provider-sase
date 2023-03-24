package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	qmgDayz "github.com/paloaltonetworks/sase-go/netsec/schema/scep/profiles"
	xlSkOUa "github.com/paloaltonetworks/sase-go/netsec/service/v1/scepprofiles"

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
	_ datasource.DataSource              = &scepProfilesListDataSource{}
	_ datasource.DataSourceWithConfigure = &scepProfilesListDataSource{}
)

func NewScepProfilesListDataSource() datasource.DataSource {
	return &scepProfilesListDataSource{}
}

type scepProfilesListDataSource struct {
	client *sase.Client
}

type scepProfilesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []scepProfilesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type scepProfilesListDsModelConfig struct {
	Algorithm             *scepProfilesListDsModelAlgorithmObject             `tfsdk:"algorithm"`
	CaIdentityName        types.String                                        `tfsdk:"ca_identity_name"`
	CertificateAttributes *scepProfilesListDsModelCertificateAttributesObject `tfsdk:"certificate_attributes"`
	Digest                types.String                                        `tfsdk:"digest"`
	Fingerprint           types.String                                        `tfsdk:"fingerprint"`
	ObjectId              types.String                                        `tfsdk:"object_id"`
	Name                  types.String                                        `tfsdk:"name"`
	ScepCaCert            types.String                                        `tfsdk:"scep_ca_cert"`
	ScepChallenge         *scepProfilesListDsModelScepChallengeObject         `tfsdk:"scep_challenge"`
	ScepClientCert        types.String                                        `tfsdk:"scep_client_cert"`
	ScepUrl               types.String                                        `tfsdk:"scep_url"`
	Subject               types.String                                        `tfsdk:"subject"`
	UseAsDigitalSignature types.Bool                                          `tfsdk:"use_as_digital_signature"`
	UseForKeyEncipherment types.Bool                                          `tfsdk:"use_for_key_encipherment"`
}

type scepProfilesListDsModelAlgorithmObject struct {
	Rsa *scepProfilesListDsModelRsaObject `tfsdk:"rsa"`
}

type scepProfilesListDsModelRsaObject struct {
	RsaNbits types.String `tfsdk:"rsa_nbits"`
}

type scepProfilesListDsModelCertificateAttributesObject struct {
	Dnsname                   types.String `tfsdk:"dnsname"`
	Rfc822name                types.String `tfsdk:"rfc822name"`
	UniformResourceIdentifier types.String `tfsdk:"uniform_resource_identifier"`
}

type scepProfilesListDsModelScepChallengeObject struct {
	DynamicValue *scepProfilesListDsModelDynamicObject `tfsdk:"dynamic_value"`
	Fixed        types.String                          `tfsdk:"fixed"`
	None         types.String                          `tfsdk:"none"`
}

type scepProfilesListDsModelDynamicObject struct {
	OtpServerUrl types.String `tfsdk:"otp_server_url"`
	Password     types.String `tfsdk:"password"`
	Username     types.String `tfsdk:"username"`
}

// Metadata returns the data source type name.
func (d *scepProfilesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scep_profiles_list"
}

// Schema defines the schema for this listing data source.
func (d *scepProfilesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
				Description:         "The max count in result entry (count per page)",
				MarkdownDescription: "The max count in result entry (count per page)",
				Optional:            true,
				Computed:            true,
			},
			"offset": dsschema.Int64Attribute{
				Description:         "The offset of the result entry",
				MarkdownDescription: "The offset of the result entry",
				Optional:            true,
				Computed:            true,
			},
			"name": dsschema.StringAttribute{
				Description:         "The name of the entry",
				MarkdownDescription: "The name of the entry",
				Optional:            true,
			},
			"folder": dsschema.StringAttribute{
				Description:         "The folder of the entry",
				MarkdownDescription: "The folder of the entry",
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
						"algorithm": dsschema.SingleNestedAttribute{
							Description:         "The `algorithm` parameter.",
							MarkdownDescription: "The `algorithm` parameter.",
							Computed:            true,
							Attributes: map[string]dsschema.Attribute{
								"rsa": dsschema.SingleNestedAttribute{
									Description:         "The `rsa` parameter.",
									MarkdownDescription: "The `rsa` parameter.",
									Computed:            true,
									Attributes: map[string]dsschema.Attribute{
										"rsa_nbits": dsschema.StringAttribute{
											Description:         "The `rsa_nbits` parameter.",
											MarkdownDescription: "The `rsa_nbits` parameter.",
											Computed:            true,
										},
									},
								},
							},
						},
						"ca_identity_name": dsschema.StringAttribute{
							Description:         "The `ca_identity_name` parameter.",
							MarkdownDescription: "The `ca_identity_name` parameter.",
							Computed:            true,
						},
						"certificate_attributes": dsschema.SingleNestedAttribute{
							Description:         "The `certificate_attributes` parameter.",
							MarkdownDescription: "The `certificate_attributes` parameter.",
							Computed:            true,
							Attributes: map[string]dsschema.Attribute{
								"dnsname": dsschema.StringAttribute{
									Description:         "The `dnsname` parameter.",
									MarkdownDescription: "The `dnsname` parameter.",
									Computed:            true,
								},
								"rfc822name": dsschema.StringAttribute{
									Description:         "The `rfc822name` parameter.",
									MarkdownDescription: "The `rfc822name` parameter.",
									Computed:            true,
								},
								"uniform_resource_identifier": dsschema.StringAttribute{
									Description:         "The `uniform_resource_identifier` parameter.",
									MarkdownDescription: "The `uniform_resource_identifier` parameter.",
									Computed:            true,
								},
							},
						},
						"digest": dsschema.StringAttribute{
							Description:         "The `digest` parameter.",
							MarkdownDescription: "The `digest` parameter.",
							Computed:            true,
						},
						"fingerprint": dsschema.StringAttribute{
							Description:         "The `fingerprint` parameter.",
							MarkdownDescription: "The `fingerprint` parameter.",
							Computed:            true,
						},
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
						"scep_ca_cert": dsschema.StringAttribute{
							Description:         "The `scep_ca_cert` parameter.",
							MarkdownDescription: "The `scep_ca_cert` parameter.",
							Computed:            true,
						},
						"scep_challenge": dsschema.SingleNestedAttribute{
							Description:         "The `scep_challenge` parameter.",
							MarkdownDescription: "The `scep_challenge` parameter.",
							Computed:            true,
							Attributes: map[string]dsschema.Attribute{
								"dynamic_value": dsschema.SingleNestedAttribute{
									Description:         "The `dynamic_value` parameter.",
									MarkdownDescription: "The `dynamic_value` parameter.",
									Computed:            true,
									Attributes: map[string]dsschema.Attribute{
										"otp_server_url": dsschema.StringAttribute{
											Description:         "The `otp_server_url` parameter.",
											MarkdownDescription: "The `otp_server_url` parameter.",
											Computed:            true,
										},
										"password": dsschema.StringAttribute{
											Description:         "The `password` parameter.",
											MarkdownDescription: "The `password` parameter.",
											Computed:            true,
										},
										"username": dsschema.StringAttribute{
											Description:         "The `username` parameter.",
											MarkdownDescription: "The `username` parameter.",
											Computed:            true,
										},
									},
								},
								"fixed": dsschema.StringAttribute{
									Description:         "The `fixed` parameter.",
									MarkdownDescription: "The `fixed` parameter.",
									Computed:            true,
								},
								"none": dsschema.StringAttribute{
									Description:         "The `none` parameter.",
									MarkdownDescription: "The `none` parameter.",
									Computed:            true,
								},
							},
						},
						"scep_client_cert": dsschema.StringAttribute{
							Description:         "The `scep_client_cert` parameter.",
							MarkdownDescription: "The `scep_client_cert` parameter.",
							Computed:            true,
						},
						"scep_url": dsschema.StringAttribute{
							Description:         "The `scep_url` parameter.",
							MarkdownDescription: "The `scep_url` parameter.",
							Computed:            true,
						},
						"subject": dsschema.StringAttribute{
							Description:         "The `subject` parameter.",
							MarkdownDescription: "The `subject` parameter.",
							Computed:            true,
						},
						"use_as_digital_signature": dsschema.BoolAttribute{
							Description:         "The `use_as_digital_signature` parameter.",
							MarkdownDescription: "The `use_as_digital_signature` parameter.",
							Computed:            true,
						},
						"use_for_key_encipherment": dsschema.BoolAttribute{
							Description:         "The `use_for_key_encipherment` parameter.",
							MarkdownDescription: "The `use_for_key_encipherment` parameter.",
							Computed:            true,
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
func (d *scepProfilesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *scepProfilesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state scepProfilesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_scep_profiles_list",
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
	svc := xlSkOUa.NewClient(d.client)
	input := xlSkOUa.ListInput{
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
	var var0 []scepProfilesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]scepProfilesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 scepProfilesListDsModelConfig
			var var3 *scepProfilesListDsModelAlgorithmObject
			if var1.Algorithm != nil {
				var3 = &scepProfilesListDsModelAlgorithmObject{}
				var var4 *scepProfilesListDsModelRsaObject
				if var1.Algorithm.Rsa != nil {
					var4 = &scepProfilesListDsModelRsaObject{}
					var4.RsaNbits = types.StringValue(var1.Algorithm.Rsa.RsaNbits)
				}
				var3.Rsa = var4
			}
			var var5 *scepProfilesListDsModelCertificateAttributesObject
			if var1.CertificateAttributes != nil {
				var5 = &scepProfilesListDsModelCertificateAttributesObject{}
				var5.Dnsname = types.StringValue(var1.CertificateAttributes.Dnsname)
				var5.Rfc822name = types.StringValue(var1.CertificateAttributes.Rfc822name)
				var5.UniformResourceIdentifier = types.StringValue(var1.CertificateAttributes.UniformResourceIdentifier)
			}
			var var6 *scepProfilesListDsModelScepChallengeObject
			if var1.ScepChallenge != nil {
				var6 = &scepProfilesListDsModelScepChallengeObject{}
				var var7 *scepProfilesListDsModelDynamicObject
				if var1.ScepChallenge.DynamicValue != nil {
					var7 = &scepProfilesListDsModelDynamicObject{}
					var7.OtpServerUrl = types.StringValue(var1.ScepChallenge.DynamicValue.OtpServerUrl)
					var7.Password = types.StringValue(var1.ScepChallenge.DynamicValue.Password)
					var7.Username = types.StringValue(var1.ScepChallenge.DynamicValue.Username)
				}
				var6.DynamicValue = var7
				var6.Fixed = types.StringValue(var1.ScepChallenge.Fixed)
				var6.None = types.StringValue(var1.ScepChallenge.None)
			}
			var2.Algorithm = var3
			var2.CaIdentityName = types.StringValue(var1.CaIdentityName)
			var2.CertificateAttributes = var5
			var2.Digest = types.StringValue(var1.Digest)
			var2.Fingerprint = types.StringValue(var1.Fingerprint)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Name = types.StringValue(var1.Name)
			var2.ScepCaCert = types.StringValue(var1.ScepCaCert)
			var2.ScepChallenge = var6
			var2.ScepClientCert = types.StringValue(var1.ScepClientCert)
			var2.ScepUrl = types.StringValue(var1.ScepUrl)
			var2.Subject = types.StringValue(var1.Subject)
			var2.UseAsDigitalSignature = types.BoolValue(var1.UseAsDigitalSignature)
			var2.UseForKeyEncipherment = types.BoolValue(var1.UseForKeyEncipherment)
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
	_ datasource.DataSource              = &scepProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &scepProfilesDataSource{}
)

func NewScepProfilesDataSource() datasource.DataSource {
	return &scepProfilesDataSource{}
}

type scepProfilesDataSource struct {
	client *sase.Client
}

type scepProfilesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`

	// Output.
	// Ref: #/components/schemas/scep-profiles
	Algorithm             *scepProfilesDsModelAlgorithmObject             `tfsdk:"algorithm"`
	CaIdentityName        types.String                                    `tfsdk:"ca_identity_name"`
	CertificateAttributes *scepProfilesDsModelCertificateAttributesObject `tfsdk:"certificate_attributes"`
	Digest                types.String                                    `tfsdk:"digest"`
	Fingerprint           types.String                                    `tfsdk:"fingerprint"`
	// input omit: ObjectId
	Name                  types.String                            `tfsdk:"name"`
	ScepCaCert            types.String                            `tfsdk:"scep_ca_cert"`
	ScepChallenge         *scepProfilesDsModelScepChallengeObject `tfsdk:"scep_challenge"`
	ScepClientCert        types.String                            `tfsdk:"scep_client_cert"`
	ScepUrl               types.String                            `tfsdk:"scep_url"`
	Subject               types.String                            `tfsdk:"subject"`
	UseAsDigitalSignature types.Bool                              `tfsdk:"use_as_digital_signature"`
	UseForKeyEncipherment types.Bool                              `tfsdk:"use_for_key_encipherment"`
}

type scepProfilesDsModelAlgorithmObject struct {
	Rsa *scepProfilesDsModelRsaObject `tfsdk:"rsa"`
}

type scepProfilesDsModelRsaObject struct {
	RsaNbits types.String `tfsdk:"rsa_nbits"`
}

type scepProfilesDsModelCertificateAttributesObject struct {
	Dnsname                   types.String `tfsdk:"dnsname"`
	Rfc822name                types.String `tfsdk:"rfc822name"`
	UniformResourceIdentifier types.String `tfsdk:"uniform_resource_identifier"`
}

type scepProfilesDsModelScepChallengeObject struct {
	DynamicValue *scepProfilesDsModelDynamicObject `tfsdk:"dynamic_value"`
	Fixed        types.String                      `tfsdk:"fixed"`
	None         types.String                      `tfsdk:"none"`
}

type scepProfilesDsModelDynamicObject struct {
	OtpServerUrl types.String `tfsdk:"otp_server_url"`
	Password     types.String `tfsdk:"password"`
	Username     types.String `tfsdk:"username"`
}

// Metadata returns the data source type name.
func (d *scepProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scep_profiles"
}

// Schema defines the schema for this listing data source.
func (d *scepProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
				Description:         "The uuid of the resource",
				MarkdownDescription: "The uuid of the resource",
				Required:            true,
			},

			// Output.
			"algorithm": dsschema.SingleNestedAttribute{
				Description:         "The `algorithm` parameter.",
				MarkdownDescription: "The `algorithm` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"rsa": dsschema.SingleNestedAttribute{
						Description:         "The `rsa` parameter.",
						MarkdownDescription: "The `rsa` parameter.",
						Computed:            true,
						Attributes: map[string]dsschema.Attribute{
							"rsa_nbits": dsschema.StringAttribute{
								Description:         "The `rsa_nbits` parameter.",
								MarkdownDescription: "The `rsa_nbits` parameter.",
								Computed:            true,
							},
						},
					},
				},
			},
			"ca_identity_name": dsschema.StringAttribute{
				Description:         "The `ca_identity_name` parameter.",
				MarkdownDescription: "The `ca_identity_name` parameter.",
				Computed:            true,
			},
			"certificate_attributes": dsschema.SingleNestedAttribute{
				Description:         "The `certificate_attributes` parameter.",
				MarkdownDescription: "The `certificate_attributes` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"dnsname": dsschema.StringAttribute{
						Description:         "The `dnsname` parameter.",
						MarkdownDescription: "The `dnsname` parameter.",
						Computed:            true,
					},
					"rfc822name": dsschema.StringAttribute{
						Description:         "The `rfc822name` parameter.",
						MarkdownDescription: "The `rfc822name` parameter.",
						Computed:            true,
					},
					"uniform_resource_identifier": dsschema.StringAttribute{
						Description:         "The `uniform_resource_identifier` parameter.",
						MarkdownDescription: "The `uniform_resource_identifier` parameter.",
						Computed:            true,
					},
				},
			},
			"digest": dsschema.StringAttribute{
				Description:         "The `digest` parameter.",
				MarkdownDescription: "The `digest` parameter.",
				Computed:            true,
			},
			"fingerprint": dsschema.StringAttribute{
				Description:         "The `fingerprint` parameter.",
				MarkdownDescription: "The `fingerprint` parameter.",
				Computed:            true,
			},
			"name": dsschema.StringAttribute{
				Description:         "The `name` parameter.",
				MarkdownDescription: "The `name` parameter.",
				Computed:            true,
			},
			"scep_ca_cert": dsschema.StringAttribute{
				Description:         "The `scep_ca_cert` parameter.",
				MarkdownDescription: "The `scep_ca_cert` parameter.",
				Computed:            true,
			},
			"scep_challenge": dsschema.SingleNestedAttribute{
				Description:         "The `scep_challenge` parameter.",
				MarkdownDescription: "The `scep_challenge` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"dynamic_value": dsschema.SingleNestedAttribute{
						Description:         "The `dynamic_value` parameter.",
						MarkdownDescription: "The `dynamic_value` parameter.",
						Computed:            true,
						Attributes: map[string]dsschema.Attribute{
							"otp_server_url": dsschema.StringAttribute{
								Description:         "The `otp_server_url` parameter.",
								MarkdownDescription: "The `otp_server_url` parameter.",
								Computed:            true,
							},
							"password": dsschema.StringAttribute{
								Description:         "The `password` parameter.",
								MarkdownDescription: "The `password` parameter.",
								Computed:            true,
							},
							"username": dsschema.StringAttribute{
								Description:         "The `username` parameter.",
								MarkdownDescription: "The `username` parameter.",
								Computed:            true,
							},
						},
					},
					"fixed": dsschema.StringAttribute{
						Description:         "The `fixed` parameter.",
						MarkdownDescription: "The `fixed` parameter.",
						Computed:            true,
					},
					"none": dsschema.StringAttribute{
						Description:         "The `none` parameter.",
						MarkdownDescription: "The `none` parameter.",
						Computed:            true,
					},
				},
			},
			"scep_client_cert": dsschema.StringAttribute{
				Description:         "The `scep_client_cert` parameter.",
				MarkdownDescription: "The `scep_client_cert` parameter.",
				Computed:            true,
			},
			"scep_url": dsschema.StringAttribute{
				Description:         "The `scep_url` parameter.",
				MarkdownDescription: "The `scep_url` parameter.",
				Computed:            true,
			},
			"subject": dsschema.StringAttribute{
				Description:         "The `subject` parameter.",
				MarkdownDescription: "The `subject` parameter.",
				Computed:            true,
			},
			"use_as_digital_signature": dsschema.BoolAttribute{
				Description:         "The `use_as_digital_signature` parameter.",
				MarkdownDescription: "The `use_as_digital_signature` parameter.",
				Computed:            true,
			},
			"use_for_key_encipherment": dsschema.BoolAttribute{
				Description:         "The `use_for_key_encipherment` parameter.",
				MarkdownDescription: "The `use_for_key_encipherment` parameter.",
				Computed:            true,
			},
		},
	}
}

// Configure prepares the struct.
func (d *scepProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *scepProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state scepProfilesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_scep_profiles",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to run the command.
	svc := xlSkOUa.NewClient(d.client)
	input := xlSkOUa.ReadInput{
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
	var var0 *scepProfilesDsModelAlgorithmObject
	if ans.Algorithm != nil {
		var0 = &scepProfilesDsModelAlgorithmObject{}
		var var1 *scepProfilesDsModelRsaObject
		if ans.Algorithm.Rsa != nil {
			var1 = &scepProfilesDsModelRsaObject{}
			var1.RsaNbits = types.StringValue(ans.Algorithm.Rsa.RsaNbits)
		}
		var0.Rsa = var1
	}
	var var2 *scepProfilesDsModelCertificateAttributesObject
	if ans.CertificateAttributes != nil {
		var2 = &scepProfilesDsModelCertificateAttributesObject{}
		var2.Dnsname = types.StringValue(ans.CertificateAttributes.Dnsname)
		var2.Rfc822name = types.StringValue(ans.CertificateAttributes.Rfc822name)
		var2.UniformResourceIdentifier = types.StringValue(ans.CertificateAttributes.UniformResourceIdentifier)
	}
	var var3 *scepProfilesDsModelScepChallengeObject
	if ans.ScepChallenge != nil {
		var3 = &scepProfilesDsModelScepChallengeObject{}
		var var4 *scepProfilesDsModelDynamicObject
		if ans.ScepChallenge.DynamicValue != nil {
			var4 = &scepProfilesDsModelDynamicObject{}
			var4.OtpServerUrl = types.StringValue(ans.ScepChallenge.DynamicValue.OtpServerUrl)
			var4.Password = types.StringValue(ans.ScepChallenge.DynamicValue.Password)
			var4.Username = types.StringValue(ans.ScepChallenge.DynamicValue.Username)
		}
		var3.DynamicValue = var4
		var3.Fixed = types.StringValue(ans.ScepChallenge.Fixed)
		var3.None = types.StringValue(ans.ScepChallenge.None)
	}
	state.Algorithm = var0
	state.CaIdentityName = types.StringValue(ans.CaIdentityName)
	state.CertificateAttributes = var2
	state.Digest = types.StringValue(ans.Digest)
	state.Fingerprint = types.StringValue(ans.Fingerprint)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.ScepCaCert = types.StringValue(ans.ScepCaCert)
	state.ScepChallenge = var3
	state.ScepClientCert = types.StringValue(ans.ScepClientCert)
	state.ScepUrl = types.StringValue(ans.ScepUrl)
	state.Subject = types.StringValue(ans.Subject)
	state.UseAsDigitalSignature = types.BoolValue(ans.UseAsDigitalSignature)
	state.UseForKeyEncipherment = types.BoolValue(ans.UseForKeyEncipherment)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &scepProfilesResource{}
	_ resource.ResourceWithConfigure   = &scepProfilesResource{}
	_ resource.ResourceWithImportState = &scepProfilesResource{}
)

func NewScepProfilesResource() resource.Resource {
	return &scepProfilesResource{}
}

type scepProfilesResource struct {
	client *sase.Client
}

type scepProfilesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Type types.String `tfsdk:"type"`

	// Request body input.
	// Ref: #/components/schemas/scep-profiles
	Algorithm             *scepProfilesRsModelAlgorithmObject             `tfsdk:"algorithm"`
	CaIdentityName        types.String                                    `tfsdk:"ca_identity_name"`
	CertificateAttributes *scepProfilesRsModelCertificateAttributesObject `tfsdk:"certificate_attributes"`
	Digest                types.String                                    `tfsdk:"digest"`
	Fingerprint           types.String                                    `tfsdk:"fingerprint"`
	ObjectId              types.String                                    `tfsdk:"object_id"`
	Name                  types.String                                    `tfsdk:"name"`
	ScepCaCert            types.String                                    `tfsdk:"scep_ca_cert"`
	ScepChallenge         *scepProfilesRsModelScepChallengeObject         `tfsdk:"scep_challenge"`
	ScepClientCert        types.String                                    `tfsdk:"scep_client_cert"`
	ScepUrl               types.String                                    `tfsdk:"scep_url"`
	Subject               types.String                                    `tfsdk:"subject"`
	UseAsDigitalSignature types.Bool                                      `tfsdk:"use_as_digital_signature"`
	UseForKeyEncipherment types.Bool                                      `tfsdk:"use_for_key_encipherment"`
}

type scepProfilesRsModelAlgorithmObject struct {
	Rsa *scepProfilesRsModelRsaObject `tfsdk:"rsa"`
}

type scepProfilesRsModelRsaObject struct {
	RsaNbits types.String `tfsdk:"rsa_nbits"`
}

type scepProfilesRsModelCertificateAttributesObject struct {
	Dnsname                   types.String `tfsdk:"dnsname"`
	Rfc822name                types.String `tfsdk:"rfc822name"`
	UniformResourceIdentifier types.String `tfsdk:"uniform_resource_identifier"`
}

type scepProfilesRsModelScepChallengeObject struct {
	DynamicValue *scepProfilesRsModelDynamicObject `tfsdk:"dynamic_value"`
	Fixed        types.String                      `tfsdk:"fixed"`
	None         types.String                      `tfsdk:"none"`
}

type scepProfilesRsModelDynamicObject struct {
	OtpServerUrl types.String `tfsdk:"otp_server_url"`
	Password     types.String `tfsdk:"password"`
	Username     types.String `tfsdk:"username"`
}

// Metadata returns the data source type name.
func (r *scepProfilesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scep_profiles"
}

// Schema defines the schema for this listing data source.
func (r *scepProfilesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"type": rsschema.StringAttribute{
				Description:         "The type of the schema node",
				MarkdownDescription: "The type of the schema node",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("cloud", "container"),
				},
			},

			"algorithm": rsschema.SingleNestedAttribute{
				Description:         "The `algorithm` parameter.",
				MarkdownDescription: "The `algorithm` parameter.",
				Optional:            true,
				Attributes: map[string]rsschema.Attribute{
					"rsa": rsschema.SingleNestedAttribute{
						Description:         "The `rsa` parameter.",
						MarkdownDescription: "The `rsa` parameter.",
						Optional:            true,
						Attributes: map[string]rsschema.Attribute{
							"rsa_nbits": rsschema.StringAttribute{
								Description:         "The `rsa_nbits` parameter.",
								MarkdownDescription: "The `rsa_nbits` parameter.",
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
			"ca_identity_name": rsschema.StringAttribute{
				Description:         "The `ca_identity_name` parameter.",
				MarkdownDescription: "The `ca_identity_name` parameter.",
				Required:            true,
			},
			"certificate_attributes": rsschema.SingleNestedAttribute{
				Description:         "The `certificate_attributes` parameter.",
				MarkdownDescription: "The `certificate_attributes` parameter.",
				Optional:            true,
				Attributes: map[string]rsschema.Attribute{
					"dnsname": rsschema.StringAttribute{
						Description:         "The `dnsname` parameter.",
						MarkdownDescription: "The `dnsname` parameter.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							DefaultString(""),
						},
					},
					"rfc822name": rsschema.StringAttribute{
						Description:         "The `rfc822name` parameter.",
						MarkdownDescription: "The `rfc822name` parameter.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							DefaultString(""),
						},
					},
					"uniform_resource_identifier": rsschema.StringAttribute{
						Description:         "The `uniform_resource_identifier` parameter.",
						MarkdownDescription: "The `uniform_resource_identifier` parameter.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							DefaultString(""),
						},
					},
				},
			},
			"digest": rsschema.StringAttribute{
				Description:         "The `digest` parameter.",
				MarkdownDescription: "The `digest` parameter.",
				Required:            true,
			},
			"fingerprint": rsschema.StringAttribute{
				Description:         "The `fingerprint` parameter.",
				MarkdownDescription: "The `fingerprint` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
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
				Validators: []validator.String{
					stringvalidator.LengthAtMost(31),
				},
			},
			"scep_ca_cert": rsschema.StringAttribute{
				Description:         "The `scep_ca_cert` parameter.",
				MarkdownDescription: "The `scep_ca_cert` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
			},
			"scep_challenge": rsschema.SingleNestedAttribute{
				Description:         "The `scep_challenge` parameter.",
				MarkdownDescription: "The `scep_challenge` parameter.",
				Optional:            true,
				Attributes: map[string]rsschema.Attribute{
					"dynamic_value": rsschema.SingleNestedAttribute{
						Description:         "The `dynamic_value` parameter.",
						MarkdownDescription: "The `dynamic_value` parameter.",
						Optional:            true,
						Attributes: map[string]rsschema.Attribute{
							"otp_server_url": rsschema.StringAttribute{
								Description:         "The `otp_server_url` parameter.",
								MarkdownDescription: "The `otp_server_url` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 255),
								},
							},
							"password": rsschema.StringAttribute{
								Description:         "The `password` parameter.",
								MarkdownDescription: "The `password` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 255),
								},
							},
							"username": rsschema.StringAttribute{
								Description:         "The `username` parameter.",
								MarkdownDescription: "The `username` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 255),
								},
							},
						},
					},
					"fixed": rsschema.StringAttribute{
						Description:         "The `fixed` parameter.",
						MarkdownDescription: "The `fixed` parameter.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							DefaultString(""),
						},
						Validators: []validator.String{
							stringvalidator.LengthBetween(0, 1024),
						},
					},
					"none": rsschema.StringAttribute{
						Description:         "The `none` parameter.",
						MarkdownDescription: "The `none` parameter.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							DefaultString(""),
						},
						Validators: []validator.String{
							stringvalidator.OneOf(""),
						},
					},
				},
			},
			"scep_client_cert": rsschema.StringAttribute{
				Description:         "The `scep_client_cert` parameter.",
				MarkdownDescription: "The `scep_client_cert` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
			},
			"scep_url": rsschema.StringAttribute{
				Description:         "The `scep_url` parameter.",
				MarkdownDescription: "The `scep_url` parameter.",
				Required:            true,
			},
			"subject": rsschema.StringAttribute{
				Description:         "The `subject` parameter.",
				MarkdownDescription: "The `subject` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
			},
			"use_as_digital_signature": rsschema.BoolAttribute{
				Description:         "The `use_as_digital_signature` parameter.",
				MarkdownDescription: "The `use_as_digital_signature` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"use_for_key_encipherment": rsschema.BoolAttribute{
				Description:         "The `use_for_key_encipherment` parameter.",
				MarkdownDescription: "The `use_for_key_encipherment` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *scepProfilesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *scepProfilesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state scepProfilesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_scep_profiles",
		"type":                        state.Type.ValueString(),
	})

	// Prepare to create the config.
	svc := xlSkOUa.NewClient(r.client)
	input := xlSkOUa.CreateInput{
		Type: state.Type.ValueString(),
	}
	var var0 qmgDayz.Config
	var var1 *qmgDayz.AlgorithmObject
	if state.Algorithm != nil {
		var1 = &qmgDayz.AlgorithmObject{}
		var var2 *qmgDayz.RsaObject
		if state.Algorithm.Rsa != nil {
			var2 = &qmgDayz.RsaObject{}
			var2.RsaNbits = state.Algorithm.Rsa.RsaNbits.ValueString()
		}
		var1.Rsa = var2
	}
	var0.Algorithm = var1
	var0.CaIdentityName = state.CaIdentityName.ValueString()
	var var3 *qmgDayz.CertificateAttributesObject
	if state.CertificateAttributes != nil {
		var3 = &qmgDayz.CertificateAttributesObject{}
		var3.Dnsname = state.CertificateAttributes.Dnsname.ValueString()
		var3.Rfc822name = state.CertificateAttributes.Rfc822name.ValueString()
		var3.UniformResourceIdentifier = state.CertificateAttributes.UniformResourceIdentifier.ValueString()
	}
	var0.CertificateAttributes = var3
	var0.Digest = state.Digest.ValueString()
	var0.Fingerprint = state.Fingerprint.ValueString()
	var0.Name = state.Name.ValueString()
	var0.ScepCaCert = state.ScepCaCert.ValueString()
	var var4 *qmgDayz.ScepChallengeObject
	if state.ScepChallenge != nil {
		var4 = &qmgDayz.ScepChallengeObject{}
		var var5 *qmgDayz.DynamicObject
		if state.ScepChallenge.DynamicValue != nil {
			var5 = &qmgDayz.DynamicObject{}
			var5.OtpServerUrl = state.ScepChallenge.DynamicValue.OtpServerUrl.ValueString()
			var5.Password = state.ScepChallenge.DynamicValue.Password.ValueString()
			var5.Username = state.ScepChallenge.DynamicValue.Username.ValueString()
		}
		var4.DynamicValue = var5
		var4.Fixed = state.ScepChallenge.Fixed.ValueString()
		var4.None = state.ScepChallenge.None.ValueString()
	}
	var0.ScepChallenge = var4
	var0.ScepClientCert = state.ScepClientCert.ValueString()
	var0.ScepUrl = state.ScepUrl.ValueString()
	var0.Subject = state.Subject.ValueString()
	var0.UseAsDigitalSignature = state.UseAsDigitalSignature.ValueBool()
	var0.UseForKeyEncipherment = state.UseForKeyEncipherment.ValueBool()
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Create(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in create", err.Error())
		return
	}

	// Store the answer to state.
	var idBuilder strings.Builder
	idBuilder.WriteString(input.Type)
	idBuilder.WriteString(IdSeparator)
	idBuilder.WriteString(ans.ObjectId)
	state.Id = types.StringValue(idBuilder.String())
	var var6 *scepProfilesRsModelAlgorithmObject
	if ans.Algorithm != nil {
		var6 = &scepProfilesRsModelAlgorithmObject{}
		var var7 *scepProfilesRsModelRsaObject
		if ans.Algorithm.Rsa != nil {
			var7 = &scepProfilesRsModelRsaObject{}
			var7.RsaNbits = types.StringValue(ans.Algorithm.Rsa.RsaNbits)
		}
		var6.Rsa = var7
	}
	var var8 *scepProfilesRsModelCertificateAttributesObject
	if ans.CertificateAttributes != nil {
		var8 = &scepProfilesRsModelCertificateAttributesObject{}
		var8.Dnsname = types.StringValue(ans.CertificateAttributes.Dnsname)
		var8.Rfc822name = types.StringValue(ans.CertificateAttributes.Rfc822name)
		var8.UniformResourceIdentifier = types.StringValue(ans.CertificateAttributes.UniformResourceIdentifier)
	}
	var var9 *scepProfilesRsModelScepChallengeObject
	if ans.ScepChallenge != nil {
		var9 = &scepProfilesRsModelScepChallengeObject{}
		var var10 *scepProfilesRsModelDynamicObject
		if ans.ScepChallenge.DynamicValue != nil {
			var10 = &scepProfilesRsModelDynamicObject{}
			var10.OtpServerUrl = types.StringValue(ans.ScepChallenge.DynamicValue.OtpServerUrl)
			var10.Password = types.StringValue(ans.ScepChallenge.DynamicValue.Password)
			var10.Username = types.StringValue(ans.ScepChallenge.DynamicValue.Username)
		}
		var9.DynamicValue = var10
		var9.Fixed = types.StringValue(ans.ScepChallenge.Fixed)
		var9.None = types.StringValue(ans.ScepChallenge.None)
	}
	state.Algorithm = var6
	state.CaIdentityName = types.StringValue(ans.CaIdentityName)
	state.CertificateAttributes = var8
	state.Digest = types.StringValue(ans.Digest)
	state.Fingerprint = types.StringValue(ans.Fingerprint)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.ScepCaCert = types.StringValue(ans.ScepCaCert)
	state.ScepChallenge = var9
	state.ScepClientCert = types.StringValue(ans.ScepClientCert)
	state.ScepUrl = types.StringValue(ans.ScepUrl)
	state.Subject = types.StringValue(ans.Subject)
	state.UseAsDigitalSignature = types.BoolValue(ans.UseAsDigitalSignature)
	state.UseForKeyEncipherment = types.BoolValue(ans.UseForKeyEncipherment)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *scepProfilesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state scepProfilesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_scep_profiles",
		"locMap":                      map[string]int{"ObjectId": 1, "Type": 0},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := xlSkOUa.NewClient(r.client)
	input := xlSkOUa.ReadInput{
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
	state.Type = types.StringValue(tokens[0])
	state.Id = idType
	var var0 *scepProfilesRsModelAlgorithmObject
	if ans.Algorithm != nil {
		var0 = &scepProfilesRsModelAlgorithmObject{}
		var var1 *scepProfilesRsModelRsaObject
		if ans.Algorithm.Rsa != nil {
			var1 = &scepProfilesRsModelRsaObject{}
			var1.RsaNbits = types.StringValue(ans.Algorithm.Rsa.RsaNbits)
		}
		var0.Rsa = var1
	}
	var var2 *scepProfilesRsModelCertificateAttributesObject
	if ans.CertificateAttributes != nil {
		var2 = &scepProfilesRsModelCertificateAttributesObject{}
		var2.Dnsname = types.StringValue(ans.CertificateAttributes.Dnsname)
		var2.Rfc822name = types.StringValue(ans.CertificateAttributes.Rfc822name)
		var2.UniformResourceIdentifier = types.StringValue(ans.CertificateAttributes.UniformResourceIdentifier)
	}
	var var3 *scepProfilesRsModelScepChallengeObject
	if ans.ScepChallenge != nil {
		var3 = &scepProfilesRsModelScepChallengeObject{}
		var var4 *scepProfilesRsModelDynamicObject
		if ans.ScepChallenge.DynamicValue != nil {
			var4 = &scepProfilesRsModelDynamicObject{}
			var4.OtpServerUrl = types.StringValue(ans.ScepChallenge.DynamicValue.OtpServerUrl)
			var4.Password = types.StringValue(ans.ScepChallenge.DynamicValue.Password)
			var4.Username = types.StringValue(ans.ScepChallenge.DynamicValue.Username)
		}
		var3.DynamicValue = var4
		var3.Fixed = types.StringValue(ans.ScepChallenge.Fixed)
		var3.None = types.StringValue(ans.ScepChallenge.None)
	}
	state.Algorithm = var0
	state.CaIdentityName = types.StringValue(ans.CaIdentityName)
	state.CertificateAttributes = var2
	state.Digest = types.StringValue(ans.Digest)
	state.Fingerprint = types.StringValue(ans.Fingerprint)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.ScepCaCert = types.StringValue(ans.ScepCaCert)
	state.ScepChallenge = var3
	state.ScepClientCert = types.StringValue(ans.ScepClientCert)
	state.ScepUrl = types.StringValue(ans.ScepUrl)
	state.Subject = types.StringValue(ans.Subject)
	state.UseAsDigitalSignature = types.BoolValue(ans.UseAsDigitalSignature)
	state.UseForKeyEncipherment = types.BoolValue(ans.UseForKeyEncipherment)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *scepProfilesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state scepProfilesRsModel
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
		"resource_name":               "sase_scep_profiles",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := xlSkOUa.NewClient(r.client)
	input := xlSkOUa.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 qmgDayz.Config
	var var1 *qmgDayz.AlgorithmObject
	if plan.Algorithm != nil {
		var1 = &qmgDayz.AlgorithmObject{}
		var var2 *qmgDayz.RsaObject
		if plan.Algorithm.Rsa != nil {
			var2 = &qmgDayz.RsaObject{}
			var2.RsaNbits = plan.Algorithm.Rsa.RsaNbits.ValueString()
		}
		var1.Rsa = var2
	}
	var0.Algorithm = var1
	var0.CaIdentityName = plan.CaIdentityName.ValueString()
	var var3 *qmgDayz.CertificateAttributesObject
	if plan.CertificateAttributes != nil {
		var3 = &qmgDayz.CertificateAttributesObject{}
		var3.Dnsname = plan.CertificateAttributes.Dnsname.ValueString()
		var3.Rfc822name = plan.CertificateAttributes.Rfc822name.ValueString()
		var3.UniformResourceIdentifier = plan.CertificateAttributes.UniformResourceIdentifier.ValueString()
	}
	var0.CertificateAttributes = var3
	var0.Digest = plan.Digest.ValueString()
	var0.Fingerprint = plan.Fingerprint.ValueString()
	var0.Name = plan.Name.ValueString()
	var0.ScepCaCert = plan.ScepCaCert.ValueString()
	var var4 *qmgDayz.ScepChallengeObject
	if plan.ScepChallenge != nil {
		var4 = &qmgDayz.ScepChallengeObject{}
		var var5 *qmgDayz.DynamicObject
		if plan.ScepChallenge.DynamicValue != nil {
			var5 = &qmgDayz.DynamicObject{}
			var5.OtpServerUrl = plan.ScepChallenge.DynamicValue.OtpServerUrl.ValueString()
			var5.Password = plan.ScepChallenge.DynamicValue.Password.ValueString()
			var5.Username = plan.ScepChallenge.DynamicValue.Username.ValueString()
		}
		var4.DynamicValue = var5
		var4.Fixed = plan.ScepChallenge.Fixed.ValueString()
		var4.None = plan.ScepChallenge.None.ValueString()
	}
	var0.ScepChallenge = var4
	var0.ScepClientCert = plan.ScepClientCert.ValueString()
	var0.ScepUrl = plan.ScepUrl.ValueString()
	var0.Subject = plan.Subject.ValueString()
	var0.UseAsDigitalSignature = plan.UseAsDigitalSignature.ValueBool()
	var0.UseForKeyEncipherment = plan.UseForKeyEncipherment.ValueBool()
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var6 *scepProfilesRsModelAlgorithmObject
	if ans.Algorithm != nil {
		var6 = &scepProfilesRsModelAlgorithmObject{}
		var var7 *scepProfilesRsModelRsaObject
		if ans.Algorithm.Rsa != nil {
			var7 = &scepProfilesRsModelRsaObject{}
			var7.RsaNbits = types.StringValue(ans.Algorithm.Rsa.RsaNbits)
		}
		var6.Rsa = var7
	}
	var var8 *scepProfilesRsModelCertificateAttributesObject
	if ans.CertificateAttributes != nil {
		var8 = &scepProfilesRsModelCertificateAttributesObject{}
		var8.Dnsname = types.StringValue(ans.CertificateAttributes.Dnsname)
		var8.Rfc822name = types.StringValue(ans.CertificateAttributes.Rfc822name)
		var8.UniformResourceIdentifier = types.StringValue(ans.CertificateAttributes.UniformResourceIdentifier)
	}
	var var9 *scepProfilesRsModelScepChallengeObject
	if ans.ScepChallenge != nil {
		var9 = &scepProfilesRsModelScepChallengeObject{}
		var var10 *scepProfilesRsModelDynamicObject
		if ans.ScepChallenge.DynamicValue != nil {
			var10 = &scepProfilesRsModelDynamicObject{}
			var10.OtpServerUrl = types.StringValue(ans.ScepChallenge.DynamicValue.OtpServerUrl)
			var10.Password = types.StringValue(ans.ScepChallenge.DynamicValue.Password)
			var10.Username = types.StringValue(ans.ScepChallenge.DynamicValue.Username)
		}
		var9.DynamicValue = var10
		var9.Fixed = types.StringValue(ans.ScepChallenge.Fixed)
		var9.None = types.StringValue(ans.ScepChallenge.None)
	}
	state.Algorithm = var6
	state.CaIdentityName = types.StringValue(ans.CaIdentityName)
	state.CertificateAttributes = var8
	state.Digest = types.StringValue(ans.Digest)
	state.Fingerprint = types.StringValue(ans.Fingerprint)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.ScepCaCert = types.StringValue(ans.ScepCaCert)
	state.ScepChallenge = var9
	state.ScepClientCert = types.StringValue(ans.ScepClientCert)
	state.ScepUrl = types.StringValue(ans.ScepUrl)
	state.Subject = types.StringValue(ans.Subject)
	state.UseAsDigitalSignature = types.BoolValue(ans.UseAsDigitalSignature)
	state.UseForKeyEncipherment = types.BoolValue(ans.UseForKeyEncipherment)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *scepProfilesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name":               "sase_scep_profiles",
		"locMap":                      map[string]int{"ObjectId": 1, "Type": 0},
		"tokens":                      tokens,
	})

	svc := xlSkOUa.NewClient(r.client)
	input := xlSkOUa.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *scepProfilesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
