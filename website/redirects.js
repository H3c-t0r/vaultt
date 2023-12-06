/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

module.exports = [
  // example redirect:
  // {
  //   source: '/vault/docs/some/path',
  //   destination: '/vault/docs/some/other/path',
  //   permanent: true,
  // },
  {
    source: '/vault/docs/plugins/plugin-portal',
    destination: '/vault/integrations',
    permanent: true,
  }
]
