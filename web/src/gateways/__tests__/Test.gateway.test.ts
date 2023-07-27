import {TRawTestResource} from '../../models/Test.model';
import {endpoints} from '../../redux/apis/Tracetest';
import TestGateway from '../Test.gateway';

const {createTest, getTestById, getTestList, runTest} = endpoints;

jest.mock('../../redux/apis/Tracetest', () => {
  const initiate = jest.fn(() => Promise.resolve());

  return {
    endpoints: {
      createTest: {initiate},
      getTestById: {initiate},
      getTestList: {initiate},
      runTest: {initiate},
    },
  };
});

describe('TestGateway', () => {
  it('should execute the create function', async () => {
    expect.assertions(1);
    const test: TRawTestResource = {type: 'Test', spec: {name: 'test', description: 'test'}};
    await TestGateway.create(test);

    expect(createTest.initiate).toBeCalledWith(test);
  });

  it('should execute the getById function', async () => {
    expect.assertions(1);
    await TestGateway.getById('testId');

    expect(getTestById.initiate).toBeCalledWith({testId: 'testId'});
  });

  it('should execute the getList function', async () => {
    expect.assertions(1);
    await TestGateway.getList();

    expect(getTestList.initiate).toBeCalledWith({});
  });

  it('should execute the runTest function', async () => {
    expect.assertions(1);
    await TestGateway.run('testId');

    expect(runTest.initiate).toBeCalledWith({testId: 'testId'});
  });
});
