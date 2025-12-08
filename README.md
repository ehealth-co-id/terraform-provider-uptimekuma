# Terraform Provider for Uptime Kuma

This Terraform provider allows you to manage [Uptime Kuma](https://github.com/louislam/uptime-kuma) resources through Terraform. Uptime Kuma is a self-hosted monitoring tool similar to Uptime Robot.

## Features

- **Monitors**: Create and manage HTTP, Ping, Port, DNS, Keyword, and other monitor types
- **Status Pages**: Create and manage status pages with monitor groups and custom domains
- **Tags**: Create and manage tags for organizing monitors
- **Direct Socket.IO Connection**: Communicates directly with Uptime Kuma v2 (no middleware required)

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23 (for development)
- Uptime Kuma v2 instance

## Version Compatibility

| Provider Version | Uptime Kuma Version | Connection Method |
|------------------|---------------------|-------------------|
| **v1.0.0+**      | **v2.x**            | Direct Socket.IO  |
| v0.x (deprecated)| v1.21.3             | HTTP Middleware   |

**⚠️ Breaking Change:** Version 1.0.0+ uses direct Socket.IO communication and is only compatible with Uptime Kuma v2. The HTTP middleware adapter is no longer required or supported.

## Quick Start

### 1. Start Uptime Kuma v2

```bash
docker compose up -d
```

The included `docker-compose.yml` automatically runs Uptime Kuma v2 on port 3001.

### 2. Setup Admin Account (First Time Only)

```bash
# Automated setup using Playwright
npm install
node scripts/setup-uptime-kuma.js

# Or manually visit http://localhost:3001 and create admin account
# Username: admin, Password: admin123
```

### 3. Configure Provider

```hcl
terraform {
  required_providers {
    uptimekuma = {
      source  = "ehealth-co-id/uptimekuma"
      version = "~> 1.0"
    }
  }
}

provider "uptimekuma" {
  base_url = "http://localhost:3001"  # Direct Uptime Kuma URL
  username = "admin"
  password = "admin123"
}

# Create an HTTP monitor
resource "uptimekuma_monitor" "website" {
  name           = "Company Website"
  type           = "http"
  url            = "https://example.com"
  interval       = 60
  retry_interval = 30
  max_retries    = 3
}

# Create a status page
resource "uptimekuma_status_page" "status" {
  slug        = "status"
  title       = "System Status"
  description = "Current status of our services"
  theme       = "dark"
  published   = true
  
  public_group_list = [
    {
      name = "Public Services"
      weight = 1
      monitor_list = [uptimekuma_monitor.website.id]
    }
  ]
}
```

See the [examples](./examples/) directory for more detailed examples.

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install
```

## Resource Documentation

### uptimekuma_monitor

The `uptimekuma_monitor` resource allows you to create and manage monitors in Uptime Kuma.

#### Example Usage

```hcl
# HTTP Monitor
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
  
  accepted_status_codes = [200, 201]
}

# Ping Monitor
resource "uptimekuma_monitor" "ping_example" {
  name           = "Ping Example"
  type           = "ping"
  hostname       = "example.com"
  interval       = 120
  retry_interval = 30
  max_retries    = 2
}

