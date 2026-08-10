// Copyright (c) HashiCorp, Inc.
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
	mockserver "github.com/mock-server/mockserver-monorepo/mockserver-client-go/v7"
)

func TestAccVPSMonarxResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVPSMonarxResourcePreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVPSMonarxResourceConfig(17923, false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"hostinger_vps_monarx.test",
						tfjsonpath.New("virtual_machine_id"),
						knownvalue.Int64Exact(17923),
					),
				},
			},
			{
				Config: testAccVPSMonarxResourceConfig(17924, true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"hostinger_vps_monarx.test",
						tfjsonpath.New("virtual_machine_id"),
						knownvalue.Int64Exact(17924),
					),
				},
			},
		},
	})
}

func testAccVPSMonarxResourceConfig(virtualMachineID int, waitForAction bool) string {
	return fmt.Sprintf(`
resource "hostinger_vps_monarx" "test" {
  virtual_machine_id = %d
  wait_for_action = %t
}`, virtualMachineID, waitForAction)
}

func testAccVPSMonarxResourcePreCheck(t *testing.T) {
	client := newMockServerClient()
	if err := client.Clear(nil, mockserver.ClearAll); err != nil {
		t.Fatalf("failed to clear mock server: %v", err)
	}
	if err := client.FreezeClock("2021-09-01T12:00:00Z"); err != nil {
		t.Fatalf("failed to freeze mock server clock: %v", err)
	}

	// language=json
	expectations := []byte(`[
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_installMonarxV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "id": 8123720,
        "name": "install_monarx",
        "created_at": "2024-09-05T07:25:36Z",
        "updated_at": "2024-09-05T07:25:36Z",
        "state": "success"
      }
    },
    "times": {
      "remainingTimes": 1
    }
  },
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_getScanMetricsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "records": 1,
        "malicious": 2,
        "compromised": 3,
        "scanned_files": 193218,
        "scan_started_at": "2025-02-27T11:54:22Z",
        "scan_ended_at": "2025-03-27T11:54:22Z"
      }
    },
    "times": {
      "remainingTimes": 2
    }
  },
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_uninstallMonarxV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "id": 8123721,
        "name": "uninstall_monarx",
        "created_at": "2024-09-05T07:25:36Z",
        "updated_at": "2024-09-05T07:25:36Z",
        "state": "success"
      }
    },
    "times": {
      "remainingTimes": 1
    }
  },
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_installMonarxV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "id": 8123722,
        "name": "install_monarx",
        "created_at": "2024-09-05T07:25:36Z",
        "updated_at": "2024-09-05T07:25:36Z",
        "state": "created"
      }
    },
    "times": {
      "remainingTimes": 1
    }
  },
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_getActionDetailsV1"
    },
    "httpResponses": [
      {
        "statusCode": 200,
        "body": {
          "id": 8123722,
          "name": "install_monarx",
          "created_at": "2024-09-05T07:25:36Z",
          "updated_at": "2024-09-05T07:25:36Z",
          "state": "created"
        }
      },
      {
        "statusCode": 200,
        "body": {
          "id": 8123722,
          "name": "install_monarx",
          "created_at": "2024-09-05T07:25:36Z",
          "updated_at": "2024-09-05T07:25:36Z",
          "state": "created"
        }
      },
      {
        "statusCode": 200,
        "body": {
          "id": 8123722,
          "name": "install_monarx",
          "created_at": "2024-09-05T07:25:36Z",
          "updated_at": "2024-09-05T07:25:36Z",
          "state": "success"
        }
      }
    ],
    "times": {
      "remainingTimes": 3
    }
  },
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_getScanMetricsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "records": 1,
        "malicious": 2,
        "compromised": 3,
        "scanned_files": 193218,
        "scan_started_at": "2025-02-27T11:54:22Z",
        "scan_ended_at": "2025-03-27T11:54:22Z"
      }
    },
    "times": {
      "remainingTimes": 1
    }
  },
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_installMonarxV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "id": 8123723,
        "name": "install_monarx",
        "created_at": "2024-09-05T07:25:36Z",
        "updated_at": "2024-09-05T07:25:36Z",
        "state": "created"
      }
    },
    "times": {
      "remainingTimes": 1
    }
  },
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_getActionDetailsV1"
    },
    "httpResponses": [
      {
        "statusCode": 200,
        "body": {
          "id": 8123723,
          "name": "install_monarx",
          "created_at": "2024-09-05T07:25:36Z",
          "updated_at": "2024-09-05T07:25:36Z",
          "state": "created"
        }
      },
      {
        "statusCode": 200,
        "body": {
          "id": 8123723,
          "name": "install_monarx",
          "created_at": "2024-09-05T07:25:36Z",
          "updated_at": "2024-09-05T07:25:36Z",
          "state": "created"
        }
      },
      {
        "statusCode": 200,
        "body": {
          "id": 8123723,
          "name": "install_monarx",
          "created_at": "2024-09-05T07:25:36Z",
          "updated_at": "2024-09-05T07:25:36Z",
          "state": "success"
        }
      }
    ],
    "times": {
      "remainingTimes": 3
    }
  },
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_uninstallMonarxV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "id": 8123724,
        "name": "uninstall_monarx",
        "created_at": "2024-09-05T07:25:36Z",
        "updated_at": "2024-09-05T07:25:36Z",
        "state": "created"
      }
    },
    "times": {
      "remainingTimes": 1
    }
  },
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_getActionDetailsV1"
    },
    "httpResponses": [
      {
        "statusCode": 200,
        "body": {
          "id": 8123724,
          "name": "uninstall_monarx",
          "created_at": "2024-09-05T07:25:36Z",
          "updated_at": "2024-09-05T07:25:36Z",
          "state": "created"
        }
      },
      {
        "statusCode": 200,
        "body": {
          "id": 8123724,
          "name": "uninstall_monarx",
          "created_at": "2024-09-05T07:25:36Z",
          "updated_at": "2024-09-05T07:25:36Z",
          "state": "created"
        }
      },
      {
        "statusCode": 200,
        "body": {
          "id": 8123724,
          "name": "uninstall_monarx",
          "created_at": "2024-09-05T07:25:36Z",
          "updated_at": "2024-09-05T07:25:36Z",
          "state": "success"
        }
      }
    ],
    "times": {
      "remainingTimes": 3
    }
  }
]`)

	req, err := http.NewRequest("PUT", "http://localhost:1234/mockserver/expectation", bytes.NewReader(expectations))
	if err != nil {
		t.Fatalf("failed to create mock server expectation request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if _, err := http.DefaultClient.Do(req); err != nil {
		t.Fatalf("failed to configure mock server expectations: %v", err)
	}
}
