import ApplicationAdapater from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default ApplicationAdapater.extend({
  namespace: 'v1',

  pathForType() {
    return 'template';
  },

  _url(backend, id) {
    let type = this.pathForType();
    let base = `${this.buildURL()}/${encodePath(backend)}/${type}`;
    if (id) {
      return `${base}/${encodePath(id)}`;
    }
    return base + '?list=true';
  },

  query(store, type, query) {
    return this.ajax(this._url(query.backend), 'GET').then(result => {
      return result;
    });
  },
});
