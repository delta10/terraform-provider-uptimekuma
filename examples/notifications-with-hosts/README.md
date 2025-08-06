# Notifications with Multiple Hosts Example

This example extends the multiple-hosts example to include notification management.

## Files Structure

```
├── main.tf              # Main Terraform configuration with notifications
├── variables.tf         # Variable definitions
├── terraform.tfvars     # Non-sensitive configuration
├── auth.tfvars          # 🔒 Sensitive authentication (don't commit!)
├── notifications.tfvars # 🔒 Sensitive notification config (don't commit!)
├── outputs.tf           # Output definitions
└── .gitignore           # Protects sensitive files
```

## Setup Instructions

### 1. Configure Authentication & Notifications

Create your sensitive configuration files:

```bash
# Create authentication file
cat > auth.tfvars << EOF
uptime_kuma_username = "your-username"
uptime_kuma_password = "your-password"
EOF

# Create notifications configuration file  
cat > notifications.tfvars << EOF
# Discord webhook for team alerts
discord_webhook_url = "https://discord.com/api/webhooks/YOUR/DISCORD/WEBHOOK"

# Slack webhook for monitoring channel
slack_webhook_url = "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"

# Email configuration
smtp_username = "your-email@gmail.com"
smtp_password = "your-app-password"
smtp_from     = "monitoring@yourcompany.com"
alert_email   = "devops@yourcompany.com"

# PagerDuty for critical alerts
pagerduty_integration_key = "your-pagerduty-integration-key"
EOF
```

### 2. Configure Monitors and Hosts

Edit `terraform.tfvars`:

```hcl
# Server configuration
uptime_kuma_url = "http://localhost:3001"

# Hosts to monitor with notification preferences
monitored_hosts = [
  {
    name     = "Production API"
    url      = "https://api.yourcompany.com"
    interval = 60
    timeout  = 30
    active   = true
    critical = true  # Gets PagerDuty alerts
  },
  {
    name     = "Company Website"
    url      = "https://yourcompany.com"
    interval = 300
    timeout  = 45
    active   = true
    critical = false
  },
  {
    name     = "Documentation Site"
    url      = "https://docs.yourcompany.com"
    interval = 600
    timeout  = 30
    active   = true
    critical = false
  },
  {
    name     = "Status Page"
    url      = "https://status.yourcompany.com"
    interval = 180
    timeout  = 30
    active   = true
    critical = false
  }
]
```

### 3. Deploy

```bash
# Initialize
tofu init

# Plan with all config files
tofu plan -var-file="auth.tfvars" -var-file="notifications.tfvars"

# Apply with parallelism=1 to avoid WebSocket issues
tofu apply -var-file="auth.tfvars" -var-file="notifications.tfvars" -parallelism=1
```

## Configuration Files

