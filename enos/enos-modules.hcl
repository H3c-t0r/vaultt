# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

module "autopilot_upgrade_storageconfig" {
  source = "./modules/autopilot_upgrade_storageconfig"
}

module "backend_consul" {
  source = "./modules/backend_consul"

  license   = var.backend_license_path == null ? null : file(abspath(var.backend_license_path))
  log_level = var.backend_log_level
}

module "backend_raft" {
  source = "./modules/backend_raft"
}

module "build_crt" {
  source = "./modules/build_crt"
}

module "build_local" {
  source = "./modules/build_local"
}

module "build_artifactory" {
  source = "./modules/vault_artifactory_artifact"
}

module "create_vpc" {
  source = "./modules/create_vpc"

  environment = "ci"
  common_tags = var.tags
}

module "ec2_info" {
  source = "./modules/ec2_info"
}

module "get_local_metadata" {
  source = "./modules/get_local_metadata"
}

module "generate_secondary_token" {
  source = "./modules/generate_secondary_token"

  vault_install_dir = var.vault_install_dir
}

module "read_license" {
  source = "./modules/read_license"
}

module "shutdown_node" {
  source = "./modules/shutdown_node"
}

module "shutdown_multiple_nodes" {
  source = "./modules/shutdown_multiple_nodes"
}

module "target_ec2_fleet" {
  source = "./modules/target_ec2_fleet"

  capacity_type    = "on-demand" // or "spot", use on-demand until we can stabilize spot fleets
  common_tags      = var.tags
  instance_mem_min = 4096
  instance_cpu_min = 2
  max_price        = "0.1432" // On-demand cost for RHEL/t3.medium on-demand in us-east
  project_name     = var.project_name
  ssh_keypair      = var.aws_ssh_keypair_name
}

module "target_ec2_spot_fleet" {
  source = "./modules/target_ec2_spot_fleet"

  common_tags      = var.tags
  instance_mem_min = 4096
  instance_cpu_min = 2
  max_price        = "0.1432" // On-demand cost for RHEL/t3.medium on-demand in us-east
  project_name     = var.project_name
  ssh_keypair      = var.aws_ssh_keypair_name
}

module "vault_agent" {
  source = "./modules/vault_agent"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_verify_agent_output" {
  source = "./modules/vault_verify_agent_output"

  vault_instance_count = var.vault_instance_count
}

module "vault_cluster" {
  source = "./modules/vault_cluster"

  install_dir    = var.vault_install_dir
  consul_license = var.backend_license_path == null ? null : file(abspath(var.backend_license_path))
  log_level      = var.vault_log_level
}

module "vault_get_cluster_ips" {
  source = "./modules/vault_get_cluster_ips"

  vault_install_dir = var.vault_install_dir
}

module "vault_unseal_nodes" {
  source = "./modules/vault_unseal_nodes"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_upgrade" {
  source = "./modules/vault_upgrade"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_verify_autopilot" {
  source = "./modules/vault_verify_autopilot"

  vault_autopilot_upgrade_status = "await-server-removal"
  vault_install_dir              = var.vault_install_dir
  vault_instance_count           = var.vault_instance_count
}

module "vault_verify_raft_auto_join_voter" {
  source = "./modules/vault_verify_raft_auto_join_voter"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_verify_undo_logs" {
  source = "./modules/vault_verify_undo_logs"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_verify_replication" {
  source = "./modules/vault_verify_replication"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_verify_ui" {
  source = "./modules/vault_verify_ui"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_verify_unsealed" {
  source = "./modules/vault_verify_unsealed"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_setup_perf_primary" {
  source = "./modules/vault_setup_perf_primary"

  vault_install_dir = var.vault_install_dir
}

module "vault_setup_perf_secondary" {
  source = "./modules/vault_setup_perf_secondary"

  vault_install_dir = var.vault_install_dir
}

module "vault_verify_read_data" {
  source = "./modules/vault_verify_read_data"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_verify_performance_replication" {
  source = "./modules/vault_verify_performance_replication"

  vault_install_dir = var.vault_install_dir
}

module "vault_verify_version" {
  source = "./modules/vault_verify_version"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_verify_write_data" {
  source = "./modules/vault_verify_write_data"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_raft_remove_peer" {
  source            = "./modules/vault_raft_remove_peer"
  vault_install_dir = var.vault_install_dir
}

module "vault_test_ui" {
  source = "./modules/vault_test_ui"

  ui_run_tests = var.ui_run_tests
}
