/**
 * @module SecretEdit
 * SecretEdit component manages the secret and model data, and displays either the create, update, empty state or show view of a KV secret.
 *
 * @example
 * ```js
 * <SecretEdit @model={{model}} @isV2/>
 * ```
/
 * @param {object} model - Model returned from secret-v2 which is generated in the secret-edit route
 * @param {string} mode - Edit, create, etc.
 * @param {string} basekey - For navigation.
 * @param {string} key - ARG TODO
 * @param {string} initalKey - ARG TODO
 * @param {function} onRefresh - ARG TODO
 * @param {function} onToggleAdvancedEdit - ARG TODO
 * @param {boolean} preferAdvancedEdit - ARG TODO
 This component is initialized from the secret-edit-layout.hbs file
 */

import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
// ARGTODO work on
// import FocusOnInsertMixin from 'vault/mixins/focus-on-insert';
// import WithNavToNearestAncestor from 'vault/mixins/with-nav-to-nearest-ancestor';
import KVObject from 'vault/lib/kv-object';
import { maybeQueryRecord } from 'vault/macros/maybe-query-record';
// https://stackoverflow.com/questions/60843983/does-ember-octane-route-class-support-using-mixins
export default class SecretEdit extends Component {
  // export default Component.extend(FocusOnInsertMixin, WithNavToNearestAncestor, {
  @service wizard;
  @service store;

  @tracked secretData = null;

  // called with a bool indicating if there's been a change in the secretData and customMetadata
  onDataChange() {}
  onRefresh() {}
  onToggleAdvancedEdit() {}

  @tracked
  isV2 = false;

  // use a named action here so we don't have to pass one in
  // this will bubble to the route
  // toggleAdvancedEdit = 'toggleAdvancedEdit';

  @tracked
  codemirrorString = null;

  constructor() {
    super(...arguments);
    let secrets = this.args.model.secretData;
    if (!secrets && this.args.model.selectedVersion) {
      this.isV2 = true;
      secrets = this.args.model.belongsTo('selectedVersion').value().secretData;
    }
    const data = KVObject.create({ content: [] }).fromJSON(secrets);
    this.secretData = data;
    this.codemirrorString = data.toJSONString();
    // if (data.isAdvanced()) {
    //   // ARG TODO cannot set an args??
    //   this.args.preferAdvancedEdit = true;
    // }
    if (this.wizard.featureState === 'details' && this.args.mode === 'create') {
      let engine = this.args.model.backend.includes('kv') ? 'kv' : this.args.model.backend;
      this.wizard.transitionFeatureMachine('details', 'CONTINUE', engine);
    }
  }

  checkSecretCapabilities() {
    return maybeQueryRecord(
      'capabilities',
      (context) => {
        // ARG TODO check context works here
        if (!context.model || context.mode === 'create') {
          return;
        }
        let backend = context.isV2 ? context.get('model.engine.id') : context.model.backend;
        let id = context.model.id;
        let path = context.isV2 ? `${backend}/data/${id}` : `${backend}/${id}`;
        return {
          id: path,
        };
      },
      'isV2',
      'model',
      'model.id',
      'mode'
    );
  }
  // these where alias
  @tracked
  canUpdateSecretData = this.checkSecretCapabilities.canUpdate;
  @tracked
  canReadSecretData = this.checkSecretCapabilities.canRead;

  checkMetadataCapabilities() {
    return maybeQueryRecord(
      'capabilities',
      (context) => {
        // ARG TODO check context work here
        if (!context.model || !context.isV2) {
          return;
        }
        let backend = context.model.backend;
        let path = `${backend}/metadata/`;
        return {
          id: path,
        };
      },
      'isV2',
      'model',
      'model.id',
      'mode'
    );
  }
  @tracked
  canDeleteSecretMetadata = this.checkMetadataCapabilities.canDelete;
  @tracked
  canUpdateSecretMetadata = this.checkMetadataCapabilities.canUpdate;
  @tracked
  canReadSecretMetadata = this.checkMetadataCapabilities.canRead;
  // this was an or
  get requestInFlight() {
    return this.args.model.isLoading || this.args.model.isReloading || this.args.model.isSaving;
  }
  // this was an or
  get buttonDisabled() {
    return this.requestInFlight || this.args.model.isFolder || this.args.model.flagsIsInvalid;
  }

  get modelForData() {
    let { model } = this.args;
    if (!model) return null;
    return this.isV2 ? model.belongsTo('selectedVersion').value() : model;
  }

  get basicModeDisabled() {
    return this.secretDataIsAdvanced || this.showAdvancedMode === false;
  }

  get secretDataAsJSON() {
    return this.secretData.toJSON();
  }

  get secretDataIsAdvanced() {
    return this.secretData.isAdvanced();
  }

  get showAdvancedMode() {
    return this.secretDataIsAdvanced || this.preferAdvancedEdit;
  }

  get isWriteWithoutRead() {
    if (!this.args.model) {
      return null;
      // ARG TODO was a return instead of null???
    }
    // if the version couldn't be read from the server
    if (this.args.isV2 && this.modelForData.failedServerRead) {
      return true;
    }
    // if the model couldn't be read from the server
    if (!this.args.isV2 && this.args.model.failedServerRead) {
      return true;
    }
    return false;
  }

  @action
  refresh() {
    this.args.nRefresh();
  }

  @action
  toggleAdvanced(bool) {
    this.args.onToggleAdvancedEdit(bool);
  }
}
