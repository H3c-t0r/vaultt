/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AdapterError from '@ember-data/adapter/error';
import AuthConfigComponent from './config';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { tracked } from '@glimmer/tracking';
import errorMessage from 'vault/utils/error-message';

/**
 * @module AuthConfigForm/Options
 * The `AuthConfigForm/Options` is options portion of the auth config form.
 *
 * @example
 * ```js
 * <AuthConfigForm::Options @modle={{this.args.model.model}} />
 * ```
 *
 * @property model=null {DS.Model} - The corresponding auth model that is being configured.
 *
 */

export default class AuthConfigOptions extends AuthConfigComponent {
  @service flashMessages;
  @service router;

  @tracked errorMessage;

  @task
  @waitFor
  *saveModel(evt) {
    evt.preventDefault();
    this.error = null;
    const data = this.args.model.config.serialize();
    data.description = this.args.model.description;

    if (this.args.model.supportsUserLockoutConfig) {
      data.user_lockout_config = {};
      this.args.model.userLockoutConfig.apiParams.forEach((attr) => {
        if (Object.keys(data).includes(attr)) {
          data.user_lockout_config[attr] = data[attr];
          delete data[attr];
        }
      });
    }

    // token_type should not be tuneable for the token auth method.
    if (this.args.model.methodType === 'token') {
      delete data.token_type;
    }

    try {
      yield this.args.model.tune(data);
    } catch (err) {
      // AdapterErrors are handled by the error-message component
      // in the form
      if (err instanceof AdapterError === false) {
        throw err;
      }
      // because we're not calling model.save the model never updates with
      // the error. Set it on the component instead.
      this.errorMessage = errorMessage(err);
      return;
    }
    this.router.transitionTo('vault.cluster.access.methods').followRedirects();
    this.flashMessages.success('The configuration was saved successfully.');
  }
}
