# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.4.3"
    }
  }
}

variable "ip_version" {
  type        = number
  description = "The IP version used for the Vault TCP listener"

  validation {
    condition     = contains([4, 6], var.ip_version)
    error_message = "The ip_version must be either 4 or 6"
  }
}

variable "primary_leader_host" {
  type = object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  })
  description = "The primary cluster leader host"
}

variable "vault_addr" {
  type        = string
  description = "The local vault API listen address"
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}


variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

locals {
  primary_leader_addr  = var.ip_version == 6 ? var.primary_leader_host.ipv6 : var.primary_leader_host.private_ip
  token_id             = random_uuid.token_id.id
  secondary_public_key = enos_remote_exec.fetch_secondary_public_key.stdout
}

resource "random_uuid" "token_id" {}

resource "enos_remote_exec" "fetch_secondary_public_key" {
  depends_on = [random_uuid.token_id]
  environment = {
    VAULT_ADDR  = var.vault_addr
    VAULT_TOKEN = var.vault_root_token
  }

  inline = ["${var.vault_install_dir}/vault write -field secondary_public_key -f sys/replication/dr/secondary/generate-public-key"]

  transport = {
    ssh = {
      host = var.primary_leader_host.public_ip
    }
  }
}

output "secondary_public_key" {
  value = local.secondary_public_key
}
