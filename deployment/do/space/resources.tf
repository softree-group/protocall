terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

resource "digitalocean_spaces_bucket" "main" {
  name   = var.name
  region = var.region
  acl = var.acl
}

output "name" {
  value = digitalocean_spaces_bucket.main.name
}

output "domain" {
  value = digitalocean_spaces_bucket.main.bucket_domain_name
}
