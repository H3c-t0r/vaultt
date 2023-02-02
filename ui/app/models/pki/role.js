import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { withModelValidations } from 'vault/decorators/model-validations';
import { withFormFields } from 'vault/decorators/model-form-fields';

const validations = {
  name: [{ type: 'presence', message: 'Name is required.' }],
};

const fieldGroups = [
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
      'allowLocalhost', // default: true (returned true by OpenApi)
      'allowAnyName',
      'enforceHostnames', // default: true (returned true by OpenApi)
    ],
  },
  {
    'Key parameters': ['keyType', 'keyBits', 'signatureBits'],
  },
  {
    'Key usage': ['keyUsage', 'extKeyUsage', 'extKeyUsageOids'],
  },
  { 'Policy identifiers': ['policyIdentifiers'] },
  {
    'Subject Alternative Name (SAN) Options': [
      'allowIpSans',
      'allowedUriSans',
      'allowUriSansTemplate',
      'allowedOtherSans',
    ],
  },
  {
    'Additional subject fields': [
      'allowedSerialNumbers',
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
];

@withFormFields(null, fieldGroups)
@withModelValidations(validations)
export default class PkiRoleModel extends Model {
  get useOpenAPI() {
    // must be a getter so it can be accessed in path-help.js
    return true;
  }
  getHelpUrl(backend) {
    return `/v1/${backend}/roles/example?help=1`;
  }

  @attr('string', { readOnly: true }) backend;

  /* Overriding OpenApi default options */
  @attr('string', {
    label: 'Role name',
    fieldValue: 'name',
    editDisabled: true,
  })
  name;

  @attr('string', {
    label: 'Issuer reference',
    detailsLabel: 'Issuer',
    defaultValue: 'default',
    subText: `Specifies the issuer that will be used to create certificates with this role. To find this, run read -field=default pki_int/config/issuers in the console. By default, we will use the mounts default issuer.`,
  })
  issuerRef;

  @attr({
    label: 'Not valid after',
    detailsLabel: 'Issued certificates expire after',
    subText:
      'The time after which this certificate will no longer be valid. This can be a TTL (a range of time from now) or a specific date.',
    editType: 'yield',
  })
  customTtl;

  @attr({
    label: 'Backdate validity',
    detailsLabel: 'Issued certificate backdating',
    helperTextDisabled: 'Vault will use the default value, 30s',
    helperTextEnabled:
      'Also called the not_before_duration property. Allows certificates to be valid for a certain time period before now. This is useful to correct clock misalignment on various systems when setting up your CA.',
    editType: 'ttl',
    defaultValue: '30s',
  })
  notBeforeDuration;

  @attr({
    label: 'Max TTL',
    helperTextDisabled:
      'The maximum Time-To-Live of certificates generated by this role. If not set, the system max lease TTL will be used.',
    editType: 'ttl',
    defaultShown: 'System default',
  })
  maxTtl;

  @attr('boolean', {
    label: 'Generate lease with certificate',
    subText:
      'Specifies if certificates issued/signed against this role will have Vault leases attached to them.',
    editType: 'boolean',
    docLink: '/vault/api-docs/secret/pki#create-update-role',
  })
  generateLease;

  @attr('boolean', {
    label: 'Do not store certificates in storage backend',
    detailsLabel: 'Store in storage backend', // template reverses value
    subText:
      'This can improve performance when issuing large numbers of certificates. However, certificates issued in this way cannot be enumerated or revoked.',
    editType: 'boolean',
    docLink: '/vault/api-docs/secret/pki#create-update-role',
  })
  noStore;

  @attr('boolean', {
    label: 'Basic constraints valid for non-CA',
    detailsLabel: 'Add basic constraints',
    subText: 'Mark Basic Constraints valid when issuing non-CA certificates.',
    editType: 'boolean',
  })
  addBasicConstraints;
  /* End of overriding default options */

  /* Overriding OpenApi Domain handling options */
  @attr({
    label: 'Allowed domains',
    subText: 'Specifies the domains this role is allowed to issue certificates for. Add one item per row.',
    editType: 'stringArray',
  })
  allowedDomains;

  @attr('boolean', {
    label: 'Allow templates in allowed domains',
  })
  allowedDomainsTemplate;
  /* End of overriding Domain handling options */

  /* Overriding OpenApi Key parameters options */
  @attr('string', {
    label: 'Key type',
    possibleValues: ['rsa', 'ec', 'ed25519', 'any'],
    defaultValue: 'rsa',
  })
  keyType;

  @attr('string', {
    label: 'Key bits',
    defaultValue: '2048',
  })
  keyBits; // no possibleValues because options are dependent on selected key type

  @attr('number', {
    label: 'Signature bits',
    subText: `Only applicable for key_type 'RSA'. Ignore for other key types.`,
    defaultValue: '0',
    possibleValues: ['0', '256', '384', '512'],
  })
  signatureBits;
  /* End of overriding Key parameters options */

  /* Overriding API Policy identifier option */
  @attr({
    label: 'Policy identifiers',
    subText: 'A comma-separated string or list of policy object identifiers (OIDs). Add one per row. ',
    editType: 'stringArray',
  })
  policyIdentifiers;
  /* End of overriding Policy identifier options */

  /* Overriding OpenApi SAN options */
  @attr('boolean', {
    label: 'Allow IP SANs',
    subText: 'Specifies if clients can request IP Subject Alternative Names.',
    editType: 'boolean',
    defaultValue: true,
  })
  allowIpSans;

  @attr({
    label: 'URI Subject Alternative Names (URI SANs)',
    subText: 'Defines allowed URI Subject Alternative Names. Add one item per row',
    editType: 'stringArray',
    docLink: '/vault/docs/concepts/policies',
  })
  allowedUriSans;

  @attr('boolean', {
    label: 'Allow URI SANs template',
    subText: 'If true, the URI SANs above may contain templates, as with ACL Path Templating.',
    editType: 'boolean',
    docLink: '/vault/docs/concepts/policies',
  })
  allowUriSansTemplate;

  @attr({
    label: 'Other SANs',
    subText: 'Defines allowed custom OID/UTF8-string SANs. Add one item per row.',
    editType: 'stringArray',
  })
  allowedOtherSans;
  /* End of overriding SAN options */

  /* Overriding OpenApi Additional subject field options */
  @attr({
    label: 'Allowed serial numbers',
    subText:
      'A list of allowed serial numbers to be requested during certificate issuance. Shell-style globbing is supported. If empty, custom-specified serial numbers will be forbidden.',
    editType: 'stringArray',
  })
  allowedSerialNumbers;

  @attr('boolean', {
    label: 'Require common name',
    subText: 'If set to false, common name will be optional when generating a certificate.',
    defaultValue: true,
  })
  requireCn;

  @attr('boolean', {
    label: 'Use CSR common name',
    subText:
      'When used with the CSR signing endpoint, the common name in the CSR will be used instead of taken from the JSON data.',
    defaultValue: true,
  })
  useCsrCommonName;

  @attr('boolean', {
    label: 'Use CSR SANs',
    subText:
      'When used with the CSR signing endpoint, the subject alternate names in the CSR will be used instead of taken from the JSON data.',
    defaultValue: true,
  })
  useCsrSans;

  @attr({
    label: 'Organization Units (OU)',
    subText:
      'A list of allowed serial numbers to be requested during certificate issuance. Shell-style globbing is supported. If empty, custom-specified serial numbers will be forbidden.',
  })
  ou;

  @attr('array', {
    defaultValue() {
      return ['DigitalSignature', 'KeyAgreement', 'KeyEncipherment'];
    },
    defaultShown: 'None',
  })
  keyUsage;

  @attr('array', {
    defaultShown: 'None',
  })
  extKeyUsage;

  @attr('array', {
    defaultShown: 'None',
  })
  extKeyUsageOids;

  @attr('string') organization;
  @attr('string') country;
  @attr('string') locality;
  @attr('string') province;
  @attr('string') streetAddress;
  @attr('string') postalCode;
  /* End of overriding Additional subject field options */

  /* CAPABILITIES
   * Default to show UI elements unless we know they can't access the given path
   */
  @lazyCapabilities(apiPath`${'backend'}/roles/${'id'}`, 'backend', 'id') updatePath;
  get canDelete() {
    return this.updatePath.get('isLoading') || this.updatePath.get('canCreate') !== false;
  }
  get canEdit() {
    return this.updatePath.get('isLoading') || this.updatePath.get('canUpdate') !== false;
  }
  get canRead() {
    return this.updatePath.get('isLoading') || this.updatePath.get('canRead') !== false;
  }

  @lazyCapabilities(apiPath`${'backend'}/issue/${'id'}`, 'backend', 'id') generatePath;
  get canGenerateCert() {
    return this.generatePath.get('isLoading') || this.generatePath.get('canUpdate') !== false;
  }
  @lazyCapabilities(apiPath`${'backend'}/sign/${'id'}`, 'backend', 'id') signPath;
  get canSign() {
    return this.signPath.get('isLoading') || this.signPath.get('canUpdate') !== false;
  }
  @lazyCapabilities(apiPath`${'backend'}/sign-verbatim/${'id'}`, 'backend', 'id') signVerbatimPath;
  get canSignVerbatim() {
    return this.signVerbatimPath.get('isLoading') || this.signVerbatimPath.get('canUpdate') !== false;
  }

  // Gets header/footer copy for specific toggle groups.
  get fieldGroupsInfo() {
    return {
      'Domain handling': {
        footer: {
          text: 'These options can interact intricately with one another. For more information,',
          docText: 'learn more here.',
          docLink: '/vault/api-docs/secret/pki#allowed_domains',
        },
      },
      'Key parameters': {
        header: {
          text: `These are the parameters for generating or validating the certificate's key material.`,
        },
      },
      'Subject Alternative Name (SAN) Options': {
        header: {
          text: `Subject Alternative Names (SANs) are identities (domains, IP addresses, and URIs) Vault attaches to the requested certificates.`,
        },
      },
      'Additional subject fields': {
        header: {
          text: `Additional identity metadata Vault can attach to the requested certificates.`,
        },
      },
    };
  }
}
