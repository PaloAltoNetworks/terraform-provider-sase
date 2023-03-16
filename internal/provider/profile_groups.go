package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	mQehIbG "github.com/paloaltonetworks/sase-go/netsec/schema/profile/groups"
	jeahrQe "github.com/paloaltonetworks/sase-go/netsec/service/v1/profilegroups"

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
	_ datasource.DataSource              = &profileGroupsListDataSource{}
	_ datasource.DataSourceWithConfigure = &profileGroupsListDataSource{}
)

func NewProfileGroupsListDataSource() datasource.DataSource {
	return &profileGroupsListDataSource{}
}

type profileGroupsListDataSource struct {
	client *sase.Client
}

type profileGroupsListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []profileGroupsListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type profileGroupsListDsModelConfig struct {
	DnsSecurity              []types.String `tfsdk:"dns_security"`
	FileBlocking             []types.String `tfsdk:"file_blocking"`
	ObjectId                 types.String   `tfsdk:"object_id"`
	Name                     types.String   `tfsdk:"name"`
	SaasSecurity             []types.String `tfsdk:"saas_security"`
	Spyware                  []types.String `tfsdk:"spyware"`
	UrlFiltering             []types.String `tfsdk:"url_filtering"`
	VirusAndWildfireAnalysis []types.String `tfsdk:"virus_and_wildfire_analysis"`
	Vulnerability            []types.String `tfsdk:"vulnerability"`
}

// Metadata returns the data source type name.
func (d *profileGroupsListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_profile_groups_list"
}

