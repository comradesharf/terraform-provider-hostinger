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
	compareValuesDiffer := statecheck.CompareValue(compare.ValuesDiffer())

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccResourceVPSFirewallPreCheck(t)
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
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id":            knownvalue.Int64Exact(24541),
								"action":        knownvalue.StringExact("accept"),
								"protocol":      knownvalue.StringExact("TCP"),
								"port":          knownvalue.StringExact("1024:2048"),
								"source":        knownvalue.StringExact("any"),
								"source_detail": knownvalue.StringExact("any"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id":            knownvalue.Int64Exact(24542),
								"action":        knownvalue.StringExact("accept"),
								"protocol":      knownvalue.StringExact("TCP"),
								"port":          knownvalue.StringExact("8080:8090"),
								"source":        knownvalue.StringExact("any"),
								"source_detail": knownvalue.StringExact("any"),
							}),
						}),
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
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id":            knownvalue.Int64Exact(24541),
								"action":        knownvalue.StringExact("accept"),
								"protocol":      knownvalue.StringExact("TCP"),
								"port":          knownvalue.StringExact("1024:2048"),
								"source":        knownvalue.StringExact("any"),
								"source_detail": knownvalue.StringExact("any"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id":            knownvalue.Int64Exact(24542),
								"action":        knownvalue.StringExact("accept"),
								"protocol":      knownvalue.StringExact("TCP"),
								"port":          knownvalue.StringExact("8080:8090"),
								"source":        knownvalue.StringExact("any"),
								"source_detail": knownvalue.StringExact("any"),
							}),
						}),
					),
				},
			},
		},
	})
}

func testAccResourceVPSFirewallPreCheck(t *testing.T) {
	client := newMockServerClient()

	if err := client.Clear(nil, mockserver.ClearAll); err != nil {
		t.Fatalf("failed to clear mock server: %v", err)
	}

	if err := client.FreezeClock("2021-09-01T12:00:00Z"); err != nil {
		t.Fatalf("failed to freeze mock server clock: %v", err)
	}

	expect1 := mockserver.Expectation{
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
	"name": "one",
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
    "id": 65225,
	"name": "two",
	"is_synced": false,
	"rules": [],
	"created_at": "2021-09-01T12:00:00Z",
	"updated_at": "2021-09-01T12:00:00Z"
}
`).
				BuildPtr(),
		},
	}

	expect2 := mockserver.Expectation{
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
	"port": "1024:2048",
	"source": "any",
	"source_detail": "any"
}
`).
				BuildPtr(),
			mockserver.Response().
				StatusCode(200).
				JSONBody(
					// language=json
					`
{
	"id": 24542,
	"action": "accept",
	"protocol": "TCP",
	"port": "8080:8090",
	"source": "any",
	"source_detail": "any"
}
`).
				BuildPtr(),
			mockserver.Response().
				StatusCode(200).
				JSONBody(
					// language=json
					`
{
	"id": 24543,
	"action": "accept",
	"protocol": "TCP",
	"port": "8080:8090",
	"source": "any",
	"source_detail": "any"
}
`).
				BuildPtr(),
			mockserver.Response().
				StatusCode(200).
				JSONBody(
					// language=json
					`
{
	"id": 24544,
	"action": "accept",
	"protocol": "TCP",
	"port": "8080:8090",
	"source": "any",
	"source_detail": "any"
}
`).
				BuildPtr(),
		},
	}

	expect3 := mockserver.Expectation{
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
	}

	expect4 := mockserver.Expectation{
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
	}

	if _, err := client.Upsert(expect1, expect2, expect3, expect4); err != nil {
		t.Fatalf("failed to upsert mock server expectation: %v", err)
	}

}

func testAccResourceVPSFirewallConfig(name string) string {
	return fmt.Sprintf(`
resource "hostinger_vps_firewall" "test" {
  name = %[1]q
  rules = [
	{
		"protocol" = "TCP"
		"port" = "1024:2048"
		"source" = "any"
		"source_detail" = "any"
    },
    {
		"protocol" = "TCP"
		"port" = "8080:8090"
		"source" = "any"
		"source_detail" = "any"
    }
  ]
}
`, name)
}
