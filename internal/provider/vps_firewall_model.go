// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
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

func (data *VPSFirewallModel) Merge(item *client.VPSV1FirewallFirewallResource) {
	data.ID = int64Value(item.Id)
	data.Name = types.StringPointerValue(item.Name)
	data.IsSynced = types.BoolPointerValue(item.IsSynced)
	data.CreatedAt = timetypes.NewRFC3339TimePointerValue(item.CreatedAt)
	data.UpdatedAt = timetypes.NewRFC3339TimePointerValue(item.UpdatedAt)

	if item.Rules != nil {
		for _, rule := range *item.Rules {
			var p VPSFirewallRuleModel
			p.ID = int64Value(rule.Id)
			p.Action = types.StringPointerValue((*string)(rule.Action))
			p.Protocol = types.StringPointerValue((*string)(rule.Protocol))
			p.Port = types.StringPointerValue(rule.Port)
			p.Source = types.StringPointerValue(rule.Source)
			p.SourceDetail = types.StringPointerValue(rule.SourceDetail)

			data.Rules = append(data.Rules, p)
		}
	}
}
