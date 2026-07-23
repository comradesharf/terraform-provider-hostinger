package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

func int32Value(value *int) types.Int32 {
	if value == nil {
		return types.Int32Null()
	}

	return types.Int32Value(int32(*value))
}

func int64Value(value *int) types.Int64 {
	if value == nil {
		return types.Int64Null()
	}

	return types.Int64Value(int64(*value))
}
