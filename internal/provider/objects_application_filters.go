package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	lhPcfTR "github.com/paloaltonetworks/sase-go/netsec/schema/objects/application/filters"
	jHKNPjP "github.com/paloaltonetworks/sase-go/netsec/service/v1/applicationfilters"

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
	_ datasource.DataSource              = &objectsApplicationFiltersListDataSource{}
	_ datasource.DataSourceWithConfigure = &objectsApplicationFiltersListDataSource{}
)

func NewObjectsApplicationFiltersListDataSource() datasource.DataSource {
	return &objectsApplicationFiltersListDataSource{}
}

type objectsApplicationFiltersListDataSource struct {
	client *sase.Client
}

type objectsApplicationFiltersListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []objectsApplicationFiltersListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type objectsApplicationFiltersListDsModelConfig struct {
	Category                []types.String                                     `tfsdk:"category"`
	Evasive                 types.Bool                                         `tfsdk:"evasive"`
	ExcessiveBandwidthUse   types.Bool                                         `tfsdk:"excessive_bandwidth_use"`
	Exclude                 []types.String                                     `tfsdk:"exclude"`
	HasKnownVulnerabilities types.Bool                                         `tfsdk:"has_known_vulnerabilities"`
	ObjectId                types.String                                       `tfsdk:"object_id"`
	IsSaas                  types.Bool                                         `tfsdk:"is_saas"`
	Name                    types.String                                       `tfsdk:"name"`
	NewAppid                types.Bool                                         `tfsdk:"new_appid"`
	Pervasive               types.Bool                                         `tfsdk:"pervasive"`
	ProneToMisuse           types.Bool                                         `tfsdk:"prone_to_misuse"`
	Risk                    []types.Int64                                      `tfsdk:"risk"`
	SaasCertifications      []types.String                                     `tfsdk:"saas_certifications"`
	SaasRisk                []types.String                                     `tfsdk:"saas_risk"`
	Subcategory             []types.String                                     `tfsdk:"subcategory"`
	Tagging                 *objectsApplicationFiltersListDsModelTaggingObject `tfsdk:"tagging"`
	Technology              []types.String                                     `tfsdk:"technology"`
	TransfersFiles          types.Bool                                         `tfsdk:"transfers_files"`
	TunnelsOtherApps        types.Bool                                         `tfsdk:"tunnels_other_apps"`
	UsedByMalware           types.Bool                                         `tfsdk:"used_by_malware"`
}

type objectsApplicationFiltersListDsModelTaggingObject struct {
	NoTag types.Bool     `tfsdk:"no_tag"`
	Tag   []types.String `tfsdk:"tag"`
}

// Metadata returns the data source type name.
func (d *objectsApplicationFiltersListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_application_filters_list"
}

