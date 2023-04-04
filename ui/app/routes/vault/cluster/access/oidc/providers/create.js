/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class OidcProvidersCreateRoute extends Route {
  @service store;

  model() {
    return this.store.createRecord('oidc/provider');
  }
}
