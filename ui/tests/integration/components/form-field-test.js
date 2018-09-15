import EmberObject from '@ember/object';
import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';
import { create } from 'ember-cli-page-object';
import sinon from 'sinon';
import formFields from '../../pages/components/form-field';

const component = create(formFields);

moduleForComponent('form-field', 'Integration | Component | form field', {
  integration: true,
  beforeEach() {
    component.setContext(this);
  },

  afterEach() {
    component.removeContext();
  },
});

const createAttr = (name, type, options) => {
  return {
    name,
    type,
    options,
  };
};

const setup = function(attr) {
  let model = EmberObject.create({});
  let spy = sinon.spy();
  this.set('onChange', spy);
  this.set('model', model);
  this.set('attr', attr);
  this.render(hbs`{{form-field attr=attr model=model onChange=onChange}}`);
  return [model, spy];
};

test('it renders', function(assert) {
  let model = EmberObject.create({});
  this.set('attr', { name: 'foo' });
  this.set('model', model);
  this.render(hbs`{{form-field attr=attr model=model}}`);

  assert.equal(component.field.labelText, 'Foo', 'renders a label');
  assert.notOk(component.hasInput, 'renders only the label');
});

test('it renders: string', function(assert) {
  let [model, spy] = setup.call(this, createAttr('foo', 'string', { defaultValue: 'default' }));
  assert.equal(component.field.labelText, 'Foo', 'renders a label');
  assert.equal(component.field.inputValue, 'default', 'renders default value');
  assert.ok(component.hasInput, 'renders input for string');
  component.field.input('bar').change();

  assert.equal(model.get('foo'), 'bar');
  assert.ok(spy.calledWith('foo', 'bar'), 'onChange called with correct args');
});

test('it renders: boolean', function(assert) {
  let [model, spy] = setup.call(this, createAttr('foo', 'boolean', { defaultValue: false }));
  assert.equal(component.field.labelText, 'Foo', 'renders a label');
  assert.notOk(component.field.inputChecked, 'renders default value');
  assert.ok(component.hasCheckbox, 'renders a checkbox for boolean');
  component.field.clickLabel();

  assert.equal(model.get('foo'), true);
  assert.ok(spy.calledWith('foo', true), 'onChange called with correct args');
});

test('it renders: number', function(assert) {
  let [model, spy] = setup.call(this, createAttr('foo', 'number', { defaultValue: 5 }));
  assert.equal(component.field.labelText, 'Foo', 'renders a label');
  assert.equal(component.field.inputValue, 5, 'renders default value');
  assert.ok(component.hasInput, 'renders input for number');
  component.field.input(8).change();

  assert.equal(model.get('foo'), 8);
  assert.ok(spy.calledWith('foo', '8'), 'onChange called with correct args');
});

test('it renders: object', function(assert) {
  setup.call(this, createAttr('foo', 'object'));
  assert.equal(component.field.labelText, 'Foo', 'renders a label');
  assert.ok(component.hasJSONEditor, 'renders the json editor');
});

test('it renders: editType textarea', function(assert) {
  let [model, spy] = setup.call(
    this,
    createAttr('foo', 'string', { defaultValue: 'goodbye', editType: 'textarea' })
  );
  assert.equal(component.field.labelText, 'Foo', 'renders a label');
  assert.ok(component.hasTextarea, 'renders a textarea');
  assert.equal(component.field.textareaValue, 'goodbye', 'renders default value');
  component.field.textarea('hello');

  assert.equal(model.get('foo'), 'hello');
  assert.ok(spy.calledWith('foo', 'hello'), 'onChange called with correct args');
});

test('it renders: editType file', function(assert) {
  setup.call(this, createAttr('foo', 'string', { editType: 'file' }));
  assert.ok(component.hasTextFile, 'renders the text-file component');
});

test('it renders: editType ttl', function(assert) {
  let [model, spy] = setup.call(this, createAttr('foo', null, { editType: 'ttl' }));
  assert.ok(component.hasTTLPicker, 'renders the ttl-picker component');
  component.field.input('3');
  component.field.select('h').change();

  assert.equal(model.get('foo'), '3h');
  assert.ok(spy.calledWith('foo', '3h'), 'onChange called with correct args');
});

test('it renders: editType stringArray', function(assert) {
  let [model, spy] = setup.call(this, createAttr('foo', 'string', { editType: 'stringArray' }));
  assert.ok(component.hasStringList, 'renders the string-list component');

  component.field.input('array').change();
  assert.deepEqual(model.get('foo'), ['array'], 'sets the value on the model');
  assert.deepEqual(spy.args[0], ['foo', ['array']], 'onChange called with correct args');
});

test('it uses a passed label', function(assert) {
  setup.call(this, createAttr('foo', 'string', { label: 'Not Foo' }));
  assert.equal(component.field.labelText, 'Not Foo', 'renders the label from options');
});

test('it renders a help tooltip', function(assert) {
  setup.call(this, createAttr('foo', 'string', { helpText: 'Here is some help text' }));
  assert.ok(component.hasTooltip, 'renders the tooltip component');
  component.tooltipTrigger();
});
