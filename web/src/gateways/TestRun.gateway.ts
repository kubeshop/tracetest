import TracetestAPI from 'redux/apis/Tracetest';
import {TRawTestSpecs} from 'models/TestSpecs.model';

const TestRunGateway = () => ({
  get(testId: string, take = 25, skip = 0) {
    return TracetestAPI.instance.endpoints.getRunList.initiate({testId, take, skip});
  },
  getById(testId: string, runId: number) {
    return TracetestAPI.instance.endpoints.getRunById.initiate({testId, runId});
  },
  reRun(testId: string, runId: number) {
    return TracetestAPI.instance.endpoints.reRun.initiate({testId, runId});
  },
  dryRun(testId: string, runId: number, testDefinition: Partial<TRawTestSpecs>) {
    return TracetestAPI.instance.endpoints.dryRun.initiate({testId, runId, testDefinition});
  },
  runTest(testId: string, variableSetId = '') {
    return TracetestAPI.instance.endpoints.runTest.initiate({testId, variableSetId});
  },
});

export default TestRunGateway();
