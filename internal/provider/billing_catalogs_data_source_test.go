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

func TestAccBillingCatalogsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccBillingCatalogsDataSourcePreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBillingCatalogsDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.hostinger_billing_catalogs.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("KVM 2"),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_billing_catalogs.test",
						tfjsonpath.New("category"),
						knownvalue.StringExact("VPS"),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_billing_catalogs.test",
						tfjsonpath.New("billing_catalogs").AtSliceIndex(0),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":       knownvalue.StringExact("hostingercom-vps-kvm2"),
							"name":     knownvalue.StringExact("KVM 2"),
							"category": knownvalue.StringExact("VPS"),
							"metadata": knownvalue.MapExact(map[string]knownvalue.Check{
								"field": knownvalue.StringExact("value"),
							}),
							"prices": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"id":                 knownvalue.StringExact("hostingercom-vps-kvm2-usd-1m"),
									"name":               knownvalue.StringExact("KVM 2 (billed every month)"),
									"currency":           knownvalue.StringExact("USD"),
									"price":              knownvalue.Int32Exact(1799),
									"first_period_price": knownvalue.Int32Exact(899),
									"period":             knownvalue.Int32Exact(1),
									"period_unit":        knownvalue.StringExact("day"),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}

func testAccBillingCatalogsDataSourcePreCheck(t *testing.T) {
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
      "operationId": "billing_getCatalogItemListV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": [
        {
          "id": "hostingercom-vps-kvm2",
          "name": "KVM 2",
          "category": "VPS",
          "metadata": {
            "field": "value"
          },
          "prices": [
            {
              "id": "hostingercom-vps-kvm2-usd-1m",
              "name": "KVM 2 (billed every month)",
              "currency": "USD",
              "price": 1799,
              "first_period_price": 899,
              "period": 1,
              "period_unit": "day"
            }
          ]
        }
      ]
    },
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

const testAccBillingCatalogsDataSourceConfig = `
data "hostinger_billing_catalogs" "test" {
	name = "KVM 2"
	category = "VPS"
}
`
