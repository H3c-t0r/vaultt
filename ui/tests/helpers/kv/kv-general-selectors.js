/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export const SELECTORS = {
  toggleJson: '[data-test-toggle-input="json"]',
  toggleMasked: '[data-test-button="toggle-masked"]',
  jsonEditor: '[data-test-component="code-mirror-modifier"] textarea',
  tooltipTrigger: '[data-test-tooltip-trigger]',
  pageTitle: '[data-test-header-title]',
  infoRowValue: (label) => `[data-test-value-div="${label}"]`,
  secretTab: (tab) => `[data-test-secrets-tab="${tab}"]`,
  emptyStateTitle: '[data-test-empty-state-title]',
  emptyStateMessage: '[data-test-empty-state-message]',
  inlineAlert: '[data-test-inline-alert]',
  messageError: '[data-test-message-error]',
};

export const parseJsonEditor = (find) => {
  return JSON.parse(find(SELECTORS.jsonEditor).innerText);
};
