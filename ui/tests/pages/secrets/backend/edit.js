/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create, visitable } from 'ember-cli-page-object';

export default create({
  visit: visitable('/vault/secrets/:backend/edit/:id'),
  visitRoot: visitable('/vault/secrets/:backend/edit'),
});
