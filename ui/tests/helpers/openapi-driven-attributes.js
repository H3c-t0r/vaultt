// The constants within this file represent the expected model attributes as parsed from OpenAPI
// if changes are made to the OpenAPI spec, that may result in changes that must be reflected
// here as well as ensured to not cause breaking changes within the UI.

/* Secret Engines */
const sshRole = {
  role: {
    editType: 'string',
    helpText: '[Required for all types] DIFFERNT Name of the role being created.',
    fieldValue: 'mutableId',
    fieldGroup: 'default',
    readOnly: true,
    label: 'Role',
    type: 'string',
  },
  algorithmSigner: {
    editType: 'string',
    helpText:
      '[Not applicable for OTP type] [Optional for CA type] When supplied, this value specifies a signing algorithm for the key. Possible values: ssh-rsa, rsa-sha2-256, rsa-sha2-512, default, or the empty string.',
    possibleValues: ['', 'default', 'ssh-rsa', 'rsa-sha2-256', 'rsa-sha2-512'],
    fieldGroup: 'default',
    label: 'Signing Algorithm',
    type: 'string',
  },
  allowBareDomains: {
    editType: 'boolean',
    helpText:
      '[Not applicable for OTP type] [Optional for CA type] If set, host certificates that are requested are allowed to use the base domains listed in "allowed_domains", e.g. "example.com". This is a separate option as in some cases this can be considered a security threat.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  allowHostCertificates: {
    editType: 'boolean',
    helpText:
      "[Not applicable for OTP type] [Optional for CA type] If set, certificates are allowed to be signed for use as a 'host'.",
    fieldGroup: 'default',
    type: 'boolean',
  },
  allowSubdomains: {
    editType: 'boolean',
    helpText:
      '[Not applicable for OTP type] [Optional for CA type] If set, host certificates that are requested are allowed to use subdomains of those listed in "allowed_domains".',
    fieldGroup: 'default',
    type: 'boolean',
  },
  allowUserCertificates: {
    editType: 'boolean',
    helpText:
      "[Not applicable for OTP type] [Optional for CA type] If set, certificates are allowed to be signed for use as a 'user'.",
    fieldGroup: 'default',
    type: 'boolean',
  },
  allowUserKeyIds: {
    editType: 'boolean',
    helpText:
      '[Not applicable for OTP type] [Optional for CA type] If true, users can override the key ID for a signed certificate with the "key_id" field. When false, the key ID will always be the token display name. The key ID is logged by the SSH server and can be useful for auditing.',
    fieldGroup: 'default',
    label: 'Allow User Key IDs',
    type: 'boolean',
  },
  allowedCriticalOptions: {
    editType: 'string',
    helpText:
      '[Not applicable for OTP type] [Optional for CA type] A comma-separated list of critical options that certificates can have when signed. To allow any critical options, set this to an empty string.',
    fieldGroup: 'default',
    type: 'string',
  },
  allowedDomains: {
    editType: 'string',
    helpText:
      '[Not applicable for OTP type] [Optional for CA type] If this option is not specified, client can request for a signed certificate for any valid host. If only certain domains are allowed, then this list enforces it.',
    fieldGroup: 'default',
    type: 'string',
  },
  allowedDomainsTemplate: {
    editType: 'boolean',
    helpText:
      '[Not applicable for OTP type] [Optional for CA type] If set, Allowed domains can be specified using identity template policies. Non-templated domains are also permitted.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  allowedExtensions: {
    editType: 'string',
    helpText:
      "[Not applicable for OTP type] [Optional for CA type] A comma-separated list of extensions that certificates can have when signed. An empty list means that no extension overrides are allowed by an end-user; explicitly specify '*' to allow any extensions to be set.",
    fieldGroup: 'default',
    type: 'string',
  },
  allowedUserKeyLengths: {
    editType: 'object',
    helpText:
      '[Not applicable for OTP type] [Optional for CA type] If set, allows the enforcement of key types and minimum key sizes to be signed.',
    fieldGroup: 'default',
    type: 'object',
  },
  allowedUsers: {
    editType: 'string',
    helpText:
      "[Optional for all types] [Works differently for CA type] If this option is not specified, or is '*', client can request a credential for any valid user at the remote host, including the admin user. If only certain usernames are to be allowed, then this list enforces it. If this field is set, then credentials can only be created for default_user and usernames present in this list. Setting this option will enable all the users with access to this role to fetch credentials for all other usernames in this list. Use with caution. N.B.: with the CA type, an empty list means that no users are allowed; explicitly specify '*' to allow any user.",
    fieldGroup: 'default',
    type: 'string',
  },
  allowedUsersTemplate: {
    editType: 'boolean',
    helpText:
      '[Not applicable for OTP type] [Optional for CA type] If set, Allowed users can be specified using identity template policies. Non-templated users are also permitted.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  cidrList: {
    editType: 'string',
    helpText:
      '[Optional for OTP type] [Not applicable for CA type] Comma separated list of CIDR blocks for which the role is applicable for. CIDR blocks can belong to more than one role.',
    fieldGroup: 'default',
    label: 'CIDR List',
    type: 'string',
  },
  defaultCriticalOptions: {
    editType: 'object',
    helpText:
      '[Not applicable for OTP type] [Optional for CA type] Critical options certificates should have if none are provided when signing. This field takes in key value pairs in JSON format. Note that these are not restricted by "allowed_critical_options". Defaults to none.',
    fieldGroup: 'default',
    type: 'object',
  },
  defaultExtensions: {
    editType: 'object',
    helpText:
      '[Not applicable for OTP type] [Optional for CA type] Extensions certificates should have if none are provided when signing. This field takes in key value pairs in JSON format. Note that these are not restricted by "allowed_extensions". Defaults to none.',
    fieldGroup: 'default',
    type: 'object',
  },
  defaultExtensionsTemplate: {
    editType: 'boolean',
    helpText:
      '[Not applicable for OTP type] [Optional for CA type] If set, Default extension values can be specified using identity template policies. Non-templated extension values are also permitted.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  defaultUser: {
    editType: 'string',
    helpText:
      "[Required for OTP type] [Optional for CA type] Default username for which a credential will be generated. When the endpoint 'creds/' is used without a username, this value will be used as default username.",
    fieldGroup: 'default',
    label: 'Default Username',
    type: 'string',
  },
  defaultUserTemplate: {
    editType: 'boolean',
    helpText:
      '[Not applicable for OTP type] [Optional for CA type] If set, Default user can be specified using identity template policies. Non-templated users are also permitted.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  excludeCidrList: {
    editType: 'string',
    helpText:
      '[Optional for OTP type] [Not applicable for CA type] Comma separated list of CIDR blocks. IP addresses belonging to these blocks are not accepted by the role. This is particularly useful when big CIDR blocks are being used by the role and certain parts of it needs to be kept out.',
    fieldGroup: 'default',
    label: 'Exclude CIDR List',
    type: 'string',
  },
  keyIdFormat: {
    editType: 'string',
    helpText:
      "[Not applicable for OTP type] [Optional for CA type] When supplied, this value specifies a custom format for the key id of a signed certificate. The following variables are available for use: '{{token_display_name}}' - The display name of the token used to make the request. '{{role_name}}' - The name of the role signing the request. '{{public_key_hash}}' - A SHA256 checksum of the public key that is being signed.",
    fieldGroup: 'default',
    label: 'Key ID Format',
    type: 'string',
  },
  keyType: {
    editType: 'string',
    helpText:
      "[Required for all types] Type of key used to login to hosts. It can be either 'otp' or 'ca'. 'otp' type requires agent to be installed in remote hosts.",
    possibleValues: ['otp', 'ca'],
    fieldGroup: 'default',
    defaultValue: 'ca',
    type: 'string',
  },
  maxTtl: {
    editType: 'ttl',
    helpText: '[Not applicable for OTP type] [Optional for CA type] The maximum allowed lease duration',
    fieldGroup: 'default',
    label: 'Max TTL',
  },
  notBeforeDuration: {
    editType: 'ttl',
    helpText:
      '[Not applicable for OTP type] [Optional for CA type] The duration that the SSH certificate should be backdated by at issuance.',
    fieldGroup: 'default',
    defaultValue: 30,
    label: 'Not before duration',
  },
  port: {
    editType: 'number',
    helpText:
      "[Optional for OTP type] [Not applicable for CA type] Port number for SSH connection. Default is '22'. Port number does not play any role in creation of OTP. For 'otp' type, this is just a way to inform client about the port number to use. Port number will be returned to client by Vault server along with OTP.",
    fieldGroup: 'default',
    defaultValue: 22,
    type: 'number',
  },
  ttl: {
    editType: 'ttl',
    helpText:
      '[Not applicable for OTP type] [Optional for CA type] The lease duration if no specific lease duration is requested. The lease duration controls the expiration of certificates issued by this backend. Defaults to the value of max_ttl.',
    fieldGroup: 'default',
    label: 'TTL',
  },
};

