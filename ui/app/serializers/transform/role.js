/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationSerializer from '../application';

export default ApplicationSerializer.extend({
  extractLazyPaginatedData(payload) {
    return payload.data.keys.map((key) => {
      const model = {
        id: key,
      };
      if (payload.backend) {
        model.backend = payload.backend;
      }
      return model;
    });
  },
});
