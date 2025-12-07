provider "uptimekuma" {
  base_url = "http://localhost:8000" # Your Uptime Kuma Web API adapter URL (not direct Uptime Kuma URL)
  username = "admin"                 # Username for authentication
  password = "password"              # Password for authentication
  # insecure_https = true            # Optional: Skip TLS certificate verification
}
