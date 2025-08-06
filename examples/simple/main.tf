terraform {
  required_providers {
    uptimekuma = {
      source = "j0r15.local/provider/uptimekuma"
      version = "~> 1.0.0"
    }
  }
}

provider "uptimekuma" {
  url      = "http://localhost:3001"
  username = "admin"
  password = "cF96H*L9LA3*HiWhx"
}

# HTTP Monitor example
resource "uptimekuma_monitor" "google" {
  name = "Google"
  type = "http"
  url  = "https://google.com"
  
  interval = 60
  timeout  = 30
  
  tags = ["production", "external"]
  active = true
  
  follow_redirect = true
  max_redirects   = 10
  
  accepted_status_codes = ["200", "301", "302"]
}

# TCP Monitor example
resource "uptimekuma_monitor" "database" {
  name = "Database"
  type = "tcp"
  hostname = "db.example.com"
  port     = 5432
  
  interval = 30
  timeout  = 10
  
  tags = ["internal", "database"]
  active = true
}

# HTTP Monitor with Basic Auth
resource "uptimekuma_monitor" "secure_api" {
  name = "Secure API"
  type = "http"
  url  = "https://api.example.com/health"
  
  interval = 120
  timeout  = 45
  
  http_method = "GET"
  
  basic_auth_user = "monitoring"
  basic_auth_pass = "secret123"
  
  ignore_tls = false
  
  tags = ["api", "secure"]
  active = true
}