// Schema defines the schema for this listing data source.
func (d *objectsApplicationFiltersListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"category": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"evasive": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"excessive_bandwidth_use": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"exclude": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"has_known_vulnerabilities": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"object_id": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"is_saas": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"new_appid": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"pervasive": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"prone_to_misuse": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"risk": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.Int64Type,
						},
						"saas_certifications": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"saas_risk": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"subcategory": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"tagging": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"no_tag": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"tag": dsschema.ListAttribute{
									Description: "",
									Computed:    true,
									ElementType: types.StringType,
								},
							},
						},
						"technology": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"transfers_files": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"tunnels_other_apps": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"used_by_malware": dsschema.BoolAttribute{
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
func (d *objectsApplicationFiltersListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *objectsApplicationFiltersListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state objectsApplicationFiltersListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_objects_application_filters_list",
		"limit":            state.Limit.ValueInt64(),
		"has_limit":        !state.Limit.IsNull(),
		"offset":           state.Offset.ValueInt64(),
		"has_offset":       !state.Offset.IsNull(),
		"name":             state.Name.ValueString(),
		"has_name":         !state.Name.IsNull(),
		"folder":           state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := jHKNPjP.NewClient(d.client)
	input := jHKNPjP.ListInput{
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
	state.Id = types.StringValue(strings.Join([]string{strconv.FormatInt(*input.Limit, 10), strconv.FormatInt(*input.Offset, 10), *input.Name, input.Folder}, IdSeparator))
	var var0 []objectsApplicationFiltersListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]objectsApplicationFiltersListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 objectsApplicationFiltersListDsModelConfig
			var var3 *objectsApplicationFiltersListDsModelTaggingObject
			if var1.Tagging != nil {
				var3 = &objectsApplicationFiltersListDsModelTaggingObject{}
				var3.NoTag = types.BoolValue(var1.Tagging.NoTag)
				var3.Tag = EncodeStringSlice(var1.Tagging.Tag)
			}
			var2.Category = EncodeStringSlice(var1.Category)
			var2.Evasive = types.BoolValue(var1.Evasive)
			var2.ExcessiveBandwidthUse = types.BoolValue(var1.ExcessiveBandwidthUse)
			var2.Exclude = EncodeStringSlice(var1.Exclude)
			var2.HasKnownVulnerabilities = types.BoolValue(var1.HasKnownVulnerabilities)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.IsSaas = types.BoolValue(var1.IsSaas)
			var2.Name = types.StringValue(var1.Name)
			var2.NewAppid = types.BoolValue(var1.NewAppid)
			var2.Pervasive = types.BoolValue(var1.Pervasive)
			var2.ProneToMisuse = types.BoolValue(var1.ProneToMisuse)
			var2.Risk = EncodeInt64Slice(var1.Risk)
			var2.SaasCertifications = EncodeStringSlice(var1.SaasCertifications)
			var2.SaasRisk = EncodeStringSlice(var1.SaasRisk)
			var2.Subcategory = EncodeStringSlice(var1.Subcategory)
			var2.Tagging = var3
			var2.Technology = EncodeStringSlice(var1.Technology)
			var2.TransfersFiles = types.BoolValue(var1.TransfersFiles)
			var2.TunnelsOtherApps = types.BoolValue(var1.TunnelsOtherApps)
			var2.UsedByMalware = types.BoolValue(var1.UsedByMalware)
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
	_ datasource.DataSource              = &objectsApplicationFiltersDataSource{}
	_ datasource.DataSourceWithConfigure = &objectsApplicationFiltersDataSource{}
)

func NewObjectsApplicationFiltersDataSource() datasource.DataSource {
	return &objectsApplicationFiltersDataSource{}
}

type objectsApplicationFiltersDataSource struct {
	client *sase.Client
}

type objectsApplicationFiltersDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`

	// Output.
	// Ref: #/components/schemas/objects-application-filters
	Category                []types.String `tfsdk:"category"`
	Evasive                 types.Bool     `tfsdk:"evasive"`
	ExcessiveBandwidthUse   types.Bool     `tfsdk:"excessive_bandwidth_use"`
	Exclude                 []types.String `tfsdk:"exclude"`
	HasKnownVulnerabilities types.Bool     `tfsdk:"has_known_vulnerabilities"`
	// input omit: ObjectId
	IsSaas             types.Bool                                     `tfsdk:"is_saas"`
	Name               types.String                                   `tfsdk:"name"`
	NewAppid           types.Bool                                     `tfsdk:"new_appid"`
	Pervasive          types.Bool                                     `tfsdk:"pervasive"`
	ProneToMisuse      types.Bool                                     `tfsdk:"prone_to_misuse"`
	Risk               []types.Int64                                  `tfsdk:"risk"`
	SaasCertifications []types.String                                 `tfsdk:"saas_certifications"`
	SaasRisk           []types.String                                 `tfsdk:"saas_risk"`
	Subcategory        []types.String                                 `tfsdk:"subcategory"`
	Tagging            *objectsApplicationFiltersDsModelTaggingObject `tfsdk:"tagging"`
	Technology         []types.String                                 `tfsdk:"technology"`
	TransfersFiles     types.Bool                                     `tfsdk:"transfers_files"`
	TunnelsOtherApps   types.Bool                                     `tfsdk:"tunnels_other_apps"`
	UsedByMalware      types.Bool                                     `tfsdk:"used_by_malware"`
}

type objectsApplicationFiltersDsModelTaggingObject struct {
	NoTag types.Bool     `tfsdk:"no_tag"`
	Tag   []types.String `tfsdk:"tag"`
}

// Metadata returns the data source type name.
func (d *objectsApplicationFiltersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_application_filters"
}

// Schema defines the schema for this listing data source.
func (d *objectsApplicationFiltersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"category": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"evasive": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"excessive_bandwidth_use": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"exclude": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"has_known_vulnerabilities": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"is_saas": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"new_appid": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"pervasive": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"prone_to_misuse": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"risk": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"saas_certifications": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"saas_risk": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"subcategory": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"tagging": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"no_tag": dsschema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"tag": dsschema.ListAttribute{
						Description: "",
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
			"technology": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"transfers_files": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"tunnels_other_apps": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"used_by_malware": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
		},
	}
}

// Configure prepares the struct.
func (d *objectsApplicationFiltersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *objectsApplicationFiltersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state objectsApplicationFiltersDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_objects_application_filters",
		"object_id":        state.ObjectId.ValueString(),
	})

	// Prepare to run the command.
	svc := jHKNPjP.NewClient(d.client)
	input := jHKNPjP.ReadInput{
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
	var var0 *objectsApplicationFiltersDsModelTaggingObject
	if ans.Tagging != nil {
		var0 = &objectsApplicationFiltersDsModelTaggingObject{}
		var0.NoTag = types.BoolValue(ans.Tagging.NoTag)
		var0.Tag = EncodeStringSlice(ans.Tagging.Tag)
	}
	state.Category = EncodeStringSlice(ans.Category)
	state.Evasive = types.BoolValue(ans.Evasive)
	state.ExcessiveBandwidthUse = types.BoolValue(ans.ExcessiveBandwidthUse)
	state.Exclude = EncodeStringSlice(ans.Exclude)
	state.HasKnownVulnerabilities = types.BoolValue(ans.HasKnownVulnerabilities)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.IsSaas = types.BoolValue(ans.IsSaas)
	state.Name = types.StringValue(ans.Name)
	state.NewAppid = types.BoolValue(ans.NewAppid)
	state.Pervasive = types.BoolValue(ans.Pervasive)
	state.ProneToMisuse = types.BoolValue(ans.ProneToMisuse)
	state.Risk = EncodeInt64Slice(ans.Risk)
	state.SaasCertifications = EncodeStringSlice(ans.SaasCertifications)
	state.SaasRisk = EncodeStringSlice(ans.SaasRisk)
	state.Subcategory = EncodeStringSlice(ans.Subcategory)
	state.Tagging = var0
	state.Technology = EncodeStringSlice(ans.Technology)
	state.TransfersFiles = types.BoolValue(ans.TransfersFiles)
	state.TunnelsOtherApps = types.BoolValue(ans.TunnelsOtherApps)
	state.UsedByMalware = types.BoolValue(ans.UsedByMalware)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &objectsApplicationFiltersResource{}
	_ resource.ResourceWithConfigure   = &objectsApplicationFiltersResource{}
	_ resource.ResourceWithImportState = &objectsApplicationFiltersResource{}
)

func NewObjectsApplicationFiltersResource() resource.Resource {
	return &objectsApplicationFiltersResource{}
}

type objectsApplicationFiltersResource struct {
	client *sase.Client
}

type objectsApplicationFiltersRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/objects-application-filters
	Category                []types.String                                 `tfsdk:"category"`
	Evasive                 types.Bool                                     `tfsdk:"evasive"`
	ExcessiveBandwidthUse   types.Bool                                     `tfsdk:"excessive_bandwidth_use"`
	Exclude                 []types.String                                 `tfsdk:"exclude"`
	HasKnownVulnerabilities types.Bool                                     `tfsdk:"has_known_vulnerabilities"`
	ObjectId                types.String                                   `tfsdk:"object_id"`
	IsSaas                  types.Bool                                     `tfsdk:"is_saas"`
	Name                    types.String                                   `tfsdk:"name"`
	NewAppid                types.Bool                                     `tfsdk:"new_appid"`
	Pervasive               types.Bool                                     `tfsdk:"pervasive"`
	ProneToMisuse           types.Bool                                     `tfsdk:"prone_to_misuse"`
	Risk                    []types.Int64                                  `tfsdk:"risk"`
	SaasCertifications      []types.String                                 `tfsdk:"saas_certifications"`
	SaasRisk                []types.String                                 `tfsdk:"saas_risk"`
	Subcategory             []types.String                                 `tfsdk:"subcategory"`
	Tagging                 *objectsApplicationFiltersRsModelTaggingObject `tfsdk:"tagging"`
	Technology              []types.String                                 `tfsdk:"technology"`
	TransfersFiles          types.Bool                                     `tfsdk:"transfers_files"`
	TunnelsOtherApps        types.Bool                                     `tfsdk:"tunnels_other_apps"`
	UsedByMalware           types.Bool                                     `tfsdk:"used_by_malware"`
}

type objectsApplicationFiltersRsModelTaggingObject struct {
	NoTag types.Bool     `tfsdk:"no_tag"`
	Tag   []types.String `tfsdk:"tag"`
}

// Metadata returns the data source type name.
func (r *objectsApplicationFiltersResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_application_filters"
}

// Schema defines the schema for this listing data source.
func (r *objectsApplicationFiltersResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

			"category": rsschema.ListAttribute{
				Description: "",
				Optional:    true,
				ElementType: types.StringType,
			},
			"evasive": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"excessive_bandwidth_use": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"exclude": rsschema.ListAttribute{
				Description: "",
				Optional:    true,
				ElementType: types.StringType,
			},
			"has_known_vulnerabilities": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"object_id": rsschema.StringAttribute{
				Description: "",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_saas": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
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
			"new_appid": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"pervasive": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"prone_to_misuse": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"risk": rsschema.ListAttribute{
				Description: "",
				Optional:    true,
				ElementType: types.Int64Type,
			},
			"saas_certifications": rsschema.ListAttribute{
				Description: "",
				Optional:    true,
				ElementType: types.StringType,
			},
			"saas_risk": rsschema.ListAttribute{
				Description: "",
				Optional:    true,
				ElementType: types.StringType,
			},
			"subcategory": rsschema.ListAttribute{
				Description: "",
				Optional:    true,
				ElementType: types.StringType,
			},
			"tagging": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"no_tag": rsschema.BoolAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"tag": rsschema.ListAttribute{
						Description: "",
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
			"technology": rsschema.ListAttribute{
				Description: "",
				Optional:    true,
				ElementType: types.StringType,
			},
			"transfers_files": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"tunnels_other_apps": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"used_by_malware": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *objectsApplicationFiltersResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *objectsApplicationFiltersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state objectsApplicationFiltersRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "sase_objects_application_filters",
		"folder":        state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := jHKNPjP.NewClient(r.client)
	input := jHKNPjP.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 lhPcfTR.Config
	var0.Category = DecodeStringSlice(state.Category)
	var0.Evasive = state.Evasive.ValueBool()
	var0.ExcessiveBandwidthUse = state.ExcessiveBandwidthUse.ValueBool()
	var0.Exclude = DecodeStringSlice(state.Exclude)
	var0.HasKnownVulnerabilities = state.HasKnownVulnerabilities.ValueBool()
	var0.IsSaas = state.IsSaas.ValueBool()
	var0.Name = state.Name.ValueString()
	var0.NewAppid = state.NewAppid.ValueBool()
	var0.Pervasive = state.Pervasive.ValueBool()
	var0.ProneToMisuse = state.ProneToMisuse.ValueBool()
	var0.Risk = DecodeInt64Slice(state.Risk)
	var0.SaasCertifications = DecodeStringSlice(state.SaasCertifications)
	var0.SaasRisk = DecodeStringSlice(state.SaasRisk)
	var0.Subcategory = DecodeStringSlice(state.Subcategory)
	var var1 *lhPcfTR.TaggingObject
	if state.Tagging != nil {
		var1 = &lhPcfTR.TaggingObject{}
		var1.NoTag = state.Tagging.NoTag.ValueBool()
		var1.Tag = DecodeStringSlice(state.Tagging.Tag)
	}
	var0.Tagging = var1
	var0.Technology = DecodeStringSlice(state.Technology)
	var0.TransfersFiles = state.TransfersFiles.ValueBool()
	var0.TunnelsOtherApps = state.TunnelsOtherApps.ValueBool()
	var0.UsedByMalware = state.UsedByMalware.ValueBool()
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Create(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in create", err.Error())
		return
	}

	// Store the answer to state.
	state.Id = types.StringValue(strings.Join([]string{input.Folder, ans.ObjectId}, IdSeparator))
	var var2 *objectsApplicationFiltersRsModelTaggingObject
	if ans.Tagging != nil {
		var2 = &objectsApplicationFiltersRsModelTaggingObject{}
		var2.NoTag = types.BoolValue(ans.Tagging.NoTag)
		var2.Tag = EncodeStringSlice(ans.Tagging.Tag)
	}
	state.Category = EncodeStringSlice(ans.Category)
	state.Evasive = types.BoolValue(ans.Evasive)
	state.ExcessiveBandwidthUse = types.BoolValue(ans.ExcessiveBandwidthUse)
	state.Exclude = EncodeStringSlice(ans.Exclude)
	state.HasKnownVulnerabilities = types.BoolValue(ans.HasKnownVulnerabilities)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.IsSaas = types.BoolValue(ans.IsSaas)
	state.Name = types.StringValue(ans.Name)
	state.NewAppid = types.BoolValue(ans.NewAppid)
	state.Pervasive = types.BoolValue(ans.Pervasive)
	state.ProneToMisuse = types.BoolValue(ans.ProneToMisuse)
	state.Risk = EncodeInt64Slice(ans.Risk)
	state.SaasCertifications = EncodeStringSlice(ans.SaasCertifications)
	state.SaasRisk = EncodeStringSlice(ans.SaasRisk)
	state.Subcategory = EncodeStringSlice(ans.Subcategory)
	state.Tagging = var2
	state.Technology = EncodeStringSlice(ans.Technology)
	state.TransfersFiles = types.BoolValue(ans.TransfersFiles)
	state.TunnelsOtherApps = types.BoolValue(ans.TunnelsOtherApps)
	state.UsedByMalware = types.BoolValue(ans.UsedByMalware)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *objectsApplicationFiltersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state objectsApplicationFiltersRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "sase_objects_application_filters",
		"locMap":        map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":        tokens,
	})

	// Prepare to read the config.
	svc := jHKNPjP.NewClient(r.client)
	input := jHKNPjP.ReadInput{
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
	var var0 *objectsApplicationFiltersRsModelTaggingObject
	if ans.Tagging != nil {
		var0 = &objectsApplicationFiltersRsModelTaggingObject{}
		var0.NoTag = types.BoolValue(ans.Tagging.NoTag)
		var0.Tag = EncodeStringSlice(ans.Tagging.Tag)
	}
	state.Category = EncodeStringSlice(ans.Category)
	state.Evasive = types.BoolValue(ans.Evasive)
	state.ExcessiveBandwidthUse = types.BoolValue(ans.ExcessiveBandwidthUse)
	state.Exclude = EncodeStringSlice(ans.Exclude)
	state.HasKnownVulnerabilities = types.BoolValue(ans.HasKnownVulnerabilities)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.IsSaas = types.BoolValue(ans.IsSaas)
	state.Name = types.StringValue(ans.Name)
	state.NewAppid = types.BoolValue(ans.NewAppid)
	state.Pervasive = types.BoolValue(ans.Pervasive)
	state.ProneToMisuse = types.BoolValue(ans.ProneToMisuse)
	state.Risk = EncodeInt64Slice(ans.Risk)
	state.SaasCertifications = EncodeStringSlice(ans.SaasCertifications)
	state.SaasRisk = EncodeStringSlice(ans.SaasRisk)
	state.Subcategory = EncodeStringSlice(ans.Subcategory)
	state.Tagging = var0
	state.Technology = EncodeStringSlice(ans.Technology)
	state.TransfersFiles = types.BoolValue(ans.TransfersFiles)
	state.TunnelsOtherApps = types.BoolValue(ans.TunnelsOtherApps)
	state.UsedByMalware = types.BoolValue(ans.UsedByMalware)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *objectsApplicationFiltersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state objectsApplicationFiltersRsModel
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
		"resource_name": "sase_objects_application_filters",
		"object_id":     state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := jHKNPjP.NewClient(r.client)
	input := jHKNPjP.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 lhPcfTR.Config
	var0.Category = DecodeStringSlice(plan.Category)
	var0.Evasive = plan.Evasive.ValueBool()
	var0.ExcessiveBandwidthUse = plan.ExcessiveBandwidthUse.ValueBool()
	var0.Exclude = DecodeStringSlice(plan.Exclude)
	var0.HasKnownVulnerabilities = plan.HasKnownVulnerabilities.ValueBool()
	var0.IsSaas = plan.IsSaas.ValueBool()
	var0.Name = plan.Name.ValueString()
	var0.NewAppid = plan.NewAppid.ValueBool()
	var0.Pervasive = plan.Pervasive.ValueBool()
	var0.ProneToMisuse = plan.ProneToMisuse.ValueBool()
	var0.Risk = DecodeInt64Slice(plan.Risk)
	var0.SaasCertifications = DecodeStringSlice(plan.SaasCertifications)
	var0.SaasRisk = DecodeStringSlice(plan.SaasRisk)
	var0.Subcategory = DecodeStringSlice(plan.Subcategory)
	var var1 *lhPcfTR.TaggingObject
	if plan.Tagging != nil {
		var1 = &lhPcfTR.TaggingObject{}
		var1.NoTag = plan.Tagging.NoTag.ValueBool()
		var1.Tag = DecodeStringSlice(plan.Tagging.Tag)
	}
	var0.Tagging = var1
	var0.Technology = DecodeStringSlice(plan.Technology)
	var0.TransfersFiles = plan.TransfersFiles.ValueBool()
	var0.TunnelsOtherApps = plan.TunnelsOtherApps.ValueBool()
	var0.UsedByMalware = plan.UsedByMalware.ValueBool()
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var2 *objectsApplicationFiltersRsModelTaggingObject
	if ans.Tagging != nil {
		var2 = &objectsApplicationFiltersRsModelTaggingObject{}
		var2.NoTag = types.BoolValue(ans.Tagging.NoTag)
		var2.Tag = EncodeStringSlice(ans.Tagging.Tag)
	}
	state.Category = EncodeStringSlice(ans.Category)
	state.Evasive = types.BoolValue(ans.Evasive)
	state.ExcessiveBandwidthUse = types.BoolValue(ans.ExcessiveBandwidthUse)
	state.Exclude = EncodeStringSlice(ans.Exclude)
	state.HasKnownVulnerabilities = types.BoolValue(ans.HasKnownVulnerabilities)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.IsSaas = types.BoolValue(ans.IsSaas)
	state.Name = types.StringValue(ans.Name)
	state.NewAppid = types.BoolValue(ans.NewAppid)
	state.Pervasive = types.BoolValue(ans.Pervasive)
	state.ProneToMisuse = types.BoolValue(ans.ProneToMisuse)
	state.Risk = EncodeInt64Slice(ans.Risk)
	state.SaasCertifications = EncodeStringSlice(ans.SaasCertifications)
	state.SaasRisk = EncodeStringSlice(ans.SaasRisk)
	state.Subcategory = EncodeStringSlice(ans.Subcategory)
	state.Tagging = var2
	state.Technology = EncodeStringSlice(ans.Technology)
	state.TransfersFiles = types.BoolValue(ans.TransfersFiles)
	state.TunnelsOtherApps = types.BoolValue(ans.TunnelsOtherApps)
	state.UsedByMalware = types.BoolValue(ans.UsedByMalware)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *objectsApplicationFiltersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name": "sase_objects_application_filters",
		"locMap":        map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":        tokens,
	})

	svc := jHKNPjP.NewClient(r.client)
	input := jHKNPjP.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *objectsApplicationFiltersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}