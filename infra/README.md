# Bocchi The Map Infrastructure

Infrastructure as Code (IaC) for Bocchi The Map using Terraform.

## Architecture

- **Google Cloud Run**: API hosting
- **Cloudflare R2**: Map tile storage (PMTiles)
- **Cloudflare Pages**: Frontend hosting
- **TiDB Serverless**: Database

## Directory Structure

```
infra/
├── main.tf              # Main configuration
├── modules/             # Reusable modules
│   ├── cloudflare/      # Cloudflare resources
│   ├── gcp/             # Google Cloud resources
│   └── tidb/            # TiDB configuration
└── environments/        # Environment-specific configs
    ├── dev/             # Development environment
    └── prod/            # Production environment
```

## Prerequisites

- Terraform >= 1.5.0
- Google Cloud SDK
- Cloudflare API Token
- TiDB Serverless account

## Setup

1. **Initialize Terraform**
   ```bash
   terraform init
   ```

2. **Create workspace for environment**
   ```bash
   terraform workspace new dev
   terraform workspace select dev
   ```

3. **Create terraform.tfvars**
   ```hcl
   gcp_project_id       = "your-project-id"
   cloudflare_api_token = "your-api-token"
   cloudflare_account_id = "your-account-id"
   ```

4. **Plan and Apply**
   ```bash
   terraform plan
   terraform apply
   ```

## Environment Variables

### Google Cloud
- `GOOGLE_APPLICATION_CREDENTIALS`: Path to service account key

### Cloudflare
- `CLOUDFLARE_API_TOKEN`: API token with R2 permissions

## Deployment

### Development
```bash
terraform workspace select dev
terraform apply -var-file=environments/dev/terraform.tfvars
```

### Production
```bash
terraform workspace select prod
terraform apply -var-file=environments/prod/terraform.tfvars
```

## Resources Created

1. **Google Cloud**
   - Cloud Run service for API
   - Secret Manager for sensitive data

2. **Cloudflare**
   - R2 bucket for map tiles
   - (Future) Pages for frontend hosting

## Cost Estimation

- **Cloud Run**: ~$0-50/month (scales to zero)
- **R2 Storage**: ~$15/month per 1TB
- **TiDB Serverless**: Free tier available