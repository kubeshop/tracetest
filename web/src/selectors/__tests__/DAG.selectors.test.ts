import DagSelectors from 'selectors/DAG.selectors';
import {IDagState} from '../../redux/slices/DAG.slice';
import {RootState} from '../../redux/store';

describe('DAGSelectors', () => {
  describe('selectEdges', () => {
    it('should return empty', () => {
      const result = DagSelectors.selectEdges({
        dag: {nodes: [], edges: []} as IDagState,
      } as RootState);
      expect(result).toStrictEqual([]);
    });
  });
  describe('selectNodes', () => {
    it('should return empty', () => {
      const result = DagSelectors.selectNodes({
        dag: {nodes: [], edges: []} as IDagState,
      } as RootState);
      expect(result).toStrictEqual([]);
    });
  });
});
