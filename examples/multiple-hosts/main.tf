terraform {
  required_providers {
    uptimekuma = {
      source = "j0r15.local/provider/uptimekuma"
      version = "~> 1.0.0"
    }
  }
}

provider "uptimekuma" {
  url      = var.uptime_kuma_url
  username = var.uptime_kuma_username
  password = var.uptime_kuma_password
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
}
