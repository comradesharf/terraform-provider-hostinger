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

func TestAccVPSPublicKeyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVPSPublicKeyResourcePreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVPSPublicKeyResourceConfig("test1", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDone"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"hostinger_vps_public_key.test",
						tfjsonpath.New("id"),
						knownvalue.Int64Exact(325),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_public_key.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("test1"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_public_key.test",
						tfjsonpath.New("key"),
						knownvalue.StringExact("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDone"),
					),
				},
			},
			{
				ResourceName:      "hostinger_vps_public_key.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccVPSPublicKeyResourceConfig("test2", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIDone"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"hostinger_vps_public_key.test",
						tfjsonpath.New("id"),
						knownvalue.Int64Exact(326),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_public_key.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("test2"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_public_key.test",
						tfjsonpath.New("key"),
						knownvalue.StringExact("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIDone"),
					),
				},
			},
		},
	})
}

func testAccVPSPublicKeyResourcePreCheck(t *testing.T) {
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
      "operationId": "VPS_createPublicKeyV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "id": 325,
        "name": "test1",
        "key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDone"
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
            "name": "test1",
            "key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDone"
          }
        ],
        "meta": {
          "total": 1,
          "current_page": 1,
          "per_page": 1
        }
      }
    },
    "times": {
      "remainingTimes": 3
    }
  },
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_deletePublicKeyV1"
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
  },
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_createPublicKeyV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "id": 326,
        "name": "test2",
        "key": "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIDone"
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
            "name": "test2",
            "key": "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIDone"
          }
        ],
        "meta": {
          "total": 1,
          "current_page": 1,
          "per_page": 1
        }
      }
    },
    "times": {
      "remainingTimes": 3
    }
  },
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_deletePublicKeyV1"
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

func testAccVPSPublicKeyResourceConfig(name, key string) string {
	return fmt.Sprintf(`
resource "hostinger_vps_public_key" "test" {
  name = %[1]q
  key  = %[2]q
}
`, name, key)
}
