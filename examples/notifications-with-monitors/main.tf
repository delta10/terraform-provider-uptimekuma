terraform {
  required_providers {
    uptimekuma = {
      source = "local/j0r15/uptimekuma"
    }
  }
}

provider "uptimekuma" {
  server_url = var.uptime_kuma_url
  username   = var.uptime_kuma_username
  password   = var.uptime_kuma_password
}

# Discord notification
resource "uptimekuma_notification" "discord" {
  name = "Discord Alerts"
  type = "Discord"
  config = {
    discordWebhookUrl = var.discord_webhook_url
    discordUsername   = "Uptime Kuma"
  }
}

# Slack notification
resource "uptimekuma_notification" "slack" {
  name = "Slack Alerts"
  type = "slack"
  config = {
    slackwebhookURL = var.slack_webhook_url
    slackchannel    = "#alerts"
    slackusername   = "Uptime Kuma"
  }
}

# Email notification
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

# Microsoft Teams notification
resource "uptimekuma_notification" "teams" {
  name = "Teams Alerts"
  type = "teams"
  config = {
    webhookUrl = var.teams_webhook_url
  }
}

# PagerDuty notification
resource "uptimekuma_notification" "pagerduty" {
  name = "PagerDuty Alerts"
  type = "PagerDuty"
  config = {
    pagerdutyIntegrationKey = var.pagerduty_integration_key
    pagerdutyPriority       = "high"
  }
}

# Website monitor with Discord and Slack notifications
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

# API endpoint monitor with all notifications
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

# Database monitor with email and PagerDuty
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

# DNS monitor with Slack notification
resource "uptimekuma_monitor" "dns" {
  name     = "DNS Resolution"
  type     = "dns"
  hostname = "example.com"
  interval = 120
  
  notification_id_list = [
    uptimekuma_notification.slack.id
  ]
}

# HTTP monitor with custom headers and authentication
resource "uptimekuma_monitor" "authenticated_api" {
  name           = "Authenticated API"
  type           = "http"
  url            = "https://api.example.com/private"
  interval       = 60
  http_method    = "GET"
  basic_auth_user = var.api_username
  basic_auth_pass = var.api_password
  
  notification_id_list = [
    uptimekuma_notification.discord.id,
    uptimekuma_notification.email.id
  ]
}
