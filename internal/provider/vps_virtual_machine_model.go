// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type VPSVirtualMachineIdentity struct {
	ID types.Int64 `tfsdk:"id"`
}

type VPSVirtualMachineIPAddressModel struct {
	ID      types.Int64       `tfsdk:"id"`
	Address iptypes.IPAddress `tfsdk:"address"`
	Ptr     types.String      `tfsdk:"ptr"`
}

func (d *VPSVirtualMachineIPAddressModel) Merge(item client.VPSV1IPAddressIPAddressResource) {
	d.ID = int64Value(item.Id)
	d.Address = iptypes.NewIPAddressPointerValue(item.Address)
	d.Ptr = types.StringPointerValue(item.Ptr)
}

type VPSVirtualMachineTemplateModel struct {
	ID            types.Int64  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Documentation types.String `tfsdk:"documentation"`
}

func (d *VPSVirtualMachineTemplateModel) Merge(item client.VPSV1TemplateTemplateResource) {
	d.ID = int64Value(item.Id)
	d.Name = types.StringPointerValue(item.Name)
	d.Description = types.StringPointerValue(item.Description)
	d.Documentation = types.StringPointerValue(item.Documentation)
}

// VPSVirtualMachineModel describes the data source data model.
type VPSVirtualMachineModel struct {
	ID              types.Int64                       `tfsdk:"id"`
	FirewallGroupID types.Int64                       `tfsdk:"firewall_group_id"`
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

func (d *VPSVirtualMachineModel) Merge(item client.VPSV1VirtualMachineVirtualMachineResource) {
	d.ID = int64Value(item.Id)
	d.FirewallGroupID = int64Value(item.FirewallGroupId)
	d.SubscriptionID = types.StringPointerValue(item.SubscriptionId)
	d.DataCenterID = int64Value(item.DataCenterId)
	d.Plan = types.StringPointerValue(item.Plan)
	d.Hostname = types.StringPointerValue(item.Hostname)
	d.State = types.StringPointerValue((*string)(item.State))
	d.ActionsLock = types.StringPointerValue((*string)(item.ActionsLock))
	d.Cpus = int64Value(item.Cpus)
	d.Memory = int64Value(item.Memory)
	d.Disk = int64Value(item.Disk)
	d.Bandwidth = int64Value(item.Bandwidth)
	d.NS1 = iptypes.NewIPAddressPointerValue(item.Ns1)
	d.NS2 = iptypes.NewIPAddressPointerValue(item.Ns2)

	if item.Ipv4 != nil {
		v, err := item.Ipv4.AsVPSV1IPAddressIPAddressCollection()
		if err == nil {
			for _, ip := range v {
				var p VPSVirtualMachineIPAddressModel
				p.Merge(ip)
				d.Ipv4 = append(d.Ipv4, p)
			}
		}
	}

	if item.Ipv6 != nil {
		v, err := item.Ipv6.AsVPSV1IPAddressIPAddressCollection()
		if err == nil {
			for _, ip := range v {
				var p VPSVirtualMachineIPAddressModel
				p.Merge(ip)
				d.Ipv6 = append(d.Ipv6, p)
			}
		}
	}

	if item.Template != nil {
		v, err := item.Template.AsVPSV1TemplateTemplateResource()
		if err == nil {
			var p VPSVirtualMachineTemplateModel
			p.Merge(v)
			d.Template = &p
		}
	}

	d.CreatedAt = timetypes.NewRFC3339TimePointerValue(item.CreatedAt)
}
