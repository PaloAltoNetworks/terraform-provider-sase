package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	ktjCEnF "github.com/paloaltonetworks/sase-go/netsec/schema/objects/services"
	eumQbRC "github.com/paloaltonetworks/sase-go/netsec/service/v1/services"

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
	_ datasource.DataSource              = &objectsServicesListDataSource{}
	_ datasource.DataSourceWithConfigure = &objectsServicesListDataSource{}
)

func NewObjectsServicesListDataSource() datasource.DataSource {
	return &objectsServicesListDataSource{}
}

type objectsServicesListDataSource struct {
	client *sase.Client
}

type objectsServicesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []objectsServicesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type objectsServicesListDsModelConfig struct {
	Description types.String                             `tfsdk:"description"`
	ObjectId    types.String                             `tfsdk:"object_id"`
	Name        types.String                             `tfsdk:"name"`
	Protocol    objectsServicesListDsModelProtocolObject `tfsdk:"protocol"`
	Tag         []types.String                           `tfsdk:"tag"`
}

type objectsServicesListDsModelProtocolObject struct {
	Tcp *objectsServicesListDsModelTcpObject `tfsdk:"tcp"`
	Udp *objectsServicesListDsModelUdpObject `tfsdk:"udp"`
}

type objectsServicesListDsModelTcpObject struct {
	Override   *objectsServicesListDsModelOverrideObject `tfsdk:"override"`
	Port       types.String                              `tfsdk:"port"`
	SourcePort types.String                              `tfsdk:"source_port"`
}

type objectsServicesListDsModelOverrideObject struct {
	HalfcloseTimeout types.Int64 `tfsdk:"halfclose_timeout"`
	Timeout          types.Int64 `tfsdk:"timeout"`
	TimewaitTimeout  types.Int64 `tfsdk:"timewait_timeout"`
}

type objectsServicesListDsModelUdpObject struct {
	Override   *objectsServicesListDsModelOverrideObject1 `tfsdk:"override"`
	Port       types.String                               `tfsdk:"port"`
	SourcePort types.String                               `tfsdk:"source_port"`
}

type objectsServicesListDsModelOverrideObject1 struct {
	Timeout types.Int64 `tfsdk:"timeout"`
}

// Metadata returns the data source type name.
func (d *objectsServicesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_services_list"
}

