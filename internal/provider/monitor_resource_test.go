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

func TestAccMonitorResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMonitorResourceConfig("HTTP Monitor", "http", "https://example.com"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uptimekuma_monitor.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("HTTP Monitor"),
					),
					statecheck.ExpectKnownValue(
						"uptimekuma_monitor.test",
						tfjsonpath.New("type"),
						knownvalue.StringExact("http"),
					),
					statecheck.ExpectKnownValue(
						"uptimekuma_monitor.test",
						tfjsonpath.New("url"),
						knownvalue.StringExact("https://example.com"),
					),
				},
			},
			// ImportState testing
			{
				ResourceName:      "uptimekuma_monitor.test",
				ImportState:       true,
				ImportStateVerify: true,
				// Certain fields may not be returned by the API and should be excluded from import verification
				ImportStateVerifyIgnore: []string{"basic_auth_pass"},
			},
			// Update and Read testing
			{
				Config: testAccMonitorResourceConfig("Updated HTTP Monitor", "http", "https://updated-example.com"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uptimekuma_monitor.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("Updated HTTP Monitor"),
					),
					statecheck.ExpectKnownValue(
						"uptimekuma_monitor.test",
						tfjsonpath.New("url"),
						knownvalue.StringExact("https://updated-example.com"),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMonitorResourceConfig(name, monitorType, url string) string {
	return fmt.Sprintf(`
provider "uptimekuma" {
  base_url = "%s"
  username = "%s"
  password = "%s"
}

resource "uptimekuma_monitor" "test" {
  name     = %[4]q
  type     = %[5]q
  url      = %[6]q
  interval = 60
  max_retries = 3
  retry_interval = 30
}
`,
		os.Getenv("UPTIMEKUMA_BASE_URL"),
		os.Getenv("UPTIMEKUMA_USERNAME"),
		os.Getenv("UPTIMEKUMA_PASSWORD"),
		name, monitorType, url)
}

func TestAccPingMonitorResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccPingMonitorResourceConfig("Ping Monitor", "example.com"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uptimekuma_monitor.ping_test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("Ping Monitor"),
					),
					statecheck.ExpectKnownValue(
						"uptimekuma_monitor.ping_test",
						tfjsonpath.New("type"),
						knownvalue.StringExact("ping"),
					),
					statecheck.ExpectKnownValue(
						"uptimekuma_monitor.ping_test",
						tfjsonpath.New("hostname"),
						knownvalue.StringExact("example.com"),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccPingMonitorResourceConfig(name, hostname string) string {
	return fmt.Sprintf(`
provider "uptimekuma" {
  base_url = "%s"
  username = "%s"
  password = "%s"
}

resource "uptimekuma_monitor" "ping_test" {
  name     = %[4]q
  type     = "ping"
  hostname = %[5]q
  interval = 60
  max_retries = 3
}
`,
		os.Getenv("UPTIMEKUMA_BASE_URL"),
		os.Getenv("UPTIMEKUMA_USERNAME"),
		os.Getenv("UPTIMEKUMA_PASSWORD"),
		name, hostname)
}
