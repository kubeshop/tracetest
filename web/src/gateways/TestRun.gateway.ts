import {endpoints} from 'redux/apis/Tracetest';
import {TRawTestSpecs} from 'models/TestSpecs.model';

const {getRunList, getRunById, reRun, dryRun, runTest} = endpoints;

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
  dryRun(testId: string, runId: string, testDefinition: Partial<TRawTestSpecs>) {
    return dryRun.initiate({testId, runId, testDefinition});
  },
  runTest(testId: string, variableSetId = '') {
    return runTest.initiate({testId, variableSetId});
  },
});

export default TestRunGateway();
