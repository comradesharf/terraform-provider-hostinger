// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// VPSPublicKeysPublicKeyModel maps a single public key from the API response.
type VPSPublicKeysPublicKeyModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Key  types.String `tfsdk:"key"`
}

func (d *VPSPublicKeysPublicKeyModel) Merge(item client.VPSV1PublicKeyPublicKeyResource) {
	d.ID = int64Value(item.Id)
	d.Name = types.StringPointerValue(item.Name)
	d.Key = types.StringPointerValue(item.Key)
}
