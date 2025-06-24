# Enable Secret Manager API
resource "google_project_service" "secret_manager_api" {
  service = "secretmanager.googleapis.com"
  disable_on_destroy = false
}

# Secret for New Relic License Key
resource "google_secret_manager_secret" "new_relic_license_key" {
  secret_id = "new-relic-license-key-${var.environment}"
  
  replication {
    auto {}
  }

  depends_on = [google_project_service.secret_manager_api]
}

# Secret for Sentry DSN
resource "google_secret_manager_secret" "sentry_dsn" {
  secret_id = "sentry-dsn-${var.environment}"
  
  replication {
    auto {}
  }

  depends_on = [google_project_service.secret_manager_api]
}

# Secret for TiDB Password
resource "google_secret_manager_secret" "tidb_password" {
  secret_id = "tidb-password-${var.environment}"

  replication {
    auto {}
  }

  depends_on = [google_project_service.secret_manager_api]
}

# IAM binding for Cloud Run service account to access secrets
resource "google_secret_manager_secret_iam_member" "new_relic_access" {
  secret_id = google_secret_manager_secret.new_relic_license_key.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.cloud_run_service_account.email}"
}

resource "google_secret_manager_secret_iam_member" "sentry_access" {
  secret_id = google_secret_manager_secret.sentry_dsn.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.cloud_run_service_account.email}"
}

resource "google_secret_manager_secret_iam_member" "tidb_access" {
  secret_id = google_secret_manager_secret.tidb_password.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.cloud_run_service_account.email}"
}

# Service account for Cloud Run
resource "google_service_account" "cloud_run_service_account" {
  account_id   = "bocchi-cloud-run-${var.environment}"
  display_name = "Bocchi Cloud Run Service Account (${var.environment})"
  description  = "Service account for Bocchi Cloud Run service in ${var.environment} environment"
}

# Grant necessary roles to the service account
resource "google_project_iam_member" "cloud_run_service_account_logging" {
  project = var.project_id
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.cloud_run_service_account.email}"
}

resource "google_project_iam_member" "cloud_run_service_account_monitoring" {
  project = var.project_id
  role    = "roles/monitoring.metricWriter"
  member  = "serviceAccount:${google_service_account.cloud_run_service_account.email}"
}

resource "google_project_iam_member" "cloud_run_service_account_trace" {
  project = var.project_id
  role    = "roles/cloudtrace.agent"
  member  = "serviceAccount:${google_service_account.cloud_run_service_account.email}"
}

# Output the service account email for use in other resources
output "cloud_run_service_account_email" {
  value = google_service_account.cloud_run_service_account.email
}

# Outputs for secret names (for manual secret value setting)
output "new_relic_secret_name" {
  value = google_secret_manager_secret.new_relic_license_key.name
  description = "Full name of the New Relic license key secret. Use this to set the secret value manually."
}

output "sentry_secret_name" {
  value = google_secret_manager_secret.sentry_dsn.name
  description = "Full name of the Sentry DSN secret. Use this to set the secret value manually."
}

output "tidb_secret_name" {
  value = google_secret_manager_secret.tidb_password.name
  description = "Full name of the TiDB password secret. Use this to set the secret value manually."
}