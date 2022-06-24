import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import handleHasManySelection from 'core/utils/search-select-has-many';

/**
 * @module OidcAssignmentForm
 * OidcAssignmentForm components are used to...
 *
 * @example
 * ```js
 * <OidcAssignmentForm @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

export default class OidcAssignmentForm extends Component {
  @service store;
  @service flashMessages;

  @tracked modelErrors;

  @task
  *save() {
    this.modelErrors = {};
    // check validity state first and abort if invalid
    const { isValid, state } = this.args.model.validate();
    if (!isValid) {
      this.modelErrors = state;
    } else {
      try {
        yield this.args.model.save();
        this.args.onSave();
      } catch (error) {
        const message = error.errors ? error.errors.join('. ') : error.message;
        this.flashMessages.danger(message);
      }
    }
  }

  @action
  cancel() {
    // revert model changes
    this.args.model.rollbackAttributes();
    this.args.onClose();
  }

  @action
  handleOperation(e) {
    let value = e.target.value;
    this.args.model.name = value;
  }

  async onEntitiesSelect(selectedIds) {
    const entityIds = await this.args.model.entityIds;
    handleHasManySelection(selectedIds, entityIds, this.store, 'entityIds');
  }
}
