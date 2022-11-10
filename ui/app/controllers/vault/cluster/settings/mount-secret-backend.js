import { inject as service } from '@ember/service';
import Controller from '@ember/controller';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
import { allEngines } from 'vault/helpers/mountable-secret-engines';
import { action } from '@ember/object';

const SUPPORTED_BACKENDS = supportedSecretBackends();

export default class MountSecretBackendController extends Controller {
  @service wizard;
  @service router;

  @action
  onMountSuccess(type, path) {
    let transition;
    if (SUPPORTED_BACKENDS.includes(type)) {
      const { engineRoute } = allEngines().findBy('type', type);
      if (engineRoute) {
        transition = this.router.transitionTo(`vault.cluster.secrets.backend.${engineRoute}`, path);
      } else if (type === 'keymgmt') {
        transition = this.router.transitionTo('vault.cluster.secrets.backend.index', path, {
          queryParams: { tab: 'provider' },
        });
      } else {
        transition = this.router.transitionTo('vault.cluster.secrets.backend.index', path);
      }
    } else {
      transition = this.router.transitionTo('vault.cluster.secrets.backends');
    }
    return transition.followRedirects().then(() => {
      this.wizard.transitionFeatureMachine(this.wizard.featureState, 'CONTINUE', type);
    });
  }
}
