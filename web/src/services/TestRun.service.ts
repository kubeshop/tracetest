import {TriggerTypes} from 'constants/Test.constants';
import TestRun from 'models/TestRun.model';
import {TSpanFlatAttribute} from 'types/Span.types';

const TestRunService = () => ({
  getResponseAttributeList({
    triggerResult: {statusCode = 0, body = '', headers = []} = {type: TriggerTypes.http, statusCode: 0},
  }: TestRun): TSpanFlatAttribute[] {
    const attributeList = [
      {key: 'body', value: body},
      {key: 'status_code', value: statusCode.toString()},
    ].concat(headers.map(({key, value}) => ({key: `headers.${key.toLowerCase()}`, value})));

    return attributeList;
  },

  getIsMissingVariablesError(data: any) {
    return !!data && data['missingVariables'];
  },
});

export default TestRunService();
