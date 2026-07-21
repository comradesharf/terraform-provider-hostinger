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
				Config: testAccReachSegmentsConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.hostinger_reach_segments.test",
						tfjsonpath.New("segments").AtSliceIndex(0),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"uuid":       knownvalue.NotNull(),
							"name":       knownvalue.NotNull(),
							"created_at": knownvalue.NotNull(),
							"updated_at": knownvalue.NotNull(),
						}),
					),
				},
			},
		},
	})
}

const testAccReachSegmentsConfig = `
data "hostinger_reach_segments" "test" {
}
`
