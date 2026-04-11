resource "devin_secret" "example" {
  name  = "DATABASE_URL"
  value = "postgresql://user:pass@host:5432/db"  # Sensitive, ForceNew
}
