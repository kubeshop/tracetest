import TracetestAPI from 'redux/apis/Tracetest';
import {TRawTestResource} from 'models/Test.model';

const {endpoints} = TracetestAPI.instance;

const {createTest, editTest, getTestById, getTestList, runTest} = endpoints;

const TestGateway = () => ({
  getList() {
    return getTestList.initiate({});
  },
  getById(testId: string) {
    return getTestById.initiate({testId});
  },
  create(test: TRawTestResource) {
    return createTest.initiate(test);
  },
  run(testId: string) {
    return runTest.initiate({testId});
  },
  edit(test: TRawTestResource, testId: string) {
    return editTest.initiate({test, testId});
  },
});

export default TestGateway();
