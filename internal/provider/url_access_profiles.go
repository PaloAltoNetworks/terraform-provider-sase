package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	wleovFf "github.com/paloaltonetworks/sase-go/netsec/schema/url/access/profiles"
	uyrOkzA "github.com/paloaltonetworks/sase-go/netsec/service/v1/urlaccessprofiles"

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
	_ datasource.DataSource              = &urlAccessProfilesListDataSource{}
	_ datasource.DataSourceWithConfigure = &urlAccessProfilesListDataSource{}
)

func NewUrlAccessProfilesListDataSource() datasource.DataSource {
	return &urlAccessProfilesListDataSource{}
}

type urlAccessProfilesListDataSource struct {
	client *sase.Client
}

type urlAccessProfilesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []urlAccessProfilesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type urlAccessProfilesListDsModelConfig struct {
	Alert                     []types.String                                                `tfsdk:"alert"`
	Allow                     []types.String                                                `tfsdk:"allow"`
	Block                     []types.String                                                `tfsdk:"block"`
	Continue                  []types.String                                                `tfsdk:"continue"`
	CredentialEnforcement     *urlAccessProfilesListDsModelCredentialEnforcementObject      `tfsdk:"credential_enforcement"`
	Description               types.String                                                  `tfsdk:"description"`
	ObjectId                  types.String                                                  `tfsdk:"object_id"`
	LogContainerPageOnly      types.Bool                                                    `tfsdk:"log_container_page_only"`
	LogHttpHdrReferer         types.Bool                                                    `tfsdk:"log_http_hdr_referer"`
	LogHttpHdrUserAgent       types.Bool                                                    `tfsdk:"log_http_hdr_user_agent"`
	LogHttpHdrXff             types.Bool                                                    `tfsdk:"log_http_hdr_xff"`
	MlavCategoryException     []types.String                                                `tfsdk:"mlav_category_exception"`
	MlavEngineUrlbasedEnabled []urlAccessProfilesListDsModelMlavEngineUrlbasedEnabledObject `tfsdk:"mlav_engine_urlbased_enabled"`
	Name                      types.String                                                  `tfsdk:"name"`
	SafeSearchEnforcement     types.Bool                                                    `tfsdk:"safe_search_enforcement"`
}

type urlAccessProfilesListDsModelCredentialEnforcementObject struct {
	Alert       []types.String                          `tfsdk:"alert"`
	Allow       []types.String                          `tfsdk:"allow"`
	Block       []types.String                          `tfsdk:"block"`
	Continue    []types.String                          `tfsdk:"continue"`
	LogSeverity types.String                            `tfsdk:"log_severity"`
	Mode        *urlAccessProfilesListDsModelModeObject `tfsdk:"mode"`
}

type urlAccessProfilesListDsModelModeObject struct {
	Disabled          types.Bool   `tfsdk:"disabled"`
	DomainCredentials types.Bool   `tfsdk:"domain_credentials"`
	GroupMapping      types.String `tfsdk:"group_mapping"`
	IpUser            types.Bool   `tfsdk:"ip_user"`
}

type urlAccessProfilesListDsModelMlavEngineUrlbasedEnabledObject struct {
	MlavPolicyAction types.String `tfsdk:"mlav_policy_action"`
	Name             types.String `tfsdk:"name"`
}

// Metadata returns the data source type name.
func (d *urlAccessProfilesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_url_access_profiles_list"
}

