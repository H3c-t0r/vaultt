import { click, settled } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/settings/configure-secret-backends/pki/index';
import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { create } from 'ember-cli-page-object';
import fm from 'vault/tests/pages/components/flash-message';
const flashMessage = create(fm);
const SELECTORS = {
  generateSigningKey: '[data-test-ssh-input="generate-signing-key-checkbox"]',
  saveConfig: '[data-test-ssh-input="configure-submit"]',
  publicKey: '[data-test-ssh-input="public-key"]',
};
module('Acceptance | settings/configure/secrets/ssh', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  test('it configures ssh ca', async function (assert) {
    const path = `ssh-${new Date().getTime()}`;
    await enablePage.enable('ssh', path);
    await settled();
    await page.visit({ backend: path });
    await settled();
    assert.dom(SELECTORS.generateSigningKey).isChecked('generate_signing_key defaults to true');
    await click(SELECTORS.generateSigningKey);
    await click(SELECTORS.saveConfig);
    assert.strictEqual(
      flashMessage.latestMessage,
      'missing public_key',
      'renders warning flash message for failed save'
    );
    await click(SELECTORS.generateSigningKey);
    await click(SELECTORS.saveConfig);
    assert.dom(SELECTORS.publicKey).exists('renders public key after saving config');
  });
});
