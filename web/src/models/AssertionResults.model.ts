import {uniqueId} from 'lodash';
import AssertionService from 'services/Assertion.service';
import {Model, TTestSchemas} from 'types/Common.types';
import AssertionResult from './AssertionResult.model';

export type TRawAssertionResults = TTestSchemas['AssertionResults'];
export type TAssertionResultEntry = {
  id: string;
  selector: string;
  originalSelector: string;
  spanIds: string[];
  resultList: AssertionResult[];
};
type AssertionResults = Model<
  TRawAssertionResults,
  {
    allPassed: boolean;
    results?: never;
    resultList: TAssertionResultEntry[];
  }
>;

const AssertionResults = ({allPassed = false, results = []}: TRawAssertionResults): AssertionResults => {
  return {
    allPassed,
    resultList: results.map(({selector: {query: selector = ''} = {}, results: resultList = []}) => ({
      id: uniqueId(),
      selector,
      spanIds: AssertionService.getSpanIds(resultList),
      originalSelector: selector,
      resultList: resultList.map(assertionResult => AssertionResult(assertionResult)),
      name: '',
    })),
  };
};

export default AssertionResults;
