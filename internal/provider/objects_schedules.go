package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	qFVQpmA "github.com/paloaltonetworks/sase-go/netsec/schema/objects/schedules"
	lNTtdgX "github.com/paloaltonetworks/sase-go/netsec/service/v1/schedules"

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
	_ datasource.DataSource              = &objectsSchedulesListDataSource{}
	_ datasource.DataSourceWithConfigure = &objectsSchedulesListDataSource{}
)

func NewObjectsSchedulesListDataSource() datasource.DataSource {
	return &objectsSchedulesListDataSource{}
}

type objectsSchedulesListDataSource struct {
	client *sase.Client
}

type objectsSchedulesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []objectsSchedulesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type objectsSchedulesListDsModelConfig struct {
	ObjectId     types.String                                  `tfsdk:"object_id"`
	Name         types.String                                  `tfsdk:"name"`
	ScheduleType objectsSchedulesListDsModelScheduleTypeObject `tfsdk:"schedule_type"`
}

type objectsSchedulesListDsModelScheduleTypeObject struct {
	NonRecurring []types.String                              `tfsdk:"non_recurring"`
	Recurring    *objectsSchedulesListDsModelRecurringObject `tfsdk:"recurring"`
}

type objectsSchedulesListDsModelRecurringObject struct {
	Daily  []types.String                           `tfsdk:"daily"`
	Weekly *objectsSchedulesListDsModelWeeklyObject `tfsdk:"weekly"`
}

type objectsSchedulesListDsModelWeeklyObject struct {
	Friday    []types.String `tfsdk:"friday"`
	Monday    []types.String `tfsdk:"monday"`
	Saturday  []types.String `tfsdk:"saturday"`
	Sunday    []types.String `tfsdk:"sunday"`
	Thursday  []types.String `tfsdk:"thursday"`
	Tuesday   []types.String `tfsdk:"tuesday"`
	Wednesday []types.String `tfsdk:"wednesday"`
}

// Metadata returns the data source type name.
func (d *objectsSchedulesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_schedules_list"
}

// Schema defines the schema for this listing data source.
func (d *objectsSchedulesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"object_id": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"schedule_type": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"non_recurring": dsschema.ListAttribute{
									Description: "",
									Computed:    true,
									ElementType: types.StringType,
								},
								"recurring": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"daily": dsschema.ListAttribute{
											Description: "",
											Computed:    true,
											ElementType: types.StringType,
										},
										"weekly": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"friday": dsschema.ListAttribute{
													Description: "",
													Computed:    true,
													ElementType: types.StringType,
												},
												"monday": dsschema.ListAttribute{
													Description: "",
													Computed:    true,
													ElementType: types.StringType,
												},
												"saturday": dsschema.ListAttribute{
													Description: "",
													Computed:    true,
													ElementType: types.StringType,
												},
												"sunday": dsschema.ListAttribute{
													Description: "",
													Computed:    true,
													ElementType: types.StringType,
												},
												"thursday": dsschema.ListAttribute{
													Description: "",
													Computed:    true,
													ElementType: types.StringType,
												},
												"tuesday": dsschema.ListAttribute{
													Description: "",
													Computed:    true,
													ElementType: types.StringType,
												},
												"wednesday": dsschema.ListAttribute{
													Description: "",
													Computed:    true,
													ElementType: types.StringType,
												},
											},
										},
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
func (d *objectsSchedulesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *objectsSchedulesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state objectsSchedulesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_objects_schedules_list",
		"limit":            state.Limit.ValueInt64(),
		"has_limit":        !state.Limit.IsNull(),
		"offset":           state.Offset.ValueInt64(),
		"has_offset":       !state.Offset.IsNull(),
		"name":             state.Name.ValueString(),
		"has_name":         !state.Name.IsNull(),
		"folder":           state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := lNTtdgX.NewClient(d.client)
	input := lNTtdgX.ListInput{
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
	var var0 []objectsSchedulesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]objectsSchedulesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 objectsSchedulesListDsModelConfig
			var var3 objectsSchedulesListDsModelScheduleTypeObject
			var var4 *objectsSchedulesListDsModelRecurringObject
			if var1.ScheduleType.Recurring != nil {
				var4 = &objectsSchedulesListDsModelRecurringObject{}
				var var5 *objectsSchedulesListDsModelWeeklyObject
				if var1.ScheduleType.Recurring.Weekly != nil {
					var5 = &objectsSchedulesListDsModelWeeklyObject{}
					var5.Friday = EncodeStringSlice(var1.ScheduleType.Recurring.Weekly.Friday)
					var5.Monday = EncodeStringSlice(var1.ScheduleType.Recurring.Weekly.Monday)
					var5.Saturday = EncodeStringSlice(var1.ScheduleType.Recurring.Weekly.Saturday)
					var5.Sunday = EncodeStringSlice(var1.ScheduleType.Recurring.Weekly.Sunday)
					var5.Thursday = EncodeStringSlice(var1.ScheduleType.Recurring.Weekly.Thursday)
					var5.Tuesday = EncodeStringSlice(var1.ScheduleType.Recurring.Weekly.Tuesday)
					var5.Wednesday = EncodeStringSlice(var1.ScheduleType.Recurring.Weekly.Wednesday)
				}
				var4.Daily = EncodeStringSlice(var1.ScheduleType.Recurring.Daily)
				var4.Weekly = var5
			}
			var3.NonRecurring = EncodeStringSlice(var1.ScheduleType.NonRecurring)
			var3.Recurring = var4
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Name = types.StringValue(var1.Name)
			var2.ScheduleType = var3
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
	_ datasource.DataSource              = &objectsSchedulesDataSource{}
	_ datasource.DataSourceWithConfigure = &objectsSchedulesDataSource{}
)

