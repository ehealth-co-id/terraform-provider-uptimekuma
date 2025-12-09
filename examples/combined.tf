# Configure the Uptime Kuma provider
terraform {
  required_providers {
    uptimekuma = {
      source  = "ehealth-co-id/uptimekuma"
      version = "~> 1.0"
    }
  }
}

provider "uptimekuma" {
  base_url = "http://localhost:3001" # Direct Uptime Kuma URL
  username = "admin"
  password = "password"
}

# Create tags for organizing monitors
resource "uptimekuma_tag" "production" {
  name  = "production"
  color = "#00FF00"
}

resource "uptimekuma_tag" "staging" {
  name  = "staging"
  color = "#FFA500"
}

resource "uptimekuma_tag" "critical" {
  name  = "critical"
  color = "#FF0000"
}

resource "uptimekuma_tag" "infrastructure" {
  name  = "infrastructure"
  color = "#0066CC"
}

# Create HTTP monitors for different services
resource "uptimekuma_monitor" "website" {
  name           = "Company Website"
  type           = "http"
  url            = "https://example.com"
  interval       = 60
  retry_interval = 30
  max_retries    = 3

  # Tag the website monitor as production and critical
  tags = [
    {
      tag_id = uptimekuma_tag.production.id
    },
    {
      tag_id = uptimekuma_tag.critical.id
      value  = "high-priority"
    }
  ]
}

resource "uptimekuma_monitor" "api" {
  name           = "API Service"
  type           = "http"
  url            = "https://api.example.com/health"
  method         = "GET"
  interval       = 30
  retry_interval = 10
  max_retries    = 3
  ignore_tls     = false

  # Optional: Custom headers for API auth
  headers = "{\"X-API-Key\":\"dummy-key\", \"Accept\":\"application/json\"}"

  # Tag the API as production
  tags = [
    {
      tag_id = uptimekuma_tag.production.id
    }
  ]
}

resource "uptimekuma_monitor" "staging_api" {
  name           = "Staging API"
  type           = "http"
  url            = "https://staging-api.example.com/health"
  method         = "GET"
  interval       = 60
  retry_interval = 30
  max_retries    = 2

  # Tag as staging environment
  tags = [
    {
      tag_id = uptimekuma_tag.staging.id
      value  = "v2-testing"
    }
  ]
}

resource "uptimekuma_monitor" "database" {
  name           = "Database Server"
  type           = "ping"
  hostname       = "db.internal.example.com"
  interval       = 60
  retry_interval = 30
  max_retries    = 2

  # Tag as infrastructure and critical
  tags = [
    {
      tag_id = uptimekuma_tag.infrastructure.id
    },
    {
      tag_id = uptimekuma_tag.critical.id
      value  = "database"
    }
  ]
}

# Create a status page with all monitors
resource "uptimekuma_status_page" "main_status" {
  slug        = "status"
  title       = "System Status"
  description = "Current status of Example Company services"
  theme       = "dark"
  published   = true

  # Group monitors by service type
  public_group_list = [
    {
      name         = "Public Services"
      weight       = 1
      monitor_list = [uptimekuma_monitor.website.id]
    },
    {
      name         = "API Services"
      weight       = 2
      monitor_list = [uptimekuma_monitor.api.id, uptimekuma_monitor.staging_api.id]
    },
    {
      name         = "Infrastructure"
      weight       = 3
      monitor_list = [uptimekuma_monitor.database.id]
    }
  ]
}
