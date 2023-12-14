/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import errorMessage from 'vault/utils/error-message';
import { inject as service } from '@ember/service';

/**
 * @module Page::CreateAndEditMessageForm
 * Page::CreateAndEditMessageForm components are used to display create and edit message form fields.
 * @example
 * ```js
 * <Page::CreateAndEditMessageForm @message={{this.message}}  />
 * ```
 * @param {model} message - message model to pass to form components
 */

export default class MessagesList extends Component {
  @service router;
  @service flashMessages;

  @tracked errorBanner = '';
  @tracked modelValidations;
  @tracked invalidFormMessage;

  willDestroy() {
    super.willDestroy();
    const noTeardown = this.store && !this.store.isDestroying;
    const { model } = this;
    if (noTeardown && model && model.get('isDirty') && !model.isDestroyed && !model.isDestroying) {
      model.rollbackAttributes();
    }
  }

  get breadcrumbs() {
    const authenticated =
      this.args.message.authenticated === undefined ? true : this.args.message.authenticated;
    return [
      { label: 'Messages', route: 'messages.index', query: { authenticated } },
      { label: 'Create Message' },
    ];
  }

  @task
  *save(event) {
    event.preventDefault();
    try {
      const { isValid, state, invalidFormMessage } = this.args.message.validate();
      this.modelValidations = isValid ? null : state;
      this.invalidFormAlert = invalidFormMessage;

      if (isValid) {
        const { isNew } = this.args.message;

        // We do these checks here since there could be a scenario where startTime and endTime are strings.
        // The model expects these attrs to be a date object, so we will need to update these attrs to be in
        // date object format.
        if (typeof this.args.message.startTime === 'string')
          this.args.message.startTime = new Date(this.args.message.startTime);
        if (typeof this.args.message.endTime === 'string')
          this.args.message.endTime = new Date(this.args.message.endTime);

        const { id } = yield this.args.message.save();
        this.flashMessages.success(`Successfully ${isNew ? 'created' : 'updated'} the message.`);
        this.router.transitionTo('vault.cluster.config-ui.messages.message.details', id);
      }
    } catch (error) {
      this.errorBanner = errorMessage(error);
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }
}
