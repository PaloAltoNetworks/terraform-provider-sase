package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	sdhSKaQ "github.com/paloaltonetworks/sase-go/netsec/schema/objects/regions"
	hhIWLbI "github.com/paloaltonetworks/sase-go/netsec/service/v1/regions"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
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
	_ datasource.DataSource              = &objectsRegionsListDataSource{}
	_ datasource.DataSourceWithConfigure = &objectsRegionsListDataSource{}
)

func NewObjectsRegionsListDataSource() datasource.DataSource {
	return &objectsRegionsListDataSource{}
}

type objectsRegionsListDataSource struct {
	client *sase.Client
}

type objectsRegionsListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []objectsRegionsListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type objectsRegionsListDsModelConfig struct {
	Address     []types.String                              `tfsdk:"address"`
	GeoLocation *objectsRegionsListDsModelGeoLocationObject `tfsdk:"geo_location"`
	ObjectId    types.String                                `tfsdk:"object_id"`
	Name        types.String                                `tfsdk:"name"`
}

type objectsRegionsListDsModelGeoLocationObject struct {
	Latitude  types.Float64 `tfsdk:"latitude"`
	Longitude types.Float64 `tfsdk:"longitude"`
}

// Metadata returns the data source type name.
func (d *objectsRegionsListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_regions_list"
}

