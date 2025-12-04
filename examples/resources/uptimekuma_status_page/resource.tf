# Copyright (c) eHealth.co.id as PT Aksara Digital Indonesia
# SPDX-License-Identifier: MPL-2.0

resource "uptimekuma_status_page" "company_status" {
  slug        = "status"
  title       = "Company Status Page"
  description = "Current status of our services"
  theme       = "dark"
  published   = true
  show_tags   = false

  # Optional: Custom domains for your status page
  domain_name_list = ["status.example.com"]

  # Optional: Custom footer and styling
  footer_text = "© 2025 Example Company"
  custom_css  = ".status-page-header { background-color: #2a3b4c; }"

  # Optional: Google Analytics tracking
  google_analytics_id = "G-EXAMPLE123"

  # Optional: Custom icon
  icon = "https://example.com/logo.png"

  # Optional: Show "Powered by Uptime Kuma" text
  show_powered_by = true

  # Group monitors on the status page
  public_group_list = [
    {
      name   = "Core Services"
      weight = 1
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