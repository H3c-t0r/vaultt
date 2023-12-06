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
  accessToken: [{ type: 'presence', message: 'Access token is required.' }],
  repositoryOwner: [{ type: 'presence', message: 'Repository owner is required.' }],
  repositoryName: [{ type: 'presence', message: 'Repository name is required.' }],
};
const displayFields = ['name', 'repositoryOwner', 'repositoryName', 'accessToken'];
const formFieldGroups = [
  { default: ['name', 'repositoryOwner', 'repositoryName'] },
  { Credentials: ['accessToken'] },
];
@withModelValidations(validations)
@withFormFields(displayFields, formFieldGroups)
export default class SyncDestinationsGithubModel extends SyncDestinationModel {
  @attr('string', { subText: 'Personal access token to authenticate to the GitHub repository.' })
  accessToken; // obfuscated, never returned by API

  @attr('string', {
    subText: 'Github organization or username that owns the repository.',
    editDisabled: true,
  })
  repositoryOwner;

  @attr('string', {
    subText: 'The name of the Github repository to connect to.',
    editDisabled: true,
  })
  repositoryName;
}