// Schema defines the schema for this listing data source.
func (d *profileGroupsListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Mobile Users Container", "Mobile Users Explicit Proxy"),
				},
			},

			// Output.
			"data": dsschema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"dns_security": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"file_blocking": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"object_id": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"saas_security": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"spyware": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"url_filtering": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"virus_and_wildfire_analysis": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"vulnerability": dsschema.ListAttribute{
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
func (d *profileGroupsListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *profileGroupsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state profileGroupsListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_profile_groups_list",
		"limit":            state.Limit.ValueInt64(),
		"has_limit":        !state.Limit.IsNull(),
		"offset":           state.Offset.ValueInt64(),
		"has_offset":       !state.Offset.IsNull(),
		"name":             state.Name.ValueString(),
		"has_name":         !state.Name.IsNull(),
		"folder":           state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := jeahrQe.NewClient(d.client)
	input := jeahrQe.ListInput{
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
	var var0 []profileGroupsListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]profileGroupsListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 profileGroupsListDsModelConfig
			var2.DnsSecurity = EncodeStringSlice(var1.DnsSecurity)
			var2.FileBlocking = EncodeStringSlice(var1.FileBlocking)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Name = types.StringValue(var1.Name)
			var2.SaasSecurity = EncodeStringSlice(var1.SaasSecurity)
			var2.Spyware = EncodeStringSlice(var1.Spyware)
			var2.UrlFiltering = EncodeStringSlice(var1.UrlFiltering)
			var2.VirusAndWildfireAnalysis = EncodeStringSlice(var1.VirusAndWildfireAnalysis)
			var2.Vulnerability = EncodeStringSlice(var1.Vulnerability)
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
	_ datasource.DataSource              = &profileGroupsDataSource{}
	_ datasource.DataSourceWithConfigure = &profileGroupsDataSource{}
)

func NewProfileGroupsDataSource() datasource.DataSource {
	return &profileGroupsDataSource{}
}

type profileGroupsDataSource struct {
	client *sase.Client
}

type profileGroupsDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`

	// Output.
	// Ref: #/components/schemas/profile-groups
	DnsSecurity  []types.String `tfsdk:"dns_security"`
	FileBlocking []types.String `tfsdk:"file_blocking"`
	// input omit: ObjectId
	Name                     types.String   `tfsdk:"name"`
	SaasSecurity             []types.String `tfsdk:"saas_security"`
	Spyware                  []types.String `tfsdk:"spyware"`
	UrlFiltering             []types.String `tfsdk:"url_filtering"`
	VirusAndWildfireAnalysis []types.String `tfsdk:"virus_and_wildfire_analysis"`
	Vulnerability            []types.String `tfsdk:"vulnerability"`
}

// Metadata returns the data source type name.
func (d *profileGroupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_profile_groups"
}

// Schema defines the schema for this listing data source.
func (d *profileGroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"dns_security": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"file_blocking": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"saas_security": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"spyware": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"url_filtering": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"virus_and_wildfire_analysis": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"vulnerability": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure prepares the struct.
func (d *profileGroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *profileGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state profileGroupsDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_profile_groups",
		"object_id":        state.ObjectId.ValueString(),
	})

	// Prepare to run the command.
	svc := jeahrQe.NewClient(d.client)
	input := jeahrQe.ReadInput{
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
	state.DnsSecurity = EncodeStringSlice(ans.DnsSecurity)
	state.FileBlocking = EncodeStringSlice(ans.FileBlocking)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.SaasSecurity = EncodeStringSlice(ans.SaasSecurity)
	state.Spyware = EncodeStringSlice(ans.Spyware)
	state.UrlFiltering = EncodeStringSlice(ans.UrlFiltering)
	state.VirusAndWildfireAnalysis = EncodeStringSlice(ans.VirusAndWildfireAnalysis)
	state.Vulnerability = EncodeStringSlice(ans.Vulnerability)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &profileGroupsResource{}
	_ resource.ResourceWithConfigure   = &profileGroupsResource{}
	_ resource.ResourceWithImportState = &profileGroupsResource{}
)

func NewProfileGroupsResource() resource.Resource {
	return &profileGroupsResource{}
}

type profileGroupsResource struct {
	client *sase.Client
}

type profileGroupsRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/profile-groups
	DnsSecurity              []types.String `tfsdk:"dns_security"`
	FileBlocking             []types.String `tfsdk:"file_blocking"`
	ObjectId                 types.String   `tfsdk:"object_id"`
	Name                     types.String   `tfsdk:"name"`
	SaasSecurity             []types.String `tfsdk:"saas_security"`
	Spyware                  []types.String `tfsdk:"spyware"`
	UrlFiltering             []types.String `tfsdk:"url_filtering"`
	VirusAndWildfireAnalysis []types.String `tfsdk:"virus_and_wildfire_analysis"`
	Vulnerability            []types.String `tfsdk:"vulnerability"`
}

// Metadata returns the data source type name.
func (r *profileGroupsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_profile_groups"
}

// Schema defines the schema for this listing data source.
func (r *profileGroupsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Mobile Users Container", "Mobile Users Explicit Proxy"),
				},
			},

			"dns_security": rsschema.ListAttribute{
				Description: "",
				Optional:    true,
				ElementType: types.StringType,
			},
			"file_blocking": rsschema.ListAttribute{
				Description: "",
				Optional:    true,
				ElementType: types.StringType,
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
			},
			"saas_security": rsschema.ListAttribute{
				Description: "",
				Optional:    true,
				ElementType: types.StringType,
			},
			"spyware": rsschema.ListAttribute{
				Description: "",
				Optional:    true,
				ElementType: types.StringType,
			},
			"url_filtering": rsschema.ListAttribute{
				Description: "",
				Optional:    true,
				ElementType: types.StringType,
			},
			"virus_and_wildfire_analysis": rsschema.ListAttribute{
				Description: "",
				Optional:    true,
				ElementType: types.StringType,
			},
			"vulnerability": rsschema.ListAttribute{
				Description: "",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure prepares the struct.
func (r *profileGroupsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *profileGroupsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state profileGroupsRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "sase_profile_groups",
		"folder":        state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := jeahrQe.NewClient(r.client)
	input := jeahrQe.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 mQehIbG.Config
	var0.DnsSecurity = DecodeStringSlice(state.DnsSecurity)
	var0.FileBlocking = DecodeStringSlice(state.FileBlocking)
	var0.Name = state.Name.ValueString()
	var0.SaasSecurity = DecodeStringSlice(state.SaasSecurity)
	var0.Spyware = DecodeStringSlice(state.Spyware)
	var0.UrlFiltering = DecodeStringSlice(state.UrlFiltering)
	var0.VirusAndWildfireAnalysis = DecodeStringSlice(state.VirusAndWildfireAnalysis)
	var0.Vulnerability = DecodeStringSlice(state.Vulnerability)
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
	state.DnsSecurity = EncodeStringSlice(ans.DnsSecurity)
	state.FileBlocking = EncodeStringSlice(ans.FileBlocking)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.SaasSecurity = EncodeStringSlice(ans.SaasSecurity)
	state.Spyware = EncodeStringSlice(ans.Spyware)
	state.UrlFiltering = EncodeStringSlice(ans.UrlFiltering)
	state.VirusAndWildfireAnalysis = EncodeStringSlice(ans.VirusAndWildfireAnalysis)
	state.Vulnerability = EncodeStringSlice(ans.Vulnerability)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *profileGroupsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state profileGroupsRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "sase_profile_groups",
		"locMap":        map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":        tokens,
	})

	// Prepare to read the config.
	svc := jeahrQe.NewClient(r.client)
	input := jeahrQe.ReadInput{
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
	state.DnsSecurity = EncodeStringSlice(ans.DnsSecurity)
	state.FileBlocking = EncodeStringSlice(ans.FileBlocking)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.SaasSecurity = EncodeStringSlice(ans.SaasSecurity)
	state.Spyware = EncodeStringSlice(ans.Spyware)
	state.UrlFiltering = EncodeStringSlice(ans.UrlFiltering)
	state.VirusAndWildfireAnalysis = EncodeStringSlice(ans.VirusAndWildfireAnalysis)
	state.Vulnerability = EncodeStringSlice(ans.Vulnerability)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *profileGroupsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state profileGroupsRsModel
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
		"resource_name": "sase_profile_groups",
		"object_id":     state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := jeahrQe.NewClient(r.client)
	input := jeahrQe.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 mQehIbG.Config
	var0.DnsSecurity = DecodeStringSlice(plan.DnsSecurity)
	var0.FileBlocking = DecodeStringSlice(plan.FileBlocking)
	var0.Name = plan.Name.ValueString()
	var0.SaasSecurity = DecodeStringSlice(plan.SaasSecurity)
	var0.Spyware = DecodeStringSlice(plan.Spyware)
	var0.UrlFiltering = DecodeStringSlice(plan.UrlFiltering)
	var0.VirusAndWildfireAnalysis = DecodeStringSlice(plan.VirusAndWildfireAnalysis)
	var0.Vulnerability = DecodeStringSlice(plan.Vulnerability)
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	state.DnsSecurity = EncodeStringSlice(ans.DnsSecurity)
	state.FileBlocking = EncodeStringSlice(ans.FileBlocking)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.SaasSecurity = EncodeStringSlice(ans.SaasSecurity)
	state.Spyware = EncodeStringSlice(ans.Spyware)
	state.UrlFiltering = EncodeStringSlice(ans.UrlFiltering)
	state.VirusAndWildfireAnalysis = EncodeStringSlice(ans.VirusAndWildfireAnalysis)
	state.Vulnerability = EncodeStringSlice(ans.Vulnerability)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *profileGroupsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name": "sase_profile_groups",
		"locMap":        map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":        tokens,
	})

	svc := jeahrQe.NewClient(r.client)
	input := jeahrQe.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *profileGroupsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
