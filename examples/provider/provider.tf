terraform {
  required_providers {
    devin = {
      source  = "hirosi1900day/devin"
      version = "~> 0.0.1"
    }
  }
}

provider "devin" {
  # Configuration options
  api_key = "test-api-key"
}
