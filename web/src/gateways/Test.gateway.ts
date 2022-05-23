import {endpoints} from '../redux/apis/TraceTest.api';
import {TTest} from '../types/Test.types';

const {createTest, getTestById, getTestList, runTest} = endpoints;

const TestGateway = () => ({
  getList() {
    return getTestList.initiate();
  },
  getById(testId: string) {
    return getTestById.initiate({testId});
  },
  create(test: Partial<TTest>) {
    return createTest.initiate(test);
  },
  run(testId: string) {
    return runTest.initiate({testId});
  },
});

export default TestGateway();
