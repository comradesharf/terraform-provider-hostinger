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

func TestAccDataSourceReachSegments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceReachSegmentsConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.hostinger_reach_segments.test",
						tfjsonpath.New("segments").AtSliceIndex(0),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"uuid":       knownvalue.StringExact("550e8400-e29b-41d4-a716-446655440000"),
							"name":       knownvalue.StringExact("Newsletter Subscribers"),
							"created_at": knownvalue.StringExact("2025-02-27T11:54:22Z"),
							"updated_at": knownvalue.StringExact("2025-02-27T11:54:22Z"),
						}),
					),
				},
			},
		},
	})
}

const testAccDataSourceReachSegmentsConfig = `
data "hostinger_reach_segments" "test" {
}
`
