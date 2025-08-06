# Multiple Hosts Monitoring - Secure Configuration

This example demonstrates how to securely manage Uptime Kuma credentials and monitor multiple hosts using Terraform variables.

## Security Features

🔒 **Separated Credentials**: Authentication details are in a separate `auth.tfvars` file  
🔒 **Sensitive Variables**: Password is marked as sensitive to prevent logging  
🔒 **Git Protection**: `.gitignore` prevents committing sensitive files  
🔒 **Environment Variables**: Alternative secure authentication methods  

## Files Structure

```
├── main.tf              # Main Terraform configuration
├── variables.tf         # Variable definitions with validation
├── terraform.tfvars     # Non-sensitive configuration
├── auth.tfvars          # 🔒 Sensitive authentication (don't commit!)
├── outputs.tf           # Output definitions
└── .gitignore           # Protects sensitive files
```

## Setup Instructions

### 1. Configure Authentication

Choose one of these secure methods:

#### Option A: Separate Auth File (Recommended)
```bash
# Edit auth.tfvars with your credentials
cat > auth.tfvars << EOF
uptime_kuma_username = "your-username"
uptime_kuma_password = "your-password"
EOF
```

#### Option B: Environment Variables
```bash
export TF_VAR_uptime_kuma_username="your-username"
export TF_VAR_uptime_kuma_password="your-password"
export TF_VAR_uptime_kuma_url="http://localhost:3001"
```

#### Option C: Interactive Input
Leave auth.tfvars empty and Terraform will prompt for credentials.

### 2. Configure Monitors

Edit `terraform.tfvars` to add/modify hosts:

```hcl
# Server configuration
uptime_kuma_url = "http://localhost:3001"

# Hosts to monitor
monitored_hosts = [
  {
    name     = "Production API"
    url      = "https://api.mycompany.com"
    interval = 60
    timeout  = 30
    active   = true
  },
  {
    name     = "Documentation"
    url      = "https://docs.mycompany.com"
    interval = 300
    timeout  = 45
    active   = true
  }
]
```

### 3. Deploy

```bash
# Initialize
tofu init

# Plan with auth file
tofu plan -var-file="auth.tfvars"

# Apply with parallelism=1 to avoid WebSocket issues
tofu apply -var-file="auth.tfvars" -parallelism=1
```

## Security Best Practices

### ✅ Do's
- ✅ Use separate auth files for credentials
- ✅ Add auth files to `.gitignore`
- ✅ Use environment variables in CI/CD
- ✅ Regularly rotate passwords
- ✅ Use strong, unique passwords
- ✅ Review who has access to credentials

### ❌ Don'ts
- ❌ Commit passwords to version control
- ❌ Hardcode credentials in main.tf
- ❌ Share auth.tfvars files via chat/email
- ❌ Use default/weak passwords
- ❌ Store credentials in plain text logs

## Usage Examples

### Development Environment
```bash
# Use local auth file
tofu apply -var-file="auth.tfvars" -parallelism=1
```

### Production/CI Environment
```bash
# Use environment variables
export TF_VAR_uptime_kuma_username="$PROD_USERNAME"
export TF_VAR_uptime_kuma_password="$PROD_PASSWORD"
tofu apply -parallelism=1
```

### Team Development
```bash
# Each team member has their own auth.tfvars
cp auth.tfvars.example auth.tfvars
# Edit auth.tfvars with personal credentials
```

## Variable Reference

### Required Variables
- `uptime_kuma_username` - Admin username (sensitive)
- `uptime_kuma_password` - Admin password (sensitive)

### Optional Variables
- `uptime_kuma_url` - Server URL (default: "http://localhost:3001")
- `monitored_hosts` - List of hosts to monitor (default: [])

### Host Configuration
Each host in `monitored_hosts` supports:
- `name` - Display name (required)
- `url` - URL to monitor (required)
- `interval` - Check interval in seconds (default: 60)
- `timeout` - Request timeout in seconds (default: 30)
- `active` - Enable monitoring (default: true)

## Troubleshooting

### Authentication Errors
1. Verify credentials in auth.tfvars
2. Check Uptime Kuma server is accessible
3. Ensure admin user exists and has correct permissions

### WebSocket Issues
Always use `-parallelism=1` to avoid concurrent WebSocket connection issues:
```bash
tofu apply -var-file="auth.tfvars" -parallelism=1
```

### Missing Variables
If you get variable errors:
```bash
# Check which method you're using
tofu plan -var-file="auth.tfvars"  # File method
# OR
echo $TF_VAR_uptime_kuma_username  # Environment method
```

## Github Release

# Generate GPG key
gpg --batch --full-generate-key <<EOF
%no-protection
Key-Type: 1
Key-Length: 4096
Subkey-Type: 1
Subkey-Length: 4096
Expire-Date: 0
Name-Comment: terraform-provider-uptimekuma
Name-Real: Your Name
Name-Email: your.email@example.com
EOF

# Export keys
gpg --armor --export-secret-keys your.email@example.com > private.key
gpg --armor --export your.email@example.com > public.key

Add to GitHub Secrets:
GPG_PRIVATE_KEY: Content of private.key
PASSPHRASE: Your GPG passphrase (if any)