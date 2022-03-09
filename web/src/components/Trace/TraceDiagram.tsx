import Text from 'antd/lib/typography/Text';
import {useMemo} from 'react';
import ReactFlow, {Background, BackgroundVariant, Handle, NodeProps, Position} from 'react-flow-renderer';
import {useDAGChart} from 'hooks/Charts';
import {ISpan} from 'types';
import * as S from './TraceDiagram.styled';

interface IPropsTraceNode extends NodeProps<ISpan> {}

interface IPropsTraceDiagram {
  spanMap: any;
  selectedSpan: any;
  onSelectSpan: (span: any) => void;
}

const TraceNode = ({id, data, selected, ...rest}: IPropsTraceNode) => {
  const systemTag = data?.attributes?.find(el => {
    if (el.key.startsWith('http')) {
      return el;
    }
    if (el.key.startsWith('db.system')) {
      return el;
    }
    if (el.key.startsWith('rpc.system')) {
      return el;
    }
    if (el.key.startsWith('messaging.system')) {
      return el;
    }
    return false;
  });

  return (
    <S.TraceNode selected={selected}>
      <S.TraceNotch system={systemTag?.key || ''}>
        <Text>{systemTag?.value?.stringValue || 'Service'}</Text>
      </S.TraceNotch>
      <Handle type="target" id={id} position={Position.Top} style={{top: 0, borderRadius: 0, visibility: 'hidden'}} />

      <Text>{`/${data?.name?.split('/')?.pop()}`}</Text>
      <Handle
        type="source"
        position={Position.Bottom}
        id={id}
        style={{bottom: 0, borderRadius: 0, visibility: 'hidden'}}
      />
    </S.TraceNode>
  );
};

const TraceDiagram = ({spanMap, selectedSpan, onSelectSpan}: IPropsTraceDiagram): JSX.Element => {
  const {
    dag,
    layout: {height},
  } = useDAGChart(spanMap);

  const handleElementClick = (event: any, element: any) => {
    onSelectSpan(spanMap[element.id]);
  };

  const dagElements = useMemo(() => {
    const dagNodes = dag.descendants().map((i: any) => {
      return {
        id: i.data.id,
        type: 'TraceNode',
        data: spanMap[i.data.id].data,
        position: {x: i.x, y: parseFloat(i.y)},
        sourcePosition: 'top',
        className: `${i.data.id === selectedSpan.id ? 'selected' : ''}`,
      };
    });

    dag.links().forEach(({source, target}: any) => {
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
  }, [spanMap, dag, selectedSpan]);

  return (
    <div style={{height: height + 100}}>
      <ReactFlow
        nodeTypes={{TraceNode}}
        defaultZoom={0.5}
        elements={dagElements as any}
        onElementClick={handleElementClick}
      >
        <Background variant={BackgroundVariant.Lines} gap={4} size={1} />
      </ReactFlow>
    </div>
  );
};

export default TraceDiagram;
