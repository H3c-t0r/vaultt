import Service from '@ember/service';

export default Service.extend({
  mode: null,

  getMode() {
    this.mode;
  },

  setMode(mode) {
    this.set('mode', mode);
  },
});