// Schema defines the schema for this listing data source.
func (d *objectsServicesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"description": dsschema.StringAttribute{
							Description:         "The `description` parameter.",
							MarkdownDescription: "The `description` parameter.",
							Computed:            true,
						},
						"object_id": dsschema.StringAttribute{
							Description:         "The `object_id` parameter.",
							MarkdownDescription: "The `object_id` parameter.",
							Computed:            true,
						},
						"name": dsschema.StringAttribute{
							Description:         "The `name` parameter.",
							MarkdownDescription: "The `name` parameter.",
							Computed:            true,
						},
						"protocol": dsschema.SingleNestedAttribute{
							Description:         "The `protocol` parameter.",
							MarkdownDescription: "The `protocol` parameter.",
							Computed:            true,
							Attributes: map[string]dsschema.Attribute{
								"tcp": dsschema.SingleNestedAttribute{
									Description:         "The `tcp` parameter.",
									MarkdownDescription: "The `tcp` parameter.",
									Computed:            true,
									Attributes: map[string]dsschema.Attribute{
										"override": dsschema.SingleNestedAttribute{
											Description:         "The `override` parameter.",
											MarkdownDescription: "The `override` parameter.",
											Computed:            true,
											Attributes: map[string]dsschema.Attribute{
												"halfclose_timeout": dsschema.Int64Attribute{
													Description:         "The `halfclose_timeout` parameter.",
													MarkdownDescription: "The `halfclose_timeout` parameter.",
													Computed:            true,
												},
												"timeout": dsschema.Int64Attribute{
													Description:         "The `timeout` parameter.",
													MarkdownDescription: "The `timeout` parameter.",
													Computed:            true,
												},
												"timewait_timeout": dsschema.Int64Attribute{
													Description:         "The `timewait_timeout` parameter.",
													MarkdownDescription: "The `timewait_timeout` parameter.",
													Computed:            true,
												},
											},
										},
										"port": dsschema.StringAttribute{
											Description:         "The `port` parameter.",
											MarkdownDescription: "The `port` parameter.",
											Computed:            true,
										},
										"source_port": dsschema.StringAttribute{
											Description:         "The `source_port` parameter.",
											MarkdownDescription: "The `source_port` parameter.",
											Computed:            true,
										},
									},
								},
								"udp": dsschema.SingleNestedAttribute{
									Description:         "The `udp` parameter.",
									MarkdownDescription: "The `udp` parameter.",
									Computed:            true,
									Attributes: map[string]dsschema.Attribute{
										"override": dsschema.SingleNestedAttribute{
											Description:         "The `override` parameter.",
											MarkdownDescription: "The `override` parameter.",
											Computed:            true,
											Attributes: map[string]dsschema.Attribute{
												"timeout": dsschema.Int64Attribute{
													Description:         "The `timeout` parameter.",
													MarkdownDescription: "The `timeout` parameter.",
													Computed:            true,
												},
											},
										},
										"port": dsschema.StringAttribute{
											Description:         "The `port` parameter.",
											MarkdownDescription: "The `port` parameter.",
											Computed:            true,
										},
										"source_port": dsschema.StringAttribute{
											Description:         "The `source_port` parameter.",
											MarkdownDescription: "The `source_port` parameter.",
											Computed:            true,
										},
									},
								},
							},
						},
						"tag": dsschema.ListAttribute{
							Description:         "The `tag` parameter.",
							MarkdownDescription: "The `tag` parameter.",
							Computed:            true,
							ElementType:         types.StringType,
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
func (d *objectsServicesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *objectsServicesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state objectsServicesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_objects_services_list",
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
	svc := eumQbRC.NewClient(d.client)
	input := eumQbRC.ListInput{
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
	var var0 []objectsServicesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]objectsServicesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 objectsServicesListDsModelConfig
			var var3 objectsServicesListDsModelProtocolObject
			var var4 *objectsServicesListDsModelTcpObject
			if var1.Protocol.Tcp != nil {
				var4 = &objectsServicesListDsModelTcpObject{}
				var var5 *objectsServicesListDsModelOverrideObject
				if var1.Protocol.Tcp.Override != nil {
					var5 = &objectsServicesListDsModelOverrideObject{}
					var5.HalfcloseTimeout = types.Int64Value(var1.Protocol.Tcp.Override.HalfcloseTimeout)
					var5.Timeout = types.Int64Value(var1.Protocol.Tcp.Override.Timeout)
					var5.TimewaitTimeout = types.Int64Value(var1.Protocol.Tcp.Override.TimewaitTimeout)
				}
				var4.Override = var5
				var4.Port = types.StringValue(var1.Protocol.Tcp.Port)
				var4.SourcePort = types.StringValue(var1.Protocol.Tcp.SourcePort)
			}
			var var6 *objectsServicesListDsModelUdpObject
			if var1.Protocol.Udp != nil {
				var6 = &objectsServicesListDsModelUdpObject{}
				var var7 *objectsServicesListDsModelOverrideObject1
				if var1.Protocol.Udp.Override != nil {
					var7 = &objectsServicesListDsModelOverrideObject1{}
					var7.Timeout = types.Int64Value(var1.Protocol.Udp.Override.Timeout)
				}
				var6.Override = var7
				var6.Port = types.StringValue(var1.Protocol.Udp.Port)
				var6.SourcePort = types.StringValue(var1.Protocol.Udp.SourcePort)
			}
			var3.Tcp = var4
			var3.Udp = var6
			var2.Description = types.StringValue(var1.Description)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Name = types.StringValue(var1.Name)
			var2.Protocol = var3
			var2.Tag = EncodeStringSlice(var1.Tag)
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
	_ datasource.DataSource              = &objectsServicesDataSource{}
	_ datasource.DataSourceWithConfigure = &objectsServicesDataSource{}
)

func NewObjectsServicesDataSource() datasource.DataSource {
	return &objectsServicesDataSource{}
}

type objectsServicesDataSource struct {
	client *sase.Client
}

type objectsServicesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`
	Folder   types.String `tfsdk:"folder"`

	// Output.
	// Ref: #/components/schemas/objects-services
	Description types.String `tfsdk:"description"`
	// input omit: ObjectId
	Name     types.String                         `tfsdk:"name"`
	Protocol objectsServicesDsModelProtocolObject `tfsdk:"protocol"`
	Tag      []types.String                       `tfsdk:"tag"`
}

type objectsServicesDsModelProtocolObject struct {
	Tcp *objectsServicesDsModelTcpObject `tfsdk:"tcp"`
	Udp *objectsServicesDsModelUdpObject `tfsdk:"udp"`
}

type objectsServicesDsModelTcpObject struct {
	Override   *objectsServicesDsModelOverrideObject `tfsdk:"override"`
	Port       types.String                          `tfsdk:"port"`
	SourcePort types.String                          `tfsdk:"source_port"`
}

type objectsServicesDsModelOverrideObject struct {
	HalfcloseTimeout types.Int64 `tfsdk:"halfclose_timeout"`
	Timeout          types.Int64 `tfsdk:"timeout"`
	TimewaitTimeout  types.Int64 `tfsdk:"timewait_timeout"`
}

type objectsServicesDsModelUdpObject struct {
	Override   *objectsServicesDsModelOverrideObject1 `tfsdk:"override"`
	Port       types.String                           `tfsdk:"port"`
	SourcePort types.String                           `tfsdk:"source_port"`
}

type objectsServicesDsModelOverrideObject1 struct {
	Timeout types.Int64 `tfsdk:"timeout"`
}

// Metadata returns the data source type name.
func (d *objectsServicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_services"
}

// Schema defines the schema for this listing data source.
func (d *objectsServicesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
				Description:         "The uuid of the resource.",
				MarkdownDescription: "The uuid of the resource.",
				Required:            true,
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
			"description": dsschema.StringAttribute{
				Description:         "The `description` parameter.",
				MarkdownDescription: "The `description` parameter.",
				Computed:            true,
			},
			"name": dsschema.StringAttribute{
				Description:         "The `name` parameter.",
				MarkdownDescription: "The `name` parameter.",
				Computed:            true,
			},
			"protocol": dsschema.SingleNestedAttribute{
				Description:         "The `protocol` parameter.",
				MarkdownDescription: "The `protocol` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"tcp": dsschema.SingleNestedAttribute{
						Description:         "The `tcp` parameter.",
						MarkdownDescription: "The `tcp` parameter.",
						Computed:            true,
						Attributes: map[string]dsschema.Attribute{
							"override": dsschema.SingleNestedAttribute{
								Description:         "The `override` parameter.",
								MarkdownDescription: "The `override` parameter.",
								Computed:            true,
								Attributes: map[string]dsschema.Attribute{
									"halfclose_timeout": dsschema.Int64Attribute{
										Description:         "The `halfclose_timeout` parameter.",
										MarkdownDescription: "The `halfclose_timeout` parameter.",
										Computed:            true,
									},
									"timeout": dsschema.Int64Attribute{
										Description:         "The `timeout` parameter.",
										MarkdownDescription: "The `timeout` parameter.",
										Computed:            true,
									},
									"timewait_timeout": dsschema.Int64Attribute{
										Description:         "The `timewait_timeout` parameter.",
										MarkdownDescription: "The `timewait_timeout` parameter.",
										Computed:            true,
									},
								},
							},
							"port": dsschema.StringAttribute{
								Description:         "The `port` parameter.",
								MarkdownDescription: "The `port` parameter.",
								Computed:            true,
							},
							"source_port": dsschema.StringAttribute{
								Description:         "The `source_port` parameter.",
								MarkdownDescription: "The `source_port` parameter.",
								Computed:            true,
							},
						},
					},
					"udp": dsschema.SingleNestedAttribute{
						Description:         "The `udp` parameter.",
						MarkdownDescription: "The `udp` parameter.",
						Computed:            true,
						Attributes: map[string]dsschema.Attribute{
							"override": dsschema.SingleNestedAttribute{
								Description:         "The `override` parameter.",
								MarkdownDescription: "The `override` parameter.",
								Computed:            true,
								Attributes: map[string]dsschema.Attribute{
									"timeout": dsschema.Int64Attribute{
										Description:         "The `timeout` parameter.",
										MarkdownDescription: "The `timeout` parameter.",
										Computed:            true,
									},
								},
							},
							"port": dsschema.StringAttribute{
								Description:         "The `port` parameter.",
								MarkdownDescription: "The `port` parameter.",
								Computed:            true,
							},
							"source_port": dsschema.StringAttribute{
								Description:         "The `source_port` parameter.",
								MarkdownDescription: "The `source_port` parameter.",
								Computed:            true,
							},
						},
					},
				},
			},
			"tag": dsschema.ListAttribute{
				Description:         "The `tag` parameter.",
				MarkdownDescription: "The `tag` parameter.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

// Configure prepares the struct.
func (d *objectsServicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *objectsServicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state objectsServicesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_objects_services",
		"object_id":                   state.ObjectId.ValueString(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := eumQbRC.NewClient(d.client)
	input := eumQbRC.ReadInput{
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
	var var0 objectsServicesDsModelProtocolObject
	var var1 *objectsServicesDsModelTcpObject
	if ans.Protocol.Tcp != nil {
		var1 = &objectsServicesDsModelTcpObject{}
		var var2 *objectsServicesDsModelOverrideObject
		if ans.Protocol.Tcp.Override != nil {
			var2 = &objectsServicesDsModelOverrideObject{}
			var2.HalfcloseTimeout = types.Int64Value(ans.Protocol.Tcp.Override.HalfcloseTimeout)
			var2.Timeout = types.Int64Value(ans.Protocol.Tcp.Override.Timeout)
			var2.TimewaitTimeout = types.Int64Value(ans.Protocol.Tcp.Override.TimewaitTimeout)
		}
		var1.Override = var2
		var1.Port = types.StringValue(ans.Protocol.Tcp.Port)
		var1.SourcePort = types.StringValue(ans.Protocol.Tcp.SourcePort)
	}
	var var3 *objectsServicesDsModelUdpObject
	if ans.Protocol.Udp != nil {
		var3 = &objectsServicesDsModelUdpObject{}
		var var4 *objectsServicesDsModelOverrideObject1
		if ans.Protocol.Udp.Override != nil {
			var4 = &objectsServicesDsModelOverrideObject1{}
			var4.Timeout = types.Int64Value(ans.Protocol.Udp.Override.Timeout)
		}
		var3.Override = var4
		var3.Port = types.StringValue(ans.Protocol.Udp.Port)
		var3.SourcePort = types.StringValue(ans.Protocol.Udp.SourcePort)
	}
	var0.Tcp = var1
	var0.Udp = var3
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.Protocol = var0
	state.Tag = EncodeStringSlice(ans.Tag)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &objectsServicesResource{}
	_ resource.ResourceWithConfigure   = &objectsServicesResource{}
	_ resource.ResourceWithImportState = &objectsServicesResource{}
)

func NewObjectsServicesResource() resource.Resource {
	return &objectsServicesResource{}
}

type objectsServicesResource struct {
	client *sase.Client
}

type objectsServicesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/objects-services
	Description types.String                         `tfsdk:"description"`
	ObjectId    types.String                         `tfsdk:"object_id"`
	Name        types.String                         `tfsdk:"name"`
	Protocol    objectsServicesRsModelProtocolObject `tfsdk:"protocol"`
	Tag         []types.String                       `tfsdk:"tag"`
}

type objectsServicesRsModelProtocolObject struct {
	Tcp *objectsServicesRsModelTcpObject `tfsdk:"tcp"`
	Udp *objectsServicesRsModelUdpObject `tfsdk:"udp"`
}

type objectsServicesRsModelTcpObject struct {
	Override   *objectsServicesRsModelOverrideObject `tfsdk:"override"`
	Port       types.String                          `tfsdk:"port"`
	SourcePort types.String                          `tfsdk:"source_port"`
}

type objectsServicesRsModelOverrideObject struct {
	HalfcloseTimeout types.Int64 `tfsdk:"halfclose_timeout"`
	Timeout          types.Int64 `tfsdk:"timeout"`
	TimewaitTimeout  types.Int64 `tfsdk:"timewait_timeout"`
}

type objectsServicesRsModelUdpObject struct {
	Override   *objectsServicesRsModelOverrideObject1 `tfsdk:"override"`
	Port       types.String                           `tfsdk:"port"`
	SourcePort types.String                           `tfsdk:"source_port"`
}

type objectsServicesRsModelOverrideObject1 struct {
	Timeout types.Int64 `tfsdk:"timeout"`
}

// Metadata returns the data source type name.
func (r *objectsServicesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_services"
}

// Schema defines the schema for this listing data source.
func (r *objectsServicesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Description:         "The folder of the entry. Value must be one of: `\"Shared\"`, `\"Mobile Users\"`, `\"Remote Networks\"`, `\"Service Connections\"`, `\"Mobile Users Container\"`, `\"Mobile Users Explicit Proxy\"`.",
				MarkdownDescription: "The folder of the entry. Value must be one of: `\"Shared\"`, `\"Mobile Users\"`, `\"Remote Networks\"`, `\"Service Connections\"`, `\"Mobile Users Container\"`, `\"Mobile Users Explicit Proxy\"`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Service Connections", "Mobile Users Container", "Mobile Users Explicit Proxy"),
				},
			},

			"description": rsschema.StringAttribute{
				Description:         "The `description` parameter. String length must be between 0 and 1023.",
				MarkdownDescription: "The `description` parameter. String length must be between 0 and 1023.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 1023),
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
			"name": rsschema.StringAttribute{
				Description:         "The `name` parameter. String length must be at most 63.",
				MarkdownDescription: "The `name` parameter. String length must be at most 63.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(63),
				},
			},
			"protocol": rsschema.SingleNestedAttribute{
				Description:         "The `protocol` parameter.",
				MarkdownDescription: "The `protocol` parameter.",
				Required:            true,
				Attributes: map[string]rsschema.Attribute{
					"tcp": rsschema.SingleNestedAttribute{
						Description:         "The `tcp` parameter.",
						MarkdownDescription: "The `tcp` parameter.",
						Optional:            true,
						Attributes: map[string]rsschema.Attribute{
							"override": rsschema.SingleNestedAttribute{
								Description:         "The `override` parameter.",
								MarkdownDescription: "The `override` parameter.",
								Optional:            true,
								Attributes: map[string]rsschema.Attribute{
									"halfclose_timeout": rsschema.Int64Attribute{
										Description:         "The `halfclose_timeout` parameter. Default: `120`. Value must be between 1 and 604800.",
										MarkdownDescription: "The `halfclose_timeout` parameter. Default: `120`. Value must be between 1 and 604800.",
										Optional:            true,
										Computed:            true,
										PlanModifiers: []planmodifier.Int64{
											DefaultInt64(120),
										},
										Validators: []validator.Int64{
											int64validator.Between(1, 604800),
										},
									},
									"timeout": rsschema.Int64Attribute{
										Description:         "The `timeout` parameter. Default: `3600`. Value must be between 1 and 604800.",
										MarkdownDescription: "The `timeout` parameter. Default: `3600`. Value must be between 1 and 604800.",
										Optional:            true,
										Computed:            true,
										PlanModifiers: []planmodifier.Int64{
											DefaultInt64(3600),
										},
										Validators: []validator.Int64{
											int64validator.Between(1, 604800),
										},
									},
									"timewait_timeout": rsschema.Int64Attribute{
										Description:         "The `timewait_timeout` parameter. Default: `15`. Value must be between 1 and 600.",
										MarkdownDescription: "The `timewait_timeout` parameter. Default: `15`. Value must be between 1 and 600.",
										Optional:            true,
										Computed:            true,
										PlanModifiers: []planmodifier.Int64{
											DefaultInt64(15),
										},
										Validators: []validator.Int64{
											int64validator.Between(1, 600),
										},
									},
								},
							},
							"port": rsschema.StringAttribute{
								Description:         "The `port` parameter. String length must be between 1 and 1023.",
								MarkdownDescription: "The `port` parameter. String length must be between 1 and 1023.",
								Required:            true,
								Validators: []validator.String{
									stringvalidator.LengthBetween(1, 1023),
								},
							},
							"source_port": rsschema.StringAttribute{
								Description:         "The `source_port` parameter. String length must be between 1 and 1023.",
								MarkdownDescription: "The `source_port` parameter. String length must be between 1 and 1023.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
								Validators: []validator.String{
									stringvalidator.LengthBetween(1, 1023),
								},
							},
						},
					},
					"udp": rsschema.SingleNestedAttribute{
						Description:         "The `udp` parameter.",
						MarkdownDescription: "The `udp` parameter.",
						Optional:            true,
						Attributes: map[string]rsschema.Attribute{
							"override": rsschema.SingleNestedAttribute{
								Description:         "The `override` parameter.",
								MarkdownDescription: "The `override` parameter.",
								Optional:            true,
								Attributes: map[string]rsschema.Attribute{
									"timeout": rsschema.Int64Attribute{
										Description:         "The `timeout` parameter. Default: `30`. Value must be between 1 and 604800.",
										MarkdownDescription: "The `timeout` parameter. Default: `30`. Value must be between 1 and 604800.",
										Optional:            true,
										Computed:            true,
										PlanModifiers: []planmodifier.Int64{
											DefaultInt64(30),
										},
										Validators: []validator.Int64{
											int64validator.Between(1, 604800),
										},
									},
								},
							},
							"port": rsschema.StringAttribute{
								Description:         "The `port` parameter. String length must be between 1 and 1023.",
								MarkdownDescription: "The `port` parameter. String length must be between 1 and 1023.",
								Required:            true,
								Validators: []validator.String{
									stringvalidator.LengthBetween(1, 1023),
								},
							},
							"source_port": rsschema.StringAttribute{
								Description:         "The `source_port` parameter. String length must be between 1 and 1023.",
								MarkdownDescription: "The `source_port` parameter. String length must be between 1 and 1023.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
								Validators: []validator.String{
									stringvalidator.LengthBetween(1, 1023),
								},
							},
						},
					},
				},
			},
			"tag": rsschema.ListAttribute{
				Description:         "The `tag` parameter.",
				MarkdownDescription: "The `tag` parameter.",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

// Configure prepares the struct.
func (r *objectsServicesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *objectsServicesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state objectsServicesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_objects_services",
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := eumQbRC.NewClient(r.client)
	input := eumQbRC.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 ktjCEnF.Config
	var0.Description = state.Description.ValueString()
	var0.Name = state.Name.ValueString()
	var var1 ktjCEnF.ProtocolObject
	var var2 *ktjCEnF.TcpObject
	if state.Protocol.Tcp != nil {
		var2 = &ktjCEnF.TcpObject{}
		var var3 *ktjCEnF.OverrideObject
		if state.Protocol.Tcp.Override != nil {
			var3 = &ktjCEnF.OverrideObject{}
			var3.HalfcloseTimeout = state.Protocol.Tcp.Override.HalfcloseTimeout.ValueInt64()
			var3.Timeout = state.Protocol.Tcp.Override.Timeout.ValueInt64()
			var3.TimewaitTimeout = state.Protocol.Tcp.Override.TimewaitTimeout.ValueInt64()
		}
		var2.Override = var3
		var2.Port = state.Protocol.Tcp.Port.ValueString()
		var2.SourcePort = state.Protocol.Tcp.SourcePort.ValueString()
	}
	var1.Tcp = var2
	var var4 *ktjCEnF.UdpObject
	if state.Protocol.Udp != nil {
		var4 = &ktjCEnF.UdpObject{}
		var var5 *ktjCEnF.OverrideObject1
		if state.Protocol.Udp.Override != nil {
			var5 = &ktjCEnF.OverrideObject1{}
			var5.Timeout = state.Protocol.Udp.Override.Timeout.ValueInt64()
		}
		var4.Override = var5
		var4.Port = state.Protocol.Udp.Port.ValueString()
		var4.SourcePort = state.Protocol.Udp.SourcePort.ValueString()
	}
	var1.Udp = var4
	var0.Protocol = var1
	var0.Tag = DecodeStringSlice(state.Tag)
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
	var var6 objectsServicesRsModelProtocolObject
	var var7 *objectsServicesRsModelTcpObject
	if ans.Protocol.Tcp != nil {
		var7 = &objectsServicesRsModelTcpObject{}
		var var8 *objectsServicesRsModelOverrideObject
		if ans.Protocol.Tcp.Override != nil {
			var8 = &objectsServicesRsModelOverrideObject{}
			var8.HalfcloseTimeout = types.Int64Value(ans.Protocol.Tcp.Override.HalfcloseTimeout)
			var8.Timeout = types.Int64Value(ans.Protocol.Tcp.Override.Timeout)
			var8.TimewaitTimeout = types.Int64Value(ans.Protocol.Tcp.Override.TimewaitTimeout)
		}
		var7.Override = var8
		var7.Port = types.StringValue(ans.Protocol.Tcp.Port)
		var7.SourcePort = types.StringValue(ans.Protocol.Tcp.SourcePort)
	}
	var var9 *objectsServicesRsModelUdpObject
	if ans.Protocol.Udp != nil {
		var9 = &objectsServicesRsModelUdpObject{}
		var var10 *objectsServicesRsModelOverrideObject1
		if ans.Protocol.Udp.Override != nil {
			var10 = &objectsServicesRsModelOverrideObject1{}
			var10.Timeout = types.Int64Value(ans.Protocol.Udp.Override.Timeout)
		}
		var9.Override = var10
		var9.Port = types.StringValue(ans.Protocol.Udp.Port)
		var9.SourcePort = types.StringValue(ans.Protocol.Udp.SourcePort)
	}
	var6.Tcp = var7
	var6.Udp = var9
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.Protocol = var6
	state.Tag = EncodeStringSlice(ans.Tag)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *objectsServicesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state objectsServicesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_objects_services",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := eumQbRC.NewClient(r.client)
	input := eumQbRC.ReadInput{
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
	var var0 objectsServicesRsModelProtocolObject
	var var1 *objectsServicesRsModelTcpObject
	if ans.Protocol.Tcp != nil {
		var1 = &objectsServicesRsModelTcpObject{}
		var var2 *objectsServicesRsModelOverrideObject
		if ans.Protocol.Tcp.Override != nil {
			var2 = &objectsServicesRsModelOverrideObject{}
			var2.HalfcloseTimeout = types.Int64Value(ans.Protocol.Tcp.Override.HalfcloseTimeout)
			var2.Timeout = types.Int64Value(ans.Protocol.Tcp.Override.Timeout)
			var2.TimewaitTimeout = types.Int64Value(ans.Protocol.Tcp.Override.TimewaitTimeout)
		}
		var1.Override = var2
		var1.Port = types.StringValue(ans.Protocol.Tcp.Port)
		var1.SourcePort = types.StringValue(ans.Protocol.Tcp.SourcePort)
	}
	var var3 *objectsServicesRsModelUdpObject
	if ans.Protocol.Udp != nil {
		var3 = &objectsServicesRsModelUdpObject{}
		var var4 *objectsServicesRsModelOverrideObject1
		if ans.Protocol.Udp.Override != nil {
			var4 = &objectsServicesRsModelOverrideObject1{}
			var4.Timeout = types.Int64Value(ans.Protocol.Udp.Override.Timeout)
		}
		var3.Override = var4
		var3.Port = types.StringValue(ans.Protocol.Udp.Port)
		var3.SourcePort = types.StringValue(ans.Protocol.Udp.SourcePort)
	}
	var0.Tcp = var1
	var0.Udp = var3
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.Protocol = var0
	state.Tag = EncodeStringSlice(ans.Tag)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *objectsServicesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state objectsServicesRsModel
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
		"resource_name":               "sase_objects_services",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := eumQbRC.NewClient(r.client)
	input := eumQbRC.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 ktjCEnF.Config
	var0.Description = plan.Description.ValueString()
	var0.Name = plan.Name.ValueString()
	var var1 ktjCEnF.ProtocolObject
	var var2 *ktjCEnF.TcpObject
	if plan.Protocol.Tcp != nil {
		var2 = &ktjCEnF.TcpObject{}
		var var3 *ktjCEnF.OverrideObject
		if plan.Protocol.Tcp.Override != nil {
			var3 = &ktjCEnF.OverrideObject{}
			var3.HalfcloseTimeout = plan.Protocol.Tcp.Override.HalfcloseTimeout.ValueInt64()
			var3.Timeout = plan.Protocol.Tcp.Override.Timeout.ValueInt64()
			var3.TimewaitTimeout = plan.Protocol.Tcp.Override.TimewaitTimeout.ValueInt64()
		}
		var2.Override = var3
		var2.Port = plan.Protocol.Tcp.Port.ValueString()
		var2.SourcePort = plan.Protocol.Tcp.SourcePort.ValueString()
	}
	var1.Tcp = var2
	var var4 *ktjCEnF.UdpObject
	if plan.Protocol.Udp != nil {
		var4 = &ktjCEnF.UdpObject{}
		var var5 *ktjCEnF.OverrideObject1
		if plan.Protocol.Udp.Override != nil {
			var5 = &ktjCEnF.OverrideObject1{}
			var5.Timeout = plan.Protocol.Udp.Override.Timeout.ValueInt64()
		}
		var4.Override = var5
		var4.Port = plan.Protocol.Udp.Port.ValueString()
		var4.SourcePort = plan.Protocol.Udp.SourcePort.ValueString()
	}
	var1.Udp = var4
	var0.Protocol = var1
	var0.Tag = DecodeStringSlice(plan.Tag)
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var6 objectsServicesRsModelProtocolObject
	var var7 *objectsServicesRsModelTcpObject
	if ans.Protocol.Tcp != nil {
		var7 = &objectsServicesRsModelTcpObject{}
		var var8 *objectsServicesRsModelOverrideObject
		if ans.Protocol.Tcp.Override != nil {
			var8 = &objectsServicesRsModelOverrideObject{}
			var8.HalfcloseTimeout = types.Int64Value(ans.Protocol.Tcp.Override.HalfcloseTimeout)
			var8.Timeout = types.Int64Value(ans.Protocol.Tcp.Override.Timeout)
			var8.TimewaitTimeout = types.Int64Value(ans.Protocol.Tcp.Override.TimewaitTimeout)
		}
		var7.Override = var8
		var7.Port = types.StringValue(ans.Protocol.Tcp.Port)
		var7.SourcePort = types.StringValue(ans.Protocol.Tcp.SourcePort)
	}
	var var9 *objectsServicesRsModelUdpObject
	if ans.Protocol.Udp != nil {
		var9 = &objectsServicesRsModelUdpObject{}
		var var10 *objectsServicesRsModelOverrideObject1
		if ans.Protocol.Udp.Override != nil {
			var10 = &objectsServicesRsModelOverrideObject1{}
			var10.Timeout = types.Int64Value(ans.Protocol.Udp.Override.Timeout)
		}
		var9.Override = var10
		var9.Port = types.StringValue(ans.Protocol.Udp.Port)
		var9.SourcePort = types.StringValue(ans.Protocol.Udp.SourcePort)
	}
	var6.Tcp = var7
	var6.Udp = var9
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.Protocol = var6
	state.Tag = EncodeStringSlice(ans.Tag)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *objectsServicesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name":               "sase_objects_services",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	svc := eumQbRC.NewClient(r.client)
	input := eumQbRC.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *objectsServicesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
