# ==============================================================================
# Cloud Scheduler
# ==============================================================================
#
# バッチ認証は BATCH_AUTH_TOKEN を Bearer トークンとして送信する。
# Cloud Run の URL はデプロイ後に設定する必要がある。
#

locals {
  cloud_run_url = google_cloud_run_v2_service.gomethod.uri
}

# サービスアカウント（Cloud Scheduler → Cloud Run 呼び出し用）
resource "google_service_account" "scheduler" {
  account_id   = "scheduler-sa"
  display_name = "Cloud Scheduler Service Account"
}

resource "google_cloud_run_v2_service_iam_member" "scheduler_invoker" {
  project  = "gomethod"
  location = "asia-northeast1"
  name     = google_cloud_run_v2_service.gomethod.name
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.scheduler.email}"
}

# --- 新月・満月の前日通知 (毎日 20:00 JST) ---
resource "google_cloud_scheduler_job" "send_moon_message_tomorrow" {
  name      = "send-moon-message-tomorrow"
  region    = "asia-northeast1"
  schedule  = "0 20 * * *"
  time_zone = "Asia/Tokyo"

  http_target {
    http_method = "POST"
    uri         = "${local.cloud_run_url}/batch/send_moon_message_tomorrow"

    oidc_token {
      service_account_email = google_service_account.scheduler.email
    }
  }
}

# --- 新月・満月の当日通知 (毎日 08:00 JST) ---
resource "google_cloud_scheduler_job" "send_moon_message_today" {
  name      = "send-moon-message-today"
  region    = "asia-northeast1"
  schedule  = "0 8 * * *"
  time_zone = "Asia/Tokyo"

  http_target {
    http_method = "POST"
    uri         = "${local.cloud_run_url}/batch/send_moon_message_today"

    oidc_token {
      service_account_email = google_service_account.scheduler.email
    }
  }
}

# --- 今日のGメッセージ (毎日 07:00 JST) ---
resource "google_cloud_scheduler_job" "send_daily_g_message" {
  name      = "send-daily-g-message"
  region    = "asia-northeast1"
  schedule  = "0 7 * * *"
  time_zone = "Asia/Tokyo"

  http_target {
    http_method = "POST"
    uri         = "${local.cloud_run_url}/batch/send_daily_g_message"

    oidc_token {
      service_account_email = google_service_account.scheduler.email
    }
  }
}

# --- サタデー動画 (毎週土曜 09:00 JST) ---
resource "google_cloud_scheduler_job" "send_weekly_g_message" {
  name      = "send-weekly-g-message"
  region    = "asia-northeast1"
  schedule  = "0 9 * * 6"
  time_zone = "Asia/Tokyo"

  http_target {
    http_method = "POST"
    uri         = "${local.cloud_run_url}/batch/send_weekly_g_message"

    oidc_token {
      service_account_email = google_service_account.scheduler.email
    }
  }
}

# --- サンデーブログ (毎週日曜 09:00 JST) ---
resource "google_cloud_scheduler_job" "send_weekly_blog_g_message" {
  name      = "send-weekly-blog-g-message"
  region    = "asia-northeast1"
  schedule  = "0 9 * * 0"
  time_zone = "Asia/Tokyo"

  http_target {
    http_method = "POST"
    uri         = "${local.cloud_run_url}/batch/send_weekly_blog_g_message"

    oidc_token {
      service_account_email = google_service_account.scheduler.email
    }
  }
}

# --- 体験談 (毎週火・木 12:00 JST) ---
resource "google_cloud_scheduler_job" "send_experience_g_message" {
  name      = "send-experience-g-message"
  region    = "asia-northeast1"
  schedule  = "0 12 * * 2,4"
  time_zone = "Asia/Tokyo"

  http_target {
    http_method = "POST"
    uri         = "${local.cloud_run_url}/batch/send_experience_g_message"

    oidc_token {
      service_account_email = google_service_account.scheduler.email
    }
  }
}

# --- 月2回のお知らせ (1日と15日 16:00 JST) ---
resource "google_cloud_scheduler_job" "send_notice" {
  name      = "send-notice"
  region    = "asia-northeast1"
  schedule  = "0 16 1,15 * *"
  time_zone = "Asia/Tokyo"

  http_target {
    http_method = "POST"
    uri         = "${local.cloud_run_url}/batch/send_notice"

    oidc_token {
      service_account_email = google_service_account.scheduler.email
    }
  }
}
