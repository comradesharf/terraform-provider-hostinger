// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type VPSDataCenterModel struct {
	ID        types.Int64  `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	City      types.String `tfsdk:"city"`
	Continent types.String `tfsdk:"continent"`
	Location  types.String `tfsdk:"location"`
}

func (d *VPSDataCenterModel) Merge(item client.VPSV1DataCenterDataCenterResource) {
	d.ID = int64Value(item.Id)
	d.Name = types.StringPointerValue(item.Name)
	d.City = types.StringPointerValue(item.City)
	d.Continent = types.StringPointerValue(item.Continent)
	d.Location = types.StringPointerValue(item.Location)
}
