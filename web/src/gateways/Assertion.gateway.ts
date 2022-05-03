import {endpoints} from '../redux/apis/Test.api';
import {TAssertion} from '../types/Assertion.types';

const {createAssertion, getAssertions, updateAssertion} = endpoints;

const AssertionGateway = () => ({
  get(testId: string) {
    return getAssertions.initiate(testId);
  },
  create(testId: string, assertion: Partial<TAssertion>) {
    return createAssertion.initiate({testId, assertion});
  },
  update(testId: string, assertionId: string, assertion: Partial<TAssertion>) {
    return updateAssertion.initiate({testId, assertion, assertionId});
  },
});

export default AssertionGateway();
