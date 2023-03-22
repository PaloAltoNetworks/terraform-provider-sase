package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	uIPtsLf "github.com/paloaltonetworks/sase-go/netsec/schema/http/header/profiles"
	wiaEZmh "github.com/paloaltonetworks/sase-go/netsec/service/v1/httpheaderprofiles"

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
	_ datasource.DataSource              = &httpHeaderProfilesListDataSource{}
	_ datasource.DataSourceWithConfigure = &httpHeaderProfilesListDataSource{}
)

func NewHttpHeaderProfilesListDataSource() datasource.DataSource {
	return &httpHeaderProfilesListDataSource{}
}

type httpHeaderProfilesListDataSource struct {
	client *sase.Client
}

type httpHeaderProfilesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []httpHeaderProfilesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type httpHeaderProfilesListDsModelConfig struct {
	Description         types.String                                             `tfsdk:"description"`
	HttpHeaderInsertion []httpHeaderProfilesListDsModelHttpHeaderInsertionObject `tfsdk:"http_header_insertion"`
	ObjectId            types.String                                             `tfsdk:"object_id"`
	Name                types.String                                             `tfsdk:"name"`
}

type httpHeaderProfilesListDsModelHttpHeaderInsertionObject struct {
	Name types.String                              `tfsdk:"name"`
	Type []httpHeaderProfilesListDsModelTypeObject `tfsdk:"type"`
}

type httpHeaderProfilesListDsModelTypeObject struct {
	Domains []types.String                               `tfsdk:"domains"`
	Headers []httpHeaderProfilesListDsModelHeadersObject `tfsdk:"headers"`
	Name    types.String                                 `tfsdk:"name"`
}

type httpHeaderProfilesListDsModelHeadersObject struct {
	Header types.String `tfsdk:"header"`
	Log    types.Bool   `tfsdk:"log"`
	Name   types.String `tfsdk:"name"`
	Value  types.String `tfsdk:"value"`
}

// Metadata returns the data source type name.
func (d *httpHeaderProfilesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_http_header_profiles_list"
}

