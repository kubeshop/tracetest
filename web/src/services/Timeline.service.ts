import {stratify} from '@visx/hierarchy';
import {INodeDataSpan, TNode} from 'types/Timeline.types';

function getHierarchyNodes(nodesData: INodeDataSpan[]) {
  return stratify<INodeDataSpan>()
    .id(d => d.id)
    .parentId(d => d.parentId)(nodesData)
    .sort((a, b) => a.data.startTime - b.data.startTime);
}

const TimelineService = () => ({
  getNodes(nodesData: INodeDataSpan[]) {
    const nodes: TNode[] = [];
    const hierarchyNodes = getHierarchyNodes(nodesData);

    hierarchyNodes.eachBefore(hierarchyNode => {
      nodes.push({
        children: hierarchyNode.children?.length ?? 0,
        data: {...hierarchyNode.data},
        depth: hierarchyNode.depth,
      });
    });

    return nodes;
  },

  getMinMax(nodes: TNode[]) {
    const startTimes = nodes.map(node => node.data.startTime);
    const endTimes = nodes.map(node => node.data.endTime);
    return [Math.min(...startTimes), Math.max(...endTimes)];
  },
});

export default TimelineService();
