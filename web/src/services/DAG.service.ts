import {coordCenter, Dag, dagStratify, layeringSimplex, sugiyama} from 'd3-dag';
import {MarkerType} from 'react-flow-renderer';

import {theme} from 'constants/Theme.constants';
import {INodeDatum} from 'types/DAG.types';

function getDagLayout<T>(nodesDatum: INodeDatum<T>[]) {
  const stratify = dagStratify();
  const dag = stratify(nodesDatum);

  const dagLayout = sugiyama()
    .layering(layeringSimplex())
    .coord(coordCenter())
    .nodeSize(() => [220, 180]);

  dagLayout(dag as never);

  return dag;
}

function getNodes<T>(dagLayout: Dag<INodeDatum<T>, undefined>) {
  return dagLayout.descendants().map(({data: {id, data, type}, x, y}) => ({
    data,
    id,
    position: {x: x ?? 0, y: y ?? 0},
    type,
  }));
}

function getEdges<T>(dagLayout: Dag<INodeDatum<T>, undefined>) {
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
  getEdgesAndNodes<T>(nodesDatum: INodeDatum<T>[]) {
    if (!nodesDatum.length) return {edges: [], nodes: []};

    const dagLayout = getDagLayout(nodesDatum);
    const edges = getEdges(dagLayout);
    const nodes = getNodes(dagLayout);

    return {edges, nodes};
  },
});

export default DAGService();
