# ☁️ Bocchi The Map Infrastructure

> **Cloud-native infrastructure with multi-provider flexibility** - Terraform, auto-scaling, and zero-downtime deployments

[![Terraform](https://img.shields.io/badge/Terraform-1.5+-7C3AED?style=flat&logo=terraform)](https://terraform.io/)
[![Google Cloud](https://img.shields.io/badge/Google_Cloud-Run-4285F4?style=flat&logo=googlecloud)](https://cloud.google.com/)
[![Vercel](https://img.shields.io/badge/Vercel-000000?style=flat&logo=vercel)](https://vercel.com/)
[![Cloudflare](https://img.shields.io/badge/Cloudflare-R2-F38020?style=flat&logo=cloudflare)](https://cloudflare.com/)
[![TiDB](https://img.shields.io/badge/TiDB-Serverless-FF6B6B?style=flat)](https://tidbcloud.com/)

**Production-ready infrastructure** that scales from zero to millions of users with minimal operational overhead. Built with Infrastructure as Code principles, multi-cloud redundancy, and cost optimization.

## ⚡ Quick Start

```bash
# Prerequisites: Terraform 1.5+, Google Cloud SDK
cd infra
cp environments/dev/terraform.tfvars.example terraform.tfvars
terraform init                 # Initialize providers
terraform plan                 # Preview changes
terraform apply                # Deploy infrastructure 🚀

# Infrastructure ready in ~5 minutes
```

## 🏗️ Architecture Overview

### Multi-Cloud Strategy

```text
🌐 Global Edge Network
├── 📱 Frontend (Vercel) - Global edge network with auto-deployments
├── 🗄️ Static Assets (Cloudflare R2) - PMTiles map storage
└── 🔒 Security (Cloudflare WAF, DDoS protection)

☁️ Compute Layer (Google Cloud)
├── 🚀 API Services (Cloud Run) - Auto-scaling containers
├── 🔐 Secrets (Secret Manager) - Encrypted configuration
└── 📊 Monitoring (Cloud Logging, Monitoring)

💾 Data Layer (TiDB Cloud)
├── 🗃️ Primary Database (Serverless) - Auto-scaling MySQL
├── 📈 Analytics (TiFlash) - OLAP workloads
└── 🔄 Backup & Recovery (Automated)
```

**Why This Architecture?**
- ⚡ **Performance** - Edge caching + regional compute
- 💰 **Cost Efficient** - Pay only for actual usage
- 🛡️ **Resilient** - Multi-provider redundancy
- 📈 **Scalable** - Zero to millions without config changes

## 🛠️ Infrastructure Components

### Terraform Modules

```hcl
# Clean, reusable infrastructure modules
module "gcp_cloud_run" {
  source = "./modules/gcp"
  
  project_id    = var.gcp_project_id
  region        = var.gcp_region
  service_name  = "bocchi-api"
  
  # Auto-scaling configuration
  min_instances = 0
  max_instances = 100
  
  # Performance tuning
  cpu_limit    = "2"
  memory_limit = "2Gi"
}

module "vercel_deployment" {
  source = "./modules/vercel"
  
  team_id     = var.vercel_team_id
  domain_name = "bocchi-map.com"
  
  # Edge optimization
  minify_js  = true
  minify_css = true
  
  # Security headers
  security_headers = true
}
```

### Environment Management

```text
environments/
├── 🧪 dev/                    # Development environment
│   ├── terraform.tfvars       # Environment-specific variables
│   └── backend.tf             # Remote state configuration
├── 🎭 staging/                # Staging environment  
│   ├── terraform.tfvars       # Pre-production testing
│   └── backend.tf             # Isolated state
└── 🚀 prod/                   # Production environment
    ├── terraform.tfvars       # Production configuration
    └── backend.tf             # Production state
```

## 🚀 Cloud Run Configuration

### Auto-Scaling API

```hcl
resource "google_cloud_run_service" "api" {
  name     = "bocchi-api"
  location = var.region

  template {
    metadata {
      annotations = {
        "autoscaling.knative.dev/minScale" = "0"
        "autoscaling.knative.dev/maxScale" = "100"
        "run.googleapis.com/cpu-throttling" = "false"
      }
    }

    spec {
      container_concurrency = 1000
      timeout_seconds       = 300

      containers {
        image = "gcr.io/${var.project_id}/bocchi-api:latest"
        
        resources {
          limits = {
            cpu    = "2"
            memory = "2Gi"
          }
        }

        env {
          name  = "TIDB_HOST"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.tidb_config.secret_id
              key  = "host"
            }
          }
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}
```

**Cloud Run Benefits:**

- **Scale to Zero** - No idle costs during low traffic
- **Fast Cold Starts** - < 500ms container startup
- **Request-based Billing** - Pay per million requests
- **Automatic HTTPs** - SSL certificates managed automatically

## 🌐 Cloudflare Edge Network

### Vercel Deployment

```hcl
resource "vercel_project" "frontend" {
  name      = "bocchi-map"
  team_id   = var.vercel_team_id
  framework = "nextjs"

  git_repository {
    type = "github"
    repo = "necofuryai/bocchi-the-map"
  }

  build_command    = "pnpm build"
  output_directory = "out"
  root_directory   = "web"

  environment = [
    {
      key    = "NEXT_PUBLIC_API_URL"
      value  = "https://api.bocchi-map.com"
      target = ["production"]
    }
  ]
}
```

### R2 Storage for Maps
```hcl
resource "cloudflare_r2_bucket" "map_tiles" {
  account_id = var.cloudflare_account_id
  name       = "bocchi-map-tiles"
  location   = "APAC"  # Optimize for Asia-Pacific users
}

# PMTiles storage for efficient map delivery
resource "cloudflare_worker" "map_worker" {
  name    = "map-tile-worker"
  content = file("${path.module}/workers/map-tiles.js")
}
```

**Vercel Advantages:**

- **Global Edge Network** - 40+ regions worldwide  
- **Automatic Deployments** - GitHub integration with preview deployments
- **Edge Functions** - Serverless compute at the edge
- **Zero Configuration** - Optimized for Next.js with built-in performance monitoring
- **Built-in Analytics** - Core Web Vitals tracking

## 💾 Database Architecture

### TiDB Serverless Configuration

```hcl
# TiDB Cloud configuration via API
resource "tidbcloud_cluster" "main" {
  project_id   = var.tidb_project_id
  name         = "bocchi-${var.environment}"
  cluster_type = "SERVERLESS"
  
  cloud_provider = "GCP"
  region         = "asia-northeast1"
  
  config {
    root_password = var.tidb_root_password
    
    # Auto-scaling configuration
    min_capacity = 0.5  # Minimum RU (Request Units)
    max_capacity = 16   # Maximum RU for cost control
    
    # Automatic backups
    backup_retention_days = 7
  }
}
```

**TiDB Benefits:**

- **MySQL Compatible** - Drop-in replacement with better scaling
- **Serverless Scaling** - Automatically handle traffic spikes
- **HTAP Workloads** - Handle both OLTP and OLAP queries
- **Global Distribution** - Multi-region replication

## 📊 Monitoring & Observability

### Google Cloud Monitoring

```hcl
resource "google_monitoring_dashboard" "api_dashboard" {
  dashboard_json = jsonencode({
    displayName = "Bocchi API Dashboard"
    mosaicLayout = {
      tiles = [
        {
          width  = 6
          height = 4
          widget = {
            title = "Request Latency"
            xyChart = {
              dataSets = [{
                timeSeriesQuery = {
                  timeSeriesFilter = {
                    filter = "resource.type=\"cloud_run_revision\""
                    aggregation = {
                      alignmentPeriod  = "60s"
                      perSeriesAligner = "ALIGN_DELTA"
                    }
                  }
                }
              }]
            }
          }
        }
      ]
    }
  })
}
```

### Health Checks & Alerts
```hcl
resource "google_monitoring_alert_policy" "high_latency" {
  display_name = "High API Latency"
  combiner     = "OR"
  
  conditions {
    display_name = "Request latency > 2s"
    
    condition_threshold {
      filter         = "resource.type=\"cloud_run_revision\""
      duration       = "300s"
      comparison     = "COMPARISON_GT"
      threshold_value = 2000
      
      aggregations {
        alignment_period   = "60s"
        per_series_aligner = "ALIGN_RATE"
      }
    }
  }
  
  notification_channels = [
    google_monitoring_notification_channel.slack.name
  ]
}
```

## 🔐 Security & Compliance

### Secret Management
```hcl
resource "google_secret_manager_secret" "app_secrets" {
  for_each = toset([
    "tidb-password",
    "jwt-secret", 
    "oauth-client-secret",
    "encryption-key"
  ])
  
  secret_id = each.key
  
  replication {
    user_managed {
      replicas {
        location = var.region
      }
    }
  }
}
```

### IAM & Access Control
```hcl
# Least privilege service account
resource "google_service_account" "api_service" {
  account_id   = "bocchi-api"
  display_name = "Bocchi API Service Account"
  description  = "Service account for Cloud Run API"
}

# Minimal required permissions
resource "google_project_iam_member" "api_permissions" {
  for_each = toset([
    "roles/secretmanager.secretAccessor",
    "roles/logging.logWriter",
    "roles/monitoring.metricWriter"
  ])
  
  project = var.project_id
  role    = each.key
  member  = "serviceAccount:${google_service_account.api_service.email}"
}
```

### Network Security
```hcl
# VPC Connector for secure database access
resource "google_vpc_access_connector" "main" {
  name          = "bocchi-connector"
  region        = var.region
  ip_cidr_range = "10.8.0.0/28"
  network       = "default"
  
  min_instances = 2
  max_instances = 3
}
```

## 💰 Cost Optimization

### Resource Tagging
```hcl
locals {
  common_tags = {
    Project     = "bocchi-the-map"
    Environment = var.environment
    ManagedBy   = "terraform"
    Owner       = "engineering-team"
    CostCenter  = "product"
  }
}
```

### Budget Alerts
```hcl
resource "google_billing_budget" "project_budget" {
  billing_account = var.billing_account_id
  display_name    = "Bocchi Map ${title(var.environment)} Budget"
  
  budget_filter {
    projects = ["projects/${var.project_id}"]
  }
  
  amount {
    specified_amount {
      currency_code = "USD"
      units         = var.monthly_budget
    }
  }
  
  threshold_rules {
    threshold_percent = 0.8
    spend_basis      = "CURRENT_SPEND"
  }
  
  threshold_rules {
    threshold_percent = 1.0
    spend_basis      = "FORECASTED_SPEND"
  }
}
```

### Cost Monitoring
| Service | Monthly Cost (Dev) | Monthly Cost (Prod) |
|---------|-------------------|-------------------|
| **Cloud Run** | $0-20 | $50-500 |
| **Cloudflare R2** | $5 | $15-50 |
| **TiDB Serverless** | $0 (Free tier) | $25-200 |
| **Networking** | $2-5 | $10-30 |
| **Total** | **$7-30** | **$100-780** |

## 🚢 Deployment Strategies

### Blue/Green Deployments
```hcl
resource "google_cloud_run_service" "api_blue" {
  count = var.deployment_strategy == "blue-green" ? 1 : 0
  name  = "bocchi-api-blue"
  # ... configuration
}

resource "google_cloud_run_service" "api_green" {
  count = var.deployment_strategy == "blue-green" ? 1 : 0
  name  = "bocchi-api-green"
  # ... configuration
}

# Traffic splitting via Load Balancer
resource "google_compute_url_map" "api_lb" {
  name = "bocchi-api-lb"
  
  default_service = google_compute_backend_service.api_backend.id
  
  host_rule {
    hosts        = ["api.bocchi-map.com"]
    path_matcher = "api-matcher"
  }
  
  path_matcher {
    name = "api-matcher"
    default_service = google_compute_backend_service.api_backend.id
    
    # Canary deployment: 10% to new version
    route_rules {
      priority = 1
      
      match_rules {
        prefix_match = "/"
        header_matches {
          name  = "X-Canary"
          exact_match = "true"
        }
      }
      
      weighted_backend_services {
        backend_service = google_compute_backend_service.api_green.id
        weight         = 10
      }
      
      weighted_backend_services {
        backend_service = google_compute_backend_service.api_blue.id
        weight         = 90
      }
    }
  }
}
```

### Environment Promotion
```bash
# Development → Staging
terraform workspace select staging
terraform apply -var-file=environments/staging/terraform.tfvars

# Staging → Production (with approval)
terraform workspace select prod
terraform plan -var-file=environments/prod/terraform.tfvars
# Manual approval required
terraform apply -var-file=environments/prod/terraform.tfvars
```

## 🧪 Testing Infrastructure

### Terratest Integration
```go
// tests/terraform_test.go
func TestTerraformGcpCloudRun(t *testing.T) {
    terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
        TerraformDir: "../modules/gcp",
        Vars: map[string]interface{}{
            "project_id":   "bocchi-test-" + strings.ToLower(random.UniqueId()),
            "region":       "us-central1", 
            "service_name": "test-api",
        },
    })
    
    defer terraform.Destroy(t, terraformOptions)
    terraform.InitAndApply(t, terraformOptions)
    
    // Validate Cloud Run service is accessible
    serviceUrl := terraform.Output(t, terraformOptions, "service_url")
    http_helper.HttpGetWithRetry(t, serviceUrl+"/health", nil, 200, "ok", 30, 5*time.Second)
}
```

### Infrastructure Validation
```bash
# Terraform validation
terraform fmt -check=true -diff=true
terraform validate
terraform plan -detailed-exitcode

# Security scanning
tfsec .
checkov -d .

# Cost estimation  
infracost breakdown --path .
```

## 📚 Operations Runbooks

### Disaster Recovery
```bash
# Database backup restoration
terraform import tidbcloud_backup.emergency ${BACKUP_ID}
terraform apply -target=tidbcloud_restore.emergency

# Application rollback
gcloud run services replace-traffic bocchi-api \
  --to-revisions=${PREVIOUS_REVISION}=100 \
  --region=asia-northeast1

# DNS failover (if needed)
cloudflare-cli dns update bocchi-map.com api \
  --content=${BACKUP_ENDPOINT} \
  --ttl=60
```

### Scaling Operations
```bash
# Emergency scaling
gcloud run services update bocchi-api \
  --max-instances=1000 \
  --region=asia-northeast1

# Database scaling
tidb-cli cluster scale-out \
  --cluster-id=${CLUSTER_ID} \
  --node-type=tikv \
  --count=3
```

## 🤝 Contributing

### Infrastructure Changes
```bash
# Setup development environment
terraform workspace new feature-xyz
terraform init
terraform plan

# Submit changes
git commit -m "feat(infra): add monitoring dashboard"
gh pr create --title "Infrastructure: Enhanced monitoring"
```

### Code Standards
- **Terraform Style** - Use `terraform fmt` and follow naming conventions
- **Module Design** - Reusable, well-documented modules
- **State Management** - Remote state with locking
- **Security** - No hardcoded secrets, with least privilege IAM

## 🎯 Roadmap

- [x] **v1.0** - Basic multi-cloud setup with auto-scaling
- [x] **v1.1** - Monitoring, alerting, and cost optimization
- [ ] **v1.2** - Advanced deployment strategies (blue/green, canary)
- [ ] **v1.3** - Multi-region deployment for global availability
- [ ] **v2.0** - Kubernetes migration for advanced orchestration

---

**☁️ Infrastructure designed for hypergrowth and operational excellence**

[📊 Cost Dashboard](https://console.cloud.google.com/billing) • [🔍 Monitoring](https://console.cloud.google.com/monitoring) • [📋 Status Page](https://status.bocchi-map.com)
