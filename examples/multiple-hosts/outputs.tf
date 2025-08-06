# Output definitions for Uptime Kuma monitors

output "monitor_details" {
  description = "Details of all created monitors"
  value = {
    for name, monitor in uptimekuma_monitor.hosts : name => {
      id   = monitor.id
      name = monitor.name
      url  = monitor.url
    }
  }
}

output "monitor_count" {
  description = "Total number of monitors created"
  value = length(uptimekuma_monitor.hosts)
}

output "active_monitors" {
  description = "List of active monitor names"
  value = [
    for name, monitor in uptimekuma_monitor.hosts : 
    monitor.name if monitor.active
  ]
}
