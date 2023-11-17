/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

const SELECTORS = {
  confirmToggle: '[data-test-confirm-action-trigger]',
  title: '[data-test-confirm-action-title]',
  message: '[data-test-confirm-action-message]',
  confirm: '[data-test-confirm-button]',
  cancel: '[data-test-confirm-cancel-button]',
};
module('Integration | Component | confirm-action', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.onConfirm = sinon.spy();
  });

  test('it renders defaults and calls onConfirmAction', async function (assert) {
    await render(hbs`
      <ConfirmAction
        @buttonText="DELETE"
        @onConfirmAction={{this.onConfirm}}
      />
      `);

    assert.dom(SELECTORS.confirmToggle).hasText('DELETE', 'renders button text');
    await click(SELECTORS.confirmToggle);
    assert.dom(SELECTORS.title).hasText('Are you sure?', 'renders default title');
    assert
      .dom(SELECTORS.message)
      .hasText('You will not be able to recover it later.', 'renders default body text');
    await click(SELECTORS.cancel);
    assert.false(this.onConfirm.called, 'does not call the action when Cancel is clicked');
    await click(SELECTORS.confirmToggle);
    await click(SELECTORS.confirm);
    assert.true(this.onConfirm.called, 'calls the action when Confirm is clicked');
    assert.dom(SELECTORS.title).doesNotExist('modal closes after confirm is clicked');
  });

  test('it renders defaults in dropdown and calls onConfirmAction', async function (assert) {
    await render(hbs`
      <ConfirmAction
        @buttonText="DELETE"
        @onConfirmAction={{this.onConfirm}}
        @isInDropdown={{true}}
      />
      `);

    assert.dom(`li ${SELECTORS.confirmToggle}`).exists('element renders inside <li>');
    assert
      .dom(SELECTORS.confirmToggle)
      .hasClass('hds-confirm-action-critical', 'button has dropdown styling');
    await click(SELECTORS.confirmToggle);
    assert.dom(SELECTORS.title).hasText('Are you sure?', 'renders default title');
    assert
      .dom(SELECTORS.message)
      .hasText('You will not be able to recover it later.', 'renders default body text');
    await click('[data-test-confirm-cancel-button]');
    assert.false(this.onConfirm.called, 'does not call the action when Cancel is clicked');
    await click(SELECTORS.confirmToggle);
    await click(SELECTORS.confirm);
    assert.true(this.onConfirm.called, 'calls the action when Confirm is clicked');
    assert.dom(SELECTORS.title).doesNotExist('modal closes after confirm is clicked');
  });

  test('it renders loading state', async function (assert) {
    await render(hbs`
      <ConfirmAction
        @buttonText="Open!"
        @onConfirmAction={{this.onConfirm}}
        @isRunning={{true}}
      />
      `);

    await click(SELECTORS.confirmToggle);

    assert.dom(SELECTORS.confirm).isDisabled('disables confirm button when loading');
    assert.dom('[data-test-confirm-button] [data-test-icon="loading"]').exists('it renders loading icon');
  });

  test('it renders passed args', async function (assert) {
    this.condition = false;
    await render(hbs`
      <ConfirmAction
        @buttonText="Open!"
        @onConfirmAction={{this.onConfirm}}
        @buttonColor="secondary"
        @confirmTitle="Do this?"
        @confirmMessage="Are you really, really sure?"
        @disabledMessage={{if this.condition "This is the reason you cannot do the thing"}}
      />
      `);

    await click(SELECTORS.confirmToggle);
    assert.dom(SELECTORS.title).hasText('Do this?', 'renders passed title');
    assert.dom(SELECTORS.message).hasText('Are you really, really sure?', 'renders passed body text');
    assert.dom(SELECTORS.confirm).hasText('Confirm');
  });

  test('it renders disabled modal', async function (assert) {
    this.condition = true;
    await render(hbs`
      <ConfirmAction
        @buttonText="Open!"
        @onConfirmAction={{this.onConfirm}}
        @buttonColor="secondary"
        @confirmTitle="Do this?"
        @confirmMessage="Are you really, really sure?"
        @disabledMessage={{if this.condition "This is the reason you cannot do the thing"}}
      />
      `);

    await click(SELECTORS.confirmToggle);
    assert.dom(SELECTORS.title).hasText('Not allowed', 'renders disabled title');
    assert
      .dom(SELECTORS.message)
      .hasText('This is the reason you cannot do the thing', 'renders disabled message as body text');
    assert.dom(SELECTORS.confirm).doesNotExist('Close');
    assert.dom(SELECTORS.cancel).hasText('Close');
  });
});
