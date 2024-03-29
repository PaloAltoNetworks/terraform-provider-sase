package provider

import (
	"context"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	wugpput "github.com/paloaltonetworks/sase-go/netsec/service/v1/jobs"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Data source (listing).
var (
	_ datasource.DataSource              = &jobsListDataSource{}
	_ datasource.DataSourceWithConfigure = &jobsListDataSource{}
)

func NewJobsListDataSource() datasource.DataSource {
	return &jobsListDataSource{}
}

type jobsListDataSource struct {
	client *sase.Client
}

type jobsListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.

	// Output.
	Data   []jobsListDsModelConfig `tfsdk:"data"`
	Limit  types.Int64             `tfsdk:"limit"`
	Offset types.Int64             `tfsdk:"offset"`
	Total  types.Int64             `tfsdk:"total"`
}

type jobsListDsModelConfig struct {
	Details    types.String `tfsdk:"details"`
	EndTs      types.String `tfsdk:"end_ts"`
	ObjectId   types.String `tfsdk:"object_id"`
	InsertTs   types.String `tfsdk:"insert_ts"`
	JobResult  types.String `tfsdk:"job_result"`
	JobStatus  types.String `tfsdk:"job_status"`
	JobType    types.String `tfsdk:"job_type"`
	LastUpdate types.String `tfsdk:"last_update"`
	OpaqueInt  types.String `tfsdk:"opaque_int"`
	OpaqueStr  types.String `tfsdk:"opaque_str"`
	Owner      types.String `tfsdk:"owner"`
	ParentId   types.String `tfsdk:"parent_id"`
	Percent    types.String `tfsdk:"percent"`
	ResultI    types.String `tfsdk:"result_i"`
	ResultStr  types.String `tfsdk:"result_str"`
	SessionId  types.String `tfsdk:"session_id"`
	StartTs    types.String `tfsdk:"start_ts"`
	StatusI    types.String `tfsdk:"status_i"`
	StatusStr  types.String `tfsdk:"status_str"`
	Summary    types.String `tfsdk:"summary"`
	TypeI      types.String `tfsdk:"type_i"`
	TypeStr    types.String `tfsdk:"type_str"`
	Uname      types.String `tfsdk:"uname"`
}

// Metadata returns the data source type name.
func (d *jobsListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jobs_list"
}

