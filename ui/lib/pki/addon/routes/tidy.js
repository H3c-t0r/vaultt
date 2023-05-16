import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfig } from 'pki/decorators/check-config';
import { hash } from 'rsvp';

@withConfig()
export default class PkiTidyRoute extends Route {
  @service store;

  model() {
    const engine = this.modelFor('application');
    return hash({
      hasConfig: this.shouldPromptConfig,
      engine,
      autoTidyConfig: this.store.queryRecord('pki/tidy', { backend: engine.id, tidyType: 'auto' }),
    });
  }
}
