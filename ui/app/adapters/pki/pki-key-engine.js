import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class PkiKeyEngineAdapter extends ApplicationAdapter {
  namespace = 'v1';

  optionsForQuery(id) {
    let data = {};
    if (!id) {
      data['list'] = true;
    }
    return { data };
  }

  urlForQuery(backend, id) {
    let url = `${this.buildURL()}/${encodePath(backend)}/keys`;
    if (id) {
      url = url + '/' + encodePath(id);
    }
    return url;
  }

  async query(store, type, query) {
    const { backend, id } = query;
    let response = await this.ajax(this.urlForQuery(backend, id), 'GET', this.optionsForQuery(id));
    return response;
  }
}
