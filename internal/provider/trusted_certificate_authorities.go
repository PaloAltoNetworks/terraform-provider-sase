package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	jOaWaMY "github.com/paloaltonetworks/sase-go/netsec/service/v1/trustedcertificateauthorities"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Data source (listing).
var (
	_ datasource.DataSource              = &trustedCertificateAuthoritiesListDataSource{}
	_ datasource.DataSourceWithConfigure = &trustedCertificateAuthoritiesListDataSource{}
)

func NewTrustedCertificateAuthoritiesListDataSource() datasource.DataSource {
	return &trustedCertificateAuthoritiesListDataSource{}
}

type trustedCertificateAuthoritiesListDataSource struct {
	client *sase.Client
}

type trustedCertificateAuthoritiesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []trustedCertificateAuthoritiesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type trustedCertificateAuthoritiesListDsModelConfig struct {
	CommonName     types.String `tfsdk:"common_name"`
	ExpiryEpoch    types.String `tfsdk:"expiry_epoch"`
	Filename       types.String `tfsdk:"filename"`
	ObjectId       types.String `tfsdk:"object_id"`
	Issuer         types.String `tfsdk:"issuer"`
	Name           types.String `tfsdk:"name"`
	NotValidAfter  types.String `tfsdk:"not_valid_after"`
	NotValidBefore types.String `tfsdk:"not_valid_before"`
	SerialNumber   types.String `tfsdk:"serial_number"`
	Subject        types.String `tfsdk:"subject"`
}

// Metadata returns the data source type name.
func (d *trustedCertificateAuthoritiesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trusted_certificate_authorities_list"
}

// Schema defines the schema for this listing data source.
func (d *trustedCertificateAuthoritiesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
				Description:         "The max count in result entry (count per page).",
				MarkdownDescription: "The max count in result entry (count per page).",
				Optional:            true,
				Computed:            true,
			},
			"offset": dsschema.Int64Attribute{
				Description:         "The offset of the result entry.",
				MarkdownDescription: "The offset of the result entry.",
				Optional:            true,
				Computed:            true,
			},
			"name": dsschema.StringAttribute{
				Description:         "The name of the entry.",
				MarkdownDescription: "The name of the entry.",
				Optional:            true,
			},
			"folder": dsschema.StringAttribute{
				Description:         "The folder of the entry. Value must be one of: `\"Shared\"`, `\"Mobile Users\"`, `\"Remote Networks\"`, `\"Service Connections\"`, `\"Mobile Users Container\"`, `\"Mobile Users Explicit Proxy\"`.",
				MarkdownDescription: "The folder of the entry. Value must be one of: `\"Shared\"`, `\"Mobile Users\"`, `\"Remote Networks\"`, `\"Service Connections\"`, `\"Mobile Users Container\"`, `\"Mobile Users Explicit Proxy\"`.",
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
						"common_name": dsschema.StringAttribute{
							Description:         "The `common_name` parameter.",
							MarkdownDescription: "The `common_name` parameter.",
							Computed:            true,
						},
						"expiry_epoch": dsschema.StringAttribute{
							Description:         "The `expiry_epoch` parameter.",
							MarkdownDescription: "The `expiry_epoch` parameter.",
							Computed:            true,
						},
						"filename": dsschema.StringAttribute{
							Description:         "The `filename` parameter.",
							MarkdownDescription: "The `filename` parameter.",
							Computed:            true,
						},
						"object_id": dsschema.StringAttribute{
							Description:         "The `object_id` parameter.",
							MarkdownDescription: "The `object_id` parameter.",
							Computed:            true,
						},
						"issuer": dsschema.StringAttribute{
							Description:         "The `issuer` parameter.",
							MarkdownDescription: "The `issuer` parameter.",
							Computed:            true,
						},
						"name": dsschema.StringAttribute{
							Description:         "The `name` parameter.",
							MarkdownDescription: "The `name` parameter.",
							Computed:            true,
						},
						"not_valid_after": dsschema.StringAttribute{
							Description:         "The `not_valid_after` parameter.",
							MarkdownDescription: "The `not_valid_after` parameter.",
							Computed:            true,
						},
						"not_valid_before": dsschema.StringAttribute{
							Description:         "The `not_valid_before` parameter.",
							MarkdownDescription: "The `not_valid_before` parameter.",
							Computed:            true,
						},
						"serial_number": dsschema.StringAttribute{
							Description:         "The `serial_number` parameter.",
							MarkdownDescription: "The `serial_number` parameter.",
							Computed:            true,
						},
						"subject": dsschema.StringAttribute{
							Description:         "The `subject` parameter.",
							MarkdownDescription: "The `subject` parameter.",
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
func (d *trustedCertificateAuthoritiesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *trustedCertificateAuthoritiesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state trustedCertificateAuthoritiesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_trusted_certificate_authorities_list",
		"terraform_provider_function": "Read",
		"limit":                       state.Limit.ValueInt64(),
		"has_limit":                   !state.Limit.IsNull(),
		"offset":                      state.Offset.ValueInt64(),
		"has_offset":                  !state.Offset.IsNull(),
		"name":                        state.Name.ValueString(),
		"has_name":                    !state.Name.IsNull(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := jOaWaMY.NewClient(d.client)
	input := jOaWaMY.ListInput{
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
	if input.Name != nil {
		idBuilder.WriteString(*input.Name)
	}
	idBuilder.WriteString(IdSeparator)
	idBuilder.WriteString(input.Folder)
	state.Id = types.StringValue(idBuilder.String())
	var var0 []trustedCertificateAuthoritiesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]trustedCertificateAuthoritiesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 trustedCertificateAuthoritiesListDsModelConfig
			var2.CommonName = types.StringValue(var1.CommonName)
			var2.ExpiryEpoch = types.StringValue(var1.ExpiryEpoch)
			var2.Filename = types.StringValue(var1.Filename)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Issuer = types.StringValue(var1.Issuer)
			var2.Name = types.StringValue(var1.Name)
			var2.NotValidAfter = types.StringValue(var1.NotValidAfter)
			var2.NotValidBefore = types.StringValue(var1.NotValidBefore)
			var2.SerialNumber = types.StringValue(var1.SerialNumber)
			var2.Subject = types.StringValue(var1.Subject)
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
