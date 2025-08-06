# Uptime Kuma Terraform Provider Examples

This directory contains examples demonstrating various features of the Uptime Kuma Terraform provider.

## Prerequisites

1. **Uptime Kuma Server**: Ensure you have Uptime Kuma running on `http://localhost:3001`
2. **Admin Credentials**: Update the provider configuration with your admin username and password
3. **OpenTofu/Terraform**: Install OpenTofu or Terraform

## Usage

### Initialize the Provider

```bash
tofu init
```

### Plan Changes

```bash
tofu plan
```

### Apply Configuration

> **Important**: Use `-parallelism=1` to avoid WebSocket concurrency issues

```bash
tofu apply -parallelism=1
```

## Examples Included

### 1. HTTP Monitor (`google`)
- **Type**: HTTP
- **URL**: https://google.com
- **Features**: Tags, follow redirects, accepted status codes
- **Interval**: 60 seconds

### 2. TCP Monitor (`database`)
- **Type**: TCP
- **Target**: db.example.com:5432
- **Features**: Tags, custom timeout
- **Interval**: 30 seconds

### 3. HTTP Monitor with Basic Auth (`secure_api`)
- **Type**: HTTP
- **URL**: https://api.example.com/health
- **Features**: Basic authentication, TLS validation, tags
- **Interval**: 120 seconds

## Supported Features

### Monitor Types
- ✅ **HTTP**: Web service monitoring
- ✅ **TCP**: Port connectivity monitoring
- 🔄 **Ping**: (Schema ready, implementation pending)

### Authentication
- ✅ **Basic Auth**: Username/password authentication
- 🔄 **Headers**: Custom headers (Schema ready)
- 🔄 **OAuth**: OAuth token authentication (Future)

### Configuration Options
- ✅ **Intervals**: Check interval, timeout, retry interval
- ✅ **Redirects**: Follow redirects, max redirects
- ✅ **Status Codes**: Accepted HTTP status codes
- ✅ **TLS**: Ignore TLS certificate errors
- ✅ **Tags**: Organize monitors with tags
- ✅ **Active/Inactive**: Enable/disable monitoring

### HTTP Options
- ✅ **Methods**: GET, POST, PUT, DELETE, etc.
- ✅ **Body**: Request body for POST/PUT
- ✅ **Headers**: Custom headers (Schema ready)

## Configuration Reference

### Provider Configuration

```hcl
provider "uptimekuma" {
  url      = "http://localhost:3001"    # Uptime Kuma server URL
  username = "admin"                    # Admin username
  password = "your-password"            # Admin password
}
```

### Monitor Resource

```hcl
resource "uptimekuma_monitor" "example" {
  name = "Example Monitor"
  type = "http"                         # http, tcp, ping
  url  = "https://example.com"
  
  # Timing
  interval        = 60                  # Check interval (seconds)
  timeout         = 30                  # Request timeout (seconds)
  retry_interval  = 60                  # Retry interval (seconds)
  max_retries     = 3                   # Maximum retries
  
  # HTTP Options
  http_method     = "GET"               # HTTP method
  follow_redirect = true                # Follow redirects
  max_redirects   = 10                  # Maximum redirects
  accepted_status_codes = ["200", "301"] # Accepted status codes
  
  # Authentication
  basic_auth_user = "username"          # Basic auth username
  basic_auth_pass = "password"          # Basic auth password
  
  # TLS
  ignore_tls      = false               # Ignore TLS errors
  
  # Organization
  tags            = ["production", "api"] # Tags
  active          = true                # Enable monitoring
}
```

### TCP Monitor Example

```hcl
resource "uptimekuma_monitor" "database" {
  name     = "Database"
  type     = "tcp"
  hostname = "db.example.com"
  port     = 5432
  
  interval = 30
  timeout  = 10
  
  tags   = ["database", "internal"]
  active = true
}
```

## Best Practices

1. **Sequential Execution**: Always use `-parallelism=1` to avoid WebSocket connection issues
2. **Sensitive Data**: Use Terraform variables for passwords and sensitive configuration
3. **Tags**: Use consistent tagging for organization and filtering
4. **Intervals**: Choose appropriate check intervals based on service criticality
5. **Timeouts**: Set realistic timeouts based on expected response times

## Troubleshooting

### WebSocket Concurrency Issues
If you encounter WebSocket connection errors, ensure you're using `-parallelism=1`:

```bash
tofu apply -parallelism=1
```

### Authentication Errors
Verify your Uptime Kuma admin credentials and server URL in the provider configuration.

### Monitor Creation Fails
Check that the target URL/hostname is accessible from your Uptime Kuma server.

## Variable-Driven Configuration

For managing multiple monitors, see the `test/` directory for an example using `terraform.tfvars` with a list of hosts.
