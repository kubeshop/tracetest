import * as d3DAG from 'd3-dag';
import {Dag} from 'd3-dag';
import _ from 'lodash';
import {TSpanMap} from '../components/Diagram/components/DAG';

export const useDAGChart = (
  spanMap: TSpanMap = {}
): void | {
  dag: Dag<{id: string; parentIds: string[]}, undefined>;
  layout: {width: number; height: number};
} => {
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