const kmipConfig = {
  defaultTlsClientKeyBits: {
    editType: 'number',
    helpText: ' on key type',
    fieldGroup: 'default',
    defaultValue: 256,
    label: 'Default TLS Client Key bits',
    type: 'number',
  },
  defaultTlsClientKeyType: {
    editType: 'string',
    helpText: 'Client certificate key type, rsa or ec',
    possibleValues: ['rsa', 'ec'],
    fieldGroup: 'default',
    defaultValue: 'ec',
    label: 'Default TLS Client Key type',
    type: 'string',
  },
  defaultTlsClientTtl: {
    editType: 'ttl',
    helpText:
      'Client certificate TTL in either an integer number of seconds (3600) or an integer time unit (1h)',
    fieldGroup: 'default',
    defaultValue: '336h',
    label: 'Default TLS Client TTL',
  },
  listenAddrs: {
    editType: 'stringArray',
    helpText:
      'A list of address:port to listen on. A bare address without port may be provided, in which case port 5696 is assumed.',
    fieldGroup: 'default',
    defaultValue: '127.0.0.1:5696',
  },
  serverHostnames: {
    editType: 'stringArray',
    helpText:
      "A list of hostnames to include in the server's TLS certificate as SAN DNS names. The first will be used as the common name (CN).",
    fieldGroup: 'default',
  },
  serverIps: {
    editType: 'stringArray',
    helpText: "A list of IP to include in the server's TLS certificate as SAN IP addresses.",
    fieldGroup: 'default',
  },
  tlsCaKeyBits: {
    editType: 'number',
    helpText: 'CA key bits, valid values depend on key type',
    fieldGroup: 'default',
    defaultValue: 256,
    label: 'TLS CA Key bits',
    type: 'number',
  },
  tlsCaKeyType: {
    editType: 'string',
    helpText: 'CA key type, rsa or ec',
    possibleValues: ['rsa', 'ec'],
    fieldGroup: 'default',
    defaultValue: 'ec',
    label: 'TLS CA Key type',
    type: 'string',
  },
  tlsMinVersion: {
    editType: 'string',
    helpText: 'Min TLS version',
    fieldGroup: 'default',
    defaultValue: 'tls12',
    label: 'Minimum TLS Version',
    type: 'string',
  },
};
const kmipRole = {
  role: {
    editType: 'string',
    helpText: 'Name of the role.',
    fieldValue: 'mutableId',
    fieldGroup: 'default',
    readOnly: true,
    label: 'Role',
    type: 'string',
  },
  operationActivate: {
    editType: 'boolean',
    helpText: 'Allow the "Activate" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Activate',
    type: 'boolean',
  },
  operationAddAttribute: {
    editType: 'boolean',
    helpText: 'Allow the "Add Attribute" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Add Attribute',
    type: 'boolean',
  },
  operationAll: {
    editType: 'boolean',
    helpText:
      'Allow ALL operations to be performed by this role. This can be overridden if other allowed operations are set to false within the same request.',
    fieldGroup: 'default',
    label: 'All',
    type: 'boolean',
  },
  operationCreate: {
    editType: 'boolean',
    helpText: 'Allow the "Create" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Create',
    type: 'boolean',
  },
  operationCreateKeyPair: {
    editType: 'boolean',
    helpText: 'Allow the "Create Key Pair" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Create Key Pair',
    type: 'boolean',
  },
  operationDecrypt: {
    editType: 'boolean',
    helpText: 'Allow the "Decrypt" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Decrypt',
    type: 'boolean',
  },
  operationDeleteAttribute: {
    editType: 'boolean',
    helpText: 'Allow the "Delete Attribute" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Delete Attribute',
    type: 'boolean',
  },
  operationDestroy: {
    editType: 'boolean',
    helpText: 'Allow the "Destroy" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Destroy',
    type: 'boolean',
  },
  operationDiscoverVersions: {
    editType: 'boolean',
    helpText: 'Allow the "Discover Versions" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Discover Versions',
    type: 'boolean',
  },
  operationEncrypt: {
    editType: 'boolean',
    helpText: 'Allow the "Encrypt" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Encrypt',
    type: 'boolean',
  },
  operationGet: {
    editType: 'boolean',
    helpText: 'Allow the "Get" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Get',
    type: 'boolean',
  },
  operationGetAttributeList: {
    editType: 'boolean',
    helpText: 'Allow the "Get Attribute List" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Get Attribute List',
    type: 'boolean',
  },
  operationGetAttributes: {
    editType: 'boolean',
    helpText: 'Allow the "Get Attributes" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Get Attributes',
    type: 'boolean',
  },
  operationImport: {
    editType: 'boolean',
    helpText: 'Allow the "Import" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Import',
    type: 'boolean',
  },
  operationLocate: {
    editType: 'boolean',
    helpText: 'Allow the "Locate" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Locate',
    type: 'boolean',
  },
  operationMac: {
    editType: 'boolean',
    helpText: 'Allow the "Mac" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Mac',
    type: 'boolean',
  },
  operationMacVerify: {
    editType: 'boolean',
    helpText: 'Allow the "Mac Verify" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Mac Verify',
    type: 'boolean',
  },
  operationModifyAttribute: {
    editType: 'boolean',
    helpText: 'Allow the "Modify Attribute" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Modify Attribute',
    type: 'boolean',
  },
  operationNone: {
    editType: 'boolean',
    helpText:
      'Allow NO operations to be performed by this role. This can be overridden if other allowed operations are set to true within the same request.',
    fieldGroup: 'default',
    label: 'None',
    type: 'boolean',
  },
  operationQuery: {
    editType: 'boolean',
    helpText: 'Allow the "Query" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Query',
    type: 'boolean',
  },
  operationRegister: {
    editType: 'boolean',
    helpText: 'Allow the "Register" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Register',
    type: 'boolean',
  },
  operationRekey: {
    editType: 'boolean',
    helpText: 'Allow the "Rekey" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Rekey',
    type: 'boolean',
  },
  operationRekeyKeyPair: {
    editType: 'boolean',
    helpText: 'Allow the "Rekey Key Pair" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Rekey Key Pair',
    type: 'boolean',
  },
  operationRevoke: {
    editType: 'boolean',
    helpText: 'Allow the "Revoke" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Revoke',
    type: 'boolean',
  },
  operationRngRetrieve: {
    editType: 'boolean',
    helpText: 'Allow the "Rng Retrieve" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Rng Retrieve',
    type: 'boolean',
  },
  operationRngSeed: {
    editType: 'boolean',
    helpText: 'Allow the "Rng Seed" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Rng Seed',
    type: 'boolean',
  },
  operationSign: {
    editType: 'boolean',
    helpText: 'Allow the "Sign" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Sign',
    type: 'boolean',
  },
  operationSignatureVerify: {
    editType: 'boolean',
    helpText: 'Allow the "Signature Verify" operation to be performed by this role',
    fieldGroup: 'default',
    label: 'Signature Verify',
    type: 'boolean',
  },
  tlsClientKeyBits: {
    editType: 'number',
    helpText: 'Client certificate key bits, valid values depend on key type',
    fieldGroup: 'default',
    defaultValue: 521,
    label: 'TLS Client Key bits',
    type: 'number',
  },
  tlsClientKeyType: {
    editType: 'string',
    helpText: 'Client certificate key type, rsa or ec',
    possibleValues: ['rsa', 'ec'],
    fieldGroup: 'default',
    defaultValue: 'ec',
    label: 'TLS Client Key type',
    type: 'string',
  },
  tlsClientTtl: {
    editType: 'ttl',
    helpText:
      'Client certificate TTL in either an integer number of seconds (10) or an integer time unit (10s)',
    fieldGroup: 'default',
    defaultValue: '86400',
    label: 'TLS Client TTL',
  },
};

