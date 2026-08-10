// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type VPSMonarxIdentity struct {
	VirtualMachineID types.Int64 `tfsdk:"virtual_machine_id"`
}
