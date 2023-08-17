import TracetestAPI from 'redux/apis/Tracetest';
import {TRawTestResource} from 'models/Test.model';

const TestGateway = () => ({
  getList() {
    return TracetestAPI.instance.endpoints.getTestList.initiate({});
  },
  getById(testId: string) {
    return TracetestAPI.instance.endpoints.getTestById.initiate({testId});
  },
  create(test: TRawTestResource) {
    return TracetestAPI.instance.endpoints.createTest.initiate(test);
  },
  run(testId: string) {
    return TracetestAPI.instance.endpoints.runTest.initiate({testId});
  },
  edit(test: TRawTestResource, testId: string) {
    return TracetestAPI.instance.endpoints.editTest.initiate({test, testId});
  },
});

export default TestGateway();
