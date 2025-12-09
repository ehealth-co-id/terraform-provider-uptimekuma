# Terraform Provider Uptime Kuma Examples

This directory contains examples for all resources provided by the Uptime Kuma Terraform provider.

## Quick Start

```bash
# Initialize terraform
terraform init

# Plan changes
terraform plan

# Apply configuration
terraform apply
```

## Prerequisites

Before using these examples, ensure you have:

1. **Uptime Kuma v2** running and accessible
2. **Provider configured** with valid credentials

```hcl
provider "uptimekuma" {
  base_url = "http://localhost:3001"  # Your Uptime Kuma instance URL
  username = "admin"
  password = "password"
}
```

## Examples

### Combined Example
[`combined.tf`](./combined.tf) - Complete example showing all resources working together including tags

### Provider Configuration
[`provider/provider.tf`](./provider/provider.tf) - Basic provider configuration

### Resources

#### Monitor Resource
[`resources/uptimekuma_monitor/resource.tf`](./resources/uptimekuma_monitor/resource.tf)

Examples of different monitor types:
- **HTTP Monitor** - Monitor web endpoints with custom headers and authentication
- **Ping Monitor** - ICMP ping monitoring
- **Port Monitor** - TCP port availability checks
- **Keyword Monitor** - Search for keywords in HTTP responses
- **Tagged Monitor** - Monitor with tags for organization

#### Status Page Resource
[`resources/uptimekuma_status_page/resource.tf`](./resources/uptimekuma_status_page/resource.tf)

Create public status pages with:
- Custom domains
- Monitor groups
- Custom styling (CSS)
- Google Analytics integration

#### Tag Resource
[`resources/uptimekuma_tag/resource.tf`](./resources/uptimekuma_tag/resource.tf)

Create and manage tags for organizing monitors:
- Color-coded tags
- Multiple tag examples (production, staging, critical, infrastructure)

## Usage Patterns

### Creating a Complete Monitoring Setup with Tags

```hcl
# 1. Create tags for organization
resource "uptimekuma_tag" "production" {
  name  = "production"
  color = "#00FF00"
}

resource "uptimekuma_tag" "critical" {
  name  = "critical"
  color = "#FF0000"
}

# 2. Create monitors with tags
resource "uptimekuma_monitor" "api" {
  name     = "API Service"
  type     = "http"
  url      = "https://api.example.com/health"
  interval = 60

  # Associate tags with the monitor
  tags = [
    {
      tag_id = uptimekuma_tag.production.id
    },
    {
      tag_id = uptimekuma_tag.critical.id
      value  = "high-priority"  # Optional value for context
    }
  ]
}

# 3. Create a status page
resource "uptimekuma_status_page" "main" {
  slug  = "status"
  title = "System Status"
  
  public_group_list = [
    {
      name         = "Services"
      weight       = 1
      monitor_list = [uptimekuma_monitor.api.id]
    }
  ]
}
```

### Using Tags for Environment Organization

```hcl
# Environment tags
resource "uptimekuma_tag" "prod" {
  name  = "prod"
  color = "#00FF00"
}

resource "uptimekuma_tag" "staging" {
  name  = "staging"
  color = "#FFA500"
}

resource "uptimekuma_tag" "dev" {
  name  = "dev"
  color = "#808080"
}

# Production monitor
resource "uptimekuma_monitor" "prod_api" {
  name = "Production API"
  type = "http"
  url  = "https://api.example.com"

  tags = [
    { tag_id = uptimekuma_tag.prod.id }
  ]
}

# Staging monitor
resource "uptimekuma_monitor" "staging_api" {
  name = "Staging API"
  type = "http"
  url  = "https://staging.api.example.com"

  tags = [
    { tag_id = uptimekuma_tag.staging.id, value = "v2-testing" }
  ]
}
```

### Using Tags with Values for Custom Metadata

```hcl
resource "uptimekuma_tag" "team" {
  name  = "team"
  color = "#0066CC"
}

resource "uptimekuma_tag" "sla" {
  name  = "sla"
  color = "#9933FF"
}

resource "uptimekuma_monitor" "payment_api" {
  name = "Payment API"
  type = "http"
  url  = "https://payment.example.com/health"

  tags = [
    { tag_id = uptimekuma_tag.team.id, value = "payments-team" },
    { tag_id = uptimekuma_tag.sla.id, value = "99.99%" }
  ]
}
```

### Import Existing Resources

```bash
# Import monitor by ID
terraform import uptimekuma_monitor.website 123

# Import status page by slug
terraform import uptimekuma_status_page.main status-slug

# Import tag by ID
terraform import uptimekuma_tag.production 1
```

## Resource Documentation

For detailed information about each resource, see:
- [Monitor Resource Documentation](../docs/resources/monitor.md)
- [Status Page Resource Documentation](../docs/resources/status_page.md)
- [Tag Resource Documentation](../docs/resources/tag.md)

## Version Compatibility

These examples are designed for:
- **Provider Version**: v1.0.0+
- **Uptime Kuma Version**: v2.x
- **Terraform Version**: >= 1.0

## Support

For issues or questions:
- [GitHub Issues](https://github.com/ehealth-co-id/terraform-provider-uptimekuma/issues)
- [Provider Documentation](https://registry.terraform.io/providers/ehealth-co-id/uptimekuma/latest/docs)
