package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	fWrszss "github.com/paloaltonetworks/sase-go/netsec/schema/tacacs/server/profiles"
	lUnrbOf "github.com/paloaltonetworks/sase-go/netsec/service/v1/tacacsserverprofiles"

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
	_ datasource.DataSource              = &tacacsServerProfilesListDataSource{}
	_ datasource.DataSourceWithConfigure = &tacacsServerProfilesListDataSource{}
)

func NewTacacsServerProfilesListDataSource() datasource.DataSource {
	return &tacacsServerProfilesListDataSource{}
}

type tacacsServerProfilesListDataSource struct {
	client *sase.Client
}

type tacacsServerProfilesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Folder types.String `tfsdk:"folder"`
	Name   types.String `tfsdk:"name"`

	// Output.
	Data []tacacsServerProfilesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type tacacsServerProfilesListDsModelConfig struct {
	ObjectId            types.String                                  `tfsdk:"object_id"`
	Protocol            types.String                                  `tfsdk:"protocol"`
	Server              []tacacsServerProfilesListDsModelServerObject `tfsdk:"server"`
	Timeout             types.Int64                                   `tfsdk:"timeout"`
	UseSingleConnection types.Bool                                    `tfsdk:"use_single_connection"`
}

type tacacsServerProfilesListDsModelServerObject struct {
	Address types.String `tfsdk:"address"`
	Name    types.String `tfsdk:"name"`
	Port    types.Int64  `tfsdk:"port"`
	Secret  types.String `tfsdk:"secret"`
}

// Metadata returns the data source type name.
func (d *tacacsServerProfilesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tacacs_server_profiles_list"
}

