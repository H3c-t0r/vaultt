terraform "k8s" {
  required_version = ">= 1.2.0"
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }

    helm = {
      source  = "hashicorp/helm"
    }
  }
}

terraform_cli "default" {
  plugin_cache_dir = var.terraform_plugin_cache_dir != null ? abspath(var.terraform_plugin_cache_dir) : null

  credentials "app.terraform.io" {
    token = var.tfc_api_token
  }
}
