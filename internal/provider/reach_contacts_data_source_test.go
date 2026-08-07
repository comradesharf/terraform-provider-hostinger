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

func TestAccReachContactsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccReachContactsDataSourcePreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReachContactsDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.hostinger_reach_contacts.test",
						tfjsonpath.New("contacts"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"uuid":                knownvalue.StringExact("550e8400-e29b-41d4-a716-446655440000"),
								"email":               knownvalue.StringExact("john.doe@example.com"),
								"name":                knownvalue.StringExact("John"),
								"surname":             knownvalue.StringExact("Doe"),
								"subscription_status": knownvalue.StringExact("subscribed"),
								"subscribed_at":       knownvalue.StringExact("2023-01-01T00:00:00Z"),
								"source":              knownvalue.StringExact("sync"),
								"note":                knownvalue.StringExact("VIP customer"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"uuid":                knownvalue.StringExact("550e8400-e29b-41d4-a716-446655440001"),
								"email":               knownvalue.StringExact("jane.doe@example.com"),
								"name":                knownvalue.StringExact("Jane"),
								"surname":             knownvalue.StringExact("Doe"),
								"subscription_status": knownvalue.StringExact("subscribed"),
								"subscribed_at":       knownvalue.StringExact("2023-01-01T00:00:00Z"),
								"source":              knownvalue.StringExact("sync"),
								"note":                knownvalue.StringExact("VIP customer"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_reach_contacts.test",
						tfjsonpath.New("group_uuid"),
						knownvalue.StringExact("550e8400-e29b-41d4-a716-446655440000"),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_reach_contacts.test",
						tfjsonpath.New("subscription_status"),
						knownvalue.StringExact("subscribed"),
					),
				},
			},
		},
	})
}

func testAccReachContactsDataSourcePreCheck(t *testing.T) {
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
      "operationId": "reach_listContactsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "uuid": "550e8400-e29b-41d4-a716-446655440000",
            "name": "John",
            "surname": "Doe",
            "email": "john.doe@example.com",
            "subscription_status": "subscribed",
            "subscribed_at": "2023-01-01T00:00:00Z",
            "source": "sync",
            "note": "VIP customer"
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
      "operationId": "reach_listContactsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "uuid": "550e8400-e29b-41d4-a716-446655440001",
            "name": "Jane",
            "surname": "Doe",
            "email": "jane.doe@example.com",
            "subscription_status": "subscribed",
            "subscribed_at": "2023-01-01T00:00:00Z",
            "source": "sync",
            "note": "VIP customer"
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
      "operationId": "reach_listContactsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "uuid": "550e8400-e29b-41d4-a716-446655440000",
            "name": "John",
            "surname": "Doe",
            "email": "john.doe@example.com",
            "subscription_status": "subscribed",
            "subscribed_at": "2023-01-01T00:00:00Z",
            "source": "sync",
            "note": "VIP customer"
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
      "operationId": "reach_listContactsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "uuid": "550e8400-e29b-41d4-a716-446655440001",
            "name": "Jane",
            "surname": "Doe",
            "email": "jane.doe@example.com",
            "subscription_status": "subscribed",
            "subscribed_at": "2023-01-01T00:00:00Z",
            "source": "sync",
            "note": "VIP customer"
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
      "operationId": "reach_listContactsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "uuid": "550e8400-e29b-41d4-a716-446655440000",
            "name": "John",
            "surname": "Doe",
            "email": "john.doe@example.com",
            "subscription_status": "subscribed",
            "subscribed_at": "2023-01-01T00:00:00Z",
            "source": "sync",
            "note": "VIP customer"
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
      "operationId": "reach_listContactsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "uuid": "550e8400-e29b-41d4-a716-446655440001",
            "name": "Jane",
            "surname": "Doe",
            "email": "jane.doe@example.com",
            "subscription_status": "subscribed",
            "subscribed_at": "2023-01-01T00:00:00Z",
            "source": "sync",
            "note": "VIP customer"
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

const testAccReachContactsDataSourceConfig = `
data "hostinger_reach_contacts" "test" {
	group_uuid = "550e8400-e29b-41d4-a716-446655440000"
	subscription_status = "subscribed"
}
`
