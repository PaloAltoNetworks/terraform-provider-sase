package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	ampruGo "github.com/paloaltonetworks/sase-go/netsec/schema/app/override/rules"
	pTgTBIe "github.com/paloaltonetworks/sase-go/netsec/service/v1/appoverriderules"

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
	_ datasource.DataSource              = &appOverrideRulesListDataSource{}
	_ datasource.DataSourceWithConfigure = &appOverrideRulesListDataSource{}
)

func NewAppOverrideRulesListDataSource() datasource.DataSource {
	return &appOverrideRulesListDataSource{}
}

type appOverrideRulesListDataSource struct {
	client *sase.Client
}

type appOverrideRulesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit    types.Int64  `tfsdk:"limit"`
	Offset   types.Int64  `tfsdk:"offset"`
	Position types.String `tfsdk:"position"`
	Folder   types.String `tfsdk:"folder"`
	Name     types.String `tfsdk:"name"`

	// Output.
	Data []appOverrideRulesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type appOverrideRulesListDsModelConfig struct {
	Application       types.String   `tfsdk:"application"`
	Description       types.String   `tfsdk:"description"`
	Destination       []types.String `tfsdk:"destination"`
	Disabled          types.Bool     `tfsdk:"disabled"`
	From              []types.String `tfsdk:"from"`
	GroupTag          types.String   `tfsdk:"group_tag"`
	ObjectId          types.String   `tfsdk:"object_id"`
	Name              types.String   `tfsdk:"name"`
	NegateDestination types.Bool     `tfsdk:"negate_destination"`
	NegateSource      types.Bool     `tfsdk:"negate_source"`
	Port              types.Int64    `tfsdk:"port"`
	Protocol          types.String   `tfsdk:"protocol"`
	Source            []types.String `tfsdk:"source"`
	Tag               []types.String `tfsdk:"tag"`
	To                []types.String `tfsdk:"to"`
}

// Metadata returns the data source type name.
func (d *appOverrideRulesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_override_rules_list"
}

