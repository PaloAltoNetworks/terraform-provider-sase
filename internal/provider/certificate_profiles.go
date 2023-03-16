package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	cozuxBy "github.com/paloaltonetworks/sase-go/netsec/schema/certificate/profiles"
	qLteaIq "github.com/paloaltonetworks/sase-go/netsec/service/v1/certificateprofiles"

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
	_ datasource.DataSource              = &certificateProfilesListDataSource{}
	_ datasource.DataSourceWithConfigure = &certificateProfilesListDataSource{}
)

func NewCertificateProfilesListDataSource() datasource.DataSource {
	return &certificateProfilesListDataSource{}
}

type certificateProfilesListDataSource struct {
	client *sase.Client
}

type certificateProfilesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []certificateProfilesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type certificateProfilesListDsModelConfig struct {
	BlockExpiredCert         types.Bool                                           `tfsdk:"block_expired_cert"`
	BlockTimeoutCert         types.Bool                                           `tfsdk:"block_timeout_cert"`
	BlockUnauthenticatedCert types.Bool                                           `tfsdk:"block_unauthenticated_cert"`
	BlockUnknownCert         types.Bool                                           `tfsdk:"block_unknown_cert"`
	CaCertificates           []certificateProfilesListDsModelCaCertificatesObject `tfsdk:"ca_certificates"`
	CertStatusTimeout        types.String                                         `tfsdk:"cert_status_timeout"`
	CrlReceiveTimeout        types.String                                         `tfsdk:"crl_receive_timeout"`
	Domain                   types.String                                         `tfsdk:"domain"`
	ObjectId                 types.String                                         `tfsdk:"object_id"`
	Name                     types.String                                         `tfsdk:"name"`
	OcspReceiveTimeout       types.String                                         `tfsdk:"ocsp_receive_timeout"`
	UseCrl                   types.Bool                                           `tfsdk:"use_crl"`
	UseOcsp                  types.Bool                                           `tfsdk:"use_ocsp"`
	UsernameField            *certificateProfilesListDsModelUsernameFieldObject   `tfsdk:"username_field"`
}

type certificateProfilesListDsModelCaCertificatesObject struct {
	DefaultOcspUrl types.String `tfsdk:"default_ocsp_url"`
	Name           types.String `tfsdk:"name"`
	OcspVerifyCert types.String `tfsdk:"ocsp_verify_cert"`
	TemplateName   types.String `tfsdk:"template_name"`
}

type certificateProfilesListDsModelUsernameFieldObject struct {
	Subject    types.String `tfsdk:"subject"`
	SubjectAlt types.String `tfsdk:"subject_alt"`
}

// Metadata returns the data source type name.
func (d *certificateProfilesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate_profiles_list"
}

