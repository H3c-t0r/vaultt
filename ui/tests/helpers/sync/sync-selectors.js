/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn } from '@ember/test-helpers';
import { SELECTORS as GENERAL } from 'vault/tests/helpers/general-selectors';

export const PAGE = {
  ...GENERAL,
  cta: {
    summary: '[data-test-cta-container] p',
    button: '[data-test-cta-button]',
  },
  associations: {
    list: {
      name: '[data-test-association-name]',
      status: '[data-test-association-status]',
      updated: '[data-test-association-updated]',
      menu: {
        sync: '[data-test-association-action="sync"]',
        view: '[data-test-association-action="view"]',
        unsync: '[data-test-association-action="unsync"]',
      },
    },
  },
  syncBadge: {
    icon: (name) => `[data-test-icon="${name}"]`,
    text: '.hds-badge__text',
  },
  selectType: (type) => `[data-test-select-destination="${type}"]`,
  cancelButton: '[data-test-cancel]',
  saveButton: '[data-test-save]',
  toolbar: (btnText) => `[data-test-toolbar="${btnText}"]`,
  form: {
    fillInByAttr: async (attr, value) => {
      // for handling more complex form input elements by attr name
      switch (attr) {
        case 'credentials':
          await click('[data-test-text-toggle]');
          return fillIn('[data-test-text-file-textarea]', value);
        case 'deploymentEnvironments':
          await click('[data-test-input="deploymentEnvironments"] input#deployment');
          await click('[data-test-input="deploymentEnvironments"] input#preview');
          return await click('[data-test-input="deploymentEnvironments"] input#production');
        default:
          return fillIn(`[data-test-input="${attr}"]`, value);
      }
    },
  },
};