// Schema defines the schema for this listing data source.
func (d *httpHeaderProfilesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Mobile Users Container", "Mobile Users Explicit Proxy"),
				},
			},

			// Output.
			"data": dsschema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"description": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"http_header_insertion": dsschema.ListNestedAttribute{
							Description: "",
							Computed:    true,
							NestedObject: dsschema.NestedAttributeObject{
								Attributes: map[string]dsschema.Attribute{
									"name": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"type": dsschema.ListNestedAttribute{
										Description: "",
										Computed:    true,
										NestedObject: dsschema.NestedAttributeObject{
											Attributes: map[string]dsschema.Attribute{
												"domains": dsschema.ListAttribute{
													Description: "",
													Computed:    true,
													ElementType: types.StringType,
												},
												"headers": dsschema.ListNestedAttribute{
													Description: "",
													Computed:    true,
													NestedObject: dsschema.NestedAttributeObject{
														Attributes: map[string]dsschema.Attribute{
															"header": dsschema.StringAttribute{
																Description: "",
																Computed:    true,
															},
															"log": dsschema.BoolAttribute{
																Description: "",
																Computed:    true,
															},
															"name": dsschema.StringAttribute{
																Description: "",
																Computed:    true,
															},
															"value": dsschema.StringAttribute{
																Description: "",
																Computed:    true,
															},
														},
													},
												},
												"name": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
											},
										},
									},
								},
							},
						},
						"object_id": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"name": dsschema.StringAttribute{
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
func (d *httpHeaderProfilesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *httpHeaderProfilesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state httpHeaderProfilesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_http_header_profiles_list",
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
	svc := wiaEZmh.NewClient(d.client)
	input := wiaEZmh.ListInput{
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
	var var0 []httpHeaderProfilesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]httpHeaderProfilesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 httpHeaderProfilesListDsModelConfig
			var var3 []httpHeaderProfilesListDsModelHttpHeaderInsertionObject
			if len(var1.HttpHeaderInsertion) != 0 {
				var3 = make([]httpHeaderProfilesListDsModelHttpHeaderInsertionObject, 0, len(var1.HttpHeaderInsertion))
				for var4Index := range var1.HttpHeaderInsertion {
					var4 := var1.HttpHeaderInsertion[var4Index]
					var var5 httpHeaderProfilesListDsModelHttpHeaderInsertionObject
					var var6 []httpHeaderProfilesListDsModelTypeObject
					if len(var4.Type) != 0 {
						var6 = make([]httpHeaderProfilesListDsModelTypeObject, 0, len(var4.Type))
						for var7Index := range var4.Type {
							var7 := var4.Type[var7Index]
							var var8 httpHeaderProfilesListDsModelTypeObject
							var var9 []httpHeaderProfilesListDsModelHeadersObject
							if len(var7.Headers) != 0 {
								var9 = make([]httpHeaderProfilesListDsModelHeadersObject, 0, len(var7.Headers))
								for var10Index := range var7.Headers {
									var10 := var7.Headers[var10Index]
									var var11 httpHeaderProfilesListDsModelHeadersObject
									var11.Header = types.StringValue(var10.Header)
									var11.Log = types.BoolValue(var10.Log)
									var11.Name = types.StringValue(var10.Name)
									var11.Value = types.StringValue(var10.Value)
									var9 = append(var9, var11)
								}
							}
							var8.Domains = EncodeStringSlice(var7.Domains)
							var8.Headers = var9
							var8.Name = types.StringValue(var7.Name)
							var6 = append(var6, var8)
						}
					}
					var5.Name = types.StringValue(var4.Name)
					var5.Type = var6
					var3 = append(var3, var5)
				}
			}
			var2.Description = types.StringValue(var1.Description)
			var2.HttpHeaderInsertion = var3
			var2.ObjectId = types.StringValue(var1.ObjectId)
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
	_ datasource.DataSource              = &httpHeaderProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &httpHeaderProfilesDataSource{}
)

func NewHttpHeaderProfilesDataSource() datasource.DataSource {
	return &httpHeaderProfilesDataSource{}
}

type httpHeaderProfilesDataSource struct {
	client *sase.Client
}

type httpHeaderProfilesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`
	Folder   types.String `tfsdk:"folder"`

	// Output.
	// Ref: #/components/schemas/http-header-profiles
	Description         types.String                                         `tfsdk:"description"`
	HttpHeaderInsertion []httpHeaderProfilesDsModelHttpHeaderInsertionObject `tfsdk:"http_header_insertion"`
	// input omit: ObjectId
	Name types.String `tfsdk:"name"`
}

type httpHeaderProfilesDsModelHttpHeaderInsertionObject struct {
	Name types.String                          `tfsdk:"name"`
	Type []httpHeaderProfilesDsModelTypeObject `tfsdk:"type"`
}

type httpHeaderProfilesDsModelTypeObject struct {
	Domains []types.String                           `tfsdk:"domains"`
	Headers []httpHeaderProfilesDsModelHeadersObject `tfsdk:"headers"`
	Name    types.String                             `tfsdk:"name"`
}

type httpHeaderProfilesDsModelHeadersObject struct {
	Header types.String `tfsdk:"header"`
	Log    types.Bool   `tfsdk:"log"`
	Name   types.String `tfsdk:"name"`
	Value  types.String `tfsdk:"value"`
}

// Metadata returns the data source type name.
func (d *httpHeaderProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_http_header_profiles"
}

// Schema defines the schema for this listing data source.
func (d *httpHeaderProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Retrieves config for a specific item.",

		Attributes: map[string]dsschema.Attribute{
			"id": dsschema.StringAttribute{
				Description: "The object ID.",
				Computed:    true,
			},

			// Input.
			"object_id": dsschema.StringAttribute{
				Description: "The uuid of the resource",
				Required:    true,
			},
			"folder": dsschema.StringAttribute{
				Description: "The folder of the entry",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Service Connections", "Mobile Users Container", "Mobile Users Explicit Proxy"),
				},
			},

			// Output.
			"description": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"http_header_insertion": dsschema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"type": dsschema.ListNestedAttribute{
							Description: "",
							Computed:    true,
							NestedObject: dsschema.NestedAttributeObject{
								Attributes: map[string]dsschema.Attribute{
									"domains": dsschema.ListAttribute{
										Description: "",
										Computed:    true,
										ElementType: types.StringType,
									},
									"headers": dsschema.ListNestedAttribute{
										Description: "",
										Computed:    true,
										NestedObject: dsschema.NestedAttributeObject{
											Attributes: map[string]dsschema.Attribute{
												"header": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"log": dsschema.BoolAttribute{
													Description: "",
													Computed:    true,
												},
												"name": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"value": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
											},
										},
									},
									"name": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
		},
	}
}

// Configure prepares the struct.
func (d *httpHeaderProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *httpHeaderProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state httpHeaderProfilesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_http_header_profiles",
		"object_id":                   state.ObjectId.ValueString(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := wiaEZmh.NewClient(d.client)
	input := wiaEZmh.ReadInput{
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
	var var0 []httpHeaderProfilesDsModelHttpHeaderInsertionObject
	if len(ans.HttpHeaderInsertion) != 0 {
		var0 = make([]httpHeaderProfilesDsModelHttpHeaderInsertionObject, 0, len(ans.HttpHeaderInsertion))
		for var1Index := range ans.HttpHeaderInsertion {
			var1 := ans.HttpHeaderInsertion[var1Index]
			var var2 httpHeaderProfilesDsModelHttpHeaderInsertionObject
			var var3 []httpHeaderProfilesDsModelTypeObject
			if len(var1.Type) != 0 {
				var3 = make([]httpHeaderProfilesDsModelTypeObject, 0, len(var1.Type))
				for var4Index := range var1.Type {
					var4 := var1.Type[var4Index]
					var var5 httpHeaderProfilesDsModelTypeObject
					var var6 []httpHeaderProfilesDsModelHeadersObject
					if len(var4.Headers) != 0 {
						var6 = make([]httpHeaderProfilesDsModelHeadersObject, 0, len(var4.Headers))
						for var7Index := range var4.Headers {
							var7 := var4.Headers[var7Index]
							var var8 httpHeaderProfilesDsModelHeadersObject
							var8.Header = types.StringValue(var7.Header)
							var8.Log = types.BoolValue(var7.Log)
							var8.Name = types.StringValue(var7.Name)
							var8.Value = types.StringValue(var7.Value)
							var6 = append(var6, var8)
						}
					}
					var5.Domains = EncodeStringSlice(var4.Domains)
					var5.Headers = var6
					var5.Name = types.StringValue(var4.Name)
					var3 = append(var3, var5)
				}
			}
			var2.Name = types.StringValue(var1.Name)
			var2.Type = var3
			var0 = append(var0, var2)
		}
	}
	state.Description = types.StringValue(ans.Description)
	state.HttpHeaderInsertion = var0
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &httpHeaderProfilesResource{}
	_ resource.ResourceWithConfigure   = &httpHeaderProfilesResource{}
	_ resource.ResourceWithImportState = &httpHeaderProfilesResource{}
)

func NewHttpHeaderProfilesResource() resource.Resource {
	return &httpHeaderProfilesResource{}
}

type httpHeaderProfilesResource struct {
	client *sase.Client
}

type httpHeaderProfilesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/http-header-profiles
	Description         types.String                                         `tfsdk:"description"`
	HttpHeaderInsertion []httpHeaderProfilesRsModelHttpHeaderInsertionObject `tfsdk:"http_header_insertion"`
	ObjectId            types.String                                         `tfsdk:"object_id"`
	Name                types.String                                         `tfsdk:"name"`
}

type httpHeaderProfilesRsModelHttpHeaderInsertionObject struct {
	Name types.String                          `tfsdk:"name"`
	Type []httpHeaderProfilesRsModelTypeObject `tfsdk:"type"`
}

type httpHeaderProfilesRsModelTypeObject struct {
	Domains []types.String                           `tfsdk:"domains"`
	Headers []httpHeaderProfilesRsModelHeadersObject `tfsdk:"headers"`
	Name    types.String                             `tfsdk:"name"`
}

type httpHeaderProfilesRsModelHeadersObject struct {
	Header types.String `tfsdk:"header"`
	Log    types.Bool   `tfsdk:"log"`
	Name   types.String `tfsdk:"name"`
	Value  types.String `tfsdk:"value"`
}

// Metadata returns the data source type name.
func (r *httpHeaderProfilesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_http_header_profiles"
}

// Schema defines the schema for this listing data source.
func (r *httpHeaderProfilesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = rsschema.Schema{
		Description: "Retrieves config for a specific item.",

		Attributes: map[string]rsschema.Attribute{
			"id": rsschema.StringAttribute{
				Description: "The object ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			// Input.
			"folder": rsschema.StringAttribute{
				Description: "The folder of the entry",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Mobile Users Container", "Mobile Users Explicit Proxy"),
				},
			},

			"description": rsschema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
			},
			"http_header_insertion": rsschema.ListNestedAttribute{
				Description: "",
				Optional:    true,
				NestedObject: rsschema.NestedAttributeObject{
					Attributes: map[string]rsschema.Attribute{
						"name": rsschema.StringAttribute{
							Description: "",
							Required:    true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
						},
						"type": rsschema.ListNestedAttribute{
							Description: "",
							Required:    true,
							NestedObject: rsschema.NestedAttributeObject{
								Attributes: map[string]rsschema.Attribute{
									"domains": rsschema.ListAttribute{
										Description: "",
										Required:    true,
										ElementType: types.StringType,
									},
									"headers": rsschema.ListNestedAttribute{
										Description: "",
										Required:    true,
										NestedObject: rsschema.NestedAttributeObject{
											Attributes: map[string]rsschema.Attribute{
												"header": rsschema.StringAttribute{
													Description: "",
													Required:    true,
													PlanModifiers: []planmodifier.String{
														DefaultString(""),
													},
												},
												"log": rsschema.BoolAttribute{
													Description: "",
													Optional:    true,
													Computed:    true,
													PlanModifiers: []planmodifier.Bool{
														DefaultBool(false),
													},
												},
												"name": rsschema.StringAttribute{
													Description: "",
													Required:    true,
													PlanModifiers: []planmodifier.String{
														DefaultString(""),
													},
												},
												"value": rsschema.StringAttribute{
													Description: "",
													Required:    true,
													PlanModifiers: []planmodifier.String{
														DefaultString(""),
													},
												},
											},
										},
									},
									"name": rsschema.StringAttribute{
										Description: "",
										Required:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
									},
								},
							},
						},
					},
				},
			},
			"object_id": rsschema.StringAttribute{
				Description: "",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": rsschema.StringAttribute{
				Description: "",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *httpHeaderProfilesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *httpHeaderProfilesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state httpHeaderProfilesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_http_header_profiles",
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := wiaEZmh.NewClient(r.client)
	input := wiaEZmh.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 uIPtsLf.Config
	var0.Description = state.Description.ValueString()
	var var1 []uIPtsLf.HttpHeaderInsertionObject
	if len(state.HttpHeaderInsertion) != 0 {
		var1 = make([]uIPtsLf.HttpHeaderInsertionObject, 0, len(state.HttpHeaderInsertion))
		for var2Index := range state.HttpHeaderInsertion {
			var2 := state.HttpHeaderInsertion[var2Index]
			var var3 uIPtsLf.HttpHeaderInsertionObject
			var3.Name = var2.Name.ValueString()
			var var4 []uIPtsLf.TypeObject
			if len(var2.Type) != 0 {
				var4 = make([]uIPtsLf.TypeObject, 0, len(var2.Type))
				for var5Index := range var2.Type {
					var5 := var2.Type[var5Index]
					var var6 uIPtsLf.TypeObject
					var6.Domains = DecodeStringSlice(var5.Domains)
					var var7 []uIPtsLf.HeadersObject
					if len(var5.Headers) != 0 {
						var7 = make([]uIPtsLf.HeadersObject, 0, len(var5.Headers))
						for var8Index := range var5.Headers {
							var8 := var5.Headers[var8Index]
							var var9 uIPtsLf.HeadersObject
							var9.Header = var8.Header.ValueString()
							var9.Log = var8.Log.ValueBool()
							var9.Name = var8.Name.ValueString()
							var9.Value = var8.Value.ValueString()
							var7 = append(var7, var9)
						}
					}
					var6.Headers = var7
					var6.Name = var5.Name.ValueString()
					var4 = append(var4, var6)
				}
			}
			var3.Type = var4
			var1 = append(var1, var3)
		}
	}
	var0.HttpHeaderInsertion = var1
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
	var var10 []httpHeaderProfilesRsModelHttpHeaderInsertionObject
	if len(ans.HttpHeaderInsertion) != 0 {
		var10 = make([]httpHeaderProfilesRsModelHttpHeaderInsertionObject, 0, len(ans.HttpHeaderInsertion))
		for var11Index := range ans.HttpHeaderInsertion {
			var11 := ans.HttpHeaderInsertion[var11Index]
			var var12 httpHeaderProfilesRsModelHttpHeaderInsertionObject
			var var13 []httpHeaderProfilesRsModelTypeObject
			if len(var11.Type) != 0 {
				var13 = make([]httpHeaderProfilesRsModelTypeObject, 0, len(var11.Type))
				for var14Index := range var11.Type {
					var14 := var11.Type[var14Index]
					var var15 httpHeaderProfilesRsModelTypeObject
					var var16 []httpHeaderProfilesRsModelHeadersObject
					if len(var14.Headers) != 0 {
						var16 = make([]httpHeaderProfilesRsModelHeadersObject, 0, len(var14.Headers))
						for var17Index := range var14.Headers {
							var17 := var14.Headers[var17Index]
							var var18 httpHeaderProfilesRsModelHeadersObject
							var18.Header = types.StringValue(var17.Header)
							var18.Log = types.BoolValue(var17.Log)
							var18.Name = types.StringValue(var17.Name)
							var18.Value = types.StringValue(var17.Value)
							var16 = append(var16, var18)
						}
					}
					var15.Domains = EncodeStringSlice(var14.Domains)
					var15.Headers = var16
					var15.Name = types.StringValue(var14.Name)
					var13 = append(var13, var15)
				}
			}
			var12.Name = types.StringValue(var11.Name)
			var12.Type = var13
			var10 = append(var10, var12)
		}
	}
	state.Description = types.StringValue(ans.Description)
	state.HttpHeaderInsertion = var10
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *httpHeaderProfilesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state httpHeaderProfilesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_http_header_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := wiaEZmh.NewClient(r.client)
	input := wiaEZmh.ReadInput{
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
	var var0 []httpHeaderProfilesRsModelHttpHeaderInsertionObject
	if len(ans.HttpHeaderInsertion) != 0 {
		var0 = make([]httpHeaderProfilesRsModelHttpHeaderInsertionObject, 0, len(ans.HttpHeaderInsertion))
		for var1Index := range ans.HttpHeaderInsertion {
			var1 := ans.HttpHeaderInsertion[var1Index]
			var var2 httpHeaderProfilesRsModelHttpHeaderInsertionObject
			var var3 []httpHeaderProfilesRsModelTypeObject
			if len(var1.Type) != 0 {
				var3 = make([]httpHeaderProfilesRsModelTypeObject, 0, len(var1.Type))
				for var4Index := range var1.Type {
					var4 := var1.Type[var4Index]
					var var5 httpHeaderProfilesRsModelTypeObject
					var var6 []httpHeaderProfilesRsModelHeadersObject
					if len(var4.Headers) != 0 {
						var6 = make([]httpHeaderProfilesRsModelHeadersObject, 0, len(var4.Headers))
						for var7Index := range var4.Headers {
							var7 := var4.Headers[var7Index]
							var var8 httpHeaderProfilesRsModelHeadersObject
							var8.Header = types.StringValue(var7.Header)
							var8.Log = types.BoolValue(var7.Log)
							var8.Name = types.StringValue(var7.Name)
							var8.Value = types.StringValue(var7.Value)
							var6 = append(var6, var8)
						}
					}
					var5.Domains = EncodeStringSlice(var4.Domains)
					var5.Headers = var6
					var5.Name = types.StringValue(var4.Name)
					var3 = append(var3, var5)
				}
			}
			var2.Name = types.StringValue(var1.Name)
			var2.Type = var3
			var0 = append(var0, var2)
		}
	}
	state.Description = types.StringValue(ans.Description)
	state.HttpHeaderInsertion = var0
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *httpHeaderProfilesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state httpHeaderProfilesRsModel
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
		"resource_name":               "sase_http_header_profiles",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := wiaEZmh.NewClient(r.client)
	input := wiaEZmh.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 uIPtsLf.Config
	var0.Description = plan.Description.ValueString()
	var var1 []uIPtsLf.HttpHeaderInsertionObject
	if len(plan.HttpHeaderInsertion) != 0 {
		var1 = make([]uIPtsLf.HttpHeaderInsertionObject, 0, len(plan.HttpHeaderInsertion))
		for var2Index := range plan.HttpHeaderInsertion {
			var2 := plan.HttpHeaderInsertion[var2Index]
			var var3 uIPtsLf.HttpHeaderInsertionObject
			var3.Name = var2.Name.ValueString()
			var var4 []uIPtsLf.TypeObject
			if len(var2.Type) != 0 {
				var4 = make([]uIPtsLf.TypeObject, 0, len(var2.Type))
				for var5Index := range var2.Type {
					var5 := var2.Type[var5Index]
					var var6 uIPtsLf.TypeObject
					var6.Domains = DecodeStringSlice(var5.Domains)
					var var7 []uIPtsLf.HeadersObject
					if len(var5.Headers) != 0 {
						var7 = make([]uIPtsLf.HeadersObject, 0, len(var5.Headers))
						for var8Index := range var5.Headers {
							var8 := var5.Headers[var8Index]
							var var9 uIPtsLf.HeadersObject
							var9.Header = var8.Header.ValueString()
							var9.Log = var8.Log.ValueBool()
							var9.Name = var8.Name.ValueString()
							var9.Value = var8.Value.ValueString()
							var7 = append(var7, var9)
						}
					}
					var6.Headers = var7
					var6.Name = var5.Name.ValueString()
					var4 = append(var4, var6)
				}
			}
			var3.Type = var4
			var1 = append(var1, var3)
		}
	}
	var0.HttpHeaderInsertion = var1
	var0.Name = plan.Name.ValueString()
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var10 []httpHeaderProfilesRsModelHttpHeaderInsertionObject
	if len(ans.HttpHeaderInsertion) != 0 {
		var10 = make([]httpHeaderProfilesRsModelHttpHeaderInsertionObject, 0, len(ans.HttpHeaderInsertion))
		for var11Index := range ans.HttpHeaderInsertion {
			var11 := ans.HttpHeaderInsertion[var11Index]
			var var12 httpHeaderProfilesRsModelHttpHeaderInsertionObject
			var var13 []httpHeaderProfilesRsModelTypeObject
			if len(var11.Type) != 0 {
				var13 = make([]httpHeaderProfilesRsModelTypeObject, 0, len(var11.Type))
				for var14Index := range var11.Type {
					var14 := var11.Type[var14Index]
					var var15 httpHeaderProfilesRsModelTypeObject
					var var16 []httpHeaderProfilesRsModelHeadersObject
					if len(var14.Headers) != 0 {
						var16 = make([]httpHeaderProfilesRsModelHeadersObject, 0, len(var14.Headers))
						for var17Index := range var14.Headers {
							var17 := var14.Headers[var17Index]
							var var18 httpHeaderProfilesRsModelHeadersObject
							var18.Header = types.StringValue(var17.Header)
							var18.Log = types.BoolValue(var17.Log)
							var18.Name = types.StringValue(var17.Name)
							var18.Value = types.StringValue(var17.Value)
							var16 = append(var16, var18)
						}
					}
					var15.Domains = EncodeStringSlice(var14.Domains)
					var15.Headers = var16
					var15.Name = types.StringValue(var14.Name)
					var13 = append(var13, var15)
				}
			}
			var12.Name = types.StringValue(var11.Name)
			var12.Type = var13
			var10 = append(var10, var12)
		}
	}
	state.Description = types.StringValue(ans.Description)
	state.HttpHeaderInsertion = var10
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *httpHeaderProfilesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name":               "sase_http_header_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	svc := wiaEZmh.NewClient(r.client)
	input := wiaEZmh.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *httpHeaderProfilesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
