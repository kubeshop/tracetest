import {createSelector} from '@reduxjs/toolkit';

import {RootState} from 'redux/store';
import {TResultAssertions} from 'types/Assertion.types';

const stateSelector = (state: RootState) => state.testDefinition;
const selectorSelector = (state: RootState, selector: string) => selector;
const spanIdSelector = (state: RootState, spanId: string) => spanId;

const selectDefinitionList = createSelector(stateSelector, ({definitionList}) => definitionList);

const selectDefinitionSelectorList = createSelector(selectDefinitionList, definitionList =>
  definitionList.map(({selector}) => selector)
);
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
            id: assertionResult.selector,
            attribute: assertion.assertion.attribute,
            assertionResult,
            label: `${assertionResult.selectorList.map(({value}) => value).join(' ')} ${
              assertionResult.pseudoSelector?.selector ?? ''
            }`,
            result: assertion.spanResults.find(spanResult => spanResult.spanId === spanId),
          }))
        )
        // Filter if it has result for the spanId
        .filter(assertion => Boolean(assertion?.result))
        // Hash items by attribute
        .reduce((prev: TResultAssertions, curr) => {
          const value = prev[curr.attribute] || {failed: [], passed: []};

          if (curr.result?.passed)
            value.passed.push({id: curr.id, label: curr.label, assertionResult: curr.assertionResult});
          else value.failed.push({id: curr.id, label: curr.label, assertionResult: curr.assertionResult});

          return {...prev, [curr.attribute]: value};
        }, {})
    );
  }
);

const TestDefinitionSelectors = () => ({
  selectDefinitionList,
  selectDefinitionSelectorList,
  selectIsSelectorExist: createSelector(selectDefinitionSelectorList, selectorSelector, (selectorList, selector) =>
    selectorList.includes(selector)
  ),
  selectIsLoading: createSelector(stateSelector, ({isLoading}) => isLoading),
  selectIsInitialized: createSelector(stateSelector, ({isInitialized}) => isInitialized),
  selectAssertionResults,
  selectDefinitionBySelector: createSelector(selectDefinitionList, selectorSelector, (definitionList, selector) =>
    definitionList.find(def => def.selector === selector)
  ),
  selectSelectedAssertion: createSelector(stateSelector, ({selectedAssertion}) => selectedAssertion),
  selectAssertionResultsBySpan,
  selectIsDraftMode: createSelector(stateSelector, ({isDraftMode}) => isDraftMode),
});

export default TestDefinitionSelectors();
