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

func TestAccVPSTemplateDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVPSTemplateDataSourcePreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: `data "hostinger_vps_template" "test" { id = 67890 }`,
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue("data.hostinger_vps_template.test", tfjsonpath.New("id"), knownvalue.Int64Exact(67890)),
				statecheck.ExpectKnownValue("data.hostinger_vps_template.test", tfjsonpath.New("name"), knownvalue.StringExact("Ubuntu 24.04 LTS")),
				statecheck.ExpectKnownValue("data.hostinger_vps_template.test", tfjsonpath.New("description"), knownvalue.StringExact("Ubuntu 24.04 LTS")),
				statecheck.ExpectKnownValue("data.hostinger_vps_template.test", tfjsonpath.New("documentation"), knownvalue.StringExact("https://ubuntu.com/server/docs")),
			},
		}},
	})
}

func testAccVPSTemplateDataSourcePreCheck(t *testing.T) {
	client := newMockServerClient()

	if err := client.Clear(nil, mockserver.ClearAll); err != nil {
		t.Fatalf("failed to clear mock server: %v", err)
	}

	// language=json
	body := []byte(`[
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_getTemplateDetailsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "id": 67890,
        "name": "Ubuntu 24.04 LTS",
        "description": "Ubuntu 24.04 LTS",
        "documentation": "https://ubuntu.com/server/docs"
      }
    },
    "times": {
      "remainingTimes": 3
    }
  }
]`)
	req, _ := http.NewRequest("PUT", "http://localhost:1234/mockserver/expectation", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if _, err := http.DefaultClient.Do(req); err != nil {
		t.Fatalf("failed to create mock expectation: %v", err)
	}
}
