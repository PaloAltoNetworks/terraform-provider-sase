package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	mfYmVgm "github.com/paloaltonetworks/sase-go/netsec/service/v1/authenticationportals"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Data source (listing).
var (
	_ datasource.DataSource              = &authenticationPortalsListDataSource{}
	_ datasource.DataSourceWithConfigure = &authenticationPortalsListDataSource{}
)

func NewAuthenticationPortalsListDataSource() datasource.DataSource {
	return &authenticationPortalsListDataSource{}
}

type authenticationPortalsListDataSource struct {
	client *sase.Client
}

type authenticationPortalsListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []authenticationPortalsListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type authenticationPortalsListDsModelConfig struct {
	AuthenticationProfile types.String `tfsdk:"authentication_profile"`
	CertificateProfile    types.String `tfsdk:"certificate_profile"`
	GpUdpPort             types.Int64  `tfsdk:"gp_udp_port"`
	ObjectId              types.String `tfsdk:"object_id"`
	IdleTimer             types.Int64  `tfsdk:"idle_timer"`
	RedirectHost          types.String `tfsdk:"redirect_host"`
	Timer                 types.Int64  `tfsdk:"timer"`
	TlsServiceProfile     types.String `tfsdk:"tls_service_profile"`
}

// Metadata returns the data source type name.
func (d *authenticationPortalsListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_portals_list"
}

// Schema defines the schema for this listing data source.
func (d *authenticationPortalsListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Retrieves a listing of config items.",

		Attributes: map[string]dsschema.Attribute{
			"id": dsschema.StringAttribute{
				Description:         "The object ID.",
				MarkdownDescription: "The object ID.",
				Computed:            true,
			},

			// Input.
			"limit": dsschema.Int64Attribute{
				Description:         "The max count in result entry (count per page)",
				MarkdownDescription: "The max count in result entry (count per page)",
				Optional:            true,
				Computed:            true,
			},
			"offset": dsschema.Int64Attribute{
				Description:         "The offset of the result entry",
				MarkdownDescription: "The offset of the result entry",
				Optional:            true,
				Computed:            true,
			},
			"folder": dsschema.StringAttribute{
				Description:         "The folder of the entry",
				MarkdownDescription: "The folder of the entry",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Service Connections", "Mobile Users Container", "Mobile Users Explicit Proxy"),
				},
			},

			// Output.
			"data": dsschema.ListNestedAttribute{
				Description:         "The `data` parameter.",
				MarkdownDescription: "The `data` parameter.",
				Computed:            true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"authentication_profile": dsschema.StringAttribute{
							Description:         "The `authentication_profile` parameter.",
							MarkdownDescription: "The `authentication_profile` parameter.",
							Computed:            true,
						},
						"certificate_profile": dsschema.StringAttribute{
							Description:         "The `certificate_profile` parameter.",
							MarkdownDescription: "The `certificate_profile` parameter.",
							Computed:            true,
						},
						"gp_udp_port": dsschema.Int64Attribute{
							Description:         "The `gp_udp_port` parameter.",
							MarkdownDescription: "The `gp_udp_port` parameter.",
							Computed:            true,
						},
						"object_id": dsschema.StringAttribute{
							Description:         "The `object_id` parameter.",
							MarkdownDescription: "The `object_id` parameter.",
							Computed:            true,
						},
						"idle_timer": dsschema.Int64Attribute{
							Description:         "The `idle_timer` parameter.",
							MarkdownDescription: "The `idle_timer` parameter.",
							Computed:            true,
						},
						"redirect_host": dsschema.StringAttribute{
							Description:         "The `redirect_host` parameter.",
							MarkdownDescription: "The `redirect_host` parameter.",
							Computed:            true,
						},
						"timer": dsschema.Int64Attribute{
							Description:         "The `timer` parameter.",
							MarkdownDescription: "The `timer` parameter.",
							Computed:            true,
						},
						"tls_service_profile": dsschema.StringAttribute{
							Description:         "The `tls_service_profile` parameter.",
							MarkdownDescription: "The `tls_service_profile` parameter.",
							Computed:            true,
						},
					},
				},
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
func (d *authenticationPortalsListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *authenticationPortalsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state authenticationPortalsListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_authentication_portals_list",
		"terraform_provider_function": "Read",
		"limit":                       state.Limit.ValueInt64(),
		"has_limit":                   !state.Limit.IsNull(),
		"offset":                      state.Offset.ValueInt64(),
		"has_offset":                  !state.Offset.IsNull(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := mfYmVgm.NewClient(d.client)
	input := mfYmVgm.ListInput{
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
	var var0 []authenticationPortalsListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]authenticationPortalsListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 authenticationPortalsListDsModelConfig
			var2.AuthenticationProfile = types.StringValue(var1.AuthenticationProfile)
			var2.CertificateProfile = types.StringValue(var1.CertificateProfile)
			var2.GpUdpPort = types.Int64Value(var1.GpUdpPort)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.IdleTimer = types.Int64Value(var1.IdleTimer)
			var2.RedirectHost = types.StringValue(var1.RedirectHost)
			var2.Timer = types.Int64Value(var1.Timer)
			var2.TlsServiceProfile = types.StringValue(var1.TlsServiceProfile)
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
