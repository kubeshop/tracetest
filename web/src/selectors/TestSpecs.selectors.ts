import {createSelector} from '@reduxjs/toolkit';

import {RootState} from 'redux/store';

const stateSelector = (state: RootState) => state.testSpecs;
const selectorSelector = (state: RootState, selector: string) => selector;
const positionIndexSelector = (state: RootState, positionIndex: number) => positionIndex;

const selectSpecs = createSelector(stateSelector, ({specs}) => specs);

const selectSpecsSelectorList = createSelector(selectSpecs, specs => specs.map(({selector}) => selector));
const selectAssertionResults = createSelector(stateSelector, ({assertionResults}) => assertionResults);

const selectTotalSpecs = createSelector(selectAssertionResults, assertionResults => {
  if (!assertionResults) return {totalFailedSpecs: 0, totalPassedSpecs: 0};

  return assertionResults.resultList.reduce<{totalFailedSpecs: number; totalPassedSpecs: number}>(
    ({totalFailedSpecs, totalPassedSpecs}, {resultList}) => {
      const someAssertionFailed = resultList.some(({allPassed}) => !allPassed);

      return {
        totalFailedSpecs: someAssertionFailed ? totalFailedSpecs + 1 : totalFailedSpecs,
        totalPassedSpecs: !someAssertionFailed ? totalPassedSpecs + 1 : totalPassedSpecs,
      };
    },
    {totalFailedSpecs: 0, totalPassedSpecs: 0}
  );
});

const TestSpecsSelectors = () => ({
  selectSpecs,
  selectSpecsSelectorList,
  selectIsSelectorExist: createSelector(selectSpecsSelectorList, selectorSelector, (selectorList, selector = '') =>
    selectorList.includes(selector)
  ),
  selectIsLoading: createSelector(stateSelector, ({isLoading}) => isLoading),
  selectIsInitialized: createSelector(stateSelector, ({isInitialized}) => isInitialized),
  selectAssertionResults,
  selectSpecBySelector: createSelector(selectSpecs, selectorSelector, (specs, selector) =>
    specs.find(spec => spec.selector === selector)
  ),
  selectAssertionBySelector: createSelector(stateSelector, selectorSelector, ({assertionResults}, selector) =>
    assertionResults?.resultList.find(def => def.selector === selector)
  ),
  selectSelectedSpanResult: createSelector(stateSelector, ({selectedSpanResult}) => selectedSpanResult),
  selectAssertionByPositionIndex: createSelector(
    stateSelector,
    positionIndexSelector,
    ({assertionResults}, positionIndex) => assertionResults?.resultList[positionIndex]
  ),
  selectSelectedSpec: createSelector(stateSelector, ({selectedSpec}) => selectedSpec),
  selectIsDraftMode: createSelector(stateSelector, ({isDraftMode}) => isDraftMode),
  selectTotalSpecs,
  selectSelectorSuggestions: createSelector(stateSelector, ({selectorSuggestions}) => selectorSuggestions),
  selectPrevSelector: createSelector(stateSelector, ({prevSelector}) => prevSelector),
});

export default TestSpecsSelectors();
