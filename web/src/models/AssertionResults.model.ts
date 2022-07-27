import {uniqueId} from 'lodash';
import AssertionService from '../services/Assertion.service';
import {TAssertionResults, TRawAssertionResults} from '../types/Assertion.types';
import AssertionResult from './AssertionResult.model';

const AssertionResults = ({allPassed = false, results = []}: TRawAssertionResults): TAssertionResults => {
  return {
    allPassed,
    resultList: results.map(({selector: {query: selector = ''} = {}, results: resultList = []}) => ({
      id: uniqueId(),
      selector,
      spanIds: AssertionService.getSpanIds(resultList),
      originalSelector: selector,
      resultList: resultList.map(assertionResult => AssertionResult(assertionResult)),
    })),
  };
};

export default AssertionResults;
