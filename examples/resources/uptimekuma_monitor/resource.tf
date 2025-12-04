# Copyright (c) eHealth.co.id as PT Aksara Digital Indonesia
# SPDX-License-Identifier: MPL-2.0

resource "uptimekuma_monitor" "http_example" {
  name           = "Example Website"
  type           = "http"
  url            = "https://example.com"
  method         = "GET"
  interval       = 60
  retry_interval = 30
  max_retries    = 3
  upside_down    = false
  ignore_tls     = false
}

resource "uptimekuma_monitor" "ping_example" {
  name           = "Ping Example"
  type           = "ping"
  hostname       = "example.com"
  interval       = 120
  retry_interval = 30
  max_retries    = 2
}

resource "uptimekuma_monitor" "keyword_example" {
  name          = "Keyword Search Example"
  type          = "keyword"
  url           = "https://example.com"
  method        = "GET"
  interval      = 300
  keyword       = "Example Domain"
  upside_down   = true # Invert status (alert when keyword is found)
  max_redirects = 3
}

resource "uptimekuma_monitor" "port_example" {
  name           = "Port Example"
  type           = "port"
  hostname       = "example.com"
  port           = 443
  interval       = 60
  retry_interval = 30
  max_retries    = 1
}

resource "uptimekuma_monitor" "authenticated_http" {
  name           = "Authenticated API"
  type           = "http"
  url            = "https://api.example.com/private"
  method         = "GET"
  interval       = 60
  retry_interval = 30
  max_retries    = 3

  # Basic authentication
  auth_method     = "basic"
  basic_auth_user = "apiuser"
  basic_auth_pass = "securepassword"

  # Custom headers (JSON formatted)
  headers = "{\"X-API-Key\":\"myapikey\", \"Accept\":\"application/json\"}"
}