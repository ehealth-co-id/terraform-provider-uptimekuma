# Terraform Provider Architecture for Uptime Kuma

## Provider Structure

```
terraform-provider-uptimekuma/
├── internal/
│   ├── provider/
│   │   ├── provider.go                     # Provider definition
│   │   ├── monitor_resource.go             # Monitor resource
│   │   ├── status_page_resource.go         # Status page resource
│   │   ├── tag_resource.go                 # Tag resource
│   │   ├── monitor_resource_test.go        # Monitor acceptance tests
│   │   ├── status_page_resource_test.go    # Status page acceptance tests
│   │   └── provider_test.go                # Provider test utilities
│   └── client/
│       └── client.go                       # Socket.IO client wrapper
└── go-uptime-kuma-client/                  # External Socket.IO library (submodule)
```

## Architecture Overview

This provider uses a **direct Socket.IO connection** to Uptime Kuma v2, eliminating the need for an HTTP middleware adapter.

### Communication Layer

```
┌─────────────┐
│  Terraform  │
└──────┬──────┘
       │ Plugin Protocol
       ▼
┌─────────────────────────────┐
│  terraform-provider-uptimekuma │
│  (internal/provider/)        │
└──────┬──────────────────────┘
       │
       ▼
┌─────────────────────────────┐
│  Client Wrapper             │
│  (internal/client/)         │
└──────┬──────────────────────┘
       │
       ▼
┌─────────────────────────────┐
│  go-uptime-kuma-client      │
│  (Socket.IO library)        │
└──────┬──────────────────────┘
       │ Socket.IO + WebSocket
       ▼
┌─────────────────────────────┐
│  Uptime Kuma v2             │
│  (Server)                   │
└─────────────────────────────┘
```

## Provider Configuration

```hcl
provider "uptimekuma" {
  base_url = "http://localhost:3001"  # Direct Uptime Kuma URL
  username = "admin"                  # Username for authentication
  password = "password"               # Password for authentication
}
```

**Note:** The `base_url` now points directly to your Uptime Kuma instance, not to a middleware adapter.

## Key Resources

### Monitor Resource

```hcl
resource "uptimekuma_monitor" "website" {
  name           = "API Health Check"
  type           = "http"
  url            = "https://api.example.com/health"
  method         = "GET"
  interval       = 60
  retry_interval = 60
  max_retries    = 3
  upside_down    = false
  ignore_tls     = false
  
  # HTTP-specific
  accepted_status_codes = [200, 201]
  headers               = jsonencode({"Authorization" = "Bearer token"})
  body                 = jsonencode({"check" = "full"})
  
  # Authentication
  auth_method        = "basic"
  basic_auth_user    = "user"
  basic_auth_pass    = "pass"
}
```

**Supported Monitor Types:**
- `http` - HTTP/HTTPS monitoring
- `ping` - ICMP ping monitoring
- `port` - TCP port monitoring
- `keyword` - HTTP keyword search monitoring

### Status Page Resource

```hcl
resource "uptimekuma_status_page" "main" {
  slug        = "main-status"
  title       = "Service Status"
  description = "Current status of our services"
  theme       = "dark"
  published   = true
  
  public_group_list = [
    {
      name = "API Services"
      weight = 1
      monitor_list = [
        uptimekuma_monitor.website.id
      ]
    }
  ]
}
```

### Tag Resource

```hcl
resource "uptimekuma_tag" "production" {
  name  = "production"
  color = "#00FF00"
}
```

## Connection and Authentication

The provider establishes a Socket.IO connection during configuration and maintains it throughout the Terraform run.

### Retry Logic

The client implements exponential backoff with jitter to handle transient connection failures:

```go
// Retry configuration in internal/client/client.go
maxRetries := 5
baseDelay := 2 * time.Second

for i := 0; i <= maxRetries; i++ {
    k, err = kuma.New(ctx, config.BaseURL, config.Username, config.Password)
    if err == nil {
        return &Client{Kuma: k}, nil
    }
    
    if i == maxRetries {
        break
    }
    
    // Exponential backoff with jitter
    backoff := float64(baseDelay) * math.Pow(2, float64(i))
    r := rand.Float64()*0.4 + 0.8  // ±20% jitter
    sleepDuration := time.Duration(backoff * r)
    
    // Cap at 30 seconds
    if sleepDuration > 30*time.Second {
        sleepDuration = 30 * time.Second
    }
    
    time.Sleep(sleepDuration)
}
```

This handles rate limiting errors (e.g., "login: Too frequently") that can occur in CI/CD environments.

### Connection Pooling (Test-Only)

For acceptance tests, the provider supports connection pooling to prevent rate limiting when running multiple test cases:

```go
// Enable via environment variable (set automatically in TestMain)
export UPTIMEKUMA_ENABLE_CONNECTION_POOL="true"
```

**How it works:**
- A global connection pool (`internal/client/pool.go`) maintains a shared Socket.IO connection
- Multiple provider instances during test execution reuse the same connection
- Reference counting prevents premature disconnection
- Automatic cleanup via `TestMain()` after all tests complete

**Benefits:**
- Reduces login frequency from ~10-20 per test run → 1 per test run
- Eliminates "login: Too frequently" errors during acceptance testing
- 10-20% faster test execution (reduced login overhead)

**Production behavior:** Connection pooling is disabled by default and only activates in test scenarios. Production usage creates a new connection for each provider instance (existing behavior).

## State Management

### Optional Fields

The provider carefully handles optional fields to prevent state drift:

1. **Empty Strings → Null**: When reading from API, empty optional strings are mapped to `null` to match Terraform configuration
2. **Default Values**: Schema defaults are set for fields like `method` (GET), `interval` (60s), `upside_down` (false)
3. **List Initialization**: Empty lists are initialized as `[]` instead of `null` when sent to Uptime Kuma v2

### Status Page Groups

Due to API limitations with `GetStatusPage` (doesn't return groups) and eventual consistency with `GetStatusPages` cache, the provider implements a "preserve state" strategy:

- If groups are found in the API/cache, they're updated in state
- If groups are not found, the existing Terraform state is preserved
- This prevents false drift detection while maintaining accuracy when data is available

## Testing

### Acceptance Tests

Run against a real Uptime Kuma v2 instance:

```bash
# Start Uptime Kuma v2
docker compose up -d

# Setup admin account (automated)
node scripts/setup-uptime-kuma.js

# Run tests
export TF_ACC=1
export UPTIMEKUMA_BASE_URL="http://localhost:3001"
export UPTIMEKUMA_USERNAME="admin"
export UPTIMEKUMA_PASSWORD="admin123"
go test -v -p 1 ./internal/provider/
```

### Coverage

- **Acceptance Tests**: 4 comprehensive tests covering Monitor, Ping, StatusPage, and StatusPage with Groups
- **Unit Tests**: Minimal (thin wrapper architecture makes unit tests of limited value)

## Design Principles

1. **Direct Connection**: No middleware required - connects directly to Uptime Kuma
2. **Resilience**: Exponential backoff with jitter for connection reliability
3. **State Accuracy**: Careful null handling to prevent drift
4. **Compatibility**: Designed for Uptime Kuma v2
5. **Simplicity**: Thin wrapper over well-tested Socket.IO library