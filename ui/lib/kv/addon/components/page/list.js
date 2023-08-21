/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { getOwner } from '@ember/application';
import { ancestorKeysForKey } from 'core/utils/key-utils';
import errorMessage from 'vault/utils/error-message';
import { pathIsDirectory } from 'kv/utils/kv-breadcrumbs';

/**
 * @module List
 * ListPage component is a component to show a list of kv/metadata secrets.
 *
 * @param {array} secrets - An array of models generated form kv/metadata query.
 * @param {string} backend - The name of the kv secret engine.
 * @param {string} pathToSecret - The directory name that the secret belongs to ex: beep/boop/
 * @param {string} pageFilter - The input on the kv-list-filter. Does not include a directory name.
 * @param {string} filterValue - The concatenation of the pathToSecret and pageFilter ex: beep/boop/my-
 * @param {boolean} noMetadataListPermissions - true if the return to query metadata LIST is 403, indicating the user does not have permissions to that endpoint.
 * @param {array} breadcrumbs - Breadcrumbs as an array of objects that contain label, route, and modelId. They are updated via the util kv-breadcrumbs to handle dynamic *pathToSecret on the list-directory route.
 * @param {string} routeName - Either list or list-directory.
 * @param {object} meta - Object with values needed for pagination, created by LazyPaginatedQuery on the store service.
 */

export default class KvListPageComponent extends Component {
  @service flashMessages;
  @service router;
  @service store;

  @tracked secretPath;

  get mountPoint() {
    // mountPoint tells transition where to start. In this case, mountPoint will always be vault.cluster.secrets.backend.kv.
    return getOwner(this).mountPoint;
  }

  get pageSizes() {
    const { total, pageSize } = this.args.meta;
    const increments = [1, 15, 30, 50, 100];
    const truncated = increments.filter((num) => num > pageSize && num < total);
    // The pageSize value must be apart of the array of options to choose from. Add it back in and sort.
    truncated.push(pageSize);
    truncated.sort((a, b) => a - b);
    return truncated;
  }

  get showPagination() {
    // only show pagination if total # of secrets are larger than the smaller of: a custom pageSize; if set || 15
    const { total, pageSize } = this.args.meta;
    return total >= Math.min(15, pageSize) ? true : false;
  }

  // callback from HDS pagination to set queryParams currentPage and currentPageSize
  get paginationQueryParams() {
    // used to set page and pageSize to the queryParams from HDS::Pagination component
    return (page, pageSize) => {
      return {
        currentPage: page,
        currentPageSize: pageSize,
      };
    };
  }

  get buttonText() {
    return pathIsDirectory(this.secretPath) ? 'View directory' : 'View secret';
  }

  @action
  async onDelete(model) {
    try {
      // The model passed in is a kv/metadata model
      await model.destroyRecord();
      this.store.clearDataset('kv/metadata'); // Clear out the store cache so that the metadata/list view is updated.
      const message = `Successfully deleted the metadata and all version data of the secret ${model.fullSecretPath}.`;
      this.flashMessages.success(message);
      // if you've deleted a secret from within a directory, transition to its parent directory.
      if (this.args.routeName === 'list-directory') {
        const ancestors = ancestorKeysForKey(model.fullSecretPath);
        const nearest = ancestors.pop();
        this.router.transitionTo(`${this.mountPoint}.list-directory`, nearest);
      } else {
        // still need to fire off a transition to refresh the model
        this.router.transitionTo(`${this.mountPoint}.list`);
      }
    } catch (error) {
      const message = errorMessage(error, 'Error deleting secret. Please try again or contact support.');
      this.flashMessages.danger(message);
    }
  }

  @action
  handleSecretPathInput(value) {
    this.secretPath = value;
  }

  @action
  transitionToSecretDetail() {
    pathIsDirectory(this.secretPath)
      ? this.router.transitionTo('vault.cluster.secrets.backend.kv.list-directory', this.secretPath)
      : this.router.transitionTo('vault.cluster.secrets.backend.kv.secret.details', this.secretPath);
  }
}
