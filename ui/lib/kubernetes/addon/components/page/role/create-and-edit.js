import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { getRules } from '../../../utils/generated-role-rules';
import { htmlSafe } from '@ember/template';
import errorMessage from 'vault/utils/error-message';

export default class CreateAndEditRolePageComponent extends Component {
  @service router;
  @service flashMessages;

  @tracked roleRulesTemplates;
  @tracked selectedTemplateId = '1';
  @tracked modelValidations;

  constructor() {
    super(...arguments);
    // first check if generatedRoleRules matches one of the templates, the user may have chosen a template and not made changes
    // in this case we need to select the corresponding template in the dropdown
    // if there is no match then replace the example rules with the user defined value for no template option
    const { generatedRoleRules } = this.args.model;
    const rulesTemplates = getRules();
    if (generatedRoleRules) {
      const template = rulesTemplates.findBy('rules', generatedRoleRules);
      if (template) {
        this.selectedTemplateId = template.id;
      } else {
        rulesTemplates.findBy('1').rules = generatedRoleRules;
      }
    }
    this.roleRulesTemplates = rulesTemplates;
  }

  get generationPreferences() {
    return [
      {
        title: 'Generate token only using existing service account',
        description:
          'Enter a service account that already exists in Kubernetes and Vault will dynamically generate a token.',
        value: 'basic',
      },
      {
        title: 'Generate token, service account, and role binding objects',
        description:
          'Enter a pre-existing role (or ClusterRole) to use. Vault will generate a token, a service account and role binding objects.',
        value: 'expanded',
      },
      {
        title: 'Generate entire Kubernetes object chain',
        description:
          'Vault will generate the entire chain— a role, a token, a service account, and role binding objects— based on rules you supply.',
        value: 'full',
      },
    ];
  }

  get extraFields() {
    return [
      {
        type: 'annotations',
        key: 'extraAnnotations',
        description: 'Attach arbitrary non-identifying metadata to objects.',
      },
      {
        type: 'labels',
        key: 'extraLabels',
        description:
          'Labels specify identifying attributes of objects that are meaningful and relevant to users.',
      },
    ];
  }

  get roleRulesHelpText() {
    const message =
      'This specifies the Role or ClusterRole rules to use when generating a role. Kubernetes documentation is';
    const link =
      '<a href="https://kubernetes.io/docs/reference/access-authn-authz/rbac/" target="_blank" rel="noopener noreferrer">available here</>';
    return htmlSafe(`${message} ${link}.`);
  }

  @action
  resetRoleRules() {
    this.roleRulesTemplates = getRules();
  }

  @action
  selectTemplate(event) {
    this.selectedTemplateId = event.target.value;
    this.args.model.generationPreferences;
  }

  @task
  @waitFor
  *save() {
    try {
      yield this.args.model.save();
      this.router.transitionTo(
        'vault.cluster.secrets.backend.kubernetes.roles.role.details',
        this.args.model.name
      );
    } catch (error) {
      const message = errorMessage(error, 'Error saving role. Please try again or contact support');
      this.flashMessages.danger(message);
    }
  }

  @action
  async onSave(event) {
    event.preventDefault();
    const { isValid, state } = await this.args.model.validate();
    if (isValid) {
      this.modelValidations = null;
      this.save.perform();
    } else {
      this.flashMessages.info('Save not performed. Check form for errors');
      this.modelValidations = state;
    }
  }

  @action
  cancel() {
    const { model } = this.args;
    const method = model.isNew ? 'unloadRecord' : 'rollbackAttributes';
    model[method]();
    this.router.transitionTo('vault.cluster.secrets.backend.kubernetes.roles');
  }
}
