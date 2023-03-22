package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	gOMQArS "github.com/paloaltonetworks/sase-go/netsec/schema/kerberos/server/profiles"
	hKcuqhS "github.com/paloaltonetworks/sase-go/netsec/service/v1/kerberosserverprofiles"

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
	_ datasource.DataSource              = &kerberosServerProfilesListDataSource{}
	_ datasource.DataSourceWithConfigure = &kerberosServerProfilesListDataSource{}
)

func NewKerberosServerProfilesListDataSource() datasource.DataSource {
	return &kerberosServerProfilesListDataSource{}
}

type kerberosServerProfilesListDataSource struct {
	client *sase.Client
}

type kerberosServerProfilesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Folder types.String `tfsdk:"folder"`
	Name   types.String `tfsdk:"name"`

	// Output.
	Data []kerberosServerProfilesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type kerberosServerProfilesListDsModelConfig struct {
	ObjectId types.String                                    `tfsdk:"object_id"`
	Server   []kerberosServerProfilesListDsModelServerObject `tfsdk:"server"`
}

type kerberosServerProfilesListDsModelServerObject struct {
	Host types.String `tfsdk:"host"`
	Name types.String `tfsdk:"name"`
	Port types.Int64  `tfsdk:"port"`
}

// Metadata returns the data source type name.
func (d *kerberosServerProfilesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kerberos_server_profiles_list"
}

