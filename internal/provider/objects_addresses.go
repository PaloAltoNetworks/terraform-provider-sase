package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	evToKLE "github.com/paloaltonetworks/sase-go/netsec/schema/objects/addresses"
	zLXjrfn "github.com/paloaltonetworks/sase-go/netsec/service/v1/addresses"

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
	_ datasource.DataSource              = &objectsAddressesListDataSource{}
	_ datasource.DataSourceWithConfigure = &objectsAddressesListDataSource{}
)

func NewObjectsAddressesListDataSource() datasource.DataSource {
	return &objectsAddressesListDataSource{}
}

type objectsAddressesListDataSource struct {
	client *sase.Client
}

type objectsAddressesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []objectsAddressesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type objectsAddressesListDsModelConfig struct {
	Description types.String   `tfsdk:"description"`
	Fqdn        types.String   `tfsdk:"fqdn"`
	ObjectId    types.String   `tfsdk:"object_id"`
	IpNetmask   types.String   `tfsdk:"ip_netmask"`
	IpRange     types.String   `tfsdk:"ip_range"`
	IpWildcard  types.String   `tfsdk:"ip_wildcard"`
	Name        types.String   `tfsdk:"name"`
	Tag         []types.String `tfsdk:"tag"`
	Type        types.String   `tfsdk:"type"`
}

// Metadata returns the data source type name.
func (d *objectsAddressesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_addresses_list"
}

