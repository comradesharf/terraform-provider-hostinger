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

func TestAccDataSourceVPSVirtualMachines(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVPSVirtualMachinesConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_virtual_machines.test",
						tfjsonpath.New("virtual_machines"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id":                knownvalue.Int64Exact(17923),
								"firewall_group_id": knownvalue.Int64Exact(260),
								"subscription_id":   knownvalue.StringExact("Azz353Uhl1xC54pR0"),
								"data_center_id":    knownvalue.Int64Exact(521),
								"plan":              knownvalue.StringExact("KVM 4"),
								"hostname":          knownvalue.StringExact("srv17923.hstgr.cloud"),
								"state":             knownvalue.StringExact("running"),
								"actions_lock":      knownvalue.StringExact("unlocked"),
								"cpus":              knownvalue.Int64Exact(4),
								"memory":            knownvalue.Int64Exact(8192),
								"disk":              knownvalue.Int64Exact(51200),
								"bandwidth":         knownvalue.Int64Exact(1073741824),
								"ns1":               knownvalue.StringExact("1.1.1.1"),
								"ns2":               knownvalue.StringExact("8.8.8.8"),
								"ipv4": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"id":      knownvalue.Int64Exact(52347),
										"address": knownvalue.StringExact("213.211.223.15"),
										"ptr":     knownvalue.StringExact("something.domain.tld"),
									}),
								}),
								"ipv6": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"id":      knownvalue.Int64Exact(52347),
										"address": knownvalue.StringExact("213.211.223.15"),
										"ptr":     knownvalue.StringExact("something.domain.tld"),
									}),
								}),
								"template": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"id":            knownvalue.Int64Exact(6523),
									"name":          knownvalue.StringExact("Ubuntu 20.04 LTS"),
									"description":   knownvalue.StringExact("Ubuntu 20.04 LTS"),
									"documentation": knownvalue.StringExact("https://docs.ubuntu.com"),
								}),
								"created_at": knownvalue.StringExact("2024-09-05T07:25:36Z"),
							}),
						}),
					),
				},
			},
		},
	})
}

const testAccDataSourceVPSVirtualMachinesConfig = `
data "hostinger_vps_virtual_machines" "test" {}
`
