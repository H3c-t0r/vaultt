import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { isAfter, isBefore, isSameDay, isSameMonth } from 'date-fns';
import getStorage from 'vault/lib/token-storage';
import { parseAPITimestamp } from 'core/utils/date-formatters';

// TODO CMB: change class and file name to Dashboard
export default class History extends Component {
  @service store;
  @service version;

  chartLegend = [
    { key: 'entity_clients', label: 'entity clients' },
    { key: 'non_entity_clients', label: 'non-entity clients' },
  ];

  // RESPONSE
  @tracked startMonthTimestamp = this.args.model.licenseStartTimestamp; // when user queries, updates to first month object of response
  @tracked endMonthTimestamp = this.args.model.currentDate; // when user queries, updates to last month object of response
  @tracked queriedActivityResponse = null;
  // track params sent to /activity request
  @tracked activityQueryParams = {
    start: { timestamp: this.args.model.licenseStartTimestamp }, // updates when user edits billing start month
    end: { timestamp: this.args.model.currentDate }, // updates when user queries end dates via calendar widget
  };

  // SEARCH SELECT
  @tracked selectedNamespace = null;
  @tracked namespaceArray = this.getActivityResponse.byNamespace
    ? this.getActivityResponse.byNamespace.map((namespace) => ({
        name: namespace.label,
        id: namespace.label,
      }))
    : [];
  @tracked selectedAuthMethod = null;
  @tracked authMethodOptions = [];

  // TEMPLATE VIEW
  @tracked showBillingStartModal = false;
  @tracked noActivityRange = '';
  @tracked isLoadingQuery = false;
  @tracked errorObject = null;

  get versionText() {
    return this.version.isEnterprise
      ? {
          label: 'Billing start month',
          description:
            'This date comes from your license, and defines when client counting starts. Without this starting point, the data shown is not reliable.',
          title: 'No billing start date found',
          message:
            'In order to get the most from this data, please enter your billing period start month. This will ensure that the resulting data is accurate.',
        }
      : {
          label: 'Client counting start date',
          description:
            'This date is when client counting starts. Without this starting point, the data shown is not reliable.',
          title: 'No start date found',
          message:
            'In order to get the most from this data, please enter a start month above. Vault will calculate new clients starting from that month.',
        };
  }

  get isDateRange() {
    return !isSameMonth(
      parseAPITimestamp(this.getActivityResponse.startTime),
      parseAPITimestamp(this.getActivityResponse.endTime)
    );
  }

  get startTimeDiscrepancy() {
    // show banner if startTime returned from activity log (response) is after the queried startTime
    const activityStartDateObject = parseAPITimestamp(this.getActivityResponse.startTime);
    const queryStartDateObject = parseAPITimestamp(this.startMonthTimestamp);
    const isLicenseStart = isSameDay(
      queryStartDateObject,
      parseAPITimestamp(this.args.model.licenseStartTimestamp)
    );

    if (isAfter(activityStartDateObject, queryStartDateObject)) {
      const message = isLicenseStart ? 'Your license start date is ' : 'You requested data from ';
      return (
        message +
        `${parseAPITimestamp(this.startMonthTimestamp, 'MMMM yyyy')}. 
        We only have data from ${parseAPITimestamp(this.getActivityResponse.startTime, 'MMMM yyyy')}, 
        and that is what is being shown here.`
      );
    } else {
      return null;
    }
  }

  get upgradeDuringActivity() {
    const versionHistory = this.args.model.versionHistory;
    if (!versionHistory || versionHistory.length === 0) {
      return null;
    }

    // filter for upgrade data of noteworthy upgrades (1.9 and/or 1.10)
    const upgradeVersionHistory = versionHistory.filter(
      (version) => version.id.match('1.9') || version.id.match('1.10')
    );
    if (!upgradeVersionHistory || upgradeVersionHistory.length === 0) {
      return null;
    }

    const activityStart = parseAPITimestamp(this.getActivityResponse.startTime);
    const activityEnd = parseAPITimestamp(this.getActivityResponse.endTime);
    // filter and return all upgrades that happened within date range of queried activity
    return upgradeVersionHistory.filter(({ timestampInstalled }) => {
      const upgradeDate = parseAPITimestamp(timestampInstalled);
      return isAfter(upgradeDate, activityStart) && isBefore(upgradeDate, activityEnd);
    });
  }

