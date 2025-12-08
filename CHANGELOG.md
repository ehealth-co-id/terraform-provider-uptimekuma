# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 1.0.1

IMPROVEMENTS:

* **Color Validation**: Added validation for Tag color field to ensure valid hex codes (e.g., #RRGGBB)

## 1.0.0

BREAKING CHANGES:

* **Provider Architecture**: Complete refactor to use Socket.IO directly instead of HTTP middleware adapter
* **Configuration**: `base_url` now points directly to Uptime Kuma instance (not middleware)
* **Removed**: `insecure_https` configuration option (not supported by Socket.IO client library)
* **Compatibility**: Now requires Uptime Kuma v2

FEATURES:

* **Direct Socket.IO Connection**: Eliminated dependency on `uptime-kuma-api` middleware
* **Connection Resilience**: Implemented exponential backoff with jitter for login retry logic
* **Connection Pooling (Test-Only)**: Added connection pool for acceptance tests to prevent rate limiting
* **Monitor Resource**: Complete HTTP/HTTPS, Ping, Port, and Keyword monitoring support
* **Status Page Resource**: Create and manage status pages with public group support
* **Tag Resource**: Tag management for organizing monitors

IMPROVEMENTS:

* **State Management**: Fixed state drift for optional fields (empty strings now correctly map to null)
* **Default Values**: Added schema defaults for boolean and numeric fields to match API behavior
* **Status Page Groups**: Implemented state preservation strategy for public groups to handle API cache consistency
* **Testing**: Updated CI workflow and acceptance tests for Uptime Kuma v2
* **Test Performance**: Connection pooling reduces test execution time by ~25% and eliminates login rate limiting

BUG FIXES:

* Fixed null value handling in Monitor resource (`AcceptedStatusCodes` initialization)
* Fixed state drift for optional string fields in Monitor and StatusPage resources
* Fixed concurrent login rate limiting with retry backoff mechanism
* Fixed "login: Too frequently" errors during acceptance testing with connection pooling
