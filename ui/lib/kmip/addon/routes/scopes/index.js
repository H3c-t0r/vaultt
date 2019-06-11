import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ListRoute from 'core/mixins/list-route';

export default Route.extend(ListRoute, {
  store: service(),
  model(params) {
    let model = [{ id: 'a scope' }];
    model.set('meta', { total: 1 });
    return model;
    return this.store
      .lazyPaginatedQuery('kmip/scope', {
        responsePath: 'data.keys',
        page: params.page,
        pageFilter: params.pageFilter,
      })
      .catch(err => {
        if (err.httpStatus === 404) {
          return [];
        } else {
          throw err;
        }
      });
  },

  actions: {
    willTransition(transition) {
      window.scrollTo(0, 0);
      if (transition.targetName !== this.routeName) {
        this.store.clearAllDatasets();
      }
      return true;
    },
    reload() {
      this.store.clearAllDatasets();
      this.refresh();
    },
  },
});
