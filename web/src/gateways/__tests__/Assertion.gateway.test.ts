import {endpoints} from '../../redux/apis/Test.api';
import AssertionGateway from '../Assertion.gateway';

const {createAssertion, getAssertions, updateAssertion} = endpoints;

jest.mock('../../redux/apis/Test.api', () => {
  const initiate = jest.fn(() => Promise.resolve());

  return {
    endpoints: {
      createAssertion: {initiate},
      getAssertions: {initiate},
      updateAssertion: {initiate},
    },
  };
});

describe('AssertionGateway', () => {
  it('should execute the getAssertions function', async () => {
    expect.assertions(1);
    await AssertionGateway.get('testId');

    expect(getAssertions.initiate).toBeCalledWith('testId');
  });

  it('should execute the createAssertion function', async () => {
    expect.assertions(1);
    const assertion = {assertionId: 'assertionId', selectors: []};
    await AssertionGateway.create('testId', assertion);

    expect(createAssertion.initiate).toBeCalledWith({testId: 'testId', assertion});
  });

  it('should execute the updateAssertion function', async () => {
    expect.assertions(1);
    const assertion = {assertionId: 'assertionId', selectors: []};
    await AssertionGateway.update('testId', 'assertionId', assertion);

    expect(updateAssertion.initiate).toBeCalledWith({testId: 'testId', assertionId: 'assertionId', assertion});
  });
});
