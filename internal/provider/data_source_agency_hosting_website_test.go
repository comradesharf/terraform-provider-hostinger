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

func TestAccDataSourceAgencyHostingWebsite(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAgencyHostingWebsiteConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.hostinger_agency_hosting_website.test",
						tfjsonpath.New("uid"),
						knownvalue.StringExact("zpwlGlp19"),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_agency_hosting_website.test",
						tfjsonpath.New("ipv4"),
						knownvalue.StringExact("192.161.10.1"),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_agency_hosting_website.test",
						tfjsonpath.New("flavor"),
						knownvalue.StringExact("wp-6.2.0"),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_agency_hosting_website.test",
						tfjsonpath.New("type"),
						knownvalue.StringExact("node-static"),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_agency_hosting_website.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Very awesome website"),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_agency_hosting_website.test",
						tfjsonpath.New("state"),
						knownvalue.StringExact("active"),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_agency_hosting_website.test",
						tfjsonpath.New("created_at"),
						knownvalue.StringExact("2024-05-29T05:49:49Z"),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_agency_hosting_website.test",
						tfjsonpath.New("domains"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"fqdn":        knownvalue.StringExact("test.com"),
								"parent_fqdn": knownvalue.StringExact("test.com"),
								"ipv6":        knownvalue.StringExact("2001:db8::1"),
								"created_at":  knownvalue.StringExact("2024-05-29T05:49:49Z"),
								"nameservers": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.StringExact("a.dns-parking.com"),
									knownvalue.StringExact("b.dns-parking.com"),
								}),
								"ssl_cert": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"created_at": knownvalue.StringExact("2024-05-29T05:49:49Z"),
									"expires_at": knownvalue.StringExact("2024-05-29T05:49:49Z"),
									"names": knownvalue.ListExact([]knownvalue.Check{
										knownvalue.StringExact("test.com"),
										knownvalue.StringExact("www.test.com"),
									}),
								}),
								"custom_ssl_cert": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"is_expired": knownvalue.Bool(false),
									"created_at": knownvalue.StringExact("2024-05-29T05:49:49Z"),
									"expires_at": knownvalue.StringExact("2024-05-29T05:49:49Z"),
								}),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_agency_hosting_website.test",
						tfjsonpath.New("preview_domain"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"fqdn":       knownvalue.StringExact("plum-bee-184082.hostingersite.com"),
							"created_at": knownvalue.StringExact("2024-05-29T05:49:49Z"),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_agency_hosting_website.test",
						tfjsonpath.New("settings"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"php": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"version": knownvalue.StringExact("8.3"),
								"workers": knownvalue.Int64Exact(4),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_agency_hosting_website.test",
						tfjsonpath.New("wordpress"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"domain":           knownvalue.StringExact("test.com"),
							"title":            knownvalue.StringExact("My Blog"),
							"language":         knownvalue.StringExact("en_US"),
							"is_config_locked": knownvalue.Bool(true),
							"created_at":       knownvalue.StringExact("2024-05-29T05:49:49Z"),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_agency_hosting_website.test",
						tfjsonpath.New("remote_access"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"mode": knownvalue.StringExact("ssh_and_sftp"),
							"ssh": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"username":            knownvalue.StringExact("u123456789_abcDeFg"),
								"host":                knownvalue.StringExact("192.161.10.1"),
								"port":                knownvalue.Int32Exact(65002),
								"is_enabled":          knownvalue.Bool(true),
								"is_password_enabled": knownvalue.Bool(true),
							}),
							"sftp": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"username":   knownvalue.StringExact("u123456789_abcDeFg"),
								"host":       knownvalue.StringExact("192.161.10.1"),
								"port":       knownvalue.Int32Exact(65002),
								"is_enabled": knownvalue.Bool(true),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_agency_hosting_website.test",
						tfjsonpath.New("server"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"hostname":     knownvalue.StringExact("us-west-1.hstgr.io"),
							"country_code": knownvalue.StringExact("us"),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_agency_hosting_website.test",
						tfjsonpath.New("order"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":         knownvalue.Int64Exact(123456),
							"status":     knownvalue.StringExact("active"),
							"created_at": knownvalue.StringExact("2024-05-29T05:49:49Z"),
							"plan": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("Hosting Single"),
								"parameters": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"disk_quota_bytes":          knownvalue.Int64Exact(21474836480),
									"inode_quota":               knownvalue.Int64Exact(10000),
									"cpu_cores":                 knownvalue.Int64Exact(2),
									"memory_quota_bytes":        knownvalue.Int64Exact(1073741824),
									"disk_iops_quota":           knownvalue.Int64Exact(100000),
									"process_quota":             knownvalue.Int64Exact(10000),
									"website_quota":             knownvalue.Int64Exact(10),
									"max_databases_per_website": knownvalue.Int64Exact(5),
									"is_cdn_available":          knownvalue.Bool(true),
								}),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_agency_hosting_website.test",
						tfjsonpath.New("user"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"username": knownvalue.StringExact("u123456789"),
							"state":    knownvalue.StringExact("active"),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_agency_hosting_website.test",
						tfjsonpath.New("staging_root"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"uid": knownvalue.StringExact("zpwlGlp19"),
						}),
					),
				},
			},
		},
	})
}

const testAccDataSourceAgencyHostingWebsiteConfig = `
data "hostinger_agency_hosting_website" "test" {
	uid = "zpwlGlp19"
}
`
