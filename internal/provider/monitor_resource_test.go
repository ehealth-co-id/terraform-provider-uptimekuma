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

// New test for Keyword monitor type.
func TestAccKeywordMonitorResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKeywordMonitorResourceConfig("Keyword Monitor", "https://example.com", "Example Domain"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uptimekuma_monitor.keyword_test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("Keyword Monitor"),
					),
					statecheck.ExpectKnownValue(
						"uptimekuma_monitor.keyword_test",
						tfjsonpath.New("type"),
						knownvalue.StringExact("keyword"),
					),
					statecheck.ExpectKnownValue(
						"uptimekuma_monitor.keyword_test",
						tfjsonpath.New("keyword"),
						knownvalue.StringExact("Example Domain"),
					),
				},
			},
		},
	})
}

func testAccKeywordMonitorResourceConfig(name, url, keyword string) string {
	return fmt.Sprintf(`
provider "uptimekuma" {
  base_url = "%s"
  username = "%s"
  password = "%s"
}

resource "uptimekuma_monitor" "keyword_test" {
  name        = %[4]q
  type        = "keyword"
  url         = %[5]q
  method      = "GET"
  keyword     = %[6]q
  interval    = 60
  max_retries = 2
}
`,
		os.Getenv("UPTIMEKUMA_BASE_URL"),
		os.Getenv("UPTIMEKUMA_USERNAME"),
		os.Getenv("UPTIMEKUMA_PASSWORD"),
		name, url, keyword)
}

// New test for Port monitor type.
func TestAccPortMonitorResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPortMonitorResourceConfig("Port Monitor", "example.com", 443),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uptimekuma_monitor.port_test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("Port Monitor"),
					),
					statecheck.ExpectKnownValue(
						"uptimekuma_monitor.port_test",
						tfjsonpath.New("type"),
						knownvalue.StringExact("port"),
					),
					statecheck.ExpectKnownValue(
						"uptimekuma_monitor.port_test",
						tfjsonpath.New("port"),
						knownvalue.Int64Exact(443),
					),
				},
			},
		},
	})
}

func testAccPortMonitorResourceConfig(name, hostname string, port int) string {
	return fmt.Sprintf(`
provider "uptimekuma" {
  base_url = "%s"
  username = "%s"
  password = "%s"
}

resource "uptimekuma_monitor" "port_test" {
  name     = %[4]q
  type     = "port"
  hostname = %[5]q
  port     = %[6]d
  interval = 60
}
`,
		os.Getenv("UPTIMEKUMA_BASE_URL"),
		os.Getenv("UPTIMEKUMA_USERNAME"),
		os.Getenv("UPTIMEKUMA_PASSWORD"),
		name, hostname, port)
}

// New test for HTTP monitor with custom headers and status codes.
func TestAccHTTPMonitorWithHeaders(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccHTTPMonitorWithHeadersConfig("API Monitor", "https://api.example.com/health"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uptimekuma_monitor.api_monitor",
						tfjsonpath.New("name"),
						knownvalue.StringExact("API Monitor"),
					),
					statecheck.ExpectKnownValue(
						"uptimekuma_monitor.api_monitor",
						tfjsonpath.New("method"),
						knownvalue.StringExact("GET"),
					),
				},
			},
		},
	})
}

func testAccHTTPMonitorWithHeadersConfig(name, url string) string {
	return fmt.Sprintf(`
provider "uptimekuma" {
  base_url = "%s"
  username = "%s"
  password = "%s"
}

resource "uptimekuma_monitor" "api_monitor" {
  name     = %[4]q
  type     = "http"
  url      = %[5]q
  method   = "GET"
  headers  = "{\"Accept\": \"application/json\"}"
  accepted_status_codes = [200, 201]
  interval = 60
  ignore_tls = false
}
`,
		os.Getenv("UPTIMEKUMA_BASE_URL"),
		os.Getenv("UPTIMEKUMA_USERNAME"),
		os.Getenv("UPTIMEKUMA_PASSWORD"),
		name, url)
}

// New test for interval and timing field updates.
func TestAccMonitorIntervalUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorIntervalUpdateConfig("Interval Monitor", 60, 30),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uptimekuma_monitor.interval_test",
						tfjsonpath.New("interval"),
						knownvalue.Int64Exact(60),
					),
					statecheck.ExpectKnownValue(
						"uptimekuma_monitor.interval_test",
						tfjsonpath.New("retry_interval"),
						knownvalue.Int64Exact(30),
					),
				},
			},
			// Update intervals
			{
				Config: testAccMonitorIntervalUpdateConfig("Interval Monitor", 120, 60),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uptimekuma_monitor.interval_test",
						tfjsonpath.New("interval"),
						knownvalue.Int64Exact(120),
					),
					statecheck.ExpectKnownValue(
						"uptimekuma_monitor.interval_test",
						tfjsonpath.New("retry_interval"),
						knownvalue.Int64Exact(60),
					),
				},
			},
		},
	})
}

func testAccMonitorIntervalUpdateConfig(name string, interval, retryInterval int) string {
	return fmt.Sprintf(`
provider "uptimekuma" {
  base_url = "%s"
  username = "%s"
  password = "%s"
}

resource "uptimekuma_monitor" "interval_test" {
  name           = %[4]q
  type           = "http"
  url            = "https://example.com"
  interval       = %[5]d
  retry_interval = %[6]d
  max_retries    = 3
}
`,
		os.Getenv("UPTIMEKUMA_BASE_URL"),
		os.Getenv("UPTIMEKUMA_USERNAME"),
		os.Getenv("UPTIMEKUMA_PASSWORD"),
		name, interval, retryInterval)
}

// New test for upside down (status inversion).
func TestAccMonitorUpsideDown(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorUpsideDownConfig("Inverted Monitor", true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uptimekuma_monitor.upside_down_test",
						tfjsonpath.New("upside_down"),
						knownvalue.Bool(true),
					),
				},
			},
			// Update to false
			{
				Config: testAccMonitorUpsideDownConfig("Inverted Monitor", false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uptimekuma_monitor.upside_down_test",
						tfjsonpath.New("upside_down"),
						knownvalue.Bool(false),
					),
				},
			},
		},
	})
}

func testAccMonitorUpsideDownConfig(name string, upsideDown bool) string {
	return fmt.Sprintf(`
provider "uptimekuma" {
  base_url = "%s"
  username = "%s"
  password = "%s"
}

resource "uptimekuma_monitor" "upside_down_test" {
  name        = %[4]q
  type        = "http"
  url         = "https://example.com"
  upside_down = %[5]t
  interval    = 60
}
`,
		os.Getenv("UPTIMEKUMA_BASE_URL"),
		os.Getenv("UPTIMEKUMA_USERNAME"),
		os.Getenv("UPTIMEKUMA_PASSWORD"),
		name, upsideDown)
}
