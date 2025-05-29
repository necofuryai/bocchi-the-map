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
      
      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
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

# Outputs
output "api_url" {
  value = google_cloud_run_v2_service.api.uri
}