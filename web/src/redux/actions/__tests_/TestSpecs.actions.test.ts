import TestSpecsActions from '../TestSpecs.actions';
import {store} from '../../store';
import {HTTP_METHOD} from '../../../constants/Common.constants';
import TestSpecsSelectors from '../../../selectors/TestSpecs.selectors';
import TestRunMock from '../../../models/__mocks__/TestRun.mock';
import TestMock from '../../../models/__mocks__/Test.mock';

jest.mock('../../../selectors/TestSpecs.selectors', () => ({
  selectSpecs: jest.fn(),
}));

const selectTestMock = TestSpecsSelectors.selectSpecs as unknown as jest.Mock;

describe('TestDefinitionActions', () => {
  beforeEach(() => {
    fetchMock.resetMocks();
  });

  describe('publish', () => {
    it('should trigger the set definition and rerun requests', async () => {
      selectTestMock.mockImplementationOnce(() => []);

      fetchMock.mockResponseOnce(JSON.stringify({}));
      fetchMock.mockResponseOnce(JSON.stringify(TestRunMock.raw()));

      await store.dispatch(
        TestSpecsActions.publish({
          test: TestMock.model(),
          testId: 'testId',
          runId: 1,
        })
      );

      const setRequest = fetchMock.mock.calls[0][0] as Request;
      const reRunRequest = fetchMock.mock.calls[1][0] as Request;

      expect(setRequest.url).toEqual('http://localhost/api/tests/testId');
      expect(setRequest.method).toEqual(HTTP_METHOD.PUT);

      expect(reRunRequest.url).toEqual('http://localhost/api/tests/testId/run/runId/rerun');
      expect(reRunRequest.method).toEqual(HTTP_METHOD.POST);

      expect(fetchMock.mock.calls.length).toBe(2);
    });
  });

  describe('dry run', () => {
    it('should trigger the dry run request', async () => {
      selectTestMock.mockImplementationOnce(() => []);

      fetchMock.mockResponseOnce(JSON.stringify(TestRunMock.raw()));
      await store.dispatch(
        TestSpecsActions.dryRun({
          testId: 'testId',
          runId: 1,
          definitionList: [],
        })
      );

      const dryRunRequest = fetchMock.mock.calls[0][0] as Request;

      expect(dryRunRequest.url).toEqual('http://localhost/api/tests/testId/run/runId/dry-run');
      expect(dryRunRequest.method).toEqual(HTTP_METHOD.PUT);

      expect(fetchMock.mock.calls.length).toBe(1);
    });
  });
});
