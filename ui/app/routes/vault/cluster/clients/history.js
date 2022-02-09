import Route from '@ember/routing/route';
import RSVP from 'rsvp';
import { action } from '@ember/object';

export default class HistoryRoute extends Route {
  async getActivity(start_time) {
    try {
      // on init ONLY make network request if we have a start time from the license
      // otherwise user needs to manually input
      return start_time ? await this.store.queryRecord('clients/activity', { start_time }) : {};
    } catch (e) {
      // ARG TODO handle
      return e;
    }
  }

  async getLicenseStartTime() {
    try {
      let license = await this.store.queryRecord('license', {});
      // if license.startTime is 'undefined' return 'null' for consistency
      return license.startTime || null;
    } catch (e) {
      // if error due to permission denied, return null so user can input date manually
      return e.httpStatus === 403 ? null : console.debug(e);
    }
  }

  async getVersionHistory() {
    try {
      let arrayOfModels = [];
      let response = await this.store.findAll('clients/version-history'); // returns a class with nested models
      response.forEach((model) => {
        arrayOfModels.push({
          id: model.id,
          perviousVersion: model.previousVersion,
          timestampInstalled: model.timestampInstalled,
        });
      });
      return arrayOfModels;
    } catch (e) {
      console.debug(e);
      return [];
    }
  }

  parseRFC3339(timestamp) {
    // convert '2021-03-21T00:00:00Z' --> ['2021', 2] (e.g. 2021 March, month is zero indexed)
    return timestamp
      ? [timestamp.split('-')[0], Number(timestamp.split('-')[1].replace(/^0+/, '')) - 1]
      : null;
  }

  async model() {
    let config = await this.store.queryRecord('clients/config', {}).catch((e) => {
      // swallowing error so activity can show if no config permissions
      console.debug(e);
      return {};
    });
    let licenseStart = await this.getLicenseStartTime();
    let activity = await this.getActivity(licenseStart);

    return RSVP.hash({
      config,
      activity,
      startTimeFromLicense: this.parseRFC3339(licenseStart),
      endTimeFromResponse: this.parseRFC3339(activity.endTime),
      versionHistory: this.getVersionHistory(),
    });
  }

  @action
  async loading(transition) {
    // eslint-disable-next-line ember/no-controller-access-in-routes
    let controller = this.controllerFor('vault.cluster.clients.history');
    if (controller) {
      controller.currentlyLoading = true;
      transition.promise.finally(function () {
        controller.currentlyLoading = false;
      });
    }
  }
}
