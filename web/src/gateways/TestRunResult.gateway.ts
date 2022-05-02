import {endpoints} from '../redux/apis/Test.api';
import {ITestAssertionResult} from '../types/Assertion.types';

const {getResultById, getResultList, updateResult} = endpoints;

const TestRunResultGateway = () => ({
  get(testId: string) {
    return getResultList.initiate(testId);
  },
  getById(testId: string, resultId: string) {
    return getResultById.initiate({testId, resultId});
  },
  update(testId: string, resultId: string, assertionResult: ITestAssertionResult) {
    return updateResult.initiate({testId, resultId, assertionResult});
  },
});

export default TestRunResultGateway();
