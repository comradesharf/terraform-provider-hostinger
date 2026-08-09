// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
)

type VPSVirtualMachineNameserversModel struct {
	NS1 iptypes.IPAddress `tfsdk:"ns1"`
	NS2 iptypes.IPAddress `tfsdk:"ns2"`
}

func (d *VPSVirtualMachineNameserversModel) Merge(item client.VPSV1VirtualMachineVirtualMachineResource) {
	d.NS1 = iptypes.NewIPAddressPointerValue(item.Ns1)
	d.NS2 = iptypes.NewIPAddressPointerValue(item.Ns2)
}
