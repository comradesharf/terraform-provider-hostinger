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

func TestAccVPSPublicKeysDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVPSPublicKeysDataSourcePreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testVPSPublicKeysDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_public_keys.test",
						tfjsonpath.New("public_keys"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id":   knownvalue.Int64Exact(325),
								"name": knownvalue.StringExact("My public key 1"),
								"key":  knownvalue.StringExact("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD..."),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id":   knownvalue.Int64Exact(326),
								"name": knownvalue.StringExact("My public key 2"),
								"key":  knownvalue.StringExact("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD..."),
							}),
						}),
					),
				},
			},
		},
	})
}

const testVPSPublicKeysDataSourceConfig = `
data "hostinger_vps_public_keys" "test" {
}
`

func testAccVPSPublicKeysDataSourcePreCheck(t *testing.T) {
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
      "operationId": "VPS_getPublicKeysV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "id": 325,
            "name": "My public key 1",
            "key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD..."
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
      "operationId": "VPS_getPublicKeysV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "id": 326,
            "name": "My public key 2",
            "key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD..."
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
      "operationId": "VPS_getPublicKeysV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "id": 325,
            "name": "My public key 1",
            "key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD..."
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
      "operationId": "VPS_getPublicKeysV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "id": 326,
            "name": "My public key 2",
            "key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD..."
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
      "operationId": "VPS_getPublicKeysV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "id": 325,
            "name": "My public key 1",
            "key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD..."
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
      "operationId": "VPS_getPublicKeysV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "id": 326,
            "name": "My public key 2",
            "key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD..."
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
