// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mock-server/mockserver-monorepo/mockserver-client-go/v7"
)

func TestAccVPSPostInstallScriptResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVPSPostInstallScriptResourcePreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVPSPostInstallScriptResourceConfig("bootstrap", "#!/bin/sh\napt-get update"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"hostinger_vps_post_install_script.test",
						tfjsonpath.New("id"),
						knownvalue.Int64Exact(781),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_post_install_script.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("bootstrap"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_post_install_script.test",
						tfjsonpath.New("content"),
						knownvalue.StringExact("#!/bin/sh\napt-get update"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_post_install_script.test",
						tfjsonpath.New("created_at"),
						knownvalue.StringExact("2021-09-01T12:00:00Z"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_post_install_script.test",
						tfjsonpath.New("updated_at"),
						knownvalue.StringExact("2021-09-01T12:00:00Z"),
					),
				},
			},
			{
				ResourceName:      "hostinger_vps_post_install_script.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccVPSPostInstallScriptResourceConfig("bootstrap-updated", "#!/bin/sh\napt-get update\napt-get install -y nginx"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"hostinger_vps_post_install_script.test",
						tfjsonpath.New("id"),
						knownvalue.Int64Exact(781),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_post_install_script.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("bootstrap-updated"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_post_install_script.test",
						tfjsonpath.New("content"),
						knownvalue.StringExact("#!/bin/sh\napt-get update\napt-get install -y nginx"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_post_install_script.test",
						tfjsonpath.New("created_at"),
						knownvalue.StringExact("2021-09-01T12:00:00Z"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_post_install_script.test",
						tfjsonpath.New("updated_at"),
						knownvalue.StringExact("2021-09-01T12:30:00Z"),
					),
				},
			},
		},
	})
}

func testAccVPSPostInstallScriptResourcePreCheck(t *testing.T) {
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
      "operationId": "VPS_createPostInstallScriptV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "id": 781,
        "name": "bootstrap",
        "content": "#!/bin/sh\napt-get update",
        "created_at": "2021-09-01T12:00:00Z",
        "updated_at": "2021-09-01T12:00:00Z"
      }
    },
    "times": {
      "remainingTimes": 1
    }
  },
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_getPostInstallScriptV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "id": 781,
        "name": "bootstrap",
        "content": "#!/bin/sh\napt-get update",
        "created_at": "2021-09-01T12:00:00Z",
        "updated_at": "2021-09-01T12:00:00Z"
      }
    },
    "times": {
      "remainingTimes": 3
    }
  },
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_updatePostInstallScriptV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "id": 781,
        "name": "bootstrap-updated",
        "content": "#!/bin/sh\napt-get update\napt-get install -y nginx",
        "created_at": "2021-09-01T12:00:00Z",
        "updated_at": "2021-09-01T12:30:00Z"
      }
    },
    "times": {
      "remainingTimes": 1
    }
  },
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_getPostInstallScriptV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "id": 781,
        "name": "bootstrap-updated",
        "content": "#!/bin/sh\napt-get update\napt-get install -y nginx",
        "created_at": "2021-09-01T12:00:00Z",
        "updated_at": "2021-09-01T12:30:00Z"
      }
    },
    "times": {
      "remainingTimes": 1
    }
  },
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_deletePostInstallScriptV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "message": "Request accepted"
      }
    },
    "times": {
      "remainingTimes": 1
    }
  }
]`)

	req, err := http.NewRequest(
		"PUT",
		"http://localhost:1234/mockserver/expectation",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("failed to create expectation request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if _, err := http.DefaultClient.Do(req); err != nil {
		t.Fatalf("failed to create expectations: %v", err)
	}
}

func testAccVPSPostInstallScriptResourceConfig(name, content string) string {
	return fmt.Sprintf(`
resource "hostinger_vps_post_install_script" "test" {
  name    = %[1]q
  content = %[2]q
}
`, name, content)
}
