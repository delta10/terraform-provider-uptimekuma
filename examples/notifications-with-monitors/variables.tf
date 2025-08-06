# Uptime Kuma Server Configuration
variable "uptime_kuma_url" {
  description = "The URL of your Uptime Kuma server"
  type        = string
}

variable "uptime_kuma_username" {
  description = "Username for Uptime Kuma authentication"
  type        = string
}

variable "uptime_kuma_password" {
  description = "Password for Uptime Kuma authentication"
  type        = string
  sensitive   = true
}

# Discord Configuration
variable "discord_webhook_url" {
  description = "Discord webhook URL for notifications"
  type        = string
  sensitive   = true
}

# Slack Configuration
variable "slack_webhook_url" {
  description = "Slack webhook URL for notifications"
  type        = string
  sensitive   = true
}

# Email Configuration
variable "smtp_host" {
  description = "SMTP server hostname"
  type        = string
}

variable "smtp_username" {
  description = "SMTP username"
  type        = string
}

variable "smtp_password" {
  description = "SMTP password"
  type        = string
  sensitive   = true
}

variable "email_from" {
  description = "Email sender address"
  type        = string
}

variable "email_to" {
  description = "Email recipient address"
  type        = string
}

# Microsoft Teams Configuration
variable "teams_webhook_url" {
  description = "Microsoft Teams webhook URL for notifications"
  type        = string
  sensitive   = true
}

# PagerDuty Configuration
variable "pagerduty_integration_key" {
  description = "PagerDuty integration key"
  type        = string
  sensitive   = true
}

# API Authentication
variable "api_username" {
  description = "Username for authenticated API monitoring"
  type        = string
}

variable "api_password" {
  description = "Password for authenticated API monitoring"
  type        = string
  sensitive   = true
}