// Schema defines the schema for this listing data source.
func (d *jobsListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Retrieves a listing of config items.",

		Attributes: map[string]dsschema.Attribute{
			"id": dsschema.StringAttribute{
				Description:         "The object ID.",
				MarkdownDescription: "The object ID.",
				Computed:            true,
			},

			// Output.
			"data": dsschema.ListNestedAttribute{
				Description:         "The `data` parameter.",
				MarkdownDescription: "The `data` parameter.",
				Computed:            true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"details": dsschema.StringAttribute{
							Description:         "The `details` parameter.",
							MarkdownDescription: "The `details` parameter.",
							Computed:            true,
						},
						"end_ts": dsschema.StringAttribute{
							Description:         "The `end_ts` parameter.",
							MarkdownDescription: "The `end_ts` parameter.",
							Computed:            true,
						},
						"object_id": dsschema.StringAttribute{
							Description:         "The `object_id` parameter.",
							MarkdownDescription: "The `object_id` parameter.",
							Computed:            true,
						},
						"insert_ts": dsschema.StringAttribute{
							Description:         "The `insert_ts` parameter.",
							MarkdownDescription: "The `insert_ts` parameter.",
							Computed:            true,
						},
						"job_result": dsschema.StringAttribute{
							Description:         "The `job_result` parameter.",
							MarkdownDescription: "The `job_result` parameter.",
							Computed:            true,
						},
						"job_status": dsschema.StringAttribute{
							Description:         "The `job_status` parameter.",
							MarkdownDescription: "The `job_status` parameter.",
							Computed:            true,
						},
						"job_type": dsschema.StringAttribute{
							Description:         "The `job_type` parameter.",
							MarkdownDescription: "The `job_type` parameter.",
							Computed:            true,
						},
						"last_update": dsschema.StringAttribute{
							Description:         "The `last_update` parameter.",
							MarkdownDescription: "The `last_update` parameter.",
							Computed:            true,
						},
						"opaque_int": dsschema.StringAttribute{
							Description:         "The `opaque_int` parameter.",
							MarkdownDescription: "The `opaque_int` parameter.",
							Computed:            true,
						},
						"opaque_str": dsschema.StringAttribute{
							Description:         "The `opaque_str` parameter.",
							MarkdownDescription: "The `opaque_str` parameter.",
							Computed:            true,
						},
						"owner": dsschema.StringAttribute{
							Description:         "The `owner` parameter.",
							MarkdownDescription: "The `owner` parameter.",
							Computed:            true,
						},
						"parent_id": dsschema.StringAttribute{
							Description:         "The `parent_id` parameter.",
							MarkdownDescription: "The `parent_id` parameter.",
							Computed:            true,
						},
						"percent": dsschema.StringAttribute{
							Description:         "The `percent` parameter.",
							MarkdownDescription: "The `percent` parameter.",
							Computed:            true,
						},
						"result_i": dsschema.StringAttribute{
							Description:         "The `result_i` parameter.",
							MarkdownDescription: "The `result_i` parameter.",
							Computed:            true,
						},
						"result_str": dsschema.StringAttribute{
							Description:         "The `result_str` parameter.",
							MarkdownDescription: "The `result_str` parameter.",
							Computed:            true,
						},
						"session_id": dsschema.StringAttribute{
							Description:         "The `session_id` parameter.",
							MarkdownDescription: "The `session_id` parameter.",
							Computed:            true,
						},
						"start_ts": dsschema.StringAttribute{
							Description:         "The `start_ts` parameter.",
							MarkdownDescription: "The `start_ts` parameter.",
							Computed:            true,
						},
						"status_i": dsschema.StringAttribute{
							Description:         "The `status_i` parameter.",
							MarkdownDescription: "The `status_i` parameter.",
							Computed:            true,
						},
						"status_str": dsschema.StringAttribute{
							Description:         "The `status_str` parameter.",
							MarkdownDescription: "The `status_str` parameter.",
							Computed:            true,
						},
						"summary": dsschema.StringAttribute{
							Description:         "The `summary` parameter.",
							MarkdownDescription: "The `summary` parameter.",
							Computed:            true,
						},
						"type_i": dsschema.StringAttribute{
							Description:         "The `type_i` parameter.",
							MarkdownDescription: "The `type_i` parameter.",
							Computed:            true,
						},
						"type_str": dsschema.StringAttribute{
							Description:         "The `type_str` parameter.",
							MarkdownDescription: "The `type_str` parameter.",
							Computed:            true,
						},
						"uname": dsschema.StringAttribute{
							Description:         "The `uname` parameter.",
							MarkdownDescription: "The `uname` parameter.",
							Computed:            true,
						},
					},
				},
			},
			"limit": dsschema.Int64Attribute{
				Description:         "The `limit` parameter.",
				MarkdownDescription: "The `limit` parameter.",
				Computed:            true,
			},
			"offset": dsschema.Int64Attribute{
				Description:         "The `offset` parameter.",
				MarkdownDescription: "The `offset` parameter.",
				Computed:            true,
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
func (d *jobsListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *jobsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state jobsListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_jobs_list",
		"terraform_provider_function": "Read",
	})

	// Prepare to run the command.
	svc := wugpput.NewClient(d.client)
	// Perform the operation.
	ans, err := svc.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error getting listing", err.Error())
		return
	}

	// Store the answer to state.
	var idBuilder strings.Builder
	idBuilder.WriteString("sase")
	state.Id = types.StringValue(idBuilder.String())
	var var0 []jobsListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]jobsListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 jobsListDsModelConfig
			var2.Details = types.StringValue(var1.Details)
			var2.EndTs = types.StringValue(var1.EndTs)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.InsertTs = types.StringValue(var1.InsertTs)
			var2.JobResult = types.StringValue(var1.JobResult)
			var2.JobStatus = types.StringValue(var1.JobStatus)
			var2.JobType = types.StringValue(var1.JobType)
			var2.LastUpdate = types.StringValue(var1.LastUpdate)
			var2.OpaqueInt = types.StringValue(var1.OpaqueInt)
			var2.OpaqueStr = types.StringValue(var1.OpaqueStr)
			var2.Owner = types.StringValue(var1.Owner)
			var2.ParentId = types.StringValue(var1.ParentId)
			var2.Percent = types.StringValue(var1.Percent)
			var2.ResultI = types.StringValue(var1.ResultI)
			var2.ResultStr = types.StringValue(var1.ResultStr)
			var2.SessionId = types.StringValue(var1.SessionId)
			var2.StartTs = types.StringValue(var1.StartTs)
			var2.StatusI = types.StringValue(var1.StatusI)
			var2.StatusStr = types.StringValue(var1.StatusStr)
			var2.Summary = types.StringValue(var1.Summary)
			var2.TypeI = types.StringValue(var1.TypeI)
			var2.TypeStr = types.StringValue(var1.TypeStr)
			var2.Uname = types.StringValue(var1.Uname)
			var0 = append(var0, var2)
		}
	}
	state.Data = var0
	state.Limit = types.Int64Value(ans.Limit)
	state.Offset = types.Int64Value(ans.Offset)
	state.Total = types.Int64Value(ans.Total)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Data source.
