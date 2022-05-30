import {useCallback, useEffect, useMemo, useRef} from 'react';
import * as d3 from 'd3';
import TraceAnalyticsService from '../../../services/Analytics/TraceAnalytics.service';
import {IDiagramProps} from '../Diagram';
import {TSpan} from '../../../types/Span.types';
import {getNotchColor} from '../../TraceNode/TraceNode.styled';
import * as S from './TimelineChart.styled';

const {onTimelineSpanClick} = TraceAnalyticsService;

const barHeight = 54;

export const TimelineChart = ({affectedSpans, trace, selectedSpan, onSelectSpan}: IDiagramProps) => {
  const svgRef = useRef<SVGSVGElement>(null);
  const treeFactory = d3.tree().size([200, 450]).nodeSize([0, 5]);

  const spanDates = trace.spans.map(span => ({
    startTime: span.startTime,
    endTime: span.endTime,
    span,
  }));

  const spanMap = trace.spans.reduce(
    (acc: {[key: string]: {id: string; parentIds: Array<string | undefined>; data: any}}, span) => {
      acc[span.id] = acc[span.id] || {
        id: span.id,
        parentIds: [],
        data: span,
        startTime: span.startTime,
        endTime: span.endTime,
      };
      acc[span.id].parentIds.push(span.parentId);
      return acc;
    },
    {}
  );

  const root = useMemo(() => {
    const dagData = Object.values(spanMap).map(({id, parentIds, ...rest}) => {
      const parent = parentIds.find(el => spanMap[el!]);

      return {id, parentId: parent, ...rest};
    });

    return d3.stratify()(dagData);
  }, [spanMap]);

  const minNano = d3.min(spanDates, s => s.startTime) as number;
  const maxNano = d3.max(spanDates, s => s.endTime) as number;

  const scaleTime = d3
    .scaleLinear()
    .domain([0, maxNano - minNano])
    .range([250, 800]);

  useEffect(() => {
    const nodes = treeFactory(root);
    const nodesSort: any[] = [];

    nodes.sort((a: any, b: any) => (a.depth === b.depth ? a.data.startTime - b.data.startTime : a.depth - b.depth));
    nodes.eachBefore(n => nodesSort.push(n));

    nodesSort.forEach((n, i) => {
      n.x = i * barHeight;
    });
    const height = barHeight * nodes.descendants().length + 50;
    const chart = d3.select(svgRef.current).attr('viewBox', `0 0 810 ${height}`);

    const xAxis = d3
      .axisTop(scaleTime)
      .ticks(10)
      .tickFormat(d => {
        return `${Number(d)}`;
      });

    const milliTicks = d3.ticks(0, maxNano - minNano, 10);

    const ticks = chart.append('g').attr('transform', 'translate(0,20)').call(xAxis);

    ticks.selectAll('text').attr('class', 'tick').style('text-anchor', 'middle');
    ticks.select('.domain').attr('stroke', 'none').attr('opacity', '0');
    ticks.selectAll('.tick line').attr('stroke', 'none');

    const grid = chart.append('g').selectAll('rect').data(milliTicks).enter();

    chart.append('text').attr('class', 'duration-ms-text').attr('x', 230).attr('y', 20).text('Duration (ms)');
    chart.append('rect').attr('class', 'cross-line').attr('y', 30);

    grid
      .append('rect')
      .attr('class', 'checkpoint-mark')
      .attr('x', d => {
        return scaleTime(d) - 0.5;
      })
      .attr('y', 20);
    chart.append('g').attr('class', 'container').attr('transform', `translate(0, 50)`);
  }, [trace]);

  const drawChart = useCallback(() => {
    const nodes = treeFactory(root);
    const nodesSort: any[] = [];
    nodes.sort((a: any, b: any) => (a.depth === b.depth ? a.data.startTime - b.data.startTime : a.depth - b.depth));
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
      .on('click', (event, d) => {
        onTimelineSpanClick(d.id);
        if (onSelectSpan) onSelectSpan(d.id);
      });

    nodeEnter
      .append('rect')
      .attr('class', d => `rect-svg ${d.id === selectedSpan?.id ? 'rect-svg-selected' : ''}`)
      .attr('rx', 3)
      .attr('ry', 3)
      .attr('x', 0)
      .attr('y', -(barHeight / 4));

    nodeEnter
      .append('g')
      .attr('class', 'chevron')
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

    nodeEnter.append('rect').attr('class', 'grey-line').attr('y', 30);

    nodeEnter
      .append('rect')
      .attr('class', 'duration-line')
      .attr('rx', 3)
      .attr('ry', 3)
      .attr('y', 25)
      .attr('x', d => {
        return scaleTime(Number(d.data.startTime || minNano) - minNano);
      })
      .attr('width', d => {
        return scaleTime(d.data.data.duration) - 250;
      })
      .attr('fill', e => {
        const span: TSpan = e.data.data;

        const color = getNotchColor(span.type);

        return color;
      });

    nodeEnter
      .append('text')
      .attr('class', 'span-name')
      .attr('y', 10)
      .text((d: any) => d.data.data?.name);

    nodeEnter
      .append('text')
      .attr('class', 'span-duration')
      .attr('y', 30)
      .text((d: any) => `${d.data.data?.duration} ms`);

    const nodeUpdate = node.merge(nodeEnter as any);

    nodeUpdate
      .attr('transform', (d: any) => `translate(${0} ,${d.x})`)
      .select('rect')
      .attr(
        'class',
        d =>
          `rect-svg ${d.id === selectedSpan?.id ? 'rect-svg-selected' : ''} ${
            affectedSpans.includes(d.id) ? 'rect-svg-affected' : ''
          }`
      );

    nodeUpdate
      .select('path')
      .attr('opacity', d => (d.children || d._children ? 1 : 0))
      .attr('d', d =>
        d._children
          ? 'M18.629 15.997l-7.083-7.081L13.462 7l8.997 8.997L13.457 25l-1.916-1.916z'
          : 'M16.003 18.626l7.081-7.081L25 13.46l-8.997 8.998-9.003-9 1.917-1.916z'
      );

    nodeUpdate.selectAll('.grey-line').attr('transform', (d: any) => `translate(${d.y + 90} ,0)`);
    nodeUpdate.selectAll('text').attr('transform', (d: any) => `translate(${d.y + 20} ,0)`);

    node.exit().remove();

    nodesSort.forEach((d: any) => {
      d.x0 = d.x;
      d.y0 = d.y;
    });
  }, [treeFactory, root, onSelectSpan, selectedSpan?.id, scaleTime, minNano, maxNano]);

  useEffect(() => {
    drawChart();
  }, [selectedSpan?.id, drawChart]);

  return (
    <S.Container barHeight={barHeight} showAffected={affectedSpans.length > 0}>
      <svg ref={svgRef} />
    </S.Container>
  );
};
