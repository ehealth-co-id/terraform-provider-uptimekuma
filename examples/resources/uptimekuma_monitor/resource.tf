# Copyright (c) eHealth.co.id as PT Aksara Digital Indonesia
# SPDX-License-Identifier: MPL-2.0

# HTTP Monitor Example
resource "uptimekuma_monitor" "http_example" {
  # Name: Display name for the monitor (string, required)
  name = "Example Website"

  # Type: Monitor type - "http", "ping", "port", "keyword" (string, required)
  type = "http"

  # URL: Target URL to monitor (string, required for http/keyword monitors)
  url = "https://example.com"

  # Method: HTTP method to use (string, default: "GET")
  # Valid values: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS
  method = "GET"

  # Interval: Check interval in seconds (number, default: 60)
  interval = 60

  # Retry Interval: Wait time before retry in seconds (number, default: 60)
  retry_interval = 30

  # Max Retries: Maximum retry attempts before marking as down (number, default: 0)
  max_retries = 3

  # Upside Down: Invert status - DOWN becomes UP (boolean, default: false)
  upside_down = false

  # Ignore TLS: Skip TLS/SSL certificate validation (boolean, default: false)
  ignore_tls = false
}

# Ping Monitor Example
resource "uptimekuma_monitor" "ping_example" {
  name = "Ping Example"
  type = "ping"

  # Hostname: Target hostname or IP address (string, required for ping/port monitors)
  hostname = "example.com"

  interval       = 120
  retry_interval = 30
  max_retries    = 2
}

# Keyword Monitor Example
resource "uptimekuma_monitor" "keyword_example" {
  name = "Keyword Search Example"
  type = "keyword"
  url  = "https://example.com"

  # Method: HTTP method for keyword monitor (string, default: "GET")
  method = "GET"

  interval = 300

  # Keyword: Text to search for in HTTP response (string, required for keyword monitors)
  keyword = "Example Domain"

  # Upside Down: Alert when keyword IS found (true) vs when NOT found (false)
  upside_down = true

  # Max Redirects: Maximum HTTP redirects to follow (number, default: 0)
  max_redirects = 3
}

# Port Monitor Example
resource "uptimekuma_monitor" "port_example" {
  name     = "Port Example"
  type     = "port"
  hostname = "example.com"

  # Port: TCP port number to check (number, required for port monitors)
  # Range: 1-65535
  port = 443

  interval       = 60
  retry_interval = 30
  max_retries    = 1
}

# Authenticated HTTP Monitor Example
resource "uptimekuma_monitor" "authenticated_http" {
  name           = "Authenticated API"
  type           = "http"
  url            = "https://api.example.com/private"
  method         = "GET"
  interval       = 60
  retry_interval = 30
  max_retries    = 3

  # Auth Method: Authentication type (string, optional)
  # Valid values: "basic", "ntlm", "mtls"
  auth_method = "basic"

  # Basic Auth User: Username for basic authentication (string, optional)
  basic_auth_user = "apiuser"

  # Basic Auth Pass: Password for basic authentication (string, sensitive, optional)
  basic_auth_pass = "securepassword"

  # Headers: Custom HTTP headers as JSON string (string, optional)
  # Format: {"Header-Name": "value", "Another-Header": "value"}
  headers = "{\"X-API-Key\":\"myapikey\", \"Accept\":\"application/json\"}"

  # Accepted Status Codes: HTTP codes considered successful (list of numbers, optional)
  # Default: All 2xx codes. Example: [200, 201, 204]
  accepted_status_codes = [200, 201]
}

# Advanced HTTP Monitor with Body
resource "uptimekuma_monitor" "http_with_body" {
  name   = "API POST Monitor"
  type   = "http"
  url    = "https://api.example.com/webhook"
  method = "POST"

  # Body: Request body for POST/PUT/PATCH requests (string, optional)
  # Can be JSON, XML, or any text format
  body = "{\"event\":\"health_check\",\"source\":\"terraform\"}"

  # Headers: Set Content-Type for the body
  headers = "{\"Content-Type\":\"application/json\"}"

  interval       = 300
  retry_interval = 60
  max_retries    = 2

  # Notification ID List: Alert notification IDs to trigger (list of numbers, optional)
  # Get IDs from Uptime Kuma notification settings
  # notification_id_list = [1, 2, 3]
}

# Monitor with Tags Example
# Tags help organize and categorize monitors for better management
resource "uptimekuma_tag" "production" {
  name  = "production"
  color = "#00FF00"
}

resource "uptimekuma_tag" "critical" {
  name  = "critical"
  color = "#FF0000"
}

resource "uptimekuma_monitor" "tagged_monitor" {
  name           = "Production API"
  type           = "http"
  url            = "https://api.example.com/health"
  interval       = 30
  retry_interval = 10
  max_retries    = 3

  # Tags: Associate tags with the monitor (list of objects, optional)
  # Each tag requires tag_id (required) and optional value
  tags = [
    {
      # tag_id: Reference to the tag resource ID (number, required)
      tag_id = uptimekuma_tag.production.id
    },
    {
      tag_id = uptimekuma_tag.critical.id
      # value: Optional custom value for this tag association (string, optional)
      # Useful for adding context like environment name, priority level, etc.
      value = "high-priority"
    }
  ]
}

# Monitor with Multiple Tags and Values
resource "uptimekuma_monitor" "multi_tagged_monitor" {
  name           = "E-commerce Website"
  type           = "http"
  url            = "https://shop.example.com"
  interval       = 60
  retry_interval = 30
  max_retries    = 3

  tags = [
    {
      tag_id = uptimekuma_tag.production.id
      value  = "main-store"
    },
    {
      tag_id = uptimekuma_tag.critical.id
      value  = "revenue-generating"
    }
  ]
}
