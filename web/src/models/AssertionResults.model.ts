import {uniqueId} from 'lodash';
import {TAssertionResults, TRawAssertionResults} from '../types/Assertion.types';
import AssertionResult from './AssertionResult.model';

const AssertionResults = ({allPassed = false, results = []}: TRawAssertionResults): TAssertionResults => {
  return {
    allPassed,
    resultList: results.map(({selector = '', results: resultList = []}) => ({
      id: uniqueId(),
      selector,
      resultList: resultList.map(assertionResult => AssertionResult(assertionResult)),
    })),
  };
};

export default AssertionResults;