const pkiAcme = {
  allowRoleExtKeyUsage: {
    editType: 'boolean',
    helpText:
      'whether the ExtKeyUsage field from a role is used, defaults to false meaning that certificate will be signed with ServerAuth.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  allowedIssuers: {
    editType: 'stringArray',
    helpText:
      'which issuers are allowed for use with ACME; by default, this will only be the primary (default) issuer',
    fieldGroup: 'default',
  },
  allowedRoles: {
    editType: 'stringArray',
    helpText:
      "which roles are allowed for use with ACME; by default via '*', these will be all roles including sign-verbatim; when concrete role names are specified, any default_directory_policy role must be included to allow usage of the default acme directories under /pki/acme/directory and /pki/issuer/:issuer_id/acme/directory.",
    fieldGroup: 'default',
  },
  defaultDirectoryPolicy: {
    editType: 'string',
    helpText:
      'the policy to be used for non-role-qualified ACME requests; by default ACME issuance will be otherwise unrestricted, equivalent to the sign-verbatim endpoint; one may also specify a role to use as this policy, as "role:<role_name>", the specified role must be allowed by allowed_roles',
    fieldGroup: 'default',
    type: 'string',
  },
  dnsResolver: {
    editType: 'string',
    helpText:
      'DNS resolver to use for domain resolution on this mount. Defaults to using the default system resolver. Must be in the format <host>:<port>, with both parts mandatory.',
    fieldGroup: 'default',
    type: 'string',
  },
  eabPolicy: {
    editType: 'string',
    helpText:
      "Specify the policy to use for external account binding behaviour, 'not-required', 'new-account-required' or 'always-required'",
    fieldGroup: 'default',
    type: 'string',
  },
  enabled: {
    editType: 'boolean',
    helpText:
      'whether ACME is enabled, defaults to false meaning that clusters will by default not get ACME support',
    fieldGroup: 'default',
    type: 'boolean',
  },
};
const pkiCertGenerate = {
  role: {
    editType: 'string',
    helpText: 'The desired role with configuration for this request',
    fieldValue: 'mutableId',
    fieldGroup: 'default',
    readOnly: true,
    label: 'Role',
    type: 'string',
  },
  altNames: {
    editType: 'string',
    helpText:
      'The requested Subject Alternative Names, if any, in a comma-delimited list. If email protection is enabled for the role, this may contain email addresses.',
    fieldGroup: 'default',
    label: 'DNS/Email Subject Alternative Names (SANs)',
    type: 'string',
  },
  commonName: {
    editType: 'string',
    helpText:
      'The requested common name; if you want more than one, specify the alternative names in the alt_names map. If email protection is enabled in the role, this may be an email address.',
    fieldGroup: 'default',
    type: 'string',
  },
  excludeCnFromSans: {
    editType: 'boolean',
    helpText:
      'If true, the Common Name will not be included in DNS or Email Subject Alternate Names. Defaults to false (CN is included).',
    fieldGroup: 'default',
    label: 'Exclude Common Name from Subject Alternative Names (SANs)',
    type: 'boolean',
  },
  format: {
    editType: 'string',
    helpText:
      'Format for returned data. Can be "pem", "der", or "pem_bundle". If "pem_bundle", any private key and issuing cert will be appended to the certificate pem. If "der", the value will be base64 encoded. Defaults to "pem".',
    possibleValues: ['pem', 'der', 'pem_bundle'],
    fieldGroup: 'default',
    defaultValue: 'pem',
    type: 'string',
  },
  ipSans: {
    editType: 'stringArray',
    helpText: 'The requested IP SANs, if any, in a comma-delimited list',
    fieldGroup: 'default',
    label: 'IP Subject Alternative Names (SANs)',
  },
  issuerRef: {
    editType: 'string',
    helpText:
      'Reference to a existing issuer; either "default" for the configured default issuer, an identifier or the name assigned to the issuer.',
    fieldGroup: 'default',
    type: 'string',
  },
  notAfter: {
    editType: 'string',
    helpText:
      'Set the not after field of the certificate with specified date value. The value format should be given in UTC format YYYY-MM-ddTHH:MM:SSZ',
    fieldGroup: 'default',
    type: 'string',
  },
  otherSans: {
    editType: 'stringArray',
    helpText:
      'Requested other SANs, in an array with the format <oid>;UTF8:<utf8 string value> for each entry.',
    fieldGroup: 'default',
    label: 'Other SANs',
  },
  privateKeyFormat: {
    editType: 'string',
    helpText:
      'Format for the returned private key. Generally the default will be controlled by the "format" parameter as either base64-encoded DER or PEM-encoded DER. However, this can be set to "pkcs8" to have the returned private key contain base64-encoded pkcs8 or PEM-encoded pkcs8 instead. Defaults to "der".',
    possibleValues: ['', 'der', 'pem', 'pkcs8'],
    fieldGroup: 'default',
    defaultValue: 'der',
    type: 'string',
  },
  removeRootsFromChain: {
    editType: 'boolean',
    helpText: 'Whether or not to remove self-signed CA certificates in the output of the ca_chain field.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  serialNumber: {
    editType: 'string',
    helpText:
      "The Subject's requested serial number, if any. See RFC 4519 Section 2.31 'serialNumber' for a description of this field. If you want more than one, specify alternative names in the alt_names map using OID 2.5.4.5. This has no impact on the final certificate's Serial Number field.",
    fieldGroup: 'default',
    type: 'string',
  },
  ttl: {
    editType: 'ttl',
    helpText:
      'The requested Time To Live for the certificate; sets the expiration date. If not specified the role default, backend default, or system default TTL is used, in that order. Cannot be larger than the role max TTL.',
    fieldGroup: 'default',
    label: 'TTL',
  },
  uriSans: {
    editType: 'stringArray',
    helpText: 'The requested URI SANs, if any, in a comma-delimited list.',
    fieldGroup: 'default',
    label: 'URI Subject Alternative Names (SANs)',
  },
  userIds: {
    editType: 'stringArray',
    helpText:
      'The requested user_ids value to place in the subject, if any, in a comma-delimited list. Restricted by allowed_user_ids. Any values are added with OID 0.9.2342.19200300.100.1.1.',
    fieldGroup: 'default',
    label: 'User ID(s)',
  },
};
const pkiCertSign = {
  role: {
    editType: 'string',
    helpText: 'The desired role with configuration for this request',
    fieldValue: 'mutableId',
    fieldGroup: 'default',
    readOnly: true,
    label: 'Role',
    type: 'string',
  },
  altNames: {
    editType: 'string',
    helpText:
      'The requested Subject Alternative Names, if any, in a comma-delimited list. If email protection is enabled for the role, this may contain email addresses.',
    fieldGroup: 'default',
    label: 'DNS/Email Subject Alternative Names (SANs)',
    type: 'string',
  },
  commonName: {
    editType: 'string',
    helpText:
      'The requested common name; if you want more than one, specify the alternative names in the alt_names map. If email protection is enabled in the role, this may be an email address.',
    fieldGroup: 'default',
    type: 'string',
  },
  csr: {
    editType: 'string',
    helpText: 'PEM-format CSR to be signed.',
    fieldGroup: 'default',
    type: 'string',
  },
  excludeCnFromSans: {
    editType: 'boolean',
    helpText:
      'If true, the Common Name will not be included in DNS or Email Subject Alternate Names. Defaults to false (CN is included).',
    fieldGroup: 'default',
    label: 'Exclude Common Name from Subject Alternative Names (SANs)',
    type: 'boolean',
  },
  format: {
    editType: 'string',
    helpText:
      'Format for returned data. Can be "pem", "der", or "pem_bundle". If "pem_bundle", any private key and issuing cert will be appended to the certificate pem. If "der", the value will be base64 encoded. Defaults to "pem".',
    possibleValues: ['pem', 'der', 'pem_bundle'],
    fieldGroup: 'default',
    defaultValue: 'pem',
    type: 'string',
  },
  ipSans: {
    editType: 'stringArray',
    helpText: 'The requested IP SANs, if any, in a comma-delimited list',
    fieldGroup: 'default',
    label: 'IP Subject Alternative Names (SANs)',
  },
  issuerRef: {
    editType: 'string',
    helpText:
      'Reference to a existing issuer; either "default" for the configured default issuer, an identifier or the name assigned to the issuer.',
    fieldGroup: 'default',
    type: 'string',
  },
  notAfter: {
    editType: 'string',
    helpText:
      'Set the not after field of the certificate with specified date value. The value format should be given in UTC format YYYY-MM-ddTHH:MM:SSZ',
    fieldGroup: 'default',
    type: 'string',
  },
  otherSans: {
    editType: 'stringArray',
    helpText:
      'Requested other SANs, in an array with the format <oid>;UTF8:<utf8 string value> for each entry.',
    fieldGroup: 'default',
    label: 'Other SANs',
  },
  privateKeyFormat: {
    editType: 'string',
    helpText:
      'Format for the returned private key. Generally the default will be controlled by the "format" parameter as either base64-encoded DER or PEM-encoded DER. However, this can be set to "pkcs8" to have the returned private key contain base64-encoded pkcs8 or PEM-encoded pkcs8 instead. Defaults to "der".',
    possibleValues: ['', 'der', 'pem', 'pkcs8'],
    fieldGroup: 'default',
    defaultValue: 'der',
    type: 'string',
  },
  removeRootsFromChain: {
    editType: 'boolean',
    helpText: 'Whether or not to remove self-signed CA certificates in the output of the ca_chain field.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  serialNumber: {
    editType: 'string',
    helpText:
      "The Subject's requested serial number, if any. See RFC 4519 Section 2.31 'serialNumber' for a description of this field. If you want more than one, specify alternative names in the alt_names map using OID 2.5.4.5. This has no impact on the final certificate's Serial Number field.",
    fieldGroup: 'default',
    type: 'string',
  },
  ttl: {
    editType: 'ttl',
    helpText:
      'The requested Time To Live for the certificate; sets the expiration date. If not specified the role default, backend default, or system default TTL is used, in that order. Cannot be larger than the role max TTL.',
    fieldGroup: 'default',
    label: 'TTL',
  },
  uriSans: {
    editType: 'stringArray',
    helpText: 'The requested URI SANs, if any, in a comma-delimited list.',
    fieldGroup: 'default',
    label: 'URI Subject Alternative Names (SANs)',
  },
  userIds: {
    editType: 'stringArray',
    helpText:
      'The requested user_ids value to place in the subject, if any, in a comma-delimited list. Restricted by allowed_user_ids. Any values are added with OID 0.9.2342.19200300.100.1.1.',
    fieldGroup: 'default',
    label: 'User ID(s)',
  },
};
const pkiCluster = {
  aiaPath: {
    editType: 'string',
    helpText:
      "Optional URI to this mount's AIA distribution point; may refer to an external non-Vault responder. This is for resolving AIA URLs and providing the {{cluster_aia_path}} template parameter and will not be used for other purposes. As such, unlike path above, this could safely be an insecure transit mechanism (like HTTP without TLS). For example: http://cdn.example.com/pr1/pki",
    fieldGroup: 'default',
    type: 'string',
  },
  path: {
    editType: 'string',
    helpText:
      "Canonical URI to this mount on this performance replication cluster's external address. This is for resolving AIA URLs and providing the {{cluster_path}} template parameter but might be used for other purposes in the future. This should only point back to this particular PR replica and should not ever point to another PR cluster. It may point to any node in the PR replica, including standby nodes, and need not always point to the active node. For example: https://pr1.vault.example.com:8200/v1/pki",
    fieldGroup: 'default',
    type: 'string',
  },
};
const pkiRole = {
  name: {
    editType: 'string',
    helpText: 'Name of the role',
    fieldValue: 'mutableId',
    fieldGroup: 'default',
    readOnly: true,
    label: 'Name',
    type: 'string',
  },
  allowAnyName: {
    editType: 'boolean',
    helpText:
      'If set, clients can request certificates for any domain, regardless of allowed_domains restrictions. See the documentation for more information.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  allowBareDomains: {
    editType: 'boolean',
    helpText:
      'If set, clients can request certificates for the base domains themselves, e.g. "example.com" of domains listed in allowed_domains. This is a separate option as in some cases this can be considered a security threat. See the documentation for more information.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  allowGlobDomains: {
    editType: 'boolean',
    helpText:
      'If set, domains specified in allowed_domains can include shell-style glob patterns, e.g. "ftp*.example.com". See the documentation for more information.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  allowIpSans: {
    editType: 'boolean',
    helpText:
      'If set, IP Subject Alternative Names are allowed. Any valid IP is accepted and No authorization checking is performed.',
    fieldGroup: 'default',
    defaultValue: true,
    label: 'Allow IP Subject Alternative Names',
    type: 'boolean',
  },
  allowLocalhost: {
    editType: 'boolean',
    helpText:
      'Whether to allow "localhost" and "localdomain" as a valid common name in a request, independent of allowed_domains value.',
    fieldGroup: 'default',
    defaultValue: true,
    type: 'boolean',
  },
  allowSubdomains: {
    editType: 'boolean',
    helpText:
      'If set, clients can request certificates for subdomains of domains listed in allowed_domains, including wildcard subdomains. See the documentation for more information.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  allowWildcardCertificates: {
    editType: 'boolean',
    helpText:
      'If set, allows certificates with wildcards in the common name to be issued, conforming to RFC 6125\'s Section 6.4.3; e.g., "*.example.net" or "b*z.example.net". See the documentation for more information.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  allowedDomains: {
    editType: 'stringArray',
    helpText:
      'Specifies the domains this role is allowed to issue certificates for. This is used with the allow_bare_domains, allow_subdomains, and allow_glob_domains to determine matches for the common name, DNS-typed SAN entries, and Email-typed SAN entries of certificates. See the documentation for more information. This parameter accepts a comma-separated string or list of domains.',
    fieldGroup: 'default',
  },
  allowedDomainsTemplate: {
    editType: 'boolean',
    helpText:
      'If set, Allowed domains can be specified using identity template policies. Non-templated domains are also permitted.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  allowedOtherSans: {
    editType: 'stringArray',
    helpText:
      'If set, an array of allowed other names to put in SANs. These values support globbing and must be in the format <oid>;<type>:<value>. Currently only "utf8" is a valid type. All values, including globbing values, must use this syntax, with the exception being a single "*" which allows any OID and any value (but type must still be utf8).',
    fieldGroup: 'default',
    label: 'Allowed Other Subject Alternative Names',
  },
  allowedSerialNumbers: {
    editType: 'stringArray',
    helpText: 'If set, an array of allowed serial numbers to put in Subject. These values support globbing.',
    fieldGroup: 'default',
  },
  allowedUriSans: {
    editType: 'stringArray',
    helpText:
      'If set, an array of allowed URIs for URI Subject Alternative Names. Any valid URI is accepted, these values support globbing.',
    fieldGroup: 'default',
    label: 'Allowed URI Subject Alternative Names',
  },
  allowedUriSansTemplate: {
    editType: 'boolean',
    helpText:
      'If set, Allowed URI SANs can be specified using identity template policies. Non-templated URI SANs are also permitted.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  allowedUserIds: {
    editType: 'stringArray',
    helpText:
      'If set, an array of allowed user-ids to put in user system login name specified here: https://www.rfc-editor.org/rfc/rfc1274#section-9.3.1',
    fieldGroup: 'default',
  },
  backend: {
    editType: 'string',
    helpText: 'Backend Type',
    fieldGroup: 'default',
    type: 'string',
  },
  basicConstraintsValidForNonCa: {
    editType: 'boolean',
    helpText: 'Mark Basic Constraints valid when issuing non-CA certificates.',
    fieldGroup: 'default',
    label: 'Basic Constraints Valid for Non-CA',
    type: 'boolean',
  },
  clientFlag: {
    editType: 'boolean',
    helpText:
      'If set, certificates are flagged for client auth use. Defaults to true. See also RFC 5280 Section 4.2.1.12.',
    fieldGroup: 'default',
    defaultValue: true,
    type: 'boolean',
  },
  cnValidations: {
    editType: 'stringArray',
    helpText:
      "List of allowed validations to run against the Common Name field. Values can include 'email' to validate the CN is a email address, 'hostname' to validate the CN is a valid hostname (potentially including wildcards). When multiple validations are specified, these take OR semantics (either email OR hostname are allowed). The special value 'disabled' allows disabling all CN name validations, allowing for arbitrary non-Hostname, non-Email address CNs.",
    fieldGroup: 'default',
    label: 'Common Name Validations',
  },
  codeSigningFlag: {
    editType: 'boolean',
    helpText:
      'If set, certificates are flagged for code signing use. Defaults to false. See also RFC 5280 Section 4.2.1.12.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  country: {
    editType: 'stringArray',
    helpText: 'If set, Country will be set to this value in certificates issued by this role.',
    fieldGroup: 'default',
  },
  emailProtectionFlag: {
    editType: 'boolean',
    helpText:
      'If set, certificates are flagged for email protection use. Defaults to false. See also RFC 5280 Section 4.2.1.12.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  enforceHostnames: {
    editType: 'boolean',
    helpText:
      'If set, only valid host names are allowed for CN and DNS SANs, and the host part of email addresses. Defaults to true.',
    fieldGroup: 'default',
    defaultValue: true,
    type: 'boolean',
  },
  extKeyUsage: {
    editType: 'stringArray',
    helpText:
      'A comma-separated string or list of extended key usages. Valid values can be found at https://golang.org/pkg/crypto/x509/#ExtKeyUsage -- simply drop the "ExtKeyUsage" part of the name. To remove all key usages from being set, set this value to an empty list. See also RFC 5280 Section 4.2.1.12.',
    fieldGroup: 'default',
    label: 'Extended Key Usage',
  },
  extKeyUsageOids: {
    editType: 'stringArray',
    helpText: 'A comma-separated string or list of extended key usage oids.',
    fieldGroup: 'default',
    label: 'Extended Key Usage OIDs',
  },
  generateLease: {
    editType: 'boolean',
    helpText:
      'If set, certificates issued/signed against this role will have Vault leases attached to them. Defaults to "false". Certificates can be added to the CRL by "vault revoke <lease_id>" when certificates are associated with leases. It can also be done using the "pki/revoke" endpoint. However, when lease generation is disabled, invoking "pki/revoke" would be the only way to add the certificates to the CRL. When large number of certificates are generated with long lifetimes, it is recommended that lease generation be disabled, as large amount of leases adversely affect the startup time of Vault.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  issuerRef: {
    editType: 'string',
    helpText: 'Reference to the issuer used to sign requests serviced by this role.',
    fieldGroup: 'default',
    type: 'string',
  },
  keyBits: {
    editType: 'number',
    helpText:
      'The number of bits to use. Allowed values are 0 (universal default); with rsa key_type: 2048 (default), 3072, or 4096; with ec key_type: 224, 256 (default), 384, or 521; ignored with ed25519.',
    fieldGroup: 'default',
    type: 'number',
  },
  keyType: {
    editType: 'string',
    helpText:
      'The type of key to use; defaults to RSA. "rsa" "ec", "ed25519" and "any" are the only valid values.',
    possibleValues: ['rsa', 'ec', 'ed25519', 'any'],
    fieldGroup: 'default',
    type: 'string',
  },
  keyUsage: {
    editType: 'stringArray',
    helpText:
      'A comma-separated string or list of key usages (not extended key usages). Valid values can be found at https://golang.org/pkg/crypto/x509/#KeyUsage -- simply drop the "KeyUsage" part of the name. To remove all key usages from being set, set this value to an empty list. See also RFC 5280 Section 4.2.1.3.',
    fieldGroup: 'default',
    defaultValue: 'DigitalSignature,KeyAgreement,KeyEncipherment',
  },
  locality: {
    editType: 'stringArray',
    helpText: 'If set, Locality will be set to this value in certificates issued by this role.',
    fieldGroup: 'default',
    label: 'Locality/City',
  },
  maxTtl: {
    editType: 'ttl',
    helpText: 'The maximum allowed lease duration. If not set, defaults to the system maximum lease TTL.',
    fieldGroup: 'default',
    label: 'Max TTL',
  },
  noStore: {
    editType: 'boolean',
    helpText:
      'If set, certificates issued/signed against this role will not be stored in the storage backend. This can improve performance when issuing large numbers of certificates. However, certificates issued in this way cannot be enumerated or revoked, so this option is recommended only for certificates that are non-sensitive, or extremely short-lived. This option implies a value of "false" for "generate_lease".',
    fieldGroup: 'default',
    type: 'boolean',
  },
  notAfter: {
    editType: 'string',
    helpText:
      'Set the not after field of the certificate with specified date value. The value format should be given in UTC format YYYY-MM-ddTHH:MM:SSZ.',
    fieldGroup: 'default',
    type: 'string',
  },
  notBeforeDuration: {
    editType: 'ttl',
    helpText: 'The duration before now which the certificate needs to be backdated by.',
    fieldGroup: 'default',
    defaultValue: 30,
  },
  organization: {
    editType: 'stringArray',
    helpText: 'If set, O (Organization) will be set to this value in certificates issued by this role.',
    fieldGroup: 'default',
  },
  ou: {
    editType: 'stringArray',
    helpText:
      'If set, OU (OrganizationalUnit) will be set to this value in certificates issued by this role.',
    fieldGroup: 'default',
    label: 'Organizational Unit',
  },
  policyIdentifiers: {
    editType: 'stringArray',
    helpText:
      'A comma-separated string or list of policy OIDs, or a JSON list of qualified policy information, which must include an oid, and may include a notice and/or cps url, using the form [{"oid"="1.3.6.1.4.1.7.8","notice"="I am a user Notice"}, {"oid"="1.3.6.1.4.1.44947.1.2.4 ","cps"="https://example.com"}].',
    fieldGroup: 'default',
  },
  postalCode: {
    editType: 'stringArray',
    helpText: 'If set, Postal Code will be set to this value in certificates issued by this role.',
    fieldGroup: 'default',
  },
  province: {
    editType: 'stringArray',
    helpText: 'If set, Province will be set to this value in certificates issued by this role.',
    fieldGroup: 'default',
    label: 'Province/State',
  },
  requireCn: {
    editType: 'boolean',
    helpText: "If set to false, makes the 'common_name' field optional while generating a certificate.",
    fieldGroup: 'default',
    label: 'Require Common Name',
    type: 'boolean',
  },
  serverFlag: {
    editType: 'boolean',
    helpText:
      'If set, certificates are flagged for server auth use. Defaults to true. See also RFC 5280 Section 4.2.1.12.',
    fieldGroup: 'default',
    defaultValue: true,
    type: 'boolean',
  },
  signatureBits: {
    editType: 'number',
    helpText:
      'The number of bits to use in the signature algorithm; accepts 256 for SHA-2-256, 384 for SHA-2-384, and 512 for SHA-2-512. Defaults to 0 to automatically detect based on key length (SHA-2-256 for RSA keys, and matching the curve size for NIST P-Curves).',
    fieldGroup: 'default',
    type: 'number',
  },
  streetAddress: {
    editType: 'stringArray',
    helpText: 'If set, Street Address will be set to this value in certificates issued by this role.',
    fieldGroup: 'default',
  },
  ttl: {
    editType: 'ttl',
    helpText:
      'The lease duration (validity period of the certificate) if no specific lease duration is requested. The lease duration controls the expiration of certificates issued by this backend. Defaults to the system default value or the value of max_ttl, whichever is shorter.',
    fieldGroup: 'default',
    label: 'TTL',
  },
  useCsrCommonName: {
    editType: 'boolean',
    helpText:
      'If set, when used with a signing profile, the common name in the CSR will be used. This does *not* include any requested Subject Alternative Names; use use_csr_sans for that. Defaults to true.',
    fieldGroup: 'default',
    defaultValue: true,
    label: 'Use CSR Common Name',
    type: 'boolean',
  },
  useCsrSans: {
    editType: 'boolean',
    helpText:
      'If set, when used with a signing profile, the SANs in the CSR will be used. This does *not* include the Common Name (cn); use use_csr_common_name for that. Defaults to true.',
    fieldGroup: 'default',
    defaultValue: true,
    label: 'Use CSR Subject Alternative Names',
    type: 'boolean',
  },
  usePss: {
    editType: 'boolean',
    helpText: 'Whether or not to use PSS signatures when using a RSA key-type issuer. Defaults to false.',
    fieldGroup: 'default',
    type: 'boolean',
  },
};
const pkiSignCsr = {
  issuerRef: {
    editType: 'string',
    helpText:
      'Reference to a existing issuer; either "default" for the configured default issuer, an identifier or the name assigned to the issuer.',
    fieldValue: 'mutableId',
    fieldGroup: 'default',
    readOnly: true,
    label: 'Issuer ref',
    type: 'string',
  },
  altNames: {
    editType: 'string',
    helpText:
      'The requested Subject Alternative Names, if any, in a comma-delimited list. May contain both DNS names and email addresses.',
    fieldGroup: 'default',
    label: 'DNS/Email Subject Alternative Names (SANs)',
    type: 'string',
  },
  commonName: {
    editType: 'string',
    helpText:
      'The requested common name; if you want more than one, specify the alternative names in the alt_names map. If not specified when signing, the common name will be taken from the CSR; other names must still be specified in alt_names or ip_sans.',
    fieldGroup: 'default',
    type: 'string',
  },
  country: {
    editType: 'stringArray',
    helpText: 'If set, Country will be set to this value.',
    fieldGroup: 'default',
  },
  csr: {
    editType: 'string',
    helpText: 'PEM-format CSR to be signed.',
    fieldGroup: 'default',
    type: 'string',
  },
  excludeCnFromSans: {
    editType: 'boolean',
    helpText:
      'If true, the Common Name will not be included in DNS or Email Subject Alternate Names. Defaults to false (CN is included).',
    fieldGroup: 'default',
    label: 'Exclude Common Name from Subject Alternative Names (SANs)',
    type: 'boolean',
  },
  format: {
    editType: 'string',
    helpText:
      'Format for returned data. Can be "pem", "der", or "pem_bundle". If "pem_bundle", any private key and issuing cert will be appended to the certificate pem. If "der", the value will be base64 encoded. Defaults to "pem".',
    possibleValues: ['pem', 'der', 'pem_bundle'],
    fieldGroup: 'default',
    defaultValue: 'pem',
    type: 'string',
  },
  ipSans: {
    editType: 'stringArray',
    helpText: 'The requested IP SANs, if any, in a comma-delimited list',
    fieldGroup: 'default',
    label: 'IP Subject Alternative Names (SANs)',
  },
  issuerName: {
    editType: 'string',
    helpText:
      "Provide a name to the generated or existing issuer, the name must be unique across all issuers and not be the reserved value 'default'",
    fieldGroup: 'default',
    type: 'string',
  },
  locality: {
    editType: 'stringArray',
    helpText: 'If set, Locality will be set to this value.',
    fieldGroup: 'default',
    label: 'Locality/City',
  },
  maxPathLength: {
    editType: 'number',
    helpText: 'The maximum allowable path length',
    fieldGroup: 'default',
    type: 'number',
  },
  notAfter: {
    editType: 'string',
    helpText:
      'Set the not after field of the certificate with specified date value. The value format should be given in UTC format YYYY-MM-ddTHH:MM:SSZ',
    fieldGroup: 'default',
    type: 'string',
  },
  notBeforeDuration: {
    editType: 'ttl',
    helpText: 'The duration before now which the certificate needs to be backdated by.',
    fieldGroup: 'default',
    defaultValue: 30,
  },
  organization: {
    editType: 'stringArray',
    helpText: 'If set, O (Organization) will be set to this value.',
    fieldGroup: 'default',
  },
  otherSans: {
    editType: 'stringArray',
    helpText:
      'Requested other SANs, in an array with the format <oid>;UTF8:<utf8 string value> for each entry.',
    fieldGroup: 'default',
    label: 'Other SANs',
  },
  ou: {
    editType: 'stringArray',
    helpText: 'If set, OU (OrganizationalUnit) will be set to this value.',
    fieldGroup: 'default',
    label: 'OU (Organizational Unit)',
  },
  permittedDnsDomains: {
    editType: 'stringArray',
    helpText:
      'Domains for which this certificate is allowed to sign or issue child certificates. If set, all DNS names (subject and alt) on child certs must be exact matches or subsets of the given domains (see https://tools.ietf.org/html/rfc5280#section-4.2.1.10).',
    fieldGroup: 'default',
    label: 'Permitted DNS Domains',
  },
  postalCode: {
    editType: 'stringArray',
    helpText: 'If set, Postal Code will be set to this value.',
    fieldGroup: 'default',
    label: 'Postal Code',
  },
  privateKeyFormat: {
    editType: 'string',
    helpText:
      'Format for the returned private key. Generally the default will be controlled by the "format" parameter as either base64-encoded DER or PEM-encoded DER. However, this can be set to "pkcs8" to have the returned private key contain base64-encoded pkcs8 or PEM-encoded pkcs8 instead. Defaults to "der".',
    possibleValues: ['', 'der', 'pem', 'pkcs8'],
    fieldGroup: 'default',
    defaultValue: 'der',
    type: 'string',
  },
  province: {
    editType: 'stringArray',
    helpText: 'If set, Province will be set to this value.',
    fieldGroup: 'default',
    label: 'Province/State',
  },
  serialNumber: {
    editType: 'string',
    helpText:
      "The Subject's requested serial number, if any. See RFC 4519 Section 2.31 'serialNumber' for a description of this field. If you want more than one, specify alternative names in the alt_names map using OID 2.5.4.5. This has no impact on the final certificate's Serial Number field.",
    fieldGroup: 'default',
    type: 'string',
  },
  signatureBits: {
    editType: 'number',
    helpText:
      'The number of bits to use in the signature algorithm; accepts 256 for SHA-2-256, 384 for SHA-2-384, and 512 for SHA-2-512. Defaults to 0 to automatically detect based on key length (SHA-2-256 for RSA keys, and matching the curve size for NIST P-Curves).',
    fieldGroup: 'default',
    type: 'number',
  },
  skid: {
    editType: 'string',
    helpText:
      "Value for the Subject Key Identifier field (RFC 5280 Section 4.2.1.2). This value should ONLY be used when cross-signing to mimic the existing certificate's SKID value; this is necessary to allow certain TLS implementations (such as OpenSSL) which use SKID/AKID matches in chain building to restrict possible valid chains. Specified as a string in hex format. Default is empty, allowing Vault to automatically calculate the SKID according to method one in the above RFC section.",
    fieldGroup: 'default',
    type: 'string',
  },
  streetAddress: {
    editType: 'stringArray',
    helpText: 'If set, Street Address will be set to this value.',
    fieldGroup: 'default',
    label: 'Street Address',
  },
  ttl: {
    editType: 'ttl',
    helpText:
      'The requested Time To Live for the certificate; sets the expiration date. If not specified the role default, backend default, or system default TTL is used, in that order. Cannot be larger than the mount max TTL. Note: this only has an effect when generating a CA cert or signing a CA cert, not when generating a CSR for an intermediate CA.',
    fieldGroup: 'default',
    label: 'TTL',
  },
  uriSans: {
    editType: 'stringArray',
    helpText: 'The requested URI SANs, if any, in a comma-delimited list.',
    fieldGroup: 'default',
    label: 'URI Subject Alternative Names (SANs)',
  },
  useCsrValues: {
    editType: 'boolean',
    helpText:
      'If true, then: 1) Subject information, including names and alternate names, will be preserved from the CSR rather than using values provided in the other parameters to this path; 2) Any key usages requested in the CSR will be added to the basic set of key usages used for CA certs signed by this path; for instance, the non-repudiation flag; 3) Extensions requested in the CSR will be copied into the issued certificate.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  usePss: {
    editType: 'boolean',
    helpText: 'Whether or not to use PSS signatures when using a RSA key-type issuer. Defaults to false.',
    fieldGroup: 'default',
    type: 'boolean',
  },
};
const pkiTidy = {
  acmeAccountSafetyBuffer: {
    editType: 'ttl',
    helpText:
      'The amount of time that must pass after creation that an account with no orders is marked revoked, and the amount of time after being marked revoked or deactivated.',
    fieldGroup: 'default',
  },
  enabled: {
    editType: 'boolean',
    helpText: 'Set to true to enable automatic tidy operations.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  intervalDuration: {
    editType: 'ttl',
    helpText:
      'Interval at which to run an auto-tidy operation. This is the time between tidy invocations (after one finishes to the start of the next). Running a manual tidy will reset this duration.',
    fieldGroup: 'default',
  },
  issuerSafetyBuffer: {
    editType: 'ttl',
    helpText:
      "The amount of extra time that must have passed beyond issuer's expiration before it is removed from the backend storage. Defaults to 8760 hours (1 year).",
    fieldGroup: 'default',
  },
  maintainStoredCertificateCounts: {
    editType: 'boolean',
    helpText:
      'This configures whether stored certificates are counted upon initialization of the backend, and whether during normal operation, a running count of certificates stored is maintained.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  pauseDuration: {
    editType: 'string',
    helpText:
      'The amount of time to wait between processing certificates. This allows operators to change the execution profile of tidy to take consume less resources by slowing down how long it takes to run. Note that the entire list of certificates will be stored in memory during the entire tidy operation, but resources to read/process/update existing entries will be spread out over a greater period of time. By default this is zero seconds.',
    fieldGroup: 'default',
    type: 'string',
  },
  publishStoredCertificateCountMetrics: {
    editType: 'boolean',
    helpText:
      'This configures whether the stored certificate count is published to the metrics consumer. It does not affect if the stored certificate count is maintained, and if maintained, it will be available on the tidy-status endpoint.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  revocationQueueSafetyBuffer: {
    editType: 'ttl',
    helpText:
      'The amount of time that must pass from the cross-cluster revocation request being initiated to when it will be slated for removal. Setting this too low may remove valid revocation requests before the owning cluster has a chance to process them, especially if the cluster is offline.',
    fieldGroup: 'default',
  },
  safetyBuffer: {
    editType: 'ttl',
    helpText:
      'The amount of extra time that must have passed beyond certificate expiration before it is removed from the backend storage and/or revocation list. Defaults to 72 hours.',
    fieldGroup: 'default',
  },
  tidyAcme: {
    editType: 'boolean',
    helpText:
      'Set to true to enable tidying ACME accounts, orders and authorizations. ACME orders are tidied (deleted) safety_buffer after the certificate associated with them expires, or after the order and relevant authorizations have expired if no certificate was produced. Authorizations are tidied with the corresponding order. When a valid ACME Account is at least acme_account_safety_buffer old, and has no remaining orders associated with it, the account is marked as revoked. After another acme_account_safety_buffer has passed from the revocation or deactivation date, a revoked or deactivated ACME account is deleted.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  tidyCertStore: {
    editType: 'boolean',
    helpText: 'Set to true to enable tidying up the certificate store',
    fieldGroup: 'default',
    type: 'boolean',
  },
  tidyCrossClusterRevokedCerts: {
    editType: 'boolean',
    helpText:
      'Set to true to enable tidying up the cross-cluster revoked certificate store. Only runs on the active primary node.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  tidyExpiredIssuers: {
    editType: 'boolean',
    helpText:
      'Set to true to automatically remove expired issuers past the issuer_safety_buffer. No keys will be removed as part of this operation.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  tidyMoveLegacyCaBundle: {
    editType: 'boolean',
    helpText:
      'Set to true to move the legacy ca_bundle from /config/ca_bundle to /config/ca_bundle.bak. This prevents downgrades to pre-Vault 1.11 versions (as older PKI engines do not know about the new multi-issuer storage layout), but improves the performance on seal wrapped PKI mounts. This will only occur if at least issuer_safety_buffer time has occurred after the initial storage migration. This backup is saved in case of an issue in future migrations. Operators may consider removing it via sys/raw if they desire. The backup will be removed via a DELETE /root call, but note that this removes ALL issuers within the mount (and is thus not desirable in most operational scenarios).',
    fieldGroup: 'default',
    type: 'boolean',
  },
  tidyRevocationList: {
    editType: 'boolean',
    helpText: "Deprecated; synonym for 'tidy_revoked_certs",
    fieldGroup: 'default',
    type: 'boolean',
  },
  tidyRevocationQueue: {
    editType: 'boolean',
    helpText:
      "Set to true to remove stale revocation queue entries that haven't been confirmed by any active cluster. Only runs on the active primary node",
    fieldGroup: 'default',
    type: 'boolean',
  },
  tidyRevokedCertIssuerAssociations: {
    editType: 'boolean',
    helpText:
      'Set to true to validate issuer associations on revocation entries. This helps increase the performance of CRL building and OCSP responses.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  tidyRevokedCerts: {
    editType: 'boolean',
    helpText:
      'Set to true to expire all revoked and expired certificates, removing them both from the CRL and from storage. The CRL will be rotated if this causes any values to be removed.',
    fieldGroup: 'default',
    type: 'boolean',
  },
};
const pkiUrls = {
  crlDistributionPoints: {
    editType: 'stringArray',
    helpText:
      'Comma-separated list of URLs to be used for the CRL distribution points attribute. See also RFC 5280 Section 4.2.1.13.',
    fieldGroup: 'default',
  },
  enableTemplating: {
    editType: 'boolean',
    helpText:
      "Whether or not to enabling templating of the above AIA fields. When templating is enabled the special values '{{issuer_id}}', '{{cluster_path}}', and '{{cluster_aia_path}}' are available, but the addresses are not checked for URI validity until issuance time. Using '{{cluster_path}}' requires /config/cluster's 'path' member to be set on all PR Secondary clusters and using '{{cluster_aia_path}}' requires /config/cluster's 'aia_path' member to be set on all PR secondary clusters.",
    fieldGroup: 'default',
    type: 'boolean',
  },
  issuingCertificates: {
    editType: 'stringArray',
    helpText:
      'Comma-separated list of URLs to be used for the issuing certificate attribute. See also RFC 5280 Section 4.2.2.1.',
    fieldGroup: 'default',
  },
  ocspServers: {
    editType: 'stringArray',
    helpText:
      'Comma-separated list of URLs to be used for the OCSP servers attribute. See also RFC 5280 Section 4.2.2.1.',
    fieldGroup: 'default',
  },
};

