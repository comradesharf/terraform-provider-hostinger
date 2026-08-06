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

func TestAccVPSFirewallDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVPSFirewallDataSourcePreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVPSFirewallDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_firewall.test",
						tfjsonpath.New("id"),
						knownvalue.Int64Exact(65224),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_firewall.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("HTTP and SSH only"),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_firewall.test",
						tfjsonpath.New("is_synced"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_firewall.test",
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
						}),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_firewall.test",
						tfjsonpath.New("created_at"),
						knownvalue.StringExact("2021-09-01T12:00:00Z"),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_firewall.test",
						tfjsonpath.New("updated_at"),
						knownvalue.StringExact("2021-09-01T12:00:00Z"),
					),
				},
			},
		},
	})
}

func testAccVPSFirewallDataSourcePreCheck(t *testing.T) {
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
      "operationId": "VPS_getFirewallDetailsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
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

const testAccVPSFirewallDataSourceConfig = `
data "hostinger_vps_firewall" "test" {
	id = 65224
}
`
