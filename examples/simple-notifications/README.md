# Simple Notifications Example

This example demonstrates how to create basic notifications without monitors. Perfect for testing notification configurations.

## Quick Start

1. **Copy the example configuration**:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```

2. **Edit terraform.tfvars** with your values:
   ```hcl
   uptime_kuma_url      = "https://your-uptime-kuma.example.com"
   uptime_kuma_username = "admin"
   uptime_kuma_password = "your-password"
   
   # Add any notification you want to test
   discord_webhook_url = "https://discord.com/api/webhooks/your-webhook"
   slack_webhook_url   = "https://hooks.slack.com/services/your-webhook"
   ```

3. **Deploy**:
   ```bash
   terraform init
   terraform apply
   ```

## Available Notifications

### Discord
```hcl
resource "uptimekuma_notification" "discord_simple" {
  name = "Discord Test"
  type = "Discord"
  config = {
    discordWebhookUrl = var.discord_webhook_url
    discordUsername   = "Uptime Kuma Bot"
  }
}
```

### Slack
```hcl
resource "uptimekuma_notification" "slack_simple" {
  name = "Slack Test"
  type = "slack"
  config = {
    slackwebhookURL = var.slack_webhook_url
    slackchannel    = "#general"
    slackusername   = "Uptime Kuma"
  }
}
```

### Email (SMTP)
```hcl
resource "uptimekuma_notification" "email_simple" {
  name = "Email Test"
  type = "smtp"
  config = {
    smtpHost     = "smtp.gmail.com"
    smtpPort     = "587"
    smtpSecure   = "false"
    smtpUsername = var.email_username
    smtpPassword = var.email_password
    emailFrom    = var.email_from
    emailTo      = var.email_to
  }
}
```

### Generic Webhook
```hcl
resource "uptimekuma_notification" "webhook_simple" {
  name = "Webhook Test"
  type = "webhook"
  config = {
    webhookURL = var.webhook_url
    httpMethod = "POST"
  }
}
```

## Testing Notifications

Once created, you can test notifications directly in the Uptime Kuma UI or by associating them with monitors.

## Variables

All notification configurations are optional - leave variables empty for notifications you don't want to create.
