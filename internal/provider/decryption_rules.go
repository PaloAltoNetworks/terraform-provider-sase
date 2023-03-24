package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	fKFKDxk "github.com/paloaltonetworks/sase-go/netsec/schema/decryption/rules"
	vWYSjCE "github.com/paloaltonetworks/sase-go/netsec/service/v1/decryptionrules"

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
	_ datasource.DataSource              = &decryptionRulesListDataSource{}
	_ datasource.DataSourceWithConfigure = &decryptionRulesListDataSource{}
)

func NewDecryptionRulesListDataSource() datasource.DataSource {
	return &decryptionRulesListDataSource{}
}

type decryptionRulesListDataSource struct {
	client *sase.Client
}

type decryptionRulesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit    types.Int64  `tfsdk:"limit"`
	Offset   types.Int64  `tfsdk:"offset"`
	Position types.String `tfsdk:"position"`
	Folder   types.String `tfsdk:"folder"`
	Name     types.String `tfsdk:"name"`

	// Output.
	Data []decryptionRulesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type decryptionRulesListDsModelConfig struct {
	Action            types.String                          `tfsdk:"action"`
	Category          []types.String                        `tfsdk:"category"`
	Description       types.String                          `tfsdk:"description"`
	Destination       []types.String                        `tfsdk:"destination"`
	DestinationHip    []types.String                        `tfsdk:"destination_hip"`
	Disabled          types.Bool                            `tfsdk:"disabled"`
	From              []types.String                        `tfsdk:"from"`
	ObjectId          types.String                          `tfsdk:"object_id"`
	LogFail           types.Bool                            `tfsdk:"log_fail"`
	LogSetting        types.String                          `tfsdk:"log_setting"`
	LogSuccess        types.Bool                            `tfsdk:"log_success"`
	Name              types.String                          `tfsdk:"name"`
	NegateDestination types.Bool                            `tfsdk:"negate_destination"`
	NegateSource      types.Bool                            `tfsdk:"negate_source"`
	Profile           types.String                          `tfsdk:"profile"`
	Service           []types.String                        `tfsdk:"service"`
	Source            []types.String                        `tfsdk:"source"`
	SourceHip         []types.String                        `tfsdk:"source_hip"`
	SourceUser        []types.String                        `tfsdk:"source_user"`
	Tag               []types.String                        `tfsdk:"tag"`
	To                []types.String                        `tfsdk:"to"`
	Type              *decryptionRulesListDsModelTypeObject `tfsdk:"type"`
}

type decryptionRulesListDsModelTypeObject struct {
	SslForwardProxy      types.Bool   `tfsdk:"ssl_forward_proxy"`
	SslInboundInspection types.String `tfsdk:"ssl_inbound_inspection"`
}

// Metadata returns the data source type name.
func (d *decryptionRulesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_decryption_rules_list"
}

