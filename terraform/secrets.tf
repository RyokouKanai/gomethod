# ==============================================================================
# Secret Manager
# ==============================================================================
#
# シークレットの「値」は Terraform で管理しない。
# 作成後に以下のコマンドで手動設定：
#   gcloud secrets versions add SECRET_NAME --data-file=-
#


resource "google_secret_manager_secret" "db_password" {
  secret_id = "db_password"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret" "line_channel_token" {
  secret_id = "line_channel_token"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret" "line_channel_secret" {
  secret_id = "line_channel_secret"
  replication {
    auto {}
  }
}

