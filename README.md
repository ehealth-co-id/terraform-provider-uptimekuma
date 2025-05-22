# Terraform Provider for Uptime Kuma

This Terraform provider allows you to manage [Uptime Kuma](https://github.com/louislam/uptime-kuma) resources through Terraform. Uptime Kuma is a self-hosted monitoring tool similar to Uptime Robot.

## Features

- **Monitors**: Create and manage HTTP, Ping, Port, DNS, Keyword, and other monitor types
- **Status Pages**: Create and manage status pages with monitor groups and custom domains

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23
- Access to an Uptime Kuma instance (self-hosted or hosted)
- Uptime Kuma Web API adapter (see setup guide below)

## Setup Guide for Uptime Kuma Web API

Uptime Kuma uses WebSocket for its API rather than a traditional REST API. This provider requires the [Uptime Kuma Web API adapter](https://github.com/MedAziz11/Uptime-Kuma-Web-API) to bridge between standard REST calls and Uptime Kuma's WebSocket API.

For a complete setup guide, refer to: https://github.com/ehealth-co-id/uptime-kuma-api-starter-pack

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install
```

## Using the provider

Configure the provider in your Terraform configuration:

```hcl
terraform {
  required_providers {
    uptimekuma = {
      source = "ehealth-co-id/uptimekuma"
    }
  }
}

provider "uptimekuma" {
  base_url = "https://localhost/api/v1/"  # Your Uptime Kuma Web API adapter URL (not direct Uptime Kuma URL)
  username = "admin"                      # Username for authentication
  password = "password"                   # Password for authentication
  # insecure_https = true                 # Optional: Skip TLS certificate verification
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

### Resource: uptimekuma_monitor

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
* `type` - (Required) The type of monitor. Valid values: `http`, `ping`, `port`, `dns`, `keyword`, `grpc-keyword`, `docker`, `push`, `steam`, `gamedig`, `mqtt`, `sqlserver`, `postgres`, `mysql`, `mongodb`, `radius`, `redis`.
* `interval` - (Optional) The interval in seconds between checks. Default: `60`.
* `retry_interval` - (Optional) The interval in seconds between retries. Default: `60`.
* `resend_interval` - (Optional) The interval in seconds for resending notifications. Default: `0`.
* `max_retries` - (Optional) The maximum number of retries. Default: `0`.
* `upside_down` - (Optional) Whether to invert status (treat DOWN as UP and vice versa). Default: `false`.
* `ignore_tls` - (Optional) Whether to ignore TLS errors. Default: `false`.

**HTTP Monitor Arguments:**
* `url` - (Required for HTTP monitors) The URL to monitor.
* `method` - (Optional) The HTTP method to use. Default: `GET`.
* `max_redirects` - (Optional) The maximum number of redirects to follow.
* `body` - (Optional) The request body for HTTP POST/PUT/PATCH requests.
* `headers` - (Optional) JSON string of request headers.
* `auth_method` - (Optional) Authentication method. Valid values: `basic`, `ntlm`, `mtls`.
* `basic_auth_user` - (Optional) Basic auth username.
* `basic_auth_pass` - (Optional) Basic auth password.

**Ping/Port Monitor Arguments:**
* `hostname` - (Required for ping/port monitors) The hostname to check.
* `port` - (Required for port monitors) The port number to check.

**Keyword Monitor Arguments:**
* `url` - (Required for keyword monitors) The URL to search for keywords.
* `keyword` - (Required for keyword monitors) The keyword to search for.

### Resource: uptimekuma_status_page

The `uptimekuma_status_page` resource allows you to create and manage status pages in Uptime Kuma.

#### Example Usage

```hcl
resource "uptimekuma_status_page" "company_status" {
  slug        = "status"
  title       = "Company Status Page"
  description = "Current status of our services"
  theme       = "dark"
  published   = true
  
  # Group monitors on the status page
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
* `theme` - (Optional) The theme for the status page. Options: `light`, `dark`.
* `published` - (Optional) Whether the status page is published. Default: `true`.
* `show_tags` - (Optional) Whether to show tags on the status page. Default: `false`.
* `domain_name_list` - (Optional) A list of custom domains for the status page.
* `footer_text` - (Optional) Custom footer text.
* `custom_css` - (Optional) Custom CSS for the status page.
* `google_analytics_id` - (Optional) Google Analytics ID.
* `icon` - (Optional) URL to a custom icon.
* `show_powered_by` - (Optional) Whether to show "Powered by Uptime Kuma" text. Default: `true`.
* `public_group_list` - (Optional) A list of monitor groups to display on the status page.
  * `name` - (Required) The name of the group.
  * `weight` - (Optional) The order/weight of the group.
  * `monitor_list` - (Optional) A list of monitor IDs to include in the group.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

### Setting Up Local Environment

1. Clone the repository
2. Install dependencies: `go mod download`

### Running Tests

This provider includes both unit tests and acceptance tests.

**Running Unit Tests**

Unit tests run against mock API endpoints and don't require a real Uptime Kuma instance:

```shell
# Run all unit tests
go test -v ./...

# Run specific test
go test -v ./internal/client -run TestClientAuthentication

# Run tests with coverage
go test -v -cover ./...
```

**Running Acceptance Tests**

Acceptance tests create real resources on your Uptime Kuma instance. These require an actual Uptime Kuma instance and valid credentials:

```shell
# Set required environment variables
export TF_ACC=1
export UPTIMEKUMA_BASE_URL="https://localhost/api/v1"
export UPTIMEKUMA_USERNAME="admin"
export UPTIMEKUMA_PASSWORD="mypassword"

# Run acceptance tests
go test -v ./internal/provider
```

*Note:* Acceptance tests create and destroy real resources. Use with caution on production instances.

For more details on testing, see [TESTING.md](./TESTING.md).

### Generate Documentation

To generate or update documentation, run:

```shell
make generate
```

## Architecture

The provider is built using the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) and follows a clean architecture pattern:

1. **Client Layer** (`internal/client/`): API client handling authentication and API calls
2. **Provider Layer** (`internal/provider/`): Terraform resource and data source implementations

For more details, see [ARCHITECTURE.md](./ARCHITECTURE.md).

## License

[MPL 2.0](LICENSE)
