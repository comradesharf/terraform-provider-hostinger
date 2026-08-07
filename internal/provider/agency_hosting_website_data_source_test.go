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

func TestAccAgencyHostingWebsiteDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccAgencyHostingWebsiteDataSourcePreCheck(t)
		},
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

func testAccAgencyHostingWebsiteDataSourcePreCheck(t *testing.T) {
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
      "operationId": "agency-hosting_getAgencyPlanWebsiteDetailsV1"
    },
    "httpResponse": {
      "statusCode": 200,
      "body": {
        "uid": "zpwlGlp19",
        "ipv4": "192.161.10.1",
        "flavor": "wp-6.2.0",
        "type": "node-static",
        "description": "Very awesome website",
        "state": "active",
        "created_at": "2024-05-29T05:49:49Z",
        "domains": [
          {
            "fqdn": "test.com",
            "parent_fqdn": "test.com",
            "ipv6": "2001:db8::1",
            "created_at": "2024-05-29T05:49:49Z",
            "nameservers": [
              "a.dns-parking.com",
              "b.dns-parking.com"
            ],
            "ssl_cert": {
              "created_at": "2024-05-29T05:49:49Z",
              "expires_at": "2024-05-29T05:49:49Z",
              "names": [
                "test.com",
                "www.test.com"
              ]
            },
            "custom_ssl_cert": {
              "is_expired": false,
              "created_at": "2024-05-29T05:49:49Z",
              "expires_at": "2024-05-29T05:49:49Z"
            }
          }
        ],
        "preview_domain": {
          "fqdn": "plum-bee-184082.hostingersite.com",
          "created_at": "2024-05-29T05:49:49Z"
        },
        "settings": {
          "php": {
            "version": "8.3",
            "workers": 4
          }
        },
        "wordpress": {
          "domain": "test.com",
          "title": "My Blog",
          "language": "en_US",
          "is_config_locked": true,
          "created_at": "2024-05-29T05:49:49Z"
        },
        "remote_access": {
          "mode": "ssh_and_sftp",
          "ssh": {
            "username": "u123456789_abcDeFg",
            "host": "192.161.10.1",
            "port": 65002,
            "is_enabled": true,
            "is_password_enabled": true
          },
          "sftp": {
            "username": "u123456789_abcDeFg",
            "host": "192.161.10.1",
            "port": 65002,
            "is_enabled": true
          }
        },
        "server": {
          "hostname": "us-west-1.hstgr.io",
          "country_code": "us"
        },
        "order": {
          "id": 123456,
          "status": "active",
          "created_at": "2024-05-29T05:49:49Z",
          "plan": {
            "name": "Hosting Single",
            "parameters": {
              "disk_quota_bytes": 21474836480,
              "inode_quota": 10000,
              "cpu_cores": 2,
              "memory_quota_bytes": 1073741824,
              "disk_iops_quota": 100000,
              "process_quota": 10000,
              "website_quota": 10,
              "max_databases_per_website": 5,
              "is_cdn_available": true
            }
          }
        },
        "user": {
          "username": "u123456789",
          "state": "active"
        },
        "staging_root": {
          "uid": "zpwlGlp19"
        }
      }
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

const testAccDataSourceAgencyHostingWebsiteConfig = `
data "hostinger_agency_hosting_website" "test" {
	uid = "zpwlGlp19"
}
`
