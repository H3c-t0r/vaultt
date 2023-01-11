import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { assert } from '@ember/debug';

/**
 * @module ObjectListInput
 * ObjectListInput components are used to render a variable number of text inputs,
 * in a single row with an "Add" button beside the last input of the row.
 * Clicking 'add' generates a new row of empty inputs. The inputs are based on the @objectKeys
 * array, in which each object corresponds to an input field in the row.
 * sample object:
 *   {
 *     label: 'Input label', // required key
 *     key: 'attrKey', // required key
 *     type: 'string',
 *     placeholder: 'Enter input here'
 *   }
 *
 * @example
 * ```js
 * <ObjectListInput @objectKeys={{this.arrayOfObjects}} @onChange={{this.handleChange}} @inputValue={{this.inputValue}}/>
 * ```
 * @param {array} objectKeys - an array of objects (sample above), the length determines the number of columns the component renders
 * @callback onChange - callback triggered when any input changes or when a row is deleted, called with array of objects containing each input's key and value ex: [ { attrKey: 'some input value' } ]
 * @param {string} [inputValue] - an array of objects to pre-fill the component inputs
 */

export default class ObjectListInput extends Component {
  @tracked inputList = [];
  @tracked inputKeys;
  @tracked disableAdd = true;

  constructor() {
    super(...arguments);
    const requiredKeys = ['label', 'key'];
    if (!this.args.objectKeys.every((obj) => requiredKeys.every((k) => Object.keys(obj).includes(k)))) {
      assert(`objects in the @objectKeys array must include keys called: ${requiredKeys.join(', ')}`);
    }
    this.inputKeys = this.args.objectKeys.map((e) => e.key);
    const emptyRow = this.createEmptyRow(this.inputKeys);
    this.inputList = this.args.inputValue ? [...this.args.inputValue, emptyRow] : [emptyRow];
  }

  createEmptyRow(keys) {
    // create a new object with from input keys with empty values
    return Object.fromEntries(keys.map((key) => [key, '']));
  }

  @action
  handleInput(idx, { target }) {
    const inputObj = this.inputList.objectAt(idx);
    inputObj[target.name] = target.value;
    this.handleChange();
  }

  @action
  addRow() {
    const emptyRow = this.createEmptyRow(this.inputKeys);
    // this.inputList.pushObject(emptyRow);
    this.inputList = [...this.inputList, emptyRow];
    this.disableAdd = true;
  }

  @action
  removeRow(idx) {
    const row = this.inputList.objectAt(idx);
    this.inputList.removeObject(row);
    this.handleChange();
  }

  @action
  handleChange() {
    // disable/enable "add" button based on last row
    const lastObject = this.inputList[this.inputList.length - 1];
    this.disableAdd = Object.values(lastObject).any((input) => input === '') ? true : false;

    // don't send an empty last object to parent
    if (Object.values(lastObject).every((input) => input === '')) {
      this.args.onChange(this.inputList.slice(0, -1));
    } else {
      this.args.onChange(this.inputList);
    }
  }
}
