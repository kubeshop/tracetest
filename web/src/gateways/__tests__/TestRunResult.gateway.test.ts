import {endpoints} from '../../redux/apis/Test.api';
import TestRunResultGateway from '../TestRunResult.gateway';

const {getResultById, getResultList, updateResult} = endpoints;

jest.mock('../../redux/apis/Test.api', () => {
  const initiate = jest.fn(() => Promise.resolve());

  return {
    endpoints: {
      getResultById: {initiate},
      getResultList: {initiate},
      updateResult: {initiate},
    },
  };
});

describe('TestRunResultGateway', () => {
  it('should execute the get function', async () => {
    expect.assertions(1);
    await TestRunResultGateway.get('testId');

    expect(getResultList.initiate).toBeCalledWith({testId: 'testId', take: 25, skip: 0});
  });

  it('should execute the getById function', async () => {
    expect.assertions(1);
    await TestRunResultGateway.getById('testId', 'resultId');

    expect(getResultById.initiate).toBeCalledWith({testId: 'testId', resultId: 'resultId'});
  });

  it('should execute the update function', async () => {
    expect.assertions(1);
    const testAssertionResult = {
      assertionResultState: true,
      assertionResult: [],
    };
    await TestRunResultGateway.update('testId', 'resultId', testAssertionResult);

    expect(updateResult.initiate).toBeCalledWith({
      testId: 'testId',
      resultId: 'resultId',
      assertionResult: testAssertionResult,
    });
  });
});
