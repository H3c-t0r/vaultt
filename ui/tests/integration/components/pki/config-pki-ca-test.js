/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { resolve } from 'rsvp';
import EmberObject from '@ember/object';
import Service from '@ember/service';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { create } from 'ember-cli-page-object';
import configPki from 'vault/tests/pages/components/pki/config-pki-ca';
import apiStub from 'vault/tests/helpers/noop-all-api-requests';

const component = create(configPki);

const storeStub = Service.extend({
  createRecord(type, args) {
    return EmberObject.create(args, {
      save() {
        return resolve(this);
      },
      destroyRecord() {},
      send() {},
      unloadRecord() {},
    });
  },
});

module('Integration | Component | config pki ca', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.server = apiStub();
    this.owner.lookup('service:flash-messages').registerTypes(['success']);
    this.owner.register('service:store', storeStub);
    this.storeService = this.owner.lookup('service:store');
  });

  hooks.afterEach(function () {
    this.server.shutdown();
  });

  const config = function (pem) {
    return EmberObject.create({
      pem: pem,
      backend: 'pki',
      caChain: 'caChain',
      der: new Blob(['der'], { type: 'text/plain' }),
    });
  };

  const setupAndRender = async function (context, onRefresh) {
    const refreshFn = onRefresh || function () {};
    context.set('config', config());
    context.set('onRefresh', refreshFn);
    await context.render(hbs`<Pki::ConfigPkiCa @onRefresh={{this.onRefresh}} @config={{this.config}} />`);
  };

  test('it renders, no pem', async function (assert) {
    await setupAndRender(this);

    assert.notOk(component.hasTitle, 'no title in the default state');
    assert.strictEqual(component.replaceCAText, 'Configure CA');
    assert.strictEqual(component.downloadLinks.length, 0, 'there are no download links');

    await component.replaceCA();
    assert.strictEqual(component.title, 'Configure CA Certificate');
    await component.back();

    await component.setSignedIntermediateBtn();
    assert.strictEqual(component.title, 'Set signed intermediate');
  });

  test('it renders, with pem', async function (assert) {
    const c = config('pem');
    this.set('config', c);
    await render(hbs`<Pki::ConfigPkiCa @config={{this.config}} />`);
    assert.notOk(component.hasTitle, 'no title in the default state');
    assert.strictEqual(component.replaceCAText, 'Add CA');
    assert.strictEqual(component.downloadLinks.length, 3, 'shows download links');
  });
});
