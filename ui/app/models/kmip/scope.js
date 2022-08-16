import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import apiPath from 'vault/utils/api-path';
import attachCapabilities from 'vault/lib/attach-capabilities';

import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

let ModelExport = Model.extend({
  name: attr('string'),
  backend: attr({ readOnly: true }),
  attrs: computed(function () {
    return expandAttributeMeta(this, ['name']);
  }),
});

export default attachCapabilities(ModelExport, {
  updatePath: apiPath`${'backend'}/scope/${'id'}`,
});
