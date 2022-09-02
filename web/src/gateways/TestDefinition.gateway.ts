import {endpoints} from 'redux/apis/TraceTest.api';
import {TRawTestSpecs} from 'types/TestSpecs.types';

const {setTestDefinition} = endpoints;

const TestDefinitionGateway = () => ({
  set(testId: string, testDefinition: Partial<TRawTestSpecs>) {
    return setTestDefinition.initiate({testId, testDefinition});
  },
});

export default TestDefinitionGateway();
