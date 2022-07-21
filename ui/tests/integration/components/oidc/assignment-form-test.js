import { module, test, skip } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn, click, findAll } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Integration | Component | oidc/assignment-form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
  });

  test('it should save new assignment', async function (assert) {
    assert.expect(6);

    this.server.post('/identity/oidc/assignment/test', (schema, req) => {
      assert.ok(true, 'Request made to save assignment');
      return JSON.parse(req.requestBody);
    });

    // so test can select an entity to remove validation error
    this.server.get('/identity/entity/id', () => ({
      data: {
        key_info: { 'f831667b-7392-7a1c-c0fc-33d48cb1c57d': { name: 'test-entity' } },
        keys: ['f831667b-7392-7a1c-c0fc-33d48cb1c57d'],
      },
    }));
    this.server.get('/identity/group/id', () => ({
      data: {
        key_info: { 'h831667b-7392-7a1c-c0fc-33d48cb1c57d': { name: 'test-group' } },
        keys: ['h831667b-7392-7a1c-c0fc-33d48cb1c57d'],
      },
    }));

    this.model = this.store.createRecord('oidc/assignment');
    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');

    await render(hbs`
      <Oidc::AssignmentForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);

    assert.dom('[data-test-oidc-assignment-title]').hasText('Create assignment', 'Form title renders');
    assert.dom('[data-test-oidc-assignment-save]').hasText('Create', 'Save button has correct label');
    await click('[data-test-oidc-assignment-save]');
    assert
      .dom('[data-test-inline-alert]')
      .hasText('Name is required.', 'Validation message is shown for name');
    assert.equal(findAll('[data-test-inline-error-message]').length, 2, `there are two validations errors.`);
    await fillIn('[data-test-input="name"]', 'test');
    await click('[data-test-component="search-select"] .ember-basic-dropdown-trigger');
    await click('.ember-power-select-option');
    await click('[data-test-oidc-assignment-save]');
  });

  skip('it should rollback attributes or unload record on cancel', async function (assert) {
    // ARG TODO WIP after finish update view
    assert.expect(5);

    this.onCancel = () => assert.ok(true, 'onCancel callback fires');

    this.model = this.store.createRecord('oidc/assignment');

    await render(hbs`
      <Oidc::AssignmentForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);

    await click('[data-test-oidc-assignment-cancel]');
    assert.true(this.model.isDestroyed, 'New model is unloaded on cancel');

    this.store.pushPayload('oidc/assignment', {
      modelName: 'oidc/assignment',
      name: 'test',
    });
    this.model = this.store.peekRecord('oidc/assignment', 'test');

    await render(hbs`
    <Oidc::AssignmentForm
      @model={{this.model}}
      @onCancel={{this.onCancel}}
      @onSave={{this.onSave}}
    />
  `);

    await fillIn('[data-test-string-list-input="0"]', 'entity-id');
    await click('[data-test-oidc-assignment-cancel]');
    // ARG TODO need to change the entity ID or group ID but need to finish the update view first.
    assert.equal(this.model.name, 'test', 'Model attributes are rolled back on cancel');
  });

  skip('it should update assignment', async function () {
    // ARG TODO in next PR. Need to modify model to show entities and groups.
  });
});
