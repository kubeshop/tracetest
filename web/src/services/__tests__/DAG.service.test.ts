import {skeletonNodeList} from 'constants/Diagram.constants';
import DAGService from '../DAG.service';

describe('DAGService', () => {
  describe('getDagInfo', () => {
    it('should return a dag instance and a list of elements', () => {
      const {dag, elementList} = DAGService.getDagInfo(skeletonNodeList);

      expect(Array.isArray(elementList)).toBeTruthy();
      expect(dag).toBeTruthy();
    });
  });
});
