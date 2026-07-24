// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type VPSVirtualMachineIPAddressModel struct {
	ID      types.Int64       `tfsdk:"id"`
	Address iptypes.IPAddress `tfsdk:"address"`
	Ptr     types.String      `tfsdk:"ptr"`
}

type VPSVirtualMachineTemplateModel struct {
	ID            types.Int64  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Documentation types.String `tfsdk:"documentation"`
}

// VPSVirtualMachineModel describes the data source data model.
type VPSVirtualMachineModel struct {
	Id              types.Int64                       `tfsdk:"id"`
	FirewallGroupId types.Int64                       `tfsdk:"firewall_group_id"`
	SubscriptionID  types.String                      `tfsdk:"subscription_id"`
	DataCenterID    types.Int64                       `tfsdk:"data_center_id"`
	Plan            types.String                      `tfsdk:"plan"`
	Hostname        types.String                      `tfsdk:"hostname"`
	State           types.String                      `tfsdk:"state"`
	ActionsLock     types.String                      `tfsdk:"actions_lock"`
	Cpus            types.Int64                       `tfsdk:"cpus"`
	Memory          types.Int64                       `tfsdk:"memory"`
	Disk            types.Int64                       `tfsdk:"disk"`
	Bandwidth       types.Int64                       `tfsdk:"bandwidth"`
	NS1             iptypes.IPAddress                 `tfsdk:"ns1"`
	NS2             iptypes.IPAddress                 `tfsdk:"ns2"`
	Ipv4            []VPSVirtualMachineIPAddressModel `tfsdk:"ipv4"`
	Ipv6            []VPSVirtualMachineIPAddressModel `tfsdk:"ipv6"`
	Template        *VPSVirtualMachineTemplateModel   `tfsdk:"template"`
	CreatedAt       timetypes.RFC3339                 `tfsdk:"created_at"`
}
