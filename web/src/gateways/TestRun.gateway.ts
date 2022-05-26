import {endpoints} from '../redux/apis/TraceTest.api';
import { TRawTestDefinition } from '../types/TestDefinition.types';

const {getRunList, getRunById, reRun, dryRun} = endpoints;

const TestRunGateway = () => ({
  get(testId: string, take = 25, skip = 0) {
    return getRunList.initiate({testId, take, skip});
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
});

export default TestRunGateway();
