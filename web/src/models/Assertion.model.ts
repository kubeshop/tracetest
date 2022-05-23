import {TRawAssertion, TAssertion} from '../types/Assertion.types';

const Assertion = ({attribute = '', comparator = '', expected = ''}: TRawAssertion): TAssertion => {
  return {
    attribute,
    comparator,
    expected,
  };
};

export default Assertion;
