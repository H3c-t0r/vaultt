import { parseCertificate } from 'vault/helpers/parse-pki-cert';
import ApplicationSerializer from '../application';

export default class PkiIssuerSerializer extends ApplicationSerializer {
  primaryKey = 'issuer_id';

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (payload.data.certificate) {
      // Parse certificate back from the API and add to payload
      const parsedCert = parseCertificate(payload.data.certificate);
      const data = { issuer_ref: payload.issuer_id, ...payload.data, ...parsedCert };
      const json = super.normalizeResponse(store, primaryModelClass, { ...payload, data }, id, requestType);
      return json;
    }
    return super.normalizeResponse(...arguments);
  }
}
