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
						knownvalue.ListPartial(map[int]knownvalue.Check{
							0: knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"id":                knownvalue.Int64Exact(17923),
								"hostname":          knownvalue.StringExact("srv17923.hstgr.cloud"),
								"state":             knownvalue.StringExact("running"),
								"cpus":              knownvalue.Int64Exact(4),
								"memory":            knownvalue.Int64Exact(8192),
								"disk":              knownvalue.Int64Exact(51200),
								"bandwidth":         knownvalue.Int64Exact(1073741824),
								"data_center_id":    knownvalue.Int64Exact(521),
								"firewall_group_id": knownvalue.Int64Exact(260),
								"os_name":           knownvalue.StringExact("Ubuntu 20.04 LTS"),
								"created_at":        knownvalue.StringExact("2024-09-05T07:25:36Z"),
								"actions_lock":      knownvalue.StringExact("unlocked"),
								"ipv4": knownvalue.ListPartial(map[int]knownvalue.Check{
									0: knownvalue.ObjectPartial(map[string]knownvalue.Check{
										"address": knownvalue.StringExact("213.331.273.15"),
										"ptr":     knownvalue.StringExact("something.domain.tld"),
										"netmask": knownvalue.Null(),
									}),
								}),
								"ipv6": knownvalue.ListPartial(map[int]knownvalue.Check{
									0: knownvalue.ObjectPartial(map[string]knownvalue.Check{
										"address": knownvalue.StringExact("213.331.273.15"),
										"ptr":     knownvalue.StringExact("something.domain.tld"),
									}),
								}),
							}),
						}),
					),
				},
			},
			{
				Config: testAccDataSourceVPSVirtualMachineConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_virtual_machines.test",
						tfjsonpath.New("virtual_machine_id"),
						knownvalue.Int64Exact(17923),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_virtual_machines.test",
						tfjsonpath.New("virtual_machines").AtSliceIndex(0),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"id":       knownvalue.Int64Exact(17923),
							"hostname": knownvalue.StringExact("srv17923.hstgr.cloud"),
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

const testAccDataSourceVPSVirtualMachineConfig = `
data "hostinger_vps_virtual_machines" "test" {
	virtual_machine_id = 17923
}
`