/* Auth Engines */
const userpassUser = {
  username: {
    editType: 'string',
    helpText: 'Username for this user.',
    fieldValue: 'mutableId',
    fieldGroup: 'default',
    readOnly: true,
    label: 'Username',
    type: 'string',
  },
  password: {
    editType: 'string',
    helpText: 'Password for this user.',
    fieldGroup: 'default',
    sensitive: true,
    type: 'string',
  },
  tokenBoundCidrs: {
    editType: 'stringArray',
    helpText:
      'A list of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.',
    fieldGroup: 'Tokens',
    label: "Generated Token's Bound CIDRs",
  },
  tokenExplicitMaxTtl: {
    editType: 'ttl',
    helpText:
      'If set, tokens created via this role carry an explicit maximum TTL. During renewal, the current maximum TTL values of the role and the mount are not checked for changes, and any updates to these values will have no effect on the token being renewed.',
    fieldGroup: 'Tokens',
    label: "Generated Token's Explicit Maximum TTL",
  },
  tokenMaxTtl: {
    editType: 'ttl',
    helpText: 'The maximum lifetime of the generated token',
    fieldGroup: 'Tokens',
    label: "Generated Token's Maximum TTL",
  },
  tokenNoDefaultPolicy: {
    editType: 'boolean',
    helpText: "If true, the 'default' policy will not automatically be added to generated tokens",
    fieldGroup: 'Tokens',
    label: "Do Not Attach 'default' Policy To Generated Tokens",
    type: 'boolean',
  },
  tokenNumUses: {
    editType: 'number',
    helpText: 'The maximum number of times a token may be used, a value of zero means unlimited',
    fieldGroup: 'Tokens',
    label: 'Maximum Uses of Generated Tokens',
    type: 'number',
  },
  tokenPeriod: {
    editType: 'ttl',
    helpText:
      'If set, tokens created via this role will have no max lifetime; instead, their renewal period will be fixed to this value. This takes an integer number of seconds, or a string duration (e.g. "24h").',
    fieldGroup: 'Tokens',
    label: "Generated Token's Period",
  },
  tokenPolicies: {
    editType: 'stringArray',
    helpText: 'A list of policies that will apply to the generated token for this user.',
    fieldGroup: 'Tokens',
    label: "Generated Token's Policies",
  },
  tokenTtl: {
    editType: 'ttl',
    helpText: 'The initial ttl of the token to generate',
    fieldGroup: 'Tokens',
    label: "Generated Token's Initial TTL",
  },
  tokenType: {
    editType: 'string',
    helpText: 'The type of token to generate, service or batch',
    fieldGroup: 'Tokens',
    label: "Generated Token's Type",
    type: 'string',
  },
};

