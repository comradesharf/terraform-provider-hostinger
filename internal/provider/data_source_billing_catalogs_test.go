// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccDataSourceBillingCatalogs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBillingCatalogsConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.hostinger_billing_catalogs.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("KVM 2"),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_billing_catalogs.test",
						tfjsonpath.New("category"),
						knownvalue.StringExact("VPS"),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_billing_catalogs.test",
						tfjsonpath.New("billing_catalogs").AtSliceIndex(0),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":       knownvalue.StringExact("hostingercom-vps-kvm2"),
							"name":     knownvalue.StringExact("KVM 2"),
							"category": knownvalue.StringExact("VPS"),
							"metadata": knownvalue.MapExact(map[string]knownvalue.Check{
								"field": knownvalue.StringExact("value"),
							}),
							"prices": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"id":                 knownvalue.StringExact("hostingercom-vps-kvm2-usd-1m"),
									"name":               knownvalue.StringExact("KVM 2 (billed every month)"),
									"currency":           knownvalue.StringExact("USD"),
									"price":              knownvalue.Int32Exact(1799),
									"first_period_price": knownvalue.Int32Exact(899),
									"period":             knownvalue.Int32Exact(1),
									"period_unit":        knownvalue.StringExact("day"),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}

const testAccDataSourceBillingCatalogsConfig = `
data "hostinger_billing_catalogs" "test" {
	name = "KVM 2"
	category = "VPS"
}
`
