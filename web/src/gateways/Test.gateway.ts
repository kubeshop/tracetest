import {endpoints} from 'redux/apis/TraceTest.api';
import {TRawTest} from 'models/Test.model';

const {createTest, editTest, getTestById, getTestList, runTest} = endpoints;

const TestGateway = () => ({
  getList() {
    return getTestList.initiate({});
  },
  getById(testId: string) {
    return getTestById.initiate({testId});
  },
  create(test: TRawTest) {
    return createTest.initiate(test);
  },
  run(testId: string) {
    return runTest.initiate({testId});
  },
  edit(test: TRawTest, testId: string) {
    return editTest.initiate({test, testId});
  },
});

export default TestGateway();
