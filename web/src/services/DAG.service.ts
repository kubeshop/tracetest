import {coordCenter, Dag, dagStratify, layeringSimplex, sugiyama} from 'd3-dag';
import {Edge, MarkerType, Node} from 'react-flow-renderer';

import {theme} from 'constants/Theme.constants';
import {INodeDatum} from 'types/DAG.types';
import {withLowPriority} from '../utils/Common';

function getDagLayout<T>(nodesDatum: INodeDatum<T>[]): Dag<INodeDatum<T>, undefined> {
  const stratify = dagStratify();
  const dag = stratify(nodesDatum);

  const dagLayout = sugiyama()
    .layering(layeringSimplex())
    .coord(coordCenter())
    .nodeSize(() => [220, 180]);

  dagLayout(dag as never);

  return dag;
}

function getNodes<T>(dagLayout: Dag<INodeDatum<T>, undefined>): Node[] {
  return dagLayout.descendants().map(({data: {id, data, type}, x, y}) => ({
    data,
    id,
    position: {x: x ?? 0, y: y ?? 0},
    type,
  }));
}

function getEdges<T>(dagLayout: Dag<INodeDatum<T>, undefined>): Edge[] {
  return dagLayout.links().map(({source, target}) => ({
    animated: false,
    id: `${source.data.id}-${target.data.id}`,
    markerEnd: {color: theme.color.border, type: MarkerType.ArrowClosed},
    source: source.data.id,
    style: {stroke: theme.color.border},
    target: target.data.id,
  }));
}

const DAGService = () => ({
  async getEdgesAndNodes<T>(nodesDatum: INodeDatum<T>[]): Promise<{edges: Edge[]; nodes: Node[]}> {
    if (!nodesDatum.length) return {edges: [], nodes: []};

    const dagLayout = await withLowPriority(getDagLayout)(nodesDatum);
    const edges = await withLowPriority(getEdges)(dagLayout);
    const nodes = await withLowPriority(getNodes)(dagLayout);

    return {edges, nodes};
  },
});

export default DAGService();
