package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	sRAOviP "github.com/paloaltonetworks/sase-go/netsec/schema/objects/external/dynamic/lists"
	iHJqznH "github.com/paloaltonetworks/sase-go/netsec/service/v1/externaldynamiclists"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
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
	_ datasource.DataSource              = &objectsExternalDynamicListsListDataSource{}
	_ datasource.DataSourceWithConfigure = &objectsExternalDynamicListsListDataSource{}
)

func NewObjectsExternalDynamicListsListDataSource() datasource.DataSource {
	return &objectsExternalDynamicListsListDataSource{}
}

type objectsExternalDynamicListsListDataSource struct {
	client *sase.Client
}

type objectsExternalDynamicListsListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []objectsExternalDynamicListsListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type objectsExternalDynamicListsListDsModelConfig struct {
	ObjectId types.String                                     `tfsdk:"object_id"`
	Name     types.String                                     `tfsdk:"name"`
	Type     objectsExternalDynamicListsListDsModelTypeObject `tfsdk:"type"`
}

type objectsExternalDynamicListsListDsModelTypeObject struct {
	Domain        *objectsExternalDynamicListsListDsModelDomainObject        `tfsdk:"domain"`
	Imei          *objectsExternalDynamicListsListDsModelImeiObject          `tfsdk:"imei"`
	Imsi          *objectsExternalDynamicListsListDsModelImsiObject          `tfsdk:"imsi"`
	Ip            *objectsExternalDynamicListsListDsModelIpObject            `tfsdk:"ip"`
	PredefinedIp  *objectsExternalDynamicListsListDsModelPredefinedIpObject  `tfsdk:"predefined_ip"`
	PredefinedUrl *objectsExternalDynamicListsListDsModelPredefinedUrlObject `tfsdk:"predefined_url"`
	Url           *objectsExternalDynamicListsListDsModelUrlObject           `tfsdk:"url"`
}

type objectsExternalDynamicListsListDsModelDomainObject struct {
	Auth               *objectsExternalDynamicListsListDsModelAuthObject     `tfsdk:"auth"`
	CertificateProfile types.String                                          `tfsdk:"certificate_profile"`
	Description        types.String                                          `tfsdk:"description"`
	ExceptionList      []types.String                                        `tfsdk:"exception_list"`
	ExpandDomain       types.Bool                                            `tfsdk:"expand_domain"`
	Recurring          objectsExternalDynamicListsListDsModelRecurringObject `tfsdk:"recurring"`
	Url                types.String                                          `tfsdk:"url"`
}

type objectsExternalDynamicListsListDsModelAuthObject struct {
	Password types.String `tfsdk:"password"`
	Username types.String `tfsdk:"username"`
}

type objectsExternalDynamicListsListDsModelRecurringObject struct {
	Daily      *objectsExternalDynamicListsListDsModelDailyObject   `tfsdk:"daily"`
	FiveMinute types.Bool                                           `tfsdk:"five_minute"`
	Hourly     types.Bool                                           `tfsdk:"hourly"`
	Monthly    *objectsExternalDynamicListsListDsModelMonthlyObject `tfsdk:"monthly"`
	Weekly     *objectsExternalDynamicListsListDsModelWeeklyObject  `tfsdk:"weekly"`
}

type objectsExternalDynamicListsListDsModelDailyObject struct {
	At types.String `tfsdk:"at"`
}

type objectsExternalDynamicListsListDsModelMonthlyObject struct {
	At         types.String `tfsdk:"at"`
	DayOfMonth types.Int64  `tfsdk:"day_of_month"`
}

type objectsExternalDynamicListsListDsModelWeeklyObject struct {
	At        types.String `tfsdk:"at"`
	DayOfWeek types.String `tfsdk:"day_of_week"`
}

type objectsExternalDynamicListsListDsModelImeiObject struct {
	Auth               *objectsExternalDynamicListsListDsModelAuthObject     `tfsdk:"auth"`
	CertificateProfile types.String                                          `tfsdk:"certificate_profile"`
	Description        types.String                                          `tfsdk:"description"`
	ExceptionList      []types.String                                        `tfsdk:"exception_list"`
	Recurring          objectsExternalDynamicListsListDsModelRecurringObject `tfsdk:"recurring"`
	Url                types.String                                          `tfsdk:"url"`
}

type objectsExternalDynamicListsListDsModelImsiObject struct {
	Auth               *objectsExternalDynamicListsListDsModelAuthObject     `tfsdk:"auth"`
	CertificateProfile types.String                                          `tfsdk:"certificate_profile"`
	Description        types.String                                          `tfsdk:"description"`
	ExceptionList      []types.String                                        `tfsdk:"exception_list"`
	Recurring          objectsExternalDynamicListsListDsModelRecurringObject `tfsdk:"recurring"`
	Url                types.String                                          `tfsdk:"url"`
}

type objectsExternalDynamicListsListDsModelIpObject struct {
	Auth               *objectsExternalDynamicListsListDsModelAuthObject     `tfsdk:"auth"`
	CertificateProfile types.String                                          `tfsdk:"certificate_profile"`
	Description        types.String                                          `tfsdk:"description"`
	ExceptionList      []types.String                                        `tfsdk:"exception_list"`
	Recurring          objectsExternalDynamicListsListDsModelRecurringObject `tfsdk:"recurring"`
	Url                types.String                                          `tfsdk:"url"`
}

type objectsExternalDynamicListsListDsModelPredefinedIpObject struct {
	Description   types.String   `tfsdk:"description"`
	ExceptionList []types.String `tfsdk:"exception_list"`
	Url           types.String   `tfsdk:"url"`
}

type objectsExternalDynamicListsListDsModelPredefinedUrlObject struct {
	Description   types.String   `tfsdk:"description"`
	ExceptionList []types.String `tfsdk:"exception_list"`
	Url           types.String   `tfsdk:"url"`
}

type objectsExternalDynamicListsListDsModelUrlObject struct {
	Auth               *objectsExternalDynamicListsListDsModelAuthObject     `tfsdk:"auth"`
	CertificateProfile types.String                                          `tfsdk:"certificate_profile"`
	Description        types.String                                          `tfsdk:"description"`
	ExceptionList      []types.String                                        `tfsdk:"exception_list"`
	Recurring          objectsExternalDynamicListsListDsModelRecurringObject `tfsdk:"recurring"`
	Url                types.String                                          `tfsdk:"url"`
}

// Metadata returns the data source type name.
func (d *objectsExternalDynamicListsListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_external_dynamic_lists_list"
}

// Schema defines the schema for this listing data source.
func (d *objectsExternalDynamicListsListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"object_id": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"type": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"domain": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"auth": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"password": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"username": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
											},
										},
										"certificate_profile": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"description": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"exception_list": dsschema.ListAttribute{
											Description: "",
											Computed:    true,
											ElementType: types.StringType,
										},
										"expand_domain": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
										"recurring": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"daily": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"at": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
													},
												},
												"five_minute": dsschema.BoolAttribute{
													Description: "",
													Computed:    true,
												},
												"hourly": dsschema.BoolAttribute{
													Description: "",
													Computed:    true,
												},
												"monthly": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"at": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
														"day_of_month": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
													},
												},
												"weekly": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"at": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
														"day_of_week": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
													},
												},
											},
										},
										"url": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
								"imei": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"auth": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"password": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"username": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
											},
										},
										"certificate_profile": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"description": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"exception_list": dsschema.ListAttribute{
											Description: "",
											Computed:    true,
											ElementType: types.StringType,
										},
										"recurring": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"daily": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"at": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
													},
												},
												"five_minute": dsschema.BoolAttribute{
													Description: "",
													Computed:    true,
												},
												"hourly": dsschema.BoolAttribute{
													Description: "",
													Computed:    true,
												},
												"monthly": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"at": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
														"day_of_month": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
													},
												},
												"weekly": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"at": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
														"day_of_week": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
													},
												},
											},
										},
										"url": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
								"imsi": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"auth": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"password": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"username": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
											},
										},
										"certificate_profile": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"description": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"exception_list": dsschema.ListAttribute{
											Description: "",
											Computed:    true,
											ElementType: types.StringType,
										},
										"recurring": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"daily": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"at": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
													},
												},
												"five_minute": dsschema.BoolAttribute{
													Description: "",
													Computed:    true,
												},
												"hourly": dsschema.BoolAttribute{
													Description: "",
													Computed:    true,
												},
												"monthly": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"at": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
														"day_of_month": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
													},
												},
												"weekly": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"at": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
														"day_of_week": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
													},
												},
											},
										},
										"url": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
								"ip": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"auth": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"password": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"username": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
											},
										},
										"certificate_profile": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"description": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"exception_list": dsschema.ListAttribute{
											Description: "",
											Computed:    true,
											ElementType: types.StringType,
										},
										"recurring": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"daily": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"at": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
													},
												},
												"five_minute": dsschema.BoolAttribute{
													Description: "",
													Computed:    true,
												},
												"hourly": dsschema.BoolAttribute{
													Description: "",
													Computed:    true,
												},
												"monthly": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"at": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
														"day_of_month": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
													},
												},
												"weekly": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"at": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
														"day_of_week": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
													},
												},
											},
										},
										"url": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
								"predefined_ip": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"description": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"exception_list": dsschema.ListAttribute{
											Description: "",
											Computed:    true,
											ElementType: types.StringType,
										},
										"url": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
								"predefined_url": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"description": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"exception_list": dsschema.ListAttribute{
											Description: "",
											Computed:    true,
											ElementType: types.StringType,
										},
										"url": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
								"url": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"auth": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"password": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"username": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
											},
										},
										"certificate_profile": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"description": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"exception_list": dsschema.ListAttribute{
											Description: "",
											Computed:    true,
											ElementType: types.StringType,
										},
										"recurring": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"daily": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"at": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
													},
												},
												"five_minute": dsschema.BoolAttribute{
													Description: "",
													Computed:    true,
												},
												"hourly": dsschema.BoolAttribute{
													Description: "",
													Computed:    true,
												},
												"monthly": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"at": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
														"day_of_month": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
													},
												},
												"weekly": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"at": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
														"day_of_week": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
													},
												},
											},
										},
										"url": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
							},
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
func (d *objectsExternalDynamicListsListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *objectsExternalDynamicListsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state objectsExternalDynamicListsListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_objects_external_dynamic_lists_list",
		"limit":            state.Limit.ValueInt64(),
		"has_limit":        !state.Limit.IsNull(),
		"offset":           state.Offset.ValueInt64(),
		"has_offset":       !state.Offset.IsNull(),
		"name":             state.Name.ValueString(),
		"has_name":         !state.Name.IsNull(),
		"folder":           state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := iHJqznH.NewClient(d.client)
	input := iHJqznH.ListInput{
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
	var var0 []objectsExternalDynamicListsListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]objectsExternalDynamicListsListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 objectsExternalDynamicListsListDsModelConfig
			var var3 objectsExternalDynamicListsListDsModelTypeObject
			var var4 *objectsExternalDynamicListsListDsModelDomainObject
			if var1.Type.Domain != nil {
				var4 = &objectsExternalDynamicListsListDsModelDomainObject{}
				var var5 *objectsExternalDynamicListsListDsModelAuthObject
				if var1.Type.Domain.Auth != nil {
					var5 = &objectsExternalDynamicListsListDsModelAuthObject{}
					var5.Password = types.StringValue(var1.Type.Domain.Auth.Password)
					var5.Username = types.StringValue(var1.Type.Domain.Auth.Username)
				}
				var var6 objectsExternalDynamicListsListDsModelRecurringObject
				var var7 *objectsExternalDynamicListsListDsModelDailyObject
				if var1.Type.Domain.Recurring.Daily != nil {
					var7 = &objectsExternalDynamicListsListDsModelDailyObject{}
					var7.At = types.StringValue(var1.Type.Domain.Recurring.Daily.At)
				}
				var var8 *objectsExternalDynamicListsListDsModelMonthlyObject
				if var1.Type.Domain.Recurring.Monthly != nil {
					var8 = &objectsExternalDynamicListsListDsModelMonthlyObject{}
					var8.At = types.StringValue(var1.Type.Domain.Recurring.Monthly.At)
					var8.DayOfMonth = types.Int64Value(var1.Type.Domain.Recurring.Monthly.DayOfMonth)
				}
				var var9 *objectsExternalDynamicListsListDsModelWeeklyObject
				if var1.Type.Domain.Recurring.Weekly != nil {
					var9 = &objectsExternalDynamicListsListDsModelWeeklyObject{}
					var9.At = types.StringValue(var1.Type.Domain.Recurring.Weekly.At)
					var9.DayOfWeek = types.StringValue(var1.Type.Domain.Recurring.Weekly.DayOfWeek)
				}
				var6.Daily = var7
				if var1.Type.Domain.Recurring.FiveMinute != nil {
					var6.FiveMinute = types.BoolValue(true)
				}
				if var1.Type.Domain.Recurring.Hourly != nil {
					var6.Hourly = types.BoolValue(true)
				}
				var6.Monthly = var8
				var6.Weekly = var9
				var4.Auth = var5
				var4.CertificateProfile = types.StringValue(var1.Type.Domain.CertificateProfile)
				var4.Description = types.StringValue(var1.Type.Domain.Description)
				var4.ExceptionList = EncodeStringSlice(var1.Type.Domain.ExceptionList)
				var4.ExpandDomain = types.BoolValue(var1.Type.Domain.ExpandDomain)
				var4.Recurring = var6
				var4.Url = types.StringValue(var1.Type.Domain.Url)
			}
			var var10 *objectsExternalDynamicListsListDsModelImeiObject
			if var1.Type.Imei != nil {
				var10 = &objectsExternalDynamicListsListDsModelImeiObject{}
				var var11 *objectsExternalDynamicListsListDsModelAuthObject
				if var1.Type.Imei.Auth != nil {
					var11 = &objectsExternalDynamicListsListDsModelAuthObject{}
					var11.Password = types.StringValue(var1.Type.Imei.Auth.Password)
					var11.Username = types.StringValue(var1.Type.Imei.Auth.Username)
				}
				var var12 objectsExternalDynamicListsListDsModelRecurringObject
				var var13 *objectsExternalDynamicListsListDsModelDailyObject
				if var1.Type.Imei.Recurring.Daily != nil {
					var13 = &objectsExternalDynamicListsListDsModelDailyObject{}
					var13.At = types.StringValue(var1.Type.Imei.Recurring.Daily.At)
				}
				var var14 *objectsExternalDynamicListsListDsModelMonthlyObject
				if var1.Type.Imei.Recurring.Monthly != nil {
					var14 = &objectsExternalDynamicListsListDsModelMonthlyObject{}
					var14.At = types.StringValue(var1.Type.Imei.Recurring.Monthly.At)
					var14.DayOfMonth = types.Int64Value(var1.Type.Imei.Recurring.Monthly.DayOfMonth)
				}
				var var15 *objectsExternalDynamicListsListDsModelWeeklyObject
				if var1.Type.Imei.Recurring.Weekly != nil {
					var15 = &objectsExternalDynamicListsListDsModelWeeklyObject{}
					var15.At = types.StringValue(var1.Type.Imei.Recurring.Weekly.At)
					var15.DayOfWeek = types.StringValue(var1.Type.Imei.Recurring.Weekly.DayOfWeek)
				}
				var12.Daily = var13
				if var1.Type.Imei.Recurring.FiveMinute != nil {
					var12.FiveMinute = types.BoolValue(true)
				}
				if var1.Type.Imei.Recurring.Hourly != nil {
					var12.Hourly = types.BoolValue(true)
				}
				var12.Monthly = var14
				var12.Weekly = var15
				var10.Auth = var11
				var10.CertificateProfile = types.StringValue(var1.Type.Imei.CertificateProfile)
				var10.Description = types.StringValue(var1.Type.Imei.Description)
				var10.ExceptionList = EncodeStringSlice(var1.Type.Imei.ExceptionList)
				var10.Recurring = var12
				var10.Url = types.StringValue(var1.Type.Imei.Url)
			}
			var var16 *objectsExternalDynamicListsListDsModelImsiObject
			if var1.Type.Imsi != nil {
				var16 = &objectsExternalDynamicListsListDsModelImsiObject{}
				var var17 *objectsExternalDynamicListsListDsModelAuthObject
				if var1.Type.Imsi.Auth != nil {
					var17 = &objectsExternalDynamicListsListDsModelAuthObject{}
					var17.Password = types.StringValue(var1.Type.Imsi.Auth.Password)
					var17.Username = types.StringValue(var1.Type.Imsi.Auth.Username)
				}
				var var18 objectsExternalDynamicListsListDsModelRecurringObject
				var var19 *objectsExternalDynamicListsListDsModelDailyObject
				if var1.Type.Imsi.Recurring.Daily != nil {
					var19 = &objectsExternalDynamicListsListDsModelDailyObject{}
					var19.At = types.StringValue(var1.Type.Imsi.Recurring.Daily.At)
				}
				var var20 *objectsExternalDynamicListsListDsModelMonthlyObject
				if var1.Type.Imsi.Recurring.Monthly != nil {
					var20 = &objectsExternalDynamicListsListDsModelMonthlyObject{}
					var20.At = types.StringValue(var1.Type.Imsi.Recurring.Monthly.At)
					var20.DayOfMonth = types.Int64Value(var1.Type.Imsi.Recurring.Monthly.DayOfMonth)
				}
				var var21 *objectsExternalDynamicListsListDsModelWeeklyObject
				if var1.Type.Imsi.Recurring.Weekly != nil {
					var21 = &objectsExternalDynamicListsListDsModelWeeklyObject{}
					var21.At = types.StringValue(var1.Type.Imsi.Recurring.Weekly.At)
					var21.DayOfWeek = types.StringValue(var1.Type.Imsi.Recurring.Weekly.DayOfWeek)
				}
				var18.Daily = var19
				if var1.Type.Imsi.Recurring.FiveMinute != nil {
					var18.FiveMinute = types.BoolValue(true)
				}
				if var1.Type.Imsi.Recurring.Hourly != nil {
					var18.Hourly = types.BoolValue(true)
				}
				var18.Monthly = var20
				var18.Weekly = var21
				var16.Auth = var17
				var16.CertificateProfile = types.StringValue(var1.Type.Imsi.CertificateProfile)
				var16.Description = types.StringValue(var1.Type.Imsi.Description)
				var16.ExceptionList = EncodeStringSlice(var1.Type.Imsi.ExceptionList)
				var16.Recurring = var18
				var16.Url = types.StringValue(var1.Type.Imsi.Url)
			}
			var var22 *objectsExternalDynamicListsListDsModelIpObject
			if var1.Type.Ip != nil {
				var22 = &objectsExternalDynamicListsListDsModelIpObject{}
				var var23 *objectsExternalDynamicListsListDsModelAuthObject
				if var1.Type.Ip.Auth != nil {
					var23 = &objectsExternalDynamicListsListDsModelAuthObject{}
					var23.Password = types.StringValue(var1.Type.Ip.Auth.Password)
					var23.Username = types.StringValue(var1.Type.Ip.Auth.Username)
				}
				var var24 objectsExternalDynamicListsListDsModelRecurringObject
				var var25 *objectsExternalDynamicListsListDsModelDailyObject
				if var1.Type.Ip.Recurring.Daily != nil {
					var25 = &objectsExternalDynamicListsListDsModelDailyObject{}
					var25.At = types.StringValue(var1.Type.Ip.Recurring.Daily.At)
				}
				var var26 *objectsExternalDynamicListsListDsModelMonthlyObject
				if var1.Type.Ip.Recurring.Monthly != nil {
					var26 = &objectsExternalDynamicListsListDsModelMonthlyObject{}
					var26.At = types.StringValue(var1.Type.Ip.Recurring.Monthly.At)
					var26.DayOfMonth = types.Int64Value(var1.Type.Ip.Recurring.Monthly.DayOfMonth)
				}
				var var27 *objectsExternalDynamicListsListDsModelWeeklyObject
				if var1.Type.Ip.Recurring.Weekly != nil {
					var27 = &objectsExternalDynamicListsListDsModelWeeklyObject{}
					var27.At = types.StringValue(var1.Type.Ip.Recurring.Weekly.At)
					var27.DayOfWeek = types.StringValue(var1.Type.Ip.Recurring.Weekly.DayOfWeek)
				}
				var24.Daily = var25
				if var1.Type.Ip.Recurring.FiveMinute != nil {
					var24.FiveMinute = types.BoolValue(true)
				}
				if var1.Type.Ip.Recurring.Hourly != nil {
					var24.Hourly = types.BoolValue(true)
				}
				var24.Monthly = var26
				var24.Weekly = var27
				var22.Auth = var23
				var22.CertificateProfile = types.StringValue(var1.Type.Ip.CertificateProfile)
				var22.Description = types.StringValue(var1.Type.Ip.Description)
				var22.ExceptionList = EncodeStringSlice(var1.Type.Ip.ExceptionList)
				var22.Recurring = var24
				var22.Url = types.StringValue(var1.Type.Ip.Url)
			}
			var var28 *objectsExternalDynamicListsListDsModelPredefinedIpObject
			if var1.Type.PredefinedIp != nil {
				var28 = &objectsExternalDynamicListsListDsModelPredefinedIpObject{}
				var28.Description = types.StringValue(var1.Type.PredefinedIp.Description)
				var28.ExceptionList = EncodeStringSlice(var1.Type.PredefinedIp.ExceptionList)
				var28.Url = types.StringValue(var1.Type.PredefinedIp.Url)
			}
			var var29 *objectsExternalDynamicListsListDsModelPredefinedUrlObject
			if var1.Type.PredefinedUrl != nil {
				var29 = &objectsExternalDynamicListsListDsModelPredefinedUrlObject{}
				var29.Description = types.StringValue(var1.Type.PredefinedUrl.Description)
				var29.ExceptionList = EncodeStringSlice(var1.Type.PredefinedUrl.ExceptionList)
				var29.Url = types.StringValue(var1.Type.PredefinedUrl.Url)
			}
			var var30 *objectsExternalDynamicListsListDsModelUrlObject
			if var1.Type.Url != nil {
				var30 = &objectsExternalDynamicListsListDsModelUrlObject{}
				var var31 *objectsExternalDynamicListsListDsModelAuthObject
				if var1.Type.Url.Auth != nil {
					var31 = &objectsExternalDynamicListsListDsModelAuthObject{}
					var31.Password = types.StringValue(var1.Type.Url.Auth.Password)
					var31.Username = types.StringValue(var1.Type.Url.Auth.Username)
				}
				var var32 objectsExternalDynamicListsListDsModelRecurringObject
				var var33 *objectsExternalDynamicListsListDsModelDailyObject
				if var1.Type.Url.Recurring.Daily != nil {
					var33 = &objectsExternalDynamicListsListDsModelDailyObject{}
					var33.At = types.StringValue(var1.Type.Url.Recurring.Daily.At)
				}
				var var34 *objectsExternalDynamicListsListDsModelMonthlyObject
				if var1.Type.Url.Recurring.Monthly != nil {
					var34 = &objectsExternalDynamicListsListDsModelMonthlyObject{}
					var34.At = types.StringValue(var1.Type.Url.Recurring.Monthly.At)
					var34.DayOfMonth = types.Int64Value(var1.Type.Url.Recurring.Monthly.DayOfMonth)
				}
				var var35 *objectsExternalDynamicListsListDsModelWeeklyObject
				if var1.Type.Url.Recurring.Weekly != nil {
					var35 = &objectsExternalDynamicListsListDsModelWeeklyObject{}
					var35.At = types.StringValue(var1.Type.Url.Recurring.Weekly.At)
					var35.DayOfWeek = types.StringValue(var1.Type.Url.Recurring.Weekly.DayOfWeek)
				}
				var32.Daily = var33
				if var1.Type.Url.Recurring.FiveMinute != nil {
					var32.FiveMinute = types.BoolValue(true)
				}
				if var1.Type.Url.Recurring.Hourly != nil {
					var32.Hourly = types.BoolValue(true)
				}
				var32.Monthly = var34
				var32.Weekly = var35
				var30.Auth = var31
				var30.CertificateProfile = types.StringValue(var1.Type.Url.CertificateProfile)
				var30.Description = types.StringValue(var1.Type.Url.Description)
				var30.ExceptionList = EncodeStringSlice(var1.Type.Url.ExceptionList)
				var30.Recurring = var32
				var30.Url = types.StringValue(var1.Type.Url.Url)
			}
			var3.Domain = var4
			var3.Imei = var10
			var3.Imsi = var16
			var3.Ip = var22
			var3.PredefinedIp = var28
			var3.PredefinedUrl = var29
			var3.Url = var30
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Name = types.StringValue(var1.Name)
			var2.Type = var3
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
	_ datasource.DataSource              = &objectsExternalDynamicListsDataSource{}
	_ datasource.DataSourceWithConfigure = &objectsExternalDynamicListsDataSource{}
)

