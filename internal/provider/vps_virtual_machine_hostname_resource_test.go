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
	mockserver "github.com/mock-server/mockserver-monorepo/mockserver-client-go/v7"
)

func TestAccVPSVirtualMachineHostnameResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVPSVirtualMachineHostnameResourcePreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVPSVirtualMachineHostnameResourceConfig(false, 17923, "srv17923.example.com"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"hostinger_vps_virtual_machine_hostname.test",
						tfjsonpath.New("virtual_machine_id"),
						knownvalue.Int64Exact(17923),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_virtual_machine_hostname.test",
						tfjsonpath.New("hostname"),
						knownvalue.StringExact("srv17923.example.com"),
					),
				},
			},
			{
				Config: testAccVPSVirtualMachineHostnameResourceConfig(false, 17923, "srv17923.updated.example.com"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"hostinger_vps_virtual_machine_hostname.test",
						tfjsonpath.New("hostname"),
						knownvalue.StringExact("srv17923.updated.example.com"),
					),
				},
			},
			{
				Config: testAccVPSVirtualMachineHostnameResourceConfig(true, 17924, "srv17924.example.com"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"hostinger_vps_virtual_machine_hostname.test",
						tfjsonpath.New("virtual_machine_id"),
						knownvalue.Int64Exact(17924),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_virtual_machine_hostname.test",
						tfjsonpath.New("hostname"),
						knownvalue.StringExact("srv17924.example.com"),
					),
				},
			},
			{
				Config: testAccVPSVirtualMachineHostnameResourceConfig(true, 17924, "srv17924.updated.example.com"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"hostinger_vps_virtual_machine_hostname.test",
						tfjsonpath.New("hostname"),
						knownvalue.StringExact("srv17924.updated.example.com"),
					),
				},
			},
		},
	})
}

func testAccVPSVirtualMachineHostnameResourceConfig(waitForAction bool, virtualMachineID int, hostname string) string {
	return fmt.Sprintf(`
resource "hostinger_vps_virtual_machine_hostname" "test" {
  virtual_machine_id = %d
  hostname           = %q
  wait_for_action    = %t
}`, virtualMachineID, hostname, waitForAction)
}

