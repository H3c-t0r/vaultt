import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { isNone } from '@ember/utils';
import { action } from '@ember/object';
import KVObject from 'vault/lib/kv-object';

/**
 * @module KvObjectEditor
 * KvObjectEditor components are called in FormFields when the editType on the model is kv.  They are used to show a key-value input field.
 *
 * @example
 * ```js
 * <KvObjectEditor
 *  @value={{get model valuePath}}
 *  @onChange={{action "setAndBroadcast" valuePath }}
 *  @label="some label"
   />
 * ```
 * @param {string} value - the value is captured from the model.
 * @param {function} onChange - function that captures the value on change
 * @param {function} [onKeyUp] - function passed in that handles the dom keyup event. Used for validation on the kv custom metadata.
 * @param {string} [label] - label displayed over key value inputs
 * @param {string} [labelClass] - override default label class in FormFieldLabel component
 * @param {string} [warning] - warning that is displayed
 * @param {string} [helpText] - helper text. In tooltip.
 * @param {string} [subText] - placed under label.
 * @param {string} [keyPlaceholder] - placeholder for key input
 * @param {string} [valuePlaceholder] - placeholder for value input
 */

export default class KvObjectEditor extends Component {
  @tracked kvData;

  constructor() {
    super(...arguments);
    this.kvData = KVObject.create({ content: [] }).fromJSON(this.args.value);
    this.addRow();
  }

  get placeholders() {
    return {
      key: this.args.keyPlaceholder || 'key',
      value: this.args.valuePlaceholder || 'value',
    };
  }
  get hasDuplicateKeys() {
    return this.kvData.uniqBy('name').length !== this.kvData.get('length');
  }

  @action
  addRow() {
    if (!isNone(this.kvData.findBy('name', ''))) {
      return;
    }
    this.kvData.addObject({ name: '', value: '' });
  }
  @action
  updateRow() {
    this.args.onChange(this.kvData.toJSON());
  }
  @action
  deleteRow(index) {
    this.kvData.removeAt(index);
    this.args.onChange(this.kvData.toJSON());
  }
  @action
  handleKeyUp(event) {
    if (this.args.onKeyUp) {
      this.args.onKeyUp(event.target.value);
    }
  }
}
