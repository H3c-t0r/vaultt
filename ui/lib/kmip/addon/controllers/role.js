/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Controller from '@ember/controller';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';

export default class RoleController extends Controller {
  @service flashMessages;
  @service router;

  @action
  async deleteRole() {
    const { id } = this.model;
    try {
      await this.model.destroyRecord();
      this.flashMessages.success(`Successfully deleted role ${id}`);
      this.router.transitionTo('vault.cluster.secrets.backend.kmip.scope.roles', this.scope);
    } catch (e) {
      this.flashMessages.danger(`There was an error deleting the role ${id}: ${e.errors.join(' ')}`);
      this.model.rollbackAttributes();
    }
  }
}
