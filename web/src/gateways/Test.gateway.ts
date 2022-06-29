import {endpoints} from 'redux/apis/TraceTest.api';
import {TDraftTest} from 'types/Plugins.types';

const {createTest, getTestById, getTestList, runTest} = endpoints;

const TestGateway = () => ({
  getList() {
    return getTestList.initiate();
  },
  getById(testId: string) {
    return getTestById.initiate({testId});
  },
  create(test: TDraftTest) {
    return createTest.initiate(test);
  },
  run(testId: string) {
    return runTest.initiate({testId});
  },
});

export default TestGateway();
