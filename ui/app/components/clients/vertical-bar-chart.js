import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { max } from 'd3-array';
// eslint-disable-next-line no-unused-vars
import { select, selectAll, node } from 'd3-selection';
import { axisLeft, axisBottom } from 'd3-axis';
import { scaleLinear, scaleBand } from 'd3-scale';
import { stack } from 'd3-shape';
import {
  GREY,
  LIGHT_AND_DARK_BLUE,
  SVG_DIMENSIONS,
  TRANSLATE,
  formatNumbers,
} from '../../utils/chart-helpers';

// TODO fill out below
/**
 * @module VerticalBarChart
 * VerticalBarChart components are used to...
 *
 * @example
 * ```js
 * <VerticalBarChart @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

export default class VerticalBarChart extends Component {
  @tracked tooltipTarget = '';
  @tracked hoveredLabel = '';

  get chartLegend() {
    return this.args.chartLegend;
  }

  @action
  registerListener(element, args) {
    // Define the chart
    let dataset = args[0];
    console.log(this.chartLegend);
    let stackFunction = stack().keys(this.chartLegend.map(l => l.key));
    let stackedData = stackFunction(dataset);

    // TODO pull out into helper? b/c same as line bar chart?
    let yScale = scaleLinear()
      .domain([0, max(dataset.map(d => d.total))]) // TODO will need to recalculate when you get the data
      .range([0, 100])
      .nice();

    let xScale = scaleBand()
      .domain(dataset.map(d => d.month))
      .range([0, SVG_DIMENSIONS.width]) // set width to fix number of pixels
      .paddingInner(0.85);

    let chartSvg = select(element);

    chartSvg.attr('viewBox', `-50 20 600 ${SVG_DIMENSIONS.height}`); // set svg dimensions

    let groups = chartSvg
      .selectAll('g')
      .data(stackedData)
      .enter()
      .append('g')
      .style('fill', (d, i) => LIGHT_AND_DARK_BLUE[i]);

    groups
      .selectAll('rect')
      .data(stackedData => stackedData)
      .enter()
      .append('rect')
      .attr('width', '7px')
      .attr('class', 'data-bar')
      .attr('height', stackedData => `${yScale(stackedData[1] - stackedData[0])}%`)
      .attr('x', ({ data }) => xScale(data.month)) // uses destructuring because was data.data.month
      .attr('y', data => `${100 - yScale(data[1])}%`); // subtract higher than 100% to give space for x axis ticks

    // MAKE AXES //
    let yAxisScale = scaleLinear()
      .domain([0, max(dataset.map(d => d.total))]) // TODO will need to recalculate when you get the data
      .range([`${SVG_DIMENSIONS.height}`, 0])
      .nice();

    let yAxis = axisLeft(yAxisScale)
      .ticks(7)
      .tickPadding(10)
      .tickSizeInner(-SVG_DIMENSIONS.width)
      .tickFormat(formatNumbers);

    let xAxis = axisBottom(xScale).tickSize(0);

    yAxis(chartSvg.append('g'));
    xAxis(chartSvg.append('g').attr('transform', `translate(0, ${SVG_DIMENSIONS.height + 10})`));

    chartSvg.selectAll('.domain').remove(); // remove domain lines

    // creating wider area for tooltip hover
    let greyBars = chartSvg
      .append('g')
      .attr('transform', `translate(${TRANSLATE.left})`)
      .style('fill', `${GREY}`)
      .style('opacity', '0')
      .style('mix-blend-mode', 'multiply');

    let tooltipRect = greyBars
      .selectAll('rect')
      .data(dataset)
      .enter()
      .append('rect')
      .style('cursor', 'pointer')
      .attr('class', 'tooltip-rect')
      .attr('height', '100%')
      .attr('width', '30px') // three times width
      .attr('y', '0') // start at bottom
      .attr('x', data => xScale(data.month)); // not data.data because this is not stacked data

    // for tooltip
    tooltipRect.on('mouseover', data => {
      this.hoveredLabel = data.month;
      // let node = chartSvg
      //   .selectAll('rect.tooltip-rect')
      //   .filter(data => data.month === this.hoveredLabel)
      //   .node();
      let node = chartSvg
        .selectAll('rect.data-bar')
        // filter for the top data bar (so y-coord !== 0) with matching month
        .filter(data => data[0] !== 0 && data.data.month === this.hoveredLabel)
        .node();
      this.tooltipTarget = node; // grab the node from the list of rects
    });
  }

  @action removeTooltip() {
    this.tooltipTarget = null;
  }
}
