package provider

import (
	"context"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	deRyMEf "github.com/paloaltonetworks/sase-go/netsec/schema/mfa/servers"
	wArkOsV "github.com/paloaltonetworks/sase-go/netsec/service/v1/mfaservers"

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

// Data source.
var (
	_ datasource.DataSource              = &mfaServersDataSource{}
	_ datasource.DataSourceWithConfigure = &mfaServersDataSource{}
)

func NewMfaServersDataSource() datasource.DataSource {
	return &mfaServersDataSource{}
}

type mfaServersDataSource struct {
	client *sase.Client
}

type mfaServersDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`
	Folder   types.String `tfsdk:"folder"`

	// Output.
	// Ref: #/components/schemas/mfa-servers
	// input omit: ObjectId
	MfaCertProfile types.String                          `tfsdk:"mfa_cert_profile"`
	MfaVendorType  *mfaServersDsModelMfaVendorTypeObject `tfsdk:"mfa_vendor_type"`
	Name           types.String                          `tfsdk:"name"`
}

type mfaServersDsModelMfaVendorTypeObject struct {
	DuoSecurityV2      *mfaServersDsModelDuoSecurityV2Object      `tfsdk:"duo_security_v2"`
	OktaAdaptiveV1     *mfaServersDsModelOktaAdaptiveV1Object     `tfsdk:"okta_adaptive_v1"`
	PingIdentityV1     *mfaServersDsModelPingIdentityV1Object     `tfsdk:"ping_identity_v1"`
	RsaSecuridAccessV1 *mfaServersDsModelRsaSecuridAccessV1Object `tfsdk:"rsa_securid_access_v1"`
}

type mfaServersDsModelDuoSecurityV2Object struct {
	DuoApiHost        types.String `tfsdk:"duo_api_host"`
	DuoBaseuri        types.String `tfsdk:"duo_baseuri"`
	DuoIntegrationKey types.String `tfsdk:"duo_integration_key"`
	DuoSecretKey      types.String `tfsdk:"duo_secret_key"`
	DuoTimeout        types.String `tfsdk:"duo_timeout"`
}

type mfaServersDsModelOktaAdaptiveV1Object struct {
	OktaApiHost types.String `tfsdk:"okta_api_host"`
	OktaBaseuri types.String `tfsdk:"okta_baseuri"`
	OktaOrg     types.String `tfsdk:"okta_org"`
	OktaTimeout types.String `tfsdk:"okta_timeout"`
	OktaToken   types.String `tfsdk:"okta_token"`
}

type mfaServersDsModelPingIdentityV1Object struct {
	PingApiHost  types.String `tfsdk:"ping_api_host"`
	PingBaseuri  types.String `tfsdk:"ping_baseuri"`
	PingOrg      types.String `tfsdk:"ping_org"`
	PingOrgAlias types.String `tfsdk:"ping_org_alias"`
	PingTimeout  types.String `tfsdk:"ping_timeout"`
	PingToken    types.String `tfsdk:"ping_token"`
}

type mfaServersDsModelRsaSecuridAccessV1Object struct {
	RsaAccessid          types.String `tfsdk:"rsa_accessid"`
	RsaAccesskey         types.String `tfsdk:"rsa_accesskey"`
	RsaApiHost           types.String `tfsdk:"rsa_api_host"`
	RsaAssurancepolicyid types.String `tfsdk:"rsa_assurancepolicyid"`
	RsaBaseuri           types.String `tfsdk:"rsa_baseuri"`
	RsaTimeout           types.String `tfsdk:"rsa_timeout"`
}

// Metadata returns the data source type name.
func (d *mfaServersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mfa_servers"
}

// Schema defines the schema for this listing data source.
func (d *mfaServersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"mfa_cert_profile": dsschema.StringAttribute{
				Description:         "The `mfa_cert_profile` parameter.",
				MarkdownDescription: "The `mfa_cert_profile` parameter.",
				Computed:            true,
			},
			"mfa_vendor_type": dsschema.SingleNestedAttribute{
				Description:         "The `mfa_vendor_type` parameter.",
				MarkdownDescription: "The `mfa_vendor_type` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"duo_security_v2": dsschema.SingleNestedAttribute{
						Description:         "The `duo_security_v2` parameter.",
						MarkdownDescription: "The `duo_security_v2` parameter.",
						Computed:            true,
						Attributes: map[string]dsschema.Attribute{
							"duo_api_host": dsschema.StringAttribute{
								Description:         "The `duo_api_host` parameter.",
								MarkdownDescription: "The `duo_api_host` parameter.",
								Computed:            true,
							},
							"duo_baseuri": dsschema.StringAttribute{
								Description:         "The `duo_baseuri` parameter.",
								MarkdownDescription: "The `duo_baseuri` parameter.",
								Computed:            true,
							},
							"duo_integration_key": dsschema.StringAttribute{
								Description:         "The `duo_integration_key` parameter.",
								MarkdownDescription: "The `duo_integration_key` parameter.",
								Computed:            true,
							},
							"duo_secret_key": dsschema.StringAttribute{
								Description:         "The `duo_secret_key` parameter.",
								MarkdownDescription: "The `duo_secret_key` parameter.",
								Computed:            true,
							},
							"duo_timeout": dsschema.StringAttribute{
								Description:         "The `duo_timeout` parameter.",
								MarkdownDescription: "The `duo_timeout` parameter.",
								Computed:            true,
							},
						},
					},
					"okta_adaptive_v1": dsschema.SingleNestedAttribute{
						Description:         "The `okta_adaptive_v1` parameter.",
						MarkdownDescription: "The `okta_adaptive_v1` parameter.",
						Computed:            true,
						Attributes: map[string]dsschema.Attribute{
							"okta_api_host": dsschema.StringAttribute{
								Description:         "The `okta_api_host` parameter.",
								MarkdownDescription: "The `okta_api_host` parameter.",
								Computed:            true,
							},
							"okta_baseuri": dsschema.StringAttribute{
								Description:         "The `okta_baseuri` parameter.",
								MarkdownDescription: "The `okta_baseuri` parameter.",
								Computed:            true,
							},
							"okta_org": dsschema.StringAttribute{
								Description:         "The `okta_org` parameter.",
								MarkdownDescription: "The `okta_org` parameter.",
								Computed:            true,
							},
							"okta_timeout": dsschema.StringAttribute{
								Description:         "The `okta_timeout` parameter.",
								MarkdownDescription: "The `okta_timeout` parameter.",
								Computed:            true,
							},
							"okta_token": dsschema.StringAttribute{
								Description:         "The `okta_token` parameter.",
								MarkdownDescription: "The `okta_token` parameter.",
								Computed:            true,
							},
						},
					},
					"ping_identity_v1": dsschema.SingleNestedAttribute{
						Description:         "The `ping_identity_v1` parameter.",
						MarkdownDescription: "The `ping_identity_v1` parameter.",
						Computed:            true,
						Attributes: map[string]dsschema.Attribute{
							"ping_api_host": dsschema.StringAttribute{
								Description:         "The `ping_api_host` parameter.",
								MarkdownDescription: "The `ping_api_host` parameter.",
								Computed:            true,
							},
							"ping_baseuri": dsschema.StringAttribute{
								Description:         "The `ping_baseuri` parameter.",
								MarkdownDescription: "The `ping_baseuri` parameter.",
								Computed:            true,
							},
							"ping_org": dsschema.StringAttribute{
								Description:         "The `ping_org` parameter.",
								MarkdownDescription: "The `ping_org` parameter.",
								Computed:            true,
							},
							"ping_org_alias": dsschema.StringAttribute{
								Description:         "The `ping_org_alias` parameter.",
								MarkdownDescription: "The `ping_org_alias` parameter.",
								Computed:            true,
							},
							"ping_timeout": dsschema.StringAttribute{
								Description:         "The `ping_timeout` parameter.",
								MarkdownDescription: "The `ping_timeout` parameter.",
								Computed:            true,
							},
							"ping_token": dsschema.StringAttribute{
								Description:         "The `ping_token` parameter.",
								MarkdownDescription: "The `ping_token` parameter.",
								Computed:            true,
							},
						},
					},
					"rsa_securid_access_v1": dsschema.SingleNestedAttribute{
						Description:         "The `rsa_securid_access_v1` parameter.",
						MarkdownDescription: "The `rsa_securid_access_v1` parameter.",
						Computed:            true,
						Attributes: map[string]dsschema.Attribute{
							"rsa_accessid": dsschema.StringAttribute{
								Description:         "The `rsa_accessid` parameter.",
								MarkdownDescription: "The `rsa_accessid` parameter.",
								Computed:            true,
							},
							"rsa_accesskey": dsschema.StringAttribute{
								Description:         "The `rsa_accesskey` parameter.",
								MarkdownDescription: "The `rsa_accesskey` parameter.",
								Computed:            true,
							},
							"rsa_api_host": dsschema.StringAttribute{
								Description:         "The `rsa_api_host` parameter.",
								MarkdownDescription: "The `rsa_api_host` parameter.",
								Computed:            true,
							},
							"rsa_assurancepolicyid": dsschema.StringAttribute{
								Description:         "The `rsa_assurancepolicyid` parameter.",
								MarkdownDescription: "The `rsa_assurancepolicyid` parameter.",
								Computed:            true,
							},
							"rsa_baseuri": dsschema.StringAttribute{
								Description:         "The `rsa_baseuri` parameter.",
								MarkdownDescription: "The `rsa_baseuri` parameter.",
								Computed:            true,
							},
							"rsa_timeout": dsschema.StringAttribute{
								Description:         "The `rsa_timeout` parameter.",
								MarkdownDescription: "The `rsa_timeout` parameter.",
								Computed:            true,
							},
						},
					},
				},
			},
			"name": dsschema.StringAttribute{
				Description:         "The `name` parameter.",
				MarkdownDescription: "The `name` parameter.",
				Computed:            true,
			},
		},
	}
}

// Configure prepares the struct.
func (d *mfaServersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *mfaServersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state mfaServersDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_mfa_servers",
		"object_id":                   state.ObjectId.ValueString(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := wArkOsV.NewClient(d.client)
	input := wArkOsV.ReadInput{
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
	var var0 *mfaServersDsModelMfaVendorTypeObject
	if ans.MfaVendorType != nil {
		var0 = &mfaServersDsModelMfaVendorTypeObject{}
		var var1 *mfaServersDsModelDuoSecurityV2Object
		if ans.MfaVendorType.DuoSecurityV2 != nil {
			var1 = &mfaServersDsModelDuoSecurityV2Object{}
			var1.DuoApiHost = types.StringValue(ans.MfaVendorType.DuoSecurityV2.DuoApiHost)
			var1.DuoBaseuri = types.StringValue(ans.MfaVendorType.DuoSecurityV2.DuoBaseuri)
			var1.DuoIntegrationKey = types.StringValue(ans.MfaVendorType.DuoSecurityV2.DuoIntegrationKey)
			var1.DuoSecretKey = types.StringValue(ans.MfaVendorType.DuoSecurityV2.DuoSecretKey)
			var1.DuoTimeout = types.StringValue(ans.MfaVendorType.DuoSecurityV2.DuoTimeout)
		}
		var var2 *mfaServersDsModelOktaAdaptiveV1Object
		if ans.MfaVendorType.OktaAdaptiveV1 != nil {
			var2 = &mfaServersDsModelOktaAdaptiveV1Object{}
			var2.OktaApiHost = types.StringValue(ans.MfaVendorType.OktaAdaptiveV1.OktaApiHost)
			var2.OktaBaseuri = types.StringValue(ans.MfaVendorType.OktaAdaptiveV1.OktaBaseuri)
			var2.OktaOrg = types.StringValue(ans.MfaVendorType.OktaAdaptiveV1.OktaOrg)
			var2.OktaTimeout = types.StringValue(ans.MfaVendorType.OktaAdaptiveV1.OktaTimeout)
			var2.OktaToken = types.StringValue(ans.MfaVendorType.OktaAdaptiveV1.OktaToken)
		}
		var var3 *mfaServersDsModelPingIdentityV1Object
		if ans.MfaVendorType.PingIdentityV1 != nil {
			var3 = &mfaServersDsModelPingIdentityV1Object{}
			var3.PingApiHost = types.StringValue(ans.MfaVendorType.PingIdentityV1.PingApiHost)
			var3.PingBaseuri = types.StringValue(ans.MfaVendorType.PingIdentityV1.PingBaseuri)
			var3.PingOrg = types.StringValue(ans.MfaVendorType.PingIdentityV1.PingOrg)
			var3.PingOrgAlias = types.StringValue(ans.MfaVendorType.PingIdentityV1.PingOrgAlias)
			var3.PingTimeout = types.StringValue(ans.MfaVendorType.PingIdentityV1.PingTimeout)
			var3.PingToken = types.StringValue(ans.MfaVendorType.PingIdentityV1.PingToken)
		}
		var var4 *mfaServersDsModelRsaSecuridAccessV1Object
		if ans.MfaVendorType.RsaSecuridAccessV1 != nil {
			var4 = &mfaServersDsModelRsaSecuridAccessV1Object{}
			var4.RsaAccessid = types.StringValue(ans.MfaVendorType.RsaSecuridAccessV1.RsaAccessid)
			var4.RsaAccesskey = types.StringValue(ans.MfaVendorType.RsaSecuridAccessV1.RsaAccesskey)
			var4.RsaApiHost = types.StringValue(ans.MfaVendorType.RsaSecuridAccessV1.RsaApiHost)
			var4.RsaAssurancepolicyid = types.StringValue(ans.MfaVendorType.RsaSecuridAccessV1.RsaAssurancepolicyid)
			var4.RsaBaseuri = types.StringValue(ans.MfaVendorType.RsaSecuridAccessV1.RsaBaseuri)
			var4.RsaTimeout = types.StringValue(ans.MfaVendorType.RsaSecuridAccessV1.RsaTimeout)
		}
		var0.DuoSecurityV2 = var1
		var0.OktaAdaptiveV1 = var2
		var0.PingIdentityV1 = var3
		var0.RsaSecuridAccessV1 = var4
	}
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.MfaCertProfile = types.StringValue(ans.MfaCertProfile)
	state.MfaVendorType = var0
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &mfaServersResource{}
	_ resource.ResourceWithConfigure   = &mfaServersResource{}
	_ resource.ResourceWithImportState = &mfaServersResource{}
)

func NewMfaServersResource() resource.Resource {
	return &mfaServersResource{}
}

type mfaServersResource struct {
	client *sase.Client
}

type mfaServersRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Position types.String `tfsdk:"position"`
	Folder   types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/mfa-servers
	ObjectId       types.String                          `tfsdk:"object_id"`
	MfaCertProfile types.String                          `tfsdk:"mfa_cert_profile"`
	MfaVendorType  *mfaServersRsModelMfaVendorTypeObject `tfsdk:"mfa_vendor_type"`
	Name           types.String                          `tfsdk:"name"`
}

type mfaServersRsModelMfaVendorTypeObject struct {
	DuoSecurityV2      *mfaServersRsModelDuoSecurityV2Object      `tfsdk:"duo_security_v2"`
	OktaAdaptiveV1     *mfaServersRsModelOktaAdaptiveV1Object     `tfsdk:"okta_adaptive_v1"`
	PingIdentityV1     *mfaServersRsModelPingIdentityV1Object     `tfsdk:"ping_identity_v1"`
	RsaSecuridAccessV1 *mfaServersRsModelRsaSecuridAccessV1Object `tfsdk:"rsa_securid_access_v1"`
}

type mfaServersRsModelDuoSecurityV2Object struct {
	DuoApiHost        types.String `tfsdk:"duo_api_host"`
	DuoBaseuri        types.String `tfsdk:"duo_baseuri"`
	DuoIntegrationKey types.String `tfsdk:"duo_integration_key"`
	DuoSecretKey      types.String `tfsdk:"duo_secret_key"`
	DuoTimeout        types.String `tfsdk:"duo_timeout"`
}

type mfaServersRsModelOktaAdaptiveV1Object struct {
	OktaApiHost types.String `tfsdk:"okta_api_host"`
	OktaBaseuri types.String `tfsdk:"okta_baseuri"`
	OktaOrg     types.String `tfsdk:"okta_org"`
	OktaTimeout types.String `tfsdk:"okta_timeout"`
	OktaToken   types.String `tfsdk:"okta_token"`
}

type mfaServersRsModelPingIdentityV1Object struct {
	PingApiHost  types.String `tfsdk:"ping_api_host"`
	PingBaseuri  types.String `tfsdk:"ping_baseuri"`
	PingOrg      types.String `tfsdk:"ping_org"`
	PingOrgAlias types.String `tfsdk:"ping_org_alias"`
	PingTimeout  types.String `tfsdk:"ping_timeout"`
	PingToken    types.String `tfsdk:"ping_token"`
}

type mfaServersRsModelRsaSecuridAccessV1Object struct {
	RsaAccessid          types.String `tfsdk:"rsa_accessid"`
	RsaAccesskey         types.String `tfsdk:"rsa_accesskey"`
	RsaApiHost           types.String `tfsdk:"rsa_api_host"`
	RsaAssurancepolicyid types.String `tfsdk:"rsa_assurancepolicyid"`
	RsaBaseuri           types.String `tfsdk:"rsa_baseuri"`
	RsaTimeout           types.String `tfsdk:"rsa_timeout"`
}

// Metadata returns the data source type name.
func (r *mfaServersResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mfa_servers"
}

// Schema defines the schema for this listing data source.
func (r *mfaServersResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"position": rsschema.StringAttribute{
				Description:         "The position of a security rule. Value must be one of: `\"pre\"`, `\"post\"`.",
				MarkdownDescription: "The position of a security rule. Value must be one of: `\"pre\"`, `\"post\"`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("pre", "post"),
				},
			},
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
			"mfa_cert_profile": rsschema.StringAttribute{
				Description:         "The `mfa_cert_profile` parameter.",
				MarkdownDescription: "The `mfa_cert_profile` parameter.",
				Required:            true,
			},
			"mfa_vendor_type": rsschema.SingleNestedAttribute{
				Description:         "The `mfa_vendor_type` parameter.",
				MarkdownDescription: "The `mfa_vendor_type` parameter.",
				Optional:            true,
				Attributes: map[string]rsschema.Attribute{
					"duo_security_v2": rsschema.SingleNestedAttribute{
						Description:         "The `duo_security_v2` parameter.",
						MarkdownDescription: "The `duo_security_v2` parameter.",
						Optional:            true,
						Attributes: map[string]rsschema.Attribute{
							"duo_api_host": rsschema.StringAttribute{
								Description:         "The `duo_api_host` parameter.",
								MarkdownDescription: "The `duo_api_host` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"duo_baseuri": rsschema.StringAttribute{
								Description:         "The `duo_baseuri` parameter.",
								MarkdownDescription: "The `duo_baseuri` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"duo_integration_key": rsschema.StringAttribute{
								Description:         "The `duo_integration_key` parameter.",
								MarkdownDescription: "The `duo_integration_key` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"duo_secret_key": rsschema.StringAttribute{
								Description:         "The `duo_secret_key` parameter.",
								MarkdownDescription: "The `duo_secret_key` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"duo_timeout": rsschema.StringAttribute{
								Description:         "The `duo_timeout` parameter.",
								MarkdownDescription: "The `duo_timeout` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
						},
					},
					"okta_adaptive_v1": rsschema.SingleNestedAttribute{
						Description:         "The `okta_adaptive_v1` parameter.",
						MarkdownDescription: "The `okta_adaptive_v1` parameter.",
						Optional:            true,
						Attributes: map[string]rsschema.Attribute{
							"okta_api_host": rsschema.StringAttribute{
								Description:         "The `okta_api_host` parameter.",
								MarkdownDescription: "The `okta_api_host` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"okta_baseuri": rsschema.StringAttribute{
								Description:         "The `okta_baseuri` parameter.",
								MarkdownDescription: "The `okta_baseuri` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"okta_org": rsschema.StringAttribute{
								Description:         "The `okta_org` parameter.",
								MarkdownDescription: "The `okta_org` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"okta_timeout": rsschema.StringAttribute{
								Description:         "The `okta_timeout` parameter.",
								MarkdownDescription: "The `okta_timeout` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"okta_token": rsschema.StringAttribute{
								Description:         "The `okta_token` parameter.",
								MarkdownDescription: "The `okta_token` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
						},
					},
					"ping_identity_v1": rsschema.SingleNestedAttribute{
						Description:         "The `ping_identity_v1` parameter.",
						MarkdownDescription: "The `ping_identity_v1` parameter.",
						Optional:            true,
						Attributes: map[string]rsschema.Attribute{
							"ping_api_host": rsschema.StringAttribute{
								Description:         "The `ping_api_host` parameter.",
								MarkdownDescription: "The `ping_api_host` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"ping_baseuri": rsschema.StringAttribute{
								Description:         "The `ping_baseuri` parameter.",
								MarkdownDescription: "The `ping_baseuri` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"ping_org": rsschema.StringAttribute{
								Description:         "The `ping_org` parameter.",
								MarkdownDescription: "The `ping_org` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"ping_org_alias": rsschema.StringAttribute{
								Description:         "The `ping_org_alias` parameter.",
								MarkdownDescription: "The `ping_org_alias` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"ping_timeout": rsschema.StringAttribute{
								Description:         "The `ping_timeout` parameter.",
								MarkdownDescription: "The `ping_timeout` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"ping_token": rsschema.StringAttribute{
								Description:         "The `ping_token` parameter.",
								MarkdownDescription: "The `ping_token` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
						},
					},
					"rsa_securid_access_v1": rsschema.SingleNestedAttribute{
						Description:         "The `rsa_securid_access_v1` parameter.",
						MarkdownDescription: "The `rsa_securid_access_v1` parameter.",
						Optional:            true,
						Attributes: map[string]rsschema.Attribute{
							"rsa_accessid": rsschema.StringAttribute{
								Description:         "The `rsa_accessid` parameter.",
								MarkdownDescription: "The `rsa_accessid` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"rsa_accesskey": rsschema.StringAttribute{
								Description:         "The `rsa_accesskey` parameter.",
								MarkdownDescription: "The `rsa_accesskey` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"rsa_api_host": rsschema.StringAttribute{
								Description:         "The `rsa_api_host` parameter.",
								MarkdownDescription: "The `rsa_api_host` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"rsa_assurancepolicyid": rsschema.StringAttribute{
								Description:         "The `rsa_assurancepolicyid` parameter.",
								MarkdownDescription: "The `rsa_assurancepolicyid` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"rsa_baseuri": rsschema.StringAttribute{
								Description:         "The `rsa_baseuri` parameter.",
								MarkdownDescription: "The `rsa_baseuri` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"rsa_timeout": rsschema.StringAttribute{
								Description:         "The `rsa_timeout` parameter.",
								MarkdownDescription: "The `rsa_timeout` parameter.",
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
			"name": rsschema.StringAttribute{
				Description:         "The `name` parameter.",
				MarkdownDescription: "The `name` parameter.",
				Required:            true,
			},
		},
	}
}

// Configure prepares the struct.
func (r *mfaServersResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *mfaServersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state mfaServersRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_mfa_servers",
		"position":                    state.Position.ValueString(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := wArkOsV.NewClient(r.client)
	input := wArkOsV.CreateInput{
		Position: state.Position.ValueString(),
		Folder:   state.Folder.ValueString(),
	}
	var var0 deRyMEf.Config
	var0.MfaCertProfile = state.MfaCertProfile.ValueString()
	var var1 *deRyMEf.MfaVendorTypeObject
	if state.MfaVendorType != nil {
		var1 = &deRyMEf.MfaVendorTypeObject{}
		var var2 *deRyMEf.DuoSecurityV2Object
		if state.MfaVendorType.DuoSecurityV2 != nil {
			var2 = &deRyMEf.DuoSecurityV2Object{}
			var2.DuoApiHost = state.MfaVendorType.DuoSecurityV2.DuoApiHost.ValueString()
			var2.DuoBaseuri = state.MfaVendorType.DuoSecurityV2.DuoBaseuri.ValueString()
			var2.DuoIntegrationKey = state.MfaVendorType.DuoSecurityV2.DuoIntegrationKey.ValueString()
			var2.DuoSecretKey = state.MfaVendorType.DuoSecurityV2.DuoSecretKey.ValueString()
			var2.DuoTimeout = state.MfaVendorType.DuoSecurityV2.DuoTimeout.ValueString()
		}
		var1.DuoSecurityV2 = var2
		var var3 *deRyMEf.OktaAdaptiveV1Object
		if state.MfaVendorType.OktaAdaptiveV1 != nil {
			var3 = &deRyMEf.OktaAdaptiveV1Object{}
			var3.OktaApiHost = state.MfaVendorType.OktaAdaptiveV1.OktaApiHost.ValueString()
			var3.OktaBaseuri = state.MfaVendorType.OktaAdaptiveV1.OktaBaseuri.ValueString()
			var3.OktaOrg = state.MfaVendorType.OktaAdaptiveV1.OktaOrg.ValueString()
			var3.OktaTimeout = state.MfaVendorType.OktaAdaptiveV1.OktaTimeout.ValueString()
			var3.OktaToken = state.MfaVendorType.OktaAdaptiveV1.OktaToken.ValueString()
		}
		var1.OktaAdaptiveV1 = var3
		var var4 *deRyMEf.PingIdentityV1Object
		if state.MfaVendorType.PingIdentityV1 != nil {
			var4 = &deRyMEf.PingIdentityV1Object{}
			var4.PingApiHost = state.MfaVendorType.PingIdentityV1.PingApiHost.ValueString()
			var4.PingBaseuri = state.MfaVendorType.PingIdentityV1.PingBaseuri.ValueString()
			var4.PingOrg = state.MfaVendorType.PingIdentityV1.PingOrg.ValueString()
			var4.PingOrgAlias = state.MfaVendorType.PingIdentityV1.PingOrgAlias.ValueString()
			var4.PingTimeout = state.MfaVendorType.PingIdentityV1.PingTimeout.ValueString()
			var4.PingToken = state.MfaVendorType.PingIdentityV1.PingToken.ValueString()
		}
		var1.PingIdentityV1 = var4
		var var5 *deRyMEf.RsaSecuridAccessV1Object
		if state.MfaVendorType.RsaSecuridAccessV1 != nil {
			var5 = &deRyMEf.RsaSecuridAccessV1Object{}
			var5.RsaAccessid = state.MfaVendorType.RsaSecuridAccessV1.RsaAccessid.ValueString()
			var5.RsaAccesskey = state.MfaVendorType.RsaSecuridAccessV1.RsaAccesskey.ValueString()
			var5.RsaApiHost = state.MfaVendorType.RsaSecuridAccessV1.RsaApiHost.ValueString()
			var5.RsaAssurancepolicyid = state.MfaVendorType.RsaSecuridAccessV1.RsaAssurancepolicyid.ValueString()
			var5.RsaBaseuri = state.MfaVendorType.RsaSecuridAccessV1.RsaBaseuri.ValueString()
			var5.RsaTimeout = state.MfaVendorType.RsaSecuridAccessV1.RsaTimeout.ValueString()
		}
		var1.RsaSecuridAccessV1 = var5
	}
	var0.MfaVendorType = var1
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
	idBuilder.WriteString(input.Position)
	idBuilder.WriteString(IdSeparator)
	idBuilder.WriteString(input.Folder)
	idBuilder.WriteString(IdSeparator)
	idBuilder.WriteString(ans.ObjectId)
	state.Id = types.StringValue(idBuilder.String())
	var var6 *mfaServersRsModelMfaVendorTypeObject
	if ans.MfaVendorType != nil {
		var6 = &mfaServersRsModelMfaVendorTypeObject{}
		var var7 *mfaServersRsModelDuoSecurityV2Object
		if ans.MfaVendorType.DuoSecurityV2 != nil {
			var7 = &mfaServersRsModelDuoSecurityV2Object{}
			var7.DuoApiHost = types.StringValue(ans.MfaVendorType.DuoSecurityV2.DuoApiHost)
			var7.DuoBaseuri = types.StringValue(ans.MfaVendorType.DuoSecurityV2.DuoBaseuri)
			var7.DuoIntegrationKey = types.StringValue(ans.MfaVendorType.DuoSecurityV2.DuoIntegrationKey)
			var7.DuoSecretKey = types.StringValue(ans.MfaVendorType.DuoSecurityV2.DuoSecretKey)
			var7.DuoTimeout = types.StringValue(ans.MfaVendorType.DuoSecurityV2.DuoTimeout)
		}
		var var8 *mfaServersRsModelOktaAdaptiveV1Object
		if ans.MfaVendorType.OktaAdaptiveV1 != nil {
			var8 = &mfaServersRsModelOktaAdaptiveV1Object{}
			var8.OktaApiHost = types.StringValue(ans.MfaVendorType.OktaAdaptiveV1.OktaApiHost)
			var8.OktaBaseuri = types.StringValue(ans.MfaVendorType.OktaAdaptiveV1.OktaBaseuri)
			var8.OktaOrg = types.StringValue(ans.MfaVendorType.OktaAdaptiveV1.OktaOrg)
			var8.OktaTimeout = types.StringValue(ans.MfaVendorType.OktaAdaptiveV1.OktaTimeout)
			var8.OktaToken = types.StringValue(ans.MfaVendorType.OktaAdaptiveV1.OktaToken)
		}
		var var9 *mfaServersRsModelPingIdentityV1Object
		if ans.MfaVendorType.PingIdentityV1 != nil {
			var9 = &mfaServersRsModelPingIdentityV1Object{}
			var9.PingApiHost = types.StringValue(ans.MfaVendorType.PingIdentityV1.PingApiHost)
			var9.PingBaseuri = types.StringValue(ans.MfaVendorType.PingIdentityV1.PingBaseuri)
			var9.PingOrg = types.StringValue(ans.MfaVendorType.PingIdentityV1.PingOrg)
			var9.PingOrgAlias = types.StringValue(ans.MfaVendorType.PingIdentityV1.PingOrgAlias)
			var9.PingTimeout = types.StringValue(ans.MfaVendorType.PingIdentityV1.PingTimeout)
			var9.PingToken = types.StringValue(ans.MfaVendorType.PingIdentityV1.PingToken)
		}
		var var10 *mfaServersRsModelRsaSecuridAccessV1Object
		if ans.MfaVendorType.RsaSecuridAccessV1 != nil {
			var10 = &mfaServersRsModelRsaSecuridAccessV1Object{}
			var10.RsaAccessid = types.StringValue(ans.MfaVendorType.RsaSecuridAccessV1.RsaAccessid)
			var10.RsaAccesskey = types.StringValue(ans.MfaVendorType.RsaSecuridAccessV1.RsaAccesskey)
			var10.RsaApiHost = types.StringValue(ans.MfaVendorType.RsaSecuridAccessV1.RsaApiHost)
			var10.RsaAssurancepolicyid = types.StringValue(ans.MfaVendorType.RsaSecuridAccessV1.RsaAssurancepolicyid)
			var10.RsaBaseuri = types.StringValue(ans.MfaVendorType.RsaSecuridAccessV1.RsaBaseuri)
			var10.RsaTimeout = types.StringValue(ans.MfaVendorType.RsaSecuridAccessV1.RsaTimeout)
		}
		var6.DuoSecurityV2 = var7
		var6.OktaAdaptiveV1 = var8
		var6.PingIdentityV1 = var9
		var6.RsaSecuridAccessV1 = var10
	}
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.MfaCertProfile = types.StringValue(ans.MfaCertProfile)
	state.MfaVendorType = var6
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *mfaServersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var idType types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &idType)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := idType.ValueString()
	tokens := strings.Split(id, IdSeparator)
	if len(tokens) != 3 {
		resp.Diagnostics.AddError("Error in resource ID format", "Expected 3 tokens")
		return
	}

	var state mfaServersRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_mfa_servers",
		"locMap":                      map[string]int{"Folder": 1, "ObjectId": 2, "Position": 0},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := wArkOsV.NewClient(r.client)
	input := wArkOsV.ReadInput{
		ObjectId: tokens[2],
		Folder:   tokens[1],
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
	state.Position = types.StringValue(tokens[0])
	state.Folder = types.StringValue(tokens[1])
	state.Id = idType
	var var0 *mfaServersRsModelMfaVendorTypeObject
	if ans.MfaVendorType != nil {
		var0 = &mfaServersRsModelMfaVendorTypeObject{}
		var var1 *mfaServersRsModelDuoSecurityV2Object
		if ans.MfaVendorType.DuoSecurityV2 != nil {
			var1 = &mfaServersRsModelDuoSecurityV2Object{}
			var1.DuoApiHost = types.StringValue(ans.MfaVendorType.DuoSecurityV2.DuoApiHost)
			var1.DuoBaseuri = types.StringValue(ans.MfaVendorType.DuoSecurityV2.DuoBaseuri)
			var1.DuoIntegrationKey = types.StringValue(ans.MfaVendorType.DuoSecurityV2.DuoIntegrationKey)
			var1.DuoSecretKey = types.StringValue(ans.MfaVendorType.DuoSecurityV2.DuoSecretKey)
			var1.DuoTimeout = types.StringValue(ans.MfaVendorType.DuoSecurityV2.DuoTimeout)
		}
		var var2 *mfaServersRsModelOktaAdaptiveV1Object
		if ans.MfaVendorType.OktaAdaptiveV1 != nil {
			var2 = &mfaServersRsModelOktaAdaptiveV1Object{}
			var2.OktaApiHost = types.StringValue(ans.MfaVendorType.OktaAdaptiveV1.OktaApiHost)
			var2.OktaBaseuri = types.StringValue(ans.MfaVendorType.OktaAdaptiveV1.OktaBaseuri)
			var2.OktaOrg = types.StringValue(ans.MfaVendorType.OktaAdaptiveV1.OktaOrg)
			var2.OktaTimeout = types.StringValue(ans.MfaVendorType.OktaAdaptiveV1.OktaTimeout)
			var2.OktaToken = types.StringValue(ans.MfaVendorType.OktaAdaptiveV1.OktaToken)
		}
		var var3 *mfaServersRsModelPingIdentityV1Object
		if ans.MfaVendorType.PingIdentityV1 != nil {
			var3 = &mfaServersRsModelPingIdentityV1Object{}
			var3.PingApiHost = types.StringValue(ans.MfaVendorType.PingIdentityV1.PingApiHost)
			var3.PingBaseuri = types.StringValue(ans.MfaVendorType.PingIdentityV1.PingBaseuri)
			var3.PingOrg = types.StringValue(ans.MfaVendorType.PingIdentityV1.PingOrg)
			var3.PingOrgAlias = types.StringValue(ans.MfaVendorType.PingIdentityV1.PingOrgAlias)
			var3.PingTimeout = types.StringValue(ans.MfaVendorType.PingIdentityV1.PingTimeout)
			var3.PingToken = types.StringValue(ans.MfaVendorType.PingIdentityV1.PingToken)
		}
		var var4 *mfaServersRsModelRsaSecuridAccessV1Object
		if ans.MfaVendorType.RsaSecuridAccessV1 != nil {
			var4 = &mfaServersRsModelRsaSecuridAccessV1Object{}
			var4.RsaAccessid = types.StringValue(ans.MfaVendorType.RsaSecuridAccessV1.RsaAccessid)
			var4.RsaAccesskey = types.StringValue(ans.MfaVendorType.RsaSecuridAccessV1.RsaAccesskey)
			var4.RsaApiHost = types.StringValue(ans.MfaVendorType.RsaSecuridAccessV1.RsaApiHost)
			var4.RsaAssurancepolicyid = types.StringValue(ans.MfaVendorType.RsaSecuridAccessV1.RsaAssurancepolicyid)
			var4.RsaBaseuri = types.StringValue(ans.MfaVendorType.RsaSecuridAccessV1.RsaBaseuri)
			var4.RsaTimeout = types.StringValue(ans.MfaVendorType.RsaSecuridAccessV1.RsaTimeout)
		}
		var0.DuoSecurityV2 = var1
		var0.OktaAdaptiveV1 = var2
		var0.PingIdentityV1 = var3
		var0.RsaSecuridAccessV1 = var4
	}
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.MfaCertProfile = types.StringValue(ans.MfaCertProfile)
	state.MfaVendorType = var0
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *mfaServersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state mfaServersRsModel
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
		"resource_name":               "sase_mfa_servers",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := wArkOsV.NewClient(r.client)
	input := wArkOsV.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 deRyMEf.Config
	var0.MfaCertProfile = plan.MfaCertProfile.ValueString()
	var var1 *deRyMEf.MfaVendorTypeObject
	if plan.MfaVendorType != nil {
		var1 = &deRyMEf.MfaVendorTypeObject{}
		var var2 *deRyMEf.DuoSecurityV2Object
		if plan.MfaVendorType.DuoSecurityV2 != nil {
			var2 = &deRyMEf.DuoSecurityV2Object{}
			var2.DuoApiHost = plan.MfaVendorType.DuoSecurityV2.DuoApiHost.ValueString()
			var2.DuoBaseuri = plan.MfaVendorType.DuoSecurityV2.DuoBaseuri.ValueString()
			var2.DuoIntegrationKey = plan.MfaVendorType.DuoSecurityV2.DuoIntegrationKey.ValueString()
			var2.DuoSecretKey = plan.MfaVendorType.DuoSecurityV2.DuoSecretKey.ValueString()
			var2.DuoTimeout = plan.MfaVendorType.DuoSecurityV2.DuoTimeout.ValueString()
		}
		var1.DuoSecurityV2 = var2
		var var3 *deRyMEf.OktaAdaptiveV1Object
		if plan.MfaVendorType.OktaAdaptiveV1 != nil {
			var3 = &deRyMEf.OktaAdaptiveV1Object{}
			var3.OktaApiHost = plan.MfaVendorType.OktaAdaptiveV1.OktaApiHost.ValueString()
			var3.OktaBaseuri = plan.MfaVendorType.OktaAdaptiveV1.OktaBaseuri.ValueString()
			var3.OktaOrg = plan.MfaVendorType.OktaAdaptiveV1.OktaOrg.ValueString()
			var3.OktaTimeout = plan.MfaVendorType.OktaAdaptiveV1.OktaTimeout.ValueString()
			var3.OktaToken = plan.MfaVendorType.OktaAdaptiveV1.OktaToken.ValueString()
		}
		var1.OktaAdaptiveV1 = var3
		var var4 *deRyMEf.PingIdentityV1Object
		if plan.MfaVendorType.PingIdentityV1 != nil {
			var4 = &deRyMEf.PingIdentityV1Object{}
			var4.PingApiHost = plan.MfaVendorType.PingIdentityV1.PingApiHost.ValueString()
			var4.PingBaseuri = plan.MfaVendorType.PingIdentityV1.PingBaseuri.ValueString()
			var4.PingOrg = plan.MfaVendorType.PingIdentityV1.PingOrg.ValueString()
			var4.PingOrgAlias = plan.MfaVendorType.PingIdentityV1.PingOrgAlias.ValueString()
			var4.PingTimeout = plan.MfaVendorType.PingIdentityV1.PingTimeout.ValueString()
			var4.PingToken = plan.MfaVendorType.PingIdentityV1.PingToken.ValueString()
		}
		var1.PingIdentityV1 = var4
		var var5 *deRyMEf.RsaSecuridAccessV1Object
		if plan.MfaVendorType.RsaSecuridAccessV1 != nil {
			var5 = &deRyMEf.RsaSecuridAccessV1Object{}
			var5.RsaAccessid = plan.MfaVendorType.RsaSecuridAccessV1.RsaAccessid.ValueString()
			var5.RsaAccesskey = plan.MfaVendorType.RsaSecuridAccessV1.RsaAccesskey.ValueString()
			var5.RsaApiHost = plan.MfaVendorType.RsaSecuridAccessV1.RsaApiHost.ValueString()
			var5.RsaAssurancepolicyid = plan.MfaVendorType.RsaSecuridAccessV1.RsaAssurancepolicyid.ValueString()
			var5.RsaBaseuri = plan.MfaVendorType.RsaSecuridAccessV1.RsaBaseuri.ValueString()
			var5.RsaTimeout = plan.MfaVendorType.RsaSecuridAccessV1.RsaTimeout.ValueString()
		}
		var1.RsaSecuridAccessV1 = var5
	}
	var0.MfaVendorType = var1
	var0.Name = plan.Name.ValueString()
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var6 *mfaServersRsModelMfaVendorTypeObject
	if ans.MfaVendorType != nil {
		var6 = &mfaServersRsModelMfaVendorTypeObject{}
		var var7 *mfaServersRsModelDuoSecurityV2Object
		if ans.MfaVendorType.DuoSecurityV2 != nil {
			var7 = &mfaServersRsModelDuoSecurityV2Object{}
			var7.DuoApiHost = types.StringValue(ans.MfaVendorType.DuoSecurityV2.DuoApiHost)
			var7.DuoBaseuri = types.StringValue(ans.MfaVendorType.DuoSecurityV2.DuoBaseuri)
			var7.DuoIntegrationKey = types.StringValue(ans.MfaVendorType.DuoSecurityV2.DuoIntegrationKey)
			var7.DuoSecretKey = types.StringValue(ans.MfaVendorType.DuoSecurityV2.DuoSecretKey)
			var7.DuoTimeout = types.StringValue(ans.MfaVendorType.DuoSecurityV2.DuoTimeout)
		}
		var var8 *mfaServersRsModelOktaAdaptiveV1Object
		if ans.MfaVendorType.OktaAdaptiveV1 != nil {
			var8 = &mfaServersRsModelOktaAdaptiveV1Object{}
			var8.OktaApiHost = types.StringValue(ans.MfaVendorType.OktaAdaptiveV1.OktaApiHost)
			var8.OktaBaseuri = types.StringValue(ans.MfaVendorType.OktaAdaptiveV1.OktaBaseuri)
			var8.OktaOrg = types.StringValue(ans.MfaVendorType.OktaAdaptiveV1.OktaOrg)
			var8.OktaTimeout = types.StringValue(ans.MfaVendorType.OktaAdaptiveV1.OktaTimeout)
			var8.OktaToken = types.StringValue(ans.MfaVendorType.OktaAdaptiveV1.OktaToken)
		}
		var var9 *mfaServersRsModelPingIdentityV1Object
		if ans.MfaVendorType.PingIdentityV1 != nil {
			var9 = &mfaServersRsModelPingIdentityV1Object{}
			var9.PingApiHost = types.StringValue(ans.MfaVendorType.PingIdentityV1.PingApiHost)
			var9.PingBaseuri = types.StringValue(ans.MfaVendorType.PingIdentityV1.PingBaseuri)
			var9.PingOrg = types.StringValue(ans.MfaVendorType.PingIdentityV1.PingOrg)
			var9.PingOrgAlias = types.StringValue(ans.MfaVendorType.PingIdentityV1.PingOrgAlias)
			var9.PingTimeout = types.StringValue(ans.MfaVendorType.PingIdentityV1.PingTimeout)
			var9.PingToken = types.StringValue(ans.MfaVendorType.PingIdentityV1.PingToken)
		}
		var var10 *mfaServersRsModelRsaSecuridAccessV1Object
		if ans.MfaVendorType.RsaSecuridAccessV1 != nil {
			var10 = &mfaServersRsModelRsaSecuridAccessV1Object{}
			var10.RsaAccessid = types.StringValue(ans.MfaVendorType.RsaSecuridAccessV1.RsaAccessid)
			var10.RsaAccesskey = types.StringValue(ans.MfaVendorType.RsaSecuridAccessV1.RsaAccesskey)
			var10.RsaApiHost = types.StringValue(ans.MfaVendorType.RsaSecuridAccessV1.RsaApiHost)
			var10.RsaAssurancepolicyid = types.StringValue(ans.MfaVendorType.RsaSecuridAccessV1.RsaAssurancepolicyid)
			var10.RsaBaseuri = types.StringValue(ans.MfaVendorType.RsaSecuridAccessV1.RsaBaseuri)
			var10.RsaTimeout = types.StringValue(ans.MfaVendorType.RsaSecuridAccessV1.RsaTimeout)
		}
		var6.DuoSecurityV2 = var7
		var6.OktaAdaptiveV1 = var8
		var6.PingIdentityV1 = var9
		var6.RsaSecuridAccessV1 = var10
	}
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.MfaCertProfile = types.StringValue(ans.MfaCertProfile)
	state.MfaVendorType = var6
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *mfaServersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var idType types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &idType)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := idType.ValueString()
	tokens := strings.Split(id, IdSeparator)
	if len(tokens) != 3 {
		resp.Diagnostics.AddError("Error in resource ID format", "Expected 3 tokens")
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource delete", map[string]any{
		"terraform_provider_function": "Delete",
		"resource_name":               "sase_mfa_servers",
		"locMap":                      map[string]int{"Folder": 1, "ObjectId": 2, "Position": 0},
		"tokens":                      tokens,
	})

	svc := wArkOsV.NewClient(r.client)
	input := wArkOsV.DeleteInput{
		ObjectId: tokens[2],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *mfaServersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
