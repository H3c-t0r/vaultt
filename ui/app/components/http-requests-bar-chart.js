import Component from '@ember/component';
import d3 from 'd3-selection';
import d3Scale from 'd3-scale';
import d3Axis from 'd3-axis';
import d3TimeFormat from 'd3-time-format';

/**
 * @module HttpRequestsBarChart
 * HttpRequestsBarChart components are used to...
 *
 * @example
 * ```js
 * <HttpRequestsBarChart @param1={param1} @param2={param2} />
 * ```
 *
 * @param param1 {String} - param1 is...
 * @param [param2=value] {String} - param2 is... //brackets mean it is optional and = sets the default value
 */

const COUNTERS = [
  {
    start_time: '2019-05-01T00:00:00Z',
    total: 50,
  },
  {
    start_time: '2019-04-01T00:00:00Z',
    total: 45,
  },
  {
    start_time: '2019-03-01T00:00:00Z',
    total: 55,
  },
];

export default Component.extend({
  tagName: '',
  data: null,
  svgContainer: null,
  margin: { top: 12, right: 12, bottom: 24, left: 24 },
  width() {
    const margin = this.margin;
    return 1344 - margin.left - margin.right;
  },
  height() {
    const margin = this.margin;
    return 240 - margin.top - margin.bottom;
  },

  didInsertElement() {
    this._super(...arguments);

    const data = COUNTERS;
    this.initBarChart(data);
  },

  initBarChart(dataIn) {
    const margin = this.margin,
      width = this.width(),
      height = this.height();

    const svgContainer = d3
      .select('.http-requests-bar-chart')
      .attr('width', width + margin.left + margin.right)
      .attr('height', height + margin.top + margin.bottom)
      .append('g')
      .attr('class', 'container')
      .attr('transform', 'translate(' + margin.left + ',' + margin.top + ')');

    this.set('svgContainer', svgContainer);

    this.barChart(dataIn);
  },

  barChart(dataIn) {
    const width = this.width(),
      height = this.height(),
      svgContainer = this.svgContainer;

    const counterTotals = dataIn.map(c => c.total);

    const yScale = d3Scale
      .scaleLinear()
      // the minimum and maximum value of the data
      .domain([0, Math.max(...counterTotals)])
      // how tall chart should be when we render it
      .range([height, 0]);

    const xScale = d3Scale
      .scaleBand()
      .domain(dataIn.map(c => c.start_time))
      // how wide it should be
      .range([0, width])
      // what % of total width it should reserve for whitespace between the bars
      .paddingInner(0.04);

    const yAxis = d3Axis.axisLeft(yScale).ticks(3);
    const xAxis = d3Axis.axisBottom(xScale).tickFormat(d3TimeFormat.timeFormat('%Y-%m-%d'));

    const xAxis_g = svgContainer
      .append('g')
      .attr('class', 'x axis')
      .attr('transform', 'translate(0,' + height + ')')
      .call(xAxis)
      .select('text');

    const yAxis_g = svgContainer
      .append('g')
      .attr('class', 'y axis')
      .attr('transform', 'translate(0,0)')
      .call(yAxis)
      .select('text');

    const bars = svgContainer.selectAll('.bar').data(dataIn);
    const labels = svgContainer.selectAll('.label').data(dataIn);

    bars
      // since the initial selection is empty (there are no bar elements yet), instantiate
      // the missing elements by appending to the enter selection
      .enter()
      .append('rect')
      .attr('class', 'bar')
      .attr('width', xScale.bandwidth())
      .attr('height', counter => height - yScale(counter.total))
      // the offset between each bar
      .attr('x', counter => xScale(counter.start_time))
      // 150 is the height of the svg
      .attr('y', counter => yScale(counter.total));

    labels
      .enter()
      .append('text')
      .attr('class', 'label')
      .attr('x', counter => xScale(counter.start_time) + xScale.bandwidth() / 2)
      .attr('y', counter => yScale(counter.total) + 15)
      .attr('dy', '.75em')
      .text(counter => counter.total)
      .style('margin', 15)
      .style('text-anchor', 'middle')
      .style('fill', '#ccc')
      .append('rect')
      .attr('width', 1)
      .attr('height', 5)
      .attr('x', counter => xScale(counter.start_time) + xScale.bandwidth() / 2)
      .attr('y', counter => height)
      .style('fill', 'blue');
  },
});