const azureConfig = {
  clientId: {
    editType: 'string',
    fieldGroup: 'default',
    helpText:
      'The OAuth2 client id to connection to Azure. This value can also be provided with the AZURE_CLIENT_ID environment variable.',
    label: 'Client ID',
    type: 'string',
  },
  clientSecret: {
    editType: 'string',
    fieldGroup: 'default',
    helpText:
      'The OAuth2 client secret to connection to Azure. This value can also be provided with the AZURE_CLIENT_SECRET environment variable.',
    type: 'string',
  },
  environment: {
    editType: 'string',
    fieldGroup: 'default',
    helpText:
      'The Azure environment name. If not provided, AzurePublicCloud is used. This value can also be provided with the AZURE_ENVIRONMENT environment variable.',
    type: 'string',
  },
  maxRetries: {
    editType: 'number',
    fieldGroup: 'default',
    helpText: 'The maximum number of attempts a failed operation will be retried before producing an error.',
    type: 'number',
  },
  maxRetryDelay: {
    editType: 'ttl',
    fieldGroup: 'default',
    helpText: 'The maximum delay allowed before retrying an operation.',
  },
  resource: {
    editType: 'string',
    fieldGroup: 'default',
    helpText:
      'The resource URL for the vault application in Azure Active Directory. This value can also be provided with the AZURE_AD_RESOURCE environment variable.',
    type: 'string',
  },
  retryDelay: {
    editType: 'ttl',
    fieldGroup: 'default',
    helpText: 'The initial amount of delay to use before retrying an operation, increasing exponentially.',
  },
  rootPasswordTtl: {
    editType: 'ttl',
    fieldGroup: 'default',
    helpText:
      'The TTL of the root password in Azure. This can be either a number of seconds or a time formatted duration (ex: 24h, 48ds)',
  },
  tenantId: {
    editType: 'string',
    fieldGroup: 'default',
    helpText:
      'The tenant id for the Azure Active Directory. This is sometimes referred to as Directory ID in AD. This value can also be provided with the AZURE_TENANT_ID environment variable.',
    label: 'Tenant ID',
    type: 'string',
  },
};

