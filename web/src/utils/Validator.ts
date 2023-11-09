import {XMLValidator} from 'fast-xml-parser';
import isEmpty from 'lodash/isEmpty';
import {BodyMode} from '../components/Inputs/Body/useBodyMode';

type Value = any;

const Validator = {
  required(value: Value) {
    return !isEmpty(value);
  },

  xml(str: Value) {
    if (str === '') return true;
    const result = XMLValidator.validate(str);
    if (result === true) return true;
    return !result.err;
  },
  getBodyType(str?: Value): BodyMode {
    if (!str) return 'none';
    if (Validator.json(str)) {
      return 'json';
    }
    if (Validator.xml(str)) {
      return 'xml';
    }
    if (str !== '') {
      return 'raw';
    }
    return 'none';
  },
  json(str: Value) {
    if (str === '') return true;
    try {
      JSON.parse(str);
      return true;
    } catch (e) {
      return false;
    }
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
