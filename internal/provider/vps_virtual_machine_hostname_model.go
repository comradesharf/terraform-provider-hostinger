// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type VPSVirtualMachineHostnameModel struct {
	Hostname types.String `tfsdk:"hostname"`
}

func (d *VPSVirtualMachineHostnameModel) Merge(item client.VPSV1VirtualMachineVirtualMachineResource) {
	d.Hostname = types.StringPointerValue(item.Hostname)
}
