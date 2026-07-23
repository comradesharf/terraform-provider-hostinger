// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccDataSourceAgencyHostingDomains(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAgencyHostingDomainsConfig,
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
								"fqdn":        knownvalue.StringExact("example.com"),
								"website_uid": knownvalue.StringExact("zpwlGlp19"),
								"created_at":  knownvalue.StringExact("2024-05-29T05:49:49Z"),
							}),
						}),
					),
				},
			},
		},
	})
}

const testAccDataSourceAgencyHostingDomainsConfig = `
data "hostinger_agency_hosting_domains" "test" {
	website_uids = [
		"535bb70f-b4bf-4250-a581-f0c8e882b1a2",
		"e3f1c5d2-4b6a-4f8e-9c3b-1a2b3c4d5e6f"
	]
}
`
