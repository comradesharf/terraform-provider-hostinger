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

func TestAccReachSegmentsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccReachSegmentsDataSourcePreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReachSegmentsDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.hostinger_reach_segments.test",
						tfjsonpath.New("segments").AtSliceIndex(0),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"uuid":       knownvalue.StringExact("550e8400-e29b-41d4-a716-446655440000"),
							"name":       knownvalue.StringExact("Newsletter Subscribers"),
							"created_at": knownvalue.StringExact("2025-02-27T11:54:22Z"),
							"updated_at": knownvalue.StringExact("2025-02-27T11:54:22Z"),
						}),
					),
				},
			},
		},
	})
}

func testAccReachSegmentsDataSourcePreCheck(t *testing.T) {
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
      "operationId": "reach_listSegmentsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": [
        {
          "uuid": "550e8400-e29b-41d4-a716-446655440000",
          "name": "Newsletter Subscribers",
          "created_at": "2025-02-27T11:54:22Z",
          "updated_at": "2025-02-27T11:54:22Z"
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

const testAccReachSegmentsDataSourceConfig = `
data "hostinger_reach_segments" "test" {
}
`
