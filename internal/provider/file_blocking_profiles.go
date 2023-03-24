package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	wMgZmmI "github.com/paloaltonetworks/sase-go/netsec/schema/file/blocking/profiles"
	fEpWCgc "github.com/paloaltonetworks/sase-go/netsec/service/v1/fileblockingprofiles"

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
	_ datasource.DataSource              = &fileBlockingProfilesListDataSource{}
	_ datasource.DataSourceWithConfigure = &fileBlockingProfilesListDataSource{}
)

func NewFileBlockingProfilesListDataSource() datasource.DataSource {
	return &fileBlockingProfilesListDataSource{}
}

type fileBlockingProfilesListDataSource struct {
	client *sase.Client
}

type fileBlockingProfilesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []fileBlockingProfilesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type fileBlockingProfilesListDsModelConfig struct {
	Description types.String                                 `tfsdk:"description"`
	ObjectId    types.String                                 `tfsdk:"object_id"`
	Name        types.String                                 `tfsdk:"name"`
	Rules       []fileBlockingProfilesListDsModelRulesObject `tfsdk:"rules"`
}

type fileBlockingProfilesListDsModelRulesObject struct {
	Action      types.String   `tfsdk:"action"`
	Application []types.String `tfsdk:"application"`
	Direction   types.String   `tfsdk:"direction"`
	FileType    []types.String `tfsdk:"file_type"`
	Name        types.String   `tfsdk:"name"`
}

// Metadata returns the data source type name.
func (d *fileBlockingProfilesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file_blocking_profiles_list"
}

