/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { BAR_WIDTH, formatNumbers } from 'vault/utils/chart-helpers';
import { formatNumber } from 'core/helpers/format-number';
import { parseAPITimestamp } from 'core/utils/date-formatters';

import type { Count, MonthlyChartData } from 'vault/vault/charts/client-counts';

interface Args {
  data: MonthlyChartData[];
  yKey: string;
  chartTitle: string;
  chartHeight?: number;
}

interface ChartData {
  x: string;
  y: number | null;
  tooltip: string;
  legendX: string;
  legendY: string;
}

/**
 * @module VerticalBarBasic
 * Renders a vertical bar chart of counts (@yKey) over time.
 *
 * @example
 <Clients::Charts::VerticalBarBasic
    @chartTitle="Secret Sync client counts"
    @data={{this.model}}
    @yKey="secret_syncs"
    @showTable={{true}}
    @chartHeight={{200}}
  />
 */
export default class VerticalBarBasic extends Component<Args> {
  barWidth = BAR_WIDTH;

  @tracked activeDatum: ChartData | null = null;

  get chartHeight() {
    return this.args.chartHeight || 190;
  }

  get chartData() {
    return this.args.data.map((d): ChartData => {
      const xValue = d.timestamp as string;
      const yValue = (d[this.args.yKey as keyof Count] as number) ?? null;
      return {
        x: parseAPITimestamp(xValue, 'M/yy') as string,
        y: yValue,
        tooltip:
          yValue === null ? 'No data' : `${formatNumber([yValue])} ${this.args.yKey.replace(/_/g, ' ')}`,
        legendX: parseAPITimestamp(xValue, 'MMMM yyyy') as string,
        legendY: (yValue ?? 'No data').toString(),
      };
    });
  }

  get yDomain() {
    const counts: number[] = this.chartData
      .map((d) => d.y)
      .flatMap((num) => (typeof num === 'number' ? [num] : []));
    const max = Math.max(...counts);
    // if max is 0, hardcode 4 because that's the y-axis tickCount
    return [0, max === 0 ? 4 : max];
  }

  get xDomain() {
    const months = this.chartData.map((d) => d.x);
    return new Set(months);
  }

  // TEMPLATE HELPERS
  barOffset = (bandwidth: number) => {
    return (bandwidth - this.barWidth) / 2;
  };

  tooltipX = (original: number, bandwidth: number) => {
    return (original + bandwidth / 2).toString();
  };

  tooltipY = (original: number) => {
    if (!original) return `0`;
    return `${original}`;
  };

  formatTicksY = (num: number): string => {
    return formatNumbers(num) || num.toString();
  };
}
