// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ReachContactsContactModel maps a single contact from the API response.
type ReachContactsContactModel struct {
	Uuid               types.String      `tfsdk:"uuid"`
	Email              types.String      `tfsdk:"email"`
	Name               types.String      `tfsdk:"name"`
	Surname            types.String      `tfsdk:"surname"`
	SubscriptionStatus types.String      `tfsdk:"subscription_status"`
	SubscribedAt       timetypes.RFC3339 `tfsdk:"subscribed_at"`
	Source             types.String      `tfsdk:"source"`
	Note               types.String      `tfsdk:"note"`
}

func (d *ReachContactsContactModel) Merge(item client.ReachV1ContactsContactResource) {
	d.Uuid = types.StringPointerValue(item.Uuid)
	d.Email = types.StringPointerValue(item.Email)
	d.Name = types.StringPointerValue(item.Name)
	d.Surname = types.StringPointerValue(item.Surname)
	d.SubscriptionStatus = types.StringPointerValue((*string)(item.SubscriptionStatus))
	d.Source = types.StringPointerValue((*string)(item.Source))
	d.Note = types.StringPointerValue(item.Note)
	d.SubscribedAt = timetypes.NewRFC3339TimePointerValue(item.SubscribedAt)
}
