# Notification Examples for Uptime Kuma Terraform Provider

This document shows examples of configuring various notification types using the Uptime Kuma Terraform provider.

## Basic Notification Structure

```hcl
resource "uptimekuma_notification" "example" {
  name          = "My Notification"
  type          = "webhook"  # See supported types below
  is_default    = false      # Apply to new monitors by default
  apply_existing = false     # Apply to existing monitors

  config = {
    # Type-specific configuration parameters
  }
}
```

## Supported Notification Types

### Discord
```hcl
resource "uptimekuma_notification" "discord" {
  name = "Discord Alerts"
  type = "discord"
  
  config = {
    discordUsername    = "Uptime Kuma"
    discordWebhookUrl  = "https://discord.com/api/webhooks/YOUR_WEBHOOK_URL"
    discordPrefixMessage = "🚨 Alert: "
  }
}
```

### Slack
```hcl
resource "uptimekuma_notification" "slack" {
  name = "Slack Alerts"
  type = "slack"
  
  config = {
    slackwebhookURL    = "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"
    slackchannel       = "#alerts"
    slackusername      = "UptimeKuma"
    slackiconemo       = ":warning:"
    slackchannelnotify = "true"
  }
}
```

### Microsoft Teams
```hcl
resource "uptimekuma_notification" "teams" {
  name = "Teams Alerts"
  type = "teams"
  
  config = {
    webhookUrl = "https://outlook.office.com/webhook/YOUR_TEAMS_WEBHOOK"
  }
}
```

### Telegram
```hcl
resource "uptimekuma_notification" "telegram" {
  name = "Telegram Alerts"
  type = "telegram"
  
  config = {
    telegramBotToken        = "YOUR_BOT_TOKEN"
    telegramChatID          = "YOUR_CHAT_ID"
    telegramSendSilently    = "false"
    telegramProtectContent  = "false"
    telegramMessageThreadID = ""
  }
}
```

### Email (SMTP)
```hcl
resource "uptimekuma_notification" "email" {
  name = "Email Alerts"
  type = "smtp"
  
  config = {
    smtpHost             = "smtp.gmail.com"
    smtpPort             = "587"
    smtpSecure           = "false"
    smtpIgnoreTLSError   = "false"
    smtpUsername         = "your-email@gmail.com"
    smtpPassword         = "your-app-password"
    smtpFrom             = "your-email@gmail.com"
    smtpTo               = "alerts@company.com"
    smtpCC               = ""
    smtpBCC              = ""
    customSubject        = "Uptime Kuma Alert"
  }
}
```

### Generic Webhook
```hcl
resource "uptimekuma_notification" "webhook" {
  name = "Custom Webhook"
  type = "webhook"
  
  config = {
    webhookURL              = "https://api.example.com/webhook"
    webhookContentType      = "application/json"
    webhookCustomBody       = jsonencode({
      "text": "Alert: [STATUS] [NAME] is [MSG]"
      "status": "[STATUS]"
      "monitor": "[NAME]"
    })
    webhookAdditionalHeaders = jsonencode({
      "Authorization": "Bearer YOUR_TOKEN"
    })
  }
}
```

### PagerDuty
```hcl
resource "uptimekuma_notification" "pagerduty" {
  name = "PagerDuty Alerts"
  type = "PagerDuty"
  
  config = {
    pagerdutyIntegrationKey = "YOUR_INTEGRATION_KEY"
    pagerdutyIntegrationUrl = "https://events.pagerduty.com/v2/enqueue"
    pagerdutyPriority       = "high"
    pagerdutyAutoResolve    = "resolve when monitor goes up"
  }
}
```

### Pushover
```hcl
resource "uptimekuma_notification" "pushover" {
  name = "Pushover Alerts"
  type = "pushover"
  
  config = {
    pushoveruserkey  = "YOUR_USER_KEY"
    pushoverapptoken = "YOUR_APP_TOKEN"
    pushoversounds   = "pushover"
    pushoverpriority = "0"
    pushovertitle    = "Uptime Kuma"
    pushoverdevice   = ""
    pushoverttl      = "3600"
  }
}
```

### Gotify
```hcl
resource "uptimekuma_notification" "gotify" {
  name = "Gotify Alerts"
  type = "gotify"
  
  config = {
    gotifyserverurl        = "https://gotify.example.com"
    gotifyapplicationToken = "YOUR_APPLICATION_TOKEN"
    gotifyPriority         = "5"
  }
}
```

## Using Notifications with Monitors

To associate notifications with monitors, include the notification IDs in the monitor configuration:

```hcl
# Create notifications
resource "uptimekuma_notification" "discord" {
  name = "Discord Alerts"
  type = "discord"
  
  config = {
    discordWebhookUrl = var.discord_webhook_url
    discordUsername   = "Uptime Kuma"
  }
}

resource "uptimekuma_notification" "email" {
  name = "Email Alerts"
  type = "smtp"
  
  config = {
    smtpHost     = "smtp.gmail.com"
    smtpPort     = "587"
    smtpUsername = var.smtp_username
    smtpPassword = var.smtp_password
    smtpFrom     = var.smtp_from
    smtpTo       = var.alert_email
  }
}

# Create monitor with notifications
resource "uptimekuma_monitor" "website" {
  name     = "Company Website"
  type     = "http"
  url      = "https://example.com"
  interval = 60
  
  # Associate notifications with monitor
  notification_id_list = [
    uptimekuma_notification.discord.id,
    uptimekuma_notification.email.id
  ]
}
```

## Default and Apply Existing

- `is_default = true`: This notification will be automatically enabled for all new monitors
- `apply_existing = true`: When creating the notification, it will be applied to all existing monitors

```hcl
resource "uptimekuma_notification" "default_slack" {
  name           = "Default Slack Alerts"
  type           = "slack"
  is_default     = true    # Apply to all new monitors
  apply_existing = true    # Apply to all existing monitors
  
  config = {
    slackwebhookURL = var.slack_webhook_url
    slackchannel    = "#monitoring"
  }
}
```

## Complete Example

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

# Discord notification
resource "uptimekuma_notification" "discord" {
  name = "Discord Alerts"
  type = "discord"
  
  config = {
    discordWebhookUrl    = var.discord_webhook_url
    discordUsername      = "Uptime Kuma"
    discordPrefixMessage = "🚨 "
  }
}

# Email notification
resource "uptimekuma_notification" "email" {
  name = "Email Alerts"
  type = "smtp"
  
  config = {
    smtpHost     = "smtp.gmail.com"
    smtpPort     = "587"
    smtpUsername = var.smtp_username
    smtpPassword = var.smtp_password
    smtpFrom     = var.smtp_from
    smtpTo       = var.alert_email
    customSubject = "🚨 Uptime Alert"
  }
}

# PagerDuty for critical services
resource "uptimekuma_notification" "pagerduty" {
  name = "PagerDuty Critical"
  type = "PagerDuty"
  
  config = {
    pagerdutyIntegrationKey = var.pagerduty_integration_key
    pagerdutyPriority       = "high"
    pagerdutyAutoResolve    = "resolve when monitor goes up"
  }
}

# Monitors with different notification strategies
resource "uptimekuma_monitor" "website" {
  name     = "Company Website"
  type     = "http"
  url      = "https://example.com"
  interval = 60
  
  notification_id_list = [
    uptimekuma_notification.discord.id,
    uptimekuma_notification.email.id
  ]
}

resource "uptimekuma_monitor" "critical_api" {
  name     = "Critical API"
  type     = "http"
  url      = "https://api.example.com/health"
  interval = 30
  
  notification_id_list = [
    uptimekuma_notification.discord.id,
    uptimekuma_notification.email.id,
    uptimekuma_notification.pagerduty.id  # Critical gets PagerDuty too
  ]
}
```

## Variables Example

```hcl
# variables.tf
variable "uptime_kuma_url" {
  description = "Uptime Kuma server URL"
  type        = string
  default     = "http://localhost:3001"
}

variable "uptime_kuma_username" {
  description = "Uptime Kuma username"
  type        = string
  sensitive   = true
}

variable "uptime_kuma_password" {
  description = "Uptime Kuma password"
  type        = string
  sensitive   = true
}

variable "discord_webhook_url" {
  description = "Discord webhook URL"
  type        = string
  sensitive   = true
}

variable "slack_webhook_url" {
  description = "Slack webhook URL"
  type        = string
  sensitive   = true
}

variable "pagerduty_integration_key" {
  description = "PagerDuty integration key"
  type        = string
  sensitive   = true
}

variable "smtp_username" {
  description = "SMTP username"
  type        = string
  sensitive   = true
}

variable "smtp_password" {
  description = "SMTP password"
  type        = string
  sensitive   = true
}

variable "smtp_from" {
  description = "SMTP from address"
  type        = string
}

variable "alert_email" {
  description = "Email address for alerts"
  type        = string
}
```

## Importing Existing Notifications

You can import existing notifications using their ID:

```bash
terraform import uptimekuma_notification.discord 1
```
