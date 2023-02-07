/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { assign } from '@ember/polyfills';
import ApplicationSerializer from './application';
import { EmbeddedRecordsMixin } from '@ember-data/serializer/rest';

export default ApplicationSerializer.extend(EmbeddedRecordsMixin, {
  attrs: {
    config: { embedded: 'always' },
  },

  normalize(modelClass, data) {
    // embedded records need a unique value to be stored
    // set id for config to uuid of secret engine
    if (data.config && !data.config.id) {
      data.config.id = data.uuid;
    }
    // move version out of options so it can be defined on secret-engine model
    data.version = data.options ? data.options.version : null;
    return this._super(modelClass, data);
  },

  normalizeBackend(path, backend) {
    let struct = {};
    for (const attribute in backend) {
      struct[attribute] = backend[attribute];
    }
    //queryRecord adds path to the response
    if (path !== null && !struct.path) {
      struct.path = path;
    }

    if (struct.data) {
      struct = assign({}, struct, struct.data);
      delete struct.data;
    }
    // strip the trailing slash off of the path so we
    // can navigate to it without getting `//` in the url
    struct.id = struct.path.slice(0, -1);
    return struct;
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const isCreate = requestType === 'createRecord';
    const isFind = requestType === 'findRecord';
    const isQueryRecord = requestType === 'queryRecord';
    let backends;
    if (isCreate) {
      backends = payload.data;
    } else if (isFind) {
      backends = this.normalizeBackend(id + '/', payload.data);
    } else if (isQueryRecord) {
      backends = this.normalizeBackend(null, payload);
    } else {
      // this is terrible, I'm sorry
      // TODO extract AWS and SSH config saving from the secret-engine model to simplify this
      if (payload.data.secret) {
        backends = Object.keys(payload.data.secret).map((id) =>
          this.normalizeBackend(id, payload.data.secret[id])
        );
      } else if (!payload.data.path) {
        backends = Object.keys(payload.data).map((id) => this.normalizeBackend(id, payload[id]));
      } else {
        backends = [this.normalizeBackend(payload.data.path, payload.data)];
      }
    }

    return this._super(store, primaryModelClass, backends, id, requestType);
  },

  serialize(snapshot) {
    const type = snapshot.record.get('engineType');
    const data = this._super(...arguments);
    // move version back to options
    data.options = data.version ? { version: data.version } : {};
    delete data.version;

    if (type !== 'kv' || data.options.version === 1) {
      // These items are on the model, but used by the kv-v2 config endpoint only
      delete data.max_versions;
      delete data.cas_required;
      delete data.delete_version_after;
    }
    // only KV uses options
    if (type !== 'kv' && type !== 'generic') {
      delete data.options;
    } else if (!data.options.version) {
      // if options.version isn't set for some reason
      // default to 2
      data.options.version = 2;
    }
    return data;
  },
});
