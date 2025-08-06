# Uptime Kuma Connection
variable "uptime_kuma_url" {
  description = "The URL of your Uptime Kuma server"
  type        = string
}

variable "uptime_kuma_username" {
  description = "Username for Uptime Kuma"
  type        = string
}

variable "uptime_kuma_password" {
  description = "Password for Uptime Kuma"
  type        = string
  sensitive   = true
}

# Discord
variable "discord_webhook_url" {
  description = "Discord webhook URL"
  type        = string
  sensitive   = true
  default     = ""
}

# Slack
variable "slack_webhook_url" {
  description = "Slack webhook URL"
  type        = string
  sensitive   = true
  default     = ""
}

# Email
variable "email_username" {
  description = "Email username"
  type        = string
  default     = ""
}

variable "email_password" {
  description = "Email password"
  type        = string
  sensitive   = true
  default     = ""
}

variable "email_from" {
  description = "Email from address"
  type        = string
  default     = ""
}

variable "email_to" {
  description = "Email to address"
  type        = string
  default     = ""
}

# Webhook
variable "webhook_url" {
  description = "Generic webhook URL"
  type        = string
  sensitive   = true
  default     = ""
}
