import { create } from 'ember-cli-page-object';
import { module, skip, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import { click, currentURL, fillIn, find, visit } from '@ember/test-helpers';
import { SELECTORS } from 'vault/tests/helpers/pki/workflow';
import { adminPolicy, readerPolicy, updatePolicy } from 'vault/tests/helpers/pki/policy-generator';

const consoleComponent = create(consoleClass);

const tokenWithPolicy = async function (name, policy) {
  await consoleComponent.runCommands([
    `write sys/policies/acl/${name} policy=${btoa(policy)}`,
    `write -field=client_token auth/token/create policies=${name}`,
  ]);
  return consoleComponent.lastLogOutput;
};

const runCommands = async function (commands) {
  try {
    await consoleComponent.runCommands(commands);
    const res = consoleComponent.lastLogOutput;
    if (res.includes('Error')) {
      throw new Error(res);
    }
    return res;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error(
      `The following occurred when trying to run the command(s):\n ${commands.join('\n')} \n\n ${
        consoleComponent.lastLogOutput
      }`
    );
    throw error;
  }
};

/**
 * This test module should test the PKI workflow, including:
 * - link between pages and confirm that the url is as expected
 * - log in as user with a policy and ensure expected UI elements are shown/hidden
 */
module('Acceptance | pki workflow', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    await authPage.login();
    // Setup PKI engine
    const mountPath = `pki-workflow-${new Date().getTime()}`;
    await enablePage.enable('pki', mountPath);
    this.mountPath = mountPath;
    await logout.visit();
  });

  hooks.afterEach(async function () {
    await logout.visit();
    await authPage.login();
    // Cleanup engine
    await runCommands([`delete sys/mounts/${this.mountPath}`]);
    await logout.visit();
  });

  test('empty state messages are correct when PKI not configured', async function (assert) {
    assert.expect(9);
    const assertEmptyState = (assert, resource) => {
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/${resource}`);
      assert
        .dom(SELECTORS.emptyStateTitle)
        .hasText(
          'PKI not configured',
          `${resource} index renders correct empty state title when PKI not configured`
        );
      assert
        .dom(SELECTORS.emptyStateMessage)
        .hasText(
          `This PKI mount hasn't yet been configured with a certificate issuer.`,
          `${resource} index empty state message correct when PKI not configured`
        );
    };
    await authPage.login(this.pkiAdminToken);
    await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);

    // TODO comment in when roles index empty state updated & update assert.expect() number
    // await click(SELECTORS.rolesTab);
    // assertEmptyState(assert, 'roles');

    await click(SELECTORS.issuersTab);
    assertEmptyState(assert, 'issuers');

    await click(SELECTORS.certsTab);
    assertEmptyState(assert, 'certificates');

    await click(SELECTORS.keysTab);
    assertEmptyState(assert, 'keys');
  });

  module('roles', function (hooks) {
    hooks.beforeEach(async function () {
      await authPage.login();
      // Setup role-specific items
      await runCommands([
        `write ${this.mountPath}/roles/some-role \
      issuer_ref="default" \
      allowed_domains="example.com" \
      allow_subdomains=true \
      max_ttl="720h"`,
      ]);
      const pki_admin_policy = adminPolicy(this.mountPath, 'roles');
      const pki_reader_policy = readerPolicy(this.mountPath, 'roles');
      const pki_editor_policy = updatePolicy(this.mountPath, 'roles');
      this.pkiRoleReader = await tokenWithPolicy('pki-reader', pki_reader_policy);
      this.pkiRoleEditor = await tokenWithPolicy('pki-editor', pki_editor_policy);
      this.pkiAdminToken = await tokenWithPolicy('pki-admin', pki_admin_policy);
      await logout.visit();
    });

    test('shows correct items if user has all permissions', async function (assert) {
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
      assert.dom(SELECTORS.rolesTab).exists('Roles tab is present');
      await click(SELECTORS.rolesTab);
      assert.dom(SELECTORS.createRoleLink).exists({ count: 1 }, 'Create role link is rendered');
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles`);
      assert.dom('.linked-block').exists({ count: 1 }, 'One role is in list');
      await click('.linked-block');
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/some-role/details`);

      assert.dom(SELECTORS.generateCertLink).exists('Generate cert link is shown');
      await click(SELECTORS.generateCertLink);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/some-role/generate`);

      // Go back to details and test all the links
      await visit(`/vault/secrets/${this.mountPath}/pki/roles/some-role/details`);
      assert.dom(SELECTORS.signCertLink).exists('Sign cert link is shown');
      await click(SELECTORS.signCertLink);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/some-role/sign`);

      await visit(`/vault/secrets/${this.mountPath}/pki/roles/some-role/details`);
      assert.dom(SELECTORS.editRoleLink).exists('Edit link is shown');
      await click(SELECTORS.editRoleLink);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/some-role/edit`);

      await visit(`/vault/secrets/${this.mountPath}/pki/roles/some-role/details`);
      assert.dom(SELECTORS.deleteRoleButton).exists('Delete role button is shown');
      await click(`${SELECTORS.deleteRoleButton} [data-test-confirm-action-trigger]`);
      await click(`[data-test-confirm-button]`);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/roles`,
        'redirects to roles list after deletion'
      );
    });

    test('it does not show toolbar items the user does not have permission to see', async function (assert) {
      await authPage.login(this.pkiRoleReader);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      assert.dom(SELECTORS.rolesTab).exists('Roles tab is present');
      await click(SELECTORS.rolesTab);
      assert.dom(SELECTORS.createRoleLink).exists({ count: 1 }, 'Create role link is rendered');
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles`);
      assert.dom('.linked-block').exists({ count: 1 }, 'One role is in list');
      await click('.linked-block');
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/some-role/details`);
      assert.dom(SELECTORS.deleteRoleButton).doesNotExist('Delete role button is not shown');
      assert.dom(SELECTORS.generateCertLink).doesNotExist('Generate cert link is not shown');
      assert.dom(SELECTORS.signCertLink).doesNotExist('Sign cert link is not shown');
      assert.dom(SELECTORS.editRoleLink).doesNotExist('Edit link is not shown');
    });

    test('it shows correct toolbar items for the user policy', async function (assert) {
      await authPage.login(this.pkiRoleEditor);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      assert.dom(SELECTORS.rolesTab).exists('Roles tab is present');
      await click(SELECTORS.rolesTab);
      assert.dom(SELECTORS.createRoleLink).exists({ count: 1 }, 'Create role link is rendered');
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles`);
      assert.dom('.linked-block').exists({ count: 1 }, 'One role is in list');
      await click('.linked-block');
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/some-role/details`);
      assert.dom(SELECTORS.deleteRoleButton).doesNotExist('Delete role button is not shown');
      assert.dom(SELECTORS.generateCertLink).exists('Generate cert link is shown');
      assert.dom(SELECTORS.signCertLink).exists('Sign cert link is shown');
      assert.dom(SELECTORS.editRoleLink).exists('Edit link is shown');
      await click(SELECTORS.editRoleLink);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/roles/some-role/edit`,
        'Links to edit view'
      );
      await click(SELECTORS.roleForm.roleCancelButton);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/roles/some-role/details`,
        'Cancel from edit goes to details'
      );
      await click(SELECTORS.generateCertLink);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/roles/some-role/generate`,
        'Generate cert button goes to generate page'
      );
      await click(SELECTORS.generateCertForm.cancelButton);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/roles/some-role/details`,
        'Cancel from generate goes to details'
      );
    });

    test('create role happy path', async function (assert) {
      const roleName = 'another-role';
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
      assert.dom(SELECTORS.rolesTab).exists('Roles tab is present');
      await click(SELECTORS.rolesTab);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles`);
      await click(SELECTORS.createRoleLink);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/create`);
      assert.dom(SELECTORS.breadcrumbContainer).exists({ count: 1 }, 'breadcrumbs are rendered');
      assert.dom(SELECTORS.breadcrumbs).exists({ count: 4 }, 'Shows 4 breadcrumbs');
      assert.dom(SELECTORS.pageTitle).hasText('Create a PKI role');

      await fillIn(SELECTORS.roleForm.roleName, roleName);
      await click(SELECTORS.roleForm.roleCreateButton);

      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/${roleName}/details`);
      assert.dom(SELECTORS.breadcrumbs).exists({ count: 4 }, 'Shows 4 breadcrumbs');
      assert.dom(SELECTORS.pageTitle).hasText(`PKI Role ${roleName}`);
    });
  });

  module('keys', function (hooks) {
    hooks.beforeEach(async function () {
      await authPage.login();
      // base config pki so empty state doesn't show
      await runCommands([`write ${this.mountPath}/root/generate/internal common_name="Hashicorp Test"`]);
      const pki_admin_policy = adminPolicy(this.mountPath);
      const pki_reader_policy = readerPolicy(this.mountPath, 'keys', true);
      const pki_editor_policy = updatePolicy(this.mountPath, 'keys');
      this.pkiKeyReader = await tokenWithPolicy('pki-reader', pki_reader_policy);
      this.pkiKeyEditor = await tokenWithPolicy('pki-editor', pki_editor_policy);
      this.pkiAdminToken = await tokenWithPolicy('pki-admin', pki_admin_policy);
      await logout.visit();
    });

    test('shows correct items if user has all permissions', async function (assert) {
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.keysTab);
      // index page
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys`);
      assert
        .dom(SELECTORS.keyPages.importKey)
        .hasAttribute(
          'href',
          `/ui/vault/secrets/${this.mountPath}/pki/keys/import`,
          'import link renders with correct url'
        );
      let keyId = find(SELECTORS.keyPages.keyId).innerText;
      assert.dom('.linked-block').exists({ count: 1 }, 'One key is in list');
      await click('.linked-block');
      // details page
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys/${keyId}/details`);
      assert.dom(SELECTORS.keyPages.downloadButton).doesNotExist('does not download button for private key');

      // edit page
      await click(SELECTORS.keyPages.keyEditLink);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys/${keyId}/edit`);
      await click(SELECTORS.keyForm.keyCancelButton);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/keys/${keyId}/details`,
        'navigates back to details on cancel'
      );
      await visit(`/vault/secrets/${this.mountPath}/pki/keys/${keyId}/edit`);
      await fillIn(SELECTORS.keyForm.keyNameInput, 'test-key');
      await click(SELECTORS.keyForm.keyCreateButton);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/keys/${keyId}/details`,
        'navigates to details after save'
      );
      await this.pauseTest;
      assert.dom(SELECTORS.keyPages.keyNameValue).hasText('test-key', 'updates key name');

      // key generate and delete navigation
      await visit(`/vault/secrets/${this.mountPath}/pki/keys`);
      await click(SELECTORS.keyPages.generateKey);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys/create`);
      await fillIn(SELECTORS.keyForm.typeInput, 'exported');
      await fillIn(SELECTORS.keyForm.keyTypeInput, 'rsa');
      await click(SELECTORS.keyForm.keyCreateButton);
      keyId = find(SELECTORS.keyPages.keyIdValue).innerText;
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys/${keyId}/details`);

      assert
        .dom(SELECTORS.alertBanner)
        .hasText(
          'Next steps This private key material will only be available once. Copy or download it now.',
          'renders banner to save private key'
        );
      assert.dom(SELECTORS.keyPages.downloadButton).exists('renders download button');
      await click(SELECTORS.keyPages.keyDeleteButton);
      await click(SELECTORS.keyPages.confirmDelete);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/keys/`,
        'navigates back to key list view on delete'
      );
    });

    test('it does not show toolbar items the user does not have permission to see', async function (assert) {
      await authPage.login(this.pkiKeyReader);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.keysTab);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys`);
      assert.dom(SELECTORS.keyPages.importKey).doesNotExist();
      assert.dom(SELECTORS.keyPages.generateKey).doesNotExist();
      assert.dom('.linked-block').exists({ count: 1 }, 'One key is in list');
      const keyId = find(SELECTORS.keyPages.keyId).innerText;
      await click('.linked-block');
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys/${keyId}/details`);
      assert.dom(SELECTORS.keyPages.keyDeleteButton).doesNotExist('Delete key button is not shown');
      assert.dom(SELECTORS.keyPages.keyEditLink).doesNotExist('Edit key button does not render');
    });

    // TODO CMB: add edit capabilities test
    skip('it shows correct toolbar items for the user policy', async function (assert) {
      await authPage.login(this.pkiKeyEditor);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.keysTab);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys`);
    });
  });
});
