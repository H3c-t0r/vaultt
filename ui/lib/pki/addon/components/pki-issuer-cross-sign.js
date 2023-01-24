import Component from '@glimmer/component';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import errorMessage from 'vault/utils/error-message';
import { waitFor } from '@ember/test-waiters';
/**
 * @module PkiIssuerCrossSign
 * PkiIssuerCrossSign components render from a parent issuer's details page to cross-sign an intermediate issuer.
 * The component reads an existing intermediate issuer, cross-signs it with a parent issuer and imports the new
 * issuer into an existing intermediate mount using three inputs from the user:
 * intermediateMount (the mount where the issuer to be cross signed lives)
 * intermediateIssuer (the name of the intermediate issuer, located in the above mount)
 * newCrossSignedIssuer (the name of the to-be-cross-signed intermediate issuer)
 *
 * The requests involved and how those inputs are used:
 * 1. Read an existing intermediate issuer
 *    -> GET /:intermediateMount/issuer/:intermediateIssuer
 * 2. Create a new CSR based on this existing issuer ID
 *    -> POST /:intermediateMount/intermediate/generate/existing
 * 3. Sign it with the new parent issuer, minting a new certificate.
 *    -> POST /this.args.parentIssuer.backend/issuer/this.args.parentIssuer.issuerName/sign-intermediate
 * 4. Import it back into the existing mount
 *    -> POST /:intermediateMount/issuers/import/bundle
 * 5. Read the imported issuer
 *    -> GET /:intermediateMount/issuer/:issuer_id
 * 5. Update this issuer with the newCrossSignedIssuer
 *    -> POST /:intermediateMount/issuer/:issuer_id
 *
 * @example
 * ```js
 * <PkiIssuerCrossSign @parentIssuer={{this.model}} />
 * ```
 * @param {object} parentIssuer - the model of the issuing certificate that will sign the issuer to-be cross-signed
 */

export default class PkiIssuerCrossSign extends Component {
  @service store;
  @tracked formData = [];
  @tracked signedIssuers = [];

  inputFields = [
    {
      label: 'Mount path',
      key: 'intermediateMount',
      placeholder: 'Mount path',
      helpText: 'The mount in which your new certificate can be found.',
    },
    {
      label: "Issuer's current name",
      key: 'intermediateIssuer',
      placeholder: 'Current issuer name',
      helpText: 'The API name of the previous intermediate which was cross-signed.',
    },
    {
      label: 'New issuer name',
      key: 'newCrossSignedIssuer',
      placeholder: 'Enter a new issuer name',
      helpText: `This is your new issuer’s name in the API.`,
    },
  ];

  @task
  @waitFor
  *submit(e) {
    e.preventDefault();
    this.signedIssuers = [];

    // iterate through submitted data and cross-sign each certificate
    for (let row = 0; row < this.formData.length; row++) {
      const { intermediateMount, intermediateIssuer, newCrossSignedIssuer } = this.formData[row];
      try {
        // returns data from existing and newly cross-signed issuers
        // { intermediateIssuer: existingIssuer, newCrossSignedIssuer: crossSignedIssuer, intermediateMount: intMount }
        const data = yield this.crossSignIntermediate(
          intermediateMount,
          intermediateIssuer,
          newCrossSignedIssuer
        );
        this.signedIssuers.addObject({ ...data, hasError: false });
      } catch (error) {
        this.signedIssuers.addObject({ ...this.formData[row], hasError: errorMessage(error) });
        continue;
      }
    }
  }

  @action
  async crossSignIntermediate(intMount, intName, newCrossSignedIssuer) {
    // 1. Fetch issuer we want to sign
    const existingIssuer = await this.store.queryRecord('pki/issuer', {
      backend: intMount,
      id: intName,
    });

    // 2. Create the new CSR
    const newCsr = await this.store
      .createRecord('pki/action', {
        keyRef: existingIssuer.keyId,
        commonName: existingIssuer.commonName,
        type: 'existing',
      })
      .save({
        adapterOptions: { actionType: 'generate-csr', mount: intMount, useIssuer: false },
      })
      .then(({ csr }) => csr);

    // 3. Sign newCSR with correct parent to create cross-signed cert
    const signedCaChain = await this.store
      .createRecord('pki/action', {
        csr: newCsr,
        commonName: existingIssuer.commonName,
      })
      .save({
        adapterOptions: {
          actionType: 'sign-intermediate',
          mount: this.args.parentIssuer.backend,
          issuerName: this.args.parentIssuer.issuerName,
        },
      })
      .then(({ caChain }) => caChain.join('\n'));

    // 4. Import the newly cross-signed cert to become an issuer
    const issuerId = await this.store
      .createRecord('pki/issuer', { pemBundle: signedCaChain })
      .save({ adapterOptions: { import: true, mount: intMount } })
      .then((importedIssuer) => {
        return Object.keys(importedIssuer.mapping).find(
          // matching key is the issuer_id
          (key) => importedIssuer.mapping[key] === existingIssuer.keyId
        );
      });

    // 5. Fetch issuer imported above by issuer_id, name and save
    const crossSignedIssuer = await this.store.queryRecord('pki/issuer', { backend: intMount, id: issuerId });
    crossSignedIssuer.issuerName = newCrossSignedIssuer;
    crossSignedIssuer.save({ adapterOptions: { mount: intMount } });
    return {
      intermediateIssuer: existingIssuer,
      newCrossSignedIssuer: crossSignedIssuer,
      intermediateMount: intMount,
    };
  }
}
