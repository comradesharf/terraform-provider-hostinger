// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/mock-server/mockserver-monorepo/mockserver-client-go/v7"
)

func TestAccVPSPublicKeysAttachAction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVPSPublicKeysAttachActionPreCheck(t)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVPSPublicKeysAttachActionConfig,
				PostApplyFunc: func() {
					client := newMockServerClient()

					err := client.VerifySequence(
						mockserver.Request().Path("/api/vps/v1/public-keys/attach/12345").Method("POST"),
						mockserver.Request().Path("/api/vps/v1/virtual-machines/12345/actions/8123715").Method("GET"),
						mockserver.Request().Path("/api/vps/v1/virtual-machines/12345/actions/8123715").Method("GET"),
						mockserver.Request().Path("/api/vps/v1/virtual-machines/12345/actions/8123715").Method("GET"),
					)
					if err != nil {
						t.Fatalf("failed to verify mock server expectations: %v", err)
					}
				},
			},
		},
	})
}

func testAccVPSPublicKeysAttachActionPreCheck(t *testing.T) {
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
      "operationId": "VPS_attachPublicKeyV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "id": 8123715,
        "name": "attach_public_key",
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
          "id": 8123715,
          "name": "attach_public_key",
          "created_at": "2024-09-05T07:25:36Z",
          "updated_at": "2024-09-05T07:25:36Z",
          "state": "created"
        }
      },
      {
        "statusCode": 200,
        "body": {
          "id": 8123715,
          "name": "attach_public_key",
          "created_at": "2024-09-05T07:25:36Z",
          "updated_at": "2024-09-05T07:25:36Z",
          "state": "created"
        }
      },
      {
        "statusCode": 200,
        "body": {
          "id": 8123715,
          "name": "attach_public_key",
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

const testAccVPSPublicKeysAttachActionConfig = `
resource "terraform_data" "test" {
	lifecycle {
		action_trigger {
			events = [after_create]
			actions = [action.hostinger_vps_public_keys_attach.test]
		}
	}
}

action "hostinger_vps_public_keys_attach" "test" {
	config {
		public_key_ids = [65224, 65225]
		virtual_machine_id = 12345
		wait_for_action = true
	}
}
`
