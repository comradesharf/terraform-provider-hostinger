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

func TestAccResourceVPSFirewallRule(t *testing.T) {
	compareValuesDiffer := statecheck.CompareValue(compare.ValuesDiffer())
	compareValuesSame := statecheck.CompareValue(compare.ValuesSame())

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccResourceVPSFirewallRulePreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVPSFirewallRuleConfig("TCP"),
				ConfigStateChecks: []statecheck.StateCheck{
					compareValuesSame.AddStateValue(
						"hostinger_vps_firewall_rule.test",
						tfjsonpath.New("firewall_id"),
					),
					compareValuesDiffer.AddStateValue(
						"hostinger_vps_firewall_rule.test",
						tfjsonpath.New("id"),
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
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceVPSFirewallRulePreCheck(t *testing.T) {
	client := newMockServerClient()

	if err := client.Clear(nil, mockserver.ClearAll); err != nil {
		t.Fatalf("failed to clear mock server: %v", err)
	}

	if err := client.FreezeClock("2021-09-01T12:00:00Z"); err != nil {
		t.Fatalf("failed to freeze mock server clock: %v", err)
	}

	if _, err := client.Upsert(
		mockserver.Expectation{
			HttpRequest: &mockserver.HttpRequest{
				SpecUrlOrPayload: "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
				OperationId:      "VPS_createNewFirewallV1",
			},
			HttpResponses: []*mockserver.HttpResponse{
				mockserver.Response().
					StatusCode(200).
					JSONBody(
						// language=json
						`
{
    "id": 65224,
	"name": "test",
	"is_synced": false,
	"rules": [],
	"created_at": "2021-09-01T12:00:00Z",
	"updated_at": "2021-09-01T12:00:00Z"
}
`).
					BuildPtr(),
			},
		},
		mockserver.Expectation{
			HttpRequest: &mockserver.HttpRequest{
				SpecUrlOrPayload: "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
				OperationId:      "VPS_createFirewallRuleV1",
			},
			HttpResponses: []*mockserver.HttpResponse{
				mockserver.Response().
					StatusCode(200).
					JSONBody(
						// language=json
						`
{
	"id": 24541,
	"action": "accept",
	"protocol": "TCP",
	"port": "80",
	"source": "custom",
	"source_detail": "any"
}
`).
					BuildPtr(),
			},
		},
		mockserver.Expectation{
			HttpRequest: &mockserver.HttpRequest{
				SpecUrlOrPayload: "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
				OperationId:      "VPS_deleteFirewallRuleV1",
			},
			HttpResponses: []*mockserver.HttpResponse{
				mockserver.Response().
					StatusCode(200).
					JSONBody(
						// language=json
						`
{
	"message": "Request accepted"
}
`).
					BuildPtr(),
			},
		},
		mockserver.Expectation{
			HttpRequest: &mockserver.HttpRequest{
				SpecUrlOrPayload: "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
				OperationId:      "VPS_deleteFirewallV1",
			},
			HttpResponses: []*mockserver.HttpResponse{
				mockserver.Response().
					StatusCode(200).
					JSONBody(
						// language=json
						`
{
	"message": "Request accepted"
}
`).
					BuildPtr(),
			},
		},
		mockserver.Expectation{
			HttpRequest: &mockserver.HttpRequest{
				SpecUrlOrPayload: "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
				OperationId:      "VPS_getFirewallDetailsV1",
			},
			HttpResponses: []*mockserver.HttpResponse{
				mockserver.Response().
					StatusCode(200).
					JSONBody(
						// language=json
						`
{
    "id": 65224,
	"name": "test",
	"is_synced": false,
	"rules": [],
	"created_at": "2021-09-01T12:00:00Z",
	"updated_at": "2021-09-01T12:00:00Z"
}
`).
					BuildPtr(),
				mockserver.Response().
					StatusCode(200).
					JSONBody(
						// language=json
						`
{
    "id": 65224,
	"name": "one",
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
`).
					BuildPtr(),
				mockserver.Response().
					StatusCode(200).
					JSONBody(
						// language=json
						`
{
    "id": 65224,
	"name": "one",
	"is_synced": false,
	"rules": [
		{
			"id": 24542,
			"action": "accept",
			"protocol": "TCP",
			"port": "8080:8090",
			"source": "any",
			"source_detail": "any"
		},
		{
	        "id": 24541,
			"action": "accept",
			"protocol": "TCP",
			"port": "1024:2048",
			"source": "any",
			"source_detail": "any"
		}
    ],
	"created_at": "2021-09-01T12:00:00Z",
	"updated_at": "2021-09-01T12:00:00Z"
}
`).
					BuildPtr(),
				mockserver.Response().
					StatusCode(200).
					JSONBody(
						// language=json
						`
{
    "id": 65225,
	"name": "two",
	"is_synced": false,
	"rules": [
		{
			"id": 24543,
			"action": "accept",
			"protocol": "TCP",
			"port": "8080:8090",
			"source": "any",
			"source_detail": "any"
		},
		{
	        "id": 24544,
			"action": "accept",
			"protocol": "TCP",
			"port": "1024:2048",
			"source": "any",
			"source_detail": "any"
		}
    ],
	"created_at": "2021-09-01T12:00:00Z",
	"updated_at": "2021-09-01T12:00:00Z"
}
`).
					BuildPtr(),
			},
		},
	); err != nil {
		t.Fatalf("failed to upsert mock server expectation: %v", err)
	}

}

func testAccResourceVPSFirewallRuleConfig(protocol string) string {
	return fmt.Sprintf(`
resource "hostinger_vps_firewall" "test" {
  name = "test"
}

resource "hostinger_vps_firewall_rule" "test" {
  protocol = %[1]q
  firewall_id = hostinger_vps_firewall.test.id
  port     = "80"
  source   = "custom"
  source_detail = "any"
}
`, protocol)
}
