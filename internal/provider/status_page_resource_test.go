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

func TestAccStatusPageResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccStatusPageResourceConfig("test-page", "Test Status Page", "System status page for testing"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uptimekuma_status_page.test",
						tfjsonpath.New("slug"),
						knownvalue.StringExact("test-page"),
					),
					statecheck.ExpectKnownValue(
						"uptimekuma_status_page.test",
						tfjsonpath.New("title"),
						knownvalue.StringExact("Test Status Page"),
					),
					statecheck.ExpectKnownValue(
						"uptimekuma_status_page.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("System status page for testing"),
					),
					statecheck.ExpectKnownValue(
						"uptimekuma_status_page.test",
						tfjsonpath.New("published"),
						knownvalue.Bool(true),
					),
				},
			},
			// ImportState testing - import by slug
			{
				ResourceName:            "uptimekuma_status_page.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateId:           "test-page",
				ImportStateVerifyIgnore: []string{"id"}, // ID might differ after import
			},
			// Update and Read testing
			{
				Config: testAccStatusPageResourceConfig("test-page", "Updated Status Page", "Updated description"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uptimekuma_status_page.test",
						tfjsonpath.New("slug"),
						knownvalue.StringExact("test-page"),
					),
					statecheck.ExpectKnownValue(
						"uptimekuma_status_page.test",
						tfjsonpath.New("title"),
						knownvalue.StringExact("Updated Status Page"),
					),
					statecheck.ExpectKnownValue(
						"uptimekuma_status_page.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Updated description"),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccStatusPageResourceConfig(slug, title, description string) string {
	return fmt.Sprintf(`
provider "uptimekuma" {
  base_url = "%s"
  username = "%s"
  password = "%s"
}

resource "uptimekuma_status_page" "test" {
  slug        = %[4]q
  title       = %[5]q
  description = %[6]q
  published   = true
  theme       = "dark"
  show_tags   = false
}
`,
		os.Getenv("UPTIMEKUMA_BASE_URL"),
		os.Getenv("UPTIMEKUMA_USERNAME"),
		os.Getenv("UPTIMEKUMA_PASSWORD"),
		slug, title, description)
}

func TestAccStatusPageResourceWithGroups(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing with monitor groups
			{
				Config: testAccStatusPageResourceWithGroupsConfig("test-page-with-groups", "Status Page With Groups"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uptimekuma_status_page.with_groups",
						tfjsonpath.New("slug"),
						knownvalue.StringExact("test-page-with-groups"),
					),
					statecheck.ExpectKnownValue(
						"uptimekuma_status_page.with_groups",
						tfjsonpath.New("title"),
						knownvalue.StringExact("Status Page With Groups"),
					),
					// We can't easily test the exact structure of nested blocks like public_group_list
					// in this basic test, but we can verify the status page itself is created
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccStatusPageResourceWithGroupsConfig(slug, title string) string {
	return fmt.Sprintf(`
provider "uptimekuma" {
  base_url = "%s"
  username = "%s"
  password = "%s"
}

resource "uptimekuma_monitor" "http1" {
  name     = "HTTP Monitor 1"
  type     = "http"
  url      = "https://example.com"
  interval = 60
}

resource "uptimekuma_monitor" "http2" {
  name     = "HTTP Monitor 2" 
  type     = "http"
  url      = "https://example.org"
  interval = 60
}

resource "uptimekuma_status_page" "with_groups" {
  slug      = %[4]q
  title     = %[5]q
  published = true
  theme     = "dark"
  
  public_group_list = [
    {
      name = "Core Services"
      weight = 1
      monitor_list = [
        uptimekuma_monitor.http1.id
      ]
    },
    {
      name = "Secondary Services"
      weight = 2
      monitor_list = [
        uptimekuma_monitor.http2.id
      ]
    }
  ]
}
`,
		os.Getenv("UPTIMEKUMA_BASE_URL"),
		os.Getenv("UPTIMEKUMA_USERNAME"),
		os.Getenv("UPTIMEKUMA_PASSWORD"),
		slug, title)
}

// New test for status page theme and customization.
func TestAccStatusPageThemeUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStatusPageThemeUpdateConfig("theme-test", "Theme Test", "dark"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uptimekuma_status_page.theme_test",
						tfjsonpath.New("theme"),
						knownvalue.StringExact("dark"),
					),
				},
			},
			// Update theme
			{
				Config: testAccStatusPageThemeUpdateConfig("theme-test", "Theme Test", "light"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uptimekuma_status_page.theme_test",
						tfjsonpath.New("theme"),
						knownvalue.StringExact("light"),
					),
				},
			},
		},
	})
}

func testAccStatusPageThemeUpdateConfig(slug, title, theme string) string {
	return fmt.Sprintf(`
provider "uptimekuma" {
  base_url = "%s"
  username = "%s"
  password = "%s"
}

resource "uptimekuma_status_page" "theme_test" {
  slug  = %[4]q
  title = %[5]q
  theme = %[6]q
}
`,
		os.Getenv("UPTIMEKUMA_BASE_URL"),
		os.Getenv("UPTIMEKUMA_USERNAME"),
		os.Getenv("UPTIMEKUMA_PASSWORD"),
		slug, title, theme)
}

// New test for custom CSS and footer.
func TestAccStatusPageCustomization(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStatusPageCustomizationConfig("custom-test", "Custom Test"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uptimekuma_status_page.custom_test",
						tfjsonpath.New("footer_text"),
						knownvalue.StringExact("© 2025 Test Company"),
					),
				},
			},
		},
	})
}

func testAccStatusPageCustomizationConfig(slug, title string) string {
	return fmt.Sprintf(`
provider "uptimekuma" {
  base_url = "%s"
  username = "%s"
  password = "%s"
}

resource "uptimekuma_status_page" "custom_test" {
  slug        = %[4]q
  title       = %[5]q
  footer_text = "© 2025 Test Company"
  custom_css  = ".header { background: #000; }"
  show_tags   = true
}
`,
		os.Getenv("UPTIMEKUMA_BASE_URL"),
		os.Getenv("UPTIMEKUMA_USERNAME"),
		os.Getenv("UPTIMEKUMA_PASSWORD"),
		slug, title)
}
