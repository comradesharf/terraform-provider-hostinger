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

func TestAccDataSourceVPSPublicKeys(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVPSPublicKeysConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_public_keys.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("public key"),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_public_keys.test",
						tfjsonpath.New("public_keys").AtSliceIndex(0).AtMapKey("id"),
						knownvalue.Int64Exact(325),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_public_keys.test",
						tfjsonpath.New("public_keys").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("My public key"),
					),
					statecheck.ExpectKnownValue(
						"data.hostinger_vps_public_keys.test",
						tfjsonpath.New("public_keys").AtSliceIndex(0).AtMapKey("key"),
						knownvalue.StringExact("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD..."),
					),
				},
			},
		},
	})
}

const testAccDataSourceVPSPublicKeysConfig = `
data "hostinger_vps_public_keys" "test" {
	name = "public key"
}
`
