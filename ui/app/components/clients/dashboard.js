import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { isSameMonth } from 'date-fns';

export default class Dashboard extends Component {
  arrayOfMonths = [
    'January',
    'February',
    'March',
    'April',
    'May',
    'June',
    'July',
    'August',
    'September',
    'October',
    'November',
    'December',
  ];
  maxNamespaces = 10;
  chartLegend = [
    { key: 'entity_clients', label: 'entity clients' },
    { key: 'non_entity_clients', label: 'non-entity clients' },
  ];

  // needed for startTime modal picker
  months = Array.from({ length: 12 }, (item, i) => {
    return new Date(0, i).toLocaleString('en-US', { month: 'long' });
  });
  years = Array.from({ length: 5 }, (item, i) => {
    return new Date().getFullYear() - i;
  });
  searchSelectModel = ['namespace'];

  @service store;

  @tracked barChartSelection = false;
  @tracked isEditStartMonthOpen = false;
  @tracked responseRangeDiffMessage = null;
  @tracked startTimeRequested = null;
  @tracked startTimeFromResponse = this.args.model.startTimeFromLicense; // ex: ['2021', 3] is April 2021 (0 indexed)
  @tracked endTimeFromResponse = this.args.model.endTimeFromResponse;
  @tracked startMonth = null;
  @tracked startYear = null;
  @tracked selectedNamespace = null;
  @tracked noActivityDate = '';

  get startTimeDisplay() {
    if (!this.startTimeFromResponse) {
      // otherwise will return date of new Date(null)
      return null;
    }
    let month = this.startTimeFromResponse[1];
    let year = this.startTimeFromResponse[0];
    return `${this.arrayOfMonths[month]} ${year}`;
  }

  get endTimeDisplay() {
    if (!this.endTimeFromResponse) {
      // otherwise will return date of new Date(null)
      return null;
    }
    let month = this.endTimeFromResponse[1];
    let year = this.endTimeFromResponse[0];
    return `${this.arrayOfMonths[month]} ${year}`;
  }

  get isDateRange() {
    return !isSameMonth(
      new Date(this.args.model.activity.startTime),
      new Date(this.args.model.activity.endTime)
    );
  }

  // API client count data by namespace for date range
  get byNamespaceActivity() {
    return this.args.model.activity?.byNamespace || null;
  }

  // top level TOTAL client counts for given date range
  get totalUsageCounts() {
    return this.selectedNamespace
      ? this.filterByNamespace(this.selectedNamespace)
      : this.args.model.activity?.total;
  }

  // total client data for horizontal bar chart in attribution component
  get totalClientsData() {
    if (this.selectedNamespace) {
      let filteredNamespace = this.filterByNamespace(this.selectedNamespace);
      return filteredNamespace.mounts ? this.filterByNamespace(this.selectedNamespace).mounts : null;
    } else {
      return this.byNamespaceActivity;
    }
  }

  get responseTimestamp() {
    return this.args.model.activity?.responseTimestamp;
  }
  // HELPERS
  areArraysTheSame(a1, a2) {
    return (
      a1 === a2 ||
      (a1 !== null &&
        a2 !== null &&
        a1.length === a2.length &&
        a1
          .map(function (val, idx) {
            return val === a2[idx];
          })
          .reduce(function (prev, cur) {
            return prev && cur;
          }, true))
    );
  }

  // ACTIONS
  @action
  async handleClientActivityQuery(month, year, dateType) {
    if (dateType === 'cancel') {
      return;
    }
    // clicked "Current Billing period" in the calendar widget
    if (dateType === 'reset') {
      this.startTimeRequested = this.args.model.startTimeFromLicense;
      this.endTimeRequested = null;
    }
    // clicked "Edit" Billing start month in Dashboard which opens a modal.
    if (dateType === 'startTime') {
      let monthIndex = this.arrayOfMonths.indexOf(month);
      this.startTimeRequested = [year.toString(), monthIndex]; // ['2021', 0] (e.g. January 2021) // TODO CHANGE TO ARRAY
      this.endTimeRequested = null;
    }
    // clicked "Custom End Month" from the calendar-widget
    if (dateType === 'endTime') {
      // use the currently selected startTime for your startTimeRequested.
      this.startTimeRequested = this.startTimeFromResponse;
      this.endTimeRequested = [year.toString(), month]; // endTime comes in as a number/index whereas startTime comes in as a month name. Hence the difference between monthIndex and month.
    }

    try {
      let response = await this.store.queryRecord('clients/activity', {
        start_time: this.startTimeRequested,
        end_time: this.endTimeRequested,
      });
      if (response.id === 'no-data') {
        // empty response is the only time we want to update the displayed date with the requested time
        this.startTimeFromResponse = this.startTimeRequested;
        this.noActivityDate = this.startTimeDisplay;
      } else {
        // note: this.startTimeDisplay (getter) is updated by this.startTimeFromResponse
        this.startTimeFromResponse = response.formattedStartTime;
        this.endTimeFromResponse = response.formattedEndTime;
      }
      // compare if the response and what you requested are the same. If they are not throw a warning.
      // this only gets triggered if the data was returned, which does not happen if the user selects a startTime after for which we have data. That's an adapter error and is captured differently.
      if (!this.areArraysTheSame(this.startTimeFromResponse, this.startTimeRequested)) {
        this.responseRangeDiffMessage = `You requested data from ${month} ${year}. We only have data from ${this.startTimeDisplay}, and that is what is being shown here.`;
      } else {
        this.responseRangeDiffMessage = null;
      }
      return response;
    } catch (e) {
      // ARG TODO handle error
    }
  }

  @action
  handleCurrentBillingPeriod() {
    this.handleClientActivityQuery(0, 0, 'reset');
  }

  @action
  selectNamespace(value) {
    // value comes in as [namespace0]
    this.selectedNamespace = value[0];
  }

  @action
  selectStartMonth(month) {
    this.startMonth = month;
  }

  @action
  selectStartYear(year) {
    this.startYear = year;
  }

  // HELPERS
  filterByNamespace(namespace) {
    // ARG TODO test with normal data if it has a slash at the end.
    return this.byNamespaceActivity.find((ns) => ns.label === `${namespace}/`);
  }
}
