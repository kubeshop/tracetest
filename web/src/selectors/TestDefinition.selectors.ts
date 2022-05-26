import {createSelector} from '@reduxjs/toolkit';
import {RootState} from '../redux/store';

const stateSelector = (state: RootState) => state.testDefinition;
const selectorSelector = (state: RootState, selector: string) => selector;

const selectDefinitionList = createSelector(stateSelector, ({definitionList}) => definitionList);

const TestSelectors = () => ({
  selectDefinitionList,
  selectIsLoading: createSelector(stateSelector, ({isLoading}) => isLoading),
  selectIsInitialized: createSelector(stateSelector, ({isInitialized}) => isInitialized),
  selectAssertionResults: createSelector(stateSelector, ({assertionResults}) => assertionResults),
  selectDefinitionBySelector: createSelector(selectDefinitionList, selectorSelector, (definitionList, selector) =>
    definitionList.find(def => def.selector === selector)
  ),
});

export default TestSelectors();
