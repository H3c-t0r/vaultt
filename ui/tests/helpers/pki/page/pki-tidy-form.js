/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export const SELECTORS = {
  inputByAttr: (attr) => `[data-test-input="${attr}"]`,
  toggleInput: (attr) => `[data-test-input="${attr}"] input`,
  intervalDuration: '[data-test-ttl-value="Automatic tidy enabled"]',
  acmeAccountSafetyBuffer: '[data-test-ttl-value="Tidy ACME enabled"]',
  toggleLabel: (label) => `[data-test-toggle-label="${label}"]`,
  tidySectionHeader: (header) => `[data-test-tidy-header="${header}"]`,
  tidySave: '[data-test-pki-tidy-button]',
  tidyCancel: '[data-test-pki-tidy-cancel]',
};