### main.tf
```hcl
terraform {
  required_providers {
    uptimekuma = {
      source = "github.com/j0r15/terraform-provider-uptimekuma"
    }
  }
}

provider "uptimekuma" {
  url      = var.uptime_kuma_url
  username = var.uptime_kuma_username
  password = var.uptime_kuma_password
}

# Notification Channels
resource "uptimekuma_notification" "discord_general" {
  name = "Discord General Alerts"
  type = "discord"
  
  config = {
    discordWebhookUrl    = var.discord_webhook_url
    discordUsername      = "Uptime Kuma"
    discordPrefixMessage = "🔔 "
  }
}

resource "uptimekuma_notification" "slack_monitoring" {
  name = "Slack Monitoring Channel"
  type = "slack"
  
  config = {
    slackwebhookURL    = var.slack_webhook_url
    slackchannel       = "#monitoring"
    slackusername      = "UptimeKuma"
    slackiconemo       = ":warning:"
    slackchannelnotify = "true"
  }
}

resource "uptimekuma_notification" "email_alerts" {
  name = "Email Alerts"
  type = "smtp"
  
  config = {
    smtpHost      = "smtp.gmail.com"
    smtpPort      = "587"
    smtpSecure    = "false"
    smtpUsername  = var.smtp_username
    smtpPassword  = var.smtp_password
    smtpFrom      = var.smtp_from
    smtpTo        = var.alert_email
    customSubject = "🚨 Uptime Kuma Alert - [STATUS] [NAME]"
  }
}

resource "uptimekuma_notification" "pagerduty_critical" {
  name = "PagerDuty Critical Alerts"
  type = "PagerDuty"
  
  config = {
    pagerdutyIntegrationKey = var.pagerduty_integration_key
    pagerdutyIntegrationUrl = "https://events.pagerduty.com/v2/enqueue"
    pagerdutyPriority       = "high"
    pagerdutyAutoResolve    = "resolve when monitor goes up"
  }
}

# Generic webhook for custom integrations
resource "uptimekuma_notification" "webhook_integration" {
  name = "Custom Webhook"
  type = "webhook"
  
  config = {
    webhookURL         = "https://api.yourcompany.com/monitoring/webhook"
    webhookContentType = "application/json"
    webhookCustomBody = jsonencode({
      "alert_type": "uptime_kuma"
      "status":     "[STATUS]"
      "monitor":    "[NAME]"
      "url":        "[URL]"
      "message":    "[MSG]"
      "timestamp":  "[TIMESTAMP]"
    })
    webhookAdditionalHeaders = jsonencode({
      "Authorization": "Bearer your-webhook-token"
      "X-Source":      "uptime-kuma-terraform"
    })
  }
}

# Monitors with different notification strategies
resource "uptimekuma_monitor" "hosts" {
  for_each = { for host in var.monitored_hosts : host.name => host }

  name     = each.value.name
  type     = "http"
  url      = each.value.url
  interval = each.value.interval
  timeout  = each.value.timeout
  active   = each.value.active

  # Standard notifications for all monitors
  notification_id_list = concat([
    uptimekuma_notification.discord_general.id,
    uptimekuma_notification.slack_monitoring.id,
    uptimekuma_notification.email_alerts.id,
    uptimekuma_notification.webhook_integration.id
  ],
  # Add PagerDuty for critical services
  each.value.critical ? [uptimekuma_notification.pagerduty_critical.id] : []
  )
}

# Create a default notification for new monitors
resource "uptimekuma_notification" "default_discord" {
  name           = "Default Discord (New Monitors)"
  type           = "discord"
  is_default     = true
  apply_existing = false  # Don't apply to existing monitors, just new ones
  
  config = {
    discordWebhookUrl    = var.discord_webhook_url
    discordUsername      = "Uptime Kuma (Default)"
    discordPrefixMessage = "📢 New Monitor Alert: "
  }
}
```

### variables.tf
```hcl
variable "uptime_kuma_url" {
  description = "The URL of the Uptime Kuma instance"
  type        = string
  default     = "http://localhost:3001"
}

variable "uptime_kuma_username" {
  description = "Username for Uptime Kuma authentication"
  type        = string
  sensitive   = true
}

variable "uptime_kuma_password" {
  description = "Password for Uptime Kuma authentication"
  type        = string
  sensitive   = true
}

variable "monitored_hosts" {
  description = "List of hosts to monitor"
  type = list(object({
    name     = string
    url      = string
    interval = number
    timeout  = number
    active   = bool
    critical = bool
  }))
  default = []

  validation {
    condition = alltrue([
      for host in var.monitored_hosts : can(regex("^https?://", host.url))
    ])
    error_message = "All URLs must start with http:// or https://."
  }

  validation {
    condition = alltrue([
      for host in var.monitored_hosts : host.interval >= 20
    ])
    error_message = "Minimum interval is 20 seconds."
  }
}

# Notification Configuration Variables
variable "discord_webhook_url" {
  description = "Discord webhook URL for notifications"
  type        = string
  sensitive   = true
}

variable "slack_webhook_url" {
  description = "Slack webhook URL for notifications"
  type        = string
  sensitive   = true
}

variable "smtp_username" {
  description = "SMTP username for email notifications"
  type        = string
  sensitive   = true
}

variable "smtp_password" {
  description = "SMTP password for email notifications"
  type        = string
  sensitive   = true
}

variable "smtp_from" {
  description = "SMTP from address"
  type        = string
}

variable "alert_email" {
  description = "Email address to send alerts to"
  type        = string
}

variable "pagerduty_integration_key" {
  description = "PagerDuty integration key for critical alerts"
  type        = string
  sensitive   = true
}
```

