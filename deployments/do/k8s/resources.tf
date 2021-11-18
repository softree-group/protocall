terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

resource "digitalocean_vpc" "main" {
  name     = "${var.namespace}-vpc"
  region   = var.region
  ip_range = var.ip_range

  timeouts {
    delete = "10m"
  }
}

resource "digitalocean_kubernetes_cluster" "main" {
  name         = "${var.namespace}-kubernetes-cluster"
  region       = var.region
  version      = var.kubernetes_version
  vpc_uuid     = digitalocean_vpc.main.id
  auto_upgrade = true

  node_pool {
    name       = "${var.namespace}-default-node-pool"
    auto_scale = false
    size       = var.size
    node_count = var.node_count
    labels     = var.labels
    tags       = var.tags
  }

  tags = var.tags
}

output "kubernetes_cluster_raw_config" {
  value     = digitalocean_kubernetes_cluster.main.kube_config[0].raw_config
  sensitive = true
}
