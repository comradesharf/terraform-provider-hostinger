// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mock-server/mockserver-monorepo/mockserver-client-go/v7"
)

func TestAccVPSVirtualMachinesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVPSVirtualMachinesDataSourcePreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVPSVirtualMachinesDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_virtual_machines.test",
						tfjsonpath.New("virtual_machines"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id":                knownvalue.Int64Exact(17923),
								"firewall_group_id": knownvalue.Int64Exact(260),
								"subscription_id":   knownvalue.StringExact("Azz353Uhl1xC54pR0"),
								"data_center_id":    knownvalue.Int64Exact(521),
								"plan":              knownvalue.StringExact("KVM 4"),
								"hostname":          knownvalue.StringExact("srv17923.hstgr.cloud"),
								"state":             knownvalue.StringExact("running"),
								"actions_lock":      knownvalue.StringExact("unlocked"),
								"cpus":              knownvalue.Int64Exact(4),
								"memory":            knownvalue.Int64Exact(8192),
								"disk":              knownvalue.Int64Exact(51200),
								"bandwidth":         knownvalue.Int64Exact(1073741824),
								"ns1":               knownvalue.StringExact("1.1.1.1"),
								"ns2":               knownvalue.StringExact("8.8.8.8"),
								"ipv4": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"id":      knownvalue.Int64Exact(52347),
										"address": knownvalue.StringExact("213.211.223.15"),
										"ptr":     knownvalue.StringExact("something.domain.tld"),
									}),
								}),
								"ipv6": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"id":      knownvalue.Int64Exact(52347),
										"address": knownvalue.StringExact("213.211.223.15"),
										"ptr":     knownvalue.StringExact("something.domain.tld"),
									}),
								}),
								"template": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"id":            knownvalue.Int64Exact(6523),
									"name":          knownvalue.StringExact("Ubuntu 20.04 LTS"),
									"description":   knownvalue.StringExact("Ubuntu 20.04 LTS"),
									"documentation": knownvalue.StringExact("https://docs.ubuntu.com"),
								}),
								"created_at": knownvalue.StringExact("2024-09-05T07:25:36Z"),
							}),
						}),
					),
				},
			},
		},
	})
}

const testAccVPSVirtualMachinesDataSourceConfig = `
data "hostinger_vps_virtual_machines" "test" {}
`

func testAccVPSVirtualMachinesDataSourcePreCheck(t *testing.T) {
	client := newMockServerClient()

	if err := client.Clear(nil, mockserver.ClearAll); err != nil {
		t.Fatalf("failed to clear mock server: %v", err)
	}

	if err := client.FreezeClock("2021-09-01T12:00:00Z"); err != nil {
		t.Fatalf("failed to freeze mock server clock: %v", err)
	}

	// language=json
	body := []byte(`[
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_getVirtualMachinesV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": [
        {
          "id": 17923,
          "firewall_group_id": 260,
          "subscription_id": "Azz353Uhl1xC54pR0",
          "data_center_id": 521,
          "plan": "KVM 4",
          "hostname": "srv17923.hstgr.cloud",
          "state": "running",
          "actions_lock": "unlocked",
          "cpus": 4,
          "memory": 8192,
          "disk": 51200,
          "bandwidth": 1073741824,
          "ns1": "1.1.1.1",
          "ns2": "8.8.8.8",
          "ipv4": [
            {
              "id": 52347,
              "address": "213.211.223.15",
              "ptr": "something.domain.tld"
            }
          ],
          "ipv6": [
            {
              "id": 52347,
              "address": "213.211.223.15",
              "ptr": "something.domain.tld"
            }
          ],
          "template": {
            "id": 6523,
            "name": "Ubuntu 20.04 LTS",
            "description": "Ubuntu 20.04 LTS",
            "documentation": "https://docs.ubuntu.com"
          },
          "created_at": "2024-09-05T07:25:36Z"
        }
      ]
    },
    "times": {
      "remainingTimes": 3
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
