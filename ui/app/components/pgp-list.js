import Ember from 'ember';

export default Ember.Component.extend({
  onDataUpdate: () => {},
  listData: Ember.computed('listLength', function() {
    let num = this.get('listLength');
    if (num) {
      num = parseInt(num, 10);
    }
    return Array(num || 0).fill(null).map(() => ({ value: '' }));
  }),
  listLength: 0,
  actions: {
    setKey(index, key) {
      let listData = this.get('listData');
      listData.replace(index, 1, key);
      this.get('onDataUpdate')(listData.compact().map(k => k.value));
    },
  },
});
