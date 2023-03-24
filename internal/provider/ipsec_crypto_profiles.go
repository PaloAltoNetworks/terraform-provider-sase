package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	bajUiwB "github.com/paloaltonetworks/sase-go/netsec/schema/ipsec/crypto/profiles"
	nufThga "github.com/paloaltonetworks/sase-go/netsec/service/v1/ipseccryptoprofiles"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Data source (listing).
var (
	_ datasource.DataSource              = &ipsecCryptoProfilesListDataSource{}
	_ datasource.DataSourceWithConfigure = &ipsecCryptoProfilesListDataSource{}
)

func NewIpsecCryptoProfilesListDataSource() datasource.DataSource {
	return &ipsecCryptoProfilesListDataSource{}
}

type ipsecCryptoProfilesListDataSource struct {
	client *sase.Client
}

type ipsecCryptoProfilesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []ipsecCryptoProfilesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type ipsecCryptoProfilesListDsModelConfig struct {
	Ah       *ipsecCryptoProfilesListDsModelAhObject       `tfsdk:"ah"`
	DhGroup  types.String                                  `tfsdk:"dh_group"`
	Esp      *ipsecCryptoProfilesListDsModelEspObject      `tfsdk:"esp"`
	ObjectId types.String                                  `tfsdk:"object_id"`
	Lifesize *ipsecCryptoProfilesListDsModelLifesizeObject `tfsdk:"lifesize"`
	Lifetime ipsecCryptoProfilesListDsModelLifetimeObject  `tfsdk:"lifetime"`
	Name     types.String                                  `tfsdk:"name"`
}

type ipsecCryptoProfilesListDsModelAhObject struct {
	Authentication []types.String `tfsdk:"authentication"`
}

type ipsecCryptoProfilesListDsModelEspObject struct {
	Authentication []types.String `tfsdk:"authentication"`
	Encryption     []types.String `tfsdk:"encryption"`
}

type ipsecCryptoProfilesListDsModelLifesizeObject struct {
	Gb types.Int64 `tfsdk:"gb"`
	Kb types.Int64 `tfsdk:"kb"`
	Mb types.Int64 `tfsdk:"mb"`
	Tb types.Int64 `tfsdk:"tb"`
}

type ipsecCryptoProfilesListDsModelLifetimeObject struct {
	Days    types.Int64 `tfsdk:"days"`
	Hours   types.Int64 `tfsdk:"hours"`
	Minutes types.Int64 `tfsdk:"minutes"`
	Seconds types.Int64 `tfsdk:"seconds"`
}

// Metadata returns the data source type name.
func (d *ipsecCryptoProfilesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipsec_crypto_profiles_list"
}

