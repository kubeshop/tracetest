import TestDefinitionActions from '../TestDefinition.actions';
import TestSelectors from '../../../selectors/Test.selectors';
import {store} from '../../store';
import TestMock from '../../../models/__mocks__/Test.mock';
import {HTTP_METHOD} from '../../../constants/Common.constants';

jest.mock('../../../selectors/Test.selectors', () => ({
  selectTest: jest.fn(),
}));

const selectTestMock = TestSelectors.selectTest as unknown as jest.Mock;

describe('TestDefinitionActions', () => {
  beforeEach(() => {
    fetchMock.resetMocks();
  });

  describe('add', () => {
    it('should add a new definition by triggering the request to the backend', async () => {
      selectTestMock.mockImplementationOnce(() => TestMock.model());

      fetchMock.mockResponse(JSON.stringify({}));
      await store.dispatch(
        TestDefinitionActions.add({testId: 'testId', definition: {selector: 'selector', assertionList: []}})
      );

      const request = fetchMock.mock.calls[0][0] as Request;

      expect(request.url).toEqual('http://localhost/api/tests/testId/definition');
      expect(request.method).toEqual(HTTP_METHOD.PUT);
      expect(fetchMock.mock.calls.length).toBe(1);
    });
  });

  describe('update', () => {
    it('should update a definition by triggering the request to the backend', async () => {
      selectTestMock.mockImplementationOnce(() => TestMock.model());

      fetchMock.mockResponse(JSON.stringify({}));
      await store.dispatch(
        TestDefinitionActions.add({testId: 'testId', definition: {selector: 'selector', assertionList: []}})
      );

      const request = fetchMock.mock.calls[0][0] as Request;

      expect(request.url).toEqual('http://localhost/api/tests/testId/definition');
      expect(request.method).toEqual(HTTP_METHOD.PUT);
      expect(fetchMock.mock.calls.length).toBe(1);
    });
  });

  describe('remove', () => {
    it('should remove a definition by triggering the request to the backend', async () => {
      selectTestMock.mockImplementationOnce(() => TestMock.model());

      fetchMock.mockResponse(JSON.stringify({}));
      await store.dispatch(TestDefinitionActions.remove({testId: 'testId', selector: 'selector'}));

      const request = fetchMock.mock.calls[0][0] as Request;

      expect(request.url).toEqual('http://localhost/api/tests/testId/definition');
      expect(request.method).toEqual(HTTP_METHOD.PUT);
      expect(fetchMock.mock.calls.length).toBe(1);
    });
  });
});
