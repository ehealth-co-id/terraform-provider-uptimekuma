# Copyright (c) eHealth.co.id as PT Aksara Digital Indonesia
# SPDX-License-Identifier: MPL-2.0

# Tags are used to organize and categorize monitors in Uptime Kuma

resource "uptimekuma_tag" "production" {
  # Name: Tag display name (string, required)
  # Used to identify and filter monitors
  name = "production"

  # Color: Tag color in hex format (string, required)
  # Format: #RRGGBB (e.g., #00FF00 for green)
  # Used for visual distinction in the UI
  color = "#00FF00"
}

resource "uptimekuma_tag" "staging" {
  name = "staging"
  # Orange color for staging environment
  color = "#FFA500"
}

resource "uptimekuma_tag" "critical" {
  name = "critical"
  # Red color for critical services
  color = "#FF0000"
}

resource "uptimekuma_tag" "infrastructure" {
  name = "infrastructure"
  # Blue color for infrastructure components
  color = "#0066CC"
}

# Additional tag color examples:
# Purple: #9933FF
# Yellow: #FFFF00
# Pink: #FF69B4
# Cyan: #00FFFF
# Dark Gray: #666666
