import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const COUNTERS = [
  {
    start_time: '2019-05-01T00:00:00Z',
    total: 50,
  },
  {
    start_time: '2019-04-01T00:00:00Z',
    total: 45,
  },
  {
    start_time: '2019-03-01T00:00:00Z',
    total: 55,
  },
];

module('Integration | Component | http-requests-table', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('counters', COUNTERS);
  });

  test('it renders', async function(assert) {
    await render(hbs`<HttpRequestsTable @counters={{counters}}/>`);

    assert.dom('.http-requests-table').exists();
  });

  test('it does not show Change column with less than one month of data', async function(assert) {
    const one_month_counter = [
      {
        start_time: '2019-05-01T00:00:00Z',
        total: 50,
      },
    ];
    await render(hbs`<HttpRequestsTable @counters={{one_month_counter}}/>`);

    assert.dom('.http-requests-table').exists();
    assert.dom('[data-test-change]').doesNotExist();
  });

  test('it shows Change column for more than one month of data', async function(assert) {
    await render(hbs`<HttpRequestsTable @counters={{counters}}/>`);

    assert.dom('[data-test-change]').exists();
  });

  test('it shows the percent change between each time window', async function(assert) {
    await render(hbs`<HttpRequestsTable @counters={{counters}}/>`);
    const expected = (((COUNTERS[1].total - COUNTERS[0].total) / COUNTERS[1].total) * 100).toFixed(1);

    assert.ok(this.element.textContent.includes(expected));
  });
});