func NewObjectsExternalDynamicListsDataSource() datasource.DataSource {
	return &objectsExternalDynamicListsDataSource{}
}

type objectsExternalDynamicListsDataSource struct {
	client *sase.Client
}

type objectsExternalDynamicListsDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`

	// Output.
	// Ref: #/components/schemas/objects-external-dynamic-lists
	// input omit: ObjectId
	Name types.String                                 `tfsdk:"name"`
	Type objectsExternalDynamicListsDsModelTypeObject `tfsdk:"type"`
}

type objectsExternalDynamicListsDsModelTypeObject struct {
	Domain        *objectsExternalDynamicListsDsModelDomainObject        `tfsdk:"domain"`
	Imei          *objectsExternalDynamicListsDsModelImeiObject          `tfsdk:"imei"`
	Imsi          *objectsExternalDynamicListsDsModelImsiObject          `tfsdk:"imsi"`
	Ip            *objectsExternalDynamicListsDsModelIpObject            `tfsdk:"ip"`
	PredefinedIp  *objectsExternalDynamicListsDsModelPredefinedIpObject  `tfsdk:"predefined_ip"`
	PredefinedUrl *objectsExternalDynamicListsDsModelPredefinedUrlObject `tfsdk:"predefined_url"`
	Url           *objectsExternalDynamicListsDsModelUrlObject           `tfsdk:"url"`
}

type objectsExternalDynamicListsDsModelDomainObject struct {
	Auth               *objectsExternalDynamicListsDsModelAuthObject     `tfsdk:"auth"`
	CertificateProfile types.String                                      `tfsdk:"certificate_profile"`
	Description        types.String                                      `tfsdk:"description"`
	ExceptionList      []types.String                                    `tfsdk:"exception_list"`
	ExpandDomain       types.Bool                                        `tfsdk:"expand_domain"`
	Recurring          objectsExternalDynamicListsDsModelRecurringObject `tfsdk:"recurring"`
	Url                types.String                                      `tfsdk:"url"`
}

type objectsExternalDynamicListsDsModelAuthObject struct {
	Password types.String `tfsdk:"password"`
	Username types.String `tfsdk:"username"`
}

type objectsExternalDynamicListsDsModelRecurringObject struct {
	Daily      *objectsExternalDynamicListsDsModelDailyObject   `tfsdk:"daily"`
	FiveMinute types.Bool                                       `tfsdk:"five_minute"`
	Hourly     types.Bool                                       `tfsdk:"hourly"`
	Monthly    *objectsExternalDynamicListsDsModelMonthlyObject `tfsdk:"monthly"`
	Weekly     *objectsExternalDynamicListsDsModelWeeklyObject  `tfsdk:"weekly"`
}

type objectsExternalDynamicListsDsModelDailyObject struct {
	At types.String `tfsdk:"at"`
}

type objectsExternalDynamicListsDsModelMonthlyObject struct {
	At         types.String `tfsdk:"at"`
	DayOfMonth types.Int64  `tfsdk:"day_of_month"`
}

type objectsExternalDynamicListsDsModelWeeklyObject struct {
	At        types.String `tfsdk:"at"`
	DayOfWeek types.String `tfsdk:"day_of_week"`
}

type objectsExternalDynamicListsDsModelImeiObject struct {
	Auth               *objectsExternalDynamicListsDsModelAuthObject     `tfsdk:"auth"`
	CertificateProfile types.String                                      `tfsdk:"certificate_profile"`
	Description        types.String                                      `tfsdk:"description"`
	ExceptionList      []types.String                                    `tfsdk:"exception_list"`
	Recurring          objectsExternalDynamicListsDsModelRecurringObject `tfsdk:"recurring"`
	Url                types.String                                      `tfsdk:"url"`
}

type objectsExternalDynamicListsDsModelImsiObject struct {
	Auth               *objectsExternalDynamicListsDsModelAuthObject     `tfsdk:"auth"`
	CertificateProfile types.String                                      `tfsdk:"certificate_profile"`
	Description        types.String                                      `tfsdk:"description"`
	ExceptionList      []types.String                                    `tfsdk:"exception_list"`
	Recurring          objectsExternalDynamicListsDsModelRecurringObject `tfsdk:"recurring"`
	Url                types.String                                      `tfsdk:"url"`
}

type objectsExternalDynamicListsDsModelIpObject struct {
	Auth               *objectsExternalDynamicListsDsModelAuthObject     `tfsdk:"auth"`
	CertificateProfile types.String                                      `tfsdk:"certificate_profile"`
	Description        types.String                                      `tfsdk:"description"`
	ExceptionList      []types.String                                    `tfsdk:"exception_list"`
	Recurring          objectsExternalDynamicListsDsModelRecurringObject `tfsdk:"recurring"`
	Url                types.String                                      `tfsdk:"url"`
}

type objectsExternalDynamicListsDsModelPredefinedIpObject struct {
	Description   types.String   `tfsdk:"description"`
	ExceptionList []types.String `tfsdk:"exception_list"`
	Url           types.String   `tfsdk:"url"`
}

type objectsExternalDynamicListsDsModelPredefinedUrlObject struct {
	Description   types.String   `tfsdk:"description"`
	ExceptionList []types.String `tfsdk:"exception_list"`
	Url           types.String   `tfsdk:"url"`
}

type objectsExternalDynamicListsDsModelUrlObject struct {
	Auth               *objectsExternalDynamicListsDsModelAuthObject     `tfsdk:"auth"`
	CertificateProfile types.String                                      `tfsdk:"certificate_profile"`
	Description        types.String                                      `tfsdk:"description"`
	ExceptionList      []types.String                                    `tfsdk:"exception_list"`
	Recurring          objectsExternalDynamicListsDsModelRecurringObject `tfsdk:"recurring"`
	Url                types.String                                      `tfsdk:"url"`
}

// Metadata returns the data source type name.
func (d *objectsExternalDynamicListsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_external_dynamic_lists"
}

// Schema defines the schema for this listing data source.
func (d *objectsExternalDynamicListsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

			// Output.
			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"type": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"domain": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"auth": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"password": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"username": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
							"certificate_profile": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"description": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"exception_list": dsschema.ListAttribute{
								Description: "",
								Computed:    true,
								ElementType: types.StringType,
							},
							"expand_domain": dsschema.BoolAttribute{
								Description: "",
								Computed:    true,
							},
							"recurring": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"daily": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"at": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
										},
									},
									"five_minute": dsschema.BoolAttribute{
										Description: "",
										Computed:    true,
									},
									"hourly": dsschema.BoolAttribute{
										Description: "",
										Computed:    true,
									},
									"monthly": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"at": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
											"day_of_month": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
										},
									},
									"weekly": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"at": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
											"day_of_week": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
										},
									},
								},
							},
							"url": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
						},
					},
					"imei": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"auth": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"password": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"username": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
							"certificate_profile": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"description": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"exception_list": dsschema.ListAttribute{
								Description: "",
								Computed:    true,
								ElementType: types.StringType,
							},
							"recurring": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"daily": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"at": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
										},
									},
									"five_minute": dsschema.BoolAttribute{
										Description: "",
										Computed:    true,
									},
									"hourly": dsschema.BoolAttribute{
										Description: "",
										Computed:    true,
									},
									"monthly": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"at": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
											"day_of_month": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
										},
									},
									"weekly": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"at": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
											"day_of_week": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
										},
									},
								},
							},
							"url": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
						},
					},
					"imsi": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"auth": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"password": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"username": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
							"certificate_profile": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"description": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"exception_list": dsschema.ListAttribute{
								Description: "",
								Computed:    true,
								ElementType: types.StringType,
							},
							"recurring": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"daily": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"at": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
										},
									},
									"five_minute": dsschema.BoolAttribute{
										Description: "",
										Computed:    true,
									},
									"hourly": dsschema.BoolAttribute{
										Description: "",
										Computed:    true,
									},
									"monthly": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"at": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
											"day_of_month": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
										},
									},
									"weekly": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"at": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
											"day_of_week": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
										},
									},
								},
							},
							"url": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
						},
					},
					"ip": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"auth": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"password": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"username": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
							"certificate_profile": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"description": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"exception_list": dsschema.ListAttribute{
								Description: "",
								Computed:    true,
								ElementType: types.StringType,
							},
							"recurring": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"daily": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"at": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
										},
									},
									"five_minute": dsschema.BoolAttribute{
										Description: "",
										Computed:    true,
									},
									"hourly": dsschema.BoolAttribute{
										Description: "",
										Computed:    true,
									},
									"monthly": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"at": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
											"day_of_month": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
										},
									},
									"weekly": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"at": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
											"day_of_week": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
										},
									},
								},
							},
							"url": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
						},
					},
					"predefined_ip": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"description": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"exception_list": dsschema.ListAttribute{
								Description: "",
								Computed:    true,
								ElementType: types.StringType,
							},
							"url": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
						},
					},
					"predefined_url": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"description": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"exception_list": dsschema.ListAttribute{
								Description: "",
								Computed:    true,
								ElementType: types.StringType,
							},
							"url": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
						},
					},
					"url": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"auth": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"password": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"username": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
							"certificate_profile": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"description": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"exception_list": dsschema.ListAttribute{
								Description: "",
								Computed:    true,
								ElementType: types.StringType,
							},
							"recurring": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"daily": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"at": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
										},
									},
									"five_minute": dsschema.BoolAttribute{
										Description: "",
										Computed:    true,
									},
									"hourly": dsschema.BoolAttribute{
										Description: "",
										Computed:    true,
									},
									"monthly": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"at": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
											"day_of_month": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
										},
									},
									"weekly": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"at": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
											"day_of_week": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
										},
									},
								},
							},
							"url": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
						},
					},
				},
			},
		},
	}
}

// Configure prepares the struct.
func (d *objectsExternalDynamicListsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *objectsExternalDynamicListsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state objectsExternalDynamicListsDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_objects_external_dynamic_lists",
		"object_id":        state.ObjectId.ValueString(),
	})

	// Prepare to run the command.
	svc := iHJqznH.NewClient(d.client)
	input := iHJqznH.ReadInput{
		ObjectId: state.ObjectId.ValueString(),
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
	state.Id = types.StringValue(idBuilder.String())
	var var0 objectsExternalDynamicListsDsModelTypeObject
	var var1 *objectsExternalDynamicListsDsModelDomainObject
	if ans.Type.Domain != nil {
		var1 = &objectsExternalDynamicListsDsModelDomainObject{}
		var var2 *objectsExternalDynamicListsDsModelAuthObject
		if ans.Type.Domain.Auth != nil {
			var2 = &objectsExternalDynamicListsDsModelAuthObject{}
			var2.Password = types.StringValue(ans.Type.Domain.Auth.Password)
			var2.Username = types.StringValue(ans.Type.Domain.Auth.Username)
		}
		var var3 objectsExternalDynamicListsDsModelRecurringObject
		var var4 *objectsExternalDynamicListsDsModelDailyObject
		if ans.Type.Domain.Recurring.Daily != nil {
			var4 = &objectsExternalDynamicListsDsModelDailyObject{}
			var4.At = types.StringValue(ans.Type.Domain.Recurring.Daily.At)
		}
		var var5 *objectsExternalDynamicListsDsModelMonthlyObject
		if ans.Type.Domain.Recurring.Monthly != nil {
			var5 = &objectsExternalDynamicListsDsModelMonthlyObject{}
			var5.At = types.StringValue(ans.Type.Domain.Recurring.Monthly.At)
			var5.DayOfMonth = types.Int64Value(ans.Type.Domain.Recurring.Monthly.DayOfMonth)
		}
		var var6 *objectsExternalDynamicListsDsModelWeeklyObject
		if ans.Type.Domain.Recurring.Weekly != nil {
			var6 = &objectsExternalDynamicListsDsModelWeeklyObject{}
			var6.At = types.StringValue(ans.Type.Domain.Recurring.Weekly.At)
			var6.DayOfWeek = types.StringValue(ans.Type.Domain.Recurring.Weekly.DayOfWeek)
		}
		var3.Daily = var4
		if ans.Type.Domain.Recurring.FiveMinute != nil {
			var3.FiveMinute = types.BoolValue(true)
		}
		if ans.Type.Domain.Recurring.Hourly != nil {
			var3.Hourly = types.BoolValue(true)
		}
		var3.Monthly = var5
		var3.Weekly = var6
		var1.Auth = var2
		var1.CertificateProfile = types.StringValue(ans.Type.Domain.CertificateProfile)
		var1.Description = types.StringValue(ans.Type.Domain.Description)
		var1.ExceptionList = EncodeStringSlice(ans.Type.Domain.ExceptionList)
		var1.ExpandDomain = types.BoolValue(ans.Type.Domain.ExpandDomain)
		var1.Recurring = var3
		var1.Url = types.StringValue(ans.Type.Domain.Url)
	}
	var var7 *objectsExternalDynamicListsDsModelImeiObject
	if ans.Type.Imei != nil {
		var7 = &objectsExternalDynamicListsDsModelImeiObject{}
		var var8 *objectsExternalDynamicListsDsModelAuthObject
		if ans.Type.Imei.Auth != nil {
			var8 = &objectsExternalDynamicListsDsModelAuthObject{}
			var8.Password = types.StringValue(ans.Type.Imei.Auth.Password)
			var8.Username = types.StringValue(ans.Type.Imei.Auth.Username)
		}
		var var9 objectsExternalDynamicListsDsModelRecurringObject
		var var10 *objectsExternalDynamicListsDsModelDailyObject
		if ans.Type.Imei.Recurring.Daily != nil {
			var10 = &objectsExternalDynamicListsDsModelDailyObject{}
			var10.At = types.StringValue(ans.Type.Imei.Recurring.Daily.At)
		}
		var var11 *objectsExternalDynamicListsDsModelMonthlyObject
		if ans.Type.Imei.Recurring.Monthly != nil {
			var11 = &objectsExternalDynamicListsDsModelMonthlyObject{}
			var11.At = types.StringValue(ans.Type.Imei.Recurring.Monthly.At)
			var11.DayOfMonth = types.Int64Value(ans.Type.Imei.Recurring.Monthly.DayOfMonth)
		}
		var var12 *objectsExternalDynamicListsDsModelWeeklyObject
		if ans.Type.Imei.Recurring.Weekly != nil {
			var12 = &objectsExternalDynamicListsDsModelWeeklyObject{}
			var12.At = types.StringValue(ans.Type.Imei.Recurring.Weekly.At)
			var12.DayOfWeek = types.StringValue(ans.Type.Imei.Recurring.Weekly.DayOfWeek)
		}
		var9.Daily = var10
		if ans.Type.Imei.Recurring.FiveMinute != nil {
			var9.FiveMinute = types.BoolValue(true)
		}
		if ans.Type.Imei.Recurring.Hourly != nil {
			var9.Hourly = types.BoolValue(true)
		}
		var9.Monthly = var11
		var9.Weekly = var12
		var7.Auth = var8
		var7.CertificateProfile = types.StringValue(ans.Type.Imei.CertificateProfile)
		var7.Description = types.StringValue(ans.Type.Imei.Description)
		var7.ExceptionList = EncodeStringSlice(ans.Type.Imei.ExceptionList)
		var7.Recurring = var9
		var7.Url = types.StringValue(ans.Type.Imei.Url)
	}
	var var13 *objectsExternalDynamicListsDsModelImsiObject
	if ans.Type.Imsi != nil {
		var13 = &objectsExternalDynamicListsDsModelImsiObject{}
		var var14 *objectsExternalDynamicListsDsModelAuthObject
		if ans.Type.Imsi.Auth != nil {
			var14 = &objectsExternalDynamicListsDsModelAuthObject{}
			var14.Password = types.StringValue(ans.Type.Imsi.Auth.Password)
			var14.Username = types.StringValue(ans.Type.Imsi.Auth.Username)
		}
		var var15 objectsExternalDynamicListsDsModelRecurringObject
		var var16 *objectsExternalDynamicListsDsModelDailyObject
		if ans.Type.Imsi.Recurring.Daily != nil {
			var16 = &objectsExternalDynamicListsDsModelDailyObject{}
			var16.At = types.StringValue(ans.Type.Imsi.Recurring.Daily.At)
		}
		var var17 *objectsExternalDynamicListsDsModelMonthlyObject
		if ans.Type.Imsi.Recurring.Monthly != nil {
			var17 = &objectsExternalDynamicListsDsModelMonthlyObject{}
			var17.At = types.StringValue(ans.Type.Imsi.Recurring.Monthly.At)
			var17.DayOfMonth = types.Int64Value(ans.Type.Imsi.Recurring.Monthly.DayOfMonth)
		}
		var var18 *objectsExternalDynamicListsDsModelWeeklyObject
		if ans.Type.Imsi.Recurring.Weekly != nil {
			var18 = &objectsExternalDynamicListsDsModelWeeklyObject{}
			var18.At = types.StringValue(ans.Type.Imsi.Recurring.Weekly.At)
			var18.DayOfWeek = types.StringValue(ans.Type.Imsi.Recurring.Weekly.DayOfWeek)
		}
		var15.Daily = var16
		if ans.Type.Imsi.Recurring.FiveMinute != nil {
			var15.FiveMinute = types.BoolValue(true)
		}
		if ans.Type.Imsi.Recurring.Hourly != nil {
			var15.Hourly = types.BoolValue(true)
		}
		var15.Monthly = var17
		var15.Weekly = var18
		var13.Auth = var14
		var13.CertificateProfile = types.StringValue(ans.Type.Imsi.CertificateProfile)
		var13.Description = types.StringValue(ans.Type.Imsi.Description)
		var13.ExceptionList = EncodeStringSlice(ans.Type.Imsi.ExceptionList)
		var13.Recurring = var15
		var13.Url = types.StringValue(ans.Type.Imsi.Url)
	}
	var var19 *objectsExternalDynamicListsDsModelIpObject
	if ans.Type.Ip != nil {
		var19 = &objectsExternalDynamicListsDsModelIpObject{}
		var var20 *objectsExternalDynamicListsDsModelAuthObject
		if ans.Type.Ip.Auth != nil {
			var20 = &objectsExternalDynamicListsDsModelAuthObject{}
			var20.Password = types.StringValue(ans.Type.Ip.Auth.Password)
			var20.Username = types.StringValue(ans.Type.Ip.Auth.Username)
		}
		var var21 objectsExternalDynamicListsDsModelRecurringObject
		var var22 *objectsExternalDynamicListsDsModelDailyObject
		if ans.Type.Ip.Recurring.Daily != nil {
			var22 = &objectsExternalDynamicListsDsModelDailyObject{}
			var22.At = types.StringValue(ans.Type.Ip.Recurring.Daily.At)
		}
		var var23 *objectsExternalDynamicListsDsModelMonthlyObject
		if ans.Type.Ip.Recurring.Monthly != nil {
			var23 = &objectsExternalDynamicListsDsModelMonthlyObject{}
			var23.At = types.StringValue(ans.Type.Ip.Recurring.Monthly.At)
			var23.DayOfMonth = types.Int64Value(ans.Type.Ip.Recurring.Monthly.DayOfMonth)
		}
		var var24 *objectsExternalDynamicListsDsModelWeeklyObject
		if ans.Type.Ip.Recurring.Weekly != nil {
			var24 = &objectsExternalDynamicListsDsModelWeeklyObject{}
			var24.At = types.StringValue(ans.Type.Ip.Recurring.Weekly.At)
			var24.DayOfWeek = types.StringValue(ans.Type.Ip.Recurring.Weekly.DayOfWeek)
		}
		var21.Daily = var22
		if ans.Type.Ip.Recurring.FiveMinute != nil {
			var21.FiveMinute = types.BoolValue(true)
		}
		if ans.Type.Ip.Recurring.Hourly != nil {
			var21.Hourly = types.BoolValue(true)
		}
		var21.Monthly = var23
		var21.Weekly = var24
		var19.Auth = var20
		var19.CertificateProfile = types.StringValue(ans.Type.Ip.CertificateProfile)
		var19.Description = types.StringValue(ans.Type.Ip.Description)
		var19.ExceptionList = EncodeStringSlice(ans.Type.Ip.ExceptionList)
		var19.Recurring = var21
		var19.Url = types.StringValue(ans.Type.Ip.Url)
	}
	var var25 *objectsExternalDynamicListsDsModelPredefinedIpObject
	if ans.Type.PredefinedIp != nil {
		var25 = &objectsExternalDynamicListsDsModelPredefinedIpObject{}
		var25.Description = types.StringValue(ans.Type.PredefinedIp.Description)
		var25.ExceptionList = EncodeStringSlice(ans.Type.PredefinedIp.ExceptionList)
		var25.Url = types.StringValue(ans.Type.PredefinedIp.Url)
	}
	var var26 *objectsExternalDynamicListsDsModelPredefinedUrlObject
	if ans.Type.PredefinedUrl != nil {
		var26 = &objectsExternalDynamicListsDsModelPredefinedUrlObject{}
		var26.Description = types.StringValue(ans.Type.PredefinedUrl.Description)
		var26.ExceptionList = EncodeStringSlice(ans.Type.PredefinedUrl.ExceptionList)
		var26.Url = types.StringValue(ans.Type.PredefinedUrl.Url)
	}
	var var27 *objectsExternalDynamicListsDsModelUrlObject
	if ans.Type.Url != nil {
		var27 = &objectsExternalDynamicListsDsModelUrlObject{}
		var var28 *objectsExternalDynamicListsDsModelAuthObject
		if ans.Type.Url.Auth != nil {
			var28 = &objectsExternalDynamicListsDsModelAuthObject{}
			var28.Password = types.StringValue(ans.Type.Url.Auth.Password)
			var28.Username = types.StringValue(ans.Type.Url.Auth.Username)
		}
		var var29 objectsExternalDynamicListsDsModelRecurringObject
		var var30 *objectsExternalDynamicListsDsModelDailyObject
		if ans.Type.Url.Recurring.Daily != nil {
			var30 = &objectsExternalDynamicListsDsModelDailyObject{}
			var30.At = types.StringValue(ans.Type.Url.Recurring.Daily.At)
		}
		var var31 *objectsExternalDynamicListsDsModelMonthlyObject
		if ans.Type.Url.Recurring.Monthly != nil {
			var31 = &objectsExternalDynamicListsDsModelMonthlyObject{}
			var31.At = types.StringValue(ans.Type.Url.Recurring.Monthly.At)
			var31.DayOfMonth = types.Int64Value(ans.Type.Url.Recurring.Monthly.DayOfMonth)
		}
		var var32 *objectsExternalDynamicListsDsModelWeeklyObject
		if ans.Type.Url.Recurring.Weekly != nil {
			var32 = &objectsExternalDynamicListsDsModelWeeklyObject{}
			var32.At = types.StringValue(ans.Type.Url.Recurring.Weekly.At)
			var32.DayOfWeek = types.StringValue(ans.Type.Url.Recurring.Weekly.DayOfWeek)
		}
		var29.Daily = var30
		if ans.Type.Url.Recurring.FiveMinute != nil {
			var29.FiveMinute = types.BoolValue(true)
		}
		if ans.Type.Url.Recurring.Hourly != nil {
			var29.Hourly = types.BoolValue(true)
		}
		var29.Monthly = var31
		var29.Weekly = var32
		var27.Auth = var28
		var27.CertificateProfile = types.StringValue(ans.Type.Url.CertificateProfile)
		var27.Description = types.StringValue(ans.Type.Url.Description)
		var27.ExceptionList = EncodeStringSlice(ans.Type.Url.ExceptionList)
		var27.Recurring = var29
		var27.Url = types.StringValue(ans.Type.Url.Url)
	}
	var0.Domain = var1
	var0.Imei = var7
	var0.Imsi = var13
	var0.Ip = var19
	var0.PredefinedIp = var25
	var0.PredefinedUrl = var26
	var0.Url = var27
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.Type = var0

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &objectsExternalDynamicListsResource{}
	_ resource.ResourceWithConfigure   = &objectsExternalDynamicListsResource{}
	_ resource.ResourceWithImportState = &objectsExternalDynamicListsResource{}
)

func NewObjectsExternalDynamicListsResource() resource.Resource {
	return &objectsExternalDynamicListsResource{}
}

type objectsExternalDynamicListsResource struct {
	client *sase.Client
}

type objectsExternalDynamicListsRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/objects-external-dynamic-lists
	ObjectId types.String                                 `tfsdk:"object_id"`
	Name     types.String                                 `tfsdk:"name"`
	Type     objectsExternalDynamicListsRsModelTypeObject `tfsdk:"type"`
}

type objectsExternalDynamicListsRsModelTypeObject struct {
	Domain        *objectsExternalDynamicListsRsModelDomainObject        `tfsdk:"domain"`
	Imei          *objectsExternalDynamicListsRsModelImeiObject          `tfsdk:"imei"`
	Imsi          *objectsExternalDynamicListsRsModelImsiObject          `tfsdk:"imsi"`
	Ip            *objectsExternalDynamicListsRsModelIpObject            `tfsdk:"ip"`
	PredefinedIp  *objectsExternalDynamicListsRsModelPredefinedIpObject  `tfsdk:"predefined_ip"`
	PredefinedUrl *objectsExternalDynamicListsRsModelPredefinedUrlObject `tfsdk:"predefined_url"`
	Url           *objectsExternalDynamicListsRsModelUrlObject           `tfsdk:"url"`
}

type objectsExternalDynamicListsRsModelDomainObject struct {
	Auth               *objectsExternalDynamicListsRsModelAuthObject     `tfsdk:"auth"`
	CertificateProfile types.String                                      `tfsdk:"certificate_profile"`
	Description        types.String                                      `tfsdk:"description"`
	ExceptionList      []types.String                                    `tfsdk:"exception_list"`
	ExpandDomain       types.Bool                                        `tfsdk:"expand_domain"`
	Recurring          objectsExternalDynamicListsRsModelRecurringObject `tfsdk:"recurring"`
	Url                types.String                                      `tfsdk:"url"`
}

type objectsExternalDynamicListsRsModelAuthObject struct {
	Password types.String `tfsdk:"password"`
	Username types.String `tfsdk:"username"`
}

type objectsExternalDynamicListsRsModelRecurringObject struct {
	Daily      *objectsExternalDynamicListsRsModelDailyObject   `tfsdk:"daily"`
	FiveMinute types.Bool                                       `tfsdk:"five_minute"`
	Hourly     types.Bool                                       `tfsdk:"hourly"`
	Monthly    *objectsExternalDynamicListsRsModelMonthlyObject `tfsdk:"monthly"`
	Weekly     *objectsExternalDynamicListsRsModelWeeklyObject  `tfsdk:"weekly"`
}

type objectsExternalDynamicListsRsModelDailyObject struct {
	At types.String `tfsdk:"at"`
}

type objectsExternalDynamicListsRsModelMonthlyObject struct {
	At         types.String `tfsdk:"at"`
	DayOfMonth types.Int64  `tfsdk:"day_of_month"`
}

type objectsExternalDynamicListsRsModelWeeklyObject struct {
	At        types.String `tfsdk:"at"`
	DayOfWeek types.String `tfsdk:"day_of_week"`
}

type objectsExternalDynamicListsRsModelImeiObject struct {
	Auth               *objectsExternalDynamicListsRsModelAuthObject     `tfsdk:"auth"`
	CertificateProfile types.String                                      `tfsdk:"certificate_profile"`
	Description        types.String                                      `tfsdk:"description"`
	ExceptionList      []types.String                                    `tfsdk:"exception_list"`
	Recurring          objectsExternalDynamicListsRsModelRecurringObject `tfsdk:"recurring"`
	Url                types.String                                      `tfsdk:"url"`
}

type objectsExternalDynamicListsRsModelImsiObject struct {
	Auth               *objectsExternalDynamicListsRsModelAuthObject     `tfsdk:"auth"`
	CertificateProfile types.String                                      `tfsdk:"certificate_profile"`
	Description        types.String                                      `tfsdk:"description"`
	ExceptionList      []types.String                                    `tfsdk:"exception_list"`
	Recurring          objectsExternalDynamicListsRsModelRecurringObject `tfsdk:"recurring"`
	Url                types.String                                      `tfsdk:"url"`
}

type objectsExternalDynamicListsRsModelIpObject struct {
	Auth               *objectsExternalDynamicListsRsModelAuthObject     `tfsdk:"auth"`
	CertificateProfile types.String                                      `tfsdk:"certificate_profile"`
	Description        types.String                                      `tfsdk:"description"`
	ExceptionList      []types.String                                    `tfsdk:"exception_list"`
	Recurring          objectsExternalDynamicListsRsModelRecurringObject `tfsdk:"recurring"`
	Url                types.String                                      `tfsdk:"url"`
}

type objectsExternalDynamicListsRsModelPredefinedIpObject struct {
	Description   types.String   `tfsdk:"description"`
	ExceptionList []types.String `tfsdk:"exception_list"`
	Url           types.String   `tfsdk:"url"`
}

type objectsExternalDynamicListsRsModelPredefinedUrlObject struct {
	Description   types.String   `tfsdk:"description"`
	ExceptionList []types.String `tfsdk:"exception_list"`
	Url           types.String   `tfsdk:"url"`
}

type objectsExternalDynamicListsRsModelUrlObject struct {
	Auth               *objectsExternalDynamicListsRsModelAuthObject     `tfsdk:"auth"`
	CertificateProfile types.String                                      `tfsdk:"certificate_profile"`
	Description        types.String                                      `tfsdk:"description"`
	ExceptionList      []types.String                                    `tfsdk:"exception_list"`
	Recurring          objectsExternalDynamicListsRsModelRecurringObject `tfsdk:"recurring"`
	Url                types.String                                      `tfsdk:"url"`
}

// Metadata returns the data source type name.
func (r *objectsExternalDynamicListsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_external_dynamic_lists"
}

// Schema defines the schema for this listing data source.
func (r *objectsExternalDynamicListsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Service Connections", "Mobile Users Container", "Mobile Users Explicit Proxy"),
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
				Validators: []validator.String{
					stringvalidator.LengthAtMost(63),
				},
			},
			"type": rsschema.SingleNestedAttribute{
				Description: "",
				Required:    true,
				Attributes: map[string]rsschema.Attribute{
					"domain": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"auth": rsschema.SingleNestedAttribute{
								Description: "",
								Optional:    true,
								Attributes: map[string]rsschema.Attribute{
									"password": rsschema.StringAttribute{
										Description: "",
										Required:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthAtMost(255),
										},
									},
									"username": rsschema.StringAttribute{
										Description: "",
										Required:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(1, 255),
										},
									},
								},
							},
							"certificate_profile": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString("None"),
								},
							},
							"description": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 255),
								},
							},
							"exception_list": rsschema.ListAttribute{
								Description: "",
								Optional:    true,
								ElementType: types.StringType,
							},
							"expand_domain": rsschema.BoolAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Bool{
									DefaultBool(false),
								},
							},
							"recurring": rsschema.SingleNestedAttribute{
								Description: "",
								Required:    true,
								Attributes: map[string]rsschema.Attribute{
									"daily": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"at": rsschema.StringAttribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString("00"),
												},
												Validators: []validator.String{
													stringvalidator.LengthBetween(2, 2),
												},
											},
										},
									},
									"five_minute": rsschema.BoolAttribute{
										Description: "",
										Optional:    true,
										Validators: []validator.Bool{
											boolvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("daily"),
												path.MatchRelative().AtParent().AtName("hourly"),
												path.MatchRelative().AtParent().AtName("monthly"),
												path.MatchRelative().AtParent().AtName("weekly"),
											),
										},
									},
									"hourly": rsschema.BoolAttribute{
										Description: "",
										Optional:    true,
										Validators: []validator.Bool{
											boolvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("daily"),
												path.MatchRelative().AtParent().AtName("five_minute"),
												path.MatchRelative().AtParent().AtName("monthly"),
												path.MatchRelative().AtParent().AtName("weekly"),
											),
										},
									},
									"monthly": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"at": rsschema.StringAttribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString("00"),
												},
												Validators: []validator.String{
													stringvalidator.LengthBetween(2, 2),
												},
											},
											"day_of_month": rsschema.Int64Attribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(0),
												},
												Validators: []validator.Int64{
													int64validator.Between(1, 31),
												},
											},
										},
									},
									"weekly": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"at": rsschema.StringAttribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString("00"),
												},
												Validators: []validator.String{
													stringvalidator.LengthBetween(2, 2),
												},
											},
											"day_of_week": rsschema.StringAttribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString(""),
												},
												Validators: []validator.String{
													stringvalidator.OneOf("sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"),
												},
											},
										},
									},
								},
							},
							"url": rsschema.StringAttribute{
								Description: "",
								Required:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString("http://"),
								},
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 255),
								},
							},
						},
					},
					"imei": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"auth": rsschema.SingleNestedAttribute{
								Description: "",
								Optional:    true,
								Attributes: map[string]rsschema.Attribute{
									"password": rsschema.StringAttribute{
										Description: "",
										Required:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthAtMost(255),
										},
									},
									"username": rsschema.StringAttribute{
										Description: "",
										Required:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(1, 255),
										},
									},
								},
							},
							"certificate_profile": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString("None"),
								},
							},
							"description": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 255),
								},
							},
							"exception_list": rsschema.ListAttribute{
								Description: "",
								Optional:    true,
								ElementType: types.StringType,
							},
							"recurring": rsschema.SingleNestedAttribute{
								Description: "",
								Required:    true,
								Attributes: map[string]rsschema.Attribute{
									"daily": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"at": rsschema.StringAttribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString("00"),
												},
												Validators: []validator.String{
													stringvalidator.LengthBetween(2, 2),
												},
											},
										},
									},
									"five_minute": rsschema.BoolAttribute{
										Description: "",
										Optional:    true,
										Validators: []validator.Bool{
											boolvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("daily"),
												path.MatchRelative().AtParent().AtName("hourly"),
												path.MatchRelative().AtParent().AtName("monthly"),
												path.MatchRelative().AtParent().AtName("weekly"),
											),
										},
									},
									"hourly": rsschema.BoolAttribute{
										Description: "",
										Optional:    true,
										Validators: []validator.Bool{
											boolvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("daily"),
												path.MatchRelative().AtParent().AtName("five_minute"),
												path.MatchRelative().AtParent().AtName("monthly"),
												path.MatchRelative().AtParent().AtName("weekly"),
											),
										},
									},
									"monthly": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"at": rsschema.StringAttribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString("00"),
												},
												Validators: []validator.String{
													stringvalidator.LengthBetween(2, 2),
												},
											},
											"day_of_month": rsschema.Int64Attribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(0),
												},
												Validators: []validator.Int64{
													int64validator.Between(1, 31),
												},
											},
										},
									},
									"weekly": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"at": rsschema.StringAttribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString("00"),
												},
												Validators: []validator.String{
													stringvalidator.LengthBetween(2, 2),
												},
											},
											"day_of_week": rsschema.StringAttribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString(""),
												},
												Validators: []validator.String{
													stringvalidator.OneOf("sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"),
												},
											},
										},
									},
								},
							},
							"url": rsschema.StringAttribute{
								Description: "",
								Required:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString("http://"),
								},
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 255),
								},
							},
						},
					},
					"imsi": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"auth": rsschema.SingleNestedAttribute{
								Description: "",
								Optional:    true,
								Attributes: map[string]rsschema.Attribute{
									"password": rsschema.StringAttribute{
										Description: "",
										Required:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthAtMost(255),
										},
									},
									"username": rsschema.StringAttribute{
										Description: "",
										Required:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(1, 255),
										},
									},
								},
							},
							"certificate_profile": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString("None"),
								},
							},
							"description": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 255),
								},
							},
							"exception_list": rsschema.ListAttribute{
								Description: "",
								Optional:    true,
								ElementType: types.StringType,
							},
							"recurring": rsschema.SingleNestedAttribute{
								Description: "",
								Required:    true,
								Attributes: map[string]rsschema.Attribute{
									"daily": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"at": rsschema.StringAttribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString("00"),
												},
												Validators: []validator.String{
													stringvalidator.LengthBetween(2, 2),
												},
											},
										},
									},
									"five_minute": rsschema.BoolAttribute{
										Description: "",
										Optional:    true,
										Validators: []validator.Bool{
											boolvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("daily"),
												path.MatchRelative().AtParent().AtName("hourly"),
												path.MatchRelative().AtParent().AtName("monthly"),
												path.MatchRelative().AtParent().AtName("weekly"),
											),
										},
									},
									"hourly": rsschema.BoolAttribute{
										Description: "",
										Optional:    true,
										Validators: []validator.Bool{
											boolvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("daily"),
												path.MatchRelative().AtParent().AtName("five_minute"),
												path.MatchRelative().AtParent().AtName("monthly"),
												path.MatchRelative().AtParent().AtName("weekly"),
											),
										},
									},
									"monthly": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"at": rsschema.StringAttribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString("00"),
												},
												Validators: []validator.String{
													stringvalidator.LengthBetween(2, 2),
												},
											},
											"day_of_month": rsschema.Int64Attribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(0),
												},
												Validators: []validator.Int64{
													int64validator.Between(1, 31),
												},
											},
										},
									},
									"weekly": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"at": rsschema.StringAttribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString("00"),
												},
												Validators: []validator.String{
													stringvalidator.LengthBetween(2, 2),
												},
											},
											"day_of_week": rsschema.StringAttribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString(""),
												},
												Validators: []validator.String{
													stringvalidator.OneOf("sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"),
												},
											},
										},
									},
								},
							},
							"url": rsschema.StringAttribute{
								Description: "",
								Required:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString("http://"),
								},
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 255),
								},
							},
						},
					},
					"ip": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"auth": rsschema.SingleNestedAttribute{
								Description: "",
								Optional:    true,
								Attributes: map[string]rsschema.Attribute{
									"password": rsschema.StringAttribute{
										Description: "",
										Required:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthAtMost(255),
										},
									},
									"username": rsschema.StringAttribute{
										Description: "",
										Required:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(1, 255),
										},
									},
								},
							},
							"certificate_profile": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString("None"),
								},
							},
							"description": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 255),
								},
							},
							"exception_list": rsschema.ListAttribute{
								Description: "",
								Optional:    true,
								ElementType: types.StringType,
							},
							"recurring": rsschema.SingleNestedAttribute{
								Description: "",
								Required:    true,
								Attributes: map[string]rsschema.Attribute{
									"daily": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"at": rsschema.StringAttribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString("00"),
												},
												Validators: []validator.String{
													stringvalidator.LengthBetween(2, 2),
												},
											},
										},
									},
									"five_minute": rsschema.BoolAttribute{
										Description: "",
										Optional:    true,
										Validators: []validator.Bool{
											boolvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("daily"),
												path.MatchRelative().AtParent().AtName("hourly"),
												path.MatchRelative().AtParent().AtName("monthly"),
												path.MatchRelative().AtParent().AtName("weekly"),
											),
										},
									},
									"hourly": rsschema.BoolAttribute{
										Description: "",
										Optional:    true,
										Validators: []validator.Bool{
											boolvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("daily"),
												path.MatchRelative().AtParent().AtName("five_minute"),
												path.MatchRelative().AtParent().AtName("monthly"),
												path.MatchRelative().AtParent().AtName("weekly"),
											),
										},
									},
									"monthly": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"at": rsschema.StringAttribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString("00"),
												},
												Validators: []validator.String{
													stringvalidator.LengthBetween(2, 2),
												},
											},
											"day_of_month": rsschema.Int64Attribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(0),
												},
												Validators: []validator.Int64{
													int64validator.Between(1, 31),
												},
											},
										},
									},
									"weekly": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"at": rsschema.StringAttribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString("00"),
												},
												Validators: []validator.String{
													stringvalidator.LengthBetween(2, 2),
												},
											},
											"day_of_week": rsschema.StringAttribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString(""),
												},
												Validators: []validator.String{
													stringvalidator.OneOf("sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"),
												},
											},
										},
									},
								},
							},
							"url": rsschema.StringAttribute{
								Description: "",
								Required:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString("http://"),
								},
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 255),
								},
							},
						},
					},
					"predefined_ip": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"description": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 255),
								},
							},
							"exception_list": rsschema.ListAttribute{
								Description: "",
								Optional:    true,
								ElementType: types.StringType,
							},
							"url": rsschema.StringAttribute{
								Description: "",
								Required:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
						},
					},
					"predefined_url": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"description": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 255),
								},
							},
							"exception_list": rsschema.ListAttribute{
								Description: "",
								Optional:    true,
								ElementType: types.StringType,
							},
							"url": rsschema.StringAttribute{
								Description: "",
								Required:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
						},
					},
					"url": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"auth": rsschema.SingleNestedAttribute{
								Description: "",
								Optional:    true,
								Attributes: map[string]rsschema.Attribute{
									"password": rsschema.StringAttribute{
										Description: "",
										Required:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthAtMost(255),
										},
									},
									"username": rsschema.StringAttribute{
										Description: "",
										Required:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(1, 255),
										},
									},
								},
							},
							"certificate_profile": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString("None"),
								},
							},
							"description": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 255),
								},
							},
							"exception_list": rsschema.ListAttribute{
								Description: "",
								Optional:    true,
								ElementType: types.StringType,
							},
							"recurring": rsschema.SingleNestedAttribute{
								Description: "",
								Required:    true,
								Attributes: map[string]rsschema.Attribute{
									"daily": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"at": rsschema.StringAttribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString("00"),
												},
												Validators: []validator.String{
													stringvalidator.LengthBetween(2, 2),
												},
											},
										},
									},
									"five_minute": rsschema.BoolAttribute{
										Description: "",
										Optional:    true,
										Validators: []validator.Bool{
											boolvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("daily"),
												path.MatchRelative().AtParent().AtName("hourly"),
												path.MatchRelative().AtParent().AtName("monthly"),
												path.MatchRelative().AtParent().AtName("weekly"),
											),
										},
									},
									"hourly": rsschema.BoolAttribute{
										Description: "",
										Optional:    true,
										Validators: []validator.Bool{
											boolvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("daily"),
												path.MatchRelative().AtParent().AtName("five_minute"),
												path.MatchRelative().AtParent().AtName("monthly"),
												path.MatchRelative().AtParent().AtName("weekly"),
											),
										},
									},
									"monthly": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"at": rsschema.StringAttribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString("00"),
												},
												Validators: []validator.String{
													stringvalidator.LengthBetween(2, 2),
												},
											},
											"day_of_month": rsschema.Int64Attribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(0),
												},
												Validators: []validator.Int64{
													int64validator.Between(1, 31),
												},
											},
										},
									},
									"weekly": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"at": rsschema.StringAttribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString("00"),
												},
												Validators: []validator.String{
													stringvalidator.LengthBetween(2, 2),
												},
											},
											"day_of_week": rsschema.StringAttribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString(""),
												},
												Validators: []validator.String{
													stringvalidator.OneOf("sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"),
												},
											},
										},
									},
								},
							},
							"url": rsschema.StringAttribute{
								Description: "",
								Required:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString("http://"),
								},
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 255),
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *objectsExternalDynamicListsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *objectsExternalDynamicListsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state objectsExternalDynamicListsRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "sase_objects_external_dynamic_lists",
		"folder":        state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := iHJqznH.NewClient(r.client)
	input := iHJqznH.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 sRAOviP.Config
	var0.Name = state.Name.ValueString()
	var var1 sRAOviP.TypeObject
	var var2 *sRAOviP.DomainObject
	if state.Type.Domain != nil {
		var2 = &sRAOviP.DomainObject{}
		var var3 *sRAOviP.AuthObject
		if state.Type.Domain.Auth != nil {
			var3 = &sRAOviP.AuthObject{}
			var3.Password = state.Type.Domain.Auth.Password.ValueString()
			var3.Username = state.Type.Domain.Auth.Username.ValueString()
		}
		var2.Auth = var3
		var2.CertificateProfile = state.Type.Domain.CertificateProfile.ValueString()
		var2.Description = state.Type.Domain.Description.ValueString()
		var2.ExceptionList = DecodeStringSlice(state.Type.Domain.ExceptionList)
		var2.ExpandDomain = state.Type.Domain.ExpandDomain.ValueBool()
		var var4 sRAOviP.RecurringObject
		var var5 *sRAOviP.DailyObject
		if state.Type.Domain.Recurring.Daily != nil {
			var5 = &sRAOviP.DailyObject{}
			var5.At = state.Type.Domain.Recurring.Daily.At.ValueString()
		}
		var4.Daily = var5
		if state.Type.Domain.Recurring.FiveMinute.ValueBool() {
			var4.FiveMinute = struct{}{}
		}
		if state.Type.Domain.Recurring.Hourly.ValueBool() {
			var4.Hourly = struct{}{}
		}
		var var6 *sRAOviP.MonthlyObject
		if state.Type.Domain.Recurring.Monthly != nil {
			var6 = &sRAOviP.MonthlyObject{}
			var6.At = state.Type.Domain.Recurring.Monthly.At.ValueString()
			var6.DayOfMonth = state.Type.Domain.Recurring.Monthly.DayOfMonth.ValueInt64()
		}
		var4.Monthly = var6
		var var7 *sRAOviP.WeeklyObject
		if state.Type.Domain.Recurring.Weekly != nil {
			var7 = &sRAOviP.WeeklyObject{}
			var7.At = state.Type.Domain.Recurring.Weekly.At.ValueString()
			var7.DayOfWeek = state.Type.Domain.Recurring.Weekly.DayOfWeek.ValueString()
		}
		var4.Weekly = var7
		var2.Recurring = var4
		var2.Url = state.Type.Domain.Url.ValueString()
	}
	var1.Domain = var2
	var var8 *sRAOviP.ImeiObject
	if state.Type.Imei != nil {
		var8 = &sRAOviP.ImeiObject{}
		var var9 *sRAOviP.AuthObject
		if state.Type.Imei.Auth != nil {
			var9 = &sRAOviP.AuthObject{}
			var9.Password = state.Type.Imei.Auth.Password.ValueString()
			var9.Username = state.Type.Imei.Auth.Username.ValueString()
		}
		var8.Auth = var9
		var8.CertificateProfile = state.Type.Imei.CertificateProfile.ValueString()
		var8.Description = state.Type.Imei.Description.ValueString()
		var8.ExceptionList = DecodeStringSlice(state.Type.Imei.ExceptionList)
		var var10 sRAOviP.RecurringObject
		var var11 *sRAOviP.DailyObject
		if state.Type.Imei.Recurring.Daily != nil {
			var11 = &sRAOviP.DailyObject{}
			var11.At = state.Type.Imei.Recurring.Daily.At.ValueString()
		}
		var10.Daily = var11
		if state.Type.Imei.Recurring.FiveMinute.ValueBool() {
			var10.FiveMinute = struct{}{}
		}
		if state.Type.Imei.Recurring.Hourly.ValueBool() {
			var10.Hourly = struct{}{}
		}
		var var12 *sRAOviP.MonthlyObject
		if state.Type.Imei.Recurring.Monthly != nil {
			var12 = &sRAOviP.MonthlyObject{}
			var12.At = state.Type.Imei.Recurring.Monthly.At.ValueString()
			var12.DayOfMonth = state.Type.Imei.Recurring.Monthly.DayOfMonth.ValueInt64()
		}
		var10.Monthly = var12
		var var13 *sRAOviP.WeeklyObject
		if state.Type.Imei.Recurring.Weekly != nil {
			var13 = &sRAOviP.WeeklyObject{}
			var13.At = state.Type.Imei.Recurring.Weekly.At.ValueString()
			var13.DayOfWeek = state.Type.Imei.Recurring.Weekly.DayOfWeek.ValueString()
		}
		var10.Weekly = var13
		var8.Recurring = var10
		var8.Url = state.Type.Imei.Url.ValueString()
	}
	var1.Imei = var8
	var var14 *sRAOviP.ImsiObject
	if state.Type.Imsi != nil {
		var14 = &sRAOviP.ImsiObject{}
		var var15 *sRAOviP.AuthObject
		if state.Type.Imsi.Auth != nil {
			var15 = &sRAOviP.AuthObject{}
			var15.Password = state.Type.Imsi.Auth.Password.ValueString()
			var15.Username = state.Type.Imsi.Auth.Username.ValueString()
		}
		var14.Auth = var15
		var14.CertificateProfile = state.Type.Imsi.CertificateProfile.ValueString()
		var14.Description = state.Type.Imsi.Description.ValueString()
		var14.ExceptionList = DecodeStringSlice(state.Type.Imsi.ExceptionList)
		var var16 sRAOviP.RecurringObject
		var var17 *sRAOviP.DailyObject
		if state.Type.Imsi.Recurring.Daily != nil {
			var17 = &sRAOviP.DailyObject{}
			var17.At = state.Type.Imsi.Recurring.Daily.At.ValueString()
		}
		var16.Daily = var17
		if state.Type.Imsi.Recurring.FiveMinute.ValueBool() {
			var16.FiveMinute = struct{}{}
		}
		if state.Type.Imsi.Recurring.Hourly.ValueBool() {
			var16.Hourly = struct{}{}
		}
		var var18 *sRAOviP.MonthlyObject
		if state.Type.Imsi.Recurring.Monthly != nil {
			var18 = &sRAOviP.MonthlyObject{}
			var18.At = state.Type.Imsi.Recurring.Monthly.At.ValueString()
			var18.DayOfMonth = state.Type.Imsi.Recurring.Monthly.DayOfMonth.ValueInt64()
		}
		var16.Monthly = var18
		var var19 *sRAOviP.WeeklyObject
		if state.Type.Imsi.Recurring.Weekly != nil {
			var19 = &sRAOviP.WeeklyObject{}
			var19.At = state.Type.Imsi.Recurring.Weekly.At.ValueString()
			var19.DayOfWeek = state.Type.Imsi.Recurring.Weekly.DayOfWeek.ValueString()
		}
		var16.Weekly = var19
		var14.Recurring = var16
		var14.Url = state.Type.Imsi.Url.ValueString()
	}
	var1.Imsi = var14
	var var20 *sRAOviP.IpObject
	if state.Type.Ip != nil {
		var20 = &sRAOviP.IpObject{}
		var var21 *sRAOviP.AuthObject
		if state.Type.Ip.Auth != nil {
			var21 = &sRAOviP.AuthObject{}
			var21.Password = state.Type.Ip.Auth.Password.ValueString()
			var21.Username = state.Type.Ip.Auth.Username.ValueString()
		}
		var20.Auth = var21
		var20.CertificateProfile = state.Type.Ip.CertificateProfile.ValueString()
		var20.Description = state.Type.Ip.Description.ValueString()
		var20.ExceptionList = DecodeStringSlice(state.Type.Ip.ExceptionList)
		var var22 sRAOviP.RecurringObject
		var var23 *sRAOviP.DailyObject
		if state.Type.Ip.Recurring.Daily != nil {
			var23 = &sRAOviP.DailyObject{}
			var23.At = state.Type.Ip.Recurring.Daily.At.ValueString()
		}
		var22.Daily = var23
		if state.Type.Ip.Recurring.FiveMinute.ValueBool() {
			var22.FiveMinute = struct{}{}
		}
		if state.Type.Ip.Recurring.Hourly.ValueBool() {
			var22.Hourly = struct{}{}
		}
		var var24 *sRAOviP.MonthlyObject
		if state.Type.Ip.Recurring.Monthly != nil {
			var24 = &sRAOviP.MonthlyObject{}
			var24.At = state.Type.Ip.Recurring.Monthly.At.ValueString()
			var24.DayOfMonth = state.Type.Ip.Recurring.Monthly.DayOfMonth.ValueInt64()
		}
		var22.Monthly = var24
		var var25 *sRAOviP.WeeklyObject
		if state.Type.Ip.Recurring.Weekly != nil {
			var25 = &sRAOviP.WeeklyObject{}
			var25.At = state.Type.Ip.Recurring.Weekly.At.ValueString()
			var25.DayOfWeek = state.Type.Ip.Recurring.Weekly.DayOfWeek.ValueString()
		}
		var22.Weekly = var25
		var20.Recurring = var22
		var20.Url = state.Type.Ip.Url.ValueString()
	}
	var1.Ip = var20
	var var26 *sRAOviP.PredefinedIpObject
	if state.Type.PredefinedIp != nil {
		var26 = &sRAOviP.PredefinedIpObject{}
		var26.Description = state.Type.PredefinedIp.Description.ValueString()
		var26.ExceptionList = DecodeStringSlice(state.Type.PredefinedIp.ExceptionList)
		var26.Url = state.Type.PredefinedIp.Url.ValueString()
	}
	var1.PredefinedIp = var26
	var var27 *sRAOviP.PredefinedUrlObject
	if state.Type.PredefinedUrl != nil {
		var27 = &sRAOviP.PredefinedUrlObject{}
		var27.Description = state.Type.PredefinedUrl.Description.ValueString()
		var27.ExceptionList = DecodeStringSlice(state.Type.PredefinedUrl.ExceptionList)
		var27.Url = state.Type.PredefinedUrl.Url.ValueString()
	}
	var1.PredefinedUrl = var27
	var var28 *sRAOviP.UrlObject
	if state.Type.Url != nil {
		var28 = &sRAOviP.UrlObject{}
		var var29 *sRAOviP.AuthObject
		if state.Type.Url.Auth != nil {
			var29 = &sRAOviP.AuthObject{}
			var29.Password = state.Type.Url.Auth.Password.ValueString()
			var29.Username = state.Type.Url.Auth.Username.ValueString()
		}
		var28.Auth = var29
		var28.CertificateProfile = state.Type.Url.CertificateProfile.ValueString()
		var28.Description = state.Type.Url.Description.ValueString()
		var28.ExceptionList = DecodeStringSlice(state.Type.Url.ExceptionList)
		var var30 sRAOviP.RecurringObject
		var var31 *sRAOviP.DailyObject
		if state.Type.Url.Recurring.Daily != nil {
			var31 = &sRAOviP.DailyObject{}
			var31.At = state.Type.Url.Recurring.Daily.At.ValueString()
		}
		var30.Daily = var31
		if state.Type.Url.Recurring.FiveMinute.ValueBool() {
			var30.FiveMinute = struct{}{}
		}
		if state.Type.Url.Recurring.Hourly.ValueBool() {
			var30.Hourly = struct{}{}
		}
		var var32 *sRAOviP.MonthlyObject
		if state.Type.Url.Recurring.Monthly != nil {
			var32 = &sRAOviP.MonthlyObject{}
			var32.At = state.Type.Url.Recurring.Monthly.At.ValueString()
			var32.DayOfMonth = state.Type.Url.Recurring.Monthly.DayOfMonth.ValueInt64()
		}
		var30.Monthly = var32
		var var33 *sRAOviP.WeeklyObject
		if state.Type.Url.Recurring.Weekly != nil {
			var33 = &sRAOviP.WeeklyObject{}
			var33.At = state.Type.Url.Recurring.Weekly.At.ValueString()
			var33.DayOfWeek = state.Type.Url.Recurring.Weekly.DayOfWeek.ValueString()
		}
		var30.Weekly = var33
		var28.Recurring = var30
		var28.Url = state.Type.Url.Url.ValueString()
	}
	var1.Url = var28
	var0.Type = var1
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
	var var34 objectsExternalDynamicListsRsModelTypeObject
	var var35 *objectsExternalDynamicListsRsModelDomainObject
	if ans.Type.Domain != nil {
		var35 = &objectsExternalDynamicListsRsModelDomainObject{}
		var var36 *objectsExternalDynamicListsRsModelAuthObject
		if ans.Type.Domain.Auth != nil {
			var36 = &objectsExternalDynamicListsRsModelAuthObject{}
			var36.Password = types.StringValue(ans.Type.Domain.Auth.Password)
			var36.Username = types.StringValue(ans.Type.Domain.Auth.Username)
		}
		var var37 objectsExternalDynamicListsRsModelRecurringObject
		var var38 *objectsExternalDynamicListsRsModelDailyObject
		if ans.Type.Domain.Recurring.Daily != nil {
			var38 = &objectsExternalDynamicListsRsModelDailyObject{}
			var38.At = types.StringValue(ans.Type.Domain.Recurring.Daily.At)
		}
		var var39 *objectsExternalDynamicListsRsModelMonthlyObject
		if ans.Type.Domain.Recurring.Monthly != nil {
			var39 = &objectsExternalDynamicListsRsModelMonthlyObject{}
			var39.At = types.StringValue(ans.Type.Domain.Recurring.Monthly.At)
			var39.DayOfMonth = types.Int64Value(ans.Type.Domain.Recurring.Monthly.DayOfMonth)
		}
		var var40 *objectsExternalDynamicListsRsModelWeeklyObject
		if ans.Type.Domain.Recurring.Weekly != nil {
			var40 = &objectsExternalDynamicListsRsModelWeeklyObject{}
			var40.At = types.StringValue(ans.Type.Domain.Recurring.Weekly.At)
			var40.DayOfWeek = types.StringValue(ans.Type.Domain.Recurring.Weekly.DayOfWeek)
		}
		var37.Daily = var38
		if ans.Type.Domain.Recurring.FiveMinute != nil {
			var37.FiveMinute = types.BoolValue(true)
		}
		if ans.Type.Domain.Recurring.Hourly != nil {
			var37.Hourly = types.BoolValue(true)
		}
		var37.Monthly = var39
		var37.Weekly = var40
		var35.Auth = var36
		var35.CertificateProfile = types.StringValue(ans.Type.Domain.CertificateProfile)
		var35.Description = types.StringValue(ans.Type.Domain.Description)
		var35.ExceptionList = EncodeStringSlice(ans.Type.Domain.ExceptionList)
		var35.ExpandDomain = types.BoolValue(ans.Type.Domain.ExpandDomain)
		var35.Recurring = var37
		var35.Url = types.StringValue(ans.Type.Domain.Url)
	}
	var var41 *objectsExternalDynamicListsRsModelImeiObject
	if ans.Type.Imei != nil {
		var41 = &objectsExternalDynamicListsRsModelImeiObject{}
		var var42 *objectsExternalDynamicListsRsModelAuthObject
		if ans.Type.Imei.Auth != nil {
			var42 = &objectsExternalDynamicListsRsModelAuthObject{}
			var42.Password = types.StringValue(ans.Type.Imei.Auth.Password)
			var42.Username = types.StringValue(ans.Type.Imei.Auth.Username)
		}
		var var43 objectsExternalDynamicListsRsModelRecurringObject
		var var44 *objectsExternalDynamicListsRsModelDailyObject
		if ans.Type.Imei.Recurring.Daily != nil {
			var44 = &objectsExternalDynamicListsRsModelDailyObject{}
			var44.At = types.StringValue(ans.Type.Imei.Recurring.Daily.At)
		}
		var var45 *objectsExternalDynamicListsRsModelMonthlyObject
		if ans.Type.Imei.Recurring.Monthly != nil {
			var45 = &objectsExternalDynamicListsRsModelMonthlyObject{}
			var45.At = types.StringValue(ans.Type.Imei.Recurring.Monthly.At)
			var45.DayOfMonth = types.Int64Value(ans.Type.Imei.Recurring.Monthly.DayOfMonth)
		}
		var var46 *objectsExternalDynamicListsRsModelWeeklyObject
		if ans.Type.Imei.Recurring.Weekly != nil {
			var46 = &objectsExternalDynamicListsRsModelWeeklyObject{}
			var46.At = types.StringValue(ans.Type.Imei.Recurring.Weekly.At)
			var46.DayOfWeek = types.StringValue(ans.Type.Imei.Recurring.Weekly.DayOfWeek)
		}
		var43.Daily = var44
		if ans.Type.Imei.Recurring.FiveMinute != nil {
			var43.FiveMinute = types.BoolValue(true)
		}
		if ans.Type.Imei.Recurring.Hourly != nil {
			var43.Hourly = types.BoolValue(true)
		}
		var43.Monthly = var45
		var43.Weekly = var46
		var41.Auth = var42
		var41.CertificateProfile = types.StringValue(ans.Type.Imei.CertificateProfile)
		var41.Description = types.StringValue(ans.Type.Imei.Description)
		var41.ExceptionList = EncodeStringSlice(ans.Type.Imei.ExceptionList)
		var41.Recurring = var43
		var41.Url = types.StringValue(ans.Type.Imei.Url)
	}
	var var47 *objectsExternalDynamicListsRsModelImsiObject
	if ans.Type.Imsi != nil {
		var47 = &objectsExternalDynamicListsRsModelImsiObject{}
		var var48 *objectsExternalDynamicListsRsModelAuthObject
		if ans.Type.Imsi.Auth != nil {
			var48 = &objectsExternalDynamicListsRsModelAuthObject{}
			var48.Password = types.StringValue(ans.Type.Imsi.Auth.Password)
			var48.Username = types.StringValue(ans.Type.Imsi.Auth.Username)
		}
		var var49 objectsExternalDynamicListsRsModelRecurringObject
		var var50 *objectsExternalDynamicListsRsModelDailyObject
		if ans.Type.Imsi.Recurring.Daily != nil {
			var50 = &objectsExternalDynamicListsRsModelDailyObject{}
			var50.At = types.StringValue(ans.Type.Imsi.Recurring.Daily.At)
		}
		var var51 *objectsExternalDynamicListsRsModelMonthlyObject
		if ans.Type.Imsi.Recurring.Monthly != nil {
			var51 = &objectsExternalDynamicListsRsModelMonthlyObject{}
			var51.At = types.StringValue(ans.Type.Imsi.Recurring.Monthly.At)
			var51.DayOfMonth = types.Int64Value(ans.Type.Imsi.Recurring.Monthly.DayOfMonth)
		}
		var var52 *objectsExternalDynamicListsRsModelWeeklyObject
		if ans.Type.Imsi.Recurring.Weekly != nil {
			var52 = &objectsExternalDynamicListsRsModelWeeklyObject{}
			var52.At = types.StringValue(ans.Type.Imsi.Recurring.Weekly.At)
			var52.DayOfWeek = types.StringValue(ans.Type.Imsi.Recurring.Weekly.DayOfWeek)
		}
		var49.Daily = var50
		if ans.Type.Imsi.Recurring.FiveMinute != nil {
			var49.FiveMinute = types.BoolValue(true)
		}
		if ans.Type.Imsi.Recurring.Hourly != nil {
			var49.Hourly = types.BoolValue(true)
		}
		var49.Monthly = var51
		var49.Weekly = var52
		var47.Auth = var48
		var47.CertificateProfile = types.StringValue(ans.Type.Imsi.CertificateProfile)
		var47.Description = types.StringValue(ans.Type.Imsi.Description)
		var47.ExceptionList = EncodeStringSlice(ans.Type.Imsi.ExceptionList)
		var47.Recurring = var49
		var47.Url = types.StringValue(ans.Type.Imsi.Url)
	}
	var var53 *objectsExternalDynamicListsRsModelIpObject
	if ans.Type.Ip != nil {
		var53 = &objectsExternalDynamicListsRsModelIpObject{}
		var var54 *objectsExternalDynamicListsRsModelAuthObject
		if ans.Type.Ip.Auth != nil {
			var54 = &objectsExternalDynamicListsRsModelAuthObject{}
			var54.Password = types.StringValue(ans.Type.Ip.Auth.Password)
			var54.Username = types.StringValue(ans.Type.Ip.Auth.Username)
		}
		var var55 objectsExternalDynamicListsRsModelRecurringObject
		var var56 *objectsExternalDynamicListsRsModelDailyObject
		if ans.Type.Ip.Recurring.Daily != nil {
			var56 = &objectsExternalDynamicListsRsModelDailyObject{}
			var56.At = types.StringValue(ans.Type.Ip.Recurring.Daily.At)
		}
		var var57 *objectsExternalDynamicListsRsModelMonthlyObject
		if ans.Type.Ip.Recurring.Monthly != nil {
			var57 = &objectsExternalDynamicListsRsModelMonthlyObject{}
			var57.At = types.StringValue(ans.Type.Ip.Recurring.Monthly.At)
			var57.DayOfMonth = types.Int64Value(ans.Type.Ip.Recurring.Monthly.DayOfMonth)
		}
		var var58 *objectsExternalDynamicListsRsModelWeeklyObject
		if ans.Type.Ip.Recurring.Weekly != nil {
			var58 = &objectsExternalDynamicListsRsModelWeeklyObject{}
			var58.At = types.StringValue(ans.Type.Ip.Recurring.Weekly.At)
			var58.DayOfWeek = types.StringValue(ans.Type.Ip.Recurring.Weekly.DayOfWeek)
		}
		var55.Daily = var56
		if ans.Type.Ip.Recurring.FiveMinute != nil {
			var55.FiveMinute = types.BoolValue(true)
		}
		if ans.Type.Ip.Recurring.Hourly != nil {
			var55.Hourly = types.BoolValue(true)
		}
		var55.Monthly = var57
		var55.Weekly = var58
		var53.Auth = var54
		var53.CertificateProfile = types.StringValue(ans.Type.Ip.CertificateProfile)
		var53.Description = types.StringValue(ans.Type.Ip.Description)
		var53.ExceptionList = EncodeStringSlice(ans.Type.Ip.ExceptionList)
		var53.Recurring = var55
		var53.Url = types.StringValue(ans.Type.Ip.Url)
	}
	var var59 *objectsExternalDynamicListsRsModelPredefinedIpObject
	if ans.Type.PredefinedIp != nil {
		var59 = &objectsExternalDynamicListsRsModelPredefinedIpObject{}
		var59.Description = types.StringValue(ans.Type.PredefinedIp.Description)
		var59.ExceptionList = EncodeStringSlice(ans.Type.PredefinedIp.ExceptionList)
		var59.Url = types.StringValue(ans.Type.PredefinedIp.Url)
	}
	var var60 *objectsExternalDynamicListsRsModelPredefinedUrlObject
	if ans.Type.PredefinedUrl != nil {
		var60 = &objectsExternalDynamicListsRsModelPredefinedUrlObject{}
		var60.Description = types.StringValue(ans.Type.PredefinedUrl.Description)
		var60.ExceptionList = EncodeStringSlice(ans.Type.PredefinedUrl.ExceptionList)
		var60.Url = types.StringValue(ans.Type.PredefinedUrl.Url)
	}
	var var61 *objectsExternalDynamicListsRsModelUrlObject
	if ans.Type.Url != nil {
		var61 = &objectsExternalDynamicListsRsModelUrlObject{}
		var var62 *objectsExternalDynamicListsRsModelAuthObject
		if ans.Type.Url.Auth != nil {
			var62 = &objectsExternalDynamicListsRsModelAuthObject{}
			var62.Password = types.StringValue(ans.Type.Url.Auth.Password)
			var62.Username = types.StringValue(ans.Type.Url.Auth.Username)
		}
		var var63 objectsExternalDynamicListsRsModelRecurringObject
		var var64 *objectsExternalDynamicListsRsModelDailyObject
		if ans.Type.Url.Recurring.Daily != nil {
			var64 = &objectsExternalDynamicListsRsModelDailyObject{}
			var64.At = types.StringValue(ans.Type.Url.Recurring.Daily.At)
		}
		var var65 *objectsExternalDynamicListsRsModelMonthlyObject
		if ans.Type.Url.Recurring.Monthly != nil {
			var65 = &objectsExternalDynamicListsRsModelMonthlyObject{}
			var65.At = types.StringValue(ans.Type.Url.Recurring.Monthly.At)
			var65.DayOfMonth = types.Int64Value(ans.Type.Url.Recurring.Monthly.DayOfMonth)
		}
		var var66 *objectsExternalDynamicListsRsModelWeeklyObject
		if ans.Type.Url.Recurring.Weekly != nil {
			var66 = &objectsExternalDynamicListsRsModelWeeklyObject{}
			var66.At = types.StringValue(ans.Type.Url.Recurring.Weekly.At)
			var66.DayOfWeek = types.StringValue(ans.Type.Url.Recurring.Weekly.DayOfWeek)
		}
		var63.Daily = var64
		if ans.Type.Url.Recurring.FiveMinute != nil {
			var63.FiveMinute = types.BoolValue(true)
		}
		if ans.Type.Url.Recurring.Hourly != nil {
			var63.Hourly = types.BoolValue(true)
		}
		var63.Monthly = var65
		var63.Weekly = var66
		var61.Auth = var62
		var61.CertificateProfile = types.StringValue(ans.Type.Url.CertificateProfile)
		var61.Description = types.StringValue(ans.Type.Url.Description)
		var61.ExceptionList = EncodeStringSlice(ans.Type.Url.ExceptionList)
		var61.Recurring = var63
		var61.Url = types.StringValue(ans.Type.Url.Url)
	}
	var34.Domain = var35
	var34.Imei = var41
	var34.Imsi = var47
	var34.Ip = var53
	var34.PredefinedIp = var59
	var34.PredefinedUrl = var60
	var34.Url = var61
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.Type = var34

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *objectsExternalDynamicListsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state objectsExternalDynamicListsRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "sase_objects_external_dynamic_lists",
		"locMap":        map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":        tokens,
	})

	// Prepare to read the config.
	svc := iHJqznH.NewClient(r.client)
	input := iHJqznH.ReadInput{
		ObjectId: tokens[1],
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
	var var0 objectsExternalDynamicListsRsModelTypeObject
	var var1 *objectsExternalDynamicListsRsModelDomainObject
	if ans.Type.Domain != nil {
		var1 = &objectsExternalDynamicListsRsModelDomainObject{}
		var var2 *objectsExternalDynamicListsRsModelAuthObject
		if ans.Type.Domain.Auth != nil {
			var2 = &objectsExternalDynamicListsRsModelAuthObject{}
			var2.Password = types.StringValue(ans.Type.Domain.Auth.Password)
			var2.Username = types.StringValue(ans.Type.Domain.Auth.Username)
		}
		var var3 objectsExternalDynamicListsRsModelRecurringObject
		var var4 *objectsExternalDynamicListsRsModelDailyObject
		if ans.Type.Domain.Recurring.Daily != nil {
			var4 = &objectsExternalDynamicListsRsModelDailyObject{}
			var4.At = types.StringValue(ans.Type.Domain.Recurring.Daily.At)
		}
		var var5 *objectsExternalDynamicListsRsModelMonthlyObject
		if ans.Type.Domain.Recurring.Monthly != nil {
			var5 = &objectsExternalDynamicListsRsModelMonthlyObject{}
			var5.At = types.StringValue(ans.Type.Domain.Recurring.Monthly.At)
			var5.DayOfMonth = types.Int64Value(ans.Type.Domain.Recurring.Monthly.DayOfMonth)
		}
		var var6 *objectsExternalDynamicListsRsModelWeeklyObject
		if ans.Type.Domain.Recurring.Weekly != nil {
			var6 = &objectsExternalDynamicListsRsModelWeeklyObject{}
			var6.At = types.StringValue(ans.Type.Domain.Recurring.Weekly.At)
			var6.DayOfWeek = types.StringValue(ans.Type.Domain.Recurring.Weekly.DayOfWeek)
		}
		var3.Daily = var4
		if ans.Type.Domain.Recurring.FiveMinute != nil {
			var3.FiveMinute = types.BoolValue(true)
		}
		if ans.Type.Domain.Recurring.Hourly != nil {
			var3.Hourly = types.BoolValue(true)
		}
		var3.Monthly = var5
		var3.Weekly = var6
		var1.Auth = var2
		var1.CertificateProfile = types.StringValue(ans.Type.Domain.CertificateProfile)
		var1.Description = types.StringValue(ans.Type.Domain.Description)
		var1.ExceptionList = EncodeStringSlice(ans.Type.Domain.ExceptionList)
		var1.ExpandDomain = types.BoolValue(ans.Type.Domain.ExpandDomain)
		var1.Recurring = var3
		var1.Url = types.StringValue(ans.Type.Domain.Url)
	}
	var var7 *objectsExternalDynamicListsRsModelImeiObject
	if ans.Type.Imei != nil {
		var7 = &objectsExternalDynamicListsRsModelImeiObject{}
		var var8 *objectsExternalDynamicListsRsModelAuthObject
		if ans.Type.Imei.Auth != nil {
			var8 = &objectsExternalDynamicListsRsModelAuthObject{}
			var8.Password = types.StringValue(ans.Type.Imei.Auth.Password)
			var8.Username = types.StringValue(ans.Type.Imei.Auth.Username)
		}
		var var9 objectsExternalDynamicListsRsModelRecurringObject
		var var10 *objectsExternalDynamicListsRsModelDailyObject
		if ans.Type.Imei.Recurring.Daily != nil {
			var10 = &objectsExternalDynamicListsRsModelDailyObject{}
			var10.At = types.StringValue(ans.Type.Imei.Recurring.Daily.At)
		}
		var var11 *objectsExternalDynamicListsRsModelMonthlyObject
		if ans.Type.Imei.Recurring.Monthly != nil {
			var11 = &objectsExternalDynamicListsRsModelMonthlyObject{}
			var11.At = types.StringValue(ans.Type.Imei.Recurring.Monthly.At)
			var11.DayOfMonth = types.Int64Value(ans.Type.Imei.Recurring.Monthly.DayOfMonth)
		}
		var var12 *objectsExternalDynamicListsRsModelWeeklyObject
		if ans.Type.Imei.Recurring.Weekly != nil {
			var12 = &objectsExternalDynamicListsRsModelWeeklyObject{}
			var12.At = types.StringValue(ans.Type.Imei.Recurring.Weekly.At)
			var12.DayOfWeek = types.StringValue(ans.Type.Imei.Recurring.Weekly.DayOfWeek)
		}
		var9.Daily = var10
		if ans.Type.Imei.Recurring.FiveMinute != nil {
			var9.FiveMinute = types.BoolValue(true)
		}
		if ans.Type.Imei.Recurring.Hourly != nil {
			var9.Hourly = types.BoolValue(true)
		}
		var9.Monthly = var11
		var9.Weekly = var12
		var7.Auth = var8
		var7.CertificateProfile = types.StringValue(ans.Type.Imei.CertificateProfile)
		var7.Description = types.StringValue(ans.Type.Imei.Description)
		var7.ExceptionList = EncodeStringSlice(ans.Type.Imei.ExceptionList)
		var7.Recurring = var9
		var7.Url = types.StringValue(ans.Type.Imei.Url)
	}
	var var13 *objectsExternalDynamicListsRsModelImsiObject
	if ans.Type.Imsi != nil {
		var13 = &objectsExternalDynamicListsRsModelImsiObject{}
		var var14 *objectsExternalDynamicListsRsModelAuthObject
		if ans.Type.Imsi.Auth != nil {
			var14 = &objectsExternalDynamicListsRsModelAuthObject{}
			var14.Password = types.StringValue(ans.Type.Imsi.Auth.Password)
			var14.Username = types.StringValue(ans.Type.Imsi.Auth.Username)
		}
		var var15 objectsExternalDynamicListsRsModelRecurringObject
		var var16 *objectsExternalDynamicListsRsModelDailyObject
		if ans.Type.Imsi.Recurring.Daily != nil {
			var16 = &objectsExternalDynamicListsRsModelDailyObject{}
			var16.At = types.StringValue(ans.Type.Imsi.Recurring.Daily.At)
		}
		var var17 *objectsExternalDynamicListsRsModelMonthlyObject
		if ans.Type.Imsi.Recurring.Monthly != nil {
			var17 = &objectsExternalDynamicListsRsModelMonthlyObject{}
			var17.At = types.StringValue(ans.Type.Imsi.Recurring.Monthly.At)
			var17.DayOfMonth = types.Int64Value(ans.Type.Imsi.Recurring.Monthly.DayOfMonth)
		}
		var var18 *objectsExternalDynamicListsRsModelWeeklyObject
		if ans.Type.Imsi.Recurring.Weekly != nil {
			var18 = &objectsExternalDynamicListsRsModelWeeklyObject{}
			var18.At = types.StringValue(ans.Type.Imsi.Recurring.Weekly.At)
			var18.DayOfWeek = types.StringValue(ans.Type.Imsi.Recurring.Weekly.DayOfWeek)
		}
		var15.Daily = var16
		if ans.Type.Imsi.Recurring.FiveMinute != nil {
			var15.FiveMinute = types.BoolValue(true)
		}
		if ans.Type.Imsi.Recurring.Hourly != nil {
			var15.Hourly = types.BoolValue(true)
		}
		var15.Monthly = var17
		var15.Weekly = var18
		var13.Auth = var14
		var13.CertificateProfile = types.StringValue(ans.Type.Imsi.CertificateProfile)
		var13.Description = types.StringValue(ans.Type.Imsi.Description)
		var13.ExceptionList = EncodeStringSlice(ans.Type.Imsi.ExceptionList)
		var13.Recurring = var15
		var13.Url = types.StringValue(ans.Type.Imsi.Url)
	}
	var var19 *objectsExternalDynamicListsRsModelIpObject
	if ans.Type.Ip != nil {
		var19 = &objectsExternalDynamicListsRsModelIpObject{}
		var var20 *objectsExternalDynamicListsRsModelAuthObject
		if ans.Type.Ip.Auth != nil {
			var20 = &objectsExternalDynamicListsRsModelAuthObject{}
			var20.Password = types.StringValue(ans.Type.Ip.Auth.Password)
			var20.Username = types.StringValue(ans.Type.Ip.Auth.Username)
		}
		var var21 objectsExternalDynamicListsRsModelRecurringObject
		var var22 *objectsExternalDynamicListsRsModelDailyObject
		if ans.Type.Ip.Recurring.Daily != nil {
			var22 = &objectsExternalDynamicListsRsModelDailyObject{}
			var22.At = types.StringValue(ans.Type.Ip.Recurring.Daily.At)
		}
		var var23 *objectsExternalDynamicListsRsModelMonthlyObject
		if ans.Type.Ip.Recurring.Monthly != nil {
			var23 = &objectsExternalDynamicListsRsModelMonthlyObject{}
			var23.At = types.StringValue(ans.Type.Ip.Recurring.Monthly.At)
			var23.DayOfMonth = types.Int64Value(ans.Type.Ip.Recurring.Monthly.DayOfMonth)
		}
		var var24 *objectsExternalDynamicListsRsModelWeeklyObject
		if ans.Type.Ip.Recurring.Weekly != nil {
			var24 = &objectsExternalDynamicListsRsModelWeeklyObject{}
			var24.At = types.StringValue(ans.Type.Ip.Recurring.Weekly.At)
			var24.DayOfWeek = types.StringValue(ans.Type.Ip.Recurring.Weekly.DayOfWeek)
		}
		var21.Daily = var22
		if ans.Type.Ip.Recurring.FiveMinute != nil {
			var21.FiveMinute = types.BoolValue(true)
		}
		if ans.Type.Ip.Recurring.Hourly != nil {
			var21.Hourly = types.BoolValue(true)
		}
		var21.Monthly = var23
		var21.Weekly = var24
		var19.Auth = var20
		var19.CertificateProfile = types.StringValue(ans.Type.Ip.CertificateProfile)
		var19.Description = types.StringValue(ans.Type.Ip.Description)
		var19.ExceptionList = EncodeStringSlice(ans.Type.Ip.ExceptionList)
		var19.Recurring = var21
		var19.Url = types.StringValue(ans.Type.Ip.Url)
	}
	var var25 *objectsExternalDynamicListsRsModelPredefinedIpObject
	if ans.Type.PredefinedIp != nil {
		var25 = &objectsExternalDynamicListsRsModelPredefinedIpObject{}
		var25.Description = types.StringValue(ans.Type.PredefinedIp.Description)
		var25.ExceptionList = EncodeStringSlice(ans.Type.PredefinedIp.ExceptionList)
		var25.Url = types.StringValue(ans.Type.PredefinedIp.Url)
	}
	var var26 *objectsExternalDynamicListsRsModelPredefinedUrlObject
	if ans.Type.PredefinedUrl != nil {
		var26 = &objectsExternalDynamicListsRsModelPredefinedUrlObject{}
		var26.Description = types.StringValue(ans.Type.PredefinedUrl.Description)
		var26.ExceptionList = EncodeStringSlice(ans.Type.PredefinedUrl.ExceptionList)
		var26.Url = types.StringValue(ans.Type.PredefinedUrl.Url)
	}
	var var27 *objectsExternalDynamicListsRsModelUrlObject
	if ans.Type.Url != nil {
		var27 = &objectsExternalDynamicListsRsModelUrlObject{}
		var var28 *objectsExternalDynamicListsRsModelAuthObject
		if ans.Type.Url.Auth != nil {
			var28 = &objectsExternalDynamicListsRsModelAuthObject{}
			var28.Password = types.StringValue(ans.Type.Url.Auth.Password)
			var28.Username = types.StringValue(ans.Type.Url.Auth.Username)
		}
		var var29 objectsExternalDynamicListsRsModelRecurringObject
		var var30 *objectsExternalDynamicListsRsModelDailyObject
		if ans.Type.Url.Recurring.Daily != nil {
			var30 = &objectsExternalDynamicListsRsModelDailyObject{}
			var30.At = types.StringValue(ans.Type.Url.Recurring.Daily.At)
		}
		var var31 *objectsExternalDynamicListsRsModelMonthlyObject
		if ans.Type.Url.Recurring.Monthly != nil {
			var31 = &objectsExternalDynamicListsRsModelMonthlyObject{}
			var31.At = types.StringValue(ans.Type.Url.Recurring.Monthly.At)
			var31.DayOfMonth = types.Int64Value(ans.Type.Url.Recurring.Monthly.DayOfMonth)
		}
		var var32 *objectsExternalDynamicListsRsModelWeeklyObject
		if ans.Type.Url.Recurring.Weekly != nil {
			var32 = &objectsExternalDynamicListsRsModelWeeklyObject{}
			var32.At = types.StringValue(ans.Type.Url.Recurring.Weekly.At)
			var32.DayOfWeek = types.StringValue(ans.Type.Url.Recurring.Weekly.DayOfWeek)
		}
		var29.Daily = var30
		if ans.Type.Url.Recurring.FiveMinute != nil {
			var29.FiveMinute = types.BoolValue(true)
		}
		if ans.Type.Url.Recurring.Hourly != nil {
			var29.Hourly = types.BoolValue(true)
		}
		var29.Monthly = var31
		var29.Weekly = var32
		var27.Auth = var28
		var27.CertificateProfile = types.StringValue(ans.Type.Url.CertificateProfile)
		var27.Description = types.StringValue(ans.Type.Url.Description)
		var27.ExceptionList = EncodeStringSlice(ans.Type.Url.ExceptionList)
		var27.Recurring = var29
		var27.Url = types.StringValue(ans.Type.Url.Url)
	}
	var0.Domain = var1
	var0.Imei = var7
	var0.Imsi = var13
	var0.Ip = var19
	var0.PredefinedIp = var25
	var0.PredefinedUrl = var26
	var0.Url = var27
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.Type = var0

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *objectsExternalDynamicListsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state objectsExternalDynamicListsRsModel
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
		"resource_name": "sase_objects_external_dynamic_lists",
		"object_id":     state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := iHJqznH.NewClient(r.client)
	input := iHJqznH.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 sRAOviP.Config
	var0.Name = plan.Name.ValueString()
	var var1 sRAOviP.TypeObject
	var var2 *sRAOviP.DomainObject
	if plan.Type.Domain != nil {
		var2 = &sRAOviP.DomainObject{}
		var var3 *sRAOviP.AuthObject
		if plan.Type.Domain.Auth != nil {
			var3 = &sRAOviP.AuthObject{}
			var3.Password = plan.Type.Domain.Auth.Password.ValueString()
			var3.Username = plan.Type.Domain.Auth.Username.ValueString()
		}
		var2.Auth = var3
		var2.CertificateProfile = plan.Type.Domain.CertificateProfile.ValueString()
		var2.Description = plan.Type.Domain.Description.ValueString()
		var2.ExceptionList = DecodeStringSlice(plan.Type.Domain.ExceptionList)
		var2.ExpandDomain = plan.Type.Domain.ExpandDomain.ValueBool()
		var var4 sRAOviP.RecurringObject
		var var5 *sRAOviP.DailyObject
		if plan.Type.Domain.Recurring.Daily != nil {
			var5 = &sRAOviP.DailyObject{}
			var5.At = plan.Type.Domain.Recurring.Daily.At.ValueString()
		}
		var4.Daily = var5
		if plan.Type.Domain.Recurring.FiveMinute.ValueBool() {
			var4.FiveMinute = struct{}{}
		}
		if plan.Type.Domain.Recurring.Hourly.ValueBool() {
			var4.Hourly = struct{}{}
		}
		var var6 *sRAOviP.MonthlyObject
		if plan.Type.Domain.Recurring.Monthly != nil {
			var6 = &sRAOviP.MonthlyObject{}
			var6.At = plan.Type.Domain.Recurring.Monthly.At.ValueString()
			var6.DayOfMonth = plan.Type.Domain.Recurring.Monthly.DayOfMonth.ValueInt64()
		}
		var4.Monthly = var6
		var var7 *sRAOviP.WeeklyObject
		if plan.Type.Domain.Recurring.Weekly != nil {
			var7 = &sRAOviP.WeeklyObject{}
			var7.At = plan.Type.Domain.Recurring.Weekly.At.ValueString()
			var7.DayOfWeek = plan.Type.Domain.Recurring.Weekly.DayOfWeek.ValueString()
		}
		var4.Weekly = var7
		var2.Recurring = var4
		var2.Url = plan.Type.Domain.Url.ValueString()
	}
	var1.Domain = var2
	var var8 *sRAOviP.ImeiObject
	if plan.Type.Imei != nil {
		var8 = &sRAOviP.ImeiObject{}
		var var9 *sRAOviP.AuthObject
		if plan.Type.Imei.Auth != nil {
			var9 = &sRAOviP.AuthObject{}
			var9.Password = plan.Type.Imei.Auth.Password.ValueString()
			var9.Username = plan.Type.Imei.Auth.Username.ValueString()
		}
		var8.Auth = var9
		var8.CertificateProfile = plan.Type.Imei.CertificateProfile.ValueString()
		var8.Description = plan.Type.Imei.Description.ValueString()
		var8.ExceptionList = DecodeStringSlice(plan.Type.Imei.ExceptionList)
		var var10 sRAOviP.RecurringObject
		var var11 *sRAOviP.DailyObject
		if plan.Type.Imei.Recurring.Daily != nil {
			var11 = &sRAOviP.DailyObject{}
			var11.At = plan.Type.Imei.Recurring.Daily.At.ValueString()
		}
		var10.Daily = var11
		if plan.Type.Imei.Recurring.FiveMinute.ValueBool() {
			var10.FiveMinute = struct{}{}
		}
		if plan.Type.Imei.Recurring.Hourly.ValueBool() {
			var10.Hourly = struct{}{}
		}
		var var12 *sRAOviP.MonthlyObject
		if plan.Type.Imei.Recurring.Monthly != nil {
			var12 = &sRAOviP.MonthlyObject{}
			var12.At = plan.Type.Imei.Recurring.Monthly.At.ValueString()
			var12.DayOfMonth = plan.Type.Imei.Recurring.Monthly.DayOfMonth.ValueInt64()
		}
		var10.Monthly = var12
		var var13 *sRAOviP.WeeklyObject
		if plan.Type.Imei.Recurring.Weekly != nil {
			var13 = &sRAOviP.WeeklyObject{}
			var13.At = plan.Type.Imei.Recurring.Weekly.At.ValueString()
			var13.DayOfWeek = plan.Type.Imei.Recurring.Weekly.DayOfWeek.ValueString()
		}
		var10.Weekly = var13
		var8.Recurring = var10
		var8.Url = plan.Type.Imei.Url.ValueString()
	}
	var1.Imei = var8
	var var14 *sRAOviP.ImsiObject
	if plan.Type.Imsi != nil {
		var14 = &sRAOviP.ImsiObject{}
		var var15 *sRAOviP.AuthObject
		if plan.Type.Imsi.Auth != nil {
			var15 = &sRAOviP.AuthObject{}
			var15.Password = plan.Type.Imsi.Auth.Password.ValueString()
			var15.Username = plan.Type.Imsi.Auth.Username.ValueString()
		}
		var14.Auth = var15
		var14.CertificateProfile = plan.Type.Imsi.CertificateProfile.ValueString()
		var14.Description = plan.Type.Imsi.Description.ValueString()
		var14.ExceptionList = DecodeStringSlice(plan.Type.Imsi.ExceptionList)
		var var16 sRAOviP.RecurringObject
		var var17 *sRAOviP.DailyObject
		if plan.Type.Imsi.Recurring.Daily != nil {
			var17 = &sRAOviP.DailyObject{}
			var17.At = plan.Type.Imsi.Recurring.Daily.At.ValueString()
		}
		var16.Daily = var17
		if plan.Type.Imsi.Recurring.FiveMinute.ValueBool() {
			var16.FiveMinute = struct{}{}
		}
		if plan.Type.Imsi.Recurring.Hourly.ValueBool() {
			var16.Hourly = struct{}{}
		}
		var var18 *sRAOviP.MonthlyObject
		if plan.Type.Imsi.Recurring.Monthly != nil {
			var18 = &sRAOviP.MonthlyObject{}
			var18.At = plan.Type.Imsi.Recurring.Monthly.At.ValueString()
			var18.DayOfMonth = plan.Type.Imsi.Recurring.Monthly.DayOfMonth.ValueInt64()
		}
		var16.Monthly = var18
		var var19 *sRAOviP.WeeklyObject
		if plan.Type.Imsi.Recurring.Weekly != nil {
			var19 = &sRAOviP.WeeklyObject{}
			var19.At = plan.Type.Imsi.Recurring.Weekly.At.ValueString()
			var19.DayOfWeek = plan.Type.Imsi.Recurring.Weekly.DayOfWeek.ValueString()
		}
		var16.Weekly = var19
		var14.Recurring = var16
		var14.Url = plan.Type.Imsi.Url.ValueString()
	}
	var1.Imsi = var14
	var var20 *sRAOviP.IpObject
	if plan.Type.Ip != nil {
		var20 = &sRAOviP.IpObject{}
		var var21 *sRAOviP.AuthObject
		if plan.Type.Ip.Auth != nil {
			var21 = &sRAOviP.AuthObject{}
			var21.Password = plan.Type.Ip.Auth.Password.ValueString()
			var21.Username = plan.Type.Ip.Auth.Username.ValueString()
		}
		var20.Auth = var21
		var20.CertificateProfile = plan.Type.Ip.CertificateProfile.ValueString()
		var20.Description = plan.Type.Ip.Description.ValueString()
		var20.ExceptionList = DecodeStringSlice(plan.Type.Ip.ExceptionList)
		var var22 sRAOviP.RecurringObject
		var var23 *sRAOviP.DailyObject
		if plan.Type.Ip.Recurring.Daily != nil {
			var23 = &sRAOviP.DailyObject{}
			var23.At = plan.Type.Ip.Recurring.Daily.At.ValueString()
		}
		var22.Daily = var23
		if plan.Type.Ip.Recurring.FiveMinute.ValueBool() {
			var22.FiveMinute = struct{}{}
		}
		if plan.Type.Ip.Recurring.Hourly.ValueBool() {
			var22.Hourly = struct{}{}
		}
		var var24 *sRAOviP.MonthlyObject
		if plan.Type.Ip.Recurring.Monthly != nil {
			var24 = &sRAOviP.MonthlyObject{}
			var24.At = plan.Type.Ip.Recurring.Monthly.At.ValueString()
			var24.DayOfMonth = plan.Type.Ip.Recurring.Monthly.DayOfMonth.ValueInt64()
		}
		var22.Monthly = var24
		var var25 *sRAOviP.WeeklyObject
		if plan.Type.Ip.Recurring.Weekly != nil {
			var25 = &sRAOviP.WeeklyObject{}
			var25.At = plan.Type.Ip.Recurring.Weekly.At.ValueString()
			var25.DayOfWeek = plan.Type.Ip.Recurring.Weekly.DayOfWeek.ValueString()
		}
		var22.Weekly = var25
		var20.Recurring = var22
		var20.Url = plan.Type.Ip.Url.ValueString()
	}
	var1.Ip = var20
	var var26 *sRAOviP.PredefinedIpObject
	if plan.Type.PredefinedIp != nil {
		var26 = &sRAOviP.PredefinedIpObject{}
		var26.Description = plan.Type.PredefinedIp.Description.ValueString()
		var26.ExceptionList = DecodeStringSlice(plan.Type.PredefinedIp.ExceptionList)
		var26.Url = plan.Type.PredefinedIp.Url.ValueString()
	}
	var1.PredefinedIp = var26
	var var27 *sRAOviP.PredefinedUrlObject
	if plan.Type.PredefinedUrl != nil {
		var27 = &sRAOviP.PredefinedUrlObject{}
		var27.Description = plan.Type.PredefinedUrl.Description.ValueString()
		var27.ExceptionList = DecodeStringSlice(plan.Type.PredefinedUrl.ExceptionList)
		var27.Url = plan.Type.PredefinedUrl.Url.ValueString()
	}
	var1.PredefinedUrl = var27
	var var28 *sRAOviP.UrlObject
	if plan.Type.Url != nil {
		var28 = &sRAOviP.UrlObject{}
		var var29 *sRAOviP.AuthObject
		if plan.Type.Url.Auth != nil {
			var29 = &sRAOviP.AuthObject{}
			var29.Password = plan.Type.Url.Auth.Password.ValueString()
			var29.Username = plan.Type.Url.Auth.Username.ValueString()
		}
		var28.Auth = var29
		var28.CertificateProfile = plan.Type.Url.CertificateProfile.ValueString()
		var28.Description = plan.Type.Url.Description.ValueString()
		var28.ExceptionList = DecodeStringSlice(plan.Type.Url.ExceptionList)
		var var30 sRAOviP.RecurringObject
		var var31 *sRAOviP.DailyObject
		if plan.Type.Url.Recurring.Daily != nil {
			var31 = &sRAOviP.DailyObject{}
			var31.At = plan.Type.Url.Recurring.Daily.At.ValueString()
		}
		var30.Daily = var31
		if plan.Type.Url.Recurring.FiveMinute.ValueBool() {
			var30.FiveMinute = struct{}{}
		}
		if plan.Type.Url.Recurring.Hourly.ValueBool() {
			var30.Hourly = struct{}{}
		}
		var var32 *sRAOviP.MonthlyObject
		if plan.Type.Url.Recurring.Monthly != nil {
			var32 = &sRAOviP.MonthlyObject{}
			var32.At = plan.Type.Url.Recurring.Monthly.At.ValueString()
			var32.DayOfMonth = plan.Type.Url.Recurring.Monthly.DayOfMonth.ValueInt64()
		}
		var30.Monthly = var32
		var var33 *sRAOviP.WeeklyObject
		if plan.Type.Url.Recurring.Weekly != nil {
			var33 = &sRAOviP.WeeklyObject{}
			var33.At = plan.Type.Url.Recurring.Weekly.At.ValueString()
			var33.DayOfWeek = plan.Type.Url.Recurring.Weekly.DayOfWeek.ValueString()
		}
		var30.Weekly = var33
		var28.Recurring = var30
		var28.Url = plan.Type.Url.Url.ValueString()
	}
	var1.Url = var28
	var0.Type = var1
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var34 objectsExternalDynamicListsRsModelTypeObject
	var var35 *objectsExternalDynamicListsRsModelDomainObject
	if ans.Type.Domain != nil {
		var35 = &objectsExternalDynamicListsRsModelDomainObject{}
		var var36 *objectsExternalDynamicListsRsModelAuthObject
		if ans.Type.Domain.Auth != nil {
			var36 = &objectsExternalDynamicListsRsModelAuthObject{}
			var36.Password = types.StringValue(ans.Type.Domain.Auth.Password)
			var36.Username = types.StringValue(ans.Type.Domain.Auth.Username)
		}
		var var37 objectsExternalDynamicListsRsModelRecurringObject
		var var38 *objectsExternalDynamicListsRsModelDailyObject
		if ans.Type.Domain.Recurring.Daily != nil {
			var38 = &objectsExternalDynamicListsRsModelDailyObject{}
			var38.At = types.StringValue(ans.Type.Domain.Recurring.Daily.At)
		}
		var var39 *objectsExternalDynamicListsRsModelMonthlyObject
		if ans.Type.Domain.Recurring.Monthly != nil {
			var39 = &objectsExternalDynamicListsRsModelMonthlyObject{}
			var39.At = types.StringValue(ans.Type.Domain.Recurring.Monthly.At)
			var39.DayOfMonth = types.Int64Value(ans.Type.Domain.Recurring.Monthly.DayOfMonth)
		}
		var var40 *objectsExternalDynamicListsRsModelWeeklyObject
		if ans.Type.Domain.Recurring.Weekly != nil {
			var40 = &objectsExternalDynamicListsRsModelWeeklyObject{}
			var40.At = types.StringValue(ans.Type.Domain.Recurring.Weekly.At)
			var40.DayOfWeek = types.StringValue(ans.Type.Domain.Recurring.Weekly.DayOfWeek)
		}
		var37.Daily = var38
		if ans.Type.Domain.Recurring.FiveMinute != nil {
			var37.FiveMinute = types.BoolValue(true)
		}
		if ans.Type.Domain.Recurring.Hourly != nil {
			var37.Hourly = types.BoolValue(true)
		}
		var37.Monthly = var39
		var37.Weekly = var40
		var35.Auth = var36
		var35.CertificateProfile = types.StringValue(ans.Type.Domain.CertificateProfile)
		var35.Description = types.StringValue(ans.Type.Domain.Description)
		var35.ExceptionList = EncodeStringSlice(ans.Type.Domain.ExceptionList)
		var35.ExpandDomain = types.BoolValue(ans.Type.Domain.ExpandDomain)
		var35.Recurring = var37
		var35.Url = types.StringValue(ans.Type.Domain.Url)
	}
	var var41 *objectsExternalDynamicListsRsModelImeiObject
	if ans.Type.Imei != nil {
		var41 = &objectsExternalDynamicListsRsModelImeiObject{}
		var var42 *objectsExternalDynamicListsRsModelAuthObject
		if ans.Type.Imei.Auth != nil {
			var42 = &objectsExternalDynamicListsRsModelAuthObject{}
			var42.Password = types.StringValue(ans.Type.Imei.Auth.Password)
			var42.Username = types.StringValue(ans.Type.Imei.Auth.Username)
		}
		var var43 objectsExternalDynamicListsRsModelRecurringObject
		var var44 *objectsExternalDynamicListsRsModelDailyObject
		if ans.Type.Imei.Recurring.Daily != nil {
			var44 = &objectsExternalDynamicListsRsModelDailyObject{}
			var44.At = types.StringValue(ans.Type.Imei.Recurring.Daily.At)
		}
		var var45 *objectsExternalDynamicListsRsModelMonthlyObject
		if ans.Type.Imei.Recurring.Monthly != nil {
			var45 = &objectsExternalDynamicListsRsModelMonthlyObject{}
			var45.At = types.StringValue(ans.Type.Imei.Recurring.Monthly.At)
			var45.DayOfMonth = types.Int64Value(ans.Type.Imei.Recurring.Monthly.DayOfMonth)
		}
		var var46 *objectsExternalDynamicListsRsModelWeeklyObject
		if ans.Type.Imei.Recurring.Weekly != nil {
			var46 = &objectsExternalDynamicListsRsModelWeeklyObject{}
			var46.At = types.StringValue(ans.Type.Imei.Recurring.Weekly.At)
			var46.DayOfWeek = types.StringValue(ans.Type.Imei.Recurring.Weekly.DayOfWeek)
		}
		var43.Daily = var44
		if ans.Type.Imei.Recurring.FiveMinute != nil {
			var43.FiveMinute = types.BoolValue(true)
		}
		if ans.Type.Imei.Recurring.Hourly != nil {
			var43.Hourly = types.BoolValue(true)
		}
		var43.Monthly = var45
		var43.Weekly = var46
		var41.Auth = var42
		var41.CertificateProfile = types.StringValue(ans.Type.Imei.CertificateProfile)
		var41.Description = types.StringValue(ans.Type.Imei.Description)
		var41.ExceptionList = EncodeStringSlice(ans.Type.Imei.ExceptionList)
		var41.Recurring = var43
		var41.Url = types.StringValue(ans.Type.Imei.Url)
	}
	var var47 *objectsExternalDynamicListsRsModelImsiObject
	if ans.Type.Imsi != nil {
		var47 = &objectsExternalDynamicListsRsModelImsiObject{}
		var var48 *objectsExternalDynamicListsRsModelAuthObject
		if ans.Type.Imsi.Auth != nil {
			var48 = &objectsExternalDynamicListsRsModelAuthObject{}
			var48.Password = types.StringValue(ans.Type.Imsi.Auth.Password)
			var48.Username = types.StringValue(ans.Type.Imsi.Auth.Username)
		}
		var var49 objectsExternalDynamicListsRsModelRecurringObject
		var var50 *objectsExternalDynamicListsRsModelDailyObject
		if ans.Type.Imsi.Recurring.Daily != nil {
			var50 = &objectsExternalDynamicListsRsModelDailyObject{}
			var50.At = types.StringValue(ans.Type.Imsi.Recurring.Daily.At)
		}
		var var51 *objectsExternalDynamicListsRsModelMonthlyObject
		if ans.Type.Imsi.Recurring.Monthly != nil {
			var51 = &objectsExternalDynamicListsRsModelMonthlyObject{}
			var51.At = types.StringValue(ans.Type.Imsi.Recurring.Monthly.At)
			var51.DayOfMonth = types.Int64Value(ans.Type.Imsi.Recurring.Monthly.DayOfMonth)
		}
		var var52 *objectsExternalDynamicListsRsModelWeeklyObject
		if ans.Type.Imsi.Recurring.Weekly != nil {
			var52 = &objectsExternalDynamicListsRsModelWeeklyObject{}
			var52.At = types.StringValue(ans.Type.Imsi.Recurring.Weekly.At)
			var52.DayOfWeek = types.StringValue(ans.Type.Imsi.Recurring.Weekly.DayOfWeek)
		}
		var49.Daily = var50
		if ans.Type.Imsi.Recurring.FiveMinute != nil {
			var49.FiveMinute = types.BoolValue(true)
		}
		if ans.Type.Imsi.Recurring.Hourly != nil {
			var49.Hourly = types.BoolValue(true)
		}
		var49.Monthly = var51
		var49.Weekly = var52
		var47.Auth = var48
		var47.CertificateProfile = types.StringValue(ans.Type.Imsi.CertificateProfile)
		var47.Description = types.StringValue(ans.Type.Imsi.Description)
		var47.ExceptionList = EncodeStringSlice(ans.Type.Imsi.ExceptionList)
		var47.Recurring = var49
		var47.Url = types.StringValue(ans.Type.Imsi.Url)
	}
	var var53 *objectsExternalDynamicListsRsModelIpObject
	if ans.Type.Ip != nil {
		var53 = &objectsExternalDynamicListsRsModelIpObject{}
		var var54 *objectsExternalDynamicListsRsModelAuthObject
		if ans.Type.Ip.Auth != nil {
			var54 = &objectsExternalDynamicListsRsModelAuthObject{}
			var54.Password = types.StringValue(ans.Type.Ip.Auth.Password)
			var54.Username = types.StringValue(ans.Type.Ip.Auth.Username)
		}
		var var55 objectsExternalDynamicListsRsModelRecurringObject
		var var56 *objectsExternalDynamicListsRsModelDailyObject
		if ans.Type.Ip.Recurring.Daily != nil {
			var56 = &objectsExternalDynamicListsRsModelDailyObject{}
			var56.At = types.StringValue(ans.Type.Ip.Recurring.Daily.At)
		}
		var var57 *objectsExternalDynamicListsRsModelMonthlyObject
		if ans.Type.Ip.Recurring.Monthly != nil {
			var57 = &objectsExternalDynamicListsRsModelMonthlyObject{}
			var57.At = types.StringValue(ans.Type.Ip.Recurring.Monthly.At)
			var57.DayOfMonth = types.Int64Value(ans.Type.Ip.Recurring.Monthly.DayOfMonth)
		}
		var var58 *objectsExternalDynamicListsRsModelWeeklyObject
		if ans.Type.Ip.Recurring.Weekly != nil {
			var58 = &objectsExternalDynamicListsRsModelWeeklyObject{}
			var58.At = types.StringValue(ans.Type.Ip.Recurring.Weekly.At)
			var58.DayOfWeek = types.StringValue(ans.Type.Ip.Recurring.Weekly.DayOfWeek)
		}
		var55.Daily = var56
		if ans.Type.Ip.Recurring.FiveMinute != nil {
			var55.FiveMinute = types.BoolValue(true)
		}
		if ans.Type.Ip.Recurring.Hourly != nil {
			var55.Hourly = types.BoolValue(true)
		}
		var55.Monthly = var57
		var55.Weekly = var58
		var53.Auth = var54
		var53.CertificateProfile = types.StringValue(ans.Type.Ip.CertificateProfile)
		var53.Description = types.StringValue(ans.Type.Ip.Description)
		var53.ExceptionList = EncodeStringSlice(ans.Type.Ip.ExceptionList)
		var53.Recurring = var55
		var53.Url = types.StringValue(ans.Type.Ip.Url)
	}
	var var59 *objectsExternalDynamicListsRsModelPredefinedIpObject
	if ans.Type.PredefinedIp != nil {
		var59 = &objectsExternalDynamicListsRsModelPredefinedIpObject{}
		var59.Description = types.StringValue(ans.Type.PredefinedIp.Description)
		var59.ExceptionList = EncodeStringSlice(ans.Type.PredefinedIp.ExceptionList)
		var59.Url = types.StringValue(ans.Type.PredefinedIp.Url)
	}
	var var60 *objectsExternalDynamicListsRsModelPredefinedUrlObject
	if ans.Type.PredefinedUrl != nil {
		var60 = &objectsExternalDynamicListsRsModelPredefinedUrlObject{}
		var60.Description = types.StringValue(ans.Type.PredefinedUrl.Description)
		var60.ExceptionList = EncodeStringSlice(ans.Type.PredefinedUrl.ExceptionList)
		var60.Url = types.StringValue(ans.Type.PredefinedUrl.Url)
	}
	var var61 *objectsExternalDynamicListsRsModelUrlObject
	if ans.Type.Url != nil {
		var61 = &objectsExternalDynamicListsRsModelUrlObject{}
		var var62 *objectsExternalDynamicListsRsModelAuthObject
		if ans.Type.Url.Auth != nil {
			var62 = &objectsExternalDynamicListsRsModelAuthObject{}
			var62.Password = types.StringValue(ans.Type.Url.Auth.Password)
			var62.Username = types.StringValue(ans.Type.Url.Auth.Username)
		}
		var var63 objectsExternalDynamicListsRsModelRecurringObject
		var var64 *objectsExternalDynamicListsRsModelDailyObject
		if ans.Type.Url.Recurring.Daily != nil {
			var64 = &objectsExternalDynamicListsRsModelDailyObject{}
			var64.At = types.StringValue(ans.Type.Url.Recurring.Daily.At)
		}
		var var65 *objectsExternalDynamicListsRsModelMonthlyObject
		if ans.Type.Url.Recurring.Monthly != nil {
			var65 = &objectsExternalDynamicListsRsModelMonthlyObject{}
			var65.At = types.StringValue(ans.Type.Url.Recurring.Monthly.At)
			var65.DayOfMonth = types.Int64Value(ans.Type.Url.Recurring.Monthly.DayOfMonth)
		}
		var var66 *objectsExternalDynamicListsRsModelWeeklyObject
		if ans.Type.Url.Recurring.Weekly != nil {
			var66 = &objectsExternalDynamicListsRsModelWeeklyObject{}
			var66.At = types.StringValue(ans.Type.Url.Recurring.Weekly.At)
			var66.DayOfWeek = types.StringValue(ans.Type.Url.Recurring.Weekly.DayOfWeek)
		}
		var63.Daily = var64
		if ans.Type.Url.Recurring.FiveMinute != nil {
			var63.FiveMinute = types.BoolValue(true)
		}
		if ans.Type.Url.Recurring.Hourly != nil {
			var63.Hourly = types.BoolValue(true)
		}
		var63.Monthly = var65
		var63.Weekly = var66
		var61.Auth = var62
		var61.CertificateProfile = types.StringValue(ans.Type.Url.CertificateProfile)
		var61.Description = types.StringValue(ans.Type.Url.Description)
		var61.ExceptionList = EncodeStringSlice(ans.Type.Url.ExceptionList)
		var61.Recurring = var63
		var61.Url = types.StringValue(ans.Type.Url.Url)
	}
	var34.Domain = var35
	var34.Imei = var41
	var34.Imsi = var47
	var34.Ip = var53
	var34.PredefinedIp = var59
	var34.PredefinedUrl = var60
	var34.Url = var61
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Name = types.StringValue(ans.Name)
	state.Type = var34

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *objectsExternalDynamicListsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name": "sase_objects_external_dynamic_lists",
		"locMap":        map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":        tokens,
	})

	svc := iHJqznH.NewClient(r.client)
	input := iHJqznH.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *objectsExternalDynamicListsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