# Port Monitor
resource "uptimekuma_monitor" "port_example" {
  name           = "Port Example"
  type           = "port"
  hostname       = "example.com"
  port           = 443
  interval       = 60
  retry_interval = 30
  max_retries    = 1
}
```

#### Argument Reference

* `name` - (Required) The name of the monitor.
* `type` - (Required) The type of monitor. Valid values: `http`, `ping`, `port`, `keyword`.
* `interval` - (Optional) The interval in seconds between checks. Default: `60`.
* `retry_interval` - (Optional) The interval in seconds between retries. Default: `60`.
* `resend_interval` - (Optional) The interval in seconds for resending notifications. Default: `0`.
* `max_retries` - (Optional) The maximum number of retries. Default: `0`.
* `upside_down` - (Optional) Whether to invert status (treat DOWN as UP and vice versa). Default: `false`.
* `ignore_tls` - (Optional) Whether to ignore TLS errors. Default: `false`.

**HTTP Monitor Arguments:**
* `url` - (Required for HTTP monitors) The URL to monitor.
* `method` - (Optional) The HTTP method to use. Default: `GET`.
* `max_redirects` - (Optional) The maximum number of redirects to follow. Default: `0`.
* `body` - (Optional) The request body for HTTP POST/PUT/PATCH requests.
* `headers` - (Optional) JSON string of request headers.
* `auth_method` - (Optional) Authentication method. Valid values: `basic`, `ntlm`, `mtls`.
* `basic_auth_user` - (Optional) Basic auth username.
* `basic_auth_pass` - (Optional) Basic auth password.
* `accepted_status_codes` - (Optional) List of accepted HTTP status codes.

**Ping/Port Monitor Arguments:**
* `hostname` - (Required for ping/port monitors) The hostname to check.
* `port` - (Required for port monitors) The port number to check.

**Keyword Monitor Arguments:**
* `url` - (Required for keyword monitors) The URL to search for keywords.
* `keyword` - (Required for keyword monitors) The keyword to search for.

### uptimekuma_status_page

The `uptimekuma_status_page` resource allows you to create and manage status pages in Uptime Kuma.

#### Example Usage

```hcl
resource "uptimekuma_status_page" "company_status" {
  slug        = "status"
  title       = "Company Status Page"
  description = "Current status of our services"
  theme       = "dark"
  published   = true
  
  public_group_list = [
    {
      name = "Core Services"
      weight = 1
      monitor_list = [
        uptimekuma_monitor.website.id
      ]
    },
    {
      name = "API Services"
      weight = 2
      monitor_list = [
        uptimekuma_monitor.api.id
      ]
    }
  ]
}
```

#### Argument Reference

* `slug` - (Required) The URL slug for the status page.
* `title` - (Required) The title of the status page.
* `description` - (Optional) The description of the status page.
* `theme` - (Optional) The theme for the status page.
* `published` - (Optional) Whether the status page is published. Default: `true`.
* `show_tags` - (Optional) Whether to show tags on the status page. Default: `false`.
* `domain_name_list` - (Optional) A list of custom domains for the status page.
* `footer_text` - (Optional) Custom footer text.
* `custom_css` - (Optional) Custom CSS for the status page.
* `google_analytics_id` - (Optional) Google Analytics ID.
* `icon` - (Optional) URL to a custom icon. Default: `/icon.svg`.
* `show_powered_by` - (Optional) Whether to show "Powered by Uptime Kuma" text. Default: `true`.
* `public_group_list` - (Optional) A list of monitor groups to display on the status page.
  * `name` - (Required) The name of the group.
  * `weight` - (Optional) The order/weight of the group.
  * `monitor_list` - (Optional) A list of monitor IDs to include in the group.

### uptimekuma_tag

The `uptimekuma_tag` resource allows you to create and manage tags in Uptime Kuma.

#### Example Usage

```hcl
resource "uptimekuma_tag" "production" {
  name  = "production"
  color = "#00FF00"
}
```

#### Argument Reference

* `name` - (Required) The name of the tag.
* `color` - (Required) The color of the tag in hex format (e.g., `#00FF00`).

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

### Running Tests

This provider includes acceptance tests that run against a real Uptime Kuma v2 instance.

```shell
# Start the local test environment
docker compose up -d

# Setup admin account (first time only)
node scripts/setup-uptime-kuma.js

# Set required environment variables
export TF_ACC=1
export UPTIMEKUMA_BASE_URL="http://localhost:3001"
export UPTIMEKUMA_USERNAME="admin"
export UPTIMEKUMA_PASSWORD="admin123"

# Run acceptance tests
go test -v -p 1 ./internal/provider/

# Clean up when done
docker compose down -v
```

*Note:* Acceptance tests create and destroy real resources. Use with caution on production instances.

### Generate Documentation

To generate or update documentation, run:

```shell
make generate
```

## Migration from v0.x

If you're upgrading from v0.x:

1. **Remove Middleware**: The HTTP middleware adapter (`uptime-kuma-api`) is no longer needed
2. **Update Configuration**: Change `base_url` to point directly to Uptime Kuma (e.g., `http://localhost:3001` instead of `http://localhost:8000`)
3. **Upgrade Uptime Kuma**: Ensure you're running Uptime Kuma v2.x
4. **Remove `insecure_https`**: This option is no longer supported

See [CHANGELOG.md](./CHANGELOG.md) for complete breaking changes.

## Architecture

The provider is built using the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) and communicates directly with Uptime Kuma via Socket.IO:

1. **Client Layer** (`internal/client/`): Socket.IO client wrapper with retry logic
2. **Provider Layer** (`internal/provider/`): Terraform resource implementations

For more details, see [ARCHITECTURE.md](./ARCHITECTURE.md).

## License

[MPL 2.0](LICENSE)