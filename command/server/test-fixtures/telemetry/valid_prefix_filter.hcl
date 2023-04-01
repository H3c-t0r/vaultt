# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

disable_mlock = true
ui            = true

telemetry {
  statsd_address = "foo"
  prefix_filter = ["-vault.expire", "-vault.audit", "+vault.expire.num_irrevocable_leases"]
}