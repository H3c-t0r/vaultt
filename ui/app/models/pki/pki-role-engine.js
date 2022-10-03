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
    label: 'Issuer reference',
    defaultValue: 'default',
    subText:
      'Specifies the issuer that will be used to create certificates with this role.  To find this, run [command]. By default, we will use the mounts default issuer.',
  })
  issuerRef;

  @attr({
    label: 'Not valid after',
    subText:
      'The time after which this certificate will no longer be valid. This can be a TTL (a range of time from now) or a specific date. If no TTL is set, the system uses "default" or the value of max_ttl, whichever is shorter. Alternatively, you can set the not_after date below.',
    editType: 'yield',
  })
  customTtl;

  @attr({
    label: 'Backdate validity',
    helperTextEnabled:
      'Also called the not_before_duration property. Allows certificates to be valid for a certain time period before now. This is useful to correct clock misalignment on various systems when setting up your CA.',
    editType: 'ttl',
    hideToggle: true,
    defaultValue: '30s', // type in API is duration which accepts both an integer and string e.g. 30 || '30s'
  })
  notBeforeDuration;

  @attr({
    label: 'Max TTL',
    helperTextDisabled:
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
  generateLease;

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
  addBasicConstraints;

  // Overriding Domain Handling options
  @attr({
    label: 'Allowed Domains',
    subText: 'Specifies the domains this role is allowed to issue certificates for. Add one item per row.',
    editType: 'stringArray',
  })
  allowedDomains;
  @attr('boolean', {
    label: 'Allow templates in allowed domains',
  })
  allowedDomainsTemplate;
  @attr('string', {
    editType: 'moreInfo',
    text: 'These options can interact intricately with one another. For more information,',
    docText: 'learn more here.',
    docLink: '/docs/concepts/password-policies',
  })
  moreInfo;

  // Overriding Key parameters options
  @attr('string', {
    label: 'Key type',
    possibleValues: ['rsa', 'ec', 'ed25519', 'any'],
    defaultValue: 'rsa',
  })
  keyType;

  @attr('number', {
    label: 'Key bits',
  })
  keyBits;
  // To set conditional "possibleValues" for the field "keyBits" that depend on the value of the selected "keyType"
  get keyBitsConditional() {
    // ARG TODO confirm this is the correct matrix (probably with core)
    const keyBitOptions = {
      rsa: [2048, 3072, 4096],
      ec: [256, 224, 384, 521],
      ed25519: [0],
      any: [0],
    };
    const attrs = expandAttributeMeta(this, ['keyBits']);
    attrs[0].options['possibleValues'] = keyBitOptions[this.keyType];
    return attrs[0];
  }

  @attr('number', {
    label: 'Signature bits',
    possibleValues: [
      {
        value: 0,
        displayName: '0 to automatically detect based on key length',
      },
      {
        value: 256,
        displayName: '256 for SHA-2-256',
      },
      {
        value: 384,
        displayName: '384 for SHA-2-384',
      },
      {
        value: 512,
        displayName: '512 for SHA-2-5124',
      },
    ],
  })
  signatureBits;

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
            'addBasicConstraints',
          ],
        },
        {
          'Domain handling': [
            'allowedDomains',
            'allowedDomainsTemplate',
            'allowBareDomains',
            'allowSubdomains',
            'allowGlobDomains',
            'allowWildcardCertificates',
            'allowLocalhost', // default: true
            'allowAnyName',
            'enforceHostnames', // default: true
            'moreInfo', // shows as helperText with icon at bottom of the options box.
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
