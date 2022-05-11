import {useCallback, useEffect, useMemo} from 'react';
import ReactFlow, {Background, FlowElement} from 'react-flow-renderer';
import {useDAGChart} from '../../../hooks/useDAGChart';
import TraceNode from '../../TraceNode';
import * as S from './DAG.styled';
import TraceDiagramAnalyticsService from '../../../services/Analytics/TraceDiagramAnalytics.service';
import {IDiagramProps} from '../Diagram';
import {ISpan} from '../../../types/Span.types';

export type TSpanInfo = {
  id: string;
  parentIds: string[];
  data: ISpan;
};

export type TSpanMap = Record<string, TSpanInfo>;

const {onClickSpan} = TraceDiagramAnalyticsService;

const Diagram: React.FC<IDiagramProps> = ({trace, selectedSpan, onSelectSpan}): JSX.Element => {
  const spanMap = useMemo<TSpanMap>(() => {
    return (
      trace?.spans?.reduce<TSpanMap>((acc, span) => {
        acc[span.spanId] = acc[span.spanId] || {id: span.spanId, parentIds: [], data: span};
        if (span.parentSpanId) acc[span.spanId].parentIds.push(span.parentSpanId);

        return acc;
      }, {}) || {}
    );
  }, [trace?.spans]);

  const dagLayout = useDAGChart(spanMap);

  const handleElementClick = useCallback(
    (event, {id}: FlowElement) => {
      onClickSpan(id);
      if (onSelectSpan) onSelectSpan(id);
    },
    [onSelectSpan]
  );

  useEffect(() => {
    if (dagLayout && dagLayout.dag) {
      const [dragNode] = dagLayout.dag.descendants();
      const span = spanMap[dragNode?.data.id];

      if (!selectedSpan && span && onSelectSpan) onSelectSpan(span.id);
    }
  }, [dagLayout, onSelectSpan, selectedSpan, spanMap]);

  const dagElements = useMemo(() => {
    if (dagLayout && dagLayout.dag) {
      const dagNodes = dagLayout.dag.descendants().map(({data, x, y}) => {
        const span = spanMap[data.id].data;

        return {
          id: data.id,
          type: 'TraceNode',
          data: span,
          position: {x, y: parseFloat(String(y))},
          selected: data.id === selectedSpan?.spanId,
          sourcePosition: 'top',
          className: `${data.id === selectedSpan?.spanId ? 'selected' : ''}`,
        };
      });

      dagLayout.dag.links().forEach(({source, target}: any) => {
        dagNodes.push({
          id: `${source.data.id}_${target.data.id}`,
          source: source.data.id,
          target: target.data.id,
          data: spanMap[source.data.id].data,
          labelShowBg: false,
          animated: false,
          arrowHeadType: 'arrowclosed',
        } as any);
      });

      return dagNodes;
    }

    return [];
  }, [dagLayout, spanMap, selectedSpan?.spanId]);

  return (
    <S.Container data-cy="diagram-dag">
      <ReactFlow
        nodeTypes={{TraceNode}}
        defaultZoom={0.5}
        elements={dagElements as any}
        onElementClick={handleElementClick}
        onLoad={({fitView}) => fitView()}
      >
        <Background gap={4} size={1} color="#FBFBFF" />
      </ReactFlow>
    </S.Container>
  );
};

export default Diagram;
