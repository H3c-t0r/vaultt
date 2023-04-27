/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfirmLeave } from 'core/decorators/confirm-leave';
import { hash } from 'rsvp';

@withConfirmLeave('model.config', ['model.urls'])
export default class PkiConfigurationCreateRoute extends Route {
  @service secretMountPath;
  @service store;

  model() {
    return hash({
      config: this.store.createRecord('pki/action'),
      urls: this.modelFor('configuration').urls,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'configure' },
    ];
  }
}
