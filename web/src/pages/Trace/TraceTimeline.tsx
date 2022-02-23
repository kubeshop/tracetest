import {useEffect, useRef} from 'react';
import styled from 'styled-components';
import * as d3 from 'd3';
import Title from 'antd/lib/typography/Title';

import data from './data.json';
import './TimelineChart.css';

const Header = styled.div`
  display: flex;
  align-items: center;
  width: 100%;
  height: 56px;
  padding: 0 16px;
  color: rgb(213, 215, 224);
`;

interface IProps {
  selectedSpan: any;
  onSelectSpan: (span: any) => void;
}

const TraceTimeline = ({selectedSpan, onSelectSpan}: IProps) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const svgRef = useRef<SVGSVGElement>(null);

  let treeFactory = d3.tree().size([200, 450]).nodeSize([0, 5]);

  const spanDates = data.resourceSpans
    .map((i: any) => i.instrumentationLibrarySpans.map((el: any) => el.spans))
    .flat(2)
    .map(span => ({
      startTime: new Date(Number(span.startTimeUnixNano) / 1000),
      endTime: new Date(Number(span.endTimeUnixNano) / 1000),
      span,
    }));

  const rootNode = data.resourceSpans
    .map((i: any) => i.instrumentationLibrarySpans.map((el: any) => el.spans))
    .flat(2)
    .sort((el1, el2) => Number(el1.startTimeUnixNano) - Number(el2.startTimeUnixNano))
    .shift();

  const spanMap = data.resourceSpans
    .map((i: any) => i.instrumentationLibrarySpans.map((el: any) => el.spans))
    .flat(2)
    .reduce((acc: {[key: string]: {id: string; parentIds: string[]; data: any}}, span) => {
      acc[span.spanId] = acc[span.spanId] || {
        id: span.spanId,
        parentIds: [],
        data: span,
        startTime: new Date(Number(span.startTimeUnixNano) / 1000),
        endTime: new Date(Number(span.endTimeUnixNano) / 1000),
      };
      acc[span.spanId].parentIds.push(span.parentSpanId);
      return acc;
    }, {});

  const dagData = Object.values(spanMap).map(({id, parentIds, ...rest}) => {
    const parents = parentIds.filter(el => spanMap[el]);
    if ((!parents || parents.length === 0) && id !== rootNode.spanId) {
      parents.push(rootNode.spanId);
    }

    return {id, parentId: parents[0], ...rest};
  });

  const root = d3.stratify()(dagData);
  const minNano = d3.min(spanDates, s => Number(s.span.startTimeUnixNano))! as number;

  const scaleTime = d3
    .scaleLinear()
    .domain([
      0,
      (d3.max(spanDates, s => Number(s.span.endTimeUnixNano))! -
        d3.min(spanDates, s => Number(s.span.startTimeUnixNano))!) as number,
    ])
    .range([250, 800]);

  const barHeight = 20;
  const theBarHeight = barHeight;

  useEffect(() => {
    let nodes = treeFactory(root);
    let nodesSort: any[] = [];
    console.log('@@nodes', nodes);
    nodes.sort((a: any, b: any) =>
      a.depth === b.depth ? a.data.startTime.getTime() - b.data.startTime.getTime() : a.depth - b.depth
    );
    nodes.eachBefore(n => nodesSort.push(n));

    nodesSort.forEach((n, i) => {
      n.x = i * barHeight;
    });
    const height = barHeight * nodes.descendants().length + 50;
    const chart = d3.select(svgRef.current).attr('viewBox', `0 0 800 ${height}`);

    const xAxis = d3
      .axisTop(scaleTime)
      .ticks(5)
      .tickFormat(d => `${Number(d) / 1000000}ms`);

    const milliTicks = d3.ticks(
      0,
      (d3.max(spanDates, d => Number(d.span.endTimeUnixNano))! -
        d3.min(spanDates, d => Number(d.span.startTimeUnixNano))!) as number,
      5
    );

    const ticks = chart.append('g').attr('transform', 'translate(0,20)').call(xAxis);

    ticks
      .selectAll('text')
      .style('text-anchor', 'middle')
      .attr('fill', '#000')
      .attr('stroke', 'none')
      .attr('font-size', 10);

    ticks.select('.domain').attr('stroke', 'none').attr('opacity', '0');
    ticks.selectAll('.tick line').attr('stroke', 'none');
    const grid = chart.append('g').selectAll('rect').data(milliTicks).enter();

    grid
      .append('rect')
      .attr('x', d => {
        return scaleTime(d) - 0.5;
      })
      .attr('y', 20)
      .attr('width', 1)
      .attr('height', height)
      .attr('stroke', 'none')
      .attr('fill', 'rgb(213, 215, 224)');
    chart.append('g').attr('class', 'container').attr('transform', `translate(0 , 30 )`);
  }, []);

  useEffect(() => {
    drawChart();
  }, [selectedSpan]);

  const drawChart = () => {
    let nodes = treeFactory(root);
    let nodesSort: any[] = [];
    nodes.sort((a: any, b: any) =>
      a.depth === b.depth ? a.data.startTime.getTime() - b.data.startTime.getTime() : a.depth - b.depth
    );
    nodes.eachBefore(n => nodesSort.push(n));

    nodesSort.forEach((n, i) => {
      n.x = i * barHeight;
    });

    const chart = d3.select(svgRef.current);

    const node = chart
      .select('g.container')
      .selectAll('g.node')
      .data(nodesSort, (d: any) => d.data.spanID);

    const nodeEnter = node
      .enter()
      .append('g')
      .attr('class', 'node')
      .attr('id', el => el.id)
      .attr('x', el => el.x)
      .attr('y', el => el.y)
      .attr('height', barHeight)
      .attr('cursor', 'pointer')
      .attr('pointer-events', 'bounding-box')
      .on('click', (event, d) => {
        onSelectSpan({data: d.data.data, id: d.id});
      });

    nodeEnter
      .append('rect')
      .attr('class', d => `rect-svg ${d.id === selectedSpan.id ? 'rect-svg-selected' : ''}`)
      .attr('rx', 3)
      .attr('ry', 3)
      .attr('x', 0)
      .attr('width', '100%')
      .attr('height', theBarHeight)
      .attr('stroke', 'none')
      .attr('fill', 'none');

    nodeEnter
      .append('g')
      .attr('transform', d => `translate(${d.y + 5},0)`)
      .append('path')
      .attr('transform', 'scale(0.5)')
      .on('click', (event, d) => {
        event.stopPropagation();
        if (d.children) {
          d._children = d.children;
          d.children = null;
        } else {
          d.children = d._children;
          d._children = null;
        }
        drawChart();
      });
    nodeEnter
      .append('rect')
      .attr('rx', 3)
      .attr('ry', 3)
      .attr('x', d => scaleTime(Number(d.data.data.startTimeUnixNano) - minNano))
      .attr('y', 5)
      .attr('width', (d: any) => scaleTime(Number(d.data.data.endTimeUnixNano) - Number(d.data.data.startTimeUnixNano)))
      .attr('height', theBarHeight / 2)
      .attr('stroke', 'none')
      .attr('fill', e => (e.depth < 2 ? 'rgb(70, 74, 102)' : 'rgb(29, 233, 182)'))
      .attr('pointer-events', 'none');

    nodeEnter
      .append('text')
      .attr('width', 180)
      .attr('y', 10)
      .attr('height', barHeight)
      .attr('fill', '#000')
      .attr('font-size', 8)
      .attr('pointer-events', 'none')
      .attr('alignment-baseline', 'middle')
      .attr('dominant-baseline', 'middle')
      .text((d: any) => d.data.data.operationName);

    let nodeUpdate = node.merge(nodeEnter as any);

    nodeUpdate
      .attr('transform', (d: any) => `translate(${0} ,${d.x})`)
      .select('rect')
      .attr('class', d => `rect-svg ${d.id === selectedSpan.id ? 'rect-svg-selected' : ''}`);

    nodeUpdate
      .select('path')
      .attr('opacity', d => (d.children || d._children ? 1 : 0))
      .attr('d', d =>
        d._children
          ? 'M18.629 15.997l-7.083-7.081L13.462 7l8.997 8.997L13.457 25l-1.916-1.916z'
          : 'M16.003 18.626l7.081-7.081L25 13.46l-8.997 8.998-9.003-9 1.917-1.916z'
      );

    nodeUpdate.selectAll('text').attr('transform', (d: any) => `translate(${d.y + 20} ,0)`);

    node.exit().remove();

    nodesSort.forEach((d: any) => {
      d.x0 = d.x;
      d.y0 = d.y;
    });
  };

  return (
    <div ref={containerRef}>
      <Header>
        <Title level={4}>Component Timeline</Title>
      </Header>
      <svg ref={svgRef} />
    </div>
  );
};

export default TraceTimeline;
