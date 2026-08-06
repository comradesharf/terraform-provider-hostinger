// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BillingCatalogsCatalogPriceModel struct {
	Currency         types.String `tfsdk:"currency"`
	FirstPeriodPrice types.Int32  `tfsdk:"first_period_price"`
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Period           types.Int32  `tfsdk:"period"`
	PeriodUnit       types.String `tfsdk:"period_unit"`
	Price            types.Int32  `tfsdk:"price"`
}

func (d *BillingCatalogsCatalogPriceModel) Merge(item client.BillingV1CatalogCatalogItemPriceResource) {
	d.ID = types.StringPointerValue(item.Id)
	d.Currency = types.StringPointerValue(item.Currency)
	d.FirstPeriodPrice = int32Value(item.FirstPeriodPrice)
	d.Name = types.StringPointerValue(item.Name)
	d.Period = int32Value(item.Period)
	d.PeriodUnit = types.StringPointerValue((*string)(item.PeriodUnit))
	d.Price = int32Value(item.Price)
}

type BillingCatalogsCatalogModel struct {
	ID       types.String                       `tfsdk:"id"`
	Category types.String                       `tfsdk:"category"`
	Name     types.String                       `tfsdk:"name"`
	Metadata map[string]types.String            `tfsdk:"metadata"`
	Prices   []BillingCatalogsCatalogPriceModel `tfsdk:"prices"`
}

func (d *BillingCatalogsCatalogModel) Merge(item client.BillingV1CatalogCatalogItemResource) {
	d.ID = types.StringPointerValue(item.Id)
	d.Category = types.StringPointerValue(item.Category)
	d.Name = types.StringPointerValue(item.Name)

	if item.Metadata != nil {
		d.Metadata = make(map[string]types.String, len(*item.Metadata))
		for k, v := range *item.Metadata {
			if s, ok := v.(string); ok {
				d.Metadata[k] = types.StringValue(s)
			}
		}
	}

	if item.Prices != nil {
		for _, price := range *item.Prices {
			var p BillingCatalogsCatalogPriceModel
			p.Merge(price)
			d.Prices = append(d.Prices, p)
		}
	}
}
