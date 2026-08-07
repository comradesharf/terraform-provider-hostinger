// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// VPSPostInstallScriptModel maps a single VPS post-install script from the API response.
type VPSPostInstallScriptModel struct {
	ID        types.Int64       `tfsdk:"id"`
	Name      types.String      `tfsdk:"name"`
	Content   types.String      `tfsdk:"content"`
	CreatedAt timetypes.RFC3339 `tfsdk:"created_at"`
	UpdatedAt timetypes.RFC3339 `tfsdk:"updated_at"`
}

func (m *VPSPostInstallScriptModel) Merge(item client.VPSV1PostInstallScriptPostInstallScriptResource) {
	m.ID = int64Value(item.Id)
	m.Name = types.StringPointerValue(item.Name)
	m.Content = types.StringPointerValue(item.Content)
	m.CreatedAt = timetypes.NewRFC3339TimePointerValue(item.CreatedAt)
	m.UpdatedAt = timetypes.NewRFC3339TimePointerValue(item.UpdatedAt)
}
