// Copyright (c) eHealth.co.id as PT Aksara Digital Indonesia
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccTagResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTagResourceConfig("production", "#00FF00"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uptimekuma_tag.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("production"),
					),
					statecheck.ExpectKnownValue(
						"uptimekuma_tag.test",
						tfjsonpath.New("color"),
						knownvalue.StringExact("#00FF00"),
					),
				},
			},
			// ImportState testing
			{
				ResourceName:      "uptimekuma_tag.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccTagResourceConfig("production", "#FF0000"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uptimekuma_tag.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("production"),
					),
					statecheck.ExpectKnownValue(
						"uptimekuma_tag.test",
						tfjsonpath.New("color"),
						knownvalue.StringExact("#FF0000"),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTagResourceNameUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with initial name
			{
				Config: testAccTagResourceConfig("staging", "#FFA500"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uptimekuma_tag.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("staging"),
					),
				},
			},
			// Update name
			{
				Config: testAccTagResourceConfig("development", "#FFA500"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uptimekuma_tag.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("development"),
					),
					statecheck.ExpectKnownValue(
						"uptimekuma_tag.test",
						tfjsonpath.New("color"),
						knownvalue.StringExact("#FFA500"),
					),
				},
			},
		},
	})
}

func TestAccTagResourceWithoutColor(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create tag without color (optional field)
			{
				Config: testAccTagResourceConfigWithoutColor("critical"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uptimekuma_tag.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("critical"),
					),
				},
			},
		},
	})
}

func testAccTagResourceConfig(name string, color string) string {
	return fmt.Sprintf(`
provider "uptimekuma" {
  base_url = %[3]q
  username = %[4]q
  password = %[5]q
}

resource "uptimekuma_tag" "test" {
  name  = %[1]q
  color = %[2]q
}
`, name, color,
		testAccGetEnv("UPTIMEKUMA_BASE_URL", "http://localhost:3001"),
		testAccGetEnv("UPTIMEKUMA_USERNAME", "admin"),
		testAccGetEnv("UPTIMEKUMA_PASSWORD", "admin123"))
}

func testAccTagResourceConfigWithoutColor(name string) string {
	return fmt.Sprintf(`
provider "uptimekuma" {
  base_url = %[2]q
  username = %[3]q
  password = %[4]q
}

resource "uptimekuma_tag" "test" {
  name = %[1]q
}
`, name,
		testAccGetEnv("UPTIMEKUMA_BASE_URL", "http://localhost:3001"),
		testAccGetEnv("UPTIMEKUMA_USERNAME", "admin"),
		testAccGetEnv("UPTIMEKUMA_PASSWORD", "admin123"))
}

func testAccGetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
