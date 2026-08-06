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

func TestAccAgencyHostingDomainsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccAgencyHostingDomainsDataSourcePreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAgencyHostingDomainsDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.hostinger_agency_hosting_domains.test",
						tfjsonpath.New("website_uids"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("535bb70f-b4bf-4250-a581-f0c8e882b1a2"),
							knownvalue.StringExact("e3f1c5d2-4b6a-4f8e-9c3b-1a2b3c4d5e6f"),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_agency_hosting_domains.test",
						tfjsonpath.New("domains"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"fqdn":        knownvalue.StringExact("example1.com"),
								"website_uid": knownvalue.StringExact("535bb70f-b4bf-4250-a581-f0c8e882b1a2"),
								"created_at":  knownvalue.StringExact("2024-05-29T05:49:49Z"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"fqdn":        knownvalue.StringExact("example2.com"),
								"website_uid": knownvalue.StringExact("e3f1c5d2-4b6a-4f8e-9c3b-1a2b3c4d5e6f"),
								"created_at":  knownvalue.StringExact("2024-05-29T05:49:49Z"),
							}),
						}),
					),
				},
			},
		},
	})
}

func testAccAgencyHostingDomainsDataSourcePreCheck(t *testing.T) {
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
      "operationId": "agency-hosting_listAgencyPlanDomainsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "fqdn": "example1.com",
            "website_uid": "535bb70f-b4bf-4250-a581-f0c8e882b1a2",
            "created_at": "2024-05-29T05:49:49Z"
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
      "operationId": "agency-hosting_listAgencyPlanDomainsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "fqdn": "example2.com",
            "website_uid": "e3f1c5d2-4b6a-4f8e-9c3b-1a2b3c4d5e6f",
            "created_at": "2024-05-29T05:49:49Z"
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
      "operationId": "agency-hosting_listAgencyPlanDomainsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "fqdn": "example1.com",
            "website_uid": "535bb70f-b4bf-4250-a581-f0c8e882b1a2",
            "created_at": "2024-05-29T05:49:49Z"
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
      "operationId": "agency-hosting_listAgencyPlanDomainsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "fqdn": "example2.com",
            "website_uid": "e3f1c5d2-4b6a-4f8e-9c3b-1a2b3c4d5e6f",
            "created_at": "2024-05-29T05:49:49Z"
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
      "operationId": "agency-hosting_listAgencyPlanDomainsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "fqdn": "example1.com",
            "website_uid": "535bb70f-b4bf-4250-a581-f0c8e882b1a2",
            "created_at": "2024-05-29T05:49:49Z"
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
      "operationId": "agency-hosting_listAgencyPlanDomainsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "data": [
          {
            "fqdn": "example2.com",
            "website_uid": "e3f1c5d2-4b6a-4f8e-9c3b-1a2b3c4d5e6f",
            "created_at": "2024-05-29T05:49:49Z"
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

const testAccAgencyHostingDomainsDataSourceConfig = `
data "hostinger_agency_hosting_domains" "test" {
	website_uids = [
		"535bb70f-b4bf-4250-a581-f0c8e882b1a2",
		"e3f1c5d2-4b6a-4f8e-9c3b-1a2b3c4d5e6f"
	]
}
`
