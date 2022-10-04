import Route from '@ember/routing/route';
import UnloadModelRoute from 'vault/mixins/unload-model-route';
import UnsavedModelRoute from 'vault/mixins/unsaved-model-route';
import { inject as service } from '@ember/service';

export default Route.extend(UnloadModelRoute, UnsavedModelRoute, {
  store: service(),

  model() {
    let itemType = this.modelFor('vault.cluster.access.identity');
    let modelType = `identity/${itemType}`;
    return this.store.createRecord(modelType);
  },
});
