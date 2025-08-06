# Variables for Uptime Kuma monitoring configuration

# Authentication variables
variable "uptime_kuma_url" {
  description = "URL of the Uptime Kuma server"
  type        = string
  default     = "http://localhost:3001"
  
  validation {
    condition     = can(regex("^https?://", var.uptime_kuma_url))
    error_message = "URL must start with http:// or https://"
  }
}

variable "uptime_kuma_username" {
  description = "Username for Uptime Kuma authentication"
  type        = string
  
  validation {
    condition     = length(var.uptime_kuma_username) > 0
    error_message = "Username cannot be empty."
  }
}

variable "uptime_kuma_password" {
  description = "Password for Uptime Kuma authentication"
  type        = string
  sensitive   = true  # Mark as sensitive to prevent logging
  
  validation {
    condition     = length(var.uptime_kuma_password) > 0
    error_message = "Password cannot be empty."
  }
}

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
