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

func TestAccVPSVirtualMachineNameserversResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVPSVirtualMachineNameserversResourcePreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVPSVirtualMachineNameserversResourceConfig(false, 17923, "1.1.1.1", "8.8.8.8"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"hostinger_vps_virtual_machine_nameservers.test",
						tfjsonpath.New("virtual_machine_id"),
						knownvalue.Int64Exact(17923),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_virtual_machine_nameservers.test",
						tfjsonpath.New("ns1"),
						knownvalue.StringExact("1.1.1.1"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_virtual_machine_nameservers.test",
						tfjsonpath.New("ns2"),
						knownvalue.StringExact("8.8.8.8"),
					),
				},
			},
			{
				Config: testAccVPSVirtualMachineNameserversResourceConfig(false, 17923, "2.2.2.2", "9.9.9.9"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"hostinger_vps_virtual_machine_nameservers.test",
						tfjsonpath.New("virtual_machine_id"),
						knownvalue.Int64Exact(17923),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_virtual_machine_nameservers.test",
						tfjsonpath.New("ns1"),
						knownvalue.StringExact("2.2.2.2"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_virtual_machine_nameservers.test",
						tfjsonpath.New("ns2"),
						knownvalue.StringExact("9.9.9.9"),
					),
				},
			},
			{
				Config: testAccVPSVirtualMachineNameserversResourceConfig(false, 17923, "2.2.2.2", "9.9.9.9"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"hostinger_vps_virtual_machine_nameservers.test",
						tfjsonpath.New("virtual_machine_id"),
						knownvalue.Int64Exact(17923),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_virtual_machine_nameservers.test",
						tfjsonpath.New("ns1"),
						knownvalue.StringExact("2.2.2.2"),
					),
					statecheck.ExpectKnownValue(
						"hostinger_vps_virtual_machine_nameservers.test",
						tfjsonpath.New("ns2"),
						knownvalue.StringExact("9.9.9.9"),
					),
				},
			},
		},
	})
}

func testAccVPSVirtualMachineNameserversResourceConfig(waitForAction bool, virtualMachineID int, ns1, ns2 string) string {
	return fmt.Sprintf(`
resource "hostinger_vps_virtual_machine_nameservers" "test" {
	virtual_machine_id = %d
	ns1 = "%s"
	ns2 = "%s"
	wait_for_action = %t
}`, virtualMachineID, ns1, ns2, waitForAction)
}

func testAccVPSVirtualMachineNameserversResourcePreCheck(t *testing.T) {
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
      "operationId": "VPS_setNameserversV1"
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
        "hostname": "srv17923.hstgr.cloud",
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
      "operationId": "VPS_setNameserversV1"
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
        "hostname": "srv17923.hstgr.cloud",
        "state": "running",
        "actions_lock": "unlocked",
        "cpus": 4,
        "memory": 8192,
        "disk": 51200,
        "bandwidth": 1073741824,
        "ns1": "2.2.2.2",
        "ns2": "9.9.9.9",
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
      "remainingTimes": 3
    }
  },
  {
    "httpRequest": {
      "specUrlOrPayload": "https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json",
      "operationId": "VPS_setNameserversV1"
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
