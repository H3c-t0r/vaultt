import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

export default class KvSecretEditRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    // TODO add model. Needs to include a queryParam for version.
    const backend = this.secretMountPath.get();
    const { name } = this.paramsFor('secrets.secret');
    return hash({
      path: name,
      backend,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'secrets' },
      { label: resolvedModel.path, route: 'secrets.secret.details', model: resolvedModel.path },
      { label: 'edit' },
    ];
  }
}
