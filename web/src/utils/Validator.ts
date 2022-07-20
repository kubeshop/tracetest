import isEmpty from 'lodash/isEmpty';

type Value = any;

const Validator = {
  required(value: Value) {
    return !isEmpty(value);
  },
  url(value: Value) {
    try {
      if (typeof value === 'string') {
        if (value.length <= 2048) {
          const innerUrl = new URL(value);
          return innerUrl.protocol === 'http:' || innerUrl.protocol === 'https:';
        }
      }
      return false;
    } catch {
      return false;
    }
  },
};

export default Validator;
