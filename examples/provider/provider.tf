terraform {
  required_providers {
    devin = {
      source  = "hirosi1900day/devin"
      version = "~> 0.0.4"
    }
  }
}

provider "devin" {
  # Configuration options
  api_key = "test-api-key"
}
