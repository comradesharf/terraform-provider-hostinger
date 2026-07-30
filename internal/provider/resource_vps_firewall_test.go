// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccResourceVPSFirewall(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccResourceVPSFirewallConfig("one"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("one"),
					),
				},
			},
			//// ImportState testing
			//{
			//	ResourceName:      "hostinger_vps_firewall.test",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//	// This is not normally necessary, but is here because this
			//	// example code does not have an actual upstream service.
			//	// Once the Read method is able to refresh information from
			//	// the upstream service, this can be removed.
			//	ImportStateVerifyIgnore: []string{"configurable_attribute", "defaulted"},
			//},
			//// Update and Read testing
			//{
			//	Config: testAccResourceVPSFirewallConfig("two"),
			//	ConfigStateChecks: []statecheck.StateCheck{
			//		statecheck.ExpectKnownValue(
			//			"hostinger_vps_firewall.test",
			//			tfjsonpath.New("id"),
			//			knownvalue.StringExact("example-id"),
			//		),
			//		statecheck.ExpectKnownValue(
			//			"hostinger_vps_firewall.test",
			//			tfjsonpath.New("defaulted"),
			//			knownvalue.StringExact("example value when not configured"),
			//		),
			//		statecheck.ExpectKnownValue(
			//			"hostinger_vps_firewall.test",
			//			tfjsonpath.New("configurable_attribute"),
			//			knownvalue.StringExact("two"),
			//		),
			//	},
			//},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccResourceVPSFirewallConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "hostinger_vps_firewall" "test" {
  name = %[1]q
}
`, configurableAttribute)
}
