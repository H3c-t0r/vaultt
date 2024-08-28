/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AdapterError from '@ember-data/adapter/error';
import { set } from '@ember/object';
import Route from '@ember/routing/route';
import type Controller from '@ember/controller';
import { service } from '@ember/service';
import { CONFIGURABLE_SECRET_ENGINES } from 'vault/helpers/mountable-secret-engines';
import errorMessage from 'vault/utils/error-message';
import { action } from '@ember/object';

import type Store from '@ember-data/store';
import type SecretEngineModel from 'vault/models/secret-engine';
import type Transition from '@ember/routing/transition';
import type { Breadcrumb } from 'vault/vault/app-types';

interface SecretsConfigurationEditController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: Record<string, string>;
  type: string;
  id: string;
}

// This route file is reused for all configurable secret engines.
// It generates config models based on the engine type.
// Saving and updating of those models are done within the engine specific components.

const CONFIG_ADAPTERS_PATHS: Record<string, string[]> = {
  aws: ['aws/lease-config', 'aws/root-config'],
  ssh: ['ssh/ca-config'],
};

export default class SecretsBackendConfigurationEdit extends Route {
  @service declare readonly store: Store;

  async model() {
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    const secretEngineRecord = this.modelFor('vault.cluster.secrets.backend') as SecretEngineModel;
    const type = secretEngineRecord.type;

    // if the engine type is not configurable, return a 404.
    if (!secretEngineRecord || !CONFIGURABLE_SECRET_ENGINES.includes(type)) {
      const error = new AdapterError();
      set(error, 'httpStatus', 404);
      throw error;
    }
    // generate the model based on the engine type.
    // and pre-set model with type and backend e.g. {type: ssh, id: ssh-123}
    const model: Record<string, unknown> = { type, id: backend };
    for (const adapterPath of CONFIG_ADAPTERS_PATHS[type] as string[]) {
      // convert the adapterPath with a name that can be passed to the components
      // ex: adapterPath = ssh/ca-config, convert to: ssh-ca-config so that you can pass to component @model={{this.model.ssh-ca-config}}
      const standardizedKey = adapterPath.replace(/\//g, '-');
      try {
        model[standardizedKey] = await this.store.queryRecord(adapterPath, {
          backend,
          type,
        });
      } catch (e: AdapterError) {
        // For most models if the adapter returns a 404, we want to create a new record.
        // The ssh secret engine however returns a 400 if the CA is not configured.
        // For ssh's 400 error, we want to create the CA config model.
        if (
          e.httpStatus === 404 ||
          (type === 'ssh' && e.httpStatus === 400 && errorMessage(e) === `keys haven't been configured yet`)
        ) {
          model[standardizedKey] = await this.store.createRecord(adapterPath, {
            backend,
            type,
          });
        } else {
          throw e;
        }
      }
    }
    return model;
  }

  @action
  willTransition() {
    // catch the transition and refresh model so the route shows the most recent model data.
    this.refresh();
  }

  setupController(
    controller: SecretsConfigurationEditController,
    resolvedModel: { backend: string; type: string; id: string },
    transition: Transition
  ) {
    super.setupController(controller, resolvedModel, transition);
    const routeStub = 'vault.cluster.secrets.backend';
    const label = resolvedModel.id || resolvedModel.backend;

    controller.breadcrumbs = [
      { label: 'Secrets', route: 'vault.cluster.secrets.backends' },
      { label: label, route: routeStub, model: resolvedModel.id }, // to landing page of the engine (ex: roles or secrets list)
      { label: 'Configuration', route: `${routeStub}.configuration.index`, model: resolvedModel.id },
      { label: 'Edit' },
    ];
  }
}
