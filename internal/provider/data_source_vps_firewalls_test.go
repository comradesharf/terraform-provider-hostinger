// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mock-server/mockserver-monorepo/mockserver-client-go/v7"
)

func TestAccDataSourceVPSFirewalls(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccDataSourceVPSFirewallsPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVPSFirewallsConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_firewalls.test",
						tfjsonpath.New("firewalls"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id":        knownvalue.Int64Exact(65224),
								"name":      knownvalue.StringExact("HTTP and SSH only"),
								"is_synced": knownvalue.Bool(false),
								"rules": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"id":            knownvalue.Int64Exact(24541),
										"action":        knownvalue.StringExact("accept"),
										"protocol":      knownvalue.StringExact("TCP"),
										"port":          knownvalue.StringExact("1024:2048"),
										"source":        knownvalue.StringExact("any"),
										"source_detail": knownvalue.StringExact("any"),
									}),
								}),
								"created_at": knownvalue.StringExact("2021-09-01T12:00:00Z"),
								"updated_at": knownvalue.StringExact("2021-09-01T12:00:00Z"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id":        knownvalue.Int64Exact(65225),
								"name":      knownvalue.StringExact("HTTP and SSH only"),
								"is_synced": knownvalue.Bool(false),
								"rules": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"id":            knownvalue.Int64Exact(24542),
										"action":        knownvalue.StringExact("accept"),
										"protocol":      knownvalue.StringExact("TCP"),
										"port":          knownvalue.StringExact("1024:2048"),
										"source":        knownvalue.StringExact("any"),
										"source_detail": knownvalue.StringExact("any"),
									}),
								}),
								"created_at": knownvalue.StringExact("2021-09-01T12:00:00Z"),
								"updated_at": knownvalue.StringExact("2021-09-01T12:00:00Z"),
							}),
						}),
					),
				},
			},
		},
	})
}

const testAccDataSourceVPSFirewallsConfig = `
data "hostinger_vps_firewalls" "test" {}
`

func testAccDataSourceVPSFirewallsPreCheck(t *testing.T) {
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
			OperationId:      "VPS_getFirewallListV1",
		},
		HttpResponses: []*mockserver.HttpResponse{
			mockserver.Response().
				StatusCode(200).
				JSONBody(
					// language=json
					`
{
	"data": [
		{
			"id": 65224,
			"name": "HTTP and SSH only",
			"is_synced": false,
			"rules": [
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
	],
	"meta": {
		"current_page": 1,
		"per_page": 1,
		"total": 2
	}
}
`).
				BuildPtr(),
			mockserver.Response().
				StatusCode(200).
				JSONBody(
					// language=json
					`
{
	"data": [
		{
			"id": 65225,
			"name": "HTTP and SSH only",
			"is_synced": false,
			"rules": [
				{
					"id": 24542,
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
	],
	"meta": {
		"current_page": 2,
		"per_page": 1,
		"total": 2
	}
}
`).
				BuildPtr(),
		},
	}

	if _, err := client.Upsert(expect1); err != nil {
		t.Fatalf("failed to upsert mock server expectation: %v", err)
	}

}
