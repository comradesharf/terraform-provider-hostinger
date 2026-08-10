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
	mockserver "github.com/mock-server/mockserver-monorepo/mockserver-client-go/v7"
)

func TestAccVPSTemplatesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVPSTemplatesDataSourcePreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: `data "hostinger_vps_templates" "test" {}`,
			ConfigStateChecks: []statecheck.StateCheck{statecheck.ExpectKnownValue(
				"data.hostinger_vps_templates.test",
				tfjsonpath.New("templates"),
				knownvalue.ListExact([]knownvalue.Check{
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"id":            knownvalue.Int64Exact(6523),
						"name":          knownvalue.StringExact("Ubuntu 20.04 LTS"),
						"description":   knownvalue.StringExact("Ubuntu 20.04 LTS"),
						"documentation": knownvalue.StringExact("https://docs.ubuntu.com"),
					}),
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"id":            knownvalue.Int64Exact(6524),
						"name":          knownvalue.StringExact("Debian 12"),
						"description":   knownvalue.StringExact("Debian 12"),
						"documentation": knownvalue.StringExact("https://www.debian.org/doc/"),
					}),
				}),
			)},
		}},
	})
}

func testAccVPSTemplatesDataSourcePreCheck(t *testing.T) {
	client := newMockServerClient()

	if err := client.Clear(nil, mockserver.ClearAll); err != nil {
		t.Fatalf("failed to clear mock server: %v", err)
	}

	// language=json
	body := []byte(`[
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_getTemplatesV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": [
        {
          "id": 6523,
          "name": "Ubuntu 20.04 LTS",
          "description": "Ubuntu 20.04 LTS",
          "documentation": "https://docs.ubuntu.com"
        },
        {
          "id": 6524,
          "name": "Debian 12",
          "description": "Debian 12",
          "documentation": "https://www.debian.org/doc/"
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
