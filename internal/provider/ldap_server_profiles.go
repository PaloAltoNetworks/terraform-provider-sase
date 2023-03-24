package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	wHHhukY "github.com/paloaltonetworks/sase-go/netsec/schema/ldap/server/profiles"
	iMdQZcj "github.com/paloaltonetworks/sase-go/netsec/service/v1/ldapserverprofiles"

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
	_ datasource.DataSource              = &ldapServerProfilesListDataSource{}
	_ datasource.DataSourceWithConfigure = &ldapServerProfilesListDataSource{}
)

func NewLdapServerProfilesListDataSource() datasource.DataSource {
	return &ldapServerProfilesListDataSource{}
}

type ldapServerProfilesListDataSource struct {
	client *sase.Client
}

type ldapServerProfilesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Folder types.String `tfsdk:"folder"`
	Name   types.String `tfsdk:"name"`

	// Output.
	Data []ldapServerProfilesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type ldapServerProfilesListDsModelConfig struct {
	Base                    types.String                                `tfsdk:"base"`
	BindDn                  types.String                                `tfsdk:"bind_dn"`
	BindPassword            types.String                                `tfsdk:"bind_password"`
	BindTimelimit           types.String                                `tfsdk:"bind_timelimit"`
	ObjectId                types.String                                `tfsdk:"object_id"`
	LdapType                types.String                                `tfsdk:"ldap_type"`
	RetryInterval           types.Int64                                 `tfsdk:"retry_interval"`
	Server                  []ldapServerProfilesListDsModelServerObject `tfsdk:"server"`
	Ssl                     types.Bool                                  `tfsdk:"ssl"`
	Timelimit               types.Int64                                 `tfsdk:"timelimit"`
	VerifyServerCertificate types.Bool                                  `tfsdk:"verify_server_certificate"`
}

type ldapServerProfilesListDsModelServerObject struct {
	Address types.String `tfsdk:"address"`
	Name    types.String `tfsdk:"name"`
	Port    types.Int64  `tfsdk:"port"`
}

// Metadata returns the data source type name.
func (d *ldapServerProfilesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ldap_server_profiles_list"
}

