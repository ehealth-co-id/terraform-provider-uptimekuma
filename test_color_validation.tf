# Copyright (c) eHealth.co.id as PT Aksara Digital Indonesia
# SPDX-License-Identifier: MPL-2.0

terraform {
  required_providers {
    uptimekuma = {
      source = "registry.terraform.io/ehealth-co-id/uptimekuma"
    }
  }
}

provider "uptimekuma" {
  base_url = "http://localhost:3001"
  username = "admin"
  password = "admin123"
}

# This should fail with validation error - invalid color format
resource "uptimekuma_tag" "invalid_color" {
  name  = "test-tag"
  color = "red"  # Invalid: not a hex code
}

# This should succeed - valid 3-digit hex
resource "uptimekuma_tag" "valid_short" {
  name  = "test-tag-short"
  color = "#FFF"
}

# This should succeed - valid 6-digit hex
resource "uptimekuma_tag" "valid_long" {
  name  = "test-tag-long"
  color = "#FF0000"
}
