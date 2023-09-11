import {createAsyncThunk} from '@reduxjs/toolkit';
import {parse, ParsedQuery, stringify} from 'query-string';
import {Params} from 'react-router-dom';
import {push} from 'redux-first-history';
import {RouterSearchFields} from 'constants/Common.constants';
import TestSpecsSelectors from 'selectors/TestSpecs.selectors';
import DAGSelectors from 'selectors/DAG.selectors';
import SpanSelectors from 'selectors/Span.selectors';
import {setSelectedSpan} from 'redux/slices/Span.slice';
import {setSelectedSpec} from 'redux/slices/TestSpecs.slice';
import {RootState} from 'redux/store';

export interface IQuery {
  search: ParsedQuery<string>;
  params: Readonly<Params<string>>;
}

const RouterActions = () => ({
  updateSelectedAssertion: createAsyncThunk<void, IQuery>(
    'router/addAssertionResult',
    async ({search}, {getState, dispatch}) => {
      const {[RouterSearchFields.SelectedAssertion]: positionIndex} = search;
      const selectedSpec = TestSpecsSelectors.selectSelectedSpec(getState() as RootState);

      if (typeof positionIndex === 'undefined' && typeof selectedSpec !== 'undefined') {
        dispatch(setSelectedSpec());
        return;
      }

      const assertionResult = TestSpecsSelectors.selectAssertionByPositionIndex(
        getState() as RootState,
        Number(positionIndex)
      );
      const isDagReady = DAGSelectors.selectNodes(getState() as RootState).length > 0;

      if (selectedSpec === assertionResult?.selector) return;
      if (assertionResult && isDagReady) dispatch(setSelectedSpec(assertionResult));
    }
  ),
  updateSelectedSpan: createAsyncThunk<void, IQuery>(
    'router/updateSelectedSpan',
    async ({params: {testId = '', runId = '0'}, search}, {getState, dispatch}) => {
      const {[RouterSearchFields.SelectedSpan]: spanId = ''} = search;
      const state = getState() as RootState;
      const span = SpanSelectors.selectSpanById(state, String(spanId), testId, Number(runId));

      if (span) dispatch(setSelectedSpan({span}));
    }
  ),
  updateSearch: createAsyncThunk<void, Partial<Record<RouterSearchFields, any>>>(
    'router/updateSearch',
    async (newSearch, {dispatch, getState}) => {
      const {
        router: {location},
      } = getState() as RootState;

      const search = parse(location?.search || '');

      const filteredSearch = Object.entries({
        ...search,
        ...newSearch,
      }).reduce((acc, [key, value]) => (value ? {...acc, [key]: value} : acc), {});

      await dispatch(
        push({
          search: stringify(filteredSearch),
        })
      );
    }
  ),
});

export default RouterActions();
