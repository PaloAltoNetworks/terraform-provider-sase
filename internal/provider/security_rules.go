package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	ffcMtmY "github.com/paloaltonetworks/sase-go/netsec/schema/security/rules"
	mPRFtcU "github.com/paloaltonetworks/sase-go/netsec/service/v1/securityrules"

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
	_ datasource.DataSource              = &securityRulesListDataSource{}
	_ datasource.DataSourceWithConfigure = &securityRulesListDataSource{}
)

func NewSecurityRulesListDataSource() datasource.DataSource {
	return &securityRulesListDataSource{}
}

type securityRulesListDataSource struct {
	client *sase.Client
}

type securityRulesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit    types.Int64  `tfsdk:"limit"`
	Offset   types.Int64  `tfsdk:"offset"`
	Position types.String `tfsdk:"position"`
	Folder   types.String `tfsdk:"folder"`
	Name     types.String `tfsdk:"name"`

	// Output.
	Data []securityRulesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type securityRulesListDsModelConfig struct {
	Action            types.String                                  `tfsdk:"action"`
	Application       []types.String                                `tfsdk:"application"`
	Category          []types.String                                `tfsdk:"category"`
	Description       types.String                                  `tfsdk:"description"`
	Destination       []types.String                                `tfsdk:"destination"`
	DestinationHip    []types.String                                `tfsdk:"destination_hip"`
	Disabled          types.Bool                                    `tfsdk:"disabled"`
	From              []types.String                                `tfsdk:"from"`
	ObjectId          types.String                                  `tfsdk:"object_id"`
	LogSetting        types.String                                  `tfsdk:"log_setting"`
	Name              types.String                                  `tfsdk:"name"`
	NegateDestination types.Bool                                    `tfsdk:"negate_destination"`
	NegateSource      types.Bool                                    `tfsdk:"negate_source"`
	ProfileSetting    *securityRulesListDsModelProfileSettingObject `tfsdk:"profile_setting"`
	Service           []types.String                                `tfsdk:"service"`
	Source            []types.String                                `tfsdk:"source"`
	SourceHip         []types.String                                `tfsdk:"source_hip"`
	SourceUser        []types.String                                `tfsdk:"source_user"`
	Tag               []types.String                                `tfsdk:"tag"`
	To                []types.String                                `tfsdk:"to"`
}

type securityRulesListDsModelProfileSettingObject struct {
	Group []types.String `tfsdk:"group"`
}

// Metadata returns the data source type name.
func (d *securityRulesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_rules_list"
}

// Schema defines the schema for this listing data source.
func (d *securityRulesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"position": dsschema.StringAttribute{
				Description: "The position of a security rule",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("pre", "post"),
				},
			},
			"folder": dsschema.StringAttribute{
				Description: "The folder of the entry",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Service Connections", "Mobile Users Container", "Mobile Users Explicit Proxy"),
				},
			},
			"name": dsschema.StringAttribute{
				Description: "The name of the entry",
				Optional:    true,
			},

			// Output.
			"data": dsschema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"action": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"application": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"category": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"description": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"destination": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"destination_hip": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"disabled": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"from": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"object_id": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"log_setting": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"negate_destination": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"negate_source": dsschema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"profile_setting": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"group": dsschema.ListAttribute{
									Description: "",
									Computed:    true,
									ElementType: types.StringType,
								},
							},
						},
						"service": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"source": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"source_hip": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"source_user": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"tag": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"to": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
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
func (d *securityRulesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *securityRulesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state securityRulesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_security_rules_list",
		"terraform_provider_function": "Read",
		"limit":                       state.Limit.ValueInt64(),
		"has_limit":                   !state.Limit.IsNull(),
		"offset":                      state.Offset.ValueInt64(),
		"has_offset":                  !state.Offset.IsNull(),
		"position":                    state.Position.ValueString(),
		"folder":                      state.Folder.ValueString(),
		"name":                        state.Name.ValueString(),
		"has_name":                    !state.Name.IsNull(),
	})

	// Prepare to run the command.
	svc := mPRFtcU.NewClient(d.client)
	input := mPRFtcU.ListInput{
		Position: state.Position.ValueString(),
		Folder:   state.Folder.ValueString(),
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
	idBuilder.WriteString(input.Position)
	idBuilder.WriteString(IdSeparator)
	idBuilder.WriteString(input.Folder)
	idBuilder.WriteString(IdSeparator)
	if input.Name != nil {
		idBuilder.WriteString(*input.Name)
	}
	state.Id = types.StringValue(idBuilder.String())
	var var0 []securityRulesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]securityRulesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 securityRulesListDsModelConfig
			var var3 *securityRulesListDsModelProfileSettingObject
			if var1.ProfileSetting != nil {
				var3 = &securityRulesListDsModelProfileSettingObject{}
				var3.Group = EncodeStringSlice(var1.ProfileSetting.Group)
			}
			var2.Action = types.StringValue(var1.Action)
			var2.Application = EncodeStringSlice(var1.Application)
			var2.Category = EncodeStringSlice(var1.Category)
			var2.Description = types.StringValue(var1.Description)
			var2.Destination = EncodeStringSlice(var1.Destination)
			var2.DestinationHip = EncodeStringSlice(var1.DestinationHip)
			var2.Disabled = types.BoolValue(var1.Disabled)
			var2.From = EncodeStringSlice(var1.From)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.LogSetting = types.StringValue(var1.LogSetting)
			var2.Name = types.StringValue(var1.Name)
			var2.NegateDestination = types.BoolValue(var1.NegateDestination)
			var2.NegateSource = types.BoolValue(var1.NegateSource)
			var2.ProfileSetting = var3
			var2.Service = EncodeStringSlice(var1.Service)
			var2.Source = EncodeStringSlice(var1.Source)
			var2.SourceHip = EncodeStringSlice(var1.SourceHip)
			var2.SourceUser = EncodeStringSlice(var1.SourceUser)
			var2.Tag = EncodeStringSlice(var1.Tag)
			var2.To = EncodeStringSlice(var1.To)
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
	_ datasource.DataSource              = &securityRulesDataSource{}
	_ datasource.DataSourceWithConfigure = &securityRulesDataSource{}
)

