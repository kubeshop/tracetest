import {stratify} from '@visx/hierarchy';
import {NodeTypesEnum} from 'constants/Visualization.constants';
import {INodeDataSpan, TNode} from 'types/Timeline.types';

function getHierarchyNodes(nodesData: INodeDataSpan[]) {
  return stratify<INodeDataSpan>()
    .id(d => d.id)
    .parentId(d => d.parentId)(nodesData)
    .sort((a, b) => a.data.startTime - b.data.startTime);
}

const TimelineService = () => ({
  getNodes(nodesData: INodeDataSpan[], type: NodeTypesEnum) {
    const nodes: TNode[] = [];
    const hierarchyNodes = getHierarchyNodes(nodesData);

    hierarchyNodes.eachBefore(hierarchyNode => {
      nodes.push({
        children: hierarchyNode.children?.length ?? 0,
        data: {...hierarchyNode.data},
        depth: hierarchyNode.depth,
        type,
      });
    });

    return nodes;
  },

  getFilteredNodes(nodes: TNode[], collapsedNodesId: string[]) {
    const filteredNodes: TNode[] = [];

    nodes.forEach(node => {
      const parentId = node.data.parentId;
      const isParentPresent = filteredNodes.some(filteredNode => filteredNode.data.id === parentId);

      if (parentId && (collapsedNodesId.includes(parentId) || !isParentPresent)) {
        return;
      }

      filteredNodes.push(node);
    });

    return filteredNodes;
  },

  getMinMax(nodes: TNode[]) {
    const startTimes = nodes.map(node => node.data.startTime);
    const endTimes = nodes.map(node => node.data.endTime);
    return [Math.min(...startTimes), Math.max(...endTimes)];
  },
});

export default TimelineService();
