// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type VPSTemplateModel struct {
	ID            types.Int64  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Documentation types.String `tfsdk:"documentation"`
}

func (m *VPSTemplateModel) Merge(item client.VPSV1TemplateTemplateResource) {
	m.ID = int64Value(item.Id)
	m.Name = types.StringPointerValue(item.Name)
	m.Description = types.StringPointerValue(item.Description)
	m.Documentation = types.StringPointerValue(item.Documentation)
}
