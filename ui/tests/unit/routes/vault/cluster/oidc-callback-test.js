/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, skip, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';
import { getParamsForCallback } from 'vault/routes/vault/cluster/oidc-callback';

module('Unit | Route | vault/cluster/oidc-callback', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.originalOpener = window.opener;
    window.opener = {
      postMessage: () => {},
    };
    this.route = this.owner.lookup('route:vault/cluster/oidc-callback');
    this.windowStub = sinon.stub(window.opener, 'postMessage');
    this.state = 'st_yOarDguU848w5YZuotLs';
    this.path = 'oidc';
    this.code = 'lTazRXEwKfyGKBUCo5TyLJzdIt39YniBJOXPABiRMkL0T';
    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: '' };
      return {
        auth_path: this.path,
        code: this.code,
      };
    };
    this.callbackUrlQueryParams = (stateParam) => {
      switch (stateParam) {
        case '':
          window.history.pushState({}, '');
          break;
        case 'stateless':
          window.history.pushState({}, '', '?' + `code=${this.code}`);
          break;
        default:
          window.history.pushState({}, '', '?' + `code=${this.code}&state=${stateParam}`);
          break;
      }
    };
  });

  hooks.afterEach(function () {
    this.windowStub.restore();
    window.opener = this.originalOpener;
    this.callbackUrlQueryParams('');
  });

  module('getParamsForCallback helper fn', function () {
    test('it parses params correctly with regular inputs no namespace', function (assert) {
      const params = {
        state: 'my-state',
        code: 'my-code',
        path: 'oidc-path',
      };
      const results = getParamsForCallback(params, '?code=my-code&state=my-state');
      assert.deepEqual(results, { source: 'oidc-callback', ...params });
    });

    test('it parses params correctly regular inputs and namespace param', function (assert) {
      const params = {
        state: 'my-state',
        code: 'my-code',
        path: 'oidc-path',
        namespace: 'my-namespace/nested',
      };
      const results = getParamsForCallback(params, '?code=my-code&state=my-state');
      assert.deepEqual(results, { source: 'oidc-callback', ...params });
    });

    test('it parses params correctly regular inputs and namespace in state', function (assert) {
      const queryString = '?code=my-code&state=my-state,ns=blah';
      const params = {
        state: 'my-state,ns', // mock what Ember does with the state QP
        code: 'my-code',
        path: 'oidc-path',
      };
      const results = getParamsForCallback(params, queryString);
      assert.deepEqual(results, { source: 'oidc-callback', ...params, state: 'my-state', namespace: 'blah' });
    });

    test('namespace in state takes precedence over namespace in route', function (assert) {
      const queryString = '?code=my-code&state=my-state,ns=some-ns/foo';
      const params = {
        state: 'my-state,ns', // mock what Ember does with the state QP
        code: 'my-code',
        path: 'oidc-path',
        namespace: 'some-ns',
      };
      const results = getParamsForCallback(params, queryString);
      assert.deepEqual(results, {
        source: 'oidc-callback',
        ...params,
        state: 'my-state',
        namespace: 'some-ns/foo',
      });
    });

    test('correctly decodes path and state', function (assert) {
      const searchString = '?code=my-code&state=spaces%20state,ns=some.ns%2Fchild';
      const params = {
        state: 'spaces state,ns',
        code: 'my-code',
        path: 'oidc-path',
      };
      const results = getParamsForCallback(params, searchString);
      assert.deepEqual(results, {
        source: 'oidc-callback',
        state: 'spaces state',
        code: 'my-code',
        path: 'oidc-path',
        namespace: 'some.ns/child',
      });
    });

    test('correctly decodes namespace in state', function (assert) {
      const queryString = '?code=my-code&state=spaces%20state,ns=spaces%20ns';
      const params = {
        state: 'spaces state,ns', // mock what Ember does with the state QP
        code: 'my-code',
        path: 'oidc-path',
      };
      const results = getParamsForCallback(params, queryString);
      assert.deepEqual(results, {
        source: 'oidc-callback',
        state: 'spaces state',
        code: 'my-code',
        path: 'oidc-path',
        namespace: 'spaces ns',
      });
    });

    test('it parses params correctly when window.location.search is empty', function (assert) {
      const params = {
        state: 'my-state',
        code: 'my-code',
        path: 'oidc-path',
        namespace: 'ns1',
      };
      const results = getParamsForCallback(params, '');
      assert.deepEqual(results, { source: 'oidc-callback', ...params, state: 'my-state' });
    });
  });

  test('it calls route', function (assert) {
    assert.ok(this.route);
  });

  skip('it uses namespace param from state instead of cluster, with custom oidc path', function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';
    this.callbackUrlQueryParams(encodeURIComponent(`${this.state},ns=test-ns`));
    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: 'admin' };
      return {
        auth_path: 'oidc-dev',
        code: this.code,
      };
    };
    this.route.afterModel();
    assert.propEqual(
      this.windowStub.lastCall.args[0],
      {
        code: this.code,
        path: 'oidc-dev',
        namespace: 'test-ns',
        state: this.state,
        source: 'oidc-callback',
      },
      'ns from state not cluster'
    );
  });

  skip('it uses namespace from cluster when state does not include ns param', function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';
    this.callbackUrlQueryParams(encodeURIComponent(this.state));
    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: 'admin' };
      return {
        auth_path: this.path,
        code: this.code,
      };
    };
    this.route.afterModel();
    assert.propEqual(
      this.windowStub.lastCall.args[0],
      {
        code: this.code,
        path: this.path,
        namespace: 'admin',
        state: this.state,
        source: 'oidc-callback',
      },
      `namespace is from cluster's namespaceQueryParam`
    );
  });

  skip('it correctly parses encoded, nested ns param from state', function (assert) {
    this.callbackUrlQueryParams(encodeURIComponent(`${this.state},ns=parent-ns/child-ns`));
    this.route.afterModel();
    assert.propEqual(
      this.windowStub.lastCall.args[0],
      {
        code: this.code,
        path: this.path,
        namespace: 'parent-ns/child-ns',
        state: this.state,
        source: 'oidc-callback',
      },
      'it has correct nested ns from state and sets as namespace param'
    );
  });

  skip('the afterModel hook returns when both cluster and route params are empty strings', function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';
    this.callbackUrlQueryParams('');
    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: '' };
      return {
        auth_path: '',
        code: '',
      };
    };
    this.route.afterModel();
    assert.propEqual(
      this.windowStub.lastCall.args[0],
      {
        path: '',
        state: '',
        code: '',
        source: 'oidc-callback',
      },
      'model hook returns with empty params'
    );
  });

  skip('the afterModel hook returns when state param does not exist', function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';
    this.callbackUrlQueryParams('stateless');
    this.route.afterModel();
    assert.propEqual(
      this.windowStub.lastCall.args[0],
      {
        code: this.code,
        path: 'oidc',
        state: '',
        source: 'oidc-callback',
      },
      'model hook returns empty string when state param nonexistent'
    );
  });

  skip('the afterModel hook returns when cluster ns exists and all route params are empty strings', function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';
    this.callbackUrlQueryParams('');
    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: 'ns1' };
      return {
        auth_path: '',
        code: '',
      };
    };
    this.route.afterModel();
    assert.propEqual(
      this.windowStub.lastCall.args[0],
      {
        code: '',
        namespace: 'ns1',
        path: '',
        source: 'oidc-callback',
        state: '',
      },
      'model hook returns with empty parameters'
    );
  });

  /*
  If authenticating to a namespace, most SSO providers return a callback url
  with a 'state' query param that includes a URI encoded namespace, example:
  '?code=BZBDVPMz0By2JTqulEMWX5-6rflW3A20UAusJYHEeFygJ&state=sst_yOarDguU848w5YZuotLs%2Cns%3Dadmin'

  Active Directory Federation Service (AD FS), instead, decodes the namespace portion:
  '?code=BZBDVPMz0By2JTqulEMWX5-6rflW3A20UAusJYHEeFygJ&state=st_yOarDguU848w5YZuotLs,ns=admin'

  'ns' isn't recognized as a separate param because there is no ampersand, so using this.paramsFor() returns
  a namespace-less state and authentication fails
  { state: 'st_yOarDguU848w5YZuotLs,ns' }
  */
  skip('it uses namespace when state param is not uri encoded', async function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';
    this.callbackUrlQueryParams(`${this.state},ns=admin`);
    this.route.afterModel();
    assert.propEqual(
      this.windowStub.lastCall.args[0],
      {
        code: this.code,
        namespace: 'admin',
        path: this.path,
        source: 'oidc-callback',
        state: this.state,
      },
      'namespace is parsed correctly'
    );
  });
});
