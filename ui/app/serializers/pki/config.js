import ApplicationSerializer from '../application';

export default class PkiConfigSerializer extends ApplicationSerializer {
  attrs = {
    formType: { serialize: false },
  };
}
