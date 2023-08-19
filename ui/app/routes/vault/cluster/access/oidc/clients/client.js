/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class OidcClientRoute extends Route {
  @service store;

  model({ name }) {
    return this.store.findRecord('oidc/client', name);
  }
}
