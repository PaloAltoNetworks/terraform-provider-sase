package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	kmfIrpR "github.com/paloaltonetworks/sase-go/netsec/service/v1/certificates"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Data source (listing).
var (
	_ datasource.DataSource              = &certificatesGetListDataSource{}
	_ datasource.DataSourceWithConfigure = &certificatesGetListDataSource{}
)

func NewCertificatesGetListDataSource() datasource.DataSource {
	return &certificatesGetListDataSource{}
}

type certificatesGetListDataSource struct {
	client *sase.Client
}

type certificatesGetListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []certificatesGetListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type certificatesGetListDsModelConfig struct {
	Algorithm      types.String `tfsdk:"algorithm"`
	Ca             types.Bool   `tfsdk:"ca"`
	CommonName     types.String `tfsdk:"common_name"`
	CommonNameInt  types.String `tfsdk:"common_name_int"`
	ExpiryEpoch    types.String `tfsdk:"expiry_epoch"`
	ObjectId       types.String `tfsdk:"object_id"`
	Issuer         types.String `tfsdk:"issuer"`
	IssuerHash     types.String `tfsdk:"issuer_hash"`
	NotValidAfter  types.String `tfsdk:"not_valid_after"`
	NotValidBefore types.String `tfsdk:"not_valid_before"`
	PublicKey      types.String `tfsdk:"public_key"`
	Subject        types.String `tfsdk:"subject"`
	SubjectHash    types.String `tfsdk:"subject_hash"`
	SubjectInt     types.String `tfsdk:"subject_int"`
}

// Metadata returns the data source type name.
func (d *certificatesGetListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificates_get_list"
}

// Schema defines the schema for this listing data source.
func (d *certificatesGetListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"name": dsschema.StringAttribute{
				Description:         "The name of the entry",
				MarkdownDescription: "The name of the entry",
				Optional:            true,
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
						"algorithm": dsschema.StringAttribute{
							Description:         "The `algorithm` parameter.",
							MarkdownDescription: "The `algorithm` parameter.",
							Computed:            true,
						},
						"ca": dsschema.BoolAttribute{
							Description:         "The `ca` parameter.",
							MarkdownDescription: "The `ca` parameter.",
							Computed:            true,
						},
						"common_name": dsschema.StringAttribute{
							Description:         "The `common_name` parameter.",
							MarkdownDescription: "The `common_name` parameter.",
							Computed:            true,
						},
						"common_name_int": dsschema.StringAttribute{
							Description:         "The `common_name_int` parameter.",
							MarkdownDescription: "The `common_name_int` parameter.",
							Computed:            true,
						},
						"expiry_epoch": dsschema.StringAttribute{
							Description:         "The `expiry_epoch` parameter.",
							MarkdownDescription: "The `expiry_epoch` parameter.",
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
						"issuer_hash": dsschema.StringAttribute{
							Description:         "The `issuer_hash` parameter.",
							MarkdownDescription: "The `issuer_hash` parameter.",
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
						"public_key": dsschema.StringAttribute{
							Description:         "The `public_key` parameter.",
							MarkdownDescription: "The `public_key` parameter.",
							Computed:            true,
						},
						"subject": dsschema.StringAttribute{
							Description:         "The `subject` parameter.",
							MarkdownDescription: "The `subject` parameter.",
							Computed:            true,
						},
						"subject_hash": dsschema.StringAttribute{
							Description:         "The `subject_hash` parameter.",
							MarkdownDescription: "The `subject_hash` parameter.",
							Computed:            true,
						},
						"subject_int": dsschema.StringAttribute{
							Description:         "The `subject_int` parameter.",
							MarkdownDescription: "The `subject_int` parameter.",
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
func (d *certificatesGetListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *certificatesGetListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state certificatesGetListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_certificates_get_list",
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
	svc := kmfIrpR.NewClient(d.client)
	input := kmfIrpR.ListInput{
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
	var var0 []certificatesGetListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]certificatesGetListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 certificatesGetListDsModelConfig
			var2.Algorithm = types.StringValue(var1.Algorithm)
			var2.Ca = types.BoolValue(var1.Ca)
			var2.CommonName = types.StringValue(var1.CommonName)
			var2.CommonNameInt = types.StringValue(var1.CommonNameInt)
			var2.ExpiryEpoch = types.StringValue(var1.ExpiryEpoch)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Issuer = types.StringValue(var1.Issuer)
			var2.IssuerHash = types.StringValue(var1.IssuerHash)
			var2.NotValidAfter = types.StringValue(var1.NotValidAfter)
			var2.NotValidBefore = types.StringValue(var1.NotValidBefore)
			var2.PublicKey = types.StringValue(var1.PublicKey)
			var2.Subject = types.StringValue(var1.Subject)
			var2.SubjectHash = types.StringValue(var1.SubjectHash)
			var2.SubjectInt = types.StringValue(var1.SubjectInt)
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
