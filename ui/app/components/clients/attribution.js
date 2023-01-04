import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { format, isSameMonth } from 'date-fns';

/**
 * @module Attribution
 * Attribution components display the top 10 total client counts for namespaces or auth methods (mounts) during a billing period.
 * A horizontal bar chart shows on the right, with the top namespace/auth method and respective client totals on the left.
 *
 * @example
 * ```js
 *  <Clients::Attribution
 *    @chartLegend={{this.chartLegend}}
 *    @totalUsageCounts={{this.totalUsageCounts}}
 *    @newUsageCounts={{this.newUsageCounts}}
 *    @totalClientAttribution={{this.totalClientAttribution}}
 *    @newClientAttribution={{this.newClientAttribution}}
 *    @selectedNamespace={{this.selectedNamespace}}
 *    @startTimestamp={{this.responseTimestamp}}
 *    @endTimestamp={{this.responseTimestamp}}
 *    @isDateRange={{this.isDateRange}}
 *    @isCurrentMonth={{false}}
 *    @responseTimestamp={{this.responseTimestamp}}
 *  />
 * ```
 * @param {array} chartLegend - (passed to child) array of objects with key names 'key' and 'label' so data can be stacked
 * @param {object} totalUsageCounts - object with total client counts for chart tooltip text
 * @param {object} newUsageCounts - object with new client counts for chart tooltip text
 * @param {array} totalClientAttribution - array of objects containing a label and breakdown of client counts for total clients
 * @param {array} newClientAttribution - array of objects containing a label and breakdown of client counts for new clients
 * @param {string} selectedNamespace - namespace selected from filter bar
 * @param {string} startTimestamp - timestamp string from activity response to render start date for CSV modal
 * @param {string} endTimestamp - timestamp string from activity response to render end date for CSV modal
 * @param {string} responseTimestamp -  ISO timestamp created in serializer to timestamp the response, renders in bottom left corner below attribution chart
 * @param {boolean} isDateRange - getter calculated in parent to relay if dataset is a date range or single month and display text accordingly
 * @param {boolean} isCurrentMonth - boolean to determine if rendering data from current month
 */

export default class Attribution extends Component {
  @tracked showCSVDownloadModal = false;
  @service download;

  get formattedStartDate() {
    if (!this.args.startTimestamp) return null;
    return parseAPITimestamp(this.args.startTimestamp, 'MMMM yyyy');
  }
  get formattedEndDate() {
    if (!this.args.startTimestamp && !this.args.endTimestamp) return null;
    // displays on CSV export modal, no need to display duplicate months and years
    const startDateObject = parseAPITimestamp(this.args.startTimestamp);
    const endDateObject = parseAPITimestamp(this.args.endTimestamp);
    return isSameMonth(startDateObject, endDateObject) ? null : format(endDateObject, 'MMMM yyyy');
  }

  get hasCsvData() {
    return this.args.totalClientAttribution ? this.args.totalClientAttribution.length > 0 : false;
  }

  get isSingleNamespace() {
    if (!this.args.totalClientAttribution) {
      return 'no data';
    }
    // if a namespace is selected, then we're viewing top 10 auth methods (mounts)
    return !!this.args.selectedNamespace;
  }

  // truncate data before sending to chart component
  get barChartTotalClients() {
    return this.args.totalClientAttribution?.slice(0, 10);
  }

  get barChartNewClients() {
    return this.args.newClientAttribution?.slice(0, 10);
  }

  get topClientCounts() {
    // get top namespace or auth method
    return this.args.totalClientAttribution ? this.args.totalClientAttribution[0] : null;
  }

  get attributionBreakdown() {
    // display text for hbs
    return this.isSingleNamespace ? 'auth method' : 'namespace';
  }

