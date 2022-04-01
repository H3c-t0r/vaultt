import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | radio-button', function (hooks) {
  setupRenderingTest(hooks);

  test('it should spread attributes on input element', async function (assert) {
    await render(
      hbs(`
      <RadioButton
        id="foo"
        name="bar"
        class="radio"
        @onChange={{fn (mut this.groupValue)}}
      />
    `)
    );
    assert.dom('input').hasAttribute('id', 'foo', 'id set on input element');
    assert.dom('input').hasAttribute('name', 'bar', 'name set on input element');
    assert.dom('input').hasClass('radio', 'class set on input element');
  });

  test('it should be checked when value and groupValue are equal', async function (assert) {
    await render(
      hbs(`
      <RadioButton
        @value="foo"
        @groupValue="foo"
        @onChange={{fn (mut this.groupValue)}}
      />
    `)
    );
    assert.true(
      this.element.querySelector('input').checked,
      'input is checked when value matches groupValue'
    );
  });

  test('it should send onChange action and mark correct radio as checked', async function (assert) {
    await render(
      hbs(`
      <RadioButton
        @value="foo"
        @groupValue={{this.groupValue}}
        @onChange={{fn (mut this.groupValue)}}
        data-test-radio-1
      />
      <RadioButton
        @value="bar"
        @groupValue={{this.groupValue}}
        @onChange={{fn (mut this.groupValue)}}
        data-test-radio-2
      />
    `)
    );
    const radio1 = this.element.querySelector('[data-test-radio-1]');
    const radio2 = this.element.querySelector('[data-test-radio-2]');
    assert.false(radio1.checked, 'radio1 is unchecked when groupValue is undefined');
    assert.false(radio2.checked, 'radio2 is unchecked when groupValue is undefined');
    await click('[data-test-radio-1]');
    assert.true(radio1.checked, 'radio1 is checked');
    assert.false(radio2.checked, 'radio2 is unchecked');
    await click('[data-test-radio-2]');
    assert.false(radio1.checked, 'radio1 is unchecked');
    assert.true(radio2.checked, 'radio2 is checked');
  });
});
