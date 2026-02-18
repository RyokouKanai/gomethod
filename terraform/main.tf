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

  # TODO: リモートバックエンドを設定する場合はここに追加
  # backend "gcs" {
  #   bucket = "your-terraform-state-bucket"
  #   prefix = "gomethod"
  # }
}

provider "google" {
  project = var.project_id
  region  = var.region
}
