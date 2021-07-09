import { computed } from '@ember/object';
import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import ModelBoundaryRoute from 'vault/mixins/model-boundary-route';

export default Route.extend(ModelBoundaryRoute, {
  auth: service(),
  controlGroup: service(),
  flashMessages: service(),
  console: service(),
  permissions: service(),
  namespaceService: service('namespace'),

  modelTypes: computed(function() {
    return ['secret', 'secret-engine'];
  }),

  async beforeModel() {
    const authType = this.auth.getAuthType();
    const baseUrl = window.location.origin;
    this.auth.deleteCurrentToken();
    this.controlGroup.deleteTokens();
    this.namespaceService.reset();
    this.console.set('isOpen', false);
    this.console.clearLog(true);
    this.flashMessages.clearMessages();
    this.permissions.reset();
    location.assign(`${baseUrl}/ui/vault/auth?with=${authType}`);
  },
});
