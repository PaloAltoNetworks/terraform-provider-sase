package provider

import (
    "context"
    "fmt"
    "strings"

    "github.com/paloaltonetworks/sase-go"
    "github.com/paloaltonetworks/sase-go/api"
    addr "github.com/paloaltonetworks/sase-go/netsec/schema/objects/addresses"
    service "github.com/paloaltonetworks/sase-go/netsec/service/v1/addresses"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    _ "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Data source.
var (
    _ datasource.DataSource = &AddressObjectListDataSource{}
    _ datasource.DataSourceWithConfigure = &AddressObjectListDataSource{}
)

func NewAddressObjectListDataSource() datasource.DataSource {
    return &AddressObjectListDataSource{}
}

type AddressObjectListDataSource struct {
    client *sase.Client
}

type AddressObjectListDsModel struct {
    // Input.
    Limit types.Int64 `tfsdk:"limit"`
    Offset types.Int64 `tfsdk:"offset"`
    Name types.String `tfsdk:"name"`
    Folder types.String `tfsdk:"folder"`

    // Output.
    Data []AddressObjectListItemDsModel `tfsdk:"data"`
    Total types.Int64 `tfsdk:"total"`
}

type AddressObjectListItemDsModel struct {
    Description types.String   `tfsdk:"description"`
    Fqdn        types.String   `tfsdk:"fqdn"`
    ObjectId    types.String   `tfsdk:"object_id"`
    IpNetmask   types.String   `tfsdk:"ip_netmask"`
    IpRange     types.String   `tfsdk:"ip_range"`
    IpWildcard  types.String   `tfsdk:"ip_wildcard"`
    Name        types.String   `tfsdk:"name"`
    Tags        []types.String `tfsdk:"tags"`
    Type        types.String   `tfsdk:"type"`
}

// Metadata returns the data source type name.
func (d *AddressObjectListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_address_object_list"
}

// Schema defines the schema for the data source.
func (d *AddressObjectListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = dsschema.Schema{
        Description: "Get a list of address objects.",

        Attributes: map[string] dsschema.Attribute{
            "limit": dsschema.Int64Attribute{
                Description: "Limit of items to return.  `-1` returns everything.",
                Optional: true,
            },
            "offset": dsschema.Int64Attribute{
                Description: "Paging offset.",
                Optional: true,
            },
            "name": dsschema.StringAttribute{
                Description: "The name.",
                Optional: true,
            },
            "folder": dsschema.StringAttribute{
                Description: "The folder.",
                Required: true,
            },
            "total": dsschema.Int64Attribute{
                Description: "Size of data.",
                Computed: true,
            },
            "data": dsschema.ListNestedAttribute{
                Description: "The data output.",
                Computed: true,
                NestedObject: dsschema.NestedAttributeObject{
                    Attributes: map[string] dsschema.Attribute{
                        "object_id": dsschema.StringAttribute{
                            Description: "The UUID.",
                            Computed: true,
                        },
                        "name": dsschema.StringAttribute{
                            Description: "The name.",
                            Computed: true,
                        },
                        "fqdn": dsschema.StringAttribute{
                            Description: "FQDN type: the FQDN.",
                            Computed: true,
                        },
                        "ip_netmask": dsschema.StringAttribute{
                            Description: "IP Netmask type: the IP/netmask.",
                            Computed: true,
                        },
                        "ip_range": dsschema.StringAttribute{
                            Description: "IP range type: the IP range.",
                            Computed: true,
                        },
                        "ip_wildcard": dsschema.StringAttribute{
                            Description: "IP wildcard type: the IP wildcard.",
                            Computed: true,
                        },
                        "type": dsschema.StringAttribute{
                            Description: "The type.",
                            Computed: true,
                        },
                        "description": dsschema.StringAttribute{
                            Description: "The description.",
                            Computed: true,
                        },
                        "tags": dsschema.ListAttribute{
                            Description: "The list of tags.",
                            ElementType: types.StringType,
                            Computed: true,
                        },
                    },
                },
            },
        },
    }
}

func (d *AddressObjectListDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    d.client = req.ProviderData.(*sase.Client)
}

func (d *AddressObjectListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var state AddressObjectListDsModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    tflog.Info(ctx, "get address object list", map[string] any{
        "folder": state.Folder.ValueString(),
        "name": state.Name.ValueString(),
        "has_name": !state.Name.IsNull(),
        "limit": state.Limit.ValueInt64(),
        "has_limit": !state.Limit.IsNull(),
        "offset": state.Offset.ValueInt64(),
        "has_offset": !state.Offset.IsNull(),
    })
    svc := service.NewClient(d.client)

    input := service.ListInput{
        Folder: state.Folder.ValueString(),
    }
    if !state.Name.IsNull() {
        input.Name = api.String(state.Name.ValueString())
    }
    if !state.Limit.IsNull() {
        input.Limit = api.Int(state.Limit.ValueInt64())
    }
    if !state.Offset.IsNull() {
        input.Offset = api.Int(state.Limit.ValueInt64())
    }

    ans, err := svc.List(ctx, input)
    if err != nil {
        resp.Diagnostics.AddError("Error getting listing", err.Error())
        return
    }

    list := make([]AddressObjectListItemDsModel, 0, len(ans.Data))
    for _, x := range ans.Data {
        list = append(list, AddressObjectListItemDsModel{
            ObjectId: types.StringValue(x.ObjectId),
            Name: types.StringValue(x.Name),
            Fqdn: types.StringValue(x.Fqdn),
            IpNetmask: types.StringValue(x.IpNetmask),
            IpRange: types.StringValue(x.IpRange),
            IpWildcard: types.StringValue(x.IpWildcard),
            Description: types.StringValue(x.Description),
            Tags: EncodeStringSlice(x.Tag),
            Type: types.StringValue(x.Type),
        })
    }
    state.Data = list
    state.Total = types.Int64Value(ans.Total)

    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
    _ resource.Resource = &AddressObjectResource{}
    _ resource.ResourceWithConfigure = &AddressObjectResource{}
    _ resource.ResourceWithImportState = &AddressObjectResource{}

)

func NewAddressObjectResource() resource.Resource {
    return &AddressObjectResource{}
}

type AddressObjectResource struct {
    client *sase.Client
}

type AddressObjectRsModel struct {
    Id types.String `tfsdk:"id"`
    Folder types.String `tfsdk:"folder"`
    Description types.String   `tfsdk:"description"`
    Fqdn        types.String   `tfsdk:"fqdn"`
    ObjectId    types.String   `tfsdk:"object_id"`
    IpNetmask   types.String   `tfsdk:"ip_netmask"`
    IpRange     types.String   `tfsdk:"ip_range"`
    IpWildcard  types.String   `tfsdk:"ip_wildcard"`
    Name        types.String   `tfsdk:"name"`
    Tags        []types.String `tfsdk:"tags"`
    Type        types.String   `tfsdk:"type"`
}

func (r *AddressObjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_address_object"
}

func (r *AddressObjectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "Manage an address object.",

        Attributes: map[string] schema.Attribute{
            "id": schema.StringAttribute{
                Description: "The Terraform resource ID.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "object_id": schema.StringAttribute{
                Description: "The UUID.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "folder": schema.StringAttribute{
                Description: "The folder.",
                Required: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "name": schema.StringAttribute{
                Description: "The name.",
                Required: true,
            },
            "fqdn": schema.StringAttribute{
                Description: "FQDN type: the FQDN.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    DefaultString(""),
                },
            },
            "ip_netmask": schema.StringAttribute{
                Description: "IP Netmask type: the IP/netmask.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    DefaultString(""),
                },
            },
            "ip_range": schema.StringAttribute{
                Description: "IP range type: the IP range.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    DefaultString(""),
                },
            },
            "ip_wildcard": schema.StringAttribute{
                Description: "IP wildcard type: the IP wildcard.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    DefaultString(""),
                },
            },
            "type": schema.StringAttribute{
                Description: "The type.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "description": schema.StringAttribute{
                Description: "The description.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    DefaultString(""),
                },
            },
            "tags": schema.ListAttribute{
                Description: "The list of tags.",
                ElementType: types.StringType,
                Optional: true,
            },
        },
    }
}

