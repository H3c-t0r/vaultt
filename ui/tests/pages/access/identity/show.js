import { create, collection, contains, text, visitable } from 'ember-cli-page-object';
import flashMessage from 'vault/tests/pages/components/flash-message';
import infoTableRow from 'vault/tests/pages/components/info-table-row';

export default create({
  visit: visitable('/vault/access/identity/:item_type/:item_id'),
  flashMessage,
  id: text('[data-test-identity-item-id]'),
  nameContains: contains('[data-test-identity-item-name]'),
  rows: collection('[data-test-component="info-table-row"]', infoTableRow),
});
