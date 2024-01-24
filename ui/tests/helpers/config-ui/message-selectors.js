/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export const PAGE = {
  // General selectors that are common between pages
  radio: (radioName) => `[data-test-radio="${radioName}"]`,
  field: (fieldName) => `[data-test-field="${fieldName}"]`,
  input: (input) => `[data-test-input="${input}"]`,
  button: (buttonName) => `[data-test-button="${buttonName}"]`,
  inlineErrorMessage: `[data-test-inline-error-message]`,
  fieldVaildation: (fieldName) => `[data-test-field-validation="${fieldName}"]`,
  modal: (name) => `[data-test-modal="${name}"]`,
  modalTitle: (title) => `[data-test-modal-title="${title}"]`,
  modalBody: (name) => `[data-test-modal-body="${name}"]`,
  modalButton: (name) => `[data-test-modal-button="${name}"]`,
  alert: (name) => `data-test-custom-alert=${name}`,
  alertTitle: (name) => `[data-test-custom-alert-title="${name}"]`,
  alertDescription: (name) => `[data-test-custom-alert-description="${name}"]`,
  badge: (name) => `[data-test-badge="${name}"]`,
};