// Schema defines the schema for this listing data source.
func (d *tacacsServerProfilesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"folder": dsschema.StringAttribute{
				Description:         "The folder of the entry",
				MarkdownDescription: "The folder of the entry",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Service Connections", "Mobile Users Container", "Mobile Users Explicit Proxy"),
				},
			},
			"name": dsschema.StringAttribute{
				Description:         "The name of the entry",
				MarkdownDescription: "The name of the entry",
				Optional:            true,
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
						"protocol": dsschema.StringAttribute{
							Description:         "The `protocol` parameter.",
							MarkdownDescription: "The `protocol` parameter.",
							Computed:            true,
						},
						"server": dsschema.ListNestedAttribute{
							Description:         "The `server` parameter.",
							MarkdownDescription: "The `server` parameter.",
							Computed:            true,
							NestedObject: dsschema.NestedAttributeObject{
								Attributes: map[string]dsschema.Attribute{
									"address": dsschema.StringAttribute{
										Description:         "The `address` parameter.",
										MarkdownDescription: "The `address` parameter.",
										Computed:            true,
									},
									"name": dsschema.StringAttribute{
										Description:         "The `name` parameter.",
										MarkdownDescription: "The `name` parameter.",
										Computed:            true,
									},
									"port": dsschema.Int64Attribute{
										Description:         "The `port` parameter.",
										MarkdownDescription: "The `port` parameter.",
										Computed:            true,
									},
									"secret": dsschema.StringAttribute{
										Description:         "The `secret` parameter.",
										MarkdownDescription: "The `secret` parameter.",
										Computed:            true,
									},
								},
							},
						},
						"timeout": dsschema.Int64Attribute{
							Description:         "The `timeout` parameter.",
							MarkdownDescription: "The `timeout` parameter.",
							Computed:            true,
						},
						"use_single_connection": dsschema.BoolAttribute{
							Description:         "The `use_single_connection` parameter.",
							MarkdownDescription: "The `use_single_connection` parameter.",
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
func (d *tacacsServerProfilesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *tacacsServerProfilesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state tacacsServerProfilesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_tacacs_server_profiles_list",
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
	svc := lUnrbOf.NewClient(d.client)
	input := lUnrbOf.ListInput{
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
	var var0 []tacacsServerProfilesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]tacacsServerProfilesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 tacacsServerProfilesListDsModelConfig
			var var3 []tacacsServerProfilesListDsModelServerObject
			if len(var1.Server) != 0 {
				var3 = make([]tacacsServerProfilesListDsModelServerObject, 0, len(var1.Server))
				for var4Index := range var1.Server {
					var4 := var1.Server[var4Index]
					var var5 tacacsServerProfilesListDsModelServerObject
					var5.Address = types.StringValue(var4.Address)
					var5.Name = types.StringValue(var4.Name)
					var5.Port = types.Int64Value(var4.Port)
					var5.Secret = types.StringValue(var4.Secret)
					var3 = append(var3, var5)
				}
			}
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Protocol = types.StringValue(var1.Protocol)
			var2.Server = var3
			var2.Timeout = types.Int64Value(var1.Timeout)
			var2.UseSingleConnection = types.BoolValue(var1.UseSingleConnection)
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
	_ datasource.DataSource              = &tacacsServerProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &tacacsServerProfilesDataSource{}
)

func NewTacacsServerProfilesDataSource() datasource.DataSource {
	return &tacacsServerProfilesDataSource{}
}

type tacacsServerProfilesDataSource struct {
	client *sase.Client
}

type tacacsServerProfilesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`
	Folder   types.String `tfsdk:"folder"`

	// Output.
	// Ref: #/components/schemas/tacacs-server-profiles
	// input omit: ObjectId
	Protocol            types.String                              `tfsdk:"protocol"`
	Server              []tacacsServerProfilesDsModelServerObject `tfsdk:"server"`
	Timeout             types.Int64                               `tfsdk:"timeout"`
	UseSingleConnection types.Bool                                `tfsdk:"use_single_connection"`
}

type tacacsServerProfilesDsModelServerObject struct {
	Address types.String `tfsdk:"address"`
	Name    types.String `tfsdk:"name"`
	Port    types.Int64  `tfsdk:"port"`
	Secret  types.String `tfsdk:"secret"`
}

// Metadata returns the data source type name.
func (d *tacacsServerProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tacacs_server_profiles"
}

// Schema defines the schema for this listing data source.
func (d *tacacsServerProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"folder": dsschema.StringAttribute{
				Description:         "The folder of the entry",
				MarkdownDescription: "The folder of the entry",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Service Connections", "Mobile Users Container", "Mobile Users Explicit Proxy"),
				},
			},

			// Output.
			"protocol": dsschema.StringAttribute{
				Description:         "The `protocol` parameter.",
				MarkdownDescription: "The `protocol` parameter.",
				Computed:            true,
			},
			"server": dsschema.ListNestedAttribute{
				Description:         "The `server` parameter.",
				MarkdownDescription: "The `server` parameter.",
				Computed:            true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"address": dsschema.StringAttribute{
							Description:         "The `address` parameter.",
							MarkdownDescription: "The `address` parameter.",
							Computed:            true,
						},
						"name": dsschema.StringAttribute{
							Description:         "The `name` parameter.",
							MarkdownDescription: "The `name` parameter.",
							Computed:            true,
						},
						"port": dsschema.Int64Attribute{
							Description:         "The `port` parameter.",
							MarkdownDescription: "The `port` parameter.",
							Computed:            true,
						},
						"secret": dsschema.StringAttribute{
							Description:         "The `secret` parameter.",
							MarkdownDescription: "The `secret` parameter.",
							Computed:            true,
						},
					},
				},
			},
			"timeout": dsschema.Int64Attribute{
				Description:         "The `timeout` parameter.",
				MarkdownDescription: "The `timeout` parameter.",
				Computed:            true,
			},
			"use_single_connection": dsschema.BoolAttribute{
				Description:         "The `use_single_connection` parameter.",
				MarkdownDescription: "The `use_single_connection` parameter.",
				Computed:            true,
			},
		},
	}
}

// Configure prepares the struct.
func (d *tacacsServerProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *tacacsServerProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state tacacsServerProfilesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_tacacs_server_profiles",
		"object_id":                   state.ObjectId.ValueString(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := lUnrbOf.NewClient(d.client)
	input := lUnrbOf.ReadInput{
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
	var var0 []tacacsServerProfilesDsModelServerObject
	if len(ans.Server) != 0 {
		var0 = make([]tacacsServerProfilesDsModelServerObject, 0, len(ans.Server))
		for var1Index := range ans.Server {
			var1 := ans.Server[var1Index]
			var var2 tacacsServerProfilesDsModelServerObject
			var2.Address = types.StringValue(var1.Address)
			var2.Name = types.StringValue(var1.Name)
			var2.Port = types.Int64Value(var1.Port)
			var2.Secret = types.StringValue(var1.Secret)
			var0 = append(var0, var2)
		}
	}
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Protocol = types.StringValue(ans.Protocol)
	state.Server = var0
	state.Timeout = types.Int64Value(ans.Timeout)
	state.UseSingleConnection = types.BoolValue(ans.UseSingleConnection)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &tacacsServerProfilesResource{}
	_ resource.ResourceWithConfigure   = &tacacsServerProfilesResource{}
	_ resource.ResourceWithImportState = &tacacsServerProfilesResource{}
)

func NewTacacsServerProfilesResource() resource.Resource {
	return &tacacsServerProfilesResource{}
}

type tacacsServerProfilesResource struct {
	client *sase.Client
}

type tacacsServerProfilesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/tacacs-server-profiles
	ObjectId            types.String                              `tfsdk:"object_id"`
	Protocol            types.String                              `tfsdk:"protocol"`
	Server              []tacacsServerProfilesRsModelServerObject `tfsdk:"server"`
	Timeout             types.Int64                               `tfsdk:"timeout"`
	UseSingleConnection types.Bool                                `tfsdk:"use_single_connection"`
}

type tacacsServerProfilesRsModelServerObject struct {
	Address types.String `tfsdk:"address"`
	Name    types.String `tfsdk:"name"`
	Port    types.Int64  `tfsdk:"port"`
	Secret  types.String `tfsdk:"secret"`
}

// Metadata returns the data source type name.
func (r *tacacsServerProfilesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tacacs_server_profiles"
}

// Schema defines the schema for this listing data source.
func (r *tacacsServerProfilesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Description:         "The folder of the entry",
				MarkdownDescription: "The folder of the entry",
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
			"protocol": rsschema.StringAttribute{
				Description:         "The `protocol` parameter.",
				MarkdownDescription: "The `protocol` parameter.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("CHAP", "PAP"),
				},
			},
			"server": rsschema.ListNestedAttribute{
				Description:         "The `server` parameter.",
				MarkdownDescription: "The `server` parameter.",
				Required:            true,
				NestedObject: rsschema.NestedAttributeObject{
					Attributes: map[string]rsschema.Attribute{
						"address": rsschema.StringAttribute{
							Description:         "The `address` parameter.",
							MarkdownDescription: "The `address` parameter.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
						},
						"name": rsschema.StringAttribute{
							Description:         "The `name` parameter.",
							MarkdownDescription: "The `name` parameter.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
						},
						"port": rsschema.Int64Attribute{
							Description:         "The `port` parameter.",
							MarkdownDescription: "The `port` parameter.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Int64{
								DefaultInt64(0),
							},
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
						"secret": rsschema.StringAttribute{
							Description:         "The `secret` parameter.",
							MarkdownDescription: "The `secret` parameter.",
							Optional:            true,
							Computed:            true,
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
				Description:         "The `timeout` parameter.",
				MarkdownDescription: "The `timeout` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					DefaultInt64(0),
				},
				Validators: []validator.Int64{
					int64validator.Between(1, 30),
				},
			},
			"use_single_connection": rsschema.BoolAttribute{
				Description:         "The `use_single_connection` parameter.",
				MarkdownDescription: "The `use_single_connection` parameter.",
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
func (r *tacacsServerProfilesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *tacacsServerProfilesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state tacacsServerProfilesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_tacacs_server_profiles",
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := lUnrbOf.NewClient(r.client)
	input := lUnrbOf.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 fWrszss.Config
	var0.Protocol = state.Protocol.ValueString()
	var var1 []fWrszss.ServerObject
	if len(state.Server) != 0 {
		var1 = make([]fWrszss.ServerObject, 0, len(state.Server))
		for var2Index := range state.Server {
			var2 := state.Server[var2Index]
			var var3 fWrszss.ServerObject
			var3.Address = var2.Address.ValueString()
			var3.Name = var2.Name.ValueString()
			var3.Port = var2.Port.ValueInt64()
			var3.Secret = var2.Secret.ValueString()
			var1 = append(var1, var3)
		}
	}
	var0.Server = var1
	var0.Timeout = state.Timeout.ValueInt64()
	var0.UseSingleConnection = state.UseSingleConnection.ValueBool()
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
	var var4 []tacacsServerProfilesRsModelServerObject
	if len(ans.Server) != 0 {
		var4 = make([]tacacsServerProfilesRsModelServerObject, 0, len(ans.Server))
		for var5Index := range ans.Server {
			var5 := ans.Server[var5Index]
			var var6 tacacsServerProfilesRsModelServerObject
			var6.Address = types.StringValue(var5.Address)
			var6.Name = types.StringValue(var5.Name)
			var6.Port = types.Int64Value(var5.Port)
			var6.Secret = types.StringValue(var5.Secret)
			var4 = append(var4, var6)
		}
	}
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Protocol = types.StringValue(ans.Protocol)
	state.Server = var4
	state.Timeout = types.Int64Value(ans.Timeout)
	state.UseSingleConnection = types.BoolValue(ans.UseSingleConnection)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *tacacsServerProfilesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state tacacsServerProfilesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_tacacs_server_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := lUnrbOf.NewClient(r.client)
	input := lUnrbOf.ReadInput{
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
	var var0 []tacacsServerProfilesRsModelServerObject
	if len(ans.Server) != 0 {
		var0 = make([]tacacsServerProfilesRsModelServerObject, 0, len(ans.Server))
		for var1Index := range ans.Server {
			var1 := ans.Server[var1Index]
			var var2 tacacsServerProfilesRsModelServerObject
			var2.Address = types.StringValue(var1.Address)
			var2.Name = types.StringValue(var1.Name)
			var2.Port = types.Int64Value(var1.Port)
			var2.Secret = types.StringValue(var1.Secret)
			var0 = append(var0, var2)
		}
	}
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Protocol = types.StringValue(ans.Protocol)
	state.Server = var0
	state.Timeout = types.Int64Value(ans.Timeout)
	state.UseSingleConnection = types.BoolValue(ans.UseSingleConnection)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *tacacsServerProfilesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state tacacsServerProfilesRsModel
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
		"resource_name":               "sase_tacacs_server_profiles",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := lUnrbOf.NewClient(r.client)
	input := lUnrbOf.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 fWrszss.Config
	var0.Protocol = plan.Protocol.ValueString()
	var var1 []fWrszss.ServerObject
	if len(plan.Server) != 0 {
		var1 = make([]fWrszss.ServerObject, 0, len(plan.Server))
		for var2Index := range plan.Server {
			var2 := plan.Server[var2Index]
			var var3 fWrszss.ServerObject
			var3.Address = var2.Address.ValueString()
			var3.Name = var2.Name.ValueString()
			var3.Port = var2.Port.ValueInt64()
			var3.Secret = var2.Secret.ValueString()
			var1 = append(var1, var3)
		}
	}
	var0.Server = var1
	var0.Timeout = plan.Timeout.ValueInt64()
	var0.UseSingleConnection = plan.UseSingleConnection.ValueBool()
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var4 []tacacsServerProfilesRsModelServerObject
	if len(ans.Server) != 0 {
		var4 = make([]tacacsServerProfilesRsModelServerObject, 0, len(ans.Server))
		for var5Index := range ans.Server {
			var5 := ans.Server[var5Index]
			var var6 tacacsServerProfilesRsModelServerObject
			var6.Address = types.StringValue(var5.Address)
			var6.Name = types.StringValue(var5.Name)
			var6.Port = types.Int64Value(var5.Port)
			var6.Secret = types.StringValue(var5.Secret)
			var4 = append(var4, var6)
		}
	}
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Protocol = types.StringValue(ans.Protocol)
	state.Server = var4
	state.Timeout = types.Int64Value(ans.Timeout)
	state.UseSingleConnection = types.BoolValue(ans.UseSingleConnection)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *tacacsServerProfilesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name":               "sase_tacacs_server_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	svc := lUnrbOf.NewClient(r.client)
	input := lUnrbOf.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *tacacsServerProfilesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
