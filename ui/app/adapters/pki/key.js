import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class PkiKeyAdapter extends ApplicationAdapter {
  namespace = 'v1';

  _urlForQuery(backend, id) {
    const url = `${this.buildURL()}/${encodePath(backend)}`;
    if (id) {
      return url + '/key/' + encodePath(id);
    }
    return url + '/keys';
  }

  query(store, type, query) {
    const { backend } = query;
    return this.ajax(this._urlForQuery(backend), 'GET', { data: { list: true } });
  }

  queryRecord(store, type, query) {
    const { backend, id } = query;
    return this.ajax(this._urlForQuery(backend, id), 'GET');
  }
}