// Schema defines the schema for this listing data source.
func (d *certificateProfilesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"block_expired_cert": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"block_timeout_cert": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"block_unauthenticated_cert": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"block_unknown_cert": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"ca_certificates": dsschema.ListNestedAttribute{
							Description: "",
							Computed:    true,
							NestedObject: dsschema.NestedAttributeObject{
								Attributes: map[string]dsschema.Attribute{
									"default_ocsp_url": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"name": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"ocsp_verify_cert": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"template_name": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
						},
						"cert_status_timeout": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"crl_receive_timeout": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"domain": dsschema.StringAttribute{
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
						"ocsp_receive_timeout": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"use_crl": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"use_ocsp": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"username_field": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"subject": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
								"subject_alt": dsschema.StringAttribute{
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
func (d *certificateProfilesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *certificateProfilesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state certificateProfilesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_certificate_profiles_list",
		"limit":            state.Limit.ValueInt64(),
		"has_limit":        !state.Limit.IsNull(),
		"offset":           state.Offset.ValueInt64(),
		"has_offset":       !state.Offset.IsNull(),
		"name":             state.Name.ValueString(),
		"has_name":         !state.Name.IsNull(),
		"folder":           state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := qLteaIq.NewClient(d.client)
	input := qLteaIq.ListInput{
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
	state.Id = types.StringValue(strings.Join([]string{strconv.FormatInt(*input.Limit, 10), strconv.FormatInt(*input.Offset, 10), *input.Name, input.Folder}, IdSeparator))
	var var0 []certificateProfilesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]certificateProfilesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 certificateProfilesListDsModelConfig
			var var3 []certificateProfilesListDsModelCaCertificatesObject
			if len(var1.CaCertificates) != 0 {
				var3 = make([]certificateProfilesListDsModelCaCertificatesObject, 0, len(var1.CaCertificates))
				for var4Index := range var1.CaCertificates {
					var4 := var1.CaCertificates[var4Index]
					var var5 certificateProfilesListDsModelCaCertificatesObject
					var5.DefaultOcspUrl = types.StringValue(var4.DefaultOcspUrl)
					var5.Name = types.StringValue(var4.Name)
					var5.OcspVerifyCert = types.StringValue(var4.OcspVerifyCert)
					var5.TemplateName = types.StringValue(var4.TemplateName)
					var3 = append(var3, var5)
				}
			}
			var var6 *certificateProfilesListDsModelUsernameFieldObject
			if var1.UsernameField != nil {
				var6 = &certificateProfilesListDsModelUsernameFieldObject{}
				var6.Subject = types.StringValue(var1.UsernameField.Subject)
				var6.SubjectAlt = types.StringValue(var1.UsernameField.SubjectAlt)
			}
			var2.BlockExpiredCert = types.BoolValue(var1.BlockExpiredCert)
			var2.BlockTimeoutCert = types.BoolValue(var1.BlockTimeoutCert)
			var2.BlockUnauthenticatedCert = types.BoolValue(var1.BlockUnauthenticatedCert)
			var2.BlockUnknownCert = types.BoolValue(var1.BlockUnknownCert)
			var2.CaCertificates = var3
			var2.CertStatusTimeout = types.StringValue(var1.CertStatusTimeout)
			var2.CrlReceiveTimeout = types.StringValue(var1.CrlReceiveTimeout)
			var2.Domain = types.StringValue(var1.Domain)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Name = types.StringValue(var1.Name)
			var2.OcspReceiveTimeout = types.StringValue(var1.OcspReceiveTimeout)
			var2.UseCrl = types.BoolValue(var1.UseCrl)
			var2.UseOcsp = types.BoolValue(var1.UseOcsp)
			var2.UsernameField = var6
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
	_ datasource.DataSource              = &certificateProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &certificateProfilesDataSource{}
)

func NewCertificateProfilesDataSource() datasource.DataSource {
	return &certificateProfilesDataSource{}
}

type certificateProfilesDataSource struct {
	client *sase.Client
}

type certificateProfilesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`

	// Output.
	// Ref: #/components/schemas/certificate-profiles
	BlockExpiredCert         types.Bool                                       `tfsdk:"block_expired_cert"`
	BlockTimeoutCert         types.Bool                                       `tfsdk:"block_timeout_cert"`
	BlockUnauthenticatedCert types.Bool                                       `tfsdk:"block_unauthenticated_cert"`
	BlockUnknownCert         types.Bool                                       `tfsdk:"block_unknown_cert"`
	CaCertificates           []certificateProfilesDsModelCaCertificatesObject `tfsdk:"ca_certificates"`
	CertStatusTimeout        types.String                                     `tfsdk:"cert_status_timeout"`
	CrlReceiveTimeout        types.String                                     `tfsdk:"crl_receive_timeout"`
	Domain                   types.String                                     `tfsdk:"domain"`
	// input omit: ObjectId
	Name               types.String                                   `tfsdk:"name"`
	OcspReceiveTimeout types.String                                   `tfsdk:"ocsp_receive_timeout"`
	UseCrl             types.Bool                                     `tfsdk:"use_crl"`
	UseOcsp            types.Bool                                     `tfsdk:"use_ocsp"`
	UsernameField      *certificateProfilesDsModelUsernameFieldObject `tfsdk:"username_field"`
}

type certificateProfilesDsModelCaCertificatesObject struct {
	DefaultOcspUrl types.String `tfsdk:"default_ocsp_url"`
	Name           types.String `tfsdk:"name"`
	OcspVerifyCert types.String `tfsdk:"ocsp_verify_cert"`
	TemplateName   types.String `tfsdk:"template_name"`
}

type certificateProfilesDsModelUsernameFieldObject struct {
	Subject    types.String `tfsdk:"subject"`
	SubjectAlt types.String `tfsdk:"subject_alt"`
}

// Metadata returns the data source type name.
func (d *certificateProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate_profiles"
}

// Schema defines the schema for this listing data source.
func (d *certificateProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"block_expired_cert": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"block_timeout_cert": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"block_unauthenticated_cert": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"block_unknown_cert": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"ca_certificates": dsschema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"default_ocsp_url": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"ocsp_verify_cert": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"template_name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
					},
				},
			},
			"cert_status_timeout": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"crl_receive_timeout": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"domain": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"ocsp_receive_timeout": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"use_crl": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"use_ocsp": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"username_field": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"subject": dsschema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"subject_alt": dsschema.StringAttribute{
						Description: "",
						Computed:    true,
					},
				},
			},
		},
	}
}

// Configure prepares the struct.
func (d *certificateProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *certificateProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state certificateProfilesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_certificate_profiles",
		"object_id":        state.ObjectId.ValueString(),
	})

	// Prepare to run the command.
	svc := qLteaIq.NewClient(d.client)
	input := qLteaIq.ReadInput{
		ObjectId: state.ObjectId.ValueString(),
	}

	// Perform the operation.
	ans, err := svc.Read(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error getting singleton", err.Error())
		return
	}

	// Store the answer to state.
	state.Id = types.StringValue(strings.Join([]string{input.ObjectId}, IdSeparator))
	var var0 []certificateProfilesDsModelCaCertificatesObject
	if len(ans.CaCertificates) != 0 {
		var0 = make([]certificateProfilesDsModelCaCertificatesObject, 0, len(ans.CaCertificates))
		for var1Index := range ans.CaCertificates {
			var1 := ans.CaCertificates[var1Index]
			var var2 certificateProfilesDsModelCaCertificatesObject
			var2.DefaultOcspUrl = types.StringValue(var1.DefaultOcspUrl)
			var2.Name = types.StringValue(var1.Name)
			var2.OcspVerifyCert = types.StringValue(var1.OcspVerifyCert)
			var2.TemplateName = types.StringValue(var1.TemplateName)
			var0 = append(var0, var2)
		}
	}
	var var3 *certificateProfilesDsModelUsernameFieldObject
	if ans.UsernameField != nil {
		var3 = &certificateProfilesDsModelUsernameFieldObject{}
		var3.Subject = types.StringValue(ans.UsernameField.Subject)
		var3.SubjectAlt = types.StringValue(ans.UsernameField.SubjectAlt)
	}
	state.BlockExpiredCert = types.BoolValue(ans.BlockExpiredCert)
	state.BlockTimeoutCert = types.BoolValue(ans.BlockTimeoutCert)
	state.BlockUnauthenticatedCert = types.BoolValue(ans.BlockUnauthenticatedCert)
	state.BlockUnknownCert = types.BoolValue(ans.BlockUnknownCert)
	state.CaCertificates = var0
	state.CertStatusTimeout = types.StringValue(ans.CertStatusTimeout)
	state.CrlReceiveTimeout = types.StringValue(ans.CrlReceiveTimeout)
	state.Domain = types.StringValue(ans.Domain)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.OcspReceiveTimeout = types.StringValue(ans.OcspReceiveTimeout)
	state.UseCrl = types.BoolValue(ans.UseCrl)
	state.UseOcsp = types.BoolValue(ans.UseOcsp)
	state.UsernameField = var3

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &certificateProfilesResource{}
	_ resource.ResourceWithConfigure   = &certificateProfilesResource{}
	_ resource.ResourceWithImportState = &certificateProfilesResource{}
)

func NewCertificateProfilesResource() resource.Resource {
	return &certificateProfilesResource{}
}

type certificateProfilesResource struct {
	client *sase.Client
}

type certificateProfilesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/certificate-profiles
	BlockExpiredCert         types.Bool                                       `tfsdk:"block_expired_cert"`
	BlockTimeoutCert         types.Bool                                       `tfsdk:"block_timeout_cert"`
	BlockUnauthenticatedCert types.Bool                                       `tfsdk:"block_unauthenticated_cert"`
	BlockUnknownCert         types.Bool                                       `tfsdk:"block_unknown_cert"`
	CaCertificates           []certificateProfilesRsModelCaCertificatesObject `tfsdk:"ca_certificates"`
	CertStatusTimeout        types.String                                     `tfsdk:"cert_status_timeout"`
	CrlReceiveTimeout        types.String                                     `tfsdk:"crl_receive_timeout"`
	Domain                   types.String                                     `tfsdk:"domain"`
	ObjectId                 types.String                                     `tfsdk:"object_id"`
	Name                     types.String                                     `tfsdk:"name"`
	OcspReceiveTimeout       types.String                                     `tfsdk:"ocsp_receive_timeout"`
	UseCrl                   types.Bool                                       `tfsdk:"use_crl"`
	UseOcsp                  types.Bool                                       `tfsdk:"use_ocsp"`
	UsernameField            *certificateProfilesRsModelUsernameFieldObject   `tfsdk:"username_field"`
}

type certificateProfilesRsModelCaCertificatesObject struct {
	DefaultOcspUrl types.String `tfsdk:"default_ocsp_url"`
	Name           types.String `tfsdk:"name"`
	OcspVerifyCert types.String `tfsdk:"ocsp_verify_cert"`
	TemplateName   types.String `tfsdk:"template_name"`
}

type certificateProfilesRsModelUsernameFieldObject struct {
	Subject    types.String `tfsdk:"subject"`
	SubjectAlt types.String `tfsdk:"subject_alt"`
}

// Metadata returns the data source type name.
func (r *certificateProfilesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate_profiles"
}

// Schema defines the schema for this listing data source.
func (r *certificateProfilesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

			"block_expired_cert": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"block_timeout_cert": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"block_unauthenticated_cert": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"block_unknown_cert": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"ca_certificates": rsschema.ListNestedAttribute{
				Description: "",
				Required:    true,
				NestedObject: rsschema.NestedAttributeObject{
					Attributes: map[string]rsschema.Attribute{
						"default_ocsp_url": rsschema.StringAttribute{
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
						"ocsp_verify_cert": rsschema.StringAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
						},
						"template_name": rsschema.StringAttribute{
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
			"cert_status_timeout": rsschema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
			},
			"crl_receive_timeout": rsschema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
			},
			"domain": rsschema.StringAttribute{
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
				Required:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(63),
				},
			},
			"ocsp_receive_timeout": rsschema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
			},
			"use_crl": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"use_ocsp": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"username_field": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"subject": rsschema.StringAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							DefaultString(""),
						},
						Validators: []validator.String{
							stringvalidator.OneOf("common-name"),
						},
					},
					"subject_alt": rsschema.StringAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							DefaultString(""),
						},
						Validators: []validator.String{
							stringvalidator.OneOf("email"),
						},
					},
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *certificateProfilesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *certificateProfilesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state certificateProfilesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "sase_certificate_profiles",
		"folder":        state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := qLteaIq.NewClient(r.client)
	input := qLteaIq.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 cozuxBy.Config
	var0.BlockExpiredCert = state.BlockExpiredCert.ValueBool()
	var0.BlockTimeoutCert = state.BlockTimeoutCert.ValueBool()
	var0.BlockUnauthenticatedCert = state.BlockUnauthenticatedCert.ValueBool()
	var0.BlockUnknownCert = state.BlockUnknownCert.ValueBool()
	var var1 []cozuxBy.CaCertificatesObject
	if len(state.CaCertificates) != 0 {
		var1 = make([]cozuxBy.CaCertificatesObject, 0, len(state.CaCertificates))
		for var2Index := range state.CaCertificates {
			var2 := state.CaCertificates[var2Index]
			var var3 cozuxBy.CaCertificatesObject
			var3.DefaultOcspUrl = var2.DefaultOcspUrl.ValueString()
			var3.Name = var2.Name.ValueString()
			var3.OcspVerifyCert = var2.OcspVerifyCert.ValueString()
			var3.TemplateName = var2.TemplateName.ValueString()
			var1 = append(var1, var3)
		}
	}
	var0.CaCertificates = var1
	var0.CertStatusTimeout = state.CertStatusTimeout.ValueString()
	var0.CrlReceiveTimeout = state.CrlReceiveTimeout.ValueString()
	var0.Domain = state.Domain.ValueString()
	var0.Name = state.Name.ValueString()
	var0.OcspReceiveTimeout = state.OcspReceiveTimeout.ValueString()
	var0.UseCrl = state.UseCrl.ValueBool()
	var0.UseOcsp = state.UseOcsp.ValueBool()
	var var4 *cozuxBy.UsernameFieldObject
	if state.UsernameField != nil {
		var4 = &cozuxBy.UsernameFieldObject{}
		var4.Subject = state.UsernameField.Subject.ValueString()
		var4.SubjectAlt = state.UsernameField.SubjectAlt.ValueString()
	}
	var0.UsernameField = var4
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Create(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in create", err.Error())
		return
	}

	// Store the answer to state.
	state.Id = types.StringValue(strings.Join([]string{input.Folder, ans.ObjectId}, IdSeparator))
	var var5 []certificateProfilesRsModelCaCertificatesObject
	if len(ans.CaCertificates) != 0 {
		var5 = make([]certificateProfilesRsModelCaCertificatesObject, 0, len(ans.CaCertificates))
		for var6Index := range ans.CaCertificates {
			var6 := ans.CaCertificates[var6Index]
			var var7 certificateProfilesRsModelCaCertificatesObject
			var7.DefaultOcspUrl = types.StringValue(var6.DefaultOcspUrl)
			var7.Name = types.StringValue(var6.Name)
			var7.OcspVerifyCert = types.StringValue(var6.OcspVerifyCert)
			var7.TemplateName = types.StringValue(var6.TemplateName)
			var5 = append(var5, var7)
		}
	}
	var var8 *certificateProfilesRsModelUsernameFieldObject
	if ans.UsernameField != nil {
		var8 = &certificateProfilesRsModelUsernameFieldObject{}
		var8.Subject = types.StringValue(ans.UsernameField.Subject)
		var8.SubjectAlt = types.StringValue(ans.UsernameField.SubjectAlt)
	}
	state.BlockExpiredCert = types.BoolValue(ans.BlockExpiredCert)
	state.BlockTimeoutCert = types.BoolValue(ans.BlockTimeoutCert)
	state.BlockUnauthenticatedCert = types.BoolValue(ans.BlockUnauthenticatedCert)
	state.BlockUnknownCert = types.BoolValue(ans.BlockUnknownCert)
	state.CaCertificates = var5
	state.CertStatusTimeout = types.StringValue(ans.CertStatusTimeout)
	state.CrlReceiveTimeout = types.StringValue(ans.CrlReceiveTimeout)
	state.Domain = types.StringValue(ans.Domain)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.OcspReceiveTimeout = types.StringValue(ans.OcspReceiveTimeout)
	state.UseCrl = types.BoolValue(ans.UseCrl)
	state.UseOcsp = types.BoolValue(ans.UseOcsp)
	state.UsernameField = var8

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *certificateProfilesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state certificateProfilesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "sase_certificate_profiles",
		"locMap":        map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":        tokens,
	})

	// Prepare to read the config.
	svc := qLteaIq.NewClient(r.client)
	input := qLteaIq.ReadInput{
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
	var var0 []certificateProfilesRsModelCaCertificatesObject
	if len(ans.CaCertificates) != 0 {
		var0 = make([]certificateProfilesRsModelCaCertificatesObject, 0, len(ans.CaCertificates))
		for var1Index := range ans.CaCertificates {
			var1 := ans.CaCertificates[var1Index]
			var var2 certificateProfilesRsModelCaCertificatesObject
			var2.DefaultOcspUrl = types.StringValue(var1.DefaultOcspUrl)
			var2.Name = types.StringValue(var1.Name)
			var2.OcspVerifyCert = types.StringValue(var1.OcspVerifyCert)
			var2.TemplateName = types.StringValue(var1.TemplateName)
			var0 = append(var0, var2)
		}
	}
	var var3 *certificateProfilesRsModelUsernameFieldObject
	if ans.UsernameField != nil {
		var3 = &certificateProfilesRsModelUsernameFieldObject{}
		var3.Subject = types.StringValue(ans.UsernameField.Subject)
		var3.SubjectAlt = types.StringValue(ans.UsernameField.SubjectAlt)
	}
	state.BlockExpiredCert = types.BoolValue(ans.BlockExpiredCert)
	state.BlockTimeoutCert = types.BoolValue(ans.BlockTimeoutCert)
	state.BlockUnauthenticatedCert = types.BoolValue(ans.BlockUnauthenticatedCert)
	state.BlockUnknownCert = types.BoolValue(ans.BlockUnknownCert)
	state.CaCertificates = var0
	state.CertStatusTimeout = types.StringValue(ans.CertStatusTimeout)
	state.CrlReceiveTimeout = types.StringValue(ans.CrlReceiveTimeout)
	state.Domain = types.StringValue(ans.Domain)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.OcspReceiveTimeout = types.StringValue(ans.OcspReceiveTimeout)
	state.UseCrl = types.BoolValue(ans.UseCrl)
	state.UseOcsp = types.BoolValue(ans.UseOcsp)
	state.UsernameField = var3

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *certificateProfilesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state certificateProfilesRsModel
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
		"resource_name": "sase_certificate_profiles",
		"object_id":     state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := qLteaIq.NewClient(r.client)
	input := qLteaIq.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 cozuxBy.Config
	var0.BlockExpiredCert = plan.BlockExpiredCert.ValueBool()
	var0.BlockTimeoutCert = plan.BlockTimeoutCert.ValueBool()
	var0.BlockUnauthenticatedCert = plan.BlockUnauthenticatedCert.ValueBool()
	var0.BlockUnknownCert = plan.BlockUnknownCert.ValueBool()
	var var1 []cozuxBy.CaCertificatesObject
	if len(plan.CaCertificates) != 0 {
		var1 = make([]cozuxBy.CaCertificatesObject, 0, len(plan.CaCertificates))
		for var2Index := range plan.CaCertificates {
			var2 := plan.CaCertificates[var2Index]
			var var3 cozuxBy.CaCertificatesObject
			var3.DefaultOcspUrl = var2.DefaultOcspUrl.ValueString()
			var3.Name = var2.Name.ValueString()
			var3.OcspVerifyCert = var2.OcspVerifyCert.ValueString()
			var3.TemplateName = var2.TemplateName.ValueString()
			var1 = append(var1, var3)
		}
	}
	var0.CaCertificates = var1
	var0.CertStatusTimeout = plan.CertStatusTimeout.ValueString()
	var0.CrlReceiveTimeout = plan.CrlReceiveTimeout.ValueString()
	var0.Domain = plan.Domain.ValueString()
	var0.Name = plan.Name.ValueString()
	var0.OcspReceiveTimeout = plan.OcspReceiveTimeout.ValueString()
	var0.UseCrl = plan.UseCrl.ValueBool()
	var0.UseOcsp = plan.UseOcsp.ValueBool()
	var var4 *cozuxBy.UsernameFieldObject
	if plan.UsernameField != nil {
		var4 = &cozuxBy.UsernameFieldObject{}
		var4.Subject = plan.UsernameField.Subject.ValueString()
		var4.SubjectAlt = plan.UsernameField.SubjectAlt.ValueString()
	}
	var0.UsernameField = var4
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var5 []certificateProfilesRsModelCaCertificatesObject
	if len(ans.CaCertificates) != 0 {
		var5 = make([]certificateProfilesRsModelCaCertificatesObject, 0, len(ans.CaCertificates))
		for var6Index := range ans.CaCertificates {
			var6 := ans.CaCertificates[var6Index]
			var var7 certificateProfilesRsModelCaCertificatesObject
			var7.DefaultOcspUrl = types.StringValue(var6.DefaultOcspUrl)
			var7.Name = types.StringValue(var6.Name)
			var7.OcspVerifyCert = types.StringValue(var6.OcspVerifyCert)
			var7.TemplateName = types.StringValue(var6.TemplateName)
			var5 = append(var5, var7)
		}
	}
	var var8 *certificateProfilesRsModelUsernameFieldObject
	if ans.UsernameField != nil {
		var8 = &certificateProfilesRsModelUsernameFieldObject{}
		var8.Subject = types.StringValue(ans.UsernameField.Subject)
		var8.SubjectAlt = types.StringValue(ans.UsernameField.SubjectAlt)
	}
	state.BlockExpiredCert = types.BoolValue(ans.BlockExpiredCert)
	state.BlockTimeoutCert = types.BoolValue(ans.BlockTimeoutCert)
	state.BlockUnauthenticatedCert = types.BoolValue(ans.BlockUnauthenticatedCert)
	state.BlockUnknownCert = types.BoolValue(ans.BlockUnknownCert)
	state.CaCertificates = var5
	state.CertStatusTimeout = types.StringValue(ans.CertStatusTimeout)
	state.CrlReceiveTimeout = types.StringValue(ans.CrlReceiveTimeout)
	state.Domain = types.StringValue(ans.Domain)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.OcspReceiveTimeout = types.StringValue(ans.OcspReceiveTimeout)
	state.UseCrl = types.BoolValue(ans.UseCrl)
	state.UseOcsp = types.BoolValue(ans.UseOcsp)
	state.UsernameField = var8

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *certificateProfilesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name": "sase_certificate_profiles",
		"locMap":        map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":        tokens,
	})

	svc := qLteaIq.NewClient(r.client)
	input := qLteaIq.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *certificateProfilesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