func (r *AddressObjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

    r.client = req.ProviderData.(*sase.Client)
}

func (r *AddressObjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var state AddressObjectRsModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    tflog.Info(ctx, "create address object", map[string] any{
        "folder": state.Folder.ValueString(),
        "name": state.Name.ValueString(),
    })
    svc := service.NewClient(r.client)

    input := service.CreateInput{
        Folder: state.Folder.ValueString(),
        Config: addr.Config{
            Name: state.Name.ValueString(),
            Fqdn: state.Fqdn.ValueString(),
            IpNetmask: state.IpNetmask.ValueString(),
            IpRange: state.IpRange.ValueString(),
            IpWildcard: state.IpWildcard.ValueString(),
            Description: state.Description.ValueString(),
            Tag: DecodeStringSlice(state.Tags),
        },
    }

    ans, err := svc.Create(ctx, input)
    if err != nil {
        resp.Diagnostics.AddError("Error creating config", err.Error())
        return
    }

    endstate := AddressObjectRsModel{
        Id: types.StringValue(fmt.Sprintf("%s:%s", state.Folder.ValueString(), ans.ObjectId)),
        ObjectId: types.StringValue(ans.ObjectId),
        Folder: types.StringValue(input.Folder),
        Name: types.StringValue(ans.Name),
        Fqdn: types.StringValue(ans.Fqdn),
        IpNetmask: types.StringValue(ans.IpNetmask),
        IpRange: types.StringValue(ans.IpRange),
        IpWildcard: types.StringValue(ans.IpWildcard),
        Description: types.StringValue(ans.Description),
        Tags: EncodeStringSlice(ans.Tag),
        Type: types.StringValue(ans.Type),
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &endstate)...)
}