const certConfig = {
  disableBinding: {
    editType: 'boolean',
    helpText:
      'If set, during renewal, skips the matching of presented client identity with the client identity used during login. Defaults to false.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  enableIdentityAliasMetadata: {
    editType: 'boolean',
    helpText:
      'If set, metadata of the certificate including the metadata corresponding to allowed_metadata_extensions will be stored in the alias. Defaults to false.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  ocspCacheSize: {
    editType: 'number',
    helpText: 'The size of the in memory OCSP response cache, shared by all configured certs',
    fieldGroup: 'default',
    type: 'number',
  },
};
const certCert = {
  name: {
    editType: 'string',
    helpText: 'The name of the certificate',
    fieldValue: 'mutableId',
    fieldGroup: 'default',
    readOnly: true,
    label: 'Name',
    type: 'string',
  },
  allowedCommonNames: {
    editType: 'stringArray',
    helpText: 'A list of names. At least one must exist in the Common Name. Supports globbing.',
    fieldGroup: 'Constraints',
  },
  allowedDnsSans: {
    editType: 'stringArray',
    helpText: 'A list of DNS names. At least one must exist in the SANs. Supports globbing.',
    fieldGroup: 'Constraints',
    label: 'Allowed DNS SANs',
  },
  allowedEmailSans: {
    editType: 'stringArray',
    helpText: 'A list of Email Addresses. At least one must exist in the SANs. Supports globbing.',
    fieldGroup: 'Constraints',
    label: 'Allowed Email SANs',
  },
  allowedMetadataExtensions: {
    editType: 'stringArray',
    helpText:
      'A list of OID extensions. Upon successful authentication, these extensions will be added as metadata if they are present in the certificate. The metadata key will be the string consisting of the OID numbers separated by a dash (-) instead of a dot (.) to allow usage in ACL templates.',
    fieldGroup: 'default',
  },
  allowedNames: {
    editType: 'stringArray',
    helpText:
      'A list of names. At least one must exist in either the Common Name or SANs. Supports globbing. This parameter is deprecated, please use allowed_common_names, allowed_dns_sans, allowed_email_sans, allowed_uri_sans.',
    fieldGroup: 'Constraints',
  },
  allowedOrganizationalUnits: {
    editType: 'stringArray',
    helpText: 'A list of Organizational Units names. At least one must exist in the OU field.',
    fieldGroup: 'Constraints',
  },
  allowedUriSans: {
    editType: 'stringArray',
    helpText: 'A list of URIs. At least one must exist in the SANs. Supports globbing.',
    fieldGroup: 'Constraints',
    label: 'Allowed URI SANs',
  },
  certificate: {
    editType: 'file',
    helpText: 'The public certificate that should be trusted. Must be x509 PEM encoded.',
    fieldGroup: 'default',
    type: 'string',
  },
  displayName: {
    editType: 'string',
    helpText: 'The display name to use for clients using this certificate.',
    fieldGroup: 'default',
    type: 'string',
  },
  ocspCaCertificates: {
    editType: 'file',
    helpText: 'Any additional CA certificates needed to communicate with OCSP servers',
    fieldGroup: 'default',
    type: 'string',
  },
  ocspEnabled: {
    editType: 'boolean',
    helpText: 'Whether to attempt OCSP verification of certificates at login',
    fieldGroup: 'default',
    type: 'boolean',
  },
  ocspFailOpen: {
    editType: 'boolean',
    helpText:
      'If set to true, if an OCSP revocation cannot be made successfully, login will proceed rather than failing. If false, failing to get an OCSP status fails the request.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  ocspQueryAllServers: {
    editType: 'boolean',
    helpText:
      'If set to true, rather than accepting the first successful OCSP response, query all servers and consider the certificate valid only if all servers agree.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  ocspServersOverride: {
    editType: 'stringArray',
    helpText:
      'A list of OCSP server addresses. If unset, the OCSP server is determined from the AuthorityInformationAccess extension on the certificate being inspected.',
    fieldGroup: 'default',
  },
  requiredExtensions: {
    editType: 'stringArray',
    helpText:
      "A list of extensions formatted as 'oid:value'. Expects the extension value to be some type of ASN1 encoded string. All values much match. Supports globbing on 'value'.",
    fieldGroup: 'default',
  },
  tokenBoundCidrs: {
    editType: 'stringArray',
    helpText:
      'A list of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.',
    fieldGroup: 'Tokens',
    label: "Generated Token's Bound CIDRs",
  },
  tokenExplicitMaxTtl: {
    editType: 'ttl',
    helpText:
      'If set, tokens created via this role carry an explicit maximum TTL. During renewal, the current maximum TTL values of the role and the mount are not checked for changes, and any updates to these values will have no effect on the token being renewed.',
    fieldGroup: 'Tokens',
    label: "Generated Token's Explicit Maximum TTL",
  },
  tokenMaxTtl: {
    editType: 'ttl',
    helpText: 'The maximum lifetime of the generated token',
    fieldGroup: 'Tokens',
    label: "Generated Token's Maximum TTL",
  },
  tokenNoDefaultPolicy: {
    editType: 'boolean',
    helpText: "If true, the 'default' policy will not automatically be added to generated tokens",
    fieldGroup: 'Tokens',
    label: "Do Not Attach 'default' Policy To Generated Tokens",
    type: 'boolean',
  },
  tokenNumUses: {
    editType: 'number',
    helpText: 'The maximum number of times a token may be used, a value of zero means unlimited',
    fieldGroup: 'Tokens',
    label: 'Maximum Uses of Generated Tokens',
    type: 'number',
  },
  tokenPeriod: {
    editType: 'ttl',
    helpText:
      'If set, tokens created via this role will have no max lifetime; instead, their renewal period will be fixed to this value. This takes an integer number of seconds, or a string duration (e.g. "24h").',
    fieldGroup: 'Tokens',
    label: "Generated Token's Period",
  },
  tokenPolicies: {
    editType: 'stringArray',
    helpText: 'A list of policies that will apply to the generated token for this user.',
    fieldGroup: 'Tokens',
    label: "Generated Token's Policies",
  },
  tokenTtl: {
    editType: 'ttl',
    helpText: 'The initial ttl of the token to generate',
    fieldGroup: 'Tokens',
    label: "Generated Token's Initial TTL",
  },
  tokenType: {
    editType: 'string',
    helpText: 'The type of token to generate, service or batch',
    fieldGroup: 'Tokens',
    label: "Generated Token's Type",
    type: 'string',
  },
};

const gcpConfig = {
  credentials: {
    editType: 'string',
    helpText:
      'Google credentials JSON that Vault will use to verify users against GCP APIs. If not specified, will use application default credentials',
    fieldGroup: 'default',
    label: 'Credentials',
    type: 'string',
  },
  customEndpoint: {
    editType: 'object',
    helpText: 'Specifies overrides for various Google API Service Endpoints used in requests.',
    fieldGroup: 'default',
    type: 'object',
  },
  gceAlias: {
    editType: 'string',
    helpText: 'Indicates what value to use when generating an alias for GCE authentications.',
    fieldGroup: 'default',
    type: 'string',
  },
  gceMetadata: {
    editType: 'stringArray',
    helpText:
      "The metadata to include on the aliases and audit logs generated by this plugin. When set to 'default', includes: instance_creation_timestamp, instance_id, instance_name, project_id, project_number, role, service_account_id, service_account_email, zone. Not editing this field means the 'default' fields are included. Explicitly setting this field to empty overrides the 'default' and means no metadata will be included. If not using 'default', explicit fields must be sent like: 'field1,field2'.",
    fieldGroup: 'default',
    defaultValue: 'field1,field2',
    label: 'gce_metadata',
  },
  iamAlias: {
    editType: 'string',
    helpText: 'Indicates what value to use when generating an alias for IAM authentications.',
    fieldGroup: 'default',
    type: 'string',
  },
  iamMetadata: {
    editType: 'stringArray',
    helpText:
      "The metadata to include on the aliases and audit logs generated by this plugin. When set to 'default', includes: project_id, role, service_account_id, service_account_email. Not editing this field means the 'default' fields are included. Explicitly setting this field to empty overrides the 'default' and means no metadata will be included. If not using 'default', explicit fields must be sent like: 'field1,field2'.",
    fieldGroup: 'default',
    defaultValue: 'field1,field2',
    label: 'iam_metadata',
  },
};

const githubConfig = {
  baseUrl: {
    editType: 'string',
    helpText:
      'The API endpoint to use. Useful if you are running GitHub Enterprise or an API-compatible authentication server.',
    fieldGroup: 'GitHub Options',
    label: 'Base URL',
    type: 'string',
  },
  organization: {
    editType: 'string',
    helpText: 'The organization users must be part of',
    fieldGroup: 'default',
    type: 'string',
  },
  organizationId: {
    editType: 'number',
    helpText: 'The ID of the organization users must be part of',
    fieldGroup: 'default',
    type: 'number',
  },
  tokenBoundCidrs: {
    editType: 'stringArray',
    helpText:
      'A list of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.',
    fieldGroup: 'Tokens',
    label: "Generated Token's Bound CIDRs",
  },
  tokenExplicitMaxTtl: {
    editType: 'ttl',
    helpText:
      'If set, tokens created via this role carry an explicit maximum TTL. During renewal, the current maximum TTL values of the role and the mount are not checked for changes, and any updates to these values will have no effect on the token being renewed.',
    fieldGroup: 'Tokens',
    label: "Generated Token's Explicit Maximum TTL",
  },
  tokenMaxTtl: {
    editType: 'ttl',
    helpText: 'The maximum lifetime of the generated token',
    fieldGroup: 'Tokens',
    label: "Generated Token's Maximum TTL",
  },
  tokenNoDefaultPolicy: {
    editType: 'boolean',
    helpText: "If true, the 'default' policy will not automatically be added to generated tokens",
    fieldGroup: 'Tokens',
    label: "Do Not Attach 'default' Policy To Generated Tokens",
    type: 'boolean',
  },
  tokenNumUses: {
    editType: 'number',
    helpText: 'The maximum number of times a token may be used, a value of zero means unlimited',
    fieldGroup: 'Tokens',
    label: 'Maximum Uses of Generated Tokens',
    type: 'number',
  },
  tokenPeriod: {
    editType: 'ttl',
    helpText:
      'If set, tokens created via this role will have no max lifetime; instead, their renewal period will be fixed to this value. This takes an integer number of seconds, or a string duration (e.g. "24h").',
    fieldGroup: 'Tokens',
    label: "Generated Token's Period",
  },
  tokenPolicies: {
    editType: 'stringArray',
    helpText: 'A list of policies that will apply to the generated token for this user.',
    fieldGroup: 'Tokens',
    label: "Generated Token's Policies",
  },
  tokenTtl: {
    editType: 'ttl',
    helpText: 'The initial ttl of the token to generate',
    fieldGroup: 'Tokens',
    label: "Generated Token's Initial TTL",
  },
  tokenType: {
    editType: 'string',
    helpText: 'The type of token to generate, service or batch',
    fieldGroup: 'Tokens',
    label: "Generated Token's Type",
    type: 'string',
  },
};

const jwtConfig = {
  boundIssuer: {
    editType: 'string',
    helpText: "The value against which to match the 'iss' claim in a JWT. Optional.",
    fieldGroup: 'default',
    type: 'string',
  },
  defaultRole: {
    editType: 'string',
    helpText:
      'The default role to use if none is provided during login. If not set, a role is required during login.',
    fieldGroup: 'default',
    type: 'string',
  },
  jwksCaPem: {
    editType: 'string',
    helpText:
      'The CA certificate or chain of certificates, in PEM format, to use to validate connections to the JWKS URL. If not set, system certificates are used.',
    fieldGroup: 'default',
    type: 'string',
  },
  jwksUrl: {
    editType: 'string',
    helpText:
      'JWKS URL to use to authenticate signatures. Cannot be used with "oidc_discovery_url" or "jwt_validation_pubkeys".',
    fieldGroup: 'default',
    type: 'string',
  },
  jwtSupportedAlgs: {
    editType: 'stringArray',
    helpText: 'A list of supported signing algorithms. Defaults to RS256.',
    fieldGroup: 'default',
  },
  jwtValidationPubkeys: {
    editType: 'stringArray',
    helpText:
      'A list of PEM-encoded public keys to use to authenticate signatures locally. Cannot be used with "jwks_url" or "oidc_discovery_url".',
    fieldGroup: 'default',
  },
  namespaceInState: {
    editType: 'boolean',
    helpText:
      'Pass namespace in the OIDC state parameter instead of as a separate query parameter. With this setting, the allowed redirect URL(s) in Vault and on the provider side should not contain a namespace query parameter. This means only one redirect URL entry needs to be maintained on the provider side for all vault namespaces that will be authenticating against it. Defaults to true for new configs.',
    fieldGroup: 'default',
    defaultValue: true,
    label: 'Namespace in OIDC state',
    type: 'boolean',
  },
  oidcClientId: {
    editType: 'string',
    helpText: 'The OAuth Client ID configured with your OIDC provider.',
    fieldGroup: 'default',
    type: 'string',
  },
  oidcClientSecret: {
    editType: 'string',
    helpText: 'The OAuth Client Secret configured with your OIDC provider.',
    fieldGroup: 'default',
    sensitive: true,
    type: 'string',
  },
  oidcDiscoveryCaPem: {
    editType: 'string',
    helpText:
      'The CA certificate or chain of certificates, in PEM format, to use to validate connections to the OIDC Discovery URL. If not set, system certificates are used.',
    fieldGroup: 'default',
    type: 'string',
  },
  oidcDiscoveryUrl: {
    editType: 'string',
    helpText:
      'OIDC Discovery URL, without any .well-known component (base path). Cannot be used with "jwks_url" or "jwt_validation_pubkeys".',
    fieldGroup: 'default',
    type: 'string',
  },
  oidcResponseMode: {
    editType: 'string',
    helpText:
      "The response mode to be used in the OAuth2 request. Allowed values are 'query' and 'form_post'.",
    fieldGroup: 'default',
    type: 'string',
  },
  oidcResponseTypes: {
    editType: 'stringArray',
    helpText: "The response types to request. Allowed values are 'code' and 'id_token'. Defaults to 'code'.",
    fieldGroup: 'default',
  },
  providerConfig: {
    editType: 'object',
    helpText: 'Provider-specific configuration. Optional.',
    fieldGroup: 'default',
    label: 'Provider Config',
    type: 'object',
  },
};

