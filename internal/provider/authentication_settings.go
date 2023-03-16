package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	uQwObPt "github.com/paloaltonetworks/sase-go/netsec/service/v1/mobileagent/authenticationsettings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Data source (listing).
var (
	_ datasource.DataSource              = &authenticationSettingsListDataSource{}
	_ datasource.DataSourceWithConfigure = &authenticationSettingsListDataSource{}
)

func NewAuthenticationSettingsListDataSource() datasource.DataSource {
	return &authenticationSettingsListDataSource{}
}

type authenticationSettingsListDataSource struct {
	client *sase.Client
}

type authenticationSettingsListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []authenticationSettingsListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type authenticationSettingsListDsModelConfig struct {
	AuthenticationProfile              types.String `tfsdk:"authentication_profile"`
	Os                                 types.String `tfsdk:"os"`
	UserCredentialOrClientCertRequired types.Bool   `tfsdk:"user_credential_or_client_cert_required"`
}

// Metadata returns the data source type name.
func (d *authenticationSettingsListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_settings_list"
}

// Schema defines the schema for this listing data source.
func (d *authenticationSettingsListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"folder": dsschema.StringAttribute{
				Description: "",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Mobile Users"),
				},
			},

			// Output.
			"data": dsschema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"authentication_profile": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"os": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"user_credential_or_client_cert_required": dsschema.BoolAttribute{
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
func (d *authenticationSettingsListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *authenticationSettingsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state authenticationSettingsListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_authentication_settings_list",
		"limit":            state.Limit.ValueInt64(),
		"has_limit":        !state.Limit.IsNull(),
		"offset":           state.Offset.ValueInt64(),
		"has_offset":       !state.Offset.IsNull(),
		"folder":           state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := uQwObPt.NewClient(d.client)
	input := uQwObPt.ListInput{
		Folder: state.Folder.ValueString(),
	}
	if !state.Limit.IsNull() {
		input.Limit = api.Int(state.Limit.ValueInt64())
	}
	if !state.Offset.IsNull() {
		input.Offset = api.Int(state.Offset.ValueInt64())
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
	idBuilder.WriteString(input.Folder)
	state.Id = types.StringValue(idBuilder.String())
	var var0 []authenticationSettingsListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]authenticationSettingsListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 authenticationSettingsListDsModelConfig
			var2.AuthenticationProfile = types.StringValue(var1.AuthenticationProfile)
			var2.Os = types.StringValue(var1.Os)
			var2.UserCredentialOrClientCertRequired = types.BoolValue(var1.UserCredentialOrClientCertRequired)
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
