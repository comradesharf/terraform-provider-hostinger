// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type VPSFirewallIdentity struct {
	ID types.Int64 `tfsdk:"id"`
}

// VPSFirewallModel maps a single firewall from the API response.
type VPSFirewallModel struct {
	VPSFirewallIdentity
	Name      types.String      `tfsdk:"name"`
	IsSynced  types.Bool        `tfsdk:"is_synced"`
	CreatedAt timetypes.RFC3339 `tfsdk:"created_at"`
	UpdatedAt timetypes.RFC3339 `tfsdk:"updated_at"`
}

func (data *VPSFirewallModel) Merge(item client.VPSV1FirewallFirewallResource) {
	data.ID = int64Value(item.Id)
	data.Name = types.StringPointerValue(item.Name)
	data.IsSynced = types.BoolPointerValue(item.IsSynced)
	data.CreatedAt = timetypes.NewRFC3339TimePointerValue(item.CreatedAt)
	data.UpdatedAt = timetypes.NewRFC3339TimePointerValue(item.UpdatedAt)
}
