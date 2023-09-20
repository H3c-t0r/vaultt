/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfig } from 'pki/decorators/check-issuers';
import { hash } from 'rsvp';
import { getCliMessage } from 'pki/routes/overview';
@withConfig()
export default class PkiRolesIndexRoute extends Route {
  @service store;
  @service secretMountPath;

  queryParams = {
    pageFilter: {
      refreshModel: true,
    },
    currentPage: {
      refreshModel: true,
    },
  };

  async fetchRoles(params) {
    try {
      return await this.store.lazyPaginatedQuery('pki/role', {
        backend: this.secretMountPath.currentPath,
        responsePath: 'data.keys',
        page: Number(params.currentPage) || 1,
        pageFilter: params.pageFilter,
      });
    } catch (e) {
      if (e.httpStatus === 404) {
        return { parentModel: this.modelFor('roles') };
      }
      throw e;
    }
  }

  model(params) {
    return hash({
      hasConfig: this.shouldPromptConfig,
      roles: this.fetchRoles(params),
      parentModel: this.modelFor('roles'),
      pageFilter: params.pageFilter,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const roles = resolvedModel.roles;

    if (roles?.length) controller.notConfiguredMessage = getCliMessage('roles');
    else controller.notConfiguredMessage = getCliMessage();
  }

  resetController(controller, isExiting) {
    if (isExiting) {
      controller.set('pageFilter', undefined);
      controller.set('currentPage', undefined);
    }
  }
}