// Schema defines the schema for this listing data source.
func (d *urlAccessProfilesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"alert": dsschema.ListAttribute{
							Description:         "The `alert` parameter.",
							MarkdownDescription: "The `alert` parameter.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"allow": dsschema.ListAttribute{
							Description:         "The `allow` parameter.",
							MarkdownDescription: "The `allow` parameter.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"block": dsschema.ListAttribute{
							Description:         "The `block` parameter.",
							MarkdownDescription: "The `block` parameter.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"continue": dsschema.ListAttribute{
							Description:         "The `continue` parameter.",
							MarkdownDescription: "The `continue` parameter.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"credential_enforcement": dsschema.SingleNestedAttribute{
							Description:         "The `credential_enforcement` parameter.",
							MarkdownDescription: "The `credential_enforcement` parameter.",
							Computed:            true,
							Attributes: map[string]dsschema.Attribute{
								"alert": dsschema.ListAttribute{
									Description:         "The `alert` parameter.",
									MarkdownDescription: "The `alert` parameter.",
									Computed:            true,
									ElementType:         types.StringType,
								},
								"allow": dsschema.ListAttribute{
									Description:         "The `allow` parameter.",
									MarkdownDescription: "The `allow` parameter.",
									Computed:            true,
									ElementType:         types.StringType,
								},
								"block": dsschema.ListAttribute{
									Description:         "The `block` parameter.",
									MarkdownDescription: "The `block` parameter.",
									Computed:            true,
									ElementType:         types.StringType,
								},
								"continue": dsschema.ListAttribute{
									Description:         "The `continue` parameter.",
									MarkdownDescription: "The `continue` parameter.",
									Computed:            true,
									ElementType:         types.StringType,
								},
								"log_severity": dsschema.StringAttribute{
									Description:         "The `log_severity` parameter.",
									MarkdownDescription: "The `log_severity` parameter.",
									Computed:            true,
								},
								"mode": dsschema.SingleNestedAttribute{
									Description:         "The `mode` parameter.",
									MarkdownDescription: "The `mode` parameter.",
									Computed:            true,
									Attributes: map[string]dsschema.Attribute{
										"disabled": dsschema.BoolAttribute{
											Description:         "The `disabled` parameter.",
											MarkdownDescription: "The `disabled` parameter.",
											Computed:            true,
										},
										"domain_credentials": dsschema.BoolAttribute{
											Description:         "The `domain_credentials` parameter.",
											MarkdownDescription: "The `domain_credentials` parameter.",
											Computed:            true,
										},
										"group_mapping": dsschema.StringAttribute{
											Description:         "The `group_mapping` parameter.",
											MarkdownDescription: "The `group_mapping` parameter.",
											Computed:            true,
										},
										"ip_user": dsschema.BoolAttribute{
											Description:         "The `ip_user` parameter.",
											MarkdownDescription: "The `ip_user` parameter.",
											Computed:            true,
										},
									},
								},
							},
						},
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
						"log_container_page_only": dsschema.BoolAttribute{
							Description:         "The `log_container_page_only` parameter.",
							MarkdownDescription: "The `log_container_page_only` parameter.",
							Computed:            true,
						},
						"log_http_hdr_referer": dsschema.BoolAttribute{
							Description:         "The `log_http_hdr_referer` parameter.",
							MarkdownDescription: "The `log_http_hdr_referer` parameter.",
							Computed:            true,
						},
						"log_http_hdr_user_agent": dsschema.BoolAttribute{
							Description:         "The `log_http_hdr_user_agent` parameter.",
							MarkdownDescription: "The `log_http_hdr_user_agent` parameter.",
							Computed:            true,
						},
						"log_http_hdr_xff": dsschema.BoolAttribute{
							Description:         "The `log_http_hdr_xff` parameter.",
							MarkdownDescription: "The `log_http_hdr_xff` parameter.",
							Computed:            true,
						},
						"mlav_category_exception": dsschema.ListAttribute{
							Description:         "The `mlav_category_exception` parameter.",
							MarkdownDescription: "The `mlav_category_exception` parameter.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"mlav_engine_urlbased_enabled": dsschema.ListNestedAttribute{
							Description:         "The `mlav_engine_urlbased_enabled` parameter.",
							MarkdownDescription: "The `mlav_engine_urlbased_enabled` parameter.",
							Computed:            true,
							NestedObject: dsschema.NestedAttributeObject{
								Attributes: map[string]dsschema.Attribute{
									"mlav_policy_action": dsschema.StringAttribute{
										Description:         "The `mlav_policy_action` parameter.",
										MarkdownDescription: "The `mlav_policy_action` parameter.",
										Computed:            true,
									},
									"name": dsschema.StringAttribute{
										Description:         "The `name` parameter.",
										MarkdownDescription: "The `name` parameter.",
										Computed:            true,
									},
								},
							},
						},
						"name": dsschema.StringAttribute{
							Description:         "The `name` parameter.",
							MarkdownDescription: "The `name` parameter.",
							Computed:            true,
						},
						"safe_search_enforcement": dsschema.BoolAttribute{
							Description:         "The `safe_search_enforcement` parameter.",
							MarkdownDescription: "The `safe_search_enforcement` parameter.",
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
func (d *urlAccessProfilesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *urlAccessProfilesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state urlAccessProfilesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_url_access_profiles_list",
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
	svc := uyrOkzA.NewClient(d.client)
	input := uyrOkzA.ListInput{
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
	var var0 []urlAccessProfilesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]urlAccessProfilesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 urlAccessProfilesListDsModelConfig
			var var3 *urlAccessProfilesListDsModelCredentialEnforcementObject
			if var1.CredentialEnforcement != nil {
				var3 = &urlAccessProfilesListDsModelCredentialEnforcementObject{}
				var var4 *urlAccessProfilesListDsModelModeObject
				if var1.CredentialEnforcement.Mode != nil {
					var4 = &urlAccessProfilesListDsModelModeObject{}
					if var1.CredentialEnforcement.Mode.Disabled != nil {
						var4.Disabled = types.BoolValue(true)
					}
					if var1.CredentialEnforcement.Mode.DomainCredentials != nil {
						var4.DomainCredentials = types.BoolValue(true)
					}
					var4.GroupMapping = types.StringValue(var1.CredentialEnforcement.Mode.GroupMapping)
					if var1.CredentialEnforcement.Mode.IpUser != nil {
						var4.IpUser = types.BoolValue(true)
					}
				}
				var3.Alert = EncodeStringSlice(var1.CredentialEnforcement.Alert)
				var3.Allow = EncodeStringSlice(var1.CredentialEnforcement.Allow)
				var3.Block = EncodeStringSlice(var1.CredentialEnforcement.Block)
				var3.Continue = EncodeStringSlice(var1.CredentialEnforcement.Continue)
				var3.LogSeverity = types.StringValue(var1.CredentialEnforcement.LogSeverity)
				var3.Mode = var4
			}
			var var5 []urlAccessProfilesListDsModelMlavEngineUrlbasedEnabledObject
			if len(var1.MlavEngineUrlbasedEnabled) != 0 {
				var5 = make([]urlAccessProfilesListDsModelMlavEngineUrlbasedEnabledObject, 0, len(var1.MlavEngineUrlbasedEnabled))
				for var6Index := range var1.MlavEngineUrlbasedEnabled {
					var6 := var1.MlavEngineUrlbasedEnabled[var6Index]
					var var7 urlAccessProfilesListDsModelMlavEngineUrlbasedEnabledObject
					var7.MlavPolicyAction = types.StringValue(var6.MlavPolicyAction)
					var7.Name = types.StringValue(var6.Name)
					var5 = append(var5, var7)
				}
			}
			var2.Alert = EncodeStringSlice(var1.Alert)
			var2.Allow = EncodeStringSlice(var1.Allow)
			var2.Block = EncodeStringSlice(var1.Block)
			var2.Continue = EncodeStringSlice(var1.Continue)
			var2.CredentialEnforcement = var3
			var2.Description = types.StringValue(var1.Description)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.LogContainerPageOnly = types.BoolValue(var1.LogContainerPageOnly)
			var2.LogHttpHdrReferer = types.BoolValue(var1.LogHttpHdrReferer)
			var2.LogHttpHdrUserAgent = types.BoolValue(var1.LogHttpHdrUserAgent)
			var2.LogHttpHdrXff = types.BoolValue(var1.LogHttpHdrXff)
			var2.MlavCategoryException = EncodeStringSlice(var1.MlavCategoryException)
			var2.MlavEngineUrlbasedEnabled = var5
			var2.Name = types.StringValue(var1.Name)
			var2.SafeSearchEnforcement = types.BoolValue(var1.SafeSearchEnforcement)
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
	_ datasource.DataSource              = &urlAccessProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &urlAccessProfilesDataSource{}
)

func NewUrlAccessProfilesDataSource() datasource.DataSource {
	return &urlAccessProfilesDataSource{}
}

type urlAccessProfilesDataSource struct {
	client *sase.Client
}

type urlAccessProfilesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`
	Folder   types.String `tfsdk:"folder"`

	// Output.
	// Ref: #/components/schemas/url-access-profiles
	Alert                 []types.String                                       `tfsdk:"alert"`
	Allow                 []types.String                                       `tfsdk:"allow"`
	Block                 []types.String                                       `tfsdk:"block"`
	Continue              []types.String                                       `tfsdk:"continue"`
	CredentialEnforcement *urlAccessProfilesDsModelCredentialEnforcementObject `tfsdk:"credential_enforcement"`
	Description           types.String                                         `tfsdk:"description"`
	// input omit: ObjectId
	LogContainerPageOnly      types.Bool                                                `tfsdk:"log_container_page_only"`
	LogHttpHdrReferer         types.Bool                                                `tfsdk:"log_http_hdr_referer"`
	LogHttpHdrUserAgent       types.Bool                                                `tfsdk:"log_http_hdr_user_agent"`
	LogHttpHdrXff             types.Bool                                                `tfsdk:"log_http_hdr_xff"`
	MlavCategoryException     []types.String                                            `tfsdk:"mlav_category_exception"`
	MlavEngineUrlbasedEnabled []urlAccessProfilesDsModelMlavEngineUrlbasedEnabledObject `tfsdk:"mlav_engine_urlbased_enabled"`
	Name                      types.String                                              `tfsdk:"name"`
	SafeSearchEnforcement     types.Bool                                                `tfsdk:"safe_search_enforcement"`
}

type urlAccessProfilesDsModelCredentialEnforcementObject struct {
	Alert       []types.String                      `tfsdk:"alert"`
	Allow       []types.String                      `tfsdk:"allow"`
	Block       []types.String                      `tfsdk:"block"`
	Continue    []types.String                      `tfsdk:"continue"`
	LogSeverity types.String                        `tfsdk:"log_severity"`
	Mode        *urlAccessProfilesDsModelModeObject `tfsdk:"mode"`
}

type urlAccessProfilesDsModelModeObject struct {
	Disabled          types.Bool   `tfsdk:"disabled"`
	DomainCredentials types.Bool   `tfsdk:"domain_credentials"`
	GroupMapping      types.String `tfsdk:"group_mapping"`
	IpUser            types.Bool   `tfsdk:"ip_user"`
}

type urlAccessProfilesDsModelMlavEngineUrlbasedEnabledObject struct {
	MlavPolicyAction types.String `tfsdk:"mlav_policy_action"`
	Name             types.String `tfsdk:"name"`
}

// Metadata returns the data source type name.
func (d *urlAccessProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_url_access_profiles"
}

// Schema defines the schema for this listing data source.
func (d *urlAccessProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"alert": dsschema.ListAttribute{
				Description:         "The `alert` parameter.",
				MarkdownDescription: "The `alert` parameter.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"allow": dsschema.ListAttribute{
				Description:         "The `allow` parameter.",
				MarkdownDescription: "The `allow` parameter.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"block": dsschema.ListAttribute{
				Description:         "The `block` parameter.",
				MarkdownDescription: "The `block` parameter.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"continue": dsschema.ListAttribute{
				Description:         "The `continue` parameter.",
				MarkdownDescription: "The `continue` parameter.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"credential_enforcement": dsschema.SingleNestedAttribute{
				Description:         "The `credential_enforcement` parameter.",
				MarkdownDescription: "The `credential_enforcement` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"alert": dsschema.ListAttribute{
						Description:         "The `alert` parameter.",
						MarkdownDescription: "The `alert` parameter.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"allow": dsschema.ListAttribute{
						Description:         "The `allow` parameter.",
						MarkdownDescription: "The `allow` parameter.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"block": dsschema.ListAttribute{
						Description:         "The `block` parameter.",
						MarkdownDescription: "The `block` parameter.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"continue": dsschema.ListAttribute{
						Description:         "The `continue` parameter.",
						MarkdownDescription: "The `continue` parameter.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"log_severity": dsschema.StringAttribute{
						Description:         "The `log_severity` parameter.",
						MarkdownDescription: "The `log_severity` parameter.",
						Computed:            true,
					},
					"mode": dsschema.SingleNestedAttribute{
						Description:         "The `mode` parameter.",
						MarkdownDescription: "The `mode` parameter.",
						Computed:            true,
						Attributes: map[string]dsschema.Attribute{
							"disabled": dsschema.BoolAttribute{
								Description:         "The `disabled` parameter.",
								MarkdownDescription: "The `disabled` parameter.",
								Computed:            true,
							},
							"domain_credentials": dsschema.BoolAttribute{
								Description:         "The `domain_credentials` parameter.",
								MarkdownDescription: "The `domain_credentials` parameter.",
								Computed:            true,
							},
							"group_mapping": dsschema.StringAttribute{
								Description:         "The `group_mapping` parameter.",
								MarkdownDescription: "The `group_mapping` parameter.",
								Computed:            true,
							},
							"ip_user": dsschema.BoolAttribute{
								Description:         "The `ip_user` parameter.",
								MarkdownDescription: "The `ip_user` parameter.",
								Computed:            true,
							},
						},
					},
				},
			},
			"description": dsschema.StringAttribute{
				Description:         "The `description` parameter.",
				MarkdownDescription: "The `description` parameter.",
				Computed:            true,
			},
			"log_container_page_only": dsschema.BoolAttribute{
				Description:         "The `log_container_page_only` parameter.",
				MarkdownDescription: "The `log_container_page_only` parameter.",
				Computed:            true,
			},
			"log_http_hdr_referer": dsschema.BoolAttribute{
				Description:         "The `log_http_hdr_referer` parameter.",
				MarkdownDescription: "The `log_http_hdr_referer` parameter.",
				Computed:            true,
			},
			"log_http_hdr_user_agent": dsschema.BoolAttribute{
				Description:         "The `log_http_hdr_user_agent` parameter.",
				MarkdownDescription: "The `log_http_hdr_user_agent` parameter.",
				Computed:            true,
			},
			"log_http_hdr_xff": dsschema.BoolAttribute{
				Description:         "The `log_http_hdr_xff` parameter.",
				MarkdownDescription: "The `log_http_hdr_xff` parameter.",
				Computed:            true,
			},
			"mlav_category_exception": dsschema.ListAttribute{
				Description:         "The `mlav_category_exception` parameter.",
				MarkdownDescription: "The `mlav_category_exception` parameter.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"mlav_engine_urlbased_enabled": dsschema.ListNestedAttribute{
				Description:         "The `mlav_engine_urlbased_enabled` parameter.",
				MarkdownDescription: "The `mlav_engine_urlbased_enabled` parameter.",
				Computed:            true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"mlav_policy_action": dsschema.StringAttribute{
							Description:         "The `mlav_policy_action` parameter.",
							MarkdownDescription: "The `mlav_policy_action` parameter.",
							Computed:            true,
						},
						"name": dsschema.StringAttribute{
							Description:         "The `name` parameter.",
							MarkdownDescription: "The `name` parameter.",
							Computed:            true,
						},
					},
				},
			},
			"name": dsschema.StringAttribute{
				Description:         "The `name` parameter.",
				MarkdownDescription: "The `name` parameter.",
				Computed:            true,
			},
			"safe_search_enforcement": dsschema.BoolAttribute{
				Description:         "The `safe_search_enforcement` parameter.",
				MarkdownDescription: "The `safe_search_enforcement` parameter.",
				Computed:            true,
			},
		},
	}
}

// Configure prepares the struct.
func (d *urlAccessProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *urlAccessProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state urlAccessProfilesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_url_access_profiles",
		"object_id":                   state.ObjectId.ValueString(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := uyrOkzA.NewClient(d.client)
	input := uyrOkzA.ReadInput{
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
	var var0 *urlAccessProfilesDsModelCredentialEnforcementObject
	if ans.CredentialEnforcement != nil {
		var0 = &urlAccessProfilesDsModelCredentialEnforcementObject{}
		var var1 *urlAccessProfilesDsModelModeObject
		if ans.CredentialEnforcement.Mode != nil {
			var1 = &urlAccessProfilesDsModelModeObject{}
			if ans.CredentialEnforcement.Mode.Disabled != nil {
				var1.Disabled = types.BoolValue(true)
			}
			if ans.CredentialEnforcement.Mode.DomainCredentials != nil {
				var1.DomainCredentials = types.BoolValue(true)
			}
			var1.GroupMapping = types.StringValue(ans.CredentialEnforcement.Mode.GroupMapping)
			if ans.CredentialEnforcement.Mode.IpUser != nil {
				var1.IpUser = types.BoolValue(true)
			}
		}
		var0.Alert = EncodeStringSlice(ans.CredentialEnforcement.Alert)
		var0.Allow = EncodeStringSlice(ans.CredentialEnforcement.Allow)
		var0.Block = EncodeStringSlice(ans.CredentialEnforcement.Block)
		var0.Continue = EncodeStringSlice(ans.CredentialEnforcement.Continue)
		var0.LogSeverity = types.StringValue(ans.CredentialEnforcement.LogSeverity)
		var0.Mode = var1
	}
	var var2 []urlAccessProfilesDsModelMlavEngineUrlbasedEnabledObject
	if len(ans.MlavEngineUrlbasedEnabled) != 0 {
		var2 = make([]urlAccessProfilesDsModelMlavEngineUrlbasedEnabledObject, 0, len(ans.MlavEngineUrlbasedEnabled))
		for var3Index := range ans.MlavEngineUrlbasedEnabled {
			var3 := ans.MlavEngineUrlbasedEnabled[var3Index]
			var var4 urlAccessProfilesDsModelMlavEngineUrlbasedEnabledObject
			var4.MlavPolicyAction = types.StringValue(var3.MlavPolicyAction)
			var4.Name = types.StringValue(var3.Name)
			var2 = append(var2, var4)
		}
	}
	state.Alert = EncodeStringSlice(ans.Alert)
	state.Allow = EncodeStringSlice(ans.Allow)
	state.Block = EncodeStringSlice(ans.Block)
	state.Continue = EncodeStringSlice(ans.Continue)
	state.CredentialEnforcement = var0
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.LogContainerPageOnly = types.BoolValue(ans.LogContainerPageOnly)
	state.LogHttpHdrReferer = types.BoolValue(ans.LogHttpHdrReferer)
	state.LogHttpHdrUserAgent = types.BoolValue(ans.LogHttpHdrUserAgent)
	state.LogHttpHdrXff = types.BoolValue(ans.LogHttpHdrXff)
	state.MlavCategoryException = EncodeStringSlice(ans.MlavCategoryException)
	state.MlavEngineUrlbasedEnabled = var2
	state.Name = types.StringValue(ans.Name)
	state.SafeSearchEnforcement = types.BoolValue(ans.SafeSearchEnforcement)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &urlAccessProfilesResource{}
	_ resource.ResourceWithConfigure   = &urlAccessProfilesResource{}
	_ resource.ResourceWithImportState = &urlAccessProfilesResource{}
)

func NewUrlAccessProfilesResource() resource.Resource {
	return &urlAccessProfilesResource{}
}

type urlAccessProfilesResource struct {
	client *sase.Client
}

type urlAccessProfilesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/url-access-profiles
	Alert                     []types.String                                            `tfsdk:"alert"`
	Allow                     []types.String                                            `tfsdk:"allow"`
	Block                     []types.String                                            `tfsdk:"block"`
	Continue                  []types.String                                            `tfsdk:"continue"`
	CredentialEnforcement     *urlAccessProfilesRsModelCredentialEnforcementObject      `tfsdk:"credential_enforcement"`
	Description               types.String                                              `tfsdk:"description"`
	ObjectId                  types.String                                              `tfsdk:"object_id"`
	LogContainerPageOnly      types.Bool                                                `tfsdk:"log_container_page_only"`
	LogHttpHdrReferer         types.Bool                                                `tfsdk:"log_http_hdr_referer"`
	LogHttpHdrUserAgent       types.Bool                                                `tfsdk:"log_http_hdr_user_agent"`
	LogHttpHdrXff             types.Bool                                                `tfsdk:"log_http_hdr_xff"`
	MlavCategoryException     []types.String                                            `tfsdk:"mlav_category_exception"`
	MlavEngineUrlbasedEnabled []urlAccessProfilesRsModelMlavEngineUrlbasedEnabledObject `tfsdk:"mlav_engine_urlbased_enabled"`
	Name                      types.String                                              `tfsdk:"name"`
	SafeSearchEnforcement     types.Bool                                                `tfsdk:"safe_search_enforcement"`
}

type urlAccessProfilesRsModelCredentialEnforcementObject struct {
	Alert       []types.String                      `tfsdk:"alert"`
	Allow       []types.String                      `tfsdk:"allow"`
	Block       []types.String                      `tfsdk:"block"`
	Continue    []types.String                      `tfsdk:"continue"`
	LogSeverity types.String                        `tfsdk:"log_severity"`
	Mode        *urlAccessProfilesRsModelModeObject `tfsdk:"mode"`
}

type urlAccessProfilesRsModelModeObject struct {
	Disabled          types.Bool   `tfsdk:"disabled"`
	DomainCredentials types.Bool   `tfsdk:"domain_credentials"`
	GroupMapping      types.String `tfsdk:"group_mapping"`
	IpUser            types.Bool   `tfsdk:"ip_user"`
}

type urlAccessProfilesRsModelMlavEngineUrlbasedEnabledObject struct {
	MlavPolicyAction types.String `tfsdk:"mlav_policy_action"`
	Name             types.String `tfsdk:"name"`
}

// Metadata returns the data source type name.
func (r *urlAccessProfilesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_url_access_profiles"
}

// Schema defines the schema for this listing data source.
func (r *urlAccessProfilesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

			"alert": rsschema.ListAttribute{
				Description:         "The `alert` parameter.",
				MarkdownDescription: "The `alert` parameter.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"allow": rsschema.ListAttribute{
				Description:         "The `allow` parameter.",
				MarkdownDescription: "The `allow` parameter.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"block": rsschema.ListAttribute{
				Description:         "The `block` parameter.",
				MarkdownDescription: "The `block` parameter.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"continue": rsschema.ListAttribute{
				Description:         "The `continue` parameter.",
				MarkdownDescription: "The `continue` parameter.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"credential_enforcement": rsschema.SingleNestedAttribute{
				Description:         "The `credential_enforcement` parameter.",
				MarkdownDescription: "The `credential_enforcement` parameter.",
				Optional:            true,
				Attributes: map[string]rsschema.Attribute{
					"alert": rsschema.ListAttribute{
						Description:         "The `alert` parameter.",
						MarkdownDescription: "The `alert` parameter.",
						Optional:            true,
						ElementType:         types.StringType,
					},
					"allow": rsschema.ListAttribute{
						Description:         "The `allow` parameter.",
						MarkdownDescription: "The `allow` parameter.",
						Optional:            true,
						ElementType:         types.StringType,
					},
					"block": rsschema.ListAttribute{
						Description:         "The `block` parameter.",
						MarkdownDescription: "The `block` parameter.",
						Optional:            true,
						ElementType:         types.StringType,
					},
					"continue": rsschema.ListAttribute{
						Description:         "The `continue` parameter.",
						MarkdownDescription: "The `continue` parameter.",
						Optional:            true,
						ElementType:         types.StringType,
					},
					"log_severity": rsschema.StringAttribute{
						Description:         "The `log_severity` parameter. Default: `%!q(*string=0xc000f53090)`.",
						MarkdownDescription: "The `log_severity` parameter. Default: `%!q(*string=0xc000f53090)`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							DefaultString("medium"),
						},
					},
					"mode": rsschema.SingleNestedAttribute{
						Description:         "The `mode` parameter.",
						MarkdownDescription: "The `mode` parameter.",
						Optional:            true,
						Attributes: map[string]rsschema.Attribute{
							"disabled": rsschema.BoolAttribute{
								Description:         "The `disabled` parameter.",
								MarkdownDescription: "The `disabled` parameter.",
								Optional:            true,
							},
							"domain_credentials": rsschema.BoolAttribute{
								Description:         "The `domain_credentials` parameter.",
								MarkdownDescription: "The `domain_credentials` parameter.",
								Optional:            true,
							},
							"group_mapping": rsschema.StringAttribute{
								Description:         "The `group_mapping` parameter.",
								MarkdownDescription: "The `group_mapping` parameter.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"ip_user": rsschema.BoolAttribute{
								Description:         "The `ip_user` parameter.",
								MarkdownDescription: "The `ip_user` parameter.",
								Optional:            true,
							},
						},
					},
				},
			},
			"description": rsschema.StringAttribute{
				Description:         "The `description` parameter. String length must be between 0 and 255.",
				MarkdownDescription: "The `description` parameter. String length must be between 0 and 255.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 255),
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
			"log_container_page_only": rsschema.BoolAttribute{
				Description:         "The `log_container_page_only` parameter. Default: `true`.",
				MarkdownDescription: "The `log_container_page_only` parameter. Default: `true`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(true),
				},
			},
			"log_http_hdr_referer": rsschema.BoolAttribute{
				Description:         "The `log_http_hdr_referer` parameter. Default: `false`.",
				MarkdownDescription: "The `log_http_hdr_referer` parameter. Default: `false`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"log_http_hdr_user_agent": rsschema.BoolAttribute{
				Description:         "The `log_http_hdr_user_agent` parameter. Default: `false`.",
				MarkdownDescription: "The `log_http_hdr_user_agent` parameter. Default: `false`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"log_http_hdr_xff": rsschema.BoolAttribute{
				Description:         "The `log_http_hdr_xff` parameter. Default: `false`.",
				MarkdownDescription: "The `log_http_hdr_xff` parameter. Default: `false`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"mlav_category_exception": rsschema.ListAttribute{
				Description:         "The `mlav_category_exception` parameter.",
				MarkdownDescription: "The `mlav_category_exception` parameter.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"mlav_engine_urlbased_enabled": rsschema.ListNestedAttribute{
				Description:         "The `mlav_engine_urlbased_enabled` parameter.",
				MarkdownDescription: "The `mlav_engine_urlbased_enabled` parameter.",
				Optional:            true,
				NestedObject: rsschema.NestedAttributeObject{
					Attributes: map[string]rsschema.Attribute{
						"mlav_policy_action": rsschema.StringAttribute{
							Description:         "The `mlav_policy_action` parameter. Value must be one of: `\"allow\"`, `\"alert\"`, `\"block\"`.",
							MarkdownDescription: "The `mlav_policy_action` parameter. Value must be one of: `\"allow\"`, `\"alert\"`, `\"block\"`.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								DefaultString(""),
							},
							Validators: []validator.String{
								stringvalidator.OneOf("allow", "alert", "block"),
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
					},
				},
			},
			"name": rsschema.StringAttribute{
				Description:         "The `name` parameter.",
				MarkdownDescription: "The `name` parameter.",
				Required:            true,
			},
			"safe_search_enforcement": rsschema.BoolAttribute{
				Description:         "The `safe_search_enforcement` parameter. Default: `false`.",
				MarkdownDescription: "The `safe_search_enforcement` parameter. Default: `false`.",
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
func (r *urlAccessProfilesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *urlAccessProfilesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state urlAccessProfilesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_url_access_profiles",
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := uyrOkzA.NewClient(r.client)
	input := uyrOkzA.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 wleovFf.Config
	var0.Alert = DecodeStringSlice(state.Alert)
	var0.Allow = DecodeStringSlice(state.Allow)
	var0.Block = DecodeStringSlice(state.Block)
	var0.Continue = DecodeStringSlice(state.Continue)
	var var1 *wleovFf.CredentialEnforcementObject
	if state.CredentialEnforcement != nil {
		var1 = &wleovFf.CredentialEnforcementObject{}
		var1.Alert = DecodeStringSlice(state.CredentialEnforcement.Alert)
		var1.Allow = DecodeStringSlice(state.CredentialEnforcement.Allow)
		var1.Block = DecodeStringSlice(state.CredentialEnforcement.Block)
		var1.Continue = DecodeStringSlice(state.CredentialEnforcement.Continue)
		var1.LogSeverity = state.CredentialEnforcement.LogSeverity.ValueString()
		var var2 *wleovFf.ModeObject
		if state.CredentialEnforcement.Mode != nil {
			var2 = &wleovFf.ModeObject{}
			if state.CredentialEnforcement.Mode.Disabled.ValueBool() {
				var2.Disabled = struct{}{}
			}
			if state.CredentialEnforcement.Mode.DomainCredentials.ValueBool() {
				var2.DomainCredentials = struct{}{}
			}
			var2.GroupMapping = state.CredentialEnforcement.Mode.GroupMapping.ValueString()
			if state.CredentialEnforcement.Mode.IpUser.ValueBool() {
				var2.IpUser = struct{}{}
			}
		}
		var1.Mode = var2
	}
	var0.CredentialEnforcement = var1
	var0.Description = state.Description.ValueString()
	var0.LogContainerPageOnly = state.LogContainerPageOnly.ValueBool()
	var0.LogHttpHdrReferer = state.LogHttpHdrReferer.ValueBool()
	var0.LogHttpHdrUserAgent = state.LogHttpHdrUserAgent.ValueBool()
	var0.LogHttpHdrXff = state.LogHttpHdrXff.ValueBool()
	var0.MlavCategoryException = DecodeStringSlice(state.MlavCategoryException)
	var var3 []wleovFf.MlavEngineUrlbasedEnabledObject
	if len(state.MlavEngineUrlbasedEnabled) != 0 {
		var3 = make([]wleovFf.MlavEngineUrlbasedEnabledObject, 0, len(state.MlavEngineUrlbasedEnabled))
		for var4Index := range state.MlavEngineUrlbasedEnabled {
			var4 := state.MlavEngineUrlbasedEnabled[var4Index]
			var var5 wleovFf.MlavEngineUrlbasedEnabledObject
			var5.MlavPolicyAction = var4.MlavPolicyAction.ValueString()
			var5.Name = var4.Name.ValueString()
			var3 = append(var3, var5)
		}
	}
	var0.MlavEngineUrlbasedEnabled = var3
	var0.Name = state.Name.ValueString()
	var0.SafeSearchEnforcement = state.SafeSearchEnforcement.ValueBool()
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
	var var6 *urlAccessProfilesRsModelCredentialEnforcementObject
	if ans.CredentialEnforcement != nil {
		var6 = &urlAccessProfilesRsModelCredentialEnforcementObject{}
		var var7 *urlAccessProfilesRsModelModeObject
		if ans.CredentialEnforcement.Mode != nil {
			var7 = &urlAccessProfilesRsModelModeObject{}
			if ans.CredentialEnforcement.Mode.Disabled != nil {
				var7.Disabled = types.BoolValue(true)
			}
			if ans.CredentialEnforcement.Mode.DomainCredentials != nil {
				var7.DomainCredentials = types.BoolValue(true)
			}
			var7.GroupMapping = types.StringValue(ans.CredentialEnforcement.Mode.GroupMapping)
			if ans.CredentialEnforcement.Mode.IpUser != nil {
				var7.IpUser = types.BoolValue(true)
			}
		}
		var6.Alert = EncodeStringSlice(ans.CredentialEnforcement.Alert)
		var6.Allow = EncodeStringSlice(ans.CredentialEnforcement.Allow)
		var6.Block = EncodeStringSlice(ans.CredentialEnforcement.Block)
		var6.Continue = EncodeStringSlice(ans.CredentialEnforcement.Continue)
		var6.LogSeverity = types.StringValue(ans.CredentialEnforcement.LogSeverity)
		var6.Mode = var7
	}
	var var8 []urlAccessProfilesRsModelMlavEngineUrlbasedEnabledObject
	if len(ans.MlavEngineUrlbasedEnabled) != 0 {
		var8 = make([]urlAccessProfilesRsModelMlavEngineUrlbasedEnabledObject, 0, len(ans.MlavEngineUrlbasedEnabled))
		for var9Index := range ans.MlavEngineUrlbasedEnabled {
			var9 := ans.MlavEngineUrlbasedEnabled[var9Index]
			var var10 urlAccessProfilesRsModelMlavEngineUrlbasedEnabledObject
			var10.MlavPolicyAction = types.StringValue(var9.MlavPolicyAction)
			var10.Name = types.StringValue(var9.Name)
			var8 = append(var8, var10)
		}
	}
	state.Alert = EncodeStringSlice(ans.Alert)
	state.Allow = EncodeStringSlice(ans.Allow)
	state.Block = EncodeStringSlice(ans.Block)
	state.Continue = EncodeStringSlice(ans.Continue)
	state.CredentialEnforcement = var6
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.LogContainerPageOnly = types.BoolValue(ans.LogContainerPageOnly)
	state.LogHttpHdrReferer = types.BoolValue(ans.LogHttpHdrReferer)
	state.LogHttpHdrUserAgent = types.BoolValue(ans.LogHttpHdrUserAgent)
	state.LogHttpHdrXff = types.BoolValue(ans.LogHttpHdrXff)
	state.MlavCategoryException = EncodeStringSlice(ans.MlavCategoryException)
	state.MlavEngineUrlbasedEnabled = var8
	state.Name = types.StringValue(ans.Name)
	state.SafeSearchEnforcement = types.BoolValue(ans.SafeSearchEnforcement)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *urlAccessProfilesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state urlAccessProfilesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_url_access_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := uyrOkzA.NewClient(r.client)
	input := uyrOkzA.ReadInput{
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
	var var0 *urlAccessProfilesRsModelCredentialEnforcementObject
	if ans.CredentialEnforcement != nil {
		var0 = &urlAccessProfilesRsModelCredentialEnforcementObject{}
		var var1 *urlAccessProfilesRsModelModeObject
		if ans.CredentialEnforcement.Mode != nil {
			var1 = &urlAccessProfilesRsModelModeObject{}
			if ans.CredentialEnforcement.Mode.Disabled != nil {
				var1.Disabled = types.BoolValue(true)
			}
			if ans.CredentialEnforcement.Mode.DomainCredentials != nil {
				var1.DomainCredentials = types.BoolValue(true)
			}
			var1.GroupMapping = types.StringValue(ans.CredentialEnforcement.Mode.GroupMapping)
			if ans.CredentialEnforcement.Mode.IpUser != nil {
				var1.IpUser = types.BoolValue(true)
			}
		}
		var0.Alert = EncodeStringSlice(ans.CredentialEnforcement.Alert)
		var0.Allow = EncodeStringSlice(ans.CredentialEnforcement.Allow)
		var0.Block = EncodeStringSlice(ans.CredentialEnforcement.Block)
		var0.Continue = EncodeStringSlice(ans.CredentialEnforcement.Continue)
		var0.LogSeverity = types.StringValue(ans.CredentialEnforcement.LogSeverity)
		var0.Mode = var1
	}
	var var2 []urlAccessProfilesRsModelMlavEngineUrlbasedEnabledObject
	if len(ans.MlavEngineUrlbasedEnabled) != 0 {
		var2 = make([]urlAccessProfilesRsModelMlavEngineUrlbasedEnabledObject, 0, len(ans.MlavEngineUrlbasedEnabled))
		for var3Index := range ans.MlavEngineUrlbasedEnabled {
			var3 := ans.MlavEngineUrlbasedEnabled[var3Index]
			var var4 urlAccessProfilesRsModelMlavEngineUrlbasedEnabledObject
			var4.MlavPolicyAction = types.StringValue(var3.MlavPolicyAction)
			var4.Name = types.StringValue(var3.Name)
			var2 = append(var2, var4)
		}
	}
	state.Alert = EncodeStringSlice(ans.Alert)
	state.Allow = EncodeStringSlice(ans.Allow)
	state.Block = EncodeStringSlice(ans.Block)
	state.Continue = EncodeStringSlice(ans.Continue)
	state.CredentialEnforcement = var0
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.LogContainerPageOnly = types.BoolValue(ans.LogContainerPageOnly)
	state.LogHttpHdrReferer = types.BoolValue(ans.LogHttpHdrReferer)
	state.LogHttpHdrUserAgent = types.BoolValue(ans.LogHttpHdrUserAgent)
	state.LogHttpHdrXff = types.BoolValue(ans.LogHttpHdrXff)
	state.MlavCategoryException = EncodeStringSlice(ans.MlavCategoryException)
	state.MlavEngineUrlbasedEnabled = var2
	state.Name = types.StringValue(ans.Name)
	state.SafeSearchEnforcement = types.BoolValue(ans.SafeSearchEnforcement)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *urlAccessProfilesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state urlAccessProfilesRsModel
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
		"resource_name":               "sase_url_access_profiles",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := uyrOkzA.NewClient(r.client)
	input := uyrOkzA.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 wleovFf.Config
	var0.Alert = DecodeStringSlice(plan.Alert)
	var0.Allow = DecodeStringSlice(plan.Allow)
	var0.Block = DecodeStringSlice(plan.Block)
	var0.Continue = DecodeStringSlice(plan.Continue)
	var var1 *wleovFf.CredentialEnforcementObject
	if plan.CredentialEnforcement != nil {
		var1 = &wleovFf.CredentialEnforcementObject{}
		var1.Alert = DecodeStringSlice(plan.CredentialEnforcement.Alert)
		var1.Allow = DecodeStringSlice(plan.CredentialEnforcement.Allow)
		var1.Block = DecodeStringSlice(plan.CredentialEnforcement.Block)
		var1.Continue = DecodeStringSlice(plan.CredentialEnforcement.Continue)
		var1.LogSeverity = plan.CredentialEnforcement.LogSeverity.ValueString()
		var var2 *wleovFf.ModeObject
		if plan.CredentialEnforcement.Mode != nil {
			var2 = &wleovFf.ModeObject{}
			if plan.CredentialEnforcement.Mode.Disabled.ValueBool() {
				var2.Disabled = struct{}{}
			}
			if plan.CredentialEnforcement.Mode.DomainCredentials.ValueBool() {
				var2.DomainCredentials = struct{}{}
			}
			var2.GroupMapping = plan.CredentialEnforcement.Mode.GroupMapping.ValueString()
			if plan.CredentialEnforcement.Mode.IpUser.ValueBool() {
				var2.IpUser = struct{}{}
			}
		}
		var1.Mode = var2
	}
	var0.CredentialEnforcement = var1
	var0.Description = plan.Description.ValueString()
	var0.LogContainerPageOnly = plan.LogContainerPageOnly.ValueBool()
	var0.LogHttpHdrReferer = plan.LogHttpHdrReferer.ValueBool()
	var0.LogHttpHdrUserAgent = plan.LogHttpHdrUserAgent.ValueBool()
	var0.LogHttpHdrXff = plan.LogHttpHdrXff.ValueBool()
	var0.MlavCategoryException = DecodeStringSlice(plan.MlavCategoryException)
	var var3 []wleovFf.MlavEngineUrlbasedEnabledObject
	if len(plan.MlavEngineUrlbasedEnabled) != 0 {
		var3 = make([]wleovFf.MlavEngineUrlbasedEnabledObject, 0, len(plan.MlavEngineUrlbasedEnabled))
		for var4Index := range plan.MlavEngineUrlbasedEnabled {
			var4 := plan.MlavEngineUrlbasedEnabled[var4Index]
			var var5 wleovFf.MlavEngineUrlbasedEnabledObject
			var5.MlavPolicyAction = var4.MlavPolicyAction.ValueString()
			var5.Name = var4.Name.ValueString()
			var3 = append(var3, var5)
		}
	}
	var0.MlavEngineUrlbasedEnabled = var3
	var0.Name = plan.Name.ValueString()
	var0.SafeSearchEnforcement = plan.SafeSearchEnforcement.ValueBool()
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var6 *urlAccessProfilesRsModelCredentialEnforcementObject
	if ans.CredentialEnforcement != nil {
		var6 = &urlAccessProfilesRsModelCredentialEnforcementObject{}
		var var7 *urlAccessProfilesRsModelModeObject
		if ans.CredentialEnforcement.Mode != nil {
			var7 = &urlAccessProfilesRsModelModeObject{}
			if ans.CredentialEnforcement.Mode.Disabled != nil {
				var7.Disabled = types.BoolValue(true)
			}
			if ans.CredentialEnforcement.Mode.DomainCredentials != nil {
				var7.DomainCredentials = types.BoolValue(true)
			}
			var7.GroupMapping = types.StringValue(ans.CredentialEnforcement.Mode.GroupMapping)
			if ans.CredentialEnforcement.Mode.IpUser != nil {
				var7.IpUser = types.BoolValue(true)
			}
		}
		var6.Alert = EncodeStringSlice(ans.CredentialEnforcement.Alert)
		var6.Allow = EncodeStringSlice(ans.CredentialEnforcement.Allow)
		var6.Block = EncodeStringSlice(ans.CredentialEnforcement.Block)
		var6.Continue = EncodeStringSlice(ans.CredentialEnforcement.Continue)
		var6.LogSeverity = types.StringValue(ans.CredentialEnforcement.LogSeverity)
		var6.Mode = var7
	}
	var var8 []urlAccessProfilesRsModelMlavEngineUrlbasedEnabledObject
	if len(ans.MlavEngineUrlbasedEnabled) != 0 {
		var8 = make([]urlAccessProfilesRsModelMlavEngineUrlbasedEnabledObject, 0, len(ans.MlavEngineUrlbasedEnabled))
		for var9Index := range ans.MlavEngineUrlbasedEnabled {
			var9 := ans.MlavEngineUrlbasedEnabled[var9Index]
			var var10 urlAccessProfilesRsModelMlavEngineUrlbasedEnabledObject
			var10.MlavPolicyAction = types.StringValue(var9.MlavPolicyAction)
			var10.Name = types.StringValue(var9.Name)
			var8 = append(var8, var10)
		}
	}
	state.Alert = EncodeStringSlice(ans.Alert)
	state.Allow = EncodeStringSlice(ans.Allow)
	state.Block = EncodeStringSlice(ans.Block)
	state.Continue = EncodeStringSlice(ans.Continue)
	state.CredentialEnforcement = var6
	state.Description = types.StringValue(ans.Description)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.LogContainerPageOnly = types.BoolValue(ans.LogContainerPageOnly)
	state.LogHttpHdrReferer = types.BoolValue(ans.LogHttpHdrReferer)
	state.LogHttpHdrUserAgent = types.BoolValue(ans.LogHttpHdrUserAgent)
	state.LogHttpHdrXff = types.BoolValue(ans.LogHttpHdrXff)
	state.MlavCategoryException = EncodeStringSlice(ans.MlavCategoryException)
	state.MlavEngineUrlbasedEnabled = var8
	state.Name = types.StringValue(ans.Name)
	state.SafeSearchEnforcement = types.BoolValue(ans.SafeSearchEnforcement)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *urlAccessProfilesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name":               "sase_url_access_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	svc := uyrOkzA.NewClient(r.client)
	input := uyrOkzA.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *urlAccessProfilesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
