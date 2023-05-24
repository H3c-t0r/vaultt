/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Controller from '@ember/controller';
import { dropTask, task } from 'ember-concurrency';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class VaultClusterAccessMethodsController extends Controller {
  @service flashMessages;

  @tracked authMethodOptions = [];
  @tracked selectedAuthType = null;
  @tracked selectedAuthName = null;

  queryParams = ['page, pageFilter'];

  page = 1;
  pageFilter = null;
  filter = null;

  get authMethodList() {
    // return an options list to filter by engine type, ex: 'kv'
    if (this.selectedAuthType) {
      // check first if the user has also filtered by name.
      if (this.selectedAuthName) {
        return this.model.filter((method) => this.selectedAuthName === method.id);
      }
      // otherwise filter by auth type
      return this.model.filter((method) => this.selectedAuthType === method.type);
    }
    // return an options list to filter by auth name, ex: 'my-userpass'
    if (this.selectedAuthName) {
      return this.model.filter((method) => this.selectedAuthName === method.id);
    }
    // no filters, return full sorted list.
    return this.model;
  }

  get authMethodArrayByType() {
    const arrayOfAllAuthTypes = this.model.map((modelObject) => modelObject.type);
    // filter out repeated auth types (e.g. [userpass, userpass] => [userpass])
    const arrayOfUniqueAuthTypes = [...new Set(arrayOfAllAuthTypes)];

    return arrayOfUniqueAuthTypes.map((authType) => ({
      name: authType,
      id: authType,
    }));
  }

  get authMethodArrayByName() {
    return this.model.map((modelObject) => ({
      name: modelObject.id,
      id: modelObject.id,
    }));
  }

  @action
  filterAuthType([type]) {
    this.selectedAuthType = type;
  }

  @action
  filterAuthName([name]) {
    this.selectedAuthName = name;
  }

  @task
  @dropTask
  *disableMethod(method) {
    const { type, path } = method;
    try {
      yield method.destroyRecord();
      this.flashMessages.success(`The ${type} Auth Method at ${path} has been disabled.`);
    } catch (err) {
      this.flashMessages.danger(
        `There was an error disabling Auth Method at ${path}: ${err.errors.join(' ')}.`
      );
    }
  }
}