// Schema defines the schema for this listing data source.
func (d *decryptionRulesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"position": dsschema.StringAttribute{
				Description:         "The position of a security rule",
				MarkdownDescription: "The position of a security rule",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("pre", "post"),
				},
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
						"action": dsschema.StringAttribute{
							Description:         "The `action` parameter.",
							MarkdownDescription: "The `action` parameter.",
							Computed:            true,
						},
						"category": dsschema.ListAttribute{
							Description:         "The `category` parameter.",
							MarkdownDescription: "The `category` parameter.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"description": dsschema.StringAttribute{
							Description:         "The `description` parameter.",
							MarkdownDescription: "The `description` parameter.",
							Computed:            true,
						},
						"destination": dsschema.ListAttribute{
							Description:         "The `destination` parameter.",
							MarkdownDescription: "The `destination` parameter.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"destination_hip": dsschema.ListAttribute{
							Description:         "The `destination_hip` parameter.",
							MarkdownDescription: "The `destination_hip` parameter.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"disabled": dsschema.BoolAttribute{
							Description:         "The `disabled` parameter.",
							MarkdownDescription: "The `disabled` parameter.",
							Computed:            true,
						},
						"from": dsschema.ListAttribute{
							Description:         "The `from` parameter.",
							MarkdownDescription: "The `from` parameter.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"object_id": dsschema.StringAttribute{
							Description:         "The `object_id` parameter.",
							MarkdownDescription: "The `object_id` parameter.",
							Computed:            true,
						},
						"log_fail": dsschema.BoolAttribute{
							Description:         "The `log_fail` parameter.",
							MarkdownDescription: "The `log_fail` parameter.",
							Computed:            true,
						},
						"log_setting": dsschema.StringAttribute{
							Description:         "The `log_setting` parameter.",
							MarkdownDescription: "The `log_setting` parameter.",
							Computed:            true,
						},
						"log_success": dsschema.BoolAttribute{
							Description:         "The `log_success` parameter.",
							MarkdownDescription: "The `log_success` parameter.",
							Computed:            true,
						},
						"name": dsschema.StringAttribute{
							Description:         "The `name` parameter.",
							MarkdownDescription: "The `name` parameter.",
							Computed:            true,
						},
						"negate_destination": dsschema.BoolAttribute{
							Description:         "The `negate_destination` parameter.",
							MarkdownDescription: "The `negate_destination` parameter.",
							Computed:            true,
						},
						"negate_source": dsschema.BoolAttribute{
							Description:         "The `negate_source` parameter.",
							MarkdownDescription: "The `negate_source` parameter.",
							Computed:            true,
						},
						"profile": dsschema.StringAttribute{
							Description:         "The `profile` parameter.",
							MarkdownDescription: "The `profile` parameter.",
							Computed:            true,
						},
						"service": dsschema.ListAttribute{
							Description:         "The `service` parameter.",
							MarkdownDescription: "The `service` parameter.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"source": dsschema.ListAttribute{
							Description:         "The `source` parameter.",
							MarkdownDescription: "The `source` parameter.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"source_hip": dsschema.ListAttribute{
							Description:         "The `source_hip` parameter.",
							MarkdownDescription: "The `source_hip` parameter.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"source_user": dsschema.ListAttribute{
							Description:         "The `source_user` parameter.",
							MarkdownDescription: "The `source_user` parameter.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"tag": dsschema.ListAttribute{
							Description:         "The `tag` parameter.",
							MarkdownDescription: "The `tag` parameter.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"to": dsschema.ListAttribute{
							Description:         "The `to` parameter.",
							MarkdownDescription: "The `to` parameter.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"type": dsschema.SingleNestedAttribute{
							Description:         "The `type` parameter.",
							MarkdownDescription: "The `type` parameter.",
							Computed:            true,
							Attributes: map[string]dsschema.Attribute{
								"ssl_forward_proxy": dsschema.BoolAttribute{
									Description:         "The `ssl_forward_proxy` parameter.",
									MarkdownDescription: "The `ssl_forward_proxy` parameter.",
									Computed:            true,
								},
								"ssl_inbound_inspection": dsschema.StringAttribute{
									Description:         "The `ssl_inbound_inspection` parameter.",
									MarkdownDescription: "The `ssl_inbound_inspection` parameter.",
									Computed:            true,
								},
							},
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
func (d *decryptionRulesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *decryptionRulesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state decryptionRulesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_decryption_rules_list",
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
	svc := vWYSjCE.NewClient(d.client)
	input := vWYSjCE.ListInput{
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
	var var0 []decryptionRulesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]decryptionRulesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 decryptionRulesListDsModelConfig
			var var3 *decryptionRulesListDsModelTypeObject
			if var1.Type != nil {
				var3 = &decryptionRulesListDsModelTypeObject{}
				if var1.Type.SslForwardProxy != nil {
					var3.SslForwardProxy = types.BoolValue(true)
				}
				var3.SslInboundInspection = types.StringValue(var1.Type.SslInboundInspection)
			}
			var2.Action = types.StringValue(var1.Action)
			var2.Category = EncodeStringSlice(var1.Category)
			var2.Description = types.StringValue(var1.Description)
			var2.Destination = EncodeStringSlice(var1.Destination)
			var2.DestinationHip = EncodeStringSlice(var1.DestinationHip)
			var2.Disabled = types.BoolValue(var1.Disabled)
			var2.From = EncodeStringSlice(var1.From)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.LogFail = types.BoolValue(var1.LogFail)
			var2.LogSetting = types.StringValue(var1.LogSetting)
			var2.LogSuccess = types.BoolValue(var1.LogSuccess)
			var2.Name = types.StringValue(var1.Name)
			var2.NegateDestination = types.BoolValue(var1.NegateDestination)
			var2.NegateSource = types.BoolValue(var1.NegateSource)
			var2.Profile = types.StringValue(var1.Profile)
			var2.Service = EncodeStringSlice(var1.Service)
			var2.Source = EncodeStringSlice(var1.Source)
			var2.SourceHip = EncodeStringSlice(var1.SourceHip)
			var2.SourceUser = EncodeStringSlice(var1.SourceUser)
			var2.Tag = EncodeStringSlice(var1.Tag)
			var2.To = EncodeStringSlice(var1.To)
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
	_ datasource.DataSource              = &decryptionRulesDataSource{}
	_ datasource.DataSourceWithConfigure = &decryptionRulesDataSource{}
)

func NewDecryptionRulesDataSource() datasource.DataSource {
	return &decryptionRulesDataSource{}
}

type decryptionRulesDataSource struct {
	client *sase.Client
}

type decryptionRulesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`
	Folder   types.String `tfsdk:"folder"`

	// Output.
	// Ref: #/components/schemas/decryption-rules
	Action         types.String   `tfsdk:"action"`
	Category       []types.String `tfsdk:"category"`
	Description    types.String   `tfsdk:"description"`
	Destination    []types.String `tfsdk:"destination"`
	DestinationHip []types.String `tfsdk:"destination_hip"`
	Disabled       types.Bool     `tfsdk:"disabled"`
	From           []types.String `tfsdk:"from"`
	// input omit: ObjectId
	LogFail           types.Bool                        `tfsdk:"log_fail"`
	LogSetting        types.String                      `tfsdk:"log_setting"`
	LogSuccess        types.Bool                        `tfsdk:"log_success"`
	Name              types.String                      `tfsdk:"name"`
	NegateDestination types.Bool                        `tfsdk:"negate_destination"`
	NegateSource      types.Bool                        `tfsdk:"negate_source"`
	Profile           types.String                      `tfsdk:"profile"`
	Service           []types.String                    `tfsdk:"service"`
	Source            []types.String                    `tfsdk:"source"`
	SourceHip         []types.String                    `tfsdk:"source_hip"`
	SourceUser        []types.String                    `tfsdk:"source_user"`
	Tag               []types.String                    `tfsdk:"tag"`
	To                []types.String                    `tfsdk:"to"`
	Type              *decryptionRulesDsModelTypeObject `tfsdk:"type"`
}

type decryptionRulesDsModelTypeObject struct {
	SslForwardProxy      types.Bool   `tfsdk:"ssl_forward_proxy"`
	SslInboundInspection types.String `tfsdk:"ssl_inbound_inspection"`
}

// Metadata returns the data source type name.
func (d *decryptionRulesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_decryption_rules"
}

// Schema defines the schema for this listing data source.
func (d *decryptionRulesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"action": dsschema.StringAttribute{
				Description:         "The `action` parameter.",
				MarkdownDescription: "The `action` parameter.",
				Computed:            true,
			},
			"category": dsschema.ListAttribute{
				Description:         "The `category` parameter.",
				MarkdownDescription: "The `category` parameter.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"description": dsschema.StringAttribute{
				Description:         "The `description` parameter.",
				MarkdownDescription: "The `description` parameter.",
				Computed:            true,
			},
			"destination": dsschema.ListAttribute{
				Description:         "The `destination` parameter.",
				MarkdownDescription: "The `destination` parameter.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"destination_hip": dsschema.ListAttribute{
				Description:         "The `destination_hip` parameter.",
				MarkdownDescription: "The `destination_hip` parameter.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"disabled": dsschema.BoolAttribute{
				Description:         "The `disabled` parameter.",
				MarkdownDescription: "The `disabled` parameter.",
				Computed:            true,
			},
			"from": dsschema.ListAttribute{
				Description:         "The `from` parameter.",
				MarkdownDescription: "The `from` parameter.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"log_fail": dsschema.BoolAttribute{
				Description:         "The `log_fail` parameter.",
				MarkdownDescription: "The `log_fail` parameter.",
				Computed:            true,
			},
			"log_setting": dsschema.StringAttribute{
				Description:         "The `log_setting` parameter.",
				MarkdownDescription: "The `log_setting` parameter.",
				Computed:            true,
			},
			"log_success": dsschema.BoolAttribute{
				Description:         "The `log_success` parameter.",
				MarkdownDescription: "The `log_success` parameter.",
				Computed:            true,
			},
			"name": dsschema.StringAttribute{
				Description:         "The `name` parameter.",
				MarkdownDescription: "The `name` parameter.",
				Computed:            true,
			},
			"negate_destination": dsschema.BoolAttribute{
				Description:         "The `negate_destination` parameter.",
				MarkdownDescription: "The `negate_destination` parameter.",
				Computed:            true,
			},
			"negate_source": dsschema.BoolAttribute{
				Description:         "The `negate_source` parameter.",
				MarkdownDescription: "The `negate_source` parameter.",
				Computed:            true,
			},
			"profile": dsschema.StringAttribute{
				Description:         "The `profile` parameter.",
				MarkdownDescription: "The `profile` parameter.",
				Computed:            true,
			},
			"service": dsschema.ListAttribute{
				Description:         "The `service` parameter.",
				MarkdownDescription: "The `service` parameter.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"source": dsschema.ListAttribute{
				Description:         "The `source` parameter.",
				MarkdownDescription: "The `source` parameter.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"source_hip": dsschema.ListAttribute{
				Description:         "The `source_hip` parameter.",
				MarkdownDescription: "The `source_hip` parameter.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"source_user": dsschema.ListAttribute{
				Description:         "The `source_user` parameter.",
				MarkdownDescription: "The `source_user` parameter.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"tag": dsschema.ListAttribute{
				Description:         "The `tag` parameter.",
				MarkdownDescription: "The `tag` parameter.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"to": dsschema.ListAttribute{
				Description:         "The `to` parameter.",
				MarkdownDescription: "The `to` parameter.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"type": dsschema.SingleNestedAttribute{
				Description:         "The `type` parameter.",
				MarkdownDescription: "The `type` parameter.",
				Computed:            true,
				Attributes: map[string]dsschema.Attribute{
					"ssl_forward_proxy": dsschema.BoolAttribute{
						Description:         "The `ssl_forward_proxy` parameter.",
						MarkdownDescription: "The `ssl_forward_proxy` parameter.",
						Computed:            true,
					},
					"ssl_inbound_inspection": dsschema.StringAttribute{
						Description:         "The `ssl_inbound_inspection` parameter.",
						MarkdownDescription: "The `ssl_inbound_inspection` parameter.",
						Computed:            true,
					},
				},
			},
		},
	}
}

// Configure prepares the struct.
func (d *decryptionRulesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *decryptionRulesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state decryptionRulesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_decryption_rules",
		"object_id":                   state.ObjectId.ValueString(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := vWYSjCE.NewClient(d.client)
	input := vWYSjCE.ReadInput{
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
	var var0 *decryptionRulesDsModelTypeObject
	if ans.Type != nil {
		var0 = &decryptionRulesDsModelTypeObject{}
		if ans.Type.SslForwardProxy != nil {
			var0.SslForwardProxy = types.BoolValue(true)
		}
		var0.SslInboundInspection = types.StringValue(ans.Type.SslInboundInspection)
	}
	state.Action = types.StringValue(ans.Action)
	state.Category = EncodeStringSlice(ans.Category)
	state.Description = types.StringValue(ans.Description)
	state.Destination = EncodeStringSlice(ans.Destination)
	state.DestinationHip = EncodeStringSlice(ans.DestinationHip)
	state.Disabled = types.BoolValue(ans.Disabled)
	state.From = EncodeStringSlice(ans.From)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.LogFail = types.BoolValue(ans.LogFail)
	state.LogSetting = types.StringValue(ans.LogSetting)
	state.LogSuccess = types.BoolValue(ans.LogSuccess)
	state.Name = types.StringValue(ans.Name)
	state.NegateDestination = types.BoolValue(ans.NegateDestination)
	state.NegateSource = types.BoolValue(ans.NegateSource)
	state.Profile = types.StringValue(ans.Profile)
	state.Service = EncodeStringSlice(ans.Service)
	state.Source = EncodeStringSlice(ans.Source)
	state.SourceHip = EncodeStringSlice(ans.SourceHip)
	state.SourceUser = EncodeStringSlice(ans.SourceUser)
	state.Tag = EncodeStringSlice(ans.Tag)
	state.To = EncodeStringSlice(ans.To)
	state.Type = var0

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &decryptionRulesResource{}
	_ resource.ResourceWithConfigure   = &decryptionRulesResource{}
	_ resource.ResourceWithImportState = &decryptionRulesResource{}
)

func NewDecryptionRulesResource() resource.Resource {
	return &decryptionRulesResource{}
}

type decryptionRulesResource struct {
	client *sase.Client
}

type decryptionRulesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Position types.String `tfsdk:"position"`
	Folder   types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/decryption-rules
	Action            types.String                      `tfsdk:"action"`
	Category          []types.String                    `tfsdk:"category"`
	Description       types.String                      `tfsdk:"description"`
	Destination       []types.String                    `tfsdk:"destination"`
	DestinationHip    []types.String                    `tfsdk:"destination_hip"`
	Disabled          types.Bool                        `tfsdk:"disabled"`
	From              []types.String                    `tfsdk:"from"`
	ObjectId          types.String                      `tfsdk:"object_id"`
	LogFail           types.Bool                        `tfsdk:"log_fail"`
	LogSetting        types.String                      `tfsdk:"log_setting"`
	LogSuccess        types.Bool                        `tfsdk:"log_success"`
	Name              types.String                      `tfsdk:"name"`
	NegateDestination types.Bool                        `tfsdk:"negate_destination"`
	NegateSource      types.Bool                        `tfsdk:"negate_source"`
	Profile           types.String                      `tfsdk:"profile"`
	Service           []types.String                    `tfsdk:"service"`
	Source            []types.String                    `tfsdk:"source"`
	SourceHip         []types.String                    `tfsdk:"source_hip"`
	SourceUser        []types.String                    `tfsdk:"source_user"`
	Tag               []types.String                    `tfsdk:"tag"`
	To                []types.String                    `tfsdk:"to"`
	Type              *decryptionRulesRsModelTypeObject `tfsdk:"type"`
}

type decryptionRulesRsModelTypeObject struct {
	SslForwardProxy      types.Bool   `tfsdk:"ssl_forward_proxy"`
	SslInboundInspection types.String `tfsdk:"ssl_inbound_inspection"`
}

// Metadata returns the data source type name.
func (r *decryptionRulesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_decryption_rules"
}

// Schema defines the schema for this listing data source.
func (r *decryptionRulesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"position": rsschema.StringAttribute{
				Description:         "The position of a security rule",
				MarkdownDescription: "The position of a security rule",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("pre", "post"),
				},
			},
			"folder": rsschema.StringAttribute{
				Description:         "The folder of the entry",
				MarkdownDescription: "The folder of the entry",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Service Connections", "Mobile Users Container", "Mobile Users Explicit Proxy"),
				},
			},

			"action": rsschema.StringAttribute{
				Description:         "The `action` parameter.",
				MarkdownDescription: "The `action` parameter.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("decrypt", "no-decrypt"),
				},
			},
			"category": rsschema.ListAttribute{
				Description:         "The `category` parameter.",
				MarkdownDescription: "The `category` parameter.",
				Required:            true,
				ElementType:         types.StringType,
			},
			"description": rsschema.StringAttribute{
				Description:         "The `description` parameter.",
				MarkdownDescription: "The `description` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
			},
			"destination": rsschema.ListAttribute{
				Description:         "The `destination` parameter.",
				MarkdownDescription: "The `destination` parameter.",
				Required:            true,
				ElementType:         types.StringType,
			},
			"destination_hip": rsschema.ListAttribute{
				Description:         "The `destination_hip` parameter.",
				MarkdownDescription: "The `destination_hip` parameter.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"disabled": rsschema.BoolAttribute{
				Description:         "The `disabled` parameter.",
				MarkdownDescription: "The `disabled` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"from": rsschema.ListAttribute{
				Description:         "The `from` parameter.",
				MarkdownDescription: "The `from` parameter.",
				Required:            true,
				ElementType:         types.StringType,
			},
			"object_id": rsschema.StringAttribute{
				Description:         "The `object_id` parameter.",
				MarkdownDescription: "The `object_id` parameter.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"log_fail": rsschema.BoolAttribute{
				Description:         "The `log_fail` parameter.",
				MarkdownDescription: "The `log_fail` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"log_setting": rsschema.StringAttribute{
				Description:         "The `log_setting` parameter.",
				MarkdownDescription: "The `log_setting` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
			},
			"log_success": rsschema.BoolAttribute{
				Description:         "The `log_success` parameter.",
				MarkdownDescription: "The `log_success` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"name": rsschema.StringAttribute{
				Description:         "The `name` parameter.",
				MarkdownDescription: "The `name` parameter.",
				Required:            true,
			},
			"negate_destination": rsschema.BoolAttribute{
				Description:         "The `negate_destination` parameter.",
				MarkdownDescription: "The `negate_destination` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"negate_source": rsschema.BoolAttribute{
				Description:         "The `negate_source` parameter.",
				MarkdownDescription: "The `negate_source` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					DefaultBool(false),
				},
			},
			"profile": rsschema.StringAttribute{
				Description:         "The `profile` parameter.",
				MarkdownDescription: "The `profile` parameter.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
			},
			"service": rsschema.ListAttribute{
				Description:         "The `service` parameter.",
				MarkdownDescription: "The `service` parameter.",
				Required:            true,
				ElementType:         types.StringType,
			},
			"source": rsschema.ListAttribute{
				Description:         "The `source` parameter.",
				MarkdownDescription: "The `source` parameter.",
				Required:            true,
				ElementType:         types.StringType,
			},
			"source_hip": rsschema.ListAttribute{
				Description:         "The `source_hip` parameter.",
				MarkdownDescription: "The `source_hip` parameter.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"source_user": rsschema.ListAttribute{
				Description:         "The `source_user` parameter.",
				MarkdownDescription: "The `source_user` parameter.",
				Required:            true,
				ElementType:         types.StringType,
			},
			"tag": rsschema.ListAttribute{
				Description:         "The `tag` parameter.",
				MarkdownDescription: "The `tag` parameter.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"to": rsschema.ListAttribute{
				Description:         "The `to` parameter.",
				MarkdownDescription: "The `to` parameter.",
				Required:            true,
				ElementType:         types.StringType,
			},
			"type": rsschema.SingleNestedAttribute{
				Description:         "The `type` parameter.",
				MarkdownDescription: "The `type` parameter.",
				Optional:            true,
				Attributes: map[string]rsschema.Attribute{
					"ssl_forward_proxy": rsschema.BoolAttribute{
						Description:         "The `ssl_forward_proxy` parameter.",
						MarkdownDescription: "The `ssl_forward_proxy` parameter.",
						Optional:            true,
					},
					"ssl_inbound_inspection": rsschema.StringAttribute{
						Description:         "The `ssl_inbound_inspection` parameter.",
						MarkdownDescription: "The `ssl_inbound_inspection` parameter.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							DefaultString(""),
						},
						Validators: []validator.String{
							stringvalidator.ConflictsWith(),
						},
					},
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *decryptionRulesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *decryptionRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state decryptionRulesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_decryption_rules",
		"position":                    state.Position.ValueString(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := vWYSjCE.NewClient(r.client)
	input := vWYSjCE.CreateInput{
		Position: state.Position.ValueString(),
		Folder:   state.Folder.ValueString(),
	}
	var var0 fKFKDxk.Config
	var0.Action = state.Action.ValueString()
	var0.Category = DecodeStringSlice(state.Category)
	var0.Description = state.Description.ValueString()
	var0.Destination = DecodeStringSlice(state.Destination)
	var0.DestinationHip = DecodeStringSlice(state.DestinationHip)
	var0.Disabled = state.Disabled.ValueBool()
	var0.From = DecodeStringSlice(state.From)
	var0.LogFail = state.LogFail.ValueBool()
	var0.LogSetting = state.LogSetting.ValueString()
	var0.LogSuccess = state.LogSuccess.ValueBool()
	var0.Name = state.Name.ValueString()
	var0.NegateDestination = state.NegateDestination.ValueBool()
	var0.NegateSource = state.NegateSource.ValueBool()
	var0.Profile = state.Profile.ValueString()
	var0.Service = DecodeStringSlice(state.Service)
	var0.Source = DecodeStringSlice(state.Source)
	var0.SourceHip = DecodeStringSlice(state.SourceHip)
	var0.SourceUser = DecodeStringSlice(state.SourceUser)
	var0.Tag = DecodeStringSlice(state.Tag)
	var0.To = DecodeStringSlice(state.To)
	var var1 *fKFKDxk.TypeObject
	if state.Type != nil {
		var1 = &fKFKDxk.TypeObject{}
		if state.Type.SslForwardProxy.ValueBool() {
			var1.SslForwardProxy = struct{}{}
		}
		var1.SslInboundInspection = state.Type.SslInboundInspection.ValueString()
	}
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
	idBuilder.WriteString(input.Position)
	idBuilder.WriteString(IdSeparator)
	idBuilder.WriteString(input.Folder)
	idBuilder.WriteString(IdSeparator)
	idBuilder.WriteString(ans.ObjectId)
	state.Id = types.StringValue(idBuilder.String())
	var var2 *decryptionRulesRsModelTypeObject
	if ans.Type != nil {
		var2 = &decryptionRulesRsModelTypeObject{}
		if ans.Type.SslForwardProxy != nil {
			var2.SslForwardProxy = types.BoolValue(true)
		}
		var2.SslInboundInspection = types.StringValue(ans.Type.SslInboundInspection)
	}
	state.Action = types.StringValue(ans.Action)
	state.Category = EncodeStringSlice(ans.Category)
	state.Description = types.StringValue(ans.Description)
	state.Destination = EncodeStringSlice(ans.Destination)
	state.DestinationHip = EncodeStringSlice(ans.DestinationHip)
	state.Disabled = types.BoolValue(ans.Disabled)
	state.From = EncodeStringSlice(ans.From)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.LogFail = types.BoolValue(ans.LogFail)
	state.LogSetting = types.StringValue(ans.LogSetting)
	state.LogSuccess = types.BoolValue(ans.LogSuccess)
	state.Name = types.StringValue(ans.Name)
	state.NegateDestination = types.BoolValue(ans.NegateDestination)
	state.NegateSource = types.BoolValue(ans.NegateSource)
	state.Profile = types.StringValue(ans.Profile)
	state.Service = EncodeStringSlice(ans.Service)
	state.Source = EncodeStringSlice(ans.Source)
	state.SourceHip = EncodeStringSlice(ans.SourceHip)
	state.SourceUser = EncodeStringSlice(ans.SourceUser)
	state.Tag = EncodeStringSlice(ans.Tag)
	state.To = EncodeStringSlice(ans.To)
	state.Type = var2

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *decryptionRulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state decryptionRulesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_decryption_rules",
		"locMap":                      map[string]int{"Folder": 1, "ObjectId": 2, "Position": 0},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := vWYSjCE.NewClient(r.client)
	input := vWYSjCE.ReadInput{
		ObjectId: tokens[2],
		Folder:   tokens[1],
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
	var var0 *decryptionRulesRsModelTypeObject
	if ans.Type != nil {
		var0 = &decryptionRulesRsModelTypeObject{}
		if ans.Type.SslForwardProxy != nil {
			var0.SslForwardProxy = types.BoolValue(true)
		}
		var0.SslInboundInspection = types.StringValue(ans.Type.SslInboundInspection)
	}
	state.Action = types.StringValue(ans.Action)
	state.Category = EncodeStringSlice(ans.Category)
	state.Description = types.StringValue(ans.Description)
	state.Destination = EncodeStringSlice(ans.Destination)
	state.DestinationHip = EncodeStringSlice(ans.DestinationHip)
	state.Disabled = types.BoolValue(ans.Disabled)
	state.From = EncodeStringSlice(ans.From)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.LogFail = types.BoolValue(ans.LogFail)
	state.LogSetting = types.StringValue(ans.LogSetting)
	state.LogSuccess = types.BoolValue(ans.LogSuccess)
	state.Name = types.StringValue(ans.Name)
	state.NegateDestination = types.BoolValue(ans.NegateDestination)
	state.NegateSource = types.BoolValue(ans.NegateSource)
	state.Profile = types.StringValue(ans.Profile)
	state.Service = EncodeStringSlice(ans.Service)
	state.Source = EncodeStringSlice(ans.Source)
	state.SourceHip = EncodeStringSlice(ans.SourceHip)
	state.SourceUser = EncodeStringSlice(ans.SourceUser)
	state.Tag = EncodeStringSlice(ans.Tag)
	state.To = EncodeStringSlice(ans.To)
	state.Type = var0

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *decryptionRulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state decryptionRulesRsModel
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
		"resource_name":               "sase_decryption_rules",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := vWYSjCE.NewClient(r.client)
	input := vWYSjCE.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 fKFKDxk.Config
	var0.Action = plan.Action.ValueString()
	var0.Category = DecodeStringSlice(plan.Category)
	var0.Description = plan.Description.ValueString()
	var0.Destination = DecodeStringSlice(plan.Destination)
	var0.DestinationHip = DecodeStringSlice(plan.DestinationHip)
	var0.Disabled = plan.Disabled.ValueBool()
	var0.From = DecodeStringSlice(plan.From)
	var0.LogFail = plan.LogFail.ValueBool()
	var0.LogSetting = plan.LogSetting.ValueString()
	var0.LogSuccess = plan.LogSuccess.ValueBool()
	var0.Name = plan.Name.ValueString()
	var0.NegateDestination = plan.NegateDestination.ValueBool()
	var0.NegateSource = plan.NegateSource.ValueBool()
	var0.Profile = plan.Profile.ValueString()
	var0.Service = DecodeStringSlice(plan.Service)
	var0.Source = DecodeStringSlice(plan.Source)
	var0.SourceHip = DecodeStringSlice(plan.SourceHip)
	var0.SourceUser = DecodeStringSlice(plan.SourceUser)
	var0.Tag = DecodeStringSlice(plan.Tag)
	var0.To = DecodeStringSlice(plan.To)
	var var1 *fKFKDxk.TypeObject
	if plan.Type != nil {
		var1 = &fKFKDxk.TypeObject{}
		if plan.Type.SslForwardProxy.ValueBool() {
			var1.SslForwardProxy = struct{}{}
		}
		var1.SslInboundInspection = plan.Type.SslInboundInspection.ValueString()
	}
	var0.Type = var1
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var2 *decryptionRulesRsModelTypeObject
	if ans.Type != nil {
		var2 = &decryptionRulesRsModelTypeObject{}
		if ans.Type.SslForwardProxy != nil {
			var2.SslForwardProxy = types.BoolValue(true)
		}
		var2.SslInboundInspection = types.StringValue(ans.Type.SslInboundInspection)
	}
	state.Action = types.StringValue(ans.Action)
	state.Category = EncodeStringSlice(ans.Category)
	state.Description = types.StringValue(ans.Description)
	state.Destination = EncodeStringSlice(ans.Destination)
	state.DestinationHip = EncodeStringSlice(ans.DestinationHip)
	state.Disabled = types.BoolValue(ans.Disabled)
	state.From = EncodeStringSlice(ans.From)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.LogFail = types.BoolValue(ans.LogFail)
	state.LogSetting = types.StringValue(ans.LogSetting)
	state.LogSuccess = types.BoolValue(ans.LogSuccess)
	state.Name = types.StringValue(ans.Name)
	state.NegateDestination = types.BoolValue(ans.NegateDestination)
	state.NegateSource = types.BoolValue(ans.NegateSource)
	state.Profile = types.StringValue(ans.Profile)
	state.Service = EncodeStringSlice(ans.Service)
	state.Source = EncodeStringSlice(ans.Source)
	state.SourceHip = EncodeStringSlice(ans.SourceHip)
	state.SourceUser = EncodeStringSlice(ans.SourceUser)
	state.Tag = EncodeStringSlice(ans.Tag)
	state.To = EncodeStringSlice(ans.To)
	state.Type = var2

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *decryptionRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"resource_name":               "sase_decryption_rules",
		"locMap":                      map[string]int{"Folder": 1, "ObjectId": 2, "Position": 0},
		"tokens":                      tokens,
	})

	svc := vWYSjCE.NewClient(r.client)
	input := vWYSjCE.DeleteInput{
		ObjectId: tokens[2],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *decryptionRulesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