func NewSecurityRulesDataSource() datasource.DataSource {
	return &securityRulesDataSource{}
}

type securityRulesDataSource struct {
	client *sase.Client
}

type securityRulesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`

	// Output.
	// Ref: #/components/schemas/security-rules
	Action         types.String   `tfsdk:"action"`
	Application    []types.String `tfsdk:"application"`
	Category       []types.String `tfsdk:"category"`
	Description    types.String   `tfsdk:"description"`
	Destination    []types.String `tfsdk:"destination"`
	DestinationHip []types.String `tfsdk:"destination_hip"`
	Disabled       types.Bool     `tfsdk:"disabled"`
	From           []types.String `tfsdk:"from"`
	// input omit: ObjectId
	LogSetting        types.String                              `tfsdk:"log_setting"`
	Name              types.String                              `tfsdk:"name"`
	NegateDestination types.Bool                                `tfsdk:"negate_destination"`
	NegateSource      types.Bool                                `tfsdk:"negate_source"`
	ProfileSetting    *securityRulesDsModelProfileSettingObject `tfsdk:"profile_setting"`
	Service           []types.String                            `tfsdk:"service"`
	Source            []types.String                            `tfsdk:"source"`
	SourceHip         []types.String                            `tfsdk:"source_hip"`
	SourceUser        []types.String                            `tfsdk:"source_user"`
	Tag               []types.String                            `tfsdk:"tag"`
	To                []types.String                            `tfsdk:"to"`
}

type securityRulesDsModelProfileSettingObject struct {
	Group []types.String `tfsdk:"group"`
}

// Metadata returns the data source type name.
func (d *securityRulesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_rules"
}

// Schema defines the schema for this listing data source.
func (d *securityRulesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"action": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"application": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"category": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"destination": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"destination_hip": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"disabled": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"from": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"log_setting": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"negate_destination": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"negate_source": dsschema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"profile_setting": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"group": dsschema.ListAttribute{
						Description: "",
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
			"service": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"source": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"source_hip": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"source_user": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"tag": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"to": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure prepares the struct.
func (d *securityRulesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *securityRulesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state securityRulesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_security_rules",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to run the command.
	svc := mPRFtcU.NewClient(d.client)
	input := mPRFtcU.ReadInput{
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
	var var0 *securityRulesDsModelProfileSettingObject
	if ans.ProfileSetting != nil {
		var0 = &securityRulesDsModelProfileSettingObject{}
		var0.Group = EncodeStringSlice(ans.ProfileSetting.Group)
	}
	state.Action = types.StringValue(ans.Action)
	state.Application = EncodeStringSlice(ans.Application)
	state.Category = EncodeStringSlice(ans.Category)
	state.Description = types.StringValue(ans.Description)
	state.Destination = EncodeStringSlice(ans.Destination)
	state.DestinationHip = EncodeStringSlice(ans.DestinationHip)
	state.Disabled = types.BoolValue(ans.Disabled)
	state.From = EncodeStringSlice(ans.From)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.LogSetting = types.StringValue(ans.LogSetting)
	state.Name = types.StringValue(ans.Name)
	state.NegateDestination = types.BoolValue(ans.NegateDestination)
	state.NegateSource = types.BoolValue(ans.NegateSource)
	state.ProfileSetting = var0
	state.Service = EncodeStringSlice(ans.Service)
	state.Source = EncodeStringSlice(ans.Source)
	state.SourceHip = EncodeStringSlice(ans.SourceHip)
	state.SourceUser = EncodeStringSlice(ans.SourceUser)
	state.Tag = EncodeStringSlice(ans.Tag)
	state.To = EncodeStringSlice(ans.To)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &securityRulesResource{}
	_ resource.ResourceWithConfigure   = &securityRulesResource{}
	_ resource.ResourceWithImportState = &securityRulesResource{}
)

func NewSecurityRulesResource() resource.Resource {
	return &securityRulesResource{}
}

type securityRulesResource struct {
	client *sase.Client
}

type securityRulesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Position types.String `tfsdk:"position"`
	Folder   types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/security-rules
	Action            types.String                              `tfsdk:"action"`
	Application       []types.String                            `tfsdk:"application"`
	Category          []types.String                            `tfsdk:"category"`
	Description       types.String                              `tfsdk:"description"`
	Destination       []types.String                            `tfsdk:"destination"`
	DestinationHip    []types.String                            `tfsdk:"destination_hip"`
	Disabled          types.Bool                                `tfsdk:"disabled"`
	From              []types.String                            `tfsdk:"from"`
	ObjectId          types.String                              `tfsdk:"object_id"`
	LogSetting        types.String                              `tfsdk:"log_setting"`
	Name              types.String                              `tfsdk:"name"`
	NegateDestination types.Bool                                `tfsdk:"negate_destination"`
	NegateSource      types.Bool                                `tfsdk:"negate_source"`
	ProfileSetting    *securityRulesRsModelProfileSettingObject `tfsdk:"profile_setting"`
	Service           []types.String                            `tfsdk:"service"`
	Source            []types.String                            `tfsdk:"source"`
	SourceHip         []types.String                            `tfsdk:"source_hip"`
	SourceUser        []types.String                            `tfsdk:"source_user"`
	Tag               []types.String                            `tfsdk:"tag"`
	To                []types.String                            `tfsdk:"to"`
}

type securityRulesRsModelProfileSettingObject struct {
	Group []types.String `tfsdk:"group"`
}

// Metadata returns the data source type name.
func (r *securityRulesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_rules"
}

// Schema defines the schema for this listing data source.
func (r *securityRulesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"position": rsschema.StringAttribute{
				Description: "The position of a security rule",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("pre", "post"),
				},
			},
			"folder": rsschema.StringAttribute{
				Description: "The folder of the entry",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Service Connections", "Mobile Users Container", "Mobile Users Explicit Proxy"),
				},
			},

			"action": rsschema.StringAttribute{
				Description: "",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
			},
			"application": rsschema.ListAttribute{
				Description: "",
				Required:    true,
				ElementType: types.StringType,
			},
			"category": rsschema.ListAttribute{
				Description: "",
				Required:    true,
				ElementType: types.StringType,
			},
			"description": rsschema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
			},
			"destination": rsschema.ListAttribute{
				Description: "",
				Required:    true,
				ElementType: types.StringType,
			},
			"destination_hip": rsschema.ListAttribute{
				Description: "",
				Optional:    true,
				ElementType: types.StringType,
			},
			"disabled": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"from": rsschema.ListAttribute{
				Description: "",
				Required:    true,
				ElementType: types.StringType,
			},
			"object_id": rsschema.StringAttribute{
				Description: "",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"log_setting": rsschema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
			},
			"name": rsschema.StringAttribute{
				Description: "",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
			},
			"negate_destination": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"negate_source": rsschema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"profile_setting": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"group": rsschema.ListAttribute{
						Description: "",
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
			"service": rsschema.ListAttribute{
				Description: "",
				Required:    true,
				ElementType: types.StringType,
			},
			"source": rsschema.ListAttribute{
				Description: "",
				Required:    true,
				ElementType: types.StringType,
			},
			"source_hip": rsschema.ListAttribute{
				Description: "",
				Optional:    true,
				ElementType: types.StringType,
			},
			"source_user": rsschema.ListAttribute{
				Description: "",
				Required:    true,
				ElementType: types.StringType,
			},
			"tag": rsschema.ListAttribute{
				Description: "",
				Optional:    true,
				ElementType: types.StringType,
			},
			"to": rsschema.ListAttribute{
				Description: "",
				Required:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure prepares the struct.
func (r *securityRulesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *securityRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state securityRulesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_security_rules",
		"position":                    state.Position.ValueString(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := mPRFtcU.NewClient(r.client)
	input := mPRFtcU.CreateInput{
		Position: state.Position.ValueString(),
		Folder:   state.Folder.ValueString(),
	}
	var var0 ffcMtmY.Config
	var0.Action = state.Action.ValueString()
	var0.Application = DecodeStringSlice(state.Application)
	var0.Category = DecodeStringSlice(state.Category)
	var0.Description = state.Description.ValueString()
	var0.Destination = DecodeStringSlice(state.Destination)
	var0.DestinationHip = DecodeStringSlice(state.DestinationHip)
	var0.Disabled = state.Disabled.ValueBool()
	var0.From = DecodeStringSlice(state.From)
	var0.LogSetting = state.LogSetting.ValueString()
	var0.Name = state.Name.ValueString()
	var0.NegateDestination = state.NegateDestination.ValueBool()
	var0.NegateSource = state.NegateSource.ValueBool()
	var var1 *ffcMtmY.ProfileSettingObject
	if state.ProfileSetting != nil {
		var1 = &ffcMtmY.ProfileSettingObject{}
		var1.Group = DecodeStringSlice(state.ProfileSetting.Group)
	}
	var0.ProfileSetting = var1
	var0.Service = DecodeStringSlice(state.Service)
	var0.Source = DecodeStringSlice(state.Source)
	var0.SourceHip = DecodeStringSlice(state.SourceHip)
	var0.SourceUser = DecodeStringSlice(state.SourceUser)
	var0.Tag = DecodeStringSlice(state.Tag)
	var0.To = DecodeStringSlice(state.To)
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Create(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in create", err.Error())
		return
	}

	// Store the answer to state.
	var idBuilder strings.Builder
	idBuilder.WriteString(input.Position)
	idBuilder.WriteString(IdSeparator)
	idBuilder.WriteString(input.Folder)
	idBuilder.WriteString(IdSeparator)
	idBuilder.WriteString(ans.ObjectId)
	state.Id = types.StringValue(idBuilder.String())
	var var2 *securityRulesRsModelProfileSettingObject
	if ans.ProfileSetting != nil {
		var2 = &securityRulesRsModelProfileSettingObject{}
		var2.Group = EncodeStringSlice(ans.ProfileSetting.Group)
	}
	state.Action = types.StringValue(ans.Action)
	state.Application = EncodeStringSlice(ans.Application)
	state.Category = EncodeStringSlice(ans.Category)
	state.Description = types.StringValue(ans.Description)
	state.Destination = EncodeStringSlice(ans.Destination)
	state.DestinationHip = EncodeStringSlice(ans.DestinationHip)
	state.Disabled = types.BoolValue(ans.Disabled)
	state.From = EncodeStringSlice(ans.From)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.LogSetting = types.StringValue(ans.LogSetting)
	state.Name = types.StringValue(ans.Name)
	state.NegateDestination = types.BoolValue(ans.NegateDestination)
	state.NegateSource = types.BoolValue(ans.NegateSource)
	state.ProfileSetting = var2
	state.Service = EncodeStringSlice(ans.Service)
	state.Source = EncodeStringSlice(ans.Source)
	state.SourceHip = EncodeStringSlice(ans.SourceHip)
	state.SourceUser = EncodeStringSlice(ans.SourceUser)
	state.Tag = EncodeStringSlice(ans.Tag)
	state.To = EncodeStringSlice(ans.To)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *securityRulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var idType types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &idType)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := idType.ValueString()
	tokens := strings.Split(id, IdSeparator)
	if len(tokens) != 3 {
		resp.Diagnostics.AddError("Error in resource ID format", "Expected 3 tokens")
		return
	}

	var state securityRulesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_security_rules",
		"locMap":                      map[string]int{"Folder": 1, "ObjectId": 2, "Position": 0},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := mPRFtcU.NewClient(r.client)
	input := mPRFtcU.ReadInput{
		ObjectId: tokens[2],
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
	state.Position = types.StringValue(tokens[0])
	state.Folder = types.StringValue(tokens[1])
	state.Id = idType
	var var0 *securityRulesRsModelProfileSettingObject
	if ans.ProfileSetting != nil {
		var0 = &securityRulesRsModelProfileSettingObject{}
		var0.Group = EncodeStringSlice(ans.ProfileSetting.Group)
	}
	state.Action = types.StringValue(ans.Action)
	state.Application = EncodeStringSlice(ans.Application)
	state.Category = EncodeStringSlice(ans.Category)
	state.Description = types.StringValue(ans.Description)
	state.Destination = EncodeStringSlice(ans.Destination)
	state.DestinationHip = EncodeStringSlice(ans.DestinationHip)
	state.Disabled = types.BoolValue(ans.Disabled)
	state.From = EncodeStringSlice(ans.From)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.LogSetting = types.StringValue(ans.LogSetting)
	state.Name = types.StringValue(ans.Name)
	state.NegateDestination = types.BoolValue(ans.NegateDestination)
	state.NegateSource = types.BoolValue(ans.NegateSource)
	state.ProfileSetting = var0
	state.Service = EncodeStringSlice(ans.Service)
	state.Source = EncodeStringSlice(ans.Source)
	state.SourceHip = EncodeStringSlice(ans.SourceHip)
	state.SourceUser = EncodeStringSlice(ans.SourceUser)
	state.Tag = EncodeStringSlice(ans.Tag)
	state.To = EncodeStringSlice(ans.To)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *securityRulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state securityRulesRsModel
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
		"resource_name":               "sase_security_rules",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := mPRFtcU.NewClient(r.client)
	input := mPRFtcU.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 ffcMtmY.Config
	var0.Action = plan.Action.ValueString()
	var0.Application = DecodeStringSlice(plan.Application)
	var0.Category = DecodeStringSlice(plan.Category)
	var0.Description = plan.Description.ValueString()
	var0.Destination = DecodeStringSlice(plan.Destination)
	var0.DestinationHip = DecodeStringSlice(plan.DestinationHip)
	var0.Disabled = plan.Disabled.ValueBool()
	var0.From = DecodeStringSlice(plan.From)
	var0.LogSetting = plan.LogSetting.ValueString()
	var0.Name = plan.Name.ValueString()
	var0.NegateDestination = plan.NegateDestination.ValueBool()
	var0.NegateSource = plan.NegateSource.ValueBool()
	var var1 *ffcMtmY.ProfileSettingObject
	if plan.ProfileSetting != nil {
		var1 = &ffcMtmY.ProfileSettingObject{}
		var1.Group = DecodeStringSlice(plan.ProfileSetting.Group)
	}
	var0.ProfileSetting = var1
	var0.Service = DecodeStringSlice(plan.Service)
	var0.Source = DecodeStringSlice(plan.Source)
	var0.SourceHip = DecodeStringSlice(plan.SourceHip)
	var0.SourceUser = DecodeStringSlice(plan.SourceUser)
	var0.Tag = DecodeStringSlice(plan.Tag)
	var0.To = DecodeStringSlice(plan.To)
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var2 *securityRulesRsModelProfileSettingObject
	if ans.ProfileSetting != nil {
		var2 = &securityRulesRsModelProfileSettingObject{}
		var2.Group = EncodeStringSlice(ans.ProfileSetting.Group)
	}
	state.Action = types.StringValue(ans.Action)
	state.Application = EncodeStringSlice(ans.Application)
	state.Category = EncodeStringSlice(ans.Category)
	state.Description = types.StringValue(ans.Description)
	state.Destination = EncodeStringSlice(ans.Destination)
	state.DestinationHip = EncodeStringSlice(ans.DestinationHip)
	state.Disabled = types.BoolValue(ans.Disabled)
	state.From = EncodeStringSlice(ans.From)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.LogSetting = types.StringValue(ans.LogSetting)
	state.Name = types.StringValue(ans.Name)
	state.NegateDestination = types.BoolValue(ans.NegateDestination)
	state.NegateSource = types.BoolValue(ans.NegateSource)
	state.ProfileSetting = var2
	state.Service = EncodeStringSlice(ans.Service)
	state.Source = EncodeStringSlice(ans.Source)
	state.SourceHip = EncodeStringSlice(ans.SourceHip)
	state.SourceUser = EncodeStringSlice(ans.SourceUser)
	state.Tag = EncodeStringSlice(ans.Tag)
	state.To = EncodeStringSlice(ans.To)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *securityRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var idType types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &idType)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := idType.ValueString()
	tokens := strings.Split(id, IdSeparator)
	if len(tokens) != 3 {
		resp.Diagnostics.AddError("Error in resource ID format", "Expected 3 tokens")
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource delete", map[string]any{
		"terraform_provider_function": "Delete",
		"resource_name":               "sase_security_rules",
		"locMap":                      map[string]int{"Folder": 1, "ObjectId": 2, "Position": 0},
		"tokens":                      tokens,
	})

	svc := mPRFtcU.NewClient(r.client)
	input := mPRFtcU.DeleteInput{
		ObjectId: tokens[2],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *securityRulesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
