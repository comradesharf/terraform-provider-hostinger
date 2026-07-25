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

func TestAccDataSourceVPSFirewalls(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVPSFirewallsConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_firewalls.test",
						tfjsonpath.New("firewalls"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id":        knownvalue.Int64Exact(65224),
								"name":      knownvalue.StringExact("HTTP and SSH only"),
								"is_synced": knownvalue.Bool(false),
								"rules": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"id":            knownvalue.Int64Exact(24541),
										"action":        knownvalue.StringExact("accept"),
										"protocol":      knownvalue.StringExact("TCP"),
										"port":          knownvalue.StringExact("1024:2048"),
										"source":        knownvalue.StringExact("any"),
										"source_detail": knownvalue.StringExact("any"),
									}),
								}),
								"created_at": knownvalue.StringExact("2021-09-01T12:00:00Z"),
								"updated_at": knownvalue.StringExact("2021-09-01T12:00:00Z"),
							}),
						}),
					),
				},
			},
		},
	})
}

const testAccDataSourceVPSFirewallsConfig = `
data "hostinger_vps_firewalls" "test" {}
`
