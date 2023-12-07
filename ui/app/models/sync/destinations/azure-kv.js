/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import SyncDestinationModel from '../destination';
import { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import { withModelValidations } from 'vault/decorators/model-validations';

const validations = {
  name: [{ type: 'presence', message: 'Name is required.' }],
  keyVaultUri: [{ type: 'presence', message: 'Key Vault URI is required.' }],
  clientId: [{ type: 'presence', message: 'Client ID is required.' }],
  clientSecret: [{ type: 'presence', message: 'Client secret is required.' }],
  tenantId: [{ type: 'presence', message: 'Tenant ID is required.' }],
};
const displayFields = ['name', 'keyVaultUri', 'tenantId', 'cloud', 'clientId', 'clientSecret'];
const formFieldGroups = [
  { default: ['name', 'keyVaultUri', 'tenantId', 'cloud', 'clientId'] },
  { Credentials: ['clientSecret'] },
];
@withModelValidations(validations)
@withFormFields(displayFields, formFieldGroups)
export default class SyncDestinationsAzureKeyVaultModel extends SyncDestinationModel {
  @attr('string', {
    label: 'Key Vault URI',
    subText:
      'URI of an existing Azure Key Vault instance. If empty, Vault will use the KEY_VAULT_URI environment variable if configured.',
    editDisabled: true,
  })
  keyVaultUri;

  @attr('string', {
    label: 'Client ID',
    subText:
      'Client ID of an Azure app registration. If empty, Vault will use the AZURE_CLIENT_ID environment variable if configured.',
  })
  clientId;

  @attr('string', {
    subText:
      'Client secret of an Azure app registration. If empty, Vault will use the AZURE_CLIENT_SECRET environment variable if configured.',
  })
  clientSecret; // obfuscated, never returned by API

  @attr('string', {
    label: 'Tenant ID',
    subText:
      'ID of the target Azure tenant. If empty, Vault will use the AZURE_TENANT_ID environment variable if configured.',
    editDisabled: true,
  })
  tenantId;

  @attr('string', {
    subText: 'Specifies a cloud for the client. The default is Azure Public Cloud.',
    defaultValue: 'cloud',
    editDisabled: true,
  })
  cloud;
}
