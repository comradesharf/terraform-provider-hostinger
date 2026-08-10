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

func TestAccVPSDataCentersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVPSDataCentersDataSourcePreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "hostinger_vps_data_centers" "test" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_data_centers.test",
						tfjsonpath.New("data_centers"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id":        knownvalue.Int64Exact(1),
								"name":      knownvalue.StringExact("Lithuania - Vilnius"),
								"city":      knownvalue.StringExact("Vilnius"),
								"continent": knownvalue.StringExact("Europe"),
								"location":  knownvalue.StringExact("LT"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id":        knownvalue.Int64Exact(2),
								"name":      knownvalue.StringExact("United States - Phoenix"),
								"city":      knownvalue.StringExact("Phoenix"),
								"continent": knownvalue.StringExact("North America"),
								"location":  knownvalue.StringExact("US"),
							}),
						}),
					),
				},
			},
		},
	})
}

func testAccVPSDataCentersDataSourcePreCheck(t *testing.T) {
	client := newMockServerClient()

	if err := client.Clear(nil, mockserver.ClearAll); err != nil {
		t.Fatalf("failed to clear mock server: %v", err)
	}

	body := []byte(`[
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_getDataCenterListV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": [
        {
          "id": 1,
          "name": "Lithuania - Vilnius",
          "city": "Vilnius",
          "continent": "Europe",
          "location": "LT"
        },
        {
          "id": 2,
          "name": "United States - Phoenix",
          "city": "Phoenix",
          "continent": "North America",
          "location": "US"
        }
      ]
    },
    "times": {
      "remainingTimes": 3
    }
  }
]`)

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
