// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AgencyHostingDomainModel maps a single domain from the API response.
type AgencyHostingDomainModel struct {
	FQDN       types.String      `tfsdk:"fqdn"`
	WebsiteUID types.String      `tfsdk:"website_uid"`
	CreatedAt  timetypes.RFC3339 `tfsdk:"created_at"`
}

func (m *AgencyHostingDomainModel) Merge(item client.AgencyHostingV1DomainsDomainResource) {
	m.FQDN = types.StringPointerValue(item.Fqdn)
	m.WebsiteUID = types.StringPointerValue(item.WebsiteUid)
	m.CreatedAt = timetypes.NewRFC3339TimePointerValue(item.CreatedAt)
}
