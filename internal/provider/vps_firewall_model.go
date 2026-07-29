// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// VPSFirewallRuleModel maps a single firewall rule from the API response.
type VPSFirewallRuleModel struct {
	ID           types.Int64  `tfsdk:"id"`
	Action       types.String `tfsdk:"action"`
	Protocol     types.String `tfsdk:"protocol"`
	Port         types.String `tfsdk:"port"`
	Source       types.String `tfsdk:"source"`
	SourceDetail types.String `tfsdk:"source_detail"`
}

// VPSFirewallModel maps a single firewall from the API response.
type VPSFirewallModel struct {
	ID        types.Int64            `tfsdk:"id"`
	Name      types.String           `tfsdk:"name"`
	IsSynced  types.Bool             `tfsdk:"is_synced"`
	Rules     []VPSFirewallRuleModel `tfsdk:"rules"`
	CreatedAt timetypes.RFC3339      `tfsdk:"created_at"`
	UpdatedAt timetypes.RFC3339      `tfsdk:"updated_at"`
}
