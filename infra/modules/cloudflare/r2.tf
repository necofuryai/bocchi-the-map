# Cloudflare R2 bucket for map tiles
resource "cloudflare_r2_bucket" "map_tiles" {
  account_id = var.cloudflare_account_id
  name       = "bocchi-map-tiles-${var.environment}"
  location   = "APAC"
}

# Variables
variable "cloudflare_account_id" {
  description = "Cloudflare Account ID"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

# Outputs
output "r2_bucket_name" {
  value = cloudflare_r2_bucket.map_tiles.name
}

output "r2_bucket_id" {
  value = cloudflare_r2_bucket.map_tiles.id
}