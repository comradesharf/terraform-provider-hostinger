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

func TestAccVPSFirewallsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVPSFirewallsDataSourcePreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVPSFirewallsDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_firewalls.test",
						tfjsonpath.New("firewalls"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id":        knownvalue.Int64Exact(65224),
								"name":      knownvalue.StringExact("HTTP and SSH only"),
								"is_synced": knownvalue.Bool(false),
								"rules": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"id":            knownvalue.Int64Exact(24541),
										"action":        knownvalue.StringExact("accept"),
										"protocol":      knownvalue.StringExact("TCP"),
										"port":          knownvalue.StringExact("1024:2048"),
										"source":        knownvalue.StringExact("any"),
										"source_detail": knownvalue.StringExact("any"),
									}),
								}),
								"created_at": knownvalue.StringExact("2021-09-01T12:00:00Z"),
								"updated_at": knownvalue.StringExact("2021-09-01T12:00:00Z"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id":        knownvalue.Int64Exact(65225),
								"name":      knownvalue.StringExact("HTTP and SSH only"),
								"is_synced": knownvalue.Bool(false),
								"rules": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"id":            knownvalue.Int64Exact(24542),
										"action":        knownvalue.StringExact("accept"),
										"protocol":      knownvalue.StringExact("TCP"),
										"port":          knownvalue.StringExact("1024:2048"),
										"source":        knownvalue.StringExact("any"),
										"source_detail": knownvalue.StringExact("any"),
									}),
								}),
								"created_at": knownvalue.StringExact("2021-09-01T12:00:00Z"),
								"updated_at": knownvalue.StringExact("2021-09-01T12:00:00Z"),
							}),
						}),
					),
				},
			},
		},
	})
}

const testAccVPSFirewallsDataSourceConfig = `
data "hostinger_vps_firewalls" "test" {}
`

func testAccVPSFirewallsDataSourcePreCheck(t *testing.T) {
	client := newMockServerClient()

	if err := client.Clear(nil, mockserver.ClearAll); err != nil {
		t.Fatalf("failed to clear mock server: %v", err)
	}

	if err := client.FreezeClock("2021-09-01T12:00:00Z"); err != nil {
		t.Fatalf("failed to freeze mock server clock: %v", err)
	}

	// language=json
	body := []byte(`
[]
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
