// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ReachSegmentModel maps a single segment from the API response.
type ReachSegmentModel struct {
	Uuid      types.String      `tfsdk:"uuid"`
	Name      types.String      `tfsdk:"name"`
	CreatedAt timetypes.RFC3339 `tfsdk:"created_at"`
	UpdatedAt timetypes.RFC3339 `tfsdk:"updated_at"`
}

func (d *ReachSegmentModel) Merge(item client.ReachV1ContactsSegmentsContactSegmentResource) {
	d.Uuid = types.StringPointerValue(item.Uuid)
	d.Name = types.StringPointerValue(item.Name)
	d.CreatedAt = timetypes.NewRFC3339TimePointerValue(item.CreatedAt)
	d.UpdatedAt = timetypes.NewRFC3339TimePointerValue(item.UpdatedAt)
}
