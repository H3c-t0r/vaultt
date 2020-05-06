import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const REPLICATION_DETAILS = {
  state: 'stream-wals',
  primaryClusterAddr: 'https://127.0.0.1:8201',
};

const REPLICATION_DETAILS_SYNCING = {
  state: 'merkle-diff',
  primaryClusterAddr: 'https://127.0.0.1:8201',
};

module('Integration | Enterprise | Component | replication-dashboard', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('clusterMode', 'secondary');
    this.set('isSecondary', true);
  });

  test('it renders', async function(assert) {
    this.set('replicationDetails', REPLICATION_DETAILS);

    await render(hbs`<ReplicationDashboard
    @replicationDetails={{replicationDetails}}
    @clusterMode={{clusterMode}}
    @isSecondary={{isSecondary}}
    />`);

    assert.dom('[data-test-replication-dashboard]').exists();
    assert.dom('[data-test-table-rows').exists();
    assert
      .dom('[data-test-primary-cluster-address]')
      .includesText(
        REPLICATION_DETAILS.primaryClusterAddr,
        `shows the correct primary cluster address value`
      );
    assert.dom('[data-test-replication-doc-link]').exists();
    assert.dom('[data-test-flash-message]').doesNotExist('no flash message is displayed on render');
  });

  test('it renders an alert banner if the dashboard is syncing', async function(assert) {
    this.set('replicationDetailsSyncing', REPLICATION_DETAILS_SYNCING);

    await render(hbs`<ReplicationDashboard 
    @replicationDetails={{replicationDetailsSyncing}} 
    @clusterMode={{clusterMode}}
    @isSecondary={{isSecondary}}
    @componentToRender='replication-secondary-card'
    />`);

    assert.dom('[data-test-isSyncing]').exists();
    assert.dom('[data-test-isReindexing]').doesNotExist();
  });
});
