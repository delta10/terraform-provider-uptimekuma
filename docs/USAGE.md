# Getting Started with the Uptime Kuma Terraform Provider

This guide shows how to use the Terraform provider to manage Uptime Kuma monitors as infrastructure-as-code, replacing your Python script approach.

## Prerequisites

1. **Uptime Kuma Instance**: Running on `http://localhost:3001` (or modify the URL)
2. **OpenTofu**: Version 1.6 or later (or Terraform 1.0+)
3. **Go**: Version 1.21 or later (for building from source)

## Installation

### Option 1: Build from Source

```bash
git clone https://github.com/j0r15/terraform-provider-uptimekuma
cd terraform-provider-uptimekuma
make build
```

### Option 2: Install Locally

```bash
make install
```

This copies the provider to your local Terraform plugins directory.

## Configuration

### Environment Variables

Set these environment variables to avoid hardcoding credentials:

```bash
export UPTIMEKUMA_URL="http://localhost:3001"
export UPTIMEKUMA_USERNAME="admin" 
export UPTIMEKUMA_PASSWORD="cF96H*L9LA3*HiWhx"
```

### Provider Configuration

```hcl
terraform {
  required_providers {
    uptimekuma = {
      source = "j0r15.local/provider/uptimekuma"  # For local build
      version = "~> 1.0"
    }
  }
}

provider "uptimekuma" {
  # Configuration will be read from environment variables
  # Or explicitly set:
  # url      = "http://localhost:3001"
  # username = "admin"
  # password = "your-password"
}
```

## Converting Your Python Script

Your original Python script:

```python
api.add_monitor(
    type=MonitorType.HTTP,
    name="Google",
    url="https://google.com",
    id=1
)
```

Becomes this Terraform configuration:

```hcl
resource "uptimekuma_monitor" "google" {
  name = "Google"
  type = "http"
  url  = "https://google.com"
  
  interval = 60
  timeout  = 30
  active   = true
}
```

## Common Use Cases

### 1. HTTP Monitor with Basic Auth

```hcl
resource "uptimekuma_monitor" "secure_api" {
  name = "Secure API"
  type = "http"
  url  = "https://api.example.com/health"
  
  basic_auth_user = "monitoring"
  basic_auth_pass = var.api_password
  
  interval = 120
  timeout  = 30
}
```

### 2. TCP Monitor

```hcl
resource "uptimekuma_monitor" "database" {
  name     = "PostgreSQL"
  type     = "tcp"
  hostname = "db.example.com"
  port     = 5432
  
  interval = 30
  timeout  = 10
}
```

### 3. Monitor with Custom Status Codes

```hcl
resource "uptimekuma_monitor" "api_endpoint" {
  name = "API Health Check"
  type = "http"
  url  = "https://api.example.com/health"
  
  accepted_status_codes = ["200", "201", "202"]
  follow_redirect       = true
  max_redirects        = 5
}
```

### 4. Monitor with Tags

```hcl
resource "uptimekuma_monitor" "production_service" {
  name = "Production Service"
  type = "http"
  url  = "https://prod.example.com"
  
  tags = ["production", "critical", "web"]
}
```

## Usage Commands

### Initialize and Apply

```bash
# Initialize Terraform
terraform init

# Plan changes
terraform plan

# Apply changes
terraform apply

# Show current state
terraform show

# Import existing monitor
terraform import uptimekuma_monitor.existing 1
```

### Manage State

```bash
# List resources
terraform state list

# Show specific resource
terraform state show uptimekuma_monitor.google

# Remove from state (without deleting)
terraform state rm uptimekuma_monitor.google
```

## Data Sources

Query existing monitors:

```hcl
data "uptimekuma_monitor" "existing" {
  id = "1"
}

output "existing_monitor_name" {
  value = data.uptimekuma_monitor.existing.name
}
```

## Troubleshooting

### Common Issues

1. **Version Parsing Error**: If you see version parsing errors with nightly builds, ensure you're using a stable Uptime Kuma version or patch the version detection.

2. **Authentication Failures**: Verify your credentials and that the Uptime Kuma instance is accessible.

3. **Provider Not Found**: Make sure you've installed the provider locally or use the correct source in your configuration.

### Debugging

Enable detailed logging:

```bash
export TF_LOG=DEBUG
terraform apply
```

### API Verification

Test the API manually:

```bash
curl -X POST http://localhost:3001/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your-password"}'
```

## Advantages over Python Script

1. **State Management**: Terraform tracks monitor state and handles updates/deletes
2. **Idempotency**: Safe to run multiple times
3. **Plan/Apply Workflow**: Preview changes before applying
4. **Integration**: Works with other Terraform resources
5. **Version Control**: Infrastructure as code with Git
6. **Rollback**: Easy to revert changes

## Next Steps

1. **Import Existing Monitors**: Use `terraform import` for existing monitors
2. **Organize with Modules**: Create reusable monitor modules
3. **CI/CD Integration**: Automate with GitHub Actions or similar
4. **Monitoring as Code**: Manage all infrastructure monitoring in Terraform

## Support

- Check the [examples/](./examples/) directory for more configurations
- Review the source code in [internal/provider/](./internal/provider/)
- Open issues for bugs or feature requests
