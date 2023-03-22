package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	rqJvTac "github.com/paloaltonetworks/sase-go/netsec/schema/saml/server/profiles"
	dZqhpfe "github.com/paloaltonetworks/sase-go/netsec/service/v1/samlserverprofiles"

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
	_ datasource.DataSource              = &samlServerProfilesListDataSource{}
	_ datasource.DataSourceWithConfigure = &samlServerProfilesListDataSource{}
)

func NewSamlServerProfilesListDataSource() datasource.DataSource {
	return &samlServerProfilesListDataSource{}
}

type samlServerProfilesListDataSource struct {
	client *sase.Client
}

type samlServerProfilesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Folder types.String `tfsdk:"folder"`
	Name   types.String `tfsdk:"name"`

	// Output.
	Data []samlServerProfilesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type samlServerProfilesListDsModelConfig struct {
	Certificate            types.String `tfsdk:"certificate"`
	EntityId               types.String `tfsdk:"entity_id"`
	ObjectId               types.String `tfsdk:"object_id"`
	MaxClockSkew           types.Int64  `tfsdk:"max_clock_skew"`
	SloBindings            types.String `tfsdk:"slo_bindings"`
	SsoBindings            types.String `tfsdk:"sso_bindings"`
	SsoUrl                 types.String `tfsdk:"sso_url"`
	ValidateIdpCertificate types.Bool   `tfsdk:"validate_idp_certificate"`
	WantAuthRequestsSigned types.Bool   `tfsdk:"want_auth_requests_signed"`
}

// Metadata returns the data source type name.
func (d *samlServerProfilesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_saml_server_profiles_list"
}

// Schema defines the schema for this listing data source.
func (d *samlServerProfilesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"certificate": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"entity_id": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"object_id": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"max_clock_skew": dsschema.Int64Attribute{
							Description: "",
							Computed:    true,
						},
						"slo_bindings": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"sso_bindings": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"sso_url": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"validate_idp_certificate": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"want_auth_requests_signed": dsschema.BoolAttribute{
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
func (d *samlServerProfilesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *samlServerProfilesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state samlServerProfilesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_saml_server_profiles_list",
		"terraform_provider_function": "Read",
		"limit":                       state.Limit.ValueInt64(),
		"has_limit":                   !state.Limit.IsNull(),
		"offset":                      state.Offset.ValueInt64(),
		"has_offset":                  !state.Offset.IsNull(),
		"folder":                      state.Folder.ValueString(),
		"name":                        state.Name.ValueString(),
		"has_name":                    !state.Name.IsNull(),
	})

	// Prepare to run the command.
	svc := dZqhpfe.NewClient(d.client)
	input := dZqhpfe.ListInput{
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
	var var0 []samlServerProfilesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]samlServerProfilesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 samlServerProfilesListDsModelConfig
			var2.Certificate = types.StringValue(var1.Certificate)
			var2.EntityId = types.StringValue(var1.EntityId)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.MaxClockSkew = types.Int64Value(var1.MaxClockSkew)
			var2.SloBindings = types.StringValue(var1.SloBindings)
			var2.SsoBindings = types.StringValue(var1.SsoBindings)
			var2.SsoUrl = types.StringValue(var1.SsoUrl)
			var2.ValidateIdpCertificate = types.BoolValue(var1.ValidateIdpCertificate)
			var2.WantAuthRequestsSigned = types.BoolValue(var1.WantAuthRequestsSigned)
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
	_ datasource.DataSource              = &samlServerProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &samlServerProfilesDataSource{}
)

func NewSamlServerProfilesDataSource() datasource.DataSource {
	return &samlServerProfilesDataSource{}
}

type samlServerProfilesDataSource struct {
	client *sase.Client
}

type samlServerProfilesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`
	Folder   types.String `tfsdk:"folder"`

	// Output.
	// Ref: #/components/schemas/saml-server-profiles
	Certificate types.String `tfsdk:"certificate"`
	EntityId    types.String `tfsdk:"entity_id"`
	// input omit: ObjectId
	MaxClockSkew           types.Int64  `tfsdk:"max_clock_skew"`
	SloBindings            types.String `tfsdk:"slo_bindings"`
	SsoBindings            types.String `tfsdk:"sso_bindings"`
	SsoUrl                 types.String `tfsdk:"sso_url"`
	ValidateIdpCertificate types.Bool   `tfsdk:"validate_idp_certificate"`
	WantAuthRequestsSigned types.Bool   `tfsdk:"want_auth_requests_signed"`
}

// Metadata returns the data source type name.
func (d *samlServerProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_saml_server_profiles"
}

// Schema defines the schema for this listing data source.
func (d *samlServerProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"certificate": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"entity_id": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"max_clock_skew": dsschema.Int64Attribute{
				Description: "",
				Computed:    true,
			},
			"slo_bindings": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"sso_bindings": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"sso_url": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"validate_idp_certificate": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"want_auth_requests_signed": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
		},
	}
}

// Configure prepares the struct.
func (d *samlServerProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *samlServerProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state samlServerProfilesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_saml_server_profiles",
		"object_id":                   state.ObjectId.ValueString(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := dZqhpfe.NewClient(d.client)
	input := dZqhpfe.ReadInput{
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
	state.Certificate = types.StringValue(ans.Certificate)
	state.EntityId = types.StringValue(ans.EntityId)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.MaxClockSkew = types.Int64Value(ans.MaxClockSkew)
	state.SloBindings = types.StringValue(ans.SloBindings)
	state.SsoBindings = types.StringValue(ans.SsoBindings)
	state.SsoUrl = types.StringValue(ans.SsoUrl)
	state.ValidateIdpCertificate = types.BoolValue(ans.ValidateIdpCertificate)
	state.WantAuthRequestsSigned = types.BoolValue(ans.WantAuthRequestsSigned)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &samlServerProfilesResource{}
	_ resource.ResourceWithConfigure   = &samlServerProfilesResource{}
	_ resource.ResourceWithImportState = &samlServerProfilesResource{}
)

func NewSamlServerProfilesResource() resource.Resource {
	return &samlServerProfilesResource{}
}

type samlServerProfilesResource struct {
	client *sase.Client
}

type samlServerProfilesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/saml-server-profiles
	Certificate            types.String `tfsdk:"certificate"`
	EntityId               types.String `tfsdk:"entity_id"`
	ObjectId               types.String `tfsdk:"object_id"`
	MaxClockSkew           types.Int64  `tfsdk:"max_clock_skew"`
	SloBindings            types.String `tfsdk:"slo_bindings"`
	SsoBindings            types.String `tfsdk:"sso_bindings"`
	SsoUrl                 types.String `tfsdk:"sso_url"`
	ValidateIdpCertificate types.Bool   `tfsdk:"validate_idp_certificate"`
	WantAuthRequestsSigned types.Bool   `tfsdk:"want_auth_requests_signed"`
}

// Metadata returns the data source type name.
func (r *samlServerProfilesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_saml_server_profiles"
}

// Schema defines the schema for this listing data source.
func (r *samlServerProfilesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
					stringvalidator.LengthAtMost(63),
				},
			},
			"entity_id": rsschema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 1024),
				},
			},
			"object_id": rsschema.StringAttribute{
				Description: "",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"max_clock_skew": rsschema.Int64Attribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					DefaultInt64(0),
				},
				Validators: []validator.Int64{
					int64validator.Between(1, 900),
				},
			},
			"slo_bindings": rsschema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("post", "redirect"),
				},
			},
			"sso_bindings": rsschema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("post", "redirect"),
				},
			},
			"sso_url": rsschema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"validate_idp_certificate": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"want_auth_requests_signed": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *samlServerProfilesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *samlServerProfilesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state samlServerProfilesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_saml_server_profiles",
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := dZqhpfe.NewClient(r.client)
	input := dZqhpfe.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 rqJvTac.Config
	var0.Certificate = state.Certificate.ValueString()
	var0.EntityId = state.EntityId.ValueString()
	var0.MaxClockSkew = state.MaxClockSkew.ValueInt64()
	var0.SloBindings = state.SloBindings.ValueString()
	var0.SsoBindings = state.SsoBindings.ValueString()
	var0.SsoUrl = state.SsoUrl.ValueString()
	var0.ValidateIdpCertificate = state.ValidateIdpCertificate.ValueBool()
	var0.WantAuthRequestsSigned = state.WantAuthRequestsSigned.ValueBool()
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
	state.Certificate = types.StringValue(ans.Certificate)
	state.EntityId = types.StringValue(ans.EntityId)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.MaxClockSkew = types.Int64Value(ans.MaxClockSkew)
	state.SloBindings = types.StringValue(ans.SloBindings)
	state.SsoBindings = types.StringValue(ans.SsoBindings)
	state.SsoUrl = types.StringValue(ans.SsoUrl)
	state.ValidateIdpCertificate = types.BoolValue(ans.ValidateIdpCertificate)
	state.WantAuthRequestsSigned = types.BoolValue(ans.WantAuthRequestsSigned)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *samlServerProfilesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state samlServerProfilesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_saml_server_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := dZqhpfe.NewClient(r.client)
	input := dZqhpfe.ReadInput{
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
	state.Certificate = types.StringValue(ans.Certificate)
	state.EntityId = types.StringValue(ans.EntityId)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.MaxClockSkew = types.Int64Value(ans.MaxClockSkew)
	state.SloBindings = types.StringValue(ans.SloBindings)
	state.SsoBindings = types.StringValue(ans.SsoBindings)
	state.SsoUrl = types.StringValue(ans.SsoUrl)
	state.ValidateIdpCertificate = types.BoolValue(ans.ValidateIdpCertificate)
	state.WantAuthRequestsSigned = types.BoolValue(ans.WantAuthRequestsSigned)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *samlServerProfilesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state samlServerProfilesRsModel
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
		"resource_name":               "sase_saml_server_profiles",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := dZqhpfe.NewClient(r.client)
	input := dZqhpfe.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 rqJvTac.Config
	var0.Certificate = plan.Certificate.ValueString()
	var0.EntityId = plan.EntityId.ValueString()
	var0.MaxClockSkew = plan.MaxClockSkew.ValueInt64()
	var0.SloBindings = plan.SloBindings.ValueString()
	var0.SsoBindings = plan.SsoBindings.ValueString()
	var0.SsoUrl = plan.SsoUrl.ValueString()
	var0.ValidateIdpCertificate = plan.ValidateIdpCertificate.ValueBool()
	var0.WantAuthRequestsSigned = plan.WantAuthRequestsSigned.ValueBool()
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	state.Certificate = types.StringValue(ans.Certificate)
	state.EntityId = types.StringValue(ans.EntityId)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.MaxClockSkew = types.Int64Value(ans.MaxClockSkew)
	state.SloBindings = types.StringValue(ans.SloBindings)
	state.SsoBindings = types.StringValue(ans.SsoBindings)
	state.SsoUrl = types.StringValue(ans.SsoUrl)
	state.ValidateIdpCertificate = types.BoolValue(ans.ValidateIdpCertificate)
	state.WantAuthRequestsSigned = types.BoolValue(ans.WantAuthRequestsSigned)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *samlServerProfilesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name":               "sase_saml_server_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	svc := dZqhpfe.NewClient(r.client)
	input := dZqhpfe.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *samlServerProfilesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
