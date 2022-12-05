import {createSelector} from '@reduxjs/toolkit';

import {RootState} from 'redux/store';
import {TResultAssertions} from 'types/Assertion.types';
import {TSpansResult} from 'types/Span.types';

const stateSelector = (state: RootState) => state.testSpecs;
const selectorSelector = (state: RootState, selector: string) => selector;
const spanIdSelector = (state: RootState, spanId: string) => spanId;
const positionIndexSelector = (state: RootState, positionIndex: number) => positionIndex;

const selectSpecs = createSelector(stateSelector, ({specs}) => specs);

const selectSpecsSelectorList = createSelector(selectSpecs, specs => specs.map(({selector}) => selector));
const selectAssertionResults = createSelector(stateSelector, ({assertionResults}) => assertionResults);

const selectAssertionResultsBySpan = createSelector(
  selectAssertionResults,
  spanIdSelector,
  (assertionResults, spanId) => {
    if (!assertionResults) return {};

    // Map and flat items in one single array
    return (
      assertionResults.resultList
        .flatMap(assertionResult =>
          assertionResult.resultList.map(assertion => ({
            id: assertionResult.selector  || 'All Spans',
            attribute: assertion.assertion,
            assertionResult,
            label: assertionResult.selector || 'All Spans',
            result: assertion.spanResults.find(spanResult => spanResult.spanId === spanId),
          }))
        )
        // Filter if it has result for the spanId
        .filter(assertion => Boolean(assertion?.result))
        // Hash items by attribute
        .reduce((prev: TResultAssertions, curr) => {
          const value = prev[curr.attribute || ''] || {failed: [], passed: []};

          if (curr.result?.passed)
            value.passed.push({id: curr.id, label: curr.label, assertionResult: curr.assertionResult});
          else value.failed.push({id: curr.id, label: curr.label, assertionResult: curr.assertionResult});

          return {...prev, [curr.attribute || '']: value};
        }, {})
    );
  }
);

const selectSpansResult = createSelector(selectAssertionResults, assertionResults => {
  if (!assertionResults) return {};

  // Map and flat items in one single array
  const results = assertionResults.resultList
    .flatMap(resultItem => resultItem.resultList)
    .flatMap(resultItem => resultItem.spanResults)
    // Hash items by spanId
    .reduce((prev: TSpansResult, curr) => {
      const value = prev[curr?.spanId] || {failed: 0, passed: 0};

      if (curr?.passed) value.passed += 1;
      else value.failed += 1;

      return {...prev, [curr?.spanId]: value};
    }, {});

  return results;
});

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
  selectAssertionByPositionIndex: createSelector(
    stateSelector,
    positionIndexSelector,
    ({assertionResults}, positionIndex) => assertionResults?.resultList[positionIndex]
  ),
  selectSelectedSpec: createSelector(stateSelector, ({selectedSpec}) => selectedSpec),
  selectAssertionResultsBySpan,
  selectIsDraftMode: createSelector(stateSelector, ({isDraftMode}) => isDraftMode),
  selectSpansResult,
  selectTotalSpecs,
  selectSelectorSuggestions: createSelector(stateSelector, ({selectorSuggestions}) => selectorSuggestions),
  selectPrevSelector: createSelector(stateSelector, ({prevSelector}) => prevSelector),
});

export default TestSpecsSelectors();
