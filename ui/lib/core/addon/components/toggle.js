/**
 * @module Toggle
 * Toggle components are used to indicate boolean values which can be toggled on or off.
 * They are a stylistic alternative to checkboxes, but still use the input[type="checkbox"] under the hood.
 *
 * @example
 * ```js
 * <Toggle @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {function} onChange - onChange is triggered on checkbox change (select, deselect). Must manually mutate checked value
 * @param {string} name - name is passed along to the form field, as well as to generate the ID of the input & "for" value of the label
 * @param {boolean} [checked=false] - checked status of the input, and must be passed in and mutated from the parent
 * @param {boolean} [disabled=false] - disabled makes the switch unclickable
 * @param {string} [size='medium'] - Sizing can be small or medium
 * @param {string} [status='normal'] - Status can be normal or success, which makes the switch have a blue background when checked=true
 */

import Component from '@glimmer/component';
import layout from '../templates/components/toggle';
import { setComponentTemplate } from '@ember/component';
import { action } from '@ember/object';

class ToggleComponent extends Component {
  get checked() {
    return this.args.checked || false;
  }

  get disabled() {
    return this.args.disabled || false;
  }

  get name() {
    return this.args.name || '';
  }

  get safeId() {
    return `toggle-${this.name.replace(/\W/g, '')}`;
  }
  get inputClasses() {
    let size = this.args.size || 'normal';
    let status = this.args.status || 'normal';
    const sizeClass = `is-${size}`;
    const statusClass = `is-${status}`;
    return `toggle ${statusClass} ${sizeClass}`;
  }

  @action
  handleChange(value) {
    this.args.onChange(value);
  }
}

export default setComponentTemplate(layout, ToggleComponent);
