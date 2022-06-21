import {dagStratify, Dag, sugiyama, layeringSimplex, decrossOpt, coordCenter} from 'd3-dag';
// import {ArrowHeadType, Elements, Position} from 'react-flow-renderer';
import {MarkerType} from 'react-flow-renderer';
import {NodeTypesEnum} from 'constants/Diagram.constants';

export interface IDAGNode<T> {
  data: T;
  id: string;
  parentIds: string[];
  type: NodeTypesEnum;
  className?: string;
}

// export type TElementList = Elements<unknown>;
export type TNode<T = unknown> = IDAGNode<T>;

function getDagLayout(nodeList: TNode[]) {
  const stratify = dagStratify();
  const dag = stratify(nodeList);

  const dagLayout = sugiyama()
    .layering(layeringSimplex())
    .decross(decrossOpt())
    .coord(coordCenter())
    .nodeSize(() => [200, 150]);

  dagLayout(dag as never);

  // console.log('### descendants', dag.descendants());
  // console.log('### links', dag.links());
  return dag;
}

function getNodes(dag: Dag<TNode, undefined>) {
  const nodes: any = dag.descendants().map(({data: {id, data, type, className}, x, y}) => {
    return {
      id,
      type,
      data,
      position: {x, y},
      // sourcePosition: Position.Top,
      className,
    };
  });
  return nodes;
}

function getEdges(dag: Dag<TNode, undefined>) {
  const edges: any = dag.links().map(({source, target}) => {
    return {
      id: `${source.data.id}-${target.data.id}`,
      source: source.data.id,
      target: target.data.id,
      animated: true, // for animated edges
      markerEnd: {type: MarkerType.ArrowClosed}, // arrow at the end of the edge
    };
  });
  return edges;
}

const DAGService = () => ({
  getNodesAndEdges(nodeList: TNode[]) {
    console.log('### running getNodesAndEdges');
    const dag = getDagLayout(nodeList);
    const nodes = getNodes(dag);
    const edges = getEdges(dag);
    // const elementList = this.getDagElementList(dag);

    return {nodes, edges};
  },
});

export default DAGService();
