terraform {
  required_providers {
    uptimekuma = {
      source = "j0r15.local/provider/uptimekuma"
      version = "~> 1.0.0"
    }
  }
}

provider "uptimekuma" {
  url      = "http://localhost:3001"
  username = "admin"
  password = "cF96H*L9LA3*HiWhx"
}

# Telegram notification for all hosts
resource "uptimekuma_notification" "telegram" {
  name = "Telegram Alerts"
  type = "telegram"
  config = {
    telegramBotToken = var.telegram_bot_token
    telegramChatID   = var.telegram_chat_id
  }
  
  # Enable the notification
  is_default = var.telegram_notification_enabled
}

# Create monitors for each host in the list
resource "uptimekuma_monitor" "hosts" {
  for_each = { for host in var.monitored_hosts : host.name => host }
  
  name     = each.value.name
  type     = "http"
  url      = each.value.url
  interval = each.value.interval
  timeout  = each.value.timeout
  active   = each.value.active
  
  # Apply Telegram notification to all hosts
  notification_id_list = [
    uptimekuma_notification.telegram.id
  ]
}