// Schema defines the schema for this listing data source.
func (d *objectsAddressesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"fqdn": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"object_id": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"ip_netmask": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"ip_range": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"ip_wildcard": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"tag": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"type": dsschema.StringAttribute{
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
func (d *objectsAddressesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *objectsAddressesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state objectsAddressesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_objects_addresses_list",
		"limit":            state.Limit.ValueInt64(),
		"has_limit":        !state.Limit.IsNull(),
		"offset":           state.Offset.ValueInt64(),
		"has_offset":       !state.Offset.IsNull(),
		"name":             state.Name.ValueString(),
		"has_name":         !state.Name.IsNull(),
		"folder":           state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := zLXjrfn.NewClient(d.client)
	input := zLXjrfn.ListInput{
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
	var var0 []objectsAddressesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]objectsAddressesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 objectsAddressesListDsModelConfig
			var2.Description = types.StringValue(var1.Description)
			var2.Fqdn = types.StringValue(var1.Fqdn)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.IpNetmask = types.StringValue(var1.IpNetmask)
			var2.IpRange = types.StringValue(var1.IpRange)
			var2.IpWildcard = types.StringValue(var1.IpWildcard)
			var2.Name = types.StringValue(var1.Name)
			var2.Tag = EncodeStringSlice(var1.Tag)
			var2.Type = types.StringValue(var1.Type)
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
	_ datasource.DataSource              = &objectsAddressesDataSource{}
	_ datasource.DataSourceWithConfigure = &objectsAddressesDataSource{}
)

func NewObjectsAddressesDataSource() datasource.DataSource {
	return &objectsAddressesDataSource{}
}

type objectsAddressesDataSource struct {
	client *sase.Client
}

type objectsAddressesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`
	Folder   types.String `tfsdk:"folder"`

	// Output.
	// Ref: #/components/schemas/objects-addresses
	Description types.String `tfsdk:"description"`
	Fqdn        types.String `tfsdk:"fqdn"`
	// input omit: ObjectId
	IpNetmask  types.String   `tfsdk:"ip_netmask"`
	IpRange    types.String   `tfsdk:"ip_range"`
	IpWildcard types.String   `tfsdk:"ip_wildcard"`
	Name       types.String   `tfsdk:"name"`
	Tag        []types.String `tfsdk:"tag"`
	Type       types.String   `tfsdk:"type"`
}

// Metadata returns the data source type name.
func (d *objectsAddressesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_addresses"
}

// Schema defines the schema for this listing data source.
func (d *objectsAddressesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"fqdn": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"ip_netmask": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"ip_range": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"ip_wildcard": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"tag": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"type": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
		},
	}
}

// Configure prepares the struct.
func (d *objectsAddressesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *objectsAddressesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state objectsAddressesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_objects_addresses",
		"object_id":        state.ObjectId.ValueString(),
		"folder":           state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := zLXjrfn.NewClient(d.client)
	input := zLXjrfn.ReadInput{
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
	state.Description = types.StringValue(ans.Description)
	state.Fqdn = types.StringValue(ans.Fqdn)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.IpNetmask = types.StringValue(ans.IpNetmask)
	state.IpRange = types.StringValue(ans.IpRange)
	state.IpWildcard = types.StringValue(ans.IpWildcard)
	state.Name = types.StringValue(ans.Name)
	state.Tag = EncodeStringSlice(ans.Tag)
	state.Type = types.StringValue(ans.Type)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &objectsAddressesResource{}
	_ resource.ResourceWithConfigure   = &objectsAddressesResource{}
	_ resource.ResourceWithImportState = &objectsAddressesResource{}
)

func NewObjectsAddressesResource() resource.Resource {
	return &objectsAddressesResource{}
}

type objectsAddressesResource struct {
	client *sase.Client
}

type objectsAddressesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/objects-addresses
	Description types.String   `tfsdk:"description"`
	Fqdn        types.String   `tfsdk:"fqdn"`
	ObjectId    types.String   `tfsdk:"object_id"`
	IpNetmask   types.String   `tfsdk:"ip_netmask"`
	IpRange     types.String   `tfsdk:"ip_range"`
	IpWildcard  types.String   `tfsdk:"ip_wildcard"`
	Name        types.String   `tfsdk:"name"`
	Tag         []types.String `tfsdk:"tag"`
	Type        types.String   `tfsdk:"type"`
}

// Metadata returns the data source type name.
func (r *objectsAddressesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_addresses"
}

// Schema defines the schema for this listing data source.
func (r *objectsAddressesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"fqdn": rsschema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("ip_netmask"),
						path.MatchRelative().AtParent().AtName("ip_range"),
						path.MatchRelative().AtParent().AtName("ip_wildcard"),
					),
				},
			},
			"object_id": rsschema.StringAttribute{
				Description: "",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ip_netmask": rsschema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("fqdn"),
						path.MatchRelative().AtParent().AtName("ip_range"),
						path.MatchRelative().AtParent().AtName("ip_wildcard"),
					),
				},
			},
			"ip_range": rsschema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("fqdn"),
						path.MatchRelative().AtParent().AtName("ip_netmask"),
						path.MatchRelative().AtParent().AtName("ip_wildcard"),
					),
				},
			},
			"ip_wildcard": rsschema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("fqdn"),
						path.MatchRelative().AtParent().AtName("ip_netmask"),
						path.MatchRelative().AtParent().AtName("ip_range"),
					),
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
			"tag": rsschema.ListAttribute{
				Description: "",
				Optional:    true,
				ElementType: types.StringType,
			},
			"type": rsschema.StringAttribute{
				Description: "",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *objectsAddressesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *objectsAddressesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state objectsAddressesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "sase_objects_addresses",
		"folder":        state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := zLXjrfn.NewClient(r.client)
	input := zLXjrfn.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 evToKLE.Config
	var0.Description = state.Description.ValueString()
	var0.Fqdn = state.Fqdn.ValueString()
	var0.IpNetmask = state.IpNetmask.ValueString()
	var0.IpRange = state.IpRange.ValueString()
	var0.IpWildcard = state.IpWildcard.ValueString()
	var0.Name = state.Name.ValueString()
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
	state.Description = types.StringValue(ans.Description)
	state.Fqdn = types.StringValue(ans.Fqdn)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.IpNetmask = types.StringValue(ans.IpNetmask)
	state.IpRange = types.StringValue(ans.IpRange)
	state.IpWildcard = types.StringValue(ans.IpWildcard)
	state.Name = types.StringValue(ans.Name)
	state.Tag = EncodeStringSlice(ans.Tag)
	state.Type = types.StringValue(ans.Type)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *objectsAddressesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state objectsAddressesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "sase_objects_addresses",
		"locMap":        map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":        tokens,
	})

	// Prepare to read the config.
	svc := zLXjrfn.NewClient(r.client)
	input := zLXjrfn.ReadInput{
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
	state.Description = types.StringValue(ans.Description)
	state.Fqdn = types.StringValue(ans.Fqdn)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.IpNetmask = types.StringValue(ans.IpNetmask)
	state.IpRange = types.StringValue(ans.IpRange)
	state.IpWildcard = types.StringValue(ans.IpWildcard)
	state.Name = types.StringValue(ans.Name)
	state.Tag = EncodeStringSlice(ans.Tag)
	state.Type = types.StringValue(ans.Type)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *objectsAddressesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state objectsAddressesRsModel
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
		"resource_name": "sase_objects_addresses",
		"object_id":     state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := zLXjrfn.NewClient(r.client)
	input := zLXjrfn.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 evToKLE.Config
	var0.Description = plan.Description.ValueString()
	var0.Fqdn = plan.Fqdn.ValueString()
	var0.IpNetmask = plan.IpNetmask.ValueString()
	var0.IpRange = plan.IpRange.ValueString()
	var0.IpWildcard = plan.IpWildcard.ValueString()
	var0.Name = plan.Name.ValueString()
	var0.Tag = DecodeStringSlice(plan.Tag)
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	state.Description = types.StringValue(ans.Description)
	state.Fqdn = types.StringValue(ans.Fqdn)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.IpNetmask = types.StringValue(ans.IpNetmask)
	state.IpRange = types.StringValue(ans.IpRange)
	state.IpWildcard = types.StringValue(ans.IpWildcard)
	state.Name = types.StringValue(ans.Name)
	state.Tag = EncodeStringSlice(ans.Tag)
	state.Type = types.StringValue(ans.Type)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *objectsAddressesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name": "sase_objects_addresses",
		"locMap":        map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":        tokens,
	})

	svc := zLXjrfn.NewClient(r.client)
	input := zLXjrfn.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *objectsAddressesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
