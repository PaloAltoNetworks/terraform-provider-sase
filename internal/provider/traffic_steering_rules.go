package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	wWVKIJO "github.com/paloaltonetworks/sase-go/netsec/service/v1/trafficsteeringrules"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Data source (listing).
var (
	_ datasource.DataSource              = &trafficSteeringRulesListDataSource{}
	_ datasource.DataSourceWithConfigure = &trafficSteeringRulesListDataSource{}
)

func NewTrafficSteeringRulesListDataSource() datasource.DataSource {
	return &trafficSteeringRulesListDataSource{}
}

type trafficSteeringRulesListDataSource struct {
	client *sase.Client
}

type trafficSteeringRulesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []trafficSteeringRulesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type trafficSteeringRulesListDsModelConfig struct {
	Action      *trafficSteeringRulesListDsModelActionObject `tfsdk:"action"`
	Category    []types.String                               `tfsdk:"category"`
	Destination []types.String                               `tfsdk:"destination"`
	Name        types.String                                 `tfsdk:"name"`
	Service     []types.String                               `tfsdk:"service"`
	Source      []types.String                               `tfsdk:"source"`
	SourceUser  []types.String                               `tfsdk:"source_user"`
}

type trafficSteeringRulesListDsModelActionObject struct {
	Forward *trafficSteeringRulesListDsModelForwardObject `tfsdk:"forward"`
	NoPbf   types.Bool                                    `tfsdk:"no_pbf"`
}

type trafficSteeringRulesListDsModelForwardObject struct {
	Target types.String `tfsdk:"target"`
}

// Metadata returns the data source type name.
func (d *trafficSteeringRulesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_traffic_steering_rules_list"
}

// Schema defines the schema for this listing data source.
func (d *trafficSteeringRulesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"action": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"forward": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"target": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
								"no_pbf": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
							},
						},
						"category": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"destination": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"service": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"source": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"source_user": dsschema.ListAttribute{
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
func (d *trafficSteeringRulesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *trafficSteeringRulesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state trafficSteeringRulesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_traffic_steering_rules_list",
		"limit":            state.Limit.ValueInt64(),
		"has_limit":        !state.Limit.IsNull(),
		"offset":           state.Offset.ValueInt64(),
		"has_offset":       !state.Offset.IsNull(),
		"name":             state.Name.ValueString(),
		"has_name":         !state.Name.IsNull(),
		"folder":           state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := wWVKIJO.NewClient(d.client)
	input := wWVKIJO.ListInput{
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
	var var0 []trafficSteeringRulesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]trafficSteeringRulesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 trafficSteeringRulesListDsModelConfig
			var var3 *trafficSteeringRulesListDsModelActionObject
			if var1.Action != nil {
				var3 = &trafficSteeringRulesListDsModelActionObject{}
				var var4 *trafficSteeringRulesListDsModelForwardObject
				if var1.Action.Forward != nil {
					var4 = &trafficSteeringRulesListDsModelForwardObject{}
					var4.Target = types.StringValue(var1.Action.Forward.Target)
				}
				var3.Forward = var4
				if var1.Action.NoPbf != nil {
					var3.NoPbf = types.BoolValue(true)
				}
			}
			var2.Action = var3
			var2.Category = EncodeStringSlice(var1.Category)
			var2.Destination = EncodeStringSlice(var1.Destination)
			var2.Name = types.StringValue(var1.Name)
			var2.Service = EncodeStringSlice(var1.Service)
			var2.Source = EncodeStringSlice(var1.Source)
			var2.SourceUser = EncodeStringSlice(var1.SourceUser)
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
