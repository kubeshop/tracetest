import {endpoints} from '../../redux/apis/TraceTest.api';
import {TRawTestDefinition} from '../../types/TestDefinition.types';
import TestDefinitionGateway from '../TestDefinition.gateway';

const {setTestDefinition, getTestDefinition} = endpoints;

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
  it('should execute the getAssertions function', async () => {
    expect.assertions(1);
    await TestDefinitionGateway.get('testId');

    expect(getTestDefinition.initiate).toBeCalledWith({testId: 'testId'});
  });

  it('should execute the createAssertion function', async () => {
    expect.assertions(1);
    const testDefinition: TRawTestDefinition = {specs: []};
    await TestDefinitionGateway.set('testId', testDefinition);

    expect(setTestDefinition.initiate).toBeCalledWith({testId: 'testId', testDefinition});
  });
});
