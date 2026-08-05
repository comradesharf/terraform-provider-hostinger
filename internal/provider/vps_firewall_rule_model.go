// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
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

func (data *VPSFirewallRuleModel) Merge(item *client.VPSV1FirewallFirewallRuleResource) {
	data.ID = int64Value(item.Id)
	data.Action = types.StringPointerValue((*string)(item.Action))
	data.Protocol = types.StringPointerValue((*string)(item.Protocol))
	data.Port = types.StringPointerValue(item.Port)
	data.Source = types.StringPointerValue(item.Source)
	data.SourceDetail = types.StringPointerValue(item.SourceDetail)
}
