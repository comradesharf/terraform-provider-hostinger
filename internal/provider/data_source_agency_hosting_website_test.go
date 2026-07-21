// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccDataSourceAgencyHostingWebsite(t *testing.T) {
	websiteUID := os.Getenv("HOSTINGER_AGENCY_WEBSITE_UID")
	if websiteUID == "" {
		t.Skip("HOSTINGER_AGENCY_WEBSITE_UID not set, skipping acceptance test")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAgencyHostingWebsiteConfig(websiteUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.hostinger_agency_hosting_website.test",
						tfjsonpath.New("website_uid"),
						knownvalue.StringExact(websiteUID),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_agency_hosting_website.test",
						tfjsonpath.New("uid"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func testAccAgencyHostingWebsiteConfig(websiteUID string) string {
	return `
data "hostinger_agency_hosting_website" "test" {
  website_uid = "` + websiteUID + `"
}
`
}
