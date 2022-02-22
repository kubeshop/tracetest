import * as d3DAG from 'd3-dag';
import data from '../pages/Trace/data.json';

interface ISpanMap {
  [key: string]: {id: string; parentIds: string[]};
}

const rootNode = data.resourceSpans
  .map((i: any) => i.instrumentationLibrarySpans.map((el: any) => el.spans))
  .flat(2)
  .sort((el1, el2) => Number(el1.startTimeUnixNano) - Number(el2.startTimeUnixNano))
  .shift();

export const useDAGChart = (spanMap: ISpanMap) => {
  const dagData = Object.values(spanMap).map(({id, parentIds}) => {
    if (spanMap[parentIds[0]] === undefined && id !== rootNode.spanId) {
      parentIds = [rootNode.spanId];
    }
    if (id === rootNode.spanId) {
      return {id, parentIds: undefined};
    }
    return {id, parentIds};
  });

  const stratify = d3DAG.dagStratify();
  const dag = stratify(dagData);

  const layout = d3DAG
    .sugiyama()
    .layering(d3DAG.layeringSimplex())
    .decross(d3DAG.decrossOpt())
    .coord(d3DAG.coordCenter())
    .nodeSize(() => [300, 150]);

  const {width, height} = layout(dag as any);

  return {dag, layout: {width, height}};
};
