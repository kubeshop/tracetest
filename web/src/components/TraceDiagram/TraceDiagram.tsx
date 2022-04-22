import {useCallback, useEffect, useMemo} from 'react';
import ReactFlow, {Background, BackgroundVariant, FlowElement} from 'react-flow-renderer';
import {useDAGChart} from 'hooks/Charts';
import TraceNode from './TraceNode';
import {TSpanInfo, TSpanMap} from '../Trace/Trace';
import {ITrace} from '../../types';
import * as S from './TraceDiagram.styled';

interface IPropsTraceDiagram {
  spanMap: TSpanMap;
  selectedSpan?: TSpanInfo;
  trace: ITrace;
  onSelectSpan(spanId: string): void;
}

const TraceDiagram = ({spanMap, trace, selectedSpan, onSelectSpan}: IPropsTraceDiagram): JSX.Element => {
  const dagLayout = useDAGChart(spanMap);

  const handleElementClick = useCallback(
    (event, {id}: FlowElement) => {
      onSelectSpan(id);
    },
    [onSelectSpan]
  );

  useEffect(() => {
    if (dagLayout && dagLayout.dag) {
      const [dragNode] = dagLayout.dag.descendants();
      const span = spanMap[dragNode?.data.id];

      if (!selectedSpan && span) {
        onSelectSpan(span.id);
      }
    }
  }, [dagLayout, onSelectSpan, selectedSpan, spanMap]);

  const dagElements = useMemo(() => {
    if (dagLayout && dagLayout.dag) {
      const dagNodes = dagLayout.dag.descendants().map(({data, x, y}) => {
        const span = spanMap[data.id].data;

        return {
          id: data.id,
          type: 'TraceNode',
          data: {span, trace},
          position: {x, y: parseFloat(String(y))},
          selected: data.id === selectedSpan?.id,
          sourcePosition: 'top',
          className: `${data.id === selectedSpan?.id ? 'selected' : ''}`,
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
  }, [dagLayout, spanMap, trace, selectedSpan?.id]);

  return (
    <S.Container style={{height: Math.max(dagLayout?.layout?.height || 0, 900) + 100}}>
      <ReactFlow
        nodeTypes={{TraceNode}}
        defaultZoom={0.5}
        elements={dagElements as any}
        onElementClick={handleElementClick}
      >
        <Background variant={BackgroundVariant.Lines} gap={4} size={1} />
      </ReactFlow>
    </S.Container>
  );
};

export default TraceDiagram;
