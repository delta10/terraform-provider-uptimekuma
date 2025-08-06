# Terraform Provider for Uptime Kuma

A Terraform provider for managing Uptime Kuma monitors as infrastructure-as-code.

## Features

- Create, update, and delete monitors (HTTP, TCP, Port, DNS, etc.)
- Create and manage notifications (Discord, Slack, Email, Teams, PagerDuty, etc.)
- Associate notifications with monitors for alerting
- Support for 50+ notification types
- Manage monitor groups and tags
- Handle authentication with Uptime Kuma instances
- Variable-driven configuration for secure credential management

## Installation

### From Source

```bash
git clone https://github.com/j0r15/terraform-provider-uptimekuma
cd terraform-provider-uptimekuma
make build
```

### Using OpenTofu/Terraform

Add the provider to your OpenTofu configuration:

```hcl
terraform {
  required_providers {
    uptimekuma = {
      source = "j0r15/uptimekuma"
      version = "~> 1.0"
    }
  }
}

provider "uptimekuma" {
  url      = "http://localhost:3001"
  username = "admin"
  password = "your-password"
}
```

## Usage

### Creating Notifications

```hcl
# Discord notification
resource "uptimekuma_notification" "discord" {
  name = "Discord Alerts"
  type = "Discord"
  config = {
    discordWebhookUrl = "https://discord.com/api/webhooks/..."
    discordUsername   = "Uptime Kuma"
  }
}

# Slack notification
resource "uptimekuma_notification" "slack" {
  name = "Slack Alerts"
  type = "slack"
  config = {
    slackwebhookURL = "https://hooks.slack.com/services/..."
    slackchannel    = "#alerts"
    slackusername   = "Uptime Kuma"
  }
}

# Email notification
resource "uptimekuma_notification" "email" {
  name = "Email Alerts"
  type = "smtp"
  config = {
    smtpHost     = "smtp.gmail.com"
    smtpPort     = "587"
    smtpSecure   = "false"
    smtpUsername = "your-email@gmail.com"
    smtpPassword = "your-app-password"
    emailFrom    = "alerts@your-domain.com"
    emailTo      = "admin@your-domain.com"
  }
}
```

### Creating an HTTP Monitor with Notifications

```hcl
resource "uptimekuma_monitor" "google" {
  name = "Google"
  type = "http"
  url  = "https://google.com"
  
  interval = 60
  timeout  = 30
  
  # Associate notifications
  notification_id_list = [
    uptimekuma_notification.discord.id,
    uptimekuma_notification.slack.id,
    uptimekuma_notification.email.id
  ]
  
  tags = ["production", "external"]
}
```

### Creating a TCP Monitor

```hcl
resource "uptimekuma_monitor" "database" {
  name = "Database"
  type = "tcp"
  hostname = "db.example.com"
  port     = 5432
  
  interval = 30
  timeout  = 10
  
  # Critical services get PagerDuty alerts
  notification_id_list = [
    uptimekuma_notification.email.id,
    uptimekuma_notification.pagerduty.id
  ]
}
```

## Supported Notification Types

The provider supports 50+ notification types including:

- **Chat Platforms**: Discord, Slack, Microsoft Teams, Telegram, Mattermost
- **Email**: SMTP, Mailgun, SendGrid, AWS SES
- **Incident Management**: PagerDuty, Opsgenie, VictorOps
- **Mobile Push**: Pushover, Pushbullet, Pushy
- **SMS**: Twilio, Clickatell, SMS Manager
- **Webhooks**: Custom webhooks, Gotify, ntfy
- **And many more...**

See the [examples/notifications-with-monitors](examples/notifications-with-monitors/) directory for comprehensive examples.

## Configuration

The provider supports the following configuration options:

- `server_url` - The URL of your Uptime Kuma instance
- `username` - Username for authentication
- `password` - Password for authentication (can be set via environment variable `UPTIMEKUMA_PASSWORD`)

## Resources

### `uptimekuma_monitor`

Manages Uptime Kuma monitors.

**Arguments:**
- `name` (Required) - The name of the monitor
- `type` (Required) - Monitor type: `http`, `tcp`, `port`, `dns`, etc.
- `url` - URL to monitor (for HTTP monitors)
- `hostname` - Hostname to monitor (for TCP/Port monitors)
- `port` - Port number (for TCP/Port monitors)
- `interval` - Check interval in seconds (default: 60)
- `timeout` - Request timeout in seconds (default: 30)
- `notification_id_list` - List of notification IDs to associate with this monitor
- `tags` - List of tags for organization
- Additional monitor-specific settings...

### `uptimekuma_notification`

Manages Uptime Kuma notifications.

**Arguments:**
- `name` (Required) - The name of the notification
- `type` (Required) - Notification type (Discord, slack, smtp, teams, PagerDuty, etc.)
- `config` (Required) - Map of configuration values specific to the notification type

**Common notification configs:**
- **Discord**: `discordWebhookUrl`, `discordUsername`
- **Slack**: `slackwebhookURL`, `slackchannel`, `slackusername` 
- **SMTP**: `smtpHost`, `smtpPort`, `smtpUsername`, `smtpPassword`, `emailFrom`, `emailTo`
- **Teams**: `webhookUrl`
- **PagerDuty**: `pagerdutyIntegrationKey`, `pagerdutyPriority`

## Development

### Requirements

- Go 1.21+
- OpenTofu 1.6+ or Terraform 1.0+

### Building

```bash
make build
```

### Install locally
```bash
make install
```


### Testing

```bash
go test ./...
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

## License

EUPL License - see LICENSE file for details.
