import {endpoints} from '../../redux/apis/TraceTest.api';
import {TRawTestSpecs} from '../../types/TestSpecs.types';
import TestDefinitionGateway from '../TestDefinition.gateway';

const {setTestDefinition} = endpoints;

jest.mock('../../redux/apis/TraceTest.api', () => {
  const initiate = jest.fn(() => Promise.resolve());

  return {
    endpoints: {
      getTestDefinition: {initiate},
      setTestDefinition: {initiate},
    },
  };
});

describe('TestDefinitionGateway', () => {
  it('should execute the createAssertion function', async () => {
    expect.assertions(1);
    const testDefinition: TRawTestSpecs = {specs: []};
    await TestDefinitionGateway.set('testId', testDefinition);

    expect(setTestDefinition.initiate).toBeCalledWith({testId: 'testId', testDefinition});
  });
});
