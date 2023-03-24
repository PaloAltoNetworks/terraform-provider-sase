package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	mtKjwYt "github.com/paloaltonetworks/sase-go/netsec/schema/wildfire/anti/virus/profiles"
	crXhgow "github.com/paloaltonetworks/sase-go/netsec/service/v1/wildfireantivirusprofiles"

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
	_ datasource.DataSource              = &wildfireAntiVirusProfilesListDataSource{}
	_ datasource.DataSourceWithConfigure = &wildfireAntiVirusProfilesListDataSource{}
)

func NewWildfireAntiVirusProfilesListDataSource() datasource.DataSource {
	return &wildfireAntiVirusProfilesListDataSource{}
}

type wildfireAntiVirusProfilesListDataSource struct {
	client *sase.Client
}

type wildfireAntiVirusProfilesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []wildfireAntiVirusProfilesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type wildfireAntiVirusProfilesListDsModelConfig struct {
	Description     types.String                                                `tfsdk:"description"`
	ObjectId        types.String                                                `tfsdk:"object_id"`
	MlavException   []wildfireAntiVirusProfilesListDsModelMlavExceptionObject   `tfsdk:"mlav_exception"`
	Name            types.String                                                `tfsdk:"name"`
	PacketCapture   types.Bool                                                  `tfsdk:"packet_capture"`
	Rules           []wildfireAntiVirusProfilesListDsModelRulesObject           `tfsdk:"rules"`
	ThreatException []wildfireAntiVirusProfilesListDsModelThreatExceptionObject `tfsdk:"threat_exception"`
}

type wildfireAntiVirusProfilesListDsModelMlavExceptionObject struct {
	Description types.String `tfsdk:"description"`
	Filename    types.String `tfsdk:"filename"`
	Name        types.String `tfsdk:"name"`
}

type wildfireAntiVirusProfilesListDsModelRulesObject struct {
	Analysis    types.String   `tfsdk:"analysis"`
	Application []types.String `tfsdk:"application"`
	Direction   types.String   `tfsdk:"direction"`
	FileType    []types.String `tfsdk:"file_type"`
	Name        types.String   `tfsdk:"name"`
}

type wildfireAntiVirusProfilesListDsModelThreatExceptionObject struct {
	Name  types.String `tfsdk:"name"`
	Notes types.String `tfsdk:"notes"`
}

// Metadata returns the data source type name.
func (d *wildfireAntiVirusProfilesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wildfire_anti_virus_profiles_list"
}

