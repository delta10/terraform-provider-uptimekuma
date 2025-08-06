# Variables for Uptime Kuma monitoring configuration

variable "monitored_hosts" {
  description = "List of hosts to monitor with Uptime Kuma"
  type = list(object({
    name     = string
    url      = string
    interval = optional(number, 60)   # Default: check every 60 seconds
    timeout  = optional(number, 30)   # Default: 30 second timeout
    active   = optional(bool, true)   # Default: monitor is active
  }))
  default = []
  
  validation {
    condition = alltrue([
      for host in var.monitored_hosts : 
      can(regex("^https?://", host.url))
    ])
    error_message = "All URLs must start with http:// or https://"
  }
  
  validation {
    condition = alltrue([
      for host in var.monitored_hosts : 
      host.interval >= 20
    ])
    error_message = "Interval must be at least 20 seconds."
  }
}

# Telegram bot token (get from @BotFather on Telegram)
variable "telegram_bot_token" {
  description = "Telegram bot token for notifications"
  type        = string
  sensitive   = true
}

# Telegram chat ID (your user ID or group chat ID)
variable "telegram_chat_id" {
  description = "Telegram chat ID to send notifications to"
  type        = string
}

# Enable/disable the Telegram notification
variable "telegram_notification_enabled" {
  description = "Whether to enable the Telegram notification as default"
  type        = bool
  default     = true
}