  get chartText() {
    const dateText = this.args.isDateRange && !this.args.isCurrentMonth ? 'date range' : 'month';
    switch (this.isSingleNamespace) {
      case true:
        return {
          description:
            'This data shows the top ten authentication methods by client count within this namespace, and can be used to understand where clients are originating. Authentication methods are organized by path.',
          newCopy: `The new clients used by the auth method for this ${dateText}. This aids in understanding which auth methods create and use new clients${
            dateText === 'date range' ? ' over time.' : '.'
          }`,
          totalCopy: `The total clients used by the auth method for this ${dateText}. This number is useful for identifying overall usage volume. `,
        };
      case false:
        return {
          description:
            'This data shows the top ten namespaces by client count and can be used to understand where clients are originating. Namespaces are identified by path. To see all namespaces, export this data.',
          newCopy: `The new clients in the namespace for this ${dateText}.
          This aids in understanding which namespaces create and use new clients${
            dateText === 'date range' ? ' over time.' : '.'
          }`,
          totalCopy: `The total clients in the namespace for this ${dateText}. This number is useful for identifying overall usage volume.`,
        };
      case 'no data':
        return {
          description: 'There is a problem gathering data',
        };
      default:
        return '';
    }
  }

  destructureCountsToArray(object) {
    // destructure the namespace object  {label: 'some-namespace', entity_clients: 171, non_entity_clients: 20, clients: 191}
    // to get integers for CSV file
    const { clients, entity_clients, non_entity_clients } = object;
    return [clients, entity_clients, non_entity_clients];
  }

  constructCsvRow(namespaceColumn, mountColumn = null, totalColumns, newColumns = null) {
    // if namespaceColumn is a string, then we're at mount level attribution, otherwise it is an object
    // if constructing a namespace row, mountColumn=null so the column is blank, otherwise it is an object
    const otherColumns = newColumns ? [...totalColumns, ...newColumns] : [...totalColumns];
    return [
      `${typeof namespaceColumn === 'string' ? namespaceColumn : namespaceColumn.label}`,
      `${mountColumn ? mountColumn.label : ''}`,
      ...otherColumns,
    ];
  }

  generateCsvData() {
    const totalAttribution = this.args.totalClientAttribution;
    const newAttribution = this.barChartNewClients ? this.args.newClientAttribution : null;
    const csvData = [];
    const csvHeader = [
      'Namespace path',
      'Authentication method',
      'Total clients',
      'Entity clients',
      'Non-entity clients',
    ];

    if (newAttribution) {
      csvHeader.push('Total new clients, New entity clients, New non-entity clients');
    }

    totalAttribution.forEach((totalClientsObject) => {
      const namespace = this.isSingleNamespace ? this.args.selectedNamespace : totalClientsObject;
      const mount = this.isSingleNamespace ? totalClientsObject : null;

      // find new client data for namespace/mount object we're iterating over
      const newClientsObject = newAttribution
        ? newAttribution.find((d) => d.label === totalClientsObject.label)
        : null;

      const totalClients = this.destructureCountsToArray(totalClientsObject);
      const newClients = newClientsObject ? this.destructureCountsToArray(newClientsObject) : null;

      csvData.push(this.constructCsvRow(namespace, mount, totalClients, newClients));
      // constructCsvRow returns an array that corresponds to a row in the csv file:
      // ['ns label', 'mount label', total client #, entity #, non-entity #, ...new client #'s]

      // only iterate through mounts if NOT viewing a single namespace
      if (!this.isSingleNamespace && namespace.mounts) {
        namespace.mounts.forEach((mount) => {
          const newMountData = newAttribution
            ? newClientsObject?.mounts.find((m) => m.label === mount.label)
            : null;
          const mountTotalClients = this.destructureCountsToArray(mount);
          const mountNewClients = newMountData ? this.destructureCountsToArray(newMountData) : null;
          csvData.push(this.constructCsvRow(namespace, mount, mountTotalClients, mountNewClients));
        });
      }
    });

    csvData.unshift(csvHeader);
    // make each nested array a comma separated string, join each array "row" in csvData with line break (\n)
    return csvData.map((d) => d.join()).join('\n');
  }

  get formattedCsvFileName() {
    const endRange = this.formattedEndDate ? `-${this.formattedEndDate}` : '';
    const csvDateRange = this.formattedStartDate + endRange;
    return this.isSingleNamespace
      ? `clients_by_auth_method_${csvDateRange}`
      : `clients_by_namespace_${csvDateRange}`;
  }

  // ACTIONS
  @action
  exportChartData(filename) {
    const contents = this.generateCsvData();
    this.download.csv(filename, contents);
    this.showCSVDownloadModal = false;
  }
}