### outputs.tf
```hcl
output "notification_summary" {
  description = "Summary of created notifications"
  value = {
    discord_general = {
      id   = uptimekuma_notification.discord_general.id
      name = uptimekuma_notification.discord_general.name
      type = uptimekuma_notification.discord_general.type
    }
    slack_monitoring = {
      id   = uptimekuma_notification.slack_monitoring.id
      name = uptimekuma_notification.slack_monitoring.name
      type = uptimekuma_notification.slack_monitoring.type
    }
    email_alerts = {
      id   = uptimekuma_notification.email_alerts.id
      name = uptimekuma_notification.email_alerts.name
      type = uptimekuma_notification.email_alerts.type
    }
    pagerduty_critical = {
      id   = uptimekuma_notification.pagerduty_critical.id
      name = uptimekuma_notification.pagerduty_critical.name
      type = uptimekuma_notification.pagerduty_critical.type
    }
    webhook_integration = {
      id   = uptimekuma_notification.webhook_integration.id
      name = uptimekuma_notification.webhook_integration.name
      type = uptimekuma_notification.webhook_integration.type
    }
  }
}

output "monitor_notification_mapping" {
  description = "Which notifications are assigned to which monitors"
  value = {
    for name, monitor in uptimekuma_monitor.hosts : name => {
      monitor_id = monitor.id
      notifications = monitor.notification_id_list
    }
  }
}

output "critical_monitors" {
  description = "List of monitors that get PagerDuty alerts"
  value = [
    for host in var.monitored_hosts : host.name if host.critical
  ]
}
```

### .gitignore
```gitignore
# Terraform files
*.tfstate
*.tfstate.*
.terraform/
.terraform.lock.hcl

# Sensitive configuration files
auth.tfvars
notifications.tfvars
*.auto.tfvars

# IDE files
.vscode/
.idea/

# Logs
*.log
```

## Usage Examples

### Development Environment
```bash
# Use local auth and notification files
tofu apply -var-file="auth.tfvars" -var-file="notifications.tfvars" -parallelism=1
```

### Production/CI Environment
```bash
# Use environment variables
export TF_VAR_uptime_kuma_username="$PROD_USERNAME"
export TF_VAR_uptime_kuma_password="$PROD_PASSWORD"
export TF_VAR_discord_webhook_url="$DISCORD_WEBHOOK"
export TF_VAR_slack_webhook_url="$SLACK_WEBHOOK"
export TF_VAR_pagerduty_integration_key="$PAGERDUTY_KEY"
# ... other notification variables

tofu apply -parallelism=1
```

### Testing Notifications
```bash
# Apply configuration
tofu apply -var-file="auth.tfvars" -var-file="notifications.tfvars" -parallelism=1

# Test by temporarily making a monitor fail
tofu apply -var='monitored_hosts=[{
  name="Test Failing Service"
  url="https://httpstat.us/500"
  interval=60
  timeout=30
  active=true
  critical=true
}]' -var-file="auth.tfvars" -var-file="notifications.tfvars" -parallelism=1
```

## Security Best Practices

### ✅ Do's
- ✅ Store sensitive notification credentials in separate files
- ✅ Use environment variables in CI/CD pipelines
- ✅ Use app-specific passwords for email notifications
- ✅ Regularly rotate webhook URLs and tokens
- ✅ Test notifications before deploying
- ✅ Use different notification channels for different severity levels

### ❌ Don'ts
- ❌ Commit webhook URLs or API keys to version control
- ❌ Use the same notification channel for all alerts
- ❌ Skip testing notification configurations
- ❌ Use plain passwords instead of app-specific passwords
- ❌ Share notification credentials via chat/email

## Troubleshooting

### Notification Not Working
1. Test the notification configuration manually in Uptime Kuma UI
2. Check webhook URLs and API keys are correct
3. Verify notification is associated with monitors
4. Check monitor is active and has valid configuration

### WebSocket Issues
Always use `-parallelism=1`:
```bash
tofu apply -var-file="auth.tfvars" -var-file="notifications.tfvars" -parallelism=1
```

### Configuration Validation
```bash
# Validate configuration
tofu validate

# Plan to see what will be created
tofu plan -var-file="auth.tfvars" -var-file="notifications.tfvars"
```

## Advanced Usage

### Dynamic Notification Assignment
```hcl
locals {
  # Define notification tiers
  notification_tiers = {
    basic = [
      uptimekuma_notification.discord_general.id,
      uptimekuma_notification.email_alerts.id
    ]
    standard = [
      uptimekuma_notification.discord_general.id,
      uptimekuma_notification.slack_monitoring.id,
      uptimekuma_notification.email_alerts.id
    ]
    critical = [
      uptimekuma_notification.discord_general.id,
      uptimekuma_notification.slack_monitoring.id,
      uptimekuma_notification.email_alerts.id,
      uptimekuma_notification.pagerduty_critical.id
    ]
  }
}

resource "uptimekuma_monitor" "hosts_dynamic" {
  for_each = { for host in var.monitored_hosts : host.name => host }

  name     = each.value.name
  type     = "http"
  url      = each.value.url
  interval = each.value.interval
  timeout  = each.value.timeout
  active   = each.value.active

  # Assign notifications based on criticality
  notification_id_list = each.value.critical ? local.notification_tiers.critical : local.notification_tiers.standard
}
```
