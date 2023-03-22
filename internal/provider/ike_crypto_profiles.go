package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	gKJgxWC "github.com/paloaltonetworks/sase-go/netsec/schema/ike/crypto/profiles"
	aZqXHLP "github.com/paloaltonetworks/sase-go/netsec/service/v1/ikecryptoprofiles"

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
	_ datasource.DataSource              = &ikeCryptoProfilesListDataSource{}
	_ datasource.DataSourceWithConfigure = &ikeCryptoProfilesListDataSource{}
)

func NewIkeCryptoProfilesListDataSource() datasource.DataSource {
	return &ikeCryptoProfilesListDataSource{}
}

type ikeCryptoProfilesListDataSource struct {
	client *sase.Client
}

type ikeCryptoProfilesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []ikeCryptoProfilesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type ikeCryptoProfilesListDsModelConfig struct {
	AuthenticationMultiple types.Int64                                 `tfsdk:"authentication_multiple"`
	DhGroup                []types.String                              `tfsdk:"dh_group"`
	Encryption             []types.String                              `tfsdk:"encryption"`
	Hash                   []types.String                              `tfsdk:"hash"`
	ObjectId               types.String                                `tfsdk:"object_id"`
	Lifetime               *ikeCryptoProfilesListDsModelLifetimeObject `tfsdk:"lifetime"`
	Name                   types.String                                `tfsdk:"name"`
}

type ikeCryptoProfilesListDsModelLifetimeObject struct {
	Days    types.Int64 `tfsdk:"days"`
	Hours   types.Int64 `tfsdk:"hours"`
	Minutes types.Int64 `tfsdk:"minutes"`
	Seconds types.Int64 `tfsdk:"seconds"`
}

// Metadata returns the data source type name.
func (d *ikeCryptoProfilesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ike_crypto_profiles_list"
}

// Schema defines the schema for this listing data source.
func (d *ikeCryptoProfilesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"authentication_multiple": dsschema.Int64Attribute{
							Description: "",
							Computed:    true,
						},
						"dh_group": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"encryption": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"hash": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"object_id": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"lifetime": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"days": dsschema.Int64Attribute{
									Description: "",
									Computed:    true,
								},
								"hours": dsschema.Int64Attribute{
									Description: "",
									Computed:    true,
								},
								"minutes": dsschema.Int64Attribute{
									Description: "",
									Computed:    true,
								},
								"seconds": dsschema.Int64Attribute{
									Description: "",
									Computed:    true,
								},
							},
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
func (d *ikeCryptoProfilesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *ikeCryptoProfilesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ikeCryptoProfilesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_ike_crypto_profiles_list",
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
	svc := aZqXHLP.NewClient(d.client)
	input := aZqXHLP.ListInput{
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
	var var0 []ikeCryptoProfilesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]ikeCryptoProfilesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 ikeCryptoProfilesListDsModelConfig
			var var3 *ikeCryptoProfilesListDsModelLifetimeObject
			if var1.Lifetime != nil {
				var3 = &ikeCryptoProfilesListDsModelLifetimeObject{}
				var3.Days = types.Int64Value(var1.Lifetime.Days)
				var3.Hours = types.Int64Value(var1.Lifetime.Hours)
				var3.Minutes = types.Int64Value(var1.Lifetime.Minutes)
				var3.Seconds = types.Int64Value(var1.Lifetime.Seconds)
			}
			var2.AuthenticationMultiple = types.Int64Value(var1.AuthenticationMultiple)
			var2.DhGroup = EncodeStringSlice(var1.DhGroup)
			var2.Encryption = EncodeStringSlice(var1.Encryption)
			var2.Hash = EncodeStringSlice(var1.Hash)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Lifetime = var3
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
	_ datasource.DataSource              = &ikeCryptoProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &ikeCryptoProfilesDataSource{}
)

func NewIkeCryptoProfilesDataSource() datasource.DataSource {
	return &ikeCryptoProfilesDataSource{}
}

type ikeCryptoProfilesDataSource struct {
	client *sase.Client
}

type ikeCryptoProfilesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`
	Folder   types.String `tfsdk:"folder"`

	// Output.
	// Ref: #/components/schemas/ike-crypto-profiles
	AuthenticationMultiple types.Int64    `tfsdk:"authentication_multiple"`
	DhGroup                []types.String `tfsdk:"dh_group"`
	Encryption             []types.String `tfsdk:"encryption"`
	Hash                   []types.String `tfsdk:"hash"`
	// input omit: ObjectId
	Lifetime *ikeCryptoProfilesDsModelLifetimeObject `tfsdk:"lifetime"`
	Name     types.String                            `tfsdk:"name"`
}

type ikeCryptoProfilesDsModelLifetimeObject struct {
	Days    types.Int64 `tfsdk:"days"`
	Hours   types.Int64 `tfsdk:"hours"`
	Minutes types.Int64 `tfsdk:"minutes"`
	Seconds types.Int64 `tfsdk:"seconds"`
}

// Metadata returns the data source type name.
func (d *ikeCryptoProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ike_crypto_profiles"
}

// Schema defines the schema for this listing data source.
func (d *ikeCryptoProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"authentication_multiple": dsschema.Int64Attribute{
				Description: "",
				Computed:    true,
			},
			"dh_group": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"encryption": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"hash": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"lifetime": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"days": dsschema.Int64Attribute{
						Description: "",
						Computed:    true,
					},
					"hours": dsschema.Int64Attribute{
						Description: "",
						Computed:    true,
					},
					"minutes": dsschema.Int64Attribute{
						Description: "",
						Computed:    true,
					},
					"seconds": dsschema.Int64Attribute{
						Description: "",
						Computed:    true,
					},
				},
			},
			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
		},
	}
}

// Configure prepares the struct.
func (d *ikeCryptoProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *ikeCryptoProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ikeCryptoProfilesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_ike_crypto_profiles",
		"object_id":                   state.ObjectId.ValueString(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := aZqXHLP.NewClient(d.client)
	input := aZqXHLP.ReadInput{
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
	var var0 *ikeCryptoProfilesDsModelLifetimeObject
	if ans.Lifetime != nil {
		var0 = &ikeCryptoProfilesDsModelLifetimeObject{}
		var0.Days = types.Int64Value(ans.Lifetime.Days)
		var0.Hours = types.Int64Value(ans.Lifetime.Hours)
		var0.Minutes = types.Int64Value(ans.Lifetime.Minutes)
		var0.Seconds = types.Int64Value(ans.Lifetime.Seconds)
	}
	state.AuthenticationMultiple = types.Int64Value(ans.AuthenticationMultiple)
	state.DhGroup = EncodeStringSlice(ans.DhGroup)
	state.Encryption = EncodeStringSlice(ans.Encryption)
	state.Hash = EncodeStringSlice(ans.Hash)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Lifetime = var0
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &ikeCryptoProfilesResource{}
	_ resource.ResourceWithConfigure   = &ikeCryptoProfilesResource{}
	_ resource.ResourceWithImportState = &ikeCryptoProfilesResource{}
)

func NewIkeCryptoProfilesResource() resource.Resource {
	return &ikeCryptoProfilesResource{}
}

type ikeCryptoProfilesResource struct {
	client *sase.Client
}

type ikeCryptoProfilesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/ike-crypto-profiles
	AuthenticationMultiple types.Int64                             `tfsdk:"authentication_multiple"`
	DhGroup                []types.String                          `tfsdk:"dh_group"`
	Encryption             []types.String                          `tfsdk:"encryption"`
	Hash                   []types.String                          `tfsdk:"hash"`
	ObjectId               types.String                            `tfsdk:"object_id"`
	Lifetime               *ikeCryptoProfilesRsModelLifetimeObject `tfsdk:"lifetime"`
	Name                   types.String                            `tfsdk:"name"`
}

type ikeCryptoProfilesRsModelLifetimeObject struct {
	Days    types.Int64 `tfsdk:"days"`
	Hours   types.Int64 `tfsdk:"hours"`
	Minutes types.Int64 `tfsdk:"minutes"`
	Seconds types.Int64 `tfsdk:"seconds"`
}

// Metadata returns the data source type name.
func (r *ikeCryptoProfilesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ike_crypto_profiles"
}

// Schema defines the schema for this listing data source.
func (r *ikeCryptoProfilesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

			"authentication_multiple": rsschema.Int64Attribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					DefaultInt64(0),
				},
				Validators: []validator.Int64{
					int64validator.AtMost(50),
				},
			},
			"dh_group": rsschema.ListAttribute{
				Description: "",
				Required:    true,
				ElementType: types.StringType,
			},
			"encryption": rsschema.ListAttribute{
				Description: "",
				Required:    true,
				ElementType: types.StringType,
			},
			"hash": rsschema.ListAttribute{
				Description: "",
				Required:    true,
				ElementType: types.StringType,
			},
			"object_id": rsschema.StringAttribute{
				Description: "",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"lifetime": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"days": rsschema.Int64Attribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							DefaultInt64(0),
						},
						Validators: []validator.Int64{
							int64validator.Between(1, 365),
						},
					},
					"hours": rsschema.Int64Attribute{
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
					"minutes": rsschema.Int64Attribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							DefaultInt64(0),
						},
						Validators: []validator.Int64{
							int64validator.Between(3, 65535),
						},
					},
					"seconds": rsschema.Int64Attribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							DefaultInt64(0),
						},
						Validators: []validator.Int64{
							int64validator.Between(180, 65535),
						},
					},
				},
			},
			"name": rsschema.StringAttribute{
				Description: "",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(31),
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *ikeCryptoProfilesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *ikeCryptoProfilesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state ikeCryptoProfilesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_ike_crypto_profiles",
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := aZqXHLP.NewClient(r.client)
	input := aZqXHLP.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 gKJgxWC.Config
	var0.AuthenticationMultiple = state.AuthenticationMultiple.ValueInt64()
	var0.DhGroup = DecodeStringSlice(state.DhGroup)
	var0.Encryption = DecodeStringSlice(state.Encryption)
	var0.Hash = DecodeStringSlice(state.Hash)
	var var1 *gKJgxWC.LifetimeObject
	if state.Lifetime != nil {
		var1 = &gKJgxWC.LifetimeObject{}
		var1.Days = state.Lifetime.Days.ValueInt64()
		var1.Hours = state.Lifetime.Hours.ValueInt64()
		var1.Minutes = state.Lifetime.Minutes.ValueInt64()
		var1.Seconds = state.Lifetime.Seconds.ValueInt64()
	}
	var0.Lifetime = var1
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
	var var2 *ikeCryptoProfilesRsModelLifetimeObject
	if ans.Lifetime != nil {
		var2 = &ikeCryptoProfilesRsModelLifetimeObject{}
		var2.Days = types.Int64Value(ans.Lifetime.Days)
		var2.Hours = types.Int64Value(ans.Lifetime.Hours)
		var2.Minutes = types.Int64Value(ans.Lifetime.Minutes)
		var2.Seconds = types.Int64Value(ans.Lifetime.Seconds)
	}
	state.AuthenticationMultiple = types.Int64Value(ans.AuthenticationMultiple)
	state.DhGroup = EncodeStringSlice(ans.DhGroup)
	state.Encryption = EncodeStringSlice(ans.Encryption)
	state.Hash = EncodeStringSlice(ans.Hash)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Lifetime = var2
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *ikeCryptoProfilesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state ikeCryptoProfilesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_ike_crypto_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := aZqXHLP.NewClient(r.client)
	input := aZqXHLP.ReadInput{
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
	var var0 *ikeCryptoProfilesRsModelLifetimeObject
	if ans.Lifetime != nil {
		var0 = &ikeCryptoProfilesRsModelLifetimeObject{}
		var0.Days = types.Int64Value(ans.Lifetime.Days)
		var0.Hours = types.Int64Value(ans.Lifetime.Hours)
		var0.Minutes = types.Int64Value(ans.Lifetime.Minutes)
		var0.Seconds = types.Int64Value(ans.Lifetime.Seconds)
	}
	state.AuthenticationMultiple = types.Int64Value(ans.AuthenticationMultiple)
	state.DhGroup = EncodeStringSlice(ans.DhGroup)
	state.Encryption = EncodeStringSlice(ans.Encryption)
	state.Hash = EncodeStringSlice(ans.Hash)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Lifetime = var0
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *ikeCryptoProfilesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ikeCryptoProfilesRsModel
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
		"resource_name":               "sase_ike_crypto_profiles",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := aZqXHLP.NewClient(r.client)
	input := aZqXHLP.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 gKJgxWC.Config
	var0.AuthenticationMultiple = plan.AuthenticationMultiple.ValueInt64()
	var0.DhGroup = DecodeStringSlice(plan.DhGroup)
	var0.Encryption = DecodeStringSlice(plan.Encryption)
	var0.Hash = DecodeStringSlice(plan.Hash)
	var var1 *gKJgxWC.LifetimeObject
	if plan.Lifetime != nil {
		var1 = &gKJgxWC.LifetimeObject{}
		var1.Days = plan.Lifetime.Days.ValueInt64()
		var1.Hours = plan.Lifetime.Hours.ValueInt64()
		var1.Minutes = plan.Lifetime.Minutes.ValueInt64()
		var1.Seconds = plan.Lifetime.Seconds.ValueInt64()
	}
	var0.Lifetime = var1
	var0.Name = plan.Name.ValueString()
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var2 *ikeCryptoProfilesRsModelLifetimeObject
	if ans.Lifetime != nil {
		var2 = &ikeCryptoProfilesRsModelLifetimeObject{}
		var2.Days = types.Int64Value(ans.Lifetime.Days)
		var2.Hours = types.Int64Value(ans.Lifetime.Hours)
		var2.Minutes = types.Int64Value(ans.Lifetime.Minutes)
		var2.Seconds = types.Int64Value(ans.Lifetime.Seconds)
	}
	state.AuthenticationMultiple = types.Int64Value(ans.AuthenticationMultiple)
	state.DhGroup = EncodeStringSlice(ans.DhGroup)
	state.Encryption = EncodeStringSlice(ans.Encryption)
	state.Hash = EncodeStringSlice(ans.Hash)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Lifetime = var2
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *ikeCryptoProfilesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name":               "sase_ike_crypto_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	svc := aZqXHLP.NewClient(r.client)
	input := aZqXHLP.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *ikeCryptoProfilesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
