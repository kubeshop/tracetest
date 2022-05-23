import {endpoints} from '../redux/apis/TraceTest.api';
import {TRawTestDefinition} from '../types/TestDefinition.types';

const {setTestDefinition, getTestDefinition} = endpoints;

const AssertionGateway = () => ({
  get(testId: string) {
    return getTestDefinition.initiate({testId});
  },
  set(testId: string, testDefinition: Partial<TRawTestDefinition>) {
    return setTestDefinition.initiate({testId, testDefinition});
  },
});

export default AssertionGateway();