// Schema defines the schema for this listing data source.
func (d *ipsecCryptoProfilesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"ah": dsschema.SingleNestedAttribute{
							Description:         "The `ah` parameter.",
							MarkdownDescription: "The `ah` parameter.",
							Computed:            true,
							Attributes: map[string]dsschema.Attribute{
								"authentication": dsschema.ListAttribute{
									Description:         "The `authentication` parameter.",
									MarkdownDescription: "The `authentication` parameter.",
									Computed:            true,
									ElementType:         types.StringType,
								},
							},
						},
						"dh_group": dsschema.StringAttribute{
							Description:         "The `dh_group` parameter.",
							MarkdownDescription: "The `dh_group` parameter.",
							Computed:            true,
						},
						"esp": dsschema.SingleNestedAttribute{
							Description:         "The `esp` parameter.",
							MarkdownDescription: "The `esp` parameter.",
							Computed:            true,
							Attributes: map[string]dsschema.Attribute{
								"authentication": dsschema.ListAttribute{
									Description:         "The `authentication` parameter.",
									MarkdownDescription: "The `authentication` parameter.",
									Computed:            true,
									ElementType:         types.StringType,
								},
								"encryption": dsschema.ListAttribute{
									Description:         "The `encryption` parameter.",
									MarkdownDescription: "The `encryption` parameter.",
									Computed:            true,
									ElementType:         types.StringType,
								},
							},
						},
						"object_id": dsschema.StringAttribute{
							Description:         "The `object_id` parameter.",
							MarkdownDescription: "The `object_id` parameter.",
							Computed:            true,
						},
						"lifesize": dsschema.SingleNestedAttribute{
							Description:         "The `lifesize` parameter.",
							MarkdownDescription: "The `lifesize` parameter.",
							Computed:            true,
							Attributes: map[string]dsschema.Attribute{
								"gb": dsschema.Int64Attribute{
									Description:         "The `gb` parameter.",
									MarkdownDescription: "The `gb` parameter.",
									Computed:            true,
								},
								"kb": dsschema.Int64Attribute{
									Description:         "The `kb` parameter.",
									MarkdownDescription: "The `kb` parameter.",
									Computed:            true,
								},
								"mb": dsschema.Int64Attribute{
									Description:         "The `mb` parameter.",
									MarkdownDescription: "The `mb` parameter.",
									Computed:            true,
								},
								"tb": dsschema.Int64Attribute{
									Description:         "The `tb` parameter.",
									MarkdownDescription: "The `tb` parameter.",
									Computed:            true,
								},
							},
						},
						"lifetime": dsschema.SingleNestedAttribute{
							Description:         "The `lifetime` parameter.",
							MarkdownDescription: "The `lifetime` parameter.",
							Computed:            true,
							Attributes: map[string]dsschema.Attribute{
								"days": dsschema.Int64Attribute{
									Description:         "The `days` parameter.",
									MarkdownDescription: "The `days` parameter.",
									Computed:            true,
								},
								"hours": dsschema.Int64Attribute{
									Description:         "The `hours` parameter.",
									MarkdownDescription: "The `hours` parameter.",
									Computed:            true,
								},
								"minutes": dsschema.Int64Attribute{
									Description:         "The `minutes` parameter.",
									MarkdownDescription: "The `minutes` parameter.",
									Computed:            true,
								},
								"seconds": dsschema.Int64Attribute{
									Description:         "The `seconds` parameter.",
									MarkdownDescription: "The `seconds` parameter.",
									Computed:            true,
								},
							},
						},
						"name": dsschema.StringAttribute{
							Description:         "The `name` parameter.",
							MarkdownDescription: "The `name` parameter.",
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
func (d *ipsecCryptoProfilesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *ipsecCryptoProfilesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ipsecCryptoProfilesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_ipsec_crypto_profiles_list",
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
	svc := nufThga.NewClient(d.client)
	input := nufThga.ListInput{
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
	var var0 []ipsecCryptoProfilesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]ipsecCryptoProfilesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 ipsecCryptoProfilesListDsModelConfig
			var var3 *ipsecCryptoProfilesListDsModelAhObject
			if var1.Ah != nil {
				var3 = &ipsecCryptoProfilesListDsModelAhObject{}
				var3.Authentication = EncodeStringSlice(var1.Ah.Authentication)
			}
			var var4 *ipsecCryptoProfilesListDsModelEspObject
			if var1.Esp != nil {
				var4 = &ipsecCryptoProfilesListDsModelEspObject{}
				var4.Authentication = EncodeStringSlice(var1.Esp.Authentication)
				var4.Encryption = EncodeStringSlice(var1.Esp.Encryption)
			}
			var var5 *ipsecCryptoProfilesListDsModelLifesizeObject
			if var1.Lifesize != nil {
				var5 = &ipsecCryptoProfilesListDsModelLifesizeObject{}
				var5.Gb = types.Int64Value(var1.Lifesize.Gb)
				var5.Kb = types.Int64Value(var1.Lifesize.Kb)
				var5.Mb = types.Int64Value(var1.Lifesize.Mb)
				var5.Tb = types.Int64Value(var1.Lifesize.Tb)
			}
			var var6 ipsecCryptoProfilesListDsModelLifetimeObject
			var6.Days = types.Int64Value(var1.Lifetime.Days)
			var6.Hours = types.Int64Value(var1.Lifetime.Hours)
			var6.Minutes = types.Int64Value(var1.Lifetime.Minutes)
			var6.Seconds = types.Int64Value(var1.Lifetime.Seconds)
			var2.Ah = var3
			var2.DhGroup = types.StringValue(var1.DhGroup)
			var2.Esp = var4
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Lifesize = var5
			var2.Lifetime = var6
			var2.Name = types.StringValue(var1.Name)
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

// Data source.
var (
	_ datasource.DataSource              = &ipsecCryptoProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &ipsecCryptoProfilesDataSource{}
)

func NewIpsecCryptoProfilesDataSource() datasource.DataSource {
	return &ipsecCryptoProfilesDataSource{}
}

type ipsecCryptoProfilesDataSource struct {
	client *sase.Client
}

type ipsecCryptoProfilesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`
	Folder   types.String `tfsdk:"folder"`

	// Output.
	// Ref: #/components/schemas/ipsec-crypto-profiles
	Ah      *ipsecCryptoProfilesDsModelAhObject  `tfsdk:"ah"`
	DhGroup types.String                         `tfsdk:"dh_group"`
	Esp     *ipsecCryptoProfilesDsModelEspObject `tfsdk:"esp"`
	// input omit: ObjectId
	Lifesize *ipsecCryptoProfilesDsModelLifesizeObject `tfsdk:"lifesize"`
	Lifetime ipsecCryptoProfilesDsModelLifetimeObject  `tfsdk:"lifetime"`
	Name     types.String                              `tfsdk:"name"`
}

type ipsecCryptoProfilesDsModelAhObject struct {
	Authentication []types.String `tfsdk:"authentication"`
}

type ipsecCryptoProfilesDsModelEspObject struct {
	Authentication []types.String `tfsdk:"authentication"`
	Encryption     []types.String `tfsdk:"encryption"`
}

type ipsecCryptoProfilesDsModelLifesizeObject struct {
	Gb types.Int64 `tfsdk:"gb"`
	Kb types.Int64 `tfsdk:"kb"`
	Mb types.Int64 `tfsdk:"mb"`
	Tb types.Int64 `tfsdk:"tb"`
}

type ipsecCryptoProfilesDsModelLifetimeObject struct {
	Days    types.Int64 `tfsdk:"days"`
	Hours   types.Int64 `tfsdk:"hours"`
	Minutes types.Int64 `tfsdk:"minutes"`
	Seconds types.Int64 `tfsdk:"seconds"`
}

// Metadata returns the data source type name.
func (d *ipsecCryptoProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipsec_crypto_profiles"
}

// Schema defines the schema for this listing data source.
func (d *ipsecCryptoProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Retrieves config for a specific item.",

		Attributes: map[string]dsschema.Attribute{
			"id": dsschema.StringAttribute{
				Description:         "The object ID.",
				MarkdownDescription: "The object ID.",
				Computed:            true,
			},

			// Input.
			"object_id": dsschema.StringAttribute{
				Description:         "The uuid of the resource",
				MarkdownDescription: "The uuid of the resource",
				Required:            true,
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
			"ah": dsschema.SingleNestedAttribute{
				Description:         "The `ah` parameter.",
				MarkdownDescription: "The `ah` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"authentication": dsschema.ListAttribute{
						Description:         "The `authentication` parameter.",
						MarkdownDescription: "The `authentication` parameter.",
						Computed:            true,
						ElementType:         types.StringType,
					},
				},
			},
			"dh_group": dsschema.StringAttribute{
				Description:         "The `dh_group` parameter.",
				MarkdownDescription: "The `dh_group` parameter.",
				Computed:            true,
			},
			"esp": dsschema.SingleNestedAttribute{
				Description:         "The `esp` parameter.",
				MarkdownDescription: "The `esp` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"authentication": dsschema.ListAttribute{
						Description:         "The `authentication` parameter.",
						MarkdownDescription: "The `authentication` parameter.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"encryption": dsschema.ListAttribute{
						Description:         "The `encryption` parameter.",
						MarkdownDescription: "The `encryption` parameter.",
						Computed:            true,
						ElementType:         types.StringType,
					},
				},
			},
			"lifesize": dsschema.SingleNestedAttribute{
				Description:         "The `lifesize` parameter.",
				MarkdownDescription: "The `lifesize` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"gb": dsschema.Int64Attribute{
						Description:         "The `gb` parameter.",
						MarkdownDescription: "The `gb` parameter.",
						Computed:            true,
					},
					"kb": dsschema.Int64Attribute{
						Description:         "The `kb` parameter.",
						MarkdownDescription: "The `kb` parameter.",
						Computed:            true,
					},
					"mb": dsschema.Int64Attribute{
						Description:         "The `mb` parameter.",
						MarkdownDescription: "The `mb` parameter.",
						Computed:            true,
					},
					"tb": dsschema.Int64Attribute{
						Description:         "The `tb` parameter.",
						MarkdownDescription: "The `tb` parameter.",
						Computed:            true,
					},
				},
			},
			"lifetime": dsschema.SingleNestedAttribute{
				Description:         "The `lifetime` parameter.",
				MarkdownDescription: "The `lifetime` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"days": dsschema.Int64Attribute{
						Description:         "The `days` parameter.",
						MarkdownDescription: "The `days` parameter.",
						Computed:            true,
					},
					"hours": dsschema.Int64Attribute{
						Description:         "The `hours` parameter.",
						MarkdownDescription: "The `hours` parameter.",
						Computed:            true,
					},
					"minutes": dsschema.Int64Attribute{
						Description:         "The `minutes` parameter.",
						MarkdownDescription: "The `minutes` parameter.",
						Computed:            true,
					},
					"seconds": dsschema.Int64Attribute{
						Description:         "The `seconds` parameter.",
						MarkdownDescription: "The `seconds` parameter.",
						Computed:            true,
					},
				},
			},
			"name": dsschema.StringAttribute{
				Description:         "The `name` parameter.",
				MarkdownDescription: "The `name` parameter.",
				Computed:            true,
			},
		},
	}
}

// Configure prepares the struct.
func (d *ipsecCryptoProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *ipsecCryptoProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ipsecCryptoProfilesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_ipsec_crypto_profiles",
		"object_id":                   state.ObjectId.ValueString(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := nufThga.NewClient(d.client)
	input := nufThga.ReadInput{
		ObjectId: state.ObjectId.ValueString(),
		Folder:   state.Folder.ValueString(),
	}

	// Perform the operation.
	ans, err := svc.Read(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error getting singleton", err.Error())
		return
	}

	// Store the answer to state.
	var idBuilder strings.Builder
	idBuilder.WriteString(input.ObjectId)
	idBuilder.WriteString(IdSeparator)
	idBuilder.WriteString(input.Folder)
	state.Id = types.StringValue(idBuilder.String())
	var var0 *ipsecCryptoProfilesDsModelAhObject
	if ans.Ah != nil {
		var0 = &ipsecCryptoProfilesDsModelAhObject{}
		var0.Authentication = EncodeStringSlice(ans.Ah.Authentication)
	}
	var var1 *ipsecCryptoProfilesDsModelEspObject
	if ans.Esp != nil {
		var1 = &ipsecCryptoProfilesDsModelEspObject{}
		var1.Authentication = EncodeStringSlice(ans.Esp.Authentication)
		var1.Encryption = EncodeStringSlice(ans.Esp.Encryption)
	}
	var var2 *ipsecCryptoProfilesDsModelLifesizeObject
	if ans.Lifesize != nil {
		var2 = &ipsecCryptoProfilesDsModelLifesizeObject{}
		var2.Gb = types.Int64Value(ans.Lifesize.Gb)
		var2.Kb = types.Int64Value(ans.Lifesize.Kb)
		var2.Mb = types.Int64Value(ans.Lifesize.Mb)
		var2.Tb = types.Int64Value(ans.Lifesize.Tb)
	}
	var var3 ipsecCryptoProfilesDsModelLifetimeObject
	var3.Days = types.Int64Value(ans.Lifetime.Days)
	var3.Hours = types.Int64Value(ans.Lifetime.Hours)
	var3.Minutes = types.Int64Value(ans.Lifetime.Minutes)
	var3.Seconds = types.Int64Value(ans.Lifetime.Seconds)
	state.Ah = var0
	state.DhGroup = types.StringValue(ans.DhGroup)
	state.Esp = var1
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Lifesize = var2
	state.Lifetime = var3
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &ipsecCryptoProfilesResource{}
	_ resource.ResourceWithConfigure   = &ipsecCryptoProfilesResource{}
	_ resource.ResourceWithImportState = &ipsecCryptoProfilesResource{}
)

func NewIpsecCryptoProfilesResource() resource.Resource {
	return &ipsecCryptoProfilesResource{}
}

type ipsecCryptoProfilesResource struct {
	client *sase.Client
}

type ipsecCryptoProfilesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/ipsec-crypto-profiles
	Ah       *ipsecCryptoProfilesRsModelAhObject       `tfsdk:"ah"`
	DhGroup  types.String                              `tfsdk:"dh_group"`
	Esp      *ipsecCryptoProfilesRsModelEspObject      `tfsdk:"esp"`
	ObjectId types.String                              `tfsdk:"object_id"`
	Lifesize *ipsecCryptoProfilesRsModelLifesizeObject `tfsdk:"lifesize"`
	Lifetime ipsecCryptoProfilesRsModelLifetimeObject  `tfsdk:"lifetime"`
	Name     types.String                              `tfsdk:"name"`
}

type ipsecCryptoProfilesRsModelAhObject struct {
	Authentication []types.String `tfsdk:"authentication"`
}

type ipsecCryptoProfilesRsModelEspObject struct {
	Authentication []types.String `tfsdk:"authentication"`
	Encryption     []types.String `tfsdk:"encryption"`
}

type ipsecCryptoProfilesRsModelLifesizeObject struct {
	Gb types.Int64 `tfsdk:"gb"`
	Kb types.Int64 `tfsdk:"kb"`
	Mb types.Int64 `tfsdk:"mb"`
	Tb types.Int64 `tfsdk:"tb"`
}

type ipsecCryptoProfilesRsModelLifetimeObject struct {
	Days    types.Int64 `tfsdk:"days"`
	Hours   types.Int64 `tfsdk:"hours"`
	Minutes types.Int64 `tfsdk:"minutes"`
	Seconds types.Int64 `tfsdk:"seconds"`
}

// Metadata returns the data source type name.
func (r *ipsecCryptoProfilesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipsec_crypto_profiles"
}

// Schema defines the schema for this listing data source.
func (r *ipsecCryptoProfilesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = rsschema.Schema{
		Description: "Retrieves config for a specific item.",

		Attributes: map[string]rsschema.Attribute{
			"id": rsschema.StringAttribute{
				Description:         "The object ID.",
				MarkdownDescription: "The object ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			// Input.
			"folder": rsschema.StringAttribute{
				Description:         "The folder of the entry",
				MarkdownDescription: "The folder of the entry",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Service Connections", "Mobile Users Container", "Mobile Users Explicit Proxy"),
				},
			},

			"ah": rsschema.SingleNestedAttribute{
				Description:         "The `ah` parameter.",
				MarkdownDescription: "The `ah` parameter.",
				Optional:            true,
				Attributes: map[string]rsschema.Attribute{
					"authentication": rsschema.ListAttribute{
						Description:         "The `authentication` parameter.",
						MarkdownDescription: "The `authentication` parameter.",
						Required:            true,
						ElementType:         types.StringType,
					},
				},
			},
			"dh_group": rsschema.StringAttribute{
				Description:         "The `dh_group` parameter.",
				MarkdownDescription: "The `dh_group` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					DefaultString("group2"),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("no-pfs", "group1", "group2", "group5", "group14", "group19", "group20"),
				},
			},
			"esp": rsschema.SingleNestedAttribute{
				Description:         "The `esp` parameter.",
				MarkdownDescription: "The `esp` parameter.",
				Optional:            true,
				Attributes: map[string]rsschema.Attribute{
					"authentication": rsschema.ListAttribute{
						Description:         "The `authentication` parameter.",
						MarkdownDescription: "The `authentication` parameter.",
						Required:            true,
						ElementType:         types.StringType,
					},
					"encryption": rsschema.ListAttribute{
						Description:         "The `encryption` parameter.",
						MarkdownDescription: "The `encryption` parameter.",
						Required:            true,
						ElementType:         types.StringType,
					},
				},
			},
			"object_id": rsschema.StringAttribute{
				Description:         "The `object_id` parameter.",
				MarkdownDescription: "The `object_id` parameter.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"lifesize": rsschema.SingleNestedAttribute{
				Description:         "The `lifesize` parameter.",
				MarkdownDescription: "The `lifesize` parameter.",
				Optional:            true,
				Attributes: map[string]rsschema.Attribute{
					"gb": rsschema.Int64Attribute{
						Description:         "The `gb` parameter.",
						MarkdownDescription: "The `gb` parameter.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Int64{
							DefaultInt64(0),
						},
						Validators: []validator.Int64{
							int64validator.Between(1, 65535),
						},
					},
					"kb": rsschema.Int64Attribute{
						Description:         "The `kb` parameter.",
						MarkdownDescription: "The `kb` parameter.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Int64{
							DefaultInt64(0),
						},
						Validators: []validator.Int64{
							int64validator.Between(1, 65535),
						},
					},
					"mb": rsschema.Int64Attribute{
						Description:         "The `mb` parameter.",
						MarkdownDescription: "The `mb` parameter.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Int64{
							DefaultInt64(0),
						},
						Validators: []validator.Int64{
							int64validator.Between(1, 65535),
						},
					},
					"tb": rsschema.Int64Attribute{
						Description:         "The `tb` parameter.",
						MarkdownDescription: "The `tb` parameter.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Int64{
							DefaultInt64(0),
						},
						Validators: []validator.Int64{
							int64validator.Between(1, 65535),
						},
					},
				},
			},
			"lifetime": rsschema.SingleNestedAttribute{
				Description:         "The `lifetime` parameter.",
				MarkdownDescription: "The `lifetime` parameter.",
				Required:            true,
				Attributes: map[string]rsschema.Attribute{
					"days": rsschema.Int64Attribute{
						Description:         "The `days` parameter.",
						MarkdownDescription: "The `days` parameter.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Int64{
							DefaultInt64(0),
						},
						Validators: []validator.Int64{
							int64validator.Between(1, 365),
						},
					},
					"hours": rsschema.Int64Attribute{
						Description:         "The `hours` parameter.",
						MarkdownDescription: "The `hours` parameter.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Int64{
							DefaultInt64(0),
						},
						Validators: []validator.Int64{
							int64validator.Between(1, 65535),
						},
					},
					"minutes": rsschema.Int64Attribute{
						Description:         "The `minutes` parameter.",
						MarkdownDescription: "The `minutes` parameter.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Int64{
							DefaultInt64(0),
						},
						Validators: []validator.Int64{
							int64validator.Between(3, 65535),
						},
					},
					"seconds": rsschema.Int64Attribute{
						Description:         "The `seconds` parameter.",
						MarkdownDescription: "The `seconds` parameter.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Int64{
							DefaultInt64(0),
						},
						Validators: []validator.Int64{
							int64validator.Between(180, 65535),
						},
					},
				},
			},
			"name": rsschema.StringAttribute{
				Description:         "The `name` parameter.",
				MarkdownDescription: "The `name` parameter.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(31),
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *ipsecCryptoProfilesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *ipsecCryptoProfilesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state ipsecCryptoProfilesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_ipsec_crypto_profiles",
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := nufThga.NewClient(r.client)
	input := nufThga.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 bajUiwB.Config
	var var1 *bajUiwB.AhObject
	if state.Ah != nil {
		var1 = &bajUiwB.AhObject{}
		var1.Authentication = DecodeStringSlice(state.Ah.Authentication)
	}
	var0.Ah = var1
	var0.DhGroup = state.DhGroup.ValueString()
	var var2 *bajUiwB.EspObject
	if state.Esp != nil {
		var2 = &bajUiwB.EspObject{}
		var2.Authentication = DecodeStringSlice(state.Esp.Authentication)
		var2.Encryption = DecodeStringSlice(state.Esp.Encryption)
	}
	var0.Esp = var2
	var var3 *bajUiwB.LifesizeObject
	if state.Lifesize != nil {
		var3 = &bajUiwB.LifesizeObject{}
		var3.Gb = state.Lifesize.Gb.ValueInt64()
		var3.Kb = state.Lifesize.Kb.ValueInt64()
		var3.Mb = state.Lifesize.Mb.ValueInt64()
		var3.Tb = state.Lifesize.Tb.ValueInt64()
	}
	var0.Lifesize = var3
	var var4 bajUiwB.LifetimeObject
	var4.Days = state.Lifetime.Days.ValueInt64()
	var4.Hours = state.Lifetime.Hours.ValueInt64()
	var4.Minutes = state.Lifetime.Minutes.ValueInt64()
	var4.Seconds = state.Lifetime.Seconds.ValueInt64()
	var0.Lifetime = var4
	var0.Name = state.Name.ValueString()
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Create(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in create", err.Error())
		return
	}

	// Store the answer to state.
	var idBuilder strings.Builder
	idBuilder.WriteString(input.Folder)
	idBuilder.WriteString(IdSeparator)
	idBuilder.WriteString(ans.ObjectId)
	state.Id = types.StringValue(idBuilder.String())
	var var5 *ipsecCryptoProfilesRsModelAhObject
	if ans.Ah != nil {
		var5 = &ipsecCryptoProfilesRsModelAhObject{}
		var5.Authentication = EncodeStringSlice(ans.Ah.Authentication)
	}
	var var6 *ipsecCryptoProfilesRsModelEspObject
	if ans.Esp != nil {
		var6 = &ipsecCryptoProfilesRsModelEspObject{}
		var6.Authentication = EncodeStringSlice(ans.Esp.Authentication)
		var6.Encryption = EncodeStringSlice(ans.Esp.Encryption)
	}
	var var7 *ipsecCryptoProfilesRsModelLifesizeObject
	if ans.Lifesize != nil {
		var7 = &ipsecCryptoProfilesRsModelLifesizeObject{}
		var7.Gb = types.Int64Value(ans.Lifesize.Gb)
		var7.Kb = types.Int64Value(ans.Lifesize.Kb)
		var7.Mb = types.Int64Value(ans.Lifesize.Mb)
		var7.Tb = types.Int64Value(ans.Lifesize.Tb)
	}
	var var8 ipsecCryptoProfilesRsModelLifetimeObject
	var8.Days = types.Int64Value(ans.Lifetime.Days)
	var8.Hours = types.Int64Value(ans.Lifetime.Hours)
	var8.Minutes = types.Int64Value(ans.Lifetime.Minutes)
	var8.Seconds = types.Int64Value(ans.Lifetime.Seconds)
	state.Ah = var5
	state.DhGroup = types.StringValue(ans.DhGroup)
	state.Esp = var6
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Lifesize = var7
	state.Lifetime = var8
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *ipsecCryptoProfilesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var idType types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &idType)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := idType.ValueString()
	tokens := strings.Split(id, IdSeparator)
	if len(tokens) != 2 {
		resp.Diagnostics.AddError("Error in resource ID format", "Expected 2 tokens")
		return
	}

	var state ipsecCryptoProfilesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_ipsec_crypto_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := nufThga.NewClient(r.client)
	input := nufThga.ReadInput{
		ObjectId: tokens[1],
		Folder:   tokens[0],
	}

	// Perform the operation.
	ans, err := svc.Read(ctx, input)
	if err != nil {
		if IsObjectNotFound(err) {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError("Error reading config", err.Error())
		}
		return
	}

	// Store the answer to state.
	state.Folder = types.StringValue(tokens[0])
	state.Id = idType
	var var0 *ipsecCryptoProfilesRsModelAhObject
	if ans.Ah != nil {
		var0 = &ipsecCryptoProfilesRsModelAhObject{}
		var0.Authentication = EncodeStringSlice(ans.Ah.Authentication)
	}
	var var1 *ipsecCryptoProfilesRsModelEspObject
	if ans.Esp != nil {
		var1 = &ipsecCryptoProfilesRsModelEspObject{}
		var1.Authentication = EncodeStringSlice(ans.Esp.Authentication)
		var1.Encryption = EncodeStringSlice(ans.Esp.Encryption)
	}
	var var2 *ipsecCryptoProfilesRsModelLifesizeObject
	if ans.Lifesize != nil {
		var2 = &ipsecCryptoProfilesRsModelLifesizeObject{}
		var2.Gb = types.Int64Value(ans.Lifesize.Gb)
		var2.Kb = types.Int64Value(ans.Lifesize.Kb)
		var2.Mb = types.Int64Value(ans.Lifesize.Mb)
		var2.Tb = types.Int64Value(ans.Lifesize.Tb)
	}
	var var3 ipsecCryptoProfilesRsModelLifetimeObject
	var3.Days = types.Int64Value(ans.Lifetime.Days)
	var3.Hours = types.Int64Value(ans.Lifetime.Hours)
	var3.Minutes = types.Int64Value(ans.Lifetime.Minutes)
	var3.Seconds = types.Int64Value(ans.Lifetime.Seconds)
	state.Ah = var0
	state.DhGroup = types.StringValue(ans.DhGroup)
	state.Esp = var1
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Lifesize = var2
	state.Lifetime = var3
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *ipsecCryptoProfilesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ipsecCryptoProfilesRsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource update", map[string]any{
		"terraform_provider_function": "Update",
		"resource_name":               "sase_ipsec_crypto_profiles",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := nufThga.NewClient(r.client)
	input := nufThga.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 bajUiwB.Config
	var var1 *bajUiwB.AhObject
	if plan.Ah != nil {
		var1 = &bajUiwB.AhObject{}
		var1.Authentication = DecodeStringSlice(plan.Ah.Authentication)
	}
	var0.Ah = var1
	var0.DhGroup = plan.DhGroup.ValueString()
	var var2 *bajUiwB.EspObject
	if plan.Esp != nil {
		var2 = &bajUiwB.EspObject{}
		var2.Authentication = DecodeStringSlice(plan.Esp.Authentication)
		var2.Encryption = DecodeStringSlice(plan.Esp.Encryption)
	}
	var0.Esp = var2
	var var3 *bajUiwB.LifesizeObject
	if plan.Lifesize != nil {
		var3 = &bajUiwB.LifesizeObject{}
		var3.Gb = plan.Lifesize.Gb.ValueInt64()
		var3.Kb = plan.Lifesize.Kb.ValueInt64()
		var3.Mb = plan.Lifesize.Mb.ValueInt64()
		var3.Tb = plan.Lifesize.Tb.ValueInt64()
	}
	var0.Lifesize = var3
	var var4 bajUiwB.LifetimeObject
	var4.Days = plan.Lifetime.Days.ValueInt64()
	var4.Hours = plan.Lifetime.Hours.ValueInt64()
	var4.Minutes = plan.Lifetime.Minutes.ValueInt64()
	var4.Seconds = plan.Lifetime.Seconds.ValueInt64()
	var0.Lifetime = var4
	var0.Name = plan.Name.ValueString()
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var5 *ipsecCryptoProfilesRsModelAhObject
	if ans.Ah != nil {
		var5 = &ipsecCryptoProfilesRsModelAhObject{}
		var5.Authentication = EncodeStringSlice(ans.Ah.Authentication)
	}
	var var6 *ipsecCryptoProfilesRsModelEspObject
	if ans.Esp != nil {
		var6 = &ipsecCryptoProfilesRsModelEspObject{}
		var6.Authentication = EncodeStringSlice(ans.Esp.Authentication)
		var6.Encryption = EncodeStringSlice(ans.Esp.Encryption)
	}
	var var7 *ipsecCryptoProfilesRsModelLifesizeObject
	if ans.Lifesize != nil {
		var7 = &ipsecCryptoProfilesRsModelLifesizeObject{}
		var7.Gb = types.Int64Value(ans.Lifesize.Gb)
		var7.Kb = types.Int64Value(ans.Lifesize.Kb)
		var7.Mb = types.Int64Value(ans.Lifesize.Mb)
		var7.Tb = types.Int64Value(ans.Lifesize.Tb)
	}
	var var8 ipsecCryptoProfilesRsModelLifetimeObject
	var8.Days = types.Int64Value(ans.Lifetime.Days)
	var8.Hours = types.Int64Value(ans.Lifetime.Hours)
	var8.Minutes = types.Int64Value(ans.Lifetime.Minutes)
	var8.Seconds = types.Int64Value(ans.Lifetime.Seconds)
	state.Ah = var5
	state.DhGroup = types.StringValue(ans.DhGroup)
	state.Esp = var6
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Lifesize = var7
	state.Lifetime = var8
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *ipsecCryptoProfilesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var idType types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &idType)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := idType.ValueString()
	tokens := strings.Split(id, IdSeparator)
	if len(tokens) != 2 {
		resp.Diagnostics.AddError("Error in resource ID format", "Expected 2 tokens")
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource delete", map[string]any{
		"terraform_provider_function": "Delete",
		"resource_name":               "sase_ipsec_crypto_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	svc := nufThga.NewClient(r.client)
	input := nufThga.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *ipsecCryptoProfilesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