var (
	_ datasource.DataSource              = &jobsDataSource{}
	_ datasource.DataSourceWithConfigure = &jobsDataSource{}
)

func NewJobsDataSource() datasource.DataSource {
	return &jobsDataSource{}
}

type jobsDataSource struct {
	client *sase.Client
}

type jobsDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	JobId types.String `tfsdk:"job_id"`

	// Output.
	// Ref: #/components/schemas/jobs
	Details    types.String `tfsdk:"details"`
	EndTs      types.String `tfsdk:"end_ts"`
	ObjectId   types.String `tfsdk:"object_id"`
	InsertTs   types.String `tfsdk:"insert_ts"`
	JobResult  types.String `tfsdk:"job_result"`
	JobStatus  types.String `tfsdk:"job_status"`
	JobType    types.String `tfsdk:"job_type"`
	LastUpdate types.String `tfsdk:"last_update"`
	OpaqueInt  types.String `tfsdk:"opaque_int"`
	OpaqueStr  types.String `tfsdk:"opaque_str"`
	Owner      types.String `tfsdk:"owner"`
	ParentId   types.String `tfsdk:"parent_id"`
	Percent    types.String `tfsdk:"percent"`
	ResultI    types.String `tfsdk:"result_i"`
	ResultStr  types.String `tfsdk:"result_str"`
	SessionId  types.String `tfsdk:"session_id"`
	StartTs    types.String `tfsdk:"start_ts"`
	StatusI    types.String `tfsdk:"status_i"`
	StatusStr  types.String `tfsdk:"status_str"`
	Summary    types.String `tfsdk:"summary"`
	TypeI      types.String `tfsdk:"type_i"`
	TypeStr    types.String `tfsdk:"type_str"`
	Uname      types.String `tfsdk:"uname"`
}

// Metadata returns the data source type name.
func (d *jobsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jobs"
}