func testAccVPSVirtualMachineHostnameResourcePreCheck(t *testing.T) {
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
      "operationId": "VPS_setHostnameV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "id": 8123712,
        "name": "action_name",
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
      "operationId": "VPS_getVirtualMachineDetailsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "id": 17923,
        "firewall_group_id": 260,
        "subscription_id": "Azz353Uhl1xC54pR0",
        "data_center_id": 521,
        "plan": "KVM 4",
        "hostname": "srv17923.example.com",
        "state": "running",
        "actions_lock": "unlocked",
        "cpus": 4,
        "memory": 8192,
        "disk": 51200,
        "bandwidth": 1073741824,
        "ns1": "1.1.1.1",
        "ns2": "8.8.8.8",
        "ipv4": [
          {
            "id": 52347,
            "address": "213.211.223.15",
            "ptr": "something.domain.tld"
          }
        ],
        "ipv6": [
          {
            "id": 52348,
            "address": "2a00:1450:4001:81b::200e",
            "ptr": "something.domain.tld"
          }
        ],
        "template": {
          "id": 6523,
          "name": "Ubuntu 20.04 LTS",
          "description": "Ubuntu 20.04 LTS",
          "documentation": "https://docs.ubuntu.com"
        },
        "created_at": "2024-09-05T07:25:36Z"
      }
    },
    "times": {
      "remainingTimes": 2
    }
  },
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_setHostnameV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "id": 8123713,
        "name": "action_name",
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
      "operationId": "VPS_getVirtualMachineDetailsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "id": 17923,
        "firewall_group_id": 260,
        "subscription_id": "Azz353Uhl1xC54pR0",
        "data_center_id": 521,
        "plan": "KVM 4",
        "hostname": "srv17923.updated.example.com",
        "state": "running",
        "actions_lock": "unlocked",
        "cpus": 4,
        "memory": 8192,
        "disk": 51200,
        "bandwidth": 1073741824,
        "ns1": "1.1.1.1",
        "ns2": "8.8.8.8",
        "ipv4": [
          {
            "id": 52347,
            "address": "213.211.223.15",
            "ptr": "something.domain.tld"
          }
        ],
        "ipv6": [
          {
            "id": 52348,
            "address": "2a00:1450:4001:81b::200e",
            "ptr": "something.domain.tld"
          }
        ],
        "template": {
          "id": 6523,
          "name": "Ubuntu 20.04 LTS",
          "description": "Ubuntu 20.04 LTS",
          "documentation": "https://docs.ubuntu.com"
        },
        "created_at": "2024-09-05T07:25:36Z"
      }
    },
    "times": {
      "remainingTimes": 2
    }
  },
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_resetHostnameV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "id": 8123714,
        "name": "action_name",
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
      "operationId": "VPS_setHostnameV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "id": 8123715,
        "name": "action_name",
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
      "operationId": "VPS_getActionDetailsV1"
    },
    "httpResponses": [
      {
        "statusCode": 200,
        "body": {
          "id": 8123715,
          "name": "action_name",
          "created_at": "2024-09-05T07:25:36Z",
          "updated_at": "2024-09-05T07:25:36Z",
          "state": "created"
        }
      },
      {
        "statusCode": 200,
        "body": {
          "id": 8123715,
          "name": "action_name",
          "created_at": "2024-09-05T07:25:36Z",
          "updated_at": "2024-09-05T07:25:36Z",
          "state": "created"
        }
      },
      {
        "statusCode": 200,
        "body": {
          "id": 8123715,
          "name": "action_name",
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
      "operationId": "VPS_getVirtualMachineDetailsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "id": 17924,
        "firewall_group_id": 260,
        "subscription_id": "Azz353Uhl1xC54pR0",
        "data_center_id": 521,
        "plan": "KVM 4",
        "hostname": "srv17924.example.com",
        "state": "running",
        "actions_lock": "unlocked",
        "cpus": 4,
        "memory": 8192,
        "disk": 51200,
        "bandwidth": 1073741824,
        "ns1": "1.1.1.1",
        "ns2": "8.8.8.8",
        "ipv4": [
          {
            "id": 52347,
            "address": "213.211.223.15",
            "ptr": "something.domain.tld"
          }
        ],
        "ipv6": [
          {
            "id": 52348,
            "address": "2a00:1450:4001:81b::200e",
            "ptr": "something.domain.tld"
          }
        ],
        "template": {
          "id": 6523,
          "name": "Ubuntu 20.04 LTS",
          "description": "Ubuntu 20.04 LTS",
          "documentation": "https://docs.ubuntu.com"
        },
        "created_at": "2024-09-05T07:25:36Z"
      }
    },
    "times": {
      "remainingTimes": 2
    }
  },
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_setHostnameV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "id": 8123716,
        "name": "action_name",
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
      "operationId": "VPS_getActionDetailsV1"
    },
    "httpResponses": [
      {
        "statusCode": 200,
        "body": {
          "id": 8123716,
          "name": "action_name",
          "created_at": "2024-09-05T07:25:36Z",
          "updated_at": "2024-09-05T07:25:36Z",
          "state": "created"
        }
      },
      {
        "statusCode": 200,
        "body": {
          "id": 8123716,
          "name": "action_name",
          "created_at": "2024-09-05T07:25:36Z",
          "updated_at": "2024-09-05T07:25:36Z",
          "state": "created"
        }
      },
      {
        "statusCode": 200,
        "body": {
          "id": 8123716,
          "name": "action_name",
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
      "operationId": "VPS_getVirtualMachineDetailsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "id": 17924,
        "firewall_group_id": 260,
        "subscription_id": "Azz353Uhl1xC54pR0",
        "data_center_id": 521,
        "plan": "KVM 4",
        "hostname": "srv17924.updated.example.com",
        "state": "running",
        "actions_lock": "unlocked",
        "cpus": 4,
        "memory": 8192,
        "disk": 51200,
        "bandwidth": 1073741824,
        "ns1": "1.1.1.1",
        "ns2": "8.8.8.8",
        "ipv4": [
          {
            "id": 52347,
            "address": "213.211.223.15",
            "ptr": "something.domain.tld"
          }
        ],
        "ipv6": [
          {
            "id": 52348,
            "address": "2a00:1450:4001:81b::200e",
            "ptr": "something.domain.tld"
          }
        ],
        "template": {
          "id": 6523,
          "name": "Ubuntu 20.04 LTS",
          "description": "Ubuntu 20.04 LTS",
          "documentation": "https://docs.ubuntu.com"
        },
        "created_at": "2024-09-05T07:25:36Z"
      }
    },
    "times": {
      "remainingTimes": 1
    }
  },
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_resetHostnameV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "id": 8123717,
        "name": "action_name",
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
      "operationId": "VPS_getActionDetailsV1"
    },
    "httpResponses": [
      {
        "statusCode": 200,
        "body": {
          "id": 8123717,
          "name": "action_name",
          "created_at": "2024-09-05T07:25:36Z",
          "updated_at": "2024-09-05T07:25:36Z",
          "state": "created"
        }
      },
      {
        "statusCode": 200,
        "body": {
          "id": 8123717,
          "name": "action_name",
          "created_at": "2024-09-05T07:25:36Z",
          "updated_at": "2024-09-05T07:25:36Z",
          "state": "created"
        }
      },
      {
        "statusCode": 200,
        "body": {
          "id": 8123717,
          "name": "action_name",
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
