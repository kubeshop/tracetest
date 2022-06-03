import * as d3DAG from 'd3-dag';
import {useMemo} from 'react';

export interface IDAGNode<T> {
  data: T;
  id: string;
  parentIds: string[];
}

export const useDAGChart = <T>(nodeList: IDAGNode<T>[] = []) => {
  return useMemo(() => {
    if (!nodeList || !nodeList.length) return {};

    const stratify = d3DAG.dagStratify();
    const dag = stratify(nodeList);

    const dagLayout = d3DAG
      .sugiyama()
      .layering(d3DAG.layeringSimplex())
      .decross(d3DAG.decrossOpt())
      .coord(d3DAG.coordCenter())
      .nodeSize(() => [200, 150]);

    const {width, height} = dagLayout(dag as never);

    return {dag, layout: {width, height}};
  }, [nodeList]);
};