func (r *AddressObjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var idType types.String
    resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &idType)...)
    if resp.Diagnostics.HasError() {
        return
    }
    id := idType.ValueString()
    tokens := strings.Split(id, IdSeparator)
    if len(tokens) != 2 {
        resp.Diagnostics.AddError("Error in resource ID", "Expected 2 tokens")
        return
    }

    tflog.Info(ctx, "read address object", map[string] any{
        "folder": tokens[0],
        "uuid": tokens[1],
    })
    svc := service.NewClient(r.client)

    input := service.ReadInput{
        Folder: tokens[0],
        ObjectId: tokens[1],
    }

    ans, err := svc.Read(ctx, input)
    if err != nil {
        if IsObjectNotFound(err) {
            resp.State.RemoveResource(ctx)
        } else {
            resp.Diagnostics.AddError("Error reading config", err.Error())
        }
        return
    }

    endstate := AddressObjectRsModel{
        Id: types.StringValue(id),
        ObjectId: types.StringValue(ans.ObjectId),
        Folder: types.StringValue(input.Folder),
        Name: types.StringValue(ans.Name),
        Fqdn: types.StringValue(ans.Fqdn),
        IpNetmask: types.StringValue(ans.IpNetmask),
        IpRange: types.StringValue(ans.IpRange),
        IpWildcard: types.StringValue(ans.IpWildcard),
        Description: types.StringValue(ans.Description),
        Tags: EncodeStringSlice(ans.Tag),
        Type: types.StringValue(ans.Type),
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &endstate)...)
}

func (r *AddressObjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var idType types.String
    resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &idType)...)
    if resp.Diagnostics.HasError() {
        return
    }
    id := idType.ValueString()
    tokens := strings.Split(id, IdSeparator)
    if len(tokens) != 2 {
        resp.Diagnostics.AddError("Error in resource ID", "Expected 2 tokens")
        return
    }

    var state AddressObjectRsModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    tflog.Info(ctx, "updating address object", map[string] any{
        "folder": tokens[0],
        "uuid": tokens[1],
    })
    svc := service.NewClient(r.client)

    input := service.UpdateInput{
        ObjectId: tokens[1],
        Config: addr.Config{
            Name: state.Name.ValueString(),
            Fqdn: state.Fqdn.ValueString(),
            IpNetmask: state.IpNetmask.ValueString(),
            IpRange: state.IpRange.ValueString(),
            IpWildcard: state.IpWildcard.ValueString(),
            Description: state.Description.ValueString(),
            Tag: DecodeStringSlice(state.Tags),
        },
    }

    ans, err := svc.Update(ctx, input)
    if err != nil {
        resp.Diagnostics.AddError("Error updating config", err.Error())
        return
    }

    endstate := AddressObjectRsModel{
        Id: types.StringValue(id),
        Folder: types.StringValue(tokens[0]),
        ObjectId: types.StringValue(tokens[1]),
        Name: types.StringValue(ans.Name),
        Fqdn: types.StringValue(ans.Fqdn),
        IpNetmask: types.StringValue(ans.IpNetmask),
        IpRange: types.StringValue(ans.IpRange),
        IpWildcard: types.StringValue(ans.IpWildcard),
        Description: types.StringValue(ans.Description),
        Tags: EncodeStringSlice(ans.Tag),
        Type: types.StringValue(ans.Type),
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &endstate)...)
}

func (r *AddressObjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state AddressObjectRsModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    tokens := strings.Split(state.Id.ValueString(), IdSeparator)
    if len(tokens) != 2 {
        resp.Diagnostics.AddError("Error in resource ID", "Expected 2 tokens")
        return
    }

    tflog.Info(ctx, "deleting address object", map[string] any{
        "folder": tokens[0],
        "uuid": tokens[1],
    })
    svc := service.NewClient(r.client)

    input := service.DeleteInput{
        ObjectId: tokens[1],
    }

    if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
        resp.Diagnostics.AddError("Error deleting config", err.Error())
        return
    }
}

func (r *AddressObjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
