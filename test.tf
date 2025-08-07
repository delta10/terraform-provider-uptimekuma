terraform {
  required_providers {
    uptimekuma = {
      source = "j0r15.local/provider/uptimekuma"
    }
  }
}

provider "uptimekuma" {
  server_url = "http://localhost:3001"
  username   = "admin"
  password   = "cF96H*L9LA3*HiWhx"
}

# Test Discord notification
resource "uptimekuma_notification" "test_discord" {
  name = "Test Discord"
  type = "discord"
  config = {
    discordWebhookUrl = "https://discord.com/api/webhooks/test"
  }
}
