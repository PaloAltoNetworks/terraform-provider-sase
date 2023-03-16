package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	zDUyfEt "github.com/paloaltonetworks/sase-go/netsec/service/v1/authenticationrules"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Data source (listing).
var (
	_ datasource.DataSource              = &authenticationRulesListDataSource{}
	_ datasource.DataSourceWithConfigure = &authenticationRulesListDataSource{}
)

func NewAuthenticationRulesListDataSource() datasource.DataSource {
	return &authenticationRulesListDataSource{}
}

type authenticationRulesListDataSource struct {
	client *sase.Client
}

type authenticationRulesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit    types.Int64  `tfsdk:"limit"`
	Offset   types.Int64  `tfsdk:"offset"`
	Position types.String `tfsdk:"position"`
	Folder   types.String `tfsdk:"folder"`
	Name     types.String `tfsdk:"name"`

	// Output.
	Data []authenticationRulesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type authenticationRulesListDsModelConfig struct {
	AuthenticationEnforcement types.String   `tfsdk:"authentication_enforcement"`
	Category                  []types.String `tfsdk:"category"`
	Description               types.String   `tfsdk:"description"`
	Destination               []types.String `tfsdk:"destination"`
	DestinationHip            []types.String `tfsdk:"destination_hip"`
	Disabled                  types.Bool     `tfsdk:"disabled"`
	From                      []types.String `tfsdk:"from"`
	GroupTag                  types.String   `tfsdk:"group_tag"`
	HipProfiles               []types.String `tfsdk:"hip_profiles"`
	ObjectId                  types.String   `tfsdk:"object_id"`
	LogAuthenticationTimeout  types.Bool     `tfsdk:"log_authentication_timeout"`
	LogSetting                types.String   `tfsdk:"log_setting"`
	Name                      types.String   `tfsdk:"name"`
	NegateDestination         types.Bool     `tfsdk:"negate_destination"`
	NegateSource              types.Bool     `tfsdk:"negate_source"`
	Service                   []types.String `tfsdk:"service"`
	Source                    []types.String `tfsdk:"source"`
	SourceHip                 []types.String `tfsdk:"source_hip"`
	SourceUser                []types.String `tfsdk:"source_user"`
	Tag                       []types.String `tfsdk:"tag"`
	Timeout                   types.Int64    `tfsdk:"timeout"`
	To                        []types.String `tfsdk:"to"`
}

// Metadata returns the data source type name.
func (d *authenticationRulesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_rules_list"
}

// Schema defines the schema for this listing data source.
func (d *authenticationRulesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"authentication_enforcement": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"category": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
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
						"destination_hip": dsschema.ListAttribute{
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
						"hip_profiles": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"object_id": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"log_authentication_timeout": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"log_setting": dsschema.StringAttribute{
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
						"source_hip": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"source_user": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"tag": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"timeout": dsschema.Int64Attribute{
							Description: "",
							Computed:    true,
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
func (d *authenticationRulesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *authenticationRulesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state authenticationRulesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_authentication_rules_list",
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
	svc := zDUyfEt.NewClient(d.client)
	input := zDUyfEt.ListInput{
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
	idBuilder.WriteString(input.Position)
	idBuilder.WriteString(IdSeparator)
	idBuilder.WriteString(input.Folder)
	idBuilder.WriteString(IdSeparator)
	if input.Name != nil {
		idBuilder.WriteString(*input.Name)
	}
	state.Id = types.StringValue(idBuilder.String())
	var var0 []authenticationRulesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]authenticationRulesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 authenticationRulesListDsModelConfig
			var2.AuthenticationEnforcement = types.StringValue(var1.AuthenticationEnforcement)
			var2.Category = EncodeStringSlice(var1.Category)
			var2.Description = types.StringValue(var1.Description)
			var2.Destination = EncodeStringSlice(var1.Destination)
			var2.DestinationHip = EncodeStringSlice(var1.DestinationHip)
			var2.Disabled = types.BoolValue(var1.Disabled)
			var2.From = EncodeStringSlice(var1.From)
			var2.GroupTag = types.StringValue(var1.GroupTag)
			var2.HipProfiles = EncodeStringSlice(var1.HipProfiles)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.LogAuthenticationTimeout = types.BoolValue(var1.LogAuthenticationTimeout)
			var2.LogSetting = types.StringValue(var1.LogSetting)
			var2.Name = types.StringValue(var1.Name)
			var2.NegateDestination = types.BoolValue(var1.NegateDestination)
			var2.NegateSource = types.BoolValue(var1.NegateSource)
			var2.Service = EncodeStringSlice(var1.Service)
			var2.Source = EncodeStringSlice(var1.Source)
			var2.SourceHip = EncodeStringSlice(var1.SourceHip)
			var2.SourceUser = EncodeStringSlice(var1.SourceUser)
			var2.Tag = EncodeStringSlice(var1.Tag)
			var2.Timeout = types.Int64Value(var1.Timeout)
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
