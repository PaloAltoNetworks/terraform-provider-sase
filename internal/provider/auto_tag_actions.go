package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	iYmUVvF "github.com/paloaltonetworks/sase-go/netsec/service/v1/autotagactions"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Data source (listing).
var (
	_ datasource.DataSource              = &autoTagActionsListDataSource{}
	_ datasource.DataSourceWithConfigure = &autoTagActionsListDataSource{}
)

func NewAutoTagActionsListDataSource() datasource.DataSource {
	return &autoTagActionsListDataSource{}
}

type autoTagActionsListDataSource struct {
	client *sase.Client
}

type autoTagActionsListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []autoTagActionsListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type autoTagActionsListDsModelConfig struct {
	Actions        []autoTagActionsListDsModelActionsObject `tfsdk:"actions"`
	Description    types.String                             `tfsdk:"description"`
	Filter         types.String                             `tfsdk:"filter"`
	ObjectId       types.String                             `tfsdk:"object_id"`
	LogType        types.String                             `tfsdk:"log_type"`
	Name           types.String                             `tfsdk:"name"`
	Quarantine     types.Bool                               `tfsdk:"quarantine"`
	SendToPanorama types.Bool                               `tfsdk:"send_to_panorama"`
}

type autoTagActionsListDsModelActionsObject struct {
	Name types.String                        `tfsdk:"name"`
	Type autoTagActionsListDsModelTypeObject `tfsdk:"type"`
}

type autoTagActionsListDsModelTypeObject struct {
	Tagging autoTagActionsListDsModelTaggingObject `tfsdk:"tagging"`
}

type autoTagActionsListDsModelTaggingObject struct {
	Action  types.String   `tfsdk:"action"`
	Tags    []types.String `tfsdk:"tags"`
	Target  types.String   `tfsdk:"target"`
	Timeout types.Int64    `tfsdk:"timeout"`
}

// Metadata returns the data source type name.
func (d *autoTagActionsListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auto_tag_actions_list"
}

// Schema defines the schema for this listing data source.
func (d *autoTagActionsListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
					stringvalidator.OneOf("Shared"),
				},
			},

			// Output.
			"data": dsschema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"actions": dsschema.ListNestedAttribute{
							Description: "",
							Computed:    true,
							NestedObject: dsschema.NestedAttributeObject{
								Attributes: map[string]dsschema.Attribute{
									"name": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"type": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"tagging": dsschema.SingleNestedAttribute{
												Description: "",
												Computed:    true,
												Attributes: map[string]dsschema.Attribute{
													"action": dsschema.StringAttribute{
														Description: "",
														Computed:    true,
													},
													"tags": dsschema.ListAttribute{
														Description: "",
														Computed:    true,
														ElementType: types.StringType,
													},
													"target": dsschema.StringAttribute{
														Description: "",
														Computed:    true,
													},
													"timeout": dsschema.Int64Attribute{
														Description: "",
														Computed:    true,
													},
												},
											},
										},
									},
								},
							},
						},
						"description": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"filter": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"object_id": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"log_type": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"quarantine": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"send_to_panorama": dsschema.BoolAttribute{
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
func (d *autoTagActionsListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *autoTagActionsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state autoTagActionsListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_auto_tag_actions_list",
		"limit":            state.Limit.ValueInt64(),
		"has_limit":        !state.Limit.IsNull(),
		"offset":           state.Offset.ValueInt64(),
		"has_offset":       !state.Offset.IsNull(),
		"name":             state.Name.ValueString(),
		"has_name":         !state.Name.IsNull(),
		"folder":           state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := iYmUVvF.NewClient(d.client)
	input := iYmUVvF.ListInput{
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
	var var0 []autoTagActionsListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]autoTagActionsListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 autoTagActionsListDsModelConfig
			var var3 []autoTagActionsListDsModelActionsObject
			if len(var1.Actions) != 0 {
				var3 = make([]autoTagActionsListDsModelActionsObject, 0, len(var1.Actions))
				for var4Index := range var1.Actions {
					var4 := var1.Actions[var4Index]
					var var5 autoTagActionsListDsModelActionsObject
					var var6 autoTagActionsListDsModelTypeObject
					var var7 autoTagActionsListDsModelTaggingObject
					var7.Action = types.StringValue(var4.Type.Tagging.Action)
					var7.Tags = EncodeStringSlice(var4.Type.Tagging.Tags)
					var7.Target = types.StringValue(var4.Type.Tagging.Target)
					var7.Timeout = types.Int64Value(var4.Type.Tagging.Timeout)
					var6.Tagging = var7
					var5.Name = types.StringValue(var4.Name)
					var5.Type = var6
					var3 = append(var3, var5)
				}
			}
			var2.Actions = var3
			var2.Description = types.StringValue(var1.Description)
			var2.Filter = types.StringValue(var1.Filter)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.LogType = types.StringValue(var1.LogType)
			var2.Name = types.StringValue(var1.Name)
			var2.Quarantine = types.BoolValue(var1.Quarantine)
			var2.SendToPanorama = types.BoolValue(var1.SendToPanorama)
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
