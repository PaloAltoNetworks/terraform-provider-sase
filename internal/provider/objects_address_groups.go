package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	nVitIaG "github.com/paloaltonetworks/sase-go/netsec/schema/objects/address/groups"
	mIAatvm "github.com/paloaltonetworks/sase-go/netsec/service/v1/addressgroups"

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
	_ datasource.DataSource              = &objectsAddressGroupsListDataSource{}
	_ datasource.DataSourceWithConfigure = &objectsAddressGroupsListDataSource{}
)

func NewObjectsAddressGroupsListDataSource() datasource.DataSource {
	return &objectsAddressGroupsListDataSource{}
}

type objectsAddressGroupsListDataSource struct {
	client *sase.Client
}

type objectsAddressGroupsListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []objectsAddressGroupsListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type objectsAddressGroupsListDsModelConfig struct {
	Description  types.String                                  `tfsdk:"description"`
	DynamicValue *objectsAddressGroupsListDsModelDynamicObject `tfsdk:"dynamic_value"`
	ObjectId     types.String                                  `tfsdk:"object_id"`
	Name         types.String                                  `tfsdk:"name"`
	Static       []types.String                                `tfsdk:"static"`
	Tag          []types.String                                `tfsdk:"tag"`
}

type objectsAddressGroupsListDsModelDynamicObject struct {
	Filter types.String `tfsdk:"filter"`
}

// Metadata returns the data source type name.
func (d *objectsAddressGroupsListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_address_groups_list"
}

// Schema defines the schema for this listing data source.
func (d *objectsAddressGroupsListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"dynamic_value": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"filter": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
							},
						},
						"object_id": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"static": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"tag": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
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
func (d *objectsAddressGroupsListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *objectsAddressGroupsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state objectsAddressGroupsListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_objects_address_groups_list",
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
	svc := mIAatvm.NewClient(d.client)
	input := mIAatvm.ListInput{
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
	var var0 []objectsAddressGroupsListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]objectsAddressGroupsListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 objectsAddressGroupsListDsModelConfig
			var var3 *objectsAddressGroupsListDsModelDynamicObject
			if var1.DynamicValue != nil {
				var3 = &objectsAddressGroupsListDsModelDynamicObject{}
				var3.Filter = types.StringValue(var1.DynamicValue.Filter)
			}
			var2.Description = types.StringValue(var1.Description)
			var2.DynamicValue = var3
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Name = types.StringValue(var1.Name)
			var2.Static = EncodeStringSlice(var1.Static)
			var2.Tag = EncodeStringSlice(var1.Tag)
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
	_ datasource.DataSource              = &objectsAddressGroupsDataSource{}
	_ datasource.DataSourceWithConfigure = &objectsAddressGroupsDataSource{}
)

func NewObjectsAddressGroupsDataSource() datasource.DataSource {
	return &objectsAddressGroupsDataSource{}
}

type objectsAddressGroupsDataSource struct {
	client *sase.Client
}

type objectsAddressGroupsDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`

	// Output.
	// Ref: #/components/schemas/objects-address-groups
	Description  types.String                              `tfsdk:"description"`
	DynamicValue *objectsAddressGroupsDsModelDynamicObject `tfsdk:"dynamic_value"`
	// input omit: ObjectId
	Name   types.String   `tfsdk:"name"`
	Static []types.String `tfsdk:"static"`
	Tag    []types.String `tfsdk:"tag"`
}

type objectsAddressGroupsDsModelDynamicObject struct {
	Filter types.String `tfsdk:"filter"`
}

// Metadata returns the data source type name.
func (d *objectsAddressGroupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_address_groups"
}

// Schema defines the schema for this listing data source.
func (d *objectsAddressGroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"description": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"dynamic_value": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"filter": dsschema.StringAttribute{
						Description: "",
						Computed:    true,
					},
				},
			},
			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"static": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"tag": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure prepares the struct.
func (d *objectsAddressGroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *objectsAddressGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state objectsAddressGroupsDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_objects_address_groups",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to run the command.
	svc := mIAatvm.NewClient(d.client)
	input := mIAatvm.ReadInput{
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
	var var0 *objectsAddressGroupsDsModelDynamicObject
	if ans.DynamicValue != nil {
		var0 = &objectsAddressGroupsDsModelDynamicObject{}
		var0.Filter = types.StringValue(ans.DynamicValue.Filter)
	}
	state.Description = types.StringValue(ans.Description)
	state.DynamicValue = var0
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.Static = EncodeStringSlice(ans.Static)
	state.Tag = EncodeStringSlice(ans.Tag)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &objectsAddressGroupsResource{}
	_ resource.ResourceWithConfigure   = &objectsAddressGroupsResource{}
	_ resource.ResourceWithImportState = &objectsAddressGroupsResource{}
)

func NewObjectsAddressGroupsResource() resource.Resource {
	return &objectsAddressGroupsResource{}
}

type objectsAddressGroupsResource struct {
	client *sase.Client
}

type objectsAddressGroupsRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/objects-address-groups
	Description  types.String                              `tfsdk:"description"`
	DynamicValue *objectsAddressGroupsRsModelDynamicObject `tfsdk:"dynamic_value"`
	ObjectId     types.String                              `tfsdk:"object_id"`
	Name         types.String                              `tfsdk:"name"`
	Static       []types.String                            `tfsdk:"static"`
	Tag          []types.String                            `tfsdk:"tag"`
}

type objectsAddressGroupsRsModelDynamicObject struct {
	Filter types.String `tfsdk:"filter"`
}

// Metadata returns the data source type name.
func (r *objectsAddressGroupsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_address_groups"
}

// Schema defines the schema for this listing data source.
func (r *objectsAddressGroupsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 1023),
				},
			},
			"dynamic_value": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"filter": rsschema.StringAttribute{
						Description: "",
						Required:    true,
						PlanModifiers: []planmodifier.String{
							DefaultString(""),
						},
						Validators: []validator.String{
							stringvalidator.LengthAtMost(2047),
						},
					},
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
			"static": rsschema.ListAttribute{
				Description: "",
				Optional:    true,
				ElementType: types.StringType,
			},
			"tag": rsschema.ListAttribute{
				Description: "",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure prepares the struct.
func (r *objectsAddressGroupsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *objectsAddressGroupsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state objectsAddressGroupsRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_objects_address_groups",
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := mIAatvm.NewClient(r.client)
	input := mIAatvm.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 nVitIaG.Config
	var0.Description = state.Description.ValueString()
	var var1 *nVitIaG.DynamicObject
	if state.DynamicValue != nil {
		var1 = &nVitIaG.DynamicObject{}
		var1.Filter = state.DynamicValue.Filter.ValueString()
	}
	var0.DynamicValue = var1
	var0.Name = state.Name.ValueString()
	var0.Static = DecodeStringSlice(state.Static)
	var0.Tag = DecodeStringSlice(state.Tag)
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
	var var2 *objectsAddressGroupsRsModelDynamicObject
	if ans.DynamicValue != nil {
		var2 = &objectsAddressGroupsRsModelDynamicObject{}
		var2.Filter = types.StringValue(ans.DynamicValue.Filter)
	}
	state.Description = types.StringValue(ans.Description)
	state.DynamicValue = var2
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.Static = EncodeStringSlice(ans.Static)
	state.Tag = EncodeStringSlice(ans.Tag)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *objectsAddressGroupsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state objectsAddressGroupsRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_objects_address_groups",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := mIAatvm.NewClient(r.client)
	input := mIAatvm.ReadInput{
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
	var var0 *objectsAddressGroupsRsModelDynamicObject
	if ans.DynamicValue != nil {
		var0 = &objectsAddressGroupsRsModelDynamicObject{}
		var0.Filter = types.StringValue(ans.DynamicValue.Filter)
	}
	state.Description = types.StringValue(ans.Description)
	state.DynamicValue = var0
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.Static = EncodeStringSlice(ans.Static)
	state.Tag = EncodeStringSlice(ans.Tag)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *objectsAddressGroupsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state objectsAddressGroupsRsModel
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
		"resource_name":               "sase_objects_address_groups",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := mIAatvm.NewClient(r.client)
	input := mIAatvm.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 nVitIaG.Config
	var0.Description = plan.Description.ValueString()
	var var1 *nVitIaG.DynamicObject
	if plan.DynamicValue != nil {
		var1 = &nVitIaG.DynamicObject{}
		var1.Filter = plan.DynamicValue.Filter.ValueString()
	}
	var0.DynamicValue = var1
	var0.Name = plan.Name.ValueString()
	var0.Static = DecodeStringSlice(plan.Static)
	var0.Tag = DecodeStringSlice(plan.Tag)
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var2 *objectsAddressGroupsRsModelDynamicObject
	if ans.DynamicValue != nil {
		var2 = &objectsAddressGroupsRsModelDynamicObject{}
		var2.Filter = types.StringValue(ans.DynamicValue.Filter)
	}
	state.Description = types.StringValue(ans.Description)
	state.DynamicValue = var2
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.Static = EncodeStringSlice(ans.Static)
	state.Tag = EncodeStringSlice(ans.Tag)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *objectsAddressGroupsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name":               "sase_objects_address_groups",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	svc := mIAatvm.NewClient(r.client)
	input := mIAatvm.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *objectsAddressGroupsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