// Schema defines the schema for this listing data source.
func (d *objectsRegionsListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"address": dsschema.ListAttribute{
							Description:         "The `address` parameter.",
							MarkdownDescription: "The `address` parameter.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"geo_location": dsschema.SingleNestedAttribute{
							Description:         "The `geo_location` parameter.",
							MarkdownDescription: "The `geo_location` parameter.",
							Computed:            true,
							Attributes: map[string]dsschema.Attribute{
								"latitude": dsschema.Float64Attribute{
									Description:         "The `latitude` parameter.",
									MarkdownDescription: "The `latitude` parameter.",
									Computed:            true,
								},
								"longitude": dsschema.Float64Attribute{
									Description:         "The `longitude` parameter.",
									MarkdownDescription: "The `longitude` parameter.",
									Computed:            true,
								},
							},
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
func (d *objectsRegionsListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *objectsRegionsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state objectsRegionsListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_objects_regions_list",
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
	svc := hhIWLbI.NewClient(d.client)
	input := hhIWLbI.ListInput{
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
	var var0 []objectsRegionsListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]objectsRegionsListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 objectsRegionsListDsModelConfig
			var var3 *objectsRegionsListDsModelGeoLocationObject
			if var1.GeoLocation != nil {
				var3 = &objectsRegionsListDsModelGeoLocationObject{}
				var3.Latitude = types.Float64Value(var1.GeoLocation.Latitude)
				var3.Longitude = types.Float64Value(var1.GeoLocation.Longitude)
			}
			var2.Address = EncodeStringSlice(var1.Address)
			var2.GeoLocation = var3
			var2.ObjectId = types.StringValue(var1.ObjectId)
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
	_ datasource.DataSource              = &objectsRegionsDataSource{}
	_ datasource.DataSourceWithConfigure = &objectsRegionsDataSource{}
)

func NewObjectsRegionsDataSource() datasource.DataSource {
	return &objectsRegionsDataSource{}
}

type objectsRegionsDataSource struct {
	client *sase.Client
}

type objectsRegionsDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`

	// Output.
	// Ref: #/components/schemas/objects-regions
	Address     []types.String                          `tfsdk:"address"`
	GeoLocation *objectsRegionsDsModelGeoLocationObject `tfsdk:"geo_location"`
	// input omit: ObjectId
	Name types.String `tfsdk:"name"`
}

type objectsRegionsDsModelGeoLocationObject struct {
	Latitude  types.Float64 `tfsdk:"latitude"`
	Longitude types.Float64 `tfsdk:"longitude"`
}

// Metadata returns the data source type name.
func (d *objectsRegionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_regions"
}

// Schema defines the schema for this listing data source.
func (d *objectsRegionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

			// Output.
			"address": dsschema.ListAttribute{
				Description:         "The `address` parameter.",
				MarkdownDescription: "The `address` parameter.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"geo_location": dsschema.SingleNestedAttribute{
				Description:         "The `geo_location` parameter.",
				MarkdownDescription: "The `geo_location` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"latitude": dsschema.Float64Attribute{
						Description:         "The `latitude` parameter.",
						MarkdownDescription: "The `latitude` parameter.",
						Computed:            true,
					},
					"longitude": dsschema.Float64Attribute{
						Description:         "The `longitude` parameter.",
						MarkdownDescription: "The `longitude` parameter.",
						Computed:            true,
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
func (d *objectsRegionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *objectsRegionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state objectsRegionsDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_objects_regions",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to run the command.
	svc := hhIWLbI.NewClient(d.client)
	input := hhIWLbI.ReadInput{
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
	var var0 *objectsRegionsDsModelGeoLocationObject
	if ans.GeoLocation != nil {
		var0 = &objectsRegionsDsModelGeoLocationObject{}
		var0.Latitude = types.Float64Value(ans.GeoLocation.Latitude)
		var0.Longitude = types.Float64Value(ans.GeoLocation.Longitude)
	}
	state.Address = EncodeStringSlice(ans.Address)
	state.GeoLocation = var0
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &objectsRegionsResource{}
	_ resource.ResourceWithConfigure   = &objectsRegionsResource{}
	_ resource.ResourceWithImportState = &objectsRegionsResource{}
)

func NewObjectsRegionsResource() resource.Resource {
	return &objectsRegionsResource{}
}

type objectsRegionsResource struct {
	client *sase.Client
}

type objectsRegionsRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/objects-regions
	Address     []types.String                          `tfsdk:"address"`
	GeoLocation *objectsRegionsRsModelGeoLocationObject `tfsdk:"geo_location"`
	ObjectId    types.String                            `tfsdk:"object_id"`
	Name        types.String                            `tfsdk:"name"`
}

type objectsRegionsRsModelGeoLocationObject struct {
	Latitude  types.Float64 `tfsdk:"latitude"`
	Longitude types.Float64 `tfsdk:"longitude"`
}

// Metadata returns the data source type name.
func (r *objectsRegionsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_regions"
}

// Schema defines the schema for this listing data source.
func (r *objectsRegionsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

			"address": rsschema.ListAttribute{
				Description:         "The `address` parameter.",
				MarkdownDescription: "The `address` parameter.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"geo_location": rsschema.SingleNestedAttribute{
				Description:         "The `geo_location` parameter.",
				MarkdownDescription: "The `geo_location` parameter.",
				Optional:            true,
				Attributes: map[string]rsschema.Attribute{
					"latitude": rsschema.Float64Attribute{
						Description:         "The `latitude` parameter. Value must be between -90 and 90.",
						MarkdownDescription: "The `latitude` parameter. Value must be between -90 and 90.",
						Required:            true,
						Validators: []validator.Float64{
							float64validator.Between(-90.000000, 90.000000),
						},
					},
					"longitude": rsschema.Float64Attribute{
						Description:         "The `longitude` parameter. Value must be between -180 and 180.",
						MarkdownDescription: "The `longitude` parameter. Value must be between -180 and 180.",
						Required:            true,
						Validators: []validator.Float64{
							float64validator.Between(-180.000000, 180.000000),
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
			"name": rsschema.StringAttribute{
				Description:         "The `name` parameter. String length must be at most 31.",
				MarkdownDescription: "The `name` parameter. String length must be at most 31.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(31),
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *objectsRegionsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *objectsRegionsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state objectsRegionsRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_objects_regions",
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := hhIWLbI.NewClient(r.client)
	input := hhIWLbI.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 sdhSKaQ.Config
	var0.Address = DecodeStringSlice(state.Address)
	var var1 *sdhSKaQ.GeoLocationObject
	if state.GeoLocation != nil {
		var1 = &sdhSKaQ.GeoLocationObject{}
		var1.Latitude = state.GeoLocation.Latitude.ValueFloat64()
		var1.Longitude = state.GeoLocation.Longitude.ValueFloat64()
	}
	var0.GeoLocation = var1
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
	var var2 *objectsRegionsRsModelGeoLocationObject
	if ans.GeoLocation != nil {
		var2 = &objectsRegionsRsModelGeoLocationObject{}
		var2.Latitude = types.Float64Value(ans.GeoLocation.Latitude)
		var2.Longitude = types.Float64Value(ans.GeoLocation.Longitude)
	}
	state.Address = EncodeStringSlice(ans.Address)
	state.GeoLocation = var2
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *objectsRegionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state objectsRegionsRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_objects_regions",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := hhIWLbI.NewClient(r.client)
	input := hhIWLbI.ReadInput{
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
	var var0 *objectsRegionsRsModelGeoLocationObject
	if ans.GeoLocation != nil {
		var0 = &objectsRegionsRsModelGeoLocationObject{}
		var0.Latitude = types.Float64Value(ans.GeoLocation.Latitude)
		var0.Longitude = types.Float64Value(ans.GeoLocation.Longitude)
	}
	state.Address = EncodeStringSlice(ans.Address)
	state.GeoLocation = var0
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *objectsRegionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state objectsRegionsRsModel
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
		"resource_name":               "sase_objects_regions",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := hhIWLbI.NewClient(r.client)
	input := hhIWLbI.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 sdhSKaQ.Config
	var0.Address = DecodeStringSlice(plan.Address)
	var var1 *sdhSKaQ.GeoLocationObject
	if plan.GeoLocation != nil {
		var1 = &sdhSKaQ.GeoLocationObject{}
		var1.Latitude = plan.GeoLocation.Latitude.ValueFloat64()
		var1.Longitude = plan.GeoLocation.Longitude.ValueFloat64()
	}
	var0.GeoLocation = var1
	var0.Name = plan.Name.ValueString()
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var2 *objectsRegionsRsModelGeoLocationObject
	if ans.GeoLocation != nil {
		var2 = &objectsRegionsRsModelGeoLocationObject{}
		var2.Latitude = types.Float64Value(ans.GeoLocation.Latitude)
		var2.Longitude = types.Float64Value(ans.GeoLocation.Longitude)
	}
	state.Address = EncodeStringSlice(ans.Address)
	state.GeoLocation = var2
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *objectsRegionsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name":               "sase_objects_regions",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	svc := hhIWLbI.NewClient(r.client)
	input := hhIWLbI.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *objectsRegionsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
