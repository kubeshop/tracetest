import {createSelector} from '@reduxjs/toolkit';

import {RootState} from '../redux/store';

const stateSelector = (state: RootState) => state.testDefinition;
const selectorSelector = (state: RootState, selector: string) => selector;

const selectDefinitionList = createSelector(stateSelector, ({definitionList}) => definitionList);
const selectDefinitionSelectorList = createSelector(selectDefinitionList, definitionList =>
  definitionList.map(({selector}) => selector)
);

const TestDefinitionSelectors = () => ({
  selectDefinitionList,
  selectDefinitionSelectorList,
  selectIsSelectorExist: createSelector(selectDefinitionSelectorList, selectorSelector, (selectorList, selector) =>
    selectorList.includes(selector)
  ),
  selectIsLoading: createSelector(stateSelector, ({isLoading}) => isLoading),
  selectIsInitialized: createSelector(stateSelector, ({isInitialized}) => isInitialized),
  selectAssertionResults: createSelector(stateSelector, ({assertionResults}) => assertionResults),
  selectDefinitionBySelector: createSelector(selectDefinitionList, selectorSelector, (definitionList, selector) =>
    definitionList.find(def => def.selector === selector)
  ),
  selectAffectedSpans: createSelector(stateSelector, ({affectedSpans}) => affectedSpans),
});

export default TestDefinitionSelectors();