  get upgradeVersionAndDate() {
    if (!this.upgradeDuringActivity || this.upgradeDuringActivity.length === 0) {
      return null;
    }
    if (this.upgradeDuringActivity.length === 2) {
      let [firstUpgrade, secondUpgrade] = this.upgradeDuringActivity;
      let firstDate = parseAPITimestamp(firstUpgrade.timestampInstalled, 'MMM d, yyyy');
      let secondDate = parseAPITimestamp(secondUpgrade.timestampInstalled, 'MMM d, yyyy');
      return `Vault was upgraded to ${firstUpgrade.id} (${firstDate}) and ${secondUpgrade.id} (${secondDate}) during this time range.`;
    } else {
      let [upgrade] = this.upgradeDuringActivity;
      return `Vault was upgraded to ${upgrade.id} on ${parseAPITimestamp(
        upgrade.timestampInstalled,
        'MMM d, yyyy'
      )}.`;
    }
  }

  get versionSpecificText() {
    if (!this.upgradeDuringActivity || this.upgradeDuringActivity.length === 0) {
      return null;
    }
    if (this.upgradeDuringActivity.length === 1) {
      let version = this.upgradeDuringActivity[0].id;
      if (version.match('1.9')) {
        return ' How we count clients changed in 1.9, so keep that in mind when looking at the data below.';
      }
      if (version.match('1.10')) {
        return ' We added monthly breakdowns and mount level attribution starting in 1.10, so keep that in mind when looking at the data below.';
      }
    }
    // return combined explanation if spans multiple upgrades
    return ' How we count clients changed in 1.9 and we added monthly breakdowns and mount level attribution starting in 1.10. Keep this in mind when looking at the data below.';
  }

  get displayStartDate() {
    if (!this.startMonthTimestamp) return null;
    return parseAPITimestamp(this.startMonthTimestamp, 'MMMM yyyy');
  }

  // GETTERS FOR RESPONSE & DATA

  // on init API response uses license start_date, getter updates when user queries dates
  get getActivityResponse() {
    return this.queriedActivityResponse || this.args.model.activity;
  }

  get byMonthActivityData() {
    if (this.selectedNamespace) {
      return this.filteredActivityByMonth;
    } else {
      return this.getActivityResponse?.byMonth;
    }
  }

  get byMonthNewClients() {
    if (this.byMonthActivityData) {
      return this.byMonthActivityData?.map((m) => m.new_clients);
    }
    return null;
  }

  get hasAttributionData() {
    if (this.selectedAuthMethod) return false;
    if (this.selectedNamespace) {
      return this.authMethodOptions.length > 0;
    }
    return !!this.totalClientAttribution && this.totalUsageCounts && this.totalUsageCounts.clients !== 0;
  }

  // (object) top level TOTAL client counts for given date range
  get totalUsageCounts() {
    return this.selectedNamespace ? this.filteredActivityByNamespace : this.getActivityResponse.total;
  }

  // (object) single month new client data with total counts + array of namespace breakdown
  get newClientCounts() {
    return this.isDateRange ? null : this.byMonthActivityData[0]?.new_clients;
  }

  // total client data for horizontal bar chart in attribution component
  get totalClientAttribution() {
    if (this.selectedNamespace) {
      return this.filteredActivityByNamespace?.mounts || null;
    } else {
      return this.getActivityResponse?.byNamespace || null;
    }
  }

  // new client data for horizontal bar chart
  get newClientAttribution() {
    // new client attribution only available in a single, historical month (not a date range)
    if (this.isDateRange) return null;

    if (this.selectedNamespace) {
      return this.newClientCounts?.mounts || null;
    } else {
      return this.newClientCounts?.namespaces || null;
    }
  }

