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

# Simple Discord notification example
resource "uptimekuma_notification" "discord_simple" {
  name = "Discord Test"
  type = "Discord"
  config = {
    discordWebhookUrl = var.discord_webhook_url
    discordUsername   = "Uptime Kuma Bot"
  }
}

# Simple Slack notification example
resource "uptimekuma_notification" "slack_simple" {
  name = "Slack Test"
  type = "slack"
  config = {
    slackwebhookURL = var.slack_webhook_url
    slackchannel    = "#general"
    slackusername   = "Uptime Kuma"
  }
}

# Basic email notification
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

# Generic webhook notification
resource "uptimekuma_notification" "webhook_simple" {
  name = "Webhook Test"
  type = "webhook"
  config = {
    webhookURL = var.webhook_url
    httpMethod = "POST"
  }
}
