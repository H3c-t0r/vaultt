import Ember from 'ember';
import ApplicationAdapter from './application';
import DS from 'ember-data';

const { AdapterError } = DS;
const { assert, inject } = Ember;

const ENDPOINTS = ['health', 'seal-status', 'tokens', 'token', 'seal', 'unseal', 'init', 'capabilities-self'];

const REPLICATION_ENDPOINTS = {
  reindex: 'reindex',
  recover: 'recover',
  status: 'status',

  primary: ['enable', 'disable', 'demote', 'secondary-token', 'revoke-secondary'],

  secondary: ['enable', 'disable', 'promote', 'update-primary'],
};

const REPLICATION_MODES = ['dr', 'performance'];
export default ApplicationAdapter.extend({
  version: inject.service(),
  shouldBackgroundReloadRecord() {
    return true;
  },
  findRecord(store, type, id, snapshot) {
    let fetches = {
      health: this.health(),
      sealStatus: this.sealStatus().catch(e => e),
    };
    if (this.get('version.isEnterprise')) {
      fetches.replicationStatus = this.replicationStatus().catch(e => e);
    }
    return Ember.RSVP.hash(fetches).then(({ health, sealStatus, replicationStatus }) => {
      let ret = {
        id,
        name: snapshot.attr('name'),
      };
      ret = Ember.assign(ret, health);
      if (sealStatus instanceof AdapterError === false) {
        ret = Ember.assign(ret, { nodes: [sealStatus] });
      }
      if (replicationStatus && replicationStatus instanceof AdapterError === false) {
        ret = Ember.assign(ret, replicationStatus.data);
      }
      return Ember.RSVP.resolve(ret);
    });
  },

  pathForType(type) {
    return type === 'cluster' ? 'clusters' : Ember.String.pluralize(type);
  },

  health() {
    return this.ajax(this.urlFor('health'), 'GET', {
      data: { standbycode: 200, sealedcode: 200, uninitcode: 200, drsecondarycode: 200 },
      unauthenticated: true,
    });
  },

  features() {
    return this.ajax(`${this.buildURL()}/license/features`, 'GET', {
      unauthenticated: true,
    });
  },

  sealStatus() {
    return this.ajax(this.urlFor('seal-status'), 'GET', { unauthenticated: true });
  },

  seal() {
    return this.ajax(this.urlFor('seal'), 'PUT');
  },

  unseal(data) {
    return this.ajax(this.urlFor('unseal'), 'PUT', {
      data,
      unauthenticated: true,
    });
  },

  initCluster(data) {
    return this.ajax(this.urlFor('init'), 'PUT', {
      data,
      unauthenticated: true,
    });
  },

  authenticate({ backend, data }) {
    const { token, password, username, path } = data;
    const url = this.urlForAuth(backend, username, path);
    const verb = backend === 'token' ? 'GET' : 'POST';
    let options = {
      unauthenticated: true,
    };
    if (backend === 'token') {
      options.headers = {
        'X-Vault-Token': token,
      };
    } else {
      options.data = token ? { token, password } : { password };
    }

    return this.ajax(url, verb, options);
  },

  urlFor(endpoint) {
    if (!ENDPOINTS.includes(endpoint)) {
      throw new Error(
        `Calls to a ${endpoint} endpoint are not currently allowed in the vault cluster adapater`
      );
    }
    return `${this.buildURL()}/${endpoint}`;
  },

  urlForAuth(type, username, path) {
    const authBackend = type.toLowerCase();
    const authURLs = {
      github: 'login',
      userpass: `login/${encodeURIComponent(username)}`,
      ldap: `login/${encodeURIComponent(username)}`,
      okta: `login/${encodeURIComponent(username)}`,
      token: 'lookup-self',
    };
    const urlSuffix = authURLs[authBackend];
    const urlPrefix = path && authBackend !== 'token' ? path : authBackend;
    if (!urlSuffix) {
      throw new Error(`There is no auth url for ${type}.`);
    }
    return `/v1/auth/${urlPrefix}/${urlSuffix}`;
  },

  urlForReplication(replicationMode, clusterMode, endpoint) {
    let suffix;
    const errString = `Calls to replication ${endpoint} endpoint are not currently allowed in the vault cluster adapater`;
    if (clusterMode) {
      assert(errString, REPLICATION_ENDPOINTS[clusterMode].includes(endpoint));
      suffix = `${replicationMode}/${clusterMode}/${endpoint}`;
    } else {
      assert(errString, REPLICATION_ENDPOINTS[endpoint]);
      suffix = `${endpoint}`;
    }
    return `${this.buildURL()}/replication/${suffix}`;
  },

  replicationStatus() {
    return this.ajax(`${this.buildURL()}/replication/status`, 'GET', { unauthenticated: true });
  },

  replicationDrPromote(data, options) {
    const verb = options && options.checkStatus ? 'GET' : 'PUT';
    return this.ajax(`${this.buildURL()}/replication/dr/secondary/promote`, verb, {
      data,
      unauthenticated: true,
    });
  },

  generateDrOperationToken(data, options) {
    const verb = options && options.checkStatus ? 'GET' : 'PUT';
    let url = `${this.buildURL()}/replication/dr/secondary/generate-operation-token/`;
    if (!data || data.pgp_key || data.otp) {
      // start the generation
      url = url + 'attempt';
    } else {
      // progress the operation
      url = url + 'update';
    }
    return this.ajax(url, verb, {
      data,
      unauthenticated: true,
    });
  },

  replicationAction(action, replicationMode, clusterMode, data) {
    assert(
      `${replicationMode} is an unsupported replication mode.`,
      replicationMode && REPLICATION_MODES.includes(replicationMode)
    );

    const url =
      action === 'recover' || action === 'reindex'
        ? this.urlForReplication(replicationMode, null, action)
        : this.urlForReplication(replicationMode, clusterMode, action);

    return this.ajax(url, 'POST', { data });
  },
});
