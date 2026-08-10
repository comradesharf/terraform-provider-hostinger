terraform {
  required_providers {
    hostinger = {
      source  = "comradesharf/hostinger"
      version = "~> 0.1.0"
    }
  }
}

provider "hostinger" {
  api_token = var.hostinger_api_token
}
