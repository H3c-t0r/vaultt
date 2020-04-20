import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const DATA = {
  clusterId: 'b829d963-6835-33eb-a903-57376024b97a',
  mode: 'primary',
  merkleRoot: 'c21c8428a0a06135cef6ae25bf8e0267ff1592a6',
};

module('Integration | Enterprise | Component | replication-table-rows', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('data', DATA);
  });

  test('it renders', async function(assert) {
    await render(hbs`<ReplicationTableRows @data={{data}}/>`);

    assert.dom('[data-test-table-rows]').exists();
  });

  test('it renders with merkle root, mode, replication set', async function(assert) {
    await render(hbs`<ReplicationTableRows @data={{data}}/>`);

    assert.dom('.empty-state').doesNotExist('does not show empty state when data is found');

    Object.keys(DATA).forEach(attr => {
      let expected = DATA[attr];
      assert.dom(`[data-test-attr="${expected}"]`).includesText(expected, `shows the correct ${attr}`);
    });
  });

  test('it renders unknown if values cannot be found', async function(assert) {
    const noAttrs = {
      clusterId: null,
      mode: null,
      merkleRoot: null,
    };
    this.set('data', noAttrs);
    await render(hbs`<ReplicationTableRows @data={{data}}/>`);

    assert.dom('[data-test-table-rows]').includesText('unknown');
  });
});
