import Route from '@ember/routing/route';
import RSVP from 'rsvp';
import { action } from '@ember/object';

export default class ClientsRoute extends Route {
  // async getMonthly() {
  //   try {
  //     return await this.store.queryRecord('clients/monthly', {});
  //   } catch (e) {
  //     // ARG TODO handle
  //     return e;
  //   }
  // }

  // async getConfig() {
  //   try {
  //     return await this.store.queryRecord('clients/config', {});
  //   } catch (e) {
  //     // ARG TODO handle
  //     return e;
  //   }
  // }

  model() {
    return RSVP.hash({
      config: this.store.queryRecord('clients/config', {}),
      monthly: this.store.queryRecord('clients/monthly', {}),
    });
  }

  @action
  async loading(transition) {
    // eslint-disable-next-line ember/no-controller-access-in-routes
    let controller = this.controllerFor('vault.cluster.clients.index');
    controller.set('currentlyLoading', true);
    transition.promise.finally(function () {
      controller.set('currentlyLoading', false);
    });
  }
}
