# Notifications with Monitors Example

This example demonstrates how to create various types of notifications and associate them with monitors in your Uptime Kuma Terraform provider.

## Features Demonstrated

- **Multiple notification types**: Discord, Slack, Email (SMTP), Microsoft Teams, PagerDuty
- **Monitor-notification associations**: Different monitors with different notification combinations
- **Various monitor types**: HTTP, Port, DNS monitoring
- **Authentication**: Basic auth for API monitoring
- **Variable-driven configuration**: All sensitive data in variables

## Notification Types Included

### Discord
```hcl
resource "uptimekuma_notification" "discord" {
  name = "Discord Alerts"
  type = "Discord"
  config = {
    discordWebhookUrl = var.discord_webhook_url
    discordUsername   = "Uptime Kuma"
  }
}
```

### Slack
```hcl
resource "uptimekuma_notification" "slack" {
  name = "Slack Alerts"
  type = "slack"
  config = {
    slackwebhookURL = var.slack_webhook_url
    slackchannel    = "#alerts"
    slackusername   = "Uptime Kuma"
  }
}
```

### Email (SMTP)
```hcl
resource "uptimekuma_notification" "email" {
  name = "Email Alerts"
  type = "smtp"
  config = {
    smtpHost     = var.smtp_host
    smtpPort     = "587"
    smtpSecure   = "false"
    smtpUsername = var.smtp_username
    smtpPassword = var.smtp_password
    emailFrom    = var.email_from
    emailTo      = var.email_to
  }
}
```

### Microsoft Teams
```hcl
resource "uptimekuma_notification" "teams" {
  name = "Teams Alerts"
  type = "teams"
  config = {
    webhookUrl = var.teams_webhook_url
  }
}
```

### PagerDuty
```hcl
resource "uptimekuma_notification" "pagerduty" {
  name = "PagerDuty Alerts"
  type = "PagerDuty"
  config = {
    pagerdutyIntegrationKey = var.pagerduty_integration_key
    pagerdutyPriority       = "high"
  }
}
```

## Monitor Examples

### Website Monitor with Multiple Notifications
```hcl
resource "uptimekuma_monitor" "website" {
  name     = "Main Website"
  type     = "http"
  url      = "https://example.com"
  interval = 60
  
  notification_id_list = [
    uptimekuma_notification.discord.id,
    uptimekuma_notification.slack.id
  ]
}
```

### API Monitor with All Notifications
```hcl
resource "uptimekuma_monitor" "api" {
  name     = "API Endpoint"
  type     = "http"
  url      = "https://api.example.com/health"
  interval = 30
  timeout  = 10
  
  notification_id_list = [
    uptimekuma_notification.discord.id,
    uptimekuma_notification.slack.id,
    uptimekuma_notification.email.id,
    uptimekuma_notification.teams.id,
    uptimekuma_notification.pagerduty.id
  ]
}
```

### Database Port Monitor
```hcl
resource "uptimekuma_monitor" "database" {
  name     = "Database Server"
  type     = "port"
  hostname = "db.example.com"
  port     = 5432
  interval = 60
  
  notification_id_list = [
    uptimekuma_notification.email.id,
    uptimekuma_notification.pagerduty.id
  ]
}
```

## Setup Instructions

1. **Copy configuration files**:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```

2. **Edit terraform.tfvars** with your actual values:
   - Uptime Kuma server URL and credentials
   - Webhook URLs for Discord, Slack, Teams
   - SMTP settings for email notifications
   - PagerDuty integration key
   - API credentials for authenticated monitoring

3. **Initialize Terraform**:
   ```bash
   terraform init
   ```

4. **Plan the deployment**:
   ```bash
   terraform plan
   ```

5. **Apply the configuration**:
   ```bash
   terraform apply
   ```

## Configuration Notes

### Discord Setup
1. Create a webhook in your Discord server
2. Copy the webhook URL to `discord_webhook_url`

### Slack Setup
1. Create an incoming webhook in your Slack workspace
2. Copy the webhook URL to `slack_webhook_url`

### Email Setup
1. Configure SMTP settings for your email provider
2. For Gmail, use an App Password instead of your regular password

### Microsoft Teams Setup
1. Add an Incoming Webhook connector to your Teams channel
2. Copy the webhook URL to `teams_webhook_url`

### PagerDuty Setup
1. Create a service in PagerDuty
2. Add an integration and copy the integration key

## Outputs

The configuration provides several useful outputs:

- `notification_ids`: Map of notification names to their IDs
- `monitor_ids`: Map of monitor names to their IDs  
- `notification_assignments`: Complete mapping of monitors to their assigned notifications

## Supported Notification Types

This provider supports 50+ notification types. Some popular ones include:

- Discord
- Slack
- Microsoft Teams
- Email (SMTP)
- PagerDuty
- Telegram
- Webhook
- Pushover
- Twilio
- And many more...

See the main documentation for a complete list of supported notification types and their configuration options.
