# ==============================================================================
# Cloud Run
# ==============================================================================

# Artifact Registry リポジトリ
resource "google_artifact_registry_repository" "gomethod" {
  location      = "asia-northeast1"
  repository_id = "gomethod"
  format        = "DOCKER"
}

# Cloud Run サービス
resource "google_cloud_run_v2_service" "gomethod" {
  name     = "gomethod"
  location = "asia-northeast1"

  template {
    scaling {
      min_instance_count = 0
      max_instance_count = 1
    }

    containers {
      image = "asia-northeast1-docker.pkg.dev/gomethod/gomethod/app:latest"

      ports {
        container_port = 8080
      }

      resources {
        cpu_idle = true # リクエスト中のみCPU割当（コスト最適化）
        limits = {
          cpu    = "1"
          memory = "256Mi"
        }
      }

      # --- アプリ設定 ---
      env {
        name  = "GIN_MODE"
        value = "release"
      }

      # --- データベース ---
      env {
        name  = "GMETHOD_DB_HOST"
        value = "/cloudsql/gomethod:asia-northeast1:gomethod"
      }
      env {
        name  = "GMETHOD_DB_NAME"
        value = "gomethod"
      }
      env {
        name  = "GMETHOD_DB_USERNAME"
        value = "gomethod"
      }
      env {
        name = "GMETHOD_DB_PASSWORD"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.db_password.secret_id
            version = "latest"
          }
        }
      }

      # --- LINE ---
      env {
        name = "LINE_CHANNEL_TOKEN"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.line_channel_token.secret_id
            version = "latest"
          }
        }
      }
      env {
        name = "LINE_CHANNEL_SECRET"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.line_channel_secret.secret_id
            version = "latest"
          }
        }
      }
    }

    # Cloud SQL 接続
    volumes {
      name = "cloudsql"
      cloud_sql_instance {
        instances = ["gomethod:asia-northeast1:gomethod"]
      }
    }
  }

  depends_on = [
    google_artifact_registry_repository.gomethod,
  ]
}

# Cloud Run を公開（LINE Webhook 用）
resource "google_cloud_run_v2_service_iam_member" "public" {
  project  = "gomethod"
  location = "asia-northeast1"
  name     = google_cloud_run_v2_service.gomethod.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
