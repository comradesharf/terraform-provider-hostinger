// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"slices"

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

func (data *VPSFirewallRuleModel) Merge(item *client.VPSV1FirewallFirewallRuleResource) {
	data.ID = int64Value(item.Id)
	data.Action = types.StringPointerValue((*string)(item.Action))
	data.Protocol = types.StringPointerValue((*string)(item.Protocol))
	data.Port = types.StringPointerValue(item.Port)
	data.Source = types.StringPointerValue(item.Source)
	data.SourceDetail = types.StringPointerValue(item.SourceDetail)
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

	if item.Rules == nil || len(*item.Rules) == 0 {
		return
	}

	slices.SortFunc(*item.Rules, func(a, b client.VPSV1FirewallFirewallRuleResource) int {
		return *a.Id - *b.Id
	})

	data.Rules = make([]VPSFirewallRuleModel, len(*item.Rules))

	for i, rule := range *item.Rules {
		var p VPSFirewallRuleModel
		p.Merge(&rule)
		data.Rules[i] = p
	}
}
