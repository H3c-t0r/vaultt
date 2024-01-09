/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { decodeString, encodeString } from 'core/utils/b64';
import ApplicationSerializer from '../application';

export default class MessageSerializer extends ApplicationSerializer {
  attrs = {
    link: { serialize: false },
    active: { serialize: false },
    start_time: { serialize: false },
    end_time: { serialize: false },
  };

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (requestType === 'query' && !payload.meta) {
      const transformed = this.mapPayload(payload);
      return super.normalizeResponse(store, primaryModelClass, transformed, id, requestType);
    }
    if (requestType === 'queryRecord') {
      const transformed = {
        ...payload.data,
        message: decodeString(payload.data.message),
        link_title: payload.data.link.title,
        link_href: payload.data.link.href,
      };
      delete transformed.link;
      return super.normalizeResponse(store, primaryModelClass, transformed, id, requestType);
    }
    return super.normalizeResponse(store, primaryModelClass, payload, id, requestType);
  }

  serialize() {
    const json = super.serialize(...arguments);
    json.message = encodeString(json.message);
    json.link = {
      title: json?.link_title || '',
      href: json?.link_href || '',
    };
    delete json?.link_title;
    delete json?.link_href;
    return json;
  }

  mapPayload(payload) {
    if (payload.data) {
      if (payload.data?.keys && Array.isArray(payload.data.keys)) {
        return payload.data.keys.map((key) => {
          const data = {
            id: key,
            linkTitle: payload.data.key_info.link?.title,
            linkHref: payload.data.key_info.link?.href,
            ...payload.data.key_info[key],
          };
          if (data.message) data.message = decodeString(data.message);
          return data;
        });
      }
      Object.assign(payload, payload.data);
      delete payload.data;
    }
    return payload;
  }

  extractLazyPaginatedData(payload) {
    return this.mapPayload(payload);
  }
}