// Schema defines the schema for this listing data source.
func (d *ldapServerProfilesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"name": dsschema.StringAttribute{
				Description:         "The name of the entry",
				MarkdownDescription: "The name of the entry",
				Optional:            true,
			},

			// Output.
			"data": dsschema.ListNestedAttribute{
				Description:         "The `data` parameter.",
				MarkdownDescription: "The `data` parameter.",
				Computed:            true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"base": dsschema.StringAttribute{
							Description:         "The `base` parameter.",
							MarkdownDescription: "The `base` parameter.",
							Computed:            true,
						},
						"bind_dn": dsschema.StringAttribute{
							Description:         "The `bind_dn` parameter.",
							MarkdownDescription: "The `bind_dn` parameter.",
							Computed:            true,
						},
						"bind_password": dsschema.StringAttribute{
							Description:         "The `bind_password` parameter.",
							MarkdownDescription: "The `bind_password` parameter.",
							Computed:            true,
						},
						"bind_timelimit": dsschema.StringAttribute{
							Description:         "The `bind_timelimit` parameter.",
							MarkdownDescription: "The `bind_timelimit` parameter.",
							Computed:            true,
						},
						"object_id": dsschema.StringAttribute{
							Description:         "The `object_id` parameter.",
							MarkdownDescription: "The `object_id` parameter.",
							Computed:            true,
						},
						"ldap_type": dsschema.StringAttribute{
							Description:         "The `ldap_type` parameter.",
							MarkdownDescription: "The `ldap_type` parameter.",
							Computed:            true,
						},
						"retry_interval": dsschema.Int64Attribute{
							Description:         "The `retry_interval` parameter.",
							MarkdownDescription: "The `retry_interval` parameter.",
							Computed:            true,
						},
						"server": dsschema.ListNestedAttribute{
							Description:         "The `server` parameter.",
							MarkdownDescription: "The `server` parameter.",
							Computed:            true,
							NestedObject: dsschema.NestedAttributeObject{
								Attributes: map[string]dsschema.Attribute{
									"address": dsschema.StringAttribute{
										Description:         "The `address` parameter.",
										MarkdownDescription: "The `address` parameter.",
										Computed:            true,
									},
									"name": dsschema.StringAttribute{
										Description:         "The `name` parameter.",
										MarkdownDescription: "The `name` parameter.",
										Computed:            true,
									},
									"port": dsschema.Int64Attribute{
										Description:         "The `port` parameter.",
										MarkdownDescription: "The `port` parameter.",
										Computed:            true,
									},
								},
							},
						},
						"ssl": dsschema.BoolAttribute{
							Description:         "The `ssl` parameter.",
							MarkdownDescription: "The `ssl` parameter.",
							Computed:            true,
						},
						"timelimit": dsschema.Int64Attribute{
							Description:         "The `timelimit` parameter.",
							MarkdownDescription: "The `timelimit` parameter.",
							Computed:            true,
						},
						"verify_server_certificate": dsschema.BoolAttribute{
							Description:         "The `verify_server_certificate` parameter.",
							MarkdownDescription: "The `verify_server_certificate` parameter.",
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
func (d *ldapServerProfilesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *ldapServerProfilesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ldapServerProfilesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_ldap_server_profiles_list",
		"terraform_provider_function": "Read",
		"limit":                       state.Limit.ValueInt64(),
		"has_limit":                   !state.Limit.IsNull(),
		"offset":                      state.Offset.ValueInt64(),
		"has_offset":                  !state.Offset.IsNull(),
		"folder":                      state.Folder.ValueString(),
		"name":                        state.Name.ValueString(),
		"has_name":                    !state.Name.IsNull(),
	})

	// Prepare to run the command.
	svc := iMdQZcj.NewClient(d.client)
	input := iMdQZcj.ListInput{
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
	idBuilder.WriteString(input.Folder)
	idBuilder.WriteString(IdSeparator)
	if input.Name != nil {
		idBuilder.WriteString(*input.Name)
	}
	state.Id = types.StringValue(idBuilder.String())
	var var0 []ldapServerProfilesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]ldapServerProfilesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 ldapServerProfilesListDsModelConfig
			var var3 []ldapServerProfilesListDsModelServerObject
			if len(var1.Server) != 0 {
				var3 = make([]ldapServerProfilesListDsModelServerObject, 0, len(var1.Server))
				for var4Index := range var1.Server {
					var4 := var1.Server[var4Index]
					var var5 ldapServerProfilesListDsModelServerObject
					var5.Address = types.StringValue(var4.Address)
					var5.Name = types.StringValue(var4.Name)
					var5.Port = types.Int64Value(var4.Port)
					var3 = append(var3, var5)
				}
			}
			var2.Base = types.StringValue(var1.Base)
			var2.BindDn = types.StringValue(var1.BindDn)
			var2.BindPassword = types.StringValue(var1.BindPassword)
			var2.BindTimelimit = types.StringValue(var1.BindTimelimit)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.LdapType = types.StringValue(var1.LdapType)
			var2.RetryInterval = types.Int64Value(var1.RetryInterval)
			var2.Server = var3
			var2.Ssl = types.BoolValue(var1.Ssl)
			var2.Timelimit = types.Int64Value(var1.Timelimit)
			var2.VerifyServerCertificate = types.BoolValue(var1.VerifyServerCertificate)
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
	_ datasource.DataSource              = &ldapServerProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &ldapServerProfilesDataSource{}
)

func NewLdapServerProfilesDataSource() datasource.DataSource {
	return &ldapServerProfilesDataSource{}
}

type ldapServerProfilesDataSource struct {
	client *sase.Client
}

type ldapServerProfilesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`
	Folder   types.String `tfsdk:"folder"`

	// Output.
	// Ref: #/components/schemas/ldap-server-profiles
	Base          types.String `tfsdk:"base"`
	BindDn        types.String `tfsdk:"bind_dn"`
	BindPassword  types.String `tfsdk:"bind_password"`
	BindTimelimit types.String `tfsdk:"bind_timelimit"`
	// input omit: ObjectId
	LdapType                types.String                            `tfsdk:"ldap_type"`
	RetryInterval           types.Int64                             `tfsdk:"retry_interval"`
	Server                  []ldapServerProfilesDsModelServerObject `tfsdk:"server"`
	Ssl                     types.Bool                              `tfsdk:"ssl"`
	Timelimit               types.Int64                             `tfsdk:"timelimit"`
	VerifyServerCertificate types.Bool                              `tfsdk:"verify_server_certificate"`
}

type ldapServerProfilesDsModelServerObject struct {
	Address types.String `tfsdk:"address"`
	Name    types.String `tfsdk:"name"`
	Port    types.Int64  `tfsdk:"port"`
}

// Metadata returns the data source type name.
func (d *ldapServerProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ldap_server_profiles"
}

// Schema defines the schema for this listing data source.
func (d *ldapServerProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"base": dsschema.StringAttribute{
				Description:         "The `base` parameter.",
				MarkdownDescription: "The `base` parameter.",
				Computed:            true,
			},
			"bind_dn": dsschema.StringAttribute{
				Description:         "The `bind_dn` parameter.",
				MarkdownDescription: "The `bind_dn` parameter.",
				Computed:            true,
			},
			"bind_password": dsschema.StringAttribute{
				Description:         "The `bind_password` parameter.",
				MarkdownDescription: "The `bind_password` parameter.",
				Computed:            true,
			},
			"bind_timelimit": dsschema.StringAttribute{
				Description:         "The `bind_timelimit` parameter.",
				MarkdownDescription: "The `bind_timelimit` parameter.",
				Computed:            true,
			},
			"ldap_type": dsschema.StringAttribute{
				Description:         "The `ldap_type` parameter.",
				MarkdownDescription: "The `ldap_type` parameter.",
				Computed:            true,
			},
			"retry_interval": dsschema.Int64Attribute{
				Description:         "The `retry_interval` parameter.",
				MarkdownDescription: "The `retry_interval` parameter.",
				Computed:            true,
			},
			"server": dsschema.ListNestedAttribute{
				Description:         "The `server` parameter.",
				MarkdownDescription: "The `server` parameter.",
				Computed:            true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"address": dsschema.StringAttribute{
							Description:         "The `address` parameter.",
							MarkdownDescription: "The `address` parameter.",
							Computed:            true,
						},
						"name": dsschema.StringAttribute{
							Description:         "The `name` parameter.",
							MarkdownDescription: "The `name` parameter.",
							Computed:            true,
						},
						"port": dsschema.Int64Attribute{
							Description:         "The `port` parameter.",
							MarkdownDescription: "The `port` parameter.",
							Computed:            true,
						},
					},
				},
			},
			"ssl": dsschema.BoolAttribute{
				Description:         "The `ssl` parameter.",
				MarkdownDescription: "The `ssl` parameter.",
				Computed:            true,
			},
			"timelimit": dsschema.Int64Attribute{
				Description:         "The `timelimit` parameter.",
				MarkdownDescription: "The `timelimit` parameter.",
				Computed:            true,
			},
			"verify_server_certificate": dsschema.BoolAttribute{
				Description:         "The `verify_server_certificate` parameter.",
				MarkdownDescription: "The `verify_server_certificate` parameter.",
				Computed:            true,
			},
		},
	}
}

// Configure prepares the struct.
func (d *ldapServerProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *ldapServerProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ldapServerProfilesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_ldap_server_profiles",
		"object_id":                   state.ObjectId.ValueString(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := iMdQZcj.NewClient(d.client)
	input := iMdQZcj.ReadInput{
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
	var var0 []ldapServerProfilesDsModelServerObject
	if len(ans.Server) != 0 {
		var0 = make([]ldapServerProfilesDsModelServerObject, 0, len(ans.Server))
		for var1Index := range ans.Server {
			var1 := ans.Server[var1Index]
			var var2 ldapServerProfilesDsModelServerObject
			var2.Address = types.StringValue(var1.Address)
			var2.Name = types.StringValue(var1.Name)
			var2.Port = types.Int64Value(var1.Port)
			var0 = append(var0, var2)
		}
	}
	state.Base = types.StringValue(ans.Base)
	state.BindDn = types.StringValue(ans.BindDn)
	state.BindPassword = types.StringValue(ans.BindPassword)
	state.BindTimelimit = types.StringValue(ans.BindTimelimit)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.LdapType = types.StringValue(ans.LdapType)
	state.RetryInterval = types.Int64Value(ans.RetryInterval)
	state.Server = var0
	state.Ssl = types.BoolValue(ans.Ssl)
	state.Timelimit = types.Int64Value(ans.Timelimit)
	state.VerifyServerCertificate = types.BoolValue(ans.VerifyServerCertificate)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &ldapServerProfilesResource{}
	_ resource.ResourceWithConfigure   = &ldapServerProfilesResource{}
	_ resource.ResourceWithImportState = &ldapServerProfilesResource{}
)

func NewLdapServerProfilesResource() resource.Resource {
	return &ldapServerProfilesResource{}
}

type ldapServerProfilesResource struct {
	client *sase.Client
}

type ldapServerProfilesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/ldap-server-profiles
	Base                    types.String                            `tfsdk:"base"`
	BindDn                  types.String                            `tfsdk:"bind_dn"`
	BindPassword            types.String                            `tfsdk:"bind_password"`
	BindTimelimit           types.String                            `tfsdk:"bind_timelimit"`
	ObjectId                types.String                            `tfsdk:"object_id"`
	LdapType                types.String                            `tfsdk:"ldap_type"`
	RetryInterval           types.Int64                             `tfsdk:"retry_interval"`
	Server                  []ldapServerProfilesRsModelServerObject `tfsdk:"server"`
	Ssl                     types.Bool                              `tfsdk:"ssl"`
	Timelimit               types.Int64                             `tfsdk:"timelimit"`
	VerifyServerCertificate types.Bool                              `tfsdk:"verify_server_certificate"`
}

type ldapServerProfilesRsModelServerObject struct {
	Address types.String `tfsdk:"address"`
	Name    types.String `tfsdk:"name"`
	Port    types.Int64  `tfsdk:"port"`
}

// Metadata returns the data source type name.
func (r *ldapServerProfilesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ldap_server_profiles"
}

// Schema defines the schema for this listing data source.
func (r *ldapServerProfilesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

			"base": rsschema.StringAttribute{
				Description:         "The `base` parameter.",
				MarkdownDescription: "The `base` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
			},
			"bind_dn": rsschema.StringAttribute{
				Description:         "The `bind_dn` parameter.",
				MarkdownDescription: "The `bind_dn` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
			},
			"bind_password": rsschema.StringAttribute{
				Description:         "The `bind_password` parameter.",
				MarkdownDescription: "The `bind_password` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(121),
				},
			},
			"bind_timelimit": rsschema.StringAttribute{
				Description:         "The `bind_timelimit` parameter.",
				MarkdownDescription: "The `bind_timelimit` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
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
			"ldap_type": rsschema.StringAttribute{
				Description:         "The `ldap_type` parameter.",
				MarkdownDescription: "The `ldap_type` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("active-directory", "e-directory", "sun", "other"),
				},
			},
			"retry_interval": rsschema.Int64Attribute{
				Description:         "The `retry_interval` parameter.",
				MarkdownDescription: "The `retry_interval` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					DefaultInt64(0),
				},
			},
			"server": rsschema.ListNestedAttribute{
				Description:         "The `server` parameter.",
				MarkdownDescription: "The `server` parameter.",
				Required:            true,
				NestedObject: rsschema.NestedAttributeObject{
					Attributes: map[string]rsschema.Attribute{
						"address": rsschema.StringAttribute{
							Description:         "The `address` parameter.",
							MarkdownDescription: "The `address` parameter.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
						},
						"name": rsschema.StringAttribute{
							Description:         "The `name` parameter.",
							MarkdownDescription: "The `name` parameter.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
						},
						"port": rsschema.Int64Attribute{
							Description:         "The `port` parameter.",
							MarkdownDescription: "The `port` parameter.",
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
			},
			"ssl": rsschema.BoolAttribute{
				Description:         "The `ssl` parameter.",
				MarkdownDescription: "The `ssl` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"timelimit": rsschema.Int64Attribute{
				Description:         "The `timelimit` parameter.",
				MarkdownDescription: "The `timelimit` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					DefaultInt64(0),
				},
			},
			"verify_server_certificate": rsschema.BoolAttribute{
				Description:         "The `verify_server_certificate` parameter.",
				MarkdownDescription: "The `verify_server_certificate` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *ldapServerProfilesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *ldapServerProfilesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state ldapServerProfilesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_ldap_server_profiles",
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := iMdQZcj.NewClient(r.client)
	input := iMdQZcj.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 wHHhukY.Config
	var0.Base = state.Base.ValueString()
	var0.BindDn = state.BindDn.ValueString()
	var0.BindPassword = state.BindPassword.ValueString()
	var0.BindTimelimit = state.BindTimelimit.ValueString()
	var0.LdapType = state.LdapType.ValueString()
	var0.RetryInterval = state.RetryInterval.ValueInt64()
	var var1 []wHHhukY.ServerObject
	if len(state.Server) != 0 {
		var1 = make([]wHHhukY.ServerObject, 0, len(state.Server))
		for var2Index := range state.Server {
			var2 := state.Server[var2Index]
			var var3 wHHhukY.ServerObject
			var3.Address = var2.Address.ValueString()
			var3.Name = var2.Name.ValueString()
			var3.Port = var2.Port.ValueInt64()
			var1 = append(var1, var3)
		}
	}
	var0.Server = var1
	var0.Ssl = state.Ssl.ValueBool()
	var0.Timelimit = state.Timelimit.ValueInt64()
	var0.VerifyServerCertificate = state.VerifyServerCertificate.ValueBool()
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
	var var4 []ldapServerProfilesRsModelServerObject
	if len(ans.Server) != 0 {
		var4 = make([]ldapServerProfilesRsModelServerObject, 0, len(ans.Server))
		for var5Index := range ans.Server {
			var5 := ans.Server[var5Index]
			var var6 ldapServerProfilesRsModelServerObject
			var6.Address = types.StringValue(var5.Address)
			var6.Name = types.StringValue(var5.Name)
			var6.Port = types.Int64Value(var5.Port)
			var4 = append(var4, var6)
		}
	}
	state.Base = types.StringValue(ans.Base)
	state.BindDn = types.StringValue(ans.BindDn)
	state.BindPassword = types.StringValue(ans.BindPassword)
	state.BindTimelimit = types.StringValue(ans.BindTimelimit)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.LdapType = types.StringValue(ans.LdapType)
	state.RetryInterval = types.Int64Value(ans.RetryInterval)
	state.Server = var4
	state.Ssl = types.BoolValue(ans.Ssl)
	state.Timelimit = types.Int64Value(ans.Timelimit)
	state.VerifyServerCertificate = types.BoolValue(ans.VerifyServerCertificate)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *ldapServerProfilesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state ldapServerProfilesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_ldap_server_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := iMdQZcj.NewClient(r.client)
	input := iMdQZcj.ReadInput{
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
	var var0 []ldapServerProfilesRsModelServerObject
	if len(ans.Server) != 0 {
		var0 = make([]ldapServerProfilesRsModelServerObject, 0, len(ans.Server))
		for var1Index := range ans.Server {
			var1 := ans.Server[var1Index]
			var var2 ldapServerProfilesRsModelServerObject
			var2.Address = types.StringValue(var1.Address)
			var2.Name = types.StringValue(var1.Name)
			var2.Port = types.Int64Value(var1.Port)
			var0 = append(var0, var2)
		}
	}
	state.Base = types.StringValue(ans.Base)
	state.BindDn = types.StringValue(ans.BindDn)
	state.BindPassword = types.StringValue(ans.BindPassword)
	state.BindTimelimit = types.StringValue(ans.BindTimelimit)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.LdapType = types.StringValue(ans.LdapType)
	state.RetryInterval = types.Int64Value(ans.RetryInterval)
	state.Server = var0
	state.Ssl = types.BoolValue(ans.Ssl)
	state.Timelimit = types.Int64Value(ans.Timelimit)
	state.VerifyServerCertificate = types.BoolValue(ans.VerifyServerCertificate)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *ldapServerProfilesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ldapServerProfilesRsModel
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
		"resource_name":               "sase_ldap_server_profiles",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := iMdQZcj.NewClient(r.client)
	input := iMdQZcj.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 wHHhukY.Config
	var0.Base = plan.Base.ValueString()
	var0.BindDn = plan.BindDn.ValueString()
	var0.BindPassword = plan.BindPassword.ValueString()
	var0.BindTimelimit = plan.BindTimelimit.ValueString()
	var0.LdapType = plan.LdapType.ValueString()
	var0.RetryInterval = plan.RetryInterval.ValueInt64()
	var var1 []wHHhukY.ServerObject
	if len(plan.Server) != 0 {
		var1 = make([]wHHhukY.ServerObject, 0, len(plan.Server))
		for var2Index := range plan.Server {
			var2 := plan.Server[var2Index]
			var var3 wHHhukY.ServerObject
			var3.Address = var2.Address.ValueString()
			var3.Name = var2.Name.ValueString()
			var3.Port = var2.Port.ValueInt64()
			var1 = append(var1, var3)
		}
	}
	var0.Server = var1
	var0.Ssl = plan.Ssl.ValueBool()
	var0.Timelimit = plan.Timelimit.ValueInt64()
	var0.VerifyServerCertificate = plan.VerifyServerCertificate.ValueBool()
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var4 []ldapServerProfilesRsModelServerObject
	if len(ans.Server) != 0 {
		var4 = make([]ldapServerProfilesRsModelServerObject, 0, len(ans.Server))
		for var5Index := range ans.Server {
			var5 := ans.Server[var5Index]
			var var6 ldapServerProfilesRsModelServerObject
			var6.Address = types.StringValue(var5.Address)
			var6.Name = types.StringValue(var5.Name)
			var6.Port = types.Int64Value(var5.Port)
			var4 = append(var4, var6)
		}
	}
	state.Base = types.StringValue(ans.Base)
	state.BindDn = types.StringValue(ans.BindDn)
	state.BindPassword = types.StringValue(ans.BindPassword)
	state.BindTimelimit = types.StringValue(ans.BindTimelimit)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.LdapType = types.StringValue(ans.LdapType)
	state.RetryInterval = types.Int64Value(ans.RetryInterval)
	state.Server = var4
	state.Ssl = types.BoolValue(ans.Ssl)
	state.Timelimit = types.Int64Value(ans.Timelimit)
	state.VerifyServerCertificate = types.BoolValue(ans.VerifyServerCertificate)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *ldapServerProfilesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name":               "sase_ldap_server_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	svc := iMdQZcj.NewClient(r.client)
	input := iMdQZcj.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *ldapServerProfilesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
