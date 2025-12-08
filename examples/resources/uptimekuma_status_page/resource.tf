# Copyright (c) eHealth.co.id as PT Aksara Digital Indonesia
# SPDX-License-Identifier: MPL-2.0

resource "uptimekuma_status_page" "company_status" {
  # Slug: URL-friendly identifier for the status page (string, required)
  # Used in URL: https://your-uptime-kuma.com/status/{slug}
  # Must be unique, lowercase, alphanumeric with hyphens
  slug = "status"

  # Title: Display title on the status page (string, required)
  title = "Company Status Page"

  # Description: Subtitle or description text (string, optional)
  description = "Current status of our services"

  # Theme: Visual theme for the page (string, optional)
  # Valid values: "light", "dark", "auto"
  theme = "dark"

  # Published: Make the status page publicly accessible (boolean, default: true)
  published = true

  # Show Tags: Display monitor tags on the page (boolean, default: false)
  show_tags = false

  # Domain Name List: Custom domains for the status page (list of strings, optional)
  # Must configure DNS to point to your Uptime Kuma instance
  # Example: ["status.example.com", "uptime.example.com"]
  domain_name_list = ["status.example.com"]

  # Footer Text: Custom footer content with HTML support (string, optional)
  footer_text = "© 2025 Example Company"

  # Custom CSS: Additional CSS styles (string, optional)
  # Can override default styles for complete customization
  custom_css = ".status-page-header { background-color: #2a3b4c; }"

  # Google Analytics ID: Track page views (string, optional)
  # Format: G-XXXXXXXXXX (GA4) or UA-XXXXXXXXX-X (Universal Analytics)
  google_analytics_id = "G-EXAMPLE123"

  # Icon: Custom icon/logo URL (string, default: "/icon.svg")
  # Absolute URL or path relative to Uptime Kuma root
  icon = "https://example.com/logo.png"

  # Show Powered By: Display "Powered by Uptime Kuma" footer (boolean, default: true)
  show_powered_by = true

  # Public Group List: Organize monitors into groups (list of objects, optional)
  # Each group creates a section on the status page
  public_group_list = [
    {
      # Name: Group display name (string, required)
      name = "Core Services"

      # Weight: Display order - lower numbers appear first (number, optional, default: 1)
      weight = 1

      # Monitor List: IDs of monitors to include (list of numbers, optional)
      # Get monitor IDs from uptimekuma_monitor resources
      monitor_list = [
        uptimekuma_monitor.http_example.id
      ]
    },
    {
      name   = "API Services"
      weight = 2
      monitor_list = [
        uptimekuma_monitor.authenticated_http.id
      ]
    },
    {
      name   = "Infrastructure"
      weight = 3
      monitor_list = [
        uptimekuma_monitor.ping_example.id,
        uptimekuma_monitor.port_example.id
      ]
    }
  ]
}

# Minimal Status Page Example
resource "uptimekuma_status_page" "minimal" {
  # Only required fields
  slug  = "minimal-status"
  title = "System Status"

  # Single group with monitors
  public_group_list = [
    {
      name         = "All Services"
      monitor_list = [uptimekuma_monitor.http_example.id]
    }
  ]
}