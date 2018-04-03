import Ember from 'ember';

const { get, inject } = Ember;

export default Ember.Component.extend({
  classNames: 'config-pki',
  flashMessages: inject.service(),

  /*
   *
   * @param String
   * @public
   * String corresponding to the route parameter for the current section
   *
   */

  section: null,

  /*
   * @param DS.Model
   * @public
   *
   * a `pki-config` model - passed in in the component useage
   *
   */
  config: null,

  /*
   * @param Function
   * @public
   *
   * function that gets called to refresh the config model
   *
   */
  onRefresh: () => {},

  loading: false,

  actions: {
    save(section) {
      this.set('loading', true);
      const config = this.get('config');
      config
        .save({
          adapterOptions: {
            method: section,
            fields: get(config, `${section}Attrs`).map(attr => attr.name),
          },
        })
        .then(() => {
          this.get('flashMessages').success(`The ${section} config for this backend has been updated.`);
          // attrs aren't persistent for Tidy
          if (section === 'tidy') {
            config.rollbackAttributes();
          }
          this.send('refresh');
        })
        .finally(() => {
          this.set('loading', false);
        });
    },
    refresh() {
      this.get('onRefresh')();
    },
  },
});
