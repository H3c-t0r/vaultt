/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfirmLeave } from 'core/decorators/confirm-leave';

@withConfirmLeave()
export default class PkiIssuersRotateRootRoute extends Route {
  @service secretMountPath;
  @service store;

  model() {
    return this.store.createRecord('pki/action', { actionType: 'rotate-root' });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'issuers', route: 'issuers.index' },
      { label: 'rotate root' },
    ];
  }
}
