terraform {
  required_providers {
    devin = {
      source  = "hirosi1900day/devin"
      version = "~> 1.0.0"
    }
  }
}

provider "devin" {
  api_key = "cog_..."  # Service User credential
  org_id  = "org_..."  # Organization ID
}
