import {TriggerTypes} from '../constants/Test.constants';
import {TSpanFlatAttribute} from '../types/Span.types';
import {TTestRun} from '../types/TestRun.types';

const TestRun = () => ({
  getResponseAttributeList({
    triggerResult: {statusCode = 0, body = '', headers = []} = {type: TriggerTypes.http, statusCode: 0},
  }: TTestRun): TSpanFlatAttribute[] {
    const attributeList = [
      {key: 'body', value: body},
      {key: 'status_code', value: statusCode.toString()},
    ].concat(headers.map(({key, value}) => ({key: `headers.${key.toLowerCase()}`, value})));

    return attributeList;
  },
});

export default TestRun();
