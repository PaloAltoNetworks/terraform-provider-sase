package provider

import (
	"context"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	pVYGFzR "github.com/paloaltonetworks/sase-go/netsec/service/v1/mobileagent/locations"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Data source (listing).
var (
	_ datasource.DataSource              = &mobileAgentLocationsListDataSource{}
	_ datasource.DataSourceWithConfigure = &mobileAgentLocationsListDataSource{}
)

func NewMobileAgentLocationsListDataSource() datasource.DataSource {
	return &mobileAgentLocationsListDataSource{}
}

type mobileAgentLocationsListDataSource struct {
	client *sase.Client
}

type mobileAgentLocationsListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data   []mobileAgentLocationsListDsModelConfig `tfsdk:"data"`
	Limit  types.Int64                             `tfsdk:"limit"`
	Offset types.Int64                             `tfsdk:"offset"`
	Total  types.Int64                             `tfsdk:"total"`
}

type mobileAgentLocationsListDsModelConfig struct {
	Region []mobileAgentLocationsListDsModelRegionObject `tfsdk:"region"`
}

type mobileAgentLocationsListDsModelRegionObject struct {
	Locations []types.String `tfsdk:"locations"`
	Name      types.String   `tfsdk:"name"`
}

// Metadata returns the data source type name.
func (d *mobileAgentLocationsListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mobile_agent_locations_list"
}

// Schema defines the schema for this listing data source.
func (d *mobileAgentLocationsListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Retrieves a listing of config items.",

		Attributes: map[string]dsschema.Attribute{
			"id": dsschema.StringAttribute{
				Description:         "The object ID.",
				MarkdownDescription: "The object ID.",
				Computed:            true,
			},

			// Input.
			"folder": dsschema.StringAttribute{
				Description:         "The `folder` parameter. Value must be one of: `\"Mobile Users\"`.",
				MarkdownDescription: "The `folder` parameter. Value must be one of: `\"Mobile Users\"`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Mobile Users"),
				},
			},

			// Output.
			"data": dsschema.ListNestedAttribute{
				Description:         "The `data` parameter.",
				MarkdownDescription: "The `data` parameter.",
				Computed:            true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"region": dsschema.ListNestedAttribute{
							Description:         "The `region` parameter.",
							MarkdownDescription: "The `region` parameter.",
							Computed:            true,
							NestedObject: dsschema.NestedAttributeObject{
								Attributes: map[string]dsschema.Attribute{
									"locations": dsschema.ListAttribute{
										Description:         "The `locations` parameter.",
										MarkdownDescription: "The `locations` parameter.",
										Computed:            true,
										ElementType:         types.StringType,
									},
									"name": dsschema.StringAttribute{
										Description:         "The `name` parameter.",
										MarkdownDescription: "The `name` parameter.",
										Computed:            true,
									},
								},
							},
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
func (d *mobileAgentLocationsListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *mobileAgentLocationsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state mobileAgentLocationsListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_mobile_agent_locations_list",
		"terraform_provider_function": "Read",
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := pVYGFzR.NewClient(d.client)
	input := pVYGFzR.ListInput{
		Folder: state.Folder.ValueString(),
	}

	// Perform the operation.
	ans, err := svc.List(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error getting listing", err.Error())
		return
	}

	// Store the answer to state.
	var idBuilder strings.Builder
	idBuilder.WriteString(input.Folder)
	state.Id = types.StringValue(idBuilder.String())
	var var0 []mobileAgentLocationsListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]mobileAgentLocationsListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 mobileAgentLocationsListDsModelConfig
			var var3 []mobileAgentLocationsListDsModelRegionObject
			if len(var1.Region) != 0 {
				var3 = make([]mobileAgentLocationsListDsModelRegionObject, 0, len(var1.Region))
				for var4Index := range var1.Region {
					var4 := var1.Region[var4Index]
					var var5 mobileAgentLocationsListDsModelRegionObject
					var5.Locations = EncodeStringSlice(var4.Locations)
					var5.Name = types.StringValue(var4.Name)
					var3 = append(var3, var5)
				}
			}
			var2.Region = var3
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