func NewObjectsSchedulesDataSource() datasource.DataSource {
	return &objectsSchedulesDataSource{}
}

type objectsSchedulesDataSource struct {
	client *sase.Client
}

type objectsSchedulesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`

	// Output.
	// Ref: #/components/schemas/objects-schedules
	// input omit: ObjectId
	Name         types.String                              `tfsdk:"name"`
	ScheduleType objectsSchedulesDsModelScheduleTypeObject `tfsdk:"schedule_type"`
}

type objectsSchedulesDsModelScheduleTypeObject struct {
	NonRecurring []types.String                          `tfsdk:"non_recurring"`
	Recurring    *objectsSchedulesDsModelRecurringObject `tfsdk:"recurring"`
}

type objectsSchedulesDsModelRecurringObject struct {
	Daily  []types.String                       `tfsdk:"daily"`
	Weekly *objectsSchedulesDsModelWeeklyObject `tfsdk:"weekly"`
}

type objectsSchedulesDsModelWeeklyObject struct {
	Friday    []types.String `tfsdk:"friday"`
	Monday    []types.String `tfsdk:"monday"`
	Saturday  []types.String `tfsdk:"saturday"`
	Sunday    []types.String `tfsdk:"sunday"`
	Thursday  []types.String `tfsdk:"thursday"`
	Tuesday   []types.String `tfsdk:"tuesday"`
	Wednesday []types.String `tfsdk:"wednesday"`
}

// Metadata returns the data source type name.
func (d *objectsSchedulesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_schedules"
}

// Schema defines the schema for this listing data source.
func (d *objectsSchedulesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"schedule_type": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"non_recurring": dsschema.ListAttribute{
						Description: "",
						Computed:    true,
						ElementType: types.StringType,
					},
					"recurring": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"daily": dsschema.ListAttribute{
								Description: "",
								Computed:    true,
								ElementType: types.StringType,
							},
							"weekly": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"friday": dsschema.ListAttribute{
										Description: "",
										Computed:    true,
										ElementType: types.StringType,
									},
									"monday": dsschema.ListAttribute{
										Description: "",
										Computed:    true,
										ElementType: types.StringType,
									},
									"saturday": dsschema.ListAttribute{
										Description: "",
										Computed:    true,
										ElementType: types.StringType,
									},
									"sunday": dsschema.ListAttribute{
										Description: "",
										Computed:    true,
										ElementType: types.StringType,
									},
									"thursday": dsschema.ListAttribute{
										Description: "",
										Computed:    true,
										ElementType: types.StringType,
									},
									"tuesday": dsschema.ListAttribute{
										Description: "",
										Computed:    true,
										ElementType: types.StringType,
									},
									"wednesday": dsschema.ListAttribute{
										Description: "",
										Computed:    true,
										ElementType: types.StringType,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure prepares the struct.
func (d *objectsSchedulesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *objectsSchedulesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state objectsSchedulesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_objects_schedules",
		"object_id":        state.ObjectId.ValueString(),
	})

	// Prepare to run the command.
	svc := lNTtdgX.NewClient(d.client)
	input := lNTtdgX.ReadInput{
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
	var var0 objectsSchedulesDsModelScheduleTypeObject
	var var1 *objectsSchedulesDsModelRecurringObject
	if ans.ScheduleType.Recurring != nil {
		var1 = &objectsSchedulesDsModelRecurringObject{}
		var var2 *objectsSchedulesDsModelWeeklyObject
		if ans.ScheduleType.Recurring.Weekly != nil {
			var2 = &objectsSchedulesDsModelWeeklyObject{}
			var2.Friday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Friday)
			var2.Monday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Monday)
			var2.Saturday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Saturday)
			var2.Sunday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Sunday)
			var2.Thursday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Thursday)
			var2.Tuesday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Tuesday)
			var2.Wednesday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Wednesday)
		}
		var1.Daily = EncodeStringSlice(ans.ScheduleType.Recurring.Daily)
		var1.Weekly = var2
	}
	var0.NonRecurring = EncodeStringSlice(ans.ScheduleType.NonRecurring)
	var0.Recurring = var1
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.ScheduleType = var0

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &objectsSchedulesResource{}
	_ resource.ResourceWithConfigure   = &objectsSchedulesResource{}
	_ resource.ResourceWithImportState = &objectsSchedulesResource{}
)

func NewObjectsSchedulesResource() resource.Resource {
	return &objectsSchedulesResource{}
}

type objectsSchedulesResource struct {
	client *sase.Client
}

type objectsSchedulesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/objects-schedules
	ObjectId     types.String                              `tfsdk:"object_id"`
	Name         types.String                              `tfsdk:"name"`
	ScheduleType objectsSchedulesRsModelScheduleTypeObject `tfsdk:"schedule_type"`
}

type objectsSchedulesRsModelScheduleTypeObject struct {
	NonRecurring []types.String                          `tfsdk:"non_recurring"`
	Recurring    *objectsSchedulesRsModelRecurringObject `tfsdk:"recurring"`
}

type objectsSchedulesRsModelRecurringObject struct {
	Daily  []types.String                       `tfsdk:"daily"`
	Weekly *objectsSchedulesRsModelWeeklyObject `tfsdk:"weekly"`
}

type objectsSchedulesRsModelWeeklyObject struct {
	Friday    []types.String `tfsdk:"friday"`
	Monday    []types.String `tfsdk:"monday"`
	Saturday  []types.String `tfsdk:"saturday"`
	Sunday    []types.String `tfsdk:"sunday"`
	Thursday  []types.String `tfsdk:"thursday"`
	Tuesday   []types.String `tfsdk:"tuesday"`
	Wednesday []types.String `tfsdk:"wednesday"`
}

// Metadata returns the data source type name.
func (r *objectsSchedulesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_schedules"
}

// Schema defines the schema for this listing data source.
func (r *objectsSchedulesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"schedule_type": rsschema.SingleNestedAttribute{
				Description: "",
				Required:    true,
				Attributes: map[string]rsschema.Attribute{
					"non_recurring": rsschema.ListAttribute{
						Description: "",
						Optional:    true,
						ElementType: types.StringType,
					},
					"recurring": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"daily": rsschema.ListAttribute{
								Description: "",
								Optional:    true,
								ElementType: types.StringType,
							},
							"weekly": rsschema.SingleNestedAttribute{
								Description: "",
								Optional:    true,
								Attributes: map[string]rsschema.Attribute{
									"friday": rsschema.ListAttribute{
										Description: "",
										Optional:    true,
										ElementType: types.StringType,
									},
									"monday": rsschema.ListAttribute{
										Description: "",
										Optional:    true,
										ElementType: types.StringType,
									},
									"saturday": rsschema.ListAttribute{
										Description: "",
										Optional:    true,
										ElementType: types.StringType,
									},
									"sunday": rsschema.ListAttribute{
										Description: "",
										Optional:    true,
										ElementType: types.StringType,
									},
									"thursday": rsschema.ListAttribute{
										Description: "",
										Optional:    true,
										ElementType: types.StringType,
									},
									"tuesday": rsschema.ListAttribute{
										Description: "",
										Optional:    true,
										ElementType: types.StringType,
									},
									"wednesday": rsschema.ListAttribute{
										Description: "",
										Optional:    true,
										ElementType: types.StringType,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *objectsSchedulesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *objectsSchedulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state objectsSchedulesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "sase_objects_schedules",
		"folder":        state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := lNTtdgX.NewClient(r.client)
	input := lNTtdgX.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 qFVQpmA.Config
	var0.Name = state.Name.ValueString()
	var var1 qFVQpmA.ScheduleTypeObject
	var1.NonRecurring = DecodeStringSlice(state.ScheduleType.NonRecurring)
	var var2 *qFVQpmA.RecurringObject
	if state.ScheduleType.Recurring != nil {
		var2 = &qFVQpmA.RecurringObject{}
		var2.Daily = DecodeStringSlice(state.ScheduleType.Recurring.Daily)
		var var3 *qFVQpmA.WeeklyObject
		if state.ScheduleType.Recurring.Weekly != nil {
			var3 = &qFVQpmA.WeeklyObject{}
			var3.Friday = DecodeStringSlice(state.ScheduleType.Recurring.Weekly.Friday)
			var3.Monday = DecodeStringSlice(state.ScheduleType.Recurring.Weekly.Monday)
			var3.Saturday = DecodeStringSlice(state.ScheduleType.Recurring.Weekly.Saturday)
			var3.Sunday = DecodeStringSlice(state.ScheduleType.Recurring.Weekly.Sunday)
			var3.Thursday = DecodeStringSlice(state.ScheduleType.Recurring.Weekly.Thursday)
			var3.Tuesday = DecodeStringSlice(state.ScheduleType.Recurring.Weekly.Tuesday)
			var3.Wednesday = DecodeStringSlice(state.ScheduleType.Recurring.Weekly.Wednesday)
		}
		var2.Weekly = var3
	}
	var1.Recurring = var2
	var0.ScheduleType = var1
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Create(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in create", err.Error())
		return
	}

	// Store the answer to state.
	state.Id = types.StringValue(strings.Join([]string{input.Folder, ans.ObjectId}, IdSeparator))
	var var4 objectsSchedulesRsModelScheduleTypeObject
	var var5 *objectsSchedulesRsModelRecurringObject
	if ans.ScheduleType.Recurring != nil {
		var5 = &objectsSchedulesRsModelRecurringObject{}
		var var6 *objectsSchedulesRsModelWeeklyObject
		if ans.ScheduleType.Recurring.Weekly != nil {
			var6 = &objectsSchedulesRsModelWeeklyObject{}
			var6.Friday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Friday)
			var6.Monday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Monday)
			var6.Saturday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Saturday)
			var6.Sunday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Sunday)
			var6.Thursday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Thursday)
			var6.Tuesday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Tuesday)
			var6.Wednesday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Wednesday)
		}
		var5.Daily = EncodeStringSlice(ans.ScheduleType.Recurring.Daily)
		var5.Weekly = var6
	}
	var4.NonRecurring = EncodeStringSlice(ans.ScheduleType.NonRecurring)
	var4.Recurring = var5
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.ScheduleType = var4

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *objectsSchedulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state objectsSchedulesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "sase_objects_schedules",
		"locMap":        map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":        tokens,
	})

	// Prepare to read the config.
	svc := lNTtdgX.NewClient(r.client)
	input := lNTtdgX.ReadInput{
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
	var var0 objectsSchedulesRsModelScheduleTypeObject
	var var1 *objectsSchedulesRsModelRecurringObject
	if ans.ScheduleType.Recurring != nil {
		var1 = &objectsSchedulesRsModelRecurringObject{}
		var var2 *objectsSchedulesRsModelWeeklyObject
		if ans.ScheduleType.Recurring.Weekly != nil {
			var2 = &objectsSchedulesRsModelWeeklyObject{}
			var2.Friday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Friday)
			var2.Monday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Monday)
			var2.Saturday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Saturday)
			var2.Sunday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Sunday)
			var2.Thursday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Thursday)
			var2.Tuesday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Tuesday)
			var2.Wednesday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Wednesday)
		}
		var1.Daily = EncodeStringSlice(ans.ScheduleType.Recurring.Daily)
		var1.Weekly = var2
	}
	var0.NonRecurring = EncodeStringSlice(ans.ScheduleType.NonRecurring)
	var0.Recurring = var1
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.ScheduleType = var0

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *objectsSchedulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state objectsSchedulesRsModel
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
		"resource_name": "sase_objects_schedules",
		"object_id":     state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := lNTtdgX.NewClient(r.client)
	input := lNTtdgX.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 qFVQpmA.Config
	var0.Name = plan.Name.ValueString()
	var var1 qFVQpmA.ScheduleTypeObject
	var1.NonRecurring = DecodeStringSlice(plan.ScheduleType.NonRecurring)
	var var2 *qFVQpmA.RecurringObject
	if plan.ScheduleType.Recurring != nil {
		var2 = &qFVQpmA.RecurringObject{}
		var2.Daily = DecodeStringSlice(plan.ScheduleType.Recurring.Daily)
		var var3 *qFVQpmA.WeeklyObject
		if plan.ScheduleType.Recurring.Weekly != nil {
			var3 = &qFVQpmA.WeeklyObject{}
			var3.Friday = DecodeStringSlice(plan.ScheduleType.Recurring.Weekly.Friday)
			var3.Monday = DecodeStringSlice(plan.ScheduleType.Recurring.Weekly.Monday)
			var3.Saturday = DecodeStringSlice(plan.ScheduleType.Recurring.Weekly.Saturday)
			var3.Sunday = DecodeStringSlice(plan.ScheduleType.Recurring.Weekly.Sunday)
			var3.Thursday = DecodeStringSlice(plan.ScheduleType.Recurring.Weekly.Thursday)
			var3.Tuesday = DecodeStringSlice(plan.ScheduleType.Recurring.Weekly.Tuesday)
			var3.Wednesday = DecodeStringSlice(plan.ScheduleType.Recurring.Weekly.Wednesday)
		}
		var2.Weekly = var3
	}
	var1.Recurring = var2
	var0.ScheduleType = var1
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var4 objectsSchedulesRsModelScheduleTypeObject
	var var5 *objectsSchedulesRsModelRecurringObject
	if ans.ScheduleType.Recurring != nil {
		var5 = &objectsSchedulesRsModelRecurringObject{}
		var var6 *objectsSchedulesRsModelWeeklyObject
		if ans.ScheduleType.Recurring.Weekly != nil {
			var6 = &objectsSchedulesRsModelWeeklyObject{}
			var6.Friday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Friday)
			var6.Monday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Monday)
			var6.Saturday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Saturday)
			var6.Sunday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Sunday)
			var6.Thursday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Thursday)
			var6.Tuesday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Tuesday)
			var6.Wednesday = EncodeStringSlice(ans.ScheduleType.Recurring.Weekly.Wednesday)
		}
		var5.Daily = EncodeStringSlice(ans.ScheduleType.Recurring.Daily)
		var5.Weekly = var6
	}
	var4.NonRecurring = EncodeStringSlice(ans.ScheduleType.NonRecurring)
	var4.Recurring = var5
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.ScheduleType = var4

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *objectsSchedulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name": "sase_objects_schedules",
		"locMap":        map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":        tokens,
	})

	svc := lNTtdgX.NewClient(r.client)
	input := lNTtdgX.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *objectsSchedulesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
