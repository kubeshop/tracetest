import {useMemo} from 'react';
import ReactFlow, {Background, BackgroundVariant} from 'react-flow-renderer';
import {useDAGChart} from '../../hooks/Charts';

interface IProps {
  spanMap: any;
  selectedSpan: any;
  onSelectSpan: (span: any) => void;
}

const TraceDiagram = ({spanMap, selectedSpan, onSelectSpan}: IProps): JSX.Element => {
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
        type: 'input',
        data: {label: <p style={{wordWrap: 'break-word'}}>{spanMap[i.data.id].data.operationName}</p>},
        position: {x: i.x, y: i.y},
        sourcePosition: 'top',
        className: `${i.data.id === selectedSpan.id ? 'selected' : ''}`,
      };
    });

    dag.links().forEach(({source, target}: any) => {
      dagNodes.push({
        id: `${source.data.id}_${target.data.id}`,
        source: source.data.id,
        target: target.data.id,
        animated: true,
      } as any);
    });
    return dagNodes;
  }, [spanMap, dag, selectedSpan]);

  return (
    <div style={{height: height + 100}}>
      <ReactFlow elements={dagElements as any} onElementClick={handleElementClick}>
        <Background variant={BackgroundVariant.Lines} gap={4} size={1} />
      </ReactFlow>
    </div>
  );
};

export default TraceDiagram;
