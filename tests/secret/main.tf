terraform {
  required_providers {
    devin = {
      source  = "hirosi1900day/devin"
      version = "1.0.0"
    }
  }
}

provider "devin" {
  api_key = "test_api_key"
  org_id  = "org-test"
}

# ===================== Secret Resources =====================

resource "devin_secret" "database_url" {
  name  = "DATABASE_URL"
  value = "postgres://user:pass@localhost:5432/db"
}

resource "devin_secret" "api_key" {
  name  = "EXTERNAL_API_KEY"
  value = "sk-test-key-12345"
}

resource "devin_secret" "webhook_secret" {
  name  = "WEBHOOK_SECRET"
  value = "whsec_test_secret"
}

# ===================== Outputs =====================

output "database_url_id" {
  value = devin_secret.database_url.id
}

output "database_url_name" {
  value = devin_secret.database_url.name
}

output "api_key_id" {
  value = devin_secret.api_key.id
}

output "webhook_secret_id" {
  value = devin_secret.webhook_secret.id
}
