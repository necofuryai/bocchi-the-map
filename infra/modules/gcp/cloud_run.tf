# Enable required APIs
resource "google_project_service" "cloud_run_api" {
  service = "run.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "compute_api" {
  service = "compute.googleapis.com"
  disable_on_destroy = false
}

# Cloud Run service for API
resource "google_cloud_run_v2_service" "api" {
  name     = "bocchi-api-${var.environment}"
  location = var.region
  
  template {
    service_account = google_service_account.cloud_run_service_account.email
    
    containers {
      image = var.api_image
      
      ports {
        container_port = 8080
      }
      
      env {
        name  = "ENV"
        value = var.environment
      }
      
      env {
        name  = "TIDB_HOST"
        value = var.tidb_host
      }
      
      env {
        name = "TIDB_PASSWORD"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.tidb_password.id
            version = "latest"
          }
        }
      }
      
      env {
        name = "NEW_RELIC_LICENSE_KEY"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.new_relic_license_key.id
            version = "latest"
          }
        }
      }
      
      env {
        name = "SENTRY_DSN"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.sentry_dsn.id
            version = "latest"
          }
        }
      }
      
      env {
        name  = "LOG_LEVEL"
        value = var.environment == "prod" ? "INFO" : "DEBUG"
      }
      
      env {
        name  = "PORT"
        value = "8080"
      }
      
      env {
        name  = "HOST"
        value = "0.0.0.0"
      }
      
      resources {
        limits = {
          cpu    = "2"
          memory = "1Gi"
        }
        requests = {
          cpu    = "0.5"
          memory = "256Mi"
        }
      }
    }
    
    scaling {
      min_instance_count = var.environment == "prod" ? 1 : 0
      max_instance_count = var.environment == "prod" ? 10 : 3
    }
  }
  
  traffic {
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
  }
  
  depends_on = [google_project_service.cloud_run_api]
}

# Secret for TiDB password
resource "google_secret_manager_secret" "tidb_password" {
  secret_id = "tidb-password-${var.environment}"
  
  replication {
    auto {}
  }
}

# IAM policy to make the service publicly accessible
resource "google_cloud_run_service_iam_member" "public" {
  location = google_cloud_run_v2_service.api.location
  service  = google_cloud_run_v2_service.api.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# Variables
variable "environment" {
  description = "Environment name"
  type        = string
  
  validation {
    condition = contains(["prod", "dev", "staging", "test"], var.environment)
    error_message = "Environment must be one of: prod, dev, staging, test."
  }
}

variable "region" {
  description = "GCP region"
  type        = string
}

variable "api_image" {
  description = "Docker image for API"
  type        = string
  default     = "gcr.io/cloudrun/placeholder"
}

variable "tidb_host" {
  description = "TiDB host"
  type        = string
}

# TODO: project_id variable defined for future use with GCP project-specific resources
variable "project_id" {
  description = "GCP project ID"
  type        = string
}

# Outputs
output "api_url" {
  value = google_cloud_run_v2_service.api.uri
}