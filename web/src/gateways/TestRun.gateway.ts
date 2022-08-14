import {endpoints} from 'redux/apis/TraceTest.api';
import {TRawTestDefinition} from 'types/TestDefinition.types';

const {getRunList, getRunById, reRun, dryRun, runTest} = endpoints;

const TestRunGateway = () => ({
  get(testId: string, take = 25, skip = 0, query = '') {
    return getRunList.initiate({testId, take, skip, query});
  },
  getById(testId: string, runId: string) {
    return getRunById.initiate({testId, runId});
  },
  reRun(testId: string, runId: string) {
    return reRun.initiate({testId, runId});
  },
  dryRun(testId: string, runId: string, testDefinition: Partial<TRawTestDefinition>) {
    return dryRun.initiate({testId, runId, testDefinition});
  },
  runTest(testId: string) {
    return runTest.initiate({testId});
  },
});

export default TestRunGateway();
