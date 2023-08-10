/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { kvMetadataPath, kvDataPath } from 'vault/utils/kv-path';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import { encodePath } from 'vault/utils/path-encoding-helpers';

const EXAMPLE_KV_DATA_GET_RESPONSE = {
  request_id: 'foobar',
  data: {
    data: { hello: 'world' },
    metadata: {
      created_time: '2023-06-20T21:26:47.592306Z',
      custom_metadata: null,
      deletion_time: '',
      destroyed: false,
      version: 1,
    },
  },
};

module('Integration | Component | kv | Page::Secret::Metadata::Version-Diff', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    const store = this.owner.lookup('service:store');
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    const metadata = this.server.create('kv-metadatum');

    metadata.id = kvMetadataPath('kv-engine', 'my-secret');
    store.pushPayload('kv/metadata', {
      modelName: 'kv/metadata',
      ...metadata,
    });
    this.metadata = store.peekRecord('kv/metadata', metadata.id);
    this.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.metadata.backend, route: 'list' },
      { label: this.metadata.path, route: 'secret.details', model: this.metadata.path },
      { label: 'version diff' },
    ];

    // compare version 4
    const dataId = kvDataPath('kv-engine', 'my-secret', 4);
    store.pushPayload('kv/data', {
      modelName: 'kv/data',
      id: dataId,
      secret_data: { foo: 'bar' },
      created_time: '2023-07-20T02:12:17.379762Z',
      custom_metadata: null,
      deletion_time: '',
      destroyed: false,
      version: 4,
    });

    this.endpoint = `${encodePath('kv-engine')}/data/${'my-secret'}`;
  });

  test('it renders compared data of the two versions and shows icons for deleted, destroyed and current', async function (assert) {
    assert.expect(12);
    this.server.get(this.endpoint, (schema, req) => {
      assert.ok('request made to the correct endpoint.');
      assert.strictEqual(
        req.queryParams.version,
        '1',
        'request includes the version flag on queryRecord and correctly only queries for the record not in the store, which is retrieved by peekRecord.'
      );
      return EXAMPLE_KV_DATA_GET_RESPONSE;
    });

    await render(
      hbs`
       <Page::Secret::Metadata::VersionDiff
        @metadata={{this.metadata}} 
        @path={{this.metadata.path}}
        @backend={{this.metadata.backend}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );
    /* eslint-disable no-useless-escape */
    assert
      .dom('[data-test-visual-diff]')
      .hasText(
        `foo"bar"\hello"world"`,
        'correctly pull in the data from version 4 and compared to version 1.'
      );
    assert
      .dom('[data-test-version-dropdown-left]')
      .hasText('Version 4', 'shows the current version for the left side default version.');
    assert
      .dom('[data-test-version-dropdown-right]')
      .hasText(
        'Version 1',
        'shows the first version that is not deleted or destroyed for the right version on init.'
      );

    await click('[data-test-version-dropdown-right]');
    for (const version in this.metadata.versions) {
      const data = this.metadata.versions[version];
      assert.dom(PAGE.versions.button(version)).exists('renders the button for each version.');

      if (data.destroyed || data.deletion_time) {
        assert
          .dom(`${PAGE.versions.button(version)} [data-test-icon="x-square-fill"]`)
          .hasClass(`${data.destroyed ? 'has-text-danger' : 'has-text-grey'}`);
      }
    }
    assert
      .dom(`${PAGE.versions.button('1')}`)
      .hasClass('is-active', 'correctly shows the selected version 1 as active giving it text blue.');
  });
});
