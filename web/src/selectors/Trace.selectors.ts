import {createSelector} from '@reduxjs/toolkit';
import {RootState} from 'redux/store';
import {ITraceState} from 'redux/slices/Trace.slice';

const traceStateSelector = (state: RootState): ITraceState => state.trace;

const TraceSelectors = () => ({
  selectEdges: createSelector(traceStateSelector, ({edges}) => edges),
  selectMatchedSpans: createSelector(traceStateSelector, ({matchedSpans}) => matchedSpans),
  selectNodes: createSelector(traceStateSelector, ({nodes}) => nodes),
  selectSelectedSpan: createSelector(traceStateSelector, ({selectedSpan}) => selectedSpan),
  selectSearchText: createSelector(traceStateSelector, ({searchText}) => searchText),
  selectSelectedAnalyzerResults: createSelector(
    traceStateSelector,
    ({selectedAnalyzerResults}) => selectedAnalyzerResults
  ),
});

export default TraceSelectors();