// Schema defines the schema for this listing data source.
func (d *kerberosServerProfilesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"object_id": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"server": dsschema.ListNestedAttribute{
							Description: "",
							Computed:    true,
							NestedObject: dsschema.NestedAttributeObject{
								Attributes: map[string]dsschema.Attribute{
									"host": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"name": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"port": dsschema.Int64Attribute{
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
func (d *kerberosServerProfilesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *kerberosServerProfilesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state kerberosServerProfilesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_kerberos_server_profiles_list",
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
	svc := hKcuqhS.NewClient(d.client)
	input := hKcuqhS.ListInput{
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
	var var0 []kerberosServerProfilesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]kerberosServerProfilesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 kerberosServerProfilesListDsModelConfig
			var var3 []kerberosServerProfilesListDsModelServerObject
			if len(var1.Server) != 0 {
				var3 = make([]kerberosServerProfilesListDsModelServerObject, 0, len(var1.Server))
				for var4Index := range var1.Server {
					var4 := var1.Server[var4Index]
					var var5 kerberosServerProfilesListDsModelServerObject
					var5.Host = types.StringValue(var4.Host)
					var5.Name = types.StringValue(var4.Name)
					var5.Port = types.Int64Value(var4.Port)
					var3 = append(var3, var5)
				}
			}
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Server = var3
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
	_ datasource.DataSource              = &kerberosServerProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &kerberosServerProfilesDataSource{}
)

func NewKerberosServerProfilesDataSource() datasource.DataSource {
	return &kerberosServerProfilesDataSource{}
}

type kerberosServerProfilesDataSource struct {
	client *sase.Client
}

type kerberosServerProfilesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`

	// Output.
	// Ref: #/components/schemas/kerberos-server-profiles
	// input omit: ObjectId
	Server []kerberosServerProfilesDsModelServerObject `tfsdk:"server"`
}

type kerberosServerProfilesDsModelServerObject struct {
	Host types.String `tfsdk:"host"`
	Name types.String `tfsdk:"name"`
	Port types.Int64  `tfsdk:"port"`
}

// Metadata returns the data source type name.
func (d *kerberosServerProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kerberos_server_profiles"
}

// Schema defines the schema for this listing data source.
func (d *kerberosServerProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"server": dsschema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"host": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"port": dsschema.Int64Attribute{
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
func (d *kerberosServerProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *kerberosServerProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state kerberosServerProfilesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_kerberos_server_profiles",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to run the command.
	svc := hKcuqhS.NewClient(d.client)
	input := hKcuqhS.ReadInput{
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
	var var0 []kerberosServerProfilesDsModelServerObject
	if len(ans.Server) != 0 {
		var0 = make([]kerberosServerProfilesDsModelServerObject, 0, len(ans.Server))
		for var1Index := range ans.Server {
			var1 := ans.Server[var1Index]
			var var2 kerberosServerProfilesDsModelServerObject
			var2.Host = types.StringValue(var1.Host)
			var2.Name = types.StringValue(var1.Name)
			var2.Port = types.Int64Value(var1.Port)
			var0 = append(var0, var2)
		}
	}
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Server = var0

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &kerberosServerProfilesResource{}
	_ resource.ResourceWithConfigure   = &kerberosServerProfilesResource{}
	_ resource.ResourceWithImportState = &kerberosServerProfilesResource{}
)

func NewKerberosServerProfilesResource() resource.Resource {
	return &kerberosServerProfilesResource{}
}

type kerberosServerProfilesResource struct {
	client *sase.Client
}

type kerberosServerProfilesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/kerberos-server-profiles
	ObjectId types.String                                `tfsdk:"object_id"`
	Server   []kerberosServerProfilesRsModelServerObject `tfsdk:"server"`
}

type kerberosServerProfilesRsModelServerObject struct {
	Host types.String `tfsdk:"host"`
	Name types.String `tfsdk:"name"`
	Port types.Int64  `tfsdk:"port"`
}

// Metadata returns the data source type name.
func (r *kerberosServerProfilesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kerberos_server_profiles"
}

// Schema defines the schema for this listing data source.
func (r *kerberosServerProfilesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

			"object_id": rsschema.StringAttribute{
				Description: "",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"server": rsschema.ListNestedAttribute{
				Description: "",
				Required:    true,
				NestedObject: rsschema.NestedAttributeObject{
					Attributes: map[string]rsschema.Attribute{
						"host": rsschema.StringAttribute{
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
						"port": rsschema.Int64Attribute{
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
					},
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *kerberosServerProfilesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *kerberosServerProfilesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state kerberosServerProfilesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_kerberos_server_profiles",
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := hKcuqhS.NewClient(r.client)
	input := hKcuqhS.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 gOMQArS.Config
	var var1 []gOMQArS.ServerObject
	if len(state.Server) != 0 {
		var1 = make([]gOMQArS.ServerObject, 0, len(state.Server))
		for var2Index := range state.Server {
			var2 := state.Server[var2Index]
			var var3 gOMQArS.ServerObject
			var3.Host = var2.Host.ValueString()
			var3.Name = var2.Name.ValueString()
			var3.Port = var2.Port.ValueInt64()
			var1 = append(var1, var3)
		}
	}
	var0.Server = var1
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
	var var4 []kerberosServerProfilesRsModelServerObject
	if len(ans.Server) != 0 {
		var4 = make([]kerberosServerProfilesRsModelServerObject, 0, len(ans.Server))
		for var5Index := range ans.Server {
			var5 := ans.Server[var5Index]
			var var6 kerberosServerProfilesRsModelServerObject
			var6.Host = types.StringValue(var5.Host)
			var6.Name = types.StringValue(var5.Name)
			var6.Port = types.Int64Value(var5.Port)
			var4 = append(var4, var6)
		}
	}
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Server = var4

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *kerberosServerProfilesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state kerberosServerProfilesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_kerberos_server_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := hKcuqhS.NewClient(r.client)
	input := hKcuqhS.ReadInput{
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
	var var0 []kerberosServerProfilesRsModelServerObject
	if len(ans.Server) != 0 {
		var0 = make([]kerberosServerProfilesRsModelServerObject, 0, len(ans.Server))
		for var1Index := range ans.Server {
			var1 := ans.Server[var1Index]
			var var2 kerberosServerProfilesRsModelServerObject
			var2.Host = types.StringValue(var1.Host)
			var2.Name = types.StringValue(var1.Name)
			var2.Port = types.Int64Value(var1.Port)
			var0 = append(var0, var2)
		}
	}
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Server = var0

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *kerberosServerProfilesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state kerberosServerProfilesRsModel
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
		"resource_name":               "sase_kerberos_server_profiles",
		"object_id":                   state.ObjectId.ValueString(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := hKcuqhS.NewClient(r.client)
	input := hKcuqhS.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
		Folder:   state.Folder.ValueString(),
	}
	var var0 gOMQArS.Config
	var var1 []gOMQArS.ServerObject
	if len(plan.Server) != 0 {
		var1 = make([]gOMQArS.ServerObject, 0, len(plan.Server))
		for var2Index := range plan.Server {
			var2 := plan.Server[var2Index]
			var var3 gOMQArS.ServerObject
			var3.Host = var2.Host.ValueString()
			var3.Name = var2.Name.ValueString()
			var3.Port = var2.Port.ValueInt64()
			var1 = append(var1, var3)
		}
	}
	var0.Server = var1
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var4 []kerberosServerProfilesRsModelServerObject
	if len(ans.Server) != 0 {
		var4 = make([]kerberosServerProfilesRsModelServerObject, 0, len(ans.Server))
		for var5Index := range ans.Server {
			var5 := ans.Server[var5Index]
			var var6 kerberosServerProfilesRsModelServerObject
			var6.Host = types.StringValue(var5.Host)
			var6.Name = types.StringValue(var5.Name)
			var6.Port = types.Int64Value(var5.Port)
			var4 = append(var4, var6)
		}
	}
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Server = var4

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *kerberosServerProfilesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name":               "sase_kerberos_server_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	svc := hKcuqhS.NewClient(r.client)
	input := hKcuqhS.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *kerberosServerProfilesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
