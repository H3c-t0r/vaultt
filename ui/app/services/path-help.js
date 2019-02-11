// Low level service that allows users to input paths to make requests to vault
// this service provides the UI synecdote to the cli commands read, write, delete, and list
import Service from '@ember/service';

import { getOwner } from '@ember/application';
import { expandOpenApiProps, combineAttributes } from 'vault/utils/openapi-to-attrs';
import { resolve } from 'rsvp';

export function sanitizePath(path) {
  //remove whitespace + remove trailing and leading slashes
  return path.trim().replace(/^\/+|\/+$/g, '');
}

export default Service.extend({
  attrs: null,
  ajax(url, options = {}) {
    let appAdapter = getOwner(this).lookup(`adapter:application`);
    let { data } = options;
    return appAdapter.ajax(url, 'GET', {
      data,
    });
  },

  getNewModel(modelType, backend, owner) {
    let name = `model:${modelType}`;
    let newModel = owner.factoryFor(name).class;
    if (newModel.merged || !newModel.prototype.useOpenAPI === true) {
      return resolve();
    }

    return this.getProps(modelType, backend).then(props => {
      if (owner.hasRegistration(name) && !newModel.merged) {
        let { attrs, newFields } = combineAttributes(newModel.attributes, props);
        newModel = newModel.extend(attrs, { newFields });
      } else {
        //generate a whole new model
      }

      newModel.reopenClass({ merged: true });
      owner.unregister(name);
      owner.register(name, newModel);
    });
  },

  getProps(modelType, backend) {
    let adapter = getOwner(this).lookup(`adapter:${modelType}`);
    let path = adapter.pathForType();
    const authMethods = [
      'auth-config/ldap',
      'auth-config/github',
      'auth-config/okta',
      'auth-config/radius',
      'auth-config/cert',
      'auth-config/gcp',
      'auth-config/azure',
      'auth-config/kubernetes',
    ];
    let helpUrl = authMethods.includes(modelType)
      ? `/v1/auth/${backend}/${path}?help=1`
      : `/v1/${backend}/${path}/example?help=1`;
    let wildcard;
    switch (path) {
      case 'roles':
        if (modelType === 'role-ssh') {
          wildcard = 'role';
        } else {
          wildcard = 'name';
        }
        break;
      case 'mounts':
        if (modelType === 'secret') {
          wildcard = 'path';
        } else {
          wildcard = 'config';
        }
        break;
      case 'sign':
      case 'issue':
        wildcard = 'role';
        break;
    }

    return this.ajax(helpUrl, backend).then(help => {
      let fullPath = wildcard ? `/${path}/{${wildcard}}` : `/${path}`;
      let props = help.openapi.paths[fullPath].post.requestBody.content['application/json'].schema.properties;
      return expandOpenApiProps(props);
    });
  },
});
