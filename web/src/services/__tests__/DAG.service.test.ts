import {skeletonNodesDatum} from 'constants/DAG.constants';
import DAGService from '../DAG.service';

describe('DAGService', () => {
  describe('getNodesAndEdges', () => {
    it('should return DAG edges and nodes', () => {
      const {edges, nodes} = DAGService.getEdgesAndNodes(skeletonNodesDatum);

      expect(Array.isArray(edges)).toBeTruthy();
      expect(Array.isArray(nodes)).toBeTruthy();
    });
  });
});