// Schema defines the schema for this listing data source.
func (d *jobsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Retrieves config for a specific item.",

		Attributes: map[string]dsschema.Attribute{
			"id": dsschema.StringAttribute{
				Description:         "The object ID.",
				MarkdownDescription: "The object ID.",
				Computed:            true,
			},

			// Input.
			"job_id": dsschema.StringAttribute{
				Description:         "The id of the job.",
				MarkdownDescription: "The id of the job.",
				Required:            true,
			},

			// Output.
			"details": dsschema.StringAttribute{
				Description:         "The `details` parameter.",
				MarkdownDescription: "The `details` parameter.",
				Computed:            true,
			},
			"end_ts": dsschema.StringAttribute{
				Description:         "The `end_ts` parameter.",
				MarkdownDescription: "The `end_ts` parameter.",
				Computed:            true,
			},
			"object_id": dsschema.StringAttribute{
				Description:         "The `object_id` parameter.",
				MarkdownDescription: "The `object_id` parameter.",
				Computed:            true,
			},
			"insert_ts": dsschema.StringAttribute{
				Description:         "The `insert_ts` parameter.",
				MarkdownDescription: "The `insert_ts` parameter.",
				Computed:            true,
			},
			"job_result": dsschema.StringAttribute{
				Description:         "The `job_result` parameter.",
				MarkdownDescription: "The `job_result` parameter.",
				Computed:            true,
			},
			"job_status": dsschema.StringAttribute{
				Description:         "The `job_status` parameter.",
				MarkdownDescription: "The `job_status` parameter.",
				Computed:            true,
			},
			"job_type": dsschema.StringAttribute{
				Description:         "The `job_type` parameter.",
				MarkdownDescription: "The `job_type` parameter.",
				Computed:            true,
			},
			"last_update": dsschema.StringAttribute{
				Description:         "The `last_update` parameter.",
				MarkdownDescription: "The `last_update` parameter.",
				Computed:            true,
			},
			"opaque_int": dsschema.StringAttribute{
				Description:         "The `opaque_int` parameter.",
				MarkdownDescription: "The `opaque_int` parameter.",
				Computed:            true,
			},
			"opaque_str": dsschema.StringAttribute{
				Description:         "The `opaque_str` parameter.",
				MarkdownDescription: "The `opaque_str` parameter.",
				Computed:            true,
			},
			"owner": dsschema.StringAttribute{
				Description:         "The `owner` parameter.",
				MarkdownDescription: "The `owner` parameter.",
				Computed:            true,
			},
			"parent_id": dsschema.StringAttribute{
				Description:         "The `parent_id` parameter.",
				MarkdownDescription: "The `parent_id` parameter.",
				Computed:            true,
			},
			"percent": dsschema.StringAttribute{
				Description:         "The `percent` parameter.",
				MarkdownDescription: "The `percent` parameter.",
				Computed:            true,
			},
			"result_i": dsschema.StringAttribute{
				Description:         "The `result_i` parameter.",
				MarkdownDescription: "The `result_i` parameter.",
				Computed:            true,
			},
			"result_str": dsschema.StringAttribute{
				Description:         "The `result_str` parameter.",
				MarkdownDescription: "The `result_str` parameter.",
				Computed:            true,
			},
			"session_id": dsschema.StringAttribute{
				Description:         "The `session_id` parameter.",
				MarkdownDescription: "The `session_id` parameter.",
				Computed:            true,
			},
			"start_ts": dsschema.StringAttribute{
				Description:         "The `start_ts` parameter.",
				MarkdownDescription: "The `start_ts` parameter.",
				Computed:            true,
			},
			"status_i": dsschema.StringAttribute{
				Description:         "The `status_i` parameter.",
				MarkdownDescription: "The `status_i` parameter.",
				Computed:            true,
			},
			"status_str": dsschema.StringAttribute{
				Description:         "The `status_str` parameter.",
				MarkdownDescription: "The `status_str` parameter.",
				Computed:            true,
			},
			"summary": dsschema.StringAttribute{
				Description:         "The `summary` parameter.",
				MarkdownDescription: "The `summary` parameter.",
				Computed:            true,
			},
			"type_i": dsschema.StringAttribute{
				Description:         "The `type_i` parameter.",
				MarkdownDescription: "The `type_i` parameter.",
				Computed:            true,
			},
			"type_str": dsschema.StringAttribute{
				Description:         "The `type_str` parameter.",
				MarkdownDescription: "The `type_str` parameter.",
				Computed:            true,
			},
			"uname": dsschema.StringAttribute{
				Description:         "The `uname` parameter.",
				MarkdownDescription: "The `uname` parameter.",
				Computed:            true,
			},
		},
	}
}

// Configure prepares the struct.
func (d *jobsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *jobsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state jobsDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_jobs",
		"job_id":                      state.JobId.ValueString(),
	})

	// Prepare to run the command.
	svc := wugpput.NewClient(d.client)
	input := wugpput.ReadInput{
		JobId: state.JobId.ValueString(),
	}

	// Perform the operation.
	ans, err := svc.Read(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error getting singleton", err.Error())
		return
	}

	// Store the answer to state.
	var idBuilder strings.Builder
	idBuilder.WriteString(input.JobId)
	state.Id = types.StringValue(idBuilder.String())
	state.Details = types.StringValue(ans.Details)
	state.EndTs = types.StringValue(ans.EndTs)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.InsertTs = types.StringValue(ans.InsertTs)
	state.JobResult = types.StringValue(ans.JobResult)
	state.JobStatus = types.StringValue(ans.JobStatus)
	state.JobType = types.StringValue(ans.JobType)
	state.LastUpdate = types.StringValue(ans.LastUpdate)
	state.OpaqueInt = types.StringValue(ans.OpaqueInt)
	state.OpaqueStr = types.StringValue(ans.OpaqueStr)
	state.Owner = types.StringValue(ans.Owner)
	state.ParentId = types.StringValue(ans.ParentId)
	state.Percent = types.StringValue(ans.Percent)
	state.ResultI = types.StringValue(ans.ResultI)
	state.ResultStr = types.StringValue(ans.ResultStr)
	state.SessionId = types.StringValue(ans.SessionId)
	state.StartTs = types.StringValue(ans.StartTs)
	state.StatusI = types.StringValue(ans.StatusI)
	state.StatusStr = types.StringValue(ans.StatusStr)
	state.Summary = types.StringValue(ans.Summary)
	state.TypeI = types.StringValue(ans.TypeI)
	state.TypeStr = types.StringValue(ans.TypeStr)
	state.Uname = types.StringValue(ans.Uname)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