// Schema defines the schema for this listing data source.
func (d *wildfireAntiVirusProfilesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"mlav_exception": dsschema.ListNestedAttribute{
							Description: "",
							Computed:    true,
							NestedObject: dsschema.NestedAttributeObject{
								Attributes: map[string]dsschema.Attribute{
									"description": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"filename": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"name": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
						},
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"packet_capture": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"rules": dsschema.ListNestedAttribute{
							Description: "",
							Computed:    true,
							NestedObject: dsschema.NestedAttributeObject{
								Attributes: map[string]dsschema.Attribute{
									"analysis": dsschema.StringAttribute{
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
						"threat_exception": dsschema.ListNestedAttribute{
							Description: "",
							Computed:    true,
							NestedObject: dsschema.NestedAttributeObject{
								Attributes: map[string]dsschema.Attribute{
									"name": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"notes": dsschema.StringAttribute{
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
func (d *wildfireAntiVirusProfilesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *wildfireAntiVirusProfilesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state wildfireAntiVirusProfilesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_wildfire_anti_virus_profiles_list",
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
	svc := crXhgow.NewClient(d.client)
	input := crXhgow.ListInput{
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
	var var0 []wildfireAntiVirusProfilesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]wildfireAntiVirusProfilesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 wildfireAntiVirusProfilesListDsModelConfig
			var var3 []wildfireAntiVirusProfilesListDsModelMlavExceptionObject
			if len(var1.MlavException) != 0 {
				var3 = make([]wildfireAntiVirusProfilesListDsModelMlavExceptionObject, 0, len(var1.MlavException))
				for var4Index := range var1.MlavException {
					var4 := var1.MlavException[var4Index]
					var var5 wildfireAntiVirusProfilesListDsModelMlavExceptionObject
					var5.Description = types.StringValue(var4.Description)
					var5.Filename = types.StringValue(var4.Filename)
					var5.Name = types.StringValue(var4.Name)
					var3 = append(var3, var5)
				}
			}
			var var6 []wildfireAntiVirusProfilesListDsModelRulesObject
			if len(var1.Rules) != 0 {
				var6 = make([]wildfireAntiVirusProfilesListDsModelRulesObject, 0, len(var1.Rules))
				for var7Index := range var1.Rules {
					var7 := var1.Rules[var7Index]
					var var8 wildfireAntiVirusProfilesListDsModelRulesObject
					var8.Analysis = types.StringValue(var7.Analysis)
					var8.Application = EncodeStringSlice(var7.Application)
					var8.Direction = types.StringValue(var7.Direction)
					var8.FileType = EncodeStringSlice(var7.FileType)
					var8.Name = types.StringValue(var7.Name)
					var6 = append(var6, var8)
				}
			}
			var var9 []wildfireAntiVirusProfilesListDsModelThreatExceptionObject
			if len(var1.ThreatException) != 0 {
				var9 = make([]wildfireAntiVirusProfilesListDsModelThreatExceptionObject, 0, len(var1.ThreatException))
				for var10Index := range var1.ThreatException {
					var10 := var1.ThreatException[var10Index]
					var var11 wildfireAntiVirusProfilesListDsModelThreatExceptionObject
					var11.Name = types.StringValue(var10.Name)
					var11.Notes = types.StringValue(var10.Notes)
					var9 = append(var9, var11)
				}
			}
			var2.Description = types.StringValue(var1.Description)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.MlavException = var3
			var2.Name = types.StringValue(var1.Name)
			var2.PacketCapture = types.BoolValue(var1.PacketCapture)
			var2.Rules = var6
			var2.ThreatException = var9
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
	_ datasource.DataSource              = &wildfireAntiVirusProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &wildfireAntiVirusProfilesDataSource{}
)

func NewWildfireAntiVirusProfilesDataSource() datasource.DataSource {
	return &wildfireAntiVirusProfilesDataSource{}
}

type wildfireAntiVirusProfilesDataSource struct {
	client *sase.Client
}

type wildfireAntiVirusProfilesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`
	Folder   types.String `tfsdk:"folder"`

	// Output.
	// Ref: #/components/schemas/wildfire-anti-virus-profiles
	Description types.String `tfsdk:"description"`
	// input omit: ObjectId
	MlavException   []wildfireAntiVirusProfilesDsModelMlavExceptionObject   `tfsdk:"mlav_exception"`
	Name            types.String                                            `tfsdk:"name"`
	PacketCapture   types.Bool                                              `tfsdk:"packet_capture"`
	Rules           []wildfireAntiVirusProfilesDsModelRulesObject           `tfsdk:"rules"`
	ThreatException []wildfireAntiVirusProfilesDsModelThreatExceptionObject `tfsdk:"threat_exception"`
}

type wildfireAntiVirusProfilesDsModelMlavExceptionObject struct {
	Description types.String `tfsdk:"description"`
	Filename    types.String `tfsdk:"filename"`
	Name        types.String `tfsdk:"name"`
}

type wildfireAntiVirusProfilesDsModelRulesObject struct {
	Analysis    types.String   `tfsdk:"analysis"`
	Application []types.String `tfsdk:"application"`
	Direction   types.String   `tfsdk:"direction"`
	FileType    []types.String `tfsdk:"file_type"`
	Name        types.String   `tfsdk:"name"`
}

type wildfireAntiVirusProfilesDsModelThreatExceptionObject struct {
	Name  types.String `tfsdk:"name"`
	Notes types.String `tfsdk:"notes"`
}

// Metadata returns the data source type name.
func (d *wildfireAntiVirusProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wildfire_anti_virus_profiles"
}

// Schema defines the schema for this listing data source.
func (d *wildfireAntiVirusProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"mlav_exception": dsschema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"description": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"filename": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
					},
				},
			},
			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"packet_capture": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"rules": dsschema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"analysis": dsschema.StringAttribute{
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
			"threat_exception": dsschema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"notes": dsschema.StringAttribute{
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
func (d *wildfireAntiVirusProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *wildfireAntiVirusProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state wildfireAntiVirusProfilesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_wildfire_anti_virus_profiles",
		"object_id":                   state.ObjectId.ValueString(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := crXhgow.NewClient(d.client)
	input := crXhgow.ReadInput{
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
	var var0 []wildfireAntiVirusProfilesDsModelMlavExceptionObject
	if len(ans.MlavException) != 0 {
		var0 = make([]wildfireAntiVirusProfilesDsModelMlavExceptionObject, 0, len(ans.MlavException))
		for var1Index := range ans.MlavException {
			var1 := ans.MlavException[var1Index]
			var var2 wildfireAntiVirusProfilesDsModelMlavExceptionObject
			var2.Description = types.StringValue(var1.Description)
			var2.Filename = types.StringValue(var1.Filename)
			var2.Name = types.StringValue(var1.Name)
			var0 = append(var0, var2)
		}
	}
	var var3 []wildfireAntiVirusProfilesDsModelRulesObject
	if len(ans.Rules) != 0 {
		var3 = make([]wildfireAntiVirusProfilesDsModelRulesObject, 0, len(ans.Rules))
		for var4Index := range ans.Rules {
			var4 := ans.Rules[var4Index]
			var var5 wildfireAntiVirusProfilesDsModelRulesObject
			var5.Analysis = types.StringValue(var4.Analysis)
			var5.Application = EncodeStringSlice(var4.Application)
			var5.Direction = types.StringValue(var4.Direction)
			var5.FileType = EncodeStringSlice(var4.FileType)
			var5.Name = types.StringValue(var4.Name)
			var3 = append(var3, var5)
		}
	}
	var var6 []wildfireAntiVirusProfilesDsModelThreatExceptionObject
	if len(ans.ThreatException) != 0 {
		var6 = make([]wildfireAntiVirusProfilesDsModelThreatExceptionObject, 0, len(ans.ThreatException))
		for var7Index := range ans.ThreatException {
			var7 := ans.ThreatException[var7Index]
			var var8 wildfireAntiVirusProfilesDsModelThreatExceptionObject
			var8.Name = types.StringValue(var7.Name)
			var8.Notes = types.StringValue(var7.Notes)
			var6 = append(var6, var8)
		}
	}
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.MlavException = var0
	state.Name = types.StringValue(ans.Name)
	state.PacketCapture = types.BoolValue(ans.PacketCapture)
	state.Rules = var3
	state.ThreatException = var6

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &wildfireAntiVirusProfilesResource{}
	_ resource.ResourceWithConfigure   = &wildfireAntiVirusProfilesResource{}
	_ resource.ResourceWithImportState = &wildfireAntiVirusProfilesResource{}
)

func NewWildfireAntiVirusProfilesResource() resource.Resource {
	return &wildfireAntiVirusProfilesResource{}
}

type wildfireAntiVirusProfilesResource struct {
	client *sase.Client
}

type wildfireAntiVirusProfilesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/wildfire-anti-virus-profiles
	Description     types.String                                            `tfsdk:"description"`
	ObjectId        types.String                                            `tfsdk:"object_id"`
	MlavException   []wildfireAntiVirusProfilesRsModelMlavExceptionObject   `tfsdk:"mlav_exception"`
	Name            types.String                                            `tfsdk:"name"`
	PacketCapture   types.Bool                                              `tfsdk:"packet_capture"`
	Rules           []wildfireAntiVirusProfilesRsModelRulesObject           `tfsdk:"rules"`
	ThreatException []wildfireAntiVirusProfilesRsModelThreatExceptionObject `tfsdk:"threat_exception"`
}

type wildfireAntiVirusProfilesRsModelMlavExceptionObject struct {
	Description types.String `tfsdk:"description"`
	Filename    types.String `tfsdk:"filename"`
	Name        types.String `tfsdk:"name"`
}

type wildfireAntiVirusProfilesRsModelRulesObject struct {
	Analysis    types.String   `tfsdk:"analysis"`
	Application []types.String `tfsdk:"application"`
	Direction   types.String   `tfsdk:"direction"`
	FileType    []types.String `tfsdk:"file_type"`
	Name        types.String   `tfsdk:"name"`
}

type wildfireAntiVirusProfilesRsModelThreatExceptionObject struct {
	Name  types.String `tfsdk:"name"`
	Notes types.String `tfsdk:"notes"`
}

// Metadata returns the data source type name.
func (r *wildfireAntiVirusProfilesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wildfire_anti_virus_profiles"
}

// Schema defines the schema for this listing data source.
func (r *wildfireAntiVirusProfilesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"mlav_exception": rsschema.ListNestedAttribute{
				Description: "",
				Optional:    true,
				NestedObject: rsschema.NestedAttributeObject{
					Attributes: map[string]rsschema.Attribute{
						"description": rsschema.StringAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
						},
						"filename": rsschema.StringAttribute{
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
					},
				},
			},
			"name": rsschema.StringAttribute{
				Description: "",
				Required:    true,
			},
			"packet_capture": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"rules": rsschema.ListNestedAttribute{
				Description: "",
				Optional:    true,
				NestedObject: rsschema.NestedAttributeObject{
					Attributes: map[string]rsschema.Attribute{
						"analysis": rsschema.StringAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
							Validators: []validator.String{
								stringvalidator.OneOf("public-cloud", "private-cloud"),
							},
						},
						"application": rsschema.ListAttribute{
							Description: "",
							Optional:    true,
							ElementType: types.StringType,
						},
						"direction": rsschema.StringAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
							Validators: []validator.String{
								stringvalidator.OneOf("download", "upload", "both"),
							},
						},
						"file_type": rsschema.ListAttribute{
							Description: "",
							Optional:    true,
							ElementType: types.StringType,
						},
						"name": rsschema.StringAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
						},
					},
				},
			},
			"threat_exception": rsschema.ListNestedAttribute{
				Description: "",
				Optional:    true,
				NestedObject: rsschema.NestedAttributeObject{
					Attributes: map[string]rsschema.Attribute{
						"name": rsschema.StringAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
						},
						"notes": rsschema.StringAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
						},
					},
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *wildfireAntiVirusProfilesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *wildfireAntiVirusProfilesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state wildfireAntiVirusProfilesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_wildfire_anti_virus_profiles",
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := crXhgow.NewClient(r.client)
	input := crXhgow.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 mtKjwYt.Config
	var0.Description = state.Description.ValueString()
	var var1 []mtKjwYt.MlavExceptionObject
	if len(state.MlavException) != 0 {
		var1 = make([]mtKjwYt.MlavExceptionObject, 0, len(state.MlavException))
		for var2Index := range state.MlavException {
			var2 := state.MlavException[var2Index]
			var var3 mtKjwYt.MlavExceptionObject
			var3.Description = var2.Description.ValueString()
			var3.Filename = var2.Filename.ValueString()
			var3.Name = var2.Name.ValueString()
			var1 = append(var1, var3)
		}
	}
	var0.MlavException = var1
	var0.Name = state.Name.ValueString()
	var0.PacketCapture = state.PacketCapture.ValueBool()
	var var4 []mtKjwYt.RulesObject
	if len(state.Rules) != 0 {
		var4 = make([]mtKjwYt.RulesObject, 0, len(state.Rules))
		for var5Index := range state.Rules {
			var5 := state.Rules[var5Index]
			var var6 mtKjwYt.RulesObject
			var6.Analysis = var5.Analysis.ValueString()
			var6.Application = DecodeStringSlice(var5.Application)
			var6.Direction = var5.Direction.ValueString()
			var6.FileType = DecodeStringSlice(var5.FileType)
			var6.Name = var5.Name.ValueString()
			var4 = append(var4, var6)
		}
	}
	var0.Rules = var4
	var var7 []mtKjwYt.ThreatExceptionObject
	if len(state.ThreatException) != 0 {
		var7 = make([]mtKjwYt.ThreatExceptionObject, 0, len(state.ThreatException))
		for var8Index := range state.ThreatException {
			var8 := state.ThreatException[var8Index]
			var var9 mtKjwYt.ThreatExceptionObject
			var9.Name = var8.Name.ValueString()
			var9.Notes = var8.Notes.ValueString()
			var7 = append(var7, var9)
		}
	}
	var0.ThreatException = var7
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
	var var10 []wildfireAntiVirusProfilesRsModelMlavExceptionObject
	if len(ans.MlavException) != 0 {
		var10 = make([]wildfireAntiVirusProfilesRsModelMlavExceptionObject, 0, len(ans.MlavException))
		for var11Index := range ans.MlavException {
			var11 := ans.MlavException[var11Index]
			var var12 wildfireAntiVirusProfilesRsModelMlavExceptionObject
			var12.Description = types.StringValue(var11.Description)
			var12.Filename = types.StringValue(var11.Filename)
			var12.Name = types.StringValue(var11.Name)
			var10 = append(var10, var12)
		}
	}
	var var13 []wildfireAntiVirusProfilesRsModelRulesObject
	if len(ans.Rules) != 0 {
		var13 = make([]wildfireAntiVirusProfilesRsModelRulesObject, 0, len(ans.Rules))
		for var14Index := range ans.Rules {
			var14 := ans.Rules[var14Index]
			var var15 wildfireAntiVirusProfilesRsModelRulesObject
			var15.Analysis = types.StringValue(var14.Analysis)
			var15.Application = EncodeStringSlice(var14.Application)
			var15.Direction = types.StringValue(var14.Direction)
			var15.FileType = EncodeStringSlice(var14.FileType)
			var15.Name = types.StringValue(var14.Name)
			var13 = append(var13, var15)
		}
	}
	var var16 []wildfireAntiVirusProfilesRsModelThreatExceptionObject
	if len(ans.ThreatException) != 0 {
		var16 = make([]wildfireAntiVirusProfilesRsModelThreatExceptionObject, 0, len(ans.ThreatException))
		for var17Index := range ans.ThreatException {
			var17 := ans.ThreatException[var17Index]
			var var18 wildfireAntiVirusProfilesRsModelThreatExceptionObject
			var18.Name = types.StringValue(var17.Name)
			var18.Notes = types.StringValue(var17.Notes)
			var16 = append(var16, var18)
		}
	}
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.MlavException = var10
	state.Name = types.StringValue(ans.Name)
	state.PacketCapture = types.BoolValue(ans.PacketCapture)
	state.Rules = var13
	state.ThreatException = var16

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *wildfireAntiVirusProfilesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state wildfireAntiVirusProfilesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_wildfire_anti_virus_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := crXhgow.NewClient(r.client)
	input := crXhgow.ReadInput{
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
	var var0 []wildfireAntiVirusProfilesRsModelMlavExceptionObject
	if len(ans.MlavException) != 0 {
		var0 = make([]wildfireAntiVirusProfilesRsModelMlavExceptionObject, 0, len(ans.MlavException))
		for var1Index := range ans.MlavException {
			var1 := ans.MlavException[var1Index]
			var var2 wildfireAntiVirusProfilesRsModelMlavExceptionObject
			var2.Description = types.StringValue(var1.Description)
			var2.Filename = types.StringValue(var1.Filename)
			var2.Name = types.StringValue(var1.Name)
			var0 = append(var0, var2)
		}
	}
	var var3 []wildfireAntiVirusProfilesRsModelRulesObject
	if len(ans.Rules) != 0 {
		var3 = make([]wildfireAntiVirusProfilesRsModelRulesObject, 0, len(ans.Rules))
		for var4Index := range ans.Rules {
			var4 := ans.Rules[var4Index]
			var var5 wildfireAntiVirusProfilesRsModelRulesObject
			var5.Analysis = types.StringValue(var4.Analysis)
			var5.Application = EncodeStringSlice(var4.Application)
			var5.Direction = types.StringValue(var4.Direction)
			var5.FileType = EncodeStringSlice(var4.FileType)
			var5.Name = types.StringValue(var4.Name)
			var3 = append(var3, var5)
		}
	}
	var var6 []wildfireAntiVirusProfilesRsModelThreatExceptionObject
	if len(ans.ThreatException) != 0 {
		var6 = make([]wildfireAntiVirusProfilesRsModelThreatExceptionObject, 0, len(ans.ThreatException))
		for var7Index := range ans.ThreatException {
			var7 := ans.ThreatException[var7Index]
			var var8 wildfireAntiVirusProfilesRsModelThreatExceptionObject
			var8.Name = types.StringValue(var7.Name)
			var8.Notes = types.StringValue(var7.Notes)
			var6 = append(var6, var8)
		}
	}
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.MlavException = var0
	state.Name = types.StringValue(ans.Name)
	state.PacketCapture = types.BoolValue(ans.PacketCapture)
	state.Rules = var3
	state.ThreatException = var6

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *wildfireAntiVirusProfilesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state wildfireAntiVirusProfilesRsModel
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
		"resource_name":               "sase_wildfire_anti_virus_profiles",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := crXhgow.NewClient(r.client)
	input := crXhgow.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 mtKjwYt.Config
	var0.Description = plan.Description.ValueString()
	var var1 []mtKjwYt.MlavExceptionObject
	if len(plan.MlavException) != 0 {
		var1 = make([]mtKjwYt.MlavExceptionObject, 0, len(plan.MlavException))
		for var2Index := range plan.MlavException {
			var2 := plan.MlavException[var2Index]
			var var3 mtKjwYt.MlavExceptionObject
			var3.Description = var2.Description.ValueString()
			var3.Filename = var2.Filename.ValueString()
			var3.Name = var2.Name.ValueString()
			var1 = append(var1, var3)
		}
	}
	var0.MlavException = var1
	var0.Name = plan.Name.ValueString()
	var0.PacketCapture = plan.PacketCapture.ValueBool()
	var var4 []mtKjwYt.RulesObject
	if len(plan.Rules) != 0 {
		var4 = make([]mtKjwYt.RulesObject, 0, len(plan.Rules))
		for var5Index := range plan.Rules {
			var5 := plan.Rules[var5Index]
			var var6 mtKjwYt.RulesObject
			var6.Analysis = var5.Analysis.ValueString()
			var6.Application = DecodeStringSlice(var5.Application)
			var6.Direction = var5.Direction.ValueString()
			var6.FileType = DecodeStringSlice(var5.FileType)
			var6.Name = var5.Name.ValueString()
			var4 = append(var4, var6)
		}
	}
	var0.Rules = var4
	var var7 []mtKjwYt.ThreatExceptionObject
	if len(plan.ThreatException) != 0 {
		var7 = make([]mtKjwYt.ThreatExceptionObject, 0, len(plan.ThreatException))
		for var8Index := range plan.ThreatException {
			var8 := plan.ThreatException[var8Index]
			var var9 mtKjwYt.ThreatExceptionObject
			var9.Name = var8.Name.ValueString()
			var9.Notes = var8.Notes.ValueString()
			var7 = append(var7, var9)
		}
	}
	var0.ThreatException = var7
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var10 []wildfireAntiVirusProfilesRsModelMlavExceptionObject
	if len(ans.MlavException) != 0 {
		var10 = make([]wildfireAntiVirusProfilesRsModelMlavExceptionObject, 0, len(ans.MlavException))
		for var11Index := range ans.MlavException {
			var11 := ans.MlavException[var11Index]
			var var12 wildfireAntiVirusProfilesRsModelMlavExceptionObject
			var12.Description = types.StringValue(var11.Description)
			var12.Filename = types.StringValue(var11.Filename)
			var12.Name = types.StringValue(var11.Name)
			var10 = append(var10, var12)
		}
	}
	var var13 []wildfireAntiVirusProfilesRsModelRulesObject
	if len(ans.Rules) != 0 {
		var13 = make([]wildfireAntiVirusProfilesRsModelRulesObject, 0, len(ans.Rules))
		for var14Index := range ans.Rules {
			var14 := ans.Rules[var14Index]
			var var15 wildfireAntiVirusProfilesRsModelRulesObject
			var15.Analysis = types.StringValue(var14.Analysis)
			var15.Application = EncodeStringSlice(var14.Application)
			var15.Direction = types.StringValue(var14.Direction)
			var15.FileType = EncodeStringSlice(var14.FileType)
			var15.Name = types.StringValue(var14.Name)
			var13 = append(var13, var15)
		}
	}
	var var16 []wildfireAntiVirusProfilesRsModelThreatExceptionObject
	if len(ans.ThreatException) != 0 {
		var16 = make([]wildfireAntiVirusProfilesRsModelThreatExceptionObject, 0, len(ans.ThreatException))
		for var17Index := range ans.ThreatException {
			var17 := ans.ThreatException[var17Index]
			var var18 wildfireAntiVirusProfilesRsModelThreatExceptionObject
			var18.Name = types.StringValue(var17.Name)
			var18.Notes = types.StringValue(var17.Notes)
			var16 = append(var16, var18)
		}
	}
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.MlavException = var10
	state.Name = types.StringValue(ans.Name)
	state.PacketCapture = types.BoolValue(ans.PacketCapture)
	state.Rules = var13
	state.ThreatException = var16

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *wildfireAntiVirusProfilesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name":               "sase_wildfire_anti_virus_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	svc := crXhgow.NewClient(r.client)
	input := crXhgow.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *wildfireAntiVirusProfilesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
