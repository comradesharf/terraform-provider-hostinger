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

func TestAccVPSFirewallRuleResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVPSFirewallRuleResourcePreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVPSFirewallRuleResourceConfig("test1", "TCP"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall_rule.test",
						tfjsonpath.New("id"),
						knownvalue.Int64Exact(24541),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall_rule.test",
						tfjsonpath.New("firewall_id"),
						knownvalue.Int64Exact(65224),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall_rule.test",
						tfjsonpath.New("protocol"),
						knownvalue.StringExact("TCP"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall_rule.test",
						tfjsonpath.New("action"),
						knownvalue.StringExact("accept"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall_rule.test",
						tfjsonpath.New("port"),
						knownvalue.StringExact("80"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall_rule.test",
						tfjsonpath.New("source"),
						knownvalue.StringExact("custom"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall_rule.test",
						tfjsonpath.New("source_detail"),
						knownvalue.StringExact("any"),
					),
				},
			},
			{
				ResourceName:      "hostinger_vps_firewall_rule.test",
				ImportState:       true,
				ImportStateId:     "65224/24541",
				ImportStateVerify: true,
			},
			{
				Config: testAccVPSFirewallRuleResourceConfig("test1", "UDP"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall_rule.test",
						tfjsonpath.New("id"),
						knownvalue.Int64Exact(24541),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall_rule.test",
						tfjsonpath.New("firewall_id"),
						knownvalue.Int64Exact(65224),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall_rule.test",
						tfjsonpath.New("protocol"),
						knownvalue.StringExact("UDP"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall_rule.test",
						tfjsonpath.New("action"),
						knownvalue.StringExact("accept"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall_rule.test",
						tfjsonpath.New("port"),
						knownvalue.StringExact("80"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall_rule.test",
						tfjsonpath.New("source"),
						knownvalue.StringExact("custom"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall_rule.test",
						tfjsonpath.New("source_detail"),
						knownvalue.StringExact("any"),
					),
				},
			},
			{
				Config: testAccVPSFirewallRuleResourceConfig("test2", "UDP"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall_rule.test",
						tfjsonpath.New("id"),
						knownvalue.Int64Exact(24542),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall_rule.test",
						tfjsonpath.New("firewall_id"),
						knownvalue.Int64Exact(65225),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall_rule.test",
						tfjsonpath.New("protocol"),
						knownvalue.StringExact("UDP"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall_rule.test",
						tfjsonpath.New("action"),
						knownvalue.StringExact("accept"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall_rule.test",
						tfjsonpath.New("port"),
						knownvalue.StringExact("80"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall_rule.test",
						tfjsonpath.New("source"),
						knownvalue.StringExact("custom"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_firewall_rule.test",
						tfjsonpath.New("source_detail"),
						knownvalue.StringExact("any"),
					),
				},
			},
		},
	})
}

func testAccVPSFirewallRuleResourcePreCheck(t *testing.T) {
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
				"name": "test1",
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
			"operationId":      "VPS_createFirewallRuleV1"
		},
		"httpResponse": {
			"statusCode": 200,
			"body": {
				"id": 24541,
				"action": "accept",
				"protocol": "TCP",
				"port": "80",
				"source": "custom",
				"source_detail": "any"
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
				"name": "test1",
				"is_synced": false,
				"rules": [
					{
						"id": 24541,
						"action": "accept",
						"protocol": "TCP",
						"port": "80",
						"source": "custom",
						"source_detail": "any"
					}
				],
				"created_at": "2021-09-01T12:00:00Z",
				"updated_at": "2021-09-01T12:00:00Z"
			}			
		},
		"times": {
			"remainingTimes": 5
		}
	},	
	{
		"httpRequest": {
			"specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
			"operationId":      "VPS_updateFirewallRuleV1"
		},
		"httpResponse": {
			"statusCode": 200,
			"body": {
				"id": 24541,
				"action": "accept",
				"protocol": "UDP",
				"port": "80",
				"source": "custom",
				"source_detail": "any"
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
				"name": "test1",
				"is_synced": false,
				"rules": [
					{
						"id": 24541,
						"action": "accept",
						"protocol": "UDP",
						"port": "80",
						"source": "custom",
						"source_detail": "any"
					}
				],
				"created_at": "2021-09-01T12:00:00Z",
				"updated_at": "2021-09-01T12:00:00Z"
			}			
		},
		"times": {
			"remainingTimes": 4
		}
	},
	{
		"httpRequest": {
			"specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
			"operationId":      "VPS_deleteFirewallRuleV1"
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
				"name": "test2",
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
			"operationId":      "VPS_createFirewallRuleV1"
		},
		"httpResponse": {
			"statusCode": 200,
			"body": {
				"id": 24542,
				"action": "accept",
				"protocol": "UDP",
				"port": "80",
				"source": "custom",
				"source_detail": "any"
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
				"name": "test2",
				"is_synced": false,
				"rules": [
					{
						"id": 24542,
						"action": "accept",
						"protocol": "UDP",
						"port": "80",
						"source": "custom",
						"source_detail": "any"
					}
				],
				"created_at": "2021-09-01T12:00:00Z",
				"updated_at": "2021-09-01T12:00:00Z"
			}			
		},
		"times": {
			"remainingTimes": 2
		}
	},	
	{
		"httpRequest": {
			"specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
			"operationId":      "VPS_deleteFirewallRuleV1"
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

func testAccVPSFirewallRuleResourceConfig(name, protocol string) string {
	return fmt.Sprintf(`
resource "hostinger_vps_firewall" "test" {
  name = %[1]q
}

resource "hostinger_vps_firewall_rule" "test" {
  protocol = %[2]q
  firewall_id = hostinger_vps_firewall.test.id
  port     = "80"
  source   = "custom"
  source_detail = "any"
}
`, name, protocol)
}
