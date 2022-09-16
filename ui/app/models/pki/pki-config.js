import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default Model.extend({
  backend: attr('string'),
  der: attr(),
  pem: attr('string'),
  caChain: attr('string'),
  attrList(keys) {
    return expandAttributeMeta(this, keys);
  },

  //urls
  urlsAttrs: computed(function () {
    let keys = ['issuingCertificates', 'crlDistributionPoints', 'ocspServers'];
    return this.attrList(keys);
  }),
  issuingCertificates: attr({
    editType: 'stringArray',
  }),
  crlDistributionPoints: attr({
    label: 'CRL Distribution Points',
    editType: 'stringArray',
  }),
  ocspServers: attr({
    label: 'OCSP Servers',
    editType: 'stringArray',
  }),

  //tidy
  tidyAttrs: computed(function () {
    let keys = ['tidyCertStore', 'tidyRevocationList', 'safetyBuffer'];
    return this.attrList(keys);
  }),
  tidyCertStore: attr('boolean', {
    defaultValue: false,
    label: 'Tidy the Certificate Store',
  }),
  tidyRevocationList: attr('boolean', {
    defaultValue: false,
    label: 'Tidy the Revocation List (CRL)',
  }),
  safetyBuffer: attr({
    defaultValue: '72h',
    editType: 'ttl',
  }),

  crlAttrs: computed(function () {
    let keys = ['expiry', 'disable'];
    return this.attrList(keys);
  }),
  //crl
  expiry: attr({
    label: 'CRL building enabled',
    helperTextEnabled: 'The CRL will expire after',
    defaultValue: '72h',
    editType: 'ttl',
    hideToggle: true, // this form field is wrapped by a toggle that sets the 'disable' attr
  }),
  disable: attr('boolean', { defaultValue: false }),
});
