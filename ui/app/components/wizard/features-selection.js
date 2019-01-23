import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { computed } from '@ember/object';
import { FEATURE_MACHINE_TIME } from 'vault/helpers/wizard-constants';

export default Component.extend({
  wizard: service(),
  version: service(),
  permissions: service(),

  init() {
    this._super(...arguments);
    this.maybeHideFeatures();
  },

  maybeHideFeatures() {
    let features = this.get('allFeatures');
    features.forEach(feat => {
      feat.disabled = !this.get(`has${feat.name}Permission`);
    });

    if (this.get('showReplication') === false) {
      let feature = this.get('allFeatures').findBy('key', 'replication');
      feature.show = false;
    }
  },

  hasSecretsPermission: computed(function() {
    return this.permissions.hasPermission('sys/mounts/example', 'update');
  }),

  hasAuthenticationPermission: computed(function() {
    const canRead = this.permissions.hasPermission('sys/auth', 'read');

    const capabilities = ['update', 'sudo'];
    const canUpdateOrCreate = capabilities.every(capability => {
      return this.permissions.hasPermission('sys/auth/example', capability);
    });

    return canRead && canUpdateOrCreate;
  }),

  hasPoliciesPermission: computed(function() {
    return this.permissions.hasPermission('sys/policies/acl', 'list');
  }),

  hasReplicationPermission: computed(function() {
    const PATHS = ['sys/replication/performance/primary/enable', 'sys/replication/dr/primary/enable'];
    return PATHS.every(path => {
      return this.permissions.hasPermission(path, 'update');
    });
  }),

  hasToolsPermission: computed(function() {
    const PATHS = ['sys/wrapping/wrap', 'sys/wrapping/lookup', 'sys/wrapping/unwrap', 'sys/wrapping/rewrap'];
    return PATHS.every(path => {
      return this.permissions.hasPermission(path, 'update');
    });
  }),

  estimatedTime: computed('selectedFeatures', function() {
    let time = 0;
    for (let feature of Object.keys(FEATURE_MACHINE_TIME)) {
      if (this.selectedFeatures.includes(feature)) {
        time += FEATURE_MACHINE_TIME[feature];
      }
    }
    return time;
  }),
  selectProgress: computed('selectedFeatures', function() {
    let bar = this.selectedFeatures.map(feature => {
      return { style: 'width:0%;', completed: false, showIcon: true, feature: feature };
    });
    if (bar.length === 0) {
      bar = [{ style: 'width:0%;', showIcon: false }];
    }
    return bar;
  }),
  allFeatures: computed(function() {
    return [
      {
        key: 'secrets',
        name: 'Secrets',
        steps: ['Enabling a secrets engine', 'Adding a secret'],
        selected: false,
        show: true,
        disabled: false,
      },
      {
        key: 'authentication',
        name: 'Authentication',
        steps: ['Enabling an auth method', 'Managing your auth method'],
        selected: false,
        show: true,
        disabled: false,
      },
      {
        key: 'policies',
        name: 'Policies',
        steps: [
          'Choosing a policy type',
          'Creating a policy',
          'Deleting your policy',
          'Other types of policies',
        ],
        selected: false,
        show: true,
        disabled: false,
      },
      {
        key: 'replication',
        name: 'Replication',
        steps: ['Setting up replication', 'Your cluster information'],
        selected: false,
        show: true,
        disabled: false,
      },
      {
        key: 'tools',
        name: 'Tools',
        steps: ['Wrapping data', 'Lookup wrapped data', 'Rewrapping your data', 'Unwrapping your data'],
        selected: false,
        show: true,
        disabled: false,
      },
    ];
  }),

  showReplication: computed('version.{hasPerfReplication,hasDRReplication}', function() {
    return this.get('version.hasPerfReplication') || this.get('version.hasDRReplication');
  }),

  selectedFeatures: computed('allFeatures.@each.selected', function() {
    return this.get('allFeatures')
      .filterBy('selected')
      .mapBy('key');
  }),

  cannotStartWizard: computed('selectedFeatures', function() {
    return !this.get('selectedFeatures').length;
  }),

  actions: {
    saveFeatures() {
      let wizard = this.get('wizard');
      wizard.saveFeatures(this.get('selectedFeatures'));
      wizard.transitionTutorialMachine('active.select', 'CONTINUE');
    },
  },
});
