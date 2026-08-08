// Copyright (c) HashiCorp, Inc.
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

func TestAccVPSPostInstallScriptDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVPSPostInstallScriptDataSourcePreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: `data "hostinger_vps_post_install_script" "test" { id = 123 }`,
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue("data.hostinger_vps_post_install_script.test",
					tfjsonpath.New("name"),
					knownvalue.StringExact("bootstrap"),
				),
				statecheck.ExpectKnownValue("data.hostinger_vps_post_install_script.test",
					tfjsonpath.New("content"),
					knownvalue.StringExact("#!/bin/sh"),
				),
			},
		}},
	})
}

func testAccVPSPostInstallScriptDataSourcePreCheck(t *testing.T) {
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
      "operationId": "VPS_getPostInstallScriptV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "id": 123,
        "name": "bootstrap",
        "content": "#!/bin/sh",
        "created_at": "2021-09-01T12:00:00Z",
        "updated_at": "2021-09-01T12:00:00Z"
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
