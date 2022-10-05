import {IKeyValue} from '../../constants/Test.constants';
import {IMockFactory} from '../../types/Common.types';
import KeyValue from '../KeyValue.model';

const KeyValueMock: IMockFactory<IKeyValue, IKeyValue> = () => ({
  raw(data = {}) {
    return {
      key: 'token',
      value:
        'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c',
      ...data,
    };
  },
  model(data = {}) {
    return KeyValue(this.raw(data));
  },
});

export default KeyValueMock();
