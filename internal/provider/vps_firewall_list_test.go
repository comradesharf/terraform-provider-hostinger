// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mock-server/mockserver-monorepo/mockserver-client-go/v7"
)

func TestAccVPSFirewallList(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVPSFirewallListPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ` `,
			},
			{
				Config:         testAccVPSFirewallListConfig,
				Query:          true,
				GenerateConfig: true,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectResourceKnownValues(
						"hostinger_vps_firewall.test",
						queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
							"id": knownvalue.Int64Exact(65224),
						}),
						[]querycheck.KnownValueCheck{
							{
								tfjsonpath.New("name"),
								knownvalue.StringExact("one"),
							},
						},
					),
					querycheck.ExpectResourceKnownValues(
						"hostinger_vps_firewall.test",
						queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
							"id": knownvalue.Int64Exact(65225),
						}),
						[]querycheck.KnownValueCheck{
							{
								tfjsonpath.New("name"),
								knownvalue.StringExact("two"),
							},
						},
					),
				},
			},
		},
	})
}

func testAccVPSFirewallListPreCheck(t *testing.T) {
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
      "operationId": "VPS_getFirewallListV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "id": 65224,
            "name": "one",
            "is_synced": false,
            "rules": [],
            "created_at": "2021-09-01T12:00:00Z",
            "updated_at": "2021-09-01T12:00:00Z"
          }
        ],
        "meta": {
          "total": 2,
          "current_page": 1,
          "per_page": 1
        }
      }
    },
    "times": {
      "remainingTimes": 1
    }
  },
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_getFirewallListV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "id": 65225,
            "name": "two",
            "is_synced": false,
            "rules": [],
            "created_at": "2021-09-01T12:00:00Z",
            "updated_at": "2021-09-01T12:00:00Z"
          }
        ],
        "meta": {
          "total": 2,
          "current_page": 2,
          "per_page": 1
        }
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

const testAccVPSFirewallListConfig = `
list "hostinger_vps_firewall" "test" {
	provider = hostinger
	include_resource = true
}
`
