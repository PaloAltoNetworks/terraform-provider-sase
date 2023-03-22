package provider

import (
	"context"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	qOFkTUB "github.com/paloaltonetworks/sase-go/netsec/service/v1/configversions"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Data source.
var (
	_ datasource.DataSource              = &candidateConfigVersionsDataSource{}
	_ datasource.DataSourceWithConfigure = &candidateConfigVersionsDataSource{}
)

func NewCandidateConfigVersionsDataSource() datasource.DataSource {
	return &candidateConfigVersionsDataSource{}
}

type candidateConfigVersionsDataSource struct {
	client *sase.Client
}

type candidateConfigVersionsDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Version types.String `tfsdk:"version"`

	// Output.
	// Ref: #/components/schemas/candidate-config-versions
	Admin       types.String `tfsdk:"admin"`
	Created     types.Int64  `tfsdk:"created"`
	Date        types.String `tfsdk:"date"`
	Deleted     types.Int64  `tfsdk:"deleted"`
	Description types.String `tfsdk:"description"`
	ObjectId    types.Int64  `tfsdk:"object_id"`
	Scope       types.String `tfsdk:"scope"`
	SwgConfig   types.String `tfsdk:"swg_config"`
	Updated     types.Int64  `tfsdk:"updated"`
	// input omit: Version
}

// Metadata returns the data source type name.
func (d *candidateConfigVersionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_candidate_config_versions"
}

// Schema defines the schema for this listing data source.
func (d *candidateConfigVersionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Retrieves config for a specific item.",

		Attributes: map[string]dsschema.Attribute{
			"id": dsschema.StringAttribute{
				Description: "The object ID.",
				Computed:    true,
			},

			// Input.
			"version": dsschema.StringAttribute{
				Description: "The version of the running config",
				Required:    true,
			},

			// Output.
			"admin": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"created": dsschema.Int64Attribute{
				Description: "",
				Computed:    true,
			},
			"date": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"deleted": dsschema.Int64Attribute{
				Description: "",
				Computed:    true,
			},
			"description": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"object_id": dsschema.Int64Attribute{
				Description: "",
				Computed:    true,
			},
			"scope": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"swg_config": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"updated": dsschema.Int64Attribute{
				Description: "",
				Computed:    true,
			},
		},
	}
}

// Configure prepares the struct.
func (d *candidateConfigVersionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *candidateConfigVersionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state candidateConfigVersionsDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_candidate_config_versions",
		"version":                     state.Version.ValueString(),
	})

	// Prepare to run the command.
	svc := qOFkTUB.NewClient(d.client)
	input := qOFkTUB.ReadInput{
		Version: state.Version.ValueString(),
	}

	// Perform the operation.
	ans, err := svc.Read(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error getting singleton", err.Error())
		return
	}

	// Store the answer to state.
	var idBuilder strings.Builder
	idBuilder.WriteString(input.Version)
	state.Id = types.StringValue(idBuilder.String())
	state.Admin = types.StringValue(ans.Admin)
	state.Created = types.Int64Value(ans.Created)
	state.Date = types.StringValue(ans.Date)
	state.Deleted = types.Int64Value(ans.Deleted)
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.Int64Value(ans.ObjectId)
	state.Scope = types.StringValue(ans.Scope)
	state.SwgConfig = types.StringValue(ans.SwgConfig)
	state.Updated = types.Int64Value(ans.Updated)
	if !state.Version.IsNull() {
		state.Version = types.StringValue(ans.Version)
	}

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