const k8sConfig = {
  disableLocalCaJwt: {
    editType: 'boolean',
    helpText:
      'Disable defaulting to the local CA cert and service account JWT when running in a Kubernetes pod',
    fieldGroup: 'default',
    label: 'Disable use of local CA and service account JWT',
    type: 'boolean',
  },
  kubernetesCaCert: {
    editType: 'string',
    helpText: 'PEM encoded CA cert for use by the TLS client used to talk with the API.',
    fieldGroup: 'default',
    label: 'Kubernetes CA Certificate',
    type: 'string',
  },
  kubernetesHost: {
    editType: 'string',
    helpText:
      'Host must be a host string, a host:port pair, or a URL to the base of the Kubernetes API server.',
    fieldGroup: 'default',
    type: 'string',
  },
  pemKeys: {
    editType: 'stringArray',
    helpText:
      'Optional list of PEM-formated public keys or certificates used to verify the signatures of kubernetes service account JWTs. If a certificate is given, its public key will be extracted. Not every installation of Kubernetes exposes these keys.',
    fieldGroup: 'default',
    label: 'Service account verification keys',
  },
  tokenReviewerJwt: {
    editType: 'string',
    helpText:
      'A service account JWT (or other token) used as a bearer token to access the TokenReview API to validate other JWTs during login. If not set the JWT used for login will be used to access the API.',
    fieldGroup: 'default',
    label: 'Token Reviewer JWT',
    type: 'string',
  },
};
const k8sRole = {
  name: {
    editType: 'string',
    helpText: 'Name of the role.',
    fieldValue: 'mutableId',
    fieldGroup: 'default',
    readOnly: true,
    label: 'Name',
    type: 'string',
  },
  aliasNameSource: {
    editType: 'string',
    helpText:
      'Source to use when deriving the Alias name. valid choices: "serviceaccount_uid" : <token.uid> e.g. 474b11b5-0f20-4f9d-8ca5-65715ab325e0 (most secure choice) "serviceaccount_name" : <namespace>/<serviceaccount> e.g. vault/vault-agent default: "serviceaccount_uid"',
    fieldGroup: 'default',
    type: 'string',
  },
  audience: {
    editType: 'string',
    helpText: 'Optional Audience claim to verify in the jwt.',
    fieldGroup: 'default',
    type: 'string',
  },
  boundServiceAccountNames: {
    editType: 'stringArray',
    helpText: 'List of service account names able to access this role. If set to "*" all names are allowed.',
    fieldGroup: 'default',
  },
  boundServiceAccountNamespaces: {
    editType: 'stringArray',
    helpText: 'List of namespaces allowed to access this role. If set to "*" all namespaces are allowed.',
    fieldGroup: 'default',
  },
  tokenBoundCidrs: {
    editType: 'stringArray',
    helpText:
      'A list of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.',
    fieldGroup: 'Tokens',
    label: "Generated Token's Bound CIDRs",
  },
  tokenExplicitMaxTtl: {
    editType: 'ttl',
    helpText:
      'If set, tokens created via this role carry an explicit maximum TTL. During renewal, the current maximum TTL values of the role and the mount are not checked for changes, and any updates to these values will have no effect on the token being renewed.',
    fieldGroup: 'Tokens',
    label: "Generated Token's Explicit Maximum TTL",
  },
  tokenMaxTtl: {
    editType: 'ttl',
    helpText: 'The maximum lifetime of the generated token',
    fieldGroup: 'Tokens',
    label: "Generated Token's Maximum TTL",
  },
  tokenNoDefaultPolicy: {
    editType: 'boolean',
    helpText: "If true, the 'default' policy will not automatically be added to generated tokens",
    fieldGroup: 'Tokens',
    label: "Do Not Attach 'default' Policy To Generated Tokens",
    type: 'boolean',
  },
  tokenNumUses: {
    editType: 'number',
    helpText: 'The maximum number of times a token may be used, a value of zero means unlimited',
    fieldGroup: 'Tokens',
    label: 'Maximum Uses of Generated Tokens',
    type: 'number',
  },
  tokenPeriod: {
    editType: 'ttl',
    helpText:
      'If set, tokens created via this role will have no max lifetime; instead, their renewal period will be fixed to this value. This takes an integer number of seconds, or a string duration (e.g. "24h").',
    fieldGroup: 'Tokens',
    label: "Generated Token's Period",
  },
  tokenPolicies: {
    editType: 'stringArray',
    helpText: 'A list of policies that will apply to the generated token for this user.',
    fieldGroup: 'Tokens',
    label: "Generated Token's Policies",
  },
  tokenTtl: {
    editType: 'ttl',
    helpText: 'The initial ttl of the token to generate',
    fieldGroup: 'Tokens',
    label: "Generated Token's Initial TTL",
  },
  tokenType: {
    editType: 'string',
    helpText: 'The type of token to generate, service or batch',
    fieldGroup: 'Tokens',
    label: "Generated Token's Type",
    type: 'string',
  },
};

const ldapConfig = {
  anonymousGroupSearch: {
    editType: 'boolean',
    helpText:
      'Use anonymous binds when performing LDAP group searches (if true the initial credentials will still be used for the initial connection test).',
    fieldGroup: 'default',
    label: 'Anonymous group search',
    type: 'boolean',
  },
  binddn: {
    editType: 'string',
    helpText: 'LDAP DN for searching for the user DN (optional)',
    fieldGroup: 'default',
    label: 'Name of Object to bind (binddn)',
    type: 'string',
  },
  bindpass: {
    editType: 'string',
    helpText: 'LDAP password for searching for the user DN (optional)',
    fieldGroup: 'default',
    sensitive: true,
    type: 'string',
  },
  caseSensitiveNames: {
    editType: 'boolean',
    helpText:
      'If true, case sensitivity will be used when comparing usernames and groups for matching policies.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  certificate: {
    editType: 'file',
    helpText:
      'CA certificate to use when verifying LDAP server certificate, must be x509 PEM encoded (optional)',
    fieldGroup: 'default',
    label: 'CA certificate',
    type: 'string',
  },
  clientTlsCert: {
    editType: 'file',
    helpText: 'Client certificate to provide to the LDAP server, must be x509 PEM encoded (optional)',
    fieldGroup: 'default',
    label: 'Client certificate',
    type: 'string',
  },
  clientTlsKey: {
    editType: 'file',
    helpText: 'Client certificate key to provide to the LDAP server, must be x509 PEM encoded (optional)',
    fieldGroup: 'default',
    label: 'Client key',
    type: 'string',
  },
  connectionTimeout: {
    editType: 'ttl',
    helpText:
      'Timeout, in seconds, when attempting to connect to the LDAP server before trying the next URL in the configuration.',
    fieldGroup: 'default',
  },
  denyNullBind: {
    editType: 'boolean',
    helpText: "Denies an unauthenticated LDAP bind request if the user's password is empty; defaults to true",
    fieldGroup: 'default',
    type: 'boolean',
  },
  dereferenceAliases: {
    editType: 'string',
    helpText:
      "When aliases should be dereferenced on search operations. Accepted values are 'never', 'finding', 'searching', 'always'. Defaults to 'never'.",
    possibleValues: ['never', 'finding', 'searching', 'always'],
    fieldGroup: 'default',
    type: 'string',
  },
  discoverdn: {
    editType: 'boolean',
    helpText: 'Use anonymous bind to discover the bind DN of a user (optional)',
    fieldGroup: 'default',
    label: 'Discover DN',
    type: 'boolean',
  },
  groupattr: {
    editType: 'string',
    helpText:
      'LDAP attribute to follow on objects returned by <groupfilter> in order to enumerate user group membership. Examples: "cn" or "memberOf", etc. Default: cn',
    fieldGroup: 'default',
    defaultValue: 'cn',
    label: 'Group Attribute',
    type: 'string',
  },
  groupdn: {
    editType: 'string',
    helpText: 'LDAP search base to use for group membership search (eg: ou=Groups,dc=example,dc=org)',
    fieldGroup: 'default',
    label: 'Group DN',
    type: 'string',
  },
  groupfilter: {
    editType: 'string',
    helpText:
      'Go template for querying group membership of user (optional) The template can access the following context variables: UserDN, Username Example: (&(objectClass=group)(member:1.2.840.113556.1.4.1941:={{.UserDN}})) Default: (|(memberUid={{.Username}})(member={{.UserDN}})(uniqueMember={{.UserDN}}))',
    fieldGroup: 'default',
    label: 'Group Filter',
    type: 'string',
  },
  insecureTls: {
    editType: 'boolean',
    helpText: 'Skip LDAP server SSL Certificate verification - VERY insecure (optional)',
    fieldGroup: 'default',
    label: 'Insecure TLS',
    type: 'boolean',
  },
  maxPageSize: {
    editType: 'number',
    helpText:
      "If set to a value greater than 0, the LDAP backend will use the LDAP server's paged search control to request pages of up to the given size. This can be used to avoid hitting the LDAP server's maximum result size limit. Otherwise, the LDAP backend will not use the paged search control.",
    fieldGroup: 'default',
    type: 'number',
  },
  requestTimeout: {
    editType: 'ttl',
    helpText:
      'Timeout, in seconds, for the connection when making requests against the server before returning back an error.',
    fieldGroup: 'default',
  },
  starttls: {
    editType: 'boolean',
    helpText: 'Issue a StartTLS command after establishing unencrypted connection (optional)',
    fieldGroup: 'default',
    label: 'Issue StartTLS',
    type: 'boolean',
  },
  tlsMaxVersion: {
    editType: 'string',
    helpText:
      "Maximum TLS version to use. Accepted values are 'tls10', 'tls11', 'tls12' or 'tls13'. Defaults to 'tls12'",
    possibleValues: ['tls10', 'tls11', 'tls12', 'tls13'],
    fieldGroup: 'default',
    label: 'Maximum TLS Version',
    type: 'string',
  },
  tlsMinVersion: {
    editType: 'string',
    helpText:
      "Minimum TLS version to use. Accepted values are 'tls10', 'tls11', 'tls12' or 'tls13'. Defaults to 'tls12'",
    possibleValues: ['tls10', 'tls11', 'tls12', 'tls13'],
    fieldGroup: 'default',
    label: 'Minimum TLS Version',
    type: 'string',
  },
  tokenBoundCidrs: {
    editType: 'stringArray',
    helpText:
      'A list of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.',
    fieldGroup: 'Tokens',
    label: "Generated Token's Bound CIDRs",
  },
  tokenExplicitMaxTtl: {
    editType: 'ttl',
    helpText:
      'If set, tokens created via this role carry an explicit maximum TTL. During renewal, the current maximum TTL values of the role and the mount are not checked for changes, and any updates to these values will have no effect on the token being renewed.',
    fieldGroup: 'Tokens',
    label: "Generated Token's Explicit Maximum TTL",
  },
  tokenMaxTtl: {
    editType: 'ttl',
    helpText: 'The maximum lifetime of the generated token',
    fieldGroup: 'Tokens',
    label: "Generated Token's Maximum TTL",
  },
  tokenNoDefaultPolicy: {
    editType: 'boolean',
    helpText: "If true, the 'default' policy will not automatically be added to generated tokens",
    fieldGroup: 'Tokens',
    label: "Do Not Attach 'default' Policy To Generated Tokens",
    type: 'boolean',
  },
  tokenNumUses: {
    editType: 'number',
    helpText: 'The maximum number of times a token may be used, a value of zero means unlimited',
    fieldGroup: 'Tokens',
    label: 'Maximum Uses of Generated Tokens',
    type: 'number',
  },
  tokenPeriod: {
    editType: 'ttl',
    helpText:
      'If set, tokens created via this role will have no max lifetime; instead, their renewal period will be fixed to this value. This takes an integer number of seconds, or a string duration (e.g. "24h").',
    fieldGroup: 'Tokens',
    label: "Generated Token's Period",
  },
  tokenPolicies: {
    editType: 'stringArray',
    helpText: 'A list of policies that will apply to the generated token for this user.',
    fieldGroup: 'Tokens',
    label: "Generated Token's Policies",
  },
  tokenTtl: {
    editType: 'ttl',
    helpText: 'The initial ttl of the token to generate',
    fieldGroup: 'Tokens',
    label: "Generated Token's Initial TTL",
  },
  tokenType: {
    editType: 'string',
    helpText: 'The type of token to generate, service or batch',
    fieldGroup: 'Tokens',
    label: "Generated Token's Type",
    type: 'string',
  },
  upndomain: {
    editType: 'string',
    helpText: 'Enables userPrincipalDomain login with [username]@UPNDomain (optional)',
    fieldGroup: 'default',
    label: 'User Principal (UPN) Domain',
    type: 'string',
  },
  url: {
    editType: 'string',
    helpText:
      'LDAP URL to connect to (default: ldap://127.0.0.1). Multiple URLs can be specified by concatenating them with commas; they will be tried in-order.',
    fieldGroup: 'default',
    label: 'URL',
    type: 'string',
  },
  usePre111GroupCnBehavior: {
    editType: 'boolean',
    helpText:
      'In Vault 1.1.1 a fix for handling group CN values of different cases unfortunately introduced a regression that could cause previously defined groups to not be found due to a change in the resulting name. If set true, the pre-1.1.1 behavior for matching group CNs will be used. This is only needed in some upgrade scenarios for backwards compatibility. It is enabled by default if the config is upgraded but disabled by default on new configurations.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  useTokenGroups: {
    editType: 'boolean',
    helpText:
      'If true, use the Active Directory tokenGroups constructed attribute of the user to find the group memberships. This will find all security groups including nested ones.',
    fieldGroup: 'default',
    type: 'boolean',
  },
  userattr: {
    editType: 'string',
    helpText: 'Attribute used for users (default: cn)',
    fieldGroup: 'default',
    defaultValue: 'cn',
    label: 'User Attribute',
    type: 'string',
  },
  userdn: {
    editType: 'string',
    helpText: 'LDAP domain to use for users (eg: ou=People,dc=example,dc=org)',
    fieldGroup: 'default',
    label: 'User DN',
    type: 'string',
  },
  userfilter: {
    editType: 'string',
    helpText:
      'Go template for LDAP user search filer (optional) The template can access the following context variables: UserAttr, Username Default: ({{.UserAttr}}={{.Username}})',
    fieldGroup: 'default',
    label: 'User Search Filter',
    type: 'string',
  },
  usernameAsAlias: {
    editType: 'boolean',
    helpText: 'If true, sets the alias name to the username',
    fieldGroup: 'default',
    type: 'boolean',
  },
};

