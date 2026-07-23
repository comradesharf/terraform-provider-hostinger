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
				Config: testAccDataSourceReachContactsConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.hostinger_reach_contacts.test",
						tfjsonpath.New("contacts").AtSliceIndex(0),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"uuid":                knownvalue.StringExact("550e8400-e29b-41d4-a716-446655440000"),
							"email":               knownvalue.StringExact("john.doe@example.com"),
							"name":                knownvalue.StringExact("John"),
							"surname":             knownvalue.StringExact("Doe"),
							"subscription_status": knownvalue.StringExact("subscribed"),
							"subscribed_at":       knownvalue.StringExact("2023-01-01T00:00:00Z"),
							"source":              knownvalue.StringExact("sync"),
							"note":                knownvalue.StringExact("VIP customer"),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_reach_contacts.test",
						tfjsonpath.New("group_uuid"),
						knownvalue.StringExact("550e8400-e29b-41d4-a716-446655440000"),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_reach_contacts.test",
						tfjsonpath.New("subscription_status"),
						knownvalue.StringExact("subscribed"),
					),
				},
			},
		},
	})
}

const testAccDataSourceReachContactsConfig = `
data "hostinger_reach_contacts" "test" {
	group_uuid = "550e8400-e29b-41d4-a716-446655440000"
	subscription_status = "subscribed"
}
`
