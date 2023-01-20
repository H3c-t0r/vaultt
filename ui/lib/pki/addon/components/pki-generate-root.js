import { action } from '@ember/object';
import { service } from '@ember/service';
import { waitFor } from '@ember/test-waiters';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { keyParamsByType } from 'pki/utils/action-params';
import errorMessage from 'vault/utils/error-message';

/**
 * @module PkiGenerateRoot
 * PkiGenerateRoot shows only the fields valid for the generate root endpoint.
 * This form handles the model save and rollback actions, and will call the passed
 * onSave and onCancel args for transition (passed from parent).
 * NOTE: this component is not TS because decorator-added parameters (eg validator and
 * formFields) aren't recognized on the model.
 *
 * @example
 * ```js
 * <PkiGenerateRoot @model={{this.model}} @onCancel={{transition-to "vault.cluster"}} @onSave={{transition-to "vault.cluster.secrets"}} @adapterOptions={{hash actionType="import" useIssuer=false}} />
 * ```
 *
 * @param {Object} model - pki/action model.
 * @callback onCancel - Callback triggered when cancel button is clicked, after model is unloaded
 * @param {Object} adapterOptions - object passed as adapterOptions on the model.save method
 */
export default class PkiGenerateRootComponent extends Component {
  @service flashMessages;
  @service router;
  @tracked showGroup = null;
  @tracked modelValidations = null;
  @tracked errorBanner = '';
  @tracked invalidFormAlert = '';

  @action
  toggleGroup(group, isOpen) {
    this.showGroup = isOpen ? group : null;
  }

  get defaultFields() {
    return [
      'type',
      'commonName',
      'issuerName',
      'customTtl',
      'notBeforeDuration',
      'format',
      'permittedDnsDomains',
      'maxPathLength',
    ];
  }
  get keyParamFields() {
    const { type } = this.args.model;
    if (!type) return null;
    const fields = keyParamsByType(type);
    return fields.map((fieldName) => {
      return this.args.model.allFields.find((attr) => attr.name === fieldName);
    });
  }

  @action cancel() {
    // Generate root form will always have a new model
    this.args.model.unloadRecord();
    this.args.onCancel();
  }

  get groups() {
    return {
      'Key parameters': this.keyParamFields,
      'Subject Alternative Name (SAN) Options': [
        'excludeCnFromSans',
        'serialNumber',
        'altNames',
        'ipSans',
        'uriSans',
        'otherSans',
      ],
      'Additional subject fields': [
        'ou',
        'organization',
        'country',
        'locality',
        'province',
        'streetAddress',
        'postalCode',
      ],
    };
  }

  @action
  checkFormValidity() {
    if (this.args.model.validate) {
      const { isValid, state, invalidFormMessage } = this.args.model.validate();
      this.modelValidations = state;
      this.invalidFormAlert = invalidFormMessage;
      return isValid;
    }
    return true;
  }

  @task
  @waitFor
  *save(event) {
    event.preventDefault();
    const continueSave = this.checkFormValidity();
    if (!continueSave) return;
    try {
      yield this.setUrls();
      const result = yield this.args.model.save({ adapterOptions: this.args.adapterOptions });
      this.flashMessages.success('Successfully generated root.');
      this.router.transitionTo(
        'vault.cluster.secrets.backend.pki.issuers.issuer.details',
        result.backend,
        result.issuerId
      );
    } catch (e) {
      this.errorBanner = errorMessage(e);
      this.invalidFormAlert = 'There was a problem generating the root.';
    }
  }

  async setUrls() {
    if (!this.args.urls || !this.args.urls.canSet || !this.args.urls.hasDirtyAttributes) return;
    return this.args.urls.save();
  }
}
