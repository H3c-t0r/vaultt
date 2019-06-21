import ApplicationAdapater from '../application';

export default ApplicationAdapater.extend({
  namespace: 'v1',
  pathForType(type) {
    return type.replace('kmip/', '');
  },

  urlForQuery() {
    return this._super(...arguments) + '?list=true';
  },

  query(store, type) {
    return this.ajax(this.buildURL(type.modelName, null, null, 'query'), 'GET');
  },

  buildURL(modelName, id, snapshot, requestType, query) {
    if (requestType === 'createRecord') {
      return this._super(...arguments);
    }
    return this._super(`${modelName}`, id, snapshot, requestType, query);
  },
});
