// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mock-server/mockserver-monorepo/mockserver-client-go/v7"
)

func TestAccVPSFirewallResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVPSFirewallResourcePreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVPSFirewallResourceConfig("one"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall.test",
						tfjsonpath.New("id"),
						knownvalue.Int64Exact(65224),
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
				},
			},
			{
				ResourceName:      "hostinger_vps_firewall.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccVPSFirewallResourceConfig("two"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall.test",
						tfjsonpath.New("id"),
						knownvalue.Int64Exact(65225),
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
				},
			},
		},
	})
}

func testAccVPSFirewallResourcePreCheck(t *testing.T) {
	client := newMockServerClient()

	if err := client.Clear(nil, mockserver.ClearAll); err != nil {
		t.Fatalf("failed to clear mock server: %v", err)
	}

	if err := client.FreezeClock("2021-09-01T12:00:00Z"); err != nil {
		t.Fatalf("failed to freeze mock server clock: %v", err)
	}

	// language=json
	body := []byte(`
[
	{
		"httpRequest": {
			"specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
			"operationId":      "VPS_createNewFirewallV1"
		},
		"httpResponse": {
			"statusCode": 200,
			"body": {
				"id": 65224,
				"name": "one",
				"is_synced": false,
				"rules": [],
				"created_at": "2021-09-01T12:00:00Z",
				"updated_at": "2021-09-01T12:00:00Z"
			}			
		},
		"times": {
			"remainingTimes": 1
		}
	},
	{
		"httpRequest": {
			"specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
			"operationId":      "VPS_getFirewallDetailsV1"
		},
		"httpResponse": {
			"statusCode": 200,
			"body": {
				"id": 65224,
				"name": "one",
				"is_synced": false,
				"rules": [],
				"created_at": "2021-09-01T12:00:00Z",
				"updated_at": "2021-09-01T12:00:00Z"
			}			
		},
		"times": {
			"remainingTimes": 3
		}
	},
	{
		"httpRequest": {
			"specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
			"operationId":      "VPS_deleteFirewallV1"
		},
		"httpResponse": {
			"statusCode": 200,
			"body": {
				"message": "Request accepted"
			}			
		},
		"times": {
			"remainingTimes": 1
		}
	},
	{
		"httpRequest": {
			"specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
			"operationId":      "VPS_createNewFirewallV1"
		},
		"httpResponse": {
			"statusCode": 200,
			"body": {
				"id": 65225,
				"name": "two",
				"is_synced": false,
				"rules": [],
				"created_at": "2021-09-01T12:00:00Z",
				"updated_at": "2021-09-01T12:00:00Z"
			}			
		},
		"times": {
			"remainingTimes": 1
		}
	},
	{
		"httpRequest": {
			"specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
			"operationId":      "VPS_getFirewallDetailsV1"
		},
		"httpResponse": {
			"statusCode": 200,
			"body": {
				"id": 65225,
				"name": "two",
				"is_synced": false,
				"rules": [],
				"created_at": "2021-09-01T12:00:00Z",
				"updated_at": "2021-09-01T12:00:00Z"
			}			
		},
		"times": {
			"remainingTimes": 1
		}
	},
	{
		"httpRequest": {
			"specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
			"operationId":      "VPS_deleteFirewallV1"
		},
		"httpResponse": {
			"statusCode": 200,
			"body": {
				"message": "Request accepted"
			}			
		},
		"times": {
			"remainingTimes": 1
		}
	}
]
`)

	req, _ := http.NewRequest(
		"PUT",
		"http://localhost:1234/mockserver/expectation",
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")
	if _, err := http.DefaultClient.Do(req); err != nil {
		t.Fatalf("failed to create expectations %v", err)
	}

}

func testAccVPSFirewallResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "hostinger_vps_firewall" "test" {
  name = %[1]q
}
`, name)
}
