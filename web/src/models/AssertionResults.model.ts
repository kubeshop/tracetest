import {uniqueId} from 'lodash';
import AssertionService from '../services/Assertion.service';
import SelectorService from '../services/Selector.service';
import {TAssertionResults, TRawAssertionResults} from '../types/Assertion.types';
import AssertionResult from './AssertionResult.model';

const AssertionResults = ({allPassed = false, results = []}: TRawAssertionResults): TAssertionResults => {
  return {
    allPassed,
    resultList: results.map(({selector: {query: selector = '', structure = []} = {}, results: resultList = []}) => ({
      id: uniqueId(),
      selector,
      isAdvancedSelector: SelectorService.getIsAdvancedSelector(selector),
      spanIds: AssertionService.getSpanIds(resultList),
      originalSelector: selector,
      pseudoSelector: SelectorService.getPseudoSelectorFromStructure(structure),
      selectorList: SelectorService.getSpanSelectorListFromStructure(structure),
      resultList: resultList.map(assertionResult => AssertionResult(assertionResult)),
    })),
  };
};

export default AssertionResults;
