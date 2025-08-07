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

# Simple Slack notification example
resource "uptimekuma_notification" "slack" {
  name = "Slack Alerts"
  type = "slack"
  config = {
    slackwebhookURL = var.slack_url
    slackchannel    = "#uptime"
    slackusername   = "Uptime Kuma"
  }
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
  
  # Apply slack notification to all hosts
  notification_id_list = [
    uptimekuma_notification.slack.id
  ]
}