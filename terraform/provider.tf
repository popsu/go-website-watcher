terraform {
  required_providers {
    aiven = {
      source  = "aiven/aiven"
      version = "2.1.13"
    }
  }
}

provider "aiven" {
  api_token = var.aiven_api_token
}
