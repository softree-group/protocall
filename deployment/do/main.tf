terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

provider "digitalocean" {}

# module "k8s" {
#   source = "./k8s"

#   namespace = "k8s"

#   kubernetes_version = "1.21.5-do.0"
#   ip_range           = "10.0.0.0/16"
#   node_count = 3
#   region             = "fra1"
#   size       = "s-2vcpu-4gb"
# }

# resource "local_file" "kubeconfig" {
#   sensitive_content = module.k8s.kubernetes_cluster_raw_config
#   filename          = "${path.root}/kube.yml"
#   file_permission   = "0644"
# }

# module "protocall-static" {
#   source = "./space"

#   name   = "protocall-static"
#   region = "fra1"
#   acl    = "public-read"
# }

# output "protocall-static-bucket" {
#   value = module.protocall-static.name
# }

# output "protocall-static-domain" {
#   value = module.protocall-static.domain
# }

module "protocall-storage" {
  source = "./space"

  name   = "protocall-storage"
  region = "fra1"
  acl    = "private"
}

output "protocall-storage-bucket" {
  value = module.protocall-storage.name
}

output "protocall-storage-domain" {
  value = module.protocall-storage.domain
}
