import Model, { attr } from '@ember-data/model';

export default class PkiTidyModel extends Model {
  @attr('boolean', { defaultValue: false }) tidyCertStore;
  @attr('boolean', { defaultValue: false }) tidyRevocationQueue;
  @attr('string', { defaultValue: '72h' }) safetyBuffer;
}
