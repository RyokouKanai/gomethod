# ==============================================================================
# Gomethod Infrastructure
# ==============================================================================

terraform {
  required_version = ">= 1.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }

  backend "gcs" {
    bucket = "gomethod-tfstate"
    prefix = "gomethod"
  }
}

provider "google" {
  project = "gomethod"
  region  = "asia-northeast1"
}
