/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfirmLeave } from 'core/decorators/confirm-leave';

import type StoreService from 'vault/services/store';

interface Params {
  type: string;
}

@withConfirmLeave()
export default class SyncSecretsDestinationsCreateDestinationRoute extends Route {
  @service declare readonly store: StoreService;

  model(params: Params) {
    const { type } = params;
    return this.store.createRecord(`sync/destinations/${type}`, { type });
  }
}
