import { assign } from '@ember/polyfills';
import { decamelize } from '@ember/string';
import DS from 'ember-data';

export default DS.RESTSerializer.extend({
  primaryKey: 'name',

  keyForAttribute: function(attr) {
    return decamelize(attr);
  },

  normalizeSecrets(payload) {
    if (payload.data.keys && Array.isArray(payload.data.keys)) {
      const secrets = payload.data.keys.map(secret => ({ name: secret }));
      return secrets;
    }
    assign(payload, payload.data);
    delete payload.data;
    return [payload];
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const nullResponses = ['updateRecord', 'createRecord', 'deleteRecord'];
    const secrets = nullResponses.includes(requestType) ? { name: id } : this.normalizeSecrets(payload);
    const { modelName } = primaryModelClass;
    let transformedPayload = { [modelName]: secrets };
    // just return the single object because ember is picky
    if (requestType === 'queryRecord') {
      let model = secrets[0];
      // timestamps for these two are in seconds...
      if (model.type === 'aes256-gcm96' || model.type === 'chacha20-poly1305') {
        for (let version in model.keys) {
          model.keys[version] = model.keys[version] * 1000;
        }
      }
      transformedPayload = { [modelName]: model };
    }

    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },

  serialize(snapshot, requestType) {
    if (requestType === 'update') {
      const min_decryption_version = snapshot.attr('minDecryptionVersion');
      const min_encryption_version = snapshot.attr('minEncryptionVersion');
      const deletion_allowed = snapshot.attr('deletionAllowed');
      return {
        min_decryption_version,
        min_encryption_version,
        deletion_allowed,
      };
    } else {
      return this._super(...arguments);
    }
  },
});