  get responseTimestamp() {
    return this.getActivityResponse.responseTimestamp;
  }

  // FILTERS
  get filteredActivityByNamespace() {
    const namespace = this.selectedNamespace;
    const auth = this.selectedAuthMethod;
    if (!namespace && !auth) {
      return this.getActivityResponse;
    }
    if (!auth) {
      return this.getActivityResponse.byNamespace.find((ns) => ns.label === namespace);
    }
    return this.getActivityResponse.byNamespace
      .find((ns) => ns.label === namespace)
      .mounts?.find((mount) => mount.label === auth);
  }

  get filteredActivityByMonth() {
    const namespace = this.selectedNamespace;
    const auth = this.selectedAuthMethod;
    if (!namespace && !auth) {
      return this.getActivityResponse?.byMonth;
    }
    const namespaceData = this.getActivityResponse?.byMonth
      .map((m) => m.namespaces_by_key[namespace])
      .filter((d) => d !== undefined);
    if (!auth) {
      return namespaceData.length === 0 ? null : namespaceData;
    }
    const mountData = namespaceData
      .map((namespace) => namespace.mounts_by_key[auth])
      .filter((d) => d !== undefined);
    return mountData.length === 0 ? null : mountData;
  }

  @action
  async handleClientActivityQuery({ dateType, monthIdx, year }) {
    this.showBillingStartModal = false;
    switch (dateType) {
      case 'cancel':
        return;
      case 'reset': // reset to initial start/end dates (current billing period)
        this.activityQueryParams.start.timestamp = this.args.model.licenseStartTimestamp;
        this.activityQueryParams.end.timestamp = this.args.model.currentDate;
        break;
      case 'startDate': // from "Edit billing start" modal
        this.activityQueryParams.start = { monthIdx, year };
        break;
      case 'endDate': // selected end date from calendar widget
        this.activityQueryParams.end = { monthIdx, year };
        break;
      default:
        break;
    }
    try {
      this.isLoadingQuery = true;
      let response = await this.store.queryRecord('clients/activity', {
        start_time: this.activityQueryParams.start,
        end_time: this.activityQueryParams.end,
      });
      if (response.id === 'no-data') {
        // if an empty response (204) the adapter returns the queried time params (instead of the backend's activity log start/end times)
        const endMonth = isSameMonth(
          parseAPITimestamp(response.startTime),
          parseAPITimestamp(response.endTime)
        )
          ? ''
          : ` to ${parseAPITimestamp(response.endTime, 'MMMM yyyy')}`;
        this.noActivityRange = `from ${parseAPITimestamp(response.startTime, 'MMMM yyyy')}` + endMonth;
      } else {
        // TODO cmb - would like to remove using the byMonth timestamps and just rely on response to return start/end times
        const { byMonth } = response;
        this.startMonthTimestamp = byMonth[0]?.timestamp || response.startTime;
        this.endMonthTimestamp = byMonth[byMonth.length - 1]?.timestamp || response.endTime;
        getStorage().setItem('vault:ui-inputted-start-date', this.startMonthTimestamp);
      }
      this.queriedActivityResponse = response;
    } catch (e) {
      this.errorObject = e;
      return e;
    } finally {
      this.isLoadingQuery = false;
    }
  }

  get hasMultipleMonthsData() {
    return this.byMonthActivityData && this.byMonthActivityData.length > 1;
  }

  @action
  selectNamespace([value]) {
    this.selectedNamespace = value;
    if (!value) {
      this.authMethodOptions = [];
      // on clear, also make sure auth method is cleared
      this.selectedAuthMethod = null;
    } else {
      // Side effect: set auth namespaces
      const mounts = this.filteredActivityByNamespace.mounts?.map((mount) => ({
        id: mount.label,
        name: mount.label,
      }));
      this.authMethodOptions = mounts;
    }
  }

  @action
  setAuthMethod([authMount]) {
    this.selectedAuthMethod = authMount;
  }
}
