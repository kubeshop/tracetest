import {createSelector} from '@reduxjs/toolkit';
import {RootState} from 'redux/store';

const spansStateSelector = (state: RootState) => state.spans;
const stateSelector = (state: RootState) => state;

const SpanSelectors = () => ({
  selectAffectedSpans: createSelector(
    spansStateSelector,
    stateSelector,
    ({affectedSpans}, {testDefinition: {assertionResults, selectedAssertion}}) => {
      if (!selectedAssertion) return affectedSpans;

      const foundAssertion = assertionResults?.resultList.find(({selector}) => selector === selectedAssertion);

      return !foundAssertion ? [] : affectedSpans;
    }
  ),
  selectSelectedSpan: createSelector(spansStateSelector, ({selectedSpan}) => selectedSpan),
  selectFocusedSpan: createSelector(spansStateSelector, ({focusedSpan}) => focusedSpan),
  selectMatchedSpans: createSelector(spansStateSelector, ({matchedSpans}) => matchedSpans),
  selectSearchText: createSelector(spansStateSelector, ({searchText}) => searchText),
});

export default SpanSelectors();
