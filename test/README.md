# Uptime Kuma Test Configuration

This test configuration creates monitors for multiple hosts with Telegram notifications.

## Setup

1. **Configure Telegram Bot**:
   - Message @BotFather on Telegram to create a new bot
   - Get your bot token
   - Message @userinfobot to get your chat ID

2. **Configure variables**:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   # Edit terraform.tfvars with your actual values
   ```

3. **Deploy**:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

## Configuration

- **Telegram notifications** will be sent to the specified chat for all monitored hosts
- **Monitors** are created for each host in the `monitored_hosts` variable
- **Customizable** intervals, timeouts, and URLs per host

## Example terraform.tfvars

```hcl
telegram_bot_token = "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz"
telegram_chat_id   = "123456789"

monitored_hosts = [
  {
    name     = "My Website"
    url      = "https://mywebsite.com"
    interval = 60
    timeout  = 30
    active   = true
  }
]
```

## Telegram Setup Guide

1. **Create Bot**:
   - Open Telegram and search for @BotFather
   - Send `/newbot` command
   - Follow instructions to create your bot
   - Save the bot token

2. **Get Chat ID**:
   - Search for @userinfobot on Telegram
   - Send any message to get your user ID
   - Use this as your chat_id

3. **Test**:
   - Start a chat with your bot
   - Your bot will send notifications when monitors go down/up
