# ==============================================================================
# Cloud SQL for MySQL
# ==============================================================================

# Cloud SQL インスタンス
resource "google_sql_database_instance" "gomethod" {
  name             = "gomethod"
  database_version = "MYSQL_8_0"
  region           = "asia-northeast1"

  deletion_protection = true

  settings {
    tier              = "db-f1-micro"
    availability_type = "ZONAL"
    disk_size         = 10
    disk_type         = "PD_SSD"
    disk_autoresize   = true

    ip_configuration {
      ipv4_enabled = true
      ssl_mode     = "ENCRYPTED_ONLY"
    }

    backup_configuration {
      enabled            = true
      binary_log_enabled = true
      start_time         = "18:00" # UTC 18:00 = JST 03:00
      backup_retention_settings {
        retained_backups = 7
      }
    }

    maintenance_window {
      day          = 7  # 日曜
      hour         = 18 # UTC 18:00 = JST 03:00
      update_track = "stable"
    }
  }
}

# データベース
resource "google_sql_database" "gomethod" {
  name      = "gomethod"
  instance  = google_sql_database_instance.gomethod.name
  charset   = "utf8mb4"
  collation = "utf8mb4_general_ci"
}

# データベースユーザー
resource "google_sql_user" "gomethod" {
  name     = "gomethod"
  instance = google_sql_database_instance.gomethod.name
  host     = "%"
}
