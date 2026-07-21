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

func TestAccDataSourceReachContacts(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReachContactsConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.hostinger_reach_contacts.test",
						tfjsonpath.New("contacts").AtSliceIndex(0),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"uuid":                knownvalue.NotNull(),
							"email":               knownvalue.NotNull(),
							"name":                knownvalue.NotNull(),
							"surname":             knownvalue.NotNull(),
							"subscription_status": knownvalue.NotNull(),
							"subscribed_at":       knownvalue.NotNull(),
						}),
					),
				},
			},
		},
	})
}

const testAccReachContactsConfig = `
data "hostinger_reach_contacts" "test" {
}
`
