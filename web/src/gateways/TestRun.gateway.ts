import {endpoints} from '../redux/apis/TraceTest.api';

const {getRunList, getRunById, reRun} = endpoints;

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
});

export default TestRunGateway();
