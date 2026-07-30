// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mock-server/mockserver-monorepo/mockserver-client-go/v7"
)

func TestAccResourceVPSFirewall(t *testing.T) {
	mockClient := newMockServerClient()

	compareValuesDiffer := statecheck.CompareValue(compare.ValuesDiffer())

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

			if err := mockClient.Clear(nil, mockserver.ClearLog); err != nil {
				t.Fatalf("failed to clear mock server logs: %v", err)
			}

			if err := mockClient.FreezeClock("2021-09-01T12:00:00Z"); err != nil {
				t.Fatalf("failed to freeze mock server clock: %v", err)
			}

			if err := mockClient.Scenario("CRUDFirewall").Set("Started"); err != nil {
				t.Fatalf("failed to set mock server scenario: %v", err)
			}
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVPSFirewallConfig("one"),
				ConfigStateChecks: []statecheck.StateCheck{
					compareValuesDiffer.AddStateValue(
						"hostinger_vps_firewall.test",
						tfjsonpath.New("id"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("one"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall.test",
						tfjsonpath.New("is_synced"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall.test",
						tfjsonpath.New("created_at"),
						knownvalue.StringExact("2021-09-01T12:00:00Z"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall.test",
						tfjsonpath.New("updated_at"),
						knownvalue.StringExact("2021-09-01T12:00:00Z"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall.test",
						tfjsonpath.New("rules"),
						knownvalue.Null(),
					),
				},
			},
			{
				ResourceName:      "hostinger_vps_firewall.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceVPSFirewallConfig("two"),
				ConfigStateChecks: []statecheck.StateCheck{
					compareValuesDiffer.AddStateValue(
						"hostinger_vps_firewall.test",
						tfjsonpath.New("id"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("two"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall.test",
						tfjsonpath.New("is_synced"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall.test",
						tfjsonpath.New("created_at"),
						knownvalue.StringExact("2021-09-01T12:00:00Z"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall.test",
						tfjsonpath.New("updated_at"),
						knownvalue.StringExact("2021-09-01T12:00:00Z"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall.test",
						tfjsonpath.New("rules"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

func testAccResourceVPSFirewallConfig(name string) string {
	return fmt.Sprintf(`
resource "hostinger_vps_firewall" "test" {
  name = %[1]q
}
`, name)
}
