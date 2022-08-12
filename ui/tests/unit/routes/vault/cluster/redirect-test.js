import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';

module('Unit | Route | vault/cluster/redirect', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.router = this.owner.lookup('service:router');
    this.originalTransition = this.router.transitionTo;
    this.router.replaceWith = sinon.stub().returns({
      followRedirects: function () {
        return {
          then: function (callback) {
            callback();
          },
        };
      },
    });
  });

  hooks.afterEach(function () {
    this.router.transitionTo = this.originalTransition;
  });

  test('it calls route', function (assert) {
    let route = this.owner.lookup('route:vault/cluster/redirect');
    assert.ok(route);
  });

  test('it redirects to auth when unauthenticated', function (assert) {
    let route = this.owner.lookup('route:vault/cluster/redirect');
    const auth = this.owner.lookup('service:auth');
    const originalToken = auth.currentToken;

    auth.currentToken = null;

    route.beforeModel({ to: { queryParams: { redirect_to: 'vault/cluster/tools' } } });

    assert.true(
      this.router.transitionTo.calledWith('vault.cluster.auth'),
      'transitions to auth when not authenticated'
    );
    auth.currentToken = originalToken;
  });

  test('it redirects to cluster when authenticated without redirect param', function (assert) {
    let route = this.owner.lookup('route:vault/cluster/redirect');
    const auth = this.owner.lookup('service:auth');
    const originalToken = auth.currentToken;

    auth.currentToken = 's.xxxxxxxxx';

    route.beforeModel({ to: { queryParams: { foo: 'bar' } } });

    assert.true(
      this.router.transitionTo.calledWith('vault.cluster'),
      'transitions to cluster when not authenticated'
    );
    auth.currentToken = originalToken;
  });

  test('it redirects to desired path when authenticated with redirect param', function (assert) {
    let route = this.owner.lookup('route:vault/cluster/redirect');
    const auth = this.owner.lookup('service:auth');
    const originalToken = auth.currentToken;

    auth.currentToken = 's.xxxxxxxxx';

    route.beforeModel({ to: { queryParams: { redirect_to: 'vault/cluster/tools' } } });

    assert.true(
      this.router.transitionTo.calledWith('vault/cluster/tools'),
      'transitions to redirect_to path when authenticated'
    );
    auth.currentToken = originalToken;
  });
});