// Schema defines the schema for this listing data source.
func (d *fileBlockingProfilesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"description": dsschema.StringAttribute{
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
						"rules": dsschema.ListNestedAttribute{
							Description: "",
							Computed:    true,
							NestedObject: dsschema.NestedAttributeObject{
								Attributes: map[string]dsschema.Attribute{
									"action": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"application": dsschema.ListAttribute{
										Description: "",
										Computed:    true,
										ElementType: types.StringType,
									},
									"direction": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"file_type": dsschema.ListAttribute{
										Description: "",
										Computed:    true,
										ElementType: types.StringType,
									},
									"name": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
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
func (d *fileBlockingProfilesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *fileBlockingProfilesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state fileBlockingProfilesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_file_blocking_profiles_list",
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
	svc := fEpWCgc.NewClient(d.client)
	input := fEpWCgc.ListInput{
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
	var var0 []fileBlockingProfilesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]fileBlockingProfilesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 fileBlockingProfilesListDsModelConfig
			var var3 []fileBlockingProfilesListDsModelRulesObject
			if len(var1.Rules) != 0 {
				var3 = make([]fileBlockingProfilesListDsModelRulesObject, 0, len(var1.Rules))
				for var4Index := range var1.Rules {
					var4 := var1.Rules[var4Index]
					var var5 fileBlockingProfilesListDsModelRulesObject
					var5.Action = types.StringValue(var4.Action)
					var5.Application = EncodeStringSlice(var4.Application)
					var5.Direction = types.StringValue(var4.Direction)
					var5.FileType = EncodeStringSlice(var4.FileType)
					var5.Name = types.StringValue(var4.Name)
					var3 = append(var3, var5)
				}
			}
			var2.Description = types.StringValue(var1.Description)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Name = types.StringValue(var1.Name)
			var2.Rules = var3
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
	_ datasource.DataSource              = &fileBlockingProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &fileBlockingProfilesDataSource{}
)

func NewFileBlockingProfilesDataSource() datasource.DataSource {
	return &fileBlockingProfilesDataSource{}
}

type fileBlockingProfilesDataSource struct {
	client *sase.Client
}

type fileBlockingProfilesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`
	Folder   types.String `tfsdk:"folder"`

	// Output.
	// Ref: #/components/schemas/file-blocking-profiles
	Description types.String `tfsdk:"description"`
	// input omit: ObjectId
	Name  types.String                             `tfsdk:"name"`
	Rules []fileBlockingProfilesDsModelRulesObject `tfsdk:"rules"`
}

type fileBlockingProfilesDsModelRulesObject struct {
	Action      types.String   `tfsdk:"action"`
	Application []types.String `tfsdk:"application"`
	Direction   types.String   `tfsdk:"direction"`
	FileType    []types.String `tfsdk:"file_type"`
	Name        types.String   `tfsdk:"name"`
}

// Metadata returns the data source type name.
func (d *fileBlockingProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file_blocking_profiles"
}

// Schema defines the schema for this listing data source.
func (d *fileBlockingProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"description": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"rules": dsschema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"action": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"application": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"direction": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"file_type": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure prepares the struct.
func (d *fileBlockingProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *fileBlockingProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state fileBlockingProfilesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_file_blocking_profiles",
		"object_id":                   state.ObjectId.ValueString(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := fEpWCgc.NewClient(d.client)
	input := fEpWCgc.ReadInput{
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
	var var0 []fileBlockingProfilesDsModelRulesObject
	if len(ans.Rules) != 0 {
		var0 = make([]fileBlockingProfilesDsModelRulesObject, 0, len(ans.Rules))
		for var1Index := range ans.Rules {
			var1 := ans.Rules[var1Index]
			var var2 fileBlockingProfilesDsModelRulesObject
			var2.Action = types.StringValue(var1.Action)
			var2.Application = EncodeStringSlice(var1.Application)
			var2.Direction = types.StringValue(var1.Direction)
			var2.FileType = EncodeStringSlice(var1.FileType)
			var2.Name = types.StringValue(var1.Name)
			var0 = append(var0, var2)
		}
	}
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.Rules = var0

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &fileBlockingProfilesResource{}
	_ resource.ResourceWithConfigure   = &fileBlockingProfilesResource{}
	_ resource.ResourceWithImportState = &fileBlockingProfilesResource{}
)

func NewFileBlockingProfilesResource() resource.Resource {
	return &fileBlockingProfilesResource{}
}

type fileBlockingProfilesResource struct {
	client *sase.Client
}

type fileBlockingProfilesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/file-blocking-profiles
	Description types.String                             `tfsdk:"description"`
	ObjectId    types.String                             `tfsdk:"object_id"`
	Name        types.String                             `tfsdk:"name"`
	Rules       []fileBlockingProfilesRsModelRulesObject `tfsdk:"rules"`
}

type fileBlockingProfilesRsModelRulesObject struct {
	Action      types.String   `tfsdk:"action"`
	Application []types.String `tfsdk:"application"`
	Direction   types.String   `tfsdk:"direction"`
	FileType    []types.String `tfsdk:"file_type"`
	Name        types.String   `tfsdk:"name"`
}

// Metadata returns the data source type name.
func (r *fileBlockingProfilesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file_blocking_profiles"
}

// Schema defines the schema for this listing data source.
func (r *fileBlockingProfilesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

			"description": rsschema.StringAttribute{
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
			},
			"rules": rsschema.ListNestedAttribute{
				Description: "",
				Optional:    true,
				NestedObject: rsschema.NestedAttributeObject{
					Attributes: map[string]rsschema.Attribute{
						"action": rsschema.StringAttribute{
							Description: "",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("alert", "block", "continue"),
							},
						},
						"application": rsschema.ListAttribute{
							Description: "",
							Required:    true,
							ElementType: types.StringType,
						},
						"direction": rsschema.StringAttribute{
							Description: "",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("download", "upload", "both"),
							},
						},
						"file_type": rsschema.ListAttribute{
							Description: "",
							Required:    true,
							ElementType: types.StringType,
						},
						"name": rsschema.StringAttribute{
							Description: "",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *fileBlockingProfilesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *fileBlockingProfilesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state fileBlockingProfilesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_file_blocking_profiles",
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := fEpWCgc.NewClient(r.client)
	input := fEpWCgc.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 wMgZmmI.Config
	var0.Description = state.Description.ValueString()
	var0.Name = state.Name.ValueString()
	var var1 []wMgZmmI.RulesObject
	if len(state.Rules) != 0 {
		var1 = make([]wMgZmmI.RulesObject, 0, len(state.Rules))
		for var2Index := range state.Rules {
			var2 := state.Rules[var2Index]
			var var3 wMgZmmI.RulesObject
			var3.Action = var2.Action.ValueString()
			var3.Application = DecodeStringSlice(var2.Application)
			var3.Direction = var2.Direction.ValueString()
			var3.FileType = DecodeStringSlice(var2.FileType)
			var3.Name = var2.Name.ValueString()
			var1 = append(var1, var3)
		}
	}
	var0.Rules = var1
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
	var var4 []fileBlockingProfilesRsModelRulesObject
	if len(ans.Rules) != 0 {
		var4 = make([]fileBlockingProfilesRsModelRulesObject, 0, len(ans.Rules))
		for var5Index := range ans.Rules {
			var5 := ans.Rules[var5Index]
			var var6 fileBlockingProfilesRsModelRulesObject
			var6.Action = types.StringValue(var5.Action)
			var6.Application = EncodeStringSlice(var5.Application)
			var6.Direction = types.StringValue(var5.Direction)
			var6.FileType = EncodeStringSlice(var5.FileType)
			var6.Name = types.StringValue(var5.Name)
			var4 = append(var4, var6)
		}
	}
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.Rules = var4

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *fileBlockingProfilesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state fileBlockingProfilesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_file_blocking_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := fEpWCgc.NewClient(r.client)
	input := fEpWCgc.ReadInput{
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
	var var0 []fileBlockingProfilesRsModelRulesObject
	if len(ans.Rules) != 0 {
		var0 = make([]fileBlockingProfilesRsModelRulesObject, 0, len(ans.Rules))
		for var1Index := range ans.Rules {
			var1 := ans.Rules[var1Index]
			var var2 fileBlockingProfilesRsModelRulesObject
			var2.Action = types.StringValue(var1.Action)
			var2.Application = EncodeStringSlice(var1.Application)
			var2.Direction = types.StringValue(var1.Direction)
			var2.FileType = EncodeStringSlice(var1.FileType)
			var2.Name = types.StringValue(var1.Name)
			var0 = append(var0, var2)
		}
	}
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.Rules = var0

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *fileBlockingProfilesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state fileBlockingProfilesRsModel
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
		"resource_name":               "sase_file_blocking_profiles",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := fEpWCgc.NewClient(r.client)
	input := fEpWCgc.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 wMgZmmI.Config
	var0.Description = plan.Description.ValueString()
	var0.Name = plan.Name.ValueString()
	var var1 []wMgZmmI.RulesObject
	if len(plan.Rules) != 0 {
		var1 = make([]wMgZmmI.RulesObject, 0, len(plan.Rules))
		for var2Index := range plan.Rules {
			var2 := plan.Rules[var2Index]
			var var3 wMgZmmI.RulesObject
			var3.Action = var2.Action.ValueString()
			var3.Application = DecodeStringSlice(var2.Application)
			var3.Direction = var2.Direction.ValueString()
			var3.FileType = DecodeStringSlice(var2.FileType)
			var3.Name = var2.Name.ValueString()
			var1 = append(var1, var3)
		}
	}
	var0.Rules = var1
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var4 []fileBlockingProfilesRsModelRulesObject
	if len(ans.Rules) != 0 {
		var4 = make([]fileBlockingProfilesRsModelRulesObject, 0, len(ans.Rules))
		for var5Index := range ans.Rules {
			var5 := ans.Rules[var5Index]
			var var6 fileBlockingProfilesRsModelRulesObject
			var6.Action = types.StringValue(var5.Action)
			var6.Application = EncodeStringSlice(var5.Application)
			var6.Direction = types.StringValue(var5.Direction)
			var6.FileType = EncodeStringSlice(var5.FileType)
			var6.Name = types.StringValue(var5.Name)
			var4 = append(var4, var6)
		}
	}
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.Rules = var4

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *fileBlockingProfilesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name":               "sase_file_blocking_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	svc := fEpWCgc.NewClient(r.client)
	input := fEpWCgc.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *fileBlockingProfilesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
