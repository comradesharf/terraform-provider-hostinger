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

func TestAccVPSPostInstallScriptsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t); testAccVPSPostInstallScriptsDataSourcePreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: `data "hostinger_vps_post_install_scripts" "test" {}`,
			ConfigStateChecks: []statecheck.StateCheck{statecheck.ExpectKnownValue(
				"data.hostinger_vps_post_install_scripts.test",
				tfjsonpath.New("scripts"),
				knownvalue.ListExact([]knownvalue.Check{
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"id":         knownvalue.Int64Exact(123),
						"name":       knownvalue.StringExact("bootstrap"),
						"content":    knownvalue.StringExact("#!/bin/sh"),
						"created_at": knownvalue.StringExact("2021-09-01T12:00:00Z"),
						"updated_at": knownvalue.StringExact("2021-09-01T12:00:00Z"),
					}),
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"id":         knownvalue.Int64Exact(124),
						"name":       knownvalue.StringExact("bootstrap"),
						"content":    knownvalue.StringExact("#!/bin/sh"),
						"created_at": knownvalue.StringExact("2021-09-01T12:00:00Z"),
						"updated_at": knownvalue.StringExact("2021-09-01T12:00:00Z"),
					}),
				}),
			)},
		}},
	})
}

func testAccVPSPostInstallScriptsDataSourcePreCheck(t *testing.T) {
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
      "operationId": "VPS_getPostInstallScriptsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "id": 123,
            "name": "bootstrap",
            "content": "#!/bin/sh",
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
      "operationId": "VPS_getPostInstallScriptsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "id": 124,
            "name": "bootstrap",
            "content": "#!/bin/sh",
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
  },
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_getPostInstallScriptsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "id": 123,
            "name": "bootstrap",
            "content": "#!/bin/sh",
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
      "operationId": "VPS_getPostInstallScriptsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "id": 124,
            "name": "bootstrap",
            "content": "#!/bin/sh",
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
  },
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_getPostInstallScriptsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "id": 123,
            "name": "bootstrap",
            "content": "#!/bin/sh",
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
      "operationId": "VPS_getPostInstallScriptsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "id": 124,
            "name": "bootstrap",
            "content": "#!/bin/sh",
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
]`)
	req, _ := http.NewRequest("PUT", "http://localhost:1234/mockserver/expectation", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if _, err := http.DefaultClient.Do(req); err != nil {
		t.Fatalf("failed to create mock expectation: %v", err)
	}
}
