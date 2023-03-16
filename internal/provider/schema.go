package provider

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Int64ToString(v int64) string {
	return fmt.Sprintf("%d", v)
}

func Float64ToString(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func StringToInt64(v string) (int64, error) {
	return strconv.ParseInt(v, 10, 64)
}

func StringToFloat64(v string) (float64, error) {
	return strconv.ParseFloat(v, 64)
}

func DecodeStringSlice(v []types.String) []string {
	if len(v) == 0 {
		return nil
	}

	ans := make([]string, 0, len(v))
	for _, x := range v {
		ans = append(ans, x.ValueString())
	}

	return ans
}

func EncodeStringSlice(v []string) []types.String {
	if len(v) == 0 {
		return nil
	}

	ans := make([]types.String, 0, len(v))
	for _, x := range v {
		ans = append(ans, types.StringValue(x))
	}

	return ans
}

func DecodeInt64Slice(v []types.Int64) []int64 {
	if len(v) == 0 {
		return nil
	}

	ans := make([]int64, 0, len(v))
	for _, x := range v {
		ans = append(ans, x.ValueInt64())
	}

	return ans
}

func EncodeInt64Slice(v []int64) []types.Int64 {
	if len(v) == 0 {
		return nil
	}

	ans := make([]types.Int64, 0, len(v))
	for _, x := range v {
		ans = append(ans, types.Int64Value(x))
	}

	return ans
}