// Schema defines the schema for this listing data source.
func (d *appOverrideRulesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"position": dsschema.StringAttribute{
				Description: "The position of a security rule",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("pre", "post"),
				},
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
						"application": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"description": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"destination": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"disabled": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"from": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"group_tag": dsschema.StringAttribute{
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
						"negate_destination": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"negate_source": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"port": dsschema.Int64Attribute{
							Description: "",
							Computed:    true,
						},
						"protocol": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"source": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"tag": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"to": dsschema.ListAttribute{
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
func (d *appOverrideRulesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *appOverrideRulesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state appOverrideRulesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_app_override_rules_list",
		"limit":            state.Limit.ValueInt64(),
		"has_limit":        !state.Limit.IsNull(),
		"offset":           state.Offset.ValueInt64(),
		"has_offset":       !state.Offset.IsNull(),
		"position":         state.Position.ValueString(),
		"folder":           state.Folder.ValueString(),
		"name":             state.Name.ValueString(),
		"has_name":         !state.Name.IsNull(),
	})

	// Prepare to run the command.
	svc := pTgTBIe.NewClient(d.client)
	input := pTgTBIe.ListInput{
		Position: state.Position.ValueString(),
		Folder:   state.Folder.ValueString(),
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
	state.Id = types.StringValue(strings.Join([]string{strconv.FormatInt(*input.Limit, 10), strconv.FormatInt(*input.Offset, 10), input.Position, input.Folder, *input.Name}, IdSeparator))
	var var0 []appOverrideRulesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]appOverrideRulesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 appOverrideRulesListDsModelConfig
			var2.Application = types.StringValue(var1.Application)
			var2.Description = types.StringValue(var1.Description)
			var2.Destination = EncodeStringSlice(var1.Destination)
			var2.Disabled = types.BoolValue(var1.Disabled)
			var2.From = EncodeStringSlice(var1.From)
			var2.GroupTag = types.StringValue(var1.GroupTag)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Name = types.StringValue(var1.Name)
			var2.NegateDestination = types.BoolValue(var1.NegateDestination)
			var2.NegateSource = types.BoolValue(var1.NegateSource)
			var2.Port = types.Int64Value(var1.Port)
			var2.Protocol = types.StringValue(var1.Protocol)
			var2.Source = EncodeStringSlice(var1.Source)
			var2.Tag = EncodeStringSlice(var1.Tag)
			var2.To = EncodeStringSlice(var1.To)
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
	_ datasource.DataSource              = &appOverrideRulesDataSource{}
	_ datasource.DataSourceWithConfigure = &appOverrideRulesDataSource{}
)

func NewAppOverrideRulesDataSource() datasource.DataSource {
	return &appOverrideRulesDataSource{}
}

type appOverrideRulesDataSource struct {
	client *sase.Client
}

type appOverrideRulesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`

	// Output.
	// Ref: #/components/schemas/app-override-rules
	Application types.String   `tfsdk:"application"`
	Description types.String   `tfsdk:"description"`
	Destination []types.String `tfsdk:"destination"`
	Disabled    types.Bool     `tfsdk:"disabled"`
	From        []types.String `tfsdk:"from"`
	GroupTag    types.String   `tfsdk:"group_tag"`
	// input omit: ObjectId
	Name              types.String   `tfsdk:"name"`
	NegateDestination types.Bool     `tfsdk:"negate_destination"`
	NegateSource      types.Bool     `tfsdk:"negate_source"`
	Port              types.Int64    `tfsdk:"port"`
	Protocol          types.String   `tfsdk:"protocol"`
	Source            []types.String `tfsdk:"source"`
	Tag               []types.String `tfsdk:"tag"`
	To                []types.String `tfsdk:"to"`
}

// Metadata returns the data source type name.
func (d *appOverrideRulesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_override_rules"
}

// Schema defines the schema for this listing data source.
func (d *appOverrideRulesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"application": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"description": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"destination": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"disabled": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"from": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"group_tag": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"negate_destination": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"negate_source": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"port": dsschema.Int64Attribute{
				Description: "",
				Computed:    true,
			},
			"protocol": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"source": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"tag": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"to": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure prepares the struct.
func (d *appOverrideRulesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *appOverrideRulesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state appOverrideRulesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_app_override_rules",
		"object_id":        state.ObjectId.ValueString(),
	})

	// Prepare to run the command.
	svc := pTgTBIe.NewClient(d.client)
	input := pTgTBIe.ReadInput{
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
	state.Application = types.StringValue(ans.Application)
	state.Description = types.StringValue(ans.Description)
	state.Destination = EncodeStringSlice(ans.Destination)
	state.Disabled = types.BoolValue(ans.Disabled)
	state.From = EncodeStringSlice(ans.From)
	state.GroupTag = types.StringValue(ans.GroupTag)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.NegateDestination = types.BoolValue(ans.NegateDestination)
	state.NegateSource = types.BoolValue(ans.NegateSource)
	state.Port = types.Int64Value(ans.Port)
	state.Protocol = types.StringValue(ans.Protocol)
	state.Source = EncodeStringSlice(ans.Source)
	state.Tag = EncodeStringSlice(ans.Tag)
	state.To = EncodeStringSlice(ans.To)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &appOverrideRulesResource{}
	_ resource.ResourceWithConfigure   = &appOverrideRulesResource{}
	_ resource.ResourceWithImportState = &appOverrideRulesResource{}
)

func NewAppOverrideRulesResource() resource.Resource {
	return &appOverrideRulesResource{}
}

type appOverrideRulesResource struct {
	client *sase.Client
}

type appOverrideRulesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Position types.String `tfsdk:"position"`
	Folder   types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/app-override-rules
	Application       types.String   `tfsdk:"application"`
	Description       types.String   `tfsdk:"description"`
	Destination       []types.String `tfsdk:"destination"`
	Disabled          types.Bool     `tfsdk:"disabled"`
	From              []types.String `tfsdk:"from"`
	GroupTag          types.String   `tfsdk:"group_tag"`
	ObjectId          types.String   `tfsdk:"object_id"`
	Name              types.String   `tfsdk:"name"`
	NegateDestination types.Bool     `tfsdk:"negate_destination"`
	NegateSource      types.Bool     `tfsdk:"negate_source"`
	Port              types.Int64    `tfsdk:"port"`
	Protocol          types.String   `tfsdk:"protocol"`
	Source            []types.String `tfsdk:"source"`
	Tag               []types.String `tfsdk:"tag"`
	To                []types.String `tfsdk:"to"`
}

// Metadata returns the data source type name.
func (r *appOverrideRulesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_override_rules"
}

// Schema defines the schema for this listing data source.
func (r *appOverrideRulesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"position": rsschema.StringAttribute{
				Description: "The position of a security rule",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("pre", "post"),
				},
			},
			"folder": rsschema.StringAttribute{
				Description: "The folder of the entry",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Service Connections", "Mobile Users Container", "Mobile Users Explicit Proxy"),
				},
			},

			"application": rsschema.StringAttribute{
				Description: "",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
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
					stringvalidator.LengthAtMost(1024),
				},
			},
			"destination": rsschema.ListAttribute{
				Description: "",
				Required:    true,
				ElementType: types.StringType,
			},
			"disabled": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"from": rsschema.ListAttribute{
				Description: "",
				Required:    true,
				ElementType: types.StringType,
			},
			"group_tag": rsschema.StringAttribute{
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
			"negate_destination": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"negate_source": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"port": rsschema.Int64Attribute{
				Description: "",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					DefaultInt64(0),
				},
				Validators: []validator.Int64{
					int64validator.Between(0, 65535),
				},
			},
			"protocol": rsschema.StringAttribute{
				Description: "",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("tcp", "udp"),
				},
			},
			"source": rsschema.ListAttribute{
				Description: "",
				Required:    true,
				ElementType: types.StringType,
			},
			"tag": rsschema.ListAttribute{
				Description: "",
				Optional:    true,
				ElementType: types.StringType,
			},
			"to": rsschema.ListAttribute{
				Description: "",
				Required:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure prepares the struct.
func (r *appOverrideRulesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *appOverrideRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state appOverrideRulesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "sase_app_override_rules",
		"position":      state.Position.ValueString(),
		"folder":        state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := pTgTBIe.NewClient(r.client)
	input := pTgTBIe.CreateInput{
		Position: state.Position.ValueString(),
		Folder:   state.Folder.ValueString(),
	}
	var var0 ampruGo.Config
	var0.Application = state.Application.ValueString()
	var0.Description = state.Description.ValueString()
	var0.Destination = DecodeStringSlice(state.Destination)
	var0.Disabled = state.Disabled.ValueBool()
	var0.From = DecodeStringSlice(state.From)
	var0.GroupTag = state.GroupTag.ValueString()
	var0.Name = state.Name.ValueString()
	var0.NegateDestination = state.NegateDestination.ValueBool()
	var0.NegateSource = state.NegateSource.ValueBool()
	var0.Port = state.Port.ValueInt64()
	var0.Protocol = state.Protocol.ValueString()
	var0.Source = DecodeStringSlice(state.Source)
	var0.Tag = DecodeStringSlice(state.Tag)
	var0.To = DecodeStringSlice(state.To)
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Create(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in create", err.Error())
		return
	}

	// Store the answer to state.
	state.Id = types.StringValue(strings.Join([]string{input.Position, input.Folder, ans.ObjectId}, IdSeparator))
	state.Application = types.StringValue(ans.Application)
	state.Description = types.StringValue(ans.Description)
	state.Destination = EncodeStringSlice(ans.Destination)
	state.Disabled = types.BoolValue(ans.Disabled)
	state.From = EncodeStringSlice(ans.From)
	state.GroupTag = types.StringValue(ans.GroupTag)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.NegateDestination = types.BoolValue(ans.NegateDestination)
	state.NegateSource = types.BoolValue(ans.NegateSource)
	state.Port = types.Int64Value(ans.Port)
	state.Protocol = types.StringValue(ans.Protocol)
	state.Source = EncodeStringSlice(ans.Source)
	state.Tag = EncodeStringSlice(ans.Tag)
	state.To = EncodeStringSlice(ans.To)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *appOverrideRulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state appOverrideRulesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "sase_app_override_rules",
		"locMap":        map[string]int{"Folder": 1, "ObjectId": 2, "Position": 0},
		"tokens":        tokens,
	})

	// Prepare to read the config.
	svc := pTgTBIe.NewClient(r.client)
	input := pTgTBIe.ReadInput{
		ObjectId: tokens[2],
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
	state.Application = types.StringValue(ans.Application)
	state.Description = types.StringValue(ans.Description)
	state.Destination = EncodeStringSlice(ans.Destination)
	state.Disabled = types.BoolValue(ans.Disabled)
	state.From = EncodeStringSlice(ans.From)
	state.GroupTag = types.StringValue(ans.GroupTag)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.NegateDestination = types.BoolValue(ans.NegateDestination)
	state.NegateSource = types.BoolValue(ans.NegateSource)
	state.Port = types.Int64Value(ans.Port)
	state.Protocol = types.StringValue(ans.Protocol)
	state.Source = EncodeStringSlice(ans.Source)
	state.Tag = EncodeStringSlice(ans.Tag)
	state.To = EncodeStringSlice(ans.To)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *appOverrideRulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state appOverrideRulesRsModel
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
		"resource_name": "sase_app_override_rules",
		"object_id":     state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := pTgTBIe.NewClient(r.client)
	input := pTgTBIe.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 ampruGo.Config
	var0.Application = plan.Application.ValueString()
	var0.Description = plan.Description.ValueString()
	var0.Destination = DecodeStringSlice(plan.Destination)
	var0.Disabled = plan.Disabled.ValueBool()
	var0.From = DecodeStringSlice(plan.From)
	var0.GroupTag = plan.GroupTag.ValueString()
	var0.Name = plan.Name.ValueString()
	var0.NegateDestination = plan.NegateDestination.ValueBool()
	var0.NegateSource = plan.NegateSource.ValueBool()
	var0.Port = plan.Port.ValueInt64()
	var0.Protocol = plan.Protocol.ValueString()
	var0.Source = DecodeStringSlice(plan.Source)
	var0.Tag = DecodeStringSlice(plan.Tag)
	var0.To = DecodeStringSlice(plan.To)
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	state.Application = types.StringValue(ans.Application)
	state.Description = types.StringValue(ans.Description)
	state.Destination = EncodeStringSlice(ans.Destination)
	state.Disabled = types.BoolValue(ans.Disabled)
	state.From = EncodeStringSlice(ans.From)
	state.GroupTag = types.StringValue(ans.GroupTag)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.NegateDestination = types.BoolValue(ans.NegateDestination)
	state.NegateSource = types.BoolValue(ans.NegateSource)
	state.Port = types.Int64Value(ans.Port)
	state.Protocol = types.StringValue(ans.Protocol)
	state.Source = EncodeStringSlice(ans.Source)
	state.Tag = EncodeStringSlice(ans.Tag)
	state.To = EncodeStringSlice(ans.To)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *appOverrideRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name": "sase_app_override_rules",
		"locMap":        map[string]int{"Folder": 1, "ObjectId": 2, "Position": 0},
		"tokens":        tokens,
	})

	svc := pTgTBIe.NewClient(r.client)
	input := pTgTBIe.DeleteInput{
		ObjectId: tokens[2],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *appOverrideRulesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
