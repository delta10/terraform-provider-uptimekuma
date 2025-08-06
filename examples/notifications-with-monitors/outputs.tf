output "notification_ids" {
  description = "Map of notification names to their IDs"
  value = {
    discord    = uptimekuma_notification.discord.id
    slack      = uptimekuma_notification.slack.id
    email      = uptimekuma_notification.email.id
    teams      = uptimekuma_notification.teams.id
    pagerduty  = uptimekuma_notification.pagerduty.id
  }
}

output "monitor_ids" {
  description = "Map of monitor names to their IDs"
  value = {
    website           = uptimekuma_monitor.website.id
    api              = uptimekuma_monitor.api.id
    database         = uptimekuma_monitor.database.id
    dns              = uptimekuma_monitor.dns.id
    authenticated_api = uptimekuma_monitor.authenticated_api.id
  }
}

output "notification_assignments" {
  description = "Monitor to notification assignments"
  value = {
    website = {
      monitor_id = uptimekuma_monitor.website.id
      notifications = [
        uptimekuma_notification.discord.id,
        uptimekuma_notification.slack.id
      ]
    }
    api = {
      monitor_id = uptimekuma_monitor.api.id
      notifications = [
        uptimekuma_notification.discord.id,
        uptimekuma_notification.slack.id,
        uptimekuma_notification.email.id,
        uptimekuma_notification.teams.id,
        uptimekuma_notification.pagerduty.id
      ]
    }
    database = {
      monitor_id = uptimekuma_monitor.database.id
      notifications = [
        uptimekuma_notification.email.id,
        uptimekuma_notification.pagerduty.id
      ]
    }
    dns = {
      monitor_id = uptimekuma_monitor.dns.id
      notifications = [
        uptimekuma_notification.slack.id
      ]
    }
    authenticated_api = {
      monitor_id = uptimekuma_monitor.authenticated_api.id
      notifications = [
        uptimekuma_notification.discord.id,
        uptimekuma_notification.email.id
      ]
    }
  }
}
