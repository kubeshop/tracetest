import {dagStratify, Dag, sugiyama, layeringSimplex, decrossOpt, coordCenter} from 'd3-dag';
import {ArrowHeadType, Elements, Position} from 'react-flow-renderer';
import {strokeColor, TraceNodes} from 'constants/Diagram.constants';

export interface IDAGNode<T> {
  data: T;
  id: string;
  parentIds: string[];
  type: TraceNodes;
}

export type TElementList = Elements<unknown>;
export type TNode<T = unknown> = IDAGNode<T>;

const DAGService = () => ({
  getDagLayout(nodeList: TNode[]) {
    const stratify = dagStratify();
    const dag = stratify(nodeList);

    const dagLayout = sugiyama()
      .layering(layeringSimplex())
      .decross(decrossOpt())
      .coord(coordCenter())
      .nodeSize(() => [200, 150]);

    dagLayout(dag as never);

    return dag;
  },
  getDagElementList(dag: Dag<TNode, undefined>) {
    const dagNodeList: TElementList = dag.descendants().map(({data: {id, data, type}, x, y}) => {
      return {
        id,
        type,
        data,
        position: {x: x!, y: parseFloat(String(y))},
        sourcePosition: Position.Top,
      };
    });

    dag.links().forEach(({source, target}) => {
      dagNodeList.push({
        id: `${source.data.id}_${target.data.id}`,
        source: source.data.id,
        target: target.data.id,
        data: source.data.data,
        labelShowBg: false,
        animated: false,
        arrowHeadType: ArrowHeadType.ArrowClosed,
        style: {stroke: strokeColor},
      });
    });

    return dagNodeList;
  },
  getDagInfo(nodeList: TNode[]) {
    const dag = this.getDagLayout(nodeList);
    const elementList = this.getDagElementList(dag);

    return {dag, elementList};
  },
});

export default DAGService();
