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
						knownvalue.ListSizeGreaterThan(0),
					),
				},
			},
		},
	})
}

func TestAccDataSourceVPSFirewallsWithID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVPSFirewallsWithIDConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_firewalls.test_with_id",
						tfjsonpath.New("firewalls"),
						knownvalue.ListSizeGreaterThan(0),
					),
				},
			},
		},
	})
}

const testAccDataSourceVPSFirewallsConfig = `
data "hostinger_vps_firewalls" "test" {}
`

const testAccDataSourceVPSFirewallsWithIDConfig = `
data "hostinger_vps_firewalls" "test_with_id" {
  firewall_id = 1
}
`
