// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

func importStatePassthroughInt64ID(ctx context.Context, attrPath path.Path, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if attrPath.Equal(path.Empty()) {
		resp.Diagnostics.AddError(
			"Resource Import Passthrough Missing Attribute Path",
			"This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				"Resource ImportState method call to ImportStatePassthroughInt64ID path must be set to a valid attribute path that can accept a string value.",
		)
		return
	}
	if req.ID == "" {
		return
	}
	v, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Resource Import Passthrough Invalid ID",
			"Failed to parse import ID as int64: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, attrPath, v)...)
}

func importStatePassthroughInt64Identity(ctx context.Context, stateAttrPath path.Path, identityAttrPath path.Path, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if stateAttrPath.Equal(path.Empty()) {
		resp.Diagnostics.AddError(
			"Resource Import Passthrough Missing State Attribute Path",
			"This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				"Resource ImportState method call to ImportStatePassthroughWithIdentity path must be set to a valid state attribute path that can accept a int64 value.",
		)
	}

	if identityAttrPath.Equal(path.Empty()) {
		resp.Diagnostics.AddError(
			"Resource Import Passthrough Missing Identity Attribute Path",
			"This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				"Resource ImportState method call to ImportStatePassthroughWithIdentity path must be set to a valid identity attribute path that is a int64 value.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// If the import is using the import identifier, (either via the "terraform import" CLI command, or a config block with the "id" attribute set)
	// pass through the ID to the designated state attribute.
	if req.ID != "" {
		v, err := strconv.ParseInt(req.ID, 10, 64)
		if err != nil {
			resp.Diagnostics.AddError(
				"Resource Import Passthrough Invalid ID",
				"Failed to parse import ID as int64: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, stateAttrPath, v)...)
		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, identityAttrPath, v)...)
		return
	}

	// The import isn't using the import identifier, so it must be using identity. Grab the designated
	// identity attribute string and set it to state.
	var identityAttrVal types.Int64
	resp.Diagnostics.Append(req.Identity.GetAttribute(ctx, identityAttrPath, &identityAttrVal)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, stateAttrPath, identityAttrVal)...)
}
