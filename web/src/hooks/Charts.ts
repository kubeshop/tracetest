import * as d3DAG from 'd3-dag';
import _ from 'lodash';

interface ISpanMap {
  [key: string]: {id: string; parentIds: string[]};
}

export const useDAGChart = (spanMap: ISpanMap = {}) => {
  if (_.isEmpty(spanMap)) {
    return;
  }

  const dagData = Object.values(spanMap).map(({id, parentIds}) => ({id, parentIds: parentIds.filter(el => el)}));
  const stratify = d3DAG.dagStratify();
  const dag = stratify(dagData);

  const layout = d3DAG
    .sugiyama()
    .layering(d3DAG.layeringSimplex())
    .decross(d3DAG.decrossOpt())
    .coord(d3DAG.coordCenter())
    .nodeSize(() => [200, 150]);

  const {width, height} = layout(dag as any);

  return {dag, layout: {width, height}};
};