const ldapGroup = {
  name: {
    editType: 'string',
    helpText: 'Name of the LDAP group.',
    fieldValue: 'mutableId',
    fieldGroup: 'default',
    readOnly: true,
    label: 'Name',
    type: 'string',
  },
  policies: {
    editType: 'stringArray',
    helpText: 'A list of policies associated to the group.',
    fieldGroup: 'default',
  },
};

const ldapUser = {
  name: {
    editType: 'string',
    helpText: 'Name of the LDAP user.',
    fieldValue: 'mutableId',
    fieldGroup: 'default',
    readOnly: true,
    label: 'Name',
    type: 'string',
  },
  groups: {
    editType: 'stringArray',
    helpText: 'A list of additional groups associated with the user.',
    fieldGroup: 'default',
  },
  policies: {
    editType: 'stringArray',
    helpText: 'A list of policies associated with the user.',
    fieldGroup: 'default',
  },
};

const oktaConfig = {
  apiToken: {
    editType: 'string',
    helpText: 'Okta API key.',
    fieldGroup: 'default',
    label: 'API Token',
    type: 'string',
  },
  baseUrl: {
    editType: 'string',
    helpText:
      'The base domain to use for the Okta API. When not specified in the configuration, "okta.com" is used.',
    fieldGroup: 'default',
    label: 'Base URL',
    type: 'string',
  },
  bypassOktaMfa: {
    editType: 'boolean',
    helpText:
      'When set true, requests by Okta for a MFA check will be bypassed. This also disallows certain status checks on the account, such as whether the password is expired.',
    fieldGroup: 'default',
    label: 'Bypass Okta MFA',
    type: 'boolean',
  },
  orgName: {
    editType: 'string',
    helpText: 'Name of the organization to be used in the Okta API.',
    fieldGroup: 'default',
    label: 'Organization Name',
    type: 'string',
  },
  tokenBoundCidrs: {
    editType: 'stringArray',
    helpText:
      'A list of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.',
    fieldGroup: 'Tokens',
    label: "Generated Token's Bound CIDRs",
  },
  tokenExplicitMaxTtl: {
    editType: 'ttl',
    helpText:
      'If set, tokens created via this role carry an explicit maximum TTL. During renewal, the current maximum TTL values of the role and the mount are not checked for changes, and any updates to these values will have no effect on the token being renewed.',
    fieldGroup: 'Tokens',
    label: "Generated Token's Explicit Maximum TTL",
  },
  tokenMaxTtl: {
    editType: 'ttl',
    helpText: 'The maximum lifetime of the generated token',
    fieldGroup: 'Tokens',
    label: "Generated Token's Maximum TTL",
  },
  tokenNoDefaultPolicy: {
    editType: 'boolean',
    helpText: "If true, the 'default' policy will not automatically be added to generated tokens",
    fieldGroup: 'Tokens',
    label: "Do Not Attach 'default' Policy To Generated Tokens",
    type: 'boolean',
  },
  tokenNumUses: {
    editType: 'number',
    helpText: 'The maximum number of times a token may be used, a value of zero means unlimited',
    fieldGroup: 'Tokens',
    label: 'Maximum Uses of Generated Tokens',
    type: 'number',
  },
  tokenPeriod: {
    editType: 'ttl',
    helpText:
      'If set, tokens created via this role will have no max lifetime; instead, their renewal period will be fixed to this value. This takes an integer number of seconds, or a string duration (e.g. "24h").',
    fieldGroup: 'Tokens',
    label: "Generated Token's Period",
  },
  tokenPolicies: {
    editType: 'stringArray',
    helpText: 'A list of policies that will apply to the generated token for this user.',
    fieldGroup: 'Tokens',
    label: "Generated Token's Policies",
  },
  tokenTtl: {
    editType: 'ttl',
    helpText: 'The initial ttl of the token to generate',
    fieldGroup: 'Tokens',
    label: "Generated Token's Initial TTL",
  },
  tokenType: {
    editType: 'string',
    helpText: 'The type of token to generate, service or batch',
    fieldGroup: 'Tokens',
    label: "Generated Token's Type",
    type: 'string',
  },
};

const oktaGroup = {
  name: {
    editType: 'string',
    helpText: 'Name of the Okta group.',
    fieldValue: 'mutableId',
    fieldGroup: 'default',
    readOnly: true,
    label: 'Name',
    type: 'string',
  },
  policies: {
    editType: 'stringArray',
    helpText: 'A list of policies associated to the group.',
    fieldGroup: 'default',
  },
};

const oktaUser = {
  name: {
    editType: 'string',
    helpText: 'Name of the user.',
    fieldValue: 'mutableId',
    fieldGroup: 'default',
    readOnly: true,
    label: 'Name',
    type: 'string',
  },
  groups: {
    editType: 'stringArray',
    helpText: 'List of groups associated with the user.',
    fieldGroup: 'default',
  },
  policies: {
    editType: 'stringArray',
    helpText: 'List of policies associated with the user.',
    fieldGroup: 'default',
  },
};

const radiusConfig = {
  dialTimeout: {
    editType: 'ttl',
    helpText: 'Number of seconds before connect times out (default: 10)',
    fieldGroup: 'default',
    defaultValue: 10,
  },
  host: {
    editType: 'string',
    helpText: 'RADIUS server host',
    fieldGroup: 'default',
    label: 'Host',
    type: 'string',
  },
  nasIdentifier: {
    editType: 'string',
    helpText: 'RADIUS NAS Identifier field (optional)',
    fieldGroup: 'default',
    label: 'NAS Identifier',
    type: 'string',
  },
  nasPort: {
    editType: 'number',
    helpText: 'RADIUS NAS port field (default: 10)',
    fieldGroup: 'default',
    defaultValue: 10,
    label: 'NAS Port',
    type: 'number',
  },
  port: {
    editType: 'number',
    helpText: 'RADIUS server port (default: 1812)',
    fieldGroup: 'default',
    defaultValue: 1812,
    type: 'number',
  },
  readTimeout: {
    editType: 'ttl',
    helpText: 'Number of seconds before response times out (default: 10)',
    fieldGroup: 'default',
    defaultValue: 10,
  },
  secret: {
    editType: 'string',
    helpText: 'Secret shared with the RADIUS server',
    fieldGroup: 'default',
    type: 'string',
  },
  tokenBoundCidrs: {
    editType: 'stringArray',
    helpText:
      'A list of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.',
    fieldGroup: 'Tokens',
    label: "Generated Token's Bound CIDRs",
  },
  tokenExplicitMaxTtl: {
    editType: 'ttl',
    helpText:
      'If set, tokens created via this role carry an explicit maximum TTL. During renewal, the current maximum TTL values of the role and the mount are not checked for changes, and any updates to these values will have no effect on the token being renewed.',
    fieldGroup: 'Tokens',
    label: "Generated Token's Explicit Maximum TTL",
  },
  tokenMaxTtl: {
    editType: 'ttl',
    helpText: 'The maximum lifetime of the generated token',
    fieldGroup: 'Tokens',
    label: "Generated Token's Maximum TTL",
  },
  tokenNoDefaultPolicy: {
    editType: 'boolean',
    helpText: "If true, the 'default' policy will not automatically be added to generated tokens",
    fieldGroup: 'Tokens',
    label: "Do Not Attach 'default' Policy To Generated Tokens",
    type: 'boolean',
  },
  tokenNumUses: {
    editType: 'number',
    helpText: 'The maximum number of times a token may be used, a value of zero means unlimited',
    fieldGroup: 'Tokens',
    label: 'Maximum Uses of Generated Tokens',
    type: 'number',
  },
  tokenPeriod: {
    editType: 'ttl',
    helpText:
      'If set, tokens created via this role will have no max lifetime; instead, their renewal period will be fixed to this value. This takes an integer number of seconds, or a string duration (e.g. "24h").',
    fieldGroup: 'Tokens',
    label: "Generated Token's Period",
  },
  tokenPolicies: {
    editType: 'stringArray',
    helpText: 'A list of policies that will apply to the generated token for this user.',
    fieldGroup: 'Tokens',
    label: "Generated Token's Policies",
  },
  tokenTtl: {
    editType: 'ttl',
    helpText: 'The initial ttl of the token to generate',
    fieldGroup: 'Tokens',
    label: "Generated Token's Initial TTL",
  },
  tokenType: {
    editType: 'string',
    helpText: 'The type of token to generate, service or batch',
    fieldGroup: 'Tokens',
    label: "Generated Token's Type",
    type: 'string',
  },
  unregisteredUserPolicies: {
    editType: 'string',
    helpText:
      'List of policies to grant upon successful RADIUS authentication of an unregistered user (default: empty)',
    fieldGroup: 'default',
    label: 'Policies for unregistered users',
    type: 'string',
  },
};

const radiusUser = {
  name: {
    editType: 'string',
    helpText: 'Name of the RADIUS user.',
    fieldValue: 'mutableId',
    fieldGroup: 'default',
    readOnly: true,
    label: 'Name',
    type: 'string',
  },
  policies: {
    editType: 'stringArray',
    helpText: 'A list of policies associated to the user.',
    fieldGroup: 'default',
  },
};

const awsConfig = {
  certName: {
    editType: 'string',
    helpText: 'Name of the certificate.',
    fieldValue: 'mutableId',
    fieldGroup: 'default',
    readOnly: true,
    label: 'Cert name',
    type: 'string',
  },
  awsPublicCert: {
    editType: 'string',
    helpText:
      'Base64 encoded AWS Public cert required to verify PKCS7 signature of the EC2 instance metadata.',
    fieldGroup: 'default',
    type: 'string',
  },
  type: {
    editType: 'string',
    helpText:
      'Takes the value of either "pkcs7" or "identity", indicating the type of document which can be verified using the given certificate. The reason is that the PKCS#7 document will have a DSA digest and the identity signature will have an RSA signature, and accordingly the public certificates to verify those also vary. Defaults to "pkcs7".',
    fieldGroup: 'default',
    type: 'string',
  },
};

export default {
  // auth
  azureConfig,
  userpassUser,
  certConfig,
  certCert,
  gcpConfig,
  githubConfig,
  jwtConfig,
  k8sConfig,
  k8sRole,
  ldapConfig,
  ldapGroup,
  ldapUser,
  oktaConfig,
  oktaGroup,
  oktaUser,
  radiusConfig,
  radiusUser,
  awsConfig,
  // secret engines
  kmipConfig,
  kmipRole,
  pkiAcme,
  pkiCertGenerate,
  pkiCertSign,
  pkiCluster,
  pkiRole,
  pkiSignCsr,
  pkiTidy,
  pkiUrls,
  sshRole,
};
