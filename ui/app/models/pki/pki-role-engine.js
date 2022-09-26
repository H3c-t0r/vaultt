import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import { withModelValidations } from 'vault/decorators/model-validations';

import fieldToAttrs from 'vault/utils/field-to-attrs';

const validations = {
  name: [{ type: 'presence', message: 'Name is required.' }],
};

@withModelValidations(validations)
export default class PkiRoleEngineModel extends Model {
  @attr('string', { readOnly: true }) backend;

  @attr('string', {
    label: 'Role name',
    fieldValue: 'name',
  })
  name;

  @attr('string', {
    label: 'Issuer Reference',
    subText:
      'Specifies the default issuer of this request. May be the value default, a name, or an issuer ID. To find this, you can view the issuer or run: read -field=default <mount-name>/config/issuers in the CLI.',
  })
  issuerRef;

  @attr({
    label: 'Not valid after',
    subText:
      'The time after which this certificate will no longer be valid. This can be a TTL (a range of time from now) or a specific date. If not set, the system uses "default" or the value of max_ttl, whichever is shorter. Alternatively, you can set the not_after date below.',
    editType: 'yield',
  })
  customTtl;

  @attr({
    label: 'Backdate validity',
    subText:
      'Also called the notBeforeDuration property. Allows certificates to be valid for a certain time period before now. This is useful to correct clock misalignment on various systems when setting up your CA.',
    editType: 'ttl',
    defaultValue: '30s',
    hideToggle: true,
  })
  notBeforeDuration;

  @attr({
    label: 'Max TTL',
    subText:
      'The maximum Time-To-Live of certificates generated by this role. If not set, the system max lease TTL will be used.',
    editType: 'ttl',
  })
  maxTtl;

  @attr('boolean', {
    label: 'Generate lease with certificate',
    subText:
      'Specifies if certificates issued/signed against this role will have Vault leases attached to them.',
    editType: 'boolean',
    docLink: '/api-docs/secret/pki#create-update-role',
  })
  generateLease; // ARG TODO confirm false by default

  @attr('boolean', {
    label: 'Do not store certificates in storage backend',
    subText:
      'This can improve performance when issuing large numbers of certificates. However, certificates issued in this way cannot be enumerated or revoked.',
    editType: 'boolean',
    docLink: '/api-docs/secret/pki#create-update-role',
  })
  noStore;

  @attr('boolean', {
    label: 'Basic constraints valid for non CA.',
    subText: 'Mark Basic Constraints valid when issuing non-CA certificates.',
    editType: 'boolean',
  })
  noStore;

  // must be a getter so it can be added to the prototype needed in the pathHelp service on the line here: if (newModel.merged || modelProto.useOpenAPI !== true) {
  get useOpenAPI() {
    return true;
  }
  getHelpUrl(backend) {
    return `/v1/${backend}/roles/example?help=1`;
  }
  @lazyCapabilities(apiPath`${'backend'}/roles/${'id'}`, 'backend', 'id') updatePath;
  get canDelete() {
    return this.updatePath.get('canCreate');
  }
  get canEdit() {
    return this.updatePath.get('canEdit');
  }
  get canRead() {
    return this.updatePath.get('canRead');
  }

  @lazyCapabilities(apiPath`${'backend'}/issue/${'id'}`, 'backend', 'id') generatePath;
  get canReadIssue() {
    // ARG TODO was duplicate name, added Issue
    return this.generatePath.get('canUpdate');
  }
  @lazyCapabilities(apiPath`${'backend'}/sign/${'id'}`, 'backend', 'id') signPath;
  get canSign() {
    return this.signPath.get('canUpdate');
  }
  @lazyCapabilities(apiPath`${'backend'}/sign-verbatim/${'id'}`, 'backend', 'id') signVerbatimPath;
  get canSignVerbatim() {
    return this.signVerbatimPath.get('canUpdate');
  }

  // Form Fields not hidden in toggle options
  _attributeMeta = null;
  get formFields() {
    if (!this._attributeMeta) {
      this._attributeMeta = expandAttributeMeta(this, ['name', 'clientType', 'redirectUris']);
    }
    return this._attributeMeta;
  }

  // Form fields hidden behind toggle options
  _fieldToAttrsGroups = null;
  // ARG TODO: I removed 'allowedDomains' but I'm fairly certain it needs to be somewhere. Confirm with design.
  get fieldGroups() {
    if (!this._fieldToAttrsGroups) {
      this._fieldToAttrsGroups = fieldToAttrs(this, [
        {
          default: [
            'name',
            'issuerRef',
            'customTtl',
            'notBeforeDuration',
            'maxTtl',
            'generateLease',
            'noStore',
          ],
        },
        {
          'Domain handling': [
            'allowedDomains',
            'allowedDomainTemplate',
            'allowBareDomains',
            'allowSubdomains',
            'allowGlobDomains',
            'allowWildcardCertificates',
            'allowLocalhost',
            'allowAnyName',
            'enforceHostnames',
          ],
        },
        {
          'Key parameters': ['keyType', 'keyBits', 'signatureBits'],
        },
        {
          'Key usage': [
            'DigitalSignature', // ARG TODO: capitalized in the docs, but should confirm
            'KeyAgreement',
            'KeyEncipherment',
            'extKeyUsage', // ARG TODO: takes a list, but we have these as checkboxes from the options on the golang site: https://pkg.go.dev/crypto/x509#ExtKeyUsage
          ],
        },
        { 'Policy identifiers': ['policyIdentifiers'] },
        {
          'Subject Alternative Name (SAN) Options': ['allowIpSans', 'allowedUriSans', 'allowedOtherSans'],
        },
        {
          'Additional subject fields': [
            'allowed_serial_numbers',
            'requireCn',
            'useCsrCommonName',
            'useCsrSans',
            'ou',
            'organization',
            'country',
            'locality',
            'province',
            'streetAddress',
            'postalCode',
          ],
        },
      ]);
    }
    return this._fieldToAttrsGroups;
  }
}
