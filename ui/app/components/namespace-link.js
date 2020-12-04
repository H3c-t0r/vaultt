import { inject as service } from '@ember/service';
import { alias } from '@ember/object/computed';
import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  namespaceService: service('namespace'),
  currentNamespace: alias('namespaceService.path'),

  tagName: '',
  //public api
  targetNamespace: null,
  showLastSegment: false,

  normalizedNamespace: computed('targetNamespace', function() {
    let ns = this.targetNamespace;
    return (ns || '').replace(/\.+/g, '/').replace(/☃/g, '.');
  }),

  namespaceDisplay: computed('normalizedNamespace', 'showLastSegment', function() {
    let ns = this.normalizedNamespace;
    let showLastSegment = this.showLastSegment;
    let parts = ns.split('/');
    if (ns === '') {
      return 'root';
    }
    return showLastSegment ? parts[parts.length - 1] : ns;
  }),

  isCurrentNamespace: computed('targetNamespace', 'currentNamespace', function() {
    return this.currentNamespace === this.targetNamespace;
  }),
});
