import {uniqueId} from 'lodash';
import AssertionService from '../services/Assertion.service';
import SelectorService from '../services/Selector.service';
import {TAssertionResults, TRawAssertionResults} from '../types/Assertion.types';
import AssertionResult from './AssertionResult.model';

const AssertionResults = ({allPassed = false, results = []}: TRawAssertionResults): TAssertionResults => {
  return {
    allPassed,
    resultList: results.map(({selector = '', results: resultList = []}) => ({
      id: uniqueId(),
      selector,
      spanCount: AssertionService.getSpanCount(resultList),
      pseudoSelector: SelectorService.getPseudoSelector(selector),
      selectorList: SelectorService.getSpanSelectorList(selector),
      resultList: resultList.map(assertionResult => AssertionResult(assertionResult)),
    })),
  };
};

export default AssertionResults;
