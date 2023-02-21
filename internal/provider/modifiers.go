package provider

import (
	"context"
	"fmt"
    "strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func DefaultString(v string) planmodifier.String {
    return &defaultValue{String: &v}
}

func DefaultBool(v bool) planmodifier.Bool {
    return &defaultValue{Bool: &v}
}

func DefaultInt64(v int64) planmodifier.Int64 {
    return &defaultValue{Int64: &v}
}

func DefaultFloat64(v float64) planmodifier.Float64 {
    return &defaultValue{Float64: &v}
}

func DefaultStringSlice(v []string) planmodifier.List {
    return &defaultValue{StringSlice: v}
}

type defaultValue struct {
    String *string
    Bool *bool
    Int64 *int64
    Float64 *float64
    StringSlice []string
}

func (o *defaultValue) Description(_ context.Context) string {
    var b strings.Builder
    prefix := "Default value: `"

    switch {
    case o.String != nil && *o.String != "":
        b.WriteString(prefix)
        b.WriteString(fmt.Sprintf("%q", *o.String))
    case o.Bool != nil && *o.Bool:
        b.WriteString(prefix)
        b.WriteString(fmt.Sprintf("%t", *o.Bool))
    case o.Int64 != nil && *o.Int64 != 0:
        b.WriteString(prefix)
        b.WriteString(fmt.Sprintf("%d", *o.Int64))
    case o.Float64 != nil && *o.Float64 != 0:
        b.WriteString(prefix)
        b.WriteString(fmt.Sprintf("%f", *o.Float64))
    case len(o.StringSlice) > 0:
        b.WriteString(prefix)
        b.WriteString("[")
        for num, x := range o.StringSlice {
            if num != 0 {
                b.WriteString(", ")
            }
            b.WriteString(fmt.Sprintf("%q", x))
        }
        b.WriteString("]")
    }

    if b.Len() > 0 {
        b.WriteString("`.")
    }

    return b.String()
}

func (o *defaultValue) MarkdownDescription(ctx context.Context) string {
    return o.Description(ctx)
}

func (o *defaultValue) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
    if !req.ConfigValue.IsNull() {
        return
    }

    // If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
    // has already been applied, and we don't want to interfere.
    if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
        return
    }

    resp.PlanValue = types.StringValue(*o.String)
}

func (o *defaultValue) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
    if !req.ConfigValue.IsNull() {
        return
    }

    // If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
    // has already been applied, and we don't want to interfere.
    if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
        return
    }

    resp.PlanValue = types.BoolValue(*o.Bool)
}

func (o *defaultValue) PlanModifyInt64(_ context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
    if !req.ConfigValue.IsNull() {
        return
    }

    // If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
    // has already been applied, and we don't want to interfere.
    if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
        return
    }

    resp.PlanValue = types.Int64Value(*o.Int64)
}

func (o *defaultValue) PlanModifyFloat64(_ context.Context, req planmodifier.Float64Request, resp *planmodifier.Float64Response) {
    if !req.ConfigValue.IsNull() {
        return
    }

    // If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
    // has already been applied, and we don't want to interfere.
    if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
        return
    }

    resp.PlanValue = types.Float64Value(*o.Float64)
}

func (o *defaultValue) PlanModifyList(_ context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
    if !req.ConfigValue.IsNull() {
        return
    }

    // If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
    // has already been applied, and we don't want to interfere.
    if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
        return
    }

    if len(o.StringSlice) == 0 {
        resp.PlanValue = basetypes.NewListNull(types.StringType)
    } else {
        list := make([]attr.Value, 0, len(o.StringSlice))
        for _, x := range o.StringSlice {
            list = append(list, types.StringValue(x))
        }
        resp.PlanValue = basetypes.NewListValueMust(types.StringType, list)
    }
}
