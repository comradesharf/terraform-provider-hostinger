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

func TestAccDataSourceVPSVirtualMachine(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVPSVirtualMachineConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_virtual_machine.test",
						tfjsonpath.New("id"),
						knownvalue.Int64Exact(17923),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_virtual_machine.test",
						tfjsonpath.New("firewall_group_id"),
						knownvalue.Int64Exact(260),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_virtual_machine.test",
						tfjsonpath.New("subscription_id"),
						knownvalue.StringExact("Azz353Uhl1xC54pR0"),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_virtual_machine.test",
						tfjsonpath.New("data_center_id"),
						knownvalue.Int64Exact(521),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_virtual_machine.test",
						tfjsonpath.New("plan"),
						knownvalue.StringExact("KVM 4"),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_virtual_machine.test",
						tfjsonpath.New("hostname"),
						knownvalue.StringExact("srv17923.hstgr.cloud"),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_virtual_machine.test",
						tfjsonpath.New("state"),
						knownvalue.StringExact("running"),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_virtual_machine.test",
						tfjsonpath.New("actions_lock"),
						knownvalue.StringExact("unlocked"),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_virtual_machine.test",
						tfjsonpath.New("cpus"),
						knownvalue.Int64Exact(4),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_virtual_machine.test",
						tfjsonpath.New("memory"),
						knownvalue.Int64Exact(8192),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_virtual_machine.test",
						tfjsonpath.New("disk"),
						knownvalue.Int64Exact(51200),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_virtual_machine.test",
						tfjsonpath.New("bandwidth"),
						knownvalue.Int64Exact(1073741824),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_virtual_machine.test",
						tfjsonpath.New("ns1"),
						knownvalue.StringExact("1.1.1.1"),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_virtual_machine.test",
						tfjsonpath.New("ns2"),
						knownvalue.StringExact("8.8.8.8"),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_virtual_machine.test",
						tfjsonpath.New("ipv4"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id":      knownvalue.Int64Exact(52347),
								"address": knownvalue.StringExact("213.211.223.15"),
								"ptr":     knownvalue.StringExact("something.domain.tld"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_virtual_machine.test",
						tfjsonpath.New("ipv4"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id":      knownvalue.Int64Exact(52347),
								"address": knownvalue.StringExact("213.211.223.15"),
								"ptr":     knownvalue.StringExact("something.domain.tld"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_virtual_machine.test",
						tfjsonpath.New("template"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":            knownvalue.Int64Exact(6523),
							"name":          knownvalue.StringExact("Ubuntu 20.04 LTS"),
							"description":   knownvalue.StringExact("Ubuntu 20.04 LTS"),
							"documentation": knownvalue.StringExact("https://docs.ubuntu.com"),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_virtual_machine.test",
						tfjsonpath.New("created_at"),
						knownvalue.StringExact("2024-09-05T07:25:36Z"),
					),
				},
			},
		},
	})
}

const testAccDataSourceVPSVirtualMachineConfig = `
data "hostinger_vps_virtual_machine" "test" {
	id = 17923
}
`
