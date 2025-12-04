# Copyright (c) eHealth.co.id as PT Aksara Digital Indonesia
# SPDX-License-Identifier: MPL-2.0

provider "uptimekuma" {
  base_url       = "https://localhost/api/v1/" # Your Uptime Kuma Web API adapter URL (not direct Uptime Kuma URL)
  username       = "admin"                     # Username for authentication
  password       = "password"                  # Password for authentication
  insecure_https = true                        # Optional: Skip TLS certificate verification
}
